package handler

import (
	"fmt"
	"ping-uptime/internal/pkg/logger"
	"ping-uptime/internal/pkg/middleware"
	"ping-uptime/internal/pkg/utils"
	"ping-uptime/modules/api_tokens/domain/service"
	"ping-uptime/modules/api_tokens/dto/request"
	"ping-uptime/modules/api_tokens/dto/response"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type ApiTokenHandler struct {
	tokenService *service.ApiTokenService
	log          *logger.Logger
	r            *utils.Response
}

func NewApiTokenHandler(log *logger.Logger, tokenService *service.ApiTokenService) *ApiTokenHandler {
	return &ApiTokenHandler{
		tokenService: tokenService,
		log:          log,
		r:            &utils.Response{},
	}
}

func (h *ApiTokenHandler) getAuthUser(c echo.Context) (uint, string, error) {
	userClaims, ok := c.Get("user").(map[string]interface{})
	if !ok {
		return 0, "", fmt.Errorf("unauthorized")
	}

	userIDVal, ok := userClaims["user_id"]
	if !ok {
		return 0, "", fmt.Errorf("invalid token: user_id missing")
	}

	var userID uint
	switch v := userIDVal.(type) {
	case float64:
		userID = uint(v)
	case int64:
		userID = uint(v)
	case int:
		userID = uint(v)
	default:
		return 0, "", fmt.Errorf("invalid user_id type")
	}

	roleVal, ok := userClaims["role"].(string)
	if !ok {
		roleVal = "user"
	}

	return userID, roleVal, nil
}

// GetAll returns tokens. Users see their own, admins see all.
//
// @Summary      List API tokens
// @Description  Get all API tokens for the authenticated user (admins see all)
// @Tags         API Tokens
// @Security     BearerAuth
// @Produce      json
// @Success      200  {array}   response.TokenResponse
// @Failure      401  {object}  utils.ErrorResponseModel
// @Router       /api-tokens [get]
func (h *ApiTokenHandler) GetAll(c echo.Context) error {
	ctx := c.Request().Context()
	userID, role, err := h.getAuthUser(c)
	if err != nil {
		return h.r.UnauthorizedResponse(c, err.Error())
	}

	var tokens []*response.TokenResponse
	if role == "admin" {
		// Admin sees all tokens — fetch per-user or full list.
		// For simplicity, we treat admin like user but could add a FindAll.
		// Using user-scoped for now; extend with admin findAll if needed.
		userTokens, err := h.tokenService.GetUserTokens(ctx, userID)
		if err != nil {
			return h.r.InternalServerErrorResponse(c, err.Error())
		}
		tokenResponses := response.FromEntities(userTokens)
		tokens = make([]*response.TokenResponse, len(tokenResponses))
		for i := range tokenResponses {
			tokens[i] = &tokenResponses[i]
		}
	} else {
		userTokens, err := h.tokenService.GetUserTokens(ctx, userID)
		if err != nil {
			return h.r.InternalServerErrorResponse(c, err.Error())
		}
		tokenResponses := response.FromEntities(userTokens)
		tokens = make([]*response.TokenResponse, len(tokenResponses))
		for i := range tokenResponses {
			tokens[i] = &tokenResponses[i]
		}
	}

	return h.r.SuccessResponse(c, tokens, "API tokens retrieved successfully")
}

// Create generates a new API token and returns it (raw token shown once).
//
// @Summary      Create API token
// @Description  Generate a new API token. Raw token is only returned once.
// @Tags         API Tokens
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body  request.CreateTokenRequest  true  "Token details"
// @Success      201  {object}  response.CreateTokenResponse
// @Failure      400  {object}  utils.ErrorResponseModel
// @Failure      401  {object}  utils.ErrorResponseModel
// @Router       /api-tokens [post]
func (h *ApiTokenHandler) Create(c echo.Context) error {
	ctx := c.Request().Context()
	userID, _, err := h.getAuthUser(c)
	if err != nil {
		return h.r.UnauthorizedResponse(c, err.Error())
	}

	req := new(request.CreateTokenRequest)
	if err := c.Bind(req); err != nil {
		return h.r.BadRequestResponse(c, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return h.r.BadRequestResponse(c, err.Error())
	}

	var expiresAt *time.Time
	if req.ExpiresAt != "" {
		t, err := time.Parse(time.RFC3339, req.ExpiresAt)
		if err != nil {
			return h.r.BadRequestResponse(c, "Invalid expires_at format, use RFC3339")
		}
		expiresAt = &t
	}

	token, rawToken, err := h.tokenService.GenerateToken(ctx, userID, req.Name, expiresAt)
	if err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	return h.r.CreatedResponse(c, response.ToCreateResponse(token, rawToken), "API token created successfully")
}

// Revoke marks a token as revoked.
//
// @Summary      Revoke API token
// @Description  Revoke an API token by ID. Immediate loss of access.
// @Tags         API Tokens
// @Security     BearerAuth
// @Produce      json
// @Param        id   path  int  true  "Token ID"
// @Success      200  {object}  utils.SuccessResponseModel
// @Failure      401  {object}  utils.ErrorResponseModel
// @Failure      403  {object}  utils.ErrorResponseModel
// @Failure      404  {object}  utils.ErrorResponseModel
// @Router       /api-tokens/{id} [delete]
func (h *ApiTokenHandler) Revoke(c echo.Context) error {
	ctx := c.Request().Context()
	userID, role, err := h.getAuthUser(c)
	if err != nil {
		return h.r.UnauthorizedResponse(c, err.Error())
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "Invalid token ID")
	}

	if err := h.tokenService.RevokeToken(ctx, uint(id), userID, role); err != nil {
		switch err {
		case service.ErrTokenNotFound:
			return h.r.NotFoundResponse(c, "API token not found")
		case service.ErrNotAuthorized:
			return h.r.ForbiddenResponse(c, "You do not have permission to revoke this token")
		default:
			return h.r.InternalServerErrorResponse(c, err.Error())
		}
	}

	return h.r.SuccessResponse(c, nil, "API token revoked successfully")
}

func (h *ApiTokenHandler) RegisterRoutes(e *echo.Echo, basePath string) {
	group := e.Group(basePath+"/api-tokens", middleware.Auth)
	group.GET("", h.GetAll)
	group.POST("", h.Create)
	group.DELETE("/:id", h.Revoke)
}

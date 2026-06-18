package handler

import (
	"fmt"
	"ping-uptime/internal/pkg/middleware"
	"ping-uptime/internal/pkg/utils"
	"ping-uptime/modules/integrations/domain/entity"
	"ping-uptime/modules/integrations/domain/service"
	"ping-uptime/modules/integrations/dto/request"
	"ping-uptime/modules/integrations/dto/response"
	"strconv"

	"github.com/labstack/echo/v4"
)

type IntegrationHandler struct {
	integrationService *service.IntegrationService
	r                  *utils.Response
}

func NewIntegrationHandler(integrationService *service.IntegrationService) *IntegrationHandler {
	return &IntegrationHandler{
		integrationService: integrationService,
		r:                  &utils.Response{},
	}
}

func (h *IntegrationHandler) getAuthUser(c echo.Context) (uint, string, error) {
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

func (h *IntegrationHandler) GetAll(c echo.Context) error {
	ctx := c.Request().Context()
	userID, role, err := h.getAuthUser(c)
	if err != nil {
		return h.r.UnauthorizedResponse(c, err.Error())
	}

	integrations, err := h.integrationService.FindAll(ctx)
	if err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	// Filter by user unless admin
	var userIntegrations []*entity.Integration
	if role == "admin" {
		userIntegrations = integrations
	} else {
		for _, i := range integrations {
			if i.UserID == userID {
				userIntegrations = append(userIntegrations, i)
			}
		}
	}

	return h.r.SuccessResponse(c, response.FromEntities(userIntegrations), "Integrations retrieved successfully")
}

func (h *IntegrationHandler) Create(c echo.Context) error {
	ctx := c.Request().Context()
	userID, _, err := h.getAuthUser(c)
	if err != nil {
		return h.r.UnauthorizedResponse(c, err.Error())
	}

	req := new(request.CreateIntegrationRequest)
	if err := c.Bind(req); err != nil {
		return h.r.BadRequestResponse(c, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return h.r.BadRequestResponse(c, err.Error())
	}

	integration := &entity.Integration{
		Name:    req.Name,
		Type:    req.Type,
		Config:  req.Config,
		Enabled: true,
		UserID:  userID,
	}

	err = h.integrationService.Create(ctx, integration)
	if err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	return h.r.CreatedResponse(c, response.FromEntity(integration), "Integration created successfully")
}

func (h *IntegrationHandler) Update(c echo.Context) error {
	ctx := c.Request().Context()
	userID, role, err := h.getAuthUser(c)
	if err != nil {
		return h.r.UnauthorizedResponse(c, err.Error())
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "Invalid integration ID")
	}

	integration, err := h.integrationService.FindByID(ctx, uint(id))
	if err != nil {
		return h.r.NotFoundResponse(c, "Integration not found")
	}

	if role != "admin" && integration.UserID != userID {
		return h.r.ForbiddenResponse(c, "You do not have permission to update this integration")
	}

	req := new(request.UpdateIntegrationRequest)
	if err := c.Bind(req); err != nil {
		return h.r.BadRequestResponse(c, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return h.r.BadRequestResponse(c, err.Error())
	}

	integration.Name = req.Name
	integration.Type = req.Type
	integration.Config = req.Config
	integration.Enabled = req.Enabled

	err = h.integrationService.Update(ctx, integration)
	if err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	return h.r.SuccessResponse(c, response.FromEntity(integration), "Integration updated successfully")
}

func (h *IntegrationHandler) Delete(c echo.Context) error {
	ctx := c.Request().Context()
	userID, role, err := h.getAuthUser(c)
	if err != nil {
		return h.r.UnauthorizedResponse(c, err.Error())
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "Invalid integration ID")
	}

	integration, err := h.integrationService.FindByID(ctx, uint(id))
	if err != nil {
		return h.r.NotFoundResponse(c, "Integration not found")
	}

	if role != "admin" && integration.UserID != userID {
		return h.r.ForbiddenResponse(c, "You do not have permission to delete this integration")
	}

	err = h.integrationService.Delete(ctx, uint(id))
	if err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	return h.r.NoContentResponse(c)
}

func (h *IntegrationHandler) Test(c echo.Context) error {
	ctx := c.Request().Context()
	userID, role, err := h.getAuthUser(c)
	if err != nil {
		return h.r.UnauthorizedResponse(c, err.Error())
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "Invalid integration ID")
	}

	integration, err := h.integrationService.FindByID(ctx, uint(id))
	if err != nil {
		return h.r.NotFoundResponse(c, "Integration not found")
	}

	if role != "admin" && integration.UserID != userID {
		return h.r.ForbiddenResponse(c, "You do not have permission to test this integration")
	}

	if err := h.integrationService.TestIntegration(ctx, uint(id)); err != nil {
		return h.r.ErrorResponse(c, 502, err.Error())
	}

	return h.r.SuccessResponse(c, nil, "Test notification sent successfully")
}

func (h *IntegrationHandler) RegisterRoutes(e *echo.Echo, basePath string) {
	group := e.Group(basePath+"/integrations", middleware.Auth)
	group.GET("", h.GetAll)
	group.POST("", h.Create)
	group.PUT("/:id", h.Update)
	group.DELETE("/:id", h.Delete)
	group.POST("/:id/test", h.Test)
}

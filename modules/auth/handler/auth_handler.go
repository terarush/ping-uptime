package handler

import (
	"fmt"
	"ping-uptime/internal/pkg/bus"
	"ping-uptime/internal/pkg/database"
	"ping-uptime/internal/pkg/jwt"
	"ping-uptime/internal/pkg/logger"
	"ping-uptime/internal/pkg/utils"
	"ping-uptime/modules/auth/domain/service"
	"ping-uptime/modules/users/domain/entity"
	"ping-uptime/modules/users/dto/request"
	"ping-uptime/modules/users/dto/response"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// AuthHandler struct handles HTTP request for auth.
type AuthHandler struct {
	authService *service.AuthService
	log         *logger.Logger
	event       *bus.EventBus
	jwt         jwt.JWT
	r           *utils.Response
}

// NewAuthHandler creates a new auth handler.
func NewAuthHandler(log *logger.Logger, event *bus.EventBus, authService *service.AuthService, jwt jwt.JWT) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		log:         log,
		event:       event,
		jwt:         jwt,
		r:           &utils.Response{},
	}
}

// Initialize Event Handle.
func (h *AuthHandler) Handle(event bus.Event) {
	fmt.Printf("User created: %v", event.Payload)
}

// Setup handles initial admin account setup.
func (h *AuthHandler) Setup(c echo.Context) error {
	h.log.Info("Handling setup request")

	req := new(request.CreateUserRequest)
	if err := c.Bind(req); err != nil {
		h.log.Error("Failed to bind request:", err)
		return h.r.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		h.log.Error("Validation failed:", err)
		return h.r.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	h.log.Debug("Request validated successfully:", req)

	user := entity.NewUser(req.Name, req.Email, req.Password)
	user.Role = "admin"
	err := h.authService.CreateUser(c.Request().Context(), user)
	if err != nil {
		if err == service.ErrEmailAlreadyUsed {
			h.log.Warn("Email already in use:", req.Email)
			return h.r.ErrorResponse(c, http.StatusConflict, "Email already in use")
		}
		if err == service.ErrRegistrationDisabled {
			h.log.Warn("Setup is disabled: setup already completed")
			return h.r.ErrorResponse(c, http.StatusForbidden, "Setup already completed. Setup route is disabled.")
		}
		h.log.Error("Failed to create admin:", err)
		return h.r.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	h.log.Debug("Admin user created successfully:", user)

	h.event.Publish(bus.Event{Type: "user.created", Payload: user})
	h.log.Debug("Event 'user.created' published successfully")

	return h.r.SuccessResponse(c, map[string]interface{}{
		"user": response.FromEntity(user),
	}, "Administrator setup successfully")
}

// SetupStatus checks if the application is set up (i.e., has at least one user).
func (h *AuthHandler) SetupStatus(c echo.Context) error {
	h.log.Info("Checking setup status")

	isSetupNeeded, err := h.authService.IsSetupNeeded(c.Request().Context())
	if err != nil {
		h.log.Error("Failed to check setup status:", err)
		return h.r.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	var settings []struct {
		Key   string
		Value string
	}
	_ = database.DB.WithContext(c.Request().Context()).Table("settings").Where("key IN ?", []string{"system_name", "allow_registration"}).Find(&settings).Error

	systemName := "ping-uptime"
	allowRegistration := "true"
	for _, s := range settings {
		if s.Key == "system_name" {
			systemName = s.Value
		} else if s.Key == "allow_registration" {
			allowRegistration = s.Value
		}
	}

	// Setup is done if setup is NOT needed
	return h.r.SuccessResponse(c, map[string]interface{}{
		"is_setup":           !isSetupNeeded,
		"system_name":        systemName,
		"allow_registration": allowRegistration == "true",
	}, "Setup status retrieved successfully")
}

// Login handles user login.
func (h *AuthHandler) Login(c echo.Context) error {
	h.log.Info("Handling login request")

	req := new(request.LoginRequest)
	if err := c.Bind(req); err != nil {
		h.log.Error("Failed to bind request:", err)
		return h.r.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		h.log.Error("Validation failed:", err)
		return h.r.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	h.log.Debug("Request validated successfully:", req)

	user, err := h.authService.ProcessLogin(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		if err == service.ErrUserNotFound || err == service.ErrInvalidPassword {
			h.log.Warn("Invalid email or password for:", req.Email)
			return h.r.ErrorResponse(c, http.StatusUnauthorized, "Invalid email or password")
		}
		h.log.Error("Failed to process login:", err)
		return h.r.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	h.log.Debug("User authenticated successfully:", user)

	accessClaims := map[string]interface{}{
		"user_id": user.ID,
		"email":   user.Email,
		"name":    user.Name,
		"role":    user.Role,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	}

	accessToken, err := h.jwt.GenerateToken(accessClaims)
	if err != nil {
		h.log.Error("Failed to generate access token:", err)
		return h.r.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	refreshClaims := map[string]interface{}{
		"user_id": user.ID,
		"email":   user.Email,
		"name":    user.Name,
		"role":    user.Role,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
	}

	refreshToken, err := h.jwt.GenerateToken(refreshClaims)
	if err != nil {
		h.log.Error("Failed to generate refresh token:", err)
		return h.r.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return h.r.SuccessResponse(c, map[string]interface{}{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
		"user":         response.FromEntity(user),
	}, "Login successful")
}

// RefreshRequest defines the structure for token refresh payload
type RefreshRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

// Refresh handles regenerating a new access token using a refresh token.
func (h *AuthHandler) Refresh(c echo.Context) error {
	h.log.Info("Handling token refresh request")

	req := new(RefreshRequest)
	if err := c.Bind(req); err != nil {
		h.log.Error("Failed to bind refresh request:", err)
		return h.r.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		h.log.Error("Validation failed for refresh request:", err)
		return h.r.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	// Validate refresh token
	valid, err := h.jwt.ValidateToken(req.RefreshToken)
	if err != nil || !valid {
		h.log.Warn("Invalid or expired refresh token")
		return h.r.ErrorResponse(c, http.StatusUnauthorized, "Invalid or expired refresh token")
	}

	// Parse refresh claims
	claims, err := h.jwt.ParseToken(req.RefreshToken)
	if err != nil {
		h.log.Error("Failed to parse refresh token claims:", err)
		return h.r.ErrorResponse(c, http.StatusUnauthorized, "Failed to parse token claims")
	}

	// Validate user still exists in DB and is not blocked
	userID := uint(0)
	if uid, ok := claims["user_id"].(float64); ok {
		userID = uint(uid)
	}
	if userID > 0 {
		if err := h.authService.ValidateUserExists(c.Request().Context(), userID); err != nil {
			h.log.Warn("Refresh denied — user not found or blocked: %v", err)
			return h.r.UnauthorizedResponse(c, "Account not available")
		}
	}

	// Generate a new access token (15 mins)
	accessClaims := map[string]interface{}{
		"user_id": claims["user_id"],
		"email":   claims["email"],
		"name":    claims["name"],
		"role":    claims["role"],
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	}

	accessToken, err := h.jwt.GenerateToken(accessClaims)
	if err != nil {
		h.log.Error("Failed to generate new access token:", err)
		return h.r.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return h.r.SuccessResponse(c, map[string]interface{}{
		"accessToken": accessToken,
	}, "Token refreshed successfully")
}

// Register handles public user registration.
func (h *AuthHandler) Register(c echo.Context) error {
	h.log.Info("Handling public registration request")

	req := new(request.CreateUserRequest)
	if err := c.Bind(req); err != nil {
		h.log.Error("Failed to bind request:", err)
		return h.r.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		h.log.Error("Validation failed:", err)
		return h.r.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	h.log.Debug("Request validated successfully:", req)

	user := entity.NewUser(req.Name, req.Email, req.Password)
	user.Role = "user"

	err := h.authService.Register(c.Request().Context(), user)
	if err != nil {
		if err == service.ErrEmailAlreadyUsed {
			h.log.Warn("Email already in use:", req.Email)
			return h.r.ErrorResponse(c, http.StatusConflict, "Email already in use")
		}
		if err == service.ErrRegistrationDisabled {
			h.log.Warn("Registration is disabled by administrator")
			return h.r.ErrorResponse(c, http.StatusForbidden, "Registration is disabled by administrator")
		}
		h.log.Error("Failed to register user:", err)
		return h.r.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	h.log.Debug("User registered successfully:", user)
	h.event.Publish(bus.Event{Type: "user.created", Payload: user})

	return h.r.SuccessResponse(c, map[string]interface{}{
		"user": response.FromEntity(user),
	}, "Registered successfully")
}

// RegisterRoutes sets up the auth routes.
func (h *AuthHandler) RegisterRoutes(e *echo.Echo, basePath string) {
	group := e.Group(basePath + "/auth")
	group.GET("/setup-status", h.SetupStatus)
	group.POST("/setup", h.Setup)
	group.POST("/login", h.Login)
	group.POST("/refresh", h.Refresh)
	group.POST("/register", h.Register)
}

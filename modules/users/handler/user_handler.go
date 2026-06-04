// internal/modules/user/interfaces/handler/user_handler.go

package handler

import (
	"fmt"
	"net/http"
	"ping-uptime/internal/pkg/bus"
	"ping-uptime/internal/pkg/logger"
	"ping-uptime/internal/pkg/middleware"
	"ping-uptime/internal/pkg/utils"
	"ping-uptime/modules/users/domain/entity"
	"ping-uptime/modules/users/domain/service"
	"ping-uptime/modules/users/dto/request"
	"ping-uptime/modules/users/dto/response"
	"strconv"

	"github.com/labstack/echo/v4"
)

// UserHandler handles HTTP requests for users
type UserHandler struct {
	userService *service.UserService
	log         *logger.Logger
	event       *bus.EventBus
	r           *utils.Response
}

// NewUserHandler creates a new user handler
func NewUserHandler(log *logger.Logger, event *bus.EventBus, userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
		log:         log,
		event:       event,
		r: &utils.Response{},
	}
}

// Event Bus Event user created
func (h *UserHandler) Handle(event bus.Event) {
	fmt.Printf("User created: %v", event.Payload)
}

// GetAllUsers gets all users
func (h *UserHandler) GetAllUsers(c echo.Context) error {
	ctx := c.Request().Context()

	users, err := h.userService.GetAllUsers(ctx)
	if err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	return h.r.SuccessResponse(c, response.FromEntities(users), "Users retrieved successfully")
}

// GetUser gets a user by ID
func (h *UserHandler) GetUser(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "Invalid user ID")
	}

	user, err := h.userService.GetUserByID(ctx, uint(id))
	if err != nil {
		if err == service.ErrUserNotFound {
			return h.r.NotFoundResponse(c, "User not found")
		}
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	return h.r.SuccessResponse(c, response.FromEntity(user), "User retrieved successfully")
}

// CreateUser creates a new user
func (h *UserHandler) CreateUser(c echo.Context) error {
	ctx := c.Request().Context()

	req := new(request.CreateUserRequest)
	if err := c.Bind(req); err != nil {
		return h.r.BadRequestResponse(c, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return h.r.BadRequestResponse(c, err.Error())
	}

	user := entity.NewUser(req.Name, req.Email, req.Password)
	err := h.userService.CreateUser(ctx, user)
	if err != nil {
		if err == service.ErrEmailAlreadyUsed {
			return h.r.ConflictResponse(c, "Email already in use")
		}
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	// event bus publish
	h.event.Publish(bus.Event{Type: "user.created", Payload: user})

	return h.r.CreatedResponse(c, response.FromEntity(user), "User created successfully")
}

// UpdateUser updates a user
func (h *UserHandler) UpdateUser(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "Invalid user ID")
	}

	req := new(request.UpdateUserRequest)
	if err := c.Bind(req); err != nil {
		return h.r.BadRequestResponse(c, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return h.r.BadRequestResponse(c, err.Error())
	}

	user, err := h.userService.GetUserByID(ctx, uint(id))
	if err != nil {
		if err == service.ErrUserNotFound {
			return h.r.NotFoundResponse(c, "User not found")
		}
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	user.Name = req.Name
	user.Email = req.Email
	if req.Password != "" {
		user.Password = req.Password
	}

	err = h.userService.UpdateUser(ctx, user)
	if err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	return h.r.SuccessResponse(c, response.FromEntity(user), "User updated successfully")
}

// DeleteUser deletes a user
func (h *UserHandler) DeleteUser(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "Invalid user ID")
	}

	err = h.userService.DeleteUser(ctx, uint(id))
	if err != nil {
		if err == service.ErrUserNotFound {
			return h.r.NotFoundResponse(c, "User not found")
		}
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	return h.r.NoContentResponse(c)
}

// Verify verifies the authenticated user's token and returns user details
func (h *UserHandler) Verify(c echo.Context) error {
	ctx := c.Request().Context()
	userClaims, ok := c.Get("user").(map[string]interface{})
	if !ok {
		return h.r.UnauthorizedResponse(c, "Unauthorized")
	}

	userIDVal, ok := userClaims["user_id"]
	if !ok {
		return h.r.UnauthorizedResponse(c, "Invalid token payload: user_id missing")
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
		return h.r.UnauthorizedResponse(c, "Invalid token payload: invalid user_id type")
	}

	user, err := h.userService.GetUserByID(ctx, userID)
	if err != nil {
		if err == service.ErrUserNotFound {
			return h.r.UnauthorizedResponse(c, "User not found")
		}
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	return h.r.SuccessResponse(c, response.FromEntity(user), "Token verified successfully")
}

// RegisterRoutes registers the user routes
func (h *UserHandler) RegisterRoutes(e *echo.Echo, basePath string) {
	// Standard auth validation route accessible to any valid session
	authGroup := e.Group(basePath+"/users", middleware.Auth)
	authGroup.GET("/verify", h.Verify)

	// Restricted admin-only routes
	adminGroup := e.Group(basePath+"/users", middleware.Auth, middleware.Admin)
	adminGroup.GET("", h.GetAllUsers)
	adminGroup.GET("/:id", h.GetUser)
	adminGroup.POST("", h.CreateUser)
	adminGroup.PUT("/:id", h.UpdateUser)
	adminGroup.DELETE("/:id", h.DeleteUser)
}

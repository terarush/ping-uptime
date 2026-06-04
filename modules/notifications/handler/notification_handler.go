package handler

import (
	"fmt"
	"ping-uptime/internal/pkg/bus"
	"ping-uptime/internal/pkg/logger"
	"ping-uptime/internal/pkg/middleware"
	"ping-uptime/internal/pkg/utils"
	"ping-uptime/modules/notifications/domain/entity"
	"ping-uptime/modules/notifications/domain/service"
	"ping-uptime/modules/notifications/dto/request"
	"ping-uptime/modules/notifications/dto/response"
	"strconv"

	"github.com/labstack/echo/v4"
)

type NotificationHandler struct {
	notificationService *service.NotificationService
	log                 *logger.Logger
	event               *bus.EventBus
	r                   *utils.Response
}

func NewNotificationHandler(log *logger.Logger, event *bus.EventBus, notificationService *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
		log:                 log,
		event:               event,
		r:                   &utils.Response{},
	}
}

func (h *NotificationHandler) getAuthUser(c echo.Context) (uint, string, error) {
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

func (h *NotificationHandler) GetAllChannels(c echo.Context) error {
	ctx := c.Request().Context()
	userID, role, err := h.getAuthUser(c)
	if err != nil {
		return h.r.UnauthorizedResponse(c, err.Error())
	}

	var channels []*entity.NotificationChannel
	if role == "admin" {
		channels, err = h.notificationService.GetAllChannels(ctx)
	} else {
		channels, err = h.notificationService.GetChannelsByUserID(ctx, userID)
	}

	if err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	return h.r.SuccessResponse(c, response.FromEntities(channels), "Notification channels retrieved successfully")
}

func (h *NotificationHandler) GetChannel(c echo.Context) error {
	ctx := c.Request().Context()
	userID, role, err := h.getAuthUser(c)
	if err != nil {
		return h.r.UnauthorizedResponse(c, err.Error())
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "Invalid channel ID")
	}

	channel, err := h.notificationService.GetChannelByID(ctx, uint(id))
	if err != nil {
		return h.r.NotFoundResponse(c, "Notification channel not found")
	}

	if role != "admin" && channel.UserID != userID {
		return h.r.ForbiddenResponse(c, "You do not have access to this notification channel")
	}

	return h.r.SuccessResponse(c, response.FromEntity(channel), "Notification channel retrieved successfully")
}

func (h *NotificationHandler) CreateChannel(c echo.Context) error {
	ctx := c.Request().Context()
	userID, _, err := h.getAuthUser(c)
	if err != nil {
		return h.r.UnauthorizedResponse(c, err.Error())
	}

	req := new(request.CreateChannelRequest)
	if err := c.Bind(req); err != nil {
		return h.r.BadRequestResponse(c, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return h.r.BadRequestResponse(c, err.Error())
	}

	channel := entity.NewNotificationChannel(req.Name, req.Type, req.Config, req.Enabled, userID)
	err = h.notificationService.CreateChannel(ctx, channel)
	if err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	h.event.Publish(bus.Event{Type: "notification.channel.created", Payload: channel})

	return h.r.CreatedResponse(c, response.FromEntity(channel), "Notification channel created successfully")
}

func (h *NotificationHandler) UpdateChannel(c echo.Context) error {
	ctx := c.Request().Context()
	userID, role, err := h.getAuthUser(c)
	if err != nil {
		return h.r.UnauthorizedResponse(c, err.Error())
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "Invalid channel ID")
	}

	channel, err := h.notificationService.GetChannelByID(ctx, uint(id))
	if err != nil {
		return h.r.NotFoundResponse(c, "Notification channel not found")
	}

	if role != "admin" && channel.UserID != userID {
		return h.r.ForbiddenResponse(c, "You do not have permission to update this notification channel")
	}

	req := new(request.UpdateChannelRequest)
	if err := c.Bind(req); err != nil {
		return h.r.BadRequestResponse(c, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return h.r.BadRequestResponse(c, err.Error())
	}

	channel.Name = req.Name
	channel.Type = req.Type
	channel.Config = req.Config
	channel.Enabled = req.Enabled

	err = h.notificationService.UpdateChannel(ctx, channel)
	if err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	h.event.Publish(bus.Event{Type: "notification.channel.updated", Payload: channel})

	return h.r.SuccessResponse(c, response.FromEntity(channel), "Notification channel updated successfully")
}

func (h *NotificationHandler) DeleteChannel(c echo.Context) error {
	ctx := c.Request().Context()
	userID, role, err := h.getAuthUser(c)
	if err != nil {
		return h.r.UnauthorizedResponse(c, err.Error())
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "Invalid channel ID")
	}

	channel, err := h.notificationService.GetChannelByID(ctx, uint(id))
	if err != nil {
		return h.r.NotFoundResponse(c, "Notification channel not found")
	}

	if role != "admin" && channel.UserID != userID {
		return h.r.ForbiddenResponse(c, "You do not have permission to delete this notification channel")
	}

	err = h.notificationService.DeleteChannel(ctx, uint(id))
	if err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	h.event.Publish(bus.Event{Type: "notification.channel.deleted", Payload: channel})

	return h.r.NoContentResponse(c)
}

func (h *NotificationHandler) RegisterRoutes(e *echo.Echo, basePath string) {
	group := e.Group(basePath+"/notification-channels", middleware.Auth)
	group.GET("", h.GetAllChannels)
	group.GET("/:id", h.GetChannel)
	group.POST("", h.CreateChannel)
	group.PUT("/:id", h.UpdateChannel)
	group.DELETE("/:id", h.DeleteChannel)
}

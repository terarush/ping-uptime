package handler

import (
	"fmt"
	"ping-uptime/internal/pkg/bus"
	"ping-uptime/internal/pkg/logger"
	"ping-uptime/internal/pkg/middleware"
	"ping-uptime/internal/pkg/utils"
	"ping-uptime/modules/settings/domain/entity"
	"ping-uptime/modules/settings/domain/service"
	"ping-uptime/modules/settings/dto/request"
	"ping-uptime/modules/settings/dto/response"

	"github.com/labstack/echo/v4"
)

type SettingHandler struct {
	settingService *service.SettingService
	log            *logger.Logger
	event          *bus.EventBus
	r              *utils.Response
}

func NewSettingHandler(log *logger.Logger, event *bus.EventBus, settingService *service.SettingService) *SettingHandler {
	return &SettingHandler{
		settingService: settingService,
		log:            log,
		event:          event,
		r:              &utils.Response{},
	}
}

func (h *SettingHandler) getAuthUser(c echo.Context) (uint, string, error) {
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

func (h *SettingHandler) GetAllSettings(c echo.Context) error {
	ctx := c.Request().Context()
	_, _, err := h.getAuthUser(c)
	if err != nil {
		return h.r.UnauthorizedResponse(c, err.Error())
	}

	settings, err := h.settingService.GetAllSettings(ctx)
	if err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	return h.r.SuccessResponse(c, response.FromEntities(settings), "Settings retrieved successfully")
}

func (h *SettingHandler) GetSetting(c echo.Context) error {
	ctx := c.Request().Context()
	_, _, err := h.getAuthUser(c)
	if err != nil {
		return h.r.UnauthorizedResponse(c, err.Error())
	}

	key := c.Param("key")
	if key == "" {
		return h.r.BadRequestResponse(c, "Key is required")
	}

	setting, err := h.settingService.GetSettingByKey(ctx, key)
	if err != nil {
		return h.r.NotFoundResponse(c, "Setting not found")
	}

	return h.r.SuccessResponse(c, response.FromEntity(setting), "Setting retrieved successfully")
}

func (h *SettingHandler) SaveSetting(c echo.Context) error {
	ctx := c.Request().Context()
	_, role, err := h.getAuthUser(c)
	if err != nil {
		return h.r.UnauthorizedResponse(c, err.Error())
	}

	if role != "admin" {
		return h.r.ForbiddenResponse(c, "Only administrators can save settings")
	}

	req := new(request.SaveSettingRequest)
	if err := c.Bind(req); err != nil {
		return h.r.BadRequestResponse(c, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return h.r.BadRequestResponse(c, err.Error())
	}

	setting := entity.NewSetting(req.Key, req.Value, req.Description)
	err = h.settingService.SetSetting(ctx, setting)
	if err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	h.event.Publish(bus.Event{Type: "setting.saved", Payload: setting})

	return h.r.SuccessResponse(c, response.FromEntity(setting), "Setting saved successfully")
}

func (h *SettingHandler) DeleteSetting(c echo.Context) error {
	ctx := c.Request().Context()
	_, role, err := h.getAuthUser(c)
	if err != nil {
		return h.r.UnauthorizedResponse(c, err.Error())
	}

	if role != "admin" {
		return h.r.ForbiddenResponse(c, "Only administrators can delete settings")
	}

	key := c.Param("key")
	if key == "" {
		return h.r.BadRequestResponse(c, "Key is required")
	}

	_, err = h.settingService.GetSettingByKey(ctx, key)
	if err != nil {
		return h.r.NotFoundResponse(c, "Setting not found")
	}

	err = h.settingService.DeleteSetting(ctx, key)
	if err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	h.event.Publish(bus.Event{Type: "setting.deleted", Payload: key})

	return h.r.NoContentResponse(c)
}

func (h *SettingHandler) RegisterRoutes(e *echo.Echo, basePath string) {
	group := e.Group(basePath+"/settings", middleware.Auth)
	group.GET("", h.GetAllSettings)
	group.GET("/:key", h.GetSetting)
	group.POST("", h.SaveSetting)
	group.DELETE("/:key", h.DeleteSetting)
}

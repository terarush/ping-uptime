package handler

import (
	"context"
	"ping-uptime/internal/pkg/bus"
	"ping-uptime/internal/pkg/logger"
	"ping-uptime/internal/pkg/middleware"
	"ping-uptime/internal/pkg/utils"
	"ping-uptime/modules/analytics/domain/entity"
	"ping-uptime/modules/analytics/domain/service"
	"strconv"

	"github.com/labstack/echo/v4"
)

type AnalyticsHandler struct {
	log            *logger.Logger
	event          *bus.EventBus
	analyticsService *service.AnalyticsService
	r              *utils.Response
}

func NewAnalyticsHandler(log *logger.Logger, event *bus.EventBus, analyticsService *service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		log:            log,
		event:          event,
		analyticsService: analyticsService,
		r:              &utils.Response{},
	}
}

func (h *AnalyticsHandler) getAuthUser(c echo.Context) (uint, string, error) {
	userClaims, ok := c.Get("user").(map[string]interface{})
	if !ok {
		return 0, "", nil
	}

	userIDVal, ok := userClaims["user_id"]
	if !ok {
		return 0, "", nil
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
		return 0, "", nil
	}

	roleVal, _ := userClaims["role"].(string)
	if roleVal == "" {
		roleVal = "user"
	}

	return userID, roleVal, nil
}

func (h *AnalyticsHandler) GetMonitorChart(c echo.Context) error {
	ctx := c.Request().Context()
	userID, role, err := h.getAuthUser(c)
	if err != nil {
		return h.r.UnauthorizedResponse(c, err.Error())
	}

	monitorID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "Invalid monitor ID")
	}

	window := c.QueryParam("window")
	if window == "" {
		window = "1m"
	}

	points, err := h.analyticsService.GetChartData(ctx, uint(monitorID), window)
	if err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	var monitorName, monitorURL string
	if role != "admin" {
		var count int64
		database.DB.WithContext(ctx).Model(nil).Where("id = ? AND user_id = ?", monitorID, userID).Count(&count)
		if count == 0 {
			return h.r.ForbiddenResponse(c, "You do not have access to this monitor")
		}
	}

	return h.r.SuccessResponse(c, map[string]interface{}{
		"monitor_id": monitorID,
		"window":     window,
		"data":       points,
	}, "Chart data retrieved successfully")
}

func (h *AnalyticsHandler) GetDashboardStats(c echo.Context) error {
	ctx := c.Request().Context()
	userID, role, err := h.getAuthUser(c)
	if err != nil {
		return h.r.UnauthorizedResponse(c, err.Error())
	}

	window := c.QueryParam("window")
	if window == "" {
		window = "1m"
	}

	var targetUserID uint = userID
	if role == "admin" {
		targetUserID = 0
	}

	stats, err := h.analyticsService.GetMonitorStats(ctx, targetUserID, window)
	if err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	return h.r.SuccessResponse(c, map[string]interface{}{
		"window": window,
		"data":   stats,
	}, "Dashboard stats retrieved successfully")
}

func (h *AnalyticsHandler) RegisterRoutes(e *echo.Echo, basePath string) {
	group := e.Group(basePath+"/analytics", middleware.Auth)
	group.GET("/monitors/:id/chart", h.GetMonitorChart)
	group.GET("/dashboard", h.GetDashboardStats)
}

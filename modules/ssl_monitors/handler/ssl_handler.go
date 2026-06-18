package handler

import (
	"fmt"
	"ping-uptime/internal/pkg/logger"
	"ping-uptime/internal/pkg/middleware"
	"ping-uptime/internal/pkg/utils"
	"ping-uptime/modules/ssl_monitors/domain/service"
	monitorService "ping-uptime/modules/monitors/domain/service"
	"strconv"

	"github.com/labstack/echo/v4"
)

type SSLHandler struct {
	sslService     *service.SSLService
	monitorService *monitorService.MonitorService
	log            *logger.Logger
	r              *utils.Response
}

func NewSSLHandler(log *logger.Logger, sslService *service.SSLService, monitorService *monitorService.MonitorService) *SSLHandler {
	return &SSLHandler{
		sslService:     sslService,
		monitorService: monitorService,
		log:            log,
		r:              &utils.Response{},
	}
}

func (h *SSLHandler) getAuthUser(c echo.Context) (uint, string, error) {
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

// GetAll returns all SSL certificate records
func (h *SSLHandler) GetAll(c echo.Context) error {
	ctx := c.Request().Context()

	certs, err := h.sslService.GetAll(ctx)
	if err != nil {
		h.log.Error("Failed to fetch SSL certs: %v", err)
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	return h.r.SuccessResponse(c, certs, "SSL certificates retrieved successfully")
}

// GetByID returns a single SSL certificate record by ID
func (h *SSLHandler) GetByID(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "Invalid SSL certificate ID")
	}

	cert, err := h.sslService.GetByID(ctx, uint(id))
	if err != nil {
		return h.r.NotFoundResponse(c, "SSL certificate not found")
	}

	return h.r.SuccessResponse(c, cert, "SSL certificate retrieved successfully")
}

// CheckMonitor triggers an SSL check for a specific monitor
func (h *SSLHandler) CheckMonitor(c echo.Context) error {
	ctx := c.Request().Context()
	userID, role, err := h.getAuthUser(c)
	if err != nil {
		return h.r.UnauthorizedResponse(c, err.Error())
	}

	monitorID, err := strconv.ParseUint(c.Param("monitorId"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "Invalid monitor ID")
	}

	// Validate monitor exists and check ownership
	monitor, err := h.monitorService.GetMonitorByID(ctx, uint(monitorID))
	if err != nil {
		return h.r.NotFoundResponse(c, "Monitor not found")
	}

	if role != "admin" && monitor.UserID != userID {
		return h.r.ForbiddenResponse(c, "You do not have access to this monitor")
	}

	cert, err := h.sslService.CheckSSL(ctx, uint(monitorID), monitor.URL)
	if err != nil {
		h.log.Error("SSL check failed for monitor %d: %v", monitorID, err)
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	return h.r.SuccessResponse(c, cert, "SSL check completed")
}

// CheckAll triggers SSL checks for all active monitors
func (h *SSLHandler) CheckAll(c echo.Context) error {
	ctx := c.Request().Context()

	monitors, err := h.monitorService.GetAllMonitors(ctx)
	if err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	type checkResult struct {
		MonitorID uint   `json:"monitor_id"`
		Domain    string `json:"domain"`
		Status    string `json:"status"`
		Error     string `json:"error,omitempty"`
	}

	results := make([]checkResult, 0, len(monitors))
	for _, mon := range monitors {
		cert, err := h.sslService.CheckSSL(ctx, mon.ID, mon.URL)
		if err != nil {
			results = append(results, checkResult{
				MonitorID: mon.ID,
				Domain:    mon.URL,
				Status:    "error",
				Error:     err.Error(),
			})
			continue
		}
		results = append(results, checkResult{
			MonitorID: mon.ID,
			Domain:    cert.Domain,
			Status:    cert.Status,
		})
	}

	return h.r.SuccessResponse(c, results, "SSL checks completed for all monitors")
}

// GetExpiring returns certificates expiring within N days
func (h *SSLHandler) GetExpiring(c echo.Context) error {
	ctx := c.Request().Context()

	daysStr := c.QueryParam("days")
	days := 30
	if daysStr != "" {
		d, err := strconv.Atoi(daysStr)
		if err == nil && d > 0 {
			days = d
		}
	}

	certs, err := h.sslService.GetExpiring(ctx, days)
	if err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	return h.r.SuccessResponse(c, certs, "Expiring SSL certificates retrieved successfully")
}

// Delete removes an SSL certificate record
func (h *SSLHandler) Delete(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "Invalid SSL certificate ID")
	}

	_, err = h.sslService.GetByID(ctx, uint(id))
	if err != nil {
		return h.r.NotFoundResponse(c, "SSL certificate not found")
	}

	if err := h.sslService.Delete(ctx, uint(id)); err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	return h.r.NoContentResponse(c)
}

// RegisterRoutes registers SSL certificate monitor routes
func (h *SSLHandler) RegisterRoutes(e *echo.Echo, basePath string) {
	group := e.Group(basePath+"/ssl-monitors", middleware.Auth)
	group.GET("", h.GetAll)
	group.GET("/:id", h.GetByID)
	group.POST("/check/:monitorId", h.CheckMonitor)
	group.POST("/check-all", h.CheckAll)
	group.GET("/expiring", h.GetExpiring)
	group.DELETE("/:id", h.Delete)
}

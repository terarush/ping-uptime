package handler

import (
	"net/http"
	"ping-uptime/internal/pkg/logger"
	"ping-uptime/internal/pkg/middleware"
	"ping-uptime/internal/pkg/utils"
	"ping-uptime/modules/maintenances/domain/entity"
	"ping-uptime/modules/maintenances/domain/service"
	"ping-uptime/modules/maintenances/dto/request"
	"ping-uptime/modules/maintenances/dto/response"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type MaintenanceHandler struct {
	svc *service.MaintenanceService
	log *logger.Logger
	r   *utils.Response
}

func NewMaintenanceHandler(log *logger.Logger, svc *service.MaintenanceService) *MaintenanceHandler {
	return &MaintenanceHandler{svc: svc, log: log, r: &utils.Response{}}
}

func (h *MaintenanceHandler) getAuthUser(c echo.Context) (uint, string, error) {
	userClaims, ok := c.Get("user").(map[string]interface{})
	if !ok {
		return 0, "", echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}
	userIDVal := userClaims["user_id"]
	var userID uint
	switch v := userIDVal.(type) {
	case float64:
		userID = uint(v)
	case int64:
		userID = uint(v)
	case int:
		userID = uint(v)
	default:
		return 0, "", echo.NewHTTPError(http.StatusUnauthorized, "invalid user_id")
	}
	roleVal, _ := userClaims["role"].(string)
	return userID, roleVal, nil
}

func (h *MaintenanceHandler) GetAll(c echo.Context) error {
	ctx := c.Request().Context()
	userID, role, err := h.getAuthUser(c)
	if err != nil {
		return h.r.UnauthorizedResponse(c, err.Error())
	}

	var items []*entity.Maintenance
	if role == "admin" {
		items, err = h.svc.GetAll(ctx)
	} else {
		items, err = h.svc.GetByUserID(ctx, userID)
	}
	if err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	// Populate monitor IDs for each maintenance
	resp := response.FromEntities(items)
	for i, m := range items {
		monitorIDs, err := h.svc.GetMonitorIDs(ctx, m.ID)
		if err != nil {
			h.log.Error("Failed to get monitor IDs:", err)
		}
		resp[i].MonitorIDs = monitorIDs
	}

	return h.r.SuccessResponse(c, resp, "Maintenances retrieved")
}

func (h *MaintenanceHandler) GetByID(c echo.Context) error {
	ctx := c.Request().Context()
	_, _, err := h.getAuthUser(c)
	if err != nil {
		return h.r.UnauthorizedResponse(c, err.Error())
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "Invalid ID")
	}

	m, err := h.svc.GetByID(ctx, uint(id))
	if err != nil {
		return h.r.NotFoundResponse(c, "Maintenance not found")
	}

	resp := response.FromEntity(m)
	monitorIDs, err := h.svc.GetMonitorIDs(ctx, m.ID)
	if err != nil {
		h.log.Error("Failed to get monitor IDs:", err)
	}
	resp.MonitorIDs = monitorIDs
	return h.r.SuccessResponse(c, resp, "Maintenance retrieved")
}

func (h *MaintenanceHandler) Create(c echo.Context) error {
	ctx := c.Request().Context()
	userID, _, err := h.getAuthUser(c)
	if err != nil {
		return h.r.UnauthorizedResponse(c, err.Error())
	}

	req := new(request.CreateMaintenanceRequest)
	if err := c.Bind(req); err != nil {
		return h.r.BadRequestResponse(c, err.Error())
	}
	if err := c.Validate(req); err != nil {
		return h.r.BadRequestResponse(c, err.Error())
	}

	startAt, err := time.Parse(time.RFC3339, req.StartAt)
	if err != nil {
		return h.r.BadRequestResponse(c, "Invalid start_at format, use RFC3339")
	}
	endAt, err := time.Parse(time.RFC3339, req.EndAt)
	if err != nil {
		return h.r.BadRequestResponse(c, "Invalid end_at format, use RFC3339")
	}

	now := time.Now()
	m := &entity.Maintenance{
		Name:        req.Name,
		Description: req.Description,
		StartAt:     startAt,
		EndAt:       endAt,
		UserID:      userID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Auto-set status
	if startAt.After(now) {
		m.Status = "scheduled"
	} else if now.After(startAt) && now.Before(endAt) {
		m.Status = "ongoing"
	} else {
		m.Status = "completed"
	}

	if err := h.svc.Create(ctx, m); err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	if len(req.MonitorIDs) > 0 {
		if err := h.svc.SetMonitorIDs(ctx, m.ID, req.MonitorIDs); err != nil {
			h.log.Error("Failed to set monitor IDs:", err)
		}
	}

	resp := response.FromEntity(m)
	resp.MonitorIDs = req.MonitorIDs
	return h.r.CreatedResponse(c, resp, "Maintenance created")
}

func (h *MaintenanceHandler) Update(c echo.Context) error {
	ctx := c.Request().Context()
	_, _, err := h.getAuthUser(c)
	if err != nil {
		return h.r.UnauthorizedResponse(c, err.Error())
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "Invalid ID")
	}

	m, err := h.svc.GetByID(ctx, uint(id))
	if err != nil {
		return h.r.NotFoundResponse(c, "Maintenance not found")
	}

	req := new(request.UpdateMaintenanceRequest)
	if err := c.Bind(req); err != nil {
		return h.r.BadRequestResponse(c, err.Error())
	}
	if err := c.Validate(req); err != nil {
		return h.r.BadRequestResponse(c, err.Error())
	}

	startAt, err := time.Parse(time.RFC3339, req.StartAt)
	if err != nil {
		return h.r.BadRequestResponse(c, "Invalid start_at format")
	}
	endAt, err := time.Parse(time.RFC3339, req.EndAt)
	if err != nil {
		return h.r.BadRequestResponse(c, "Invalid end_at format")
	}

	m.Name = req.Name
	m.Description = req.Description
	m.StartAt = startAt
	m.EndAt = endAt
	m.UpdatedAt = time.Now()

	now := time.Now()
	if startAt.After(now) {
		m.Status = "scheduled"
	} else if now.After(startAt) && now.Before(endAt) {
		m.Status = "ongoing"
	} else {
		m.Status = "completed"
	}

	if err := h.svc.Update(ctx, m); err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	if err := h.svc.SetMonitorIDs(ctx, m.ID, req.MonitorIDs); err != nil {
		h.log.Error("Failed to set monitor IDs:", err)
	}

	resp := response.FromEntity(m)
	resp.MonitorIDs = req.MonitorIDs
	return h.r.SuccessResponse(c, resp, "Maintenance updated")
}

func (h *MaintenanceHandler) Delete(c echo.Context) error {
	ctx := c.Request().Context()
	_, _, err := h.getAuthUser(c)
	if err != nil {
		return h.r.UnauthorizedResponse(c, err.Error())
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "Invalid ID")
	}

	if err := h.svc.Delete(ctx, uint(id)); err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	return h.r.NoContentResponse(c)
}

func (h *MaintenanceHandler) RegisterRoutes(e *echo.Echo, basePath string) {
	group := e.Group(basePath+"/maintenances", middleware.Auth)
	group.GET("", h.GetAll)
	group.GET("/:id", h.GetByID)
	group.POST("", h.Create)
	group.PUT("/:id", h.Update)
	group.DELETE("/:id", h.Delete)
}

package handler

import (
	"fmt"
	"ping-uptime/internal/pkg/bus"
	"ping-uptime/internal/pkg/logger"
	"ping-uptime/internal/pkg/middleware"
	"ping-uptime/internal/pkg/utils"
	"ping-uptime/modules/incidents/domain/entity"
	"ping-uptime/modules/incidents/domain/service"
	"ping-uptime/modules/incidents/dto/request"
	"ping-uptime/modules/incidents/dto/response"
	monitorService "ping-uptime/modules/monitors/domain/service"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type IncidentHandler struct {
	incidentService *service.IncidentService
	monitorService  *monitorService.MonitorService
	log             *logger.Logger
	event           *bus.EventBus
	r               *utils.Response
}

func NewIncidentHandler(
	log *logger.Logger,
	event *bus.EventBus,
	incidentService *service.IncidentService,
	monitorService *monitorService.MonitorService,
) *IncidentHandler {
	return &IncidentHandler{
		incidentService: incidentService,
		monitorService:  monitorService,
		log:             log,
		event:           event,
		r:               &utils.Response{},
	}
}

func (h *IncidentHandler) getAuthUser(c echo.Context) (uint, string, error) {
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

func (h *IncidentHandler) GetAllIncidents(c echo.Context) error {
	ctx := c.Request().Context()
	userID, role, err := h.getAuthUser(c)
	if err != nil {
		return h.r.UnauthorizedResponse(c, err.Error())
	}

	var incidents []*entity.Incident
	var hErr error
	if role == "admin" {
		incidents, hErr = h.incidentService.GetAllIncidents(ctx)
	} else {
		incidents, hErr = h.incidentService.GetIncidentsByUserID(ctx, userID)
	}

	if hErr != nil {
		return h.r.InternalServerErrorResponse(c, hErr.Error())
	}

	return h.r.SuccessResponse(c, response.FromEntities(incidents), "Incidents retrieved successfully")
}

func (h *IncidentHandler) GetIncident(c echo.Context) error {
	ctx := c.Request().Context()
	userID, role, err := h.getAuthUser(c)
	if err != nil {
		return h.r.UnauthorizedResponse(c, err.Error())
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "Invalid incident ID")
	}

	incident, err := h.incidentService.GetIncidentByID(ctx, uint(id))
	if err != nil {
		return h.r.NotFoundResponse(c, "Incident not found")
	}

	if role != "admin" && incident.UserID != userID {
		return h.r.ForbiddenResponse(c, "You do not have access to this incident")
	}

	return h.r.SuccessResponse(c, response.FromEntity(incident), "Incident retrieved successfully")
}

func (h *IncidentHandler) GetIncidentsByMonitor(c echo.Context) error {
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
		return h.r.ForbiddenResponse(c, "You do not have access to this monitor's incidents")
	}

	incidents, err := h.incidentService.GetIncidentsByMonitorID(ctx, uint(monitorID))
	if err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	return h.r.SuccessResponse(c, response.FromEntities(incidents), "Monitor incidents retrieved successfully")
}

func (h *IncidentHandler) CreateIncident(c echo.Context) error {
	ctx := c.Request().Context()
	userID, role, err := h.getAuthUser(c)
	if err != nil {
		return h.r.UnauthorizedResponse(c, err.Error())
	}

	req := new(request.CreateIncidentRequest)
	if err := c.Bind(req); err != nil {
		return h.r.BadRequestResponse(c, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return h.r.BadRequestResponse(c, err.Error())
	}

	// Check monitor ownership
	monitor, err := h.monitorService.GetMonitorByID(ctx, req.MonitorID)
	if err != nil {
		return h.r.BadRequestResponse(c, "Invalid monitor ID: monitor does not exist")
	}

	if role != "admin" && monitor.UserID != userID {
		return h.r.ForbiddenResponse(c, "You do not have permission to create an incident for this monitor")
	}

	incident := entity.NewIncident(req.MonitorID, monitor.UserID, req.Status, req.ErrorMessage, req.Latency)
	err = h.incidentService.CreateIncident(ctx, incident)
	if err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	h.event.Publish(bus.Event{Type: "incident.created", Payload: incident})

	return h.r.CreatedResponse(c, response.FromEntity(incident), "Incident created successfully")
}

func (h *IncidentHandler) UpdateIncident(c echo.Context) error {
	ctx := c.Request().Context()
	userID, role, err := h.getAuthUser(c)
	if err != nil {
		return h.r.UnauthorizedResponse(c, err.Error())
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "Invalid incident ID")
	}

	incident, err := h.incidentService.GetIncidentByID(ctx, uint(id))
	if err != nil {
		return h.r.NotFoundResponse(c, "Incident not found")
	}

	if role != "admin" && incident.UserID != userID {
		return h.r.ForbiddenResponse(c, "You do not have permission to update this incident")
	}

	req := new(request.UpdateIncidentRequest)
	if err := c.Bind(req); err != nil {
		return h.r.BadRequestResponse(c, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return h.r.BadRequestResponse(c, err.Error())
	}

	oldStatus := incident.Status
	incident.Status = req.Status
	incident.ErrorMessage = req.ErrorMessage
	incident.Latency = req.Latency

	if req.Status == "resolved" && oldStatus != "resolved" {
		now := time.Now()
		incident.ResolvedAt = &now
	}

	err = h.incidentService.UpdateIncident(ctx, incident)
	if err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	h.event.Publish(bus.Event{Type: "incident.updated", Payload: incident})

	return h.r.SuccessResponse(c, response.FromEntity(incident), "Incident updated successfully")
}

func (h *IncidentHandler) DeleteIncident(c echo.Context) error {
	ctx := c.Request().Context()
	userID, role, err := h.getAuthUser(c)
	if err != nil {
		return h.r.UnauthorizedResponse(c, err.Error())
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "Invalid incident ID")
	}

	incident, err := h.incidentService.GetIncidentByID(ctx, uint(id))
	if err != nil {
		return h.r.NotFoundResponse(c, "Incident not found")
	}

	if role != "admin" && incident.UserID != userID {
		return h.r.ForbiddenResponse(c, "You do not have permission to delete this incident")
	}

	err = h.incidentService.DeleteIncident(ctx, uint(id))
	if err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	h.event.Publish(bus.Event{Type: "incident.deleted", Payload: incident})

	return h.r.NoContentResponse(c)
}

func (h *IncidentHandler) RegisterRoutes(e *echo.Echo, basePath string) {
	group := e.Group(basePath+"/incidents", middleware.Auth)
	group.GET("", h.GetAllIncidents)
	group.GET("/:id", h.GetIncident)
	group.GET("/monitor/:monitorId", h.GetIncidentsByMonitor)
	group.POST("", h.CreateIncident)
	group.PUT("/:id", h.UpdateIncident)
	group.DELETE("/:id", h.DeleteIncident)
}

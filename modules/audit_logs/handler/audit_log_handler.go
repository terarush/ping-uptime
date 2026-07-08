package handler

import (
	"context"
	"encoding/json"
	"ping-uptime/internal/pkg/bus"
	"ping-uptime/internal/pkg/logger"
	"ping-uptime/internal/pkg/middleware"
	"ping-uptime/internal/pkg/utils"
	"ping-uptime/modules/audit_logs/domain/repository"
	"ping-uptime/modules/audit_logs/domain/service"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type AuditLogHandler struct {
	svc *service.AuditLogService
	log *logger.Logger
	r   *utils.Response
}

func NewAuditLogHandler(log *logger.Logger, svc *service.AuditLogService) *AuditLogHandler {
	return &AuditLogHandler{svc: svc, log: log, r: &utils.Response{}}
}

// @Summary      [Admin] List audit logs
// @Description  Retrieve audit logs with optional filtering. Requires admin privileges.
// @Tags         AuditLogs
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        user_id      query  int     false  "Filter by user ID"
// @Param        entity_type  query  string  false  "Filter by entity type (e.g. monitor, incident)"
// @Param        action       query  string  false  "Filter by action (e.g. created, updated, deleted)"
// @Param        from         query  string  false  "Start date/time (RFC3339)"
// @Param        to           query  string  false  "End date/time (RFC3339)"
// @Param        limit        query  int     false  "Max results to return"
// @Success      200  {object}  utils.SuccessResponseModel
// @Failure      400  {object}  utils.ErrorResponseModel
// @Failure      401  {object}  utils.ErrorResponseModel
// @Failure      403  {object}  utils.ErrorResponseModel
// @Failure      500  {object}  utils.ErrorResponseModel
// @Router       /api/audit-logs [get]
func (h *AuditLogHandler) getAll(c echo.Context) error {
	ctx := c.Request().Context()

	filter := repository.AuditLogFilter{}

	if uidStr := c.QueryParam("user_id"); uidStr != "" {
		if uid, err := strconv.ParseUint(uidStr, 10, 32); err == nil {
			v := uint(uid)
			filter.UserID = &v
		}
	}
	if et := c.QueryParam("entity_type"); et != "" {
		filter.EntityType = &et
	}
	if a := c.QueryParam("action"); a != "" {
		filter.Action = &a
	}
	if f := c.QueryParam("from"); f != "" {
		filter.From = &f
	}
	if t := c.QueryParam("to"); t != "" {
		filter.To = &t
	}
	if lStr := c.QueryParam("limit"); lStr != "" {
		if l, err := strconv.Atoi(lStr); err == nil {
			filter.Limit = l
		}
	}

	items, err := h.svc.FindFiltered(ctx, filter)
	if err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.SuccessResponse(c, items, "Audit logs retrieved")
}

// @Summary      [Admin] Search audit logs
// @Description  Retrieve audit logs filtered by entity type and entity ID. Requires admin privileges.
// @Tags         AuditLogs
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        entity_type  query  string  true  "Entity type (e.g. monitor, incident)"
// @Param        entity_id    query  int     true  "Entity ID"
// @Success      200  {object}  utils.SuccessResponseModel
// @Failure      400  {object}  utils.ErrorResponseModel
// @Failure      401  {object}  utils.ErrorResponseModel
// @Failure      403  {object}  utils.ErrorResponseModel
// @Failure      500  {object}  utils.ErrorResponseModel
// @Router       /api/audit-logs/search [get]
func (h *AuditLogHandler) getByEntity(c echo.Context) error {
	ctx := c.Request().Context()
	entityType := c.QueryParam("entity_type")
	entityIDStr := c.QueryParam("entity_id")
	if entityType == "" || entityIDStr == "" {
		return h.r.BadRequestResponse(c, "entity_type and entity_id required")
	}
	entityID, err := strconv.ParseUint(entityIDStr, 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "Invalid entity_id")
	}
	items, err := h.svc.GetByEntity(ctx, entityType, uint(entityID))
	if err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.SuccessResponse(c, items, "Audit logs retrieved")
}

func (h *AuditLogHandler) RegisterRoutes(e *echo.Echo, basePath string) {
	group := e.Group(basePath+"/audit-logs", middleware.Auth, middleware.Admin)
	group.GET("", h.getAll)
	group.GET("/search", h.getByEntity)
}

func (h *AuditLogHandler) SubscribeToEvents(event *bus.EventBus, svc *service.AuditLogService) {
	subscribe := func(eventType string) {
		event.SubscribeFunc(eventType, func(ev bus.Event) {
			ctx := context.Background()
			payloadJSON, _ := json.Marshal(ev.Payload)
			details := string(payloadJSON)

			parts := splitEventType(ev.Type)
			if len(parts) < 2 {
				return
			}
			entityType := parts[0]
			action := parts[1]

			var pm map[string]interface{}
			if err := json.Unmarshal(payloadJSON, &pm); err != nil {
				return
			}

			entityID := uint(0)
			if id, ok := pm["id"].(float64); ok {
				entityID = uint(id)
			}

			userID := uint(0)
			if uid, ok := pm["user_id"].(float64); ok {
				userID = uint(uid)
			}

			if err := svc.Log(ctx, userID, action, entityType, entityID, details); err != nil {
				h.log.Error("failed to log audit event", "type", ev.Type, "error", err)
			}
		})
	}

	events := []string{
		"monitor.created", "monitor.updated", "monitor.deleted",
		"incident.created", "incident.resolved",
		"status_page.created", "status_page.updated", "status_page.deleted",
		"notification_channel.created", "notification_channel.updated", "notification_channel.deleted",
	}

	for _, e := range events {
		subscribe(e)
	}
}

func splitEventType(t string) []string {
	idx := strings.Index(t, ".")
	if idx == -1 {
		return []string{t}
	}
	return []string{t[:idx], t[idx+1:]}
}

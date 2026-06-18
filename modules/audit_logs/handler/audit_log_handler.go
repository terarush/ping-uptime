package handler

import (
	"context"
	"encoding/json"
	"ping-uptime/internal/pkg/bus"
	"ping-uptime/internal/pkg/logger"
	"ping-uptime/internal/pkg/middleware"
	"ping-uptime/internal/pkg/utils"
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

func (h *AuditLogHandler) getAll(c echo.Context) error {
	ctx := c.Request().Context()
	items, err := h.svc.GetAll(ctx)
	if err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.SuccessResponse(c, items, "Audit logs retrieved")
}

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

			// Extract entity ID from payload
			payloadMap, ok := ev.Payload.(map[string]interface{})
			if !ok {
				// Try to marshal/unmarshal to get generic map
				var pm map[string]interface{}
				if err := json.Unmarshal(payloadJSON, &pm); err != nil {
					return
				}
				payloadMap = pm
			}

			entityID := uint(0)
			if id, ok := payloadMap["id"].(float64); ok {
				entityID = uint(id)
			}

			if err := svc.Log(ctx, 0, action, entityType, entityID, details); err != nil {
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

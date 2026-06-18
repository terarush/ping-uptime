package handler

import (
	"fmt"
	"ping-uptime/internal/pkg/bus"
	"ping-uptime/internal/pkg/logger"
	"ping-uptime/internal/pkg/middleware"
	"ping-uptime/internal/pkg/utils"
	monitorEntity "ping-uptime/modules/monitors/domain/entity"
	monitorService "ping-uptime/modules/monitors/domain/service"
	"ping-uptime/modules/status_pages/domain/entity"
	"ping-uptime/modules/status_pages/domain/service"
	"ping-uptime/modules/status_pages/dto/request"
	"ping-uptime/modules/status_pages/dto/response"
	"strconv"

	"github.com/labstack/echo/v4"
)

type StatusPageHandler struct {
	statusPageService *service.StatusPageService
	monitorService    *monitorService.MonitorService
	log               *logger.Logger
	event             *bus.EventBus
	r                 *utils.Response
}

func NewStatusPageHandler(
	log *logger.Logger,
	event *bus.EventBus,
	statusPageService *service.StatusPageService,
	monitorService *monitorService.MonitorService,
) *StatusPageHandler {
	return &StatusPageHandler{
		statusPageService: statusPageService,
		monitorService:    monitorService,
		log:               log,
		event:             event,
		r:                 &utils.Response{},
	}
}

func (h *StatusPageHandler) getAuthUser(c echo.Context) (uint, string, error) {
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

func (h *StatusPageHandler) GetPublicStatusPage(c echo.Context) error {
	ctx := c.Request().Context()
	slug := c.Param("slug")
	if slug == "" {
		return h.r.BadRequestResponse(c, "Slug is required")
	}

	page, err := h.statusPageService.GetStatusPageBySlug(ctx, slug)
	if err != nil {
		return h.r.NotFoundResponse(c, "Status page not found")
	}

	return h.r.SuccessResponse(c, response.FromEntity(page), "Status page retrieved successfully")
}

func (h *StatusPageHandler) GetAllStatusPages(c echo.Context) error {
	ctx := c.Request().Context()
	userID, role, err := h.getAuthUser(c)
	if err != nil {
		return h.r.UnauthorizedResponse(c, err.Error())
	}

	var pages []*entity.StatusPage
	if role == "admin" {
		pages, err = h.statusPageService.GetAllStatusPages(ctx)
	} else {
		pages, err = h.statusPageService.GetStatusPagesByUserID(ctx, userID)
	}

	if err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	return h.r.SuccessResponse(c, response.FromEntities(pages), "Status pages retrieved successfully")
}

func (h *StatusPageHandler) GetStatusPage(c echo.Context) error {
	ctx := c.Request().Context()
	userID, role, err := h.getAuthUser(c)
	if err != nil {
		return h.r.UnauthorizedResponse(c, err.Error())
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "Invalid status page ID")
	}

	page, err := h.statusPageService.GetStatusPageByID(ctx, uint(id))
	if err != nil {
		return h.r.NotFoundResponse(c, "Status page not found")
	}

	if role != "admin" && page.UserID != userID {
		return h.r.ForbiddenResponse(c, "You do not have access to this status page")
	}

	return h.r.SuccessResponse(c, response.FromEntity(page), "Status page retrieved successfully")
}

func (h *StatusPageHandler) CreateStatusPage(c echo.Context) error {
	ctx := c.Request().Context()
	userID, role, err := h.getAuthUser(c)
	if err != nil {
		return h.r.UnauthorizedResponse(c, err.Error())
	}

	req := new(request.CreateStatusPageRequest)
	if err := c.Bind(req); err != nil {
		return h.r.BadRequestResponse(c, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return h.r.BadRequestResponse(c, err.Error())
	}

	var monitors []*monitorEntity.Monitor
	for _, mID := range req.MonitorIDs {
		mon, err := h.monitorService.GetMonitorByID(ctx, mID)
		if err != nil {
			return h.r.BadRequestResponse(c, fmt.Sprintf("Monitor with ID %d not found", mID))
		}
		if role != "admin" && mon.UserID != userID {
			return h.r.ForbiddenResponse(c, fmt.Sprintf("You do not own monitor with ID %d", mID))
		}
		monitors = append(monitors, mon)
	}

	page := entity.NewStatusPage(req.Name, req.Slug, req.Description, userID)
	page.Monitors = monitors

	err = h.statusPageService.CreateStatusPage(ctx, page)
	if err != nil {
		if err == service.ErrSlugAlreadyTaken {
			return h.r.BadRequestResponse(c, err.Error())
		}
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	h.event.Publish(bus.Event{Type: "status_page.created", Payload: page})

	return h.r.CreatedResponse(c, response.FromEntity(page), "Status page created successfully")
}

func (h *StatusPageHandler) UpdateStatusPage(c echo.Context) error {
	ctx := c.Request().Context()
	userID, role, err := h.getAuthUser(c)
	if err != nil {
		return h.r.UnauthorizedResponse(c, err.Error())
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "Invalid status page ID")
	}

	page, err := h.statusPageService.GetStatusPageByID(ctx, uint(id))
	if err != nil {
		return h.r.NotFoundResponse(c, "Status page not found")
	}

	if role != "admin" && page.UserID != userID {
		return h.r.ForbiddenResponse(c, "You do not have permission to update this status page")
	}

	req := new(request.UpdateStatusPageRequest)
	if err := c.Bind(req); err != nil {
		return h.r.BadRequestResponse(c, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return h.r.BadRequestResponse(c, err.Error())
	}

	var monitors []*monitorEntity.Monitor
	for _, mID := range req.MonitorIDs {
		mon, err := h.monitorService.GetMonitorByID(ctx, mID)
		if err != nil {
			return h.r.BadRequestResponse(c, fmt.Sprintf("Monitor with ID %d not found", mID))
		}
		if role != "admin" && mon.UserID != userID {
			return h.r.ForbiddenResponse(c, fmt.Sprintf("You do not own monitor with ID %d", mID))
		}
		monitors = append(monitors, mon)
	}

	page.Name = req.Name
	page.Slug = req.Slug
	page.Description = req.Description
	page.Monitors = monitors

	err = h.statusPageService.UpdateStatusPage(ctx, page)
	if err != nil {
		if err == service.ErrSlugAlreadyTaken {
			return h.r.BadRequestResponse(c, err.Error())
		}
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	h.event.Publish(bus.Event{Type: "status_page.updated", Payload: page})

	return h.r.SuccessResponse(c, response.FromEntity(page), "Status page updated successfully")
}

func (h *StatusPageHandler) DeleteStatusPage(c echo.Context) error {
	ctx := c.Request().Context()
	userID, role, err := h.getAuthUser(c)
	if err != nil {
		return h.r.UnauthorizedResponse(c, err.Error())
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "Invalid status page ID")
	}

	page, err := h.statusPageService.GetStatusPageByID(ctx, uint(id))
	if err != nil {
		return h.r.NotFoundResponse(c, "Status page not found")
	}

	if role != "admin" && page.UserID != userID {
		return h.r.ForbiddenResponse(c, "You do not have permission to delete this status page")
	}

	err = h.statusPageService.DeleteStatusPage(ctx, uint(id))
	if err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	h.event.Publish(bus.Event{Type: "status_page.deleted", Payload: page})

	return h.r.NoContentResponse(c)
}

func (h *StatusPageHandler) BadgeSVG(c echo.Context) error {
	slug := c.Param("slug")
	if slug == "" {
		return c.String(400, "")
	}

	ctx := c.Request().Context()
	page, err := h.statusPageService.GetStatusPageBySlug(ctx, slug)
	if err != nil {
		return c.String(404, "")
	}

	label := "Operational"
	color := "#22c55e"

	if len(page.Monitors) > 0 {
		downCount := 0
		for _, m := range page.Monitors {
			if m.UptimeStatus == "down" {
				downCount++
			}
		}
		if downCount == len(page.Monitors) {
			label = "Outage"
			color = "#ef4444"
		} else if downCount > 0 {
			label = "Partial Outage"
			color = "#f59e0b"
		}
	}

	svg := fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="140" height="20">
  <rect rx="3" width="140" height="20" fill="#555"/>
  <rect rx="3" width="75" height="20" fill="#444"/>
  <rect rx="3" x="75" width="65" height="20" fill="%[1]s"/>
  <g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="11">
    <text x="38" y="14" fill="#010101" fill-opacity=".3">uptime</text>
    <text x="38" y="13">uptime</text>
    <text x="107" y="14" fill="#010101" fill-opacity=".3">%[2]s</text>
    <text x="107" y="13">%[2]s</text>
  </g>
</svg>`, color, label)

	c.Response().Header().Set(echo.HeaderContentType, "image/svg+xml")
	c.Response().Header().Set(echo.HeaderCacheControl, "public, max-age=60")
	return c.String(200, svg)
}

func (h *StatusPageHandler) RegisterRoutes(e *echo.Echo, basePath string) {
	// Public routes
	e.GET(basePath+"/status-pages/slug/:slug", h.GetPublicStatusPage)
	e.GET(basePath+"/status-pages/:slug/badge.svg", h.BadgeSVG)

	// Authenticated routes
	group := e.Group(basePath+"/status-pages", middleware.Auth)
	group.GET("", h.GetAllStatusPages)
	group.GET("/:id", h.GetStatusPage)
	group.POST("", h.CreateStatusPage)
	group.PUT("/:id", h.UpdateStatusPage)
	group.DELETE("/:id", h.DeleteStatusPage)
}

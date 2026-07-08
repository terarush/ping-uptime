package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"ping-uptime/internal/pkg/bus"
	"ping-uptime/internal/pkg/database"
	"ping-uptime/internal/pkg/logger"
	"ping-uptime/internal/pkg/middleware"
	"ping-uptime/internal/pkg/utils"
	incidentEntity "ping-uptime/modules/incidents/domain/entity"
	"ping-uptime/modules/monitors/domain/entity"
	"ping-uptime/modules/monitors/domain/service"
	"ping-uptime/modules/monitors/dto/request"
	"ping-uptime/modules/monitors/dto/response"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type ClientChan chan string

type SSEHub struct {
	clients    map[ClientChan]bool
	register   chan ClientChan
	unregister chan ClientChan
	broadcast  chan string
}

func NewSSEHub() *SSEHub {
	return &SSEHub{
		clients:    make(map[ClientChan]bool),
		register:   make(chan ClientChan),
		unregister: make(chan ClientChan),
		broadcast:  make(chan string, 100),
	}
}

func (h *SSEHub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			delete(h.clients, client)
			close(client)
		case msg := <-h.broadcast:
			for client := range h.clients {
				select {
				case client <- msg:
				default:
					// drop message if client buffer is full/blocked
				}
			}
		}
	}
}

type MonitorHandler struct {
	monitorService *service.MonitorService
	log            *logger.Logger
	event          *bus.EventBus
	r              *utils.Response
	sseHub         *SSEHub
}

func NewMonitorHandler(log *logger.Logger, event *bus.EventBus, monitorService *service.MonitorService) *MonitorHandler {
	hub := NewSSEHub()
	go hub.Run()

	h := &MonitorHandler{
		monitorService: monitorService,
		log:            log,
		event:          event,
		r:              &utils.Response{},
		sseHub:         hub,
	}

	// Subscribe to monitor status events to broadcast to all clients
	event.SubscribeFunc("monitor.created", h.broadcastEvent)
	event.SubscribeFunc("monitor.updated", h.broadcastEvent)
	event.SubscribeFunc("monitor.deleted", h.broadcastEvent)
	event.SubscribeFunc("monitor.checked", h.broadcastEvent)
	event.SubscribeFunc("incident.created", h.broadcastEvent)
	event.SubscribeFunc("incident.resolved", h.broadcastEvent)

	return h
}

func (h *MonitorHandler) getAuthUser(c echo.Context) (uint, string, error) {
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

// @Summary      List all monitors
// @Description  Retrieves all monitors. Admins see all, regular users see only their own.
// @Tags         Monitors
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Success      200  {object}  utils.SuccessResponseModel{data=[]response.MonitorResponse}
// @Failure      401  {object}  utils.ErrorResponseModel
// @Failure      500  {object}  utils.ErrorResponseModel
// @Router       /api/monitors [get]
func (h *MonitorHandler) GetAllMonitors(c echo.Context) error {
	ctx := c.Request().Context()
	userID, role, err := h.getAuthUser(c)
	if err != nil {
		return h.r.UnauthorizedResponse(c, err.Error())
	}

	var monitors []*entity.Monitor
	if role == "admin" {
		monitors, err = h.monitorService.GetAllMonitors(ctx)
	} else {
		monitors, err = h.monitorService.GetMonitorsByUserID(ctx, userID)
	}

	if err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	return h.r.SuccessResponse(c, response.FromEntities(monitors), "Monitors retrieved successfully")
}

// @Summary      Get monitor by ID
// @Description  Retrieves a single monitor by its ID. Users can only access their own monitors; admins can access any.
// @Tags         Monitors
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Monitor ID"
// @Success      200  {object}  utils.SuccessResponseModel{data=response.MonitorResponse}
// @Failure      400  {object}  utils.ErrorResponseModel
// @Failure      401  {object}  utils.ErrorResponseModel
// @Failure      403  {object}  utils.ErrorResponseModel
// @Failure      404  {object}  utils.ErrorResponseModel
// @Router       /api/monitors/{id} [get]
func (h *MonitorHandler) GetMonitor(c echo.Context) error {
	ctx := c.Request().Context()
	userID, role, err := h.getAuthUser(c)
	if err != nil {
		return h.r.UnauthorizedResponse(c, err.Error())
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "Invalid monitor ID")
	}

	monitor, err := h.monitorService.GetMonitorByID(ctx, uint(id))
	if err != nil {
		return h.r.NotFoundResponse(c, "Monitor not found")
	}

	// Access control: only owner or admin can view
	if role != "admin" && monitor.UserID != userID {
		return h.r.ForbiddenResponse(c, "You do not have access to this monitor")
	}

	return h.r.SuccessResponse(c, response.FromEntity(monitor), "Monitor retrieved successfully")
}

// @Summary      Create a monitor
// @Description  Creates a new uptime monitor. Non-admin users get enforced defaults (60s interval, 10s timeout). Admin users can set custom values.
// @Tags         Monitors
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body      request.CreateMonitorRequest  true  "Monitor details"
// @Success      201   {object}  utils.SuccessResponseModel{data=response.MonitorResponse}
// @Failure      400   {object}  utils.ErrorResponseModel
// @Failure      401   {object}  utils.ErrorResponseModel
// @Failure      500   {object}  utils.ErrorResponseModel
// @Router       /api/monitors [post]
func (h *MonitorHandler) CreateMonitor(c echo.Context) error {
	ctx := c.Request().Context()
	userID, role, err := h.getAuthUser(c)
	if err != nil {
		return h.r.UnauthorizedResponse(c, err.Error())
	}

	req := new(request.CreateMonitorRequest)
	if err := c.Bind(req); err != nil {
		return h.r.BadRequestResponse(c, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return h.r.BadRequestResponse(c, err.Error())
	}

	// Enforce defaults for non-admins to prevent resource exhaustion
	interval := req.Interval
	timeout := req.Timeout
	if role != "admin" {
		interval = 60
		timeout = 10
	}

	monitor := entity.NewMonitor(req.Name, req.URL, req.Type, interval, timeout, userID)
	monitor.CheckSSL = req.CheckSSL
	err = h.monitorService.CreateMonitor(ctx, monitor)
	if err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	// Attach tags if provided
	if len(req.TagIDs) > 0 {
		h.attachTags(ctx, monitor.ID, req.TagIDs)
		// Reload to include tags
		monitor, _ = h.monitorService.GetMonitorByID(ctx, monitor.ID)
	}

	h.event.Publish(bus.Event{Type: "monitor.created", Payload: monitor})

	// Run initial check immediately in the background
	checkCtx, checkCancel := context.WithTimeout(context.Background(), time.Duration(monitor.Timeout+10)*time.Second)
	go func() {
		defer checkCancel()
		h.monitorService.PerformCheck(checkCtx, monitor)
	}()

	return h.r.CreatedResponse(c, response.FromEntity(monitor), "Monitor created successfully")
}

// @Summary      Update a monitor
// @Description  Updates an existing monitor by ID. Non-admin users get enforced default interval/timeout. Only the owner or an admin can update.
// @Tags         Monitors
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id    path      int                        true  "Monitor ID"
// @Param        body  body      request.UpdateMonitorRequest  true  "Updated monitor details"
// @Success      200   {object}  utils.SuccessResponseModel{data=response.MonitorResponse}
// @Failure      400   {object}  utils.ErrorResponseModel
// @Failure      401   {object}  utils.ErrorResponseModel
// @Failure      403   {object}  utils.ErrorResponseModel
// @Failure      404   {object}  utils.ErrorResponseModel
// @Failure      500   {object}  utils.ErrorResponseModel
// @Router       /api/monitors/{id} [put]
func (h *MonitorHandler) UpdateMonitor(c echo.Context) error {
	ctx := c.Request().Context()
	userID, role, err := h.getAuthUser(c)
	if err != nil {
		return h.r.UnauthorizedResponse(c, err.Error())
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "Invalid monitor ID")
	}

	monitor, err := h.monitorService.GetMonitorByID(ctx, uint(id))
	if err != nil {
		return h.r.NotFoundResponse(c, "Monitor not found")
	}

	// Access control
	if role != "admin" && monitor.UserID != userID {
		return h.r.ForbiddenResponse(c, "You do not have permission to update this monitor")
	}

	req := new(request.UpdateMonitorRequest)
	if err := c.Bind(req); err != nil {
		return h.r.BadRequestResponse(c, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return h.r.BadRequestResponse(c, err.Error())
	}

	monitor.Name = req.Name
	monitor.URL = req.URL
	monitor.Type = req.Type
	monitor.CheckSSL = req.CheckSSL

	// Only admins can modify interval and timeout
	if role == "admin" {
		monitor.Interval = req.Interval
		monitor.Timeout = req.Timeout
	} else {
		// Enforce defaults for non-admins if they somehow bypass or had custom values
		monitor.Interval = 60
		monitor.Timeout = 10
	}
	monitor.Status = req.Status

	err = h.monitorService.UpdateMonitor(ctx, monitor)
	if err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	// Attach tags if provided (replaces existing)
	if req.TagIDs != nil {
		h.detachAllTags(ctx, monitor.ID)
		if len(req.TagIDs) > 0 {
			h.attachTags(ctx, monitor.ID, req.TagIDs)
		}
		// Reload to include tags
		monitor, _ = h.monitorService.GetMonitorByID(ctx, monitor.ID)
	}

	h.event.Publish(bus.Event{Type: "monitor.updated", Payload: monitor})

	// Run check immediately in the background
	checkCtx, checkCancel := context.WithTimeout(context.Background(), time.Duration(monitor.Timeout+10)*time.Second)
	go func() {
		defer checkCancel()
		h.monitorService.PerformCheck(checkCtx, monitor)
	}()

	return h.r.SuccessResponse(c, response.FromEntity(monitor), "Monitor updated successfully")
}

// @Summary      Delete a monitor
// @Description  Deletes a monitor by ID. Only the owner or an admin can delete.
// @Tags         Monitors
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Monitor ID"
// @Success      204  {object}  utils.SuccessResponseModel
// @Failure      400  {object}  utils.ErrorResponseModel
// @Failure      401  {object}  utils.ErrorResponseModel
// @Failure      403  {object}  utils.ErrorResponseModel
// @Failure      404  {object}  utils.ErrorResponseModel
// @Failure      500  {object}  utils.ErrorResponseModel
// @Router       /api/monitors/{id} [delete]
func (h *MonitorHandler) DeleteMonitor(c echo.Context) error {
	ctx := c.Request().Context()
	userID, role, err := h.getAuthUser(c)
	if err != nil {
		return h.r.UnauthorizedResponse(c, err.Error())
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "Invalid monitor ID")
	}

	monitor, err := h.monitorService.GetMonitorByID(ctx, uint(id))
	if err != nil {
		return h.r.NotFoundResponse(c, "Monitor not found")
	}

	// Access control
	if role != "admin" && monitor.UserID != userID {
		return h.r.ForbiddenResponse(c, "You do not have permission to delete this monitor")
	}

	err = h.monitorService.DeleteMonitor(ctx, uint(id))
	if err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	h.event.Publish(bus.Event{Type: "monitor.deleted", Payload: monitor})

	return h.r.NoContentResponse(c)
}

func (h *MonitorHandler) attachTags(ctx context.Context, monitorID uint, tagIDs []uint) {
	for _, tid := range tagIDs {
		database.DB.WithContext(ctx).Exec(
			"INSERT OR IGNORE INTO monitor_tags (monitor_id, tag_id, created_at) VALUES (?, ?, datetime('now'))",
			monitorID, tid,
		)
	}
}

func (h *MonitorHandler) detachAllTags(ctx context.Context, monitorID uint) {
	database.DB.WithContext(ctx).Exec("DELETE FROM monitor_tags WHERE monitor_id = ?", monitorID)
}

func (h *MonitorHandler) broadcastEvent(event bus.Event) {
	payload, err := json.Marshal(map[string]interface{}{
		"type":    event.Type,
		"payload": event.Payload,
	})
	if err == nil {
		h.sseHub.broadcast <- string(payload)
	}
}

// @Summary      Stream monitor events (SSE)
// @Description  Server-sent events stream for real-time monitor and incident updates. No authentication required.
// @Tags         Monitors
// @Accept       json
// @Produce      text/event-stream
// @Success      200  {object}  utils.SuccessResponseModel
// @Router       /api/monitors/events [get]
func (h *MonitorHandler) StreamEvents(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	c.Response().Header().Set(echo.HeaderCacheControl, "no-cache")
	c.Response().Header().Set(echo.HeaderConnection, "keep-alive")
	c.Response().WriteHeader(http.StatusOK)

	clientChan := make(ClientChan, 10)
	h.sseHub.register <- clientChan

	defer func() {
		h.sseHub.unregister <- clientChan
	}()

	ctx := c.Request().Context()

	// Send initial ping to establish connection
	fmt.Fprint(c.Response().Writer, "event: connected\ndata: {}\n\n")
	c.Response().Flush()

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-clientChan:
			if !ok {
				return nil
			}
			fmt.Fprintf(c.Response().Writer, "data: %s\n\n", msg)
			c.Response().Flush()
		}
	}
}

// @Summary      Receive heartbeat ping
// @Description  Public endpoint to register a heartbeat ping from a monitor via its heartbeat token. Updates uptime status to "up" and resolves any active incidents.
// @Tags         Monitors
// @Accept       json
// @Produce      json
// @Param        token  path      string  true  "Heartbeat token"
// @Success      200    {object}  utils.SuccessResponseModel{data=response.MonitorResponse}
// @Failure      400    {object}  utils.ErrorResponseModel
// @Failure      404    {object}  utils.ErrorResponseModel
// @Router       /api/heartbeat/{token} [post]
func (h *MonitorHandler) HandleHeartbeat(c echo.Context) error {
	token := c.Param("token")
	if token == "" {
		return h.r.BadRequestResponse(c, "Missing heartbeat token")
	}

	ctx := c.Request().Context()
	var mon entity.Monitor
	err := database.DB.WithContext(ctx).Where("heartbeat_token = ?", token).First(&mon).Error
	if err != nil {
		return h.r.NotFoundResponse(c, "Invalid heartbeat token")
	}

	if mon.Status != "active" {
		return h.r.BadRequestResponse(c, "Monitor is not active")
	}

	now := time.Now()
	mon.UptimeStatus = "up"
	mon.LastCheckedAt = &now
	mon.LastLatency = 0

	database.DB.WithContext(ctx).Save(&mon)
	database.DB.WithContext(ctx).Create(entity.NewCheckRecord(mon.ID, true, 0, 0))

	h.event.Publish(bus.Event{Type: "monitor.checked", Payload: &mon})

	// If monitor was down, resolve incidents
	var activeIncidents []*incidentEntity.Incident
	database.DB.WithContext(ctx).Where("monitor_id = ? AND status = ?", mon.ID, "active").Find(&activeIncidents)
	for _, inc := range activeIncidents {
		inc.Status = "resolved"
		inc.ResolvedAt = &now
		database.DB.WithContext(ctx).Save(inc)
		h.event.Publish(bus.Event{Type: "incident.resolved", Payload: inc})
	}

	return h.r.SuccessResponse(c, response.FromEntity(&mon), "Heartbeat received")
}

// @Summary      Get public daily chart data
// @Description  Returns daily uptime chart data for a public monitor. No authentication required.
// @Tags         Monitors
// @Accept       json
// @Produce      json
// @Param        id    path      int     true  "Monitor ID"
// @Param        days  query     int     false  "Number of days (1-365, default 30)"
// @Success      200   {object}  utils.SuccessResponseModel{data=[]entity.DailyChartPoint}
// @Failure      400   {object}  utils.ErrorResponseModel
// @Failure      500   {object}  utils.ErrorResponseModel
// @Router       /api/monitors/public/{id}/daily [get]
func (h *MonitorHandler) GetPublicDailyChart(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "Invalid monitor ID")
	}

	days := 30
	if dStr := c.QueryParam("days"); dStr != "" {
		if d, err := strconv.Atoi(dStr); err == nil && d > 0 && d <= 365 {
			days = d
		}
	}

	points, err := h.monitorService.GetDailyChart(ctx, uint(id), days)
	if err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	return h.r.SuccessResponse(c, points, "Daily chart data retrieved")
}

func (h *MonitorHandler) RegisterRoutes(e *echo.Echo, basePath string) {
	e.GET(basePath+"/monitors/events", h.StreamEvents)

	// Public endpoints (no auth)
	e.POST(basePath+"/heartbeat/:token", h.HandleHeartbeat)
	e.GET(basePath+"/monitors/public/:id/daily", h.GetPublicDailyChart)

	group := e.Group(basePath+"/monitors", middleware.Auth)
	group.GET("", h.GetAllMonitors)
	group.GET("/:id", h.GetMonitor)
	group.POST("", h.CreateMonitor)
	group.PUT("/:id", h.UpdateMonitor)
	group.DELETE("/:id", h.DeleteMonitor)
}

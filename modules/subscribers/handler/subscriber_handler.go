package handler

import (
	"ping-uptime/internal/pkg/logger"
	"ping-uptime/internal/pkg/middleware"
	"ping-uptime/internal/pkg/utils"
	"ping-uptime/modules/subscribers/domain/service"
	"strconv"

	"github.com/labstack/echo/v4"
)

type SubscriberHandler struct {
	svc *service.SubscriberService
	log *logger.Logger
	r   *utils.Response
}

func NewSubscriberHandler(log *logger.Logger, svc *service.SubscriberService) *SubscriberHandler {
	return &SubscriberHandler{svc: svc, log: log, r: &utils.Response{}}
}

type subscribeRequest struct {
	Email        string `json:"email" validate:"required,email"`
	StatusPageID uint   `json:"status_page_id" validate:"required"`
}

func (h *SubscriberHandler) Subscribe(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(subscribeRequest)
	if err := c.Bind(req); err != nil {
		return h.r.BadRequestResponse(c, err.Error())
	}
	if err := c.Validate(req); err != nil {
		return h.r.BadRequestResponse(c, err.Error())
	}

	sub, err := h.svc.Subscribe(ctx, req.Email, req.StatusPageID)
	if err != nil {
		if err == service.ErrAlreadySubscribed {
			return h.r.ConflictResponse(c, "Already subscribed")
		}
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	return h.r.SuccessResponse(c, sub, "Subscription successful. Please check your email to verify.")
}

type verifyQuery struct {
	Token string `json:"token" query:"token"`
}

func (h *SubscriberHandler) Verify(c echo.Context) error {
	ctx := c.Request().Context()
	token := c.QueryParam("token")
	if token == "" {
		return h.r.BadRequestResponse(c, "Missing verification token")
	}

	if err := h.svc.Verify(ctx, token); err != nil {
		return h.r.BadRequestResponse(c, "Invalid or expired token")
	}

	return h.r.SuccessResponse(c, nil, "Email verified successfully")
}

type unsubscribeRequest struct {
	Email        string `json:"email" validate:"required,email"`
	StatusPageID uint   `json:"status_page_id" validate:"required"`
}

func (h *SubscriberHandler) Unsubscribe(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(unsubscribeRequest)
	if err := c.Bind(req); err != nil {
		return h.r.BadRequestResponse(c, err.Error())
	}

	if err := h.svc.Unsubscribe(ctx, req.Email, req.StatusPageID); err != nil {
		return h.r.NotFoundResponse(c, "Subscription not found")
	}

	return h.r.SuccessResponse(c, nil, "Unsubscribed successfully")
}

func (h *SubscriberHandler) Count(c echo.Context) error {
	ctx := c.Request().Context()
	pageIDStr := c.Param("pageID")
	pageID, err := strconv.ParseUint(pageIDStr, 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "Invalid page ID")
	}

	count, err := h.svc.CountByPageID(ctx, uint(pageID))
	if err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	return h.r.SuccessResponse(c, map[string]int64{"count": count}, "Subscriber count retrieved")
}

func (h *SubscriberHandler) RegisterRoutes(e *echo.Echo, basePath string) {
	// Public routes (no auth)
	e.POST(basePath+"/status-pages/subscribe", h.Subscribe)
	e.GET(basePath+"/status-pages/subscribe/verify", h.Verify)
	e.POST(basePath+"/status-pages/unsubscribe", h.Unsubscribe)

	// Authenticated routes
	group := e.Group(basePath+"/status-pages", middleware.Auth)
	group.GET("/:pageID/subscribers/count", h.Count)
}

package handler

import (
	"ping-uptime/internal/pkg/logger"
	"ping-uptime/internal/pkg/middleware"
	"ping-uptime/internal/pkg/utils"
	"ping-uptime/modules/notifications/domain/repository"

	"github.com/labstack/echo/v4"
)

type NotificationLogHandler struct {
	logRepo repository.NotificationLogRepository
	log     *logger.Logger
	r       *utils.Response
}

func NewNotificationLogHandler(log *logger.Logger, logRepo repository.NotificationLogRepository) *NotificationLogHandler {
	return &NotificationLogHandler{logRepo: logRepo, log: log, r: &utils.Response{}}
}

func (h *NotificationLogHandler) getAuthUser(c echo.Context) (uint, string, error) {
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

// @Summary      List notification logs
// @Description  Get paginated notification logs for the authenticated user (admin sees all)
// @Tags         Notification Logs
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        channel_type  query  string  false  "Filter by channel type"
// @Param        status        query  string  false  "Filter by status (sent/failed)"
// @Param        event_type    query  string  false  "Filter by event type"
// @Param        date_from     query  string  false  "Filter logs from date (RFC3339)"
// @Param        date_to       query  string  false  "Filter logs to date (RFC3339)"
// @Param        limit         query  int     false  "Results per page"
// @Param        offset        query  int     false  "Result offset"
// @Success      200  {object}  utils.SuccessResponseModel{data=[]entity.NotificationLog}
// @Failure      401  {object}  utils.ErrorResponseModel
// @Failure      500  {object}  utils.ErrorResponseModel
// @Router       /notification-logs [get]
func (h *NotificationLogHandler) GetAll(c echo.Context) error {
	ctx := c.Request().Context()
	userID, role, err := h.getAuthUser(c)
	if err != nil {
		return h.r.UnauthorizedResponse(c, err.Error())
	}

	isAdmin := role == "admin"

	filter := map[string]string{}
	if channelType := c.QueryParam("channel_type"); channelType != "" {
		filter["channel_type"] = channelType
	}
	if status := c.QueryParam("status"); status != "" {
		filter["status"] = status
	}
	if eventType := c.QueryParam("event_type"); eventType != "" {
		filter["event_type"] = eventType
	}
	if dateFrom := c.QueryParam("date_from"); dateFrom != "" {
		filter["date_from"] = dateFrom
	}
	if dateTo := c.QueryParam("date_to"); dateTo != "" {
		filter["date_to"] = dateTo
	}
	if limit := c.QueryParam("limit"); limit != "" {
		filter["limit"] = limit
	}
	if offset := c.QueryParam("offset"); offset != "" {
		filter["offset"] = offset
	}

	logs, err := h.logRepo.FindAll(ctx, userID, isAdmin, filter)
	if err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	return h.r.SuccessResponse(c, logs, "Notification logs retrieved")
}

func (h *NotificationLogHandler) RegisterRoutes(e *echo.Echo, basePath string) {
	group := e.Group(basePath+"/notification-logs", middleware.Auth)
	group.GET("", h.GetAll)
}

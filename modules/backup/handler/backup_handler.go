package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"ping-uptime/internal/pkg/middleware"
	"ping-uptime/internal/pkg/utils"
	"ping-uptime/modules/backup/domain/service"

	"github.com/labstack/echo/v4"
)

type BackupHandler struct {
	svc *service.BackupService
	r   *utils.Response
}

func NewBackupHandler(svc *service.BackupService) *BackupHandler {
	return &BackupHandler{svc: svc, r: &utils.Response{}}
}

func (h *BackupHandler) getAuthUser(c echo.Context) (uint, error) {
	userClaims, ok := c.Get("user").(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("unauthorized")
	}
	userIDVal, ok := userClaims["user_id"]
	if !ok {
		return 0, fmt.Errorf("invalid token: user_id missing")
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
		return 0, fmt.Errorf("invalid user_id type")
	}
	return userID, nil
}

// Export godoc
// GET /api/backup/export
func (h *BackupHandler) Export(c echo.Context) error {
	ctx := c.Request().Context()

	data, err := h.svc.Export(ctx)
	if err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	// Record the export
	userID, err := h.getAuthUser(c)
	if err != nil {
		userID = 0
	}

	// Store a history entry for the export (no file on disk, just a record)
	if err := h.svc.CreateRecord(ctx, "manual-export-"+time.Now().Format("2006-01-02")+".json", 0, userID); err != nil {
		// Non-fatal — log but don't block the download
		c.Logger().Error("failed to record backup history: " + err.Error())
	}

	return h.r.SuccessResponse(c, data, "Backup exported successfully")
}

// Import godoc
// POST /api/backup/import
func (h *BackupHandler) Import(c echo.Context) error {
	ctx := c.Request().Context()

	// Parse multipart form, max 50MB
	if err := c.Request().ParseMultipartForm(50 << 20); err != nil {
		return h.r.BadRequestResponse(c, "Failed to parse upload: "+err.Error())
	}

	file, header, err := c.Request().FormFile("file")
	if err != nil {
		return h.r.BadRequestResponse(c, "Missing file field: "+err.Error())
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return h.r.BadRequestResponse(c, "Failed to read file: "+err.Error())
	}

	if err := h.svc.Import(ctx, data); err != nil {
		return h.r.BadRequestResponse(c, "Import failed: "+err.Error())
	}

	// Record import
	userID, err := h.getAuthUser(c)
	if err != nil {
		userID = 0
	}
	if err := h.svc.CreateRecord(ctx, "import-"+header.Filename, header.Size, userID); err != nil {
		c.Logger().Error("failed to record import history: " + err.Error())
	}

	return h.r.SuccessResponse(c, nil, "Configuration imported successfully")
}

// Download godoc
// GET /api/backup/download
func (h *BackupHandler) Download(c echo.Context) error {
	ctx := c.Request().Context()

	data, err := h.svc.Export(ctx)
	if err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	filename := "ping-uptime-backup-" + time.Now().Format("2006-01-02-150405") + ".json"
	c.Response().Header().Set(echo.HeaderContentDisposition, "attachment; filename="+filename)
	c.Response().Header().Set(echo.HeaderContentType, "application/json")
	return c.Blob(http.StatusOK, "application/json", jsonData)
}

// GetHistory godoc
// GET /api/backup/history
func (h *BackupHandler) GetHistory(c echo.Context) error {
	ctx := c.Request().Context()

	records, err := h.svc.GetHistory(ctx)
	if err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	return h.r.SuccessResponse(c, records, "Backup history retrieved")
}

// DeleteRecord godoc
// DELETE /api/backup/history/:id
func (h *BackupHandler) DeleteRecord(c echo.Context) error {
	ctx := c.Request().Context()

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return h.r.BadRequestResponse(c, "Invalid ID")
	}

	if err := h.svc.DeleteRecord(ctx, uint(id)); err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	return h.r.SuccessResponse(c, nil, "Backup record deleted")
}

func (h *BackupHandler) RegisterRoutes(e *echo.Echo, basePath string) {
	group := e.Group(basePath+"/backup", middleware.Auth, middleware.Admin)
	group.GET("/export", h.Export)
	group.GET("/download", h.Download)
	group.POST("/import", h.Import)
	group.GET("/history", h.GetHistory)
	group.DELETE("/history/:id", h.DeleteRecord)
}

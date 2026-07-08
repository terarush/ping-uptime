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

// @Summary      [Admin] Export backup data
// @Description  Export all configuration data as a JSON backup. Requires admin privileges.
// @Tags         Backup
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Success      200  {object}  utils.SuccessResponseModel
// @Failure      401  {object}  utils.ErrorResponseModel
// @Failure      403  {object}  utils.ErrorResponseModel
// @Failure      500  {object}  utils.ErrorResponseModel
// @Router       /api/backup/export [get]
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

// @Summary      [Admin] Import backup data
// @Description  Import configuration from a JSON backup file upload (multipart). Requires admin privileges.
// @Tags         Backup
// @Security     BearerAuth
// @Accept       multipart/form-data
// @Produce      json
// @Param        file  formData  file  true  "Backup JSON file to import"
// @Success      200   {object}  utils.SuccessResponseModel
// @Failure      400   {object}  utils.ErrorResponseModel
// @Failure      401   {object}  utils.ErrorResponseModel
// @Failure      403   {object}  utils.ErrorResponseModel
// @Failure      500   {object}  utils.ErrorResponseModel
// @Router       /api/backup/import [post]
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

// @Summary      [Admin] Download backup as file
// @Description  Download a JSON backup file as an attachment. Requires admin privileges.
// @Tags         Backup
// @Security     BearerAuth
// @Accept       json
// @Produce      application/json
// @Success      200  {file}  binary  "Backup JSON file download"
// @Failure      401  {object}  utils.ErrorResponseModel
// @Failure      403  {object}  utils.ErrorResponseModel
// @Failure      500  {object}  utils.ErrorResponseModel
// @Router       /api/backup/download [get]
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

// @Summary      [Admin] List backup history
// @Description  Retrieve backup/import history records. Requires admin privileges.
// @Tags         Backup
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Success      200  {object}  utils.SuccessResponseModel
// @Failure      401  {object}  utils.ErrorResponseModel
// @Failure      403  {object}  utils.ErrorResponseModel
// @Failure      500  {object}  utils.ErrorResponseModel
// @Router       /api/backup/history [get]
func (h *BackupHandler) GetHistory(c echo.Context) error {
	ctx := c.Request().Context()

	records, err := h.svc.GetHistory(ctx)
	if err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	return h.r.SuccessResponse(c, records, "Backup history retrieved")
}

// @Summary      [Admin] Delete backup history record
// @Description  Delete a backup/import history record by ID. Requires admin privileges.
// @Tags         Backup
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path  int  true  "Backup record ID"
// @Success      200  {object}  utils.SuccessResponseModel
// @Failure      400  {object}  utils.ErrorResponseModel
// @Failure      401  {object}  utils.ErrorResponseModel
// @Failure      403  {object}  utils.ErrorResponseModel
// @Failure      500  {object}  utils.ErrorResponseModel
// @Router       /api/backup/history/{id} [delete]
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

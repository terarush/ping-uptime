package handler

import (
	"ping-uptime/internal/pkg/middleware"
	"ping-uptime/internal/pkg/utils"
	"ping-uptime/modules/tags/domain/service"
	"strconv"

	"github.com/labstack/echo/v4"
)

type TagHandler struct {
	svc *service.TagService
	r   *utils.Response
}

func NewTagHandler(svc *service.TagService) *TagHandler {
	return &TagHandler{svc: svc, r: &utils.Response{}}
}

type createTagRequest struct {
	Name  string `json:"name" validate:"required"`
	Color string `json:"color"`
}

type updateTagRequest struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type attachTagsRequest struct {
	TagIDs []uint `json:"tag_ids" validate:"required"`
}

// @Summary      List all tags
// @Description  Get all tags
// @Tags         Tags
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Success      200  {object}  utils.SuccessResponseModel{data=[]entity.Tag}
// @Failure      401  {object}  utils.ErrorResponseModel
// @Router       /api/tags [get]
func (h *TagHandler) GetAll(c echo.Context) error {
	ctx := c.Request().Context()
	tags, err := h.svc.GetAll(ctx)
	if err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.SuccessResponse(c, tags, "Tags retrieved successfully")
}

// @Summary      Create a tag
// @Description  Create a new tag
// @Tags         Tags
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  createTagRequest  true  "Tag details"
// @Success      201  {object}  utils.SuccessResponseModel{data=entity.Tag}
// @Failure      400  {object}  utils.ErrorResponseModel
// @Failure      401  {object}  utils.ErrorResponseModel
// @Router       /api/tags [post]
func (h *TagHandler) Create(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(createTagRequest)
	if err := c.Bind(req); err != nil {
		return h.r.BadRequestResponse(c, err.Error())
	}
	if err := c.Validate(req); err != nil {
		return h.r.BadRequestResponse(c, err.Error())
	}

	tag, err := h.svc.Create(ctx, req.Name, req.Color)
	if err != nil {
		if err == service.ErrTagNameEmpty || err == service.ErrTagNameTaken {
			return h.r.BadRequestResponse(c, err.Error())
		}
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.CreatedResponse(c, tag, "Tag created successfully")
}

// @Summary      Update a tag
// @Description  Update an existing tag
// @Tags         Tags
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path    int                true  "Tag ID"
// @Param        request  body    updateTagRequest   true  "Updated tag details"
// @Success      200  {object}  utils.SuccessResponseModel{data=entity.Tag}
// @Failure      400  {object}  utils.ErrorResponseModel
// @Failure      401  {object}  utils.ErrorResponseModel
// @Failure      404  {object}  utils.ErrorResponseModel
// @Router       /api/tags/{id} [put]
func (h *TagHandler) Update(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "Invalid tag ID")
	}

	req := new(updateTagRequest)
	if err := c.Bind(req); err != nil {
		return h.r.BadRequestResponse(c, err.Error())
	}

	tag, err := h.svc.Update(ctx, uint(id), req.Name, req.Color)
	if err != nil {
		if err == service.ErrTagNotFound {
			return h.r.NotFoundResponse(c, "Tag not found")
		}
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.SuccessResponse(c, tag, "Tag updated successfully")
}

// @Summary      Delete a tag
// @Description  Delete a tag by ID
// @Tags         Tags
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path    int  true  "Tag ID"
// @Success      204  "No Content"
// @Failure      400  {object}  utils.ErrorResponseModel
// @Failure      401  {object}  utils.ErrorResponseModel
// @Failure      404  {object}  utils.ErrorResponseModel
// @Router       /api/tags/{id} [delete]
func (h *TagHandler) Delete(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "Invalid tag ID")
	}

	if err := h.svc.Delete(ctx, uint(id)); err != nil {
		if err == service.ErrTagNotFound {
			return h.r.NotFoundResponse(c, "Tag not found")
		}
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.NoContentResponse(c)
}

// @Summary      Attach tags to monitor
// @Description  Attach tags to a monitor
// @Tags         Tags
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path    int                  true  "Monitor ID"
// @Param        request  body    attachTagsRequest    true  "Tag IDs to attach"
// @Success      200  {object}  utils.SuccessResponseModel
// @Failure      400  {object}  utils.ErrorResponseModel
// @Failure      401  {object}  utils.ErrorResponseModel
// @Failure      404  {object}  utils.ErrorResponseModel
// @Router       /api/monitors/{id}/tags [post]
func (h *TagHandler) AttachTags(c echo.Context) error {
	ctx := c.Request().Context()
	monitorID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "Invalid monitor ID")
	}

	req := new(attachTagsRequest)
	if err := c.Bind(req); err != nil {
		return h.r.BadRequestResponse(c, err.Error())
	}
	if err := c.Validate(req); err != nil {
		return h.r.BadRequestResponse(c, err.Error())
	}

	if err := h.svc.AttachTags(ctx, uint(monitorID), req.TagIDs); err != nil {
		if err == service.ErrTagNotFound {
			return h.r.NotFoundResponse(c, "One or more tags not found")
		}
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.SuccessResponse(c, nil, "Tags attached successfully")
}

// @Summary      Get monitor tags
// @Description  Get all tags attached to a monitor
// @Tags         Tags
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path    int  true  "Monitor ID"
// @Success      200  {object}  utils.SuccessResponseModel{data=[]entity.Tag}
// @Failure      400  {object}  utils.ErrorResponseModel
// @Failure      401  {object}  utils.ErrorResponseModel
// @Router       /api/monitors/{id}/tags [get]
func (h *TagHandler) GetMonitorTags(c echo.Context) error {
	ctx := c.Request().Context()
	monitorID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "Invalid monitor ID")
	}

	tags, err := h.svc.GetMonitorTags(ctx, uint(monitorID))
	if err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.SuccessResponse(c, tags, "Monitor tags retrieved successfully")
}

// @Summary      Detach tag from monitor
// @Description  Remove a tag from a monitor
// @Tags         Tags
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id     path    int  true  "Monitor ID"
// @Param        tagID  path    int  true  "Tag ID"
// @Success      204  "No Content"
// @Failure      400  {object}  utils.ErrorResponseModel
// @Failure      401  {object}  utils.ErrorResponseModel
// @Router       /api/monitors/{id}/tags/{tagID} [delete]
func (h *TagHandler) DetachTag(c echo.Context) error {
	ctx := c.Request().Context()
	monitorID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "Invalid monitor ID")
	}

	tagID, err := strconv.ParseUint(c.Param("tagID"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "Invalid tag ID")
	}

	if err := h.svc.DetachTag(ctx, uint(monitorID), uint(tagID)); err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}
	return h.r.NoContentResponse(c)
}

func (h *TagHandler) RegisterRoutes(e *echo.Echo, basePath string) {
	group := e.Group(basePath+"/tags", middleware.Auth)
	group.GET("", h.GetAll)
	group.POST("", h.Create)
	group.PUT("/:id", h.Update)
	group.DELETE("/:id", h.Delete)

	// Monitor-tag association routes
	monitorGroup := e.Group(basePath+"/monitors", middleware.Auth)
	monitorGroup.POST("/:id/tags", h.AttachTags)
	monitorGroup.GET("/:id/tags", h.GetMonitorTags)
	monitorGroup.DELETE("/:id/tags/:tagID", h.DetachTag)
}

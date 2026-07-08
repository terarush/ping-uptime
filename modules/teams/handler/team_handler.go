package handler

import (
	"fmt"
	"ping-uptime/internal/pkg/middleware"
	"ping-uptime/internal/pkg/utils"
	"ping-uptime/modules/teams/domain/entity"
	"ping-uptime/modules/teams/domain/service"
	"ping-uptime/modules/teams/dto/request"
	"ping-uptime/modules/teams/dto/response"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type TeamHandler struct {
	svc *service.TeamService
	db  *gorm.DB
	r   *utils.Response
}

func NewTeamHandler(svc *service.TeamService, db *gorm.DB) *TeamHandler {
	return &TeamHandler{svc: svc, db: db, r: &utils.Response{}}
}

func (h *TeamHandler) getAuthUser(c echo.Context) (uint, string, error) {
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

// @Summary      List teams
// @Description  Get all teams. Admin sees all; user sees own teams.
// @Tags         Teams
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Success      200  {object}  utils.SuccessResponseModel{data=[]response.TeamResponse}
// @Failure      401  {object}  utils.ErrorResponseModel
// @Router       /api/teams [get]
func (h *TeamHandler) GetAll(c echo.Context) error {
	ctx := c.Request().Context()
	userID, role, err := h.getAuthUser(c)
	if err != nil {
		return h.r.UnauthorizedResponse(c, err.Error())
	}

	var teams []*entity.Team
	if role == "admin" {
		teams, err = h.svc.GetAll(ctx)
	} else {
		// For non-admin: get all teams and we could filter by membership
		allTeams, errAll := h.svc.GetAll(ctx)
		if errAll != nil {
			return h.r.InternalServerErrorResponse(c, errAll.Error())
		}
		// Filter to teams where user is a member
		var userTeams []*entity.Team
		for _, t := range allTeams {
			isAdmin, _ := h.svc.IsTeamAdmin(ctx, t.ID, userID)
			if isAdmin {
				userTeams = append(userTeams, t)
				continue
			}
			members, _ := h.svc.GetMembers(ctx, t.ID)
			for _, m := range members {
				if m.UserID == userID {
					userTeams = append(userTeams, t)
					break
				}
			}
		}
		teams = userTeams
	}
	if err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	resp := response.TeamFromEntities(teams)

	// Attach member count
	for i, t := range teams {
		members, err := h.svc.GetMembers(ctx, t.ID)
		if err == nil {
			count := int64(0)
			for _, m := range members {
				if m.Status == "accepted" {
					count++
				}
			}
			resp[i].MemberCount = count
		}
	}

	return h.r.SuccessResponse(c, resp, "Teams retrieved successfully")
}

// @Summary      Create a team
// @Description  Create a new team. Creator becomes team admin.
// @Tags         Teams
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  request.CreateTeamRequest  true  "Team details"
// @Success      201  {object}  utils.SuccessResponseModel{data=response.TeamResponse}
// @Failure      400  {object}  utils.ErrorResponseModel
// @Failure      401  {object}  utils.ErrorResponseModel
// @Router       /api/teams [post]
func (h *TeamHandler) Create(c echo.Context) error {
	ctx := c.Request().Context()
	userID, _, err := h.getAuthUser(c)
	if err != nil {
		return h.r.UnauthorizedResponse(c, err.Error())
	}

	req := new(request.CreateTeamRequest)
	if err := c.Bind(req); err != nil {
		return h.r.BadRequestResponse(c, err.Error())
	}
	if err := c.Validate(req); err != nil {
		return h.r.BadRequestResponse(c, err.Error())
	}

	team := &entity.Team{Name: req.Name}
	if err := h.svc.Create(ctx, team, userID); err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	return h.r.CreatedResponse(c, response.TeamFromEntity(team), "Team created successfully")
}

// @Summary      Update a team
// @Description  Update a team's name. Only team admins.
// @Tags         Teams
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path    int                        true  "Team ID"
// @Param        request  body    request.UpdateTeamRequest  true  "Updated team details"
// @Success      200  {object}  utils.SuccessResponseModel{data=response.TeamResponse}
// @Failure      400  {object}  utils.ErrorResponseModel
// @Failure      401  {object}  utils.ErrorResponseModel
// @Failure      403  {object}  utils.ErrorResponseModel
// @Failure      404  {object}  utils.ErrorResponseModel
// @Router       /api/teams/{id} [put]
func (h *TeamHandler) Update(c echo.Context) error {
	ctx := c.Request().Context()
	userID, _, err := h.getAuthUser(c)
	if err != nil {
		return h.r.UnauthorizedResponse(c, err.Error())
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "Invalid team ID")
	}

	// Check admin permission
	isAdmin, _ := h.svc.IsTeamAdmin(ctx, uint(id), userID)
	if !isAdmin {
		return h.r.ForbiddenResponse(c, "Only team admins can update the team")
	}

	req := new(request.UpdateTeamRequest)
	if err := c.Bind(req); err != nil {
		return h.r.BadRequestResponse(c, err.Error())
	}
	if err := c.Validate(req); err != nil {
		return h.r.BadRequestResponse(c, err.Error())
	}

	team, err := h.svc.Update(ctx, uint(id), req.Name)
	if err != nil {
		if err == service.ErrTeamNotFound {
			return h.r.NotFoundResponse(c, "Team not found")
		}
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	return h.r.SuccessResponse(c, response.TeamFromEntity(team), "Team updated successfully")
}

// @Summary      [Admin] Delete team
// @Description  Delete a team. Only admins.
// @Tags         Teams
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path    int  true  "Team ID"
// @Success      204  "No Content"
// @Failure      400  {object}  utils.ErrorResponseModel
// @Failure      401  {object}  utils.ErrorResponseModel
// @Failure      403  {object}  utils.ErrorResponseModel
// @Router       /api/teams/{id} [delete]
func (h *TeamHandler) Delete(c echo.Context) error {
	ctx := c.Request().Context()
	_, role, err := h.getAuthUser(c)
	if err != nil {
		return h.r.UnauthorizedResponse(c, err.Error())
	}

	if role != "admin" {
		return h.r.ForbiddenResponse(c, "Only admins can delete teams")
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "Invalid team ID")
	}

	if err := h.svc.Delete(ctx, uint(id)); err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	return h.r.NoContentResponse(c)
}

// @Summary      List team members
// @Description  Get all members of a team. Accessible by team members and admins.
// @Tags         Teams
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path    int  true  "Team ID"
// @Success      200  {object}  utils.SuccessResponseModel{data=[]response.TeamMemberResponse}
// @Failure      400  {object}  utils.ErrorResponseModel
// @Failure      401  {object}  utils.ErrorResponseModel
// @Failure      403  {object}  utils.ErrorResponseModel
// @Router       /api/teams/{id}/members [get]
func (h *TeamHandler) GetMembers(c echo.Context) error {
	ctx := c.Request().Context()
	userID, role, err := h.getAuthUser(c)
	if err != nil {
		return h.r.UnauthorizedResponse(c, err.Error())
	}

	teamID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "Invalid team ID")
	}

	// Check access: admin or team member
	if role != "admin" {
		isAdmin, _ := h.svc.IsTeamAdmin(ctx, uint(teamID), userID)
		if !isAdmin {
			members, _ := h.svc.GetMembers(ctx, uint(teamID))
			isMember := false
			for _, m := range members {
				if m.UserID == userID {
					isMember = true
					break
				}
			}
			if !isMember {
				return h.r.ForbiddenResponse(c, "You are not a member of this team")
			}
		}
	}

	members, err := h.svc.GetMembers(ctx, uint(teamID))
	if err != nil {
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	resp := response.TeamMemberFromEntities(members)

	// Enrich with user info (email, name)
	for i, m := range members {
		var user struct {
			ID    uint
			Email string
			Name  string
		}
		if err := h.db.Table("users").Select("id, email, name").Where("id = ?", m.UserID).Scan(&user).Error; err == nil {
			resp[i] = response.TeamMemberWithUserFromEntity(m, user.ID, user.Email, user.Name)
		}
	}

	return h.r.SuccessResponse(c, resp, "Team members retrieved successfully")
}

// @Summary      Invite team member
// @Description  Invite a user to join a team. Only team admins.
// @Tags         Teams
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path    int                         true  "Team ID"
// @Param        request  body    request.InviteMemberRequest  true  "Invitation details"
// @Success      201  {object}  utils.SuccessResponseModel{data=response.TeamMemberResponse}
// @Failure      400  {object}  utils.ErrorResponseModel
// @Failure      401  {object}  utils.ErrorResponseModel
// @Failure      403  {object}  utils.ErrorResponseModel
// @Failure      409  {object}  utils.ErrorResponseModel
// @Router       /api/teams/{id}/members [post]
func (h *TeamHandler) InviteMember(c echo.Context) error {
	ctx := c.Request().Context()
	userID, _, err := h.getAuthUser(c)
	if err != nil {
		return h.r.UnauthorizedResponse(c, err.Error())
	}

	teamID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "Invalid team ID")
	}

	// Only team admin can invite
	isAdmin, _ := h.svc.IsTeamAdmin(ctx, uint(teamID), userID)
	if !isAdmin {
		return h.r.ForbiddenResponse(c, "Only team admins can invite members")
	}

	req := new(request.InviteMemberRequest)
	if err := c.Bind(req); err != nil {
		return h.r.BadRequestResponse(c, err.Error())
	}
	if err := c.Validate(req); err != nil {
		return h.r.BadRequestResponse(c, err.Error())
	}

	member, err := h.svc.InviteMember(ctx, uint(teamID), userID, req.UserID, req.Role)
	if err != nil {
		if err == service.ErrAlreadyMember {
			return h.r.ConflictResponse(c, "User is already a member of this team")
		}
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	return h.r.CreatedResponse(c, response.TeamMemberFromEntity(member), "Invitation sent successfully")
}

// @Summary      Update team member role
// @Description  Update a team member's role. Only team admins.
// @Tags         Teams
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id        path    int                         true  "Team ID"
// @Param        memberID  path    int                         true  "Member ID"
// @Param        request   body    request.UpdateMemberRequest  true  "Updated role"
// @Success      200  {object}  utils.SuccessResponseModel{data=response.TeamMemberResponse}
// @Failure      400  {object}  utils.ErrorResponseModel
// @Failure      401  {object}  utils.ErrorResponseModel
// @Failure      403  {object}  utils.ErrorResponseModel
// @Failure      404  {object}  utils.ErrorResponseModel
// @Router       /api/teams/{id}/members/{memberID} [put]
func (h *TeamHandler) UpdateMember(c echo.Context) error {
	ctx := c.Request().Context()
	userID, _, err := h.getAuthUser(c)
	if err != nil {
		return h.r.UnauthorizedResponse(c, err.Error())
	}

	teamID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "Invalid team ID")
	}

	memberID, err := strconv.ParseUint(c.Param("memberID"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "Invalid member ID")
	}

	// Only team admin can update roles
	isAdmin, _ := h.svc.IsTeamAdmin(ctx, uint(teamID), userID)
	if !isAdmin {
		return h.r.ForbiddenResponse(c, "Only team admins can update member roles")
	}

	req := new(request.UpdateMemberRequest)
	if err := c.Bind(req); err != nil {
		return h.r.BadRequestResponse(c, err.Error())
	}
	if err := c.Validate(req); err != nil {
		return h.r.BadRequestResponse(c, err.Error())
	}

	member, err := h.svc.UpdateMemberRole(ctx, uint(teamID), uint(memberID), req.Role)
	if err != nil {
		if err == service.ErrMemberNotFound {
			return h.r.NotFoundResponse(c, "Member not found")
		}
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	return h.r.SuccessResponse(c, response.TeamMemberFromEntity(member), "Member role updated successfully")
}

// @Summary      Remove team member
// @Description  Remove a member from a team. Only team admins.
// @Tags         Teams
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id        path    int  true  "Team ID"
// @Param        memberID  path    int  true  "Member ID"
// @Success      204  "No Content"
// @Failure      400  {object}  utils.ErrorResponseModel
// @Failure      401  {object}  utils.ErrorResponseModel
// @Failure      403  {object}  utils.ErrorResponseModel
// @Failure      404  {object}  utils.ErrorResponseModel
// @Router       /api/teams/{id}/members/{memberID} [delete]
func (h *TeamHandler) RemoveMember(c echo.Context) error {
	ctx := c.Request().Context()
	userID, _, err := h.getAuthUser(c)
	if err != nil {
		return h.r.UnauthorizedResponse(c, err.Error())
	}

	teamID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "Invalid team ID")
	}

	memberID, err := strconv.ParseUint(c.Param("memberID"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "Invalid member ID")
	}

	// Only team admin can remove members
	isAdmin, _ := h.svc.IsTeamAdmin(ctx, uint(teamID), userID)
	if !isAdmin {
		return h.r.ForbiddenResponse(c, "Only team admins can remove members")
	}

	if err := h.svc.RemoveMember(ctx, uint(teamID), uint(memberID)); err != nil {
		if err == service.ErrMemberNotFound {
			return h.r.NotFoundResponse(c, "Member not found")
		}
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	return h.r.NoContentResponse(c)
}

// @Summary      Accept team invitation
// @Description  Accept a pending team invitation
// @Tags         Teams
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path    int  true  "Team ID"
// @Success      200  {object}  utils.SuccessResponseModel
// @Failure      400  {object}  utils.ErrorResponseModel
// @Failure      401  {object}  utils.ErrorResponseModel
// @Failure      404  {object}  utils.ErrorResponseModel
// @Router       /api/teams/{id}/members/accept [post]
func (h *TeamHandler) AcceptInvitation(c echo.Context) error {
	ctx := c.Request().Context()
	userID, _, err := h.getAuthUser(c)
	if err != nil {
		return h.r.UnauthorizedResponse(c, err.Error())
	}

	teamID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "Invalid team ID")
	}

	if err := h.svc.AcceptInvitation(ctx, uint(teamID), userID); err != nil {
		if err == service.ErrMemberNotFound {
			return h.r.NotFoundResponse(c, "Invitation not found")
		}
		if err == service.ErrNotPending {
			return h.r.BadRequestResponse(c, "Invitation is not pending")
		}
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	return h.r.SuccessResponse(c, nil, "Invitation accepted successfully")
}

// @Summary      Reject team invitation
// @Description  Reject a pending team invitation
// @Tags         Teams
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path    int  true  "Team ID"
// @Success      200  {object}  utils.SuccessResponseModel
// @Failure      400  {object}  utils.ErrorResponseModel
// @Failure      401  {object}  utils.ErrorResponseModel
// @Failure      404  {object}  utils.ErrorResponseModel
// @Router       /api/teams/{id}/members/reject [post]
func (h *TeamHandler) RejectInvitation(c echo.Context) error {
	ctx := c.Request().Context()
	userID, _, err := h.getAuthUser(c)
	if err != nil {
		return h.r.UnauthorizedResponse(c, err.Error())
	}

	teamID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return h.r.BadRequestResponse(c, "Invalid team ID")
	}

	if err := h.svc.RejectInvitation(ctx, uint(teamID), userID); err != nil {
		if err == service.ErrMemberNotFound {
			return h.r.NotFoundResponse(c, "Invitation not found")
		}
		if err == service.ErrNotPending {
			return h.r.BadRequestResponse(c, "Invitation is not pending")
		}
		return h.r.InternalServerErrorResponse(c, err.Error())
	}

	return h.r.SuccessResponse(c, nil, "Invitation rejected successfully")
}

func (h *TeamHandler) RegisterRoutes(e *echo.Echo, basePath string) {
	group := e.Group(basePath+"/teams", middleware.Auth)
	group.GET("", h.GetAll)
	group.POST("", h.Create)
	group.PUT("/:id", h.Update)
	group.DELETE("/:id", h.Delete)
	group.GET("/:id/members", h.GetMembers)
	group.POST("/:id/members", h.InviteMember)
	group.PUT("/:id/members/:memberID", h.UpdateMember)
	group.DELETE("/:id/members/:memberID", h.RemoveMember)
	group.POST("/:id/members/accept", h.AcceptInvitation)
	group.POST("/:id/members/reject", h.RejectInvitation)
}

package response

import (
	"ping-uptime/modules/teams/domain/entity"
	"time"
)

type TeamResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	MemberCount int64  `json:"member_count,omitempty"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type TeamMemberResponse struct {
	ID        uint              `json:"id"`
	TeamID    uint              `json:"team_id"`
	UserID    uint              `json:"user_id"`
	Role      string            `json:"role"`
	Status    string            `json:"status"`
	InvitedBy uint              `json:"invited_by"`
	User      *TeamMemberUser   `json:"user,omitempty"`
	CreatedAt string            `json:"created_at"`
	UpdatedAt string            `json:"updated_at"`
}

type TeamMemberUser struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func TeamFromEntity(t *entity.Team) *TeamResponse {
	return &TeamResponse{
		ID:        t.ID,
		Name:      t.Name,
		CreatedAt: t.CreatedAt.Format(time.RFC3339),
		UpdatedAt: t.UpdatedAt.Format(time.RFC3339),
	}
}

func TeamFromEntities(teams []*entity.Team) []*TeamResponse {
	resp := make([]*TeamResponse, len(teams))
	for i, t := range teams {
		resp[i] = TeamFromEntity(t)
	}
	return resp
}

func TeamMemberFromEntity(m *entity.TeamMember) *TeamMemberResponse {
	return &TeamMemberResponse{
		ID:        m.ID,
		TeamID:    m.TeamID,
		UserID:    m.UserID,
		Role:      m.Role,
		Status:    m.Status,
		InvitedBy: m.InvitedBy,
		CreatedAt: m.CreatedAt.Format(time.RFC3339),
		UpdatedAt: m.UpdatedAt.Format(time.RFC3339),
	}
}

func TeamMemberFromEntities(members []*entity.TeamMember) []*TeamMemberResponse {
	resp := make([]*TeamMemberResponse, len(members))
	for i, m := range members {
		resp[i] = TeamMemberFromEntity(m)
	}
	return resp
}

func TeamMemberWithUserFromEntity(m *entity.TeamMember, userID uint, email, name string) *TeamMemberResponse {
	r := TeamMemberFromEntity(m)
	r.User = &TeamMemberUser{
		ID:    userID,
		Email: email,
		Name:  name,
	}
	return r
}

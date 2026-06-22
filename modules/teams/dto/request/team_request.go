package request

type CreateTeamRequest struct {
	Name string `json:"name" validate:"required,max=100"`
}

type UpdateTeamRequest struct {
	Name string `json:"name" validate:"required,max=100"`
}

type InviteMemberRequest struct {
	UserID uint   `json:"user_id" validate:"required"`
	Role   string `json:"role" validate:"omitempty,oneof=admin member"`
}

type UpdateMemberRequest struct {
	Role string `json:"role" validate:"required,oneof=admin member"`
}

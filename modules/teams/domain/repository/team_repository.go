package repository

import (
	"context"
	"ping-uptime/modules/teams/domain/entity"
)

type TeamRepository interface {
	// Team
	Create(ctx context.Context, team *entity.Team) error
	FindAll(ctx context.Context) ([]*entity.Team, error)
	FindByID(ctx context.Context, id uint) (*entity.Team, error)
	Update(ctx context.Context, team *entity.Team) error
	Delete(ctx context.Context, id uint) error

	// TeamMember
	CreateMember(ctx context.Context, member *entity.TeamMember) error
	FindMembersByTeamID(ctx context.Context, teamID uint) ([]*entity.TeamMember, error)
	FindMemberByUserIDAndTeam(ctx context.Context, userID, teamID uint) (*entity.TeamMember, error)
	FindMemberByID(ctx context.Context, id uint) (*entity.TeamMember, error)
	UpdateMember(ctx context.Context, member *entity.TeamMember) error
	DeleteMember(ctx context.Context, id uint) error
	FindPendingByUserID(ctx context.Context, userID uint) ([]*entity.TeamMember, error)
}

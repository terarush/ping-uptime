package repository

import (
	"context"
	"ping-uptime/internal/pkg/database"
	"ping-uptime/modules/teams/domain/entity"
)

type TeamRepositoryImpl struct{}

func NewTeamRepositoryImpl() TeamRepository {
	return TeamRepositoryImpl{}
}

// Team

func (r TeamRepositoryImpl) Create(ctx context.Context, team *entity.Team) error {
	return database.DB.WithContext(ctx).Create(team).Error
}

func (r TeamRepositoryImpl) FindAll(ctx context.Context) ([]*entity.Team, error) {
	var teams []*entity.Team
	result := database.DB.WithContext(ctx).Find(&teams)
	if result.Error != nil {
		return nil, result.Error
	}
	return teams, nil
}

func (r TeamRepositoryImpl) FindByID(ctx context.Context, id uint) (*entity.Team, error) {
	var team entity.Team
	result := database.DB.WithContext(ctx).First(&team, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &team, nil
}

func (r TeamRepositoryImpl) Update(ctx context.Context, team *entity.Team) error {
	return database.DB.WithContext(ctx).Save(team).Error
}

func (r TeamRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return database.DB.WithContext(ctx).Delete(&entity.Team{}, id).Error
}

// TeamMember

func (r TeamRepositoryImpl) CreateMember(ctx context.Context, member *entity.TeamMember) error {
	return database.DB.WithContext(ctx).Create(member).Error
}

func (r TeamRepositoryImpl) FindMembersByTeamID(ctx context.Context, teamID uint) ([]*entity.TeamMember, error) {
	var members []*entity.TeamMember
	result := database.DB.WithContext(ctx).Where("team_id = ?", teamID).Find(&members)
	if result.Error != nil {
		return nil, result.Error
	}
	return members, nil
}

func (r TeamRepositoryImpl) FindMemberByUserIDAndTeam(ctx context.Context, userID, teamID uint) (*entity.TeamMember, error) {
	var member entity.TeamMember
	result := database.DB.WithContext(ctx).
		Where("user_id = ? AND team_id = ?", userID, teamID).
		First(&member)
	if result.Error != nil {
		return nil, result.Error
	}
	return &member, nil
}

func (r TeamRepositoryImpl) FindMemberByID(ctx context.Context, id uint) (*entity.TeamMember, error) {
	var member entity.TeamMember
	result := database.DB.WithContext(ctx).First(&member, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &member, nil
}

func (r TeamRepositoryImpl) UpdateMember(ctx context.Context, member *entity.TeamMember) error {
	return database.DB.WithContext(ctx).Save(member).Error
}

func (r TeamRepositoryImpl) DeleteMember(ctx context.Context, id uint) error {
	return database.DB.WithContext(ctx).Delete(&entity.TeamMember{}, id).Error
}

func (r TeamRepositoryImpl) FindPendingByUserID(ctx context.Context, userID uint) ([]*entity.TeamMember, error) {
	var members []*entity.TeamMember
	result := database.DB.WithContext(ctx).
		Where("user_id = ? AND status = ?", userID, "pending").
		Find(&members)
	if result.Error != nil {
		return nil, result.Error
	}
	return members, nil
}

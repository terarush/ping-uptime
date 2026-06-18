package service

import (
	"context"
	"errors"
	"ping-uptime/modules/teams/domain/entity"
	"ping-uptime/modules/teams/domain/repository"
)

var (
	ErrTeamNotFound      = errors.New("team not found")
	ErrMemberNotFound    = errors.New("team member not found")
	ErrAlreadyMember     = errors.New("user is already a member of this team")
	ErrNotTeamAdmin      = errors.New("only team admins can perform this action")
	ErrNotPending        = errors.New("invitation is not in pending status")
	ErrCannotRemoveOwner = errors.New("cannot remove the last admin from the team")
)

type TeamService struct {
	repo repository.TeamRepository
}

func NewTeamService(repo repository.TeamRepository) *TeamService {
	return &TeamService{repo: repo}
}

// Team CRUD

func (s *TeamService) Create(ctx context.Context, team *entity.Team, ownerID uint) error {
	if err := s.repo.Create(ctx, team); err != nil {
		return err
	}
	// Add creator as admin member with accepted status
	member := &entity.TeamMember{
		TeamID:    team.ID,
		UserID:    ownerID,
		Role:      "admin",
		InvitedBy: ownerID,
		Status:    "accepted",
	}
	return s.repo.CreateMember(ctx, member)
}

func (s *TeamService) GetAll(ctx context.Context) ([]*entity.Team, error) {
	return s.repo.FindAll(ctx)
}

func (s *TeamService) GetByID(ctx context.Context, id uint) (*entity.Team, error) {
	team, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrTeamNotFound
	}
	return team, nil
}

func (s *TeamService) Update(ctx context.Context, id uint, name string) (*entity.Team, error) {
	team, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrTeamNotFound
	}
	team.Name = name
	if err := s.repo.Update(ctx, team); err != nil {
		return nil, err
	}
	return team, nil
}

func (s *TeamService) Delete(ctx context.Context, id uint) error {
	// Delete all members first, then team
	members, err := s.repo.FindMembersByTeamID(ctx, id)
	if err != nil {
		return err
	}
	for _, m := range members {
		if err := s.repo.DeleteMember(ctx, m.ID); err != nil {
			return err
		}
	}
	return s.repo.Delete(ctx, id)
}

// Member management

func (s *TeamService) GetMembers(ctx context.Context, teamID uint) ([]*entity.TeamMember, error) {
	return s.repo.FindMembersByTeamID(ctx, teamID)
}

func (s *TeamService) InviteMember(ctx context.Context, teamID, invitedByUserID, userID uint, role string) (*entity.TeamMember, error) {
	if role == "" {
		role = "member"
	}

	// Check if already a member
	existing, err := s.repo.FindMemberByUserIDAndTeam(ctx, userID, teamID)
	if err == nil && existing != nil {
		return nil, ErrAlreadyMember
	}

	member := &entity.TeamMember{
		TeamID:    teamID,
		UserID:    userID,
		Role:      role,
		InvitedBy: invitedByUserID,
		Status:    "pending",
	}
	if err := s.repo.CreateMember(ctx, member); err != nil {
		return nil, err
	}
	return member, nil
}

func (s *TeamService) AcceptInvitation(ctx context.Context, teamID, userID uint) error {
	member, err := s.repo.FindMemberByUserIDAndTeam(ctx, userID, teamID)
	if err != nil {
		return ErrMemberNotFound
	}
	if member.Status != "pending" {
		return ErrNotPending
	}
	member.Status = "accepted"
	return s.repo.UpdateMember(ctx, member)
}

func (s *TeamService) RejectInvitation(ctx context.Context, teamID, userID uint) error {
	member, err := s.repo.FindMemberByUserIDAndTeam(ctx, userID, teamID)
	if err != nil {
		return ErrMemberNotFound
	}
	if member.Status != "pending" {
		return ErrNotPending
	}
	member.Status = "rejected"
	return s.repo.UpdateMember(ctx, member)
}

func (s *TeamService) UpdateMemberRole(ctx context.Context, teamID, memberID uint, role string) (*entity.TeamMember, error) {
	member, err := s.repo.FindMemberByID(ctx, memberID)
	if err != nil {
		return nil, ErrMemberNotFound
	}
	if member.TeamID != teamID {
		return nil, ErrMemberNotFound
	}
	member.Role = role
	if err := s.repo.UpdateMember(ctx, member); err != nil {
		return nil, err
	}
	return member, nil
}

func (s *TeamService) RemoveMember(ctx context.Context, teamID, memberID uint) error {
	member, err := s.repo.FindMemberByID(ctx, memberID)
	if err != nil {
		return ErrMemberNotFound
	}
	if member.TeamID != teamID {
		return ErrMemberNotFound
	}
	return s.repo.DeleteMember(ctx, memberID)
}

// Authorization helpers

func (s *TeamService) IsTeamAdmin(ctx context.Context, teamID, userID uint) (bool, error) {
	member, err := s.repo.FindMemberByUserIDAndTeam(ctx, userID, teamID)
	if err != nil {
		return false, nil
	}
	return member.Role == "admin" && member.Status == "accepted", nil
}

func (s *TeamService) GetTeamsByUserID(ctx context.Context, userID uint) ([]*entity.Team, error) {
	_, err := s.repo.FindPendingByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	// For now, get all teams; caller should filter
	return s.repo.FindAll(ctx)
}

func (s *TeamService) GetMemberByID(ctx context.Context, id uint) (*entity.TeamMember, error) {
	return s.repo.FindMemberByID(ctx, id)
}

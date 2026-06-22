package service

import (
	"context"
	"errors"
	"ping-uptime/modules/tags/domain/entity"
	"ping-uptime/modules/tags/domain/repository"
)

var (
	ErrTagNotFound     = errors.New("tag not found")
	ErrTagNameEmpty    = errors.New("tag name cannot be empty")
	ErrTagNameTaken    = errors.New("tag name already exists")
	ErrMonitorNotFound = errors.New("monitor not found")
)

type TagService struct {
	repo repository.TagRepository
}

func NewTagService(repo repository.TagRepository) *TagService {
	return &TagService{repo: repo}
}

func (s *TagService) Create(ctx context.Context, name, color string) (*entity.Tag, error) {
	if name == "" {
		return nil, ErrTagNameEmpty
	}

	existing, err := s.repo.FindByName(ctx, name)
	if err == nil && existing != nil {
		return nil, ErrTagNameTaken
	}

	if color == "" {
		color = "#6366f1"
	}

	tag := &entity.Tag{
		Name:  name,
		Color: color,
	}
	err = s.repo.Create(ctx, tag)
	return tag, err
}

func (s *TagService) GetAll(ctx context.Context) ([]*entity.Tag, error) {
	return s.repo.FindAll(ctx)
}

func (s *TagService) FindByID(ctx context.Context, id uint) (*entity.Tag, error) {
	tag, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrTagNotFound
	}
	return tag, nil
}

func (s *TagService) Update(ctx context.Context, id uint, name, color string) (*entity.Tag, error) {
	tag, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrTagNotFound
	}

	if name != "" {
		tag.Name = name
	}
	if color != "" {
		tag.Color = color
	}

	err = s.repo.Update(ctx, tag)
	return tag, err
}

func (s *TagService) Delete(ctx context.Context, id uint) error {
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return ErrTagNotFound
	}
	return s.repo.Delete(ctx, id)
}

func (s *TagService) AttachTags(ctx context.Context, monitorID uint, tagIDs []uint) error {
	for _, tagID := range tagIDs {
		_, err := s.repo.FindByID(ctx, tagID)
		if err != nil {
			return ErrTagNotFound
		}
	}

	// Clear existing then attach new — replace semantics
	if err := s.repo.DeleteMonitorTags(ctx, monitorID); err != nil {
		return err
	}

	for _, tagID := range tagIDs {
		if err := s.repo.AttachToMonitor(ctx, monitorID, tagID); err != nil {
			return err
		}
	}
	return nil
}

func (s *TagService) DetachTag(ctx context.Context, monitorID, tagID uint) error {
	return s.repo.DetachFromMonitor(ctx, monitorID, tagID)
}

func (s *TagService) GetMonitorTags(ctx context.Context, monitorID uint) ([]*entity.Tag, error) {
	return s.repo.GetTagsByMonitorID(ctx, monitorID)
}

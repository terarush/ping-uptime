package service

import (
	"context"
	"errors"
	"ping-uptime/modules/status_pages/domain/entity"
	"ping-uptime/modules/status_pages/domain/repository"
)

var (
	ErrStatusPageNotFound = errors.New("status page not found")
	ErrSlugAlreadyTaken   = errors.New("slug already taken")
)

type StatusPageService struct {
	pageRepo repository.StatusPageRepository
}

func NewStatusPageService(pageRepo repository.StatusPageRepository) *StatusPageService {
	return &StatusPageService{
		pageRepo: pageRepo,
	}
}

func (s *StatusPageService) CreateStatusPage(ctx context.Context, page *entity.StatusPage) error {
	if page.Name == "" || page.Slug == "" {
		return errors.New("name and slug cannot be empty")
	}

	// Check if slug is unique
	existing, _ := s.pageRepo.FindBySlug(ctx, page.Slug)
	if existing != nil {
		return ErrSlugAlreadyTaken
	}

	return s.pageRepo.Create(ctx, page)
}

func (s *StatusPageService) DeleteStatusPage(ctx context.Context, id uint) error {
	return s.pageRepo.Delete(ctx, id)
}

func (s *StatusPageService) GetAllStatusPages(ctx context.Context) ([]*entity.StatusPage, error) {
	return s.pageRepo.FindAll(ctx)
}

func (s *StatusPageService) GetStatusPagesByUserID(ctx context.Context, userID uint) ([]*entity.StatusPage, error) {
	return s.pageRepo.FindByUserID(ctx, userID)
}

func (s *StatusPageService) GetStatusPageByID(ctx context.Context, id uint) (*entity.StatusPage, error) {
	return s.pageRepo.FindByID(ctx, id)
}

func (s *StatusPageService) GetStatusPageBySlug(ctx context.Context, slug string) (*entity.StatusPage, error) {
	return s.pageRepo.FindBySlug(ctx, slug)
}

func (s *StatusPageService) UpdateStatusPage(ctx context.Context, page *entity.StatusPage) error {
	// Check if slug is updated and taken by another page
	existing, _ := s.pageRepo.FindBySlug(ctx, page.Slug)
	if existing != nil && existing.ID != page.ID {
		return ErrSlugAlreadyTaken
	}

	return s.pageRepo.Update(ctx, page)
}

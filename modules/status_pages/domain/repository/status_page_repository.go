package repository

import (
	"context"
	"ping-uptime/modules/status_pages/domain/entity"
)

type StatusPageRepository interface {
	Create(ctx context.Context, statusPage *entity.StatusPage) error
	Delete(ctx context.Context, id uint) error
	FindAll(ctx context.Context) ([]*entity.StatusPage, error)
	FindByUserID(ctx context.Context, userID uint) ([]*entity.StatusPage, error)
	FindByID(ctx context.Context, id uint) (*entity.StatusPage, error)
	FindBySlug(ctx context.Context, slug string) (*entity.StatusPage, error)
	Update(ctx context.Context, statusPage *entity.StatusPage) error
}

package repository

import (
	"context"
	"ping-uptime/modules/subscribers/domain/entity"
)

type SubscriberRepository interface {
	Create(ctx context.Context, s *entity.Subscriber) error
	FindByToken(ctx context.Context, token string) (*entity.Subscriber, error)
	FindByEmailAndPage(ctx context.Context, email string, pageID uint) (*entity.Subscriber, error)
	FindByPageID(ctx context.Context, pageID uint) ([]*entity.Subscriber, error)
	Delete(ctx context.Context, id uint) error
	Update(ctx context.Context, s *entity.Subscriber) error
	CountByPageID(ctx context.Context, pageID uint) (int64, error)
}

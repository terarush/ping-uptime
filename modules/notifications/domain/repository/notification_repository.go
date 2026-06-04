package repository

import (
	"context"
	"ping-uptime/modules/notifications/domain/entity"
)

type NotificationRepository interface {
	Create(ctx context.Context, channel *entity.NotificationChannel) error
	Delete(ctx context.Context, id uint) error
	FindAll(ctx context.Context) ([]*entity.NotificationChannel, error)
	FindByUserID(ctx context.Context, userID uint) ([]*entity.NotificationChannel, error)
	FindByID(ctx context.Context, id uint) (*entity.NotificationChannel, error)
	Update(ctx context.Context, channel *entity.NotificationChannel) error
}

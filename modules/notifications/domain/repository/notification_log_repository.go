package repository

import (
	"context"
	"ping-uptime/modules/notifications/domain/entity"
)

type NotificationLogRepository interface {
	Create(ctx context.Context, log *entity.NotificationLog) error
	FindAll(ctx context.Context, userID uint, isAdmin bool, filter map[string]string) ([]*entity.NotificationLog, error)
	FindByID(ctx context.Context, id uint) (*entity.NotificationLog, error)
}

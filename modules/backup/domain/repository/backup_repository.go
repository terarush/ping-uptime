package repository

import (
	"context"
	"ping-uptime/modules/backup/domain/entity"
)

type BackupRepository interface {
	Create(ctx context.Context, r *entity.BackupRecord) error
	FindAll(ctx context.Context) ([]*entity.BackupRecord, error)
	Delete(ctx context.Context, id uint) error
}

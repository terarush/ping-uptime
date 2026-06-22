package repository

import (
	"context"
	"ping-uptime/internal/pkg/database"
	"ping-uptime/modules/backup/domain/entity"
)

type BackupRepositoryImpl struct{}

func (r *BackupRepositoryImpl) Create(ctx context.Context, record *entity.BackupRecord) error {
	return database.DB.WithContext(ctx).Create(record).Error
}

func (r *BackupRepositoryImpl) FindAll(ctx context.Context) ([]*entity.BackupRecord, error) {
	var items []*entity.BackupRecord
	err := database.DB.WithContext(ctx).Order("created_at DESC").Find(&items).Error
	return items, err
}

func (r *BackupRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return database.DB.WithContext(ctx).Delete(&entity.BackupRecord{}, id).Error
}

func NewBackupRepositoryImpl() BackupRepository {
	return &BackupRepositoryImpl{}
}

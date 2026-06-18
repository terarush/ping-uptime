package repository

import (
	"context"
	"ping-uptime/internal/pkg/database"
	"ping-uptime/modules/audit_logs/domain/entity"
)

type AuditLogRepositoryImpl struct{}

func (r AuditLogRepositoryImpl) Create(ctx context.Context, log *entity.AuditLog) error {
	return database.DB.WithContext(ctx).Create(log).Error
}

func (r AuditLogRepositoryImpl) FindAll(ctx context.Context) ([]*entity.AuditLog, error) {
	var items []*entity.AuditLog
	err := database.DB.WithContext(ctx).Order("created_at DESC").Limit(200).Find(&items).Error
	return items, err
}

func (r AuditLogRepositoryImpl) FindByUserID(ctx context.Context, userID uint) ([]*entity.AuditLog, error) {
	var items []*entity.AuditLog
	err := database.DB.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Limit(200).Find(&items).Error
	return items, err
}

func (r AuditLogRepositoryImpl) FindByEntity(ctx context.Context, entityType string, entityID uint) ([]*entity.AuditLog, error) {
	var items []*entity.AuditLog
	err := database.DB.WithContext(ctx).Where("entity_type = ? AND entity_id = ?", entityType, entityID).Order("created_at DESC").Limit(100).Find(&items).Error
	return items, err
}

func NewAuditLogRepositoryImpl() AuditLogRepository {
	return AuditLogRepositoryImpl{}
}

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

func (r AuditLogRepositoryImpl) FindFiltered(ctx context.Context, filter AuditLogFilter) ([]*entity.AuditLog, error) {
	query := database.DB.WithContext(ctx).Model(&entity.AuditLog{}).Order("created_at DESC")

	if filter.UserID != nil && *filter.UserID > 0 {
		query = query.Where("user_id = ?", *filter.UserID)
	}
	if filter.EntityType != nil && *filter.EntityType != "" {
		query = query.Where("entity_type = ?", *filter.EntityType)
	}
	if filter.Action != nil && *filter.Action != "" {
		query = query.Where("action = ?", *filter.Action)
	}
	if filter.From != nil && *filter.From != "" {
		query = query.Where("created_at >= ?", *filter.From)
	}
	if filter.To != nil && *filter.To != "" {
		query = query.Where("created_at <= ?", *filter.To)
	}

	limit := filter.Limit
	if limit <= 0 || limit > 500 {
		limit = 200
	}

	var items []*entity.AuditLog
	err := query.Limit(limit).Find(&items).Error
	return items, err
}

func NewAuditLogRepositoryImpl() AuditLogRepository {
	return AuditLogRepositoryImpl{}
}

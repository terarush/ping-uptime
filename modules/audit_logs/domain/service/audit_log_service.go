package service

import (
	"context"
	"ping-uptime/modules/audit_logs/domain/entity"
	"ping-uptime/modules/audit_logs/domain/repository"
)

type AuditLogService struct {
	repo repository.AuditLogRepository
}

func NewAuditLogService(repo repository.AuditLogRepository) *AuditLogService {
	return &AuditLogService{repo: repo}
}

func (s *AuditLogService) Log(ctx context.Context, userID uint, action, entityType string, entityID uint, details string) error {
	return s.repo.Create(ctx, &entity.AuditLog{
		UserID:     userID,
		Action:     action,
		EntityType: entityType,
		EntityID:   entityID,
		Details:    details,
	})
}

func (s *AuditLogService) GetAll(ctx context.Context) ([]*entity.AuditLog, error) {
	return s.repo.FindAll(ctx)
}

func (s *AuditLogService) GetByUserID(ctx context.Context, userID uint) ([]*entity.AuditLog, error) {
	return s.repo.FindByUserID(ctx, userID)
}

func (s *AuditLogService) GetByEntity(ctx context.Context, entityType string, entityID uint) ([]*entity.AuditLog, error) {
	return s.repo.FindByEntity(ctx, entityType, entityID)
}

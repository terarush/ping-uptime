package repository

import (
	"context"
	"ping-uptime/modules/audit_logs/domain/entity"
)

type AuditLogRepository interface {
	Create(ctx context.Context, log *entity.AuditLog) error
	FindAll(ctx context.Context) ([]*entity.AuditLog, error)
	FindByUserID(ctx context.Context, userID uint) ([]*entity.AuditLog, error)
	FindByEntity(ctx context.Context, entityType string, entityID uint) ([]*entity.AuditLog, error)
}

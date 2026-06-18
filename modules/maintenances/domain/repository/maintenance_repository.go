package repository

import (
	"context"
	"ping-uptime/modules/maintenances/domain/entity"
)

type MaintenanceRepository interface {
	Create(ctx context.Context, m *entity.Maintenance) error
	Delete(ctx context.Context, id uint) error
	FindAll(ctx context.Context) ([]*entity.Maintenance, error)
	FindByUserID(ctx context.Context, userID uint) ([]*entity.Maintenance, error)
	FindByID(ctx context.Context, id uint) (*entity.Maintenance, error)
	Update(ctx context.Context, m *entity.Maintenance) error
	FindActiveByMonitorID(ctx context.Context, monitorID uint) ([]*entity.Maintenance, error)
	SetMonitorIDs(ctx context.Context, maintenanceID uint, monitorIDs []uint) error
	GetMonitorIDs(ctx context.Context, maintenanceID uint) ([]uint, error)
}

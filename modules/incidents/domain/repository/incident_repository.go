package repository

import (
	"context"
	"ping-uptime/modules/incidents/domain/entity"
)

type IncidentRepository interface {
	Create(ctx context.Context, incident *entity.Incident) error
	Delete(ctx context.Context, id uint) error
	FindAll(ctx context.Context) ([]*entity.Incident, error)
	FindByUserID(ctx context.Context, userID uint) ([]*entity.Incident, error)
	FindByID(ctx context.Context, id uint) (*entity.Incident, error)
	FindByMonitorID(ctx context.Context, monitorID uint) ([]*entity.Incident, error)
	Update(ctx context.Context, incident *entity.Incident) error
}

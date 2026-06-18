package repository

import (
	"context"
	"ping-uptime/modules/integrations/domain/entity"
)

type IntegrationRepository interface {
	Create(ctx context.Context, integration *entity.Integration) error
	FindAll(ctx context.Context) ([]*entity.Integration, error)
	FindByID(ctx context.Context, id uint) (*entity.Integration, error)
	FindByType(ctx context.Context, integrationType string) ([]*entity.Integration, error)
	Update(ctx context.Context, integration *entity.Integration) error
	Delete(ctx context.Context, id uint) error
}

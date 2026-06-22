package repository

import (
	"context"
	"ping-uptime/modules/ssl_monitors/domain/entity"
)

type SSLRepository interface {
	Create(ctx context.Context, cert *entity.SSLCert) error
	Update(ctx context.Context, cert *entity.SSLCert) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*entity.SSLCert, error)
	FindByMonitorID(ctx context.Context, monitorID uint) (*entity.SSLCert, error)
	FindAll(ctx context.Context) ([]*entity.SSLCert, error)
	FindExpiring(ctx context.Context, days int) ([]*entity.SSLCert, error)
}

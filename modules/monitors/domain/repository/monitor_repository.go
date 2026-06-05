package repository

import (
	"context"
	"ping-uptime/modules/monitors/domain/entity"
)

type MonitorRepository interface {
	Create(ctx context.Context, monitor *entity.Monitor) error
	Delete(ctx context.Context, id uint) error
	FindAll(ctx context.Context) ([]*entity.Monitor, error)
	FindByUserID(ctx context.Context, userID uint) ([]*entity.Monitor, error)
	FindByID(ctx context.Context, id uint) (*entity.Monitor, error)
	Update(ctx context.Context, monitor *entity.Monitor) error
	CreateCheckRecord(ctx context.Context, rec *entity.CheckRecord) error
	GetDailyChart(ctx context.Context, monitorID uint, days int) ([]entity.DailyChartPoint, error)
}

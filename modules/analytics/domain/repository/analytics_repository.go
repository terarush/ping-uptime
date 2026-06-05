package repository

import (
	"context"
	"ping-uptime/modules/analytics/domain/entity"
)

type AnalyticsRepository interface {
	GetChartData(ctx context.Context, monitorID uint, window string) ([]entity.ChartDataPoint, error)
	GetMonitorStats(ctx context.Context, userID uint, window string) ([]entity.MonitorStats, error)
}

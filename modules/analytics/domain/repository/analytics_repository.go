package repository

import "context"

type AnalyticsRepository interface {
	GetChartData(ctx context.Context, monitorID uint, window string) ([]entity.ChartDataPoint, error)
	GetMonitorStats(ctx context.Context, userID uint, window string) ([]entity.MonitorStats, error)
}

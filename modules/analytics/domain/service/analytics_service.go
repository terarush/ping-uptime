package service

import (
	"context"
	"ping-uptime/internal/pkg/bus"
	"ping-uptime/modules/analytics/domain/entity"
	"ping-uptime/modules/analytics/domain/repository"
)

type AnalyticsService struct {
	analyticsRepo repository.AnalyticsRepository
	event         *bus.EventBus
}

func NewAnalyticsService(analyticsRepo repository.AnalyticsRepository, event *bus.EventBus) *AnalyticsService {
	return &AnalyticsService{
		analyticsRepo: analyticsRepo,
		event:         event,
	}
}

func (s *AnalyticsService) GetChartData(ctx context.Context, monitorID uint, window string) ([]entity.ChartDataPoint, error) {
	return s.analyticsRepo.GetChartData(ctx, monitorID, window)
}

func (s *AnalyticsService) GetMonitorStats(ctx context.Context, userID uint, window string) ([]entity.MonitorStats, error) {
	return s.analyticsRepo.GetMonitorStats(ctx, userID, window)
}

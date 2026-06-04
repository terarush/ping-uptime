package service

import (
	"context"
	"errors"
	"ping-uptime/internal/pkg/bus"
	"ping-uptime/modules/monitors/domain/entity"
	"ping-uptime/modules/monitors/domain/repository"
)

var (
	ErrMonitorNotFound = errors.New("monitor not found")
)

type MonitorService struct {
	monitorRepo repository.MonitorRepository
	event       *bus.EventBus
}

func NewMonitorService(monitorRepo repository.MonitorRepository, event *bus.EventBus) *MonitorService {
	return &MonitorService{
		monitorRepo: monitorRepo,
		event:       event,
	}
}

func (s *MonitorService) CreateMonitor(ctx context.Context, monitor *entity.Monitor) error {
	if monitor.Name == "" || monitor.URL == "" {
		return errors.New("name and URL cannot be empty")
	}
	return s.monitorRepo.Create(ctx, monitor)
}

func (s *MonitorService) DeleteMonitor(ctx context.Context, id uint) error {
	return s.monitorRepo.Delete(ctx, id)
}

func (s *MonitorService) GetAllMonitors(ctx context.Context) ([]*entity.Monitor, error) {
	return s.monitorRepo.FindAll(ctx)
}

func (s *MonitorService) GetMonitorsByUserID(ctx context.Context, userID uint) ([]*entity.Monitor, error) {
	return s.monitorRepo.FindByUserID(ctx, userID)
}

func (s *MonitorService) GetMonitorByID(ctx context.Context, id uint) (*entity.Monitor, error) {
	return s.monitorRepo.FindByID(ctx, id)
}

func (s *MonitorService) UpdateMonitor(ctx context.Context, monitor *entity.Monitor) error {
	return s.monitorRepo.Update(ctx, monitor)
}

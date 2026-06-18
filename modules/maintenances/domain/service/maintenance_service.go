package service

import (
	"context"
	"errors"
	"ping-uptime/modules/maintenances/domain/entity"
	"ping-uptime/modules/maintenances/domain/repository"
)

var (
	ErrMaintenanceNotFound = errors.New("maintenance not found")
)

type MaintenanceService struct {
	repo repository.MaintenanceRepository
}

func NewMaintenanceService(repo repository.MaintenanceRepository) *MaintenanceService {
	return &MaintenanceService{repo: repo}
}

func (s *MaintenanceService) Create(ctx context.Context, m *entity.Maintenance) error {
	if m.Name == "" {
		return errors.New("name cannot be empty")
	}
	return s.repo.Create(ctx, m)
}

func (s *MaintenanceService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func (s *MaintenanceService) GetAll(ctx context.Context) ([]*entity.Maintenance, error) {
	return s.repo.FindAll(ctx)
}

func (s *MaintenanceService) GetByUserID(ctx context.Context, userID uint) ([]*entity.Maintenance, error) {
	return s.repo.FindByUserID(ctx, userID)
}

func (s *MaintenanceService) GetByID(ctx context.Context, id uint) (*entity.Maintenance, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *MaintenanceService) Update(ctx context.Context, m *entity.Maintenance) error {
	return s.repo.Update(ctx, m)
}

func (s *MaintenanceService) IsMonitorInMaintenance(ctx context.Context, monitorID uint) (bool, error) {
	items, err := s.repo.FindActiveByMonitorID(ctx, monitorID)
	return len(items) > 0, err
}

func (s *MaintenanceService) SetMonitorIDs(ctx context.Context, maintenanceID uint, monitorIDs []uint) error {
	return s.repo.SetMonitorIDs(ctx, maintenanceID, monitorIDs)
}

func (s *MaintenanceService) GetMonitorIDs(ctx context.Context, maintenanceID uint) ([]uint, error) {
	return s.repo.GetMonitorIDs(ctx, maintenanceID)
}

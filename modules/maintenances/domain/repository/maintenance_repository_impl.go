package repository

import (
	"context"
	"ping-uptime/internal/pkg/database"
	"ping-uptime/modules/maintenances/domain/entity"
	"time"
)

type MaintenanceRepositoryImpl struct{}

func (r MaintenanceRepositoryImpl) Create(ctx context.Context, m *entity.Maintenance) error {
	return database.DB.WithContext(ctx).Create(m).Error
}

func (r MaintenanceRepositoryImpl) Delete(ctx context.Context, id uint) error {
	database.DB.WithContext(ctx).Where("maintenance_id = ?", id).Delete(&entity.MaintenanceMonitor{})
	return database.DB.WithContext(ctx).Delete(&entity.Maintenance{}, id).Error
}

func (r MaintenanceRepositoryImpl) FindAll(ctx context.Context) ([]*entity.Maintenance, error) {
	var items []*entity.Maintenance
	result := database.DB.WithContext(ctx).Order("start_at DESC").Find(&items)
	return items, result.Error
}

func (r MaintenanceRepositoryImpl) FindByUserID(ctx context.Context, userID uint) ([]*entity.Maintenance, error) {
	var items []*entity.Maintenance
	result := database.DB.WithContext(ctx).Where("user_id = ?", userID).Order("start_at DESC").Find(&items)
	return items, result.Error
}

func (r MaintenanceRepositoryImpl) FindByID(ctx context.Context, id uint) (*entity.Maintenance, error) {
	var m entity.Maintenance
	result := database.DB.WithContext(ctx).First(&m, id)
	return &m, result.Error
}

func (r MaintenanceRepositoryImpl) Update(ctx context.Context, m *entity.Maintenance) error {
	return database.DB.WithContext(ctx).Save(m).Error
}

func (r MaintenanceRepositoryImpl) FindActiveByMonitorID(ctx context.Context, monitorID uint) ([]*entity.Maintenance, error) {
	now := time.Now()
	var items []*entity.Maintenance
	err := database.DB.WithContext(ctx).
		Joins("JOIN maintenance_monitors ON maintenance_monitors.maintenance_id = maintenances.id").
		Where("maintenance_monitors.monitor_id = ?", monitorID).
		Where("start_at <= ? AND end_at >= ?", now, now).
		Find(&items).Error
	return items, err
}

func (r MaintenanceRepositoryImpl) SetMonitorIDs(ctx context.Context, maintenanceID uint, monitorIDs []uint) error {
	database.DB.WithContext(ctx).Where("maintenance_id = ?", maintenanceID).Delete(&entity.MaintenanceMonitor{})
	for _, mid := range monitorIDs {
		database.DB.WithContext(ctx).Create(&entity.MaintenanceMonitor{MaintenanceID: maintenanceID, MonitorID: mid})
	}
	return nil
}

func (r MaintenanceRepositoryImpl) GetMonitorIDs(ctx context.Context, maintenanceID uint) ([]uint, error) {
	var rows []entity.MaintenanceMonitor
	err := database.DB.WithContext(ctx).Where("maintenance_id = ?", maintenanceID).Find(&rows).Error
	if err != nil {
		return nil, err
	}
	ids := make([]uint, len(rows))
	for i, r := range rows {
		ids[i] = r.MonitorID
	}
	return ids, nil
}

func NewMaintenanceRepositoryImpl() MaintenanceRepository {
	return MaintenanceRepositoryImpl{}
}

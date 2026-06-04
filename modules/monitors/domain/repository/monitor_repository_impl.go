package repository

import (
	"context"
	"errors"
	"ping-uptime/internal/pkg/database"
	"ping-uptime/modules/monitors/domain/entity"
)

var (
	ERR_RECORD_NOT_FOUND = errors.New("record not found")
)

type MonitorRepositoryImpl struct{}

func (r MonitorRepositoryImpl) Create(ctx context.Context, monitor *entity.Monitor) error {
	return database.DB.WithContext(ctx).Create(monitor).Error
}

func (r MonitorRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return database.DB.WithContext(ctx).Delete(&entity.Monitor{}, id).Error
}

func (r MonitorRepositoryImpl) FindAll(ctx context.Context) ([]*entity.Monitor, error) {
	var monitors []*entity.Monitor
	result := database.DB.WithContext(ctx).Find(&monitors)
	if result.Error != nil {
		return nil, result.Error
	}
	return monitors, nil
}

func (r MonitorRepositoryImpl) FindByUserID(ctx context.Context, userID uint) ([]*entity.Monitor, error) {
	var monitors []*entity.Monitor
	result := database.DB.WithContext(ctx).Where("user_id = ?", userID).Find(&monitors)
	if result.Error != nil {
		return nil, result.Error
	}
	return monitors, nil
}

func (r MonitorRepositoryImpl) FindByID(ctx context.Context, id uint) (*entity.Monitor, error) {
	var monitor entity.Monitor
	result := database.DB.WithContext(ctx).First(&monitor, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &monitor, nil
}

func (r MonitorRepositoryImpl) Update(ctx context.Context, monitor *entity.Monitor) error {
	return database.DB.WithContext(ctx).Save(monitor).Error
}

func NewMonitorRepositoryImpl() MonitorRepository {
	return MonitorRepositoryImpl{}
}

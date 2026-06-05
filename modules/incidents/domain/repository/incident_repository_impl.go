package repository

import (
	"context"
	"errors"
	"ping-uptime/internal/pkg/database"
	"ping-uptime/modules/incidents/domain/entity"
)

var (
	ERR_RECORD_NOT_FOUND = errors.New("record not found")
)

type IncidentRepositoryImpl struct{}

func (r IncidentRepositoryImpl) Create(ctx context.Context, incident *entity.Incident) error {
	return database.DB.WithContext(ctx).Create(incident).Error
}

func (r IncidentRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return database.DB.WithContext(ctx).Delete(&entity.Incident{}, id).Error
}

func (r IncidentRepositoryImpl) FindAll(ctx context.Context) ([]*entity.Incident, error) {
	var incidents []*entity.Incident
	result := database.DB.WithContext(ctx).Order("created_at DESC").Find(&incidents)
	if result.Error != nil {
		return nil, result.Error
	}
	return incidents, nil
}

func (r IncidentRepositoryImpl) FindByUserID(ctx context.Context, userID uint) ([]*entity.Incident, error) {
	var incidents []*entity.Incident
	result := database.DB.WithContext(ctx).Where("user_id = ?", userID).Find(&incidents)
	if result.Error != nil {
		return nil, result.Error
	}
	return incidents, nil
}

func (r IncidentRepositoryImpl) FindByID(ctx context.Context, id uint) (*entity.Incident, error) {
	var incident entity.Incident
	result := database.DB.WithContext(ctx).First(&incident, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &incident, nil
}

func (r IncidentRepositoryImpl) FindByMonitorID(ctx context.Context, monitorID uint) ([]*entity.Incident, error) {
	var incidents []*entity.Incident
	result := database.DB.WithContext(ctx).Where("monitor_id = ?", monitorID).Find(&incidents)
	if result.Error != nil {
		return nil, result.Error
	}
	return incidents, nil
}

func (r IncidentRepositoryImpl) Update(ctx context.Context, incident *entity.Incident) error {
	return database.DB.WithContext(ctx).Save(incident).Error
}

func NewIncidentRepositoryImpl() IncidentRepository {
	return IncidentRepositoryImpl{}
}

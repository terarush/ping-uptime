package repository

import (
	"context"
	"errors"
	"ping-uptime/internal/pkg/database"
	"ping-uptime/modules/integrations/domain/entity"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type IntegrationRepositoryImpl struct{}

func (r IntegrationRepositoryImpl) Create(ctx context.Context, integration *entity.Integration) error {
	return database.DB.WithContext(ctx).Create(integration).Error
}

func (r IntegrationRepositoryImpl) FindAll(ctx context.Context) ([]*entity.Integration, error) {
	var integrations []*entity.Integration
	result := database.DB.WithContext(ctx).Find(&integrations)
	if result.Error != nil {
		return nil, result.Error
	}
	return integrations, nil
}

func (r IntegrationRepositoryImpl) FindByID(ctx context.Context, id uint) (*entity.Integration, error) {
	var integration entity.Integration
	result := database.DB.WithContext(ctx).First(&integration, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &integration, nil
}

func (r IntegrationRepositoryImpl) FindByType(ctx context.Context, integrationType string) ([]*entity.Integration, error) {
	var integrations []*entity.Integration
	result := database.DB.WithContext(ctx).Where("type = ?", integrationType).Find(&integrations)
	if result.Error != nil {
		return nil, result.Error
	}
	return integrations, nil
}

func (r IntegrationRepositoryImpl) Update(ctx context.Context, integration *entity.Integration) error {
	return database.DB.WithContext(ctx).Save(integration).Error
}

func (r IntegrationRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return database.DB.WithContext(ctx).Delete(&entity.Integration{}, id).Error
}

func NewIntegrationRepositoryImpl() IntegrationRepository {
	return IntegrationRepositoryImpl{}
}

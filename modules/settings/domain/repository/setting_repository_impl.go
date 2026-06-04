package repository

import (
	"context"
	"errors"
	"ping-uptime/internal/pkg/database"
	"ping-uptime/modules/settings/domain/entity"
)

var (
	ERR_RECORD_NOT_FOUND = errors.New("record not found")
)

type SettingRepositoryImpl struct{}

func (r SettingRepositoryImpl) CreateOrUpdate(ctx context.Context, setting *entity.Setting) error {
	return database.DB.WithContext(ctx).Save(setting).Error
}

func (r SettingRepositoryImpl) Delete(ctx context.Context, key string) error {
	return database.DB.WithContext(ctx).Delete(&entity.Setting{}, "key = ?", key).Error
}

func (r SettingRepositoryImpl) FindAll(ctx context.Context) ([]*entity.Setting, error) {
	var settings []*entity.Setting
	result := database.DB.WithContext(ctx).Find(&settings)
	if result.Error != nil {
		return nil, result.Error
	}
	return settings, nil
}

func (r SettingRepositoryImpl) FindByKey(ctx context.Context, key string) (*entity.Setting, error) {
	var setting entity.Setting
	result := database.DB.WithContext(ctx).Where("key = ?", key).First(&setting)
	if result.Error != nil {
		return nil, result.Error
	}
	return &setting, nil
}

func NewSettingRepositoryImpl() SettingRepository {
	return SettingRepositoryImpl{}
}

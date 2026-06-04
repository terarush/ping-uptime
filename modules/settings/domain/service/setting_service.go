package service

import (
	"context"
	"errors"
	"ping-uptime/modules/settings/domain/entity"
	"ping-uptime/modules/settings/domain/repository"
)

var (
	ErrSettingNotFound = errors.New("setting not found")
)

type SettingService struct {
	settingRepo repository.SettingRepository
}

func NewSettingService(settingRepo repository.SettingRepository) *SettingService {
	return &SettingService{
		settingRepo: settingRepo,
	}
}

func (s *SettingService) SetSetting(ctx context.Context, setting *entity.Setting) error {
	if setting.Key == "" {
		return errors.New("key cannot be empty")
	}
	return s.settingRepo.CreateOrUpdate(ctx, setting)
}

func (s *SettingService) DeleteSetting(ctx context.Context, key string) error {
	return s.settingRepo.Delete(ctx, key)
}

func (s *SettingService) GetAllSettings(ctx context.Context) ([]*entity.Setting, error) {
	return s.settingRepo.FindAll(ctx)
}

func (s *SettingService) GetSettingByKey(ctx context.Context, key string) (*entity.Setting, error) {
	return s.settingRepo.FindByKey(ctx, key)
}

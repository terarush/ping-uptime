package repository

import (
	"context"
	"errors"
	"ping-uptime/internal/pkg/database"
	"ping-uptime/modules/notifications/domain/entity"
)

var (
	ERR_RECORD_NOT_FOUND = errors.New("record not found")
)

type NotificationRepositoryImpl struct{}

func (r NotificationRepositoryImpl) Create(ctx context.Context, channel *entity.NotificationChannel) error {
	return database.DB.WithContext(ctx).Create(channel).Error
}

func (r NotificationRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return database.DB.WithContext(ctx).Delete(&entity.NotificationChannel{}, id).Error
}

func (r NotificationRepositoryImpl) FindAll(ctx context.Context) ([]*entity.NotificationChannel, error) {
	var channels []*entity.NotificationChannel
	result := database.DB.WithContext(ctx).Find(&channels)
	if result.Error != nil {
		return nil, result.Error
	}
	return channels, nil
}

func (r NotificationRepositoryImpl) FindByUserID(ctx context.Context, userID uint) ([]*entity.NotificationChannel, error) {
	var channels []*entity.NotificationChannel
	result := database.DB.WithContext(ctx).Where("user_id = ?", userID).Find(&channels)
	if result.Error != nil {
		return nil, result.Error
	}
	return channels, nil
}

func (r NotificationRepositoryImpl) FindByID(ctx context.Context, id uint) (*entity.NotificationChannel, error) {
	var channel entity.NotificationChannel
	result := database.DB.WithContext(ctx).First(&channel, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &channel, nil
}

func (r NotificationRepositoryImpl) Update(ctx context.Context, channel *entity.NotificationChannel) error {
	return database.DB.WithContext(ctx).Save(channel).Error
}

func NewNotificationRepositoryImpl() NotificationRepository {
	return NotificationRepositoryImpl{}
}

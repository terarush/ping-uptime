package service

import (
	"context"
	"errors"
	"ping-uptime/modules/notifications/domain/entity"
	"ping-uptime/modules/notifications/domain/repository"
)

var (
	ErrChannelNotFound = errors.New("notification channel not found")
)

type NotificationService struct {
	channelRepo repository.NotificationRepository
}

func NewNotificationService(channelRepo repository.NotificationRepository) *NotificationService {
	return &NotificationService{
		channelRepo: channelRepo,
	}
}

func (s *NotificationService) CreateChannel(ctx context.Context, channel *entity.NotificationChannel) error {
	if channel.Name == "" || channel.Type == "" {
		return errors.New("name and type cannot be empty")
	}
	return s.channelRepo.Create(ctx, channel)
}

func (s *NotificationService) DeleteChannel(ctx context.Context, id uint) error {
	return s.channelRepo.Delete(ctx, id)
}

func (s *NotificationService) GetAllChannels(ctx context.Context) ([]*entity.NotificationChannel, error) {
	return s.channelRepo.FindAll(ctx)
}

func (s *NotificationService) GetChannelsByUserID(ctx context.Context, userID uint) ([]*entity.NotificationChannel, error) {
	return s.channelRepo.FindByUserID(ctx, userID)
}

func (s *NotificationService) GetChannelByID(ctx context.Context, id uint) (*entity.NotificationChannel, error) {
	return s.channelRepo.FindByID(ctx, id)
}

func (s *NotificationService) UpdateChannel(ctx context.Context, channel *entity.NotificationChannel) error {
	return s.channelRepo.Update(ctx, channel)
}

package repository

import (
	"context"
	"ping-uptime/internal/pkg/database"
	"ping-uptime/modules/subscribers/domain/entity"
)

type SubscriberRepositoryImpl struct{}

func (r SubscriberRepositoryImpl) Create(ctx context.Context, s *entity.Subscriber) error {
	return database.DB.WithContext(ctx).Create(s).Error
}

func (r SubscriberRepositoryImpl) FindByToken(ctx context.Context, token string) (*entity.Subscriber, error) {
	var s entity.Subscriber
	err := database.DB.WithContext(ctx).Where("token = ?", token).First(&s).Error
	return &s, err
}

func (r SubscriberRepositoryImpl) FindByEmailAndPage(ctx context.Context, email string, pageID uint) (*entity.Subscriber, error) {
	var s entity.Subscriber
	err := database.DB.WithContext(ctx).Where("email = ? AND status_page_id = ?", email, pageID).First(&s).Error
	return &s, err
}

func (r SubscriberRepositoryImpl) FindByPageID(ctx context.Context, pageID uint) ([]*entity.Subscriber, error) {
	var items []*entity.Subscriber
	err := database.DB.WithContext(ctx).Where("status_page_id = ? AND verified = ?", pageID, true).Find(&items).Error
	return items, err
}

func (r SubscriberRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return database.DB.WithContext(ctx).Delete(&entity.Subscriber{}, id).Error
}

func (r SubscriberRepositoryImpl) Update(ctx context.Context, s *entity.Subscriber) error {
	return database.DB.WithContext(ctx).Save(s).Error
}

func (r SubscriberRepositoryImpl) CountByPageID(ctx context.Context, pageID uint) (int64, error) {
	var count int64
	err := database.DB.WithContext(ctx).Model(&entity.Subscriber{}).Where("status_page_id = ? AND verified = ?", pageID, true).Count(&count).Error
	return count, err
}

func NewSubscriberRepositoryImpl() SubscriberRepository {
	return SubscriberRepositoryImpl{}
}

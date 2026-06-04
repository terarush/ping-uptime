package repository

import (
	"context"
	"errors"
	"ping-uptime/internal/pkg/database"
	"ping-uptime/modules/status_pages/domain/entity"
)

var (
	ERR_RECORD_NOT_FOUND = errors.New("record not found")
)

type StatusPageRepositoryImpl struct{}

func (r StatusPageRepositoryImpl) Create(ctx context.Context, statusPage *entity.StatusPage) error {
	return database.DB.WithContext(ctx).Create(statusPage).Error
}

func (r StatusPageRepositoryImpl) Delete(ctx context.Context, id uint) error {
	var sp entity.StatusPage
	if err := database.DB.WithContext(ctx).First(&sp, id).Error; err == nil {
		database.DB.WithContext(ctx).Model(&sp).Association("Monitors").Clear()
	}
	return database.DB.WithContext(ctx).Delete(&entity.StatusPage{}, id).Error
}

func (r StatusPageRepositoryImpl) FindAll(ctx context.Context) ([]*entity.StatusPage, error) {
	var pages []*entity.StatusPage
	result := database.DB.WithContext(ctx).Preload("Monitors").Find(&pages)
	if result.Error != nil {
		return nil, result.Error
	}
	return pages, nil
}

func (r StatusPageRepositoryImpl) FindByUserID(ctx context.Context, userID uint) ([]*entity.StatusPage, error) {
	var pages []*entity.StatusPage
	result := database.DB.WithContext(ctx).Where("user_id = ?", userID).Preload("Monitors").Find(&pages)
	if result.Error != nil {
		return nil, result.Error
	}
	return pages, nil
}

func (r StatusPageRepositoryImpl) FindByID(ctx context.Context, id uint) (*entity.StatusPage, error) {
	var page entity.StatusPage
	result := database.DB.WithContext(ctx).Preload("Monitors").First(&page, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &page, nil
}

func (r StatusPageRepositoryImpl) FindBySlug(ctx context.Context, slug string) (*entity.StatusPage, error) {
	var page entity.StatusPage
	result := database.DB.WithContext(ctx).Where("slug = ?", slug).Preload("Monitors").First(&page)
	if result.Error != nil {
		return nil, result.Error
	}
	return &page, nil
}

func (r StatusPageRepositoryImpl) Update(ctx context.Context, statusPage *entity.StatusPage) error {
	// GORM updates association join table entries automatically on Save if the slice is loaded and modified.
	// But we can also replace associations explicitly to avoid duplicates.
	err := database.DB.WithContext(ctx).Model(statusPage).Association("Monitors").Replace(statusPage.Monitors)
	if err != nil {
		return err
	}
	return database.DB.WithContext(ctx).Save(statusPage).Error
}

func NewStatusPageRepositoryImpl() StatusPageRepository {
	return StatusPageRepositoryImpl{}
}

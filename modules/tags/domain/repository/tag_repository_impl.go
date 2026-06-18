package repository

import (
	"context"
	"ping-uptime/internal/pkg/database"
	"ping-uptime/modules/tags/domain/entity"
)

type TagRepositoryImpl struct{}

func (r TagRepositoryImpl) Create(ctx context.Context, tag *entity.Tag) error {
	return database.DB.WithContext(ctx).Create(tag).Error
}

func (r TagRepositoryImpl) FindAll(ctx context.Context) ([]*entity.Tag, error) {
	var tags []*entity.Tag
	err := database.DB.WithContext(ctx).Order("name ASC").Find(&tags).Error
	return tags, err
}

func (r TagRepositoryImpl) FindByID(ctx context.Context, id uint) (*entity.Tag, error) {
	var tag entity.Tag
	err := database.DB.WithContext(ctx).First(&tag, id).Error
	return &tag, err
}

func (r TagRepositoryImpl) FindByName(ctx context.Context, name string) (*entity.Tag, error) {
	var tag entity.Tag
	err := database.DB.WithContext(ctx).Where("name = ?", name).First(&tag).Error
	return &tag, err
}

func (r TagRepositoryImpl) Update(ctx context.Context, tag *entity.Tag) error {
	return database.DB.WithContext(ctx).Save(tag).Error
}

func (r TagRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return database.DB.WithContext(ctx).Delete(&entity.Tag{}, id).Error
}

func (r TagRepositoryImpl) AttachToMonitor(ctx context.Context, monitorID, tagID uint) error {
	mt := &entity.MonitorTag{
		MonitorID: monitorID,
		TagID:     tagID,
	}
	return database.DB.WithContext(ctx).Create(mt).Error
}

func (r TagRepositoryImpl) DetachFromMonitor(ctx context.Context, monitorID, tagID uint) error {
	return database.DB.WithContext(ctx).
		Where("monitor_id = ? AND tag_id = ?", monitorID, tagID).
		Delete(&entity.MonitorTag{}).Error
}

func (r TagRepositoryImpl) GetMonitorTags(ctx context.Context, monitorID uint) ([]*entity.MonitorTag, error) {
	var mts []*entity.MonitorTag
	err := database.DB.WithContext(ctx).
		Where("monitor_id = ?", monitorID).
		Find(&mts).Error
	return mts, err
}

func (r TagRepositoryImpl) GetTagsByMonitorID(ctx context.Context, monitorID uint) ([]*entity.Tag, error) {
	var tags []*entity.Tag
	err := database.DB.WithContext(ctx).
		Joins("JOIN monitor_tags ON monitor_tags.tag_id = tags.id").
		Where("monitor_tags.monitor_id = ?", monitorID).
		Find(&tags).Error
	return tags, err
}

func (r TagRepositoryImpl) DeleteMonitorTags(ctx context.Context, monitorID uint) error {
	return database.DB.WithContext(ctx).
		Where("monitor_id = ?", monitorID).
		Delete(&entity.MonitorTag{}).Error
}

func NewTagRepositoryImpl() TagRepository {
	return TagRepositoryImpl{}
}

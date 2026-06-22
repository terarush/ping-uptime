package repository

import (
	"context"
	"strconv"
	"ping-uptime/internal/pkg/database"
	"ping-uptime/modules/notifications/domain/entity"
)

type NotificationLogRepositoryImpl struct{}

func (r NotificationLogRepositoryImpl) Create(ctx context.Context, log *entity.NotificationLog) error {
	return database.DB.WithContext(ctx).Create(log).Error
}

func (r NotificationLogRepositoryImpl) FindAll(ctx context.Context, userID uint, isAdmin bool, filter map[string]string) ([]*entity.NotificationLog, error) {
	query := database.DB.WithContext(ctx).Model(&entity.NotificationLog{})

	if !isAdmin {
		query = query.Where("user_id = ?", userID)
	}

	if channelType, ok := filter["channel_type"]; ok && channelType != "" {
		query = query.Where("channel_type = ?", channelType)
	}
	if status, ok := filter["status"]; ok && status != "" {
		query = query.Where("status = ?", status)
	}
	if eventType, ok := filter["event_type"]; ok && eventType != "" {
		query = query.Where("event_type = ?", eventType)
	}
	if dateFrom, ok := filter["date_from"]; ok && dateFrom != "" {
		query = query.Where("sent_at >= ?", dateFrom)
	}
	if dateTo, ok := filter["date_to"]; ok && dateTo != "" {
		query = query.Where("sent_at <= ?", dateTo+"T23:59:59Z")
	}

	query = query.Order("sent_at DESC")

	limit := 100
	if l, ok := filter["limit"]; ok && l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}
	query = query.Limit(limit)

	if offsetStr, ok := filter["offset"]; ok && offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			query = query.Offset(offset)
		}
	}

	var logs []*entity.NotificationLog
	err := query.Find(&logs).Error
	return logs, err
}

func (r NotificationLogRepositoryImpl) FindByID(ctx context.Context, id uint) (*entity.NotificationLog, error) {
	var log entity.NotificationLog
	err := database.DB.WithContext(ctx).First(&log, id).Error
	return &log, err
}

func NewNotificationLogRepositoryImpl() NotificationLogRepository {
	return NotificationLogRepositoryImpl{}
}

var _ NotificationLogRepository = (*NotificationLogRepositoryImpl)(nil)

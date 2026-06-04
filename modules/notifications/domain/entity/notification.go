package entity

import (
	"time"
)

type NotificationChannel struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type" gorm:"type:varchar(50)"` // email, webhook, telegram, discord, slack
	Config    string    `json:"config" gorm:"type:text"`      // JSON config string
	Enabled   bool      `json:"enabled" gorm:"default:true"`
	UserID    uint      `json:"user_id" gorm:"index"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (*NotificationChannel) TableName() string {
	return "notification_channels"
}

func NewNotificationChannel(name, channelType, config string, enabled bool, userID uint) *NotificationChannel {
	now := time.Now()
	return &NotificationChannel{
		Name:      name,
		Type:      channelType,
		Config:    config,
		Enabled:   enabled,
		UserID:    userID,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

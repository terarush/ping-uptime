package entity

import "time"

type NotificationLog struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	UserID         uint      `gorm:"index" json:"user_id"`
	NotificationID uint      `json:"notification_id"`
	ChannelType    string    `gorm:"type:varchar(50)" json:"channel_type"`
	MonitorID      *uint     `json:"monitor_id"`
	IncidentID     *uint     `json:"incident_id"`
	EventType      string    `gorm:"type:varchar(20)" json:"event_type"`
	Status         string    `gorm:"type:varchar(20)" json:"status"`
	ErrorMessage   string    `gorm:"type:text" json:"error_message"`
	Recipient      string    `gorm:"type:varchar(255)" json:"recipient"`
	SentAt         time.Time `json:"sent_at"`
}

func (*NotificationLog) TableName() string {
	return "notification_logs"
}

func NewNotificationLog(userID, notificationID uint, channelType string, monitorID, incidentID *uint, eventType, status, errMsg, recipient string) *NotificationLog {
	return &NotificationLog{
		UserID:         userID,
		NotificationID: notificationID,
		ChannelType:    channelType,
		MonitorID:      monitorID,
		IncidentID:     incidentID,
		EventType:      eventType,
		Status:         status,
		ErrorMessage:   errMsg,
		Recipient:      recipient,
		SentAt:         time.Now(),
	}
}

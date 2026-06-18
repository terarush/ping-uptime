package entity

import "time"

type AuditLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `json:"user_id" gorm:"index"`
	Action     string    `json:"action" gorm:"type:varchar(20)"`      // create, update, delete
	EntityType string    `json:"entity_type" gorm:"type:varchar(50)"` // monitor, incident, status_page, notification_channel
	EntityID   uint      `json:"entity_id"`
	Details    string    `json:"details" gorm:"type:text"` // JSON metadata
	CreatedAt  time.Time `json:"created_at"`
}

func (*AuditLog) TableName() string {
	return "audit_logs"
}

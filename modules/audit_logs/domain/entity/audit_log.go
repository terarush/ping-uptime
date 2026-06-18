package entity

import "time"

type AuditLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `json:"user_id" gorm:"index"`
	Action     string    `json:"action" gorm:"type:varchar(20)"`      // create, update, delete
	EntityType string    `json:"entity_type" gorm:"type:varchar(50);index:idx_entity_type_id,priority:1"` // monitor, incident, status_page, notification_channel
	EntityID   uint      `json:"entity_id" gorm:"index:idx_entity_type_id,priority:2"`
	Details    string    `json:"details" gorm:"type:text"` // JSON metadata
	CreatedAt  time.Time `json:"created_at" gorm:"index:idx_audit_log_created_at"`
}

func (*AuditLog) TableName() string {
	return "audit_logs"
}

package entity

import "time"

type CheckRecord struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	MonitorID  uint      `gorm:"index" json:"monitor_id"`
	Success    bool      `json:"success"`
	Latency    int       `json:"latency"`
	StatusCode int       `json:"status_code"`
	CheckedAt  time.Time `gorm:"index" json:"checked_at"`
}

func (*CheckRecord) TableName() string {
	return "check_records"
}

func NewCheckRecord(monitorID uint, success bool, latency, statusCode int) *CheckRecord {
	return &CheckRecord{
		MonitorID:  monitorID,
		Success:    success,
		Latency:    latency,
		StatusCode: statusCode,
		CheckedAt:  time.Now(),
	}
}

package entity

import "time"

type DailyCheck struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	MonitorID    uint      `gorm:"index" json:"monitor_id"`
	CheckDate    time.Time `gorm:"index;type:date" json:"check_date"`
	Status       string    `gorm:"type:varchar(20)" json:"status"` // up, down, mixed
	ChecksTotal  int       `json:"checks_total"`
	ChecksFailed int       `json:"checks_failed"`
	CreatedAt    time.Time `json:"created_at"`
}

func (*DailyCheck) TableName() string {
	return "daily_checks"
}

func NewDailyCheck(monitorID uint, checkDate time.Time, status string, total, failed int) *DailyCheck {
	return &DailyCheck{
		MonitorID:    monitorID,
		CheckDate:    checkDate,
		Status:       status,
		ChecksTotal:  total,
		ChecksFailed: failed,
		CreatedAt:    time.Now(),
	}
}

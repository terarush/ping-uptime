package entity

import (
	"time"
)

type Monitor struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	Name          string     `json:"name"`
	URL           string     `json:"url"`
	Type          string     `json:"type" gorm:"type:varchar(50);default:'http'"`
	Interval      int        `json:"interval" gorm:"default:60"`
	Timeout       int        `json:"timeout" gorm:"default:10"`
	Status        string     `json:"status" gorm:"type:varchar(50);default:'active'"`        // active, paused
	UptimeStatus  string     `json:"uptime_status" gorm:"type:varchar(50);default:'unknown'"` // up, down, unknown
	LastCheckedAt *time.Time `json:"last_checked_at,omitempty"`
	UserID        uint       `json:"user_id" gorm:"index"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

func (*Monitor) TableName() string {
	return "monitors"
}

func NewMonitor(name, url, monitorType string, interval, timeout int, userID uint) *Monitor {
	now := time.Now()
	return &Monitor{
		Name:         name,
		URL:          url,
		Type:         monitorType,
		Interval:     interval,
		Timeout:      timeout,
		Status:       "active",
		UptimeStatus: "unknown",
		UserID:       userID,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

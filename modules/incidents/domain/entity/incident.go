package entity

import (
	"time"
)

type Incident struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	MonitorID    uint       `json:"monitor_id" gorm:"index"`
	UserID       uint       `json:"user_id" gorm:"index"`
	Status       string     `json:"status" gorm:"type:varchar(50);default:'active'"` // active, resolved
	ErrorMessage string     `json:"error_message"`
	Latency      int        `json:"latency"`
	CreatedAt    time.Time  `json:"created_at"`
	ResolvedAt   *time.Time `json:"resolved_at,omitempty"`
}

func (*Incident) TableName() string {
	return "incidents"
}

func NewIncident(monitorID, userID uint, status, errorMessage string, latency int) *Incident {
	now := time.Now()
	return &Incident{
		MonitorID:    monitorID,
		UserID:       userID,
		Status:       status,
		ErrorMessage: errorMessage,
		Latency:      latency,
		CreatedAt:    now,
	}
}

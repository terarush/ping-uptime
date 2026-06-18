package entity

import "time"

type Maintenance struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	StartAt     time.Time `json:"start_at"`
	EndAt       time.Time `json:"end_at"`
	Status      string    `json:"status" gorm:"type:varchar(20);default:'scheduled'"` // scheduled, ongoing, completed
	UserID      uint      `json:"user_id" gorm:"index"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (*Maintenance) TableName() string {
	return "maintenances"
}

type MaintenanceMonitor struct {
	MaintenanceID uint `gorm:"primaryKey"`
	MonitorID     uint `gorm:"primaryKey"`
}

func (*MaintenanceMonitor) TableName() string {
	return "maintenance_monitors"
}

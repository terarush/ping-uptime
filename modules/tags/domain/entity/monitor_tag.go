package entity

import "time"

type MonitorTag struct {
	MonitorID uint      `gorm:"primaryKey" json:"monitor_id"`
	TagID     uint      `gorm:"primaryKey" json:"tag_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (*MonitorTag) TableName() string { return "monitor_tags" }

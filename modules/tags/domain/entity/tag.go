package entity

import "time"

type Tag struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(50);uniqueIndex" json:"name"`
	Color     string    `gorm:"type:varchar(7);default:#6366f1" json:"color"`
	CreatedAt time.Time `json:"created_at"`
}

func (*Tag) TableName() string { return "tags" }

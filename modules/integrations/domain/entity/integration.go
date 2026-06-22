package entity

import "time"

type Integration struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(100)" json:"name"`
	Type      string    `gorm:"type:varchar(50)" json:"type"` // slack, discord, webhook, github, pagerduty
	Config    string    `gorm:"type:text" json:"config"`      // JSON config
	Enabled   bool      `gorm:"default:true" json:"enabled"`
	UserID    uint      `gorm:"index" json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (*Integration) TableName() string { return "integrations" }

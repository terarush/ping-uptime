package entity

import (
	"time"
)

type Setting struct {
	Key         string    `gorm:"primaryKey;type:varchar(100)" json:"key"`
	Value       string    `json:"value" gorm:"type:text"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (*Setting) TableName() string {
	return "settings"
}

func NewSetting(key, value, description string) *Setting {
	now := time.Now()
	return &Setting{
		Key:         key,
		Value:       value,
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

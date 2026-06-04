package entity

import (
	"ping-uptime/modules/monitors/domain/entity"
	"time"
)

type StatusPage struct {
	ID          uint              `gorm:"primaryKey" json:"id"`
	Name        string            `json:"name"`
	Slug        string            `json:"slug" gorm:"uniqueIndex;type:varchar(100)"`
	Description string            `json:"description"`
	UserID      uint              `json:"user_id" gorm:"index"`
	Monitors    []*entity.Monitor `json:"monitors,omitempty" gorm:"many2many:status_page_monitors;"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

func (*StatusPage) TableName() string {
	return "status_pages"
}

func NewStatusPage(name, slug, description string, userID uint) *StatusPage {
	now := time.Now()
	return &StatusPage{
		Name:        name,
		Slug:        slug,
		Description: description,
		UserID:      userID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

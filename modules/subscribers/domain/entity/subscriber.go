package entity

import "time"

type Subscriber struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Email        string    `json:"email" gorm:"type:varchar(255);uniqueIndex:idx_status_email"`
	StatusPageID uint      `json:"status_page_id" gorm:"index;uniqueIndex:idx_status_email"`
	Verified     bool      `json:"verified" gorm:"default:false"`
	Token        string    `json:"-" gorm:"type:varchar(64);uniqueIndex"`
	CreatedAt    time.Time `json:"created_at"`
}

func (*Subscriber) TableName() string {
	return "subscribers"
}

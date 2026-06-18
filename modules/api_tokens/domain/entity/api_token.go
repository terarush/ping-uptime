package entity

import "time"

type ApiToken struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	UserID      uint       `gorm:"index;not null" json:"user_id"`
	Name        string     `gorm:"type:varchar(100)" json:"name"`
	TokenPrefix string     `gorm:"type:varchar(8)" json:"token_prefix"`
	TokenHash   string     `gorm:"type:varchar(64);uniqueIndex" json:"-"`
	LastUsedAt  *time.Time `json:"last_used_at"`
	ExpiresAt   *time.Time `json:"expires_at"`
	IsRevoked   bool       `gorm:"default:false" json:"is_revoked"`
	CreatedAt   time.Time  `json:"created_at"`
}

func (*ApiToken) TableName() string { return "api_tokens" }

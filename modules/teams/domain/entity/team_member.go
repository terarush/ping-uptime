package entity

import "time"

type TeamMember struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	TeamID    uint      `gorm:"index;not null" json:"team_id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	Role      string    `gorm:"type:varchar(20);default:member" json:"role"`
	InvitedBy uint      `json:"invited_by"`
	Status    string    `gorm:"type:varchar(20);default:pending" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (*TeamMember) TableName() string { return "team_members" }

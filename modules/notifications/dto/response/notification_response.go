package response

import (
	"ping-uptime/modules/notifications/domain/entity"
	"time"
)

type ChannelResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Config    string    `json:"config"`
	Enabled   bool      `json:"enabled"`
	UserID    uint      `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func FromEntity(c *entity.NotificationChannel) *ChannelResponse {
	return &ChannelResponse{
		ID:        c.ID,
		Name:      c.Name,
		Type:      c.Type,
		Config:    c.Config,
		Enabled:   c.Enabled,
		UserID:    c.UserID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func FromEntities(channels []*entity.NotificationChannel) []*ChannelResponse {
	responses := make([]*ChannelResponse, len(channels))
	for i, c := range channels {
		responses[i] = FromEntity(c)
	}
	return responses
}

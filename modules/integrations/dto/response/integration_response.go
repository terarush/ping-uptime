package response

import (
	"ping-uptime/modules/integrations/domain/entity"
	"time"
)

type IntegrationResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Config    string    `json:"config"`
	Enabled   bool      `json:"enabled"`
	UserID    uint      `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func FromEntity(i *entity.Integration) *IntegrationResponse {
	return &IntegrationResponse{
		ID:        i.ID,
		Name:      i.Name,
		Type:      i.Type,
		Config:    i.Config,
		Enabled:   i.Enabled,
		UserID:    i.UserID,
		CreatedAt: i.CreatedAt,
		UpdatedAt: i.UpdatedAt,
	}
}

func FromEntities(integrations []*entity.Integration) []*IntegrationResponse {
	responses := make([]*IntegrationResponse, len(integrations))
	for idx, i := range integrations {
		responses[idx] = FromEntity(i)
	}
	return responses
}

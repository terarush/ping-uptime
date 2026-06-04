package response

import (
	"ping-uptime/modules/settings/domain/entity"
	"time"
)

type SettingResponse struct {
	Key         string    `json:"key"`
	Value       string    `json:"value"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func FromEntity(s *entity.Setting) *SettingResponse {
	return &SettingResponse{
		Key:         s.Key,
		Value:       s.Value,
		Description: s.Description,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}

func FromEntities(settings []*entity.Setting) []*SettingResponse {
	responses := make([]*SettingResponse, len(settings))
	for i, s := range settings {
		responses[i] = FromEntity(s)
	}
	return responses
}

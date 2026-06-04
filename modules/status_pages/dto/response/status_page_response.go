package response

import (
	monitorResponse "ping-uptime/modules/monitors/dto/response"
	"ping-uptime/modules/status_pages/domain/entity"
	"time"
)

type StatusPageResponse struct {
	ID          uint                              `json:"id"`
	Name        string                            `json:"name"`
	Slug        string                            `json:"slug"`
	Description string                            `json:"description"`
	UserID      uint                              `json:"user_id"`
	Monitors    []*monitorResponse.MonitorResponse `json:"monitors,omitempty"`
	CreatedAt   time.Time                         `json:"created_at"`
	UpdatedAt   time.Time                         `json:"updated_at"`
}

func FromEntity(s *entity.StatusPage) *StatusPageResponse {
	var monitors []*monitorResponse.MonitorResponse
	if len(s.Monitors) > 0 {
		monitors = monitorResponse.FromEntities(s.Monitors)
	} else {
		// Ensure empty array instead of null in JSON response
		monitors = []*monitorResponse.MonitorResponse{}
	}

	return &StatusPageResponse{
		ID:          s.ID,
		Name:        s.Name,
		Slug:        s.Slug,
		Description: s.Description,
		UserID:      s.UserID,
		Monitors:    monitors,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}

func FromEntities(pages []*entity.StatusPage) []*StatusPageResponse {
	responses := make([]*StatusPageResponse, len(pages))
	for i, s := range pages {
		responses[i] = FromEntity(s)
	}
	return responses
}

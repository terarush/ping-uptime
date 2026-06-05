package response

import (
	"ping-uptime/modules/incidents/domain/entity"
	"time"
)

type IncidentResponse struct {
	ID           uint       `json:"id"`
	MonitorID    uint       `json:"monitor_id"`
	UserID       uint       `json:"user_id"`
	Status       string     `json:"status"`
	ErrorMessage string     `json:"error_message"`
	Latency      int        `json:"latency"`
	CreatedAt    time.Time  `json:"created_at"`
	ResolvedAt   *time.Time `json:"resolved_at,omitempty"`
}

func FromEntity(i *entity.Incident) *IncidentResponse {
	var resolved *time.Time
	if i.ResolvedAt != nil {
		localTime := i.ResolvedAt.Local()
		resolved = &localTime
	}
	return &IncidentResponse{
		ID:           i.ID,
		MonitorID:    i.MonitorID,
		UserID:       i.UserID,
		Status:       i.Status,
		ErrorMessage: i.ErrorMessage,
		Latency:      i.Latency,
		CreatedAt:    i.CreatedAt.Local(),
		ResolvedAt:   resolved,
	}
}

func FromEntities(incidents []*entity.Incident) []*IncidentResponse {
	responses := make([]*IncidentResponse, len(incidents))
	for i, inc := range incidents {
		responses[i] = FromEntity(inc)
	}
	return responses
}

package response

import (
	"ping-uptime/modules/maintenances/domain/entity"
	"time"
)

type MaintenanceResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	StartAt     time.Time `json:"start_at"`
	EndAt       time.Time `json:"end_at"`
	Status      string    `json:"status"`
	UserID      uint      `json:"user_id"`
	MonitorIDs  []uint    `json:"monitor_ids,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func FromEntity(m *entity.Maintenance) *MaintenanceResponse {
	return &MaintenanceResponse{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		StartAt:     m.StartAt.Local(),
		EndAt:       m.EndAt.Local(),
		Status:      m.Status,
		UserID:      m.UserID,
		CreatedAt:   m.CreatedAt.Local(),
		UpdatedAt:   m.UpdatedAt.Local(),
	}
}

func FromEntities(items []*entity.Maintenance) []*MaintenanceResponse {
	out := make([]*MaintenanceResponse, len(items))
	for i, m := range items {
		out[i] = FromEntity(m)
	}
	return out
}

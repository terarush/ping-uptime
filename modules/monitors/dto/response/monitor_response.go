package response

import (
	"ping-uptime/modules/monitors/domain/entity"
	"time"
)

type MonitorResponse struct {
	ID             uint       `json:"id"`
	Name           string     `json:"name"`
	URL            string     `json:"url"`
	Type           string     `json:"type"`
	Interval       int        `json:"interval"`
	Timeout        int        `json:"timeout"`
	Status         string     `json:"status"`
	UptimeStatus   string     `json:"uptime_status"`
	LastCheckedAt  *time.Time `json:"last_checked_at,omitempty"`
	LastLatency    int        `json:"last_latency"`
	CheckSSL       bool       `json:"check_ssl"`
	SSLExpiresAt   *time.Time `json:"ssl_expires_at,omitempty"`
	HeartbeatToken string     `json:"heartbeat_token,omitempty"`
	UserID         uint       `json:"user_id"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

func FromEntity(m *entity.Monitor) *MonitorResponse {
	var lastChecked *time.Time
	if m.LastCheckedAt != nil {
		localTime := m.LastCheckedAt.Local()
		lastChecked = &localTime
	}
	var sslExpires *time.Time
	if m.SSLExpiresAt != nil {
		localTime := m.SSLExpiresAt.Local()
		sslExpires = &localTime
	}
	return &MonitorResponse{
		ID:             m.ID,
		Name:           m.Name,
		URL:            m.URL,
		Type:           m.Type,
		Interval:       m.Interval,
		Timeout:        m.Timeout,
		Status:         m.Status,
		UptimeStatus:   m.UptimeStatus,
		LastCheckedAt:  lastChecked,
		LastLatency:    m.LastLatency,
		CheckSSL:       m.CheckSSL,
		SSLExpiresAt:   sslExpires,
		HeartbeatToken: m.HeartbeatToken,
		UserID:         m.UserID,
		CreatedAt:      m.CreatedAt.Local(),
		UpdatedAt:      m.UpdatedAt.Local(),
	}
}

func FromEntities(monitors []*entity.Monitor) []*MonitorResponse {
	responses := make([]*MonitorResponse, len(monitors))
	for i, m := range monitors {
		responses[i] = FromEntity(m)
	}
	return responses
}

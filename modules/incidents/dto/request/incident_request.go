package request

type CreateIncidentRequest struct {
	MonitorID    uint   `json:"monitor_id" validate:"required"`
	Status       string `json:"status" validate:"required"`
	ErrorMessage string `json:"error_message"`
	Latency      int    `json:"latency"`
}

type UpdateIncidentRequest struct {
	Status       string `json:"status" validate:"required"`
	ErrorMessage string `json:"error_message"`
	Latency      int    `json:"latency"`
}

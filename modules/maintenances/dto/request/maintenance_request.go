package request

type CreateMaintenanceRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	StartAt     string `json:"start_at" validate:"required"`
	EndAt       string `json:"end_at" validate:"required"`
	MonitorIDs  []uint `json:"monitor_ids"`
}

type UpdateMaintenanceRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	StartAt     string `json:"start_at" validate:"required"`
	EndAt       string `json:"end_at" validate:"required"`
	MonitorIDs  []uint `json:"monitor_ids"`
}

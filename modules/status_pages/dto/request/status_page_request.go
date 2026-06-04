package request

type CreateStatusPageRequest struct {
	Name        string `json:"name" validate:"required"`
	Slug        string `json:"slug" validate:"required"`
	Description string `json:"description"`
	MonitorIDs  []uint `json:"monitor_ids"`
}

type UpdateStatusPageRequest struct {
	Name        string `json:"name" validate:"required"`
	Slug        string `json:"slug" validate:"required"`
	Description string `json:"description"`
	MonitorIDs  []uint `json:"monitor_ids"`
}

package request

type CreateMonitorRequest struct {
	Name     string `json:"name" validate:"required"`
	URL      string `json:"url" validate:"required,url"`
	Type     string `json:"type" validate:"required"`
	Interval int    `json:"interval" validate:"required,min=1"`
	Timeout  int    `json:"timeout" validate:"required,min=5"`
}

type UpdateMonitorRequest struct {
	Name     string `json:"name" validate:"required"`
	URL      string `json:"url" validate:"required,url"`
	Type     string `json:"type" validate:"required"`
	Interval int    `json:"interval" validate:"required,min=1"`
	Timeout  int    `json:"timeout" validate:"required,min=5"`
	Status   string `json:"status" validate:"required"`
}

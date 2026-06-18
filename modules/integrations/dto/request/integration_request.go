package request

type CreateIntegrationRequest struct {
	Name   string `json:"name" validate:"required,max=100"`
	Type   string `json:"type" validate:"required,oneof=slack discord webhook github pagerduty"`
	Config string `json:"config" validate:"required"`
}

type UpdateIntegrationRequest struct {
	Name    string `json:"name" validate:"required,max=100"`
	Type    string `json:"type" validate:"required,oneof=slack discord webhook github pagerduty"`
	Config  string `json:"config" validate:"required"`
	Enabled bool   `json:"enabled"`
}

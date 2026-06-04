package request

type CreateChannelRequest struct {
	Name    string `json:"name" validate:"required"`
	Type    string `json:"type" validate:"required"`
	Config  string `json:"config" validate:"required"` // JSON config
	Enabled bool   `json:"enabled"`
}

type UpdateChannelRequest struct {
	Name    string `json:"name" validate:"required"`
	Type    string `json:"type" validate:"required"`
	Config  string `json:"config" validate:"required"`
	Enabled bool   `json:"enabled"`
}

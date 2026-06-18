package request

type CreateTokenRequest struct {
	Name      string `json:"name" validate:"required,max=100"`
	ExpiresAt string `json:"expires_at"` // optional RFC3339
}

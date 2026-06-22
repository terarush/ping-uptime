package response

import (
	"ping-uptime/modules/api_tokens/domain/entity"
	"time"
)

type TokenResponse struct {
	ID          uint       `json:"id"`
	Name        string     `json:"name"`
	TokenPrefix string     `json:"token_prefix"`
	LastUsedAt  *time.Time `json:"last_used_at"`
	ExpiresAt   *time.Time `json:"expires_at"`
	IsRevoked   bool       `json:"is_revoked"`
	CreatedAt   time.Time  `json:"created_at"`
}

type CreateTokenResponse struct {
	TokenResponse
	RawToken string `json:"raw_token"` // only returned once on create
}

func FromEntity(t *entity.ApiToken) TokenResponse {
	return TokenResponse{
		ID:          t.ID,
		Name:        t.Name,
		TokenPrefix: t.TokenPrefix,
		LastUsedAt:  t.LastUsedAt,
		ExpiresAt:   t.ExpiresAt,
		IsRevoked:   t.IsRevoked,
		CreatedAt:   t.CreatedAt,
	}
}

func FromEntities(tokens []*entity.ApiToken) []TokenResponse {
	result := make([]TokenResponse, len(tokens))
	for i, t := range tokens {
		result[i] = FromEntity(t)
	}
	return result
}

func ToCreateResponse(t *entity.ApiToken, rawToken string) CreateTokenResponse {
	return CreateTokenResponse{
		TokenResponse: FromEntity(t),
		RawToken:      rawToken,
	}
}

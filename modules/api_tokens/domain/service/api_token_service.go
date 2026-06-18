package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"ping-uptime/modules/api_tokens/domain/entity"
	"ping-uptime/modules/api_tokens/domain/repository"
	"time"
)

var (
	ErrTokenNotFound = errors.New("api token not found")
	ErrTokenRevoked  = errors.New("api token is revoked")
	ErrTokenExpired  = errors.New("api token has expired")
	ErrNotAuthorized = errors.New("not authorized to manage this token")
)

type ApiTokenService struct {
	tokenRepo repository.ApiTokenRepository
}

func NewApiTokenService(tokenRepo repository.ApiTokenRepository) *ApiTokenService {
	return &ApiTokenService{
		tokenRepo: tokenRepo,
	}
}

// GenerateToken creates a new API token, stores its hash, and returns the raw token string.
// The raw token is only returned once and cannot be retrieved later.
func (s *ApiTokenService) GenerateToken(ctx context.Context, userID uint, name string, expiresAt *time.Time) (*entity.ApiToken, string, error) {
	rawBytes := make([]byte, 40)
	if _, err := rand.Read(rawBytes); err != nil {
		return nil, "", fmt.Errorf("failed to generate token bytes: %w", err)
	}

	rawToken := "pu_" + hex.EncodeToString(rawBytes)

	hash := sha256.Sum256([]byte(rawToken))
	tokenHash := hex.EncodeToString(hash[:])

	prefix := rawToken[:10] // "pu_" + first 6 hex chars

	token := &entity.ApiToken{
		UserID:      userID,
		Name:        name,
		TokenPrefix: prefix,
		TokenHash:   tokenHash,
		ExpiresAt:   expiresAt,
	}

	if err := s.tokenRepo.Create(ctx, token); err != nil {
		return nil, "", fmt.Errorf("failed to store token: %w", err)
	}

	return token, rawToken, nil
}

// GetUserTokens returns all tokens for a given user (without hash).
func (s *ApiTokenService) GetUserTokens(ctx context.Context, userID uint) ([]*entity.ApiToken, error) {
	return s.tokenRepo.FindByUserID(ctx, userID)
}

// GetTokenByHash looks up a token by its SHA256 hash, checking revocation and expiry.
func (s *ApiTokenService) GetTokenByHash(ctx context.Context, hash string) (*entity.ApiToken, error) {
	token, err := s.tokenRepo.FindByHash(ctx, hash)
	if err != nil {
		return nil, ErrTokenNotFound
	}

	if token.IsRevoked {
		return nil, ErrTokenRevoked
	}

	if token.ExpiresAt != nil && token.ExpiresAt.Before(time.Now()) {
		return nil, ErrTokenExpired
	}

	return token, nil
}

// RevokeToken marks a token as revoked, only if caller is the owner or an admin.
func (s *ApiTokenService) RevokeToken(ctx context.Context, id uint, callerID uint, callerRole string) error {
	token, err := s.tokenRepo.FindByID(ctx, id)
	if err != nil {
		return ErrTokenNotFound
	}

	if callerRole != "admin" && token.UserID != callerID {
		return ErrNotAuthorized
	}

	token.IsRevoked = true
	return s.tokenRepo.Update(ctx, token)
}

// DeleteExpired removes all expired tokens from the database.
func (s *ApiTokenService) DeleteExpired(ctx context.Context) error {
	return s.tokenRepo.DeleteExpired(ctx)
}

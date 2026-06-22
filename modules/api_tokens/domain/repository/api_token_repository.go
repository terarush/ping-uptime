package repository

import (
	"context"
	"ping-uptime/modules/api_tokens/domain/entity"
)

type ApiTokenRepository interface {
	Create(ctx context.Context, token *entity.ApiToken) error
	Update(ctx context.Context, token *entity.ApiToken) error
	FindByHash(ctx context.Context, hash string) (*entity.ApiToken, error)
	FindByUserID(ctx context.Context, userID uint) ([]*entity.ApiToken, error)
	FindByID(ctx context.Context, id uint) (*entity.ApiToken, error)
	Delete(ctx context.Context, id uint) error
	DeleteExpired(ctx context.Context) error
}

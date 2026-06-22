package repository

import (
	"context"
	"errors"
	"ping-uptime/internal/pkg/database"
	"ping-uptime/modules/api_tokens/domain/entity"
	"time"
)

type ApiTokenRepositoryImpl struct{}

func (r ApiTokenRepositoryImpl) Create(ctx context.Context, token *entity.ApiToken) error {
	return database.DB.WithContext(ctx).Create(token).Error
}

func (r ApiTokenRepositoryImpl) Update(ctx context.Context, token *entity.ApiToken) error {
	return database.DB.WithContext(ctx).Save(token).Error
}

func (r ApiTokenRepositoryImpl) FindByHash(ctx context.Context, hash string) (*entity.ApiToken, error) {
	var token entity.ApiToken
	result := database.DB.WithContext(ctx).Where("token_hash = ?", hash).First(&token)
	if result.Error != nil {
		return nil, result.Error
	}
	return &token, nil
}

func (r ApiTokenRepositoryImpl) FindByUserID(ctx context.Context, userID uint) ([]*entity.ApiToken, error) {
	var tokens []*entity.ApiToken
	result := database.DB.WithContext(ctx).Where("user_id = ?", userID).Order("created_at desc").Find(&tokens)
	if result.Error != nil {
		return nil, result.Error
	}
	return tokens, nil
}

func (r ApiTokenRepositoryImpl) FindByID(ctx context.Context, id uint) (*entity.ApiToken, error) {
	var token entity.ApiToken
	result := database.DB.WithContext(ctx).First(&token, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &token, nil
}

func (r ApiTokenRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return database.DB.WithContext(ctx).Delete(&entity.ApiToken{}, id).Error
}

func (r ApiTokenRepositoryImpl) DeleteExpired(ctx context.Context) error {
	return database.DB.WithContext(ctx).
		Where("expires_at IS NOT NULL AND expires_at < ?", time.Now()).
		Delete(&entity.ApiToken{}).Error
}

var _ ApiTokenRepository = (*ApiTokenRepositoryImpl)(nil)

func NewApiTokenRepositoryImpl() ApiTokenRepository {
	return ApiTokenRepositoryImpl{}
}

// Ensure ErrRecordNotFound is accessible for callers.
var ErrRecordNotFound = errors.New("record not found")

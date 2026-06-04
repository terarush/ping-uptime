package repository

import (
	"context"
	"ping-uptime/modules/settings/domain/entity"
)

type SettingRepository interface {
	CreateOrUpdate(ctx context.Context, setting *entity.Setting) error
	Delete(ctx context.Context, key string) error
	FindAll(ctx context.Context) ([]*entity.Setting, error)
	FindByKey(ctx context.Context, key string) (*entity.Setting, error)
}

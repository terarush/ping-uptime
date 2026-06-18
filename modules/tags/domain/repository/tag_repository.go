package repository

import (
	"context"
	"ping-uptime/modules/tags/domain/entity"
)

type TagRepository interface {
	Create(ctx context.Context, tag *entity.Tag) error
	FindAll(ctx context.Context) ([]*entity.Tag, error)
	FindByID(ctx context.Context, id uint) (*entity.Tag, error)
	FindByName(ctx context.Context, name string) (*entity.Tag, error)
	Update(ctx context.Context, tag *entity.Tag) error
	Delete(ctx context.Context, id uint) error
	AttachToMonitor(ctx context.Context, monitorID, tagID uint) error
	DetachFromMonitor(ctx context.Context, monitorID, tagID uint) error
	GetMonitorTags(ctx context.Context, monitorID uint) ([]*entity.MonitorTag, error)
	GetTagsByMonitorID(ctx context.Context, monitorID uint) ([]*entity.Tag, error)
	DeleteMonitorTags(ctx context.Context, monitorID uint) error
}

package repository

import (
	"context"
	"ping-uptime/internal/pkg/database"
	"ping-uptime/modules/ssl_monitors/domain/entity"
)

type SSLRepositoryImpl struct{}

func (r SSLRepositoryImpl) Create(ctx context.Context, cert *entity.SSLCert) error {
	return database.DB.WithContext(ctx).Create(cert).Error
}

func (r SSLRepositoryImpl) Update(ctx context.Context, cert *entity.SSLCert) error {
	return database.DB.WithContext(ctx).Save(cert).Error
}

func (r SSLRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return database.DB.WithContext(ctx).Delete(&entity.SSLCert{}, id).Error
}

func (r SSLRepositoryImpl) FindByID(ctx context.Context, id uint) (*entity.SSLCert, error) {
	var cert entity.SSLCert
	result := database.DB.WithContext(ctx).First(&cert, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &cert, nil
}

func (r SSLRepositoryImpl) FindByMonitorID(ctx context.Context, monitorID uint) (*entity.SSLCert, error) {
	var cert entity.SSLCert
	result := database.DB.WithContext(ctx).Where("monitor_id = ?", monitorID).First(&cert)
	if result.Error != nil {
		return nil, result.Error
	}
	return &cert, nil
}

func (r SSLRepositoryImpl) FindAll(ctx context.Context) ([]*entity.SSLCert, error) {
	var certs []*entity.SSLCert
	result := database.DB.WithContext(ctx).Order("checked_at DESC").Find(&certs)
	if result.Error != nil {
		return nil, result.Error
	}
	return certs, nil
}

func (r SSLRepositoryImpl) FindExpiring(ctx context.Context, days int) ([]*entity.SSLCert, error) {
	var certs []*entity.SSLCert
	result := database.DB.WithContext(ctx).
		Where("days_remaining <= ? AND days_remaining >= 0", days).
		Order("days_remaining ASC").
		Find(&certs)
	if result.Error != nil {
		return nil, result.Error
	}
	return certs, nil
}

func NewSSLRepositoryImpl() SSLRepository {
	return SSLRepositoryImpl{}
}

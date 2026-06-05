package repository

import (
	"context"
	"errors"
	"fmt"
	"ping-uptime/internal/pkg/database"
	"ping-uptime/modules/monitors/domain/entity"
)

var (
	ERR_RECORD_NOT_FOUND = errors.New("record not found")
)

type MonitorRepositoryImpl struct{}

func (r MonitorRepositoryImpl) Create(ctx context.Context, monitor *entity.Monitor) error {
	return database.DB.WithContext(ctx).Create(monitor).Error
}

func (r MonitorRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return database.DB.WithContext(ctx).Delete(&entity.Monitor{}, id).Error
}

func (r MonitorRepositoryImpl) FindAll(ctx context.Context) ([]*entity.Monitor, error) {
	var monitors []*entity.Monitor
	result := database.DB.WithContext(ctx).Find(&monitors)
	if result.Error != nil {
		return nil, result.Error
	}
	return monitors, nil
}

func (r MonitorRepositoryImpl) FindByUserID(ctx context.Context, userID uint) ([]*entity.Monitor, error) {
	var monitors []*entity.Monitor
	result := database.DB.WithContext(ctx).Where("user_id = ?", userID).Find(&monitors)
	if result.Error != nil {
		return nil, result.Error
	}
	return monitors, nil
}

func (r MonitorRepositoryImpl) FindByID(ctx context.Context, id uint) (*entity.Monitor, error) {
	var monitor entity.Monitor
	result := database.DB.WithContext(ctx).First(&monitor, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &monitor, nil
}

func (r MonitorRepositoryImpl) Update(ctx context.Context, monitor *entity.Monitor) error {
	return database.DB.WithContext(ctx).Save(monitor).Error
}

func (r MonitorRepositoryImpl) CreateCheckRecord(ctx context.Context, rec *entity.CheckRecord) error {
	return database.DB.WithContext(ctx).Create(rec).Error
}

func (r MonitorRepositoryImpl) GetDailyChart(ctx context.Context, monitorID uint, days int) ([]entity.DailyChartPoint, error) {
	var points []entity.DailyChartPoint
	err := database.DB.WithContext(ctx).
		Table("check_records").
		Select("DATE(checked_at) as date, SUM(CASE WHEN success THEN 0 ELSE 1 END) as failed_count, COUNT(*) as total_count").
		Where("monitor_id = ? AND checked_at >= ?", monitorID, fmt.Sprintf("now() - interval '%d days'", days)).
		Group("DATE(checked_at)").
		Order("date ASC").
		Scan(&points).Error
	if err != nil {
		return nil, err
	}
	out := make([]entity.DailyChartPoint, 0, len(points))
	for _, p := range points {
		failed := p.Failed
		if failed < 0 {
			failed = 0
		}
		total := p.Total
		if total <= 0 {
			total = 1
		}
		pct := 100.0 - (float64(failed) / float64(total) * 100.0)
		status := "up"
		if failed > 0 {
			status = "down"
		}
		out = append(out, entity.DailyChartPoint{
			Date:      p.Date,
			UptimePct: pct,
			Status:    status,
			Failed:    failed,
			Total:     total,
		})
	}
	return out, nil
}

func NewMonitorRepositoryImpl() MonitorRepository {
	return MonitorRepositoryImpl{}
}

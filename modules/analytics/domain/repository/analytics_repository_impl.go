package repository

import (
	"context"
	"ping-uptime/internal/pkg/database"
	"ping-uptime/modules/analytics/domain/entity"
	"time"
)

type AnalyticsRepositoryImpl struct{}

func (r AnalyticsRepositoryImpl) GetChartData(ctx context.Context, monitorID uint, window string) ([]entity.ChartDataPoint, error) {
	var points []entity.ChartDataPoint

	startTime := getWindowStartTime(window)
	if startTime.IsZero() {
		startTime = time.Now().AddDate(0, -1, 0)
	}

	err := database.DB.WithContext(ctx).
		Table("check_records").
		Select("DATE(checked_at) as date, SUM(CASE WHEN success THEN 0 ELSE 1 END) as failed_count, COUNT(*) as total_count").
		Where("monitor_id = ? AND checked_at >= ?", monitorID, startTime).
		Group("DATE(checked_at)").
		Order("date ASC").
		Scan(&points).Error

	if err != nil {
		return nil, err
	}

	out := make([]entity.ChartDataPoint, 0, len(points))
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
		out = append(out, entity.ChartDataPoint{
			Date:   p.Date,
			Status: status,
			Uptime: pct,
			Failed: failed,
			Total:  total,
		})
	}
	return out, nil
}

func (r AnalyticsRepositoryImpl) GetMonitorStats(ctx context.Context, userID uint, window string) ([]entity.MonitorStats, error) {
	var monitors []struct {
		ID       uint
		Name     string
		URL      string
		Status   string
	}

	database.DB.WithContext(ctx).Model(nil).Select("id, name, url, status").Scan(&monitors)

	stats := make([]entity.MonitorStats, 0, len(monitors))
	for _, m := range monitors {
		points, err := r.GetChartData(ctx, m.ID, window)
		if err != nil {
			continue
		}

		var totalChecks, failedChecks int
		var uptimePct float64
		status := "operational"

		if len(points) > 0 {
			for _, p := range points {
				totalChecks += p.Total
				failedChecks += p.Failed
			}
			uptimePct = 100.0 - (float64(failedChecks) / float64(totalChecks) * 100.0)
			if failedChecks > 0 {
				status = "degraded"
			}
			if len(points) > 1 && points[len(points)-1].Status == "down" {
				status = "outage"
			}
		} else {
			uptimePct = 100.0
		}

		stats = append(stats, entity.MonitorStats{
			MonitorID:    m.ID,
			MonitorName:  m.Name,
			MonitorURL:   m.URL,
			Window:       window,
			UptimePct:    uptimePct,
			TotalChecks:  totalChecks,
			FailedChecks: failedChecks,
			Points:       points,
			Status:       status,
		})
	}
	return stats, nil
}

func getWindowStartTime(window string) time.Time {
	now := time.Now()
	switch window {
	case "1h":
		return now.Add(-time.Hour)
	case "1d", "1w":
		return now.AddDate(0, 0, -7)
	case "1m":
		return now.AddDate(0, -1, 0)
	case "1y":
		return now.AddDate(-1, 0, 0)
	case "all":
		return time.Time{}
	default:
		return now.AddDate(0, -1, 0)
	}
}

func NewAnalyticsRepositoryImpl() AnalyticsRepository {
	return AnalyticsRepositoryImpl{}
}

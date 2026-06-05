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

	var selectExpr, groupExpr string
	switch window {
	case "1h":
		selectExpr = "strftime('%Y-%m-%d %H:%M:00', checked_at) as date, SUM(CASE WHEN success THEN 0 ELSE 1 END) as failed, COUNT(*) as total, AVG(latency) as latency"
		groupExpr = "strftime('%Y-%m-%d %H:%M:00', checked_at)"
	case "1d":
		selectExpr = "strftime('%Y-%m-%d %H:00:00', checked_at) as date, SUM(CASE WHEN success THEN 0 ELSE 1 END) as failed, COUNT(*) as total, AVG(latency) as latency"
		groupExpr = "strftime('%Y-%m-%d %H:00:00', checked_at)"
	case "1w", "1m":
		selectExpr = "DATE(checked_at) as date, SUM(CASE WHEN success THEN 0 ELSE 1 END) as failed, COUNT(*) as total, AVG(latency) as latency"
		groupExpr = "DATE(checked_at)"
	case "1y", "all":
		selectExpr = "strftime('%Y-W%W', checked_at) as date, SUM(CASE WHEN success THEN 0 ELSE 1 END) as failed, COUNT(*) as total, AVG(latency) as latency"
		groupExpr = "strftime('%Y-W%W', checked_at)"
	default:
		selectExpr = "DATE(checked_at) as date, SUM(CASE WHEN success THEN 0 ELSE 1 END) as failed, COUNT(*) as total, AVG(latency) as latency"
		groupExpr = "DATE(checked_at)"
	}

	err := database.DB.WithContext(ctx).
		Table("check_records").
		Select(selectExpr).
		Where("monitor_id = ? AND checked_at >= ?", monitorID, startTime).
		Group(groupExpr).
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
			Date:    p.Date,
			Status:  status,
			Uptime:  pct,
			Failed:  failed,
			Total:   total,
			Latency: p.Latency,
		})
	}
	return out, nil
}

func (r AnalyticsRepositoryImpl) GetMonitorStats(ctx context.Context, userID uint, window string) ([]entity.MonitorStats, error) {
	var monitors []struct {
		ID     uint
		Name   string
		URL    string
		Status string
	}

	db := database.DB.WithContext(ctx).Table("monitors").Select("id, name, url, status")
	if userID != 0 {
		db = db.Where("user_id = ?", userID)
	}
	db.Scan(&monitors)

	stats := make([]entity.MonitorStats, 0, len(monitors))
	for _, m := range monitors {
		points, err := r.GetChartData(ctx, m.ID, window)
		if err != nil {
			continue
		}

		var totalChecks, failedChecks int
		var uptimePct float64
		status := "operational"

		var latStats struct {
			Avg float64
			Min float64
			Max float64
		}

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

			// Query latency stats for successful checks
			database.DB.WithContext(ctx).
				Table("check_records").
				Select("COALESCE(AVG(latency), 0) as avg, COALESCE(MIN(latency), 0) as min, COALESCE(MAX(latency), 0) as max").
				Where("monitor_id = ? AND checked_at >= ? AND success = ? AND latency > 0", m.ID, getWindowStartTime(window), true).
				Scan(&latStats)
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
			AvgLatency:   latStats.Avg,
			MinLatency:   latStats.Min,
			MaxLatency:   latStats.Max,
		})
	}
	return stats, nil
}

func getWindowStartTime(window string) time.Time {
	now := time.Now()
	switch window {
	case "1h":
		return now.Add(-time.Hour)
	case "1d":
		return now.AddDate(0, 0, -1)
	case "1w":
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

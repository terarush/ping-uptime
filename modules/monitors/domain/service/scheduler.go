package service

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"ping-uptime/internal/pkg/bus"
	"ping-uptime/internal/pkg/database"
	incidentEntity "ping-uptime/modules/incidents/domain/entity"
	"ping-uptime/modules/monitors/domain/entity"
	"strings"
	"time"
)

// StartScheduler starts the background cron loop to process pending checks.
func (s *MonitorService) StartScheduler(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	// Initial run
	s.runPendingChecks(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.runPendingChecks(ctx)
		}
	}
}

func (s *MonitorService) runPendingChecks(ctx context.Context) {
	var monitors []*entity.Monitor
	err := database.DB.WithContext(ctx).Where("status = ?", "active").Find(&monitors).Error
	if err != nil {
		return
	}

	now := time.Now()
	for _, mon := range monitors {
		isDue := mon.LastCheckedAt == nil || mon.LastCheckedAt.Add(time.Duration(mon.Interval)*time.Second).Before(now)
		if isDue {
			// Run concurrently in a separate goroutine
			go s.PerformCheck(ctx, mon)
		}
	}
}

func (s *MonitorService) PerformCheck(ctx context.Context, mon *entity.Monitor) {
	var success bool
	var latency int
	var errMsg string

	if strings.ToLower(mon.Type) == "ping" {
		success, latency, errMsg = performPingCheck(mon.URL, mon.Timeout)
	} else {
		success, latency, errMsg = performHTTPCheck(mon.URL, mon.Timeout)
	}

	newStatus := "up"
	if !success {
		newStatus = "down"
	}

	oldStatus := mon.UptimeStatus
	now := time.Now()

	// Update monitor
	mon.UptimeStatus = newStatus
	mon.LastCheckedAt = &now

	err := database.DB.WithContext(ctx).Save(mon).Error
	if err != nil {
		return
	}

	// Trigger Incidents transition
	if newStatus == "down" {
		if oldStatus == "up" {
			if errMsg == "" {
				errMsg = "Connection failed"
			}
			inc := incidentEntity.NewIncident(mon.ID, mon.UserID, "active", errMsg, latency)
			database.DB.WithContext(ctx).Create(inc)
			s.event.Publish(bus.Event{Type: "incident.created", Payload: inc})
		} else {
			var activeIncidents []*incidentEntity.Incident
			err := database.DB.WithContext(ctx).Where("monitor_id = ? AND status = ?", mon.ID, "active").Find(&activeIncidents).Error
			if err == nil && len(activeIncidents) > 0 {
				activeIncidents[0].Latency = latency
				database.DB.WithContext(ctx).Save(activeIncidents[0])
			}
		}
	} else if newStatus == "up" && oldStatus == "down" {
		resolved := incidentEntity.NewIncident(mon.ID, mon.UserID, "resolved", "", latency)
		resolved.ResolvedAt = &now
		database.DB.WithContext(ctx).Create(resolved)
		s.event.Publish(bus.Event{Type: "incident.resolved", Payload: resolved})
	}

	s.event.Publish(bus.Event{Type: "monitor.checked", Payload: mon})
}

func performHTTPCheck(urlStr string, timeoutSec int) (bool, int, string) {
	// Ensure protocol is present
	if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
		urlStr = "http://" + urlStr
	}

	client := &http.Client{
		Timeout: time.Duration(timeoutSec) * time.Second,
	}

	start := time.Now()
	resp, err := client.Get(urlStr)
	latency := int(time.Since(start).Milliseconds())

	if err != nil {
		return false, latency, err.Error()
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return false, latency, fmt.Sprintf("HTTP Status Code: %d", resp.StatusCode)
	}

	return true, latency, ""
}

func performPingCheck(urlStr string, timeoutSec int) (bool, int, string) {
	// Parse URL host
	u, err := url.Parse(urlStr)
	var host string
	if err != nil || u.Host == "" {
		// Attempt parsing direct address
		host = urlStr
	} else {
		host = u.Host
	}

	// Ensure port mapping
	hostOnly := host
	port := "80"
	if strings.Contains(host, ":") {
		h, p, err := net.SplitHostPort(host)
		if err == nil {
			hostOnly = h
			port = p
		}
	} else {
		if u != nil && u.Scheme == "https" {
			port = "443"
		}
	}

	// Filter protocols prefixes if host contains it
	hostOnly = strings.TrimPrefix(hostOnly, "http://")
	hostOnly = strings.TrimPrefix(hostOnly, "https://")

	target := net.JoinHostPort(hostOnly, port)
	start := time.Now()
	conn, err := net.DialTimeout("tcp", target, time.Duration(timeoutSec)*time.Second)
	latency := int(time.Since(start).Milliseconds())

	if err != nil {
		return false, latency, err.Error()
	}
	conn.Close()

	return true, latency, ""
}

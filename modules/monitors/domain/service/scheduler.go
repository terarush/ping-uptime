package service

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"ping-uptime/internal/pkg/bus"
	"ping-uptime/internal/pkg/database"
	incidentEntity "ping-uptime/modules/incidents/domain/entity"
	"ping-uptime/modules/monitors/domain/entity"
	"strings"
	"syscall"
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
	var statusCode int
	var latency int
	var errMsg string

	switch strings.ToLower(mon.Type) {
	case "ping":
		success, latency, errMsg = performPingCheck(mon.URL, mon.Timeout)
		statusCode = 0
	case "heartbeat":
		if mon.LastCheckedAt != nil {
			deadline := mon.LastCheckedAt.Add(time.Duration(mon.Interval) * 2 * time.Second)
			if time.Now().After(deadline) {
				success = false
				errMsg = "No heartbeat received — exceeded interval"
				latency = 0
			} else {
				success = true
				latency = 0
			}
		} else {
			return
		}
		statusCode = 0
	default:
		success, statusCode, latency, errMsg = performHTTPCheck(mon.URL, mon.Timeout)
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
	if latency > 0 {
		mon.LastLatency = latency
	}

	// SSL certificate check — only for HTTPS monitors with check_ssl enabled
	if success && mon.CheckSSL && strings.HasPrefix(strings.ToLower(mon.URL), "https://") {
		sslExpiry, sslErr := checkSSLCert(mon.URL, mon.Timeout)
		if sslErr == nil {
			mon.SSLExpiresAt = sslExpiry
		}
	}

	// Check if monitor is in active maintenance window — skip incident triggers
	inMaintenance := false
	if newStatus == "down" {
		var count int64
		database.DB.WithContext(ctx).
			Table("maintenance_monitors").
			Joins("JOIN maintenances ON maintenances.id = maintenance_monitors.maintenance_id").
			Where("maintenance_monitors.monitor_id = ?", mon.ID).
			Where("maintenances.start_at <= ? AND maintenances.end_at >= ?", time.Now(), time.Now()).
			Count(&count)
		inMaintenance = count > 0
	}

	err := database.DB.WithContext(ctx).Save(mon).Error
	if err != nil {
		return
	}

	_ = s.monitorRepo.CreateCheckRecord(ctx, entity.NewCheckRecord(mon.ID, success, latency, statusCode))

	// Trigger Incidents transition — skip if in maintenance window
	if newStatus == "down" && !inMaintenance {
		if oldStatus == "up" || oldStatus == "unknown" {
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

func performHTTPCheck(urlStr string, timeoutSec int) (bool, int, int, string) {
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
		return false, 0, latency, friendlyHTTPError(err, timeoutSec)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 500 {
		return false, resp.StatusCode, latency, fmt.Sprintf("HTTP Status Code: %d", resp.StatusCode)
	}

	return true, resp.StatusCode, latency, ""
}

func friendlyHTTPError(err error, timeoutSec int) string {
	var urlErr *url.Error
	if errors.As(err, &urlErr) {
		err = urlErr.Err
	}

	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return fmt.Sprintf("Request timed out after %d seconds — the server may be down or blocking connections", timeoutSec)
	}

	var opErr *net.OpError
	if errors.As(err, &opErr) {
		var syscallErr syscall.Errno
		if errors.As(opErr.Err, &syscallErr) {
			switch syscallErr {
			case syscall.ECONNREFUSED:
				return fmt.Sprintf("Connection refused — no service listening at this address")
			case syscall.ECONNRESET:
				return fmt.Sprintf("Connection reset by server — the service may have crashed or restarted")
			case syscall.ENETUNREACH:
				return fmt.Sprintf("Network unreachable — check your internet connection or firewall")
			case syscall.ETIMEDOUT:
				return fmt.Sprintf("Connection timed out — the server is not responding")
			case syscall.EHOSTUNREACH:
				return fmt.Sprintf("Host unreachable — no route to this server")
			case syscall.EHOSTDOWN:
				return fmt.Sprintf("Host is down — the server is powered off or disconnected")
			}
		}
	}

	if errors.Is(err, os.ErrDeadlineExceeded) {
		return fmt.Sprintf("Request timed out after %d seconds", timeoutSec)
	}

	errStr := err.Error()
	if len(errStr) > 120 {
		errStr = errStr[:120]
	}
	return errStr
}

func performPingCheck(urlStr string, timeoutSec int) (bool, int, string) {
	u, err := url.Parse(urlStr)
	var host string
	if err != nil || u.Host == "" {
		host = urlStr
	} else {
		host = u.Host
	}

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

	hostOnly = strings.TrimPrefix(hostOnly, "http://")
	hostOnly = strings.TrimPrefix(hostOnly, "https://")

	target := net.JoinHostPort(hostOnly, port)
	start := time.Now()
	conn, err := net.DialTimeout("tcp", target, time.Duration(timeoutSec)*time.Second)
	latency := int(time.Since(start).Milliseconds())

	if err != nil {
		return false, latency, friendlyPingError(err, hostOnly, port, timeoutSec)
	}
	conn.Close()

	return true, latency, ""
}

func friendlyPingError(err error, host, port string, timeoutSec int) string {
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return fmt.Sprintf("Connection to %s:%s timed out after %d seconds", host, port, timeoutSec)
	}

	var opErr *net.OpError
	if errors.As(err, &opErr) {
		var syscallErr syscall.Errno
		if errors.As(opErr.Err, &syscallErr) {
			switch syscallErr {
			case syscall.ECONNREFUSED:
				return fmt.Sprintf("Connection refused — no service listening at %s:%s", host, port)
			case syscall.ECONNRESET:
				return "Connection reset by server"
			case syscall.ENETUNREACH:
				return fmt.Sprintf("Network unreachable — cannot reach %s", host)
			case syscall.ETIMEDOUT:
				return fmt.Sprintf("Connection to %s timed out", host)
			case syscall.EHOSTUNREACH:
				return fmt.Sprintf("Host unreachable — no route to %s", host)
			case syscall.EHOSTDOWN:
				return fmt.Sprintf("Host %s is down", host)
			case syscall.ENXIO:
				return fmt.Sprintf("No such device or address — invalid host %s", host)
			}
		}
	}

	if errors.Is(err, os.ErrDeadlineExceeded) {
		return fmt.Sprintf("Connection to %s:%s timed out after %d seconds", host, port, timeoutSec)
	}

	errStr := err.Error()
	if len(errStr) > 120 {
		errStr = errStr[:120]
	}
	return errStr
}

func checkSSLCert(urlStr string, timeoutSec int) (*time.Time, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	host := u.Host
	if host == "" {
		return nil, fmt.Errorf("invalid URL: no host")
	}

	dialer := &net.Dialer{Timeout: time.Duration(timeoutSec) * time.Second}
	conn, err := tls.DialWithDialer(dialer, "tcp", host, &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	certs := conn.ConnectionState().PeerCertificates
	if len(certs) == 0 {
		return nil, fmt.Errorf("no peer certificates")
	}
	return &certs[0].NotAfter, nil
}

package service

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"

	"ping-uptime/modules/ssl_monitors/domain/entity"
	"ping-uptime/modules/ssl_monitors/domain/repository"
)

var (
	ErrSSLCertNotFound = errors.New("ssl certificate record not found")
)

type SSLService struct {
	sslRepo repository.SSLRepository
}

func NewSSLService(sslRepo repository.SSLRepository) *SSLService {
	return &SSLService{sslRepo: sslRepo}
}

// CheckSSL dials the monitor's host, parses the TLS certificate, and upserts the record.
func (s *SSLService) CheckSSL(ctx context.Context, monitorID uint, targetURL string) (*entity.SSLCert, error) {
	domain, err := extractDomain(targetURL)
	if err != nil {
		return s.upsertError(ctx, monitorID, targetURL, fmt.Sprintf("invalid URL: %v", err))
	}

	host := domain
	port := "443"
	if h, p, err := net.SplitHostPort(domain); err == nil {
		host = h
		port = p
	}

	dialer := &net.Dialer{Timeout: 10 * time.Second}
	conn, err := tls.DialWithDialer(dialer, "tcp", net.JoinHostPort(host, port), &tls.Config{
		InsecureSkipVerify: false,
	})
	if err != nil {
		return s.upsertError(ctx, monitorID, domain, fmt.Sprintf("tls dial failed: %v", err))
	}
	defer conn.Close()

	certs := conn.ConnectionState().PeerCertificates
	if len(certs) == 0 {
		return s.upsertError(ctx, monitorID, domain, "no peer certificates presented")
	}

	leaf := certs[0]
	now := time.Now()
	daysRemaining := int(leaf.NotAfter.Sub(now).Hours() / 24)
	if daysRemaining < 0 {
		daysRemaining = 0
	}

	status := classifyStatus(daysRemaining)

	cert := entity.NewSSLCert(
		monitorID,
		domain,
		leaf.Issuer.CommonName,
		leaf.NotBefore,
		leaf.NotAfter,
		daysRemaining,
		status,
		"",
	)
	cert.CheckedAt = now

	// Upsert: try to find existing record for this monitor
	existing, findErr := s.sslRepo.FindByMonitorID(ctx, monitorID)
	if findErr != nil {
		// Create new
		if err := s.sslRepo.Create(ctx, cert); err != nil {
			return nil, fmt.Errorf("failed to create ssl cert record: %w", err)
		}
		return cert, nil
	}

	// Update existing
	cert.ID = existing.ID
	cert.CreatedAt = existing.CreatedAt
	if err := s.sslRepo.Update(ctx, cert); err != nil {
		return nil, fmt.Errorf("failed to update ssl cert record: %w", err)
	}
	return cert, nil
}

func (s *SSLService) GetAll(ctx context.Context) ([]*entity.SSLCert, error) {
	return s.sslRepo.FindAll(ctx)
}

func (s *SSLService) GetExpiring(ctx context.Context, days int) ([]*entity.SSLCert, error) {
	if days <= 0 {
		days = 30
	}
	return s.sslRepo.FindExpiring(ctx, days)
}

func (s *SSLService) GetByMonitorID(ctx context.Context, monitorID uint) (*entity.SSLCert, error) {
	return s.sslRepo.FindByMonitorID(ctx, monitorID)
}

func (s *SSLService) GetByID(ctx context.Context, id uint) (*entity.SSLCert, error) {
	return s.sslRepo.FindByID(ctx, id)
}

func (s *SSLService) Delete(ctx context.Context, id uint) error {
	return s.sslRepo.Delete(ctx, id)
}

// upsertError creates or updates an SSL record with error status
func (s *SSLService) upsertError(ctx context.Context, monitorID uint, domain, errMsg string) (*entity.SSLCert, error) {
	cert := entity.NewSSLCert(
		monitorID,
		domain,
		"",
		time.Time{},
		time.Time{},
		0,
		"error",
		errMsg,
	)
	cert.CheckedAt = time.Now()

	existing, findErr := s.sslRepo.FindByMonitorID(ctx, monitorID)
	if findErr != nil {
		if createErr := s.sslRepo.Create(ctx, cert); createErr != nil {
			return nil, createErr
		}
		return cert, nil
	}

	cert.ID = existing.ID
	cert.CreatedAt = existing.CreatedAt
	if updateErr := s.sslRepo.Update(ctx, cert); updateErr != nil {
		return nil, updateErr
	}
	return cert, nil
}

func extractDomain(rawURL string) (string, error) {
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		rawURL = "https://" + rawURL
	}
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	if u.Host == "" {
		return "", errors.New("no host in url")
	}
	return u.Host, nil
}

func classifyStatus(daysRemaining int) string {
	switch {
	case daysRemaining <= 0:
		return "expired"
	case daysRemaining < 14:
		return "expiring_soon"
	case daysRemaining <= 30:
		return "expiring_soon"
	default:
		return "valid"
	}
}

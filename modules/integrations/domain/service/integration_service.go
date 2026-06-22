package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"ping-uptime/modules/integrations/domain/entity"
	"ping-uptime/modules/integrations/domain/repository"
	"time"
)

var (
	ErrIntegrationNotFound = errors.New("integration not found")
)

type IntegrationService struct {
	integrationRepo repository.IntegrationRepository
}

func NewIntegrationService(integrationRepo repository.IntegrationRepository) *IntegrationService {
	return &IntegrationService{
		integrationRepo: integrationRepo,
	}
}

func (s *IntegrationService) Create(ctx context.Context, integration *entity.Integration) error {
	if integration.Name == "" || integration.Type == "" {
		return errors.New("name and type cannot be empty")
	}
	return s.integrationRepo.Create(ctx, integration)
}

func (s *IntegrationService) FindAll(ctx context.Context) ([]*entity.Integration, error) {
	return s.integrationRepo.FindAll(ctx)
}

func (s *IntegrationService) FindByID(ctx context.Context, id uint) (*entity.Integration, error) {
	return s.integrationRepo.FindByID(ctx, id)
}

func (s *IntegrationService) Update(ctx context.Context, integration *entity.Integration) error {
	return s.integrationRepo.Update(ctx, integration)
}

func (s *IntegrationService) Delete(ctx context.Context, id uint) error {
	return s.integrationRepo.Delete(ctx, id)
}

// TestIntegration sends a test payload to the integration's webhook URL.
// Extracts webhook_url from config JSON. Returns error on failure.
func (s *IntegrationService) TestIntegration(ctx context.Context, id uint) error {
	integration, err := s.integrationRepo.FindByID(ctx, id)
	if err != nil {
		return ErrIntegrationNotFound
	}

	var config struct {
		WebhookURL string `json:"webhook_url"`
	}
	if err := json.Unmarshal([]byte(integration.Config), &config); err != nil {
		return fmt.Errorf("failed to parse integration config: %w", err)
	}

	if config.WebhookURL == "" {
		return errors.New("webhook_url not found in integration config")
	}

	payload := map[string]interface{}{
		"event": "test",
		"integration": map[string]interface{}{
			"id":   integration.ID,
			"name": integration.Name,
			"type": integration.Type,
		},
		"message": "This is a test notification from Ping Uptime.",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal test payload: %w", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(config.WebhookURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to send test request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook returned non-2xx status: %d", resp.StatusCode)
	}

	return nil
}

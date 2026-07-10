package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"ping-uptime/internal/pkg/database"
	"ping-uptime/modules/backup/domain/entity"
	"ping-uptime/modules/backup/domain/repository"
)

type BackupData struct {
	Version    string        `json:"version"`
	ExportedAt string        `json:"exported_at"`
	Data       BackupPayload `json:"data"`
}

type BackupPayload struct {
	Monitors            []map[string]interface{} `json:"monitors,omitempty"`
	CheckRecords        []map[string]interface{} `json:"check_records,omitempty"`
	Incidents           []map[string]interface{} `json:"incidents,omitempty"`
	StatusPages         []map[string]interface{} `json:"status_pages,omitempty"`
	Settings            []map[string]interface{} `json:"settings,omitempty"`
	NotificationChannels []map[string]interface{} `json:"notification_channels,omitempty"`
	Tags                []map[string]interface{} `json:"tags,omitempty"`
	Integrations        []map[string]interface{} `json:"integrations,omitempty"`
	SSLCerts            []map[string]interface{} `json:"ssl_certs,omitempty"`
	ApiTokens           []map[string]interface{} `json:"api_tokens,omitempty"`
}

type BackupService struct {
	repo repository.BackupRepository
}

func NewBackupService(repo repository.BackupRepository) *BackupService {
	return &BackupService{repo: repo}
}

func (s *BackupService) Export(ctx context.Context) (*BackupData, error) {
	payload, err := s.dumpAllTables(ctx)
	if err != nil {
		return nil, fmt.Errorf("dump tables: %w", err)
	}

	return &BackupData{
		Version:    "1.0",
		ExportedAt: time.Now().UTC().Format(time.RFC3339),
		Data:       *payload,
	}, nil
}

func (s *BackupService) Import(ctx context.Context, data []byte, importUserID uint) error {
	var backup BackupData
	if err := json.Unmarshal(data, &backup); err != nil {
		return fmt.Errorf("invalid json: %w", err)
	}

	if backup.Version != "1.0" {
		return fmt.Errorf("unsupported backup version: %s", backup.Version)
	}

	if err := s.restoreTable(ctx, "monitors", backup.Data.Monitors, importUserID); err != nil {
		return fmt.Errorf("restore monitors: %w", err)
	}
	// check_records has no user_id; monitor_id preserved, so no remap.
	if err := s.restoreTable(ctx, "check_records", backup.Data.CheckRecords, 0); err != nil {
		return fmt.Errorf("restore check_records: %w", err)
	}
	if err := s.restoreTable(ctx, "incidents", backup.Data.Incidents, importUserID); err != nil {
		return fmt.Errorf("restore incidents: %w", err)
	}
	if err := s.restoreTable(ctx, "status_pages", backup.Data.StatusPages, importUserID); err != nil {
		return fmt.Errorf("restore status_pages: %w", err)
	}
	if err := s.restoreTable(ctx, "settings", backup.Data.Settings, 0); err != nil {
		return fmt.Errorf("restore settings: %w", err)
	}
	if err := s.restoreTable(ctx, "notification_channels", backup.Data.NotificationChannels, importUserID); err != nil {
		return fmt.Errorf("restore notification_channels: %w", err)
	}
	if err := s.restoreTable(ctx, "tags", backup.Data.Tags, 0); err != nil {
		return fmt.Errorf("restore tags: %w", err)
	}
	if err := s.restoreTable(ctx, "integrations", backup.Data.Integrations, importUserID); err != nil {
		return fmt.Errorf("restore integrations: %w", err)
	}
	if err := s.restoreTable(ctx, "ssl_certs", backup.Data.SSLCerts, 0); err != nil {
		return fmt.Errorf("restore ssl_certs: %w", err)
	}
	if err := s.restoreTable(ctx, "api_tokens", backup.Data.ApiTokens, importUserID); err != nil {
		return fmt.Errorf("restore api_tokens: %w", err)
	}

	return nil
}

func (s *BackupService) GetHistory(ctx context.Context) ([]*entity.BackupRecord, error) {
	return s.repo.FindAll(ctx)
}

func (s *BackupService) DeleteRecord(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func (s *BackupService) CreateRecord(ctx context.Context, fileName string, fileSize int64, userID uint) error {
	record := entity.NewBackupRecord(fileName, fileSize, userID)
	return s.repo.Create(ctx, record)
}

func (s *BackupService) dumpAllTables(ctx context.Context) (*BackupPayload, error) {
	payload := &BackupPayload{}
	var err error

	payload.Monitors, err = s.queryTable(ctx, "monitors")
	if err != nil {
		return nil, err
	}

	payload.CheckRecords, err = s.queryTable(ctx, "check_records")
	if err != nil {
		return nil, err
	}

	payload.Incidents, err = s.queryTable(ctx, "incidents")
	if err != nil {
		return nil, err
	}

	payload.StatusPages, err = s.queryTable(ctx, "status_pages")
	if err != nil {
		return nil, err
	}

	payload.Settings, err = s.queryTable(ctx, "settings")
	if err != nil {
		return nil, err
	}

	payload.NotificationChannels, err = s.queryTable(ctx, "notification_channels")
	if err != nil {
		return nil, err
	}

	payload.Tags, err = s.queryTable(ctx, "tags")
	if err != nil {
		return nil, err
	}

	payload.Integrations, err = s.queryTable(ctx, "integrations")
	if err != nil {
		return nil, err
	}

	payload.SSLCerts, err = s.queryTable(ctx, "ssl_certs")
	if err != nil {
		return nil, err
	}

	payload.ApiTokens, err = s.queryTable(ctx, "api_tokens")
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func (s *BackupService) queryTable(ctx context.Context, table string) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	err := database.DB.WithContext(ctx).Table(table).Find(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (s *BackupService) restoreTable(ctx context.Context, table string, records []map[string]interface{}, importUserID uint) error {
	if len(records) == 0 {
		return nil
	}

	// Truncate table before restore
	if err := database.DB.WithContext(ctx).Exec(fmt.Sprintf("DELETE FROM %s", table)).Error; err != nil {
		return fmt.Errorf("truncate %s: %w", table, err)
	}

	needsUserRemap := importUserID > 0

	for _, record := range records {
		if needsUserRemap {
			delete(record, "user_id")
			record["user_id"] = importUserID
		}
		if err := database.DB.WithContext(ctx).Table(table).Create(record).Error; err != nil {
			return fmt.Errorf("insert into %s: %w", table, err)
		}
	}

	return nil
}

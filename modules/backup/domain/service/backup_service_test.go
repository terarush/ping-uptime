package service

import (
	"context"
	"encoding/json"
	"testing"

	"ping-uptime/internal/pkg/database"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Verifies check_records survive an export -> import roundtrip.
// Regression: import previously dropped check_records, so analytics showed no
// data after restoring a backup.
func TestExportImportRoundtripPreservesCheckRecords(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	database.DB = db

	// Export dumps every table; create them all (only check_records seeded).
	for _, tbl := range []string{
		"monitors", "check_records", "incidents", "status_pages", "settings",
		"notification_channels", "tags", "integrations", "ssl_certs", "api_tokens",
	} {
		if err := db.Exec(`CREATE TABLE ` + tbl + ` (id INTEGER PRIMARY KEY, monitor_id INTEGER, success BOOL, latency INTEGER, status_code INTEGER, checked_at DATETIME, user_id INTEGER)`).Error; err != nil {
			t.Fatalf("create %s: %v", tbl, err)
		}
	}
	if err := db.Exec(`INSERT INTO check_records (monitor_id, success, latency, status_code, checked_at) VALUES (1, 1, 42, 200, '2026-07-09 10:00:00+00:00')`).Error; err != nil {
		t.Fatalf("seed: %v", err)
	}

	svc := NewBackupService(nil)
	ctx := context.Background()

	backup, err := svc.Export(ctx)
	if err != nil {
		t.Fatalf("export: %v", err)
	}
	if len(backup.Data.CheckRecords) != 1 {
		t.Fatalf("export check_records = %d, want 1", len(backup.Data.CheckRecords))
	}

	// Wipe then import.
	if err := db.Exec(`DELETE FROM check_records`).Error; err != nil {
		t.Fatalf("wipe: %v", err)
	}
	raw, _ := json.Marshal(backup)
	if err := svc.Import(ctx, raw, 0); err != nil {
		t.Fatalf("import: %v", err)
	}

	var count int64
	db.Table("check_records").Count(&count)
	if count != 1 {
		t.Fatalf("after import check_records = %d, want 1", count)
	}
}

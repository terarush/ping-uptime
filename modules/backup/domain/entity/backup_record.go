package entity

import "time"

type BackupRecord struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	FileName  string    `gorm:"type:varchar(255)" json:"file_name"`
	FileSize  int64     `json:"file_size"`
	UserID    uint      `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (*BackupRecord) TableName() string { return "backup_records" }

func NewBackupRecord(fileName string, fileSize int64, userID uint) *BackupRecord {
	return &BackupRecord{FileName: fileName, FileSize: fileSize, UserID: userID, CreatedAt: time.Now()}
}

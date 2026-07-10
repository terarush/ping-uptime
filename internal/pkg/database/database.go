package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

type DBModel struct {
	ServerMode   string `config:"server_mode"`
	Name         string `config:"db_name"`
	ConnLifeTime int    `config:"conn_lifetime"`
}

func (c *DBModel) OpenDB() (*gorm.DB, *error) {
	dbPath := c.Name
	if dbPath == "" {
		dbPath = "ping-uptime"
	}
	// Append .db extension if no extension is present
	if filepath.Ext(dbPath) == "" {
		dbPath = dbPath + ".db"
	}

	// Make sure parent directory exists if a custom path is specified
	dir := filepath.Dir(dbPath)
	if dir != "." && dir != "/" {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			log.Fatalf("Cannot create directory for SQLite database: %s", err.Error())
			return nil, &err
		}
	}

	// Use DELETE journal mode for single-file DB (no .db-wal/.db-shm).
	// Switch to _journal_mode=WAL if concurrent read/write throughput matters.
	dsn := fmt.Sprintf("%s?_journal_mode=DELETE&_busy_timeout=5000", dbPath)
	connection := sqlite.Open(dsn)

	db, err := gorm.Open(connection, &gorm.Config{})
	if err != nil {
		log.Fatalf("Cannot Connect to SQLite DB With Message %s", err.Error())
		return nil, &err
	}

	conPool, err := db.DB()
	if err != nil {
		log.Fatalf("Cannot Create Connection Pool to DB With Message %s", err.Error())
		return nil, &err
	}

	// Optimize connection pool for SQLite
	conPool.SetMaxIdleConns(2)
	conPool.SetMaxOpenConns(5)
	conPool.SetConnMaxLifetime(time.Duration(c.ConnLifeTime) * time.Minute)

	return db, nil
}

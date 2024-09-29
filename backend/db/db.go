package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type SQLiteOptions struct {
	MaxIdleConnections    int
	MaxOpenConnections    int
	MaxConnectionLifeTime time.Duration
	LogLevel              int
	WALEnabled            bool
	ForeignKeys           bool
}

func (o *SQLiteOptions) DSN(dbPath string) string {
	dsn := fmt.Sprintf("file:%s?cache=shared", dbPath)
	if o.WALEnabled {
		dsn += "&_journal_mode=WAL"
	}
	if o.ForeignKeys {
		dsn += "&_foreign_keys=on"
	}
	return dsn
}

func dbHomePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	dbPath := filepath.Join(homeDir, ".KTT", "db", "sqlite.db")
	return dbPath, nil
}

func getDBFilePath() (string, error) {
	dbPath, err := dbHomePath()
	if err != nil {
		return "", err
	}
	// Ensure the directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory %s: %w", dir, err)
	}
	// Check if the file exists and is writable
	if _, err := os.Stat(dbPath); err == nil {
		// File exists, check if it's writable
		file, err := os.OpenFile(dbPath, os.O_WRONLY, 0666)
		if err != nil {
			return "", fmt.Errorf("database file exists but is not writable: %w", err)
		}
		file.Close()
	} else if os.IsNotExist(err) {
		// File doesn't exist, try to create it
		file, err := os.Create(dbPath)
		if err != nil {
			return "", fmt.Errorf("failed to create database file: %w", err)
		}
		file.Close()
	} else {
		return "", fmt.Errorf("error checking database file: %w", err)
	}

	return dbPath, nil
}

func NewSQLite(opts *SQLiteOptions) (*sql.DB, error) {
	dbPath, err := getDBFilePath()
	if err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite3", opts.DSN(dbPath))
	if err != nil {
		return nil, err
	}
	return db, nil
}

// func NewSQLiteGORMDB(opts *SQLiteOptions) (*gorm.DB, error) {
// 	dbPath, err := getDBFilePath()
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get database file path: %w", err)
// 	}

// 	logLevel := logger.Silent
// 	if opts.LogLevel != 0 {
// 		logLevel = logger.LogLevel(opts.LogLevel)
// 	}

// 	db, err := gorm.Open(sqlite.Open(opts.DSN(dbPath)), &gorm.Config{
// 		Logger: logger.Default.LogMode(logLevel),
// 	})
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to open database: %w", err)
// 	}

// 	sqlDB, err := db.DB()
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get database instance: %w", err)
// 	}

// 	sqlDB.SetMaxOpenConns(opts.MaxOpenConnections)
// 	sqlDB.SetConnMaxLifetime(opts.MaxConnectionLifeTime)
// 	sqlDB.SetMaxIdleConns(opts.MaxIdleConnections)

// 	// Enable busy timeout to handle concurrent access
// 	if _, err := sqlDB.Exec("PRAGMA busy_timeout = 5000;"); err != nil {
// 		return nil, fmt.Errorf("failed to set busy timeout: %w", err)
// 	}

// 	return db, nil
// }

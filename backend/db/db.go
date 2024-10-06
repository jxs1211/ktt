package db

import (
	"database/sql"
	"fmt"
	"ktt/backend/utils/tool"
	"time"

	_ "github.com/mattn/go-sqlite3"
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

func NewSQLite(opts *SQLiteOptions) (*sql.DB, error) {
	dbPath, err := tool.GetDBFilePath()
	if err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite3", opts.DSN(dbPath))
	if err != nil {
		return nil, err
	}
	return db, nil
}

// type RedkaOptions struct {
// 	MaxIdleConnections    int
// 	MaxOpenConnections    int
// 	MaxConnectionLifeTime time.Duration
// 	LogLevel              int
// 	WALEnabled            bool
// 	ForeignKeys           bool
// }

// func (o RedkaOptions) DSN(dbPath string) string {
// 	dsn := fmt.Sprintf("file:%s", dbPath)
// 	return dsn
// }

// func NewRedka(o *redka.Options) (*redka.DB, error) {
// 	// Open or create the data.db file.
// 	dbPath, err := getDBFilePath()
// 	if err != nil {
// 		return nil, err
// 	}
// 	db, err := redka.Open(dbPath, o)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return db, nil
// }

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

package db

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"ktt/backend/utils/log"
)

//go:embed store/schema.sql
var ddl string

// initStore reads the db configuration, creates a gorm.DB instance, and initializes the miniblog store layer.
func InitStore() (*sql.DB, error) {
	dbOptions := &SQLiteOptions{
		MaxIdleConnections:    100,
		MaxOpenConnections:    100,
		MaxConnectionLifeTime: 10 * time.Second,
		LogLevel:              4,
		WALEnabled:            true,
		ForeignKeys:           true,
	}
	fmt.Printf("dbOptions:\n%+v\n", dbOptions)
	db, err := NewSQLite(dbOptions)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	// _ = store.NewStore(ins)
	// Check if the database is accessible
	if err := db.Ping(); err != nil {
		return nil, err
	}
	// create tables
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		if !strings.Contains(err.Error(), "already exist") {
			return nil, err
		}
		log.Info("InitStore", "err", err.Error())
	}

	return db, nil
}

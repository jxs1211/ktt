package db

import (
	"context"
	_ "embed"
	"fmt"
	"strings"
	"time"

	"ktt/backend/utils/log"
)

//go:embed store/schema.sql
var ddl string

// initStore reads the db configuration, creates a gorm.DB instance, and initializes the miniblog store layer.
func InitStore() error {
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
		return err
	}
	ctx := context.Background()
	// _ = store.NewStore(ins)
	// Check if the database is accessible
	if err := db.Ping(); err != nil {
		return err
	}
	// create tables
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		if !strings.Contains(err.Error(), "already exist") {
			return err
		}
		log.Info("InitStore", "err", err.Error())
	}

	return nil
}

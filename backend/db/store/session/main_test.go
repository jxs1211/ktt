package session

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	"ktt/backend/db"

	_ "github.com/mattn/go-sqlite3"
)

var (
	testDB  *sql.DB
	testCtx context.Context
)

func TestMain(m *testing.M) {
	var err error
	testDB, err = db.NewSQLite(&db.SQLiteOptions{
		MaxIdleConnections:    100,
		MaxOpenConnections:    100,
		MaxConnectionLifeTime: 10 * time.Second,
		LogLevel:              4,
		WALEnabled:            true,
		ForeignKeys:           true,
	})
	if err != nil {
		log.Fatalf("failed to initialize test database: %v", err)
	}

	testCtx = context.Background()
	// Run the tests
	code := m.Run()

	// Teardown
	if err := testDB.Close(); err != nil {
		log.Fatalf("failed to close test database: %v", err)
	}

	// Exit with the code from the tests
	os.Exit(code)
}

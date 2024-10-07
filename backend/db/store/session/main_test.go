package session

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"ktt/backend/utils/tool"
	"log"
	"os"
	"strings"
	"testing"

	// "ktt/backend/db"

	_ "github.com/mattn/go-sqlite3"
)

var (
	testDB    *sql.DB
	testCtx   context.Context
	testCmds  = strings.Join([]string{"zsh"}, ",")
	testTable = "test_sessions"
	testUpDdl = `
CREATE TABLE sessions (
  id         INTEGER  PRIMARY KEY,
  cluster_name TEXT NOT NULL,
  address     TEXT NOT NULL,
  port        TEXT NOT NULL,
  cmds        TEXT NOT NULL,
  created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`
	testDownDdl = `
DROP TABLE sessions;
`
)

func TestMain(m *testing.M) {
	var err error
	dbPath, err := tool.GetDBFilePath()
	if err != nil {
		log.Fatal(err)
	}
	testDB, err = sql.Open("sqlite3", fmt.Sprintf("file:%s?cache=shared", dbPath))
	if err != nil {
		log.Fatalf("failed to initialize test database: %v", err)
	}

	testCtx = context.Background()
	// create tables
	if _, err := testDB.ExecContext(testCtx, testUpDdl); err != nil {
		if !strings.Contains(err.Error(), "already exist") {
			log.Fatal(err)
		}
		log.Println("table already exists")
	}
	// Run the tests
	code := m.Run()

	// drop tables
	// if _, err := testDB.ExecContext(testCtx, testDownDdl); err != nil {
	// 	log.Fatal(err)
	// }
	// Teardown
	if err := testDB.Close(); err != nil {
		log.Fatalf("failed to close test database: %v", err)
	}

	// Exit with the code from the tests
	os.Exit(code)
}

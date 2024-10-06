package db

import (
	_ "embed"
	"path/filepath"
	"testing"

	_ "github.com/mattn/go-sqlite3"

	"ktt/backend/utils/log"
	strutil "ktt/backend/utils/string"
)

func TestInitStore(t *testing.T) {
	log.Init(filepath.Join(strutil.RootPath(), "logs"))
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name: "base",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// db, err := InitStore()
			// assert.NoError(t, err)
			// assert.NotNil(t, db)

			// t.Log(db.NewRedka(nil))

		})
	}
}

// func TestRedka(t *testing.T) {
// 	dbPath := "data.db"
// 	defer func() {
// 		err := os.Remove(dbPath)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 	}()
// 	db, err := redka.Open(dbPath, nil)
// 	assert.NoError(t, err)
// 	ok, err := db.Hash().SetMany("session-1211", map[string]any{
// 		"cluster_name": "test",
// 		"address":      "0.0.0.0",
// 		"port":         "1211",
// 		"cmds":         "zsh",
// 	})
// 	assert.NoError(t, err)
// 	t.Log("set ok: ", ok)
// 	fieldSet := []string{"cluster_name", "address", "port", "cmds"}
// 	v, err := db.Hash().GetMany("session-1211", fieldSet...)
// 	assert.NoError(t, err)
// 	slog.Info("hash getMany", "session-1211 value", v)
// 	// extend more field
// 	ok, err = db.Hash().SetMany("session-1212", map[string]any{
// 		"cluster_name": "test",
// 		"address":      "0.0.0.0",
// 		"port":         "1211",
// 		"cmds":         "zsh",
// 		"namespace":    "default",
// 	})
// 	assert.NoError(t, err)
// 	t.Log("set ok: ", ok)
// 	fieldSet = []string{"cluster_name", "address", "port", "cmds", "namespace"}
// 	v, err = db.Hash().GetMany("session-1212", fieldSet...)
// 	assert.NoError(t, err)
// 	slog.Info("hash getMany", "session-1212 value", v)
// }

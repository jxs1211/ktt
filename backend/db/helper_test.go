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
			if err := InitStore(); (err != nil) != tt.wantErr {
				t.Errorf("InitStore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

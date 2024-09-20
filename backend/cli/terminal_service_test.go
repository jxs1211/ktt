package cli

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"ktt/backend/types"
	"ktt/backend/utils/log"
	strutil "ktt/backend/utils/string"
)

func TestMain(m *testing.M) {
	log.Init(filepath.Join(strutil.RootPath(), "logs"))
	code := m.Run()
	os.Exit(code)
}

func TestTerminalService_StartTerminal(t *testing.T) {
	tests := []struct {
		name     string
		ctx      context.Context
		terminal terminal
		want     types.JSResp
	}{
		{
			name: "base",
			ctx:  context.Background(),
			want: types.JSResp{Success: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewTerminalService()
			s.Start(tt.ctx)
			// add testify to assert the result
			got := s.StartTerminal()
			assert.Equal(t, tt.want, got)
			err := s.CloseTerminal()
			assert.Equal(t, nil, err)
		})
	}
}

func TestTerminalService_startTerminal(t *testing.T) {
	type fields struct {
		ctx      context.Context
		terminal terminal
		mutex    sync.Mutex
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &TerminalService{
				ctx:      tt.fields.ctx,
				terminal: tt.fields.terminal,
			}
			if err := s.startTerminal(); (err != nil) != tt.wantErr {
				t.Errorf("TerminalService.startTerminal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

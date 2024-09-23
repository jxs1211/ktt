package cli

import (
	"context"
	"os"
	"path/filepath"
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
		terminal terminal
		resp     types.JSResp
	}{
		{
			name:     "gotty terminal",
			terminal: NewGottyTerminal(context.Background()),
			resp:     types.JSResp{Success: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewTerminalService()
			jsResp := s.StartTerminal()
			assert.True(t, jsResp.Success)
			err := s.CloseTerminal()
			assert.NoError(t, err)
		})
	}
}

func TestTerminalService_startTerminal(t *testing.T) {
	tests := []struct {
		name     string
		terminal terminal
		err      error
	}{
		{
			name:     "gotty terminal",
			terminal: NewGottyTerminal(context.Background()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewTerminalService()
			err := s.startTerminal()
			assert.NoError(t, err)
			err = s.CloseTerminal()
			assert.NoError(t, err)
		})
	}
}

func TestGottyTerminal_Start(t *testing.T) {
	tests := []struct {
		name    string
		ctx     context.Context
		wantErr bool
	}{
		{
			name: "start gotty terminal",
			ctx:  context.Background(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gottyTerminal := NewGottyTerminal(tt.ctx)
			err := gottyTerminal.Start()
			if (err != nil) != tt.wantErr {
				t.Errorf("GottyTerminal.Start() error = %v, wantErr %v", err, tt.wantErr)
			}
			// Clean up by closing the terminal
			err = gottyTerminal.Close()
			assert.NoError(t, err)
		})
	}
}

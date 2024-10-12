package cli

import (
	"context"
	"ktt/backend/db/store/session"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// func TestTerminalService_StartTerminal(t *testing.T) {
// 	tests := []struct {
// 		name    string
// 		address string
// 		port    string
// 		cmds    []string
// 		resp    types.JSResp
// 	}{
// 		{
// 			name:    "cli server terminal",
// 			address: "0.0.0.0",
// 			port:    "8888",
// 			cmds:    []string{"bash"},
// 			resp:    types.JSResp{Success: true},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := NewTerminalService()
// 			s.Start(context.Background())
// 			jsResp := s.StartTerminal(tt.address, tt.port, tt.cmds)
// 			assert.True(t, jsResp.Success)
// 			time.Sleep(time.Second * 2)
// 			err := s.CloseTerminal(tt.address, tt.port)
// 			assert.NoError(t, err)
// 		})
// 	}
// }

func TestTerminalService_startTerminal(t *testing.T) {
	tests := []struct {
		name        string
		clusterName string
		address     string
		port        string
		cmds        string
		eventEmit   bool
		err         error
	}{
		{
			name:        "gotty terminal",
			clusterName: "kind-test",
			address:     "127.0.0.1",
			port:        "8888",
			cmds:        "zsh",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewTerminalService(testDB, tt.eventEmit)
			s.Start(context.Background())
			session, err := s.q.CreateSession(context.Background(), session.CreateSessionParams{
				ClusterName: tt.clusterName, Address: tt.address, Port: tt.port, Cmds: tt.cmds,
			})
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				err = s.q.DeleteSession(context.Background(), session.ID)
				if err != nil {
					t.Fatal(err)
				}
			}()
			err = s.startTerminal(tt.clusterName, tt.address, tt.port, tt.cmds)
			assert.NoError(t, err)
			time.Sleep(5 * time.Second)
			err = s.closeTerminal(session.ID, tt.clusterName, tt.address, tt.port, tt.cmds)
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

// func TestTerminalService_startCliServer(t *testing.T) {
// 	tests := []struct {
// 		name        string
// 		ctx         context.Context
// 		terminalMap map[string]terminal
// 		mutex       sync.Mutex
// 		q           session.Queries
// 		address     string
// 		port        string
// 		cmds        string
// 		want        *CliServer
// 		wantErr     bool
// 	}{
// 		{
// 			name: "base",
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s, err := NewTerminalService(testDB)
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 			got, err := s.startCliServer(tt.address, tt.port, tt.cmds)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("TerminalService.startCliServer() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("TerminalService.startCliServer() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

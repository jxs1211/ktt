package cli

import (
	"context"
	"testing"
	"time"
)

func TestNewCliServer(t *testing.T) {
	tests := []struct {
		name    string
		want    *CliServer
		wantErr bool
	}{
		{
			name: "base",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv, err := NewCliServer(context.Background(), "localhost", "8888", []string{"zsh"})
			if err != nil {
				t.Fatal(err)
			}
			err = srv.Start()
			if err != nil {
				t.Fatal(err)
			}
			time.Sleep(5 * time.Second)
			t.Log("stop server")
			err = srv.Close()
			if err != nil {
				t.Fatal(err)
			}
			if err := <-srv.Errs; err != nil {
				t.Fatal(err)
			}
			// if (err != nil) != tt.wantErr {
			// 	t.Errorf("NewCliServer() error = %v, wantErr %v", err, tt.wantErr)
			// 	return
			// }
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("NewCliServer() = %v, want %v", got, tt.want)
			// }
		})
	}
}

func TestCliServer_Start(t *testing.T) {
	tests := []struct {
		name    string
		address string
		port    string
		cmds    []string
		wantErr bool
	}{
		{
			name:    "base",
			address: "0.0.0.0",
			port:    "8888",
			cmds:    []string{"bash"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewCliServer(context.Background(), tt.address, tt.port, tt.cmds)
			t.Logf("server: %+v", s)
			if err != nil {
				t.Fatal(err)
			}
			// if err := s.Start(); (err != nil) != tt.wantErr {
			// 	t.Fatalf("CliServer.Start() error = %v, wantErr %v", err, tt.wantErr)
			// }
			s.Start()
			time.Sleep(5 * time.Second)
			if err := s.Close(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

// func TestCliServer_Close(t *testing.T) {
// 	type fields struct {
// 		app        *cli.App
// 		ctx        context.Context
// 		ctxCancel  context.CancelFunc
// 		gCtx       context.Context
// 		gCtxCancel context.CancelFunc
// 		errs       chan error
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := &CliServer{
// 				app:        tt.fields.app,
// 				ctx:        tt.fields.ctx,
// 				ctxCancel:  tt.fields.ctxCancel,
// 				gCtx:       tt.fields.gCtx,
// 				gCtxCancel: tt.fields.gCtxCancel,
// 				errs:       tt.fields.errs,
// 			}
// 			if err := s.Close(); (err != nil) != tt.wantErr {
// 				t.Errorf("CliServer.Close() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

// func TestCliServer_Restart(t *testing.T) {
// 	type fields struct {
// 		app        *cli.App
// 		ctx        context.Context
// 		ctxCancel  context.CancelFunc
// 		gCtx       context.Context
// 		gCtxCancel context.CancelFunc
// 		errs       chan error
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := &CliServer{
// 				app:        tt.fields.app,
// 				ctx:        tt.fields.ctx,
// 				ctxCancel:  tt.fields.ctxCancel,
// 				gCtx:       tt.fields.gCtx,
// 				gCtxCancel: tt.fields.gCtxCancel,
// 				errs:       tt.fields.errs,
// 			}
// 			if err := s.Restart(); (err != nil) != tt.wantErr {
// 				t.Errorf("CliServer.Restart() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

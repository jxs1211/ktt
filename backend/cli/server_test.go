package cli

import (
	"context"
	"reflect"
	"testing"

	cli "github.com/urfave/cli/v2"
)

func TestNewCliServer(t *testing.T) {
	tests := []struct {
		name    string
		want    *CliServer
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCliServer()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCliServer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCliServer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCliServer_Start(t *testing.T) {
	tests := []struct {
		name       string
		app        *cli.App
		ctx        context.Context
		ctxCancel  context.CancelFunc
		gCtx       context.Context
		gCtxCancel context.CancelFunc
		errs       chan error
		wantErr    bool
	}{
		{
			name: "base",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// s := &CliServer{
			// 	app:        tt.app,
			// 	ctx:        tt.ctx,
			// 	ctxCancel:  tt.ctxCancel,
			// 	gCtx:       tt.gCtx,
			// 	gCtxCancel: tt.gCtxCancel,
			// 	errs:       tt.errs,
			// }
			s, err := NewCliServer()
			if err != nil {
				t.Fatal(err)
			}
			if err := s.Start(); (err != nil) != tt.wantErr {
				t.Errorf("CliServer.Start() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCliServer_Close(t *testing.T) {
	type fields struct {
		app        *cli.App
		ctx        context.Context
		ctxCancel  context.CancelFunc
		gCtx       context.Context
		gCtxCancel context.CancelFunc
		errs       chan error
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &CliServer{
				app:        tt.fields.app,
				ctx:        tt.fields.ctx,
				ctxCancel:  tt.fields.ctxCancel,
				gCtx:       tt.fields.gCtx,
				gCtxCancel: tt.fields.gCtxCancel,
				errs:       tt.fields.errs,
			}
			if err := s.Close(); (err != nil) != tt.wantErr {
				t.Errorf("CliServer.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCliServer_Restart(t *testing.T) {
	type fields struct {
		app        *cli.App
		ctx        context.Context
		ctxCancel  context.CancelFunc
		gCtx       context.Context
		gCtxCancel context.CancelFunc
		errs       chan error
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
			s := &CliServer{
				app:        tt.fields.app,
				ctx:        tt.fields.ctx,
				ctxCancel:  tt.fields.ctxCancel,
				gCtx:       tt.fields.gCtx,
				gCtxCancel: tt.fields.gCtxCancel,
				errs:       tt.fields.errs,
			}
			if err := s.Restart(); (err != nil) != tt.wantErr {
				t.Errorf("CliServer.Restart() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

package session

import (
	"context"
	"reflect"
	"strings"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"

	"ktt/backend/db"
)

func TestQueries_CreateSession(t *testing.T) {
	inst, err := db.NewSQLite(&db.SQLiteOptions{
		MaxIdleConnections:    100,
		MaxOpenConnections:    100,
		MaxConnectionLifeTime: 10 * time.Second,
		LogLevel:              4,
		WALEnabled:            true,
		ForeignKeys:           true,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := inst.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()
	cmds := strings.Join([]string{"zsh"}, ",")
	tests := []struct {
		name    string
		db      DBTX
		ctx     context.Context
		arg     CreateSessionParams
		want    Session
		wantErr bool
	}{
		{
			name: "base",
			db:   inst,
			ctx:  context.Background(),
			arg: CreateSessionParams{
				ClusterName: "test",
				Address:     "0.0.0.0",
				Port:        "1211",
				Cmds:        cmds,
			},
			want: Session{
				ClusterName: "test",
				Address:     "0.0.0.0",
				Port:        "1211",
				Cmds:        cmds,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Queries{
				db: tt.db,
			}
			got, err := q.CreateSession(tt.ctx, tt.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Queries.CreateSession() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want.ClusterName, got.ClusterName)
			assert.Equal(t, tt.want.Address, got.Address)
			assert.Equal(t, tt.want.Port, got.Port)
			assert.Equal(t, tt.want.Cmds, got.Cmds)
			assert.NotNil(t, got.CreatedAt)
			assert.True(t, got.CreatedAt.Valid)
			assert.NotNil(t, got.UpdatedAt)
			assert.True(t, got.UpdatedAt.Valid)
			assert.NotNil(t, got.ID)
		})
	}
}

func TestQueries_DeleteSession(t *testing.T) {
	type fields struct {
		db DBTX
	}
	type args struct {
		ctx context.Context
		id  int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Queries{
				db: tt.fields.db,
			}
			if err := q.DeleteSession(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("Queries.DeleteSession() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestQueries_GetSession(t *testing.T) {
	type fields struct {
		db DBTX
	}
	type args struct {
		ctx context.Context
		id  int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Session
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Queries{
				db: tt.fields.db,
			}
			got, err := q.GetSession(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Queries.GetSession() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Queries.GetSession() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueries_GetSessionsByClusterName(t *testing.T) {
	type fields struct {
		db DBTX
	}
	type args struct {
		ctx         context.Context
		clusterName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []Session
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Queries{
				db: tt.fields.db,
			}
			got, err := q.GetSessionsByClusterName(tt.args.ctx, tt.args.clusterName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Queries.GetSessionsByClusterName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Queries.GetSessionsByClusterName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueries_ListSessions(t *testing.T) {
	type fields struct {
		db DBTX
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []Session
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Queries{
				db: tt.fields.db,
			}
			got, err := q.ListSessions(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Queries.ListSessions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Queries.ListSessions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueries_UpdateSession(t *testing.T) {
	type fields struct {
		db DBTX
	}
	type args struct {
		ctx context.Context
		arg UpdateSessionParams
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Session
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Queries{
				db: tt.fields.db,
			}
			got, err := q.UpdateSession(tt.args.ctx, tt.args.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Queries.UpdateSession() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Queries.UpdateSession() = %v, want %v", got, tt.want)
			}
		})
	}
}

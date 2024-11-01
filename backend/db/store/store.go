package store

//go:generate mockgen -destination mock_store.go -package store ktt/backend/db/store IStore,SessionStore

import (
	"ktt/backend/db/store/session"
	"sync"

	"database/sql"
)

var (
	once sync.Once
	// 全局变量，方便其它包直接调用已初始化好的 S 实例.
	S *datastore
)

// IStore 定义了 Store 层需要实现的方法.
type IStore interface {
	DB() *sql.DB
	Sessions() session.Queries
}

// datastore 是 IStore 的一个具体实现.
type datastore struct {
	db *sql.DB
}

// 确保 datastore 实现了 IStore 接口.
var _ IStore = (*datastore)(nil)

// NewStore 创建一个 IStore 类型的实例.
func NewStore(db *sql.DB) *datastore {
	// 确保 S 只被初始化一次
	once.Do(func() {
		S = &datastore{db}
	})

	return S
}

func (ds *datastore) DB() *sql.DB {
	return ds.db
}

func (ds *datastore) Sessions() session.Queries {
	return *session.New(ds.db)
}

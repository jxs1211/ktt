package db

import (
	"context"
	"database/sql"
	"ktt/backend/db/store"
	"ktt/backend/db/store/session"
	"ktt/backend/types"
)

type DBService struct {
	ctx   context.Context
	store store.IStore
}

func NewDBService(db *sql.DB) *DBService {
	return &DBService{
		store: store.NewStore(db),
	}
}

func (s *DBService) Start(ctx context.Context) {
	s.ctx = ctx
}

func (s *DBService) GetSessionsByClusterName(name string) types.JSResp {
	sess := s.store.Sessions()
	items, err := sess.GetSessionsByClusterName(s.ctx, name)
	if err != nil {
		return types.FailedResp(err.Error())
	}
	if items == nil {
		items = []session.Session{}
	}
	return types.JSResp{
		Success: true,
		Data:    items,
	}
}

func (s *DBService) CreateSession(clusterName, address, port, cmds string) types.JSResp {
	sess := s.store.Sessions()
	session, err := sess.CreateSession(s.ctx, session.CreateSessionParams{
		ClusterName: clusterName, Address: address,
		Port: port, Cmds: cmds,
	})
	if err != nil {
		return types.FailedResp(err.Error())
	}
	return types.JSResp{Success: true, Data: session}
}

func (s *DBService) DeleteSession(id int64) types.JSResp {
	sess := s.store.Sessions()
	err := sess.DeleteSession(s.ctx, id)
	if err != nil {
		return types.FailedResp(err.Error())
	}
	return types.JSResp{Success: true}
}

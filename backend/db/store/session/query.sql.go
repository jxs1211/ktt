// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package session

import (
	"context"
)

const createSession = `-- name: CreateSession :one
INSERT INTO sessions (
  cluster_name, address, port, cmds
) VALUES (
  ?, ?, ?, ?
)
RETURNING id, cluster_name, address, port, cmds, created_at, updated_at
`

type CreateSessionParams struct {
	ClusterName string
	Address     string
	Port        string
	Cmds        string
}

func (q *Queries) CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error) {
	row := q.db.QueryRowContext(ctx, createSession,
		arg.ClusterName,
		arg.Address,
		arg.Port,
		arg.Cmds,
	)
	var i Session
	err := row.Scan(
		&i.ID,
		&i.ClusterName,
		&i.Address,
		&i.Port,
		&i.Cmds,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteSession = `-- name: DeleteSession :exec
DELETE FROM sessions
WHERE id = ?
`

func (q *Queries) DeleteSession(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteSession, id)
	return err
}

const getSession = `-- name: GetSession :one
SELECT id, cluster_name, address, port, cmds, created_at, updated_at
FROM sessions
WHERE id = ? LIMIT 1
`

func (q *Queries) GetSession(ctx context.Context, id int64) (Session, error) {
	row := q.db.QueryRowContext(ctx, getSession, id)
	var i Session
	err := row.Scan(
		&i.ID,
		&i.ClusterName,
		&i.Address,
		&i.Port,
		&i.Cmds,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getSessionByAddrAndPort = `-- name: GetSessionByAddrAndPort :one
SELECT id, cluster_name, address, port, cmds, created_at, updated_at
FROM sessions
WHERE address = ? AND port = ? LIMIT 1
`

type GetSessionByAddrAndPortParams struct {
	Address string
	Port    string
}

func (q *Queries) GetSessionByAddrAndPort(ctx context.Context, arg GetSessionByAddrAndPortParams) (Session, error) {
	row := q.db.QueryRowContext(ctx, getSessionByAddrAndPort, arg.Address, arg.Port)
	var i Session
	err := row.Scan(
		&i.ID,
		&i.ClusterName,
		&i.Address,
		&i.Port,
		&i.Cmds,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getSessionByClusterAddrPortCmds = `-- name: GetSessionByClusterAddrPortCmds :one
SELECT id, cluster_name, address, port, cmds, created_at, updated_at
FROM sessions
WHERE cluster_name = ? AND address = ? AND port = ? AND cmds = ?
LIMIT 1
`

type GetSessionByClusterAddrPortCmdsParams struct {
	ClusterName string
	Address     string
	Port        string
	Cmds        string
}

func (q *Queries) GetSessionByClusterAddrPortCmds(ctx context.Context, arg GetSessionByClusterAddrPortCmdsParams) (Session, error) {
	row := q.db.QueryRowContext(ctx, getSessionByClusterAddrPortCmds,
		arg.ClusterName,
		arg.Address,
		arg.Port,
		arg.Cmds,
	)
	var i Session
	err := row.Scan(
		&i.ID,
		&i.ClusterName,
		&i.Address,
		&i.Port,
		&i.Cmds,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getSessionByParams = `-- name: GetSessionByParams :one
SELECT id, cluster_name, address, port, cmds, created_at, updated_at
FROM sessions
WHERE address = ? AND port = ? AND cmds = ? LIMIT 1
`

type GetSessionByParamsParams struct {
	Address string
	Port    string
	Cmds    string
}

func (q *Queries) GetSessionByParams(ctx context.Context, arg GetSessionByParamsParams) (Session, error) {
	row := q.db.QueryRowContext(ctx, getSessionByParams, arg.Address, arg.Port, arg.Cmds)
	var i Session
	err := row.Scan(
		&i.ID,
		&i.ClusterName,
		&i.Address,
		&i.Port,
		&i.Cmds,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getSessionsByClusterName = `-- name: GetSessionsByClusterName :many
SELECT id, cluster_name, address, port, cmds, created_at, updated_at
FROM sessions
WHERE cluster_name = ?
ORDER BY created_at
`

func (q *Queries) GetSessionsByClusterName(ctx context.Context, clusterName string) ([]Session, error) {
	rows, err := q.db.QueryContext(ctx, getSessionsByClusterName, clusterName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Session
	for rows.Next() {
		var i Session
		if err := rows.Scan(
			&i.ID,
			&i.ClusterName,
			&i.Address,
			&i.Port,
			&i.Cmds,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listSessions = `-- name: ListSessions :many
SELECT id, cluster_name, address, port, cmds, created_at, updated_at
FROM sessions
ORDER BY cluster_name
`

func (q *Queries) ListSessions(ctx context.Context) ([]Session, error) {
	rows, err := q.db.QueryContext(ctx, listSessions)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Session
	for rows.Next() {
		var i Session
		if err := rows.Scan(
			&i.ID,
			&i.ClusterName,
			&i.Address,
			&i.Port,
			&i.Cmds,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateSession = `-- name: UpdateSession :one
UPDATE sessions
SET cluster_name = ?,
    address = ?,
    port = ?,
    cmds = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING id, cluster_name, address, port, cmds, created_at, updated_at
`

type UpdateSessionParams struct {
	ClusterName string
	Address     string
	Port        string
	Cmds        string
	ID          int64
}

func (q *Queries) UpdateSession(ctx context.Context, arg UpdateSessionParams) (Session, error) {
	row := q.db.QueryRowContext(ctx, updateSession,
		arg.ClusterName,
		arg.Address,
		arg.Port,
		arg.Cmds,
		arg.ID,
	)
	var i Session
	err := row.Scan(
		&i.ID,
		&i.ClusterName,
		&i.Address,
		&i.Port,
		&i.Cmds,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

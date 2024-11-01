// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package session

import (
	"database/sql"
)

type Session struct {
	ID          int64          `json:"id"`
	ClusterName string         `json:"cluster_name"`
	Address     string         `json:"address"`
	Port        string         `json:"port"`
	Cmds        string         `json:"cmds"`
	CreatedAt   sql.NullTime   `json:"created_at"`
	UpdatedAt   sql.NullTime   `json:"updated_at"`
}

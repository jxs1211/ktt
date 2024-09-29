-- name: GetSessionsByClusterName :many
SELECT id, cluster_name, address, port, cmds, created_at, updated_at
FROM sessions
WHERE cluster_name = ?
ORDER BY created_at;

-- name: GetSession :one
SELECT id, cluster_name, address, port, cmds, created_at, updated_at
FROM sessions
WHERE id = ? LIMIT 1;

-- name: ListSessions :many
SELECT id, cluster_name, address, port, cmds, created_at, updated_at
FROM sessions
ORDER BY cluster_name;

-- name: CreateSession :one
INSERT INTO sessions (
  cluster_name, address, port, cmds
) VALUES (
  ?, ?, ?, ?
)
RETURNING *;

-- name: UpdateSession :one
UPDATE sessions
SET cluster_name = ?,
    address = ?,
    port = ?,
    cmds = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING id, cluster_name, address, port, cmds, created_at, updated_at;

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE id = ?;

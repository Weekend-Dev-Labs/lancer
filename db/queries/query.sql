-- name: CreateSession :one
INSERT INTO sessions (file_size, chunk_size, max_chunk, file_name, temp_path, owner_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: ListSessions :many
SELECT * FROM sessions;

-- name: ListSessionDetails :many
SELECT id, file_name, owner_id, created_at FROM sessions;

-- name: UpdateSession :exec
UPDATE sessions
SET file_name = $1, temp_path = $2, current_chunk = $3
WHERE id = $4;

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE id = $1;

-- name: FindSessionsByOwner :many
SELECT * FROM sessions
WHERE owner_id = $1;

-- name: FindSessionsAfterDate :many
SELECT * FROM sessions
WHERE created_at > $1;

-- name: FindLargeSessions :many
SELECT * FROM sessions
WHERE file_size > $1;

-- name: CountSessions :one
SELECT COUNT(*) AS total_sessions FROM sessions;

-- name: TotalFileSize :one
SELECT SUM(file_size) AS total_file_size FROM sessions;

-- name: AverageChunkSize :one
SELECT AVG(chunk_size) AS average_chunk_size FROM sessions;

-- name: PaginateSessions :many
SELECT * FROM sessions
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: FindMissingTempPaths :many
SELECT * FROM sessions
WHERE temp_path IS NULL OR temp_path = '';

-- name: FindSessionsByOwnerIndexed :many
SELECT * FROM sessions
WHERE owner_id = $1;

-- name: FindSessionsByOwners :many
SELECT * FROM sessions
WHERE owner_id IN ($1); -- Requires using a slice/array in Go

-- name: FindSessionsByDate :many
SELECT * FROM sessions
WHERE DATE(created_at) = $1;

-- name: FindDuplicateFileNames :many
SELECT file_name, COUNT(*) AS duplicate_count
FROM sessions
GROUP BY file_name
HAVING COUNT(*) > 1;

-- name: FindSessionById :one
SELECT * FROM sessions
WHERE id = $1;

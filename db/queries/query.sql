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

-- name: CreateUploadedFile :one
INSERT INTO uploaded_files (
    file_name, file_path, file_size, file_type, uploaded_by, checksum, description, provider, provider_metadata
)
VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: ListUploadedFiles :many
SELECT * FROM uploaded_files;

-- name: GetUploadedFileByID :one
SELECT * FROM uploaded_files
WHERE id = $1;

-- name: UpdateUploadedFileMetadata :exec
UPDATE uploaded_files
SET description = $1, provider_metadata = $2
WHERE id = $3;

-- name: SoftDeleteUploadedFile :exec
UPDATE uploaded_files
SET is_deleted = TRUE
WHERE id = $1;

-- name: FindFilesByUser :many
SELECT * FROM uploaded_files
WHERE uploaded_by = $1;

-- name: ListActiveFiles :many
SELECT * FROM uploaded_files
WHERE is_deleted = FALSE;

-- name: FindFilesAfterDate :many
SELECT * FROM uploaded_files
WHERE uploaded_at > $1;

-- name: FindFilesByProvider :many
SELECT * FROM uploaded_files
WHERE provider = $1;

-- name: CountTotalFiles :one
SELECT COUNT(*) AS total_files FROM uploaded_files;

-- name: CountFilesByUser :one
SELECT COUNT(*) AS user_file_count
FROM uploaded_files
WHERE uploaded_by = $1;

-- name: TotalUploadFileSize :one
SELECT SUM(file_size) AS total_file_size FROM uploaded_files;

-- name: GetLargestFile :one
SELECT * FROM uploaded_files
ORDER BY file_size DESC
LIMIT 1;

-- name: PaginateUploadedFiles :many
SELECT * FROM uploaded_files
WHERE is_deleted = FALSE
ORDER BY uploaded_at DESC
LIMIT $1 OFFSET $2;

-- name: FindDuplicateFilesByChecksum :many
SELECT checksum, COUNT(*) AS duplicate_count
FROM uploaded_files
GROUP BY checksum
HAVING COUNT(*) > 1;

-- name: HardDeleteUploadedFile :exec
DELETE FROM uploaded_files
WHERE id = $1;

-- name: ListDeletedFiles :many
SELECT * FROM uploaded_files
WHERE is_deleted = TRUE;

-- name: SearchProviderMetadata :many
SELECT * FROM uploaded_files
WHERE provider_metadata @> $1;


-- name: FindFilesByMetadataKey :many
SELECT * FROM uploaded_files
WHERE provider_metadata ? $1;

-- name: CreateUser :one
INSERT INTO users (email, password)
VALUES ($1, $2)
RETURNING *;

-- name: ListUsers :many
SELECT * FROM users;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: UpdateLastLogin :exec
UPDATE users
SET last_login = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: DeleteUserByID :exec
DELETE FROM users
WHERE id = $1;

-- name: FindUsersAfterDate :many
SELECT * FROM users
WHERE created_at > $1;

-- name: GetRecentUsers :many
SELECT * FROM users
ORDER BY created_at DESC
LIMIT $1;

-- name: GetInactiveUsers :many
SELECT * FROM users
WHERE last_login < $1;

-- name: CountTotalUsers :one
SELECT COUNT(*) AS total_users FROM users;

-- name: CountUsersAfterDate :one
SELECT COUNT(*) AS recent_users
FROM users
WHERE created_at > $1;

-- name: PaginateUsers :many
SELECT * FROM users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CheckEmailExists :one
SELECT EXISTS (
    SELECT 1 FROM users WHERE email = $1
) AS email_exists;

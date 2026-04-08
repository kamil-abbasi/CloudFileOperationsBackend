-- name: FindFiles :many
SELECT *
FROM files
WHERE directory_id = $1;

-- name: FindFilesByLocation :many
SELECT *
FROM files
WHERE location LIKE $1 || '%';

-- name: FindFileById :one
SELECT *
FROM files
WHERE id = $1;

-- name: FindFileByNameAndDirectoryId :one
SELECT *
FROM files
WHERE name = $1
AND directory_id IS NOT DISTINCT FROM $2;

-- name: SaveFile :exec
INSERT INTO files(id, user_id, directory_id, name, size, location, checksum)
VALUES($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT(id)
DO UPDATE SET
	user_id = EXCLUDED.user_id,
	directory_id = EXCLUDED.directory_id,
	name = EXCLUDED.name,
	size = EXCLUDED.size,
	location = EXCLUDED.location;

-- name: RemoveFileRows :execrows
DELETE FROM files WHERE id = $1;

-- name: SaveFileItem :exec
INSERT INTO directory_items(id, user_id, type)
VALUES($1, $2, 'file')
ON CONFLICT(id)
DO NOTHING;

-- name: RemoveFileItem :exec
DELETE FROM directory_items
WHERE id = $1 AND type = 'file';

-- name: ListDirectoryItems :many
SELECT id, name, user_id, 'directory' AS type
FROM directories d
WHERE d.location = $1
UNION
SELECT id, name, user_id, 'file' AS type
FROM files f
WHERE f.location = $1;

-- name: FindDirectories :many
SELECT *
FROM directories;

-- name: FindDirectoryById :one
SELECT *
FROM directories
WHERE id = $1;

-- name: FindDirectoryByNameAndParentId :one
SELECT *
FROM directories
WHERE name = $1 AND parent_id IS NOT DISTINCT FROM $2;

-- name: SaveDirectory :exec
INSERT INTO directories(id, user_id, parent_id, name, location)
VALUES($1, $2, $3, $4, $5)
ON CONFLICT(id)
DO UPDATE SET
	user_id = EXCLUDED.user_id,
	parent_id = EXCLUDED.parent_id,
	name = EXCLUDED.name,
	location = EXCLUDED.location;

-- name: SaveDirectoryItem :exec
INSERT INTO directory_items(id, user_id, type)
VALUES($1, $2, 'directory')
ON CONFLICT(id)
DO NOTHING;

-- name: RemoveDirectoryRows :execrows
DELETE FROM directories WHERE id = $1;

-- name: RemoveDirectoryItem :exec
DELETE FROM directory_items
WHERE id = $1 AND type = 'directory';

-- name: FindUserByName :one
SELECT *
FROM users
WHERE name = $1;

-- name: SaveUser :exec
INSERT INTO users(id, name, password_hash)
VALUES($1, $2, $3)
ON CONFLICT(id)
DO UPDATE SET
name = EXCLUDED.name;

-- name: RemoveUserRows :execrows
DELETE FROM users WHERE id = $1;
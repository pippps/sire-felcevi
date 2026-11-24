-- name: Createfile :one
INSERT INTO files (id, created_at, updated_at, content_type, name, size, etag ,user_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;

-- name: GetFiles :many
SELECT
    *
FROM files
WHERE user_id = $1;

-- name: GetFile :one
SELECT 
    * 
FROM files
WHERE id = $1;

-- name: DeleteFile :exec
DELETE FROM files
WHERE id = $1;
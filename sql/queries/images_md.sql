-- name: CreateImageMeta :one
INSERT INTO images_md (id, created_at, updated_at, title, image_url , description ,user_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: GetImagesMeta :many
SELECT
    *
FROM images_md
WHERE user_id = $1;

-- name: GetImageMeta :one
SELECT 
    * 
FROM images_md
WHERE id = $1;

-- name: DeleteImageMeta :exec
DELETE FROM images_md
WHERE id = $1;


-- name: GetByID :one
SELECT id FROM tasks
WHERE id = $1
ORDER BY created_at LIMIT 1;

-- name: GetAll :many
SELECT 
  *
FROM tasks
ORDER BY created_at;

-- name: Create :one
INSERT INTO tasks (
    title, description
) VALUES (
    $1, $2
) RETURNING *;

-- name: Update :exec
UPDATE tasks
SET
    title = $2,
    description = $3
WHERE id = $1;

-- name: Delete :exec
DELETE FROM tasks
WHERE id = $1;

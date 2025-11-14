-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name)
VALUES (
    $1,
    NOW (),
    NOW (),
    $2
)
RETURNING *;

-- name: GetUsers :many
SELECT *
from users;

-- name: GetUser :one
SELECT * 
from users
WHERE name = $1;

-- name: Reset :exec
DELETE from users;
-- name: GetUsers :many
SELECT * FROM users;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: CreateUser :one
INSERT INTO users(id, created_at, updated_at, email, user_name, hashed_password)
VALUES (gen_random_uuid (), NOW(), NOW(), $1, $2, $3)
RETURNING *;
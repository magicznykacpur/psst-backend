-- name: GetChatById :one
SELECT * FROM chats WHERE id = $1;

-- name: GetChatsByUser :many
SELECT * FROM chats WHERE sender_id = $1;

-- name: CreateChatWith :one
INSERT INTO chats (id, created_at, updated_at, sender_id, receiver_id)
VALUES (
    gen_random_uuid (), NOW(), NOW(), $1, $2
)
RETURNING *;
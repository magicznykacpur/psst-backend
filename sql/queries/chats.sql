-- name: GetChatById :one
SELECT * FROM chats WHERE id = $1;

-- name: GetChatsByUser :many
SELECT * FROM chats WHERE chat_by = $1;

-- name: CreateChatWith :one
INSERT INTO chats (id, created_at, updated_at, chat_by, chat_with)
VALUES (
    gen_random_uuid (), NOW(), NOW(), $1, $2
)
RETURNING *;
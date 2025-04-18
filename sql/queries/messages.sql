-- name: GetMessageById :one
SELECT * FROM messages WHERE id = $1;

-- name: GetMessagesByChatId :many
SELECT * FROM messages WHERE chat_id = $1;

-- name: CreateMessage :one
INSERT INTO messages (id, created_at, updated_at, body, chat_id, sender_id, receiver_id)
VALUES (
    gen_random_uuid (),
    NOW(),
    NOW(),
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: DeleteMessage :exec
DELETE FROM messages WHERE id = $1 AND chat_id = $2;

-- name: GetMessageWhereChatAndUser :one
SELECT * FROM messages WHERE chat_id = $1 and sender_id = $2;
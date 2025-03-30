-- name: GetChatById :one
SELECT * FROM chats WHERE id = $1;

-- name: GetChatsByUser :many
SELECT chats.id, chats.created_at, chats.updated_at,
 receiver_id, users1.user_name as receiver_username,
 sender_id, users2.user_name as sender_username
FROM chats
JOIN users users1 ON chats.receiver_id = users1.id
JOIN users users2 ON chats.sender_id = users2.id
WHERE sender_id = $1 OR receiver_id = $1;

-- name: CreateChatWith :one
INSERT INTO chats (id, created_at, updated_at, sender_id, receiver_id)
VALUES (
    gen_random_uuid (), NOW(), NOW(), $1, $2
)
RETURNING *;

-- name: DeleteChat :exec
DELETE FROM chats WHERE id = $1 AND sender_id = $2;

-- name: GetChatByIdAndSender :one
SELECT * FROM chats WHERE id = $1 AND sender_id = $2;

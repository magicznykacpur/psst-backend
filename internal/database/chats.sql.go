// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: chats.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createChatWith = `-- name: CreateChatWith :one
INSERT INTO chats (id, created_at, updated_at, sender_id, receiver_id)
VALUES (
    gen_random_uuid (), NOW(), NOW(), $1, $2
)
RETURNING id, created_at, updated_at, sender_id, receiver_id
`

type CreateChatWithParams struct {
	SenderID   uuid.UUID
	ReceiverID uuid.UUID
}

func (q *Queries) CreateChatWith(ctx context.Context, arg CreateChatWithParams) (Chat, error) {
	row := q.db.QueryRowContext(ctx, createChatWith, arg.SenderID, arg.ReceiverID)
	var i Chat
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.SenderID,
		&i.ReceiverID,
	)
	return i, err
}

const deleteChat = `-- name: DeleteChat :exec
DELETE FROM chats WHERE id = $1 AND sender_id = $2
`

type DeleteChatParams struct {
	ID       uuid.UUID
	SenderID uuid.UUID
}

func (q *Queries) DeleteChat(ctx context.Context, arg DeleteChatParams) error {
	_, err := q.db.ExecContext(ctx, deleteChat, arg.ID, arg.SenderID)
	return err
}

const getChatById = `-- name: GetChatById :one
SELECT id, created_at, updated_at, sender_id, receiver_id FROM chats WHERE id = $1
`

func (q *Queries) GetChatById(ctx context.Context, id uuid.UUID) (Chat, error) {
	row := q.db.QueryRowContext(ctx, getChatById, id)
	var i Chat
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.SenderID,
		&i.ReceiverID,
	)
	return i, err
}

const getChatByIdAndSender = `-- name: GetChatByIdAndSender :one
SELECT id, created_at, updated_at, sender_id, receiver_id FROM chats WHERE id = $1 AND sender_id = $2
`

type GetChatByIdAndSenderParams struct {
	ID       uuid.UUID
	SenderID uuid.UUID
}

func (q *Queries) GetChatByIdAndSender(ctx context.Context, arg GetChatByIdAndSenderParams) (Chat, error) {
	row := q.db.QueryRowContext(ctx, getChatByIdAndSender, arg.ID, arg.SenderID)
	var i Chat
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.SenderID,
		&i.ReceiverID,
	)
	return i, err
}

const getChatsByUser = `-- name: GetChatsByUser :many
SELECT chats.id, chats.created_at, chats.updated_at,
 receiver_id, users1.user_name as receiver_username,
 sender_id, users2.user_name as sender_username
FROM chats
JOIN users users1 ON chats.receiver_id = users1.id
JOIN users users2 ON chats.sender_id = users2.id
WHERE sender_id = $1 OR receiver_id = $1
`

type GetChatsByUserRow struct {
	ID               uuid.UUID
	CreatedAt        time.Time
	UpdatedAt        time.Time
	ReceiverID       uuid.UUID
	ReceiverUsername string
	SenderID         uuid.UUID
	SenderUsername   string
}

func (q *Queries) GetChatsByUser(ctx context.Context, senderID uuid.UUID) ([]GetChatsByUserRow, error) {
	rows, err := q.db.QueryContext(ctx, getChatsByUser, senderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetChatsByUserRow
	for rows.Next() {
		var i GetChatsByUserRow
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.ReceiverID,
			&i.ReceiverUsername,
			&i.SenderID,
			&i.SenderUsername,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

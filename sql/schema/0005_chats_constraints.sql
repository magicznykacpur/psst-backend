-- +goose Up
ALTER TABLE chats
ADD CONSTRAINT unique_chat_between UNIQUE(sender_id, receiver_id);

-- +goose Down
ALTER TABLE chats
DROP CONSTRAINT unique_chat_between;
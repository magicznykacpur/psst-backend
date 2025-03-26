-- +goose Up
ALTER TABLE messages
ADD COLUMN sender_id UUID NOT NULL;
ALTER TABLE messages
ADD COLUMN receiver_id UUID NOT NULL;

-- +goose Down
ALTER TABLE messages
DROP COLUMN sender_id;
ALTER TABLE messages
DROP COLUMN receiver_id;
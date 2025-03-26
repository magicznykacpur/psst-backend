-- +goose Up
CREATE TABLE messages (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    body TEXT NOT NULL,
    chat_id UUID NOT NULL,
    FOREIGN KEY (chat_id) REFERENCES chats (id)
);

-- +goose Down
DROP TABLE messages;
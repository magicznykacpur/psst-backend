-- +goose Up
CREATE TABLE chats (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    chat_by UUID NOT NULL,
    chat_with UUID NOT NULL,
    FOREIGN KEY (chat_by) REFERENCES users (id),
    FOREIGN KEY (chat_with) REFERENCES users (id)
);

-- +goose Down
DROP TABLE chats;
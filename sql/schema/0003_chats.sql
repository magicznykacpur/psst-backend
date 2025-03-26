-- +goose Up
CREATE TABLE chats (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    sender_id UUID NOT NULL,
    receiver_id UUID NOT NULL,
    FOREIGN KEY (sender_id) REFERENCES users (id),
    FOREIGN KEY (receiver_id) REFERENCES users (id)
);

-- +goose Down
DROP TABLE chats;
-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    email TEXT NOT NULL,
    user_name TEXT NOT NULL,
    hashed_password TEXT NOT NULL
);

-- +goose Down
DROP TABLE users;
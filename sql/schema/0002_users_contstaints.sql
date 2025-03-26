-- +goose Up
ALTER TABLE users
ADD CONSTRAINT email_unq UNIQUE(email);
ALTER TABLE users
ADD CONSTRAINT user_name_unq UNIQUE(user_name);

-- +goose DOwn
ALTER TABLE users
DROP CONSTRAINT email_unq;
ALTER TABLE users
DROP CONSTRAINT user_name_unq;
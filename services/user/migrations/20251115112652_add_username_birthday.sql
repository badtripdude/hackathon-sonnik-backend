-- +goose Up
ALTER TABLE users
ADD COLUMN username TEXT NOT NULL DEFAULT '',
ADD COLUMN birth_date DATE NOT NULL DEFAULT '1970-01-01';

-- +goose Down
ALTER TABLE users
DROP COLUMN birth_date,
DROP COLUMN username;

-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
ADD COLUMN avatar BYTEA;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
DROP COLUMN avatar;
-- +goose StatementEnd

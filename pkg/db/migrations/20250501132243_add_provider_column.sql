-- +goose Up
ALTER TABLE transaction ADD COLUMN provider VARCHAR(255);

-- +goose Down
ALTER TABLE transaction DROP COLUMN provider;

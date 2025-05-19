-- +goose Up
ALTER TABLE notifications ADD COLUMN type VARCHAR(255);
ALTER TABLE notifications ADD COLUMN status VARCHAR(255);

-- +goose Down
ALTER TABLE notifications DROP COLUMN type;
ALTER TABLE notifications DROP COLUMN status;


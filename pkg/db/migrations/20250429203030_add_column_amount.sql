-- +goose Up
ALTER TABLE transaction ADD COLUMN amount NUMERIC(20, 2) NOT NULL default 0.00;

-- +goose Down
ALTER TABLE transaction DROP COLUMN amount;


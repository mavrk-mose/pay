-- +goose Up
ALTER TABLE transaction DROP COLUMN credit_amount;
ALTER TABLE transaction DROP COLUMN debit_amount;

-- +goose Down
ALTER TABLE transaction ADD COLUMN debit_amount NUMERIC(20, 2) NOT NULL default 0.00;
ALTER TABLE transaction ADD COLUMN credit_amount NUMERIC(20, 2) NOT NULL default 0.00;

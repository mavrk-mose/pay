-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS templates
(
    id      VARCHAR(255) PRIMARY KEY,
    title   TEXT,
    message TEXT,
    type    VARCHAR(255)
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS templates;
-- +goose StatementEnd

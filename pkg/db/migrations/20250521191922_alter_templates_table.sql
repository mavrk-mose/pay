-- +goose Up
-- +goose StatementBegin
ALTER TABLE templates
    ADD COLUMN subject TEXT, -- optional, for email subjects
    ADD COLUMN channel VARCHAR(50) NOT NULL DEFAULT 'email', -- web, push, sms, email
    ADD COLUMN variables JSONB DEFAULT '[]'::jsonb, -- list of expected variables like ["user_name"]
    ADD COLUMN metadata JSONB DEFAULT '{}'::jsonb, -- optional extra info (e.g. language)
    ADD COLUMN created_at TIMESTAMPTZ DEFAULT NOW(),
    ADD COLUMN updated_at TIMESTAMPTZ DEFAULT NOW();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE templates
    DROP COLUMN IF EXISTS subject,
    DROP COLUMN IF EXISTS channel,
    DROP COLUMN IF EXISTS variables,
    DROP COLUMN IF EXISTS metadata,
    DROP COLUMN IF EXISTS created_at,
    DROP COLUMN IF EXISTS updated_at;
-- +goose StatementEnd

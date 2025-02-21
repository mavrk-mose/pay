-- Create custom enums for entry type and transaction status
CREATE TYPE entry_type AS ENUM ('debit', 'credit');
CREATE TYPE transaction_status AS ENUM ('pending', 'confirmed', 'failed');
CREATE TYPE wallet_status AS ENUM ('active','terminated')

-- Create custom enum types (if you prefer enums)
CREATE TYPE transaction_type AS ENUM ('withdrawal', 'deposit', 'transfer', 'charge');
CREATE TYPE transaction_status AS ENUM ('pending', 'confirmed', 'failed');

-- Create the combined transaction records table
CREATE TABLE transaction
(
    id               SERIAL PRIMARY KEY,
    external_ref     VARCHAR(100) UNIQUE NOT NULL,
    type             transaction_type    NOT NULL,
    status           transaction_status  NOT NULL DEFAULT 'pending',
    details          TEXT, -- or JSONB if you want to store JSON-formatted extra details
    currency         VARCHAR(10)         NOT NULL,
    debit_wallet_id  BIGINT              NOT NULL,
    debit_amount     NUMERIC(20, 2)      NOT NULL,
    credit_wallet_id BIGINT              NOT NULL,
    credit_amount    NUMERIC(20, 2)      NOT NULL,
    created_at       TIMESTAMPTZ                  DEFAULT NOW(),
    updated_at       TIMESTAMPTZ                  DEFAULT NOW()
);

CREATE TABLE product_configurations
(
    id             SERIAL PRIMARY KEY,
    product_name   VARCHAR(255)  NOT NULL,
    fee_percentage DECIMAL(5, 2) NOT NULL, -- Fee percentage (e.g., 2.5%)
    created_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE users
(
    id            SERIAL PRIMARY KEY,
    google_id     TEXT UNIQUE NOT NULL,
    name          TEXT        NOT NULL,
    email         TEXT UNIQUE NOT NULL,
    avatar_url    TEXT,
    location      TEXT,
    language      TEXT,
    currency      TEXT,
    notification_channel TEXT NOT NULL, -- e.g., "push", "sms", "email"
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_login_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE wallets (
    id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES users(google_id) ON DELETE CASCADE,
    balance NUMERIC DEFAULT 0.00,
    status wallet_status NOT NULL DEFAULT 'active',
    currency VARCHAR(10) DEFAULT 'USD',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

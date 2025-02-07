-- Create custom enums for entry type and transaction status
CREATE TYPE entry_type AS ENUM ('debit', 'credit');
CREATE TYPE transaction_status AS ENUM ('pending', 'confirmed', 'failed');

-- Create custom enum types (if you prefer enums)
CREATE TYPE transaction_type AS ENUM ('withdrawal', 'deposit', 'transfer', 'charge');
CREATE TYPE transaction_status AS ENUM ('pending', 'confirmed', 'failed');

-- Create the combined transaction records table
CREATE TABLE transaction (
    id SERIAL PRIMARY KEY,
    external_ref VARCHAR(100) UNIQUE NOT NULL,
    type transaction_type NOT NULL,
    status transaction_status NOT NULL DEFAULT 'pending',
    details TEXT,  -- or JSONB if you want to store JSON-formatted extra details
    currency VARCHAR(10) NOT NULL,
    debit_wallet_id BIGINT NOT NULL,
    debit_amount NUMERIC(20,2) NOT NULL,
    credit_wallet_id BIGINT NOT NULL,
    credit_amount NUMERIC(20,2) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

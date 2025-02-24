-- Create custom enums for entry type and transaction status
CREATE TYPE entry_type AS ENUM ('debit', 'credit');
CREATE TYPE transaction_status AS ENUM ('pending', 'confirmed', 'failed');

-- custom enum for wallet 
CREATE TYPE wallet_status AS ENUM ('active','terminated')

-- custom enum for discount types
CREATE TYPE discount_type AS ENUM (
    'merchant',     -- Merchant-specific discount
    'referral',     -- Discount awarded for successful referrals
    'loyalty',      -- Reward for repeat or long-term customers
    'seasonal',     -- Special discount during holiday or seasonal events
    'volume',       -- Discount based on bulk purchasing
    'promotional'   -- Time-limited marketing promotion discount
);

-- Create custom enum types (if you prefer enums)
CREATE TYPE transaction_type AS ENUM ('withdrawal', 'deposit', 'transfer', 'charge','tax');
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

CREATE TABLE referrals (
    id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL,
    referral_code TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE referral_usages (
    id SERIAL PRIMARY KEY,
    applied_user_id TEXT NOT NULL,
    referral_code TEXT NOT NULL,
    bonus NUMERIC DEFAULT 5.0,
    applied_at TIMESTAMP DEFAULT NOW()
);

-- Table for storing vouchers
CREATE TABLE vouchers (
    id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL,
    code TEXT NOT NULL UNIQUE,
    amount NUMERIC NOT NULL,
    currency TEXT NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    redeemed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Table for merchant discounts
CREATE TABLE discounts (
    id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL,
    type discount_type NOT NULL DEFAULT ,
    discount_percentage NUMERIC NOT NULL,
    valid_from TIMESTAMP NOT NULL,
    valid_until TIMESTAMP NOT NULL
);

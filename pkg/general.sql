-- Enums
CREATE TYPE entry_type AS ENUM (
    'debit',
    'credit'
    );

CREATE TYPE wallet_status AS ENUM (
    'active', -- Wallet is fully operational
    'suspended', -- Wallet temporarily restricted due to fraud checks, verification issues, or policy violations
    'terminated', -- Wallet is permanently closed and cannot be reactivated
    'frozen', -- Wallet is locked due to security reasons, disputes, or legal compliance
    'pending', -- Wallet is created but awaiting verification or activation
    'restricted' -- Wallet has limited functionality (e.g., can't withdraw, but can receive funds)
    );

CREATE TYPE discount_type AS ENUM (
    'merchant', -- Merchant-specific discount
    'referral', -- Discount awarded for successful referrals
    'loyalty', -- Reward for repeat or long-term customers
    'seasonal', -- Special discount during holiday or seasonal events
    'volume', -- Discount based on bulk purchasing
    'promotional' -- Time-limited marketing promotion discount
    );

CREATE TYPE transaction_type AS ENUM (
    'withdrawal', -- Money withdrawn from an account
    'deposit', -- Money added to an account
    'transfer', -- Money moved between accounts
    'charge', -- Fee charged for services
    'tax', -- Tax deduction from a transaction
    'refund', -- Money returned to a customer
    'reversal', -- A previously completed transaction is reversed
    'cashback', -- Rewards given to a user for transactions
    'fee', -- Additional processing or service fees
    'payout', -- Funds disbursed to merchants or users
    'hold', -- Funds temporarily locked for verification
    'release', -- Previously held funds are made available
    'adjustment', -- Manual correction to an account balance
    'subscription' -- Recurring payment for a service
    );

CREATE TYPE transaction_status AS ENUM (
    'pending', -- Transaction is initiated but not yet processed
    'processing', -- Transaction is currently being processed
    'confirmed', -- Successfully completed transaction
    'failed', -- Transaction failed due to an error
    'reversed', -- Funds returned after a successful transaction was reversed
    'refunded', -- Transaction amount has been refunded to the payer
    'chargeback', -- Payment disputed and refunded by financial institution
    'on_hold', -- Temporarily held for fraud checks or verification
    'expired', -- Transaction expired before completion
    'canceled', -- Transaction was manually canceled by user or system
    'partially_paid', -- Partial payment received (e.g., installment payments)
    'disputed' -- Transaction is under review due to a dispute
    );

-- Tables
CREATE TABLE IF NOT EXISTS transaction
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
    checksum         TEXT                NOT NULL,
    created_at       TIMESTAMPTZ                  DEFAULT NOW(),
    updated_at       TIMESTAMPTZ                  DEFAULT NOW()
    --     pid VAR VARCHAR(255) NOT NULL,
    --     src_user_id VARCHAR(255) NOT NULL,
    --     dst_user_id VARCHAR(255) NOT NULL,
    --     src_wallet_id INT NOT NULL,
    --     dst_wallet_id INT NOT NULL,
    --     src_account_id INT NOT NULL,
    --     dst_account_id INT NOT NULL,
    --     src_account_type VARCHAR(255) NOT NULL,
    --     dst_account_type VARCHAR(255) NOT NULL,
    --     final_dst_merchant_wallet_id INT,
    --     amount INT NOT NULL,
);

CREATE TABLE IF NOT EXISTS product_configurations
(
    id             SERIAL PRIMARY KEY,
    product_name   VARCHAR(255)  NOT NULL,
    fee_percentage DECIMAL(5, 2) NOT NULL, -- Fee percentage (e.g., 2.5%)
    created_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS users
(
    id                   SERIAL PRIMARY KEY,
    google_id            TEXT UNIQUE NOT NULL,
    name                 TEXT        NOT NULL,
    email                TEXT UNIQUE NOT NULL,
    avatar_url           TEXT,
    location             TEXT,
    language             TEXT,
    currency             TEXT,
    role                 TEXT,
    notification_channel TEXT        NOT NULL, -- e.g., "push", "sms", "email"
    created_at           TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at           TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_login_at        TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    onboarded            BOOLEAN   DEFAULT FALSE,
    provider             TEXT      DEFAULT 'google'
);

CREATE TABLE IF NOT EXISTS wallets
(
    id         SERIAL PRIMARY KEY,
    user_id    VARCHAR(255) REFERENCES users (google_id) ON DELETE CASCADE,
    balance    NUMERIC                DEFAULT 0.00,
    status     wallet_status NOT NULL DEFAULT 'active',
    currency   VARCHAR(10)            DEFAULT 'USD',
    created_at TIMESTAMP              DEFAULT NOW(),
    updated_at TIMESTAMP              DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS referrals
(
    id            SERIAL PRIMARY KEY,
    user_id       TEXT NOT NULL,
    referral_code TEXT NOT NULL UNIQUE,
    created_at    TIMESTAMP DEFAULT NOW()

);

CREATE TABLE IF NOT EXISTS referral_usages
(
    id              SERIAL PRIMARY KEY,
    applied_user_id TEXT NOT NULL,
    referral_code   TEXT NOT NULL,
    bonus           NUMERIC   DEFAULT 5.0,
    applied_at      TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS vouchers
(
    id         SERIAL PRIMARY KEY,
    user_id    TEXT      NOT NULL,
    code       TEXT      NOT NULL UNIQUE,
    amount     NUMERIC   NOT NULL,
    currency   TEXT      NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    redeemed   BOOLEAN   DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS discounts
(
    id                  SERIAL PRIMARY KEY,
    user_id             TEXT          NOT NULL,
    type                discount_type NOT NULL DEFAULT 'merchant',
    discount_percentage NUMERIC       NOT NULL,
    valid_from          TIMESTAMP     NOT NULL,
    valid_until         TIMESTAMP     NOT NULL
);

CREATE TABLE IF NOT EXISTS notifications
(
    id         SERIAL PRIMARY KEY,
    user_id    UUID        NOT NULL,
    channel    VARCHAR(10) NOT NULL, -- 'sms', 'web', 'email', 'push'
    title      TEXT,
    message    TEXT,
    metadata   JSONB,                -- optional: store extra info
    is_read    BOOLEAN   DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT now()
);
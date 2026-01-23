CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    login TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
    );

CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY,
    number TEXT UNIQUE NOT NULL,
    user_id UUID NOT NULL,
    status TEXT NOT NULL,
    accrual NUMERIC,
    uploaded_at TIMESTAMPTZ NOT NULL DEFAULT now()
    );

CREATE TABLE IF NOT EXISTS balances (
    user_id UUID PRIMARY KEY,
    current NUMERIC NOT NULL DEFAULT 0,
    withdrawn NUMERIC NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS withdrawals (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    order_number TEXT NOT NULL,
    sum NUMERIC NOT NULL,
    processed_at TIMESTAMPTZ NOT NULL DEFAULT now()
    );

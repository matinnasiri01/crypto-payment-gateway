-- ============================================================
-- 000001_init.up.sql
-- ============================================================

CREATE TABLE IF NOT EXISTS users (
                                     id               UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    email            VARCHAR(255) NOT NULL UNIQUE,
    password_hash    TEXT         NOT NULL,
    withdraw_address TEXT         NOT NULL,
    balance          NUMERIC(36,18) NOT NULL DEFAULT 0,
    created_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW()
    );

CREATE TABLE IF NOT EXISTS invoices (
                                        id              UUID           PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID           NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status          TEXT           NOT NULL DEFAULT 'pending'
    CHECK (status IN ('pending', 'paid', 'cancelled', 'expired')),
    amount          NUMERIC(36,18) NOT NULL,
    description     TEXT,
    callback_url    TEXT,
    pay_to_address  TEXT           NOT NULL,
    paid_by_address TEXT,
    overpayment     NUMERIC(36,18) DEFAULT 0,
    created_at      TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    expired_at      TIMESTAMPTZ    NOT NULL
    );

CREATE INDEX IF NOT EXISTS idx_invoices_user_id    ON invoices(user_id);
CREATE INDEX IF NOT EXISTS idx_invoices_status     ON invoices(status);
CREATE INDEX IF NOT EXISTS idx_invoices_expired_at ON invoices(expired_at);
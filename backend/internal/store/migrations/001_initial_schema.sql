-- 001_initial_schema.sql
-- Lead Automation MVP - Initial Database Schema

-- Enable UUID extension (optional, we use BIGSERIAL for now)
-- CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ========================
-- Users (SaaS customers)
-- ========================
CREATE TABLE IF NOT EXISTS users (
    id            BIGSERIAL PRIMARY KEY,
    email         VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name     VARCHAR(255) NOT NULL,
    company_name  VARCHAR(255) DEFAULT '',
    plan          VARCHAR(50) NOT NULL DEFAULT 'starter',
    is_active     BOOLEAN NOT NULL DEFAULT TRUE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);

-- ========================
-- Channels (WA / IG / FB)
-- ========================
CREATE TABLE IF NOT EXISTS channels (
    id            BIGSERIAL PRIMARY KEY,
    user_id       BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    platform      VARCHAR(50) NOT NULL, -- 'whatsapp', 'instagram', 'facebook'
    account_id    VARCHAR(255) NOT NULL,
    account_name  VARCHAR(255) DEFAULT '',
    access_token  TEXT NOT NULL DEFAULT '',
    refresh_token TEXT DEFAULT '',
    token_expiry  TIMESTAMPTZ,
    is_active     BOOLEAN NOT NULL DEFAULT TRUE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, platform, account_id)
);

CREATE INDEX idx_channels_user ON channels(user_id);

-- ========================
-- Contacts (Leads)
-- ========================
CREATE TABLE IF NOT EXISTS contacts (
    id                  BIGSERIAL PRIMARY KEY,
    user_id             BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    channel_id          BIGINT REFERENCES channels(id) ON DELETE SET NULL,
    platform            VARCHAR(50) NOT NULL,
    platform_user_id    VARCHAR(255) NOT NULL,
    name                VARCHAR(255) DEFAULT '',
    phone               VARCHAR(50) DEFAULT '',
    email               VARCHAR(255) DEFAULT '',
    budget              VARCHAR(100) DEFAULT '',
    preferred_location  VARCHAR(255) DEFAULT '',
    purchase_timeline   VARCHAR(100) DEFAULT '',
    tags                TEXT[] DEFAULT '{}',
    is_hot_lead         BOOLEAN NOT NULL DEFAULT FALSE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, platform, platform_user_id)
);

CREATE INDEX idx_contacts_user ON contacts(user_id);
CREATE INDEX idx_contacts_hot ON contacts(user_id, is_hot_lead) WHERE is_hot_lead = TRUE;

-- ========================
-- Messages
-- ========================
CREATE TABLE IF NOT EXISTS messages (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    channel_id      BIGINT REFERENCES channels(id) ON DELETE SET NULL,
    contact_id      BIGINT NOT NULL REFERENCES contacts(id) ON DELETE CASCADE,
    platform        VARCHAR(50) NOT NULL,
    direction       VARCHAR(10) NOT NULL, -- 'inbound', 'outbound'
    content         TEXT NOT NULL DEFAULT '',
    message_type    VARCHAR(50) NOT NULL DEFAULT 'text',
    platform_msg_id VARCHAR(255) DEFAULT '',
    status          VARCHAR(20) NOT NULL DEFAULT 'sent',
    is_automated    BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_messages_contact ON messages(contact_id);
CREATE INDEX idx_messages_user_created ON messages(user_id, created_at DESC);

-- ========================
-- Automations (Rules)
-- ========================
CREATE TABLE IF NOT EXISTS automations (
    id           BIGSERIAL PRIMARY KEY,
    user_id      BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name         VARCHAR(255) NOT NULL,
    trigger_type VARCHAR(50) NOT NULL DEFAULT 'keyword', -- 'keyword', 'first_message'
    keywords     TEXT[] DEFAULT '{}',
    reply_text   TEXT NOT NULL DEFAULT '',
    reply_media  TEXT DEFAULT '',
    delay_ms     INT NOT NULL DEFAULT 0,
    is_active    BOOLEAN NOT NULL DEFAULT TRUE,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_automations_user ON automations(user_id);

-- ========================
-- Broadcasts
-- ========================
CREATE TABLE IF NOT EXISTS broadcasts (
    id           BIGSERIAL PRIMARY KEY,
    user_id      BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name         VARCHAR(255) NOT NULL,
    content      TEXT NOT NULL DEFAULT '',
    media_url    TEXT DEFAULT '',
    status       VARCHAR(20) NOT NULL DEFAULT 'draft',
    total_sent   INT NOT NULL DEFAULT 0,
    total_failed INT NOT NULL DEFAULT 0,
    scheduled_at TIMESTAMPTZ,
    sent_at      TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_broadcasts_user ON broadcasts(user_id);

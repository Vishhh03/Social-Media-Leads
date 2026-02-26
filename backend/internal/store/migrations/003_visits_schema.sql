-- 003_visits_schema.sql
-- Lead Automation MVP - Visits Booking Table

CREATE TABLE IF NOT EXISTS visits (
    id                  BIGSERIAL PRIMARY KEY,
    user_id             BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    contact_id          BIGINT NOT NULL REFERENCES contacts(id) ON DELETE CASCADE,
    project_name        VARCHAR(255) NOT NULL,
    visit_time          TIMESTAMPTZ NOT NULL,
    status              VARCHAR(50) NOT NULL DEFAULT 'confirmed', -- 'confirmed', 'rescheduled', 'completed', 'cancelled'
    lead_source_channel VARCHAR(50) NOT NULL, -- 'ig', 'wa'
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_visits_user ON visits(user_id);
CREATE INDEX idx_visits_contact ON visits(contact_id);
CREATE INDEX idx_visits_time ON visits(visit_time);

-- Add tracking fields to contacts
ALTER TABLE contacts ADD COLUMN IF NOT EXISTS booking_state VARCHAR(50) NOT NULL DEFAULT 'new';
ALTER TABLE contacts ADD COLUMN IF NOT EXISTS bot_paused BOOLEAN NOT NULL DEFAULT FALSE;

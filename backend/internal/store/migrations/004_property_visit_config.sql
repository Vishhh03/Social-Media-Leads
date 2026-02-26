-- 004_property_visit_config.sql
-- Stores the output of the TemplateWizard per tenant.
-- Intentionally flat. No polymorphism. No future-proofing.

CREATE TABLE IF NOT EXISTS property_visit_configs (
    id           BIGSERIAL PRIMARY KEY,
    user_id      BIGINT NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    project_name TEXT NOT NULL,
    brochure_url TEXT NOT NULL DEFAULT '',
    agent_phone  TEXT NOT NULL DEFAULT '',
    is_active    BOOLEAN NOT NULL DEFAULT TRUE,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_pvc_user_id ON property_visit_configs(user_id);

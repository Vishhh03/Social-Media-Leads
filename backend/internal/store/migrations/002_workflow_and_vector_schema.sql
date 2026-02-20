-- 002_workflow_and_vector_schema.sql
-- Lead Automation MVP - AI Workflow Orchestrator & Vector DB

-- Enable pgvector extension
CREATE EXTENSION IF NOT EXISTS vector;

-- ========================
-- Knowledge Base (RAG)
-- ========================
CREATE TABLE IF NOT EXISTS knowledge_base (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title       VARCHAR(255) NOT NULL,
    content     TEXT NOT NULL,
    embedding   VECTOR(1536), -- Assuming OpenAI ada-002 dimensions
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_knowledge_base_user ON knowledge_base(user_id);
-- HNSW Index for fast similarity search
CREATE INDEX idx_knowledge_base_embedding ON knowledge_base USING hnsw (embedding vector_cosine_ops);

-- ========================
-- Workflows (Orchestrator Blueprint)
-- ========================
CREATE TABLE IF NOT EXISTS workflows (
    id           BIGSERIAL PRIMARY KEY,
    user_id      BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name         VARCHAR(255) NOT NULL,
    trigger_type VARCHAR(100) NOT NULL, -- e.g., 'meta_dm_received', 'tag_added'
    status       VARCHAR(50) NOT NULL DEFAULT 'draft', -- 'draft', 'published'
    prompt       TEXT, -- The original AI prompt used to generate this (if any)
    nodes        JSONB NOT NULL DEFAULT '[]'::jsonb,
    edges        JSONB NOT NULL DEFAULT '[]'::jsonb,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_workflows_user ON workflows(user_id);
CREATE INDEX idx_workflows_user_trigger ON workflows(user_id, trigger_type);

-- ========================
-- Workflow Executions (Running State)
-- ========================
CREATE TABLE IF NOT EXISTS workflow_executions (
    id              BIGSERIAL PRIMARY KEY,
    workflow_id     BIGINT NOT NULL REFERENCES workflows(id) ON DELETE CASCADE,
    contact_id      BIGINT NOT NULL REFERENCES contacts(id) ON DELETE CASCADE,
    current_node_id VARCHAR(255) NOT NULL,
    status          VARCHAR(50) NOT NULL DEFAULT 'running', -- 'running', 'waiting', 'completed', 'failed'
    state_data      JSONB NOT NULL DEFAULT '{}'::jsonb, -- Context payload passed between nodes
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_workflow_executions_workflow_contact ON workflow_executions(workflow_id, contact_id);
CREATE INDEX idx_workflow_executions_status ON workflow_executions(status);

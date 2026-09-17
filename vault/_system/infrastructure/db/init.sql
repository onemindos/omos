-- =============================================================================
-- OmniDB Schema for OneMind OS
-- PostgreSQL 16 + TimescaleDB + pgvector
-- =============================================================================
-- ONE database, THREE superpowers:
--   1. PostgreSQL 16  — full SQL, JSONB, triggers, foreign keys
--   2. TimescaleDB    — hypertables, compression, continuous aggregates
--   3. pgvector       — vector embeddings, semantic search, cosine similarity
--
-- This is the MIND'S MEMORY — both temporal and semantic in one place.
-- =============================================================================

-- Create extensions
CREATE EXTENSION IF NOT EXISTS timescaledb;
CREATE EXTENSION IF NOT EXISTS vector;

-- ============================================
-- EVENTS TABLE (Hypertable — TimescaleDB)
-- All events archived from NATS JetStream
-- ============================================
CREATE TABLE events (
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    subject TEXT NOT NULL,
    data JSONB NOT NULL,
    source TEXT,                        -- Which service produced the event
    PRIMARY KEY (timestamp, subject)
);

-- Convert to hypertable (automatic time partitioning)
-- Chunks by 1 day for efficient querying and compression
SELECT create_hypertable('events', 'timestamp', chunk_time_interval => INTERVAL '1 day');

-- ============================================
-- INDEXES FOR COMMON QUERIES
-- ============================================

-- Subject queries: "Show me all heart rate events"
CREATE INDEX idx_events_subject ON events (subject, timestamp DESC);

-- JSONB queries: "Find events where value > 100"
CREATE INDEX idx_events_data ON events USING GIN (data);

-- Source queries: "What events came from Home Assistant?"
CREATE INDEX idx_events_source ON events (source, timestamp DESC);

-- ============================================
-- COMPRESSION POLICY
-- Compress chunks older than 7 days (95%+ compression)
-- ============================================
ALTER TABLE events SET (
    timescaledb.compress,
    timescaledb.compress_segmentby = 'subject'
);

SELECT add_compression_policy('events', INTERVAL '7 days');

-- ============================================
-- RETENTION POLICY (Optional)
-- Uncomment to auto-delete data older than 2 years
-- ============================================
-- SELECT add_retention_policy('events', INTERVAL '2 years');

-- ============================================
-- CONTINUOUS AGGREGATES
-- Pre-computed rollups for fast dashboard queries
-- ============================================

-- Daily event counts by subject
CREATE MATERIALIZED VIEW daily_event_counts
WITH (timescaledb.continuous) AS
SELECT
    time_bucket('1 day', timestamp) AS day,
    subject,
    COUNT(*) as event_count
FROM events
GROUP BY day, subject
WITH NO DATA;

-- Refresh policy: update every hour
SELECT add_continuous_aggregate_policy('daily_event_counts',
    start_offset => INTERVAL '3 days',
    end_offset => INTERVAL '1 hour',
    schedule_interval => INTERVAL '1 hour');

-- Daily biometric averages (for health dashboard)
CREATE MATERIALIZED VIEW daily_biometrics
WITH (timescaledb.continuous) AS
SELECT
    time_bucket('1 day', timestamp) AS day,
    subject,
    AVG((data->>'value')::float) as avg_value,
    MIN((data->>'value')::float) as min_value,
    MAX((data->>'value')::float) as max_value,
    COUNT(*) as reading_count
FROM events
WHERE subject LIKE 'events.biometric.%'
GROUP BY day, subject
WITH NO DATA;

SELECT add_continuous_aggregate_policy('daily_biometrics',
    start_offset => INTERVAL '3 days',
    end_offset => INTERVAL '1 hour',
    schedule_interval => INTERVAL '1 hour');

-- ============================================
-- HELPER FUNCTIONS
-- ============================================

-- Get events for a date range
CREATE OR REPLACE FUNCTION get_events(
    start_time TIMESTAMPTZ,
    end_time TIMESTAMPTZ,
    subject_pattern TEXT DEFAULT '%'
)
RETURNS TABLE (
    event_timestamp TIMESTAMPTZ,
    event_subject TEXT,
    event_data JSONB
)
LANGUAGE SQL
AS $$
    SELECT timestamp, subject, data
    FROM events
    WHERE
        timestamp BETWEEN start_time AND end_time
        AND subject LIKE subject_pattern
    ORDER BY timestamp DESC;
$$;

-- Get daily summary for a subject
CREATE OR REPLACE FUNCTION get_daily_summary(
    subject_pattern TEXT,
    days_back INTEGER DEFAULT 30
)
RETURNS TABLE (
    day DATE,
    total_events BIGINT
)
LANGUAGE SQL
AS $$
    SELECT
        day::DATE,
        event_count as total_events
    FROM daily_event_counts
    WHERE
        subject LIKE subject_pattern
        AND day >= NOW() - (days_back || ' days')::INTERVAL
    ORDER BY day DESC;
$$;

-- ============================================
-- PERMISSIONS
-- ============================================
GRANT ALL ON events TO onemind;
GRANT ALL ON daily_event_counts TO onemind;
GRANT ALL ON daily_biometrics TO onemind;

-- =============================================================================
-- PGVECTOR TABLES — Semantic Memory Layer
-- =============================================================================

-- ============================================
-- EMBEDDINGS: Agent long-term memory
-- ============================================
-- Stores vector embeddings for semantic recall.
-- OpenClaw has its own SQLite memory, but this gives
-- ALL agents (Legacy, HP, LE, GE) a shared semantic memory
-- that persists across restarts and is queryable via SQL.
-- ============================================
CREATE TABLE embeddings (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    agent_id TEXT NOT NULL,             -- 'legacy', 'hp', 'le', 'ge', 'system'
    content TEXT NOT NULL,              -- The original text content
    embedding vector(1536),             -- OpenAI text-embedding-3-small (1536d)
    metadata JSONB DEFAULT '{}',        -- Arbitrary metadata (source, tags, etc.)
    domain TEXT DEFAULT 'general'       -- 'ui', 'hp', 'le', 'ge', 'system'
);

-- Cosine similarity index (IVFFlat for <100K rows, switch to HNSW at scale)
CREATE INDEX idx_embeddings_vector ON embeddings
    USING ivfflat (embedding vector_cosine_ops)
    WITH (lists = 100);

-- Agent + domain filtering
CREATE INDEX idx_embeddings_agent ON embeddings (agent_id, domain, created_at DESC);

-- Metadata GIN index for JSONB queries
CREATE INDEX idx_embeddings_metadata ON embeddings USING GIN (metadata);

-- ============================================
-- KNOWLEDGE CHUNKS: Codex document embeddings
-- ============================================
-- Embedded chunks from Codex vault markdown files
-- for RAG (Retrieval-Augmented Generation).
-- OpenClaw's native hybrid search covers workspace files,
-- but this enables cross-agent shared knowledge and
-- advanced semantic queries via SQL.
-- ============================================
CREATE TABLE knowledge_chunks (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    file_path TEXT NOT NULL,            -- e.g. '00-24 UI/02 Agents/LEGACY.md'
    chunk_index INTEGER NOT NULL,       -- Position within the file
    content TEXT NOT NULL,              -- The text chunk
    embedding vector(1536),             -- Embedding of the chunk
    metadata JSONB DEFAULT '{}',        -- File metadata, headings, tags
    domain TEXT DEFAULT 'general',      -- Which OneMind domain
    UNIQUE (file_path, chunk_index)
);

CREATE INDEX idx_knowledge_vector ON knowledge_chunks
    USING ivfflat (embedding vector_cosine_ops)
    WITH (lists = 100);

CREATE INDEX idx_knowledge_file ON knowledge_chunks (file_path);
CREATE INDEX idx_knowledge_domain ON knowledge_chunks (domain, updated_at DESC);

-- ============================================
-- CORTEX STATE: Agent runtime state snapshots
-- ============================================
-- Inspired by OM1's cortex tick loop.
-- Captures agent state at each heartbeat for
-- debugging, replay, and pattern analysis.
-- ============================================
CREATE TABLE cortex_ticks (
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    agent_id TEXT NOT NULL,             -- 'legacy', 'hp', 'le', 'ge'
    tick_num BIGINT NOT NULL,           -- Monotonic tick counter
    mode TEXT DEFAULT 'default',        -- Current operating mode
    inputs JSONB DEFAULT '{}',          -- Fused inputs at this tick
    actions JSONB DEFAULT '[]',         -- Actions decided this tick
    latency_ms FLOAT,                   -- LLM response time
    metadata JSONB DEFAULT '{}',        -- Extra context
    PRIMARY KEY (timestamp, agent_id)
);

SELECT create_hypertable('cortex_ticks', 'timestamp', chunk_time_interval => INTERVAL '1 day');

CREATE INDEX idx_cortex_agent ON cortex_ticks (agent_id, timestamp DESC);
CREATE INDEX idx_cortex_mode ON cortex_ticks (mode, timestamp DESC);

ALTER TABLE cortex_ticks SET (
    timescaledb.compress,
    timescaledb.compress_segmentby = 'agent_id'
);

SELECT add_compression_policy('cortex_ticks', INTERVAL '7 days');

-- ============================================
-- SEMANTIC SEARCH FUNCTIONS
-- ============================================

-- Find similar embeddings by vector cosine distance
CREATE OR REPLACE FUNCTION semantic_search(
    query_embedding vector(1536),
    match_count INTEGER DEFAULT 10,
    match_threshold FLOAT DEFAULT 0.7,
    filter_agent TEXT DEFAULT NULL,
    filter_domain TEXT DEFAULT NULL
)
RETURNS TABLE (
    id BIGINT,
    content TEXT,
    agent_id TEXT,
    domain TEXT,
    similarity FLOAT,
    metadata JSONB,
    created_at TIMESTAMPTZ
)
LANGUAGE SQL
AS $$
    SELECT
        e.id,
        e.content,
        e.agent_id,
        e.domain,
        1 - (e.embedding <=> query_embedding) as similarity,
        e.metadata,
        e.created_at
    FROM embeddings e
    WHERE
        (filter_agent IS NULL OR e.agent_id = filter_agent)
        AND (filter_domain IS NULL OR e.domain = filter_domain)
        AND 1 - (e.embedding <=> query_embedding) > match_threshold
    ORDER BY e.embedding <=> query_embedding
    LIMIT match_count;
$$;

-- Find relevant knowledge chunks for RAG
CREATE OR REPLACE FUNCTION knowledge_search(
    query_embedding vector(1536),
    match_count INTEGER DEFAULT 5,
    filter_domain TEXT DEFAULT NULL
)
RETURNS TABLE (
    file_path TEXT,
    chunk_index INTEGER,
    content TEXT,
    similarity FLOAT,
    metadata JSONB
)
LANGUAGE SQL
AS $$
    SELECT
        k.file_path,
        k.chunk_index,
        k.content,
        1 - (k.embedding <=> query_embedding) as similarity,
        k.metadata
    FROM knowledge_chunks k
    WHERE
        (filter_domain IS NULL OR k.domain = filter_domain)
    ORDER BY k.embedding <=> query_embedding
    LIMIT match_count;
$$;

-- ============================================
-- PERMISSIONS (pgvector tables)
-- ============================================
GRANT ALL ON embeddings TO onemind;
GRANT ALL ON knowledge_chunks TO onemind;
GRANT ALL ON cortex_ticks TO onemind;
GRANT ALL ON SEQUENCE embeddings_id_seq TO onemind;
GRANT ALL ON SEQUENCE knowledge_chunks_id_seq TO onemind;

-- Success message
DO $$
BEGIN
    RAISE NOTICE 'OmniDB schema initialized: PostgreSQL 16 + TimescaleDB + pgvector';
    RAISE NOTICE 'Tables: events, embeddings, knowledge_chunks, cortex_ticks';
    RAISE NOTICE 'Functions: get_events, get_daily_summary, semantic_search, knowledge_search';
END $$;

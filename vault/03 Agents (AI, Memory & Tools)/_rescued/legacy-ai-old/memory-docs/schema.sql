-- ============================================================
-- Legacy AI - Phase 2 Supabase Schema
-- Real-time conversation logging and analytics
-- ============================================================

-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ============================================================
-- Sessions Table
-- Tracks LiveKit rooms/sessions
-- ============================================================

CREATE TABLE IF NOT EXISTS sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    session_id TEXT UNIQUE NOT NULL,  -- LiveKit room name
    user_id TEXT,  -- User identifier
    started_at TIMESTAMPTZ DEFAULT NOW(),
    ended_at TIMESTAMPTZ,
    mode TEXT CHECK (mode IN ('online', 'offline', 'auto')),
    metadata JSONB DEFAULT '{}'::jsonb,
    
    -- Indexes
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_sessions_session_id ON sessions(session_id);
CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_started_at ON sessions(started_at DESC);

-- ============================================================
-- Conversation Messages Table
-- All conversation messages with metadata
-- ============================================================

CREATE TABLE IF NOT EXISTS conversation_messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    session_id TEXT NOT NULL REFERENCES sessions(session_id) ON DELETE CASCADE,
    role TEXT NOT NULL CHECK (role IN ('user', 'assistant', 'system')),
    content TEXT NOT NULL,
    
    -- Performance metadata
    metadata JSONB DEFAULT '{}'::jsonb,  -- latency, model, tokens, etc.
    
    -- Timestamps
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_messages_session_id ON conversation_messages(session_id);
CREATE INDEX idx_messages_created_at ON conversation_messages(created_at DESC);
CREATE INDEX idx_messages_role ON conversation_messages(role);

-- GIN index for metadata queries
CREATE INDEX idx_messages_metadata ON conversation_messages USING GIN (metadata);

-- ============================================================
-- Tool Executions Table
-- Analytics for MCP tool usage
-- ============================================================

CREATE TABLE IF NOT EXISTS tool_executions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    session_id TEXT NOT NULL REFERENCES sessions(session_id) ON DELETE CASCADE,
    tool_name TEXT NOT NULL,
    
    -- Execution details
    inputs JSONB NOT NULL,
    outputs JSONB,
    success BOOLEAN NOT NULL,
    error_message TEXT,
    duration_ms NUMERIC(10, 2),
    
    -- Timestamps
    executed_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_tools_session_id ON tool_executions(session_id);
CREATE INDEX idx_tools_tool_name ON tool_executions(tool_name);
CREATE INDEX idx_tools_executed_at ON tool_executions(executed_at DESC);
CREATE INDEX idx_tools_success ON tool_executions(success);

-- ============================================================
-- Mode Transitions Table
-- Track online ↔ offline transitions
-- ============================================================

CREATE TABLE IF NOT EXISTS mode_transitions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    session_id TEXT NOT NULL REFERENCES sessions(session_id) ON DELETE CASCADE,
    from_mode TEXT NOT NULL CHECK (from_mode IN ('online', 'offline')),
    to_mode TEXT NOT NULL CHECK (to_mode IN ('online', 'offline')),
    reason TEXT,  -- "connectivity_lost", "manual", "auto_detect"
    
    -- Timestamps
    transitioned_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_transitions_session_id ON mode_transitions(session_id);
CREATE INDEX idx_transitions_transitioned_at ON mode_transitions(transitioned_at DESC);

-- ============================================================
-- Row Level Security (RLS)
-- ============================================================

-- Enable RLS on all tables
ALTER TABLE sessions ENABLE ROW LEVEL SECURITY;
ALTER TABLE conversation_messages ENABLE ROW LEVEL SECURITY;
ALTER TABLE tool_executions ENABLE ROW LEVEL SECURITY;
ALTER TABLE mode_transitions ENABLE ROW LEVEL SECURITY;

-- Policy: Service role can do anything (for server-side)
CREATE POLICY "Service role full access" ON sessions
    FOR ALL USING (auth.role() = 'service_role');

CREATE POLICY "Service role full access" ON conversation_messages
    FOR ALL USING (auth.role() = 'service_role');

CREATE POLICY "Service role full access" ON tool_executions
    FOR ALL USING (auth.role() = 'service_role');

CREATE POLICY "Service role full access" ON mode_transitions
    FOR ALL USING (auth.role() = 'service_role');

-- Policy: Authenticated users can only see their own data
CREATE POLICY "Users see own sessions" ON sessions
    FOR SELECT USING (auth.uid()::text = user_id);

CREATE POLICY "Users see own messages" ON conversation_messages
    FOR SELECT USING (
        session_id IN (SELECT session_id FROM sessions WHERE user_id = auth.uid()::text)
    );

CREATE POLICY "Users see own tool executions" ON tool_executions
    FOR SELECT USING (
        session_id IN (SELECT session_id FROM sessions WHERE user_id = auth.uid()::text)
    );

CREATE POLICY "Users see own transitions" ON mode_transitions
    FOR SELECT USING (
        session_id IN (SELECT session_id FROM sessions WHERE user_id = auth.uid()::text)
    );

-- ============================================================
-- Realtime Publication
-- Enable realtime subscriptions for live updates
-- ============================================================

ALTER PUBLICATION supabase_realtime ADD TABLE conversation_messages;
ALTER PUBLICATION supabase_realtime ADD TABLE tool_executions;
ALTER PUBLICATION supabase_realtime ADD TABLE mode_transitions;

-- ============================================================
-- Analytics Views
-- ============================================================

-- View: Session statistics
CREATE OR REPLACE VIEW session_stats AS
SELECT 
    s.session_id,
    s.user_id,
    s.started_at,
    s.ended_at,
    EXTRACT(EPOCH FROM (s.ended_at - s.started_at)) AS duration_seconds,
    COUNT(DISTINCT m.id) AS message_count,
    COUNT(DISTINCT t.id) AS tool_execution_count,
    AVG((m.metadata->>'latency_ms')::numeric) AS avg_latency_ms
FROM sessions s
LEFT JOIN conversation_messages m ON s.session_id = m.session_id
LEFT JOIN tool_executions t ON s.session_id = t.session_id
GROUP BY s.id, s.session_id, s.user_id, s.started_at, s.ended_at;

-- View: Tool usage statistics
CREATE OR REPLACE VIEW tool_usage_stats AS
SELECT 
    tool_name,
    COUNT(*) AS total_executions,
    COUNT(*) FILTER (WHERE success = true) AS successful_executions,
    COUNT(*) FILTER (WHERE success = false) AS failed_executions,
    ROUND(100.0 * COUNT(*) FILTER (WHERE success = true) / COUNT(*), 2) AS success_rate,
    AVG(duration_ms) AS avg_duration_ms,
    MAX(duration_ms) AS max_duration_ms
FROM tool_executions
GROUP BY tool_name
ORDER BY total_executions DESC;

-- ============================================================
-- Functions
-- ============================================================

-- Function: End session
CREATE OR REPLACE FUNCTION end_session(p_session_id TEXT)
RETURNS VOID AS $$
BEGIN
    UPDATE sessions
    SET ended_at = NOW()
    WHERE session_id = p_session_id
      AND ended_at IS NULL;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- Function: Log mode transition
CREATE OR REPLACE FUNCTION log_mode_transition(
    p_session_id TEXT,
    p_from_mode TEXT,
    p_to_mode TEXT,
    p_reason TEXT DEFAULT NULL
)
RETURNS UUID AS $$
DECLARE
    v_transition_id UUID;
BEGIN
    INSERT INTO mode_transitions (session_id, from_mode, to_mode, reason)
    VALUES (p_session_id, p_from_mode, p_to_mode, p_reason)
    RETURNING id INTO v_transition_id;
    
    -- Update session mode
    UPDATE sessions
    SET mode = p_to_mode
    WHERE session_id = p_session_id;
    
    RETURN v_transition_id;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- ============================================================
-- Sample Data (Development Only)
-- ============================================================

-- Uncomment for local testing
-- INSERT INTO sessions (session_id, user_id, mode) VALUES
--     ('dev-session-1', 'user-123', 'online');

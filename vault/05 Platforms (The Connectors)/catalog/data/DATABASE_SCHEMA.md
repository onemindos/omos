# OneMind OS v2 - Database Schema Documentation

## Table of Contents

1. [Overview](#overview)
2. [Core Tables](#core-tables)
3. [AgentOS Native Tables](#agentos-native-tables)
4. [Indexes](#indexes)
5. [Relationships](#relationships)
6. [Data Types](#data-types)
7. [Seed Data](#seed-data)
8. [Triggers](#triggers)
9. [Migrations](#migrations)
10. [Usage Examples](#usage-examples)

---

## Overview

### Database Information

- **Database Name**: `onemind_os`
- **RDBMS**: PostgreSQL 16+
- **Extensions**:
  - `pgvector` - Vector similarity search for embeddings
  - `timescaledb` (planned) - Time-series data for event streaming
- **Connection URL**: `postgresql+psycopg://onemind:onemind_dev_password@localhost:5432/onemind_os`
- **Version**: v2.0
- **Schema File**: `/infra/db/init.sql`

### Architecture Highlights

OneMind OS v2 uses a **hybrid database architecture**:

1. **Custom Configuration Tables**: Agent/team configs, knowledge bases, memories
2. **AgentOS Native Tables**: Sessions, runs, evaluations (auto-managed by framework)
3. **Vector Storage**: pgvector for semantic search on documents and memories
4. **Event Sourcing**: Regular events table (will upgrade to TimescaleDB hypertable)

### Purpose

The database serves as:
- **Configuration Store**: Persistent agent and team configurations
- **Knowledge Repository**: RAG (Retrieval Augmented Generation) document storage
- **Memory System**: Long-term agent memories with semantic search
- **Session Archive**: Conversation history and summaries
- **Event Log**: System events and analytics

---

## Core Tables

### 1. agent_configs

Stores agent configurations with 40+ parameters for fine-grained control over agent behavior.

#### Schema

```sql
CREATE TABLE agent_configs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    role VARCHAR(255),

    -- Model configuration
    model_id VARCHAR(255) NOT NULL,
    model_name VARCHAR(255) NOT NULL,

    -- Instructions
    instructions TEXT[],
    markdown BOOLEAN DEFAULT TRUE,

    -- Tools
    tools TEXT[],

    -- Knowledge & RAG
    knowledge_bases TEXT[],
    add_knowledge_to_context BOOLEAN DEFAULT FALSE,

    -- Memory
    update_memory_on_run BOOLEAN DEFAULT FALSE,
    enable_agentic_memory BOOLEAN DEFAULT FALSE,
    add_memories_to_context BOOLEAN DEFAULT FALSE,

    -- History & Context
    add_history_to_context BOOLEAN DEFAULT TRUE,
    num_history_runs INTEGER DEFAULT 5,
    num_history_messages INTEGER,
    add_datetime_to_context BOOLEAN DEFAULT TRUE,

    -- Session Management
    enable_session_summaries BOOLEAN DEFAULT FALSE,
    add_session_summary_to_context BOOLEAN DEFAULT FALSE,
    enable_agentic_state BOOLEAN DEFAULT FALSE,

    -- Advanced
    tool_call_limit INTEGER,
    reasoning_steps INTEGER,

    -- Metadata
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
```

#### Column Descriptions

| Column | Type | Description |
|--------|------|-------------|
| `id` | UUID | Primary key, auto-generated |
| `name` | VARCHAR(255) | Unique agent name (required) |
| `description` | TEXT | Human-readable agent description |
| `role` | VARCHAR(255) | Agent role (e.g., "Research Specialist") |
| `model_id` | VARCHAR(255) | Model identifier (e.g., `us.anthropic.claude-sonnet-4-5-20250929-v1:0`) |
| `model_name` | VARCHAR(255) | Human-readable model name (e.g., "Claude Sonnet 4.5") |
| `instructions` | TEXT[] | Array of instruction strings for agent behavior |
| `markdown` | BOOLEAN | Whether to render responses in markdown (default: true) |
| `tools` | TEXT[] | Array of tool names (e.g., `['DuckDuckGoTools', 'FileTools']`) |
| `knowledge_bases` | TEXT[] | Array of knowledge base IDs (UUIDs as strings) |
| `add_knowledge_to_context` | BOOLEAN | Include KB content in context (RAG) |
| `update_memory_on_run` | BOOLEAN | Save new memories after each run |
| `enable_agentic_memory` | BOOLEAN | Use semantic memory retrieval |
| `add_memories_to_context` | BOOLEAN | Include memories in agent context |
| `add_history_to_context` | BOOLEAN | Include conversation history (default: true) |
| `num_history_runs` | INTEGER | Number of previous runs to include (default: 5) |
| `num_history_messages` | INTEGER | Number of messages to include (alternative to runs) |
| `add_datetime_to_context` | BOOLEAN | Add current date/time to context (default: true) |
| `enable_session_summaries` | BOOLEAN | Generate session summaries |
| `add_session_summary_to_context` | BOOLEAN | Include session summary in context |
| `enable_agentic_state` | BOOLEAN | Maintain persistent state across sessions |
| `tool_call_limit` | INTEGER | Max tool calls per run (optional limit) |
| `reasoning_steps` | INTEGER | Number of reasoning steps (for advanced models) |
| `created_at` | TIMESTAMPTZ | Record creation timestamp |
| `updated_at` | TIMESTAMPTZ | Last update timestamp (auto-updated via trigger) |

#### Indexes

- Primary key index on `id`
- Unique index on `name`

#### Example

```sql
INSERT INTO agent_configs (
    name, description, role, model_id, model_name,
    instructions, tools, add_history_to_context, num_history_runs
) VALUES (
    'Assistant',
    'Helpful assistant for everyday tasks',
    'General Purpose AI',
    'us.anthropic.claude-sonnet-4-5-20250929-v1:0',
    'Claude Sonnet 4.5',
    ARRAY['You are a helpful AI assistant for OneMind OS v2.'],
    ARRAY['DuckDuckGoTools', 'FileTools'],
    TRUE,
    5
);
```

---

### 2. team_configs

Stores team configurations that coordinate multiple agents.

#### Schema

```sql
CREATE TABLE team_configs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,

    -- Model configuration (for team lead)
    model_id VARCHAR(255) NOT NULL,
    model_name VARCHAR(255) NOT NULL,

    -- Team members (references to agent_configs)
    member_ids UUID[] NOT NULL,

    -- AgentOS Team coordination flags
    respond_directly BOOLEAN DEFAULT FALSE,
    determine_input_for_members BOOLEAN DEFAULT TRUE,
    delegate_to_all_members BOOLEAN DEFAULT FALSE,

    -- Instructions
    instructions TEXT[],
    markdown BOOLEAN DEFAULT TRUE,

    -- History
    add_history_to_context BOOLEAN DEFAULT TRUE,

    -- Metadata
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
```

#### Column Descriptions

| Column | Type | Description |
|--------|------|-------------|
| `id` | UUID | Primary key, auto-generated |
| `name` | VARCHAR(255) | Unique team name (required) |
| `description` | TEXT | Team purpose and capabilities |
| `model_id` | VARCHAR(255) | Model for team lead orchestration |
| `model_name` | VARCHAR(255) | Human-readable model name |
| `member_ids` | UUID[] | Array of agent_config UUIDs |
| `respond_directly` | BOOLEAN | Team lead responds directly (vs delegating) |
| `determine_input_for_members` | BOOLEAN | Lead determines member inputs (default: true) |
| `delegate_to_all_members` | BOOLEAN | Delegate tasks to all members simultaneously |
| `instructions` | TEXT[] | Team-level instructions |
| `markdown` | BOOLEAN | Render responses in markdown |
| `add_history_to_context` | BOOLEAN | Include conversation history |
| `created_at` | TIMESTAMPTZ | Record creation timestamp |
| `updated_at` | TIMESTAMPTZ | Last update timestamp |

#### Indexes

- Primary key index on `id`
- Unique index on `name`

#### Example

```sql
INSERT INTO team_configs (
    name, description, model_id, model_name, member_ids,
    respond_directly, determine_input_for_members
) VALUES (
    'Research Team',
    'Comprehensive research with analysis',
    'us.anthropic.claude-sonnet-4-5-20250929-v1:0',
    'Claude Sonnet 4.5',
    ARRAY[
        'a1b2c3d4-5678-90ab-cdef-123456789abc'::UUID,  -- Researcher
        'b2c3d4e5-6789-01bc-def0-234567890bcd'::UUID   -- Analyst
    ],
    FALSE,
    TRUE
);
```

---

### 3. knowledge_bases

Stores metadata for knowledge base collections used in RAG.

#### Schema

```sql
CREATE TABLE knowledge_bases (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    document_count INTEGER DEFAULT 0,

    -- Vector store configuration
    embedding_model VARCHAR(255) DEFAULT 'amazon.titan-embed-text-v2:0',
    vector_dimension INTEGER DEFAULT 1024,

    -- Metadata
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
```

#### Column Descriptions

| Column | Type | Description |
|--------|------|-------------|
| `id` | UUID | Primary key, auto-generated |
| `name` | VARCHAR(255) | Unique knowledge base name |
| `description` | TEXT | KB purpose and content description |
| `document_count` | INTEGER | Number of documents (manually tracked) |
| `embedding_model` | VARCHAR(255) | AWS Bedrock embedding model |
| `vector_dimension` | INTEGER | Embedding vector dimension (1024 for Titan v2) |
| `created_at` | TIMESTAMPTZ | Record creation timestamp |
| `updated_at` | TIMESTAMPTZ | Last update timestamp |

#### Indexes

- Primary key index on `id`
- Unique index on `name`

#### Example

```sql
INSERT INTO knowledge_bases (name, description, document_count) VALUES
('Documentation', 'Project documentation and guides', 0),
('Research Papers', 'Academic research collection', 0);
```

---

### 4. kb_documents

Stores documents with vector embeddings for RAG semantic search.

#### Schema

```sql
CREATE TABLE kb_documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    kb_id UUID NOT NULL REFERENCES knowledge_bases(id) ON DELETE CASCADE,

    -- Document content
    content TEXT NOT NULL,
    metadata JSONB,

    -- Vector embedding (using pgvector)
    embedding vector(1024),

    -- Metadata
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX kb_documents_embedding_idx
ON kb_documents USING ivfflat (embedding vector_cosine_ops)
WITH (lists = 100);
```

#### Column Descriptions

| Column | Type | Description |
|--------|------|-------------|
| `id` | UUID | Primary key, auto-generated |
| `kb_id` | UUID | Foreign key to knowledge_bases (cascading delete) |
| `content` | TEXT | Document text content (chunked) |
| `metadata` | JSONB | Flexible metadata (source, author, date, etc.) |
| `embedding` | vector(1024) | pgvector embedding for semantic search |
| `created_at` | TIMESTAMPTZ | Record creation timestamp |

#### Indexes

- Primary key index on `id`
- Foreign key index on `kb_id` (implicit)
- **IVFFlat index** on `embedding` for fast similarity search

#### Vector Search Performance

The IVFFlat index provides approximate nearest neighbor (ANN) search:
- **Index type**: `ivfflat` (Inverted File with Flat compression)
- **Distance metric**: `vector_cosine_ops` (cosine similarity)
- **Lists parameter**: 100 (number of clusters)
- **Trade-off**: Speed vs accuracy (100 lists = good balance)

#### Example

```sql
-- Insert document with embedding
INSERT INTO kb_documents (kb_id, content, metadata, embedding) VALUES (
    'kb-uuid-here',
    'OneMind OS is a personal AI command center...',
    '{"source": "README.md", "section": "Overview"}'::JSONB,
    '[0.123, -0.456, 0.789, ...]'::vector(1024)
);

-- Semantic search query
SELECT id, content, metadata,
       1 - (embedding <=> query_vector) AS similarity
FROM kb_documents
WHERE kb_id = 'kb-uuid-here'
ORDER BY embedding <=> query_vector
LIMIT 5;
```

---

### 5. agent_memories

Stores persistent agent memories with vector embeddings for semantic recall.

#### Schema

```sql
CREATE TABLE agent_memories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    agent_id UUID NOT NULL REFERENCES agent_configs(id) ON DELETE CASCADE,

    -- Memory content
    memory TEXT NOT NULL,
    memory_type VARCHAR(50) DEFAULT 'user',

    -- Vector embedding for semantic search
    embedding vector(1024),

    -- Metadata
    session_id VARCHAR(255),
    run_id VARCHAR(255),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX agent_memories_embedding_idx
ON agent_memories USING ivfflat (embedding vector_cosine_ops)
WITH (lists = 100);

CREATE INDEX agent_memories_agent_id_idx
ON agent_memories(agent_id);
```

#### Column Descriptions

| Column | Type | Description |
|--------|------|-------------|
| `id` | UUID | Primary key, auto-generated |
| `agent_id` | UUID | Foreign key to agent_configs (cascading delete) |
| `memory` | TEXT | Memory content (fact, preference, context) |
| `memory_type` | VARCHAR(50) | Memory type: `user`, `system`, `tool`, etc. |
| `embedding` | vector(1024) | pgvector embedding for semantic retrieval |
| `session_id` | VARCHAR(255) | Optional session reference |
| `run_id` | VARCHAR(255) | Optional run reference |
| `created_at` | TIMESTAMPTZ | Memory creation timestamp |

#### Indexes

- Primary key index on `id`
- IVFFlat index on `embedding` (semantic search)
- B-tree index on `agent_id` (filter by agent)

#### Example

```sql
-- Insert memory
INSERT INTO agent_memories (agent_id, memory, memory_type, embedding) VALUES (
    'agent-uuid-here',
    'User prefers concise responses without excessive detail',
    'user',
    '[0.234, -0.567, 0.890, ...]'::vector(1024)
);

-- Retrieve relevant memories
SELECT memory, memory_type, created_at,
       1 - (embedding <=> query_vector) AS relevance
FROM agent_memories
WHERE agent_id = 'agent-uuid-here'
ORDER BY embedding <=> query_vector
LIMIT 10;
```

---

### 6. session_summaries

Stores session summaries for conversation context compression.

#### Schema

```sql
CREATE TABLE session_summaries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    agent_id UUID NOT NULL REFERENCES agent_configs(id) ON DELETE CASCADE,
    session_id VARCHAR(255) NOT NULL,

    -- Summary content
    summary TEXT NOT NULL,
    message_count INTEGER DEFAULT 0,

    -- Timestamps
    session_start TIMESTAMPTZ,
    session_end TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX session_summaries_agent_id_idx
ON session_summaries(agent_id);

CREATE INDEX session_summaries_session_id_idx
ON session_summaries(session_id);
```

#### Column Descriptions

| Column | Type | Description |
|--------|------|-------------|
| `id` | UUID | Primary key, auto-generated |
| `agent_id` | UUID | Foreign key to agent_configs (cascading delete) |
| `session_id` | VARCHAR(255) | Session identifier |
| `summary` | TEXT | Generated summary of conversation |
| `message_count` | INTEGER | Number of messages in session |
| `session_start` | TIMESTAMPTZ | Session start time |
| `session_end` | TIMESTAMPTZ | Session end time |
| `created_at` | TIMESTAMPTZ | Summary generation timestamp |

#### Indexes

- Primary key index on `id`
- B-tree index on `agent_id` (filter by agent)
- B-tree index on `session_id` (lookup sessions)

#### Example

```sql
INSERT INTO session_summaries (
    agent_id, session_id, summary, message_count, session_start, session_end
) VALUES (
    'agent-uuid-here',
    'session-123',
    'User asked about database schema. Agent explained tables and relationships.',
    8,
    '2026-02-06 10:00:00'::TIMESTAMPTZ,
    '2026-02-06 10:15:00'::TIMESTAMPTZ
);
```

---

### 7. agent_states

Stores persistent agent state as JSONB for agentic behavior.

#### Schema

```sql
CREATE TABLE agent_states (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    agent_id UUID NOT NULL REFERENCES agent_configs(id) ON DELETE CASCADE,

    -- State data (stored as JSON)
    state_data JSONB NOT NULL,

    -- Metadata
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,

    UNIQUE(agent_id)
);
```

#### Column Descriptions

| Column | Type | Description |
|--------|------|-------------|
| `id` | UUID | Primary key, auto-generated |
| `agent_id` | UUID | Foreign key to agent_configs (one state per agent) |
| `state_data` | JSONB | Flexible JSON state (counters, flags, context) |
| `created_at` | TIMESTAMPTZ | State creation timestamp |
| `updated_at` | TIMESTAMPTZ | Last state update timestamp |

#### Constraints

- **UNIQUE(agent_id)**: One state per agent

#### Example

```sql
-- Insert agent state
INSERT INTO agent_states (agent_id, state_data) VALUES (
    'agent-uuid-here',
    '{
        "task_count": 42,
        "last_topic": "database schema",
        "user_preferences": {
            "verbosity": "concise",
            "code_language": "python"
        }
    }'::JSONB
);

-- Update agent state
UPDATE agent_states
SET state_data = state_data || '{"task_count": 43}'::JSONB
WHERE agent_id = 'agent-uuid-here';
```

---

### 8. events

Stores system events for logging and analytics (future: TimescaleDB hypertable).

#### Schema

```sql
CREATE TABLE events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    timestamp TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    type TEXT NOT NULL,
    source TEXT NOT NULL,
    priority INT DEFAULT 1,
    data JSONB,
    correlation_id TEXT
);

CREATE INDEX idx_events_timestamp ON events (timestamp DESC);
CREATE INDEX idx_events_type ON events (type);
CREATE INDEX idx_events_source ON events (source);
CREATE INDEX idx_events_correlation ON events (correlation_id) WHERE correlation_id IS NOT NULL;
CREATE INDEX idx_events_data ON events USING GIN (data);
```

#### Column Descriptions

| Column | Type | Description |
|--------|------|-------------|
| `id` | UUID | Primary key, auto-generated |
| `timestamp` | TIMESTAMPTZ | Event occurrence time (default: now) |
| `type` | TEXT | Event type (e.g., `agent.run`, `team.execution`) |
| `source` | TEXT | Event source (e.g., agent/team name) |
| `priority` | INT | Event priority (1=low, 5=high) |
| `data` | JSONB | Event payload (flexible structure) |
| `correlation_id` | TEXT | For tracing related events |

#### Indexes

- Primary key index on `id`
- B-tree index on `timestamp` (DESC for recent events)
- B-tree index on `type` (event filtering)
- B-tree index on `source` (source filtering)
- Partial B-tree index on `correlation_id` (tracing)
- GIN index on `data` (JSONB queries)

#### Future: TimescaleDB Hypertable

This table will be converted to a TimescaleDB hypertable for:
- Automatic time-based partitioning
- Efficient compression
- Continuous aggregates
- Data retention policies

#### Example

```sql
INSERT INTO events (type, source, priority, data, correlation_id) VALUES (
    'agent.run.started',
    'Assistant',
    2,
    '{
        "run_id": "run-123",
        "session_id": "session-456",
        "model": "claude-sonnet-4-5",
        "message": "Hello, world!"
    }'::JSONB,
    'trace-abc-def'
);
```

---

## AgentOS Native Tables

These tables are **auto-managed by AgentOS framework** and should not be modified directly:

### 1. agent_sessions

Stores conversation sessions with message history.

- **Purpose**: Session management and conversation tracking
- **Managed by**: AgentOS `/sessions` API
- **Columns**: session_id, agent_id, user_id, messages, created_at, updated_at

### 2. agent_runs

Stores agent execution logs.

- **Purpose**: Track agent invocations and responses
- **Managed by**: AgentOS `/agents/{id}/runs` API
- **Columns**: run_id, agent_id, session_id, input, output, tool_calls, created_at

### 3. agent_eval_runs

Stores agent evaluation results.

- **Purpose**: Agent performance benchmarking
- **Managed by**: AgentOS evaluation framework
- **Columns**: eval_id, agent_id, test_case, expected, actual, score, created_at

### 4. Other AgentOS Tables

- `agent_knowledge_associations` - Agent-to-knowledge-base mappings
- `agent_tool_calls` - Detailed tool execution logs
- `agent_metrics` - Performance and cost metrics

**Note**: AgentOS manages these tables automatically. Do not create foreign key constraints to these tables from custom tables, as AgentOS may modify the schema.

---

## Indexes

### Summary

OneMind OS v2 uses 13+ indexes for performance optimization:

| Table | Index Name | Type | Purpose |
|-------|-----------|------|---------|
| `kb_documents` | `kb_documents_embedding_idx` | IVFFlat | Vector similarity search |
| `agent_memories` | `agent_memories_embedding_idx` | IVFFlat | Vector similarity search |
| `agent_memories` | `agent_memories_agent_id_idx` | B-tree | Filter by agent |
| `session_summaries` | `session_summaries_agent_id_idx` | B-tree | Filter by agent |
| `session_summaries` | `session_summaries_session_id_idx` | B-tree | Lookup sessions |
| `events` | `idx_events_timestamp` | B-tree (DESC) | Recent events first |
| `events` | `idx_events_type` | B-tree | Filter by event type |
| `events` | `idx_events_source` | B-tree | Filter by source |
| `events` | `idx_events_correlation` | B-tree (partial) | Trace correlated events |
| `events` | `idx_events_data` | GIN | JSONB queries |

### Vector Indexes (IVFFlat)

**Purpose**: Fast approximate nearest neighbor (ANN) search for semantic similarity.

**Configuration**:
- **Index type**: `ivfflat` (pgvector)
- **Distance operator**: `<=>` (cosine distance)
- **Lists parameter**: 100 (number of clusters)
- **Trade-off**: 100 lists provides good balance between speed and accuracy

**Usage**:
```sql
-- Find similar documents
SELECT content, 1 - (embedding <=> query_vector) AS similarity
FROM kb_documents
ORDER BY embedding <=> query_vector
LIMIT 5;
```

**Performance Notes**:
- IVFFlat requires `lists` parameter (typically sqrt(num_rows))
- Build index AFTER inserting data (not on empty table)
- Index is approximate (faster but may miss some results)
- For exact search, omit index temporarily

### Foreign Key Indexes

Foreign keys automatically create indexes on referencing columns:
- `kb_documents.kb_id` → `knowledge_bases.id`
- `agent_memories.agent_id` → `agent_configs.id`
- `session_summaries.agent_id` → `agent_configs.id`
- `agent_states.agent_id` → `agent_configs.id`

### Unique Indexes

Created automatically for UNIQUE constraints:
- `agent_configs.name`
- `team_configs.name`
- `knowledge_bases.name`
- `agent_states.agent_id`

---

## Relationships

### Entity Relationship Diagram (ASCII)

```
┌─────────────────┐
│  agent_configs  │
│  (id, name, ...)│
└────────┬────────┘
         │
         │ 1:N (cascading delete)
         │
         ├────────────────────────────────┐
         │                                │
         ▼                                ▼
┌────────────────┐              ┌──────────────────┐
│ agent_memories │              │ session_summaries│
│ (agent_id FK)  │              │  (agent_id FK)   │
└────────────────┘              └──────────────────┘
         │                                │
         │                                │
         ▼                                ▼
┌────────────────┐              ┌──────────────────┐
│  agent_states  │              │  agent_sessions  │
│ (agent_id FK)  │              │  (AgentOS-native)│
└────────────────┘              └──────────────────┘


┌─────────────────┐
│  team_configs   │
│ (member_ids[])  │
└────────┬────────┘
         │
         │ N:M (UUID array)
         │
         ▼
┌─────────────────┐
│  agent_configs  │
│     (id)        │
└─────────────────┘


┌──────────────────┐
│ knowledge_bases  │
│  (id, name, ...)│
└────────┬─────────┘
         │
         │ 1:N (cascading delete)
         │
         ▼
┌──────────────────┐
│  kb_documents    │
│  (kb_id FK)      │
│  embedding[1024] │
└──────────────────┘
```

### Foreign Key Constraints

| Table | Column | References | On Delete |
|-------|--------|-----------|-----------|
| `kb_documents` | `kb_id` | `knowledge_bases(id)` | CASCADE |
| `agent_memories` | `agent_id` | `agent_configs(id)` | CASCADE |
| `session_summaries` | `agent_id` | `agent_configs(id)` | CASCADE |
| `agent_states` | `agent_id` | `agent_configs(id)` | CASCADE |

### Cascade Behavior

**ON DELETE CASCADE**: Deleting an agent or knowledge base automatically deletes related records:

```sql
-- Delete agent → cascades to:
DELETE FROM agent_configs WHERE id = 'agent-uuid';
-- Automatically deletes:
-- - agent_memories (all memories for this agent)
-- - session_summaries (all session summaries)
-- - agent_states (agent's persistent state)

-- Delete knowledge base → cascades to:
DELETE FROM knowledge_bases WHERE id = 'kb-uuid';
-- Automatically deletes:
-- - kb_documents (all documents in this KB)
```

### Many-to-Many Relationships

**Teams ↔ Agents**: Implemented via UUID array (not junction table)

```sql
-- team_configs.member_ids contains array of agent_configs.id
SELECT * FROM team_configs WHERE 'agent-uuid'::UUID = ANY(member_ids);
```

**Agents ↔ Knowledge Bases**: Implemented via TEXT array

```sql
-- agent_configs.knowledge_bases contains array of KB IDs (as strings)
SELECT * FROM agent_configs WHERE 'kb-uuid' = ANY(knowledge_bases);
```

---

## Data Types

### UUID

**Usage**: Primary keys for all core tables
**Format**: RFC 4122 compliant (8-4-4-4-12 hex digits)
**Generation**: `gen_random_uuid()` (PostgreSQL native)

```sql
id UUID PRIMARY KEY DEFAULT gen_random_uuid()
```

**Benefits**:
- Globally unique (no collisions)
- No sequence coordination needed
- Works well with distributed systems
- Hard to guess (security)

### TIMESTAMPTZ

**Usage**: All timestamp columns
**Format**: ISO 8601 with timezone
**Default**: `CURRENT_TIMESTAMP` (UTC)

```sql
created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
```

**Benefits**:
- Timezone-aware (stores UTC, displays local)
- Avoids daylight saving issues
- Proper temporal sorting

**Example**:
```sql
SELECT created_at AT TIME ZONE 'America/New_York' FROM agent_configs;
```

### TEXT[]

**Usage**: Array columns for instructions, tools, knowledge bases
**Format**: PostgreSQL native array type

```sql
instructions TEXT[]
tools TEXT[]
knowledge_bases TEXT[]
member_ids UUID[]
```

**Benefits**:
- No junction table needed
- Array functions available
- GIN indexing support

**Operations**:
```sql
-- Check membership
WHERE 'DuckDuckGoTools' = ANY(tools)

-- Array concatenation
UPDATE agent_configs SET tools = tools || ARRAY['NewTool']

-- Array length
SELECT array_length(tools, 1) FROM agent_configs
```

### JSONB

**Usage**: Flexible structured data (metadata, state, events)
**Format**: Binary JSON (faster than JSON)
**Indexing**: GIN indexes for queries

```sql
metadata JSONB
state_data JSONB
data JSONB
```

**Benefits**:
- Schema flexibility
- Rich query operators
- Efficient storage (binary)
- Indexable (GIN)

**Operations**:
```sql
-- Extract field
SELECT metadata->>'source' FROM kb_documents

-- Query nested field
WHERE data->'user_preferences'->>'verbosity' = 'concise'

-- Update field
UPDATE agent_states
SET state_data = jsonb_set(state_data, '{task_count}', '43')
```

### vector(1024)

**Usage**: Embeddings for semantic search (pgvector extension)
**Format**: Fixed-dimension float array
**Dimension**: 1024 (AWS Bedrock Titan Embed v2)

```sql
embedding vector(1024)
```

**Benefits**:
- Native vector operations
- Similarity search operators
- IVFFlat/HNSW indexing

**Distance Operators**:
- `<->` : Euclidean distance (L2)
- `<#>` : Negative inner product
- `<=>` : Cosine distance (most common)

**Example**:
```sql
-- Cosine similarity search
SELECT content, 1 - (embedding <=> query_vector) AS similarity
FROM kb_documents
ORDER BY embedding <=> query_vector
LIMIT 5;
```

---

## Seed Data

### Default Agents (4)

OneMind OS v2 ships with 4 pre-configured agents:

#### 1. Assistant
- **Role**: General Purpose AI
- **Model**: Claude Sonnet 4.5
- **Tools**: DuckDuckGo, FileTools
- **Description**: Helpful assistant for everyday tasks

#### 2. Researcher
- **Role**: Research Specialist
- **Model**: Claude Sonnet 4.5
- **Tools**: DuckDuckGo, Wikipedia, Arxiv
- **Description**: Deep research with academic sources
- **Special Instructions**: Always cite sources

#### 3. Analyst
- **Role**: Data Analyst
- **Model**: Claude Sonnet 4.5
- **Tools**: DuckDuckGo, Calculator
- **Description**: Data analysis and insights
- **Special Instructions**: Structured thinking, actionable recommendations

#### 4. Coder
- **Role**: Programming Specialist
- **Model**: Claude Sonnet 4.5
- **Tools**: DuckDuckGo, GitHub
- **Description**: Code writing and explanation
- **Special Instructions**: Clean code, documentation, best practices

### Default Knowledge Bases (3)

#### 1. Documentation
- **Purpose**: Project documentation and guides
- **Document Count**: 0 (empty)

#### 2. Research Papers
- **Purpose**: Academic research collection
- **Document Count**: 0 (empty)

#### 3. Code Examples
- **Purpose**: Programming examples and patterns
- **Document Count**: 0 (empty)

### Seed Data SQL

```sql
-- Insert default agents
INSERT INTO agent_configs (
    name, description, role, model_id, model_name,
    instructions, tools, add_history_to_context, num_history_runs
) VALUES
('Assistant', 'Helpful assistant for everyday tasks', 'General Purpose AI',
 'us.anthropic.claude-sonnet-4-5-20250929-v1:0', 'Claude Sonnet 4.5',
 ARRAY['You are a helpful AI assistant for OneMind OS v2.'],
 ARRAY['DuckDuckGoTools', 'FileTools'], TRUE, 5),

('Researcher', 'Deep research with academic sources', 'Research Specialist',
 'us.anthropic.claude-sonnet-4-5-20250929-v1:0', 'Claude Sonnet 4.5',
 ARRAY['You are a research specialist.', 'Always cite your sources.'],
 ARRAY['DuckDuckGoTools', 'WikipediaTools', 'ArxivTools'], TRUE, 5),

('Analyst', 'Data analysis and insights', 'Data Analyst',
 'us.anthropic.claude-sonnet-4-5-20250929-v1:0', 'Claude Sonnet 4.5',
 ARRAY['You are a data analyst.', 'Provide actionable recommendations.'],
 ARRAY['DuckDuckGoTools', 'CalculatorTools'], TRUE, 5),

('Coder', 'Code writing and explanation', 'Programming Specialist',
 'us.anthropic.claude-sonnet-4-5-20250929-v1:0', 'Claude Sonnet 4.5',
 ARRAY['You are a programming specialist.', 'Write clean, well-documented code.'],
 ARRAY['DuckDuckGoTools', 'GitHubTools'], TRUE, 5)
ON CONFLICT (name) DO NOTHING;

-- Insert default knowledge bases
INSERT INTO knowledge_bases (name, description, document_count) VALUES
('Documentation', 'Project documentation and guides', 0),
('Research Papers', 'Academic research collection', 0),
('Code Examples', 'Programming examples and patterns', 0)
ON CONFLICT (name) DO NOTHING;
```

---

## Triggers

### update_updated_at_column()

Automatically updates `updated_at` timestamp on row modification.

#### Function Definition

```sql
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';
```

#### Triggers

```sql
-- Agent configs
CREATE TRIGGER update_agent_configs_updated_at
BEFORE UPDATE ON agent_configs
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Team configs
CREATE TRIGGER update_team_configs_updated_at
BEFORE UPDATE ON team_configs
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Knowledge bases
CREATE TRIGGER update_knowledge_bases_updated_at
BEFORE UPDATE ON knowledge_bases
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Agent states
CREATE TRIGGER update_agent_states_updated_at
BEFORE UPDATE ON agent_states
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
```

#### Usage

```sql
-- Update automatically sets updated_at
UPDATE agent_configs
SET description = 'New description'
WHERE id = 'agent-uuid';
-- updated_at is now CURRENT_TIMESTAMP
```

---

## Migrations

### Strategy

OneMind OS v2 uses **SQL-based migrations** for schema changes.

#### Migration Files

Location: `/infra/db/migrations/`

```
migrations/
├── 001_initial_schema.sql      # Base schema (init.sql)
├── 002_add_agent_tags.sql      # Add tags column (example)
├── 003_add_cost_tracking.sql   # Add cost columns (example)
└── README.md                    # Migration instructions
```

#### Migration Template

```sql
-- Migration: 002_add_agent_tags.sql
-- Description: Add tags column to agent_configs for categorization
-- Date: 2026-02-06
-- Author: OneMind Team

BEGIN;

-- Add tags column
ALTER TABLE agent_configs
ADD COLUMN tags TEXT[] DEFAULT ARRAY[]::TEXT[];

-- Create GIN index for tag queries
CREATE INDEX agent_configs_tags_idx ON agent_configs USING GIN (tags);

-- Backfill default tags
UPDATE agent_configs SET tags = ARRAY['default'] WHERE tags IS NULL;

COMMIT;
```

#### Migration Checklist

Before applying migration:
- [ ] Test in development environment
- [ ] Backup production database
- [ ] Review rollback plan
- [ ] Check for breaking changes
- [ ] Update documentation

#### Rollback Template

```sql
-- Rollback: 002_add_agent_tags.sql

BEGIN;

DROP INDEX IF EXISTS agent_configs_tags_idx;
ALTER TABLE agent_configs DROP COLUMN IF EXISTS tags;

COMMIT;
```

#### Future: Migration Tool

Consider using migration tools:
- **Alembic** (SQLAlchemy-based, Python)
- **Flyway** (Java-based, SQL files)
- **Liquibase** (XML/YAML/SQL)

---

## Usage Examples

### 1. Create Agent with Knowledge Base

```sql
-- Create knowledge base
INSERT INTO knowledge_bases (name, description)
VALUES ('Python Docs', 'Official Python documentation')
RETURNING id;
-- Returns: kb-uuid-123

-- Create agent with KB
INSERT INTO agent_configs (
    name, description, model_id, model_name,
    tools, knowledge_bases, add_knowledge_to_context
) VALUES (
    'Python Expert',
    'Expert in Python programming',
    'us.anthropic.claude-sonnet-4-5-20250929-v1:0',
    'Claude Sonnet 4.5',
    ARRAY['DuckDuckGoTools', 'PythonTools'],
    ARRAY['kb-uuid-123'],
    TRUE
);
```

### 2. Add Documents to Knowledge Base

```sql
-- Insert document (assuming embedding generated elsewhere)
INSERT INTO kb_documents (kb_id, content, metadata, embedding)
VALUES (
    'kb-uuid-123',
    'Python is a high-level, interpreted programming language...',
    '{"source": "python.org", "section": "Introduction"}'::JSONB,
    generate_embedding('Python is a high-level...')  -- Custom function
);
```

### 3. Semantic Search

```sql
-- Find similar documents
WITH query AS (
    SELECT generate_embedding('What is Python?') AS embedding
)
SELECT
    d.content,
    d.metadata->>'source' AS source,
    1 - (d.embedding <=> q.embedding) AS similarity
FROM kb_documents d, query q
WHERE d.kb_id = 'kb-uuid-123'
ORDER BY d.embedding <=> q.embedding
LIMIT 5;
```

### 4. Create Team

```sql
-- Get agent IDs
WITH agents AS (
    SELECT id, name FROM agent_configs
    WHERE name IN ('Researcher', 'Analyst')
)
INSERT INTO team_configs (
    name, description, model_id, model_name, member_ids
)
SELECT
    'Research Team',
    'Comprehensive research with analysis',
    'us.anthropic.claude-sonnet-4-5-20250929-v1:0',
    'Claude Sonnet 4.5',
    array_agg(id)
FROM agents;
```

### 5. Agent Memory Management

```sql
-- Add memory
INSERT INTO agent_memories (agent_id, memory, memory_type, embedding)
VALUES (
    'agent-uuid-123',
    'User prefers Python over JavaScript',
    'user',
    generate_embedding('User prefers Python over JavaScript')
);

-- Retrieve relevant memories
SELECT memory, created_at
FROM agent_memories
WHERE agent_id = 'agent-uuid-123'
ORDER BY embedding <=> generate_embedding('programming language preference')
LIMIT 10;
```

### 6. Session Summary

```sql
-- Create session summary
INSERT INTO session_summaries (
    agent_id, session_id, summary, message_count,
    session_start, session_end
) VALUES (
    'agent-uuid-123',
    'session-456',
    'User asked about database schema. Agent explained tables, relationships, and vector indexes.',
    12,
    '2026-02-06 10:00:00+00',
    '2026-02-06 10:30:00+00'
);

-- Get recent summaries
SELECT session_id, summary, message_count, created_at
FROM session_summaries
WHERE agent_id = 'agent-uuid-123'
ORDER BY created_at DESC
LIMIT 5;
```

### 7. Event Logging

```sql
-- Log agent run
INSERT INTO events (type, source, data, correlation_id)
VALUES (
    'agent.run.completed',
    'Python Expert',
    '{
        "run_id": "run-789",
        "session_id": "session-456",
        "duration_ms": 1234,
        "token_count": 567,
        "tool_calls": 2
    }'::JSONB,
    'trace-abc-def'
);

-- Query events by type
SELECT
    timestamp,
    source,
    data->>'run_id' AS run_id,
    data->>'duration_ms' AS duration_ms
FROM events
WHERE type = 'agent.run.completed'
  AND timestamp > NOW() - INTERVAL '1 hour'
ORDER BY timestamp DESC;
```

### 8. Agent State Management

```sql
-- Initialize agent state
INSERT INTO agent_states (agent_id, state_data)
VALUES (
    'agent-uuid-123',
    '{
        "conversation_count": 0,
        "favorite_topic": null,
        "user_preferences": {}
    }'::JSONB
);

-- Update state
UPDATE agent_states
SET state_data = jsonb_set(
    state_data,
    '{conversation_count}',
    ((state_data->>'conversation_count')::int + 1)::text::jsonb
)
WHERE agent_id = 'agent-uuid-123';

-- Query state
SELECT
    state_data->>'conversation_count' AS count,
    state_data->>'favorite_topic' AS topic
FROM agent_states
WHERE agent_id = 'agent-uuid-123';
```

---

## Performance Considerations

### Vector Index Tuning

**IVFFlat `lists` parameter**:
- **Rule of thumb**: `lists = sqrt(num_rows)`
- **Small dataset** (< 1000 rows): 10-50 lists
- **Medium dataset** (1000-100K rows): 100-500 lists
- **Large dataset** (> 100K rows): 1000+ lists

**Rebuild index after bulk inserts**:
```sql
DROP INDEX kb_documents_embedding_idx;
CREATE INDEX kb_documents_embedding_idx
ON kb_documents USING ivfflat (embedding vector_cosine_ops)
WITH (lists = 500);  -- Adjust based on row count
```

### Query Optimization

**Use EXPLAIN ANALYZE**:
```sql
EXPLAIN ANALYZE
SELECT content
FROM kb_documents
ORDER BY embedding <=> query_vector
LIMIT 5;
```

**Vacuum and Analyze**:
```sql
VACUUM ANALYZE kb_documents;
VACUUM ANALYZE agent_memories;
```

### Connection Pooling

Use connection pooler (PgBouncer, SQLAlchemy pool) for:
- Reduced connection overhead
- Better resource utilization
- Handling connection spikes

### Partitioning (Future)

Consider table partitioning for:
- `events` (by timestamp)
- `agent_memories` (by agent_id)
- `kb_documents` (by kb_id)

---

## Security Considerations

### Row-Level Security (RLS)

Future enhancement: Implement RLS for multi-tenant isolation:

```sql
-- Enable RLS
ALTER TABLE agent_configs ENABLE ROW LEVEL SECURITY;

-- Create policy
CREATE POLICY agent_isolation ON agent_configs
FOR ALL
USING (owner_id = current_user_id());
```

### Sensitive Data

**Never store in database**:
- LLM API keys (use environment variables)
- User passwords (hash with bcrypt/argon2)
- PII without encryption

**Use JSONB for metadata, not PII**:
```sql
-- Good: Non-sensitive metadata
metadata: {"source": "docs.python.org", "section": "Tutorial"}

-- Bad: Sensitive data
metadata: {"user_email": "user@example.com", "ip_address": "1.2.3.4"}
```

### Prepared Statements

Always use parameterized queries to prevent SQL injection:

```python
# Good
cursor.execute(
    "SELECT * FROM agent_configs WHERE name = %s",
    (agent_name,)
)

# Bad (SQL injection risk)
cursor.execute(
    f"SELECT * FROM agent_configs WHERE name = '{agent_name}'"
)
```

---

## Backup and Recovery

### Backup Strategy

**Daily full backup**:
```bash
pg_dump -h localhost -U onemind -d onemind_os -F c -f backup_$(date +%Y%m%d).dump
```

**Continuous archiving** (WAL):
```sql
-- postgresql.conf
archive_mode = on
archive_command = 'cp %p /mnt/backups/wal/%f'
```

### Recovery

**Restore from dump**:
```bash
pg_restore -h localhost -U onemind -d onemind_os backup_20260206.dump
```

**Point-in-time recovery** (PITR):
```sql
-- recovery.conf
restore_command = 'cp /mnt/backups/wal/%f %p'
recovery_target_time = '2026-02-06 10:30:00'
```

---

## Monitoring

### Key Metrics

**Database size**:
```sql
SELECT pg_size_pretty(pg_database_size('onemind_os'));
```

**Table sizes**:
```sql
SELECT
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;
```

**Index usage**:
```sql
SELECT
    schemaname,
    tablename,
    indexname,
    idx_scan AS scans
FROM pg_stat_user_indexes
ORDER BY idx_scan DESC;
```

**Slow queries** (enable pg_stat_statements):
```sql
SELECT
    query,
    calls,
    total_time,
    mean_time
FROM pg_stat_statements
ORDER BY mean_time DESC
LIMIT 10;
```

---

## References

- **AgentOS Documentation**: https://docs.agno.com/agentos
- **pgvector Guide**: https://github.com/pgvector/pgvector
- **PostgreSQL Docs**: https://www.postgresql.org/docs/
- **TimescaleDB**: https://docs.timescale.com/

---

**Document Version**: 1.0
**Last Updated**: 2026-02-06
**Author**: OneMind Team
**Status**: Complete



---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]

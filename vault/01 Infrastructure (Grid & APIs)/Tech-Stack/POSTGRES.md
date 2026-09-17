# PostgreSQL with pgvector

> Relational database with vector similarity search for embeddings and AI applications.

---

## Overview

| Attribute         | Value                                                   |
| ----------------- | ------------------------------------------------------- |
| **Type**          | Relational Database + Vector Store                      |
| **Default Port**  | 5432                                                    |
| **Extensions**    | pgvector, pg_trgm, uuid-ossp                            |
| **License**       | PostgreSQL License (permissive)                         |
| **Documentation** | [postgresql.org/docs](https://www.postgresql.org/docs/) |

---

## What is PostgreSQL + pgvector?

**PostgreSQL** is a powerful, open-source relational database known for:

- ACID compliance
- Complex queries and joins
- JSON/JSONB support
- Extensibility

**pgvector** is an extension that adds:

- Vector data type (embeddings)
- Similarity search (cosine, L2, inner product)
- HNSW and IVFFlat indexes for fast search

---

## OneMind Usage

### Primary Use Cases

```
┌─────────────────────────────────────────────────────────────────┐
│                     OneMind PostgreSQL                           │
│                                                                  │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │                    Core Data                              │   │
│  │  Users • Sessions • Agents • Tasks • Events • Logs       │   │
│  └──────────────────────────────────────────────────────────┘   │
│                                                                  │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │                    Vector Search (pgvector)              │   │
│  │  Knowledge embeddings • Semantic search • RAG retrieval  │   │
│  └──────────────────────────────────────────────────────────┘   │
│                                                                  │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │                    Memory Storage                         │   │
│  │  Letta memory blocks • Conversation history              │   │
│  └──────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
```

| Use Case           | Description                           |
| ------------------ | ------------------------------------- |
| **Core Data**      | Users, sessions, agent configs, tasks |
| **Vector Search**  | Knowledge base semantic search        |
| **RAG Retrieval**  | Document chunks with embeddings       |
| **Memory Storage** | Letta persistent memory               |
| **Event Log**      | Audit trail, consciousness events     |
| **Metrics**        | Intelligence scores, analytics        |

---

## Docker Compose

```yaml
services:
  postgres:
    image: pgvector/pgvector:pg16
    container_name: onemind-postgres
    restart: unless-stopped
    environment:
      POSTGRES_USER: onemind
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
      POSTGRES_DB: onemind
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./init.sql:/docker-entrypoint-initdb.d/init.sql:ro
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U onemind"]
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - onemind-network

volumes:
  postgres_data:
```

### Initialization Script

```sql
-- init.sql
CREATE EXTENSION IF NOT EXISTS vector;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Knowledge base with embeddings
CREATE TABLE knowledge_chunks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    content TEXT NOT NULL,
    embedding vector(1536),  -- OpenAI ada-002 dimension
    metadata JSONB DEFAULT '{}',
    source TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Create HNSW index for fast similarity search
CREATE INDEX ON knowledge_chunks
    USING hnsw (embedding vector_cosine_ops)
    WITH (m = 16, ef_construction = 64);
```

---

## pgvector Basics

### Vector Data Type

```sql
-- Create table with vector column
CREATE TABLE items (
    id SERIAL PRIMARY KEY,
    content TEXT,
    embedding vector(1536)  -- Dimension must match your embeddings
);

-- Insert vector
INSERT INTO items (content, embedding)
VALUES ('Hello world', '[0.1, 0.2, 0.3, ...]');
```

### Similarity Search

```sql
-- Cosine similarity (most common for text embeddings)
SELECT id, content, 1 - (embedding <=> $1) AS similarity
FROM items
ORDER BY embedding <=> $1
LIMIT 10;

-- L2 distance (Euclidean)
SELECT id, content, embedding <-> $1 AS distance
FROM items
ORDER BY embedding <-> $1
LIMIT 10;

-- Inner product (for normalized vectors)
SELECT id, content, embedding <#> $1 AS neg_inner_product
FROM items
ORDER BY embedding <#> $1
LIMIT 10;
```

### Operators

| Operator | Distance Function       |
| -------- | ----------------------- |
| `<->`    | L2 (Euclidean) distance |
| `<=>`    | Cosine distance         |
| `<#>`    | Negative inner product  |

---

## Index Types

### HNSW (Recommended)

Hierarchical Navigable Small World — fast, accurate, builds incrementally.

```sql
-- Create HNSW index
CREATE INDEX ON items
    USING hnsw (embedding vector_cosine_ops)
    WITH (m = 16, ef_construction = 64);

-- Tune search accuracy vs speed
SET hnsw.ef_search = 100;  -- Higher = more accurate, slower
```

### IVFFlat

Inverted File Index — faster to build, requires more tuning.

```sql
-- Create IVFFlat index (after data is loaded)
CREATE INDEX ON items
    USING ivfflat (embedding vector_cosine_ops)
    WITH (lists = 100);

-- Set probes for search
SET ivfflat.probes = 10;  -- Higher = more accurate, slower
```

### Index Selection

| Factor              | HNSW    | IVFFlat             |
| ------------------- | ------- | ------------------- |
| Build time          | Slower  | Faster              |
| Query speed         | Faster  | Slower              |
| Memory              | Higher  | Lower               |
| Incremental inserts | ✅ Good | ❌ Requires rebuild |
| Accuracy            | Higher  | Lower (tunable)     |

---

## OneMind Schema

### Core Tables

```sql
-- Users and sessions
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    email TEXT UNIQUE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id),
    agent_id TEXT NOT NULL,
    started_at TIMESTAMPTZ DEFAULT NOW(),
    ended_at TIMESTAMPTZ,
    metadata JSONB DEFAULT '{}'
);

-- Agent configurations
CREATE TABLE agents (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    model TEXT NOT NULL,
    instructions TEXT[],
    tools TEXT[],
    config JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Consciousness events
CREATE TABLE consciousness_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_type TEXT NOT NULL,
    source TEXT NOT NULL,
    data JSONB NOT NULL,
    awareness_level TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX ON consciousness_events (event_type, created_at DESC);
CREATE INDEX ON consciousness_events USING gin (data);
```

### Knowledge Base

```sql
-- Document sources
CREATE TABLE knowledge_sources (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    source_type TEXT NOT NULL,  -- 'file', 'url', 'manual'
    path TEXT,
    metadata JSONB DEFAULT '{}',
    last_synced TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Chunked content with embeddings
CREATE TABLE knowledge_chunks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    source_id UUID REFERENCES knowledge_sources(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    embedding vector(1536),
    chunk_index INTEGER,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX ON knowledge_chunks
    USING hnsw (embedding vector_cosine_ops)
    WITH (m = 16, ef_construction = 64);

CREATE INDEX ON knowledge_chunks (source_id);
CREATE INDEX ON knowledge_chunks USING gin (to_tsvector('english', content));
```

### Memory (Letta Integration)

```sql
-- Agent memory blocks
CREATE TABLE agent_memory (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    agent_id TEXT NOT NULL,
    block_type TEXT NOT NULL,  -- 'persona', 'human', 'system'
    content TEXT NOT NULL,
    embedding vector(1536),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Archival memory
CREATE TABLE archival_memory (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    agent_id TEXT NOT NULL,
    content TEXT NOT NULL,
    embedding vector(1536),
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX ON archival_memory
    USING hnsw (embedding vector_cosine_ops);
```

---

## Python Integration

### Connection

```python
import asyncpg
from sqlalchemy.ext.asyncio import create_async_engine, AsyncSession

# asyncpg (for async operations)
async def get_pool():
    return await asyncpg.create_pool(
        "postgresql://onemind:password@localhost:5432/onemind"
    )

# SQLAlchemy async
engine = create_async_engine(
    "postgresql+asyncpg://onemind:password@localhost:5432/onemind",
    echo=True
)
```

### Vector Search

```python
import numpy as np
from openai import OpenAI

# Generate embedding
client = OpenAI()

def get_embedding(text: str) -> list[float]:
    response = client.embeddings.create(
        input=text,
        model="text-embedding-ada-002"
    )
    return response.data[0].embedding

# Search similar documents
async def search_knowledge(query: str, limit: int = 5):
    embedding = get_embedding(query)

    async with pool.acquire() as conn:
        rows = await conn.fetch("""
            SELECT id, content, 1 - (embedding <=> $1) AS similarity
            FROM knowledge_chunks
            ORDER BY embedding <=> $1
            LIMIT $2
        """, embedding, limit)

        return [dict(row) for row in rows]
```

### RAG Pattern

```python
async def rag_query(question: str):
    # 1. Search for relevant context
    chunks = await search_knowledge(question, limit=5)

    # 2. Build context
    context = "\n\n".join([
        f"[Source: {c['source']}]\n{c['content']}"
        for c in chunks
    ])

    # 3. Query LLM with context
    prompt = f"""Answer based on the following context:

{context}

Question: {question}

Answer:"""

    response = await agent.run(prompt)
    return response
```

### Bulk Insert with Embeddings

```python
async def index_documents(documents: list[dict]):
    """Batch insert documents with embeddings."""

    # Generate embeddings in batch
    texts = [doc["content"] for doc in documents]
    embeddings = client.embeddings.create(
        input=texts,
        model="text-embedding-ada-002"
    )

    # Prepare records
    records = [
        (
            doc["content"],
            embeddings.data[i].embedding,
            doc.get("metadata", {}),
            doc.get("source")
        )
        for i, doc in enumerate(documents)
    ]

    # Bulk insert
    async with pool.acquire() as conn:
        await conn.executemany("""
            INSERT INTO knowledge_chunks (content, embedding, metadata, source)
            VALUES ($1, $2, $3, $4)
        """, records)
```

---

## Hybrid Search

Combine vector similarity with full-text search:

```sql
-- Create text search index
CREATE INDEX ON knowledge_chunks
    USING gin (to_tsvector('english', content));

-- Hybrid search query
WITH semantic AS (
    SELECT id, content, 1 - (embedding <=> $1) AS semantic_score
    FROM knowledge_chunks
    ORDER BY embedding <=> $1
    LIMIT 20
),
keyword AS (
    SELECT id, content, ts_rank(to_tsvector('english', content), plainto_tsquery($2)) AS keyword_score
    FROM knowledge_chunks
    WHERE to_tsvector('english', content) @@ plainto_tsquery($2)
)
SELECT
    COALESCE(s.id, k.id) AS id,
    COALESCE(s.content, k.content) AS content,
    COALESCE(s.semantic_score, 0) * 0.7 + COALESCE(k.keyword_score, 0) * 0.3 AS combined_score
FROM semantic s
FULL OUTER JOIN keyword k ON s.id = k.id
ORDER BY combined_score DESC
LIMIT 10;
```

---

## Performance Tuning

### PostgreSQL Configuration

```ini
# postgresql.conf

# Memory
shared_buffers = 256MB
effective_cache_size = 1GB
work_mem = 64MB
maintenance_work_mem = 256MB

# WAL
wal_buffers = 16MB
max_wal_size = 2GB

# Parallel queries
max_parallel_workers_per_gather = 4
max_parallel_workers = 8

# For vector search
random_page_cost = 1.1
effective_io_concurrency = 200
```

### Index Tuning

```sql
-- Tune HNSW search
SET hnsw.ef_search = 100;  -- Default 40, higher = more accurate

-- Tune IVFFlat probes
SET ivfflat.probes = 20;   -- Default 1, higher = more accurate

-- Analyze after bulk inserts
ANALYZE knowledge_chunks;
```

### Query Optimization

```sql
-- Check query plan
EXPLAIN ANALYZE
SELECT id, content, embedding <=> $1 AS distance
FROM knowledge_chunks
ORDER BY embedding <=> $1
LIMIT 10;

-- Ensure index is used
-- Look for "Index Scan using knowledge_chunks_embedding_idx"
```

---

## Backup & Recovery

### pg_dump

```bash
# Full backup
pg_dump -h localhost -U onemind -d onemind > backup.sql

# Custom format (compressed)
pg_dump -h localhost -U onemind -Fc -d onemind > backup.dump

# Restore
pg_restore -h localhost -U onemind -d onemind backup.dump
```

### Docker Backup

```bash
# Backup
docker exec onemind-postgres pg_dump -U onemind onemind > backup.sql

# Restore
cat backup.sql | docker exec -i onemind-postgres psql -U onemind onemind
```

---

## Monitoring

### Key Queries

```sql
-- Database size
SELECT pg_size_pretty(pg_database_size('onemind'));

-- Table sizes
SELECT
    relname AS table_name,
    pg_size_pretty(pg_total_relation_size(relid)) AS total_size
FROM pg_stat_user_tables
ORDER BY pg_total_relation_size(relid) DESC;

-- Index usage
SELECT
    indexrelname AS index_name,
    idx_scan AS scans,
    pg_size_pretty(pg_relation_size(indexrelid)) AS size
FROM pg_stat_user_indexes
ORDER BY idx_scan DESC;

-- Slow queries
SELECT query, calls, mean_exec_time, total_exec_time
FROM pg_stat_statements
ORDER BY mean_exec_time DESC
LIMIT 10;
```

### Docker Health Check

```yaml
healthcheck:
  test: ["CMD-SHELL", "pg_isready -U onemind"]
  interval: 10s
  timeout: 5s
  retries: 5
```

---

## Troubleshooting

### Vector Index Not Used

```sql
-- Check index exists
\di knowledge_chunks_embedding_idx

-- Force index usage
SET enable_seqscan = off;

-- Rebuild index
REINDEX INDEX knowledge_chunks_embedding_idx;
```

### Slow Similarity Search

```sql
-- Check if HNSW parameters are set
SHOW hnsw.ef_search;

-- Increase for better accuracy
SET hnsw.ef_search = 200;

-- Consider IVFFlat for very large datasets
```

### Connection Issues

```bash
# Check PostgreSQL is running
docker logs onemind-postgres

# Test connection
psql -h localhost -U onemind -d onemind -c "SELECT 1"

# Check pg_hba.conf for auth issues
```

---

## Resources

- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [pgvector GitHub](https://github.com/pgvector/pgvector)
- [pgvector Documentation](https://github.com/pgvector/pgvector#readme)
- [asyncpg Documentation](https://magicstack.github.io/asyncpg/)

---

_Last updated: January 2026_

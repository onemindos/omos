# PostgreSQL

> **Primary Relational Database**

## Overview

PostgreSQL is the primary database for OneMind OS, storing user data, agent memory, sessions, and application state. Enhanced with pgvector for AI embeddings.

## Status

- **Active**: ✅ Yes
- **Version**: 16.x
- **Extensions**: pgvector, pg_cron
- **Last Updated**: January 25, 2026

## Use Cases

1. **User Data** - Profiles, settings
2. **Agent Memory** - Conversation history
3. **Vector Storage** - AI embeddings (pgvector)
4. **Sessions** - User sessions
5. **Audit Logs** - Activity tracking

## Key Features

- ACID compliance
- JSON/JSONB support
- Full-text search
- Vector similarity (pgvector)
- Extensible with extensions
- Robust query optimizer

## Installation

### Docker

```bash
docker run -d \
  --name postgres \
  -p 5432:5432 \
  -e POSTGRES_USER=onemind \
  -e POSTGRES_PASSWORD=secret \
  -e POSTGRES_DB=onemind \
  -v postgres-data:/var/lib/postgresql/data \
  pgvector/pgvector:pg16
```

### Docker Compose

```yaml
services:
  postgres:
    image: pgvector/pgvector:pg16
    ports:
      - "5432:5432"
    environment:
      POSTGRES_USER: onemind
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
      POSTGRES_DB: onemind
    volumes:
      - postgres-data:/var/lib/postgresql/data
      - ./init.sql:/docker-entrypoint-initdb.d/init.sql
    healthcheck:
      test: ["CMD", "pg_isready", "-U", "onemind"]
      interval: 10s
      timeout: 5s
      retries: 5

volumes:
  postgres-data:
```

## Schema Design

### Users Table

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);
```

### Agent Memory Table

```sql
CREATE TABLE agent_memory (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID NOT NULL,
    role VARCHAR(50) NOT NULL,  -- 'user', 'assistant', 'system'
    content TEXT NOT NULL,
    metadata JSONB DEFAULT '{}',
    embedding vector(1536),  -- OpenAI embedding dimension
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_memory_session ON agent_memory(session_id);
CREATE INDEX idx_memory_embedding ON agent_memory
    USING ivfflat (embedding vector_cosine_ops) WITH (lists = 100);
```

### Sessions Table

```sql
CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    title VARCHAR(255),
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
```

## Python Integration

### SQLAlchemy

```python
from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker

DATABASE_URL = "postgresql://onemind:secret@localhost:5432/onemind"

engine = create_engine(DATABASE_URL)
SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)

def get_db():
    db = SessionLocal()
    try:
        yield db
    finally:
        db.close()
```

### Async with asyncpg

```python
import asyncpg

async def create_pool():
    return await asyncpg.create_pool(
        host='localhost',
        port=5432,
        user='onemind',
        password='secret',
        database='onemind',
        min_size=5,
        max_size=20
    )

async def get_user(pool, user_id: str):
    async with pool.acquire() as conn:
        return await conn.fetchrow(
            "SELECT * FROM users WHERE id = $1",
            user_id
        )
```

## pgvector for AI

### Enable Extension

```sql
CREATE EXTENSION IF NOT EXISTS vector;
```

### Vector Operations

```python
import asyncpg
import numpy as np

async def store_embedding(pool, content: str, embedding: list[float]):
    async with pool.acquire() as conn:
        await conn.execute('''
            INSERT INTO agent_memory (content, embedding)
            VALUES ($1, $2)
        ''', content, embedding)

async def search_similar(pool, query_embedding: list[float], limit: int = 5):
    async with pool.acquire() as conn:
        return await conn.fetch('''
            SELECT content, 1 - (embedding <=> $1) as similarity
            FROM agent_memory
            ORDER BY embedding <=> $1
            LIMIT $2
        ''', query_embedding, limit)
```

### Distance Functions

| Operator | Distance Type   | Use Case           |
| -------- | --------------- | ------------------ |
| `<=>`    | Cosine distance | Text similarity    |
| `<->`    | L2 distance     | General similarity |
| `<#>`    | Inner product   | Normalized vectors |

## Full-Text Search

### Setup

```sql
-- Add search vector column
ALTER TABLE agent_memory ADD COLUMN search_vector tsvector;

-- Create trigger to auto-update
CREATE FUNCTION update_search_vector() RETURNS trigger AS $$
BEGIN
    NEW.search_vector := to_tsvector('english', NEW.content);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER search_vector_update
    BEFORE INSERT OR UPDATE ON agent_memory
    FOR EACH ROW EXECUTE FUNCTION update_search_vector();

-- Create index
CREATE INDEX idx_memory_search ON agent_memory USING gin(search_vector);
```

### Search Query

```sql
SELECT content, ts_rank(search_vector, query) as rank
FROM agent_memory, to_tsquery('english', 'meeting & calendar') query
WHERE search_vector @@ query
ORDER BY rank DESC
LIMIT 10;
```

## JSONB Operations

```sql
-- Store JSON
INSERT INTO sessions (metadata) VALUES ('{"source": "web", "tags": ["work"]}');

-- Query JSON
SELECT * FROM sessions WHERE metadata->>'source' = 'web';

-- Array contains
SELECT * FROM sessions WHERE metadata->'tags' ? 'work';

-- Update nested JSON
UPDATE sessions
SET metadata = jsonb_set(metadata, '{tags}', '["work", "important"]')
WHERE id = '...';
```

## Migrations

### Alembic Setup

```bash
pip install alembic
alembic init migrations
```

### Create Migration

```bash
alembic revision --autogenerate -m "add users table"
alembic upgrade head
```

### Migration Example

```python
def upgrade():
    op.create_table(
        'users',
        sa.Column('id', sa.UUID(), primary_key=True),
        sa.Column('email', sa.String(255), unique=True, nullable=False),
        sa.Column('created_at', sa.DateTime(timezone=True), server_default=sa.func.now()),
    )

def downgrade():
    op.drop_table('users')
```

## Backup & Recovery

### Backup

```bash
# Full backup
pg_dump -U onemind onemind > backup.sql

# Custom format (compressed)
pg_dump -U onemind -Fc onemind > backup.dump

# Only specific tables
pg_dump -U onemind -t users -t sessions onemind > partial.sql
```

### Restore

```bash
# From SQL
psql -U onemind onemind < backup.sql

# From custom format
pg_restore -U onemind -d onemind backup.dump
```

## Monitoring

### Key Queries

```sql
-- Active connections
SELECT count(*) FROM pg_stat_activity;

-- Slow queries
SELECT query, calls, mean_time, total_time
FROM pg_stat_statements
ORDER BY total_time DESC
LIMIT 10;

-- Table sizes
SELECT tablename, pg_size_pretty(pg_total_relation_size(tablename::regclass))
FROM pg_tables WHERE schemaname = 'public';

-- Index usage
SELECT indexrelname, idx_scan, idx_tup_read
FROM pg_stat_user_indexes;
```

## Best Practices

- ✅ Use connection pooling (PgBouncer or asyncpg pool)
- ✅ Create indexes for frequently queried columns
- ✅ Use JSONB instead of JSON for indexing
- ✅ Set up regular backups
- ✅ Monitor slow queries
- ✅ Use migrations for schema changes
- ✅ Enable pg_stat_statements for query analysis

## Resources

- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [pgvector Guide](https://github.com/pgvector/pgvector)
- [SQLAlchemy Documentation](https://docs.sqlalchemy.org/)

---

**Related**: [Redis](./Redis.md) | [LanceDB](./LanceDB.md)

---
title: "CLI — psql (PostgreSQL)"
type: reference
status: active
created: 2026-04-10
domain: SO/01
tags: [type/reference, domain/so, topic/cli, topic/postgresql, topic/database]
ai_summary: "psql CLI quick reference for OmniDB (PostgreSQL + TimescaleDB + pgvector) running on VPS onemind-main at port 5432."
---

# psql — PostgreSQL CLI

> OmniDB runs on VPS `onemind-main` — PostgreSQL 16 + TimescaleDB + pgvector.
> Container: `onemind-omnidb` | Port: `5432`

---

## Connect

```bash
# Connect via Docker (from VPS)
docker exec -it onemind-omnidb psql -U onemind -d onemind

# Connect from Mac via SSH tunnel
ssh -L 5432:localhost:5432 root@104.248.4.212 -N &
psql -h localhost -U onemind -d onemind

# Connect via Tailscale
psql -h 100.77.52.93 -U onemind -d onemind -p 5432
```

---

## psql Meta-Commands

```bash
\l                    # list databases
\c database_name      # connect to database
\dt                   # list tables in current db
\d table_name         # describe table (columns, types)
\di                   # list indexes
\dn                   # list schemas
\du                   # list roles/users
\x                    # toggle expanded display (good for wide rows)
\q                    # quit
\?                    # help for psql commands
\h SELECT             # help for SQL command
```

---

## Common Queries

```sql
-- Check TimescaleDB version
SELECT extversion FROM pg_extension WHERE extname = 'timescaledb';

-- Check pgvector
SELECT extversion FROM pg_extension WHERE extname = 'vector';

-- List all tables with row counts
SELECT schemaname, tablename,
       (xpath('/row/c/text()', query_to_xml(
         format('select count(*) as c from %I.%I', schemaname, tablename),
         false, true, '')))[1]::text::int AS rows
FROM pg_tables WHERE schemaname = 'public';

-- Check active connections
SELECT count(*), state FROM pg_stat_activity GROUP BY state;

-- Kill idle connections
SELECT pg_terminate_backend(pid)
FROM pg_stat_activity
WHERE state = 'idle' AND query_start < now() - interval '1 hour';
```

---

## pgvector (RAG / Embeddings)

```sql
-- Create vector extension
CREATE EXTENSION IF NOT EXISTS vector;

-- Create table with vector column (1536 dims = OpenAI text-embedding-3-small)
CREATE TABLE knowledge_chunks (
  id SERIAL PRIMARY KEY,
  content TEXT,
  embedding vector(1536),
  metadata JSONB,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Cosine similarity search (top 5)
SELECT content, 1 - (embedding <=> '[0.1,0.2,...]'::vector) AS similarity
FROM knowledge_chunks
ORDER BY embedding <=> '[0.1,0.2,...]'::vector
LIMIT 5;

-- Create HNSW index for fast search
CREATE INDEX ON knowledge_chunks USING hnsw (embedding vector_cosine_ops);
```

---

## TimescaleDB

```sql
-- Create hypertable (time-series)
SELECT create_hypertable('metrics', 'created_at');

-- Query recent data
SELECT * FROM metrics WHERE created_at > NOW() - INTERVAL '24 hours';

-- Continuous aggregates
CREATE MATERIALIZED VIEW hourly_stats
WITH (timescaledb.continuous) AS
SELECT time_bucket('1 hour', created_at) AS bucket, avg(value)
FROM metrics GROUP BY bucket;
```

---

## Backup & Restore

```bash
# Dump database
docker exec onemind-omnidb pg_dump -U onemind onemind > onemind_backup_$(date +%Y%m%d).sql

# Restore from dump
docker exec -i onemind-omnidb psql -U onemind onemind < onemind_backup.sql
```


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]

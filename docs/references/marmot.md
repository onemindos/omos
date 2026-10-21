# Marmot v2 — Leaderless Distributed SQLite

**Repo:** github.com/maxpert/marmot  
**Language:** Go  
**Role in OneMind:** Edge-deployable distributed relational DB — leaderless, MySQL-compatible, built-in RAG vector search

Leaderless distributed SQLite replication using gossip protocol with eventual consistency. Any node accepts writes. MySQL wire protocol compatible — connects with any MySQL client.

## Key features

- Leaderless: no single point of failure, any node accepts writes
- MySQL Protocol: DBeaver, mysql CLI, any MySQL client works out of the box
- CDC-based row-level replication
- Distributed transactions (Percolator-style write intents + conflict detection)
- DDL replication (schema changes cluster-wide)
- **Built-in vector search**: local IVF/PQ indexes for RAG workloads + exact rerank from base table
- Multi-database per cluster

## Why it matters for OneMind

Complement to [[ha-sqlite-cluster]] — different trade-off:
- HA: leader-based, stronger consistency, NATS JetStream replication
- Marmot: leaderless, eventual consistency, gossip protocol, MySQL wire (broader client support)

The built-in vector search makes Marmot a candidate for OneMind's agent RAG store — structured data + vector indexes in one SQLite node, replicated to all edge nodes.

## Integration point

Mission records, entity registry, task queue + vector embeddings for agent memory — all in one node, replicated via gossip to homestead fleet.

# HA — Highly Available SQLite Cluster

**Repo:** github.com/litesql/ha  
**Language:** Go  
**License:** Apache 2.0  
**Role in OneMind:** Local-first distributed database — lightweight alternative to ClickHouse for edge/homestead deployments

Distributed SQLite cluster powered by embedded NATS JetStream. Leader-based or leaderless cluster modes, full CDC, cross-shard queries, and transaction undo.

## Protocols supported

- HTTP API
- gRPC API
- Go `database/sql` driver (`github.com/litesql/go-ha`)
- JDBC driver
- MySQL Wire Protocol
- PostgreSQL Wire Protocol

## Why it matters for OneMind

- NATS JetStream is the replication backbone — fits directly into OneMind's CORE account
- Embeddable SQLite for edge nodes that can't run ClickHouse
- Read/write replicas via `go-ha` driver — local reads, NATS-replicated writes
- CDC = real-time streaming of database changes onto NATS subjects
- Phase 0 candidate: replace single-node SQLite with clustered HA for mission state persistence

## Integration point

```
omos-core NATS JetStream (CORE account)
    ↕ HA replication
SQLite at each edge node
    ↓ CDC
evt.data.{table}.changed  (NATS subject for downstream consumers)
```

## Related

- [[CLICKHOUSE-TIMEMACHINE]] — primary time-series store for telemetry
- HA is for structured relational state (mission records, entity registry, task queue)

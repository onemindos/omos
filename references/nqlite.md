# nqlite — SQLite Replicated via NATS JetStream

**Repo:** github.com/voxoco/nqlite  
**Language:** Deno/TypeScript  
**Role in OneMind:** Simplest path to distributed SQLite at the edge — NATS JetStream as replication layer

SQLite with NATS JetStream for replication and persistence. No leader election, no node-to-node awareness — nodes only need a NATS connection. Writes publish to `nqlite.push` JetStream subject; reads are local. Snapshots to NATS Object Store.

## Why it matters for OneMind

The simplest distributed relational DB for edge nodes. OneMind services already have NATS — nqlite adds persistent structured state with zero additional infrastructure:

```bash
nqlite --nats-host=nats://onemind-mac:4222 --creds=./core.creds
```

## HTTP API

```bash
# Write
POST /db/query  ["INSERT INTO tasks(id, status) VALUES('t1', 'pending')"]

# Read (local, no NATS)
GET /db/query?q=SELECT * FROM tasks WHERE status='pending'
```

## JetStream subjects used

```
nqlite.push      # SQL write replication subject
nqlite           # JetStream stream name + Object Store bucket
```

## See Also

- [[marmot]] — leaderless gossip-based SQLite, MySQL wire, built-in vector search
- [[ha-sqlite-cluster]] — HA SQLite with NATS JetStream, multi-protocol (gRPC/MySQL/PostgreSQL wire)

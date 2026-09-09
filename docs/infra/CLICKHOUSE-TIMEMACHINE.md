# CLICKHOUSE-TIMEMACHINE.md — ClickHouse Time Machine Migration: events.> → evt.>

> **Card:** P2-T | **Phase:** 2 | **Branch:** phase2/p2t-clickhouse-rebind
> **Author:** coder (kanban worker) | **Date:** 2026-06-10
> **Reviewer:** Legacy (required before merge)

---

## 0. Purpose

ClickHouse is the OneMind Time Machine — the permanent analytical record of every event
on the bus. It was wired to the legacy subject `events.>`. The NATS Doctrine (NATS-DOCTRINE.md §2)
defines the canonical subjects as `evt.>`, `tel.>`, `cmd.>`. This card rebinds ingestion
to the semantic subjects **additively** (dual-ingest, zero data loss, fully reversible).

---

## 1. Pre-migration State Capture (SHOW CREATE — captured before any change)

### onemind.events (durable MergeTree sink — NEVER DROP)

```sql
CREATE TABLE onemind.events
(
    `timestamp` DateTime64(3) DEFAULT now64(3),
    `subject` String,
    `payload` String
)
ENGINE = MergeTree
PARTITION BY toYYYYMM(timestamp)
ORDER BY (subject, timestamp)
TTL timestamp + toIntervalYear(10)
SETTINGS index_granularity = 8192
```

### onemind.nats_events (DEPRECATED — keep running, events.> consumer)

```sql
CREATE TABLE onemind.nats_events
(
    `timestamp` DateTime64(3) DEFAULT now64(3),
    `subject` String,
    `payload` String
)
ENGINE = NATS
SETTINGS nats_url = 'localhost:4222', nats_subjects = 'events.>',
         nats_format = 'JSONEachRow', nats_num_consumers = 1
```

### onemind.nats_events_queue (DEPRECATED — keep running, raw blob consumer)

```sql
CREATE TABLE onemind.nats_events_queue
(
    `message` String
)
ENGINE = NATS
SETTINGS nats_url = '127.0.0.1:4222', nats_subjects = 'events.>',
         nats_format = 'RawBLOB', nats_num_consumers = 1, nats_reconnect_wait = 5000
```

### onemind.nats_to_events (DEPRECATED — keep running, nats_events_queue → events MV)

```sql
CREATE MATERIALIZED VIEW onemind.nats_to_events TO onemind.events
(
    `timestamp` DateTime64(3),
    `subject` String,
    `payload` String
)
AS SELECT
    now64(3)  AS timestamp,
    _subject  AS subject,
    message   AS payload
FROM onemind.nats_events_queue
```

### Pre-migration row count

```
SELECT count() FROM onemind.events
→ 31265
```

---

## 2. Migration Design

### 2.1 Approach: Additive Dual-Ingest

We do NOT alter or drop any existing objects. We add new objects:

```
[NATS bus: events.>] --> nats_events_queue (RawBLOB) --> nats_to_events (MV) --> events (MergeTree)
                    --> nats_events (JSONEachRow)                                [LEGACY — keep running]

[NATS bus: evt.>  ] --> nats_evt (RawBLOB)           --> evt_to_events (MV)  --> events (MergeTree)  [bridge]
                                                      --> evt_to_events_v2 (MV)--> events_v2 (MergeTree) [NEW]
```

During transition, `events` receives from BOTH subject spaces. `events_v2` receives only
from `evt.>`. Full history can be queried as:

```sql
SELECT timestamp AS ts, subject, payload AS data FROM onemind.events
UNION ALL
SELECT ts, subject, data FROM onemind.events_v2
ORDER BY ts DESC;
```

### 2.2 Shape Decision: Two Tables

**Option A (rejected):** Keep flat 3-column shape for evt.> (timestamp, subject, payload).
- Con: throws away the rich P0-E envelope already flowing on evt.>
- Con: evolving the column set later requires schema migration on a growing archive table

**Option B (chosen):** `events` = legacy flat archive, `events_v2` = envelope-aligned.
- `events` keeps the 3-column (timestamp, subject, payload) shape unchanged — historical
  continuity, no backfill risk, existing analytics queries unbroken.
- `events_v2` uses the P0-E target shape: (v, id, ts, type, source, subject, data)
  as documented in `schemas/README.md`. This is the forward-compatible record.
- Both tables coexist with 10-year TTL. UNION ALL gives full history.
- When P0-E envelope (card P0-E) is finalized in om-schema KV, `events_v2` columns can be
  reconciled without touching `events`.

### 2.3 NATS Table Format: RawBLOB

`nats_evt` uses `RawBLOB` format (mirrors the proven `nats_events_queue` pattern) rather
than `JSONEachRow` with 3 columns. Why:

1. Envelope producers emit `{"v","id","ts","type","source","subject","data"}` — no "payload"
   key — which would cause a 3-column JSONEachRow table to silently write empty payload strings.
2. RawBLOB stores the full raw message in `message`; the MVs do extraction via JSONExtractXxx.
3. Forward-compatible: any envelope field can be extracted without changing the NATS table DDL.

### 2.4 P0-E Envelope Status

`om-schema` KV bucket exists but is empty (P0-E card not yet landed). `events_v2` uses
the target shape from `schemas/README.md`. JSONExtract fields will silently return defaults
for missing keys — no ingestion failures. A TODO comment marks where to reconcile once
P0-E lands.

---

## 3. Commands Run (in order)

All DDL applied via ClickHouse HTTP interface (port 8124), no shell secrets.
ClickHouse is local, unauthed, on non-default ports (HTTP 8124, TCP 9001).

```bash
# Create nats_evt (RawBLOB NATS engine, evt.>)
curl -s http://localhost:8124/ --data \
  "CREATE TABLE IF NOT EXISTS onemind.nats_evt (message String)
   ENGINE = NATS
   SETTINGS nats_url='localhost:4222', nats_subjects='evt.>',
            nats_format='RawBLOB', nats_num_consumers=1, nats_reconnect_wait=5000"
# → 200 OK

# Create events_v2 (envelope-aligned MergeTree)
curl -s http://localhost:8124/ --data \
  "CREATE TABLE IF NOT EXISTS onemind.events_v2 (
     v Int32 DEFAULT 1, id String DEFAULT '', ts DateTime64(3) DEFAULT now64(3),
     type String DEFAULT '', source String DEFAULT '', subject String, data String)
   ENGINE = MergeTree PARTITION BY toYYYYMM(ts) ORDER BY (subject, ts)
   TTL ts + toIntervalYear(10) SETTINGS index_granularity=8192"
# → 200 OK

# Create evt_to_events (bridge MV: nats_evt → events)
curl -s http://localhost:8124/ --data \
  "CREATE MATERIALIZED VIEW IF NOT EXISTS onemind.evt_to_events TO onemind.events AS
   SELECT now64(3) AS timestamp, _subject AS subject, message AS payload
   FROM onemind.nats_evt"
# → 200 OK

# Create evt_to_events_v2 (envelope MV: nats_evt → events_v2)
curl -s http://localhost:8124/ --data \
  "CREATE MATERIALIZED VIEW IF NOT EXISTS onemind.evt_to_events_v2 TO onemind.events_v2 AS
   SELECT JSONExtractInt(message,'v') AS v, JSONExtractString(message,'id') AS id,
          coalesce(toDateTime64OrNull(JSONExtractString(message,'ts'),3), now64(3)) AS ts,
          JSONExtractString(message,'type') AS type, JSONExtractString(message,'source') AS source,
          _subject AS subject, JSONExtractString(message,'data') AS data
   FROM onemind.nats_evt"
# → 200 OK
```

DDL is idempotent and lives in `nats/clickhouse/timemachine.sql`.

---

## 4. Verification Transcript (real output)

### 4.1 Before: historical rows intact

```
clickhouse client --port 9001 --query "SELECT count() FROM onemind.events"
→ 31265
```

### 4.2 Test publish on evt.sys.timemachine.test (full P0-E envelope shape)

```bash
nats pub evt.sys.timemachine.test \
  '{"v":1,"id":"test-001-1781117510","ts":"2026-06-10T18:51:50.000Z",
    "type":"evt.sys.timemachine.test","source":"coder-worker",
    "subject":"evt.sys.timemachine.test","data":"{\"ok\":1}"}'
# → 14:51:50 Published 180 bytes to "evt.sys.timemachine.test"
```

### 4.3 Row present in events_v2 (envelope-aligned sink)

```
clickhouse client --port 9001 --query \
  "SELECT v, id, ts, type, source, subject, data
   FROM onemind.events_v2
   WHERE subject='evt.sys.timemachine.test'
   ORDER BY ts DESC LIMIT 1 FORMAT Vertical"

Row 1:
──────
v:       1
id:      test-001-1781117510
ts:      2026-06-10 14:51:50.945
type:    evt.sys.timemachine.test
source:  coder-worker
subject: evt.sys.timemachine.test
data:    {"ok":1}
```

### 4.4 Row also present in events (legacy bridge working)

```
clickhouse client --port 9001 --query \
  "SELECT timestamp, subject, payload
   FROM onemind.events
   WHERE subject='evt.sys.timemachine.test'
   ORDER BY timestamp DESC LIMIT 1 FORMAT Vertical"

Row 1:
──────
timestamp: 2026-06-10 14:51:50.931
subject:   evt.sys.timemachine.test
payload:   {"v":1,"id":"test-001-1781117510",...}
```

### 4.5 Historical rows still present

```
clickhouse client --port 9001 --query "SELECT count() FROM onemind.events"
→ 31274   (31265 original + test messages — no rows lost)
```

### 4.6 SHOW TABLES FROM onemind (all objects present, all legacy objects intact)

```
events
events_v2
evt_to_events
evt_to_events_v2
metrics
nats_events          ← DEPRECATED, still running
nats_events_queue    ← DEPRECATED, still running
nats_evt             ← NEW
nats_to_events       ← DEPRECATED, still running
netdata_metrics
system_metrics
```

---

## 5. Rollback Procedure

This migration is additive. To roll back (remove the new objects):

```sql
-- Safe to drop — these are streaming engines or empty MergeTree with no durable rows
DROP TABLE IF EXISTS onemind.evt_to_events_v2;
DROP TABLE IF EXISTS onemind.evt_to_events;
DROP TABLE IF EXISTS onemind.events_v2;
DROP TABLE IF EXISTS onemind.nats_evt;
```

The legacy objects (`nats_events`, `nats_events_queue`, `nats_to_events`, `events`)
are unaffected. NEVER include those in a rollback.

---

## 6. Post-cutover Removal of Legacy Objects (future — not this card)

When ALL producers have cut to `evt.>` and Legacy has approved:

1. Confirm `nats_events` is no longer receiving rows (count stable for ≥7 days).
2. Run:
   ```sql
   DROP TABLE onemind.nats_to_events;
   DROP TABLE onemind.nats_events_queue;
   DROP TABLE onemind.nats_events;
   ```
3. `onemind.events` stays forever as the historical archive (10-year TTL).

---

## 7. Open TODOs

- **P0-E reconcile:** Once the P0-E envelope schema is published to `om-schema` KV,
  reconcile `events_v2` column types/names against it. Update `evt_to_events_v2` MV if
  needed.
- **tel.> and cmd.>:** The task body mentions optionally wiring `tel.>` and `cmd.>`.
  These follow the same pattern: add `nats_tel` (RawBLOB, tel.>) and `nats_cmd` (RawBLOB,
  cmd.>) NATS tables with corresponding MVs → `events_v2` (or dedicated sink tables).
  Deferred to a follow-on card.
- **Legacy deprecation tracking:** Add a monitoring alert on `nats_events` insert rate
  so Legacy can confirm when it's safe to remove.

---

## 8. Connection Reference

ClickHouse is local, unauthed (no credentials), non-default ports:
- HTTP API: `http://localhost:8124/`
- TCP client: `clickhouse client --port 9001`

NATS: `localhost:4222` (unauthenticated on local CORE account)

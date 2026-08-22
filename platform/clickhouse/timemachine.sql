-- =============================================================================
-- timemachine.sql — OneMind ClickHouse Time Machine DDL (target state)
-- =============================================================================
-- Idempotent: every statement uses IF NOT EXISTS.
-- Apply with:
--   clickhouse client --port 9001 --queries-file nats/clickhouse/timemachine.sql
--
-- Hard rules (never relax without Legacy approval):
--   - NEVER DROP onemind.events or DELETE any rows.
--   - NEVER DROP nats_events, nats_events_queue, or nats_to_events in this phase.
--     Dual-ingest stays up. Mark deprecated; remove only after all producers
--     have cut to evt.> and Legacy approves.
--   - No secrets — ClickHouse here is local/unauthed on port 9001.
--   - If any CREATE fails: stop, capture the error, leave existing objects
--     untouched, escalate.
-- =============================================================================

-- ---------------------------------------------------------------------------
-- 0. Database (idempotent)
-- ---------------------------------------------------------------------------
CREATE DATABASE IF NOT EXISTS onemind;

-- ---------------------------------------------------------------------------
-- 1. LEGACY objects (kept, dual-ingest, DEPRECATED — do not modify or drop)
-- ---------------------------------------------------------------------------
--
-- DEPRECATED: nats_events
--   Legacy NATS engine table bound to events.>. Keep running until all
--   producers cut to evt.> and Legacy approves removal.
--   Removal trigger: all producers emit on evt.*; nats_events row-count
--   stops growing for ≥7 days; Legacy sign-off.
--
-- DEPRECATED: nats_events_queue
--   Raw-blob NATS table also bound to events.>. Used by nats_to_events MV.
--   Remove at the same time as nats_events.
--
-- DEPRECATED: nats_to_events
--   Materialized view: nats_events_queue -> events (flat 3-column shape).
--   Remove at the same time as nats_events_queue.
--
-- (No DDL here — these already exist and must not be recreated or altered.)
-- ---------------------------------------------------------------------------

-- ---------------------------------------------------------------------------
-- 2. NEW: NATS-engine table bound to evt.> (RawBLOB)
-- ---------------------------------------------------------------------------
-- Uses RawBLOB format (same as nats_events_queue) so the full raw JSON blob
-- is stored in `message`. The MVs below do all field extraction, making this
-- table compatible with both flat legacy payloads AND the P0-E envelope.
--
-- DESIGN DECISION: RawBLOB chosen over JSONEachRow because:
--   (a) Envelope producers send {"v","id","ts","type","source","subject","data"} —
--       no "payload" key — which would cause JSONEachRow on a 3-column table
--       to silently write empty payload strings.
--   (b) RawBLOB + JSONExtract in MVs is the proven pattern (nats_events_queue).
--   (c) Retains full message fidelity for future schema evolution.
--
-- TODO (card P0-E): once the P0-E envelope schema is published to om-schema KV,
-- verify JSONExtract field names/types match and update accordingly.

CREATE TABLE IF NOT EXISTS onemind.nats_evt
(
    `message` String
)
ENGINE = NATS
SETTINGS
    nats_url           = 'localhost:4222',
    nats_subjects      = 'evt.>',
    nats_format        = 'RawBLOB',
    nats_num_consumers = 1,
    nats_reconnect_wait = 5000;

-- ---------------------------------------------------------------------------
-- 3. NEW: events_v2 — envelope-aligned MergeTree sink for evt.> events
-- ---------------------------------------------------------------------------
-- Column alignment follows the P0-E target shape documented in schemas/README.md:
--   v       Int32            — envelope version (default 1)
--   id      String           — UUIDv4 message id
--   ts      DateTime64(3)    — RFC3339 event timestamp (authoritative)
--   type    String           — event type (e.g. evt.task.issue.created)
--   source  String           — originating actor id
--   subject String           — NATS subject the message arrived on
--   data    String           — JSON payload blob (opaque until P0-E schemas land)
--
-- WHY A SECOND TABLE (not alter events)?
--   `events` (flat, 3-column) is the durable legacy archive. Altering its
--   schema while it holds 31K+ rows risks a backfill burden and breaks
--   existing analytics queries. `events_v2` carries the forward-compatible
--   shape; the two tables coexist and can be UNION ALL'd for full history.
--
-- TODO (card P0-E): reconcile column types against om-schema KV once the
-- envelope schema is published there; version-bump events_v2 if needed.

CREATE TABLE IF NOT EXISTS onemind.events_v2
(
    `v`       Int32          DEFAULT 1,
    `id`      String         DEFAULT '',
    `ts`      DateTime64(3)  DEFAULT now64(3),
    `type`    String         DEFAULT '',
    `source`  String         DEFAULT '',
    `subject` String,
    `data`    String
)
ENGINE = MergeTree
PARTITION BY toYYYYMM(ts)
ORDER BY (subject, ts)
TTL ts + toIntervalYear(10)
SETTINGS index_granularity = 8192;

-- ---------------------------------------------------------------------------
-- 4. NEW: nats_evt -> events (legacy bridge MV — keeps events growing)
-- ---------------------------------------------------------------------------
-- Routes every evt.> message into the legacy flat `events` table so that
-- legacy queries on `events` see both events.> and evt.> traffic during the
-- dual-ingest transition. Uses _subject (ClickHouse NATS virtual column)
-- for the canonical subject value.

CREATE MATERIALIZED VIEW IF NOT EXISTS onemind.evt_to_events
TO onemind.events
AS
SELECT
    now64(3) AS timestamp,
    _subject AS subject,
    message  AS payload
FROM onemind.nats_evt;

-- ---------------------------------------------------------------------------
-- 5. NEW: nats_evt -> events_v2 (envelope-aligned MV)
-- ---------------------------------------------------------------------------
-- Routes evt.> messages into events_v2, extracting P0-E envelope fields.
-- Producers SHOULD emit the full P0-E envelope; any missing field falls back
-- to a safe default so the row is never rejected.
--
-- TODO (card P0-E): once all producers emit the full envelope, tighten the
-- JSONExtract types and remove the DEFAULT fallbacks from the SELECT.

CREATE MATERIALIZED VIEW IF NOT EXISTS onemind.evt_to_events_v2
TO onemind.events_v2
AS
SELECT
    JSONExtractInt(message,    'v')       AS v,
    JSONExtractString(message, 'id')      AS id,
    coalesce(
        toDateTime64OrNull(
            JSONExtractString(message, 'ts'), 3
        ),
        now64(3)
    )                                     AS ts,
    JSONExtractString(message, 'type')    AS type,
    JSONExtractString(message, 'source')  AS source,
    _subject                              AS subject,
    JSONExtractString(message, 'data')    AS data
FROM onemind.nats_evt;

-- =============================================================================
-- END timemachine.sql
-- =============================================================================

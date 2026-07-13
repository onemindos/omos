-- OneMind Events Table — time machine for all NATS intel + entity events
--
-- Run:
--   kubectl exec -n clickhouse chi-onemind-onemind-0-0-0 -- \
--     clickhouse-client --user onemind --password onemind2026 \
--     --queries-file /dev/stdin < platform/clickhouse/migrate-events.sql
--
-- Or from within the pod:
--   clickhouse-client --user onemind --password onemind2026
--   :) SOURCE platform/clickhouse/migrate-events.sql

CREATE TABLE IF NOT EXISTS onemind.events
(
    timestamp    DateTime64(3) DEFAULT now64(3),
    subject      LowCardinality(String),
    entity_id    String,
    entity_type  LowCardinality(String),
    source       LowCardinality(String),
    lat          Nullable(Float64),
    lon          Nullable(Float64),
    payload      String
)
ENGINE = MergeTree
PARTITION BY (toYYYYMM(timestamp), entity_type)
ORDER BY (entity_type, entity_id, timestamp)
TTL toDateTime(timestamp) + INTERVAL 10 YEAR
SETTINGS index_granularity = 8192;

-- Useful query examples:
--
-- All weather alerts last 7 days:
--   SELECT timestamp, entity_id, subject, JSONExtractString(payload,'properties','event') as event,
--          lat, lon FROM onemind.events
--   WHERE entity_type='weather_alert' AND timestamp > now() - INTERVAL 7 DAY
--   ORDER BY timestamp DESC;
--
-- Earthquakes M5+ last 30 days near a point (lat 35, lon 139 — Tokyo):
--   SELECT timestamp, entity_id, JSONExtractFloat(payload,'properties','magnitude') as mag,
--          lat, lon FROM onemind.events
--   WHERE entity_type='earthquake'
--     AND timestamp > now() - INTERVAL 30 DAY
--     AND greatCircleDistance(lon, lat, 139.0, 35.0) < 500000
--   ORDER BY mag DESC;
--
-- All events by source last 24h:
--   SELECT source, entity_type, count() as cnt FROM onemind.events
--   WHERE timestamp > now() - INTERVAL 1 DAY
--   GROUP BY source, entity_type ORDER BY cnt DESC;
--
-- Fires near a coordinate (lat 37.5, lon -120 — California):
--   SELECT timestamp, lat, lon, JSONExtractString(payload,'properties','confidence') as conf
--   FROM onemind.events
--   WHERE entity_type='wildfire'
--     AND timestamp > now() - INTERVAL 7 DAY
--     AND greatCircleDistance(lon, lat, -120.0, 37.5) < 200000
--   ORDER BY timestamp DESC;

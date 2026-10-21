# ORP — Open Reality Protocol

**Repo:** https://github.com/shieldofsteel/orp  
**Language:** Rust (14 workspace crates)  
**License:** Apache 2.0  
**Binary:** ~45 MB, statically linked, ARM64-capable

---

## What It Is

ORP is a single Rust binary that fuses 40+ real-time protocols into one cryptographically-signed, queryable entity graph — with a built-in COP map, ORP-QL graph query language, and mTLS federation mesh. No JVM, no Postgres, no Kubernetes. `./orp start` and ingesting in 30 seconds.

The positioning: *SQLite-style single-binary simplicity, pushed up the stack to where Lattice OS / Maven OS / Palantir AIP currently live — but open and Apache 2.0.*

---

## Architecture

Single binary bundles:

- **Connectors** — 40 protocol adapters, each a Tokio task
- **Stream Processor** — dedup, change detection, entity resolution, batch insertion
- **Storage** — DuckDB (OLAP analytics) + RocksDB (dedup/DLQ event store) + in-memory graph projection (Cypher-style, refreshed every 30s)
- **ORP-QL** — query language that compiles to DuckDB SQL or graph traversal
- **API** — Axum HTTP + WebSocket, REST + real-time subscriptions
- **Frontend** — React SPA (Leaflet map, entity inspector, alert feed) served as embedded static assets
- **Security** — OIDC JWKS (RS256/ES256), ABAC, Argon2id keystore, Ed25519-signed hash-chained audit log, mTLS federation

---

## 14 Workspace Crates

| Crate | Purpose |
|-------|---------|
| orp-core | Orchestrator, config loading, component wiring, Tokio bootstrap |
| orp-connector | 40 protocol adapters |
| orp-stream | Dedup, entity resolution, batch pipeline |
| orp-storage | DuckDB + RocksDB + graph projection |
| orp-query | ORP-QL parser, planner, hybrid router |
| orp-entity | Entity model, type registry |
| orp-geospatial | Geo indexing, bounding-box queries |
| orp-tak | TAK/CoT protocol integration |
| orp-security | OIDC, ABAC, Argon2id, Ed25519, mTLS |
| orp-audit | Hash-chained audit log, at-rest encryption |
| orp-config | Config schema, templates, validation |
| orp-proto | Protobuf definitions |
| orp-ml | ML inference hooks |
| orp-testbed | Integration test harness |

---

## Protocol Adapters (40+)

AIS, AISStream.io, ADS-B, MAVLink, CoT/TAK, KLV (MISB ST 0601), CCSDS+SGP4, OPC-UA, MQTT, Modbus, DNP3, BACnet, CAN bus, Kafka, NATS, NMEA, NMEA 2000, LoRaWAN, ACARS, ASTERIX, CAP, CEF, GRIB, GeoJSON, GTFS, HL7 v2.5, HTTP poller, Netflow, PCAP, SparkplugB, STIX, Syslog, Zeek, CSV watcher, WebSocket client, generic REST API

---

## Storage Detail

- **DuckDB** — OLAP analytics, ORP-QL to SQL, entity history
- **RocksDB** — dedup state, dead-letter queue, event store
- **Graph** — in-memory projection over DuckDB, Cypher-style traversal, refreshed every 30s

---

## OneMind Integration Points

- **NATS connector** (orp-connector/nats.rs) — ingests ent.* and tel.* subjects from the OneMind bus
- **CoT/TAK bridge** (orp-tak) — receives from tak.cot.out, injects to tak.cot.inject
- **MAVLink adapter** — feeds drone telemetry alongside om-overwatch-relay
- **ORP runs as a NEX workload** — confirmed running via nex-node on zeus-mac (nexus: onemind)
- **Config** lives at integrations/ai-agents/orp/config.yaml in this repo
- **Data** at ~/.local/share/orp/ (audit key) + local DuckDB/RocksDB paths

---

## Comparison

| | SQLite | Postgres | Lattice OS | ORP |
|--|--------|----------|------------|-----|
| Single binary | yes | no | no | yes |
| Real-time ingest | no | partial | yes | yes 50k ev/s |
| Protocol adapters | 0 | 0 | proprietary | 40+ open |
| Graph query language | no | ext only | proprietary | ORP-QL |
| Federation mesh | no | replication | closed | open |
| Edge / ARM64 | yes | partial | no | yes 45 MB |
| Built-in COP map | no | no | yes | yes |
| License | public domain | PostgreSQL | proprietary | Apache 2.0 |

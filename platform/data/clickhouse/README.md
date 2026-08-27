# omos-observ

Lean, ClickHouse-native observability for the OneMind bus. Supersedes the old
`services/observability/` (Grafana + Loki + Jaeger) stack.

## Why this shape

OneMind already committed to its observability backend: **ClickHouse ingests off NATS
natively** via the NATS table engine (`onemind.nats_evt` → `evt_to_events` → `events`,
subject `evt.>`). This stack extends that one pattern instead of bolting on parallel
stores.

Four heavy stores collapse to one:

| Concern  | Old            | omos-observ                                             |
|----------|----------------|---------------------------------------------------------|
| Logs     | Loki           | ClickHouse `onemind.logs`  (NATS engine on `log.>`)     |
| Metrics  | (Prometheus)   | ClickHouse `onemind.metrics` (scraped from NATS `:8222`) |
| Traces   | Jaeger         | dropped until something emits OTLP                      |
| Events   | —              | ClickHouse `onemind.events` (existing `evt.>` pipeline) |
| Viz      | Grafana        | Grafana (provisioned, ClickHouse datasource)           |
| Uptime   | —              | Uptime Kuma (blackbox pings)                            |

## Native launchd, not NEX

The collectors run as **native launchd** agents, not NEX workloads. Observability must
not depend on the mesh it watches — if NEX dies, the watchtower has to stay up to debug
it. This matches the tier NATS / ClickHouse / nex-node already live in.

- Native launchd: ClickHouse, `omos.observ.surveyor` (:7777), `omos.observ.metrics` (this repo).
- Docker (vendor UIs only): Grafana, Uptime Kuma.

## Layout

```
omos-observ/
├── clickhouse/observ.sql              # idempotent: log.> + metric.> tables/MVs + metrics
├── collectors/
│   ├── nats-metrics-scrape.py         # NATS :8222 + Surveyor :7777 -> onemind.metrics (15s)
│   └── kuma-seed.py                   # idempotent Uptime Kuma monitor seeder
├── launchd/
│   ├── omos.observ.surveyor.plist     # NATS Surveyor w/ $SYS creds -> :7777
│   └── omos.observ.metrics.plist      # ClickHouse metric collector
├── grafana/provisioning/
│   ├── datasources/observ.yml         # ClickHouse datasource
│   └── dashboards/
│       ├── nats-observability.json    # connections, JetStream, per-stream messages
│       └── onemind-logs.json          # logs + events explorer, volume, top subjects
├── docker-compose.yml                 # Grafana + Uptime Kuma
└── scripts/{up,install-launchd,teardown-deadweight}.sh
```

## Bring up

```bash
cp .env.example .env      # optional: set GRAFANA_ADMIN_PASSWORD
bash scripts/up.sh        # schema + collector + Grafana/Kuma
```

- Grafana: http://localhost:3000 (admin / onemind)
- Uptime Kuma: http://localhost:3001 (admin / onemind2026)
- Surveyor: http://localhost:7777/metrics

Then retire the old stack:

```bash
bash scripts/teardown-deadweight.sh
```

## Metrics scraped into `onemind.metrics`

From NATS monitoring `:8222` (auth-free):
- `nats_varz_*` — connections, subscriptions, slow_consumers, in/out msgs+bytes, mem, cpu, routes
- `nats_jsz_*` — streams, consumers, messages, bytes, memory, storage, api_total, api_errors
- `nats_stream_*` — per-stream messages / bytes / consumers / seqs (labels: `stream=NAME`)

From NATS Surveyor `:7777` (system-account survey):
- `nats_core_*` — server-level connections, mem, cpu, subscriptions, JetStream cluster info
- `nats_jetstream_*` — account/stream/consumer JetStream detail

## Uptime Kuma monitors (seeded by `collectors/kuma-seed.py`)

`scripts/up.sh` runs the seeder automatically. It sets up the admin account on first run
and upserts these monitors (idempotent):

| Monitor             | Type | Target                                   |
|---------------------|------|------------------------------------------|
| NATS client         | TCP  | host.docker.internal:4222                |
| NATS monitoring     | HTTP | http://host.docker.internal:8222/varz    |
| ClickHouse          | HTTP | http://host.docker.internal:8123/ping    |
| Grafana             | HTTP | http://host.docker.internal:3000/api/health |
| TAK Server (device) | TCP  | host.docker.internal:8089                |
| WebTAK              | TCP  | host.docker.internal:8446 (mTLS -> port check) |
| GeoServer           | TCP  | host.docker.internal:8600                |
| Murmur (voice)      | TCP  | host.docker.internal:64738               |

Re-run manually: `KUMA_PASSWORD=onemind2026 ./.venv/bin/python collectors/kuma-seed.py`

## Surveyor auth (resolved)

NATS runs with `auth_required` and a `$SYS` system account (user `sys`). The old NEX
`nats-surveyor` workload connected unauthenticated (→ `OMOS` account) so it surveyed
nothing (`nats_survey_surveyed_count 0`). It's replaced by the native
`omos.observ.surveyor` launchd agent, which authenticates as `sys` and surveys the
server. Its metrics are folded into ClickHouse via the `SURVEYOR` env in the metrics plist.

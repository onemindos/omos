# services/ — runtime service code

OneMind OS runtime services: Python bus subscribers & provisioning scripts, plus the
Go services monorepo absorbed from the former `omos-services` repo (merged
2026-08-30, history preserved via git-filter-repo subdirectory merge).

## Python subscribers (live operational scripts)

| Script | Purpose |
|--------|---------|
| `om_task.py` | Task Engine subscriber (Python predecessor of `services/task-engine/`) |
| `om_presence.py` | Presence / voice bridge subscriber |
| `notify_router.py` | Routes `evt.notify.>` events to notification channels |
| `graphite_bridge.py` | Bridges NATS telemetry to Graphite/Carbon |
| `netdata_exporter.py` | Exports NATS metrics to Netdata |
| `subscriber.py` | Base subscriber class / generic bus subscriber |
| `subscriber_runner.py` | Process supervisor — runs and restarts subscriber processes |
| `create_streams.py` | Provisioning: creates all JetStream streams |
| `create_buckets.py` | Provisioning: creates all JetStream KV buckets |
| `seed_rules.py` | Seeds initial rules/config into KV buckets |
| `ctl.py` | CLI control tool — start/stop/status subscribers |
| `launch.sh` / `launch_all.sh` | Start single / all subscribers |

## Go services (from omos-services)

| Service | Path | Purpose |
|---------|------|---------|
| **task-engine** | `task-engine/` | NATS-native task orchestration — create, dispatch, track tasks |
| **omd** | `omd/` | OneMind Daemon — runs on every node, supervises local agent harnesses |
| **agentops** | `agentops/` | Agent operations (Go): agents as map entities and TAK contacts, ATAK GeoChat to agents, deterministic watch jobs that propose through policy, AOI briefs with an offline model fallback |
| **bridge** | `bridge/` | Universal protocol bridge (MAVLink, MQTT, HA, TAK, webhook) → NATS |
| **nex** | `nex/` | NATS Execution Engine workloads (overwatch bridges, om-task) |
| **agents/meta** | `agents/meta/` | Meta agent service |
| **cot2constellation** | `bridge/bridges/tak/` | TAK CoT → Overwatch bridge |

## Development

Go services — each has its own `go.mod`, build individually:
```bash
cd services/task-engine && go build ./...
cd services/omd && go build ./...
```

Python subscribers — see launchers above. All services connect to
`nats://nats.nats.svc.cluster.local:4222` (K8s) or via edge leaf node, speak the
OneMind subject grammar (`cmd.*`, `tel.*`, `evt.*`, `svc.*`), and deploy via
`ghcr.io/onemindos/{service}:latest` with manifests in `platform/`.

In K8s the NATS servers accept clients only from the namespaces listed in
`platform/nats/network-policy.yaml`: a service in another namespace needs a rule
there (see `platform/nats/OPERATIONS.md`), or its connections time out silently.

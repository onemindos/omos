# omos-services

OneMind OS — NATS-native Go services monorepo.

## Services

| Service | Path | Purpose |
|---------|------|---------|
| **task-engine** | `task-engine/` | NATS-native task orchestration — create, dispatch, track tasks; fires real-world actions to drones, robots, agents |
| **omd** | `omd/` | OneMind Daemon — runs on every node (Mac, Pi, robot, drone); registers actor, subscribes to `cmd.node.{id}.>`, supervises local agent harnesses |
| **bridge** | `bridge/` | Universal protocol bridge — connects every device and protocol (MAVLink, MQTT, Home Assistant, TAK, webhook) to the NATS fabric |

## Architecture

All services are NATS-native Go binaries. Each:
- Connects to `nats://nats.nats.svc.cluster.local:4222` (K8s) or via leaf node (edge)
- Speaks the OneMind subject grammar (`cmd.*`, `tel.*`, `evt.*`, `svc.*`)
- Wraps all messages in the P0-E event envelope (see `omos-core/schemas/envelope/`)
- Deploys via `ghcr.io/onemindos/{service}:latest` — manifests in `omos-infra/platform/`

## Subject Grammar

```
tel.{node}.{class}.{metric}       telemetry
evt.{domain}.{entity}.{event}     semantic events
cmd.{domain}.{entity}.{action}    commands
task.create / task.dispatch.*     task engine RPC
cmd.node.{node-id}.>              omd command channel
```

## Development

Each service has its own `go.mod`. Build individually:
```bash
cd task-engine && go build ./...
cd omd && go build ./...
cd bridge/bridges/mavlink && go build ./...
```

Images are built and pushed to GHCR via GitHub Actions. K8s manifests live in `omos-infra/platform/`.

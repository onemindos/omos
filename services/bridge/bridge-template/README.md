# bridge-template

Copy this directory to `bridges/your-protocol` to start a new bridge.

## Structure

```
bridges/your-protocol/
├── main.go          ← entry point — wire your adapter here
├── adapter.go       ← implement BridgeAdapter interface
├── go.mod           ← module: github.com/onemindos/omos-bridge/bridges/your-protocol
├── go.sum
├── .env.example     ← document all required env vars
├── bridge.yaml      ← NEX Nexfile for edge deployment
└── README.md        ← protocol overview, setup, env vars, example output
```

## Quickstart

```bash
# 1. Copy template
cp -r bridge-template bridges/your-protocol
cd bridges/your-protocol

# 2. Rename module
# Edit go.mod: change module path to github.com/onemindos/omos-bridge/bridges/your-protocol

# 3. Implement BridgeAdapter interface in adapter.go

# 4. Configure
cp .env.example .env

# 5. Run
NATS_URL=nats://localhost:4222 go run .
```

## Standard Contract

All bridges must comply. Non-negotiable.

**Subjects:**
```
omos.{domain}.{entity_id}.telemetry   ← you publish here
omos.{domain}.{entity_id}.command     ← you subscribe here
omos.bridge.{name}.health             ← heartbeat every 30s
omos.bridge.{name}.errors             ← structured errors
```

**Envelope:**
```json
{
  "id": "ULID",
  "ts": 1718123456789,
  "source": "your-bridge",
  "entity": "domain.entity-id",
  "domain": "domain",
  "payload": {}
}
```

**Process:**
- All config from env vars
- Heartbeat every 30s
- Graceful SIGTERM shutdown with `nc.Drain()`
- Auto-reconnect — never crash on disconnect

Full spec: [../../README.md](../../README.md#the-standard-contract)

# Contributing to omos-bridge

First: thank you. Every bridge you write connects more of the world to the fabric.

## The One Rule

Your bridge must implement the **standard contract** — subject pattern, message envelope, process contract. Everything else (language choice within supported SDKs, internal architecture, logging style) is yours.

No contract compliance = PR does not merge. No exceptions.

---

## Before You Start

1. Check the [Bridge Catalog](./README.md#bridge-catalog) — is your protocol already listed or in progress?
2. Open an issue: `[Bridge Proposal] protocol-name` — describe the protocol, the domain it belongs to, and your intended language.
3. Wait for a maintainer to label it `accepted` before opening a PR. Avoids duplicate work.

---

## Creating a Bridge

### Step 1 — Copy the template

```bash
cp -r bridge-template bridges/your-protocol
cd bridges/your-protocol
```

### Step 2 — Implement the adapter

In Go:
```go
// Implement the BridgeAdapter interface in adapter.go
type BridgeAdapter interface {
    Connect(ctx context.Context) error
    Subscribe(ctx context.Context, out chan<- BridgeMessage) error
    SendCommand(ctx context.Context, cmd Command) error
    Disconnect() error
}
```

### Step 3 — Test it

```bash
# Start a local NATS server
nats-server -js

# Run your bridge
NATS_URL=nats://localhost:4222 go run .

# Verify output
nats sub "omos.{your-domain}.*.telemetry"
nats sub "omos.bridge.{your-bridge}.health"
```

### Step 4 — Update the catalog

Add your bridge to the catalog table in [README.md](./README.md).

### Step 5 — Open a PR

Title: `[Bridge] Add {protocol} bridge`

PR must include:
- [ ] Standard contract compliance (subject pattern, envelope, process contract)
- [ ] `.env.example` with all required env vars documented
- [ ] `bridge.yaml` (NEX Nexfile)
- [ ] `README.md` for the bridge (protocol overview, setup, env vars, example output)
- [ ] At least one working test

---

## Standards Reference

Full standard contract: [README.md — The Standard Contract](./README.md#the-standard-contract)

Subject pattern, message envelope, and process contract are **locked**. If you think a field needs to change, open an issue first.

---

## Code of Conduct

Be direct. Be useful. Ship working code.

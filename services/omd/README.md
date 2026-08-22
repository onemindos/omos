# `omd/`

**OneMind Daemon** — the L3 Compute layer of OneMind OS. A Go binary that runs on every node in the fleet.

> Status: 🔜 **Under construction** — skeleton is in place, implementation in progress.

---

## What it does

On boot, `omd`:
1. Connects to NATS (local socket, or Tailscale to core)
2. Registers the node as an **actor** and advertises its **capabilities** (`fs.read`, `fs.write`, `shell.exec`, `agent.run:hermes`, `gpu`, `camera`, …)
3. Subscribes to its command subject `cmd.node.{node-id}.>`
4. Discovers and supervises local agent harnesses via the **Synadia Agent SDK v0.3** (Hermes, Claude Code, Codex, …)

Every node in the fleet — Mac, barn Pi, office box, robot, drone — runs `omd`. This is what makes them first-class actors on the NATS nervous system.

---

## Structure

```
cmd/omd/        main entrypoint
internal/
  actor/        actor registration + capability advertisement
  bus/          NATS connection management
  config/       config loading
go.mod
go.sum
```

---

## Build

```bash
cd omd
go build ./cmd/omd
./omd --nats nats://localhost:4222
```

---

## Architecture reference

See `docs/ARCHITECTURE.md §3 — L3 Compute` for the full design spec.

---
name: OpenClaw CLI — Node (Headless)
updated: 2026-05-11
---

# `openclaw node`

Deploys and manages a **headless node agent** — a lightweight OpenClaw process that can receive and execute remote commands without running a full gateway.

**Official docs:** [docs.openclaw.ai/cli/node](https://docs.openclaw.ai/cli/node.md)

> Note: This is `openclaw node` (singular — headless agent deployment). For managing paired OpenClaw nodes (remote gateways), see [`nodes.md`](nodes.md).

---

## Commands

```bash
# Start headless node
openclaw node start                           # Start in foreground
openclaw node start --daemon                  # Start as background daemon
openclaw node start --gateway wss://my-host  # Connect to specific gateway

# Stop node
openclaw node stop

# Status
openclaw node status

# Register with a gateway
openclaw node register --gateway wss://my-host --token <tok>
openclaw node register --name "edge-node-1"

# Execute a command on this node (from another machine)
openclaw nodes invoke <nodeId> --command "openclaw status"
```

---

## Use Cases

- Edge device that runs local tools but delegates AI to a remote gateway
- Air-gapped or resource-constrained environments
- Secondary execution host for distributing agent workloads

---

## Headless vs Full Gateway

| Feature | Headless Node | Full Gateway |
|---|---|---|
| Receives commands | ✅ | ✅ |
| Hosts agents | ❌ | ✅ |
| Channel integrations | ❌ | ✅ |
| Local tool execution | ✅ | ✅ |
| Web UI | ❌ | ✅ |
| Resource footprint | Low | Full |

---

## Notes

- Headless nodes appear in `openclaw nodes list` on the gateway they registered with
- Use `--daemon` for persistent deployments; manage with `openclaw node stop/status`
- Pair with `openclaw nodes invoke` from the gateway host to send commands to the node

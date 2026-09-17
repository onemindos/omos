---
name: OpenClaw — Gateway
updated: 2026-05-10
---

# Gateway

Gateway configuration, health checks, architecture, VPS deployment, logging, and node setup for omos-core and omos-sentinel.

**Official docs:** [docs.openclaw.ai/gateway](https://docs.openclaw.ai/gateway/index.md)

## What's Here

| File | Description |
|---|---|
| [OPENCLAW-ARCHITECTURE.md](OPENCLAW-ARCHITECTURE.md) | Gateway architecture overview |
| [OPENCLAW_ARCHITECTURE_REFERENCE.md](OPENCLAW_ARCHITECTURE_REFERENCE.md) | Detailed architecture reference |
| [OPENCLAW-CONFIG-REFERENCE.md](OPENCLAW-CONFIG-REFERENCE.md) | Full configuration reference |
| [DEPLOYMENT.md](DEPLOYMENT.md) | Deployment runbook |
| [NODES.md](NODES.md) | Node setup and management |
| [VPS-GATEWAY.md](VPS-GATEWAY.md) | VPS gateway configuration |
| [VPS-SERVICES.md](VPS-SERVICES.md) | VPS system services overview |
| [DISASTER-RECOVERY.md](DISASTER-RECOVERY.md) | Disaster recovery procedures |

## Key Official Pages

- [Gateway runbook](https://docs.openclaw.ai/gateway/index.md)
- [Configuration reference](https://docs.openclaw.ai/gateway/configuration-reference.md)
- [Health checks](https://docs.openclaw.ai/gateway/health.md)
- [Secrets management](https://docs.openclaw.ai/gateway/secrets.md)
- [OpenTelemetry export](https://docs.openclaw.ai/gateway/opentelemetry.md)
- [Tailscale](https://docs.openclaw.ai/gateway/tailscale.md)

## Our Nodes

| Node | Role | Tailscale IP |
|---|---|---|
| omos-core | Primary | `100.123.126.34` |
| omos-sentinel | Failover | — |

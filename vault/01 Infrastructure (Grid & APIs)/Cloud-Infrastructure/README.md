# OneMind Cloud Infrastructure Documentation

*Complete reference for the OneMind Cloud deployment - January 2026*

## Documentation Index

| Document | Purpose | When to Use |
|----------|---------|-------------|
| [ONEMIND-CLOUD-INFRASTRUCTURE.md](./ONEMIND-CLOUD-INFRASTRUCTURE.md) | Master overview of the entire cloud stack | Understanding the big picture |
| [CLOUDFLARE-TUNNEL-CONFIG.md](./CLOUDFLARE-TUNNEL-CONFIG.md) | Cloudflare Tunnel and Zero Trust setup | Adding routes, configuring auth |
| [DOCKER-SERVICES-REFERENCE.md](./DOCKER-SERVICES-REFERENCE.md) | Detailed container configurations | Debugging containers, health checks |
| [OPERATIONS-RUNBOOK.md](./OPERATIONS-RUNBOOK.md) | Day-to-day procedures and commands | Daily ops, troubleshooting, emergencies |
| [ACCESS-PROTOCOLS.md](./ACCESS-PROTOCOLS.md) | Security tiers and authentication | Access issues, adding users/devices |

---

## Quick Start

### Check Server Status
```bash
ssh ubuntu@100.87.87.69 "docker ps --format 'table {{.Names}}\t{{.Status}}'"
```

### Access Services

**Via Tailscale (Fast):**
- Portainer: http://100.87.87.69:9000
- N8N: http://100.87.87.69:5678

**Via Cloudflare (Secure):**
- Portainer: https://portainer.onemindos.dev
- N8N: https://n8n.onemindos.dev

### Common Commands
```bash
# Restart a service
docker restart onemind-n8n

# View logs
docker logs onemind-agno --tail 100 -f

# Check health
docker ps --filter "health=unhealthy"
```

---

## Architecture Summary

```
                    ┌─────────────────┐
                    │   USER ACCESS   │
                    └────────┬────────┘
                             │
              ┌──────────────┴──────────────┐
              │                             │
     ┌────────▼────────┐         ┌─────────▼─────────┐
     │    TAILSCALE    │         │    CLOUDFLARE     │
     │  100.87.87.69   │         │  *.onemindos.dev  │
     │    (Private)    │         │    (Public)       │
     └────────┬────────┘         └─────────┬─────────┘
              │                             │
              │    ┌─────────────────┐      │
              │    │   ZERO TRUST    │      │
              │    │  (6 protected   │◄─────┘
              │    │     apps)       │
              │    └────────┬────────┘
              │             │
              └──────┬──────┘
                     │
        ┌────────────▼────────────┐
        │   AWS LIGHTSAIL         │
        │   onemind-cloud-1       │
        │   16GB RAM / 4 vCPU     │
        ├─────────────────────────┤
        │                         │
        │   ┌─────────────────┐   │
        │   │ DOCKER COMPOSE  │   │
        │   │                 │   │
        │   │ • Agno (AI)     │   │
        │   │ • Letta (Mem)   │   │
        │   │ • N8N (Auto)    │   │
        │   │ • Kestra (Data) │   │
        │   │ • Portainer     │   │
        │   │ • PostgreSQL    │   │
        │   │ • Redis         │   │
        │   │ • Qdrant        │   │
        │   │ • + 5 more      │   │
        │   └─────────────────┘   │
        │                         │
        └─────────────────────────┘
```

---

## Key Credentials (1Password)

| Item | Vault | Contains |
|------|-------|----------|
| Cloudflare - OneMind OS | OneMind - UI | API keys, Account ID, Zone ID, Tunnel ID |
| AWS Lightsail | OneMind - UI | Console access |
| Tailscale | OneMind - UI | Admin console access |

---

## Service Status Quick Reference

| Service | Tailscale | Cloudflare | Protected |
|---------|-----------|------------|-----------|
| Agno | :8080 | agno.onemindos.dev | Yes |
| Letta | :8283 | letta.onemindos.dev | No |
| N8N | :5678 | n8n.onemindos.dev | Yes |
| Kestra | :8088 | kestra.onemindos.dev | Yes |
| Portainer | :9000 | portainer.onemindos.dev | Yes |
| TOC | :3002 | toc.onemindos.dev | Yes |
| VNC | :6080 | vnc.onemindos.dev | Yes |
| Uptime Kuma | :3001 | status.onemindos.dev | No |
| NATS | :8222 | nats.onemindos.dev | No |
| Actions | :8001 | actions.onemindos.dev | No |

---

## Emergency Contacts

- **AWS Support:** Through console
- **Cloudflare Support:** https://dash.cloudflare.com (Support tab)
- **Documentation Issues:** Update files in this folder

---

## Changelog

| Date | Author | Change |
|------|--------|--------|
| 2026-01-18 | Claude Code | Initial documentation suite created |
| 2026-01-18 | Claude Code | Zero Trust configured for 6 applications |
| 2026-01-18 | Claude Code | Added toc, vnc, actions Cloudflare routes |

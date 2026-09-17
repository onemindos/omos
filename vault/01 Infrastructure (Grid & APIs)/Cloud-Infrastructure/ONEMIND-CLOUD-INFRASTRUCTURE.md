# OneMind Cloud Infrastructure

*Master Reference Document - January 2026*

## Executive Summary

OneMind Cloud is a self-hosted AI infrastructure running on AWS Lightsail, providing the computational backbone for the LEGACY AI system. It hosts 13+ Docker containers orchestrated via Docker Compose, accessible through dual pathways: Tailscale (private VPN) and Cloudflare Tunnel (public HTTPS).

---

## Infrastructure Stack

```
┌─────────────────────────────────────────────────────────────────┐
│                     ONEMIND CLOUD ARCHITECTURE                  │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│   ┌─────────────┐     ┌─────────────┐     ┌─────────────┐      │
│   │   MOBILE    │     │   DESKTOP   │     │    DEV      │      │
│   │  (Browser)  │     │  (Browser)  │     │   (SSH)     │      │
│   └──────┬──────┘     └──────┬──────┘     └──────┬──────┘      │
│          │                   │                   │              │
│          ▼                   ▼                   ▼              │
│   ┌─────────────────────────────────────────────────────┐      │
│   │              ACCESS LAYER                            │      │
│   │  ┌─────────────────┐  ┌─────────────────┐           │      │
│   │  │   CLOUDFLARE    │  │    TAILSCALE    │           │      │
│   │  │  *.onemindos.dev│  │  100.87.87.69   │           │      │
│   │  │  (Public HTTPS) │  │  (Private VPN)  │           │      │
│   │  └────────┬────────┘  └────────┬────────┘           │      │
│   │           │                    │                     │      │
│   │           ▼                    ▼                     │      │
│   │  ┌─────────────────────────────────────────┐        │      │
│   │  │        ZERO TRUST AUTH LAYER            │        │      │
│   │  │   (Email verification for Zeus only)    │        │      │
│   │  └─────────────────────────────────────────┘        │      │
│   └─────────────────────────────────────────────────────┘      │
│                              │                                  │
│                              ▼                                  │
│   ┌─────────────────────────────────────────────────────┐      │
│   │              AWS LIGHTSAIL INSTANCE                  │      │
│   │              onemind-cloud-1 (Ubuntu)                │      │
│   │              16GB RAM / 4 vCPU                       │      │
│   ├─────────────────────────────────────────────────────┤      │
│   │                                                      │      │
│   │  ┌─────────────────────────────────────────────┐    │      │
│   │  │           DOCKER COMPOSE STACK              │    │      │
│   │  │                                             │    │      │
│   │  │  ┌─────────┐ ┌─────────┐ ┌─────────┐       │    │      │
│   │  │  │  AGNO   │ │  LETTA  │ │   N8N   │       │    │      │
│   │  │  │ :8080   │ │ :8283   │ │ :5678   │       │    │      │
│   │  │  └─────────┘ └─────────┘ └─────────┘       │    │      │
│   │  │                                             │    │      │
│   │  │  ┌─────────┐ ┌─────────┐ ┌─────────┐       │    │      │
│   │  │  │ KESTRA  │ │PORTAINER│ │  UPTIME │       │    │      │
│   │  │  │ :8088   │ │ :9000   │ │ :3001   │       │    │      │
│   │  │  └─────────┘ └─────────┘ └─────────┘       │    │      │
│   │  │                                             │    │      │
│   │  │  ┌─────────┐ ┌─────────┐ ┌─────────┐       │    │      │
│   │  │  │   TOC   │ │   VNC   │ │ ACTIONS │       │    │      │
│   │  │  │ :3002   │ │ :6080   │ │ :8001   │       │    │      │
│   │  │  └─────────┘ └─────────┘ └─────────┘       │    │      │
│   │  │                                             │    │      │
│   │  │  ┌─────────┐ ┌─────────┐ ┌─────────┐       │    │      │
│   │  │  │  NATS   │ │ QDRANT  │ │POSTGRES │       │    │      │
│   │  │  │ :4222   │ │ :6333   │ │ :5432   │       │    │      │
│   │  │  └─────────┘ └─────────┘ └─────────┘       │    │      │
│   │  │                                             │    │      │
│   │  └─────────────────────────────────────────────┘    │      │
│   │                                                      │      │
│   └─────────────────────────────────────────────────────┘      │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## Server Details

| Property | Value |
|----------|-------|
| **Provider** | AWS Lightsail |
| **Hostname** | onemind-cloud-1 |
| **OS** | Ubuntu Linux |
| **RAM** | 16 GB |
| **vCPUs** | 4 |
| **Public IP** | 54.165.53.218 |
| **Tailscale IP** | 100.87.87.69 |
| **Docker Version** | 29.1.5 |
| **Compose Path** | `/home/ubuntu/onemind-os/docker_substrate_runtime/cloud/` |

---

## Network Architecture

### Dual Access Model

OneMind uses a **dual-access architecture** providing both private and public pathways:

| Pathway | Technology | Domain/IP | Use Case |
|---------|------------|-----------|----------|
| **Private** | Tailscale VPN | 100.87.87.69 | Development, debugging, direct access |
| **Public** | Cloudflare Tunnel | *.onemindos.dev | Mobile, production, external APIs |

### Why Both?

```
TAILSCALE (Private)                    CLOUDFLARE (Public)
─────────────────────                  ─────────────────────
✓ Direct P2P connection (~5ms)         ✓ Works from any browser
✓ Full port access                     ✓ Automatic HTTPS/SSL
✓ No auth prompts (device-based)       ✓ Zero Trust authentication
✓ Bypass rate limits                   ✓ DDoS protection
✓ Access internal ports                ✓ Mobile-friendly
✗ Requires Tailscale app               ✗ Adds ~50ms latency
✗ Device must be on tailnet            ✗ Only configured ports
```

**Rule of Thumb:**
- Use **Tailscale** when on your own devices for development
- Use **Cloudflare** when on untrusted networks or mobile

---

## Service Inventory

### AI & Automation Services

| Service | Port | Purpose | Health |
|---------|------|---------|--------|
| **Agno (LEGACY Brain)** | 8080 | Primary AI agent runtime | ✅ |
| **Letta** | 8283 | Long-term AI memory server | ✅ |
| **N8N** | 5678 | Workflow automation | ✅ |
| **Kestra** | 8088 | Data pipeline orchestration | ✅ |
| **Actions** | 8001 | Custom action gateway API | ✅ |

### Infrastructure Services

| Service | Port | Purpose | Health |
|---------|------|---------|--------|
| **Portainer** | 9000 | Docker management UI | ✅ |
| **Uptime Kuma** | 3001 | Service monitoring | ✅ |
| **TOC** | 3002 | Table of Contents / Dashboard | ✅ |
| **VNC** | 6080 | Remote desktop (noVNC) | ✅ |
| **NATS** | 4222/8222 | Message queue (data/monitor) | ✅ |

### Data Services

| Service | Port | Purpose | Health |
|---------|------|---------|--------|
| **PostgreSQL** | 5432 | Primary database | ✅ |
| **Qdrant** | 6333 | Vector database | ✅ |
| **Redis** | 6379 | Cache & pub/sub | ✅ |

### Stopped Services

| Service | Port | Status | Reason |
|---------|------|--------|--------|
| **Voice** | - | ⏹ Stopped | Python event loop bug |
| **Agent-UI** | 3000 | ⏹ Stopped | Not configured |

---

## File System Layout

```
/home/ubuntu/onemind-os/
├── docker_substrate_runtime/
│   └── cloud/
│       ├── docker-compose.yml      # Main compose file
│       ├── .env                    # Environment variables
│       └── data/                   # Persistent volumes
│           ├── postgres/
│           ├── qdrant/
│           ├── n8n/
│           ├── letta/
│           └── ...
├── docs_memory_banks/
│   └── ZEUS-ACCESS-PROTOCOL.md     # Access documentation
└── logs/
    └── *.log                       # Application logs
```

---

## Systemd Services

Beyond Docker, two critical services run as systemd units:

### Cloudflare Tunnel
```bash
# Service file: /etc/systemd/system/cloudflared.service
sudo systemctl status cloudflared
sudo systemctl restart cloudflared
sudo journalctl -u cloudflared -f
```

### Tailscale
```bash
# Built-in service
sudo systemctl status tailscaled
tailscale status
```

---

## Quick Reference

### Common Operations

```bash
# SSH to server
ssh ubuntu@100.87.87.69

# Check all containers
docker ps --format 'table {{.Names}}\t{{.Status}}\t{{.Ports}}'

# View specific logs
docker logs onemind-agno --tail 100 -f

# Restart a service
docker restart onemind-n8n

# Check disk space
df -h

# Check memory
free -h

# View compose config
cat /home/ubuntu/onemind-os/docker_substrate_runtime/cloud/docker-compose.yml
```

### Emergency Commands

```bash
# Restart all containers
cd /home/ubuntu/onemind-os/docker_substrate_runtime/cloud
docker compose restart

# Full stack restart
docker compose down && docker compose up -d

# Check for zombie processes
docker ps -a --filter "status=exited"

# Nuclear option (removes all stopped containers)
docker container prune -f
```

---

## Related Documentation

| Document | Purpose |
|----------|---------|
| [Cloudflare Configuration](./CLOUDFLARE-TUNNEL-CONFIG.md) | Tunnel & Zero Trust setup |
| [Docker Services Reference](./DOCKER-SERVICES-REFERENCE.md) | Detailed container configs |
| [Operations Runbook](./OPERATIONS-RUNBOOK.md) | Day-to-day procedures |
| [Access Protocols](./ACCESS-PROTOCOLS.md) | Security & authentication |

---

## Changelog

| Date | Change |
|------|--------|
| 2026-01 | Initial infrastructure documentation |
| 2026-01 | Added Zero Trust Access for 6 apps |
| 2026-01 | Added toc, vnc, actions routes to Cloudflare |
| 2026-01 | Fixed container health checks |
| 2026-01 | Removed redundant Caddy container |

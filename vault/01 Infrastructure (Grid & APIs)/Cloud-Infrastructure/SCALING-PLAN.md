# OneMind OS - VM Architecture Scaling Plan

## Overview

This document outlines the scaling path from a single VM to multi-region deployment. Each phase is designed to be a natural progression - only move to the next phase when you hit the limits of the current one.

## Scaling Phases

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                         SCALING PROGRESSION                                      │
├─────────────────────────────────────────────────────────────────────────────────┤
│                                                                                  │
│  PHASE 1              PHASE 2              PHASE 3              PHASE 4         │
│  Single VM            Two VMs              Multi-VM             Multi-Region    │
│  (Start Here)         (When Needed)        (Scale)              (Global)        │
│                                                                                  │
│  ┌─────────┐         ┌─────────┐          ┌─────────┐          ┌─────────┐     │
│  │ 8GB VM  │         │ AI VM   │          │ AI VM   │          │ US-East │     │
│  │         │   ───►  │ 4GB     │   ───►   │ 4GB     │   ───►   │ Primary │     │
│  │ All     │         ├─────────┤          ├─────────┤          ├─────────┤     │
│  │Services │         │Utility  │          │Utility  │          │ EU-West │     │
│  │         │         │ VM 4GB  │          │ VM 8GB  │          │  DR     │     │
│  └─────────┘         └─────────┘          ├─────────┤          └─────────┘     │
│                                           │ Edge x3 │                           │
│  ~$40/mo             ~$40/mo              │ Jetsons │          ~$120/mo        │
│                                           └─────────┘                           │
│                                           ~$60/mo                               │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

---

## Phase 1: Single VM (Current - Start Here)

**When**: Starting out, development, low traffic

**Architecture**:
```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                              SINGLE VM (8GB RAM)                                 │
│                              Azure VM                        │
│                              ~$40/month                                          │
├─────────────────────────────────────────────────────────────────────────────────┤
│                                                                                  │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │                         ALL PROFILES                                     │   │
│  │                                                                          │   │
│  │  CORE              OBSERVABILITY        NOTIFICATIONS      UTILITY      │   │
│  │  ────              ─────────────        ─────────────      ───────      │   │
│  │  Postgres          Netdata              Novu               Kestra       │   │
│  │  Redis             VictoriaMetrics      MongoDB            Guacamole    │   │
│  │  NATS              Grafana                                 EMQX         │   │
│  │  Agno              Uptime Kuma                             Redpanda     │   │
│  │  UI                                                                      │   │
│  │                                                                          │   │
│  │  RAM: ~1.5GB       RAM: ~600MB          RAM: ~400MB        RAM: ~800MB  │   │
│  │                                                                          │   │
│  │                         TOTAL: ~3.3GB                                    │   │
│  │                         AVAILABLE: ~4.7GB buffer                         │   │
│  │                                                                          │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                  │
│  EDGE LOCATIONS (Connected via NATS Leaf)                                       │
│  ─────────────────────────────────────────                                       │
│  Home ──► Home Assistant + Frigate + Mosquitto ──► EMQX (cloud) ──► NATS       │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

**Startup Command**:
```bash
docker compose --profile core --profile observability --profile notifications --profile utility up -d
```

**When to Move to Phase 2**:
- RAM consistently >70% (~5.6GB)
- CPU consistently >60%
- Response times degrading
- Need better isolation

---

## Phase 2: Two VMs (Recommended for Production)

**When**: Production use, need reliability, hitting resource limits

**Architecture**:
```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                              TWO VM ARCHITECTURE                                 │
│                              ~$40-50/month total                                 │
├─────────────────────────────────────────────────────────────────────────────────┤
│                                                                                  │
│  ┌─────────────────────────────┐         ┌─────────────────────────────┐       │
│  │       AI VM (4GB RAM)       │         │    UTILITY VM (4GB RAM)     │       │
│  │       ~$20/month            │         │       ~$20/month            │       │
│  │                             │         │                             │       │
│  │  CORE PROFILE               │         │  OBSERVABILITY PROFILE      │       │
│  │  ──────────────             │         │  ────────────────────       │       │
│  │  • Postgres + TimescaleDB   │         │  • Netdata                  │       │
│  │  • Redis                    │         │  • VictoriaMetrics          │       │
│  │  • NATS (Hub)               │◄───────►│  • Grafana                  │       │
│  │  • Agno (Legacy)            │  NATS   │  • Uptime Kuma              │       │
│  │  • UI                       │         │                             │       │
│  │                             │         │  UTILITY PROFILE            │       │
│  │  RAM: ~1.5GB                │         │  ───────────────            │       │
│  │  Purpose: Fast, critical    │         │  • Kestra                   │       │
│  │                             │         │  • Guacamole                │       │
│  │                             │         │  • EMQX                     │       │
│  │                             │         │  • Redpanda Connect         │       │
│  │                             │         │                             │       │
│  │                             │         │  NOTIFICATIONS PROFILE      │       │
│  │                             │         │  ──────────────────         │       │
│  │                             │         │  • Novu                     │       │
│  │                             │         │  • MongoDB                  │       │
│  │                             │         │                             │       │
│  │                             │         │  RAM: ~2.5GB                │       │
│  │                             │         │  Purpose: Heavy, restartable│       │
│  └─────────────────────────────┘         └─────────────────────────────┘       │
│                                                                                  │
│  Benefits:                                                                       │
│  • AI VM stays fast (not competing for resources)                               │
│  • Utility VM can be restarted without affecting core                           │
│  • Can scale each independently                                                  │
│  • Better failure isolation                                                      │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

**Docker Compose Split**:

AI VM (`docker-compose.ai.yml`):
```yaml
profiles: [core]
# postgres, redis, nats, agno, ui
```

Utility VM (`docker-compose.utility.yml`):
```yaml
profiles: [observability, notifications, utility]
# All other services + NATS connection to AI VM
```

**When to Move to Phase 3**:
- Adding multiple edge locations
- Need dedicated edge processing
- Want redundancy

---

## Phase 3: Multi-VM with Edge (Full Deployment)

**When**: Multiple locations, edge processing, full omnipresent mode

**Architecture**:
```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                           MULTI-VM WITH EDGE                                     │
│                           ~$60-80/month + edge hardware                          │
├─────────────────────────────────────────────────────────────────────────────────┤
│                                                                                  │
│                              CLOUD LAYER                                         │
│  ┌─────────────────────────────┐         ┌─────────────────────────────┐       │
│  │       AI VM (4GB)           │         │    UTILITY VM (8GB)         │       │
│  │                             │         │                             │       │
│  │  • Postgres + TimescaleDB   │         │  • Netdata                  │       │
│  │  • Redis                    │◄───────►│  • VictoriaMetrics          │       │
│  │  • NATS (Hub)               │  NATS   │  • Grafana                  │       │
│  │  • Agno (Legacy)            │         │  • Kestra                   │       │
│  │  • UI                       │         │  • Novu + MongoDB           │       │
│  │                             │         │  • EMQX                     │       │
│  │                             │         │  • Redpanda Connect         │       │
│  │                             │         │  • Guacamole                │       │
│  └──────────────┬──────────────┘         └─────────────────────────────┘       │
│                 │                                                                │
│            NATS Hub (:7422 leaf connections)                                    │
│                 │                                                                │
│  ┌──────────────┴──────────────┬─────────────────────┬────────────────────┐    │
│  │                             │                     │                    │    │
│  ▼                             ▼                     ▼                    ▼    │
│                              EDGE LAYER                                         │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐  ┌──────────┐  │
│  │   HOME EDGE     │  │  OFFICE EDGE    │  │   FARM EDGE     │  │   CAR    │  │
│  │   (Jetson 4GB)  │  │  (Jetson 4GB)   │  │  (Jetson 8GB)   │  │  (Pi 5)  │  │
│  │                 │  │                 │  │                 │  │          │  │
│  │ • Home Assistant│  │ • Home Assistant│  │ • Home Assistant│  │ • GPS    │  │
│  │ • Frigate       │  │ • Frigate       │  │ • FarmBot       │  │ • OBD-II │  │
│  │ • Ollama        │  │ • Ollama        │  │ • Ollama        │  │ • NATS   │  │
│  │ • RemoteAgent   │  │ • RemoteAgent   │  │ • RemoteAgent   │  │   Leaf   │  │
│  │ • Mosquitto     │  │ • Mosquitto     │  │ • Mosquitto     │  │          │  │
│  │ • NATS Leaf     │  │ • NATS Leaf     │  │ • NATS Leaf     │  │          │  │
│  │                 │  │                 │  │                 │  │          │  │
│  │ Cameras: 6      │  │ Cameras: 2      │  │ Sensors: many   │  │          │  │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘  └──────────┘  │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

**When to Move to Phase 4**:
- Global users needing low latency
- Compliance requirements (data residency)
- Disaster recovery requirements

---

## Phase 4: Multi-Region (Enterprise Scale)

**When**: Global presence, disaster recovery, compliance

**Architecture**:
```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                           MULTI-REGION DEPLOYMENT                                │
│                           ~$120-200/month                                        │
├─────────────────────────────────────────────────────────────────────────────────┤
│                                                                                  │
│     US-EAST (Primary)                         EU-WEST (Secondary/DR)            │
│  ┌─────────────────────────────┐         ┌─────────────────────────────┐       │
│  │       AI VM (4-8GB)         │         │       AI VM (4GB)           │       │
│  │                             │         │                             │       │
│  │  • Postgres (Primary)       │◄───────►│  • Postgres (Replica)       │       │
│  │  • NATS (Cluster)           │  NATS   │  • NATS (Cluster)           │       │
│  │  • Agno (Legacy)            │ Cluster │  • Agno (Standby)           │       │
│  │  • UI                       │         │  • UI                       │       │
│  └─────────────────────────────┘         └─────────────────────────────┘       │
│                                                                                  │
│  ┌─────────────────────────────┐         ┌─────────────────────────────┐       │
│  │    UTILITY VM (8GB)         │         │    UTILITY VM (4GB)         │       │
│  │                             │         │                             │       │
│  │  • Full observability       │◄───────►│  • Minimal observability    │       │
│  │  • Kestra (Primary)         │  Sync   │  • Kestra (Standby)         │       │
│  │  • Novu                     │         │                             │       │
│  │  • EMQX (Cluster)           │         │  • EMQX (Cluster)           │       │
│  └─────────────────────────────┘         └─────────────────────────────┘       │
│                                                                                  │
│                              EDGE DEVICES                                        │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │                                                                          │   │
│  │  US Edges ──► US-EAST NATS        EU Edges ──► EU-WEST NATS             │   │
│  │  (Auto-route to nearest)          (Auto-route to nearest)                │   │
│  │                                                                          │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                  │
│  Benefits:                                                                       │
│  • <50ms latency for all users                                                  │
│  • Automatic failover                                                            │
│  • Data residency compliance                                                     │
│  • True high availability                                                        │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

---

## Migration Checklist

### Phase 1 → Phase 2

```bash
# 1. Provision new Utility VM
# 2. Install Docker, Tailscale
# 3. Configure NATS to connect to AI VM

# On Utility VM:
export NATS_URL=nats://100.x.x.x:4222  # AI VM Tailscale IP

# 4. Move services one by one (test each):
#    - Move Netdata first (low risk)
#    - Move Grafana (low risk)
#    - Move Kestra (medium risk)
#    - Move EMQX (higher risk - test HA bridge)
#    - Move Novu last (test notifications)

# 5. Update DNS/firewall as needed
# 6. Monitor for 48 hours before decommissioning old config
```

### Phase 2 → Phase 3

```bash
# 1. Prepare Jetson images with edge profile
# 2. Generate NATS credentials for each edge
# 3. Configure HA Mosquitto → EMQX bridge
# 4. Deploy edge by edge, test each:
#    - Deploy Home edge
#    - Test HA → NATS flow
#    - Test Frigate → NATS flow
#    - Test RemoteAgent communication
# 5. Repeat for each location
```

### Phase 3 → Phase 4

```bash
# 1. Set up Postgres streaming replication
# 2. Configure NATS super-cluster
# 3. Set up global load balancer (Cloudflare)
# 4. Test failover scenarios
# 5. Configure edge devices for multi-region
# 6. Update DNS with geo-routing
```

---

## Resource Requirements by Phase

| Phase | VMs | Total RAM | Monthly Cost | Complexity |
|-------|-----|-----------|--------------|------------|
| **1** | 1 | 8GB | ~$40 | Low |
| **2** | 2 | 8GB | ~$40-50 | Medium |
| **3** | 2 + edges | 8GB + edge | ~$60-80 | Medium-High |
| **4** | 4+ | 16GB+ | ~$120-200 | High |

---

## Signals to Scale

| Signal | Action |
|--------|--------|
| RAM >70% sustained | Add RAM or split VMs |
| CPU >60% sustained | Split VMs or upgrade |
| Response time >500ms | Split to 2 VMs |
| Adding 3+ edge locations | Move to Phase 3 |
| Global users complaining | Move to Phase 4 |
| Compliance requirements | Move to Phase 4 |

---

## Rollback Plan

Each phase can be rolled back:

| From | To | Steps |
|------|-----|-------|
| Phase 2 | Phase 1 | Merge docker-compose, move all to single VM |
| Phase 3 | Phase 2 | Keep cloud VMs, disconnect edges |
| Phase 4 | Phase 3 | Decommission secondary region |

---

## Recommended Starting Point

**Start with Phase 1 (Single 8GB VM)** because:
1. Simpler to manage
2. All services on same network (no cross-VM latency)
3. Easy to debug
4. Profiles give flexibility to disable services if needed
5. Can scale to Phase 2 in ~1 hour when needed

**Move to Phase 2 when**:
- You've been running stable for 1+ month
- Resource usage consistently >60%
- You want better isolation for production

---
id: ui-02-20260111053925
title: "🎯 **Legacy Ai   Deployment Architecture Deep Dive**"
type: concept
domain: SO
domain: 00-24
folder: "00-24 SO (Sovereign Operations)"
status: active
created: 2026-01-11
modified: 2026-01-11
owner: zeus
tags: [domain/so, folder/02]
todoist_project: "00-24 SO (Sovereign Operations)"
ai_summary: "🎯 **Legacy Ai   Deployment Architecture Deep Dive**"
---

## Overview

  

This document provides a comprehensive guide for deploying Legacy AI across single and multi-VM architectures, with global redundancy, edge device management, and Portainer orchestration.

  

---

  

## **SINGLE VM vs. MULTI-VM - THE DEFINITIVE ANSWER**

  

### **✅ RECOMMENDED: Single Larger VM** (For Your Use Case)

  

**Why Single VM is Better for Legacy AI:**

  

1. **Self-Access & Tool Calling** 🔧

   - Legacy AI needs to call her own APIs (Dify, n8n, Home Assistant, etc.)

   - Single VM = `localhost` or `docker network` connections (instant, no authentication complexity)

   - Multi-VM = External HTTP calls, API keys, firewall rules, network latency

  

2. **Portainer Management** 🐳

   - Single VM: Portainer manages all stacks in one interface

   - Multi-VM: Portainer **can** manage multiple VMs via "Environments" feature, BUT:

     - Requires Portainer Agent on each VM

     - More complex to configure

     - Still need to switch between environments in UI

  

3. **Maintenance Overhead** 🛠️

   - Single VM: One OS to patch, one Docker daemon, one backup script

   - Multi-VM: 3-4 VMs to maintain, network between them, distributed monitoring

  

4. **Networking Simplicity** 🌐

   - Single VM: Docker bridge network, service discovery automatic

   - Multi-VM: Firewalls, VPN/VPC peering, service mesh (Consul/Istio), DNS setup

  

5. **Cost Efficiency** 💰

   - Single 64GB VM: ~$160/month (AWS Lightsail)

   - 3x 16-24GB VMs: ~$240/month total + networking costs

  

**When You SHOULD Split to Multi-VM:**

- ❌ Not for "separation of concerns" (Docker does that)

- ❌ Not for "security" (Docker networks + firewalls handle that)

- ✅ YES if you hit **hardware limits** (e.g., need 128GB RAM)

- ✅ YES for **geographic distribution** (different regions)

- ✅ YES for **specialized hardware** (GPU VM separate from CPU VM)

  

---

  

## **RECOMMENDED ARCHITECTURE: Single VM + Regional Mirrors**

  

### **Primary Architecture**

  

```

┌─────────────────────────────────────────────────────────────┐

│  AWS Lightsail 64GB/16vCPU VM (us-east-1)                  │

│  Cost: ~$160/month                                          │

├─────────────────────────────────────────────────────────────┤

│                                                              │

│  🧠 CORE Services (Docker Compose)                          │

│  ├─ Letta (brain)                   4GB RAM                 │

│  ├─ LiveKit (senses transport)      4GB RAM                 │

│  ├─ PostgreSQL (memory)             6GB RAM                 │

│  ├─ Redis (cache)                   2GB RAM                 │

│  └─ Supabase Sync                   2GB RAM                 │

│                                                              │

│  🔧 CAPABILITIES Services                                   │

│  ├─ Dify (workflow automation)      6GB RAM                 │

│  ├─ n8n (task automation)           2GB RAM                 │

│  ├─ Home Assistant (smart home)     4GB RAM                 │

│  ├─ Node-RED (visual flows)         1GB RAM                 │

│  ├─ EMQX (MQTT)                     2GB RAM                 │

│  ├─ Plane (project management)      8GB RAM                 │

│  ├─ Airbyte (data pipelines)        8GB RAM                 │

│  └─ ByteBot (browser automation)    4GB RAM                 │

│                                                              │

│  📊 OBSERVABILITY Services (Separate Docker Network)        │

│  ├─ Prometheus (metrics)            4GB RAM                 │

│  ├─ Grafana (dashboards)            2GB RAM                 │

│  ├─ Loki (log aggregation)          2GB RAM                 │

│  └─ Uptime Kuma (uptime)            0.5GB RAM               │

│                                                              │

│  🎛️ COMMAND-CENTER Services                                │

│  ├─ Portainer (container mgmt)      1GB RAM                 │

│  ├─ Dashy (service dashboard)       0.5GB RAM               │

│  └─ Caddy (reverse proxy + SSL)     0.5GB RAM               │

│                                                              │

│  📚 KNOWLEDGE Services                                      │

│  ├─ Obsidian Sync (CouchDB)         1GB RAM                 │

│  └─ Minio (S3-compatible storage)   1GB RAM                 │

│                                                              │

│  TOTAL: ~60GB RAM used (4GB buffer)                         │

└─────────────────────────────────────────────────────────────┘

```

  

### **Should You Separate Telemetry?**

  

**❌ NO - Keep Telemetry on Same VM** (with separate Docker network)

  

**Why:**

1. **Self-Monitoring**: Legacy AI should monitor herself without external dependencies

2. **Immediate Metrics**: Prometheus scrapes services on `localhost` (no network latency)

3. **Debugging**: When services fail, telemetry MUST stay up to diagnose

4. **Resource Efficiency**: Telemetry uses <10GB RAM total, fits easily

  

**How to Isolate:**

```yaml

# docker-compose.yml

networks:

  core_network:          # Letta, Dify, n8n, etc.

    driver: bridge

  observability_network: # Prometheus, Grafana

    driver: bridge

  public_network:        # Caddy (exposed to internet)

    driver: bridge

  

services:

  letta:

    networks:

      - core_network

      - observability_network  # Allow Prometheus to scrape

  prometheus:

    networks:

      - observability_network  # Can scrape all services

      - public_network         # Grafana accessible via Caddy

```

  

**Telemetry stays self-contained, but isolated from public exposure.**

  

---

  

## **GLOBAL REDUNDANCY & LOW-LATENCY ARCHITECTURE**

  

### **Strategy: Active-Active Multi-Region with Cloudflare**

  

```

┌────────────────────────────────────────────────────────────────┐

│                    Cloudflare Global Network                   │

│  ┌──────────────────────────────────────────────────────────┐  │

│  │  Cloudflare Load Balancer (DNS + Geo-Routing)           │  │

│  │  - Health checks every 60s                               │  │

│  │  - Geo-steering (sends traffic to nearest region)       │  │

│  │  - Automatic failover                                    │  │

│  └──────────────────────────────────────────────────────────┘  │

└────────────────────────────────────────────────────────────────┘

           │                      │                      │

           ▼                      ▼                      ▼

    ┌──────────┐          ┌──────────┐          ┌──────────┐

    │ US-EAST  │          │ EU-WEST  │          │ AP-SOUTH │

    │ Primary  │          │ Mirror   │          │ Mirror   │

    └──────────┘          └──────────┘          └──────────┘

    AWS Lightsail         AWS Lightsail         AWS Lightsail

    Virginia              London                Mumbai

    $160/month            $160/month            $160/month

```

  

### **How It Works**

  

**1. User Access** (e.g., you in Europe)

```

User in Paris

    ↓

Cloudflare DNS (geo-routing)

    ↓

Routes to EU-WEST (London) VM

    ↓

<50ms latency (vs. 150ms to US-EAST)

```

  

**2. Regional VM Setup** (Identical Deployments)

  

Each region runs **identical** Docker Compose stack:

  

```yaml

# ENVIRONMENT VARIABLES DIFFER PER REGION

# .env.us-east

REGION=us-east-1

DATABASE_URL=postgresql://us-east-db.internal

REDIS_URL=redis://us-east-redis.internal

  

# .env.eu-west

REGION=eu-west-1

DATABASE_URL=postgresql://eu-west-db.internal

REDIS_URL=redis://eu-west-redis.internal

```

  

**3. Database Replication Strategy**

  

**Option A: PostgreSQL Multi-Master (via pglogical or BDR)**

- Each region has read-write PostgreSQL

- Asynchronous replication (eventual consistency)

- Conflict resolution via "last-write-wins" or custom logic

  

**Option B: Supabase Multi-Region (Recommended)**

- Supabase Pro plan supports read replicas

- Primary in US-EAST (writes)

- Read replicas in EU-WEST, AP-SOUTH

- Automatic failover promotion

  

**Option C: CockroachDB (Best for Multi-Region)**

- Distributed SQL database

- Multi-region by design

- Automatic geo-partitioning

- Survives region failures

  

**4. Shared State Synchronization**

  

**Redis Cluster (for session data):**

```

US-EAST Redis Cluster (3 nodes)

    ↕ (replication)

EU-WEST Redis Cluster (3 nodes)

    ↕ (replication)

AP-SOUTH Redis Cluster (3 nodes)

```

  

**Letta Agent State:**

- Store conversation state in PostgreSQL (replicated)

- Store session cache in Redis (replicated)

- Store file uploads in S3-compatible storage (MinIO or AWS S3 with cross-region replication)

  

---

  

## **PORTAINER MULTI-VM MANAGEMENT**

  

### **How Portainer Manages Multiple VMs**

  

**Architecture:**

```

┌─────────────────────────────────────────┐

│  Portainer Server (on US-EAST VM)      │

│  Web UI: https://portainer.legacy.ai   │

└─────────────────────────────────────────┘

           │

           ├─── Manages US-EAST (local)

           │

           ├─── Manages EU-WEST (remote) ──► Portainer Agent

           │                                   (port 9001)

           │

           └─── Manages AP-SOUTH (remote) ─► Portainer Agent

                                              (port 9001)

```

  

**Setup Steps:**

  

**1. Install Portainer Server (US-EAST VM)**

```yaml

# docker-compose.yml (US-EAST)

services:

  portainer:

    image: portainer/portainer-ce:latest

    ports:

      - "9443:9443"  # HTTPS UI

      - "8000:8000"  # Edge agent tunnel

    volumes:

      - /var/run/docker.sock:/var/run/docker.sock

      - portainer_data:/data

    restart: always

```

  

**2. Install Portainer Agent (EU-WEST, AP-SOUTH VMs)**

```yaml

# docker-compose.yml (EU-WEST, AP-SOUTH)

services:

  portainer_agent:

    image: portainer/agent:latest

    ports:

      - "9001:9001"

    volumes:

      - /var/run/docker.sock:/var/run/docker.sock

      - /var/lib/docker/volumes:/var/lib/docker/volumes

    restart: always

```

  

**3. Add Remote Environments in Portainer UI**

```

Portainer UI → Environments → Add Environment

    → Agent

    → Name: EU-WEST

    → Endpoint URL: https://eu-west-vm.internal:9001

    → Add

```

  

**4. Manage All VMs from One Interface**

- Switch between environments via dropdown

- Deploy same stack to multiple environments

- Monitor resource usage across all VMs

  

---

  

## **EDGE DEVICE SYNCHRONIZATION**

  

### **How Edge Devices Connect to Regional VMs**

  

**Architecture:**

```

┌─────────────────────────────────────────────────────────┐

│  Edge Device (Jetson Nano at your home in US)          │

│  - Docker containers (offline-capable)                  │

│  - PowerSync client                                     │

│  - Cloudflare Tunnel (zerotrust connectivity)          │

└─────────────────────────────────────────────────────────┘

              │

              ▼

      Cloudflare Tunnel (encrypted)

              │

              ▼

      Cloudflare Load Balancer (geo-routing)

              │

              ▼

      Routes to nearest VM (US-EAST)

              │

              ▼

┌─────────────────────────────────────────────────────────┐

│  US-EAST VM (Primary)                                   │

│  - PowerSync server                                     │

│  - PostgreSQL (master)                                  │

│  - Letta API                                            │

└─────────────────────────────────────────────────────────┘

```

  

**When You Travel to Europe:**

```

┌─────────────────────────────────────────────────────────┐

│  Edge Device (Jetson Nano travels with you to Europe)  │

└─────────────────────────────────────────────────────────┘

              │

              ▼

      Cloudflare Tunnel (detects new location)

              │

              ▼

      Cloudflare Load Balancer (geo-routing)

              │

              ▼

      Routes to EU-WEST VM (closer)

              │

              ▼

┌─────────────────────────────────────────────────────────┐

│  EU-WEST VM (Mirror)                                    │

│  - PowerSync server                                     │

│  - PostgreSQL (replica)                                 │

│  - Letta API                                            │

└─────────────────────────────────────────────────────────┘

```

  

**Automatic Failover:**

- Cloudflare detects fastest route

- Edge device doesn't need reconfiguration

- Data syncs to nearest VM

- VM replicates to other regions

  

---

  

## **PORTAINER EDGE AGENT FOR JETSON NANO**

  

### **How to Manage Edge Devices via Portainer**

  

**Architecture:**

```

┌─────────────────────────────────────────┐

│  Portainer Server (US-EAST VM)          │

│  https://portainer.legacy.ai            │

└─────────────────────────────────────────┘

           │

           ├─── Manages US-EAST VM (local)

           ├─── Manages EU-WEST VM (agent)

           ├─── Manages AP-SOUTH VM (agent)

           │

           └─── Manages Jetson Nano (edge agent) ◄── NAT/Firewall OK

                     ↕ (polling connection)

                Jetson Nano at home

```

  

**Portainer Edge Agent Features:**

- ✅ **No inbound ports needed** (edge device polls Portainer server)

- ✅ **Works behind NAT/firewall** (no port forwarding)

- ✅ **Intermittent connectivity** (stores commands, executes when online)

- ✅ **Deploy stacks remotely** to Jetson Nano from Portainer UI

- ✅ **Monitor resource usage** (CPU, RAM, disk)

  

**Setup on Jetson Nano:**

```bash

# 1. Generate edge agent key in Portainer UI

# Portainer UI → Environments → Add Environment → Edge Agent

  

# 2. Deploy edge agent on Jetson Nano

docker run -d \

  -v /var/run/docker.sock:/var/run/docker.sock \

  -v /var/lib/docker/volumes:/var/lib/docker/volumes \

  -v /:/host \

  -v portainer_agent_data:/data \

  --restart always \

  -e EDGE=1 \

  -e EDGE_ID=jetson-nano-home \

  -e EDGE_KEY=<your-edge-key> \

  -e EDGE_INSECURE_POLL=1 \

  --name portainer_edge_agent \

  portainer/agent:latest

```

  

**Deploy Docker Compose to Jetson Nano from Portainer:**

```

Portainer UI → Environments → jetson-nano-home → Stacks → Add Stack

    → Name: legacy-ai-edge

    → Paste docker-compose-edge.yml

    → Deploy

```

  

**Portainer automatically:**

- Pulls images to Jetson Nano

- Starts containers

- Monitors health

- Shows logs

- Allows exec into containers

  

---

  

## **CLOUDFLARE CONFIGURATION**

  

### **Setup for Global Load Balancing**

  

**1. Cloudflare Load Balancer**

```

Cloudflare Dashboard → Traffic → Load Balancing

  

Create Pool: us-east-pool

  - Origin: us-east.legacy.internal

  - Health check: https://us-east.legacy.internal/health

  - Region: North America

  

Create Pool: eu-west-pool

  - Origin: eu-west.legacy.internal

  - Health check: https://eu-west.legacy.internal/health

  - Region: Western Europe

  

Create Pool: ap-south-pool

  - Origin: ap-south.legacy.internal

  - Health check: https://ap-south.legacy.internal/health

  - Region: Southern Asia

  

Create Load Balancer:

  - Hostname: legacy.ai

  - Default Pool: us-east-pool

  - Geo-Steering: ON

    - North America → us-east-pool

    - Europe → eu-west-pool

    - Asia → ap-south-pool

  - Fallback Pool: us-east-pool

```

  

**2. Cloudflare Tunnels (for Edge Devices)**

```bash

# On Jetson Nano

cloudflared tunnel create jetson-nano-home

cloudflared tunnel route dns jetson-nano-home jetson.legacy.ai

cloudflared tunnel run jetson-nano-home

```

  

**Config file (`~/.cloudflared/config.yml`):**

```yaml

tunnel: jetson-nano-home

credentials-file: /root/.cloudflared/<tunnel-id>.json

  

ingress:

  - hostname: jetson.legacy.ai

    service: http://localhost:8080  # Local Letta instance

  - service: http_status:404

```

  

**Benefits:**

- ✅ No public IP needed for Jetson Nano

- ✅ Encrypted tunnel to Cloudflare

- ✅ Access Jetson Nano from anywhere via `jetson.legacy.ai`

- ✅ Works behind NAT/firewall

  

---

  

## **FINAL RECOMMENDATIONS**

  

### **Phase 2A (Now): Single VM Setup**

```

✅ Deploy to single 64GB/16vCPU VM (US-EAST)

✅ All services in Docker Compose

✅ Portainer for management

✅ Telemetry on same VM (separate network)

✅ Caddy for reverse proxy + SSL

✅ Daily backups to S3

```

  

**Cost:** ~$160/month  

**Complexity:** Low  

**Maintenance:** Minimal  

  

### **Phase 2B (3-6 months): Add Edge Device**

```

🚀 Deploy Jetson Nano at home

🚀 Portainer Edge Agent for remote management

🚀 Cloudflare Tunnel for connectivity

🚀 PowerSync for offline-first data sync

🚀 Same Docker Compose stack (edge-optimized)

```

  

**Additional Cost:** ~$100 one-time (Jetson Nano)  

**Complexity:** Medium  

**Maintenance:** Low (managed via Portainer)  

  

### **Phase 3 (6-12 months): Add Regional Mirrors**

```

🌍 Deploy EU-WEST mirror (64GB VM)

🌍 Deploy AP-SOUTH mirror (64GB VM)

🌍 Cloudflare Load Balancer (geo-routing)

🌍 PostgreSQL multi-region replication

🌍 Portainer manages all VMs

```

  

**Additional Cost:** ~$320/month (2 more VMs)  

**Complexity:** High (database replication)  

**Maintenance:** Medium (3 VMs to maintain)

![[DECISION MATRIX.png]]

| Requirement | Single VM | Multi-VM | Multi-Region |

|-------------|-----------|----------|--------------|

| **Legacy AI self-access** | ✅ Excellent (localhost) | ⚠️ Complex (HTTP) | ⚠️ Complex (HTTP) |

| **Portainer management** | ✅ Simple (one env) | ⚠️ Moderate (agents) | ⚠️ Moderate (agents) |

| **Maintenance** | ✅ Low | ⚠️ Medium | ❌ High |

| **Cost** | ✅ $160/month | ⚠️ $240/month | ❌ $480/month |

| **Latency (global)** | ❌ High | ❌ High | ✅ Low (<50ms) |

| **Redundancy** | ❌ Single point of failure | ⚠️ Regional failure | ✅ Multi-region |

| **Edge device sync** | ✅ PowerSync + Tunnel | ✅ PowerSync + Tunnel | ✅ Geo-routed |

  

---

  

## **MY RECOMMENDATION**

  

### **Start with Single VM (Phase 2A)**

- Self-host everything on single 64GB VM

- Use Portainer for container management

- Keep telemetry on same VM (isolated network)

- Use Cloudflare for SSL + DDoS protection

- Use cloud APIs for LLMs (GitHub Models, Anthropic)

  

**When to expand:**

- ✅ Add Jetson Nano (Phase 2B) when you want offline capability

- ✅ Add regional mirrors (Phase 3) when you travel frequently to other continents

- ✅ Add dedicated GPU VM when you want self-hosted 33B+ LLMs

  

**You won't outgrow a single VM until you need:**

- More than 64GB RAM

- GPU acceleration for local LLMs

- <50ms global latency

- Multi-region redundancy

  

**For now, single VM is perfect.** 🎯
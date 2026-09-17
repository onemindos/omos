---
id: ui-15-20260111053926
title: "Neon vs Supabase for Letta.com AI Agents: Comprehensive Edge Deployment Analysis"
type: concept
domain: SO
domain: 00-24
folder: "00-24 SO (Sovereign Operations)"
status: active
created: 2026-01-11
modified: 2026-01-11
owner: zeus
tags: [domain/so, folder/15]
todoist_project: "00-24 SO (Sovereign Operations)"
ai_summary: "Neon vs Supabase for Letta.com AI Agents: Comprehensive Edge Deployment Analysis"
---

# Neon vs Supabase for Letta.com AI Agents: Comprehensive Edge Deployment Analysis

## Executive Summary

**Both Neon and Supabase are officially compatible with Letta.com**, but neither platform is architecturally designed for true edge, offline, or self-hosted deployment at scale. Neon excels at serverless cloud deployments with instant database branching, while Supabase provides a mature self-hosting option and superior real-time capabilities. **However, for your specific requirements—offline robots, AR glasses, off-grid scenarios, and sub-50ms edge latency—both platforms will require substantial additional infrastructure** using local-first sync engines like PowerSync or Electric SQL.

**Critical Finding:** Letta requires PostgreSQL with pgvector, and both platforms support this natively. The real challenge lies in your edge computing requirements, which neither platform addresses out-of-the-box.

---

## 1. Letta.com AI Agent Compatibility

### Official Support Status

**Letta explicitly endorses both platforms**. From Letta's blog post "AI Agents Stack":

> "Postgres-based companies like **Neon (serverless Postgres) and Supabase** also offer embedding-based search and storage for agents."

### Database Requirements

**Mandatory for Letta:**
- PostgreSQL (production-required; SQLite only for local development)
- **pgvector extension** (non-negotiable for agent memory and RAG)
- Database migrations handled automatically by Letta
- Persistent storage for agent state, messages, and archival memory

**Setup Simplicity:**
Both platforms require identical minimal configuration:

```bash
# For Neon or Supabase
CREATE EXTENSION vector;

# Set connection string
LETTA_PG_URI=postgresql://user:pass@host:5432/db
```

Letta handles all table creation and schema migrations automatically when connecting to either platform.

### Community Verification

**Supabase:** Multiple detailed community guides exist (Medium, December 2024). Users report setup takes "minutes" with automatic migrations.

**Neon:** Officially mentioned in Letta documentation. Works seamlessly with Neon's serverless architecture, popular for AI agent deployments (80% of Neon databases created by AI agents).

### Why Letta Needs These Features

Letta stores agent memories as vector embeddings in the database, requiring:
- **Vector similarity search** via pgvector HNSW indexes
- **ACID transactions** for consistent agent state
- **Full-text search** for conversation history
- **JSON support** for structured metadata
- **High-frequency reads/writes** for active agent conversations

**Verdict:** Both platforms are equally compatible with Letta's core requirements. Compatibility is not a differentiator.

---

## 2. Stack Component Coverage Comparison

Your Airbyte architecture requires 8 components. Here's how each platform addresses them:

| Component | Neon | Supabase | Winner |
|-----------|------|----------|--------|
| **1. Primary OLTP DB** | ✅ Native Postgres 14-17 | ✅ Native Postgres 14.6+ | **TIE** |
| **2. pgvector Store** | ✅ Native, limited benchmarks | ✅ Native, 1800 QPS @ 16XL proven | **Supabase** |
| **3. Redis Cache/Queue** | ❌ External required | ❌ External required (Wrappers available) | **TIE** |
| **4. Object Storage** | ⚠️ Internal S3, apps need external | ✅ Native (Supabase Storage, S3-compatible) | **Supabase** |
| **5. OpenSearch** | ❌ External required | ⚠️ Postgres FTS native, not OpenSearch | **Neither** |
| **6. CDC for Airbyte** | ✅ Logical replication, official connector | ✅ Logical replication, official connector | **TIE** |
| **7. Analytics (Iceberg/Trino)** | ❌ External (OLTP-only) | ⚠️ Logflare-based (requires setup) | **Supabase** |
| **8. Observability** | ⚠️ Basic console + external | ✅ Logflare + Prometheus + Log Drains | **Supabase** |

### Detailed Stack Analysis

**Postgres OLTP (Native - Both Platforms)**
- Neon: Separated storage/compute, auto-scales 0.25-56 CU, scale-to-zero
- Supabase: Traditional Postgres, predictable performance, always-on paid tier
- Both support full SQL, ACID, connection pooling

**pgvector (Native - Supabase Stronger)**
- **Neon**: Extension available, suitable for moderate workloads (\<10M vectors), no comprehensive benchmarks
- **Supabase**: Extensively benchmarked—1800 QPS at 16XL compute, 1M vectors tested, HNSW optimization guides, 28x lower latency vs Pinecone

**Redis Cache/Queue (External - Both Platforms)**
- Neither provides native Redis
- Must integrate Upstash, ElastiCache, or self-hosted Redis
- Supabase offers Wrappers extension to query Redis from Postgres (partial mitigation)
- **Impact:** Adds complexity and cost to your n8n worker queue architecture

**Object Storage (Supabase Advantage)**
- **Neon**: Uses S3 internally, but apps need separate S3 buckets
- **Supabase**: Native Storage API with S3-compatible backend, RLS policies, CDN integration, image transformations
- For your "files, docs, artifacts, exports" requirement, Supabase provides integrated solution

**OpenSearch (Neither Platform)**
- **Neon**: No search beyond basic Postgres full-text
- **Supabase**: Native Postgres FTS (GIN/GiST indexes, `tsvector`), but NOT OpenSearch/Elasticsearch
- Both require external OpenSearch service for your "full-text + logs, hybrid vector/text search" requirement

**CDC for Airbyte (Both Platforms)**
- Both support Postgres logical replication natively
- Official Airbyte connectors available for both
- Enables your "Postgres source → S3 (Parquet)" pipeline
- Requires logical replication setup (replication slots, publications)

**Analytics (Supabase Slight Edge)**
- **Neon**: OLTP-only, must replicate to external analytics DB
- **Supabase**: Logflare-based analytics with BigQuery backend (cloud), but self-hosted setup complex
- Neither replaces your planned "Iceberg on S3 + Trino" architecture

**Observability (Supabase Stronger)**
- **Neon**: Basic console metrics, logs available, requires external APM
- **Supabase**: Comprehensive Logflare integration, Prometheus metrics, Log Drains to Datadog, 200+ Grafana metrics, AI log analysis
- For Langfuse integration, both work equally well as external services

### Stack Coverage Score
- **Neon:** 3/8 native, 2/8 integration, 3/8 external = **37.5% coverage**
- **Supabase:** 5/8 native, 2/8 integration, 1/8 external = **62.5% coverage**

**Supabase provides significantly more integrated stack components**, reducing the number of separate services you must manage.

---

## 3. Edge Computing & Offline Deployment Reality Check

**Critical Finding:** Neither platform is designed for true edge deployment or offline operation. Your requirements for "robots, AR glasses, off-grid scenarios" fundamentally conflict with both platforms' cloud-centric architectures.

### Edge Deployment Comparison

| Capability | Neon | Supabase | Reality |
|------------|------|----------|---------|
| **Database at Edge** | ❌ Cloud-only in AWS regions | ❌ Cloud-only, single region per project | **Neither supports** |
| **Edge Compute** | ✅ HTTP/WebSocket driver for edge runtimes | ✅ Edge Functions (29 regions) | **Both support compute, not data** |
| **Cross-Region Read Replicas** | ⚠️ Same-region only, planned future | ✅ Read replicas (Pro+), async replication | **Supabase limited** |
| **Offline Operation** | ❌ Requires internet | ❌ Requires internet | **Neither supports natively** |
| **Local-First Sync** | ❌ No native support | ❌ No native support | **Both require PowerSync/Electric** |
| **Multi-Region Writes** | ❌ Single primary region | ❌ Single primary region | **Neither supports** |

### What "Edge" Actually Means for Each Platform

**Neon's Edge Capability:**
- **Edge driver** (@neondatabase/serverless) runs on Cloudflare Workers, Vercel Edge
- Uses HTTP/WebSocket instead of TCP for edge runtime compatibility
- **Database remains centralized** in AWS us-west-2, us-east-1, or eu-west-1
- Result: Edge compute with cloud database = 100-300ms cross-region latency

**Supabase's Edge Capability:**
- **Edge Functions** (Deno-based) deploy globally to 29 regions
- Database remains in chosen region (14+ AWS regions available)
- Read replicas enable geo-distributed reads (Pro+ tier)
- Result: Edge functions with regional database = 100-250ms for distant users

**Ground Truth:** Both platforms provide "edge compute" but not "edge data." Your robots and AR glasses would query cloud databases hundreds of milliseconds away.

### Offline & Local-First Requirements

For your stated requirements—**"offline/online, edge devices (robots, glasses), off-grid scenarios"**—both platforms require identical additional architecture:

**Required Third-Party Solutions:**

1. **PowerSync** (Recommended for Production)
   - Service layer between Postgres and edge devices
   - SQLite on device with automatic sync
   - Bidirectional sync with conflict resolution
   - Works with both Neon and Supabase
   - Cost: Additional service + subscription ($50+/month)

2. **Electric SQL** (Modern Alternative)
   - Postgres CDC → PGlite/SQLite on edge
   - HTTP-based shape sync
   - Works with both platforms
   - More open-source friendly

3. **ObjectBox** (Resource-Constrained Devices)
   - Embedded database for robots/IoT
   - On-device vector search
   - P2P sync capabilities
   - Completely separate from Neon/Supabase

**Architecture Reality for Your Use Case:**

```
Edge Device (Robot/Glasses)
  ├─ Local DB (SQLite via PowerSync or ObjectBox)
  ├─ AI Model (TensorFlow Lite)
  ├─ Sync Engine (PowerSync/Electric SQL)
  └─ Network: Intermittent
       ↓ When online
Cloud Database (Neon or Supabase)
  ├─ Postgres with pgvector
  ├─ Logical replication enabled
  └─ Sync service (PowerSync/Electric SQL)
```

**Bottom Line:** Your edge/offline requirements add $5,000-$15,000 in development time plus ongoing service costs, **regardless of choosing Neon or Supabase**. This is the most significant architectural consideration.

---

## 4. Self-Hosting & Off-Grid Deployment

This is where the platforms dramatically diverge.

### Self-Hosting Maturity

| Factor | Neon | Supabase | Winner |
|--------|------|----------|--------|
| **Feasibility** | Possible but not production-ready | Production-ready | **Supabase** |
| **Documentation** | Minimal, unsupported | Excellent, official guides | **Supabase** |
| **Deployment Method** | Build from source or Kubernetes operator (alpha) | Docker Compose (minutes) | **Supabase** |
| **Components Included** | Pageserver, Safekeepers, Compute, Storage Controller | All services (Auth, Storage, Realtime, Postgres) | **Supabase** |
| **Infrastructure Complexity** | High (requires S3, multiple components) | Low (single Docker Compose) | **Supabase** |
| **Official Support** | None for self-hosted | Community support, self-hosting guides | **Supabase** |
| **Update Process** | Manual, complex | Manual but documented | **Supabase** |
| **Cost (8GB RAM VPS)** | $40-80/month (if you can get it working) | $40-80/month | **TIE** |

### Neon Self-Hosting Reality

**Status:** Technically possible, practically not recommended for production.

- Open-source Apache 2.0 license (code available)
- Requires building from source with Rust toolchain, protobuf, PostgreSQL
- **Architecture complexity:** Pageserver (storage backend) + 3+ Safekeepers (redundant WAL) + Compute nodes + Storage Controller + S3-compatible storage
- Community Kubernetes operator exists but **early alpha stage** (not production-ready)
- **No official support** for self-hosted deployments
- Missing features: Automated upgrades, monitoring integrations, orchestration

**Verdict:** Neon is fundamentally designed as a managed cloud service. Self-hosting undermines its core value proposition (serverless, instant provisioning) and introduces operational complexity comparable to managing a distributed database cluster.

### Supabase Self-Hosting Reality

**Status:** Production-ready with mature tooling.

**Deployment Options:**
- **Docker Compose:** Official method, setup in minutes
- **Kubernetes:** Community Helm charts available
- **Cloud providers:** Templates for AWS, GCP, Azure

**What You Get:**
All components self-hosted: Postgres, PostgREST (REST API), GoTrue (Auth), Realtime (WebSockets), Storage (S3-compatible), Kong (API Gateway), postgres-meta (admin), Supavisor (connection pooler)

**Operational Requirements:**
- Manual updates (no auto-updates like cloud)
- DIY backups and monitoring
- Analytics requires additional BigQuery or Postgres backend
- No auto-scaling (manual configuration)

**Cost Comparison:**
- Self-hosted VPS (8GB RAM): $40-80/month
- Supabase Cloud Pro: $25/month + usage (typically $100-300/month)
- **Savings at scale:** $100-220/month for moderate workloads
- **Hidden cost:** DevOps time for maintenance

### Off-Grid & Air-Gapped Scenarios

**For completely offline deployments (ships, remote facilities, secure environments):**

**Supabase:** Fully supports air-gapped deployment via self-hosting. All components can run on local hardware without internet. Ideal for your "off-grid scenarios" requirement.

**Neon:** Theoretically possible but requires significant infrastructure expertise. Not recommended for air-gapped deployments due to complexity and lack of official support.

**Recommendation:** If self-hosting or off-grid deployment is a requirement, **Supabase is the only viable option between these two platforms**.

---

## 5. Real-Time Performance for Edge Devices

Your requirement: **"Sub-50ms latency for real-time robot/glasses control"**

### Latency Benchmarks

**Same-Region Performance (Best Case):**

| Scenario | Neon | Supabase | Achieves \<50ms? |
|----------|------|----------|-----------------|
| **Hot query (TCP)** | 2ms avg, 9ms p99 | 1ms avg, 4ms p99 | ✅ Both |
| **Hot query (HTTP)** | 5ms avg, 19ms p99 | 38ms avg (with auth) | ⚠️ Neon only |
| **Cold start** | 500ms avg, 1200ms p99 | N/A (paid tier) | ❌ Neon fails |
| **WebSocket message** | N/A (no native) | 6ms median, 28ms p95 | ✅ Supabase |
| **Realtime broadcast** | N/A | 6ms median (224k msgs/sec) | ✅ Supabase |

**Cross-Region Performance (Your Edge Devices):**

| Location | Neon (from US-East DB) | Supabase (from US-East DB) | Achieves \<50ms? |
|----------|------------------------|---------------------------|-----------------|
| **Same region** | 5-15ms | 11-26ms | ✅ Both |
| **US-West** | 60-80ms | ~70ms | ❌ Both fail |
| **Europe** | 100-150ms | 100-150ms | ❌ Both fail |
| **Asia** | 200-300ms | 200-300ms | ❌ Both fail |

### Critical Findings

**Sub-50ms latency is achievable ONLY when:**
1. ✅ Application and database are in the **same AWS region**
2. ✅ Database compute is **warm** (not cold start)
3. ✅ Using **optimized connection** (TCP pooled or HTTP)
4. ✅ Queries are **properly indexed**

**Sub-50ms latency is NOT achievable when:**
1. ❌ Edge devices in different regions than database
2. ❌ Neon database experiences cold start (500ms penalty)
3. ❌ Using Supabase HTTP with auth middleware (38ms overhead)
4. ❌ Sequential queries (5 queries × 250ms = 1250ms total)

### Real-World Scenario: Robot in Factory

**Scenario:** Robot with sensor requiring live updates every 50ms for control feedback.

**With Neon (cloud database in US-East):**
- Robot in US-East factory: ✅ 5-15ms (works)
- Robot in Europe factory: ❌ 150ms (fails requirement)
- Robot on autonomous vehicle: ❌ 200ms+ with cellular (fails)
- Cold start after idle: ❌ 500ms (critical failure for real-time control)

**With Supabase (cloud database in US-East):**
- Robot in US-East factory: ✅ 11ms (works)
- Robot in Europe factory: ❌ 150ms (fails requirement)
- Robot on autonomous vehicle: ❌ 200ms+ with cellular (fails)
- Realtime WebSocket: ✅ 6ms median for state updates (if in same region)

**With Local-First Architecture (Required for Your Use Case):**
- Local database (SQLite/ObjectBox) on robot: ✅ 1-5ms (works anywhere)
- Background sync to cloud: Eventual consistency (acceptable)
- Works offline completely: ✅ Critical for intermittent connectivity

### Live Data Updates & Sensor Readings

**Telemetry Streaming Capabilities:**

**Supabase Realtime:**
- 224,000 messages/sec (Broadcast mode)
- 10,000 inserts/sec → 80,000 concurrent clients (DB triggers)
- 6ms median latency for WebSocket delivery
- Ideal for dashboards, monitoring, live coordination

**Neon:**
- No built-in real-time streaming
- Requires external service (Ably, Pusher, Socket.io)
- Logical replication available for CDC to external systems
- Suitable for batch analytics, not live updates

**For Your Robot/Glasses Sensors:**
- High-frequency sensor data (100 Hz+): **Requires local processing**, sync batches
- Live telemetry dashboard: **Supabase Realtime** excels
- Control commands: **Must be local-first** (\<50ms not achievable cloud-only)

### Verdict: Sub-50ms Edge Latency

**Neither platform achieves sub-50ms globally**. For real-time robot/glasses control, you **must implement local-first architecture** with on-device databases (ObjectBox, SQLite, PGlite) and background sync to cloud (Neon or Supabase as backend).

---

## 6. pgvector Performance at Edge

Your use case requires vector embeddings for "RAG + long-term memory."

### pgvector Benchmarks

**Supabase (Extensively Tested):**
- **1M OpenAI embeddings** (1536 dimensions):
  - 16XL compute: 1800 QPS @ 10 probes, 91% accuracy
  - 16XL compute: 670 QPS @ 40 probes, 98% accuracy
  - Latency: 5-15ms median, 50-100ms p99 under load
- **Optimization:** HNSW index, lists = rows/200, pre-warm with 10k queries
- **Memory:** 50M vectors require ~850GB RAM (index must fit in memory)
- **Proven scale:** Production workloads with 1.6M+ embeddings

**Neon (Limited Benchmarks):**
- pgvector extension available
- Performance comparable to vanilla Postgres with pgvector
- Suitable for moderate workloads (\<10M vectors)
- No comprehensive public benchmarks at scale
- Remote storage architecture may introduce slight overhead

**Third-Party Comparisons:**
- Supabase pgvector: 4x better QPS than Pinecone p1.x1
- Supabase pgvector: 28x lower p95 latency vs Pinecone s1 (50M vectors)
- Qdrant: 15x faster than pgvector for pure vector workloads (specialized DB advantage)

### Edge Device Constraints

**Your robots and AR glasses likely have:**
- ARM processors (1-4 GHz)
- 512MB - 8GB RAM
- 8GB - 64GB storage
- Battery constraints

**pgvector on Edge:**
- HNSW index for 50k vectors: ~150MB RAM
- Query latency with PGlite (WASM Postgres): 3-10ms
- Suitable for local RAG with limited context windows

**Recommended Edge Vector Architecture:**
```
Edge Device:
  ├─ PGlite or SQLite with pgvector
  ├─ Sync relevant embeddings only (shape-based)
  ├─ Keep index \<100K vectors
  └─ Query locally (1-5ms latency)

Cloud (Neon or Supabase):
  ├─ Full embedding corpus (millions)
  ├─ HNSW indexed
  └─ Sync subsets to edge devices
```

### CDC for Edge Sync

Both platforms support CDC for syncing vector embeddings to edge:

**Logical Replication Performance:**
- Initial sync: 50M vectors = 10-30 minutes (depends on bandwidth)
- Incremental sync: \<1s for typical workloads (\<1000 changes/sec)
- Electric SQL or PowerSync can sync pgvector tables

**Conflict Resolution:**
- Vector embeddings typically append-only (no conflicts)
- Use last-write-wins for metadata updates
- Server reconciliation for critical updates

**Verdict:** For edge vector search, implement local embeddings with sync to cloud. Both Neon and Supabase support this pattern equally via logical replication.

---

## 7. Cost Analysis for Edge AI Agents

### Cloud Cost Comparison

**Neon Pricing (Scale-to-Zero Model):**
- Free tier: 0.5 GB, 190 compute hours @ 0.25 CU
- Launch ($19/month minimum): $0.14/CU-hour, $0.35/GB storage
- Scale: $0.26/CU-hour
- **Advantage:** Only pay for active compute (scale-to-zero)

**Supabase Pricing (Always-On Model):**
- Free tier: 500 MB, pauses after 7 days
- Pro ($25/month): 8 GB database, no pause
- Compute add-ons: Small (4GB) $50/mo, Medium (8GB) $100/mo
- **Advantage:** Predictable pricing, included auth/storage/realtime

**Example: 10 AI Agents**

| Scenario | Neon Cost | Supabase Cost | Winner |
|----------|-----------|---------------|--------|
| **Intermittent (2 hrs/day each)** | ~$150/mo (scale-to-zero saves 90%) | ~$125/mo (Pro + Small compute) | **Neon** |
| **Always-On (24/7)** | ~$1,000/mo (no scale-to-zero benefit) | ~$225/mo (Pro + Medium compute) | **Supabase** |
| **Variable Workload** | ~$300/mo (auto-scales) | ~$225/mo (fixed tier) | **Supabase** |

**Edge Architecture Costs (Additional for Both):**
- PowerSync service: $50-200/month
- Redis (Upstash): $20-100/month
- OpenSearch (managed): $100-500/month
- ObjectBox (for devices): One-time or per-device license
- **Total additional:** $200-800/month for edge/offline capabilities

### Self-Hosting Cost

**VPS Requirements for Your Stack:**
- Development: 2-4 GB RAM, $10-20/month
- Production: 16-32 GB RAM, $80-150/month
- High-performance: 64 GB RAM, $200-400/month

**Self-Hosting Trade-offs:**
- **Savings:** $100-300/month vs cloud at scale
- **Hidden cost:** DevOps time (10-20 hours/month)
- **Risk:** No auto-scaling, manual backups, downtime during updates

### TCO Recommendation

**Choose Neon Cloud if:**
- Highly variable workloads (scale-to-zero saves money)
- Development/staging environments with low activity
- Want instant database branching for testing

**Choose Supabase Cloud if:**
- Steady always-on workloads (more cost-effective at scale)
- Want included auth, storage, realtime (reduces separate service costs)
- Prefer predictable monthly bills

**Choose Supabase Self-Hosted if:**
- Cost savings priority at scale (\>$300/month cloud spend)
- Data sovereignty requirements
- Have DevOps expertise
- Air-gapped deployment needed

---

## 8. Final Recommendation Matrix

### Choose Neon If You Need:

✅ **Cloud-Native AI Agents**
- Instant database provisioning for ephemeral agents
- Scale-to-zero for cost efficiency
- Database branching for testing/staging
- 80% of Neon databases created by AI agents—optimized for this

✅ **Development Velocity**
- Copy-on-write branches in \<1 second
- Preview environments for every PR
- Schema changes without data duplication

✅ **Variable Workloads**
- Auto-scaling from 0.25 to 56 CU
- Pay only for active time
- Cost optimization for intermittent agents

❌ **Do NOT Choose Neon If You Need:**
- Self-hosting (not production-ready)
- Built-in real-time features (requires external service)
- Always-on low latency (cold starts problematic)
- Mature edge deployment (data stays central)
- Extensive pgvector benchmarks (limited public data)

### Choose Supabase If You Need:

✅ **Full-Stack Platform**
- Postgres + Auth + Storage + Realtime in one service
- Mature self-hosting option (Docker Compose)
- Comprehensive observability and monitoring

✅ **Real-Time Capabilities**
- Built-in WebSocket pub/sub (800k+ msgs/sec)
- Presence tracking for agent coordination
- Database change streaming with RLS

✅ **Vector Search at Scale**
- Extensively benchmarked pgvector (1800 QPS proven)
- Production-tested with 1.6M+ embeddings
- Detailed optimization guides

✅ **Self-Hosting Requirements**
- Production-ready Docker deployment
- Air-gapped capability
- Full feature parity with cloud

❌ **Do NOT Choose Supabase If You Need:**
- Scale-to-zero (always-on on paid plans)
- Instant database branching (limited to schema)
- Absolute lowest cost for variable workloads (Neon wins)

### For Your Specific Requirements

**Your Stack: "OLTP + pgvector + Redis + S3 + OpenSearch + CDC + Analytics + Observability"**

| Component | Winner | Reasoning |
|-----------|--------|-----------|
| OLTP Postgres | TIE | Both excellent |
| pgvector | **Supabase** | Proven benchmarks, optimization guides |
| Redis | TIE (both external) | Neither provides native |
| Object Storage | **Supabase** | Native Storage API vs external only |
| OpenSearch | TIE (both external) | Neither provides OpenSearch |
| CDC/Airbyte | TIE | Both support logical replication |
| Analytics | **Supabase** | Logflare vs nothing |
| Observability | **Supabase** | Comprehensive vs basic |

**Winner for Stack Coverage: Supabase (62.5% vs 37.5%)**

**Your Requirements: "Offline/online, edge devices (robots, glasses), off-grid, sub-50ms"**

| Requirement | Winner | Reasoning |
|-------------|--------|-----------|
| Edge database | **NEITHER** | Both keep data centralized in cloud |
| Offline operation | **NEITHER** | Both require PowerSync/Electric SQL |
| Self-hosting | **Supabase** | Production-ready vs alpha/unsupported |
| Sub-50ms latency | **NEITHER** | Only achievable with local-first architecture |
| Real-time updates | **Supabase** | Built-in WebSocket vs external service needed |

**Winner for Edge Requirements: Supabase (marginally better due to self-hosting)**

---

## 9. Recommended Architecture for Your Use Case

Given your requirements for "all levels of Letta AI agent deployment, offline/online, edge devices (robots, glasses), off-grid scenarios, and real-time control + live updates," here's the optimal architecture:

### Cloud Backend Layer
```
Cloud (Choose One):
├─ Supabase (Recommended)
│  ├─ Postgres with pgvector (agent memory)
│  ├─ Realtime (agent coordination)
│  ├─ Storage (files, artifacts)
│  ├─ Logical replication enabled
│  └─ Self-hostable for off-grid scenarios
│
OR
│
├─ Neon (Alternative)
│  ├─ Postgres with pgvector (agent memory)
│  ├─ Instant branching (dev/test)
│  └─ Scale-to-zero (cost optimization)
```

### Sync Layer (Required for Edge/Offline)
```
Sync Service (Choose One):
├─ PowerSync (Production-Ready)
│  ├─ Postgres CDC → SQLite on devices
│  ├─ Server reconciliation
│  ├─ Conflict resolution
│  └─ Cost: $50-200/month + service
│
OR
│
├─ Electric SQL (Open-Source)
│  ├─ Shape-based sync
│  ├─ PGlite on edge
│  └─ HTTP-based (scalable)
```

### Edge Device Layer (Required for Your Use Case)
```
Robots/AR Glasses:
├─ Local DB (choose based on constraints)
│  ├─ PGlite (Postgres-compatible, WASM)
│  ├─ SQLite (universal, lightweight)
│  └─ ObjectBox (resource-constrained, vectors)
│
├─ AI Model (TensorFlow Lite, ONNX)
├─ Local pgvector (\<100K embeddings)
├─ Sync engine (PowerSync/Electric client)
└─ Control loop (1-5ms latency, local-only)
```

### Additional Services (External to Both Platforms)
```
Required Services:
├─ Redis: Upstash (edge-compatible) or ElastiCache
├─ OpenSearch: Amazon OpenSearch or self-hosted
├─ Airbyte: For CDC to S3/Parquet
└─ Langfuse: Docker deployment for LLM observability
```

### Architecture Decision Tree

**For Cloud-Only Letta Deployment (No Edge Requirements):**
→ Use **Neon** for cost efficiency with scale-to-zero

**For Cloud + Development/Staging Branching:**
→ Use **Neon** for instant branch creation

**For Cloud + Real-Time + Full-Stack:**
→ Use **Supabase** for integrated platform

**For Self-Hosted/Off-Grid:**
→ Use **Supabase** (only production-ready option)

**For Edge/Offline/Robots/Glasses:**
→ Use **Supabase or Neon** as backend + **PowerSync/Electric SQL** + **Local database** on devices (required architecture regardless of choice)

---

## 10. Bottom Line: Definitive Recommendation

**For your specific requirements—Letta AI agents with offline robots, AR glasses, off-grid scenarios, and sub-50ms edge latency:**

### Primary Recommendation: **Supabase**

**Why Supabase:**
1. **Stack coverage:** 62.5% of components native vs 37.5% for Neon
2. **Self-hosting maturity:** Production-ready Docker deployment for off-grid scenarios
3. **Real-time capabilities:** Built-in WebSocket pub/sub for agent coordination
4. **pgvector performance:** Extensively benchmarked and proven at scale
5. **Observability:** Comprehensive monitoring for LLM operations
6. **Total cost:** Lower at always-on workloads after edge services added

**Supabase enables your edge/offline requirements through:**
- Mature self-hosting (air-gapped deployment possible)
- Native Storage API (reduces external services)
- Built-in real-time (agent coordination without external WebSocket service)
- Better analytics and observability (Langfuse integration + native logs)

### Alternative: **Neon + Significant Additional Architecture**

**Neon is viable if:**
- Cost optimization is primary concern (scale-to-zero for variable workloads)
- You heavily use database branching for development workflows
- You're willing to add external services for real-time, storage, analytics
- Your deployment is cloud-only (not self-hosted)

### The Hard Truth About Edge Computing

**Neither Neon nor Supabase solves your edge/offline requirements natively**. For robots, AR glasses, and off-grid scenarios requiring sub-50ms latency and offline operation, you **must** implement:

1. **Local-first architecture** (PowerSync, Electric SQL, or ObjectBox)
2. **On-device databases** (SQLite, PGlite, ObjectBox)
3. **Sync service layer** ($50-200/month additional cost)
4. **Development complexity** ($10,000-$20,000 additional development)

Both platforms serve as cloud backends in this architecture—**the edge capability comes from your local-first sync layer, not the database platform**.

### Implementation Roadmap

**Phase 1: Cloud MVP (2-4 weeks)**
- Deploy Supabase (cloud or self-hosted)
- Set up Letta with pgvector
- Integrate Langfuse for observability
- Deploy external Redis and OpenSearch

**Phase 2: Edge Sync (4-8 weeks)**
- Integrate PowerSync or Electric SQL
- Implement local databases on test devices
- Set up CDC and conflict resolution
- Test offline scenarios

**Phase 3: Production Edge (4-6 weeks)**
- Deploy to robot/glasses hardware
- Optimize vector embeddings for edge
- Implement real-time control loops
- Off-grid testing and validation

**Total Timeline:** 10-18 weeks for complete edge deployment

### Final Word

**Supabase is the recommended platform** for your Letta AI agent deployment due to superior stack coverage, production-ready self-hosting, and mature real-time capabilities. However, achieving your edge computing goals requires substantial additional architecture beyond either platform's native capabilities. Budget accordingly for local-first sync infrastructure, which will cost more in development time and ongoing services than the database choice itself.

The database platform decision is important but secondary to your edge architecture strategy. **Invest your energy in designing robust local-first sync patterns**, as this is where the real complexity and value of your edge AI agent system lies.
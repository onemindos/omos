---
id: so-01-langfuse-20260527
title: "Langfuse — LLM Observability Platform"
type: service
domain: SO
folder: "00-24 SO (Sovereign Operations)/01 Infrastructure (Grid & APIs)"
status: active
created: 2026-05-27
modified: 2026-05-27
owner: zeus
tags: [domain/so, folder/01, service/langfuse, ai/observability, ai/agents]
---

# Langfuse — LLM Observability Platform

> Meta-layer for managing, observing, evaluating, and promoting agents. Connects to all OneMind agents via SDK instrumentation.

---

## Access

| Interface | URL |
|-----------|-----|
| **Web UI** | http://localhost:3004 |
| **Worker** | http://localhost:3032 (internal) |
| **MinIO Console** | http://localhost:9091 (internal) |

---

## Stack Location

```
~/onemind/langfuse/
├── docker-compose.yml   ← stack definition (web on :3004)
└── .env                 ← secrets (NEXTAUTH_SECRET, SALT, ENCRYPTION_KEY)
```

---

## Services (Docker)

| Container | Image | Port | Purpose |
|-----------|-------|------|---------|
| langfuse-web | langfuse/langfuse:3 | 3004 | UI + API |
| langfuse-worker | langfuse/langfuse-worker:3 | 3032 | Background jobs |
| clickhouse | clickhouse/clickhouse-server | localhost only | Analytics DB |
| postgres | postgres:17 | localhost only | Primary DB |
| redis | redis:7 | localhost only | Queue |
| minio | chainguard/minio | 9090 (S3 API) | Blob storage |

---

## Commands

```bash
# Start
cd ~/onemind/langfuse && docker compose up -d

# Stop
cd ~/onemind/langfuse && docker compose down

# Logs
cd ~/onemind/langfuse && docker compose logs -f langfuse-web

# Status
cd ~/onemind/langfuse && docker compose ps
```

---

## Agent Integration

### Environment Variables (add to each agent)
```bash
LANGFUSE_PUBLIC_KEY=pk-lf-...      # from Langfuse UI → Settings → API Keys
LANGFUSE_SECRET_KEY=sk-lf-...
LANGFUSE_HOST=http://localhost:3004
```

### Python — Drop-in (zero code change)
```python
from langfuse.anthropic import anthropic  # replaces: import anthropic
```

### Python — Decorator (full control)
```python
from langfuse.decorators import observe

@observe()
def zenith_agent(input: str):
    # existing agent code unchanged
    ...
```

### Hermes Orchestrator Pattern
```python
@observe()
def hermes(input: str):           # parent trace
    zenith_result = zenith_agent(input)    # child span
    unity_result = unity_agent(input)      # child span
    return synthesize(zenith_result, unity_result)
```

---

## Features In Use

| Feature | Purpose |
|---------|---------|
| **Tracing** | Every agent call logged — inputs, outputs, latency, tokens, cost |
| **Playground** | Test prompts against live models before deploying |
| **Prompt Management** | Version-controlled prompts, fetch at runtime |
| **Datasets** | Build test sets from production traces |
| **Evals** | LLM-as-judge scoring (hallucination, relevance, correctness) |
| **Experiments** | Compare agent v1 vs v2 on same dataset |
| **Dashboards** | Cost, latency, performance over time |
| **Human Annotation** | Manual review queue for edge cases |

---

## Playground → Production Pipeline

```
Agent in playground (localhost:3004)
    → Langfuse captures trace
    → Auto-eval scores output
    → Human review for edge cases
    → Build dataset from failures
    → Run experiment (v1 vs v2)
    → Promote winner to production agents
```

---

## Cloudflare Tunnel DNS

To expose Langfuse externally via `onemindos.dev`:

Add to `~/.cloudflared/config.yml`:
```yaml
  - hostname: langfuse.onemindos.dev
    service: http://localhost:3004
    originRequest:
      httpHostHeader: langfuse.onemindos.dev
```

Then add CNAME in Cloudflare DNS:
```
langfuse.onemindos.dev  →  <tunnel-id>.cfargotunnel.com
```

---

## Related

- [[NETWORK]] — NetBird mesh, Cloudflare Tunnel config
- [[LAN-SERVICES]] — All local services map
- Agents: Legacy, Zenith, Unity, Eternal, Sovereignty, Sentinel, Hermes

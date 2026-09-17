---
title: "Langfuse"
type: service-card
status: active
created: 2026-05-01
updated: 2026-06-14
owner: zeus
category: monitoring
tags: [platform, ai, tools, observability, tracing, self-hosted, vps, docker]
ai_summary: "Self-hosted LLM observability platform (traces, evals, cost). Runs on omos-cloud (Lightsail) Docker at :3000 — 6 containers, not on the Mac."
---
# Langfuse

> LLM observability platform. Trace, evaluate, and monitor AI applications. Self-hosted on Docker.

## Hosting
| Field | Value |
|-------|-------|
| Host | **omos-cloud** (AWS Lightsail) — *not the Mac* |
| Tailscale IP | 100.111.189.116 |
| Public IP | 3.224.147.141 |
| Runtime | Docker Compose |
| URL (mesh) | `http://100.111.189.116:3000` |
| Port | 3000 |
| Uptime | stable (7+ days as of 2026-06-14) |

## Containers (6)
| Container | Role | State |
|-----------|------|-------|
| langfuse-langfuse-web-1 | web UI / API (:3000) | up |
| langfuse-langfuse-worker-1 | async worker | up |
| langfuse-postgres-1 | Postgres | healthy |
| langfuse-redis-1 | Redis | healthy |
| langfuse-clickhouse-1 | ClickHouse (analytics) | healthy |
| langfuse-minio-1 | MinIO (object store) | healthy |

> Note: the local `/Users/onemind/langfuse` folder is just a repo clone — the live stack runs on omos-cloud, not the Mac.

## Key Features
| Feature | Notes |
|---------|-------|
| Traces | Full request/response logging per LLM call |
| Scores | Evaluate quality (manual + automated) |
| Datasets | Create eval datasets from production traces |
| Prompts | Version-controlled prompt management |
| Cost tracking | Token usage + cost per model |
| Dashboards | Latency, error rates, costs |
| Sessions | Group traces by user session |

## Integration
- OpenAI SDK wrapper (drop-in)
- LangChain callback handler
- Hermes: trace all agent tool calls + LLM completions
- NATS events → Langfuse for fleet-wide observability

---

## Related
- [[00-24 SO (Sovereign Operations)/01 The Fabric (Node Registry)/01 The Fabric|The Fabric]]

---
title: "Kestra"
type: reference
status: inactive
tags: [platform, infrastructure, orchestration, workflows]
---
# Kestra

> Open-source workflow orchestration platform — the modern Airflow/n8n alternative.

## What It Does

- **Visual workflow builder** — DAGs, branching, loops, error handling
- **Language-agnostic** — tasks can be Python, Node, Shell, SQL, HTTP, anything
- **Event-driven triggers** — webhook, schedule, file watcher, message queue, flow trigger
- **Real-time execution** — not batch-only like Airflow
- **Built-in secrets management** — KV store, namespace isolation
- **Plugins for everything** — dbt, Spark, Kafka, S3, GCS, Postgres, NATS(!), MQTT, etc.
- **Self-hosted** — single binary or Docker, Java-based (JVM)
- **Git-based version control** — flows are YAML, live in git

## Architecture

```
Kestra Server (JVM)
├── Scheduler — cron triggers
├── Executor — runs tasks (local, Docker, K8s, remote)
├── Queue — internal task queue (JDBC, Kafka, or NATS as backend!)
├── Storage — flow definitions (YAML) + execution logs
└── UI — visual flow designer + monitoring
```

## When You'd Need It

- Complex multi-step workflows with conditional branching
- Orchestrating cross-system pipelines (ETL, data sync, backups)
- Workflows that need retries, timeouts, SLAs, notifications
- Team collaboration on workflows with approval gates

## Why It Might Be Redundant

| OneMind OS Has | Replaces Kestra's |
|---|---|
| NATS subscribers + rules KV | Event-driven workflows |
| Hermes cron jobs | Scheduled tasks |
| Hermes agent orchestration | Multi-step execution |
| NATS JetStream | Reliable message delivery |

## Decision

**Keeping in stack for reference.** If workflows grow beyond what NATS subscribers + Hermes cron can handle (e.g., complex branching, approval gates, visual debugging), Kestra is the upgrade path. Its NATS plugin means it can integrate natively.

## Resources

- Site: https://kestra.io
- Docs: https://kestra.io/docs
- GitHub: https://github.com/kestra-io/kestra
- License: Apache 2.0
- Deploy: `docker run -p 8080:8080 kestra/kestra:latest server local`

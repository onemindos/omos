---
title: "Synadia Platform Research"
type: reference
domain: SO
status: active
created: 2026-05-26
updated: 2026-05-26
tags: [nats, synadia, infrastructure, messaging, jetstream, kv-store]
---

# Synadia Platform Research

Synadia Communications — the creators and maintainers of NATS.io (41K+ GitHub stars, 11K+ Slack members).

Source: https://www.synadia.com/platform

---

## Products

### Synadia Platform
Single-tenant, fully managed NATS in your cloud. Enterprise control without operating it yourself.
- Two deployment modes: Managed (BYOC) and Self-hosted
- 99.99% uptime SLA (managed)
- Zero downtime upgrades
- NATS experts watching the dashboard
- 60-day trial available

### Synadia Cloud
Multi-tenant, hosted NATS. Build quickly without operating it yourself.

### Synadia Deploy for Kubernetes
Self-serve NATS & Synadia Platform deployment for Kubernetes environments.

### Synadia Insights
Granular, NATS-native observability including 100+ audit checks from Synadia experts.

### Synadia Protect
A security gateway for NATS. Enforce policy on every connection and message.

---

## Core Capabilities (Synadia Platform)

### Secure
- Authenticate users with any OIDC identity provider
- Encrypt and manage secrets with any KMS
- Define precise access policies for teams
- Full audit log for resource changes
- FIPS-compliant NATS builds available

### Manage
- Create JetStream streams, key-value stores, and object stores
- Mirror data across regions, clouds, and edge locations
- Real-time stream and consumer metrics monitoring
- Share streams securely across accounts and teams

### Observe
- Visualize entire topology (cloud to edge)
- Analyze metrics across servers, clients, and streams
- Monitor account usage across teams and domains
- Custom alert rules for any resource

### Extend
- **Connectors**: Move data between NATS and external systems. Out-of-the-box MongoDB, AWS, Azure, GCP integrations. Runtime-agnostic.
- **Workloads**: Run distributed compute and deploy containers or JS functions anywhere NATS runs

---

## Support Tiers
- Three response time SLA options (up to 24/7 global)
- Onboarding and consulting assistance
- Dedicated Slack channel for ad-hoc questions
- Recurring Health Checks and Architectural Reviews

---

## Solutions / Use Cases
- Distributed & Agentic AI
- Connected Car
- Financial Services
- Microservices
- Smart Manufacturing

---

## NATS Key/Value (KV) Store

Built directly on top of NATS JetStream. Client-side abstraction that models message subjects as keys — familiar KV operations backed by JetStream's stream persistence, replication, and durability.

Reference: https://docs.nats.io/nats-concepts/jetstream/key-value-store

### Key Characteristics

**Built-in Change Notifications**
Supports watchers for real-time updates when keys are added, updated, or deleted. Ideal for:
- Dynamic configuration management
- Cache invalidation

**Optimistic Concurrency Control**
Tracks a revision number per entry. The `Update` method accepts an expected revision — prevents non-deterministic "last writer wins" scenarios during concurrent modifications.

**History and TTL**
- Configurable history versions per key (default: 1)
- TTL support for automatic entry expiration
- Use case: session storage, ephemeral data

**Unified Infrastructure**
KV storage, object storage, and message brokering all within the same NATS server binary. No separate infrastructure needed — reduces deployment and monitoring complexity.

### OneMind Relevance
- **Agent state**: KV buckets store agent-state, zeus-context, config (see `NATS-INFRASTRUCTURE.md` — 6 active KV buckets)
- **Config distribution**: Push config changes to agents in real-time via watchers
- **Session data**: TTL-based entries for ephemeral agent context
- **Concurrency-safe writes**: Revision-based updates for concurrent agent modifications

---

## Related Notes
- `NATS-INFRASTRUCTURE.md` — full OneMind NATS infra doc (792 lines)
- `services/nats/NATS.md` — service-level NATS config
- `services/nats/hub.conf` — NATS server config

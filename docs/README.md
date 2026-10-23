# `docs/`

The canonical architecture and doctrine documents for OneMind OS. **Read these before touching anything else.** They define what we are building, how the bus is structured, and in what order work happens.

---

## Design documents (under review)

`architecture/` holds the OneMind OS design documents: `OneMind-OS-Architecture.md` (unified architecture), `ONEMIND-MASTER-PICTURE.md` (system synthesis) and `OMOS-DATA-PICTURE.md` (data pipeline), plus `decisions/`. They are being reviewed and updated; until then, `STACK.md` and `../CLAUDE.md` win on live facts. `architecture/zero-trust-plan.md` is the zero-trust plan (current state, gaps, ordered PRs).

---

## Canon (read in order)

| File | Purpose |
|------|---------|
| `ARCHITECTURE.md` | The 5-layer system design — substrate, memory, compute, task engine, interface. Actor model. Design principles. |
| `NATS-DOCTRINE.md` | Binding NATS rules — subject grammar, account topology, JetStream stream definitions, security model. |
| `MASTER-PLAN.md` | Build program — phased roadmap, current state, known issues, delegation packets for worker agents. |
| `STACK.md` | Tiered inventory of the live cluster: every Argo app and Helm release with state, backup, exposure and GitOps status; GitOps gaps. |
| `REPO-LAYOUT.md` | What each top-level directory is, active versus legacy, proposed consolidation. |
| `PLATFORM-EXTRACTION.md` | Platform extraction notes — how OneMind OS maps to the Synadia commercial suite (Control Plane, Insights, Protect, Connect, Workloads, Agents). |

---

## `infra/`

Runbooks and setup guides for the live infrastructure:

| File | Purpose |
|------|---------|
| `CLUSTER-SERVICE-INVENTORY.md` | Dated live inventory of Argo CD ownership, Kubernetes Services, zero-replica Deployments, PVCs, and investigation candidates |
| `CLICKHOUSE-TIMEMACHINE.md` | ClickHouse Time Machine setup, schema, and query patterns |
| `NATS-SETUP.md` | NATS server install, account config, JetStream setup on macOS |
| `TAK-INTEGRATION-REFERENCE.md` | TAK Server + NATS integration reference (ADS-B, AIS, CoT bridging) |
| `CONSTELLATION-OVERWATCH.md` | Constellation Overwatch COP — setup, bridge, Pulsar, fleet management |

---

## `patterns/`

Reusable NATS code patterns:

| File | Purpose |
|------|---------|
| `nats-request-reply-go.md` | Request/reply pattern in Go — reference for service implementations |

---

## `references/`

Research notes on external projects and the NATS ecosystem (`nats-ecosystem-references.md`, `nats/`, `ai/`, `defense/` and about 30 single-project notes). The top-level `references/` folder was merged in here on 2026-10-17. Not code, not deployed; review checklist in `research/2026-10-17-references-review.md`.

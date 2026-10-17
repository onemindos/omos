# `docs/`

The canonical architecture and doctrine documents for OneMind OS. **Read these before touching anything else.** They define what we are building, how the bus is structured, and in what order work happens.

---

## Design documents (under review)

`architecture/` holds the OneMind OS design documents: `OneMind-OS-Architecture.md` (unified architecture), `ONEMIND-MASTER-PICTURE.md` (system synthesis) and `OMOS-DATA-PICTURE.md` (data pipeline), plus `decisions/`. They are being reviewed and updated; until then, `STACK.md` and `../CLAUDE.md` win on live facts.

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

Research and ecosystem notes:

| File | Purpose |
|------|---------|
| `nats-ecosystem-references.md` | NATS ecosystem repo index and notes |
| `agri_tactical_nats_repos.md` | Agricultural + tactical NATS use cases |
| `nats_simulation_quick_reference.md` | NATS simulation quick reference |
| `nats_simulation_report.md` | Full NATS simulation research report |
| `nats_simulation_urls_for_forking.txt` | Source URLs from simulation research |
| `onemind-recovery/` | Recovery procedures and backup notes |

---
title: "NEX — NATS Execution Engine"
section: "13-SYNADIA"
type: note
tags: [synadia, nex, execution-engine, workloads, edge, runtime, nexlet, distributed]
source: https://github.com/synadia-io/nex
created: 2026-06-14
---

# NEX — NATS Execution Engine

**NEX** (`synadia-io/nex`, ★310, Go) is a lightweight way to **run distributed workloads over NATS** — deploy and execute binaries, functions, or services across a fleet **without standing up Kubernetes or heavy orchestration**. The bus *is* the control plane. This is one of the most strategically interesting projects for an edge fleet.

See also: [[NATS/13-SYNADIA/synadia-products|Synadia Products]] · [[NATS/11-ADVANCED/embedded-nats|Embedded NATS]] · [[NATS/06-TOPOLOGY-SCALING/leaf-nodes|Leaf Nodes]]

---

## The idea

You already have NATS connecting every node in your fleet. NEX uses that fabric to also **schedule and run work**: a node advertises capacity, you submit a workload, and NEX places it on a node and runs it — with the workload getting **scoped NATS credentials** automatically (so it can immediately do pub/sub, request/reply, JetStream).

- **Nodes** run a NEX process; **nexlets** (runtime agents) execute the actual workloads.
- Nodes **auction** work to nexlets — imperative, lightweight placement, no etcd/scheduler stack.
- **Native executables work out of the box**; bring your own runtime by writing a pluggable **nexlet** for other workload types.
- **Location transparency** — a workload runs wherever connectivity exists: edge node, retail box, connected car, cloud. Workloads can run **locally on edge nodes even when disconnected** from the cloud.

It was unveiled at KubeCon NA 2024 as Synadia's answer to "distributed/edge workloads without Kubernetes complexity."

---

## Quickstart (shape of it)

```bash
# install
go install github.com/synadia-io/nex/cmd/nex@latest

# 1. have a NATS server
nats-server -a 127.0.0.1 -p 4222

# 2. bring up a NEX node (native nexlet by default)
nex -s nats://127.0.0.1:4222 node up
nex -s nats://127.0.0.1:4222 node list      # verify registration

# 3. deploy a workload via a Nexfile (YAML)
```

```yaml
# Nexfile
name: "system-reporter"
description: "Reports system stats every 30 seconds"
type: native
lifecycle: service        # service (long-running) or function (reactive)
start_request:
  # ...command, args, env, scoped creds...
```

The workload comes up on a node, gets NATS creds, and starts publishing — all coordinated over the bus.

---

## Why NEX fits edge fleets (vs Kubernetes)

| | Kubernetes | NEX |
|--|-----------|-----|
| Control plane | etcd + API server + scheduler | **NATS** (already there) |
| Footprint | heavy | **single binary**, lightweight |
| Offline edge | struggles | **runs locally when disconnected** |
| Workload comms | sidecars/service mesh | **native NATS creds, instant** |
| Model | declarative, complex | **imperative, simple** |

NEX trades Kubernetes' vast feature set for **operational simplicity on constrained, intermittently-connected nodes** — exactly the edge profile Kubernetes handles badly.

---

## NEX CE & editions
- **`synadia-io/nex`** — the open-source engine (Apache-2.0, community).
- **`synadia-io/nex-ce`** — "Nex CE" public assets (community edition packaging).
- NEX is also supported within **Synadia Platform** for enterprise users.

So NEX itself is **open source** — usable in a sovereign, self-hosted fleet with no subscription.

---

## OneMind Relevance (strategic)
- The fleet already runs NATS on **hundreds of planned nodes** (House/Barn/Office, vehicles, robots, drone). NEX is the **natural way to dispatch workloads** onto them **without adopting Kubernetes** — which directly serves the locality/portability/sovereignty values.
- **Offline-capable edge execution** matches the leaf-node + local-JetStream doctrine: a barn node keeps running its assigned workload during a hub outage, coordinated over the local bus.
- **Scoped creds per workload** dovetails with our NKEY/JWT identity model ([[NATS/10-PROTOCOL-INTERNALS/nuid-and-nkeys|NUID & NKEYS]]) — every dispatched job is an authenticated fabric citizen.
- **Status:** on the evaluation radar. It's OSS, so piloting it doesn't compromise sovereignty. Candidate for how agents/jobs land on physical nodes.

## Sources
- NEX — https://github.com/synadia-io/nex
- NEX CE assets — https://github.com/synadia-io/nex-ce
- Synadia on NEX/edge — https://www.efficientlyconnected.com/synadia-reimagines-distributed-workloads-with-nats-and-nex-a-vision-for-edge-and-cloud-native-computing/
- RethinkConn 2025 — https://www.synadia.com/lp/rethink-conn-2025

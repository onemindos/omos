---
title: "Beszel — Fleet Monitor"
type: research-candidate
status: evaluating
created: 2026-06-08
updated: 2026-06-08
owner: zeus
category: infrastructure-monitoring
tags: [research, monitoring, fleet, infrastructure, self-hosted, candidate]
ai_summary: "Lightweight server fleet monitoring. Hub + agent architecture — tiny agent per node, one central hub shows CPU/RAM/disk/network across every machine. Evaluate when fleet grows beyond current nodes."
---

# Beszel — Fleet Monitor

> **Status:** Candidate — evaluate when fleet expands. Lightweight alternative to Grafana/Prometheus stack.
> When adopted, move to `SO/03 Infrastructure/services/`.

## What It Is

Lightweight self-hosted server monitoring with a hub-and-agent architecture. A tiny agent binary runs on each node and reports metrics back to a central hub. Single dashboard shows the health of every machine in the fleet.

## Architecture

```
Beszel Hub (omos-mbp or dedicated node)
    ↑ SSH tunnels
Agent  Agent  Agent  Agent
(MBP)  (Pi)   (Pi)   (VPS)
```

Agents communicate back to the hub over SSH — no open ports required on agent nodes.

## Key Features

| Feature | Description |
|---|---|
| **Lightweight** | Tiny resource footprint vs Grafana/Prometheus |
| **Hub + agent** | One hub, unlimited nodes |
| **SSH transport** | No inbound ports needed on monitored nodes |
| **Metrics** | CPU, RAM, disk, network per node |
| **Alerts** | Threshold-based alerts |
| **Docker stats** | Container-level metrics per node |
| **Web UI** | Clean dashboard, no Grafana config required |

## Why It's Interesting for OneMind OS

- Fleet growing toward hundreds of nodes (house/barn/office + robots + VPS)
- Currently no unified fleet health view
- Much lighter than full Prometheus + Grafana + exporters stack
- SSH-based transport aligns with existing Tailscale mesh — no new firewall rules
- Docker stats coverage useful across all containerized services

## Current Monitoring Stack

- Langfuse `:3000` — AI observability
- worldwideview `:3001` — world monitor
- No unified fleet health dashboard yet

## Evaluate When

- [ ] Fleet reaches 5+ active nodes
- [ ] Current ad-hoc monitoring becomes insufficient
- [ ] Compare against Netdata (already documented in SO/03 — `NETDATA-MONITORING.md`)

## Links

- **GitHub:** https://github.com/henrygd/beszel
- **Docs:** https://beszel.dev

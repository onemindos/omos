---
title: Network Mesh Comparison — Tailscale/Headscale vs NetBird vs Twingate vs WireGuard
date: 2026-06-05
tags: [networking, mesh, wireguard, iot, nats, infrastructure]
status: decision-locked
decision: Headscale + Tailscale clients
---

# Network Mesh Comparison
## Tailscale/Headscale · NetBird · Twingate · WireGuard

**Date:** June 2026  
**Scope:** Private fleet mesh — 50+ nodes, cloud/local/edge/IoT, NATS overlay  
**Decision:** **Headscale + Tailscale clients** — confirmed, rationale below

---

## Executive Summary

**Headscale + Tailscale clients is the right long-term answer for OneMind OS.** WireGuard is the foundation every one of these tools builds on — you're not choosing *instead* of WireGuard, you're choosing *how* to manage it. Twingate is a hard no (cloud-locked, per-user pricing at fleet scale). NetBird is architecturally solid but you've already lived through its fragility. Raw WireGuard is too manual at 50+ nodes. Headscale + Tailscale clients gets you permanent sovereignty, zero cost at any scale, best-in-class client apps, and NATS rides on top of it cleanly.

---

## What Each One Actually Is

| | WireGuard | Tailscale / Headscale | NetBird | Twingate |
|---|---|---|---|---|
| **Layer** | Protocol only | Control plane + WireGuard | Control plane + WireGuard | ZTNA proxy + WireGuard |
| **Model** | Manual key exchange | Coordination server (cloud or self-hosted) | Management server (cloud or self-hosted) | Cloud control plane only |
| **Self-hosted** | N/A (it's a kernel module) | Headscale (1 binary) | Yes (4 containers) | ❌ No — cloud-only |
| **Free at 50+ nodes** | Yes (no cost, all manual) | **Yes (Headscale)** | Yes (self-hosted) | ❌ No — per-user pricing |
| **IoT / embedded support** | Manual/kernel | Tailscale app (Linux ARM down to Pi Zero) | NetBird agent (Linux ARM) | ❌ No native IoT |
| **DNS / MagicDNS** | None | ✅ Best-in-class | ✅ Partial | Split-tunnel DNS only |
| **Control plane sovereignty** | You own everything | **You own everything (Headscale)** | You own everything | ❌ Their cloud |
| **Mobile clients** | Manual (no official app) | ✅ Best iOS/Android apps | ✅ Good apps | ✅ Good apps |
| **Complexity** | Very high (manual) | Low (1 binary + existing apps) | Medium (4 containers) | Low (SaaS) |

---

## Deep Dive: Each Option

### 1. Raw WireGuard
**What it is:** A kernel-level VPN protocol. That's it. No coordination, no DNS, no key management, no dashboard.

**Strengths:**
- Absolute minimum overhead (~300 lines of kernel code)
- You own 100% of every key and config
- Runs on anything — routers, MCUs with WireGuard ports, embedded Linux
- Protocol is stable, peer-reviewed, mature

**Weaknesses:**
- Every node addition requires manual key exchange on EVERY peer
- No dynamic routing, no DNS, no auto-reconnect
- No mesh topology — you manually define which peers can reach which
- 50 nodes = ~2,500 key pairs to manage manually
- No key rotation automation, no revocation at scale

**NATS fit:** WireGuard provides the tunnel; NATS rides on top natively. No integration needed — NATS subjects traverse the WireGuard interface like any TCP/UDP traffic.

**OneMind fit:** ❌ Too much manual overhead at 50+ nodes. It's the right answer *inside* a single site, not across a growing fleet. WireGuard is the protocol you're using whether you pick Headscale, NetBird, or Twingate — the others are just automating the key exchange and routing.

---

### 2. Tailscale (managed) / Headscale (self-hosted)

**What it is:** Tailscale built best-in-class WireGuard mesh clients and a coordination server. Headscale re-implements only the coordination server — you keep the Tailscale clients.

```
Tailscale managed:  Tailscale client → login.tailscale.com (their cloud)
Headscale:          Tailscale client → your server (hs.onemindos.dev:8085)
```

**Strengths:**
- **Best client apps in the ecosystem** — iOS, Android, macOS, Linux all first-class
- **MagicDNS** — every node gets `omos-<name>.mesh.onemindos.dev` automatically
- **Zero cost at any scale** with Headscale — 50, 500, 5000 nodes, same $0
- **Permanent domain** — `onemindos.dev` is yours forever, no vendor subdomain
- **1 binary** for the control plane — no Postgres, no Redis, SQLite only
- **DERP relay built-in** — nodes behind NAT/CGNAT still connect (drones in field, robots at remote sites)
- **Subnet routing** — MBP advertises `192.168.1.0/24`, all mesh nodes reach LAN directly
- **If vendor disappears:** nothing changes — Headscale is yours, Tailscale client apps are free
- **Posture checks:** not in Headscale, but ACL/tag policies available via hujson

**Weaknesses:**
- **Headscale is a community project** — not by Tailscale Inc., protocol compatibility could drift (unlikely, has been stable for 3+ years)
- **No Funnel** (public URL exposure without port forwarding) — Tailscale managed only
- **No posture checks** — Headscale doesn't enforce device health/trust
- **Headplane is the only UI** — community-built, not official

**NATS fit:** ✅ Perfect. NATS listens on `0.0.0.0:4222`. Any mesh node can connect to `omos-mbp.mesh.onemindos.dev:4222`. Mesh IPs are stable (Headscale assigns sequentially). You can pin NATS routes to mesh FQDNs and they never break. Long-term: NATS JetStream can distribute across sites over the mesh tunnel.

**IoT/Edge fit:** ✅ Tailscale runs on Pi Zero W (64MB RAM), ARM Cortex-A7, most embedded Linux. Pre-auth keys + tags let you auto-enroll sensors with no user interaction. DIMOS, OM1, and future drones can all run Tailscale agents. Off-mesh devices (ESP32, Meshtastic, Zigbee) bridge through a Pi or HA node already on mesh.

**Scale:** ~10K nodes per Headscale instance — orders of magnitude beyond your fleet target.

---

### 3. NetBird

**What it is:** WireGuard mesh with a management plane that includes dashboard, setup keys, posture checks, and access policies. Self-hosted via 4 Docker containers (management, signal, relay/STUN, dashboard).

**Strengths:**
- **Full dashboard UI** — better than Headplane for day-to-day fleet management
- **Setup keys** — auto-enroll new nodes (similar to Tailscale pre-auth keys)
- **Peer groups and fine-grained ACL** — more granular than Headscale's hujson
- **Posture checks** — verify kernel version, OS, firewall state before allowing access
- **Network routes** — equivalent to subnet routing
- **Active development** — commercial company, NetBird Cloud available as fallback

**Weaknesses:**
- **4 containers + management server is a fragile stack** — any one component failing breaks the mesh (you lived this in June 2026)
- **STUN cannot run behind Cloudflare Tunnel** — UDP/3478 is silently broken, forces relay-only connections (root cause of your June 2026 mesh failure)
- **Management server is Linux-only** — no macOS binary, Docker required on MBP
- **Signal + STUN + relay** are separate services — complexity vs Headscale's embedded DERP
- **If NetBird Cloud goes down and your self-hosted breaks** — no fallback
- **DNS wildcards broken with certain configs** — your diagnosis doc confirmed this

**NATS fit:** ✅ Same as Headscale — NATS rides on top, mesh IPs are stable. NetBird's nameserver groups can point to the NATS broker by hostname if configured. No deeper integration.

**IoT/Edge fit:** ✅ NetBird agent runs on embedded Linux. Similar ARM support to Tailscale. Setup keys work headlessly. But the control plane fragility is a problem — if the management server is down, new devices can't enroll and existing peers may lose coordination.

**Scale:** ~5K+ nodes (self-hosted). Sufficient for your fleet.

**Verdict for you:** You already ran NetBird, diagnosed its failure modes, and replaced it. The June 2026 diagnosis (STUN down, iPhone/HA disconnected) was not a configuration error — it was an architectural brittleness. 4 containers is 4 failure points. You moved to Headscale for the right reasons.

---

### 4. Twingate

**What it is:** Zero-Trust Network Access (ZTNA) — a different model from mesh VPN. Instead of all nodes joining a shared network, Twingate brokers individual resource access via its cloud control plane. No network join — users/services authenticate per-resource.

**Pricing (2026):**
- Starter: Free, up to 5 users
- Home: $15/month flat (up to 7 users)
- Teams: $5/user/month
- Business: $10/user/month
- Enterprise: Custom

**Strengths:**
- **Best-in-class zero-trust UX** for enterprise B2B access control
- **No "network join" model** — users get per-resource access, not mesh membership
- **Posture checks, MFA, device trust** — excellent compliance story
- **No client port forwarding needed** — connector sits behind firewall, outbound-only
- **Service accounts** — non-human access for CI/CD, service mesh use cases

**Weaknesses for OneMind OS:**
- ❌ **Cloud-only control plane** — there is no self-hosted option. Period. Their control plane is their cloud.
- ❌ **Per-user pricing at fleet scale** — 50 nodes = 50 users = $250/month minimum at Teams tier. 500 nodes = $2,500/month.
- ❌ **No IoT support** — designed for human users accessing resources, not autonomous devices forming a mesh
- ❌ **No NATS integration path** — NATS operates at the network level; Twingate's model is resource proxying, not IP mesh. NATS subjects don't map to "resources" cleanly.
- ❌ **Vendor lock-in** — if Twingate shuts down or changes pricing, your entire network access model breaks
- ❌ **Not a mesh** — nodes can't reach each other directly; all traffic proxied through connectors

**NATS fit:** ❌ Fundamentally incompatible model. NATS requires nodes to connect directly at the IP level. Twingate proxies access to specific resources — you could expose `:4222` as a Twingate resource, but then every node needs a Twingate client and per-user license, and P2P NATS cluster gossip doesn't work through proxy connectors.

**IoT/Edge fit:** ❌ Twingate has no concept of device-to-device mesh. It's human users → protected resources. Drones, robots, sensors do not fit this model.

**Verdict:** Wrong category for OneMind OS. Twingate solves "how do employees securely access company resources from anywhere." You're solving "how do autonomous nodes in a fleet communicate with each other." Completely different problem spaces.

---

## NATS Integration Comparison

This is a critical axis since NATS is your long-term data plane.

| | WireGuard | Headscale/Tailscale | NetBird | Twingate |
|---|---|---|---|---|
| **NATS over mesh** | ✅ Yes (manual IPs) | ✅ Yes (stable MagicDNS FQDNs) | ✅ Yes (stable mesh IPs) | ❌ Incompatible model |
| **Stable hostname for NATS server** | ❌ Manual config | ✅ `omos-mbp.mesh.onemindos.dev:4222` | ✅ `<peer-ip>:4222` | ❌ |
| **IoT devices reach NATS** | Manual | ✅ Pre-auth key → mesh join → NATS | ✅ Setup key → agent → NATS | ❌ |
| **Multi-site NATS cluster** | Manual tunnel | ✅ JetStream across sites via mesh | ✅ Same | ❌ |
| **NATS leaf nodes at edge** | Manual | ✅ Any mesh node runs leaf node | ✅ Same | ❌ |
| **Off-mesh devices → NATS** | N/A | ✅ Bridge node (Pi) → NATS | ✅ Same | ❌ |

**Long-term NATS architecture on Headscale:**
```
Cloud VPS (leaf node)
    ↓ WireGuard via mesh
omos-mbp (NATS hub, :4222)
    ↓ mesh
omos-ha, omos-dell, omos-sensors (NATS clients)
    ↓ bridge
ESP32/Meshtastic/Zigbee → Pi bridge → NATS
```

Every mesh node can resolve `omos-mbp.mesh.onemindos.dev:4222` — permanent, never breaks, no IP changes.

---

## IoT + Edge + Cloud Matrix

| Scenario | WireGuard | Headscale | NetBird | Twingate |
|---|---|---|---|---|
| Pi / ARM Linux node | ✅ Manual | ✅ Auto-enroll | ✅ Auto-enroll | ❌ |
| Docker container on node | ✅ Manual | ✅ Tailscale sidecar | ✅ NetBird sidecar | ⚠️ Connector |
| Drone companion PC (Linux) | ✅ Manual | ✅ Pre-auth key | ✅ Setup key | ❌ |
| ESP32 / bare metal MCU | ❌ (no WG stack) | ❌ (no client) | ❌ (no client) | ❌ |
| Home Assistant OS | ✅ Add-on | ✅ Tailscale add-on | ✅ NetBird agent | ❌ |
| Cloud VPS | ✅ Manual | ✅ Linux client | ✅ Linux agent | ⚠️ Connector |
| macOS laptop | ✅ Manual | ✅ App | ✅ App | ✅ App |
| iOS device | ❌ (no kernel access) | ✅ Official app | ✅ App | ✅ App |
| DIMOS / OM1 robot | ✅ Manual (Linux) | ✅ ARM agent | ✅ ARM agent | ❌ |
| Barn sensors (50+ nodes) | ❌ Too manual | ✅ Bulk pre-auth keys | ✅ Setup keys | ❌ |

**Gap for all options:** True embedded (ESP32, Pixhawk, Meshtastic LoRa) cannot run WireGuard. Bridge pattern: `device → Pi/companion → mesh` is universal. This is solved the same way regardless of Headscale vs NetBird.

---

## Long-Term Architecture Fit (50+ nodes → 500+ nodes)

### Headscale + Tailscale clients

```
omos-mbp (Headscale control plane, LaunchDaemon)
├── Registration: tailscale up --login-server=http://hs.onemindos.dev:8085 --authkey=<key>
├── MagicDNS: omos-<name>.mesh.onemindos.dev (automatic)
├── Subnet route: 192.168.1.0/24 advertised by omos-mbp
├── NATS hub: omos-mbp:4222 (reachable by all mesh nodes by FQDN)
└── Scale: 10K nodes, $0 forever
```

Adding a new node (robot, sensor, VPS):
1. `tailscale up --login-server=... --authkey=<reusable-key> --hostname=omos-<name>`
2. NATS client connects to `omos-mbp.mesh.onemindos.dev:4222`
3. Done. No dashboard changes, no key distribution, no manual config.

### What changes at scale
- At 50+ nodes: add ACL policies in `/etc/headscale/acl.hujson` to segment robots from laptops
- At 100+ nodes: consider a second Headscale instance per site cluster for resilience
- At 500+ nodes: evaluate NATS JetStream clustering across sites (Headscale mesh handles the transit)

---

## Feature Gap: What Headscale Is Missing vs Alternatives

| Feature | Missing from Headscale | Workaround |
|---|---|---|
| Posture checks | ✅ Missing | Tag-based segmentation in ACL; eBPF host monitoring via Netdata |
| Device trust scoring | ✅ Missing | Out of scope for private fleet — all nodes are yours |
| Funnel (public URLs) | ✅ Missing | Not needed — Cloudflare Tunnel or port forward for specific services |
| Official support SLA | ✅ Missing | Community project — acceptable for self-run infrastructure |
| Node auto-approval | ✅ Missing | Pre-auth keys with `--reusable` covers headless enrollment |
| Multi-org / tenants | ✅ Missing | Single user namespace — correct for OneMind |

None of these gaps affect your use case.

---

## Cost at Scale

| Nodes | WireGuard | Headscale | NetBird self | NetBird cloud | Twingate Teams |
|---|---|---|---|---|---|
| 10 | $0 | $0 | $0 | Free tier | $50/mo |
| 50 | $0 | $0 | $0 | ~$50/mo | $250/mo |
| 100 | $0 | $0 | $0 | ~$100/mo | $500/mo |
| 500 | $0 | $0 | $0 | ~$500/mo | $2,500/mo |
| 1000 | $0 | $0 | $0 | ~$1,000/mo | $5,000/mo |

**At your fleet target (100-500 nodes), only WireGuard, Headscale, and NetBird self-hosted are viable cost-wise.** NetBird cloud and Twingate both become prohibitive.

---

## Decision Summary

### Primary: **Headscale + Tailscale clients** ✅ (Already deployed)

**Why it wins long-term:**
1. **$0 forever** — fleet can grow to 10K nodes, cost doesn't change
2. **Permanent URLs** — `onemindos.dev` is yours. No vendor can revoke it or change it
3. **1 binary** — Headscale is one process, SQLite database, no Postgres, no Redis, no signal server, no STUN server. DERP relay is embedded or you point to Tailscale's free relay infrastructure
4. **Best client ecosystem** — Tailscale iOS/Android/macOS apps are first-class. They work with Headscale unchanged
5. **NATS fits perfectly** — stable MagicDNS FQDNs mean NATS routes never break even as nodes come and go
6. **IoT enrollment at scale** — reusable pre-auth keys + tag auto-assignment → robots and sensors self-enroll

### Not Recommended

| Option | Verdict | Reason |
|---|---|---|
| **Twingate** | ❌ Hard no | Cloud-locked, per-user pricing, wrong model for device mesh, NATS incompatible |
| **NetBird self-hosted** | ❌ Replaced | Fragile 4-container stack, STUN/Cloudflare incompatibility, you already hit the failure modes |
| **NetBird cloud** | ⚠️ Emergency fallback | Fine as temp fallback if Headscale has issues, but adds vendor dependency |
| **Raw WireGuard** | ⚠️ For single-site LAN only | Too manual at fleet scale; use as embedded protocol not management tool |

---

## One Thing to Watch

**OpenZiti** is the only architecture that genuinely scales to millions of endpoints and handles cloud/edge/IoT natively with zero-trust built into the protocol (dark services — no listening ports). If OneMind OS grows to thousands of nodes across multiple sites with strict zero-trust enforcement requirements, OpenZiti is worth a fresh evaluation. It pairs with NATS well (OpenZiti as network fabric, NATS as messaging layer).

For your current and near-term fleet (10-500 nodes), Headscale handles it cleanly. Revisit OpenZiti when fleet exceeds 500+ heterogeneous nodes across 5+ physical sites.

---

## Quick Reference

```bash
# Headscale health
curl -s http://localhost:8085/health

# List mesh nodes
sudo headscale nodes list

# Enroll new node (on target device)
tailscale up \
  --login-server=http://hs.onemindos.dev:8085 \
  --authkey=hskey-auth-cf37dbaff1366bd... \
  --hostname=omos-<name>

# NATS address from any mesh node
omos-mbp.mesh.onemindos.dev:4222
```

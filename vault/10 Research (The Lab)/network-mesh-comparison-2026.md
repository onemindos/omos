---
id: network-mesh-comparison-2026
title: Network Mesh Platform Comparison — 2026
tags: [infrastructure, networking, vpn, wireguard, tailscale, netbird, twingate, nats, iot, mesh]
owner: legacy
status: decision-locked
updated: 2026-06-07
---

# Network Mesh Platform Comparison — 2026

**Scope:** OneMind OS fleet — omos-mbp, omos-cloud, IoT nodes (DIMOS, OM1, greenhouse, barn), future drones/robots  
**Requirements:** NATS compatibility, cloud + local + edge + IoT, sovereignty, long-term scalability  
**Decision date:** 2026-06-07

---

## TL;DR — Decision

| Layer | Platform | Reason |
|-------|----------|--------|
| **Transport (L3/4)** | WireGuard (kernel) | Foundation — fastest, most battle-tested; used by all tools below |
| **Control Plane — Now** | Tailscale cloud | Already deployed, working, zero friction. Keep it. |
| **Control Plane — Future** | NetBird self-hosted | Full sovereignty, unlimited nodes, BSD-3, NATS-compatible topology |
| **Messaging (L7)** | NATS JetStream | Runs OVER the mesh — pub/sub for IoT/agent comms |
| **Skip entirely** | Twingate | User-to-resource model, SaaS-only, not a mesh |
| **Skip as control plane** | Raw WireGuard | Missing management layer — too much DIY |

**Migration trigger:** Fleet > 15 nodes OR any hard vendor dependency hits with Tailscale cloud

---

## The Stack Mental Model

```
┌─────────────────────────────────────────────────────────┐
│  L7 — NATS JetStream  (pub/sub, agent comms, IoT data)  │
├─────────────────────────────────────────────────────────┤
│  L3/4 — NetBird / Tailscale  (encrypted overlay mesh)   │
├─────────────────────────────────────────────────────────┤
│  L0 — WireGuard kernel  (raw transport — all use this)  │
└─────────────────────────────────────────────────────────┘
```

NATS and the mesh are NOT competitors. NATS is the messaging layer. The mesh is the encrypted network. They stack.

---

## Platform Comparison Matrix

| Dimension | Tailscale (cloud) | Headscale (self-hosted) | NetBird (self-hosted) | Twingate | WireGuard (raw) |
|-----------|:-----------------:|:-----------------------:|:---------------------:|:--------:|:---------------:|
| **Transport** | WireGuard | WireGuard | WireGuard | TLS (NOT WireGuard) | WireGuard |
| **Control plane** | Tailscale SaaS | Self-hosted (3rd party) | Self-hosted (official) | SaaS-only ❌ | None — DIY |
| **Model** | Mesh | Mesh | Mesh + ZTNA | User→Resource ❌ | P2P primitive |
| **IoT support** | ✅ First-class | ✅ via Headscale | ✅ First-class | ⚠️ Headless Linux only | ✅ Embedded Linux |
| **ARM/MIPS binaries** | ✅ (Pi, ARM64, static) | ✅ via TS client | ✅ (ARMv6, ARM64, MIPS) | ⚠️ Linux only | ✅ Kernel module |
| **Node limits** | 50 tagged free + $1/mo | Unlimited | Unlimited (self-hosted) | 50–300 resources | None |
| **Pricing** | Per-user + $1/IoT device | Free OSS | Free OSS (CE) | Per-user only | Free |
| **Self-hosted control plane** | ❌ (need Headscale) | ✅ | ✅ Official | ❌ Never | N/A |
| **ACLs / Zero Trust** | ✅ All plans | ✅ subset | ✅ Comprehensive | ✅ User-centric | ❌ IP allowlists only |
| **Auth keys (headless)** | ✅ | ✅ | ✅ | ⚠️ Service Accounts | ❌ |
| **NATS compatibility** | Transparent (L3 tunnel) | Transparent | Transparent | Transparent | Transparent |
| **Device-to-device mesh** | ✅ Native | ✅ Native | ✅ Native | ❌ NOT supported | ✅ Manual |
| **Subnet router (non-WG devices)** | ✅ (don't count vs limits) | ✅ | ✅ | ❌ Connector only | ❌ |
| **MagicDNS / auto-DNS** | ✅ | ✅ | ✅ | ❌ No DNS management | ❌ |
| **SCIM / IdP** | ✅ Standard+ | ✅ partial | ✅ Team+ | ✅ Business+ | ❌ |
| **Device posture checks** | ✅ Standard+ | ❌ | ✅ Business+ | ✅ Business+ | ❌ |
| **GitOps ACL management** | ✅ | ✅ | ⚠️ API-based | ❌ | ❌ |
| **Vendor lock-in** | ⚠️ Medium | Low | Low | ❌ High | None |
| **Sovereignty** | ⚠️ Coordination via SaaS | ✅ | ✅ | ❌ | ✅ |
| **Status** | ✅ Production | ✅ Stable | ✅ Production (v0.72.2) | ✅ Production | ✅ Stable (in kernel) |

---

## Deep Dive: Each Platform

### 1. WireGuard (raw)

**What it is:** A modern VPN kernel module. The transport layer that ALL other platforms in this comparison use internally.

**Strengths:**
- Fastest VPN available: ~1 Gbps throughput, 0.4ms latency (4× faster than OpenVPN)
- In Linux kernel since 5.6 (March 2020) — zero external install on modern distros
- Ultra-minimal codebase (~4,000 lines) — easy to audit, easy to port
- Runs on: Raspberry Pi, OpenWrt, Buildroot, EdgeOS, Termux, Android, iOS
- Free, GPLv2

**What's missing (intentionally):**
- No key management or distribution
- No ACLs beyond IP-based allowlists
- No user identity / SSO / MFA
- No dynamic peer discovery (static config only)
- No control plane or management UI
- No automatic mesh topology (full mesh = configure every peer on every peer)
- No automatic NAT traversal (has keepalive for hold-open; no STUN/TURN built in)
- No certificate revocation

**Verdict:** WireGuard is the foundation. Run it everywhere as a kernel module. Never manage it directly — let Tailscale or NetBird manage it.

---

### 2. Tailscale (cloud) + Headscale (self-hosted control plane)

**What it is:** The leading WireGuard mesh VPN product. Tailscale cloud = SaaS control plane. Headscale = open-source 3rd-party reimplementation of the control plane.

**Strengths:**
- Already deployed in Zeus's fleet (tail717ba0.ts.net) — zero migration needed
- Best NAT traversal in class — handles the most complex firewall scenarios
- Edge/IoT product tier: 4via6 routing for overlapping IP subnets (critical for robotics/vehicles)
- Subnet router: bridge non-WG devices (sensors, PLCs, cameras) without installing WG on them
- Excellent ARM support: Pi packages, static ARM/MIPS binaries, robotics testimonials
- Auth keys for headless device enrollment (factory-deploy)
- MagicDNS auto-DNS across all tailnet nodes
- 50 tagged resources free; add more at $1/device/mo
- GitOps ACL management

**Weaknesses:**
- Coordination still goes through Tailscale SaaS even with Headscale (DERP relays = Tailscale's)
- Headscale is NOT made by Tailscale — 3rd party project (juanfont/headscale)
- Headscale doesn't support all features (no Tailscale Funnel, no Aperture AI, MDM limited)
- IoT scale (large fleet) requires Enterprise = contact sales / custom pricing
- Per-user + per-device pricing adds up at scale
- Headscale dropped from Zeus's fleet June 2026 (complexity, STUN issues behind Cloudflare Tunnel)

**NATS:**
Tailscale is L3/L4. NATS runs over the Tailscale mesh transparently. Your NATS server at `:4222` is accessible to all tailnet peers via MagicDNS. No special config — it just works.

**Verdict:** Keep Tailscale cloud for now. It's running, it works, MagicDNS is in use across the fleet. Migration trigger = fleet >15 nodes or per-device pricing becomes painful.

---

### 3. NetBird (self-hosted — recommended migration target)

**What it is:** Open-source, self-hosted ZTNA mesh built on WireGuard. Officially maintained (not 3rd party). BSD-3 license. Full control plane included.

**Strengths:**
- Self-hosted control plane is OFFICIAL — Management + Signal + Relay + Dashboard
- Community Edition: FREE, unlimited nodes, BSD-3 license → true sovereignty
- Single Docker Compose deploy — runs on omos-cloud (Lightsail) in minutes
- Same WireGuard kernel transport — identical performance to Tailscale
- Excellent IoT: ARM64, ARMv6 (Pi Zero/1), MIPS (OpenWrt routers) official binaries
- DeltaQuad (drone manufacturer) is a production customer — validates IoT/robotics use case
- Comprehensive ACLs: group policies, posture checks, device approval
- SCIM + IdP integration on paid plans
- 25.8k GitHub stars, v0.72.2 released June 6 2026 — very active
- ISO 27001 + GDPR + DORA certified (SaaS plans)

**Weaknesses:**
- Self-hosted HA requires Commercial License (paid) — Community Edition is single-instance
- SCIM/IdP provisioning automation requires paid tier (Community = manual groups)
- Newer project than Tailscale — fewer edge cases battle-tested
- WebSocket relay (replacing Coturn TURN since v0.29.0) is still in transition
- ~1.3k open issues (normal for scale, but worth watching)

**Self-hosted stack:**
```
omos-cloud (Lightsail) — 4c/16GB already exists
├── netbird-management   (gRPC control plane)
├── netbird-signal       (WebRTC-style P2P signaling)
├── netbird-relay        (fallback TURN)
├── netbird-dashboard    (web UI)
└── Caddy / Traefik      (TLS reverse proxy)
```

**NATS:**
NetBird creates WireGuard interfaces on each peer. NATS runs over those interfaces. For NATS leafnodes across locations: peer1 (omos-mbp) can reach peer2 (omos-cloud) NATS server via NetBird mesh IP. Leafnode config = standard NATS `remotes:` pointing to the NetBird mesh IP of each hub.

**Verdict:** Best long-term choice. Self-hosted = sovereignty. Unlimited nodes = never outgrow it. WireGuard transport = same performance. ARM/MIPS = IoT ready. Deploy on omos-cloud when ready to migrate.

---

### 4. Twingate

**What it is:** Zero Trust Network Access (ZTNA) product. NOT a mesh VPN.

**Architecture:** Client → Relay → Connector → Resource (always through a proxy; never peer-to-peer mesh).

**Why it doesn't fit OneMind OS:**
1. **No mesh** — Twingate explicitly says it's not a mesh VPN (their own docs). Device-to-device mesh is not supported.
2. **SaaS-only control plane** — no self-hosted option. Ever. On any plan.
3. **User-to-resource model** — designed for humans accessing corporate resources, not device fleets
4. **Per-user pricing** — doesn't scale for IoT device fleets
5. **NATS incompatible model** — pub/sub requires peers to connect to each other; Twingate's connector model would require routing every NATS message through connectors
6. **Not WireGuard** — uses TLS tunnels through connectors (different performance characteristics)

**Verdict: Skip entirely.** Wrong architecture for fleet mesh + NATS + IoT.

---

## NATS Integration Deep Dive

All four platforms are L3/L4 overlays — they all carry any TCP/UDP protocol transparently, including NATS.

**How NATS maps to the mesh:**

```
NATS Cluster / Supercluster topology:
─────────────────────────────────────
omos-mbp (NATS hub, :4222)
  └── WireGuard mesh IP: 100.x.x.x (Tailscale) or 10.x.x.x (NetBird)

omos-cloud (NATS hub, :4222)
  └── WireGuard mesh IP: 100.y.y.y (Tailscale) or 10.y.y.y (NetBird)

IoT node / DIMOS (NATS client or leafnode)
  └── WireGuard mesh IP: auto-assigned

NATS Leafnode config:
  remotes:
    - url: nats://<mesh-ip-of-omos-mbp>:7422

NATS JetStream KV / object store replicates across the mesh automatically
```

**NetBird + NATS advantage over Tailscale:**
- NetBird assigns deterministic mesh IPs (configurable CIDR)
- Tailscale assigns 100.x.x.x (CGNAT space — can't change)
- For NATS cluster routes, deterministic IPs are cleaner than MagicDNS hostnames
- Self-hosted NetBird DNS gives you `<node>.omos.mesh` style internal DNS

---

## Migration Plan: Tailscale → NetBird

When ready (not urgent — Tailscale works fine now):

```
Phase 1: Deploy NetBird self-hosted on omos-cloud
  - Docker Compose stack on 100.111.189.116
  - Domain: net.onemindos.dev → NetBird management UI
  - Estimated time: 2 hours

Phase 2: Enroll omos-mbp in NetBird
  - Run netbird-ui or netbird service
  - Both Tailscale + NetBird can coexist during transition
  - Estimated time: 30 min

Phase 3: Enroll IoT nodes one by one
  - ARM64/ARMv6/MIPS binaries from NetBird releases
  - Use Setup Keys for headless/automated enrollment
  - Estimated time: 30 min per node

Phase 4: Update NATS remotes to use NetBird mesh IPs
  - Update server.conf routes

Phase 5: Decommission Tailscale
  - Only after all nodes verified on NetBird mesh
  - Keep Tailscale account active as backup for 30 days
```

---

## Platform Summary Cards

### Tailscale — Keep Now
```
Status:   ✅ ACTIVE (already deployed)
Use for:  Current fleet connectivity
Limit:    Cloud SaaS control plane, $1/IoT device
Keep:     Yes — don't migrate until trigger hit
```

### NetBird — Migrate To
```
Status:   ✅ Ready to deploy
Use for:  Long-term fleet mesh — sovereignty, unlimited IoT
Deploy:   omos-cloud → Docker Compose
License:  BSD-3, free unlimited (Community Edition)
```

### Twingate — Skip
```
Status:   ❌ Wrong architecture
Reason:   User-to-resource, SaaS-only, no mesh, no IoT fleet support
```

### WireGuard (raw) — Foundation
```
Status:   ✅ Running (kernel module on all Linux nodes)
Role:     Transport layer — let Tailscale/NetBird manage it
Never:    Manage directly — too much DIY
```

---

## References

- Tailscale Edge/IoT: https://tailscale.com/use-cases/edge-iot
- Tailscale Pricing: https://tailscale.com/pricing
- Headscale: https://headscale.net
- NetBird: https://netbird.io
- NetBird Self-hosted docs: https://docs.netbird.io/selfhosted/selfhosted-guide
- Twingate: https://www.twingate.com
- WireGuard: https://www.wireguard.com
- NATS Leafnode docs: https://docs.nats.io/running-a-nats-service/configuration/leafnodes
- Prior research: `networking-platforms-all-in-one-2025.md`
- Prior research: `research-unified-networking-mesh.md`

# All-in-One Networking Platforms for IoT + Edge + Cloud + Fleet Management at Scale

## Research Date: May 2025

---

## Executive Summary

**No single platform does everything.** But some come much closer than others. The closest to "one platform" are:

1. **OpenZiti** — most complete zero-trust overlay (network + identity + service discovery + edge SDK)
2. **Cloudflare Zero Trust** — most complete managed solution (but vendor lock-in, not self-hosted)
3. **Pangolin + WireGuard** — best emerging open-source alternative for tunneled reverse proxy + identity
4. **Tailscale/Headscale + Caddy** — pragmatic choice for <10K nodes

For **massive IoT** (100K+ devices), you still need purpose-built device management (Eclipse Hono/Ditto or equivalent) layered on top of a network fabric.

---

## 1. Tailscale for Business / Headscale

**What it is:** WireGuard-based mesh VPN with coordination server (Tailscale SaaS or Headscale self-hosted)

### Features
- **DNS:** MagicDNS (auto-assigns `device.tailnet.ts.net`), split DNS, custom nameservers
- **ACLs:** Tag-based, group-based, fine-grained per-port/protocol policies (JSON/HuJSON)
- **Funnel:** Expose local services to public internet via Tailscale's edge (HTTPS only)
- **Serve:** Reverse proxy local ports/files to your tailnet (private) or Funnel (public)
- **Subnet routers:** Bridge non-Tailscale networks
- **Exit nodes:** Route all traffic through a node
- **Headscale:** Open-source control plane, supports ACLs, DNS, DERP, pre-auth keys

### Coverage Matrix
| Capability | Status |
|---|---|
| DNS + Service Discovery | ✅ MagicDNS |
| Reverse Proxy | ⚠️ Funnel/Serve (basic, no load balancing) |
| Docker integration | ✅ Sidecar container or `--network=host` |
| Edge/IoT | ⚠️ Works but no device lifecycle management |
| Multi-location | ✅ Core strength |
| Auto-scaling | ❌ No orchestration |

### Scale
- Tailscale: Used by companies with 10K+ nodes. DERP relay architecture handles coordination well.
- Headscale: Community reports stable at ~5K nodes. Single-binary, SQLite/Postgres.
- **Not designed for 100K+ IoT devices** — each device needs a full WireGuard key + coordination overhead.

### Verdict
**Best for:** Dev teams, hybrid cloud, remote access. 
**Gap:** No container orchestration, no auto-scaling, no IoT device provisioning. Still need: Docker/K8s orchestrator + reverse proxy (Caddy/Traefik) + monitoring.

---

## 2. Cloudflare Zero Trust + Tunnels

**What it is:** Managed zero-trust platform — reverse proxy, DNS, identity, network tunnels, WARP client

### Features
- **Cloudflare Tunnel (cloudflared):** Outbound-only connections from your infra to CF edge. No public IPs needed.
- **Access:** Identity-aware proxy (OIDC, SAML, mTLS, service tokens)
- **Gateway:** DNS filtering, HTTP inspection, DLP
- **WARP:** Device agent for full tunnel or split tunnel
- **Magic WAN / Magic Transit:** L3 connectivity between sites
- **Workers + R2 + D1:** Edge compute + storage

### Coverage Matrix
| Capability | Status |
|---|---|
| DNS + Service Discovery | ✅ Authoritative DNS, Gateway DNS |
| Reverse Proxy | ✅ Full L7, with caching, WAF, load balancing |
| Docker integration | ⚠️ cloudflared sidecar only |
| Edge/IoT | ⚠️ WARP not designed for constrained IoT |
| Multi-location | ✅ Global anycast edge |
| Auto-scaling | ❌ (Workers scale, but not your infra) |

### Scale
- Cloudflare runs ~20% of the internet. Tunnel can handle millions of connections.
- Used by enterprises with thousands of internal apps exposed via tunnels.

### Verdict
**Closest to all-in-one for web services.** DNS + reverse proxy + identity + DDoS + CDN in one. 
**Gaps:** Not self-hosted. WARP client too heavy for constrained IoT. No device fleet management. No container orchestration. Vendor lock-in.

---

## 3. Pangolin (fosrl/pangolin)

**What it is:** Identity-aware VPN and tunneled reverse proxy based on WireGuard. Open-source Cloudflare Tunnel alternative. **20.9K GitHub stars**, very active (6,608 commits).

### Features
- Self-hosted reverse proxy via WireGuard tunnels
- Identity-aware access (built-in auth, SSO)
- Web UI for managing sites/resources
- Automatic HTTPS (Let's Encrypt)
- Docker-first deployment
- CLI for managing sites
- Multi-site support with a single control plane ("Gerbil" agent on remote sites)

### Coverage Matrix
| Capability | Status |
|---|---|
| DNS + Service Discovery | ⚠️ Basic — you point DNS at Pangolin, it routes |
| Reverse Proxy | ✅ Core feature |
| Docker integration | ✅ Docker-native deployment |
| Edge/IoT | ⚠️ Agent ("Newt") runs on Linux, could work on edge |
| Multi-location | ✅ Multiple sites connected to one panel |
| Auto-scaling | ❌ |

### Scale
- New project (2024-2025), rapidly growing. Not yet proven at massive scale.
- Best for: Homelab to mid-size deployments replacing Cloudflare Tunnels.

### Verdict
**Excellent CF Tunnel replacement** for self-hosters. Not an all-in-one platform — no orchestration, limited service discovery, no IoT device management. Think of it as "self-hosted Cloudflare Tunnel + Access" only.

---

## 4. Teleport

**What it is:** Unified access plane for infrastructure (SSH, K8s, databases, web apps, Windows)

### Features
- Certificate-based access (short-lived certs, no static keys)
- Session recording and audit
- RBAC with identity provider integration
- Kubernetes access, database access proxy
- Machine ID (for service-to-service)
- Device Trust

### Coverage Matrix
| Capability | Status |
|---|---|
| DNS + Service Discovery | ❌ Not its job |
| Reverse Proxy | ⚠️ Application access only (not general reverse proxy) |
| Docker integration | ✅ Via K8s access |
| Edge/IoT | ⚠️ Agent-based, but heavyweight for IoT |
| Multi-location | ✅ Trusted clusters |
| Auto-scaling | ❌ |

### Scale
- Enterprise deployments with thousands of nodes. Used by Nasdaq, Samsung, etc.

### Verdict
**Access plane only** — solves "who can access what" beautifully, but doesn't solve networking, service discovery, or orchestration. You still need everything else. Complement, not replacement.

---

## 5. ZeroTier

**What it is:** Peer-to-peer encrypted virtual network (L2/L3 overlay)

### Features
- Virtual Ethernet (L2) — can bridge entire networks
- Central controller (ZeroTier Central SaaS or self-hosted `ztncui`)
- Flow rules (firewall-like, but for the overlay)
- Multicast support
- Runs on: Linux, macOS, Windows, FreeBSD, iOS, Android, **OpenWrt, RTOS (partial)**

### Coverage Matrix
| Capability | Status |
|---|---|
| DNS + Service Discovery | ❌ (Must layer mDNS or external) |
| Reverse Proxy | ❌ |
| Docker integration | ⚠️ Manual sidecar |
| Edge/IoT | ✅ Lightweight (~15MB RAM), runs on OpenWrt/ARM |
| Multi-location | ✅ Core design |
| Auto-scaling | ❌ |

### Scale
- **ZeroTier claims 9M+ nodes** across all users. Single networks can have 10K+ members.
- IoT: Relatively lightweight but still needs ~15MB RAM. Not for ultra-constrained MCUs.
- **Better IoT fit than Tailscale** due to lighter footprint and L2 support.

### Verdict
**Great mesh for IoT-capable Linux devices** (Raspberry Pi, routers, gateways). But it's ONLY networking — no service discovery, no reverse proxy, no orchestration. Need to layer everything else on top.

---

## 6. Rancher / SUSE Edge

**What it is:** Kubernetes management platform with edge-specific distributions (K3s, RKE2)

### Features
- **K3s:** Lightweight K8s for edge (~70MB binary, runs on RPi)
- **Rancher:** Multi-cluster management UI
- **Fleet:** GitOps-based cluster management at scale (designed for 1M+ clusters)
- **Elemental:** Immutable OS provisioning for edge
- **Akri:** Discover and use leaf devices (cameras, sensors) as K8s resources

### Coverage Matrix
| Capability | Status |
|---|---|
| DNS + Service Discovery | ✅ CoreDNS + K8s services |
| Reverse Proxy | ✅ Ingress controllers (Traefik built into K3s) |
| Docker/Container | ✅ Core function |
| Edge/IoT | ✅ Designed for it (K3s) |
| Multi-location | ✅ Fleet manages 1000s of clusters |
| Auto-scaling | ✅ HPA, KEDA, cluster autoscaler |

### Scale
- **SUSE/Rancher claims Fleet tested to 1 million clusters.** Real deployments: Telcos, retail (thousands of edge locations).
- K3s runs on devices with 512MB RAM.

### Verdict
**Closest to all-in-one for container workloads at edge.** If your entire stack is containerized, K3s + Fleet + Traefik gives you networking + service discovery + reverse proxy + orchestration + auto-scaling + multi-location. **Still missing:** WireGuard mesh/VPN between sites (add Tailscale/ZeroTier), IoT device-level management for non-K8s devices.

---

## 7. Eclipse IoT Stack (Hono, Ditto, Kura)

**What it is:** Purpose-built open-source IoT platform components

### Components
- **Eclipse Hono:** Device connectivity layer (MQTT, AMQP, HTTP). Handles millions of device connections.
- **Eclipse Ditto:** Digital twins — mirrors device state, provides REST/WebSocket API
- **Eclipse Kura:** IoT gateway framework (Java-based, runs on gateways)
- **Eclipse Hawkbit:** OTA firmware updates
- **Eclipse Vorto:** Device capability modeling

### Coverage Matrix
| Capability | Status |
|---|---|
| DNS + Service Discovery | ❌ (Relies on K8s/infra) |
| Reverse Proxy | ❌ |
| Docker integration | ✅ Deploys on K8s |
| Edge/IoT | ✅ **Purpose-built**, handles millions of devices |
| Multi-location | ⚠️ Needs infrastructure layer |
| Auto-scaling | ⚠️ Hono scales horizontally on K8s |

### Scale
- **Eclipse Hono: Tested to 1M+ concurrent device connections** (Bosch production use).
- Ditto: Used by BMW, Bosch for digital twin scenarios.

### Verdict
**Best for pure IoT device management at massive scale** but it's a component stack, not all-in-one. You need K8s underneath, networking on top, and your own reverse proxy + DNS. It solves the IoT-specific problems (protocol translation, device auth, OTA, twins) that nothing else here addresses.

---

## 8. Self-Hosted Equivalents of AWS IoT Greengrass / Azure IoT Edge

### Self-hosted alternatives:
- **EdgeX Foundry** (Linux Foundation): Microservice-based IoT edge platform. Plugin architecture. Handles device protocols, local analytics, data export.
- **ThingsBoard:** Open-source IoT platform (device management, dashboards, rule engine). CE edition is free.
- **MainFlux / Magistrala:** Open-source IoT messaging platform (MQTT, CoAP, HTTP normalized).
- **EMQ X / EMQX:** MQTT broker at massive scale (100M+ connections claimed). Open-source.

### Best Self-Hosted Combo:
**EMQX (broker) + ThingsBoard (management/UI) + K3s (orchestration)** = self-hosted IoT Greengrass equivalent.

### Scale
- EMQX: **100M MQTT connections** in single cluster (benchmarked)
- ThingsBoard: Tens of thousands of devices in production
- EdgeX: Reference implementations, growing adoption

---

## 9. OpenZiti (by NetFoundry)

**What it is:** Open-source, programmable zero-trust network overlay. **4.2K stars, very active development.**

### Features
- **Zero trust from the ground up** — no listening ports, no inbound firewall rules needed
- **Application-embedded SDKs** (Go, C, Python, Node.js, .NET, Java, Swift) — network is IN the app
- **Controller + Router architecture** — controllers manage policy, routers forward traffic
- **Service-based model** — define services, bind to them from anywhere
- **Identity management** built in (x509 certificates, enrollment)
- **Smart routing** — picks best path through mesh dynamically
- **Dark services** — services invisible to anyone without identity + policy
- **BrowZer:** Zero-trust access from browsers without client install
- **Ziti Desktop Edge / Mobile Edge:** Client apps for end users

### Coverage Matrix
| Capability | Status |
|---|---|
| DNS + Service Discovery | ✅ Intercept DNS, service-based discovery |
| Reverse Proxy | ✅ (via tunnelers — services appear local) |
| Docker integration | ✅ SDK-embedded or sidecar tunneler |
| Edge/IoT | ✅ C SDK for constrained devices, designed for it |
| Multi-location | ✅ Core design — mesh of routers |
| Auto-scaling | ⚠️ Routers can be added dynamically, but no orchestration |

### Scale
- NetFoundry (commercial arm) runs **millions of endpoints** on OpenZiti fabric.
- Designed for IoT from inception — C SDK is small enough for gateways.

### Verdict
**🏆 Most architecturally complete zero-trust network platform.** It's the only one where the network is truly programmable and embeddable. Services are dark by default. Identity is mandatory. SDK approach means no sidecars needed. **Biggest gap:** No container orchestration (pair with K3s). No IoT device lifecycle management (pair with Eclipse stack). But for the pure networking + identity + service mesh layer, it's the most "all-in-one."

---

## 10. Ockam

**What it is:** End-to-end encrypted, mutually authenticated channels between any two points

### Features
- **Secure channels** through any transport (TCP, UDP, Kafka, relays)
- **Vault-based key management** (TPM, HSM, cloud KMS)
- **Credential-based access control** (attribute-based, verifiable credentials)
- **Portals:** Expose TCP services through encrypted relay (like magic wormholes)
- **Kafka encrypted topics** — E2E encryption through Kafka without Kafka seeing plaintext
- **Enrollment and credential issuance**

### Coverage Matrix
| Capability | Status |
|---|---|
| DNS + Service Discovery | ❌ |
| Reverse Proxy | ⚠️ Portals can proxy TCP |
| Docker integration | ✅ CLI-based sidecar |
| Edge/IoT | ✅ Designed for IoT (Rust, small binary) |
| Multi-location | ✅ Relay-based |
| Auto-scaling | ❌ |

### Scale
- Focused use cases: Industrial IoT, supply chain. Production deployments exist but not published numbers.

### Verdict
**Niche but powerful** for E2E encryption requirements (compliance, zero-trust data pipelines). Not a general networking platform. Complements other tools rather than replacing them.

---

## 11. Netmaker Pro

**What it is:** WireGuard-based overlay network with enterprise management features

### Features
- Auto-generates WireGuard configs for mesh or hub-spoke
- **Egress/Ingress gateways** (bridge external networks)
- **ACLs** at the network level
- **DNS** (private DNS for the overlay)
- **Remote Access Gateways** (VPN-like access for clients)
- **Pro:** SSO, HA, metrics, Kubernetes operator, branding
- **Netclient:** Lightweight agent

### Coverage Matrix
| Capability | Status |
|---|---|
| DNS + Service Discovery | ✅ Built-in DNS for overlay |
| Reverse Proxy | ❌ |
| Docker integration | ⚠️ Netclient in container |
| Edge/IoT | ⚠️ Agent needs Linux, ~25MB |
| Multi-location | ✅ |
| Auto-scaling | ❌ |

### Scale
- Smaller community than Tailscale/ZeroTier. Used in enterprise pilots.
- **Hundreds to low thousands of nodes** in practice.

### Verdict
**WireGuard management layer** — similar space to Tailscale but more network-engineer-friendly (supports multiple networks, relay servers, egress). Still only networking. Need everything else on top.

---

## 12. NEW 2025-2026 Projects

### **Defguard** (2024-2025)
- Open-source WireGuard VPN + MFA + LDAP + OIDC provider + enrollment
- All-in-one identity + VPN for organizations
- Not IoT-scale but interesting for unified identity+network

### **Firezone** (2024-2025, rewritten)
- WireGuard-based zero-trust access. Complete rewrite in Rust/Elixir.
- Replaces corporate VPN. Policy engine, IdP integration.
- Focus: remote access, not IoT

### **Pomerium** (2024-2025)
- Identity-aware reverse proxy (like CF Access, self-hosted)
- Integrates with IdPs, per-request authorization
- Pairs well with any mesh VPN

### **Sidero/Omni by Talos** (2025)
- SaaS + self-hosted bare-metal Kubernetes management
- Talos Linux (immutable, API-managed) + Omni (multi-cluster management)
- Interesting for edge: immutable OS + K8s from bare metal up

### **Defined Networking (Nebula-based)** (2025)
- Commercial offering on top of Slack's Nebula overlay
- Lightweight mesh, good for edge. Management UI.

### **Netbird** (2024-2025)
- Open-source WireGuard mesh + identity + access control
- Combines network + SSO + ACL in one platform
- Growing rapidly, 12K+ GitHub stars
- Peer-to-peer with relay fallback, NetBird agent is lightweight

---

## Comparative Summary

| Platform | DNS | Svc Disc | Rev Proxy | Docker | IoT | Multi-Loc | AutoScale | Self-Hosted | Max Scale |
|---|---|---|---|---|---|---|---|---|---|
| Tailscale/Headscale | ✅ | ⚠️ | ⚠️ | ✅ | ⚠️ | ✅ | ❌ | ✅ | ~10K |
| Cloudflare ZT | ✅ | ✅ | ✅ | ⚠️ | ❌ | ✅ | ⚠️ | ❌ | Millions |
| Pangolin | ⚠️ | ⚠️ | ✅ | ✅ | ⚠️ | ✅ | ❌ | ✅ | ~1K (early) |
| Teleport | ❌ | ❌ | ⚠️ | ✅ | ❌ | ✅ | ❌ | ✅ | ~10K |
| ZeroTier | ❌ | ❌ | ❌ | ⚠️ | ✅ | ✅ | ❌ | ✅ | ~10K/net |
| **Rancher/K3s/Fleet** | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | **1M clusters** |
| Eclipse IoT | ❌ | ⚠️ | ❌ | ✅ | ✅ | ⚠️ | ⚠️ | ✅ | **1M+ devices** |
| **OpenZiti** | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ⚠️ | ✅ | **Millions** |
| Ockam | ❌ | ❌ | ⚠️ | ✅ | ✅ | ✅ | ❌ | ✅ | Unknown |
| Netmaker | ✅ | ⚠️ | ❌ | ⚠️ | ⚠️ | ✅ | ❌ | ✅ | ~1K |
| Netbird | ⚠️ | ⚠️ | ❌ | ✅ | ⚠️ | ✅ | ❌ | ✅ | ~5K |

---

## 🏆 Recommended Architectures

### Architecture A: "Closest to All-in-One" (Container-centric)
```
K3s + Fleet (orchestration, service discovery, ingress, autoscaling)
+ Tailscale/Headscale or OpenZiti (secure mesh between locations)
+ Traefik (reverse proxy, built into K3s)
= 2-3 tools covering everything for containerized workloads
```

### Architecture B: "Maximum IoT Scale"
```
OpenZiti (zero-trust network fabric, identity, service mesh)
+ Eclipse Hono/Ditto (device connectivity, digital twins)
+ K3s + Fleet (edge orchestration)
+ EMQX (MQTT at scale)
= 4 tools but handles millions of IoT devices
```

### Architecture C: "Pragmatic Mid-Scale" (<10K nodes)
```
Pangolin (reverse proxy + tunnels + identity)
+ Tailscale/Headscale (mesh between sites)
+ Docker Swarm or K3s (orchestration)
= 3 tools, simple, works today
```

### Architecture D: "If You Must Choose ONE"
**OpenZiti** is the single platform that comes closest to doing it all for the network layer. It handles: zero-trust connectivity, service discovery, DNS interception, identity, policy, and works from cloud to IoT (via C SDK). You still need orchestration (K3s) and IoT device management (Hono/ThingsBoard) on top, but the networking is unified.

---

## Key Insight

**The "one platform" doesn't exist yet.** The problem space spans:
1. **Network connectivity** (mesh, tunnels, routing)
2. **Identity & access** (who can reach what)
3. **Service discovery & DNS** (finding things)
4. **Reverse proxy / ingress** (exposing things)
5. **Container orchestration** (running things)
6. **IoT device lifecycle** (provisioning, OTA, monitoring constrained devices)
7. **Auto-scaling** (responding to load)

No single project covers all 7. The minimum viable stack is **2-3 tools**:
- **OpenZiti or Tailscale** (covers 1-3)
- **K3s + Fleet** (covers 4-5, 7)
- **Eclipse Hono or EMQX** (covers 6, if you have real IoT)

The industry trend is convergence — watch **OpenZiti** (adding more app-layer features) and **Rancher/Fleet** (adding more network-layer features) as they grow toward each other.

# Unified Networking & Service Mesh Platform Research

## Use Case Summary
Growing fleet of 1→100+ nodes across home, office, barn, cloud VMs, edge devices, drones, IoT. Mix of Docker containers and native services. Need: automatic DNS (*.omos.lan), Docker auto-discovery, NAT-traversing mesh, edge/IoT support, self-hosted, no vendor lock-in, never outgrow it.

---

## Platform Comparison Matrix

| Platform | DNS | Discovery | Proxy/LB | Mesh | Edge/IoT | Docker | Certs | Self-Hosted |
|----------|-----|-----------|----------|------|----------|--------|-------|-------------|
| **Nomad + Consul + Traefik** | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| Consul Connect (standalone) | ✅ | ✅ | ⚠️ | ✅ | ✅ | ⚠️ | ✅ | ✅ |
| Headscale + Traefik | ⚠️ | ❌ | ✅ | ✅ | ✅ | ⚠️ | ⚠️ | ✅ |
| K3s | ✅ | ✅ | ✅ | ✅ | ⚠️ | ✅ | ✅ | ✅ |
| Netmaker | ✅ | ❌ | ❌ | ✅ | ✅ | ❌ | ⚠️ | ✅ |
| Nebula + Consul | ⚠️ | ✅ | ❌ | ✅ | ✅ | ⚠️ | ✅ | ✅ |
| NetBird + Consul + Traefik | ⚠️ | ✅ | ✅ | ✅ | ✅ | ⚠️ | ✅ | ✅ |
| WireGuard + CoreDNS + Traefik | ✅ | ❌ | ✅ | ✅ | ⚠️ | ❌ | ⚠️ | ✅ |
| Cilium | ✅ | ✅ | ✅ | ✅ | ❌ | ✅ | ✅ | ✅ |
| Istio/Linkerd | ✅ | ✅ | ✅ | ⚠️ | ❌ | ✅ | ✅ | ✅ |

---

## Detailed Analysis

### 1. 🏆 Nomad + Consul + Traefik (HashiCorp Stack) — RECOMMENDED

**What it handles:** Everything.
- **DNS**: Consul provides `.consul` DNS interface; integrates with CoreDNS/dnsmasq for custom domains (*.omos.lan)
- **Service Discovery**: Consul's core function — automatic registration, health checks, catalog
- **Reverse Proxy/LB**: Traefik watches Consul catalog, auto-configures routes with zero config per service
- **Docker**: Nomad's Docker driver auto-registers containers in Consul with tags for Traefik routing
- **Mesh**: Consul Connect provides mTLS service mesh; WAN federation via mesh gateways across datacenters
- **Edge/IoT**: Nomad client is a single ~75MB binary, runs on ARM/Raspberry Pi; Consul agent is similarly lightweight
- **Certs**: Consul Connect has built-in CA; Traefik handles Let's Encrypt for external certs
- **Multi-location**: Consul WAN federation or cluster peering connects sites; mesh gateways handle NAT traversal

**What it DOESN'T handle (gaps):**
- NAT traversal between locations needs mesh gateways exposed on public IPs OR pair with WireGuard/Nebula for the underlay network
- No built-in VPN/overlay for the physical network layer (Consul mesh is L7, not L3)
- **Fix**: Add Nebula or WireGuard as the L3 underlay; Consul/Nomad/Traefik ride on top

**Resource usage:**
- Consul server: ~200MB RAM per node (3-node cluster recommended)
- Nomad server: ~150MB RAM per node
- Nomad client (edge): ~75MB RAM, single binary
- Traefik: ~50MB RAM
- Total control plane: ~1GB on 3 servers; edge nodes: ~100-150MB each

**Complexity:** Medium. Three binaries to deploy, but all are single-binary, config-file-driven. Excellent documentation. No Kubernetes complexity.

**Scaling ceiling:** Consul tested to 10,000+ nodes. Nomad handles 10,000+ nodes, 2M+ containers. LinkedIn, Cloudflare, Roblox use at scale. You will never outgrow this.

**Heterogeneous support:** ✅ Linux (x86/ARM), macOS, Windows, Raspberry Pi, bare metal, VMs, containers. Single binary per platform.

**Real-world examples:** Cloudflare (edge network), Roblox (game servers), CircleCI, Pandora, Jet.com, Trivago.

**Graduate from it?** No. This is the endgame stack. Scales from 1 node (dev) to planet-scale.

---

### 2. Headscale + Traefik (+ Consul optional)

**What it handles:**
- **Mesh**: Headscale (self-hosted Tailscale) provides WireGuard mesh with NAT traversal, MagicDNS
- **DNS**: MagicDNS gives hostnames to all nodes automatically
- **Proxy**: Traefik on each location provides HTTP routing
- **Edge**: Tailscale client runs on everything including IoT

**Gaps:**
- No service discovery (only node-level DNS, not per-container)
- No Docker auto-registration
- No orchestration or scheduling
- Need to manually configure Traefik routes or add Consul for discovery
- MagicDNS is node-level, not service-level

**Resource usage:** Very light — Tailscale client ~30MB RAM per node, Headscale server ~100MB

**Complexity:** Low for mesh; medium when adding service discovery

**Scaling:** Headscale supports thousands of nodes. L3 mesh only though.

**Heterogeneous:** ✅ Excellent — iOS, Android, Linux, macOS, Windows, embedded

**Verdict:** Excellent L3 underlay, but needs Consul+Traefik on top for service-level features. Good companion to the Nomad stack.

---

### 3. K3s (Lightweight Kubernetes)

**What it handles:**
- DNS (CoreDNS built-in), service discovery (native), proxy/LB (Traefik included), Docker/containerd, certs (cert-manager)
- Multi-node within a cluster

**Gaps:**
- Multi-location/WAN: Kubernetes was NOT designed for high-latency WAN links between nodes. Split-brain issues.
- Edge: K3s runs on Pi, but Kubernetes control plane assumes reliable networking
- IoT: Kubelet is ~100MB minimum; too heavy for constrained devices
- NAT traversal: None built-in; need separate VPN underlay
- Multi-cluster federation is complex (requires tools like Submariner, Liqo, or Skupper)

**Resource usage:** Server node: 512MB+ RAM minimum; agent: 256MB+. Heavier than Nomad.

**Complexity:** High. Kubernetes is inherently complex. YAML sprawl. CRDs. Operators. Networking plugins.

**Scaling:** Within a cluster: thousands of nodes. Across locations: painful.

**Heterogeneous:** Linux only for nodes (ARM + x86). No macOS/Windows workers.

**Real-world:** Many edge deployments (Rancher, SUSE). But always within a single "cluster" concept.

**Graduate from it?** No, but you might wish you hadn't started with it. Kubernetes adds operational overhead that the Nomad stack avoids for this use case.

---

### 4. Consul Connect (Standalone, no Nomad)

**What it handles:**
- DNS, service discovery, health checks, service mesh (mTLS), multi-datacenter WAN federation
- Built-in Connect CA for certificate management

**Gaps:**
- No container orchestration (you manage Docker yourself with docker-compose or systemd)
- No reverse proxy (need to add Traefik/Nginx)
- No workload scheduling
- Services must self-register or use registrator for Docker

**Resource usage:** ~200MB per server, ~75MB per client

**Complexity:** Low-medium. Single binary. Very well documented.

**Scaling:** 10,000+ nodes proven

**Heterogeneous:** ✅ All platforms

**Verdict:** Great service mesh and discovery layer. Add Traefik for proxy. If you don't need orchestration (manual docker-compose is fine), this + Traefik is simpler than full Nomad stack.

---

### 5. NetBird + Consul + Traefik

**What it handles:**
- **Mesh**: NetBird (open-source, self-hosted) provides WireGuard mesh with NAT traversal, ACLs
- **Discovery**: Consul provides service registration
- **Proxy**: Traefik routes based on Consul catalog

**Gaps:**
- Three separate systems to maintain
- NetBird doesn't do service discovery or DNS (node-level only)
- No container orchestration
- No integrated DNS (need separate CoreDNS or Consul DNS)

**Resource usage:** Light per component, but 3 systems to run

**Complexity:** Medium-high (3 independent systems)

**Scaling:** Each component scales well individually

**Heterogeneous:** ✅ NetBird has wide platform support

**Verdict:** Works but is the "duct tape" approach. The Nomad+Consul+Traefik stack is equally 3 tools but far more integrated.

---

### 6. Nebula + Consul

**What it handles:**
- **Mesh**: Nebula (created at Slack) provides encrypted overlay network, NAT traversal via "lighthouses"
- **Discovery**: Consul on top of Nebula network

**Gaps:**
- No proxy/LB (add Traefik)
- No orchestration
- No Docker auto-registration
- Nebula DNS is basic (host-level)
- Certificate management is manual (Nebula has its own CA but it's for mesh certs only)

**Resource usage:** Nebula: ~15MB RAM, extremely lightweight. Excellent for IoT.

**Complexity:** Medium. Nebula is simple but Consul adds complexity.

**Scaling:** Nebula tested at Slack with 50,000+ nodes. Consul scales to 10,000+.

**Heterogeneous:** ✅ Nebula runs on everything including very constrained devices

**Real-world:** Slack's entire infrastructure

**Verdict:** Best L3 underlay for IoT/constrained devices. Pair with Consul+Traefik for full stack.

---

### 7. Netmaker

**What it handles:**
- WireGuard mesh with management UI, NAT traversal (TURN relay), DNS, ACLs
- Node auto-join, network segmentation

**Gaps:**
- No service discovery (network-level only)
- No reverse proxy
- No Docker integration
- No container orchestration
- DNS is node-level, not service-level
- Limited health checking

**Resource usage:** Server: ~200MB; clients: ~50MB

**Complexity:** Low for mesh setup (nice UI)

**Scaling:** Hundreds of nodes; not tested at thousands

**Heterogeneous:** Linux, macOS, Docker, embedded (WireGuard kernel support needed)

**Verdict:** Good WireGuard mesh manager but only handles L3. Still need all application-layer tooling on top. Less mature than Nebula or Headscale.

---

### 8. Defined Networking (from Nebula creators)

**What it handles:**
- Commercial/hosted version of Nebula with management plane
- Same mesh capabilities as Nebula but with UI and API

**Gaps:** Same as Nebula — only L3 overlay. Everything else must be added.

**Verdict:** If you want managed Nebula, but doesn't solve the full problem.

---

### 9. Cilium (eBPF-based)

**What it handles:**
- Advanced L3-L7 networking, load balancing, service mesh (sidecar-free via eBPF)
- Network policies, observability, encryption
- Cluster Mesh for multi-cluster Kubernetes

**Gaps:**
- **Requires Kubernetes** — not standalone
- No edge/IoT support outside K8s
- No NAT traversal / WAN mesh
- Linux-only (needs eBPF, kernel 4.19+)
- No heterogeneous device support

**Resource usage:** ~300MB RAM per node; requires modern Linux kernel

**Complexity:** High (Kubernetes prerequisite + eBPF complexity)

**Scaling:** Excellent within Kubernetes

**Verdict:** Best-in-class for Kubernetes networking but completely irrelevant for heterogeneous edge/IoT fleets outside Kubernetes.

---

### 10. Istio / Linkerd

**What it handles:**
- L7 service mesh, mTLS, traffic management, observability
- Built for Kubernetes microservices

**Gaps:**
- **Kubernetes-only**
- Heavy sidecar overhead (Istio: ~100MB per pod; Linkerd lighter at ~20MB)
- No edge/IoT, no multi-location mesh, no NAT traversal
- No heterogeneous device support
- No Docker standalone support

**Resource usage:** Istio control plane: 2-4GB. Linkerd: ~500MB. Plus per-pod sidecars.

**Complexity:** Very high (Istio is notoriously complex)

**Scaling:** Fine within a cluster

**Verdict:** Wrong tool entirely for this use case. These are Kubernetes-internal mesh tools, not cross-location fleet tools.

---

### 11. WireGuard + CoreDNS + Traefik (Manual)

**What it handles:**
- Mesh: WireGuard point-to-point tunnels
- DNS: CoreDNS with custom zones (*.omos.lan)
- Proxy: Traefik with file/docker provider

**Gaps:**
- No NAT traversal (manual port forwarding needed)
- No service discovery (manual DNS records or scripts)
- No auto-registration of containers
- No health checking
- Scaling WireGuard config is O(n²) — every node needs config for every other node
- No management plane

**Resource usage:** Minimal per component

**Complexity:** Starts simple, becomes unmanageable past ~10 nodes. WireGuard config explosion.

**Scaling:** Breaks down at ~20+ nodes without automation tooling

**Heterogeneous:** ✅ WireGuard runs everywhere

**Verdict:** Fine for 3-5 nodes. Becomes a maintenance nightmare. You WILL need to graduate from this.

---

## Final Recommendation

### The Winning Stack: **Nebula + Nomad + Consul + Traefik**

This four-tool combination covers every requirement with no gaps:

| Layer | Tool | Role |
|-------|------|------|
| **L3 Overlay Network** | **Nebula** | WireGuard-style encrypted mesh, NAT traversal, works on IoT/constrained devices (15MB) |
| **Orchestration** | **Nomad** | Schedules Docker containers + raw binaries + Java + anything. Single binary, ARM-ready. |
| **Service Mesh & Discovery** | **Consul** | DNS (*.omos.lan via forwarding), service catalog, health checks, mTLS, multi-DC federation |
| **Reverse Proxy** | **Traefik** | Auto-discovers services from Consul catalog, Let's Encrypt certs, HTTP/gRPC routing |

### Why This Combination Wins:

1. **Never outgrow it**: Each component is proven at 10,000+ nodes
2. **Truly heterogeneous**: Runs on Raspberry Pi, drones, cloud VMs, Macs, bare metal
3. **Self-hosted & sovereign**: All open-source (BSL for Hashi tools but source-available and self-hostable)
4. **Minimal resource usage**: Edge node total: ~150MB RAM (Nebula 15MB + Nomad client 75MB + Consul client 50MB)
5. **Single-binary deploys**: Each tool is one binary. No JVM, no runtime deps.
6. **Docker auto-discovery**: Nomad registers Docker tasks in Consul automatically; Traefik picks them up
7. **Custom DNS**: Consul DNS + dnsmasq/CoreDNS forwarding = `service.omos.lan` for everything
8. **NAT traversal**: Nebula lighthouses punch through NATs; all nodes talk directly
9. **Zero-config new services**: Deploy a Nomad job → appears in Consul → Traefik routes to it → DNS resolves it
10. **Graceful scaling**: Start with 1 node running everything, split to dedicated servers as you grow

### Deployment Architecture:

```
┌─────────────────────────────────────────────────────────────────┐
│                        CONTROL PLANE (3 nodes)                   │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐                      │
│  │ Consul   │  │ Nomad    │  │ Nebula   │  (can be same 3 VMs) │
│  │ Server   │  │ Server   │  │Lighthouse│                       │
│  └──────────┘  └──────────┘  └──────────┘                      │
└─────────────────────────────────────────────────────────────────┘

┌──────────────┐  ┌──────────────┐  ┌──────────────┐
│  HOME NODE   │  │  OFFICE NODE │  │  CLOUD VM    │
│  Nebula      │  │  Nebula      │  │  Nebula      │
│  Nomad Client│  │  Nomad Client│  │  Nomad Client│
│  Consul Agent│  │  Consul Agent│  │  Consul Agent│
│  Traefik     │  │  Traefik     │  │  Traefik     │
│  [containers]│  │  [containers]│  │  [containers]│
└──────────────┘  └──────────────┘  └──────────────┘

┌──────────────┐  ┌──────────────┐
│  RASPBERRY PI│  │  IoT DEVICE  │
│  Nebula      │  │  Nebula only │
│  Nomad Client│  │  (15MB)      │
│  Consul Agent│  │              │
└──────────────┘  └──────────────┘
```

### Alternative: Simplified Stack (if no orchestration needed)

If you manage containers manually (docker-compose per node) and don't need scheduling:

**Nebula + Consul + Traefik** (3 tools instead of 4)
- Drop Nomad, use docker-compose + [Registrator](https://github.com/gliderlabs/registrator) for Docker→Consul auto-registration
- Simpler but loses scheduling, rolling deploys, and resource bin-packing

### Licensing Note (2024+ HashiCorp BSL):
Consul and Nomad switched to Business Source License. For your use case (self-hosted infrastructure, not selling as a service), this is completely fine. The BSL only restricts offering these as a managed service to third parties. OpenTofu-style forks exist (OpenBao for Vault) if this concerns you, but Consul/Nomad BSL is non-restrictive for internal use.

---

## Quick-Start Path

1. **Week 1**: Install Nebula on all nodes (lighthouse on a cloud VM with public IP)
2. **Week 2**: Deploy Consul cluster (3 servers on reliable nodes)
3. **Week 3**: Deploy Nomad (servers co-located with Consul; clients on all worker nodes)
4. **Week 4**: Deploy Traefik as a Nomad job watching Consul catalog
5. **Week 5**: Configure DNS forwarding (dnsmasq: `*.omos.lan` → Consul DNS)
6. **Done**: Every new Nomad job auto-gets DNS + routing + discovery + mesh encryption

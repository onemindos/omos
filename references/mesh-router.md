# Mesh Router (NSL.SH)

**Repo:** github.com/mesh-router (mesh-router-root monorepo)  
**Language:** Go  
**Role in OneMind:** WireGuard mesh routing for homestead nodes behind NAT

Distributed mesh routing that enables secure traffic between public servers and private homestead nodes. Supports direct IP routing (low latency for public nodes) and WireGuard VPN tunneling (for nodes behind NAT with no public IP).

## Components

- `mesh-router-gateway` — public entry point, routes to direct or VPN path
- Direct Route — agent-based, for nodes with public IPs
- VPN Tunnel — WireGuard tunnel for nodes behind NAT/CGNAT

## Why it matters for OneMind

Homestead nodes (Mac mini, Raspberry Pi edge nodes, field devices) often sit behind CGNAT or home routers with no public IP. Mesh Router + WireGuard provides:
- Inbound connectivity to NATS leaf nodes without port forwarding
- Secure mesh between homestead fleet members
- Low-latency path when public IPs are available, automatic VPN fallback

## Integration point

Mesh Router sits beneath the NATS transport layer — NATS leaf node connections run through the WireGuard tunnel, giving full NATS connectivity to NAT-ed homestead nodes.

## See Also

- `natpunch-go` (NAT puncher for WireGuard) — in [[nats-tactical-repos]]
- `ikto` (NATS-based WireGuard mesh) — in [[nats-tactical-repos]]

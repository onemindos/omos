---
title: "Tailscale"
type: reference
status: active
tags: [platform, infrastructure, networking, vpn, mesh, wireguard]
---
# Tailscale

> WireGuard-based mesh VPN. Zero-config, connects all devices into one private network.

## Key Features
| Feature | Notes |
|---------|-------|
| Mesh VPN | All devices connect peer-to-peer |
| Zero config | No port forwarding, no firewall rules |
| WireGuard | Fast, modern crypto |
| MagicDNS | Access devices by name (e.g. `macbook.tail`) |
| Exit nodes | Route traffic through any node |
| Subnet routes | Expose LAN devices to tailnet |
| ACLs | Fine-grained access control |
| Funnel | Expose services to public internet |

## CLI
```bash
tailscale status          # Show connected devices
tailscale ping node       # Test connectivity
tailscale up              # Connect
tailscale down            # Disconnect
tailscale ip              # Show Tailscale IP
tailscale ssh user@node   # SSH over tailnet
```

## Integration
- All OneMind nodes on same tailnet
- SSH without port forwarding
- Access Home Assistant remotely
- Fleet connectivity backbone

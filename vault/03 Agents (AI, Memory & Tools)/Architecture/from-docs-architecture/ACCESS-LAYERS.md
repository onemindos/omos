# OneMind OS - Access Layers Guide

## Overview

OneMind OS uses three complementary access layers that serve different purposes. This document explains how Cloudflare, Tailscale, and Apache Guacamole work together.

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                              ACCESS ARCHITECTURE                                 │
├─────────────────────────────────────────────────────────────────────────────────┤
│                                                                                  │
│                           INTERNET                                               │
│                              │                                                   │
│                              ▼                                                   │
│  ┌───────────────────────────────────────────────────────────────────────────┐  │
│  │                     LAYER 1: CLOUDFLARE                                   │  │
│  │                     (Public Access Gateway)                               │  │
│  │                                                                           │  │
│  │   Purpose: Public-facing services with security                          │  │
│  │                                                                           │  │
│  │   ┌─────────────┐    ┌─────────────┐    ┌─────────────┐                  │  │
│  │   │    DNS      │    │    CDN      │    │   DDoS      │                  │  │
│  │   │  Routing    │    │   Cache     │    │ Protection  │                  │  │
│  │   └─────────────┘    └─────────────┘    └─────────────┘                  │  │
│  │                              │                                            │  │
│  │                    Cloudflare Tunnel                                      │  │
│  │                      (No open ports)                                      │  │
│  │                              │                                            │  │
│  │   Routes:                    ▼                                            │  │
│  │   • onemindos.dev ──────────────────► UI (:3000)                         │  │
│  │   • api.onemindos.dev ──────────────► Agno API (:8080)                   │  │
│  │                                                                           │  │
│  └───────────────────────────────────────────────────────────────────────────┘  │
│                                                                                  │
│                                                                                  │
│  ┌───────────────────────────────────────────────────────────────────────────┐  │
│  │                     LAYER 2: TAILSCALE                                    │  │
│  │                     (Private Mesh VPN)                                    │  │
│  │                                                                           │  │
│  │   Purpose: Direct private access to all services                         │  │
│  │                                                                           │  │
│  │   ┌─────────────────────────────────────────────────────────────────┐    │  │
│  │   │                     TAILSCALE NETWORK                           │    │  │
│  │   │                      100.x.x.x/24                                │    │  │
│  │   │                                                                  │    │  │
│  │   │  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐        │    │  │
│  │   │  │ Your Mac │  │ Cloud VM │  │ Jetson 1 │  │ Jetson 2 │        │    │  │
│  │   │  │100.x.x.1 │  │100.x.x.2 │  │100.x.x.3 │  │100.x.x.4 │        │    │  │
│  │   │  └────┬─────┘  └────┬─────┘  └────┬─────┘  └────┬─────┘        │    │  │
│  │   │       │             │             │             │               │    │  │
│  │   │       └─────────────┴─────────────┴─────────────┘               │    │  │
│  │   │                         WireGuard                                │    │  │
│  │   │                       (Encrypted)                                │    │  │
│  │   └─────────────────────────────────────────────────────────────────┘    │  │
│  │                                                                           │  │
│  │   Access from your devices:                                              │  │
│  │   • ssh user@100.x.x.2           (SSH to cloud)                         │  │
│  │   • http://100.x.x.2:19999       (Netdata)                              │  │
│  │   • http://100.x.x.2:5432        (Postgres direct)                      │  │
│  │   • http://100.x.x.3:5000        (Frigate on Jetson)                    │  │
│  │                                                                           │  │
│  └───────────────────────────────────────────────────────────────────────────┘  │
│                                                                                  │
│                                                                                  │
│  ┌───────────────────────────────────────────────────────────────────────────┐  │
│  │                     LAYER 3: APACHE GUACAMOLE                             │  │
│  │                     (Remote Desktop Gateway)                              │  │
│  │                                                                           │  │
│  │   Purpose: Browser-based GUI access with audit trail                     │  │
│  │                                                                           │  │
│  │   ┌─────────────────────────────────────────────────────────────────┐    │  │
│  │   │                    GUACAMOLE GATEWAY                            │    │  │
│  │   │                                                                  │    │  │
│  │   │   Browser ────► Guacamole ────► Target Machine                  │    │  │
│  │   │                     │                                            │    │  │
│  │   │              ┌──────┴──────┐                                     │    │  │
│  │   │              │  Protocols  │                                     │    │  │
│  │   │              ├─────────────┤                                     │    │  │
│  │   │              │ • VNC       │──► Jetson Desktop                   │    │  │
│  │   │              │ • RDP       │──► Windows Server                   │    │  │
│  │   │              │ • SSH       │──► Any Linux                        │    │  │
│  │   │              │ • Telnet    │──► Network devices                  │    │  │
│  │   │              └─────────────┘                                     │    │  │
│  │   │                                                                  │    │  │
│  │   │   Features:                                                      │    │  │
│  │   │   • Session recording                                            │    │  │
│  │   │   • Multi-user access                                            │    │  │
│  │   │   • Audit logging                                                │    │  │
│  │   │   • LDAP/TOTP auth                                               │    │  │
│  │   │   • Works through firewalls                                      │    │  │
│  │   └─────────────────────────────────────────────────────────────────┘    │  │
│  │                                                                           │  │
│  └───────────────────────────────────────────────────────────────────────────┘  │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

## Comparison Table

| Feature | Cloudflare | Tailscale | Guacamole |
|---------|------------|-----------|-----------|
| **Purpose** | Public access | Private access | Desktop access |
| **Network Layer** | DNS/HTTP | VPN (WireGuard) | Application |
| **Who can access** | Anyone (with auth) | Your devices only | Authorized users |
| **Protocols** | HTTP/HTTPS/WS | Any TCP/UDP | VNC/RDP/SSH |
| **Open ports** | None (tunnel) | None (NAT traversal) | One (HTTPS) |
| **DDoS protection** | ✅ Yes | N/A | N/A |
| **CDN/Caching** | ✅ Yes | N/A | N/A |
| **GUI access** | N/A | No (network only) | ✅ Yes |
| **Session recording** | N/A | N/A | ✅ Yes |
| **Audit logging** | ✅ Yes | ✅ Yes | ✅ Yes |
| **Cost** | Free tier | Free (100 devices) | Free (self-hosted) |

## When to Use Each

### Use Cloudflare When:
- Users need to access OneMind UI from anywhere
- Public API access is required
- You need DDoS protection
- You want CDN caching for static assets
- You don't want to expose server ports

### Use Tailscale When:
- You need direct SSH access
- You're accessing internal services (Netdata, Grafana)
- You need database access
- You're debugging or developing
- You want encrypted device-to-device communication

### Use Guacamole When:
- You need GUI/desktop access to remote machines
- You want browser-only access (no client install)
- You need session recording for compliance
- Multiple users need access to same machines
- You're accessing from a restricted network

## Setup Instructions

### 1. Cloudflare Tunnel Setup

```bash
# Install cloudflared on server
curl -L https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-amd64 -o cloudflared
chmod +x cloudflared
sudo mv cloudflared /usr/local/bin/

# Authenticate
cloudflared tunnel login

# Create tunnel
cloudflared tunnel create onemind

# Configure tunnel
cat > ~/.cloudflared/config.yml << EOF
tunnel: <TUNNEL_ID>
credentials-file: /root/.cloudflared/<TUNNEL_ID>.json

ingress:
  - hostname: onemindos.dev
    service: http://localhost:3000
  - hostname: api.onemindos.dev
    service: http://localhost:8080
  - service: http_status:404
EOF

# Run as service
sudo cloudflared service install
sudo systemctl start cloudflared
sudo systemctl enable cloudflared
```

### 2. Tailscale Setup

```bash
# Install on each machine
curl -fsSL https://tailscale.com/install.sh | sh

# Authenticate (first time)
sudo tailscale up

# On servers, enable as exit node (optional)
sudo tailscale up --advertise-exit-node

# Enable subnet routing (optional)
sudo tailscale up --advertise-routes=192.168.1.0/24

# Check status
tailscale status
```

### 3. Apache Guacamole Setup

```yaml
# Add to docker-compose.yml
services:
  guacamole:
    image: guacamole/guacamole
    container_name: guacamole
    ports:
      - "8080:8080"
    environment:
      - GUACD_HOSTNAME=guacd
      - POSTGRES_HOSTNAME=guac-postgres
      - POSTGRES_DATABASE=guacamole_db
      - POSTGRES_USER=guacamole
      - POSTGRES_PASSWORD=${GUAC_DB_PASSWORD}
    depends_on:
      - guacd
      - guac-postgres
    restart: unless-stopped

  guacd:
    image: guacamole/guacd
    container_name: guacd
    restart: unless-stopped

  guac-postgres:
    image: postgres:15-alpine
    container_name: guac-postgres
    environment:
      - POSTGRES_DB=guacamole_db
      - POSTGRES_USER=guacamole
      - POSTGRES_PASSWORD=${GUAC_DB_PASSWORD}
    volumes:
      - guac_db:/var/lib/postgresql/data
      - ./guac-initdb.sql:/docker-entrypoint-initdb.d/init.sql
    restart: unless-stopped

volumes:
  guac_db:
```

Initialize database:
```bash
# Generate init script
docker run --rm guacamole/guacamole /opt/guacamole/bin/initdb.sh --postgres > guac-initdb.sql
```

### 4. DNS Configuration

| Domain | Type | Target | Proxy |
|--------|------|--------|-------|
| `onemindos.dev` | CNAME | `<tunnel-id>.cfargotunnel.com` | ✅ |
| `api.onemindos.dev` | CNAME | `<tunnel-id>.cfargotunnel.com` | ✅ |
| `guac.onemindos.dev` | CNAME | `<tunnel-id>.cfargotunnel.com` | ✅ |

## Access Matrix

| Service | Cloudflare URL | Tailscale URL | Guacamole |
|---------|----------------|---------------|-----------|
| OneMind UI | `https://onemindos.dev` | `http://100.x.x.2:3000` | N/A |
| Agno API | `https://api.onemindos.dev` | `http://100.x.x.2:8080` | N/A |
| Netdata | N/A | `http://100.x.x.2:19999` | N/A |
| Grafana | N/A | `http://100.x.x.2:3001` | N/A |
| Postgres | N/A | `100.x.x.2:5432` | N/A |
| Jetson Desktop | N/A | N/A | VNC via Guacamole |
| Server SSH | N/A | `ssh user@100.x.x.2` | SSH via Guacamole |

## Security Considerations

### Cloudflare
- Enable "Always Use HTTPS"
- Configure WAF rules
- Enable bot protection
- Use Access policies for sensitive routes

### Tailscale
- Use ACLs to restrict device access
- Enable MagicDNS for easier access
- Review connected devices regularly
- Use tags for device groups

### Guacamole
- Enable TOTP 2FA
- Use LDAP for user management
- Enable session recording for audits
- Regularly rotate connection passwords

## Example Workflows

### Daily Development
```
1. SSH via Tailscale: ssh user@100.x.x.2
2. Check Netdata: http://100.x.x.2:19999
3. Access app: https://onemindos.dev
```

### Debug Jetson GUI
```
1. Open Guacamole: https://guac.onemindos.dev
2. Login with your credentials
3. Click "Home Jetson" connection
4. Full desktop appears in browser
```

### Share Access with Team Member
```
1. Add user to Guacamole
2. Assign specific connections
3. They access via browser only
4. All sessions are recorded
```

### Emergency Access from Hotel
```
1. Guacamole works through any firewall
2. Access: https://guac.onemindos.dev
3. No VPN client needed
4. Full server access in browser
```

## Troubleshooting

### Cloudflare Tunnel Not Working
```bash
# Check tunnel status
cloudflared tunnel info onemind

# Check service logs
sudo journalctl -u cloudflared -f
```

### Tailscale Not Connecting
```bash
# Check status
tailscale status

# Re-authenticate
sudo tailscale up --reset
```

### Guacamole Connection Failed
```bash
# Check guacd logs
docker logs guacd

# Verify target is reachable from guacd container
docker exec guacd nc -zv 192.168.1.10 22
```

# Cloudflare Tunnel & Zero Trust Configuration

*OneMind Cloud Network Layer - January 2026*

## Overview

OneMind uses **Cloudflare Tunnel** (formerly Argo Tunnel) to expose services securely without opening firewall ports. Combined with **Cloudflare Zero Trust Access**, sensitive applications require email-based authentication.

---

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                    CLOUDFLARE NETWORK FLOW                      │
│                                                                 │
│   User Request                                                  │
│        │                                                        │
│        ▼                                                        │
│   ┌─────────────────────────────────────────────┐              │
│   │         CLOUDFLARE DNS (onemindos.dev)      │              │
│   │                                             │              │
│   │   n8n.onemindos.dev ──► CNAME ──► tunnel    │              │
│   │   toc.onemindos.dev ──► CNAME ──► tunnel    │              │
│   │   etc...                                    │              │
│   └─────────────────────────────────────────────┘              │
│                         │                                       │
│                         ▼                                       │
│   ┌─────────────────────────────────────────────┐              │
│   │           CLOUDFLARE EDGE                   │              │
│   │                                             │              │
│   │   • SSL/TLS termination                     │              │
│   │   • DDoS protection                         │              │
│   │   • WAF (Web Application Firewall)          │              │
│   │   • Rate limiting                           │              │
│   └─────────────────────────────────────────────┘              │
│                         │                                       │
│                         ▼                                       │
│   ┌─────────────────────────────────────────────┐              │
│   │        ZERO TRUST ACCESS CHECK              │              │
│   │                                             │              │
│   │   Protected apps require:                   │              │
│   │   • Email: zeusdelacruz@gmail.com          │              │
│   │   • Email: zeus.vaults@pm.me               │              │
│   │                                             │              │
│   │   Unprotected apps: pass through            │              │
│   └─────────────────────────────────────────────┘              │
│                         │                                       │
│                         ▼                                       │
│   ┌─────────────────────────────────────────────┐              │
│   │     CLOUDFLARE TUNNEL (cloudflared)         │              │
│   │                                             │              │
│   │   Running as systemd service on             │              │
│   │   onemind-cloud-1 (AWS Lightsail)          │              │
│   │                                             │              │
│   │   Maintains outbound connection to          │              │
│   │   Cloudflare edge (no inbound ports)        │              │
│   └─────────────────────────────────────────────┘              │
│                         │                                       │
│                         ▼                                       │
│   ┌─────────────────────────────────────────────┐              │
│   │        LOCAL DOCKER SERVICES                │              │
│   │                                             │              │
│   │   localhost:8080 ──► Agno                   │              │
│   │   localhost:5678 ──► N8N                    │              │
│   │   localhost:3002 ──► TOC                    │              │
│   │   etc...                                    │              │
│   └─────────────────────────────────────────────┘              │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## Credentials (1Password)

**Item Name:** "Cloudflare - OneMind OS"

| Field | Value |
|-------|-------|
| **Account Email** | zeus.vaults@pm.me |
| **Account ID** | ce13703e0b76938c179ebc51e2945c7e |
| **Zone ID** | 32beabc2a2f2debece7b5c55cfdd9b2a |
| **Tunnel ID** | 7ddfa54f-c672-4dd0-ad50-19b431b5ad72 |
| **Tunnel Name** | onemind-os-tunnel |
| **Global API Key** | (stored in 1Password) |

---

## Tunnel Configuration

### Current Ingress Rules

The tunnel routes traffic based on hostname to local services:

| Hostname | Service | Port | Notes |
|----------|---------|------|-------|
| command.onemindos.dev | http://localhost:3000 | 3000 | Agent UI (not running) |
| api.onemindos.dev | http://localhost:8000 | 8000 | API Gateway |
| agno.onemindos.dev | http://localhost:8080 | 8080 | LEGACY Brain |
| letta.onemindos.dev | http://localhost:8283 | 8283 | Memory Server |
| n8n.onemindos.dev | http://localhost:5678 | 5678 | Workflow Automation |
| kestra.onemindos.dev | http://localhost:8088 | 8088 | Data Pipelines |
| status.onemindos.dev | http://localhost:3001 | 3001 | Uptime Kuma |
| portainer.onemindos.dev | http://localhost:9000 | 9000 | Docker Management |
| nats.onemindos.dev | http://localhost:8222 | 8222 | NATS Monitor |
| toc.onemindos.dev | http://localhost:3002 | 3002 | Dashboard |
| vnc.onemindos.dev | http://localhost:6080 | 6080 | Remote Desktop |
| actions.onemindos.dev | http://localhost:8001 | 8001 | Action Gateway |
| * (catch-all) | http_status:404 | - | Default 404 |

### Tunnel Service (systemd)

```bash
# Service location
/etc/systemd/system/cloudflared.service

# Check status
sudo systemctl status cloudflared

# View logs
sudo journalctl -u cloudflared -f

# Restart tunnel
sudo systemctl restart cloudflared
```

### Token-Based Configuration

The tunnel uses **token-based authentication** (not config file). This means:
- Configuration is managed via Cloudflare Dashboard/API
- Token stored in systemd service file
- No local config.yml needed

---

## Zero Trust Access Configuration

### Protected Applications

Six applications require Zero Trust authentication:

| App Name | Domain | App ID |
|----------|--------|--------|
| OneMind - Portainer | portainer.onemindos.dev | 84f346b2-4cf0-44ad-b8fd-3707438b0b5d |
| OneMind - N8N | n8n.onemindos.dev | 109d04aa-baff-4775-83d9-c07f823c3721 |
| OneMind - Kestra | kestra.onemindos.dev | e0f8f1f0-fb13-4b52-9fac-667b739db2a1 |
| OneMind - LEGACY Brain | agno.onemindos.dev | 29d10ad4-3722-4153-a215-1d04596783aa |
| OneMind - TOC | toc.onemindos.dev | 31dc61e8-c5ea-4382-9128-ca91c3a6405c |
| OneMind - VNC | vnc.onemindos.dev | 5de3b10f-2bae-4788-b352-532e23e2a0d3 |

### Access Policy: "Zeus Only"

Each protected app has an "Allow" policy with:

```json
{
  "name": "Zeus Only",
  "decision": "allow",
  "include": [
    {
      "email": {
        "email": "zeusdelacruz@gmail.com"
      }
    },
    {
      "email": {
        "email": "zeus.vaults@pm.me"
      }
    }
  ]
}
```

### Unprotected Applications

These services are intentionally open (read-only monitoring):

| Service | Domain | Reason |
|---------|--------|--------|
| Uptime Kuma | status.onemindos.dev | Public status page |
| NATS Monitor | nats.onemindos.dev | Read-only metrics |
| Letta | letta.onemindos.dev | API requires auth token |
| Actions API | actions.onemindos.dev | API requires auth |

---

## API Management

### Adding a New Route

Using the Cloudflare API to add a tunnel route:

```bash
# Get current config
curl -X GET "https://api.cloudflare.com/client/v4/accounts/${ACCOUNT_ID}/cfd_tunnel/${TUNNEL_ID}/configurations" \
  -H "X-Auth-Email: zeus.vaults@pm.me" \
  -H "X-Auth-Key: ${GLOBAL_API_KEY}" \
  -H "Content-Type: application/json"

# Update config (PUT entire ingress array)
curl -X PUT "https://api.cloudflare.com/client/v4/accounts/${ACCOUNT_ID}/cfd_tunnel/${TUNNEL_ID}/configurations" \
  -H "X-Auth-Email: zeus.vaults@pm.me" \
  -H "X-Auth-Key: ${GLOBAL_API_KEY}" \
  -H "Content-Type: application/json" \
  -d '{
    "config": {
      "ingress": [
        {"hostname": "newservice.onemindos.dev", "service": "http://localhost:PORT"},
        ... existing routes ...,
        {"service": "http_status:404"}
      ]
    }
  }'
```

### Creating a Zero Trust Access App

```bash
# Create application
curl -X POST "https://api.cloudflare.com/client/v4/accounts/${ACCOUNT_ID}/access/apps" \
  -H "X-Auth-Email: zeus.vaults@pm.me" \
  -H "X-Auth-Key: ${GLOBAL_API_KEY}" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "OneMind - ServiceName",
    "domain": "service.onemindos.dev",
    "type": "self_hosted",
    "session_duration": "24h"
  }'

# Get app ID from response, then create policy
curl -X POST "https://api.cloudflare.com/client/v4/accounts/${ACCOUNT_ID}/access/apps/${APP_ID}/policies" \
  -H "X-Auth-Email: zeus.vaults@pm.me" \
  -H "X-Auth-Key: ${GLOBAL_API_KEY}" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Zeus Only",
    "decision": "allow",
    "include": [
      {"email": {"email": "zeusdelacruz@gmail.com"}},
      {"email": {"email": "zeus.vaults@pm.me"}}
    ]
  }'
```

---

## Authentication Flow

When a user accesses a protected app:

```
1. User navigates to https://n8n.onemindos.dev
                    │
                    ▼
2. Cloudflare checks: Is this a Zero Trust app?
   YES → Redirect to auth page
                    │
                    ▼
3. User sees "Enter your email" page
   Enters: zeusdelacruz@gmail.com
                    │
                    ▼
4. Cloudflare sends one-time code to email
                    │
                    ▼
5. User enters code from email
                    │
                    ▼
6. Cloudflare sets session cookie (24h)
   Redirects to original URL
                    │
                    ▼
7. User accesses n8n.onemindos.dev
   (No more auth prompts for 24h)
```

---

## Troubleshooting

### "Access Denied" Error
- Email not in allow list
- Check: zeusdelacruz@gmail.com or zeus.vaults@pm.me

### "502 Bad Gateway"
- Target container not running
- Check: `docker ps | grep <service>`

### "Connection Timed Out"
- Tunnel not running
- Check: `sudo systemctl status cloudflared`

### Routes Not Working
- DNS not propagated (wait 5 min)
- Check: `dig service.onemindos.dev`

### Auth Loop (keeps asking for email)
- Clear browser cookies for onemindos.dev
- Try incognito window

---

## Security Considerations

1. **No open ports** - All traffic goes through Cloudflare Tunnel (outbound only)
2. **Zero Trust** - Sensitive apps require email verification
3. **Session limits** - 24-hour session duration
4. **Audit logs** - All access logged in Cloudflare dashboard
5. **WAF protection** - Cloudflare blocks common attacks

---

## Dashboard Access

- **Cloudflare Dashboard**: https://dash.cloudflare.com
- **Zero Trust Dashboard**: https://one.dash.cloudflare.com
- **Tunnel Management**: Networks → Tunnels
- **Access Apps**: Access → Applications

Login with: zeus.vaults@pm.me

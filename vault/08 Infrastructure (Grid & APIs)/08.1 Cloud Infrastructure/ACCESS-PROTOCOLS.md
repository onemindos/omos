# OneMind Cloud Access Protocols

*Security Tiers and Authentication Guide - January 2026*

## Overview

OneMind Cloud implements a layered security model with two access pathways (Tailscale/Cloudflare) and three protection tiers (Open, Protected, Restricted). This document defines who can access what, and how.

---

## Security Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                    SECURITY LAYERS                              │
│                                                                 │
│   ┌─────────────────────────────────────────────────────┐      │
│   │  LAYER 1: NETWORK ACCESS                            │      │
│   │                                                     │      │
│   │  ┌─────────────────┐    ┌─────────────────┐        │      │
│   │  │   TAILSCALE     │    │   CLOUDFLARE    │        │      │
│   │  │   (Private)     │    │   (Public)      │        │      │
│   │  │                 │    │                 │        │      │
│   │  │ • Device-based  │    │ • Email-based   │        │      │
│   │  │ • WireGuard VPN │    │ • Zero Trust    │        │      │
│   │  │ • Direct access │    │ • Edge proxy    │        │      │
│   │  └─────────────────┘    └─────────────────┘        │      │
│   └─────────────────────────────────────────────────────┘      │
│                              │                                  │
│                              ▼                                  │
│   ┌─────────────────────────────────────────────────────┐      │
│   │  LAYER 2: APPLICATION AUTH                          │      │
│   │                                                     │      │
│   │  ┌───────────┐ ┌───────────┐ ┌───────────┐         │      │
│   │  │   OPEN    │ │ PROTECTED │ │ RESTRICTED│         │      │
│   │  │           │ │           │ │           │         │      │
│   │  │ No auth   │ │ Zero Trust│ │ API Key + │         │      │
│   │  │ needed    │ │ email     │ │ Zero Trust│         │      │
│   │  └───────────┘ └───────────┘ └───────────┘         │      │
│   └─────────────────────────────────────────────────────┘      │
│                              │                                  │
│                              ▼                                  │
│   ┌─────────────────────────────────────────────────────┐      │
│   │  LAYER 3: SERVICE AUTH                              │      │
│   │                                                     │      │
│   │  Individual services may have additional auth:      │      │
│   │  • N8N: username/password                          │      │
│   │  • Portainer: username/password                    │      │
│   │  • Letta: API token                                │      │
│   │  • etc.                                            │      │
│   └─────────────────────────────────────────────────────┘      │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## Access Tiers

### Tier 1: Open (Public Monitoring)

Services that are intentionally accessible without authentication.

| Service | Cloudflare URL | Purpose | Risk Level |
|---------|----------------|---------|------------|
| Uptime Kuma | status.onemindos.dev | Public status page | Low |
| NATS Monitor | nats.onemindos.dev | Read-only metrics | Low |

**Rationale:**
- Status page should be accessible to check if services are up
- NATS monitor is read-only, no sensitive data exposed
- These help with debugging without needing auth

---

### Tier 2: Protected (Zero Trust Required)

Services requiring Cloudflare Zero Trust email authentication.

| Service | Cloudflare URL | Zero Trust App |
|---------|----------------|----------------|
| Portainer | portainer.onemindos.dev | OneMind - Portainer |
| N8N | n8n.onemindos.dev | OneMind - N8N |
| Kestra | kestra.onemindos.dev | OneMind - Kestra |
| Agno (LEGACY) | agno.onemindos.dev | OneMind - LEGACY Brain |
| TOC | toc.onemindos.dev | OneMind - TOC |
| VNC | vnc.onemindos.dev | OneMind - VNC |

**Authentication Flow:**
1. Navigate to protected URL
2. Cloudflare prompts: "Enter your email"
3. Enter allowed email address
4. Receive one-time code via email
5. Enter code → Access granted for 24 hours

**Allowed Email Addresses:**
- zeusdelacruz@gmail.com
- zeus.vaults@pm.me

---

### Tier 3: Restricted (API Key + Protected)

Services that have additional API authentication beyond Zero Trust.

| Service | Cloudflare URL | Additional Auth |
|---------|----------------|-----------------|
| Letta | letta.onemindos.dev | Bearer API token |
| Actions | actions.onemindos.dev | API key header |

**Note:** These services accept API calls with proper tokens even without browser-based Zero Trust auth.

---

## Access Methods by Device

### From Zeus's Devices (Tailscale Installed)

**Best Option:** Use Tailscale IP directly

```
http://100.87.87.69:PORT
```

| Advantage | Disadvantage |
|-----------|--------------|
| Fastest (~5ms latency) | Requires Tailscale app |
| No auth prompts | Device must be on tailnet |
| Access to all ports | |
| Works offline from cloud | |

**Tailscale URLs:**

| Service | URL |
|---------|-----|
| Portainer | http://100.87.87.69:9000 |
| N8N | http://100.87.87.69:5678 |
| Kestra | http://100.87.87.69:8088 |
| Agno | http://100.87.87.69:8080 |
| Letta | http://100.87.87.69:8283 |
| TOC | http://100.87.87.69:3002 |
| VNC | http://100.87.87.69:6080 |
| Uptime Kuma | http://100.87.87.69:3001 |
| NATS | http://100.87.87.69:8222 |
| Actions | http://100.87.87.69:8001 |

---

### From Any Browser (Public Internet)

**Best Option:** Use Cloudflare URLs

```
https://service.onemindos.dev
```

| Advantage | Disadvantage |
|-----------|--------------|
| Works from any browser | ~50ms added latency |
| Mobile-friendly | May require Zero Trust auth |
| DDoS protection | Only configured routes |
| Automatic HTTPS | |

**Cloudflare URLs:**

| Service | URL | Auth Required |
|---------|-----|---------------|
| Portainer | https://portainer.onemindos.dev | Yes (Zero Trust) |
| N8N | https://n8n.onemindos.dev | Yes (Zero Trust) |
| Kestra | https://kestra.onemindos.dev | Yes (Zero Trust) |
| Agno | https://agno.onemindos.dev | Yes (Zero Trust) |
| Letta | https://letta.onemindos.dev | No |
| TOC | https://toc.onemindos.dev | Yes (Zero Trust) |
| VNC | https://vnc.onemindos.dev | Yes (Zero Trust) |
| Uptime Kuma | https://status.onemindos.dev | No |
| NATS | https://nats.onemindos.dev | No |
| Actions | https://actions.onemindos.dev | No |

---

### From Mobile (iPhone)

**Option A: Tailscale App (Recommended for trusted device)**
1. Install Tailscale from App Store
2. Sign in with Tailscale account
3. Use Tailscale IPs: `http://100.87.87.69:PORT`

**Option B: Cloudflare URLs (For quick access)**
1. Open Safari/Chrome
2. Navigate to `https://service.onemindos.dev`
3. Complete Zero Trust auth if prompted

---

## Tailscale Setup

### Adding a New Device

1. **Install Tailscale**
   - Mac: `brew install tailscale`
   - iOS: App Store → Tailscale
   - Windows: https://tailscale.com/download
   - Linux: https://tailscale.com/download/linux

2. **Sign In**
   ```bash
   tailscale login
   ```
   Sign in with the account that owns the tailnet

3. **Verify Connection**
   ```bash
   tailscale status
   # Should show onemind-cloud-1 in list
   ```

4. **Test Access**
   ```bash
   curl http://100.87.87.69:3001
   # Should return Uptime Kuma HTML
   ```

### Tailnet Details

| Property | Value |
|----------|-------|
| Tailnet | tail717ba0.ts.net |
| Server Hostname | onemind-cloud-1 |
| Server IP | 100.87.87.69 |
| MagicDNS | onemind-cloud-1.tail717ba0.ts.net |

---

## Zero Trust Setup

### Login Process

```
1. Navigate to https://n8n.onemindos.dev
                    │
                    ▼
2. Cloudflare Access Page
   ┌─────────────────────────────────────┐
   │                                     │
   │   Sign in to OneMind - N8N          │
   │                                     │
   │   Enter your email address          │
   │   ┌─────────────────────────────┐   │
   │   │ zeusdelacruz@gmail.com      │   │
   │   └─────────────────────────────┘   │
   │                                     │
   │        [Send me a code]             │
   │                                     │
   └─────────────────────────────────────┘
                    │
                    ▼
3. Check email for one-time code
                    │
                    ▼
4. Enter code
   ┌─────────────────────────────────────┐
   │                                     │
   │   Enter the code sent to            │
   │   zeusdelacruz@gmail.com            │
   │                                     │
   │   ┌─────────────────────────────┐   │
   │   │ 847291                       │   │
   │   └─────────────────────────────┘   │
   │                                     │
   │           [Verify]                  │
   │                                     │
   └─────────────────────────────────────┘
                    │
                    ▼
5. Access granted! Session valid for 24h
```

### Adding a New Protected App

To add Zero Trust to a new service:

1. **Via Cloudflare Dashboard:**
   - Go to https://one.dash.cloudflare.com
   - Access → Applications → Add an Application
   - Choose "Self-hosted"
   - Set domain and add policy

2. **Via API:**
   ```bash
   # Create app
   curl -X POST "https://api.cloudflare.com/client/v4/accounts/${ACCOUNT_ID}/access/apps" \
     -H "X-Auth-Email: zeus.vaults@pm.me" \
     -H "X-Auth-Key: ${API_KEY}" \
     -d '{
       "name": "OneMind - NewService",
       "domain": "newservice.onemindos.dev",
       "type": "self_hosted"
     }'

   # Add policy
   curl -X POST "https://api.cloudflare.com/client/v4/accounts/${ACCOUNT_ID}/access/apps/${APP_ID}/policies" \
     -H "X-Auth-Email: zeus.vaults@pm.me" \
     -H "X-Auth-Key: ${API_KEY}" \
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

## Service-Level Authentication

Some services have their own authentication layer in addition to network/Zero Trust:

### N8N
- **Type:** Username/Password
- **Location:** First-time setup wizard
- **Recovery:** Reset via database or container

### Portainer
- **Type:** Username/Password
- **Location:** First-time setup wizard
- **Recovery:** Delete data volume, recreate

### Kestra
- **Type:** Optional auth via config
- **Current:** No additional auth configured

### Letta
- **Type:** API Bearer Token
- **Header:** `Authorization: Bearer <token>`
- **Get Token:** Create via Letta admin

---

## Security Best Practices

### Do's

- ✅ Use Tailscale for development work
- ✅ Use Cloudflare for mobile/public access
- ✅ Keep email addresses up to date in Zero Trust
- ✅ Use unique, strong passwords for service logins
- ✅ Review access logs periodically

### Don'ts

- ❌ Don't share Tailscale auth key
- ❌ Don't bypass Zero Trust for sensitive apps
- ❌ Don't expose ports directly (no firewall rules)
- ❌ Don't use same password across services
- ❌ Don't leave sessions logged in on shared devices

---

## Troubleshooting Access Issues

### "Access Denied" on Zero Trust

1. Check you're using the correct email
2. Verify email is in allow list
3. Try incognito window (clear cookies)
4. Check Cloudflare Access logs

### Tailscale Connection Issues

```bash
# Check status
tailscale status

# If disconnected
tailscale up

# If device not visible
tailscale ping onemind-cloud-1
```

### Can't Reach Service via Tailscale

```bash
# Test connectivity
ping 100.87.87.69

# Test specific port
nc -zv 100.87.87.69 8080

# Check if service is running
ssh ubuntu@100.87.87.69 "docker ps | grep <service>"
```

### Zero Trust Auth Loop

1. Clear all cookies for onemindos.dev
2. Clear browser cache
3. Try different browser
4. Check session hasn't expired (24h limit)

---

## Audit and Compliance

### Access Logging

- **Cloudflare:** All Zero Trust access logged in dashboard
- **Tailscale:** Connection logs in admin console
- **Services:** Individual service audit logs

### Review Schedule

| Frequency | Action |
|-----------|--------|
| Weekly | Check for failed access attempts |
| Monthly | Review Zero Trust policy members |
| Quarterly | Audit all service accounts |
| Annually | Full security review |

---

## Emergency Access

If both Tailscale and Cloudflare are down:

1. **AWS Lightsail Console**
   - URL: https://lightsail.aws.amazon.com
   - Use browser-based SSH

2. **Direct SSH (if firewall allows)**
   - May need to add firewall rule temporarily
   - `ssh ubuntu@54.165.53.218`

3. **Emergency Contact**
   - AWS Support (if instance issue)
   - Cloudflare Support (if tunnel issue)
   - Tailscale Support (if VPN issue)

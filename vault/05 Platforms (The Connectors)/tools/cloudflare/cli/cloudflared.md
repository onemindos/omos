---
title: "CLI — cloudflared (Cloudflare Tunnel)"
type: reference
status: active
created: 2026-04-10
domain: SO/01
tags: [type/reference, domain/so, topic/cli, topic/cloudflare, topic/network]
ai_summary: "cloudflared CLI quick reference — Cloudflare Tunnel setup for exposing VPS services without opening firewall ports."
---

# cloudflared — Cloudflare Tunnel CLI

> Exposes VPS services to the internet through Cloudflare's network — no open firewall ports needed.
> Install: `brew install cloudflared` (Mac), `curl -L --output cloudflared.deb .../cloudflared-linux-amd64.deb && dpkg -i ...` (Linux)

---

## Auth

```bash
# Authenticate (opens browser for Cloudflare OAuth)
cloudflared tunnel login

# Verify login
cloudflared tunnel list
```

---

## Create & Run a Tunnel

```bash
# Create new tunnel
cloudflared tunnel create onemind-tunnel

# List tunnels
cloudflared tunnel list

# Route DNS to tunnel
cloudflared tunnel route dns onemind-tunnel app.onemindos.com

# Run tunnel (reads config from ~/.cloudflared/config.yml)
cloudflared tunnel run onemind-tunnel

# Run tunnel with inline config
cloudflared tunnel --url http://localhost:3456 run onemind-tunnel
```

---

## Config File

Path: `~/.cloudflared/config.yml` (Mac) or `/etc/cloudflared/config.yml` (VPS)

```yaml
tunnel: <tunnel-id>
credentials-file: /root/.cloudflared/<tunnel-id>.json

ingress:
  - hostname: app.onemindos.com
    service: http://localhost:3456
  - hostname: n8n.onemindos.com
    service: http://localhost:5678
  - hostname: omnidb.onemindos.com
    service: http://localhost:8080
  # catchall
  - service: http_status:404
```

---

## Run as a Service (VPS)

```bash
# Install as systemd service
cloudflared service install

# Start service
systemctl start cloudflared

# Check service status
systemctl status cloudflared

# View logs
journalctl -u cloudflared -f
```

---

## Cleanup

```bash
# Delete a tunnel
cloudflared tunnel delete onemind-tunnel

# Cleanup DNS routes (done manually via Cloudflare dashboard or API)
```

---

## Troubleshoot

```bash
# Check tunnel metrics
cloudflared tunnel info onemind-tunnel

# Test connectivity
cloudflared tunnel --loglevel debug run onemind-tunnel

# Check which tunnels are active
cloudflared tunnel list --output json | jq '.[] | {name, status}'
```


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]

# Hermes + Matrix Client Configuration
# For Zeus's devices — Recommended Matrix client config

## Recommended Clients by Platform

| Platform | Client | Why |
|----------|--------|-----|
| **macOS** | Element (Desktop) | Full feature, E2EE, bridges |
| **iOS** | Element X | Native iOS, fast sync |
| **Android** | Element X or SchildiChat | Material design, ATAK-friendly |
| **Linux** | Element (Desktop) | Cross-platform |
| **Web** | Element Web (element.onemindos.com) | No install, full access |
| **Terminal** | gomuks or weechat-matrix | CLI access |

## Server Config for All Devices

```json
{
  "default_server_config": {
    "m.homeserver": {
      "base_url": "https://matrix.onemindos.com",
      "server_name": "onemindos.com"
    }
  }
}
```

## Recommended Chat Rooms Structure

### 1. Infrastructure Alerts (PUBLIC - join via link)
- #alerts:onemindos.com — Critical system alerts
- #infra-health:onemindos.com — Health checks, metrics

### 2. Operations (PRIVATE - invite only)
- #ops-oncall:onemindos.com — On-call rotation, urgent
- #deployments:onemindos.com — CI/CD, rollouts

### 3. Intelligence (PRIVATE - E2EE)
- #intel-fabric:onemindos.com — ORP, NATS, fusion
- #tak-alerts:onemindos.com — CoT, tactical

### 4. Admin (PRIVATE - admins only)
- #admin:root:onemindos.com — Server admin
- #legacy-admin:onemindos.com — Hermes/Legacy control

## Security
- E2EE enabled on private rooms
- All public rooms readable via permalink
- Sensitive rooms: encrypted + invite-only

## Mobile Push Notifications
- Use Synapse's push gateway
- Or configure Tiiime/Matrixpush for custom push

## Federation
- Public rooms federate with matrix.org
- Private rooms remain on onemindos.com only

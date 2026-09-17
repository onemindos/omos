# OneMind OS - Tailscale Network Configuration

> **Last Updated**: 2026-02-03
> **Status**: Production (Tailscale-only)

## Overview

OneMind OS uses **Tailscale** as the sole networking layer. This provides:
- **Secure mesh VPN** - All traffic encrypted via WireGuard
- **Multi-device access** - Same backend from Mac, iPhone, iPad
- **Low latency** - Direct peer-to-peer connections
- **No public DNS** - Private network, no Cloudflare/DNS complexity

## Network Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                    TAILSCALE MESH NETWORK                       │
│                                                                 │
│  ┌─────────────────┐                                           │
│  │   Your Mac      │                                           │
│  │  100.114.142.55 │────┐                                      │
│  └─────────────────┘    │                                      │
│                         │    WireGuard                         │
│  ┌─────────────────┐    │    Encrypted                         │
│  │   Your iPhone   │────┼──────────────►  ┌─────────────────┐  │
│  │  100.122.203.68 │    │                 │  Azure Server   │  │
│  └─────────────────┘    │                 │  100.102.21.44  │  │
│                         │                 │                 │  │
│  ┌─────────────────┐    │                 │  ┌───────────┐  │  │
│  │   Your iPad     │────┘                 │  │ Backend   │  │  │
│  │   100.x.x.x     │                      │  │ :8080     │  │  │
│  └─────────────────┘                      │  ├───────────┤  │  │
│                                           │  │ Postgres  │  │  │
│                                           │  │ :5432     │  │  │
│                                           │  ├───────────┤  │  │
│                                           │  │ NATS      │  │  │
│                                           │  │ :9222     │  │  │
│                                           │  └───────────┘  │  │
│                                           └─────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
```

## Device IPs

| Device | Tailscale IP | Role |
|--------|-------------|------|
| **Azure Server** | `100.102.21.44` | Backend services |
| **Mac** | `100.114.142.55` | Development |
| **iPhone** | `100.122.203.68` | Mobile app |

## Services

All services run on the Azure server (`100.102.21.44`):

| Service | Port | URL | Purpose |
|---------|------|-----|---------|
| Backend API | 8080 | `http://100.102.21.44:8080` | Agno agents, chat |
| PostgreSQL | 5432 | `100.102.21.44:5432` | Database |
| NATS WebSocket | 9222 | `ws://100.102.21.44:9222` | Real-time events |
| MQTT | 1883 | `100.102.21.44:1883` | IoT/Wearables |
| Flutter Web | 3000 | `http://100.102.21.44:3000` | Web UI |

## Flutter Configuration

All URLs are configured in a single file:

**File**: `frontend/lib/platform/config/environment.dart`

```dart
class Environment {
  static const String _tailscaleIp = '100.102.21.44';

  static String get apiUrl => 'http://$_tailscaleIp:8080';
  static String get natsUrl => 'ws://$_tailscaleIp:9222';
  static String get mqttBroker => _tailscaleIp;
  static String get gatewayUrl => 'http://$_tailscaleIp:8000';
}
```

## Multi-Device Streaming

Both Flutter web (Mac) and iOS app connect to the **same backend**:

```
┌──────────────┐     ┌──────────────────────┐
│  Mac Browser │────►│                      │
│  Flutter Web │     │   Backend            │
└──────────────┘     │   100.102.21.44:8080 │
                     │                      │
┌──────────────┐     │   - Same sessions    │
│  iPhone App  │────►│   - Same agents      │
│  Flutter iOS │     │   - Same database    │
└──────────────┘     │   - Shared memory    │
                     └──────────────────────┘
```

**Benefits of shared backend**:
- **Conversation continuity** - Start chat on Mac, continue on iPhone
- **Unified agent memory** - Agents remember you across devices
- **Real-time sync** - Changes reflect immediately on all devices
- **Better streaming** - Direct connection, no proxy hops

## Testing Connectivity

### Quick Health Check

```bash
# From any device on Tailscale
curl http://100.102.21.44:8080/health
# Expected: {"status":"ok","instantiated_at":"..."}
```

### From Mac Terminal

```bash
# Check Tailscale status
tailscale status

# Ping server
ping 100.102.21.44

# Test backend
curl http://100.102.21.44:8080/agents | jq '.[0].name'
```

### From iPhone

1. Open Safari
2. Go to: `http://100.102.21.44:8080/health`
3. Should see JSON response

## Troubleshooting

### "Connection refused"

1. Check Tailscale is connected:
   ```bash
   tailscale status
   ```

2. Check backend is running:
   ```bash
   # Via Azure CLI
   az vm run-command invoke --resource-group ONEMIND-PRODUCTION \
     --name onemind-vm --command-id RunShellScript \
     --scripts "docker ps | grep agno"
   ```

### "Backend not responding"

Restart the container:
```bash
az vm run-command invoke --resource-group ONEMIND-PRODUCTION \
  --name onemind-vm --command-id RunShellScript \
  --scripts "docker restart onemind-agno"
```

### iPhone can't connect

1. Ensure Tailscale app is installed and connected (VPN icon visible)
2. Try toggling Tailscale off/on
3. Check if server is reachable: `ping 100.102.21.44` in Terminal on Mac

## Migration from Cloudflare

**Previous setup** (deprecated):
- Cloudflare DNS: `*.onemindos.dev`
- Cloudflare Tunnel to backend
- Complex routing, occasional disconnects

**Current setup** (Tailscale-only):
- No public DNS records
- Direct mesh VPN connections
- Lower latency, more reliable
- 16 DNS records deleted from Cloudflare

## Security

Tailscale provides:
- **WireGuard encryption** - Military-grade encryption
- **Zero trust** - Only your devices can access
- **No public exposure** - Server not reachable from internet
- **MagicDNS** - Optional friendly hostnames

The backend listens on `0.0.0.0:8080` but is only accessible via Tailscale IPs (100.x.x.x range).

## Quick Reference

| Action | Command/URL |
|--------|-------------|
| Check server health | `curl http://100.102.21.44:8080/health` |
| List agents | `curl http://100.102.21.44:8080/agents` |
| Tailscale status | `tailscale status` |
| Restart backend | Via Azure CLI (see above) |
| Run Flutter web | `cd frontend && flutter run -d chrome` |
| Run Flutter iOS | `cd frontend && flutter run` |

---

**Maintained by**: Zeus
**Related docs**:
- [INFRASTRUCTURE.md](./INFRASTRUCTURE.md) - Azure VM details
- [CLAUDE.md](../../../CLAUDE.md) - Development guide

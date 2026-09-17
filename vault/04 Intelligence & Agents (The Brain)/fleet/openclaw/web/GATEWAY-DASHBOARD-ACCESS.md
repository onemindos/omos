---
id: gateway-dashboard-access
title: "Gateway Dashboard Access"
type: reference
status: active
created: 2026-05-03
modified: 2026-05-03
owner: zeus
tags: [type/reference, domain/infrastructure, tool/openclaw]
---
# Gateway Dashboard Access

How to access OpenClaw gateway dashboards across the fleet.

---

## Dashboard URLs & Tokens

| Node | Dashboard URL | Auth Token |
|------|--------------|------------|
| **omos-core** | `https://omos-core.tail717ba0.ts.net` | `Vd3kHXi7HSAgWPtEGBiQC42UEy98ZyGH` |
| **omos-sentinel** | `https://omos-sentinel.tail717ba0.ts.net` | *(check sentinel openclaw.json)* |

Both dashboards are Tailscale-only — not reachable from the public internet.

---

## Device Pairing Fix (Applied 2026-05-03)

OpenClaw v2026.5.2 introduced a device pairing handshake that requires each browser to be individually approved before accessing the dashboard. In a single-operator setup this creates a deadlock: the CLI needs `operator.admin` scope to approve devices, but device-pair blocks the CLI from getting that scope.

### The Fix

Set `gateway.controlUi.dangerouslyDisableDeviceAuth: true` in `openclaw.json`:

```json
{
  "gateway": {
    "controlUi": {
      "allowedOrigins": ["https://omos-core.tail717ba0.ts.net"],
      "dangerouslyDisableDeviceAuth": true
    }
  }
}
```

**Applied to both nodes:**
- omos-core (`/home/ubuntu/.openclaw/openclaw.json`)
- omos-sentinel (`/home/ubuntu/.openclaw/openclaw.json`)

### What This Does

- Skips per-browser device fingerprint + approval handshake on the Control UI
- Does NOT remove token auth — the gateway token is still required to connect
- Does NOT affect Tailscale network isolation — dashboard is still private-network only
- The name sounds dangerous but it's appropriate for single-operator setups

### Security Layers Still Active

1. **Tailscale** — dashboard only reachable on private tailnet
2. **Token auth** (`gateway.auth.mode: "token"`) — requires token to connect
3. **allowedOrigins** — restricts which origins can load the Control UI

### When You'd Want Device Pairing Back

If you add other operators who should have limited access to specific gateways, re-enable by setting `dangerouslyDisableDeviceAuth: false` (or removing the key). Each new browser will then require explicit approval via `openclaw devices approve`.

---

## Troubleshooting

### "origin not allowed"
Add your dashboard URL to `gateway.controlUi.allowedOrigins`:
```json
"controlUi": {
  "allowedOrigins": ["https://omos-core.tail717ba0.ts.net"]
}
```

### "unauthorized: gateway token missing"
Paste the gateway token into the Control UI settings when prompted.

### "device pairing required" / "role upgrade pending approval"
If `dangerouslyDisableDeviceAuth` is not set, you hit the device pairing loop. Fix:
1. Set `dangerouslyDisableDeviceAuth: true` in `gateway.controlUi`
2. Clear pending devices: `echo '{}' > ~/.openclaw/devices/pending.json`
3. Restart gateway: `systemctl --user restart openclaw-gateway.service`

---

## Mac Client Config

The Mac connects to Core as a remote client. Config at `~/.openclaw/openclaw.json`:

```json
{
  "gateway": {
    "mode": "remote",
    "remote": {
      "url": "wss://omos-core.tail717ba0.ts.net",
      "token": "Vd3kHXi7HSAgWPtEGBiQC42UEy98ZyGH",
      "transport": "direct"
    }
  }
}
```

**Critical:** `remote.token` must match Core's `gateway.auth.token` exactly, or you'll get pairing/auth errors.


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]

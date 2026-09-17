# OneMind OS — LiveKit Setup

> Real-time audio/video · Primary: LiveKit Cloud · Self-hosted: fallback only · **NOT currently running**

---

## Status

LiveKit is defined in the Docker Compose stack (profile: `full`) but **not actively running**. The primary setup uses **LiveKit Cloud** (managed). The self-hosted container is a fallback option.

---

## Architecture Role

LiveKit provides real-time audio/video (WebRTC) infrastructure for:

- **Voice mode** — Legacy's live voice conversations with Zeus
- **VisionClaw** — camera feed processing and real-time vision
- **Voice avatar** — future embodied avatar rendering

```
iOS App / Mac ──WebRTC──▶ LiveKit Cloud ──events──▶ Legacy (OpenClaw)
                          (or self-hosted)
```

---

## Primary: LiveKit Cloud

https://livekit.io — free tier: 50 hrs/month real-time usage

Required env vars (in `/opt/onemind/.env.sa` or via 1Password):

| Var | Description |
|-----|-------------|
| `LIVEKIT_URL` | Cloud URL, e.g. `wss://yourproject.livekit.cloud` |
| `LIVEKIT_API_KEY` | API key from LiveKit Cloud dashboard |
| `LIVEKIT_API_SECRET` | API secret from LiveKit Cloud dashboard |

Legacy uses these via the `LIVEKIT_API_KEY` and `LIVEKIT_API_SECRET` env vars already present in `.env.tpl`.

---

## Fallback: Self-Hosted Docker

```yaml
# From /opt/onemind/docker-compose.yml (profile: full)
onemind-livekit:
  image: livekit/livekit-server:latest
  container_name: onemind-livekit
  profiles: [full]
  restart: unless-stopped
  ports:
    - "7880:7880"       # HTTP API
    - "7881:7881"       # WebRTC TCP
    - "7882:7882/udp"   # WebRTC UDP
  environment:
    LIVEKIT_KEYS: "${LIVEKIT_API_KEY:-devkey}:${LIVEKIT_API_SECRET:-devsecret}"
```

---

## Ports (self-hosted only)

| Port | Protocol | Use |
|------|----------|-----|
| 7880 | TCP | HTTP API |
| 7881 | TCP | WebRTC (TCP fallback) |
| 7882 | UDP | WebRTC (primary) |

---

## Start / Manage (self-hosted fallback)

```bash
# Start LiveKit only
cd /opt/onemind
docker compose --profile full up -d livekit

# Stop LiveKit
docker compose stop livekit

# View logs
docker logs -f onemind-livekit
```

---

## iOS / Swift Integration

LiveKit SDK: `LiveKit-Swift` (via SPM)

Config in iOS app uses:
- `LIVEKIT_URL` → connect URL
- `LIVEKIT_API_KEY` + `LIVEKIT_API_SECRET` → token generation (server-side)

---

## Decision: Cloud vs Self-Hosted

Use **LiveKit Cloud** unless:
- Cloud free tier (50 hrs/mo) is exhausted, or
- You need on-premise data isolation

To switch to self-hosted: start the Docker container, update `LIVEKIT_URL` to `ws://104.248.4.212:7880`.

---

## Related Files

- `infra/services/livekit/LIVEKIT.md` — general LiveKit overview
- `infra/services/livekit/LIVEKIT-NATS-BOUNDARIES.md` — when to use LiveKit vs NATS
- `infra/services/livekit/LIVEKIT-SELF-HOSTED-REALITY.md` — self-hosting tradeoffs

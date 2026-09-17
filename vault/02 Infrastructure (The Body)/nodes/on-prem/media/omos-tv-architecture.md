# OMOS TV — Architecture, Development & Network Reference
**Created:** 2026-04-18
**Author:** Legacy ⚡
**Domain:** SO 02 Infrastructure + GE 75 Empire (content delivery)

---

## 1. Device Identity

| Field | Value |
|-------|-------|
| Name | OMOS TV |
| Model | Hisense HiSmart 4K ATV4 |
| OS | Android TV |
| Cast | Chromecast built-in (v1.68.410410) |
| Tailscale IP | `100.117.247.62` |
| Local LAN IP | `192.168.1.172` |
| Tailscale DNS | `omos-tv.tail717ba0.ts.net` |
| Cast Ports | 8008 (HTTP), 8009 (Cast TLS) |
| ADB | Disabled (ports closed) |

---

## 2. System Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    ONE MIND OS NETWORK                      │
│                                                             │
│  ┌──────────────────────┐    ┌──────────────────────────┐  │
│  │  onemindos-cloud-1   │    │   Zeus's MacBook Pro     │  │
│  │  VPS · 64.227.8.23   │    │   192.168.1.171          │  │
│  │  Tailscale: .25.10   │◄──►│   Tailscale: .114.55     │  │
│  │                      │    │                          │  │
│  │  omos-tv-receiver ✅  │    │  catt (cast bridge)      │  │
│  │  :8889 dashboard     │    │  pychromecast            │  │
│  │  :8888 media server  │    │                          │  │
│  └──────────────────────┘    └──────────┬───────────────┘  │
│                                         │ same LAN         │
│                              ┌──────────▼───────────────┐  │
│                              │       OMOS TV            │  │
│                              │  192.168.1.172           │  │
│                              │  Tailscale: .117.62      │  │
│                              │  Hisense HiSmart 4K ATV4 │  │
│                              └──────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

### Why VPS Can't Cast Directly (Today)
Google Cast uses **mDNS** (multicast DNS) for device discovery — broadcast packets that only travel within the same LAN segment. The VPS is in a DigitalOcean datacenter; it cannot hear mDNS from Zeus's home network. The Tailscale IP is reachable but the Cast TLS handshake requires a local-network connection to initiate.

**Workaround (current):** Mac acts as the cast bridge — it's on the same LAN as the TV and can initiate the Cast session. VPS serves media and state, Mac triggers the cast.

**Permanent fix (future):** Tailscale subnet router — see Section 5.

---

## 3. Current Cast Stack

### VPS Components

| Component | Path | Purpose |
|-----------|------|---------|
| `cast.py` | `skills/omos-tv/scripts/cast.py` | Pychromecast controller |
| `push_state.py` | `skills/omos-tv/scripts/push_state.py` | Push display state to TV |
| `receiver_server.py` | `skills/omos-tv/scripts/receiver_server.py` | Serves dashboard + state API |
| `receiver/index.html` | `skills/omos-tv/receiver/index.html` | One Mind OS TV app |
| `omos-tv-receiver.service` | `/etc/systemd/system/` | Keeps receiver server alive |

### Mac Components

| Component | Location | Purpose |
|-----------|----------|---------|
| `catt` | `~/Library/Python/3.9/bin/catt` | Cast CLI (mDNS discovery) |
| `pychromecast` | pip3 | Python cast library |

### Service Endpoints

| Endpoint | URL | Purpose |
|----------|-----|---------|
| Dashboard | `http://100.98.25.10:8889/dashboard` | TV receiver app |
| State API GET | `http://100.98.25.10:8889/api/state` | Read current display state |
| State API POST | `http://100.98.25.10:8889/api/state` | Push new state (agents use this) |
| Health | `http://100.98.25.10:8889/health` | Liveness check |
| Media Server | `http://100.98.25.10:8888/` | Serve media files to TV |

---

## 4. Agent Usage — How to Control the TV

### Push a message to the TV
```bash
python3 /root/.openclaw/workspace/skills/omos-tv/scripts/push_state.py \
  message "Today's focus: One Mind OS launch" --agent "Legacy"
```

### Push a stats dashboard
```bash
python3 /root/.openclaw/workspace/skills/omos-tv/scripts/push_state.py stats \
  '[{"label":"Tasks Done","value":"12","unit":"today"},
    {"label":"Karma","value":"4821","unit":"pts"},
    {"label":"Streak","value":"7","unit":"days"},
    {"label":"Focus","value":"GE 75","unit":"domain"}]'
```

### Push a presentation slide
```bash
python3 /root/.openclaw/workspace/skills/omos-tv/scripts/push_state.py \
  presentation "The One Mind Method" \
  "CODEX: Capture → Organize → Direct → Execute → Multiply" \
  --slide "1 / 5"
```

### Cast a YouTube video (requires Mac online)
```bash
/Users/zeuslegacy/Library/Python/3.9/bin/catt -d "OMOS TV" cast "https://youtube.com/watch?v=VIDEO_ID"
```

### Cast a rendered video from VPS media server (requires Mac online)
```bash
/Users/zeuslegacy/Library/Python/3.9/bin/catt -d "OMOS TV" \
  cast "http://100.98.25.10:8888/videos/rendered/intro.mp4"
```

### Clear TV / reset to default
```bash
python3 /root/.openclaw/workspace/skills/omos-tv/scripts/push_state.py clear
```

---

## 5. Future Architecture — Tailscale Subnet Router

### What It Is
A **subnet router** is one machine on a local network that advertises its entire LAN to Tailscale. Once configured, every device on the VPS (onemindos-cloud-1) can reach every device on that LAN — no Tailscale app needed on each device.

### Why This Matters
- VPS can cast directly to TV — no Mac dependency
- Legacy can reach printers, smart displays, NAS, any device
- Works even when Mac is off
- One router per physical location = full agent control of that space

### Setup (when ready)
On the machine that will be the subnet router (e.g., office Mac, Pi, NUC):
```bash
# Advertise the local subnet
tailscale up --advertise-routes=192.168.1.0/24

# On Tailscale admin console: approve the route
# https://login.tailscale.com/admin/machines

# On VPS: enable route acceptance
tailscale up --accept-routes
```

After this: `ping 192.168.1.172` works from VPS directly.

### Office Network Vision
```
NYC Office LAN (192.168.1.0/24)
├── Office Mac / NUC  ← Tailscale subnet router
├── OMOS TV
├── Secondary displays
├── NAS / storage
├── Smart home devices
└── Any new device → Legacy sees it automatically

        ↕ Tailscale (subnet routed)

VPS onemindos-cloud-1
  └── Legacy casts, controls, monitors entire office LAN
```

---

## 6. What is mDNS?

**mDNS** (Multicast DNS, RFC 6762) is a zero-configuration protocol where devices announce themselves on a local network without a central DNS server.

- Your phone opens YouTube and sees "OMOS TV" available to cast — that's mDNS
- Devices broadcast to `224.0.0.251:5353` every few seconds: "I'm OMOS TV at 192.168.1.172"
- Only works within the same LAN segment (doesn't cross routers/internet)
- Also used by: AirPlay, AirDrop, Bonjour, network printers, HomeKit

**Why the VPS can't use it:** mDNS packets don't travel across the internet. The VPS is on a different network. A subnet router solves this by relaying mDNS through Tailscale.

---

## 7. Tailscale: Exit Node vs Subnet Router

These are two different and commonly confused features:

### Subnet Router
**Routes traffic TO a local network** — agents/devices outside can reach devices inside.
- Use case: VPS reaches TV at 192.168.1.172
- Direction: Outside → Inside LAN
- LAN devices don't need Tailscale installed

### Exit Node
**Routes ALL your internet traffic THROUGH a machine** — like a VPN.
- Use case: Zeus's iPhone browses the internet through the VPS (appears as NYC IP)
- Direction: Your device's internet traffic → Exit Node → Internet
- Hides your real IP, bypasses geo-restrictions
- NOT useful for reaching local devices

### Summary Table

| Feature | Subnet Router | Exit Node |
|---------|--------------|-----------|
| Purpose | Reach LAN devices | Route internet traffic |
| Direction | Outside → LAN | Device → Internet |
| Use case | Agent controls TV | Browse as VPS IP |
| Needs app on LAN devices? | No | No |
| What you want for office | ✅ Yes | 🟡 Optional |

---

## 8. Device Discovery — What Legacy Can See

### Today (Tailscale only)
Legacy can only see devices with the Tailscale app installed:
- onemindos-cloud-1 (VPS) ✅
- zeus-macbook-pro ✅
- zeus-iphone-17 ✅
- omos-tv ✅ (has Tailscale)

### With Subnet Router (future)
Legacy can reach **any device on the LAN** — no app needed:
- Smart TVs (all of them)
- Printers
- NAS devices
- Raspberry Pis
- Smart home hubs
- Any device that responds to TCP/IP

### Finding devices on a LAN
```bash
# From subnet-routed VPS:
nmap -sn 192.168.1.0/24          # Ping scan — find all live hosts
nmap -p 8008,8009 192.168.1.0/24  # Find all Cast devices
arp -a                             # Show ARP table (known devices)
```

---

## 9. TV Dashboard — Display Modes

The One Mind OS TV app (`receiver/index.html`) supports three modes:

### Message Mode (default)
Full-screen text message from an agent. Good for:
- Notifications from Legacy
- Focus reminders
- Quick announcements

### Stats Mode
4-column grid of metric cards. Good for:
- Daily productivity stats
- System health metrics
- Todoist task counts
- Business KPIs

### Presentation Mode
Large title + body text. Good for:
- Student lectures with Legacy
- One Mind OS framework presentations
- Live teaching sessions
- Slide-by-slide content

---

## 10. Roadmap / Future Development

- [ ] **Custom Cast Receiver App** — register with Google Cast SDK, run natively on TV as a Cast app (no polling needed — push via Cast messaging)
- [ ] **Tailscale subnet router** — VPS casts directly, Mac not required
- [ ] **ADB access** — enable ADB on TV for deeper Android control (install apps, screenshots, input events)
- [ ] **Media server as systemd service** — persistent, auto-restart
- [ ] **HyperFrames → TV pipeline** — Empire renders video → auto-cast to TV
- [ ] **Voice trigger** — "Legacy, cast to TV" from iPhone triggers cast
- [ ] **Multi-room** — second TV in office joins cast group
- [ ] **Live stats feed** — real-time Todoist/system data pushed to TV every 60s
- [ ] **Student mode** — structured presentation flow with Legacy narrating via voice

---

*Last updated: 2026-04-18 | Legacy ⚡ | SO 02 Infrastructure*


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]

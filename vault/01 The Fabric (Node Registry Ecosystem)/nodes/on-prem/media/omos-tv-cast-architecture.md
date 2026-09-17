# OMOS TV — Cast Architecture & Development Reference
**Domain:** SO 01 Intelligence  
**Status:** Phase 1 complete — Cast proven working  
**Last Updated:** 2026-04-18  

---

## Device Profile

| Field | Value |
|-------|-------|
| **Name** | OMOS TV |
| **Model** | Hisense HiSmart 4K ATV4 |
| **OS** | Android TV |
| **Cast Version** | 1.68.410410 |
| **Local IP** | 192.168.1.172 |
| **Tailscale IP** | 100.117.247.62 |
| **Cast Ports** | 8008 (HTTP), 8009 (Cast TLS) |
| **ADB** | Not enabled (ports 5554/5555 closed) |
| **Tailscale Node** | omos-tv.tail717ba0.ts.net |

---

## Network Architecture

### Current (Phase 1)
```
Legacy (VPS 64.227.8.23)
        ↓ sessions_send / exec host=node
Mac (zeus-macbook-pro 192.168.1.171 / Tailscale 100.114.142.55)
        ↓ pychromecast / catt (mDNS on local LAN)
OMOS TV (192.168.1.172)
```

**Why VPS can't cast directly:**  
Google Cast uses mDNS (multicast DNS) for device discovery — broadcast packets that only travel within the same LAN. The VPS is on a DigitalOcean network in NYC; it cannot hear local broadcast traffic. Tailscale provides IP routing but not LAN broadcast relay.

### Future (Phase 2 — Office/Home Node)
```
Legacy (VPS)
        ↓ Tailscale (direct LAN access via subnet router)
Office Node (Raspberry Pi / Mac Mini / any machine)
        ↓ subnet router advertises 192.168.x.0/24
All LAN devices (TVs, displays, printers, smart devices)
```

With a Tailscale subnet router running on one machine per physical space:
- VPS gets direct access to every device on that LAN
- No Mac dependency — works even if Mac is off
- Add any device to the office — Legacy sees it automatically
- One node per location = full agent control of that space

---

## Protocol Stack

### mDNS — What It Is
Devices on a LAN broadcast "I'm here!" announcements every few seconds using multicast DNS (port 5353). This is how your phone finds the TV when you open YouTube — both devices hear the same LAN broadcasts. The VPS is too far away to hear them. A subnet router relays these through Tailscale.

### Pychromecast — Full Protocol Access
`pychromecast` is the Python library that implements the full Google Cast v2 protocol. It gives agents:

- Connect to any Cast device by IP (bypasses mDNS when IP is known)
- Launch any app on the TV (YouTube, Netflix, custom receiver)
- Full media queue control (add, skip, loop, seek)
- Real-time playback status callbacks
- **Custom Cast receiver apps** — webpages that run natively ON the TV
- Multi-room audio groups
- Device reboot / factory reset
- Volume, mute, display name, locale, timezone

### catt — Convenience Wrapper
`catt` wraps pychromecast for simple CLI use. Good for quick casts but limited control. Use pychromecast directly for agent automation.

---

## Custom Cast Receiver App — The Power Move

A Custom Cast Receiver is a webpage hosted on a URL that runs natively on the TV when cast to it. The TV becomes a display agent.

### What It Enables
- **Live presentations** — Legacy casts a lesson to the TV; students see it
- **One Mind Dashboard** — real-time stats, task counts, system health on the big screen
- **Video playback** — HyperFrames-rendered MP4s cast directly
- **Agent-pushed notifications** — Legacy sends a message, it appears on screen
- **Slideshow/content** — Empire pushes brand content to the TV

### Architecture
```
VPS serves receiver webpage
        ↓ (static HTML/JS at http://100.98.25.10:PORT)
Mac triggers cast to OMOS TV
        ↓ (tells TV: load this URL as receiver)
TV loads and renders the page
        ↓
Legacy pushes data via WebSocket to the receiver
TV updates in real-time
```

### Receiver App Plan
- Host at: `/root/onemind-media/receiver/` (served via VPS HTTP)
- Stack: HTML + CSS + WebSocket client
- Modes:
  - `dashboard` — system stats, Todoist tasks, weather
  - `video` — fullscreen MP4 player
  - `presentation` — slide-style content for students/demos
  - `notification` — overlay alert from Legacy

---

## Media Pipeline

### VPS Media Server
- Directory: `/root/onemind-media/`
- Serve command: `python3 -m http.server 8888 --directory /root/onemind-media`
- Accessible at: `http://100.98.25.10:8888`
- Systemd service: **TODO — create persistent media server service**

### End-to-End Flow
```
1. Empire renders video via HyperFrames
   → saves to /root/onemind-media/videos/rendered/

2. Dropbox sync (rclone-dropbox.service)
   → Mac: ~/Library/CloudStorage/Dropbox.../onemind-media/

3. VPS HTTP server serves the file
   → http://100.98.25.10:8888/videos/rendered/output.mp4

4. Legacy sends command to Mac node
   → exec host=node: python3 cast.py "http://100.98.25.10:8888/videos/rendered/output.mp4"

5. Mac pychromecast casts URL to OMOS TV
   → TV plays the video
```

---

## Tailscale Network Map

| Device | Tailscale IP | Local IP | Status |
|--------|-------------|----------|--------|
| onemindos-cloud-1 (VPS) | 100.98.25.10 | 64.227.8.23 | ✅ Always on |
| zeus-macbook-pro | 100.114.142.55 | 192.168.1.171 | ✅ Active |
| zeus-iphone-17 | 100.122.203.68 | — | ✅ Online |
| OMOS TV | 100.117.247.62 | 192.168.1.172 | ✅ When on |

### Tailscale Concepts

**Subnet Router**  
A machine that runs `tailscale up --advertise-routes=192.168.1.0/24`. It tells the Tailscale network "I can reach everything on this subnet." Once set up, the VPS can reach any device on that LAN directly — no app needed on each device.

**Exit Node**  
A machine that routes ALL internet traffic for other Tailscale nodes through itself. If the VPS is set as an exit node, your phone's internet traffic goes VPS → internet. Useful for: privacy, bypassing geo-restrictions, forcing all traffic through a known point. Different from subnet router — exit node is about internet routing, subnet router is about LAN access.

**MagicDNS**  
Tailscale's DNS layer. Instead of `100.98.25.10` you can use `onemindos-cloud-1` as the hostname. Automatically resolves across all your devices.

**Key Expiry**  
OMOS TV key expires: 2026-10-14. Renew before then or TV drops off the network.

---

## Development Roadmap

### Phase 1 — ✅ Complete
- [x] Device identified: Hisense HiSmart 4K ATV4
- [x] Cast proven working (YouTube via catt from Mac)
- [x] VPS media server tested
- [x] Tailscale connectivity confirmed
- [x] Architecture documented

### Phase 2 — Cast Skill
- [ ] Create `omos-tv` skill in Legacy workspace
- [ ] Pychromecast script: `cast_to_tv.py` on Mac node
- [ ] Commands: `cast <url>`, `cast_video <file>`, `volume <n>`, `stop`, `status`
- [ ] Legacy can trigger with: "cast this to OMOS TV"
- [ ] Create persistent VPS media server (systemd service)

### Phase 3 — Custom Receiver App
- [ ] Build One Mind receiver webpage (HTML/CSS/JS/WebSocket)
- [ ] Modes: dashboard, video, presentation, notification
- [ ] Host on VPS, cast via Mac
- [ ] Legacy pushes real-time data to TV display
- [ ] Live student presentation mode

### Phase 4 — Office Node (when ready)
- [ ] Set up Tailscale subnet router on office machine
- [ ] VPS gets direct LAN access — Mac no longer required as bridge
- [ ] Add all office devices to Legacy's awareness
- [ ] Expand to home, other locations

---

## Files & Paths

| Path | Purpose |
|------|---------|
| `/root/onemind-media/videos/rendered/` | HyperFrames output |
| `/root/onemind-media/receiver/` | Cast receiver app (TODO) |
| `~/.openclaw/workspace/skills/omos-tv/` | Cast skill (TODO) |
| `/Users/zeuslegacy/Library/Python/3.9/bin/catt` | catt CLI on Mac |
| `100.98.25.10:8888` | VPS media server |

---

## Notes
- Cast uses TLS on port 8009 — pychromecast handles this automatically
- TV goes idle and stops mDNS after ~10min; wake it first or use direct IP
- hotspot BSSID: FA:8F:38:8D:2F:C5 (Hisense)
- Tailscale relay: nyc (same region as VPS — low latency)
- catt is a wrapper over pychromecast — for the skill, use pychromecast directly


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]

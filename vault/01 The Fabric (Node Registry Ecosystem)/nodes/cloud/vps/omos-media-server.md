# OMOS Media Server

**Status:** ✅ Live (2026-04-24)
**Host:** onemindos-cloud-1 (VPS, 64.227.8.23 / Tailscale 100.98.25.10)
**Port:** 8888
**URL:** `http://100.98.25.10:8888/`
**Root dir:** `/root/media/`
**Service:** `omos-media-server.service` (systemd, auto-start on boot)

---

## Purpose

Serves local media files (MP4, images, audio) over HTTP so the OMOS TV and other nodes can fetch and play them without external hosting. Primary use case: cast videos to OMOS TV via Pychromecast.

---

## Service Config

**File:** `/etc/systemd/system/omos-media-server.service`

```ini
[Unit]
Description=One Mind OS — Media HTTP Server
After=network.target
Wants=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/root/media
ExecStart=/usr/bin/python3 -m http.server 8888 --bind 0.0.0.0 --directory /root/media
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

**Commands:**
```bash
systemctl status omos-media-server    # Check status
systemctl restart omos-media-server   # Restart
systemctl stop omos-media-server      # Stop
journalctl -u omos-media-server -f    # Logs
```

---

## Directory Structure

```
/root/media/
├── testing-tv/          ← Videos from Dropbox "testing tv" folder
│   ├── composio-cli-promo.mp4
│   ├── onemind-intro-v2.mp4
│   ├── openclaw-dreaming.mp4
│   └── tailscale-explainer.mp4
```

---

## Dropbox Sync

Videos are pulled from Dropbox via rclone (configured on VPS, remote: `dropbox`).

```bash
# Sync a Dropbox folder to media server
rclone copy "dropbox:folder name" /root/media/dest-folder/

# List Dropbox remotes
rclone lsd dropbox:
```

---

## Casting to OMOS TV

Once a file is in `/root/media/`, cast it directly:

```bash
cd /root/.openclaw/workspace/skills/omos-tv/scripts
python3 cast.py cast_media http://100.98.25.10:8888/testing-tv/onemind-intro-v2.mp4
python3 cast.py kill    # Stop cast, clear TV (works even when no active session)
```

---

## Notes

- Media server binds to `0.0.0.0:8888` — accessible from Tailscale network and LAN
- No auth — internal network only, not exposed to public internet
- Port 8888 is separate from the OMOS TV receiver (port 8889)
- Add files to `/root/media/` and they're immediately servable — no restart needed

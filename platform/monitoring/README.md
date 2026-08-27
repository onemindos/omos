# Murmur (Mumble Voice Server)

**Live on:** DO Droplet `146.190.143.195:64738` (TCP + UDP)  
**Role:** Voice I/O layer for the Hermes agent pipeline — Murmur → omos-voice-bridge → NATS → Legacy  
**Supervised by:** Docker on the watchdog droplet (alongside Uptime Kuma)

## Config
`murmur.ini` is the canonical server config. On the droplet it lives at `/data/murmur.ini` (mounted into the container).

## Deploy (fresh droplet)
```bash
docker run -d \
  --name murmur \
  --restart unless-stopped \
  -p 64738:64738/tcp \
  -p 64738:64738/udp \
  -v /data/murmur:/data \
  mumblevoip/mumble-server:latest
```

## Notes
- SQLite DB lives on the droplet at `/data/mumble-server.sqlite` — back up before any droplet rebuild
- SuperUser password stored in 1Password: `op://00-24 SO/murmur-superuser/password`
- iTAK Mumble plugin should point at `146.190.143.195:64738`
- omos-voice-bridge needs `MUMBLE_HOST=146.190.143.195`

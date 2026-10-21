# voice/ — Murmur (Mumble) tactical voice

Sovereign push-to-talk voice for the fleet, alongside TAK. Runs as the `murmur` container
(image `mumblevoip/mumble-server`), port **64738 tcp+udp**.

## Layout
```
voice/
├── config/murmur.ini   # committed reference config (welcome text, bandwidth, users)
└── data/               # runtime: mumble-server.sqlite, generated config, certs (GITIGNORED)
```

## Important notes
- The official image is configured by **env vars** (`MUMBLE_SUPERUSER_PASSWORD`, `MUMBLE_CONFIG_*`),
  set in `deploy/.env`. The prior container used the wrong `MURMUR_` prefix and a random
  SuperUser password was generated at first boot — the migrate step fixes this.
- `config/murmur.ini` is a **reference** — the running server uses the env-generated
  `data/mumble_server_config.ini`. Keep them consistent or drive everything from env.
- `data/` holds the user database — **not committed**. Back it up separately.

## First-time move into this folder
```bash
deploy/scripts/migrate-murmur.sh    # moves ~/murmur/data here, recreates container
```

## Client connection
```
Server:   192.168.1.183 (LAN)  |  100.74.182.117 (Tailscale)
Port:     64738
```
Register each operator as a Mumble user; don't run day-to-day as SuperUser.
Roster check: `deploy/scripts/murmur-roster.sh`.

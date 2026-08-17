# deploy/ — Tactical stack control plane

Single source of truth for running the local tactical stack: **TAK Server + TAK DB + Murmur voice**.
NATS runs natively on the host (launchd `com.nats.server`, `:4222`); the TAK NATS bridge plugin
connects out to it. This folder holds structure + runbooks only — **no secrets, no certs** (gitignored).

## Bring up / manage

```bash
cd deploy
cp .env.example .env          # then set TAK_DB_PASSWORD + MUMBLE_SUPERUSER_PASSWORD
docker compose --env-file .env up -d
```

| Action | Command |
|---|---|
| Status of whole stack | `deploy/scripts/status.sh` |
| Rebind TAK to this folder's paths | `deploy/scripts/stabilize-tak-binds.sh` |
| Move Murmur into ATAK home (one-time cutover) | `deploy/scripts/migrate-murmur.sh` |
| Issue a device cert + TAK user | `deploy/scripts/tak-add-device.sh <cn> <user> [group] [role]` |
| List connected TAK clients | `TAK_ADMIN_P12_PASS=... deploy/scripts/tak-list-devices.sh` |
| Murmur user roster | `deploy/scripts/murmur-roster.sh` |

## Ports (host → container)

| Host | Container | Service | Notes |
|---|---|---|---|
| 8089 | 8089 | ATAK/iTAK device TLS | primary device connection (CoT) |
| 8450 | 8443 | WebTAK UI / REST API | host 8443 is taken, so WebTAK is on 8450 → `https://100.74.182.117:8450` |
| 8444 | 8444 | Federation | |
| 8446 | 8446 | **Cert enrollment** | **standard port** — iTAK/ATAK auto-enroll here; do NOT remap |
| 64738 tcp+udp | — | Murmur voice | |
| 4222 / 8222 | — | NATS client / monitor | native host, not docker |

> ⚠️ Enrollment must stay on host **8446** (1:1 to container 8446). TAK clients hard-assume the
> enrollment connector is at `<server>:8446` — they enter the CoT port (8089) and enroll on 8446
> automatically. An earlier non-standard mapping (`8448→8446`) broke "Connect with Credentials".
> See [enrollment runbook](docs/enrollment.md).

## Client connection methods

1. **Connect with Credentials (recommended for iTAK)** — enter server `100.74.182.117`, port `8089`,
   SSL, username + password. Client auto-enrolls a cert on 8446. See [docs/enrollment.md](docs/enrollment.md).
2. **Data package** — one-tap `.zip` (`deploy/scripts/build-datapackage.sh`). Works on ATAK; iTAK is
   picky about layout (must be flat root, no `cert/` subfolder).
3. **Manual cert import** — import a `.p12` + set server manually.

## Two cutovers to run once (brief restarts)

1. **TAK rebind** — the live `takserver` was created with bind paths under a now-missing
   `omos-core/ATAK/...`. `stabilize-tak-binds.sh` recreates it from this real tree
   (`/Users/zeus/ATAK/server/...`). DB volume is preserved.
2. **Murmur move** — `migrate-murmur.sh` moves `~/murmur/data` → `voice/data` and recreates
   the container from this compose (fixes the `MURMUR_` vs `MUMBLE_` env name too).

After both, everything the stack needs lives under `/Users/zeus/ATAK`.

## Addressing

- LAN: `192.168.1.183`
- Tailscale: `100.74.182.117` (preferred off-LAN)

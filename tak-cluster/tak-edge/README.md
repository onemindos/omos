# TAK Server Infrastructure

Sovereign TAK Server deployment, configuration, and integration layer for OneMind OS.

## Structure

```
tak-server-infrastructure/
├── deploy/              — TAK Server docker-compose, scripts, enrollment procedures
│   ├── docker-compose.yml
│   ├── .env.example
│   ├── scripts/         — operational utilities (status, provisioning, migration)
│   └── docs/            — enrollment walkthrough
├── devices/             — Device provisioning manifests, registry, loadouts
│   ├── registry.md      — authoritative device roster
│   ├── *.yaml           — per-device manifest (iphone17, ipad, tablet)
│   ├── loadouts/        — platform-specific client configs
│   └── scripts/         — provisioning automation
├── plugins/             ⚠️ MOVED — NATS bridge source now at omos-infra/tak-cluster/plugin-src/
├── voice/               — Murmur voice server config (tactical PTT)
│   └── config/murmur.ini
├── docs/                — TAK/ATAK reference (official PDFs + test scenarios)
│   ├── TAK-Server-Configuration-Guide.pdf
│   ├── CoT-Base-Schema.pdf
│   └── overwatch-sentinel-test-scenarios.md
└── client/              ⚠️ NOT IN THIS REPO (too large for git)
    → See /Users/zeus/ATAK/client/ (4.1GB: base APKs + 30 plugins + templates)
```

## Quick Reference

### Deploy TAK Server locally (docker-compose)
```bash
cd deploy/
cp .env.example .env
# Edit .env with your network details
docker-compose up -d
```

### Provision a device
```bash
# Generate enrollment package
./devices/scripts/build-datapackage.sh zeus-iphone17
# Result: enrollment.zip → distribute to device

# Or manual provisioning
./deploy/scripts/tak-add-device.sh <device-name> <callsign>
```

### List devices
```bash
./deploy/scripts/tak-list-devices.sh
```

### Check status
```bash
./deploy/scripts/status.sh
```

## Device Roles

From `devices/registry.md`:

| Device | Role | Base Client | Key Plugins |
|--------|------|-------------|------------|
| zeus-iphone17 | Comms Operator | ATAK-civ | taktalk, VANS/Beartooth |
| zeus-ipad | ISR Operator | ATAK-civ | camera, ADS-B, reports |
| zeus-tablet | Drone Pilot | ATAK-civ | UAS_Tool, camera |

## Client APKs & Plugins

⚠️ **NOT stored in this repo** (too large, 3.6GB total)

**Location:** `/Users/zeus/ATAK/client/`

**Structure:**
```
ATAK/client/
├── base/                — Base ATAK clients (pick ONE per device)
│   ├── ATAK-5.7.0.9-civ-release.apk (351M)
│   ├── ATAK-5.7.0.9-civSmall-release.apk (107M)
│   └── TakTracker-8.1.0-release.apk (66M)
├── plugins/             — 30 plugin APKs + UAS_Tool suite (2.6GB)
├── templates/           — Device role provisioning manifests
│   ├── drone-pilot/MANIFEST.md
│   ├── comms-operator/MANIFEST.md
│   └── isr-operator/MANIFEST.md
└── README.md            — Complete provisioning guide
```

**Why not in git?** 3.6GB of binaries bloats the repo. APKs are:
- Immutable (static releases from tak.gov)
- Available independently (downloaded fresh when provisioning)
- Referenced by manifest, not committed

**To add new plugins:**
1. Download APK from tak.gov or build from source
2. Place in `/Users/zeus/ATAK/client/plugins/`
3. Update `templates/*/MANIFEST.md` to reference new plugin
4. Document plugin purpose in README

## Integration Points

### NATS Bridge Plugin
**Source:** `omos-infra/tak-cluster/plugin-src/` (Java/Gradle)

Bridges all CoT (Cursor-on-Target) traffic between TAK Server and NATS:
- TAK → NATS: `tak.cot.out` (all events)
- NATS → TAK: `tak.cot.inject` (inject events)

Deployed in k8s as `takserver-plugins` pod.

### Voice Server
**Config:** `voice/config/murmur.ini`

Murmur (Mumble) tactical voice server:
- Port: 64738
- Max users: 100
- Codec: Opus
- Bandwidth: 72 kbit/s

Run locally with deploy/docker-compose or on k8s.

### Device Management
**Provisioning:** `devices/scripts/`

Enrollment workflow:
1. Create device manifest: `devices/<device-name>.yaml`
2. Generate enrollment package: `deploy/scripts/build-datapackage.sh <device-name>`
3. Distribute enrollment.zip to device
4. Device imports cert + joins TAK Server
5. Update registry: `devices/registry.md`

## Deployment Targets

- **Local:** `deploy/docker-compose.yml` (development/testing)
- **k8s:** CloudTAK helm charts (production, onemindos fork)
- **LAN:** 192.168.1.183:8089 (Tailscale fallback: 100.74.182.117:8089)

## Reference

- **TAK Server docs:** `docs/TAK-Server-Configuration-Guide.pdf`
- **CoT protocol:** `docs/CoT-Base-Schema.pdf`
- **Test scenarios:** `docs/overwatch-sentinel-test-scenarios.md`
- **Device enrollment:** `deploy/docs/enrollment.md`
- **Murmur config:** `voice/config/murmur.ini`

## Notes

- All configs are gitignored for secrets (certs, keys, .env)
- `.gitignore` blocks `exports/*`, `*.p12`, `.env` — verified safe to commit
- NATS bridge requires matching TAK Server version (5.7.0+)
- Murmur integrates with taktalk plugin for unified voice

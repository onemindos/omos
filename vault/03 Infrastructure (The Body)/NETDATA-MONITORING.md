---
domain: 02
name: Netdata Monitoring
status: active
updated: 2026-05-22
owner: zeus
tags: [monitoring, netdata, infrastructure, observability]
---

# Netdata Monitoring — One Mind

Real-time infrastructure monitoring across all OneMind nodes via Netdata Cloud.

---

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│              Netdata Cloud (app.netdata.cloud)               │
│              Room: 88629640-0b89-4296-bc1e-828df73bf690      │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │  omos-mbp    │  │  omos-dell   │  │  omos-ha     │     │
│  │  Mac M2 Max  │  │  Dell Linux  │  │  HA OS       │     │
│  │  :19999      │  │  :19999      │  │  :19999      │     │
│  │  ACLK ✅     │  │  ACLK ✅     │  │  ACLK ✅     │     │
│  └──────────────┘  └──────────────┘  └──────────────┘     │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

---

## Nodes

| Node | OS | IP (Tailscale) | Local Dashboard | Claimed |
|------|-----|----------------|-----------------|---------|
| **omos-mbp** | macOS (M2 Max) | 100.98.178.23 | http://localhost:19999 | ✅ |
| **omos-dell** | Ubuntu 26.04 | 100.66.84.71 | http://omos-dell:19999 | ✅ |
| **omos-ha** | Home Assistant OS (Alpine) | 100.64.25.12 | http://192.168.1.180:19999 | ✅ |

---

## Cloud Configuration

| Field | Value |
|-------|-------|
| Cloud URL | https://app.netdata.cloud |
| Room ID | `88629640-0b89-4296-bc1e-828df73bf690` |
| Claim Token | `rtVvPEA3gze...` (in 1Password) |
| Connection | ACLK (Agent-Cloud Link) |

---

## Installation Method

All nodes installed via the official kickstart script:

```bash
curl -sSL https://get.netdata.cloud/kickstart.sh > /tmp/netdata-kickstart.sh
sh /tmp/netdata-kickstart.sh --dont-wait --stable-channel \
  --claim-token <TOKEN> \
  --claim-rooms 88629640-0b89-4296-bc1e-828df73bf690 \
  --claim-url https://app.netdata.cloud
```

### Node-Specific Notes

**Mac (omos-mbp):**
- Installed via Homebrew or kickstart
- Runs as LaunchAgent

**Dell (omos-dell):**
- Standard Linux install via kickstart
- Runs as systemd service: `netdata.service`
- SSH: `zeus@omos-dell` (sudo required)

**Home Assistant (omos-ha):**
- Static binary install (Alpine Linux, no package manager)
- SSH: `hassio@192.168.1.180` (passwordless sudo via SSH add-on)
- HA OS is minimal — kickstart uses static method automatically

---

## What's Monitored

Per-node metrics collected automatically:
- CPU, RAM, disk I/O, network traffic
- Process list and resource consumption
- Docker containers (where applicable)
- Disk space and inodes
- Network interfaces and connections
- System services health

---

## Access

| Method | URL |
|--------|-----|
| Cloud dashboard | https://app.netdata.cloud |
| Mac local | http://localhost:19999 |
| Dell local | http://omos-dell:19999 |
| HA local | http://192.168.1.180:19999 |

---

## Future

- [ ] Add VPS Atlas (3.234.227.222) to Netdata Cloud
- [ ] Configure alerting (Telegram notifications for disk/CPU thresholds)
- [ ] Custom dashboards for NATS, Docker, agent health
- [ ] Netdata → NATS bridge (publish metrics to `metrics.{node}.*` subjects)

---

## Related

- [[03 Infrastructure (The Body)/NATS-INFRASTRUCTURE|NATS Infrastructure]]
- [[03 Infrastructure (The Body)/nodes/edge/omos-ha|omos-ha node]]
- Netdata docs: https://learn.netdata.cloud

# OneMind OS — Observability & Monitoring

> Last Updated: 2026-03-14
> Datadog-based monitoring for both VPSes

---

## Overview

All monitoring runs through **Datadog** (site: `us5.datadoghq.com`). Each VPS runs a Datadog Agent that reports system metrics, process data, and logs.

```
  onemind-main (104.248.4.212)        onemind-sentinel (159.203.122.175)
  DD Agent v7.76.3                    DD Agent v7.76.3
  - System metrics                    - System metrics
  - Process monitor                   - Process monitor
  - Docker stats                      - Logs
  - Logs
          |                                    |
          +----------------+-------------------+
                           v
                  Datadog Cloud (us5)
                  11 Monitors
                  Alerts -> zeus@zeusdelacruz.me
```

---

## Datadog Config

| Setting | Value |
|---------|-------|
| **Site** | `us5.datadoghq.com` |
| **API Key** | `7295caf2fa7bca00eebe5e29fbc54d0e` |
| **App Key** | `d200090c3db1d9a669a125b6b6e00ce28c577382` |
| **Agent version** | 7.76.3 |
| **Config file** | `/etc/datadog-agent/datadog.yaml` (both VPSes) |
| **Process monitoring** | Enabled |
| **Log collection** | Enabled |

### Agent Config (both VPSes)

Key settings in `/etc/datadog-agent/datadog.yaml`:

```yaml
api_key: 7295caf2fa7bca00eebe5e29fbc54d0e
site: us5.datadoghq.com
tags:
  - env:production
  - role:<main|sentinel>
  - project:onemindos
  - host:<hostname>

process_config:
  process_collection:
    enabled: true

logs_enabled: true
```

---

## Monitors (11 total)

All monitors alert to `zeus@zeusdelacruz.me`.

### Host Health

| ID | Monitor | Target | Trigger |
|----|---------|--------|---------|
| 18513582 | Main VPS Host Unreachable | onemind-main | Agent stops reporting for 5 min |
| 18513583 | Sentinel VPS Host Unreachable | onemind-sentinel | Agent stops reporting for 5 min |

### CPU

| ID | Monitor | Target | Warning | Critical |
|----|---------|--------|---------|----------|
| 18513584 | Main VPS High CPU | onemind-main | >70% (5m avg) | >85% (5m avg) |
| 18513585 | Sentinel VPS High CPU | onemind-sentinel | >70% (5m avg) | >85% (5m avg) |

### Memory

| ID | Monitor | Target | Warning | Critical |
|----|---------|--------|---------|----------|
| 18513586 | Main VPS High Memory | onemind-main | <20% available | <10% available |
| 18513587 | Sentinel VPS High Memory | onemind-sentinel | <30% available | <20% available |

### Disk

| ID | Monitor | Target | Warning | Critical |
|----|---------|--------|---------|----------|
| 18513588 | Main VPS Disk Usage | onemind-main /dev/vda1 | >70% | >80% |
| 18513589 | Sentinel VPS Disk Usage | onemind-sentinel /dev/vda1 | >60% | >75% |

### Services

| ID | Monitor | Target | Trigger | Renotify |
|----|---------|--------|---------|----------|
| 18513590 | OpenClaw Down on Main | onemind-main | Process openclaw-gatewa not running | Every 15 min |
| 18513591 | OpenClaw Down on Sentinel | onemind-sentinel | Process openclaw-gatewa not running | Every 15 min |
| 18513592 | NATS Container Down on Main | onemind-main | Docker NATS container count < 1 | Every 15 min |

---

## Sentinel Datadog Startup

On the Sentinel VPS, Datadog has a delayed start to avoid memory pressure during OpenClaw boot:

**Override file**: `/etc/systemd/system/datadog-agent.service.d/after-openclaw.conf`

```ini
[Unit]
After=openclaw-sentinel.service
Requires=openclaw-sentinel.service

[Service]
ExecStartPre=/bin/sleep 45
```

---

## Agent Management

```bash
# Check agent status (either VPS)
datadog-agent status 2>&1 | grep -A3 'API Keys'

# Restart agent
systemctl restart datadog-agent

# Check agent config
grep -E '^api_key|^site|^tags' /etc/datadog-agent/datadog.yaml

# Check what processes are monitored
datadog-agent status 2>&1 | grep -A10 'Process Agent'
```

---

## Dashboard Access

- **Monitors**: https://us5.datadoghq.com/monitors/manage
- **Infrastructure**: https://us5.datadoghq.com/infrastructure
- **Host Map**: https://us5.datadoghq.com/infrastructure/map
- **Processes**: https://us5.datadoghq.com/process
- **Logs**: https://us5.datadoghq.com/logs

---

## API Usage

```bash
# Validate API key
curl -s "https://us5.datadoghq.com/api/v1/validate" \
  -H "DD-API-KEY: 7295caf2fa7bca00eebe5e29fbc54d0e"

# List monitors
curl -s "https://us5.datadoghq.com/api/v1/monitor" \
  -H "DD-API-KEY: 7295caf2fa7bca00eebe5e29fbc54d0e" \
  -H "DD-APPLICATION-KEY: d200090c3db1d9a669a125b6b6e00ce28c577382"

# Create a new monitor
curl -s -X POST "https://us5.datadoghq.com/api/v1/monitor" \
  -H "DD-API-KEY: 7295caf2fa7bca00eebe5e29fbc54d0e" \
  -H "DD-APPLICATION-KEY: d200090c3db1d9a669a125b6b6e00ce28c577382" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "[OneMind] My New Monitor",
    "type": "metric alert",
    "query": "avg(last_5m):avg:system.cpu.user{host:onemind-main} > 90",
    "message": "Description.\n\n@zeus@zeusdelacruz.me",
    "tags": ["env:production", "team:onemind"],
    "options": {"thresholds": {"critical": 90, "warning": 80}}
  }'
```

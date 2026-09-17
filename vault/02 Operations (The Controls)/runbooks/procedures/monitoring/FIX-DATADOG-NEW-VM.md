# Datadog Setup — onemindos-cloud-1

> Last Updated: 2026-04-15
> Status: COMPLETE — agent active, monitors live

## Current Setup

| Field | Value |
|-------|-------|
| Host | `onemindos-cloud-1` (DigitalOcean, `64.227.8.23`) |
| Site | `us5.datadoghq.com` |
| Agent | Datadog Agent 7 — systemd (NOT Docker) |
| API Key | 1Password → `Datadog Api key (onemindos)` |
| App Key | 1Password → `Datadog` (field: `app key`) |

> **Credentials**: Never hardcode API keys. Always read from 1Password:
> `op read "op://00-24 SO (Sovereign Operations)/Datadog Api key (onemindos)/credential"`

---

## Quick Reference — Common Tasks

### Check Agent Status

```bash
ssh root@64.227.8.23 'datadog-agent status'
```

### Restart Agent

```bash
ssh root@64.227.8.23 'systemctl restart datadog-agent datadog-agent-security'
```

### View Agent Logs

```bash
ssh root@64.227.8.23 'journalctl -u datadog-agent -f'
```

### Check Config

```bash
ssh root@64.227.8.23 'datadog-agent configcheck'
```

---

## Reinstall from Scratch

If you ever need to reinstall the agent on a new droplet:

```bash
# 1. Get API key from 1Password
DD_KEY=$(op read "op://00-24 SO (Sovereign Operations)/Datadog Api key (onemindos)/credential")

# 2. Install
ssh root@NEW_IP "DD_API_KEY=$DD_KEY DD_SITE=us5.datadoghq.com DD_HOSTNAME=onemindos-cloud-1 \
  bash -c '\$(curl -L https://install.datadoghq.com/scripts/install_script_agent7.sh)'"

# 3. Push config files (from _system/openclaw or re-create per DATADOG-OBSERVABILITY-PLAN.md)
# 4. Restart: systemctl restart datadog-agent datadog-agent-security
# 5. Recreate monitors via API (see monitor IDs in DATADOG-OBSERVABILITY-PLAN.md)
```

---

## Active Monitors (April 2026)

| ID | Name | Threshold |
|----|------|-----------|
| 19174064 | [onemindos-cloud-1] CPU High | Warn >70%, Alert >85% |
| 19174073 | [onemindos-cloud-1] Memory High | Warn <25%, Alert <15% usable |
| 19174086 | [onemindos-cloud-1] Disk Usage High | Warn >75%, Alert >85% |
| 19174174 | [onemindos-cloud-1] OpenClaw Down | Alert if process < 1 for 2min |

---

## History / What Was Fixed (April 15, 2026)

- **Old problem**: 9 stale monitors (onemind-main, Sentinel VPS) firing alerts for a deleted droplet
- **Fix**: Deleted all 9 stale monitors via API
- **Old API key** (`...14ca`) was revoked — replaced with new key (`...516a`) in 1Password
- **Agent installed** fresh on `onemindos-cloud-1` via systemd (not Docker)
- **Security agent** enabled with runtime security + network monitoring
- **Process check** configured for `openclaw-gateway`
- **Log collection** via journald for OpenClaw, ssh, fail2ban, tailscaled, cron

---

## ~~Old Content Below — OBSOLETE~~

> The original content of this file referenced Docker-based agent setup on `onemind-vm`/`onemind-os`
> and a hardcoded API key (`d66b033a...`) that is revoked. All of that is no longer relevant.
> Current agent runs as a native systemd service, not a Docker container.

## Quick Fix (5 minutes)

### Step 1: SSH into your VM

```bash
ssh your-vm-host
cd /path/to/onemind-os
```

### Step 2: Set Environment Variables

```bash
# Required for agent
export DD_API_KEY=d66b033a1148351d62ae65c048706395
export DD_SITE=us5.datadoghq.com
export DD_ENV=prod

# Required for monitor updates (get from Datadog dashboard)
export DD_APP_KEY=your_app_key_here
```

> **Get DD_APP_KEY**: https://us5.datadoghq.com/organization-settings/application-keys

### Step 3: Check Current Status

```bash
./infra/scripts/fix-datadog.sh
```

This will:
- ✅ Check if Datadog agent is running
- ✅ Verify environment variables
- ✅ Check Docker autodiscovery labels
- ✅ Check integration status (PostgreSQL, Redis, Docker)
- ✅ Identify stale monitors

### Step 4: Start/Restart Datadog Agent

If the agent is not running:

```bash
docker-compose up -d datadog
```

Or restart if it's having issues:

```bash
docker restart onemind-datadog
```

### Step 5: Verify Agent Health

```bash
# Check agent status
docker exec onemind-datadog agent status

# Check integration status
docker exec onemind-datadog agent status | grep -A 5 "postgres"
docker exec onemind-datadog agent status | grep -A 5 "redis"
docker exec onemind-datadog agent status | grep -A 5 "docker"
```

### Step 6: Recreate Monitors (with correct hostname)

```bash
# Make sure DD_APP_KEY is set!
export DD_APP_KEY=your_app_key_here

# Recreate all monitors with correct hostname
./infra/scripts/datadog-setup.sh
```

This creates:
- **[P1] OneMind API Down** - HTTP check fails 3 times
- **[P2] High CPU Usage** - CPU > 80% for 5 min on `onemind-os`
- **[P2] High Memory Usage** - Memory > 85% for 5 min on `onemind-os`
- **[P2] High Disk Usage** - Disk > 90% on `onemind-os`
- **[P2] Container Restarting** - Containers < 5 running
- **[P3] PostgreSQL High Connections** - Connections > 80
- **[P3] Redis High Memory** - Memory > 80%

Plus synthetic tests:
- **OneMind API Health** - `https://api.onemindos.dev/health` every 60s
- **OneMind UI Available** - `https://command.onemindos.dev` every 5m
- **OneMind Agents Endpoint** - `https://api.onemindos.dev/agents` every 5m

## Verify Setup

### 1. Check Datadog Dashboard

Go to: https://us5.datadoghq.com/infrastructure

You should see:
- Host: `onemind-os`
- All containers visible (postgres, redis, nats, agno, flutter-web, datadog)
- Metrics flowing

### 2. Check Monitors

Go to: https://us5.datadoghq.com/monitors/manage

You should see:
- All 7 monitors with status (OK/Alert/No Data)
- Correct hostname: `onemind-os` (not `onemind-vm`)

### 3. Check Synthetic Tests

Go to: https://us5.datadoghq.com/synthetics/tests

You should see:
- 3 tests (API Health, UI Available, Agents Endpoint)
- All passing (green)

### 4. Check APM Traces

Go to: https://us5.datadoghq.com/apm/services

You should see:
- Service: `onemind-agno`
- Traces flowing from API calls

## Troubleshooting

### Agent Not Starting

```bash
# Check logs
docker logs onemind-datadog --tail 100

# Common issues:
# 1. DD_API_KEY not set in docker-compose environment
# 2. DD_API_KEY invalid
# 3. Docker socket permission issues

# Fix: Ensure DD_API_KEY is in .env or docker-compose.yml
```

### No Metrics Flowing

```bash
# Check if containers have Datadog labels
docker inspect onemind-postgres | grep datadoghq
docker inspect onemind-redis | grep datadoghq

# If missing, restart with docker-compose:
docker-compose down
docker-compose up -d
```

### Integrations Not Working

```bash
# Check integration configs
docker exec onemind-datadog agent configcheck

# View integration logs
docker exec onemind-datadog agent status | grep -A 20 "postgres"

# Common issue: Password mismatch
# Fix: Ensure POSTGRES_PASSWORD in docker-compose matches autodiscovery label
```

### Still Getting Alerts

1. **Wait 5-10 minutes** - Metrics take time to flow
2. **Check monitor configuration** - Go to each monitor and verify query
3. **Mute old monitors** - If you see `onemind-vm` monitors, mute or delete them
4. **Check notification channels** - Ensure `@discord-webhook`, `@telegram` are configured

## Cleanup Old Monitors

If you see duplicate monitors with old hostname:

```bash
# List all monitors
curl -s -X GET "https://api.us5.datadoghq.com/api/v1/monitor" \
    -H "DD-API-KEY: $DD_API_KEY" \
    -H "DD-APPLICATION-KEY: $DD_APP_KEY" \
    | jq '.[] | select(.query | contains("onemind-vm")) | {id, name}'

# Delete a monitor by ID
curl -X DELETE "https://api.us5.datadoghq.com/api/v1/monitor/{monitor_id}" \
    -H "DD-API-KEY: $DD_API_KEY" \
    -H "DD-APPLICATION-KEY: $DD_APP_KEY"
```

## Environment Variables Checklist

Make sure these are set in your VM environment:

```bash
# In Elestio Dashboard → CI/CD → Environment Variables
DD_API_KEY=d66b033a1148351d62ae65c048706395
DD_SITE=us5.datadoghq.com
DD_ENV=prod

# Or in .env file on VM
echo "DD_API_KEY=d66b033a1148351d62ae65c048706395" >> .env
echo "DD_SITE=us5.datadoghq.com" >> .env
echo "DD_ENV=prod" >> .env
```

## Quick Links

| Resource | URL |
|----------|-----|
| **Infrastructure** | https://us5.datadoghq.com/infrastructure |
| **Monitors** | https://us5.datadoghq.com/monitors/manage |
| **Synthetics** | https://us5.datadoghq.com/synthetics/tests |
| **APM Services** | https://us5.datadoghq.com/apm/services |
| **Logs** | https://us5.datadoghq.com/logs |
| **Containers** | https://us5.datadoghq.com/containers |
| **Get App Key** | https://us5.datadoghq.com/organization-settings/application-keys |

## Summary

1. ✅ **GitHub Workflows**: Already stopped (no action needed)
2. ✅ **Datadog Scripts**: Updated with correct hostname (`onemind-os`)
3. ⚠️ **Action Required**: Run fix script on VM to verify agent and update monitors

Run this on your VM:

```bash
# 1. Export keys
export DD_API_KEY=d66b033a1148351d62ae65c048706395
export DD_APP_KEY=your_app_key_here  # Get from Datadog dashboard

# 2. Check status
./infra/scripts/fix-datadog.sh

# 3. Restart agent if needed
docker-compose up -d datadog

# 4. Recreate monitors
./infra/scripts/datadog-setup.sh
```

That's it! Alerts should stop once monitors are updated. 🎉


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]

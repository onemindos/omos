# Gateway Watchdog — External Recovery System

**Status:** ACTIVE ✅  
**Last Updated:** 2026-04-28

---

## Architecture

```
VPS systemd timer (every 2 min)
  → /usr/local/bin/gateway-watchdog.sh
    → curl http://localhost:18789/health
    → If DOWN: POST to Letta API (Legacy-AI)
      → Legacy-AI executes recovery runbook
        → SSH restart → verify → notify Zeus
          → Discord alert + Telegram alert
            → call_zeus if 10 failures
```

## Components

### VPS Watchdog (primary trigger)
- **Service:** /etc/systemd/system/gateway-watchdog.service
- **Timer:** /etc/systemd/system/gateway-watchdog.timer (every 2 min)
- **Script:** /usr/local/bin/gateway-watchdog.sh
- **Log:** /var/log/gateway-watchdog.log
- **State file:** /var/tmp/gateway-watchdog-state (tracks first-down timestamp)

### Letta Legacy-AI Agent (executor)
- **Agent ID:** agent-da585f6f-189f-472a-974c-29b9578f9b55
- **Platform:** Letta Cloud (https://app.letta.com)
- **Model:** claude-opus-4-5
- **Memory block:** gateway_watchdog (full runbook)
- **Tools:** Bash, send_discord_alert, telegram_send_alert, call_zeus, comms, fetch_webpage, BashOutput, run_code_with_tools, letta, todoist

### Recovery Runbook (in Legacy-AI memory)
1. SSH health check — confirm outage
2. Notify Zeus immediately: "Working on it..."
3. Restart gateway: systemctl --user restart openclaw-gateway
4. Wait 15s, check health, if OK notify "Restored"
5. If still down: pull logs, diagnose, retry every 2 min
6. Max 10 attempts, then phone call escalation (call_zeus)
7. Always send final report

## Notification Chain
1. Discord — send_discord_alert (Zeus #alerts)
2. Telegram — telegram_send_alert (backup)
3. Phone call — call_zeus (10 failures / VPS unreachable)

## Commands
systemctl status gateway-watchdog.timer
tail -f /var/log/gateway-watchdog.log
systemctl start gateway-watchdog.service  (manual test)
systemctl list-timers gateway-watchdog.timer

## Pending: Datadog External Monitor
Need Datadog API key for offsite synthetic check (VPS-down scenario).
Log into https://us5.datadoghq.com > Organization Settings > API Keys
Create key named gateway-watchdog, store in 1Password SO vault Datadog item.
Grid will configure the rest.


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/operations (Ops Procedures)/security (Security Ops)/security|security]]

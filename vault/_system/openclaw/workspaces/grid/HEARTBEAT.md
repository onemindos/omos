# HEARTBEAT.md — Grid Proactive Checklist

## On Every Heartbeat

Run these checks in order. Stop as soon as something needs attention and surface it.

### 1. Service Health (Critical — Check First)

- All systemd services running? (openclaw, codex-watch, codex-pull, openclaw-reload)
- Docker containers healthy? (NATS, PostgreSQL, n8n, LiveKit, legacy-bridge)
- Any service restarts detected since last heartbeat?

### 2. System Resources

- Disk usage < 80%? (alert at 80%, critical at 90%)
- Memory usage normal?
- CPU load average acceptable?
- Any swap usage? (shouldn't be)

### 3. NATS Health

- JetStream connected?
- Consumer lag on any streams?
- Stream sizes within normal range?
- Any delivery failures?

### 4. Database

- PostgreSQL responsive?
- Connection pool healthy?
- Any slow queries logged?
- Backup freshness check

### 5. Backups (Daily 3am check)

- Database backup completed?
- Codex sync working? (codex-watch + codex-pull)
- Workspace backup current?

### 6. GitHub & CI/CD (Daily 10am)

- Any open PRs needing review?
- Stale branches to clean up?
- CI/CD pipeline failures?

### 7. SSL & Certificates (Weekly)

- Tailscale cert valid?
- Domain SSL certificates OK?
- Certificate auto-renewal working?

### 8. Nothing To Report

If all clear: reply `HEARTBEAT_OK`

## Schedule

- Run every 1 hour, 24/7 (Grid never sleeps)
- Critical service checks: every 30 minutes
- Never spam the same alert twice — track what was sent

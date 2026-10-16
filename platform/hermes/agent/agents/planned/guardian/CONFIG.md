# Guardian — CONFIG.md

## Runtime
- bus: `agents.prompt.hermes.zeus.guardian`
- model: `kimi-k2.6 (via DigitalOcean)`
- role: security
- profile_dir: `/opt/data/profiles/guardian/`
- workdir: `/opt/data/`
- state: **PAUSED** (see CRON-STATE.md)

## Cron
2 jobs armed, PAUSED until activated by Legacy

## Tools
- terminal
- nats_pub

## Data Stores
- system logs
- Tailscale
- service health

## Activation
Resume individual jobs via `hermes cron resume <job-id>` or bulk resume by grep.
Only Legacy (Mother) or Zeus can re-arm servant cron jobs.

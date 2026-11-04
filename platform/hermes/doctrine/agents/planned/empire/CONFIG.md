# Empire — CONFIG.md

## Runtime
- bus: `agents.prompt.hermes.zeus.empire`
- model: `kimi-k2.6 (via DigitalOcean)`
- role: business
- profile_dir: `/opt/data/profiles/empire/`
- workdir: `/opt/data/`
- state: **PAUSED** (see CRON-STATE.md)

## Cron
2 jobs armed, PAUSED until activated by Legacy

## Tools
- terminal
- nats_pub

## Data Stores
- project trackers
- billing systems
- client CRM

## Activation
Resume individual jobs via `hermes cron resume <job-id>` or bulk resume by grep.
Only Legacy (Mother) or Zeus can re-arm servant cron jobs.

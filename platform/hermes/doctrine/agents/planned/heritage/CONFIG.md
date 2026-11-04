# Heritage — CONFIG.md

## Runtime
- bus: `agents.prompt.hermes.zeus.heritage`
- model: `kimi-k2.6 (via DigitalOcean)`
- role: records
- profile_dir: `/opt/data/profiles/heritage/`
- workdir: `/opt/data/`
- state: **PAUSED** (see CRON-STATE.md)

## Cron
2 jobs armed, PAUSED until activated by Legacy

## Tools
- terminal
- nats_pub

## Data Stores
- Obsidian vault
- document archives
- date registry

## Activation
Resume individual jobs via `hermes cron resume <job-id>` or bulk resume by grep.
Only Legacy (Mother) or Zeus can re-arm servant cron jobs.

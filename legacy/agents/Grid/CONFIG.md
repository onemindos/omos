# Grid — CONFIG.md

## Runtime
- bus: `agents.prompt.hermes.zeus.grid`
- model: `claude-haiku-4-5 (via Bedrock)`
- role: infrastructure
- profile_dir: `/opt/data/profiles/grid/`
- workdir: `/opt/data/`
- state: **PAUSED** (see CRON-STATE.md)

## Cron
2 jobs armed, PAUSED until activated by Legacy

## Tools
- terminal
- nats_pub

## Data Stores
- DOKS API
- Tailscale
- resource metrics

## Activation
Resume individual jobs via `hermes cron resume <job-id>` or bulk resume by grep.
Only Legacy (Mother) or Zeus can re-arm servant cron jobs.

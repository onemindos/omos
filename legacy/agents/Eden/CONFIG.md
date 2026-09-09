# Eden — CONFIG.md

## Runtime
- bus: `agents.prompt.hermes.zeus.eden`
- model: `claude-haiku-4-5 (via Bedrock)`
- role: sustainability
- profile_dir: `/opt/data/profiles/eden/`
- workdir: `/opt/data/`
- state: **PAUSED** (see CRON-STATE.md)

## Cron
2 jobs armed, PAUSED until activated by Legacy

## Tools
- terminal
- web
- nats_pub

## Data Stores
- RADS sensors
- weather API
- greenhouse logs

## Activation
Resume individual jobs via `hermes cron resume <job-id>` or bulk resume by grep.
Only Legacy (Mother) or Zeus can re-arm servant cron jobs.

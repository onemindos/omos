# Haven — CONFIG.md

## Runtime
- bus: `agents.prompt.hermes.zeus.haven`
- model: `claude-haiku-4-5 (via Bedrock)`
- role: home
- profile_dir: `/opt/data/profiles/haven/`
- workdir: `/opt/data/`
- state: **PAUSED** (see CRON-STATE.md)

## Cron
2 jobs armed, PAUSED until activated by Legacy

## Tools
- terminal
- nats_pub

## Data Stores
- Home Assistant
- device registries
- automation logs

## Activation
Resume individual jobs via `hermes cron resume <job-id>` or bulk resume by grep.
Only Legacy (Mother) or Zeus can re-arm servant cron jobs.

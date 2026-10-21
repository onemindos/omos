# Spartan — CONFIG.md

## Runtime
- bus: `agents.prompt.hermes.zeus.spartan`
- model: `claude-haiku-4-5 (via Bedrock)`
- role: performance
- profile_dir: `/opt/data/profiles/spartan/`
- workdir: `/opt/data/`
- state: **PAUSED** (see CRON-STATE.md)

## Cron
2 jobs armed, PAUSED until activated by Legacy

## Tools
- terminal
- web
- nats_pub

## Data Stores
- health APIs
- schedule data
- activity logs

## Activation
Resume individual jobs via `hermes cron resume <job-id>` or bulk resume by grep.
Only Legacy (Mother) or Zeus can re-arm servant cron jobs.

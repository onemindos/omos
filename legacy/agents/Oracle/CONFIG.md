# Oracle — CONFIG.md

## Runtime
- bus: `agents.prompt.hermes.zeus.oracle`
- model: `qwen3-coder-flash (via DO Inference)`
- role: intelligence
- profile_dir: `/opt/data/profiles/oracle/`
- workdir: `/opt/data/`
- state: **PAUSED** (see CRON-STATE.md)

## Cron
2 jobs armed, PAUSED until activated by Legacy

## Tools
- terminal
- web
- nats_pub

## Data Stores
- ClickHouse
- NATS streams
- weather API

## Activation
Resume individual jobs via `hermes cron resume <job-id>` or bulk resume by grep.
Only Legacy (Mother) or Zeus can re-arm servant cron jobs.

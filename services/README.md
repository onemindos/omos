# `subscribers/`

Python NATS bus subscribers and provisioning scripts. These are the live operational scripts that connect to the bus and do real work — event routing, metrics export, stream/bucket setup, and service management.

---

## Scripts

| Script | Purpose |
|--------|---------|
| `om_task.py` | Task Engine subscriber (mirrors `services/task-engine/`) |
| `om_presence.py` | Presence / voice bridge subscriber (mirrors `services/connect/`) |
| `notify_router.py` | Routes `evt.notify.>` events to the right notification channel (push, SMS, etc.) |
| `graphite_bridge.py` | Bridges NATS telemetry to Graphite/Carbon for metrics storage |
| `netdata_exporter.py` | Exports NATS metrics to Netdata |
| `subscriber.py` | Base subscriber class / generic bus subscriber |
| `subscriber_runner.py` | Process supervisor — runs and restarts subscriber processes |
| `create_streams.py` | Provisioning: creates all JetStream streams per NATS-DOCTRINE |
| `create_buckets.py` | Provisioning: creates all JetStream KV buckets |
| `seed_rules.py` | Seeds initial rules/config into KV buckets |
| `ctl.py` | CLI control tool — start/stop/status subscribers |

---

## Launchers

```bash
# Start a single subscriber
./subscribers/launch.sh om_task

# Start all subscribers
./subscribers/launch_all.sh
```

---

## Logs

All subscriber logs go to `subscribers/logs/`. Each subscriber has three log files:
- `{name}.log` — combined
- `{name}-stdout.log`
- `{name}-stderr.log`

---

## Provisioning (first-time setup)

```bash
python3 subscribers/create_streams.py   # create JetStream streams
python3 subscribers/create_buckets.py   # create KV buckets
python3 subscribers/seed_rules.py       # seed initial config
```

> **Note:** `subscribers/` is the operational home for scripts in active development. Once a service matures it graduates to `services/{name}/`.

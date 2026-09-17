---
title: "The nats CLI"
section: "09-CLI-TOOLING"
type: note
tags: [nats, cli, natscli, streams, consumers, kv, pub, sub]
source: https://github.com/nats-io/natscli
created: 2026-06-14
---

# The `nats` CLI

The primary command-line interface to NATS. Repo: **`nats-io/natscli`** (binary: `nats`). One tool covers messaging, JetStream, KV/Object, server admin, benchmarking, and schema.

See also: [[NATS/09-CLI-TOOLING/contexts|Contexts]] · [[NATS/09-CLI-TOOLING/other-tools|nats-top & bench]] · [[NATS/14-REFERENCE/cli-cheatsheet|CLI Cheatsheet]]

---

## Messaging (Core)

```shell
nats pub orders.new '{"id":1}'                 # publish
nats pub orders.new --count 100 'msg #{{Count}}'  # templated bulk publish
nats sub 'orders.*'                            # subscribe (wildcard)
nats sub orders.new --queue workers            # queue subscribe
nats request time.now ''                        # request-reply (one)
nats reply time.now '{{Time}}'                  # serve replies
```

---

## JetStream — Streams

```shell
nats stream ls                                  # list streams
nats stream add ORDERS --subjects 'orders.>' --storage file --retention limits
nats stream info ORDERS                          # config + state
nats stream report                               # all streams: msgs/bytes/consumers
nats stream view ORDERS                          # browse messages
nats stream purge ORDERS                         # drop messages (keep stream)
nats stream rm ORDERS                            # delete stream
nats stream backup ORDERS ./bak                  # backup
nats stream restore ./bak                        # restore
nats stream edit ORDERS --max-age=168h           # change config
```

## JetStream — Consumers

```shell
nats consumer ls ORDERS
nats consumer add ORDERS dispatch --pull --ack explicit --filter orders.new
nats consumer info ORDERS dispatch               # ack floor, pending, redelivered
nats consumer next ORDERS dispatch --count 10    # fetch + ack a batch
nats consumer rm ORDERS dispatch
```

---

## KV & Object Store

```shell
nats kv add CONFIG --history=5 --ttl=1h
nats kv put CONFIG version 2.14.2
nats kv get CONFIG version
nats kv watch CONFIG                              # live changes
nats kv ls CONFIG ; nats kv history CONFIG version

nats object add ASSETS
nats object put ASSETS ./firmware.bin
nats object get ASSETS firmware.bin
nats object ls ASSETS
```

---

## Server Admin

```shell
nats server ls                                    # list servers (system account)
nats server info                                  # server snapshot
nats server report jetstream                      # JS report (needs SYS creds!)
nats server report connections
nats server check connection --server localhost:4222   # health check (alerting)
nats server check jetstream --stream ORDERS
nats server ping                                   # round-trip to all servers
```

> **Pitfall (verified):** `nats server report jetstream` and `nats server ls` use the **system account** — supply SYS creds/context or they return nothing.

---

## Micro Services

```shell
nats micro ls                                     # discover services
nats micro info calc                              # endpoints + metadata
nats micro stats calc                             # request/error/latency stats
```

---

## Benchmarking & Schema

```shell
nats bench foo --pub 2 --sub 2 --msgs 100000      # throughput/latency test
nats bench foo --js --pub 1 --msgs 100000          # JetStream bench
nats schema ls                                     # NATS API schemas
nats events                                         # live advisory/event stream
```

---

## Account / Auth Helpers

```shell
nats account info                                  # JetStream account limits/usage
nats server passwd                                 # generate a bcrypted password
```

---

## Global Flags You'll Use Constantly

| Flag | Purpose |
|------|---------|
| `-s, --server` | Server URL(s) |
| `--creds` | JWT credentials file |
| `--user / --password` | User auth |
| `--context` | Use a saved [[NATS/09-CLI-TOOLING/contexts|context]] |
| `--js-domain` | Target a specific JetStream [[NATS/06-TOPOLOGY-SCALING/domains|domain]] |
| `-h, --help` | Help on any subcommand (the CLI is deeply self-documenting) |

> **Tip:** almost every command has rich `--help`. `nats stream add --help` lists every stream config flag — it's the fastest "config reference" you have.

## Sources
- natscli — https://github.com/nats-io/natscli
- CLI docs — https://docs.nats.io/using-nats/nats-tools/nats_cli
- JetStream CLI — https://docs.nats.io/using-nats/nats-tools/nats_cli/nats_cli_jetstream
- Bench — https://docs.nats.io/using-nats/nats-tools/nats_cli/natsbench

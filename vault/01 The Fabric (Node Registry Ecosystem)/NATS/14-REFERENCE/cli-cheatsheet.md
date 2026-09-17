---
title: "CLI Cheatsheet"
section: "14-REFERENCE"
type: note
tags: [nats, cli, cheatsheet, nats-cli, nsc, nk, commands, reference]
source: https://github.com/nats-io/natscli
created: 2026-06-14
---

# CLI Cheatsheet

The `nats`, `nsc`, and `nk` commands you actually type, grouped by job. Full tool tour in [[NATS/09-CLI-TOOLING/README|Section 09]].

See also: [[NATS/09-CLI-TOOLING/nats-cli|nats CLI]] · [[NATS/14-REFERENCE/config-cheatsheet|Config Cheatsheet]] · [[NATS/14-REFERENCE/monitoring-endpoints|Monitoring Endpoints]]

---

## Contexts (connection profiles)

```bash
nats context add hub --server nats://localhost:4222 --creds ~/.nkeys/app.creds
nats context add edge --server nats://barn:4222 --description "barn leaf"
nats context ls                     # list
nats context select hub             # set default
nats context info                   # show current
nats context validate hub           # test connectivity
```
Every command below honors `--context <name>` (or the selected default). See [[NATS/09-CLI-TOOLING/contexts|Contexts]].

---

## Core pub/sub & RPC

```bash
nats pub orders.new '{"id":1}'              # publish
nats pub orders.new --count 100 'msg {{Count}}'   # 100 msgs, template
nats sub 'orders.>'                          # subscribe (wildcard)
nats sub orders.new --queue workers          # queue-group subscribe
nats request svc.echo 'ping' --timeout 1s    # request-reply
nats reply svc.echo 'pong {{Request}}'       # serve replies (responder)
nats rtt                                      # round-trip latency
```

---

## JetStream — streams

```bash
nats stream ls                               # list streams
nats stream add ORDERS \
  --subjects 'orders.>' --storage file \
  --replicas 3 --retention limits            # create (or: --config file.json)
nats stream info ORDERS
nats stream view ORDERS                       # browse messages
nats stream purge ORDERS                      # delete all messages
nats stream rm ORDERS
nats stream backup ORDERS ./bak              # backup
nats stream restore ./bak                     # restore
```

## JetStream — consumers

```bash
nats consumer ls ORDERS
nats consumer add ORDERS dispatch \
  --pull --ack explicit --filter orders.new --deliver all
nats consumer info ORDERS dispatch
nats consumer next ORDERS dispatch --count 10   # pull a batch
nats consumer rm ORDERS dispatch
```

## JetStream — publish/consume helpers

```bash
nats pub orders.new 'x' -H Nats-Msg-Id:abc123   # dedup id (exactly-once)
nats stream get ORDERS 42                         # get message by seq
```

---

## Key-Value

```bash
nats kv add CONFIG --history 5 --replicas 3
nats kv put CONFIG feature.flag on
nats kv get CONFIG feature.flag
nats kv ls CONFIG
nats kv watch CONFIG                          # live updates
nats kv rm CONFIG feature.flag
nats kv history CONFIG feature.flag
```
See [[NATS/03-KV-OBJECT-STORE/key-value-store|KV Store]].

## Object Store

```bash
nats object add ASSETS
nats object put ASSETS ./big.bin
nats object get ASSETS big.bin --output ./out.bin
nats object ls ASSETS
nats object rm ASSETS big.bin
```

---

## Monitoring & diagnostics

```bash
nats server ls                       # discovered servers ($SYS)
nats server info                     # this server's vitals
nats server report jetstream         # JS usage across cluster
nats server report connections       # connection report
nats server check                    # health checks
nats server ping                     # ping all servers
nats account info                    # account limits/usage
nats events                          # stream server advisory events
nats bench foo --pub 5 --sub 5 --msgs 1M   # benchmark
```
See [[NATS/11-ADVANCED/benchmarking|Benchmarking]] · [[NATS/12-ECOSYSTEM/monitoring-exporters|Monitoring]].

---

## Schema / inspection

```bash
nats schema ls
nats stream report                   # all streams summary table
nats consumer report ORDERS          # consumer ack/pending table
```

---

## `nsc` — JWT/NKEY auth management

```bash
nsc add operator --name OmosOperator --sys      # create operator (+ SYS)
nsc add account APP
nsc add user --account APP app                  # creates app.creds
nsc edit user app --account APP \
  --allow-pub 'orders.>' --allow-sub '_INBOX.>'
nsc generate creds --account APP --name app > app.creds
nsc push -A                                      # push JWTs to resolver
nsc describe account APP
```
See [[NATS/05-SECURITY/jwt-nkeys-decentralized|JWT & NKEYS]].

## `nk` — raw NKEY tool

```bash
nk -gen user -pubout                 # generate a user NKEY (seed + pub)
nk -gen account                      # account key
nk -sign payload.txt -inkey seed.txt # sign with a seed
```

---

## Handy globals
| Flag | Effect |
|------|--------|
| `--context <name>` | use a saved context |
| `-s <url>` | server URL (override) |
| `--creds <file>` | creds file |
| `-H k:v` | add a header (pub) |
| `--js-domain <d>` | target a JetStream domain |
| `-v` / `--trace` | verbose / protocol trace |

## OneMind Relevance
- `nats context` per environment (`hub`, each `omos-<node>` leaf) — switch with `--context`, never hand-type URLs/creds.
- Fleet stream/consumer creation should be **scripted** (checked-in `nats stream add --config *.json`), not typed ad-hoc — reproducible subject law.
- `nats server report jetstream` + `nats account info` are the quick fleet-health reflexes; deeper telemetry via Surveyor.

## Sources
- natscli — https://github.com/nats-io/natscli
- nats CLI docs — https://docs.nats.io/using-nats/nats-tools/nats_cli
- nsc — https://docs.nats.io/using-nats/nats-tools/nsc

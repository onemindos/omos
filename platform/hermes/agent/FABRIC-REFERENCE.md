# Fabric reference

Stack facts that used to be written into Legacy's soul. They change, so they live here and not in `agents/legacy/SOUL.md`. The authoritative inventory is `docs/STACK.md`; this file is the short version an agent needs, with each fact checked against the live cluster on 2026-10-16. If something here is wrong, fix it here.

## NATS

- Cluster: 3-node JetStream StatefulSet in namespace `nats`. Cloud workloads connect to `nats://nats.nats.svc.cluster.local:4222`.
- Ports on the `nats` service: `4222` client, `7422` leaf nodes, `8080` WebSocket. Monitoring is `8222`.
- Accounts: **`ONEMIND`** and **`SYS`** only. (An earlier soul listed `SYS`, `CORE`, `FABRIC` and `INTERFACE`; those do not exist.)
- Edge machines connect through a local leaf node to the `omos-nats` Tailscale load balancer (leaf port 7422). The leaf is a forwarding pipe, not a peer bus.
- Streams and KeyValue buckets are defined in `platform/nats/streams/` (39 CRs) and reconciled by the NACK controller.
- Subject patterns are locked: see `docs/NATS-DOCTRINE.md`.
- Agent prompts: `agents.prompt.<framework>.<owner>.<session>`, for example `agents.prompt.hermes.zeus.legacy`.

## Other services

| What | Facts |
|------|-------|
| TAK Server | Public load balancer exposes `8089` (CoT), `8443`, `8446` (enrolment/WebTAK), `9101`, `9102`, `64738` (Mumble). CoT in and out of NATS: `tak.cot.inject` and `tak.cot.>` |
| ClickHouse | Namespace `clickhouse`, service `clickhouse-onemind`: HTTP `8123`, native `9000` |
| Admin UIs | Tailnet only: Argo CD, Grafana, Hermes dashboard, n8n, Headlamp, Hubble, the NATS UIs, Backstage preview, OpenClaw |
| Backups | Hermes volume: VolSync nightly 01:00. See `docs/operations/hermes-recoverability-plan.md` |

## Handy commands

```bash
nats micro ls --server nats://localhost:4222            # agents registered on the bus
nats req agents.prompt.hermes.zeus.legacy-edge "<task>"  # prompt an edge node
nats sub "agents.>"                                      # watch agent traffic
```

## Corrections to older documents

The July and August genome documents (now archived) and the earlier soul carried values that are no longer true:

| Said | Now |
|------|-----|
| Model provider AWS Bedrock | DigitalOcean inference (see `AGENTS.yaml` for each agent's model) |
| NATS accounts SYS, CORE, FABRIC, INTERFACE | `ONEMIND` and `SYS` |
| ClickHouse port 9001 | 8123 (HTTP) and 9000 (native) |
| `om_presence` is a live voice bridge | The voice bridge was archived in August |
| Pod volume is NFS | Block storage PVC `hermes-data` |
| Cluster `do-nyc3-onemind-hub`, 3 nodes of 4 vCPU | Cluster `omos-cloud`, 5 nodes across two pools |
| Repo `omos-core`, repo `omos-legacy` as the genome | Both retired; this directory in `omos` is the source |

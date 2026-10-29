# Owner proposal: stream caps, COT retention, CMD discard policy and the DET stream (A22)

Date: 2026-10-29. Feature A22 of the OneMind XR build (BUS-DESIGN section 4.1 and section 11).
Status: proposed in a draft PR. **Merging is activation.** `platform-nats-streams` is an
automated Argo Application (`selfHeal: true`), so NACK applies every change here within a
few minutes of the merge. Nothing in this PR has been applied to the cluster; the cluster was
only read.

## What changes

Live numbers are from `jsz?streams=1&config=1` on nats-0, 2026-10-29 11:29 UTC (read-only).
Byte caps are per replica: every R3 stream counts fully on every node.

| Stream | Live today (messages / bytes) | maxAge before → after | maxBytes after | Other | Data removed on sync |
|---|---|---|---|---|---|
| `COT` | 9,531,820 / **4,948,843,022** (first 2026-09-29) | 720h → **72h** | 3 GiB | — | **Yes: about 4.45 GB per replica** (decision A22-P1) |
| `CMD` | 2 / 166 | 720h → **8760h** | 512 MiB | `discard: old` → **`new`** | No (the 2 messages are kept longer) (A22-P2) |
| `TEL` | 0 / 0 | 168h → 72h | 4 GiB | + `tel.feed.>`, `maxMsgsPerSubject: 2000`, `allowDirect` | No |
| `DET` (new) | does not exist | 168h | 2 GiB | `det.>` | No (A22-P3) |
| `INTEL` | 92,129 / 103,612,881 (first 2026-08-04) | 2160h (keep) | 1 GiB | `duplicateWindow` 2m → **1h** | No: 104 MB is 10% of the cap (A22-P4) |
| `GEO` | 0 / 0 | 720h → 168h | 512 MiB | — | No (A22-P5) |
| `AGENTS` | 20,118 / 3,881,988 | 168h (keep) | 1 GiB | — | No |
| `EVT` | 9 / 5,292 | 2160h (keep) | 2 GiB | — | No |
| `ENV`, `MISSION`, `SYS` | 0 / 0 | keep | 1 GiB each | — | No |
| `SAAS`, `VOICE`, `ZEUS` | 0 / 0 | keep | 512 MiB each | — | No |
| `WORLD` | 0 / 0 | 24h (keep) | 2 GB (keep) | — | No change |

Also added: alert `NATSCmdStreamNearFull` (80% of the CMD cap, on `nats_stream_total_bytes`,
label `stream="CMD"`, a metric confirmed present in Prometheus read-only) and the contract test
`scripts/tests/test_gitops_stream_caps.py`.

Not in this PR: `ent.yaml` and `kv-entity-correlation.yaml` (PR #552),
`kv-ingest-locks.yaml` and `kv-ingest-state.yaml` (PR #553), `evt.entity.>` on `EVT`
(BUS-P5, filed as A4-P5 in #552), `maxBytes` on the other KV buckets (BUS-DESIGN section 11
"other 21 kv-*.yaml"; follow-up). `TEL` gets `tel.feed.>` with the same two lines as #553,
so either PR can merge first without a conflict (checked with `git merge-tree`).

## Per-node budget

Sum of all Stream caps with #552's `ENT` included: **21.36 GiB**. Adding the object-store
budget (9 GiB) and the KV budget (2 GiB) gives 32.36 GiB of the 43.14 GiB usable per node
(`max_file_store` 45 GiB minus `reserved_storage` 2,001,048,576 bytes). The contract test
fails CI if a new Stream CR lacks `maxBytes` or the sum stops fitting. The engine test also
creates every CR on a local server whose `max_file_store` is 45 GiB, which only succeeds if
the reservations fit.

## Decisions for the owner

### A22-P1: cut COT from 30 days to 72 hours (removes data)

- **What is lost.** On sync the server removes every `COT` message older than 72 h at once (the
  local test shows a shorter `max_age` on an existing stream removes older messages within a
  second). Today that is about 8.6 M messages and about 4.45 GB per replica, about 13 GB of
  disk over the three nodes. About 0.5 GB stays (4.95 GB / 30 days ≈ 165 MB per day × 3).
- **What the data is.** One subject, `tak.cot.inject`: MTA Metro-North/LIRR CoT from
  Redpanda Connect `feed-mta` (about 3.1 msg/s), ISS-group satellites from `feed-satellites`
  and `omos-tak-projector` output (omos-live-state §2-§4). It has no JetStream consumer, and
  TAK does not read it (the `tak-nats-leaf` path carries nothing). It is regenerated
  telemetry, not an operator record.
- **Who might want it.** Replaying `COT` to rebuild entities (#552, entity.md section 10) needs
  only the last 72 h. Any longer-term analysis of transit/satellite history would lose 27 days.
- **Recovery copy.** `backup-nats-streams` archives all streams nightly and last succeeded at
  2026-10-29 03:34 UTC (CronJob status, read-only). Archive retention in the bucket is
  **unverified**: keep the last archive taken before the merge if the history matters.
- **The 3 GiB cap** is six times the expected 72 h volume, so it only acts if the CoT rate
  rises roughly sixfold (for example aircraft and ships returning without A0 quotas).
- **Alternative.** Keep 720h with a 3 GiB cap: the server then drops the oldest 1.95 GB at
  sync instead, and the stream keeps about 18 days. Say so and I change one line.

### A22-P2: CMD keeps a year of commands and refuses when full

`discard: new` is a deliberate deviation from every other stream. When `CMD` reaches 512 MiB
it refuses new publishes (error 10077 "maximum bytes exceeded", shown in the local test)
instead of deleting the oldest audit records. The policy service treats the failed publish as
a refusal, so no command leaves without its audit record (run rule 6). The cost: a full `CMD`
stops all commands until messages age out at 365 days or the cap is raised. The 80% alert
gives warning; at today's volume (2 messages) that is far away. Alertmanager has no receiver
yet (docs/STACK.md), so the alert shows only in Grafana and Prometheus until one is set.

### A22-P3: new DET stream

`det.{domain}.{type}.{id}` is in the locked grammar but no stream binds it, so detections are
core-only. `DET` stores them for 7 days under a 2 GiB cap (6 GiB of disk across the cluster
at most). Nothing publishes `det.*` today, so it starts empty. The Hermes leaf already holds
a `det.>` interest.

### A22-P4: INTEL dedup window 2m to 1h

`feed-intel` already sets `Nats-Msg-Id` to the source event id (`platform/omos/intel-scripts.yaml`
lines 528, 769, 1005, 1255, 1498). A one-hour window also drops repeats that a restarted
collector re-sends within the hour. Risk: if a source reuses one id for an updated event
within the hour, the update is dropped as a duplicate. The window costs a little server
memory for the id table (INTEL writes are sparse).

### A22-P5: GEO to 7 days

`GEO` is empty and is a retirement candidate (BUS-P10). This shrinks it only; it is not deleted.

## Activation order

1. **Before merging (owner):** confirm the 2026-10-29 (or later) `backup-nats-streams` archive
   and its `.verified.json` marker exist in the backup bucket, and keep it if the 30-day COT
   history matters. A fresh one-off run is `kubectl -n nats create job --from=cronjob/backup-nats-streams backup-nats-streams-pre-a22`
   (a cluster write: owner only).
2. **Merge this PR.** Argo syncs `platform-nats-streams`; NACK updates 14 streams and creates
   `DET`. The COT trim happens at that point. If you approve everything except P1, revert
   commit "feat(nats): COT retention 720h -> 72h" on the branch before merging, or ask for
   it to be split out.
3. **#553** (omos-ingest) can merge before or after: its `TEL` lines are identical, and it adds
   the `ingest-*` KV buckets.
4. **#552** (`ENT`) after its own owner approval; its CR stays excluded until then.
5. Server 2.14.2 → 2.14.9 before relying on per-message TTL in `ENT` (BUS-DESIGN 4.1
   prerequisite 1). Not part of this PR.

## Checks after the sync (read-only)

```bash
kubectl --context do-nyc3-omos-cloud -n nats get streams.jetstream.nats.io   # all Ready, det present
kubectl --context do-nyc3-omos-cloud get --raw \
  '/api/v1/namespaces/nats/pods/nats-0:8222/proxy/jsz?streams=1&config=1' \
  | python3 -c 'import json,sys; [print(s["name"], s["config"].get("max_age"), s["config"].get("max_bytes"), s["config"].get("discard"), s["state"]["bytes"]) for a in json.load(sys.stdin)["account_details"] for s in a["stream_detail"] if not s["name"].startswith("KV_")]'
```

Expected: `COT` `max_age` 259200000000000 ns, `max_bytes` 3221225472, bytes about 0.5 GB;
`CMD` `discard` `new`; `TEL` 9 subjects including `tel.feed.>`; `DET` present with
`max_bytes` 2147483648; `NATSJetStreamStorageHigh` baseline drops from about 10% to about 1%.

## Rollback

Reverting the PR restores the old limits for future messages. Revert does **not** bring back
COT messages removed by the 72 h cut; only the backup archive does (`platform/backups/`).
`preventDelete: true` stays on every CR, so no rollback path deletes a stream.

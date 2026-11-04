# Reference architecture proposals RA1 … RA4

Source: onemind-xr `docs/design/REFERENCE-ARCHITECTURE.md` §13 (branch `run/m0`). Index and
ruling order: [README.md](README.md).

---

## RA1. `evt.video.>` is not bound in `EVT`

**Proposed.** The video registry publishes `evt.video.{entity_id}.{online|offline}` when a
camera or drone stream appears or drops in MediaMTX. Add `evt.video.>` to `EVT`
(`platform/nats/streams/evt.yaml`) in the same CR PR as `evt.entity.>` (BUS-P5).

**Why the grammar does not cover it.** The subject fits `evt.{domain}.{entity}.{event}`, but
`EVT` binds ten domains (`evt.yaml` lines 11–20) and not `video`, so the events would be
core-only (onemind-xr `docs/research/video-and-plugins.md` item 6).

**Alternatives considered.**
- No events: clients read the entity's `media[]` field and KV `OVERWATCH_VIDEO`
  (`platform/nats/streams/kv-overwatch-video.yaml`) only. Works for current state, but leaves no
  history of outages. This is the interim.
- `evt.entity.{entity_id}.video_online` under BUS-P5's domain: keeps one filter for entity facts,
  but mixes registry state into entity lifecycle events and depends on BUS-P5. Rejected.

**Interim reading.** Entity `media[]` and KV `OVERWATCH_VIDEO`; nothing publishes `evt.video.>`.

**Depends on it.** The video registry (K lane; no open PR). No open PR publishes it.

```subjects
proposed  evt.video.{entity_id}.{online|offline}
filter    evt.video.>
filter    evt.entity.>
current   evt.{domain}.{entity}.{event}
alternative  evt.entity.{entity_id}.video_online
```

- [ ] Approve as proposed
- [ ] Approve with changes:
- [ ] Reject

---

## RA2. Policy names and audit subject used everywhere

Two halves, each filed once: the RPC names under
[ONT-P14](ontology.md#ont-p14-policy-rpc-names) and the audit subject under
[ONT-P10](ontology.md#ont-p10-one-audit-subject). The design set adopted both on 2026-10-29
(ONTOLOGY §8, asset-command.md §1, mission-task.md §6, FEATURES D3, D5, I2, I4), and #554 uses
exactly these names. The Apollo mapping in onemind-xr `docs/research/palantir-patterns.md` that
wrote `evt.command.{result}` is superseded by ONT-P10.

```subjects
proposed  svc.policy.command.propose
proposed  svc.policy.task.propose
proposed  evt.so.action.{event}
legacy    evt.command.{result}
alternative  svc.policy.propose
```

- [ ] Ruled with ONT-P10 and ONT-P14

---

## RA3. Site upgrade as a command (information only)

**Recorded, not proposed.** onemind-xr `docs/research/palantir-patterns.md` §8 suggested site
upgrades as commands, which would need a new verb under `cmd.fleet.node.{site}.{verb}` outside
the SL5 list. Not proposed: the site agent pulls upgrades inside its change window and the
approval is the merge of the release PR. Only status is published, on
`tel.fleet.node.{site}.status` (BUS-P1).

```subjects
current   tel.fleet.node.{site}.status
alternative  cmd.fleet.node.{site}.upgrade
alternative  cmd.fleet.node.{site}.{verb}
```

- [ ] Noted (no ruling needed)

---

## RA4. KLV subject form

The video lane wrote `tel.{node}.drone.sensor` (onemind-xr `docs/research/video-and-plugins.md`
item 7; `docs/research/standards.md` line 207); BUS-DESIGN uses `tel.fleet.drone.{id}.sensor`.
The form is ruled by [BUS-P1](bus.md#bus-p1-telemetry-token-two-is-a-class) (class second,
fleet in §9.1 form) and the metric name by
[BUS-P7](bus.md#bus-p7-fleet-telemetry-metrics-beyond-position-and-status).

```subjects
proposed  tel.fleet.drone.{id}.sensor
alternative  tel.{node}.drone.sensor
```

- [ ] Ruled with BUS-P1 and BUS-P7

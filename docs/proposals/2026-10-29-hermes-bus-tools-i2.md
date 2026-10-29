# Owner proposals: Hermes bus tools (FEATURES I2, D5, D11)

Status: **proposed, for the owner**. Date: 2026-10-29. PR: `feat/hermes-nats-tools`
(`platform/hermes/agent/plugins/onemind-bus`). Filed under onemind-xr decision 0004 item 8.

**No new subject pattern is proposed.** Every subject the plugin uses is in omos
`docs/NATS-DOCTRINE.md` §9.1/§9.7 or in the `svc.{service}.{endpoint}` grammar, and is already
listed in onemind-xr `docs/design/BUS-DESIGN.md` §3.3–§3.4. The items below are permission,
payload and identity decisions. They belong in the A24 register (omos PR #556) as rows of
`cross-reference.md` under this PR, all "in PR" (no grammar ruling needed); this file is the
proposal text, not a second register. The rows are not in #556 yet; this PR does not edit
#556, so they are handed off to its owner to paste into the Mapping table:

```
| #564 `feat/hermes-nats-tools` | `docs/proposals/2026-10-29-hermes-bus-tools-i2.md` |   (PR table)
| #564 I2-P1 `hermes` may request `svc.layer.snapshot` | permission | — | in PR |
| #564 I2-P2 map-action acks from CLIENT back to AGENTS (service import) | permission | — | in PR |
| #564 I2-P3 agent map actions recorded in `CMD` | stream capture | — | in PR |
| #564 I2-P4 `cmd.map.*` payload v1 and ack | payload | — | in PR |
| #564 I2-P5 `on_behalf_of` for agent proposals | identity | — | in PR |
```

```subjects
current   svc.entity.get
current   svc.entity.query
current   svc.layer.list
current   svc.layer.meta
current   svc.layer.snapshot
current   svc.policy.command.propose
current   svc.policy.task.propose
current   svc.policy.action.get
current   cmd.map.viewport.set
current   cmd.map.filter.set
current   cmd.map.highlight.{uid}
current   cmd.map.annotation.{id}.{verb}
```

## I2-P1: let Hermes read `svc.layer.snapshot` (bounded)

BUS-DESIGN §3.4 gives Hermes only `svc.layer.list` and `svc.layer.meta`; the A19 layout (omos
PR #562, `identities.yaml` user `hermes`, `AGENTS` imports) follows it. FEATURES I2 asks for
`svc.layer.*` reads, and an area question ("what is in this field") is cheaper as one snapshot
than as entity query pages. The plugin always sends an area (bbox or zone), `fields: render`
and `max ≤ 200`, so the read stays small and carries only render fields (the prompt-injection
surface of security-and-trust.md T1 is the callsign, as for `svc.entity.*`).

Change if approved: `hermes` publish allow `svc.layer.snapshot`; `AGENTS` import
`svc.layer.snapshot` from `ONEMIND` with `share: true`. The layer service's export already uses
`response_type: stream` (PR #561 acceptance test), which the request-many reply needs.
Until then the tool is off: `ONEMIND_BUS_LAYER_SNAPSHOT` defaults to `0` and the tool is not
registered, so the model never sees it. If approved: set `ONEMIND_BUS_LAYER_SNAPSHOT=1` after the
A19 change is live. If rejected: nothing to do.

- [ ] approve  - [ ] approve with changes  - [ ] reject

## I2-P2: carry map-action acks from CLIENT back to AGENTS

D11 says Hermes reports only actions that succeeded, so a map tool needs the client's ack. In the
A19 layout `AGENTS` exports `cmd.map.>` to `CLIENT` as a **stream**. Proved on nats-server
v2.15.2 (`tests/test_account_layout.py`): the client receives the action, but its reply lands in
`CLIENT` and never reaches Hermes, so every map tool would report `no_client_ack`.

Change if approved: `CLIENT` exports `cmd.map.>` as a **service** with `response_type: stream`
to `AGENTS`, and `AGENTS` imports it (replacing the stream export and import). Every client
subscriber still receives each action (core fan-out inside `CLIENT`), and each can answer.
Proved in the same test (`test_i2p2_service_import_returns_the_ack`).

- [ ] approve  - [ ] approve with changes  - [ ] reject

## I2-P3: record agent map actions in `CMD`

Agent map actions are level A1: "no gate; logged" (agents-that-act.md §8). Under A19 they stay in
`AGENTS`/`CLIENT`, and `CMD` (`cmd.>`, account `ONEMIND`) never sees them. The plugin logs every
tool call to the Hermes log, but the bus record is missing.

Change if approved: `AGENTS` also exports `cmd.map.>` as a stream to `ONEMIND`, which imports it,
so `CMD` stores Hermes' map actions (`Nats-Msg-Id` = the payload `id`). Volume is small (one
message per user-visible action).

- [ ] approve  - [ ] approve with changes  - [ ] reject

## I2-P4: ratify the `cmd.map.*` payload v1 and the ack

The doctrine §9.7 payloads are examples. The plugin needs more: `viewport.set` modes `fly_to`,
`frame` and `tour`; `filter.set` `show`/`hide` deltas (the §9.7 `layers` full list stays valid);
`highlight` `select`/`clear`; annotation kinds `note`, `line`, `area`, `boundary`, `measure` with
GeoJSON geometry; and a reply `{v, id, ok, client, applied | error}`. Measure and camera tour use
existing subjects (annotation and viewport) instead of new ones. The full table is in the plugin
README ("`cmd.map.*` payload v1 and the client ack").

Asked: approve the payload as the contract the clients implement (Unreal `OneMindBus`, L0 web
client J1). No client subscribes to `cmd.map.>` today; implementing it is a follow-up for those
lanes.

- [ ] approve  - [ ] approve with changes  - [ ] reject

## I2-P5: whose behalf an agent proposal is on

The policy service (omos PR #554) requires an agent proposal to name a human in
`requested_by.on_behalf_of`, and checks that both are `proposer` and the human is `operator` of
the asset. Hermes does not have a verified per-message sender yet (Synadia Agent Protocol §13
signed identity is not enabled; agents-that-act.md fact 23). Interim: a fixed
`ONEMIND_BUS_ON_BEHALF_OF` per Hermes deployment (`user:zeus` for Legacy). The human confirm in
the UI is unchanged, so a wrong `on_behalf_of` can at most create a proposal that the named human
then sees and rejects.

Asked: accept the static setting until §13 sender identity is on, then take `on_behalf_of` from
the verified prompt sender.

- [ ] approve  - [ ] approve with changes  - [ ] reject

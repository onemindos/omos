# Grammar and stream owner proposals: index (FEATURES A24)

Status: **proposed, for the owner**. Date: 2026-10-29 (updated in the integration pass
the same day). Base: `main` at `62e3d22b`.
Filed under onemind-xr decision 0004 item 8 ("use the existing grammar only; if something truly
has no pattern, use the closest existing pattern and write a proposal in omos `docs/` for the
owner"). **This PR changes no live subject, stream, Application or workload.** It adds
documents and a read-only checker, `scripts/check_proposal_subjects.py`.

Source register: onemind-xr `docs/design/BUS-DESIGN.md` §10.1 (branch `run/m0`), which combines
BUS-DESIGN §10, ONTOLOGY §9, REFERENCE-ARCHITECTURE §13, `contracts/cot-projection.md` §11 and
JADC2-MAP §6 F1. Every ID there has an entry here. Duplicates are argued once, under the first
ID; the other IDs point to it and are ruled with it. Proposals in the open omos PRs (#550 … #569)
and in the onemind-xr client lanes (#4, #5) are mapped in [cross-reference.md](cross-reference.md),
so nothing is asked twice. Proposals raised after this register was first written are added
**by reference**: their text stays in the PR that raised them, and this register adds a mapping
row and, where they need a grammar or stream ruling, a row in the index below.

| File | Entries |
|---|---|
| [bus.md](bus.md) | BUS-P1 … BUS-P11 |
| [ontology.md](ontology.md) | ONT-P1 … ONT-P14 |
| [reference-architecture.md](reference-architecture.md) | RA1 … RA4 |
| [cot-projection.md](cot-projection.md) | CP-1 … CP-4 |
| [f1-intercept-example.md](f1-intercept-example.md) | F1 |
| [cross-reference.md](cross-reference.md) | Every PR proposal against this register: A4-P*, A0-*, POL-*, H2 P*, A7 P*, A1-P*, A5-P*, A6-P*, A19-P*, A21-P*, A22-P*, D1-P*, D2-P*, H3-P*, I1-P*, I2-P*, J2-P*, K-P*, WEB-P*, UE-P* |
| [integration-resolutions.md](integration-resolutions.md) | INT-R1 … INT-R4 (cross-PR blockers from the integration pass) |
| `i1-agents-watch-briefs.md` | I1-P1 … I1-P7, filed by #567 into this folder; arrives when #567 merges |

## Rule in this order (what blocks activation)

"Activation" means syncing a service or stream that an open PR has built, or turning on
behaviour it ships switched off. Within a tier, the order is the order to rule.

| Tier | ID | Decision asked | Blocks | Duplicates ruled with it |
|---|---|---|---|---|
| **1. Blocks syncing a built PR** | ONT-P1 | Approve the `ENT` stream CR | #552 sync (A4-P1), and #561, #563, #569 stacked on it | — |
| | INT-R2 | The entity service fills `capabilities` from the bridge's `commands_accepted` | every asset command in the cluster (#554 denies `not_capable`) | WEB-P5, H3-P4 |
| | INT-R1 | One `params_hash` form, over the action the message executes | every task-dispatched asset command (#568 with #550, #566) | D2-P1, H3-P2 |
| | BUS-P2 | Asset commands are `cmd.fleet.{class}.{id}.{verb}`; mission verb `cancel` | #554 and #550 command path; #566 | ONT-P4 |
| | ONT-P14 | Policy RPC names `svc.policy.{command\|task}.{propose\|confirm}`, `svc.policy.action.get`, `svc.policy.check` | #554 sync | RA2 (naming half) |
| | ONT-P10 | One audit subject `evt.so.action.{event}` | #554 sync | RA2 (audit half) |
| | ONT-P7 | `CMD` 365 d, 512 MiB; discard policy as corrected by A22-P2 (`discard: old` for now) | #554 audit retention; #555 | A22-P2 |
| | BUS-P1 | Telemetry token two is a class; fleet uses doctrine §9.1; `tel.feed.>` in `TEL` | #550, #551 status, #553 heartbeat binding (A0-2), #561 site monitor, #566 | ONT-P5, ONT-P6, RA4 (form), A5-P5, A1-P3 (status) |
| | INT-R3 | One `CLIENT` import list in #562 | browser and Unreal clients under #562; Hermes map acks | WEB-P2, UE-P3, I2-P2, A5-P7, D2-P5 (exports), K-P7 |
| | CP-4 | `tak.cot.mesh.rx`/`.tx` are the mesh pair; grammar line "CoT to and from TAK" | #551 at a site | — |
| **2. Ratifies what open PRs already publish** | F1 | Replace the `cmd.drone.x47.intercept` example; forbid the verb | none (safety; doc only) | — |
| | BUS-P9 | Dead letters on `evt.sys.{source}.rejected` | none (`SYS` stores it today) | A1-P3 (dead letters) |
| | ONT-P9 | `cmd` domains `entity`, `team`, `watch`, `feed` | #554 `cmd.entity.>` | I1-P7 (information) |
| | ONT-P8 | `entity_type` vocabulary (and `satellite`) | #552 (A4-P4) | I1-P5 (class `agent`, with BUS-P1) |
| | BUS-P3 | Live 5–6-token `intel.>` subjects kept; new collectors four tokens | #553 (A0-3) | — |
| | CP-3 | `tak.cot.mesh.tx` senders use core publish, no reply | #551 | — |
| | CP-2 | Gateway skips OneMind's own CoT when `MESH_TO_TAK_SERVER=true` | #551 (flag stays `false`) | A19-P8 |
| | BUS-P7 | Any metric name in `tel.fleet.{class}.{id}.{metric}` | K5 KLV extractor (no PR) | RA4 (metric) |
| | D1-P2 | Zone command results on `evt.map.geofence.{created\|updated\|deleted\|refused}` | #569 closing zone actions | — |
| | D1-P7 | Zone verb `remove` canonical, `delete` alias | #569 | — |
| | D2-P4 | `svc.mission.{get\|list\|report\|apply}`; refusals on `evt.mission.{id}.rejected` (one refusal word with D1-P2) | #568 | — |
| | A5-P4 | `svc.layer.{list\|meta\|snapshot\|watch\|unwatch}`, streamed replies | #561 | — |
| | I1-P4 | Briefs on `evt.so.brief.published`; `svc.brief.{generate\|latest}` | #567 | — |
| | I1-P5 | Agents are fleet class `agent`, `tel.fleet.agent.{id}.position` | #567 | — |
| | K-P7 | `svc.registry.video.{get\|list}` | #559 | — |
| **3. Unlocks behaviour shipped off** | BUS-P5 | Domain `entity`; `evt.entity.>` in `EVT` | #552 lifecycle events; #554 closing `cmd.entity.>` actions (POL-14); #567 stored alerts | ONT-P2 |
| | D1-P3 | Alerts on `evt.so.alert.{raised\|updated\|resolved}` in `EVT` | #569 D6 alert stream; #567 alerts | I1-P3 (one stored alert fact; see cross-reference) |
| | ONT-P11 | Asset tasks in the `mission` domain, never `evt.task.>` | #554 tasks; #568 task service | — |
| | ONT-P13 | `evt.vehicle.>` in `EVT` | #554 vehicle outcomes (POL-4) | — |
| | BUS-P6 | World classes `rail`, `transit` | #558 grammar pipelines (replicas 0) | A1-P1 |
| | RA1 | `evt.video.>` in `EVT` | #559 registry events (core publish only) | K-P1 |
| | INT-R4 | The projector serves `svc.cot.send` (kinds `geochat`, `shape`) | web and Unreal chat and shapes; agentops' exit from I1-P2 | WEB-P4, UE-P1 |
| | I1-P1 | Agent contacts in TAK (cot-projection §4.1 exception) | #567 contact PLI (ships off) | — |
| | I1-P2 | Agentops as an interim second publisher of the TAK outbound subjects | #567 GeoChat replies (ship off) | ends with INT-R4 |
| | UE-P5 | The client's own TAK stream carries only its SA and human chat | onemind-xr #4 shapes on its own stream (off) | — |
| **4. Blocks later features** | ONT-P3 | `DET` stream CR | #555 `DET` CR (A22-P3); detectors (H8, Remote ID); takml move (K-P5); #552 skips `DET` while absent | A22-P3, K-P5 |
| | BUS-P10 | Shrink, then retire `GEO`; no new `evt.tak.cot.{uid}` publisher | #555 (A22-P5) | A22-P5 |
| | CP-1 | `cmd.projector.{instance}.{disable\|enable}` | #563 keeps the legacy form until ruled | recorded with BUS-P4; A19-P9 |
| | BUS-P4 | Out-of-grammar legacy subjects recorded, not extended | none | CP-1 |
| | BUS-P11 | Doctrine §4: buffering needs a site stream | site leaf (A23) | — |
| | UE-P2 | `svc.tak.mission.>` served by tak-client; `evt.mission.{tak-token}.changed` | C7 Data Sync (needs a #560 follow-up) | — |
| | K-P9 | KV `plugin-registry` and `svc.plugins.list` | plugin loader (no PR) | — |
| **5. Information only** | BUS-P8, ONT-P12, RA3 | Noted, no ruling needed | — | — |

Tier 1 holds every stream or naming decision without which a built service cannot be synced
safely; nothing in tiers 2–5 has to wait for tier 1. INT-R2 and INT-R1 sit high in tier 1 because without them
no asset command passes the policy service and the bridge in the cluster, whatever else is
ruled. Two more cross-PR blockers are service items, not grammar, and are ruled in their own
PR: D1-P4 (`zone.*` in #554's catalogue, needed for any zone write; WEB-P1 is the same ask) and
I1-P6 (the `agentops` user, needed to sync #567). F1 is listed first in tier 2 because it is
a one-line documentation fix that removes a stop-line conflict and depends on nothing.

Subjects named on this page (checked by the script):

```subjects
proposed  cmd.fleet.{class}.{id}.{verb}
proposed  svc.policy.{command|task}.{propose|confirm}
proposed  svc.policy.action.get
proposed  svc.policy.check
proposed  evt.so.action.{event}
proposed  tel.fleet.{class}.{id}.{metric}
proposed  tak.cot.mesh.rx
proposed  tak.cot.mesh.tx
proposed  evt.sys.{source}.rejected
proposed  cmd.projector.{instance}.{disable|enable}
filter    tel.feed.>
filter    evt.entity.>
filter    evt.vehicle.>
filter    evt.video.>
filter    cmd.entity.>
filter    intel.>
proposed  evt.map.geofence.{created|updated|deleted|refused}
proposed  svc.mission.{get|list|report|apply}
proposed  evt.mission.{id}.rejected
proposed  svc.layer.{list|meta|snapshot|watch|unwatch}
proposed  evt.so.brief.published
proposed  svc.brief.{generate|latest}
proposed  tel.fleet.agent.{id}.position
proposed  svc.registry.video.{get|list}
proposed  evt.so.alert.{raised|updated|resolved}
proposed  svc.cot.send
proposed  evt.mission.{tak-token}.changed
proposed  svc.plugins.list
filter    svc.tak.mission.>
legacy    cmd.drone.x47.intercept
legacy    evt.task.>
legacy    evt.tak.cot.{uid}
```

## How to rule

Each entry ends with three checkboxes (approve, approve with changes, reject). A duplicate has
one box, "ruled with …". Ticking a box changes nothing live: each approved item becomes its own
small PR (a stream CR edit, a `CLAUDE.md` or doctrine line), applied by `platform-nats-streams`
or merged as documentation.

## Checking the subjects

Every subject in these files is listed in a fenced `subjects` block with a tag (`proposed`,
`interim`, `current`, `filter`, `legacy`). The checker matches each against the `CLAUDE.md`
grammar block and `docs/NATS-DOCTRINE.md` §2 and §9.1, applies the stop-line word list and the
closed verb list (SL5), and reports which stream on the checkout stores the subject today:

```bash
python scripts/check_proposal_subjects.py --self-test
```

It fails if a proposed, interim, current or filter subject is outside the grammar, if a
proposed or interim subject has a stop-line word, or if a subject in the text is missing from
the blocks. Only subjects tagged `legacy` may be off-grammar.

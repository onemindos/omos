# Grammar and stream owner proposals: index (FEATURES A24)

Status: **proposed, for the owner**. Date: 2026-10-29. Base: `main` at `62e3d22b`.
Filed under onemind-xr decision 0004 item 8 ("use the existing grammar only; if something truly
has no pattern, use the closest existing pattern and write a proposal in omos `docs/` for the
owner"). **This PR changes no live subject, stream, Application or workload.** It adds
documents and a read-only checker, `scripts/check_proposal_subjects.py`.

Source register: onemind-xr `docs/design/BUS-DESIGN.md` §10.1 (branch `run/m0`), which combines
BUS-DESIGN §10, ONTOLOGY §9, REFERENCE-ARCHITECTURE §13, `contracts/cot-projection.md` §11 and
JADC2-MAP §6 F1. Every ID there has an entry here. Duplicates are argued once, under the first
ID; the other IDs point to it and are ruled with it. Proposals in the open omos PRs (#550, #551,
#552, #553, #554) are mapped in [cross-reference.md](cross-reference.md), so nothing is asked twice.

| File | Entries |
|---|---|
| [bus.md](bus.md) | BUS-P1 … BUS-P11 |
| [ontology.md](ontology.md) | ONT-P1 … ONT-P14 |
| [reference-architecture.md](reference-architecture.md) | RA1 … RA4 |
| [cot-projection.md](cot-projection.md) | CP-1 … CP-4 |
| [f1-intercept-example.md](f1-intercept-example.md) | F1 |
| [cross-reference.md](cross-reference.md) | A4-P*, A0-*, POL-*, H2 P*, A7 P* against this register |

## Rule in this order (what blocks activation)

"Activation" means syncing a service or stream that an open PR has built, or turning on
behaviour it ships switched off. Within a tier, the order is the order to rule.

| Tier | ID | Decision asked | Blocks | Duplicates ruled with it |
|---|---|---|---|---|
| **1. Blocks syncing a built PR** | ONT-P1 | Approve the `ENT` stream CR | #552 sync (A4-P1) | — |
| | BUS-P2 | Asset commands are `cmd.fleet.{class}.{id}.{verb}`; mission verb `cancel` | #554 and #550 command path | ONT-P4 |
| | ONT-P14 | Policy RPC names `svc.policy.{command\|task}.{propose\|confirm}`, `svc.policy.action.get`, `svc.policy.check` | #554 sync | RA2 (naming half) |
| | ONT-P10 | One audit subject `evt.so.action.{event}` | #554 sync | RA2 (audit half) |
| | ONT-P7 | `CMD` 365 d, 512 MiB, `discard: new` | #554 fail-closed audit (POL-10) | — |
| | BUS-P1 | Telemetry token two is a class; fleet uses doctrine §9.1; `tel.feed.>` in `TEL` | #550, #551 status, #553 heartbeat binding (A0-2) | ONT-P5, ONT-P6, RA4 (form) |
| | CP-4 | `tak.cot.mesh.rx`/`.tx` are the mesh pair; grammar line "CoT to and from TAK" | #551 at a site | — |
| **2. Ratifies what open PRs already publish** | F1 | Replace the `cmd.drone.x47.intercept` example; forbid the verb | none (safety; doc only) | — |
| | BUS-P9 | Dead letters on `evt.sys.{source}.rejected` | none (`SYS` stores it today) | — |
| | ONT-P9 | `cmd` domains `entity`, `team`, `watch`, `feed` | #554 `cmd.entity.>` | — |
| | ONT-P8 | `entity_type` vocabulary (and `satellite`) | #552 (A4-P4) | — |
| | BUS-P3 | Live 5–6-token `intel.>` subjects kept; new collectors four tokens | #553 (A0-3) | — |
| | CP-3 | `tak.cot.mesh.tx` senders use core publish, no reply | #551 | — |
| | CP-2 | Gateway skips OneMind's own CoT when `MESH_TO_TAK_SERVER=true` | #551 (flag stays `false`) | — |
| | BUS-P7 | Any metric name in `tel.fleet.{class}.{id}.{metric}` | K5 KLV extractor (no PR) | RA4 (metric) |
| **3. Unlocks behaviour shipped off** | BUS-P5 | Domain `entity`; `evt.entity.>` in `EVT` | #552 lifecycle events; #554 closing `cmd.entity.>` actions (POL-14) | ONT-P2 |
| | ONT-P11 | Asset tasks in the `mission` domain, never `evt.task.>` | #554 tasks; H4 task service | — |
| | ONT-P13 | `evt.vehicle.>` in `EVT` | #554 vehicle outcomes (POL-4) | — |
| **4. Blocks later features (no open PR)** | ONT-P3 | `DET` stream CR | detectors (H8, Remote ID); #552 skips `DET` while absent | — |
| | RA1 | `evt.video.>` in `EVT` | video registry (K lane) | — |
| | BUS-P10 | Shrink, then retire `GEO`; no new `evt.tak.cot.{uid}` publisher | stream caps PR (A22) | — |
| | BUS-P6 | World classes `rail`, `transit` | transit collector | — |
| | CP-1 | `cmd.projector.{instance}.{disable\|enable}` | projector rework (A6) | recorded with BUS-P4 |
| | BUS-P4 | Out-of-grammar legacy subjects recorded, not extended | none | CP-1 |
| | BUS-P11 | Doctrine §4: buffering needs a site stream | site leaf (A23) | — |
| **5. Information only** | BUS-P8, ONT-P12, RA3 | Noted, no ruling needed | — | — |

Tier 1 holds every stream or naming decision without which a built service cannot be synced
safely; nothing in tiers 2–5 has to wait for tier 1. F1 is listed first in tier 2 because it is
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

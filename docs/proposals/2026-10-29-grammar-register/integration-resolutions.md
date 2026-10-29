# Integration-pass resolutions INT-R1 … INT-R4

Status: **proposed, for the owner**. Date: 2026-10-29. Filed under onemind-xr decision 0004
item 8. Wiring the M4 PRs together (onemind-xr `docs/STATUS.md`, M4 row) found five cross-PR
blockers. Four of them were proposed several times, in different words, by different PRs. Each
entry below combines those proposals into one ruling and names the PR that carries the change
(the fifth blocker, `zone.*` in the policy catalogue, is D1-P4, now carried by #554 `f4519b17`). Ruling
an INT entry also rules every proposal listed in "Rules with". None of them adds a subject
pattern.

The code that applies a resolution lives in the carrying PR, not here. Status on 2026-10-29
(re-checked after the fix lanes pushed):

| ID | Carrying commit | State |
|---|---|---|
| INT-R1 | #554 `f4519b17` (POL-16 = D2-P1, shared vectors); #550 `2b724093`; #566 `44de8068` | on the PR heads |
| INT-R2 | #552 `73823045` | on the PR head |
| INT-R3 | #562: the A19 change (A19-P11) is **not pushed yet**; its working copy differs from the table below (see "Open contradiction") | open |
| INT-R4 | #563: the projector's `svc.cot.send` is **not pushed yet**; its working copy takes two bodies, not one (see "Open contradiction") | open |

Check the PR before you tick a box. The fix lanes' evidence is in onemind-xr
`docs/evidence/INT/` (branch `run/m0`).

| ID | Resolution | Rules with | Carried by | Tier |
|---|---|---|---|---|
| INT-R1 | A command carries the `params_hash` of the action that was approved; bridges recompute only without `task_id` | D2-P1 (= POL-16 in #554), H3-P2 | #554, #550, #566 (clients test against the shared vectors) | 1 |
| INT-R2 | The entity service fills `capabilities` on `sim-*` fleet entities from the bridge status, bounded by fleet class | WEB-P5, H3-P4 | #552 (entity) | 1 |
| INT-R3 | One consolidated `CLIENT` import list (stream imports are an open owner choice) | WEB-P2, UE-P3, I2-P2, A5-P7, D2-P5 (exports), K-P7, A19-P11 | #562 (identities) | 1 |
| INT-R4 | The projector serves `svc.cot.send` (body form is an open owner choice: one body or two) | WEB-P4, UE-P1, B-P4, I1-P2 (exit), I1-P1 (publisher) | #563 (projector) | 3 |

---

## INT-R1. `params_hash` canonical form

**Problem.** #554 stamped the fleet command of a `task.assign` with the hash of the
`task.assign` action. #550 recomputed the hash as `asset.{verb}` over the command `data`, so it
refused every command that a task dispatched (`not_authorized`). This was reproduced end to end
in onemind-xr `docs/evidence/M4/D2/e2e-unpatched.txt`. Three clients (web, Unreal, Hermes) and
two bridges (MAVLink, ROS 2) each compute the hash, and the only shared reference was #554's
code.

**Resolution (as carried by #554 `f4519b17`, #550 `2b724093`, #566 `44de8068`).**
1. `params_hash` is `sha256:` followed by the lowercase hex SHA-256 of the RFC 8785 (JCS)
   canonical JSON of `{action_type, action_version, target, params}`. This is the formula in
   onemind-xr `contracts/asset-command.md` §6.1, unchanged.
2. A command carries the hash of the action that was approved (asset-command.md §4, §5, §6.3;
   mission-task.md §6 step 6). Without `task_id`, that action is `asset.{verb}` over the
   command `data`, and the bridge recomputes the hash. With `task_id`, it is `task.assign` or
   `task.abort` over the task parameters. The bridge cannot recompute that hash, so it checks
   the `tsk-` id and the `sha256:` form and relies on the Ed25519 policy signature, which
   covers `data` and `approval`. An unsigned task-bound command is refused when signatures are
   required, and logged with a warning otherwise. The policy service is unchanged, and #568's
   `services/mission/test/interop/policy-task-fleet-params-hash.patch` is not needed.
   This is POL-16 (D2-P1) in #554's
   `docs/proposals/2026-10-29-policy-service-owner-proposals.md`.
3. The reference implementation is #554 `services/policy/internal/policy/jcs.go`. The one
   cross-implementation vector file is
   `services/schemas/command/testdata/params-hash-vectors.json`, carried identically by #554,
   #550 and #566 and tested by all three (H3-P2). The clients' golden values (onemind-xr
   `docs/evidence/M4/UE/jcs-golden.txt`, `docs/evidence/M4/WEB/jcs-vectors.txt`) match it.

**Not chosen.** Every command on `cmd.fleet.>`, including a task-dispatched one, carrying the
`asset.{verb}` hash over its own `data` (the #568 patch applied to #554). That would let a
bridge recompute every hash without the signature, but the command would no longer carry the
hash the human approved, which §6.3 requires. To choose it instead, amend §6.3 to name a
per-command hash, apply the #568 patch to #554, and revert the task-bound branch in #550 and
#566.

- [ ] Approve as carried (approved action's own hash; signature for task-bound commands)
- [ ] Choose the per-command `asset.{verb}` hash instead (amend §6.3)
- [ ] Reject

## INT-R2. `capabilities` on fleet entities

**Problem.** The policy service denies an asset verb (`not_capable`) unless the target entity
lists the matching task type. The entity service (#552) never set `capabilities`, so in the
cluster every `goto` would be denied. The web, Unreal, ROS 2 and mission e2e runs each used an
entity stand-in to get past this. Meanwhile the bridges already publish the verbs they accept:
#550 lists `commands_accepted` in its status `meta`, and #566 does the same.

**Resolution (as carried by #552 `73823045`).** The entity service maps the verbs in the
latest `tel.fleet.{class}.{id}.status` `meta.commands_accepted` to
`capabilities: ["onemind.task.{verb}.v1", …]` on that fleet entity, with two limits that make
it narrower than a plain mapping:
- **`sim-*` only.** Only simulated fleet ids get capabilities (stop line, decision 0004 #13).
  A non-sim fleet id, a feed, and a CoT or TAK record never get any, and a record that stops
  being simulated, `sim-*` and taskable loses them.
- **Bounded by class.** `land` and `orbit` are air verbs (drones only); `observe` needs an
  observable asset, which nothing is today. A robot or vehicle gets `goto`, `hold` and `rtl` at
  most. `estop` needs no capability (it is unblockable), and the other verbs are not taskable.

It never invents a capability the bridge did not list, and it gives no class defaults. With no
status, the entity gets no capabilities, so the policy service fails closed. A later registry
entry may narrow the list but may not widen it. The policy service is unchanged. Clients offer
a verb only when the entity lists it (the WEB follow-up). The same commit reads `meta.kind` for
drones and robots (H3-P4) and takes `health.connection` from the bridge's `meta.connection`.

**Not chosen.** Making the policy service read `commands_accepted` directly. That would give
the policy service a second source of entity facts, against the single-writer rule
(REFERENCE-ARCHITECTURE R3).

- [ ] Approve as carried
- [ ] Approve with changes:
- [ ] Reject

## INT-R3. One `CLIENT` import list

**Problem.** #562 gives a `CLIENT` user requests on `svc.policy.>`, `svc.entity.>`,
`svc.layer.>`, `svc.cot.send` and `agents.prompt.hermes.>`, plus the `cmd.map.>` stream from
`AGENTS`. It gives nothing else. Five lanes each asked for more, separately: WEB-P2, UE-P3,
I2-P2, A5-P7 and D2-P5. K-P7 asked for the video registry, which BUS-DESIGN §7.1 lists as
`svc.registry.>` but #562 leaves out.

**Resolution.** The list below goes into #562's `CLIENT` and `AGENTS` blocks as one change.
The services part matches the A19 working copy (A19-P11, not pushed yet); the stream imports are
where the two disagree (next section).

| Import into `CLIENT` | Kind | From | Asked by | Note |
|---|---|---|---|---|
| `svc.layer.>` | service, `response_type: stream`, `response_threshold: 10m`, `share: true` | `ONEMIND` | A5-P7 | already on #562's head |
| `svc.registry.video.>` | service, `share: true` | `ONEMIND` | K-P7 | video registry reads; narrower than BUS-DESIGN §7.1's `svc.registry.>`, as A19-P11 has it |
| `svc.mission.{get\|list\|report}` | service, `share: true` | `ONEMIND` | D2-P5 | `svc.mission.apply` stays policy-only; A19-P11 also gives `AGENTS` `get` and `list` |
| `svc.geofence.alerts` | service, `share: true` | `ONEMIND` | D1-P3 | open-alert list |
| `cmd.map.>` | **service** with `response_type: stream` (`response_threshold: 5s` in A19-P11), exported by `CLIENT` to `AGENTS` | `CLIENT` → `AGENTS` | I2-P2 | replaces the stream export, so client acks reach Hermes |
| `svc.tak.mission.>` | service, `share: true` | `ONEMIND` | UE-P2, UE-P3 | A19-P11 adds it now (no responder until tak-client serves UE-P2); the register had it waiting for UE-P2 |

**Open contradiction: stream imports (A19-P11 against this register).** The A19 working copy
in `platform/nats/identities/identities.yaml` adds these read-only stream imports into
`CLIENT`. Clients subscribe to them and cannot publish into `ONEMIND` through them.

| Stream import in A19-P11 | For | Register before this change |
|---|---|---|
| `evt.so.action.>` | live action status and audit (D3) | poll `svc.policy.action.get` until ruled |
| `evt.mission.>` | mission and task facts (D2, C7) | `evt.mission.*.changed` waits for UE-P2 |
| `evt.map.geofence.>` | zone crossings (D1) | not listed |
| `tel.feed.*.status` | feed freshness (B20) | `tel.feed.>` |
| `tak.cot.out`, `tak.cot.mesh.rx` | GeoChat and CoT from TAK and the mesh (C8) | open choice against a `svc.cot.chat` history endpoint |

`tel.feed.*.status` is narrower than `tel.feed.>` and carries only heartbeats, so the register
takes A19-P11's form. The exposure risk applies to the rest. A stream import shows every
user's actions, every mission and every zone crossing to every `CLIENT` user, and per-user
filtering on an import is unverified (WEB-P2 (a)). The two TAK imports show all TAK CoT,
including PLI from EUDs, to every authenticated human. The owner chooses for each row.

- [ ] Approve the services table
- [ ] Stream imports: ☐ all of A19-P11 ☐ only `tel.feed.*.status`, `evt.map.geofence.>` ☐ none (poll `svc.policy.action.get`)
- [ ] Chat receive: ☐ stream import of `tak.cot.out`, `tak.cot.mesh.rx` (A19-P11) ☐ `svc.cot.chat` from the projector (WEB-P2 (c))
- [ ] Reject

## INT-R4. `svc.cot.send`

**Problem.** onemind-xr `contracts/cot-projection.md` §9 names `svc.cot.send` and leaves its
body to FEATURES C2. #562 already imports it into `CLIENT`, but no service on a PR head serves
it. Three things fill the gap differently. The web client sends a structured GeoChat body
(WEB-P4). The Unreal client sends client-built CoT XML in the body of
`client/OneMind/Docs/svc-cot-send.md` (B-P4), now including shapes (UE-P1). Agentops publishes
CoT itself on the TAK outbound subjects (I1-P1, I1-P2).

**Common ground.** The projector (#563) serves `svc.cot.send` and is the only service that
publishes outbound CoT, so the marking gate G7, the kill switch G1 and the echo convention
apply to everything sent. It takes the caller's identity from the server-set
`Nats-Request-Info` header of the `CLIENT` import and refuses a same-account caller. Output goes
on `tak.cot.inject` for `server` and `tak.cot.mesh.tx` for `mesh`, with a core publish and no
reply (CP-3, B-P3). Error codes are `invalid`, `forbidden`, `rate_limited` and
`target_failed`. When the endpoint is live, agentops moves its contact PLI and its GeoChat
replies onto it and loses its publish grant on the two TAK subjects (I1-P2 exit). The client's
own TAK stream is a separate decision (UE-P5).

**Open contradiction: one body or two.** The register proposed one body, with kinds `geochat`
and `shape`, the client never sending XML, and one error shape. The COT lane's working copy
(`services/projector/internal/projector/cotsend.go`, not pushed to #563 yet; evidence in
onemind-xr `docs/evidence/INT/COT/evidence.log`) takes two:

| | CoT form (Unreal, B-P4, UE-P1) | Structured GeoChat (web, WEB-P4) |
|---|---|---|
| Request | `{v, kind: marker\|spi\|route\|geofence\|shape\|chat, targets[], client_id, client_ref, marking{tlp, releasable_to}, dest_uid, cot}` | `{v, kind: "geochat", target: server\|mesh\|both, message_uid, room{id, name, kind}, text, sender{uid, callsign}, marking?}` |
| Who builds the XML | the client; the projector re-checks all of it against an allowlist (reserved `OM.`/`INTEL.` uids, uid bound to the client, kind-to-type table, `<detail>` allowlist, no person-track data, time window, size) and publishes the checked bytes | the projector |
| Error reply | `{v, ok: false, error: "<code>", detail}` | `{v, ok: false, error: {code, message}}` |
| Idempotency | `client_ref` | `message_uid` |

Two bodies keep both built clients working unchanged. The cost is a CoT parser on a
client-facing endpoint and two error shapes. One body removes both but needs the Unreal
client to stop sending XML.

- [ ] Two bodies as in the COT working copy (adopt B-P4 and WEB-P4; one error shape later)
- [ ] One body, projector builds all CoT (kinds `geochat` and `shape`; Unreal client changes)
- [ ] Reject

```subjects
filter    cmd.fleet.>
filter    evt.so.action.>
current   svc.policy.action.get
filter    svc.policy.>
filter    svc.entity.>
filter    svc.layer.>
filter    svc.registry.video.>
alternative  svc.registry.>
current   tel.fleet.{class}.{id}.status
proposed  svc.mission.{get|list|report}
current   svc.mission.apply
proposed  svc.geofence.alerts
filter    tel.feed.*.status
filter    tel.feed.>
filter    evt.map.geofence.>
filter    evt.mission.>
filter    cmd.map.>
filter    agents.prompt.hermes.>
proposed  svc.cot.send
alternative  svc.cot.chat
current   tak.cot.out
current   tak.cot.mesh.rx
current   tak.cot.inject
current   tak.cot.mesh.tx
filter    svc.tak.mission.>
filter    evt.mission.*.changed
action    task.assign
action    task.abort
action    asset.{verb}
```

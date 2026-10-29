# CoT projection proposals CP-1 … CP-4

Source: onemind-xr `docs/contracts/cot-projection.md` §11 (branch `run/m0`). Index and ruling
order: [README.md](README.md). #551 (TAK mesh gateway) files CP-4 as its P1 and CP-2 as part of
its P7; they are ruled here once.

---

## CP-1. Projector pause and resume subjects

**Proposed.** `cmd.projector.{instance}.{disable|enable}`, with `{instance}` = `hub` or the site
id. The policy service publishes it after an approved `projector.disable` or `projector.enable`
action (an A2 action in onemind-xr `docs/contracts/asset-command.md` §1 "other actions").

**Why the grammar does not cover it.** The live projector subscribes three-token
`cmd.projector.disable` and `cmd.projector.enable` (`platform/omos/tak-projector.yaml` lines
27–28 and 412), outside `cmd.{domain}.{entity}.{action}` (`CLAUDE.md` line 88). With a hub and
site projectors, a pause also needs a target.

**Alternatives considered.** Keep three tokens and pause every projector at once (no per-site
control); a KV flag the projector watches (state, but no `CMD` audit receipt). Rejected.

**Interim reading.** The live three-token subjects stay, recorded with the BUS-P4 legacy list.

**Depends on it.** Projector rework (FEATURES A6) and the policy catalogue (A2). The #554
catalogue has no projector action yet. No open PR publishes either form.

```subjects
proposed  cmd.projector.{instance}.{disable|enable}
proposed  cmd.projector.hub.disable
legacy    cmd.projector.disable
current   cmd.{domain}.{entity}.{action}
legacy    cmd.projector.enable
action    projector.disable
action    projector.enable
```

- [ ] Approve as proposed
- [ ] Approve with changes:
- [ ] Reject

---

## CP-2. The gateway must not forward OneMind's own CoT to TAK Server

**Proposed.** When `MESH_TO_TAK_SERVER=true`, the gateway copies mesh traffic to
`tak.cot.inject` but skips uids `OM.*` (projected entities), `INTEL.*` and its own presence SA,
so OneMind markers do not loop back into TAK Server as if a device had sent them. The flag stays
`false` everywhere (decision 0004 #6).

**Why.** Not a subject gap: a loop in the data path the grammar allows
(cot-projection.md §7).

**Alternatives considered.** Rely on TAK Server deduplication by uid (the server would still
rebroadcast a stale copy); never allow the flag. The first was rejected; the second is the
interim.

**Interim reading.** Flag `false`. #551 already implements the skip list (its P7: with `true`,
`OM.*`, `INTEL.*` and `omos-gw-*` uids are never copied) and tests it.

**Depends on it.** #551 (P7). Same as: #551 P7 (skip list part).

```subjects
current   tak.cot.inject
proposed  tak.cot.mesh.rx
```

- [ ] Approve as proposed
- [ ] Approve with changes:
- [ ] Reject

---

## CP-3. `tak.cot.mesh.tx` sender semantics

**Proposed.** Senders (only the site projector) use a core publish with no reply subject.
Delivery is confirmed by the gateway's `tx` counter on `tel.fleet.edge.{id}.status`, not by a
reply.

**Why.** The gateway answers any message that has a reply subject with `ok` or an error
(#551 `services/bridge/bridges/tak/main.go` lines 10 and 380–381). `SITE_COT` (and `COT` at the
hub) captures `tak.cot.mesh.tx`, and JetStream also answers a stored message's reply subject
with its PubAck, so a request would get two racing replies on one inbox (reasoned from how
JetStream acknowledges, not run; onemind-xr `docs/design/BUS-DESIGN.md` §12). BUS-DESIGN §3.3 is
aligned with this.

**Alternatives considered.** JetStream publish with PubAck as the receipt (proves storage, not
delivery to the mesh); request/reply to the gateway on a subject no stream captures (a new
subject). Rejected.

**Interim reading.** As proposed.

**Depends on it.** #551 (gateway), site projector (A6 at a site).

```subjects
proposed  tak.cot.mesh.tx
current   tel.fleet.edge.{id}.status
```

- [ ] Approve as proposed
- [ ] Approve with changes:
- [ ] Reject

---

## CP-4. The mesh pair under `tak.cot.>`

**Proposed.** Ratify `tak.cot.mesh.rx` (mesh to bus) and `tak.cot.mesh.tx` (bus to mesh) as the
mesh pair. Change the `CLAUDE.md` grammar line 94 description from "CoT from TAK (cot-tap
output)" to "CoT to and from TAK", and state that `COT` at the hub and `SITE_COT` at a site
store both directions under the `COT` caps in BUS-DESIGN §4.1 (72 h, 3 GiB).

**Why the grammar does not cover it.** Line 94 defines `tak.cot.>` as "CoT from TAK".
`tak.cot.mesh.tx` carries CoT toward TAK devices (the inject direction), and `COT` captures all
of `tak.cot.>` (`platform/nats/streams/cot.yaml` line 11), so it stores mesh output too. #551
already publishes both subjects and adds a paragraph about them to `docs/NATS-DOCTRINE.md` §9.8.

**Alternatives considered.** Put the bus-to-mesh direction beside `tak.cot.inject` as
`tak.cot.inject.mesh` (inject keeps its meaning, but `tak.cot.inject` is the TAK Server plugin
input and the plugin subscription must not see mesh traffic); a separate `mesh.>` prefix (new
pattern). Rejected.

**Interim reading.** The subjects as #551 publishes them.

**Depends on it.** #551 (P1). Site template (A23) for `SITE_COT`.

```subjects
proposed  tak.cot.mesh.rx
proposed  tak.cot.mesh.tx
filter    tak.cot.>
current   tak.cot.inject
alternative  tak.cot.inject.mesh
alternative  mesh.cot.tx
```

- [ ] Approve as proposed (rules #551 P1)
- [ ] Approve with changes:
- [ ] Reject

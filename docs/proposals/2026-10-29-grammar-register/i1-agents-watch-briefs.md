# Agent operations proposals I1-P1 … I1-P7 (FEATURES I1, I3, D10)

Status: **proposed, for the owner**. Date: 2026-10-29. Filed under onemind-xr decision 0004
item 8. This file belongs to the grammar register of PR #556
(`docs/proposals/2026-10-29-grammar-register/`); it adds entries by reference and re-argues
nothing the register already holds. When #556 merges, add one row per entry to its
`README.md` index (tier given below). Nothing here changes a live subject, stream,
Application or workload. The service is `services/agentops` (this PR).

| ID | Decision asked | Tier (register README) | Ruled with |
|---|---|---|---|
| I1-P1 | Agent contacts in TAK: an exception to cot-projection §4.1 (agents "never" projected) | 3 (unlocks behaviour shipped off) | — |
| I1-P2 | agentops may publish GeoChat replies (`b-t-f` only) until the projector serves the ephemeral-send endpoint | 3 | CP-3 (core publish, no reply) |
| I1-P3 | Alerts are `evt.entity.{entity_id}.alert`; a feed's alert id is `feed-{source}` | 3 | **BUS-P5** (stores them) |
| I1-P4 | AOI briefs on `evt.so.brief.published`, reads on `svc.brief.generate` and `svc.brief.latest` | 2 (ratifies what this PR publishes) | ONT-P10 reasoning (the so domain is stored by `EVT`) |
| I1-P5 | Agents are fleet class `agent` with uid `agent-{agent}-{owner}-{name}`; two small entity-service changes | 2 | BUS-P1, ONT-P8 |
| I1-P6 | NATS user `agentops` (identities file of PR #562) and one OpenFGA tuple | 1 for sync of this PR | — |
| I1-P7 | Watch definitions stay in Git (ConfigMap) until the `watch` command domain is ruled | 5 (information) | **ONT-P9** |

---

## I1-P1. Agent contacts in TAK

**Proposed.** FEATURES I1 says "Hermes shows as a contact" in TAK, but onemind-xr
`docs/contracts/cot-projection.md` §4.1 row `person, team, agent, node, station` says agents
are **never** projected (gates G4, G6), and §3.1 sets `indicators.egress_tak: false` for them.
Allow one exception: a reachable agent is announced to TAK as a contact PLI with uid
`OM.agent.{entity_id}` (the `OM.` prefix is the projector's echo convention, so the entity
service and tak-client drop it on the way back), type `a-f-G-U-C`, `<contact callsign
endpoint="*:-1:stcp">`, `<__group>`, `<takv>` and its status and current task in `<remarks>`
(`services/agentops/internal/cot/cot.go` `RenderContact`). The point is the agent's anchor with
`ce`/`le` unknown, so nobody reads it as a fix.

**Who publishes it.** Long term the projector, through its ephemeral-send endpoint
(cot-projection §9). Until that exists, agentops publishes it on the two TAK outbound
subjects when `CONTACT_TARGETS` names them (`server`, `mesh`); the manifest ships it empty.

**Unverified.** That ATAK lists such a PLI as a messageable contact when it arrives through TAK
Server rather than the mesh, and how it draws `a-f-G-U-C` for a non-human. The owner check is in
onemind-xr `docs/evidence/M4/I1/OWNER-STEPS.md` step 3.

```subjects
current   tak.cot.inject
current   tak.cot.mesh.tx
alternative  svc.cot.send
```

## I1-P2. GeoChat replies from agents

**Proposed.** The ATAK-to-agent path (`services/agentops/internal/chat`) reads GeoChat from
`tak.cot.out` (tak-client, PR #560) and `tak.cot.mesh.rx` (mesh gateway, PR #551), prompts the
addressed agent on its Synadia prompt subject, and answers with a directed GeoChat on the
matching outbound subject: `tak.cot.inject` with `<marti><dest callsign>` for TAK Server,
`tak.cot.mesh.tx` with header `Tak-Dest-Uid` for the mesh, core publish without reply (CP-3).
Security-and-trust T4 makes the projector the only publisher of those subjects; this asks for
a second publisher, used by agentops only for `b-t-f` replies and the I1-P1 contact PLI, until
the projector serves `svc.cot.send`, after which agentops moves to it and loses the grant. NATS
permissions cannot restrict by CoT type, so the grant itself is subject-wide; that is why this
needs a ruling rather than a config change.

Inbound text is untrusted: the sender's TAK uid and callsign are claimed, not verified, and are
passed to the agent as `trust: claimed`; per-sender (6/min) and global (60/min) limits; chats
older than 2 min are dropped, not answered late; the agent's own echoes are ignored; a §7
mid-stream query (an approval prompt) is never answered from TAK. The manifest ships
`REPLY_TARGETS` empty.

```subjects
current   tak.cot.out
current   tak.cot.mesh.rx
current   agents.prompt.{agent}.{owner}.{name}
```

## I1-P3. Alert subject and feed ids

**Proposed.** Watch jobs publish each alert on `evt.entity.{entity_id}.alert`, the `alert`
event that entity.md §1 and BUS-DESIGN §3.2 already list for the `entity` event domain, with a
second revision of the same `alert_id` carrying the model's explanation and the proposal
outcome (`Nats-Msg-Id: {alert_id}:{revision}`). Geofence crossings also go on doctrine §9.1
`evt.map.geofence.{id}.{enter|exit}` (stored by `MISSION` today). A feed is not an entity, so a
feed outage uses the id `feed-{source}` (ONTOLOGY §3 object type `feed`, keyed by source).

**Stored by.** `evt.entity.>` in `EVT` is **BUS-P5**; rule it there. Until then agentops
publishes alerts with a core publish (it asks the server whether a stream binds the subject
first), so they reach live subscribers but are not stored.

**Follow-up for PR #552 (not a grammar question).** ONTOLOGY §3 makes alerts entities
(`ent.alert.{entity_id}.state`); the entity service would need an adapter for
`evt.entity.*.alert`. Not built in either PR.

```subjects
proposed  evt.entity.{entity_id}.alert
current   evt.map.geofence.{id}.{enter|exit}
filter    evt.entity.>
filter    evt.entity.*.alert
current   ent.alert.{entity_id}.state
```

## I1-P4. Briefs

**Proposed.** Each AOI brief is published on `evt.so.brief.published` (the `so` domain is
stored by `EVT` today through `evt.so.>`, the reasoning BUS-DESIGN §3.3 used for
`evt.so.action.*`), `Nats-Msg-Id` = `brief_id`. Reads are request/reply on
`svc.brief.generate` (`{aoi_id}` or an inline `{aoi}`) and `svc.brief.latest` (`{aoi_id}`),
both inside the grammar line svc.{service}.{endpoint}. A brief reads `svc.entity.query` (PR #552) and the
feed heartbeats `tel.feed.{source}.status` (PR #553; stored once `TEL` binds `tel.feed.>`,
BUS-P1).

**Alternative considered.** A new domain, `evt.brief.{aoi}.published`: no stream binds it, and
it would need an `EVT` filter for one publisher. Rejected for now.

```subjects
proposed  evt.so.brief.published
filter    evt.so.>
current   evt.so.action.*
proposed  svc.brief.generate
proposed  svc.brief.latest
current   svc.entity.query
current   tel.feed.{source}.status
filter    tel.feed.>
alternative  evt.brief.{aoi}.published
```

## I1-P5. Agents as fleet class `agent`

**Proposed.** agentops follows Synadia heartbeats on `agents.hb.{agent}.{owner}.{name}`
(stored by `AGENTS`) and publishes each live agent as a doctrine §9.2 geo-envelope on
`tel.fleet.agent.{id}.position` (doctrine §9.1 class `agent`, §9.3 type `fleet.agent`) with
uid `agent-{agent}-{owner}-{name}`. The entity service turns it into
`ent.agent.{entity_id}.state`; proved against the PR #552 binary
(`services/agentops/internal/e2e`). An agent may add an optional `onemind` block to its
heartbeat (`mode`, `activity`, `current_task`, `lat`/`lon`); protocol §8.3 lets receivers
ignore unknown fields, so other callers are unaffected.

Two small changes are asked of PR #552 (not grammar):
1. Map a fleet `agent` uid to `entity_id` unchanged (today any non-`sim-` fleet uid becomes
   `fleet-{uid}`; agentops' manifest works around it with `ENTITY_LINKS`
   `{"fleet_uid:agent-…": "agent-…"}`).
2. Copy `meta.current_task` into the contract field `current_task` and `meta.activity` into
   `status.activity` (today both reach clients only inside `meta`).

Both items belong to PR #552 (`feat/entity-service`, `services/entity`) and are tracked there:
until #552 adopts them, agent entity ids need the `ENTITY_LINKS` workaround and Hermes'
current task and activity reach clients only in `meta.*`. I1 is therefore built and tested
here, but stays **awaiting owner rulings** (I1-P1, I1-P2, I1-P5) rather than done; ATAK to
Hermes messaging ships off (`REPLY_TARGETS` empty) until I1-P1 and I1-P2 are ruled.

```subjects
proposed  tel.fleet.agent.{id}.position
current   agents.hb.{agent}.{owner}.{name}
current   ent.agent.{entity_id}.state
```

## I1-P6. Identity and authorization for agentops

**Proposed.** Add this user to `platform/nats/identities/identities.yaml` (PR #562 layout), and
seal its password as `omos-agentops-nats` in namespace `omos`. Until then the Deployment uses
the shared `omos-nats-creds`, as the entity service does.

```yaml
  - user: agentops
    account: ONEMIND
    stage: A
    purpose: "omos-agentops (I1, I3, D10): agent presence, ATAK chat to agents, watch jobs, briefs; proposes, never confirms"
    secret: {namespace: omos, name: omos-agentops-nats, user_key: user, password_key: password}
    connection_types: [STANDARD]
    publish:
      - "tel.fleet.agent.>"
      - "evt.entity.*.alert"
      - "evt.map.geofence.>"
      - "evt.so.brief.published"
      - "agents.hb.onemind-watch.>"
      - "agents.hb.onemind-brief.>"
      - "agents.prompt.hermes.>"
      - "svc.entity.query"
      - "svc.policy.task.propose"
      - "tak.cot.inject"        # I1-P1/P2 only; drop when svc.cot.send exists
      - "tak.cot.mesh.tx"       # I1-P1/P2 only
    publish_deny:
      - "svc.policy.*.confirm"
      - "cmd.>"
    subscribe:
      - "svc.brief.>"
      - "agents.status.onemind-watch.>"
      - "agents.status.onemind-brief.>"
      - "tak.cot.out"
      - "tak.cot.mesh.rx"
      - "tel.feed.*.status"
      - "evt.entity.*.alert"
      - "_INBOX.>"
    allow_responses: true
    jetstream:
      info: true
      names: true             # StreamNameBySubject before each fact publish
      consume: [ENT, TEL, AGENTS]
      kv_read: [geofences]
```

OpenFGA (policy model `platform/security/openfga/onemind-actions.json`, PR #554): one tuple
`agent:onemind-watch/zeus/hub proposer action_type:task.assign`. The watch's human owner
(`user:zeus`) must also be `proposer` of `task.assign` and `operator` of each asset; the policy
service checks both (agent and on-behalf-of intersection). While agentops is in the policy
service's own account its proposals are recorded `actor_verified: false`; stage B of #562
moves the proposal identity into `AGENTS`.

```subjects
current   svc.policy.task.propose
action    task.assign
legacy    svc.policy.*.confirm
legacy    cmd.>
```

## I1-P7. Watch definitions

**Information.** Watch jobs are data (agents-that-act §4): today a YAML file in the
`omos-agentops` ConfigMap, validated at start (explicit inputs, caps, a closed list of
responses, a human owner). ONTOLOGY §6 names the action `watch.create`/`watch.delete` on a
`watch` command domain, which is **ONT-P9**; once ruled, watches can move to a KV bucket written
only through policy. No new subject is used now.

```subjects
alternative  cmd.watch.{id}.{create|delete}
```

# omos agentops

Agent operations for OneMind: FEATURES **I1** (agents are entities and contacts), **I3** (watch
jobs) and **D10** (AOI briefs with an offline fallback). One Go binary, four components, each
switchable with `COMPONENTS`. Deployment: `platform/omos/agentops/` (Argo Application
`platform-omos-agentops`, manual sync). Owner proposals:
`docs/proposals/2026-10-29-grammar-register/i1-agents-watch-briefs.md` (I1-P1 … I1-P7).

Design inputs (onemind-xr, branch `run/m0`): `docs/research/agents-that-act.md` §4, §8;
`docs/design/ONTOLOGY.md` §3 rows `agent`, `alert`, `feed`; `docs/contracts/entity.md`;
`docs/contracts/cot-projection.md` §4.1, §8.2; `docs/design/BUS-DESIGN.md` §3.

## What each component does

| Component | Reads | Writes | Rule |
|---|---|---|---|
| `presence` | `agents.hb.>` (`AGENTS`, last per subject, then live) | `tel.fleet.agent.{uid}.position` (`TEL`), doctrine §9.2 geo-envelope, type `fleet.agent` | Never writes `ent.*`: the entity service (PR #552) turns it into `ent.agent.{entity_id}.state` |
| `chat` | `tak.cot.out`, `tak.cot.mesh.rx` (GeoChat `b-t-f*` only) | prompt on `agents.prompt.{agent}.{owner}.{name}`; reply on `tak.cot.inject` / `tak.cot.mesh.tx` and agent contact PLIs, **only** if `REPLY_TARGETS` / `CONTACT_TARGETS` say so (I1-P1, I1-P2) | Sender identity is claimed, never trusted |
| `watch` | the inputs each watch declares: `ent.{type}.*.state` (`ENT`), KV `geofences`, `tel.feed.*.status` | `evt.map.geofence.{zone}.{enter\|exit}` (`MISSION`), `evt.entity.{id}.alert` (`EVT` once BUS-P5), one request to `svc.policy.task.propose` per allowed proposal | Detection is deterministic; the model explains and picks from a closed list; **never confirms** |
| `brief` | `svc.entity.query` (PR #552), `tel.feed.*.status` (PR #553), `evt.entity.*.alert` | `evt.so.brief.published` (`EVT`); serves `svc.brief.generate`, `svc.brief.latest` | Facts gathered by code; model output checked (citations, gaps) |

The watch runner and the briefer are agents themselves: each beats on
`agents.hb.onemind-watch.zeus.hub` / `agents.hb.onemind-brief.zeus.hub` with its current task,
so `presence` shows them on the map next to Hermes.

Facts are published with a JetStream ack and `Nats-Msg-Id` when a stream binds the subject;
when none does yet (owner-proposed bindings), with a core publish. The publisher asks the
server (`$JS.API.STREAM.NAMES`, cached a minute) before each new subject, so core subscribers
never get a duplicate and nothing waits for an ack that cannot come.

## I1: agents on the map and in TAK

```
Hermes ──agents.hb.hermes.zeus.legacy {…, "onemind":{mode, activity, current_task}}──► AGENTS
   presence ──tel.fleet.agent.agent-hermes-zeus-legacy.position──► TEL ──► entity service
   entity service ──ent.agent.agent-hermes-zeus-legacy.state──► ENT ──► layer service, clients, Hermes
ATAK ──GeoChat to "LEGACY"──► TAK Server ──tak-client──► tak.cot.out ─┐
ATAK ──GeoChat (mesh)──► mesh gateway ──► tak.cot.mesh.rx ────────────┤
   chat ──agents.prompt.hermes.zeus.legacy (Synadia §5, trust: claimed)──► Hermes
   chat ◄── response chunks + terminator (§6) ──
   chat ──b-t-f reply──► tak.cot.inject (marti dest) / tak.cot.mesh.tx (Tak-Dest-Uid) ──► ATAK
   chat ──contact PLI OM.agent.{uid}──► same subjects, every CONTACT_EVERY
```

- **Position.** An agent has none. It is drawn at the position its heartbeat reports, else
  `AGENT_ANCHORS[ref]`, else `SITE_LAT`/`SITE_LON`; `meta.position_kind` says which, and `ce`
  is 50 m so nobody takes it for a fix. `stale` is 3 × the heartbeat interval (protocol §8.2
  offline rule); an offline agent is published as `mode: offline` for 10 minutes, then dropped.
- **Entity id.** The uid is `agent-{agent}-{owner}-{name}`. PR #552 prefixes non-`sim-` fleet
  uids with `fleet-`; set `ENTITY_LINKS={"fleet_uid:agent-hermes-zeus-legacy":"agent-hermes-zeus-legacy"}`
  on the entity service until I1-P5 item 1 lands. Status and current task arrive in `meta`
  (`meta.mode`, `meta.activity`, `meta.current_task`); I1-P5 item 2 copies them into the
  contract fields.
- **Contact uid** is `OM.agent.{uid}`: the entity service and tak-client drop `OM.*` as the
  projector's echo, so the contact never comes back as a `person`.
- **Chat safety.** Text is framed as data from a claimed TAK identity; at most 1000 characters
  forwarded, 1000 (mesh 600) returned; 6 messages per sender and 60 overall per minute; chats
  older than 2 minutes dropped; the agent's own echo ignored; a protocol §7 query (approval
  prompt) is refused, because approvals are only given in the OneMind app.

## I3: watch file

```yaml
watches:
  - id: sim-keep-in                 # ^[a-z0-9][a-z0-9_-]{0,62}$
    owner: user:zeus                # the human proposals are made for (required)
    level: warning                  # advisory | caution | warning
    rule: {kind: zone, entity_types: [drone], entity_ids: ["sim-*"], zone: zone-home, alert: outside}
    inputs: [ent.drone.*.state, "kv:geofences"]   # exactly what the rule reads
    rate: {cooldown_s: 300, max_alerts_per_hour: 12, max_llm_per_hour: 6, max_proposals_per_hour: 3}
    on_trigger:
      explain: true
      propose: {task_types: [onemind.task.rtl.v1, onemind.task.hold.v1], assignee: about}
```

| Kind | Fires when | Fields |
|---|---|---|
| `zone` | an entity is outside (or inside) a KV `geofences` zone (polygon with holes, or circle); crossings also emit doctrine `evt.map.geofence.*` | `zone`, `alert` |
| `stale` | no update for `max_age_s`, or lifecycle `stale`/`lost`, or a final `live: false` record | `max_age_s` |
| `threshold` | a numeric field (dotted path) crosses `value`; stays active until `clear` | `field`, `op` (`lt le gt ge`), `value`, `clear` |
| `feed_stale` | a feed heartbeat says stale, failing, or stops (3 × its limit); expected feeds that never report | `sources`, `reasons` |

Load-time checks (the service refuses to start otherwise): unknown fields; inputs must cover
what the rule reads and nothing more, and may not be wildcards wider than one entity type;
`owner` must be `user:…`; rates within caps (cooldown ≥ 10 s, ≤ 600 alerts, ≤ 60 model calls,
≤ 20 proposals per hour); `propose.action_type` is only `task.assign`; task types only
`rtl`, `hold`, `manual`, `agent_request` (types whose spec needs no coordinates from a model);
stop-line words refused.

On a rising edge (after cooldown and hourly cap) the runner:
1. publishes the alert (revision 1, no model involved);
2. if `explain` or `propose`: asks the model chain for `{explanation, proposal|null}` with the
   alert, evidence and the allowed list as data;
3. sends at most one `task.assign` to `svc.policy.task.propose` with `requested_by {kind: agent,
   id: agent:onemind-watch/zeus/hub, on_behalf_of: <owner>}`, `request_id` and `corr` as UUIDs
   derived from the alert id (idempotent), and the hashes `prompt_sha256`/`tool_input_sha256`;
   a choice outside the list is recorded as `not_sent`;
4. publishes revision 2 of the alert with the explanation and the policy answer
   (`awaiting_confirmation` for an A3 asset task). A human confirms in the OneMind app; the
   runner has no confirm code path, and its NATS user (I1-P6) is denied
   `svc.policy.*.confirm` and `cmd.>`.

## D10: briefs and the offline fallback

`MODEL_PRIMARY`, then `MODEL_LOCAL`, then the built-in `stub`. Specs: `none`, `stub`,
`agent:agents.prompt.{agent}.{owner}.{name}` (Hermes over the bus), or
`openai:{base url}|{model}` (any OpenAI-compatible server: Ollama, llama.cpp `llama-server`).
No API key is ever read or sent, so a paid hosted API cannot be called; an OpenAI-compatible
URL must be loopback, private or cluster-internal unless `MODEL_ALLOW_REMOTE=true`.

A brief reads up to 400 entities in the AOI box (paging `svc.entity.query`), drops any outside
the box, ranks them (own simulated assets, alert subjects, agents and people, events, then
confidence and recency), offers the top 25 to the model with strings sanitised (no brackets or
markup), and checks the answer: ids in `[brackets]` that were not offered are replaced with
`[unverified]`; fewer than 3 citations or an unnamed degraded feed is repaired from the facts
(`validation.repaired`). `model.offline` is true when no remote model wrote it. Candidate local
models and their licences: [LOCAL-MODELS.md](LOCAL-MODELS.md).

## Configuration (environment)

| Variable | Default | Meaning |
|---|---|---|
| `NATS_URL`, `NATS_USER`, `NATS_PASSWORD` | `nats://127.0.0.1:4222` | connection (SealedSecret in the cluster) |
| `COMPONENTS` | `presence,chat,watch,brief` | which components run |
| `PUBLISH_MODE` | `auto` | `auto`, `jetstream` or `core` |
| `SITE_LAT`, `SITE_LON` | — | agent anchor; required by `presence` |
| `AGENT_ANCHORS`, `AGENT_CALLSIGNS` | — | JSON by agent ref `agent/owner/name` |
| `AGENTS_STREAM`, `ENT_STREAM`, `TEL_STREAM`, `ZONE_BUCKET` | `AGENTS`, `ENT`, `TEL`, `geofences` | inputs |
| `CHAT_AGENTS` | `hermes/zeus/legacy` | agents reachable from TAK |
| `REPLY_TARGETS`, `CONTACT_TARGETS` | empty | `server`, `mesh` or both (I1-P1, I1-P2) |
| `CONTACT_EVERY`, `AGENT_COT_TYPE`, `AGENT_TAK_TEAM`, `AGENT_TAK_ROLE` | `60s`, `a-f-G-U-C`, `Cyan`, `HQ` | contact PLI |
| `WATCH_FILE`, `WATCH_AGENT` | `/config/watches.yaml`, `onemind-watch/zeus/hub` | watch runner |
| `BRIEF_AOIS`, `BRIEF_EVERY`, `BRIEF_AGENT` | `[]`, `15m`, `onemind-brief/zeus/hub` | briefer |
| `EXPECTED_FEEDS`, `FEED_LAYERS` | `{}` | `{source: staleAfterSec}`, `{source: [layer]}` |
| `MODEL_PRIMARY`, `MODEL_LOCAL`, `MODEL_TIMEOUT`, `MODEL_ALLOW_REMOTE` | `none`, `none`, `60s`, `false` | model chain |
| `HTTP_ADDR` | `:8080` | `/healthz`, `/metrics` |

## Tests

```bash
cd services/agentops
go vet ./...
NATS_SERVER_BIN=$(command -v nats-server) go test -race -count=1 ./...
```

Integration tests start their own nats-server (free port, temporary store) with the streams of
the omos CRs plus the proposed bindings, and cover: the I3 acceptance (a watch job raises an
alert and proposes a task that waits for a human; no confirm, no `cmd.>`); a model that answers
outside the envelope or carries injected text gets nothing sent; the core fallback delivers
exactly once; the D10 acceptance (brief with the primary and local models unreachable, at least
3 cited entities, gaps named, stored in `EVT`) and a local model whose hallucinated id is
removed; the I1 path (heartbeat → `tel.fleet.agent.*`, contact PLI, GeoChat over the TAK Server
and mesh paths, echo, stale, foreign, rate limit, agent offline).

`internal/e2e` runs against the real services of the open PRs and skips without them:

```bash
# from omos worktrees of feat/entity-service and feat/policy-service
(cd services/entity && go build -o /tmp/entity-service ./cmd/entity-service)
(cd services/policy && go build -o /tmp/policy ./cmd/policy)
ENTITY_SERVICE_BIN=/tmp/entity-service POLICY_SERVICE_BIN=/tmp/policy \
POLICY_SCHEMA_DIR=<policy worktree>/services/schemas \
POLICY_MODEL_FILE=<policy worktree>/platform/security/openfga/onemind-actions.json \
go test -race -count=1 ./internal/e2e
```

## Not in this slice

- Alerts as entities (`ent.alert.*`): needs an entity-service adapter (I1-P3 follow-up).
- Watch definitions through policy and KV (ONT-P9, I1-P7); today a ConfigMap.
- Hermes answering chats and explaining alerts with a real model: Hermes' own configuration,
  checked by the owner (onemind-xr `docs/evidence/M4/I1/OWNER-STEPS.md`).
- Multi-replica watch runner (leader lease); one replica today, `Recreate`.

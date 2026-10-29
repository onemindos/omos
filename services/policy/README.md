# Policy and audit service (`services/policy`)

FEATURES I4 (onemind-xr). Every action that changes the world goes through this service:
request, policy check, human confirmation, signed command, executor, outcome event. It is the
only publisher of `cmd.*` (other than `cmd.map.>`, which clients and agents publish directly)
and the single writer of the audit subject `evt.so.action.*`.

Contracts it implements (onemind-xr, branch `run/m0`): `docs/contracts/asset-command.md`
(sections 1–9), `docs/contracts/mission-task.md` (sections 4–6), `docs/design/BUS-DESIGN.md`
(sections 3.3, 3.4 and 7), `docs/design/ONTOLOGY.md` (sections 6 and 7).

## Endpoints (NATS request/reply, grammar `svc.{service}.{endpoint}`)

| Subject | Who calls it | Body | Reply |
|---|---|---|---|
| `svc.policy.command.propose` | clients, Hermes, watch jobs | ActionRequest (asset-command.md §6.1) for every non-task action | `{action_id, status, level, params_hash, confirm_by, approvals_needed, denial}` |
| `svc.policy.task.propose` | same | ActionRequest for `task.*` and `mission.*` | same |
| `svc.policy.command.confirm` / `svc.policy.task.confirm` | **humans only** | `{v, action_id, params_hash, decision: approve\|reject, reason}` | `{action_id, status, approvals_needed, approvals_received, command_id, denial}` |
| `svc.policy.action.get` | clients, Hermes (own actions only) | `{v, action_id}` | `{action, denial}` |
| `svc.policy.check` | other services; clients and agents for the dry run only | `{v, user, relation, object, context?}` (callers in the policy service's own account only, so a client cannot map the permission graph) or `{v, action: ActionRequest}` (dry run as the caller, nothing stored or sent) | `{allowed, ...}` |

## Rules it enforces (each has a test in `internal/policy/*_test.go`)

| Rule | Code on failure |
|---|---|
| Request envelope and action parameters validate against `services/schemas/actions/{type}.v{n}.schema.json`; task specs against `services/schemas/tasks/`; unknown fields refused | `schema_invalid` |
| Closed verb allowlist (`Catalogue` in `catalogue.go`); anything else refused | `unknown_action` |
| Civilian stop line: `strike`, `engage`, `intercept`, `designate`, `attack`, `kill`, `weapon`, `hostile`, follow/track a person, face or identity recognition, payload release, CoT `t-k-*` are refused by name, in action types, task types and parameter values | `forbidden_action` |
| Orbit, observe or a task pointed at a `person` entity or a person-class `track` | `forbidden_action` |
| Free text (`reason`, task instructions, agent questions, any string parameter) that asks to follow, track, tail or chase a person, or names a strike-class action (`TextStopLine` in `catalogue.go`; defence in depth behind the closed verb list). The safety verbs `estop`, `rtl`, `land` and `hold` are exempt, so "bird strike, stop now" still stops the asset; the wording is recorded as `text_stop_line` in their audit events | `forbidden_action` |
| A5 actions (credentials, cluster, foundation) | `forbidden_action` |
| Commandable targets must be `sim-*` and not marked non-simulated | `not_simulated` |
| The asset advertises the matching task type in `capabilities` | `not_capable` |
| The asset is online (no command is ever held for an offline asset), the point is outside every keep-out and inside a keep-in zone (KV `geofences`); `estop` skips these so a stop is never blocked | `precondition_failed` (`detail`) |
| A keep-in or keep-out zone the service cannot read (corridor, MultiPolygon, malformed coordinates, missing kind) is kept as *unevaluable* and logged; movement verbs (`goto`, `land`, `orbit`, `arm`, `takeoff`, movement tasks) are denied for the assets it applies to until the shape is supported. `estop`, `rtl` and `hold` are not blocked | `precondition_failed` (`detail: geofence_unevaluable`) |
| OpenFGA: `proposer` on `action_type:{name}`; an agent needs a human `on_behalf_of` and both must be proposers (intersection); the human must be `operator` of the asset. Confirm: `approver` (step 1, plus `operator`) or `second_approver` (step 2) | `not_authorized` |
| OpenFGA, the entity service, the geofence watcher or the audit stream unreachable: deny (fail closed) | `policy_unavailable` |
| Levels A0/A1 run at once (logged); A2/A3 need one human; A4 two distinct humans; A5 never | `duplicate_confirmation` for the same human twice |
| Agents may propose; any agent identity calling confirm is refused (and the server also denies the publish) | `agent_cannot_confirm` |
| Confirmation needs a verified human identity: the server-set `Nats-Request-Info` header from a `share: true` service import (CLIENT account). Same-account callers are recorded `actor_verified: false` and cannot confirm | `not_authorized` |
| An approval is bound to one `action_id` and its `params_hash` (SHA-256 of the RFC 8785 JCS form of `{action_type, action_version, target, params}`); a hash from any other action, or the wrong endpoint family | `conflict` |
| An approval is used once: confirming an approved, dispatched or finished action | `replayed` |
| No confirmation within `confirm_by` (120 s): the action expires | `expired` |
| One requester may hold at most 10 actions awaiting confirmation (`estop` is never limited). The count is keyed on the verified identity; an unverified same-account caller is counted in its own bucket for the id it claims, so it cannot use up a real principal's slots | `precondition_failed` (`detail: rate_limited`) |
| Parameters are stored, shown, hashed and sent in one canonical (RFC 8785) form, so duplicate keys cannot show the confirmer one value and send another | — |
| `request_id` is an idempotency key per requester: the same requester with the same body gets the same action; with a different body, `conflict`. Another principal reusing the id goes through every check itself | `conflict` |
| Preconditions are checked again at dispatch | `precondition_failed` |
| Commands carry `expires_at` (30 s), `simulated`, the approval and an Ed25519 `sig` over the JCS envelope; they are published with a JetStream ack (the `CMD` PubAck is the record) | — |
| Task assignment publishes `cmd.mission.{tsk}.create` first and sends the asset command only after `evt.mission.{tsk}.created` | `no_ack` |
| No `acknowledged` from the bridge within 10 s; acknowledged but no `executed`/`failed` within 10 min | `no_ack` |
| Entity, mission and person/team/agent task commands close on the entity or task service result (`evt.entity.{id}.*` / `evt.mission.{id}.*` matched by `causation` = command id or `data.action_id`; `evt.mission.{tsk}.created` for a task); none within 30 s | `no_ack` |

Every proposal, denial, refused confirmation, approval, rejection, expiry, dispatch, bridge
outcome and entity/task result is published as one event on `evt.so.action.{event}` in the omos event envelope v1
(`Nats-Msg-Id` = envelope id; the `EVT` stream stores `evt.so.>`). Agent sessions add
`agent_session`, `model_alias`, `prompt_sha256` and `tool_input_sha256` (hashes only).
Denial and refused-confirmation events are rate limited to 60 per identity **and code** per
minute (the excess is logged, not published), so a client flooding bad requests cannot flood
`EVT`, and flooding one kind of bad request cannot hide a different attempt (a stop-line verb,
an agent confirm) from the audit. A refused confirmation is audited even when the action id
does not exist (`action_id: act-unknown` when it is not a valid id).

Concurrency: requests are handled concurrently (up to 128 in flight). No lock is held across
an OpenFGA check, an entity lookup, a publish or the wait for `evt.mission.{tsk}.created`;
each action has its own operation lock, so a slow dependency on one action never delays a
proposal or confirmation of another (an `estop` is never queued behind a task assignment).

## Run the tests

```
cd services/policy
go vet ./...
go test -race -count=1 ./...
```

The tests start an in-process `nats-server` (v2.14.8, the live 2.14 line) with three accounts
(ONEMIND, CLIENT, AGENTS) and per-identity permissions from
`docs/proposals/2026-10-29-nats-per-service-permissions.md`, random per-run passwords, the
CMD/EVT/MISSION streams and KV `geofences`, stub entity, bridge and task services, and the
OpenFGA model from `platform/security/openfga/onemind-actions.json` (in-process evaluator, plus
the HTTP client against a stub OpenFGA server and a closed port). The last check of
`TestPolicyService` is the whole-run invariant: no `cmd.fleet.>` message without a human
`evt.so.action.approved`.

## Red-team suite (FEATURES I6)

`internal/policy/redteam_test.go` runs 18 adversarial cases (RT01–RT18) against the same
in-process stack: prompt-injected proposals and an injection corpus, an agent confirming,
forged and replayed approvals, parameter swaps after approval, stop-line verbs with synonyms and
encodings, follow-a-person requests, non-`sim-` ids, expired approvals and commands, offline
assets, OpenFGA down, bus permissions, flooding, audit completeness with the whole-run
invariant, the forbidden-name check over `services/schemas` and the catalogue
(`forbidden.go`; pre-existing hits are counted per file and word in `forbiddenBaseline`, so any
extra occurrence fails; compounds such as `airstrike` match), safety verbs whose reason uses
stop-line wording (RT17) and a quota flood under a claimed identity (RT18). Each case is mapped
to onemind-xr `docs/design/JADC2-MAP.md` §5.3 and `docs/contracts/asset-command.md` §9.
`REDTEAM_REPORT=/path/table.md go test ./internal/policy/ -run 'TestRedTeam|TestForbiddenNames'`
writes the pass/fail table.

## Configuration (environment)

| Variable | Default | Meaning |
|---|---|---|
| `NATS_URL`, `NATS_USER`, `NATS_PASSWORD`, `NATS_CA_FILE` | `nats://127.0.0.1:4222` | the `policy` NATS user (password from a SealedSecret) |
| `POLICY_AUTHZ` | `openfga` | `openfga` (HTTP Check API) or `local` (model JSON + tuples file, local stack only) |
| `OPENFGA_URL`, `OPENFGA_STORE_ID`, `OPENFGA_MODEL_ID`, `OPENFGA_TOKEN` | — | OpenFGA; the token from a SealedSecret |
| `POLICY_MODEL_FILE`, `POLICY_TUPLES_FILE` | `/config/onemind-actions.json`, `/config/tuples.json` | for `local` |
| `POLICY_SCHEMA_DIR` | `/schemas` | copied into the image from `services/schemas` |
| `POLICY_PRINCIPALS_FILE` | `/config/principals.json` | NATS identity to principal map |
| `POLICY_SIGNING_KEY_FILE`, `POLICY_SIGNING_KID` | unset, `policy-1` | Ed25519 PKCS#8 PEM; unset means unsigned commands (logged) |
| `POLICY_GEOFENCE_BUCKET` | `geofences` | KV bucket watched for zones |
| `HEALTH_ADDR` | `:8080` | `/healthz` |

## Not done here (follow-ups)

- Action records live in memory. A restart forgets pending actions (they can no longer be
  confirmed, which fails safe) and the idempotency memory. Denials are not held. Every other
  state has a deadline, so every action finishes and is dropped from memory 24 h later (the
  audit trail in `EVT` keeps them). The contract's Postgres decision
  row with an outbox (asset-command.md §5) is a follow-up.
- One replica (`Recreate`); no leader election.
- The bridge-side checks (signature, nonce, expiry, own geofence) are H6 and the MAVLink
  bridge (H2); the stub bridge in the tests only mirrors them.
- Composite task types (`survey_area`, `follow_route`, `patrol`, `inspect_structure`,
  `capture_media`, `relay`, `deliver`) are decomposed by the task service and are refused here
  as `unknown_action` until then.

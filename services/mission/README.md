# Task service (`services/mission`)

FEATURES D2 (missions) and H4 (tasking model), onemind-xr. It implements
`docs/contracts/mission-task.md` (onemind-xr, branch `run/m0`): missions and tasks as records
in KV `mission-state`, a task state machine with typed assignees (asset, person, team,
agent), and the facts on `evt.mission.*`. It is the "task service" the policy service
(omos #554) waits on and the MAVLink bridge (omos #550) reports to.

It never commands an asset: only the policy service publishes `cmd.*` (the service's NATS
user has no publish grant on `cmd.>`).

## Bus

| Direction | Subject or key | Notes |
|---|---|---|
| in | `cmd.mission.{msn}.{create\|waypoint\|complete\|abort}`, `cmd.mission.{tsk}.{create\|abort\|complete}` | from the policy service after a human confirm; durable consumer `mission-commands` on `CMD`. Expired commands (`expires_at`) and commands without an approval are refused, never applied late |
| in | `evt.so.action.dispatched` | `task.assign` sets `sent`; the asset command id (and the abort hold's) goes into the report guard. Durable consumer `mission-events` on `EVT` |
| in | `evt.drone.*.task`, `evt.robot.*.task` | asset status reports (section 7); the subject id must be the task's assignee |
| in | `tel.fleet.*.*.position` (core) | arrival check for asset `goto` tasks (owner proposal D2-P2; `MISSION_ARRIVAL`) |
| out | KV `mission-state` keys `msn-…`, `tsk-…` | only writer; every write is compare-and-set on the KV revision, re-read and re-applied on a conflict |
| out | `evt.mission.{tsk}.{state}` | one fact per state change; `type` `evt.mission.task.{state}`; `data.status_version` = the record's `version.status`; `Nats-Msg-Id` `{task_id}:{status_version}`. `created` carries the whole record (section 6.1 step 5) |
| out | `evt.mission.{msn}.{created\|waypoint\|updated\|assigned\|unassigned\|completed\|aborted}` | `data` = the mission record plus `event` (and `action_id`, `task_id` where relevant) |
| out | `evt.mission.{id}.rejected` | a refused command: `data` `{verb, code, message, action_id, command_id, at_t}` |
| serve | `svc.mission.get` `{v, id}` → `{record, revision}` | |
| serve | `svc.mission.list` `{kind?, state?, assignee?, parent_id?, bbox?, limit?}` → `{items, count, truncated}` | from the watched copy of the bucket |
| serve | `svc.mission.report` (people and agents) | section 7 `data`; the caller must be the verified assignee (server-set `Nats-Request-Info` from a CLIENT or AGENTS import) |
| serve | `svc.mission.apply` (policy service only) | `{v, action_id, op: assign\|unassign, mission_id, to: "{type}.{id}", role?, by?}` → `evt.mission.{msn}.assigned`; refuses callers from other accounts |

Every result of a command carries `causation` = the command id, so the policy service
closes the action on it (its POL-14): `created`, `waypoint`, `completed`, `aborted` and
`cancel_requested` become `evt.so.action.executed`; `rejected` becomes `failed` with `code`.

## State machine and ordering

States, rollups and terminal rules are those of section 5. A task record exists only after a
human approved it. Missions go `planned` → `active` (first child executes) → `completed`
(`complete`, or every child `done_ok`) or `aborted` (`abort`; every open child is cancelled at
once if not yet sent, else `cancel_requested`).

`version.status` is the service's own counter: +1 per accepted change. Reports are ordered by
`report_guard` (owner proposal D2-P3), because a bridge numbers its reports per task from 1
for the life of its process (#550 P7) and knows nothing of `version.status`:

- `report_guard.commands` lists the asset command ids of the task in dispatch order (from
  `evt.so.action.dispatched`, or from a report that arrives before that audit).
- A bridge report is ordered by (index of its `causation` in `commands`, its
  `status_version`). A report from a newer command wins even when a restarted bridge started
  again at 1; inside one command `status_version` must increase; anything else is dropped and
  counted (`mission_reports_dropped_stale_total`). A report for a terminal task is dropped and
  counted too (`mission_reports_dropped_terminal_total`).
- People and agents (`svc.mission.report`) report one sequence per task: `status_version`
  must be greater than the last accepted one.

Abort (section 6.2): `cancel_requested` keeps its state against progress from the original
command; a report from the hold that says `executing` or `done_ok` ends the task
`done_not_ok` + `cancelled`; if the asset refuses the hold, the task returns to its previous
state with `error.code = rejected`. A hold is recognised by what it is, not by arrival order:
its command id is in `report_guard.abort_commands` (from the `task.abort` dispatch audit), or
it is any command after the task's one `task.assign` command. The abort command (`CMD`), the
hold's audit and the hold's reports (both `EVT`) can arrive in any order: a hold report
before the abort keeps the task's state and is kept in `report_guard.hold`, and is applied
when the abort lands.

Redelivery: every record keeps the newest 32 applied `cmd.mission` command ids
(`applied_commands`). A redelivered command already applied (ack lost) is recognised before
the expiry and precondition checks: it never produces a `rejected` fact (a create reports
`created` again so the action can close; anything else is a no-op).

Timeouts (rule 4): no `received` by `delivery.before_t` (default creation + 60 s) or past
`stale_t` (default `execution.complete_before_t`, else creation + 1 h) ends `done_not_ok` +
`timeout`. Nothing is re-sent.

## Tests

```
cd services/mission
go vet ./... && go vet -tags e2e ./...
go test -race -count=1 ./...
```

The unit tests start an in-process nats-server with the A19 account layout (ONEMIND, CLIENT,
AGENTS), a least-privilege `mission` user (any JetStream call it is not granted fails the
test), the `CMD`, `EVT`, `MISSION` streams and KV `mission-state`, publish commands and audit
events exactly as `services/policy` does and reports exactly as the MAVLink bridge does, and
validate every fact against `services/schemas/envelope/envelope.schema.json`. They cover
mission-task.md section 10 tests 1 (task-service half), 2, 4 (defence in depth: a create whose
target or target snapshot is a person is refused `forbidden_action`) and 6, every ordering of
abort, hold audit and hold reports, redelivery after expiry, list paging, eviction, the bridge-restart and
early-report cases, abort, rejections, idempotent redelivery, missions, the RPCs, arrival,
restart and a forced compare-and-set conflict.

End to end, with the real policy service (#554), the real MAVLink bridge (#550) and the
onemind-xr scripted drone (section 10 tests 1, 3, 4 and 6, plus the whole-run invariant that no
`cmd.fleet.>` exists without a human approval):

```
SIM_DRONE=../onemind-xr/scripts/sim/mavlink_sim_drone.py PYTHON=python3 ./test/e2e.sh
```

The drone flies the task goto and the task ends `done_ok` (POL-16 / D2-P1: the bridge accepts
the `task.assign` `params_hash` on a task-bearing command).

## Configuration (environment)

| Variable | Default | Meaning |
|---|---|---|
| `NATS_URL`, `NATS_USER`, `NATS_PASSWORD`, `NATS_CA_FILE` | `nats://127.0.0.1:4222` | the `mission` NATS user (password from a SealedSecret) |
| `MISSION_BUCKET` | `mission-state` | KV bucket |
| `MISSION_CMD_STREAM`, `MISSION_EVT_STREAM` | `CMD`, `EVT` | streams of the two durable consumers |
| `MISSION_CONSUMER` | `mission` | consumer name prefix (`mission-commands`, `mission-events`) |
| `MISSION_REPORT_CLASSES` | `drone,robot` | asset classes whose `evt.{class}.*.task` `EVT` binds (`vehicle` waits for ONT-P13) |
| `MISSION_PRINCIPALS_FILE` | `/config/principals.json` | identity directory for `svc.mission.report` (missing: every report is refused) |
| `MISSION_ARRIVAL`, `MISSION_ARRIVAL_RADIUS_M`, `MISSION_ARRIVAL_VERTICAL_M` | `false`, `5`, `5` | goto arrival from telemetry (D2-P2; off until the owner ratifies it, since it trusts any `tel.fleet.*.*.position` publisher) |
| `MISSION_CACHE_TERMINAL_FOR` | `24h` | how long terminal records stay in the in-memory copy (`svc.mission.list`) |
| `MISSION_SWEEP_EVERY` | `1s` | timeout sweep period |
| `MISSION_SOURCE` | `service:task` | envelope `source` and `updated_by` of service-made changes |
| `HEALTH_ADDR` | `:8080` | `/healthz`, `/metrics` (counters) |

## Memory and retention

- KV `mission-state` keeps every record (history 5); the `MISSION` stream keeps 90 days of
  facts. Nothing here deletes a record.
- The in-memory copy (KV watch) holds every open record and terminal records until
  `MISSION_CACHE_TERMINAL_FOR` (24 h) after their last change; an eviction pass runs every
  minute. `svc.mission.list` shows only that copy; `svc.mission.get` reads KV, so an evicted
  record is still served. The timeout sweep walks open tasks only. Writers of one key are
  serialised by 64 striped locks (no per-key allocation).
- On start the watch loads every KV record once, then evicts old terminal ones. If the
  bucket grows past what 128Mi holds at start, the owner should add a KV `max_age` or an
  archival job (follow-up; not decided here).
- `svc.mission.list` replies are capped by count (`limit`, max 1000) and by size (about
  900 KB, below the 1 MB NATS default `max_payload`); a capped reply has `truncated: true` and
  `next`, which the caller passes as `after` for the next page.

## Not done here (follow-ups, see the owner proposals)

- Composite task types (`survey_area`, `follow_route`, …) are accepted as records but not
  decomposed into child tasks yet (H5); the policy service refuses them today.
- The signature on `cmd.mission.*` is not verified; only the `policy` user may publish
  `cmd.mission.>` (A19 stage B) (D2-P8).
- The entity record's `current_task` is not updated (the entity service owns `ent.*`) (D2-P9).
- TAK Server mission sync (`tak_mission`) is the tak-client bridge's (#560); this service
  keeps the field.

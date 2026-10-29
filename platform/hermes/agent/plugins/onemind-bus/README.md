# onemind-bus (Hermes plugin)

Hermes tools over the OneMind NATS bus: onemind-xr FEATURES **I2** (agent tools over NATS),
**D5** (agent actions go through policy; every asset command needs a human confirm; an agent
proposes and never confirms) and **D11** (reversible map actions over `cmd.map.>`, reporting
only actions that succeeded). Design sources (onemind-xr, branch `run/m0`):
`docs/research/agents-that-act.md` §1 and §8, `docs/design/BUS-DESIGN.md` §3.3, §3.4 and §7,
`docs/contracts/asset-command.md` §1–§2, omos `docs/NATS-DOCTRINE.md` §9.7.

It replaces the generic bus pipe: when this plugin is enabled, the `nats_publish` tool of the
`nats-context` plugin is blocked (see "Retiring `nats_publish`").

## Tools

| Tool | Subject | Kind |
|---|---|---|
| `onemind_entity_get` | `svc.entity.get` | read (A0) |
| `onemind_entity_query` | `svc.entity.query` (pages; `near`, `name`, `moving_only` filtered in the tool) | read (A0) |
| `onemind_layer_list`, `onemind_layer_meta` | `svc.layer.list`, `svc.layer.meta` | read (A0) |
| `onemind_layer_snapshot` | `svc.layer.snapshot` (area required, at most 200 items, `fields: render`) | read (A0); proposal I2-P1 |
| `onemind_command_propose` | `svc.policy.command.propose` | proposal; waits for a human |
| `onemind_task_propose` | `svc.policy.task.propose` | proposal; waits for a human |
| `onemind_action_status` | `svc.policy.action.get` | read of own actions |
| `onemind_map_fly_to`, `onemind_map_frame`, `onemind_map_camera_tour` | `cmd.map.viewport.set` | reversible UI (A1) |
| `onemind_map_layers` | `cmd.map.filter.set` | reversible UI (A1) |
| `onemind_map_select` | `cmd.map.highlight.{entity_id}` | reversible UI (A1) |
| `onemind_map_annotate`, `onemind_map_measure` | `cmd.map.annotation.{id}.create` | reversible UI (A1) |
| `onemind_map_remove_annotation` | `cmd.map.annotation.{id}.remove` | reversible UI (A1) |

There is **no confirm tool**, no generic publish and no tool argument that names a subject. The
bus client (`bus.py`) refuses any subject outside `ALLOWED_REQUESTS` and the `cmd.map` pattern
before sending, so no tool can reach `cmd.fleet.>`, `svc.policy.*.confirm`, `tak.cot.*` or
`ent.*`. Geofence edits (`cmd.map.geofence.*`) are A2 actions that write KV `geofences`, so
they are not a direct agent map action. The server-side guarantee is the per-agent NATS user
(A19, omos PR #562): the `hermes` user cannot publish those subjects at all.

Proposals carry `requested_by: {kind: agent, id: ONEMIND_BUS_AGENT_ID, on_behalf_of:
ONEMIND_BUS_ON_BEHALF_OF}` and the audit hashes `agent_session` (the Hermes task id),
`prompt_sha256` and `tool_input_sha256`; raw prompts never leave the pod. The `request_id` is
derived from the session and the tool input, so a retried tool call returns the same action
instead of creating a second one. The tool result tells the model that the action is only
proposed and awaits a human.

Every tool call writes one JSON log line (`event: onemind_bus_tool`, tool, outcome, duration,
session, argument hash) to the `onemind_bus` logger.

## `cmd.map.*` payload v1 and the client ack

No new subject: the subjects are omos `docs/NATS-DOCTRINE.md` §9.1/§9.7. The payloads extend
the §9.7 examples (proposal I2-P4). Every payload has:

```jsonc
{ "v": 1, "id": "<32 hex>", "ts": "RFC 3339", "source": "agent:hermes/zeus/legacy",
  "corr": "<hermes session>", "reason": "short text for the operator" }
```

| Subject | Fields |
|---|---|
| `cmd.map.viewport.set` | `mode: fly_to` with `lat, lon, range_m, heading_deg, pitch_deg, duration_s` and §9.7 `zoom` (derived from `range_m`); `mode: frame` with `bbox [w,s,e,n]`, `entity_ids`; `mode: tour` with `stops [{lat, lon, range_m, heading_deg, pitch_deg, dwell_s, duration_s, label}]`, `loop` |
| `cmd.map.filter.set` | `show: [layer]`, `hide: [layer]` (layer ids from `svc.layer.list`); clients keep accepting the §9.7 `layers` full list |
| `cmd.map.highlight.{uid}` | `select: true, ring, label, color`, or `clear: true` |
| `cmd.map.annotation.{id}.create` | `kind: note\|line\|area\|boundary\|measure`, `label`, `note`, `geometry` (GeoJSON Point, LineString or Polygon), `measure: {distance_m, area_m2}`, `entity_id(s)`. `boundary` geometry is copied from the entity record (`geo.coords`), never written by the model |
| `cmd.map.annotation.{id}.remove` | — |

The message is published with a reply subject and `Nats-Msg-Id` = `id`. A map client that
applied it replies on that subject:

```jsonc
{ "v": 1, "id": "<same id>", "ok": true, "client": "<client id>", "applied": { /* resulting view */ } }
{ "v": 1, "id": "<same id>", "ok": false, "client": "<client id>", "error": "why" }
```

A tool reports `ok: true` only when at least one client acked `ok: true` within
`ONEMIND_BUS_MAP_ACK_WINDOW_S`. Otherwise it returns `ok: false` with `error: no_client_ack` or
`rejected_by_client` and a `say` line telling the model not to claim the action. A JetStream
PubAck from a stream that captures `cmd.>` is recorded as `audit: {stream, seq}` and never
counted as an ack. In the A19 account layout the ack cannot cross from CLIENT back to AGENTS
over the stream import; proposal I2-P2 fixes it (proved in `tests/test_account_layout.py`).

No client implements this yet: the Unreal client (onemind-xr `run/m2-client`) and the L0 web
client (onemind-xr PR #2) need a `cmd.map.>` subscriber that applies and acks (follow-up for
those lanes).

## Configuration (environment of the Hermes container)

| Variable | Default | Meaning |
|---|---|---|
| `ONEMIND_BUS_URL` | `NATS_URL`, else `nats://localhost:4222` | in the pod: the `nats-leaf` sidecar |
| `ONEMIND_BUS_USER`, `ONEMIND_BUS_PASSWORD`, `ONEMIND_BUS_CREDS` | unset | only for a direct connection; in the pod the leaf carries the identity |
| `ONEMIND_BUS_AGENT_ID` | `agent:hermes/zeus/legacy` | must match the policy principals file entry for the `hermes` user |
| `ONEMIND_BUS_ON_BEHALF_OF` | unset | the human the agent proposes for (for example `user:zeus`); unset means every proposal is refused locally (`not_configured`). Proposal I2-P5 |
| `ONEMIND_BUS_MODEL_ALIAS` | unset | recorded in the audit event |
| `ONEMIND_BUS_TIMEOUT_S` | `5` | request timeout |
| `ONEMIND_BUS_MAP_ACK_WINDOW_S` | `1.5` | how long a map tool waits for a client ack |
| `ONEMIND_BUS_LAYER_SNAPSHOT` | `1` | `0` removes the snapshot read (until I2-P1 is ruled) |
| `ONEMIND_BUS_RETIRE_NATS_PUBLISH` | `1` | block `nats_publish` (and `ONEMIND_BUS_RETIRED_TOOLS`) through `pre_tool_call` |
| `ONEMIND_BUS_RETIRED_TOOLS` | unset | extra tool names to block, comma-separated (for example `tak_publish_cot`) |

Dependencies: `nats-py` (already installed into `/gateway-deps` by the `install-gateway-deps`
init container in `deployment.yaml`) and the Python standard library.

## Retiring `nats_publish`

The plugin registers a `pre_tool_call` hook that returns `{"action": "block", ...}` for
`nats_publish`. That hook contract is Hermes `hermes_cli/plugins.py`
(`get_pre_tool_call_block_message`), read in the 0.14.0 source; the live image is
`v2026.9.24`, where it is **unverified** until the owner step below runs. The flag is on by
default, so enabling the plugin retires the pipe; `ONEMIND_BUS_RETIRE_NATS_PUBLISH=0` rolls it
back without disabling the new tools. The permanent removal is the A19 `hermes` NATS user,
which cannot publish outside its allow-list whatever tool asks.

## Install (owner; nothing here changes the live pod)

This PR adds source and tests only. The plugin is installed the way `onemind-vision` is (a
GitOps gap recorded in `../README.md`):

1. Copy this folder to the Hermes volume as `/opt/data/plugins/onemind-bus/` (without `tests/`).
2. Add `onemind-bus` to `plugins.enabled` in the Hermes `config.yaml` on the volume.
3. Set `ONEMIND_BUS_ON_BEHALF_OF` (and optionally `ONEMIND_BUS_MODEL_ALIAS`) in the
   `hermes-env` SealedSecret or as plain `env` in `deployment.yaml` (no secret value).
4. Restart Hermes; check the log for `onemind_bus_registered` and that a `nats_publish` call is
   answered "retired".

The exact commands and expected output are in onemind-xr `docs/evidence/M4/I2/OWNER-STEPS.md`.

## Tests

```bash
cd platform/hermes/agent/plugins/onemind-bus/tests
pip install nats-py pytest pyyaml
NATS_SERVER_BIN=$(command -v nats-server) python -m pytest -v
# D5 against the real policy service (omos PR #554):
POLICY_BIN=/path/to/policy POLICY_SRC=/path/to/omos-with-#554 python -m pytest -v test_policy_integration.py
```

| File | What it proves |
|---|---|
| `test_tools.py` | each tool against stub entity, layer, policy and map services on a disposable server; the subject allowlist; the retire hook; no tool reaches `cmd.fleet.>` or confirm even with hostile arguments; injected callsigns come back as data |
| `test_scripted_model.py` | D11: ten scripted text requests, each checked against the stub map client's view; with no ack the script reports failure. I2: "What's moving near the north field?" answered with bus tools only. D5: the follow-a-person refusal |
| `test_policy_integration.py` | D5 with the real policy service: a Hermes proposal waits; Hermes cannot confirm (server permission violation; `agent_cannot_confirm` for an agent identity that reaches confirm); only a human confirm dispatches; follow-a-person is denied `forbidden_action` and audited; no `cmd.fleet.>` without a human approval |
| `test_account_layout.py` | the cross-account ack problem of the A19 layout and the I2-P2 fix |

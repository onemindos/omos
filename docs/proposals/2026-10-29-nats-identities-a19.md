# Owner proposals: per-service NATS identities (FEATURES A19)

Status: **proposals for the owner**, 2026-10-29. Filed with the identity config in
`platform/nats/identities/` (runbook `platform/nats/IDENTITIES.md`). Nothing is live. No new
subject pattern is introduced: every subject granted already exists in the locked grammar,
in the wave-1 PRs (#550 #551 #552 #553 #554) or in live legacy use. Server-internal subjects
(`_INBOX`, `_R_`, `$JS`, `$KV`, `$SRV`) are not grammar subjects. Each item gives the reading
used until it is ruled.

| ID | Decision needed | Evidence | Proposed ruling | Used until ruled |
|---|---|---|---|---|
| A19-P1 | Hermes moves to account `AGENTS` with its own leaf user `hermes` in stage A2, not a restricted user inside `ONEMIND` first | Policy treats every `AGENTS` identity as an agent and reads the caller from the server-set `Nats-Request-Info` (#554 `principals.json`). Inside `ONEMIND` a client can forge that header on a proposal (#554 proposal, "Residual risk"). permtest `TestHermesLeaf`: from its leaf, Hermes is seen as `{"acc":"AGENTS","user":"hermes"}` | Move Hermes to `AGENTS` directly | `AGENTS` |
| A19-P2 | Hermes loses `cmd.ha.*` (Home Assistant control, `platform/hermes/agent/RUNBOOK.md`) and the raw feeds it receives today (`tel.>`, `intel.>`, `tak.cot.>`, `ent.>`, `evt.>`, `det.>`, `cmd.alert.>`, `sys.>`; onemind-xr `omos-live-state.md` §6) | BUS-DESIGN §7.2 denies Hermes every `cmd.*` except `cmd.map.>` and the firehoses (prompt-injection surface, security-and-trust.md T1); FEATURES A19 test | Home Assistant actions become a policy action (proposed through `svc.policy.command.propose`, executed by the HA bridge). Hermes reads entities through `svc.entity.*` | Denied |
| A19-P3 | `svc.layer.watch` streams batches "on the same reply subject until the lease ends" and renewals arrive as new requests (onemind-xr `contracts/layer-delta.md` §3) | Server limits: `allow_responses` and the export `response_threshold` bound how long replies may follow one request (nats-server 2.14.2 defaults: 1 reply, 2 min). permtest `TestCrossAccountRequests` proves 5 streamed replies through the CLIENT import with `max: -1, expires: 10m` | Each renewal re-arms the stream: the layer service sends later batches to the newest renewal's reply subject. Lease stays at most 300 s | `layer`: `allow_responses {max: -1, expires: 10m}`; export `response_type: stream`, `response_threshold: 10m` |
| A19-P4 | `cmd.map.>` flows `AGENTS → CLIENT` only; `ONEMIND` neither sends nor receives it. BUS-DESIGN §7.1 lists `ONEMIND` exporting stream `cmd.map.>` | No ONEMIND service consumes or produces `cmd.map.>` in #550–#554. Two-way stream exports of one subject between accounts risk loops and duplicates | Keep it one-way. D1 geofence drawing reaches the policy service as a proposal, not as a `cmd.map.>` message ONEMIND reads | One-way |
| A19-P5 | Reply privacy inside an account | Every requester needs subscribe `_INBOX.>`, so it can read other users' replies in the same account (CLIENT users can read each other's replies) | Per-identity inbox prefixes (`nats.CustomInboxPrefix("_INBOX.<user>")` in each service and client) and subscribe `_INBOX.<user>.>`; or stage 2 (JWT, auth callout) | `_INBOX.>` |
| A19-P6 | Agents outside Hermes: `legacy-edge` (Hermes on the Mac, `agents.prompt.hermes.zeus.legacy-edge`) and delegation to other frameworks (`agents.prompt.<framework>.…`, legacy-edge `SOUL.md`) | `AGENTS` imports nothing from `ONEMIND` under `agents.prompt.>`, so the cloud Hermes cannot prompt agents that live in `ONEMIND`. Prompts to `legacy-edge` from ONEMIND still reach the Mac leaf locally. A wildcard import both ways would form a service-import cycle | Each agent framework moves into `AGENTS` with its own user, or gets a non-overlapping import (`agents.prompt.<framework>.>`) on request | No cross-account delegation |
| A19-P7 | Per-account JetStream limits | Live `ONEMIND` limits are all `-1` (nats-config read 2026-10-29); the server-level `max_file_store` is 45Gi | Set `ONEMIND` `max_file` from the A22 cap arithmetic once the stream caps merge; `CLIENT` and `AGENTS` stay without JetStream | `ONEMIND` `-1`; no JetStream for `CLIENT`/`AGENTS` |
| A19-P8 | Mesh gateway forwarding to TAK Server | Decision 0004 #6: `MESH_TO_TAK_SERVER=false` by default (#551) | If enabled, add `tak.cot.inject` to `mesh-gateway` and `leaf-site-mac` and allow it in the checker (`scripts/nats_identities.py`, projector-only rule) | Not granted |
| A19-P9 | Legacy subject `cmd.projector.{enable,disable}` (3 tokens; the grammar's commands have 4) | Live kill switch in `platform/omos/tak-projector.yaml` | Keep for the projector only; rename when the projector is rebuilt | Granted to `projector` (subscribe) |
| A19-P10 | What stage B closes for the shared logins `onemind` and `leaf_tak`, and when the rest closes | Stage B denies only `cmd.{fleet,drone,robot,vehicle,mission,entity}.>`, `tak.cot.inject`, `tak.cot.mesh.tx`, `ent.*.*.state`, `evt.so.action.>`, `svc.policy.>` (`identities.yaml` `legacy_stage_b`). Still open until stage C: `cmd.projector.disable` (live kill switch, also reachable with the browser-exposed `ws.onemindos.dev` password), `cmd.ha.*`, `cmd.alert.*`, future command domains, and forged `tak.cot.out`, `tel.fleet.drone.sim-x1.*`, `evt.drone.sim-x1.*` (entity and policy consume them). The #554 migration step 3 ("deny `cmd.>` except `cmd.map.>` for every remaining shared user") says more than stage B does | (a) Move the projector toggle behind policy (`cmd.projector.>` published only by `policy` after a human confirm), then add `cmd.projector.>` to the stage-B deny list. (b) Add a stage B2 that denies `onemind` and `leaf_tak` `cmd.>` except `cmd.map.>`, plus `tak.cot.out`, `tel.fleet.>` and `evt.drone.>`, once the TAK leaf is on `leaf-tak` and the MAVLink bridge and mesh gateway are on their own users (check `connz?auth=1` and `leafz` first). (c) Stage C retires both logins | Stage B as listed; the open subjects are recorded in `platform/nats/IDENTITIES.md` "Limits" |

## Information only (no ruling needed)

- **Hermes refusals are server permission violations at its sidecar.** A hub only filters a
  leaf: it advertises just the subjects the leaf user may publish, so a forbidden publish
  vanished on the spoke with no error and no log line (first permtest run). The rendered Hermes
  `leaf.conf` (`nats_identities.py leaf-conf`) now gives the agent's local user the hub
  allow-lists, so `cmd.fleet.drone.sim-x1.goto`, `tak.cot.inject` and every other denied subject
  returns `Permissions Violation` to the agent and logs `Publish Violation` in the `nats-leaf`
  container (`TestHermesLeaf`). The FEATURES A19 wording is met as written; no ruling needed.
- **Legacy users get connection types from A2:** `onemind` STANDARD and WEBSOCKET, `sys` and
  `monitor` STANDARD, so removing `leafnodes.authorization` does not let them open leaf links.

- **Stage A2 needs a pod restart.** Leaf users carry permissions only as account users, and
  nats-server 2.14.2 cannot reload `leafnodes.authorization.users`. The runbook rolls the
  three pods one at a time.
- **Drift found, not fixed:** `platform/nats/values.template.yaml` has
  `jetstream.fileStore.pvc.size: 20Gi` (rendered `max_file_store: 20Gi`), the live ConfigMap
  has `45Gi`. The identity values do not touch it. Follow-up for whoever adopts the NATS
  release into Argo.
- **Live config structure verified read-only** on 2026-10-29 (passwords redacted before
  output): accounts `ONEMIND{onemind}`, `SYS{sys, monitor}`, `leaf_tak` in
  `leafnodes.authorization`, WebSocket `no_tls`, server 2.14.2, chart `nats-2.14.2`.

# Per-service NATS identities (FEATURES A19)

Status, 2026-10-29: **not live**. This directory gives the bus one user per service with
publish and subscribe allow-lists, a `CLIENT` account for viewers, an `AGENTS` account for
Hermes, and separate leaf users for the TAK leaf and the Mac site. Activating it is an owner
step. Each stage below is a change to the `nats/nats-config` ConfigMap; there is never a
`helm upgrade` (it would overwrite the Git-owned StatefulSet, see `OPERATIONS.md`).

Why: today every application shares one user, `onemind`, with no `permissions` block, and one
leaf user, `leaf_tak` (TAK leaf, Hermes sidecar and `zeus-mac` all use it). Anyone holding the
shared password can publish `cmd.>` and `tak.cot.inject`, also through the public
`ws.onemindos.dev`. Design: onemind-xr `docs/design/BUS-DESIGN.md` §7,
`docs/research/security-and-trust.md` Rec. 1 and the permission proposal
`docs/proposals/2026-10-29-nats-per-service-permissions.md` (PR #554).

## Files

| File | What it is |
|---|---|
| `identities/identities.yaml` | Source of truth: accounts, users, permissions, Secret locations. No credentials |
| `../../scripts/nats_identities.py` | Renders the spec (Helm values, live ConfigMap, test config), checks invariants, generates passwords outside the repo |
| `values.identities.stage-{a1,a2,b}.template.yaml` | Generated Helm values with `__<USER>_PW__` placeholders, for parity with `values.secret.yaml` |
| `identities/permtest/` | Go tests: the rendered config on in-process nats-server 2.14.2 (the live version) |
| `../../scripts/tests/test_nats_identities.py` | Renderer and invariant unit tests |
| `../../docs/proposals/2026-10-29-nats-identities-a19.md` | Decisions that need the owner |

`python3 scripts/nats_identities.py check` must pass after any edit of the spec. CI
(`.github/workflows/test-nats-identities.yml`) runs it, the Python unit tests and the Go
permission tests against the checksum-verified nats-server v2.14.2 release binary (see Tests).

## Identities

Only the essentials; `identities.yaml` has the full lists. `_INBOX.>` subscribe is granted to
every requester. Responders use `allow_responses`, so no service may publish on `_INBOX.>`.

| User | Account | Publishes | Subscribes / serves | Password goes to |
|---|---|---|---|---|
| `policy` | ONEMIND | `cmd.>`, `evt.so.action.>`, `svc.entity.get`, `svc.mission.apply` | `svc.policy.>`, `evt.{drone,robot,vehicle,mission,entity}.>`, `evt.map.geofence.>` (zone command results `evt.map.geofence.{created,updated,deleted,refused}`, D1-P2; without it zone actions end `no_ack`); reads KV `geofences` (writes no bucket) | SealedSecret `onemind-core/policy-service-runtime` (`NATS_USER`, `NATS_PASSWORD`; it also holds `OPENFGA_TOKEN`, #554) |
| `entity` | ONEMIND | `ent.*.*.state`, `evt.entity.>`, `evt.sys.entity.rejected`, KV `entity-correlation` | `svc.entity.>`, `cmd.entity.>`, `$SRV` (micro); consumes COT TEL WORLD INTEL DET ENT | `omos/omos-entity-nats` (`user`, `password`) |
| `ingest` | ONEMIND | `intel.>`, `evt.world.>`, `tel.feed.>`, `evt.sys.*.rejected`, KV `ingest-locks`, `ingest-state` | reads KV `feed-config` | `omos/omos-ingest-nats` |
| `layer` (#561) | ONEMIND | KV `world-state`, `fleet-state` | `svc.layer.>` (streamed replies), `tel.feed.*.status`, `tel.fleet.node.*.status`; consumes ENT, last value from TEL, lists streams (A19-P15) | `omos/omos-layer-nats` |
| `projector` (live Python) | ONEMIND | `tak.cot.inject` | `ent.*.*.state` (only with `ENT_PROJECTION=true`), `cmd.projector.>`; consumes INTEL | `omos/omos-tak-projector-nats` |
| `projector-gated` (#563) | ONEMIND | `tak.cot.inject`, `tak.cot.mesh.tx` | `svc.cot.send`, `ent.*.*.state`, `tak.cot.out`, `tak.cot.mesh.rx` (tasking guard), `cmd.projector.>`; consumes ENT; reads KV `geofences`, `mission-state` | `omos/omos-tak-projector-gated-nats` |
| `tak-client` (#560) | ONEMIND | `tak.cot.out` (acked into COT) | `tak.cot.inject` (queue `tak-cot-inject`) | `takserver/tak-client-nats` |
| `geofence` (#569) | ONEMIND | `evt.map.geofence.>`, `evt.so.alert.>`, KV `geofences` (only writer) | `svc.geofence.alerts`, `tel.feed.*.status`; consumes CMD ENT MISSION EVT | `omos/omos-geofence-nats` |
| `mission` (#568) | ONEMIND | `evt.mission.>`, KV `mission-state` (only writer) | `svc.mission.>`, `tel.fleet.*.*.position`; consumes CMD EVT | `onemind-core/mission-service-runtime` (`NATS_USER`, `NATS_PASSWORD`) |
| `agentops` (#567) | ONEMIND | `tel.fleet.agent.>`, `evt.entity.*.alert`, `evt.map.geofence.>`, `evt.so.brief.published`, its `agents.hb.*`, requests `svc.entity.query`, `svc.policy.task.propose`, `svc.cot.send`, `agents.prompt.hermes.>` (never `tak.cot.inject`, never confirm) | `svc.brief.>`, `tak.cot.out`, `tak.cot.mesh.rx`, `tel.feed.*.status`, `evt.entity.*.alert`; consumes ENT TEL AGENTS; reads KV `geofences` | `omos/omos-agentops-nats` |
| `video-registry` (#559) | ONEMIND | `evt.video.>`, KV `OVERWATCH_VIDEO` | `svc.registry.video.>` | `onetak/video-registry-nats` |
| `bridge-mavlink-home` | ONEMIND | `tel.fleet.drone.sim-x1.>`, `evt.drone.sim-x1.>` | `cmd.fleet.drone.sim-x1.>` | edge `.env` / keychain (direct connection only) |
| `bridge-ros2-home` (#566) | ONEMIND | `tel.fleet.robot.sim-tb1.>`, `evt.robot.sim-tb1.>` | `cmd.fleet.robot.sim-tb1.>` | edge `.env` / keychain (direct connection only) |
| `mesh-gateway` | ONEMIND | `tak.cot.mesh.rx`, `tel.fleet.edge.tak-mesh.status` | `tak.cot.mesh.tx` | edge `.env` / keychain (direct connection only) |
| `agent-fabric` | ONEMIND | `agents.prompt.>`, `$SRV.>` | `agents.>` | `agent-fabric/agent-fabric-nats` (`NATS_USER`, `NATS_PASSWORD`) |
| `leaf-tak` (leaf only) | ONEMIND | `tak.cot.out` | `tak.cot.inject` | inside `leaf.conf` of `takserver/tak-nats-leaf` |
| `leaf-site-mac` (leaf only) | ONEMIND | the site's bridges: `tel.fleet.drone.sim-x1.>`, `evt.drone.sim-x1.>`, `tel.fleet.robot.sim-tb1.>`, `evt.robot.sim-tb1.>`, `tak.cot.mesh.rx`, `tel.fleet.edge.tak-mesh.status` | `cmd.fleet.drone.sim-x1.>`, `cmd.fleet.robot.sim-tb1.>`, `tak.cot.mesh.tx`, `ent.*.*.state` | Mac leaf config (`infra/edge/mac-leaf`, keychain) |
| `hermes` (leaf only) | AGENTS | `agents.*.hermes.>`, `cmd.map.>` (acked by clients, I2-P2), requests `svc.entity.>`, `svc.layer.{list,meta}`, `svc.policy.{command,task}.propose`, `svc.policy.action.get`, `svc.mission.{get,list}`; replies `_R_.>` | `agents.prompt.hermes.>`, `svc.bridge.>`, `$SRV` | inside `leaf.conf` of `hermes/hermes-leaf-conf` (rendered by `leaf-conf`; its local user carries the same lists) |
| `zeus` | CLIENT | requests `svc.{layer,entity,policy}.>`, `svc.cot.send`, `svc.mission.{get,list,report}`, `svc.geofence.alerts`, `svc.registry.video.{get,list}`, `svc.tak.mission.>`, `agents.prompt.hermes.>`; `cmd.map.>`; acks with `allow_responses {max: -1, expires: 5s}` | `cmd.map.>`, `evt.so.action.>`, `evt.mission.>`, `evt.map.geofence.>`, `tel.feed.*.status`, `tak.cot.out`, `tak.cot.mesh.rx` (stream imports, A19-P11) | owner keychain on each client device |
| `alice` (test only) | CLIENT | same lists as `zeus` | same as `zeus` | none: `test_only`, rendered only by `conf --test-users` (permtest, E2E harness); never in the values files or the live config |
| `auth-callout` (#565) | AUTH | nothing (answers through `allow_responses {max: 1, expires: 5s}`) | `$SYS.REQ.USER.AUTH` | `onemind-core/auth-callout-runtime` (`NATS_USER`, `NATS_PASSWORD`) |
| `onemind`, `leaf_tak` | ONEMIND | everything until stage B; then denied `cmd.{fleet,drone,robot,vehicle,mission,entity}.>`, `tak.cot.inject`, `tak.cot.mesh.tx`, `ent.*.*.state`, `evt.so.action.>`, `svc.policy.>` (other `cmd.*` domains and forged inputs stay open until stage C, see Limits) | everything; from stage B not `svc.{policy,entity,layer,cot,mission,geofence,registry}.>` and not `svc.mission.apply` (cannot impersonate a service or assign tasks) | existing `values.secret.yaml` |

Accounts: `ONEMIND` keeps every stream. It exports the policy API to `CLIENT` (all of it) and
to `AGENTS` (propose and action lookup only, never confirm), `svc.entity.>` to both, and
`svc.layer.>` to `CLIENT` with streamed responses. `svc.mission.{get,list}` go to both,
`svc.mission.report`, `svc.geofence.alerts`, `svc.registry.video.>` and `svc.tak.mission.>` to
`CLIENT` only. `svc.mission.apply` is never exported. For the clients' displays it also exports
the read-only streams `evt.so.action.>`, `evt.mission.>`, `evt.map.geofence.>`,
`tel.feed.*.status`, `tak.cot.out` and `tak.cot.mesh.rx` to `CLIENT` (A19-P11). `CLIENT`,
`AGENTS` and `AUTH` have no JetStream. `AGENTS` exports Hermes' prompts, `svc.bridge.>` and
`$SRV` discovery to ONEMIND (and prompts to CLIENT), and its heartbeats and status to the
ONEMIND stream `AGENTS`. `CLIENT` exports `cmd.map.>` to `AGENTS` as a service with streamed
responses, so Hermes' map actions reach every client and the acks come back (#564 I2-P2).
`AUTH` holds only the auth callout responder; the `authorization.auth_callout` block is still
enabled by hand (`platform/auth-callout/nats-callout.conf`, after A1, not with A2: J2-P1).

Role modes (#565): callout-issued CLIENT users carry their role's allow lists in the user JWT.
`python3 scripts/nats_identities.py check-roles --roles platform/security/roles/onemind-roles.json`
checks that file against this account (A19-P11).

## What changes for running clients

| Client | Today | After its move |
|---|---|---|
| Hermes pod | leaf as `leaf_tak` in ONEMIND. The hub forwards `tel.>`, `intel.>`, `tak.cot.>`, `ent.>`, `evt.>`, `det.>`, `cmd.alert.>`, `sys.>` to it, and it can publish anything (onemind-xr `docs/research/omos-live-state.md` §6) | leaf as `hermes` in AGENTS: no raw feeds, no `cmd.*` other than `cmd.map.>` (so no `cmd.ha.*` Home Assistant control, proposal A19-P2), entity data through `svc.entity.*`, actions through `svc.policy.*.propose` |
| TAK leaf | `leaf_tak`, unrestricted | `leaf-tak`: `tak.cot.out` out, `tak.cot.inject` in |
| `zeus-mac` leaf | `leaf_tak` | stays on `leaf_tak` (HA, cameras, legacy subjects) until its bridges are moved; `leaf-site-mac` is the target for the OneMind site bridges |
| `omos-tak-projector` | `omos-nats-creds` (`onemind`) | `projector`; must move before stage B or its `tak.cot.inject` is refused |
| `omos-tak-projector-gated`, geofence, task, agentops, layer, video-registry (#559–#569) | `omos-nats-creds` or not deployed | their own users above; seal each Secret before the first sync |
| agent-console / agent-tap | `onemind` (`agent-fabric-nats`) | `agent-fabric` |
| `ws.onemindos.dev` browsers, `onemind-core`, `feed-*`, `sink-clickhouse`, NACK | `onemind` | unchanged in stage 1 (stage B removes their command rights); retire `onemind` in stage C |

## Activation (owner)

Prerequisites: the `do-nyc3-omos-cloud` context, `kubeseal`, the `nats` CLI (for bcrypt),
Python 3 with PyYAML, and a private working directory **outside the repository**, for example
`~/onemind-private/nats`. The renderer refuses to read or write password files inside the repo.

### 0. Save the current config (rollback point)

```bash
W=~/onemind-private/nats; mkdir -p "$W"; chmod 700 "$W"
kubectl --context do-nyc3-omos-cloud -n nats get cm nats-config -o yaml > "$W/nats-config.before.yaml"
kubectl --context do-nyc3-omos-cloud -n nats get cm nats-config -o jsonpath='{.data.nats\.conf}' > "$W/live.conf"
```

`live.conf` holds bcrypt hashes, not passwords; keep it private anyway.

### 1. Generate passwords for the new users

```bash
python3 scripts/nats_identities.py gen-passwords --out "$W/identities.json" \
  --users policy,entity,ingest,layer,projector,projector-gated,tak-client,geofence,mission,agentops,video-registry,bridge-mavlink-home,bridge-ros2-home,mesh-gateway,agent-fabric,leaf-tak,leaf-site-mac,hermes,zeus,auth-callout
```

`python3 scripts/nats_identities.py secrets` prints where each password is sealed.

Store the file in the password manager. Existing users (`sys`, `monitor`, `onemind`,
`leaf_tak`) keep their passwords: `live` copies their hashes from `live.conf`.

### 2. Seal one Secret per workload

Pattern (repeat per row of the table; namespace-wide scope as in the repo `CLAUDE.md`):

```bash
kubeseal --fetch-cert --controller-name=sealed-secrets-controller \
  --controller-namespace=kube-system > /tmp/pub-cert.pem
U=entity NS=omos NAME=omos-entity-nats
kubectl create secret generic "$NAME" -n "$NS" \
  --from-literal=user="$U" \
  --from-file=password=<(python3 -c 'import json,sys;sys.stdout.write(json.load(open(sys.argv[1]))[sys.argv[2]])' "$W/identities.json" "$U") \
  --dry-run=client -o json \
  | kubeseal --cert /tmp/pub-cert.pem --format yaml --scope namespace-wide \
  > platform/omos/entity/omos-entity-nats.sealedsecret.yaml
```

- `policy-service-runtime` uses keys `NATS_USER` / `NATS_PASSWORD` and also `OPENFGA_TOKEN`
  (`platform/policy/README.md` in #554): seal all three in one Secret.
- `agent-fabric-nats` keeps its keys `NATS_USER` / `NATS_PASSWORD`.
- Leaf credentials go inside `leaf.conf`. TAK leaf: the remote URL user changes from
  `leaf_tak` to `leaf-tak`; re-seal `platform/takserver/deploy/tak-nats-leaf.sealedsecret.yaml`
  (key `leaf.conf`). Hermes: render the whole sidecar config instead of hand-editing
  `platform/hermes/agent/SECRETS.md` §2, then seal it as
  `platform/hermes/agent/hermes-leaf-conf.sealedsecret.yaml` (instead of the
  `kubectl create secret` in that file):

  ```bash
  python3 scripts/nats_identities.py leaf-conf --user hermes \
    --passwords "$W/identities.json" --out "$W/hermes-leaf.conf"
  kubectl create secret generic hermes-leaf-conf -n hermes \
    --from-file=leaf.conf="$W/hermes-leaf.conf" --dry-run=client -o json \
    | kubeseal --cert /tmp/pub-cert.pem --format yaml --scope namespace-wide \
    > platform/hermes/agent/hermes-leaf-conf.sealedsecret.yaml
  ```

  The rendered `leaf.conf` keeps `server_name: hermes-legacy-leaf`, ports 4222/8222, the
  password-less `local` user and the TLS CA path, and gives `local` the **same publish and
  subscribe allow-lists as the hub user `hermes`**. The hub alone would only filter: it
  advertises to a leaf just the subjects that leaf user may publish, so a forbidden publish
  from the agent would be dropped on the spoke with no error and no log line. With the local
  allow-list the sidecar refuses it: the agent gets `Permissions Violation for Publish to
  "cmd.fleet.drone.sim-x1.goto"` and the `nats-leaf` container logs `Publish Violation - Subject
  "cmd.fleet.drone.sim-x1.goto"` (FEATURES A19 test; proved in `TestHermesLeaf`). Watch it with
  `kubectl -n hermes logs deploy/hermes -c nats-leaf | grep Violation`.
- `mission-service-runtime` and `auth-callout-runtime` use keys `NATS_USER` / `NATS_PASSWORD`
  (#568, #565); `video-registry-nats` is in namespace `onetak` (#559).
- Edge and client passwords (`bridge-mavlink-home`, `bridge-ros2-home`, `mesh-gateway`,
  `leaf-site-mac`, `zeus`) go to the owner keychain, not to the cluster.

Commit the SealedSecrets in a PR, but **do not repoint any workload until stage A1 is live**: a
workload that logs in as a user the server does not know yet stops.

### 3. Stage A1: add accounts and users (config reload, no restart)

```bash
python3 scripts/nats_identities.py live --live-conf "$W/live.conf" \
  --passwords "$W/identities.json" --stage A1 --out "$W/nats.conf.A1"
diff "$W/live.conf" "$W/nats.conf.A1" | less   # only accounts change at A1
kubectl --context do-nyc3-omos-cloud -n nats create configmap nats-config \
  --from-file=nats.conf="$W/nats.conf.A1" --dry-run=client -o yaml \
  | kubectl --context do-nyc3-omos-cloud -n nats apply -f -
```

The `reloader` sidecar on each pod reloads the config (the kubelet can take a minute to update
the mounted file). Check:

```bash
for p in nats-0 nats-1 nats-2; do kubectl --context do-nyc3-omos-cloud -n nats logs $p -c nats --since=5m | grep -E 'Reloaded|reload'; done
```

Expected: `Reloaded server configuration` on all three pods, and no error. Existing clients
and leafs stay connected; nothing is restricted yet for them. Then repoint workloads to their
new Secrets one at a time (separate PRs; check each pod's log for `Authorization Violation`
or `Permissions Violation`): `projector`, `agent-fabric`, then the services from #552, #553,
#554 as they deploy.

Keep `values.secret.yaml` in step: copy `values.identities.stage-a1.template.yaml` into it
with the bcrypt hashes from `nats.conf.A1` for each `__<USER>_PW__`. Do not run Helm.

### 4. Stage A2: leaf users get permissions (needs a pod restart)

From A2 there is no `leafnodes.authorization`, so any account user without
`allowed_connection_types` could open a leaf link on 7422. The spec gives the remaining
legacy users connection types from A2: `onemind` STANDARD and WEBSOCKET, `sys` and `monitor`
STANDARD (`TestConnectionTypes`). No current leaf uses them: while `leafnodes.authorization`
exists (today and at A1), only its user `leaf_tak` can log in as a leaf (`server/auth.go`,
`isLeafNodeAuthorized`). Before applying A2, confirm that no client of these users needs
another connection type (MQTT is not enabled):

```bash
curl -s 'http://127.0.0.1:8222/connz?auth=1&limit=1024' | python3 -c 'import sys,json,collections;print(collections.Counter((c.get("authorized_user"),c.get("kind","Client"),c.get("type")) for c in json.load(sys.stdin)["connections"]))'
```

nats-server 2.14.2 cannot reload a change of `leafnodes.authorization.users`
(`server/reload.go`: "config reload not supported for LeafNode: field Users"), and while that
block exists only its users can log in as leafs (`server/auth.go`, `isLeafNodeAuthorized`).
Proved locally in `TestActivationAndRollback` and `TestStageA1Additive`.

```bash
python3 scripts/nats_identities.py live --live-conf "$W/live.conf" \
  --passwords "$W/identities.json" --stage A2 --out "$W/nats.conf.A2"
# apply as in step 3; the reloader logs the refusal, the pods keep running A1.
```

Then restart one pod at a time, highest ordinal first, as `OPERATIONS.md` describes (the PDB
allows one voluntary disruption). After each restart, wait for Ready and for JetStream to be
current before the next one:

```bash
kubectl --context do-nyc3-omos-cloud -n nats delete pod nats-2
kubectl --context do-nyc3-omos-cloud -n nats wait --for=condition=Ready pod/nats-2 --timeout=300s
# nats-box, as the sys user: every stream current, no lagging replica
nats --context <sys> server report jetstream
nats --context <sys> stream report
# repeat for nats-1, then nats-0
```

During the roll the servers differ only in leaf authentication. The existing leafs reconnect as
`leaf_tak`, which is now an ONEMIND account user restricted to leaf connections. Afterwards,
move the TAK leaf to `leaf-tak` and Hermes to `hermes` (re-sealed `leaf.conf`, step 2), each
in its own PR. Check `leafz` shows each leaf with its new user and account:

```bash
kubectl --context do-nyc3-omos-cloud -n nats port-forward pod/nats-0 8222 &
curl -s 'http://127.0.0.1:8222/leafz' | python3 -c 'import sys,json;[print(l.get("name"),l["account"]) for l in json.load(sys.stdin)["leafs"]]'
```

Expected: the TAK leaf in `ONEMIND`, `hermes-legacy-leaf` in `AGENTS`, `zeus-mac` in `ONEMIND`.

### 5. Stage B: close the shared logins (config reload)

Only after `omos-tak-projector` uses `projector` and Hermes uses `hermes`; otherwise the
projector's `tak.cot.inject` is refused.

```bash
python3 scripts/nats_identities.py live --live-conf "$W/live.conf" \
  --passwords "$W/identities.json" --stage B --out "$W/nats.conf.B"
# apply as in step 3
for p in nats-0 nats-1 nats-2; do kubectl --context do-nyc3-omos-cloud -n nats logs $p -c nats --since=10m | grep -E 'Violation'; done
```

Every `Publish Violation` / `Subscription Violation` line names the client and subject. A
legitimate client on that list means roll back to A2 (below) and move that client first.
Acceptance checks (nats-box or a client device; expected output from the local run in
onemind-xr `docs/evidence/M1/A19/`):

```bash
nats --user zeus     pub cmd.fleet.drone.sim-x1.goto '{}'   # Permissions Violation for Publish to "cmd.fleet.drone.sim-x1.goto"
nats --user zeus     pub ent.drone.sim-x1.state '{}'        # Permissions Violation for Publish to "ent.drone.sim-x1.state"
nats --user onemind  pub tak.cot.inject '<event/>'          # Permissions Violation for Publish to "tak.cot.inject"
nats --user hermes   pub agents.hb.x '{}'                   # Authorization Violation (leaf-only credential)
```

### 6. Stage C: retire the shared users (later)

When `leafz` and `connz?auth=1` show no client on `onemind` or `leaf_tak` (NACK needs a
`nack` user with JetStream admin first), render `--stage C`, apply, and remove the passwords.
Removing `leaf_tak` is a leaf-user change only if it is still in `leafnodes.authorization`;
after A2 it is an account user and the change reloads.

## Rollback

| From | To | How |
|---|---|---|
| A1 | today | apply `nats-config.before.yaml` (`kubectl apply -f`), reload. Workloads already moved to new users must move back first |
| A2 | A1 | apply `nats.conf.A1`, then the same one-pod-at-a-time restart (leaf users again) |
| B | A2 | apply `nats.conf.A2`; reload restores the shared users at once (open connections get the new permissions without reconnecting, proved in `TestActivationAndRollback`) |

The reloader keeps the last good config running when a reload is refused, so a bad apply does
not take the bus down. A server that restarts reads whatever the ConfigMap holds: keep the
ConfigMap and the running stage the same before any pod restart.

## Adding a human

Asset actions at A4 (`asset.arm`, `asset.takeoff`) need two distinct humans: the second
approver must be a different principal from the first (asset-command.md). The live layout has one
human, `zeus`. Until browsers log in through Authentik and the auth callout (#565, where a human
is an Authentik user and no NATS user is configured), a human is a CLIENT user in
`identities.yaml`:

1. Copy the `alice` entry, rename it (lowercase; no `;`, no `ak:` prefix), remove
   `test_only: true` and add a `secret` location (for a person: their device keychain, as
   `zeus`). `publish: *client_human_publish` and `subscribe: *client_human_subscribe` keep the
   lists identical to `zeus`.
2. `python3 scripts/nats_identities.py values` and commit the spec and the three values files.
3. Generate the password (`gen-passwords --users <name>`), give it to that person, and add the
   bcrypt hash through `live` (step 3 below). Never commit it.
4. Map the user in the policy service's `principals.json` (`{"account": "CLIENT", "user":
   "<name>", "kind": "user", "id": "user:<name>"}`), in the layer service's principals (groups),
   and give `user:<name>` the OpenFGA `approver` / `second_approver` tuples for the action types
   that person may confirm (#554).

`alice` is the test-only second human. She has no credentials anywhere: `nats_identities.py`
renders her only with `conf --test-users`, which the permtest harness passes and the E2E harness
can pass. `check` refuses a `test_only` user outside CLIENT or with a `secret`, and the Python
tests fail if she shows up in a values file or a `live` render. `TestTwoDistinctHumans` proves
both humans reach `svc.policy.{command,task}.confirm` and that the policy service sees two
different server-set identities (`CLIENT/zeus`, `CLIENT/alice`). For a test, the policy
`principals.json` needs her mapped as `user:alice` (as in onemind-xr
`docs/evidence/INT/E2E/harness/main.go`).

## Local and stand-in Hermes (tests)

The `hermes` credential is LEAFNODE only (`allowed_connection_types: [LEAFNODE]`). A STANDARD
login as `hermes` gets `Authorization Violation` (`TestConnectionTypes`). The onemind-xr E2E run
worked around that by allowing STANDARD for its scripted stand-in. That skips the sidecar,
where Hermes' refusals are logged. Connect a local or scripted Hermes the way the pod does:

```bash
# disposable hub with a leafnode listener, e.g. leafnodes { host: 127.0.0.1, port: 17422 }
python3 scripts/nats_identities.py leaf-conf --user hermes --passwords "$W/passwords.json" \
  --out "$W/hermes-standin.conf" --remote 127.0.0.1:17422 --plain --port 14222 --http-port 18222
nats-server -c "$W/hermes-standin.conf"     # binds AGENTS through the leaf, loopback only
# the agent (scripted or real) connects with no credentials:
NATS_URL=nats://127.0.0.1:14222  <agent>
```

`--plain` drops the CA (plain `nats://` to a local hub) and listens on 127.0.0.1. `--port`
moves the client port off 4222, which a local nats-server may already use. Without these flags
the output is the pod's `leaf.conf` (step 2). The agent's identity at the hub is
`AGENTS/hermes`. A forbidden publish fails at the sidecar with `Permissions Violation` (client
error and sidecar log line). `TestHermesStandInSidecar` runs exactly this route.

## Limits that remain in stage 1

- Any user in an account that may subscribe `_INBOX.>` can read replies sent to other users'
  inboxes in that account. Per-identity inbox prefixes fix it (proposal A19-P5).
- A JetStream read grant is per stream: a consumer grant on `ENT` reads all of `ENT`.
- A leaf spoke drops a forbidden publish silently unless the spoke itself has permissions. The
  Hermes sidecar's rendered `leaf.conf` (step 2) gives the agent a local user with the hub
  allow-lists, so Hermes' attempts are refused and logged at the sidecar. The TAK leaf and the
  Mac site leaf have no such local user yet: there, a forbidden publish is still only filtered
  (nothing reaches the hub, no error, no log line).
- Stage B closes only the commands in policy's domains (`cmd.fleet/drone/robot/vehicle/mission/entity.>`),
  `tak.cot.inject`, `tak.cot.mesh.tx`, `ent.*.*.state`, `evt.so.action.>` and `svc.policy.>`
  for `onemind` and `leaf_tak`. Until stage C those shared logins (`onemind` also through the
  browser-exposed `ws.onemindos.dev`) can still publish other command domains:
  `cmd.projector.disable` (the projector's live kill switch), `cmd.ha.*`, `cmd.alert.*` and any
  future domain. They can also forge inputs that services trust: `tak.cot.out` (the entity
  service turns it into `ent` state, which the projector sends to TAK), `tel.fleet.drone.sim-x1.*`
  and `evt.drone.sim-x1.*` (policy reads `evt.drone.*` as command outcomes). Proposal A19-P10.
- Passwords, not nkeys or JWTs. Operator/JWT mode with revocation is stage 2 (onemind-xr
  BUS-DESIGN §7, doctrine §5).
- ONEMIND keeps unlimited per-account JetStream limits, bounded by the server's
  `max_file_store` (proposal A19-P7).

## Tests

CI runs all three on every PR that touches the identities (the Go tests with `-tags natsbin`
against the nats-server v2.14.2 release binary, checksum-verified in the workflow). Locally:

```bash
python3 scripts/nats_identities.py check
python3 -m unittest scripts/tests/test_nats_identities.py
cd platform/nats/identities/permtest
go test -race -count=1 -v ./...                      # nats-server 2.14.2 compiled in (go.mod)
NATS_SERVER_BIN=/path/to/nats-server-v2.14.2 go test -tags natsbin -race -count=1 -v ./...
# bcrypt like live (needs the nats CLI): NATS_IDENTITY_HASH=1 go test -count=1 ./...
```

The Go harness renders the configs with `conf --test-users` (two humans) by running
`scripts/nats_identities.py`, so `python3` must
have PyYAML (`pip install pyyaml==6.0.2`), or set `PYTHON` to an interpreter that has it.
Without it every Go test fails in setup with `ModuleNotFoundError: No module named 'yaml'`.

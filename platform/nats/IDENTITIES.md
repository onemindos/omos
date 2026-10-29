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

`python3 scripts/nats_identities.py check` must pass after any edit of the spec (CI runs it with
the tests: `.github/workflows/test-nats-identities.yml`).

## Identities

Only the essentials; `identities.yaml` has the full lists. `_INBOX.>` subscribe is granted to
every requester. Responders use `allow_responses`, so no service may publish on `_INBOX.>`.

| User | Account | Publishes | Subscribes / serves | Password goes to |
|---|---|---|---|---|
| `policy` | ONEMIND | `cmd.>`, `evt.so.action.>`, `svc.entity.get`, KV `geofences`, `mission-state` | `svc.policy.>`, `evt.{drone,robot,vehicle,mission,entity}.>` | SealedSecret `onemind-core/policy-service-runtime` (`NATS_USER`, `NATS_PASSWORD`; it also holds `OPENFGA_TOKEN`, #554) |
| `entity` | ONEMIND | `ent.*.*.state`, `evt.entity.>`, `evt.sys.entity.rejected`, KV `entity-correlation` | `svc.entity.>`, `cmd.entity.>`, `$SRV` (micro); consumes COT TEL WORLD INTEL DET ENT | `omos/omos-entity-nats` (`user`, `password`) |
| `ingest` | ONEMIND | `intel.>`, `evt.world.>`, `tel.feed.>`, `evt.sys.*.rejected`, KV `ingest-locks`, `ingest-state` | reads KV `feed-config` | `omos/omos-ingest-nats` |
| `layer` (A5, not built) | ONEMIND | KV `world-state`, `fleet-state` | `svc.layer.>` (streamed replies), `tel.feed.>`; consumes ENT TEL | `omos/omos-layer-nats` |
| `projector` | ONEMIND | `tak.cot.inject`, `tak.cot.mesh.tx` | `svc.cot.send`, `ent.*.*.state`, `cmd.projector.>`; consumes INTEL ENT | `omos/omos-tak-projector-nats` |
| `tak-client` (A21, not built) | ONEMIND | `tak.cot.out` | — | `takserver/tak-client-nats` |
| `bridge-mavlink-home` | ONEMIND | `tel.fleet.drone.sim-x1.>`, `evt.drone.sim-x1.>` | `cmd.fleet.drone.sim-x1.>` | edge `.env` / keychain (direct connection only) |
| `mesh-gateway` | ONEMIND | `tak.cot.mesh.rx`, `tel.fleet.edge.tak-mesh.status` | `tak.cot.mesh.tx` | edge `.env` / keychain (direct connection only) |
| `agent-fabric` | ONEMIND | `agents.prompt.>`, `$SRV.>` | `agents.>` | `agent-fabric/agent-fabric-nats` (`NATS_USER`, `NATS_PASSWORD`) |
| `leaf-tak` (leaf only) | ONEMIND | `tak.cot.out` | `tak.cot.inject` | inside `leaf.conf` of `takserver/tak-nats-leaf` |
| `leaf-site-mac` (leaf only) | ONEMIND | the site's bridges: `tel.fleet.drone.sim-x1.>`, `evt.drone.sim-x1.>`, `tak.cot.mesh.rx`, `tel.fleet.edge.tak-mesh.status` | `cmd.fleet.drone.sim-x1.>`, `tak.cot.mesh.tx`, `ent.*.*.state` | Mac leaf config (`infra/edge/mac-leaf`, keychain) |
| `hermes` (leaf only) | AGENTS | `agents.*.hermes.>`, `cmd.map.>`, requests `svc.entity.>`, `svc.layer.{list,meta}`, `svc.policy.{command,task}.propose`, `svc.policy.action.get`; replies `_R_.>` | `agents.prompt.hermes.>`, `svc.bridge.>`, `$SRV` | inside `leaf.conf` of `hermes/hermes-leaf-conf` |
| `zeus` | CLIENT | requests `svc.{layer,entity,policy}.>`, `svc.cot.send`, `agents.prompt.hermes.>`; `cmd.map.>` | `cmd.map.>` | owner keychain on each client device |
| `onemind`, `leaf_tak` | ONEMIND | everything until stage B; then denied commands, `tak.cot.inject`, `ent.*.*.state`, `evt.so.action.>`, `svc.policy.>` | everything; from stage B not `svc.{policy,entity,layer,cot}.>` (cannot impersonate a service) | existing `values.secret.yaml` |

Accounts: `ONEMIND` keeps every stream. It exports the policy API to `CLIENT` (all of it) and
to `AGENTS` (propose and action lookup only, never confirm), `svc.entity.>` to both, and
`svc.layer.>` to `CLIENT` with streamed responses. `CLIENT` and `AGENTS` have no JetStream.
`AGENTS` exports Hermes' prompts, `svc.bridge.>` and `$SRV` discovery to ONEMIND (and prompts
to CLIENT), its heartbeats and status to the ONEMIND stream `AGENTS`, and `cmd.map.>` to CLIENT.

## What changes for running clients

| Client | Today | After its move |
|---|---|---|
| Hermes pod | leaf as `leaf_tak` in ONEMIND. The hub forwards `tel.>`, `intel.>`, `tak.cot.>`, `ent.>`, `evt.>`, `det.>`, `cmd.alert.>`, `sys.>` to it, and it can publish anything (onemind-xr `docs/research/omos-live-state.md` §6) | leaf as `hermes` in AGENTS: no raw feeds, no `cmd.*` other than `cmd.map.>` (so no `cmd.ha.*` Home Assistant control, proposal A19-P2), entity data through `svc.entity.*`, actions through `svc.policy.*.propose` |
| TAK leaf | `leaf_tak`, unrestricted | `leaf-tak`: `tak.cot.out` out, `tak.cot.inject` in |
| `zeus-mac` leaf | `leaf_tak` | stays on `leaf_tak` (HA, cameras, legacy subjects) until its bridges are moved; `leaf-site-mac` is the target for the OneMind site bridges |
| `omos-tak-projector` | `omos-nats-creds` (`onemind`) | `projector`; must move before stage B or its `tak.cot.inject` is refused |
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
  --users policy,entity,ingest,layer,projector,tak-client,bridge-mavlink-home,mesh-gateway,agent-fabric,leaf-tak,leaf-site-mac,hermes,zeus
```

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
  (key `leaf.conf`). Hermes: same `leaf.conf` as `platform/hermes/agent/SECRETS.md` §2 with
  user `hermes`; commit it as `platform/hermes/agent/hermes-leaf-conf.sealedsecret.yaml`
  instead of the `kubectl create secret` in that file.
- Edge and client passwords (`bridge-mavlink-home`, `mesh-gateway`, `leaf-site-mac`, `zeus`)
  go to the owner keychain, not to the cluster.

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

## Limits that remain in stage 1

- Any user in an account that may subscribe `_INBOX.>` can read replies sent to other users'
  inboxes in that account. Per-identity inbox prefixes fix it (proposal A19-P5).
- A JetStream read grant is per stream: a consumer grant on `ENT` reads all of `ENT`.
- A leaf spoke drops a forbidden publish silently: Hermes' client gets no error, and nothing
  reaches the hub. Direct connections get `Permissions Violation` from the server.
- Passwords, not nkeys or JWTs. Operator/JWT mode with revocation is stage 2 (onemind-xr
  BUS-DESIGN §7, doctrine §5).
- ONEMIND keeps unlimited per-account JetStream limits, bounded by the server's
  `max_file_store` (proposal A19-P7).

## Tests

```bash
python3 scripts/nats_identities.py check
python3 -m unittest scripts/tests/test_nats_identities.py
cd platform/nats/identities/permtest && go test -race -count=1 -v ./...
# bcrypt like live (needs the nats CLI): NATS_IDENTITY_HASH=1 go test -count=1 ./...
```

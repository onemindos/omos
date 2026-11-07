# Tool gateway runbook (agentgateway)

Design: [`docs/architecture/onemind-tool-gateway.md`](../architecture/onemind-tool-gateway.md). State (verified live 2026-10-30):
**Phases 1 to 3d are live** (core, policy plane, Kubernetes, data, MediaMTX); the Phase 4 execute gate is merged as code but not deployed.
The gateway serves nine routes (`/ops`, `/dev`, `/cloud`, `/data`, `/tak`, `/home`, `/work`, `/identity`, `/execute`, all `/<group>/mcp`) and
six targets: `onemind` and `kubernetes` (`ops`), `stac`, `core` and `clickhouse` (`data`), `mediamtx` (`tak`). Per-backend runbooks:
[`tool-gateway-kubernetes.md`](tool-gateway-kubernetes.md), [`tool-gateway-data.md`](tool-gateway-data.md). `dev`, `cloud`, `home` and `work`
answer with an empty tool list until their backends join (they wait on owner credentials); `/identity` (deferred) and `/execute` (closed
until Phase 4 is deployed) always do. Hermes `onemind-hermes` uses `ops`, `data` and `tak` (Phase 5, config in `platform/hermes/onemind-hermes/`, live config updated with `mcp/apply.sh`); `/execute` is not wired until Phase 4 is deployed. The Hermes cron scripts call `/ops` too; LiteLLM no longer serves MCP (Phase 6).

## LiteLLM and in-cluster Ollama retired (2026-11-06)

Owner decision: agentgateway replaces LiteLLM as the hub, and Hermes uses Ollama Cloud (an external service, not the in-cluster Ollama). LiteLLM logged no traffic and its only consumer was `onemind-mcp`. This change removed from Git: Applications `ecosystem-litellm` and `ecosystem-ollama`, their values and chart-lock entry, the `litellm-clients` and `cpu-model-gateway-only` NetworkPolicies, `helm-oci-litellm`, the `onemind-model-clients-v1` bootstrap Job, the six `onemind-ai--litellm-*` sealed secrets, the `onemind-cpu-model` VolSync source, the `onemind-mcp` `ollama`/`litellm` adapters, env vars and egress rules, and the `litellm_models`, `ollama_models` and `ollama_running_models` gateway tools.

Left in place on purpose (owner decides later):

- The `litellm` database and role are being dropped: both are set to `ensure: absent` in `platform/ecosystem/resources/` (deleting the CR would not drop them: `databaseReclaimPolicy: retain` and the Application has `prune: false`). A follow-up PR removes the absent entries once the drop is verified. The live Secret `onemind-data/db-litellm` is not pruned by Argo; delete it by hand with owner approval.
- The DigitalOcean volume behind PVC `onemind-cpu-model` (20 GiB, `do-block-storage-retain`) stays; its restic repository secret `volume-aeb2c9c4e472-restic` is still in `platform/backups/volsync/sealed/` so existing backups stay restorable.
- `onemind-core` `model_summary` (capability `summarize-site`) is disabled: `MODEL_GATEWAY_URL` and `MODEL_GATEWAY_KEY` were removed from the sealed `core-runtime` secret, and without `MODEL_GATEWAY_URL` the summary returns a fixed "disabled" message plus the evidence. Set both keys to re-enable it against any OpenAI-compatible endpoint.
- Live cleanup (the parent does not prune; `ecosystem-litellm` carries a finalizer so deleting it cascades to its resources, `ecosystem-ollama` has none and would leave the PVC): delete Applications, the Deployments/Services, leftover Secrets and the PVC object by hand with owner approval; the PV is Retain.

## What runs

| Thing | Where |
|---|---|
| Argo CD Application | `platform-tool-gateway` (`platform/argocd/platform-tool-gateway.yaml`), three sources: the chart, `platform/tool-gateway/manifests/` (recursed), the repo as `$values` |
| Chart | `agentgateway-standalone` v1.7.0 (OCI `cr.agentgateway.dev/charts`, digest `sha256:f8e74f989c6705e3e1aaaaabacf554d594e2a24867eb4097b216e019a88c14a6`) |
| Image | `cr.agentgateway.dev/agentgateway:v1.7.0@sha256:60216aafb03d37f94e952ff3d51d2f630a65b10ecaebdc0a864882f92610c45b` |
| Pod | Deployment `agentgateway` in `tool-gateway`, 1 replica, on `onemind-ecosystem-cpu` |
| Endpoint | `http://agentgateway.tool-gateway.svc.cluster.local:8080/<group>/mcp`; metrics `:15020`, readiness `:15021/healthz/ready` |
| Client key | Secret `tool-gateway/gateway-clients` key `hermes`; the same value is sealed as `hermes/onemind-gateway-client` key `ONEMIND_GATEWAY_KEY` |
| Hand-written values | `platform/tool-gateway/values.yaml` |
| Generated values | `platform/tool-gateway/values.generated.yaml` (never edit; see below) |

## How the config is built

Helm replaces lists instead of merging them, so everything list-shaped (routes, targets, `extraEnv`, authorization rules) is generated:

```
platform/tool-gateway/gateway.yaml                          listener port, clients, groups, anchor   (shared; policy phase)
platform/tool-gateway/targets/<name>.yaml                   one fragment per backend target          (owned by the backend's PR)
platform/tool-gateway/policy/classification/<name>.yaml     one reviewed classification per target   (owned by the backend's PR)
        |  python scripts/gateway/render_gateway_config.py
        v
platform/tool-gateway/values.generated.yaml
```

CI (`scripts/tests/test_gitops_tool_gateway.py` and the `Tool gateway policy drift` step in `quality.yml`) fails when the committed output is
stale, when a classification file is invalid or has no target (or the reverse), or while a `proposed` list is unreviewed. After merging a
branch that changed `targets/` or `policy/`, re-run the renderer and commit the result.

### Policy: how tools become visible

- Every target has a classification file. Each tool is `read`, `write`, `destructive` or `physical` (design section 3.1). **A tool in no list is
  destructive: hidden and uncallable on every route.**
- A group route (`kind: group`, the default) gets one CEL `allow` rule per target covering its `read` and `write` tools. Nothing else is
  listed or callable (spike A1). A group with no exposed tool gets `require: "false"`, because a route with no rule lists everything.
- `destructive` and `physical` tools belong to `/execute/mcp`. That route is rendered with `require: "false"` (deny-all) and carries every
  target that has such tools (derived; a fragment never lists `execute`). Phase 4 adds visibility `allow` rules plus a call-only extAuthz rule
  (a `require` on `extauthz.decision` empties `tools/list`, spike A5). `identity` (`kind: deferred`) is deny-all with no targets.
- `anchor` (gateway.yaml) names a live target. A route with no live target answers `initialize` with a 500, so a group with nothing to
  expose borrows the anchor: it answers with an empty tool list because its deny-all rule hides and denies every tool.
- Every route keeps `apiKey.mode: strict`.

### Classification file format

```yaml
target: grafana                 # = file name; same as the fragment name
read: [search_dashboards, query_*]   # exact tool names (as the backend reports them, no <target>_ prefix) or globs ('*' only)
write: [create_annotation]
destructive: [delete_dashboard]      # execute route only
physical: []                         # execute route only
```

**There are no argument guards.** In agentgateway v1.7.0 a CEL predicate on `mcp.tool.arguments.*` inside an `allow` rule never fires
(verified against the pinned image by the Kubernetes backend work; details in [`tool-gateway-kubernetes.md`](tool-gateway-kubernetes.md)), so
a "destructive when `replicas == 0`" rule would expose the tool with no protection at all. The `guards:` key is therefore rejected: the
parser raises a `PolicyError` and CI fails. Classify a tool that is dangerous for some arguments by its worst case (destructive if any
argument can destroy), and put the argument-level rule in the Phase 4 execute gate, which sees the arguments.

A gated name that a read/write glob would swallow, or a name in two classes, is a render error. `proposed:` is written by the classifier and
must be deleted by a human after review.

`scripts/gateway/classify_tools.py` proposes the file. It calls `tools/list` on a backend and puts each new tool in a class by name
(`delete`, `remove`, `scale`, `restart`, `kill` and similar are destructive; `lock`, `unlock`, `door`, `siren`, `valve` and similar are
physical; `get`, `list`, `query` are read; `create`, `update`, `sync`, `annotate` are write; anything it does not recognise is destructive).
Those are proposals: read every one.

```bash
# discover: port-forward the backend (a gateway route only lists tools that are already classified), then
kubectl --context do-nyc3-omos-cloud -n <ns> port-forward svc/<backend> 18000:8000 &
python scripts/gateway/classify_tools.py --target grafana --url http://127.0.0.1:18000/mcp --write
#   with a credential: add --key-env GRAFANA_TOKEN (header "Authorization: Bearer <value>", read from the environment, never argv);
#   --header-env 'X-Api-Key={}' changes the header form
# review platform/tool-gateway/policy/classification/grafana.yaml, move tools to their real class, delete `proposed`

# drift check of a live route (exit 1 on new or stale tools): --url http://127.0.0.1:18080/ops/mcp --strip-prefix grafana_ --key-env GATEWAY_KEY
python scripts/gateway/classify_tools.py --target grafana --url ...

python scripts/gateway/classify_tools.py --check      # offline, what CI runs
```

OpenAPI targets: the generated tools are named after the operations (`list_items`, `delete_item_by_id`), which the name heuristics handle;
check each against the HTTP method in the spec (`GET/HEAD` read, `POST/PUT/PATCH` write, `DELETE` destructive) and move strays.

### OpenAPI targets

The gateway cannot filter the operations of an OpenAPI spec, so the spec it reads is trimmed in Git to the classified operations. Tool names
are the `operationId`s (`mediamtx_pathsList` on the route). Per target:

1. Write `policy/classification/<name>.yaml` with the `operationId`s (exact names, no globs); unlisted upstream operations stay out.
2. Generate `manifests/openapi-<name>.configmap.yaml` (ConfigMap `gateway-openapi-<name>`, key `openapi.yaml`) with
   `python scripts/gateway/trim_openapi.py --name <name> --source <upstream spec> --upstream "<repo vX.Y.Z path>" --server http://<svc>.<ns>.svc.cluster.local:<port> --classification platform/tool-gateway/policy/classification/<name>.yaml --out platform/tool-gateway/manifests/openapi-<name>.configmap.yaml`
3. Add a volume and mount for it to `values.yaml` (`extraVolumes`, `extraVolumeMounts`, mount path `/specs/<name>`) and point the fragment at
   `openapi: {schema: {file: /specs/<name>/openapi.yaml}, host: <svc>.<ns>.svc.cluster.local:<port>}`.
4. `scripts/tests/test_gitops_tool_gateway_openapi.py` fails when the spec, classification, fragment, mount, sealed secret or egress disagree.

A spec-only change does not roll the pod: `kubectl -n tool-gateway rollout restart deploy/agentgateway`.

### Target fragment format

```yaml
name: grafana              # [a-z][a-z0-9]*  (no "-", "." or "_"); the file must be <name>.yaml
groups: [ops]              # group routes (keys of gateway.yaml `groups`, kind group) that expose this target
env:                       # optional; every ${VAR} used in `target` must be declared here
  TARGET_GRAFANA_TOKEN:    # must start with TARGET_<NAME>_
    secret: gateway-target-grafana   # Secret in tool-gateway (a SealedSecret under manifests/)
    key: token
target:                    # verbatim agentgateway MCP target, WITHOUT `name`
  mcp: {host: http://mcp-grafana.tool-gateway.svc.cluster.local:8000/mcp}   # or sse / openapi
  policies:
    backendAuth:
      key: "${TARGET_GRAFANA_TOKEN}"    # "Authorization: Bearer <value>"
    # other headers: requestHeaderModifier: {set: {x-api-key: "${TARGET_X_KEY}"}}
```

Rules enforced by the renderer: a literal dollar is written `$$`; any bare `$` or an undeclared `${VAR}` is an error (the gateway expands the
whole config text and a missing variable is a startup error). The Host header sent to a backend is the target host without its port
(`requestHeaderModifier` cannot change it), so a backend that checks `Host` must allow the port-less name (see `ONEMIND_ALLOWED_HOSTS`).

## Recipe: add a backend (one PR per backend)

A Phase 3 PR adds only its own files and regenerates `values.generated.yaml`. It does not touch `gateway.yaml`, other targets, or the
single-writer files (`docs/STACK.md`, `platform/hermes/mcp/coverage.yaml`, `PLATFORM-ROADMAP.md`, `platform/backstage/catalog/apps.yaml`).

1. **Target fragment** `platform/tool-gateway/targets/<name>.yaml` (format above): `groups` per design section 4, `env` for each credential,
   `target` pointing at the backend Service (or the SaaS URL with `backendTLS`).
2. **Classification** `platform/tool-gateway/policy/classification/<name>.yaml`: run `classify_tools.py --write` against the backend, then
   review. Put deletes, scale-to-zero, drops, session kills, user removal and anything unrecoverable in `destructive`; locks, doors,
   garage, alarm, valves and `cmd.*` publishes in `physical`. A tool that is dangerous for only some arguments takes its worst-case class
   (guards do not work in v1.7.0 and are rejected). Anything unsure stays destructive. Give the provider token the least scope too (design section 3.1).
3. **Sealed secret** `platform/tool-gateway/manifests/gateway-target-<name>.sealedsecret.yaml`: one Secret named `gateway-target-<name>` in
   `tool-gateway`, namespace-wide scope, the `kubeseal` pattern in `CLAUDE.md`, value fed on stdin, never in argv. Keys must match the
   fragment's `env`. Tell the owner which token to create (design section 3.3).
4. **Egress** `platform/tool-gateway/manifests/egress-<name>.yaml`: a NetworkPolicy selecting the gateway pod
   (`app.kubernetes.io/name: agentgateway`) with egress to the backend only (copy `egress-onemind.yaml`), and the matching ingress rule in
   the backend's own namespace. For a SaaS target the rule is egress to 443 to the provider.
5. **In-cluster backend** (only when it runs in `tool-gateway`): Deployment, Service and NetworkPolicy in
   `platform/tool-gateway/manifests/backends/<name>/`, image pinned by digest, requests and limits, read-only root filesystem, restricted pod
   security. A backend that checks `Host` must accept the port-less service name.
6. **Render and test**: `python scripts/gateway/render_gateway_config.py`, then `python -m unittest discover -s scripts/tests -p 'test_gitops_tool_gateway.py'`,
   `python scripts/gateway/classify_tools.py --check`, `python scripts/validate_changed_yaml.py`, `python scripts/validate_manifests.py --all`,
   `python scripts/validate_helm_renders.py --all`. Commit `values.generated.yaml` with the PR; after a merge that overlaps another backend
   PR, rebase and re-render (the file is generated, never hand-merged).
7. **Verify live after merge** (Operate, below): the group's `tools/list` equals the read and write list of the classification; one read; one
   safe write if one exists; one destructive tool is absent from the list and a direct call returns `unknown resource`; `/execute/mcp` still
   lists nothing.

## Execute gate (Phase 4)

`/execute/mcp` lists the destructive and physical tools of every enabled target and asks the `execute-gate` (AuthZEN, `extAuthz`) before each
`tools/call`. Rules on the route: one `allow` per target (visibility), then `require: mcp.methodName != "tools/call" || extauthz.decision == true`.
The rule is keyed on the method because a `require` on `extauthz.decision` alone also runs for `tools/list`, where no decision exists, and
empties the list. Only the `execute` group has `authz:` in `gateway.yaml`; the renderer adds nothing to any other route.

Verified 2026-11-03 against the pinned image, rendered `values.generated.yaml` (backends replaced by a mock): `/ops` lists and calls normally;
`/execute` lists the gated tools; a call with the gate Service absent returns `403` and never reaches the backend; a destructive tool called on
`/ops` is `unknown resource`; no key is `401`. The AuthZEN body the gate receives is `{subject:{id:"hermes"}, resource:{id:"<target>/<tool>"},
action:{name:"tools/call"}, context:{arguments:{...}}}`. A tool called without `arguments` needs the `has()` fallback in the rendered body, or the
CEL errors and the call is denied before the gate is asked.

**Off until switched on.** The Deployment has `replicas: 0`: the pod needs three SealedSecrets (`gate-matrix` token and room_id, `gate-nats`
user, password and ca.crt, `ghcr` pull secret). `scripts/gateway/keys-session.sh` seals them and sets `replicas: 1`; a unit test fails if
replicas and the three files disagree. The gateway pod does not read any of them, so with the gate off `/ops`, `/data` and `/tak` are untouched.
The gate's policy file (`manifests/execute-gate/policy.configmap.yaml`) and the pod's `policy-sha256` annotation are generated by
`render_gateway_config.py` from the classification files; physical tools need `override <id>`, everything else `execute <id>`.

Live verification after switch-on (design section 7, Phase 4 (a) to (h)): call `kubernetes_resources_delete` on `tool-gateway/gateway-selftest`
through `/execute`; a request id appears in the room `approvals`; reply `execute <id>` as the owner; the retry succeeds once; a second retry and a
changed argument are denied; scale `execute-gate` to 0 and calls are denied. Then add the Hermes server:
`platform/hermes/onemind-hermes/mcp/apply.sh --include execute` (held back by default) and restart Hermes.

## Owner key session (execute gate and key-based groups)

Everything is in Git and switched on by sealing keys: no cluster command is typed by hand. The owner runs, in a terminal on a machine with
`kubectl` (context `do-nyc3-omos-cloud`) and `kubeseal`, from the repo root on a fresh branch:

```bash
scripts/gateway/keys-session.sh list                          # steps, output files, what each needs
scripts/gateway/keys-session.sh --dry-run all                 # rehearsal on a temp copy; nothing real is read or written
scripts/gateway/keys-session.sh ghcr nats gate                # the execute gate
scripts/gateway/keys-session.sh grafana github ...            # any subset; re-run one step to replace its key
```

Each step prints where to create the credential (UI path and the least-privilege scope), offers `s` to skip, reads the value with `read -s`
(no echo, no history, no argv), pipes `kubectl create secret ... --dry-run=client -o json | kubeseal --scope namespace-wide` and prints only
the file name and the length of each value. It never commits: review the diff, run the tests, commit, open the PR.

**What "switch on" means.** A target fragment with `enabled_by: [gateway-target-<x>]` is not rendered while that SealedSecret file is absent,
so a missing key cannot break a route or the pod. Its deployable manifests (backend Deployment, egress, rules in the app namespace) wait in
`platform/tool-gateway/pending/<x>/manifests/`. Sealing the key makes the script move them into `manifests/`, re-render
`values.generated.yaml` and (for the execute gate) set `replicas: 1` once `gate-matrix`, `gate-nats` and `ghcr` all exist. Unit tests fail
when a pending directory and its sealed file disagree.

The `gate` step runs `register_new_matrix_user` inside the Synapse pod for `@gate:onemindos.dev` (not admin, a random password that is never
displayed, read from the pod's own `/config/homeserver.yaml`, so the rotated Matrix Postgres password does not matter), logs in, creates the
private unencrypted room `approvals`, invites `@zeus:onemindos.dev` and seals `gate-matrix` (`token`, `room_id`). The owner then accepts the
invitation in Element. Retry after a half-failed run with `GATE_USER=gate2`.

Order of operations for the session:

1. Merge the execute-gate PR, then the targets PR (it also adds Argo CD account `legacy`, applied by `platform-argocd`; wait for Synced).
2. `keys-session.sh ghcr nats gate`, then the key-based steps the owner wants, in one branch; review the diff; run `python -m unittest discover -s scripts/tests`.
3. Commit and merge that PR. Argo applies the sealed files, starts `execute-gate` and the activated backends.
4. Verify the execute gate (section above, steps (a) to (h)), then `platform/hermes/onemind-hermes/mcp/apply.sh --include execute`, restart Hermes.
5. For each activated target: port-forward its backend (or use its URL) and run
   `python scripts/gateway/classify_tools.py --target <x> --url <mcp url> --key-env <VAR>`; review `proposed`, move tools to their class,
   delete `proposed`, re-render, PR. The committed classification is a conservative read-only start; it is not final.

Karakeep has a key step but no target: `@karakeep/mcp` is stdio-only and no image serves it over HTTP yet (`pending/karakeep/README.md`).
TAK Portal and the TAK Server client certificate are stored only: the portal API has no read operations and there is no TAK Server adapter.

## Operate

```bash
kubectl --context do-nyc3-omos-cloud -n tool-gateway get pods,svc
kubectl --context do-nyc3-omos-cloud -n tool-gateway logs deploy/agentgateway --tail=50     # JSON access log
```

Check from inside the cluster. The Hermes pod has no `kubectl` (it does have `ONEMIND_GATEWAY_KEY` in its environment since Phase 5), so read the Secret on the workstation and
pipe it in as a curl config on stdin (the key is never printed or put in argv):

```bash
KEY=$(kubectl --context do-nyc3-omos-cloud -n hermes get secret onemind-gateway-client -o jsonpath='{.data.ONEMIND_GATEWAY_KEY}' | base64 -d)
printf 'header = "Authorization: Bearer %s"\n' "$KEY" | kubectl --context do-nyc3-omos-cloud -n hermes exec -i deploy/onemind-hermes -c hermes -- \
  curl -s -K - -H 'content-type: application/json' -H 'accept: application/json, text/event-stream' \
  -d '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-06-18","capabilities":{},"clientInfo":{"name":"check","version":"0"}}}' \
  -D - http://agentgateway.tool-gateway.svc.cluster.local:8080/ops/mcp
unset KEY
```

Without the `Authorization` header the same call returns `401 authentication failed`. `tools/list` needs the `mcp-session-id` response header
from `initialize` and a `notifications/initialized` first.

## Rotate the client key

Generate a new value, seal it into both namespaces (`tool-gateway/gateway-clients` key `hermes`, `hermes/onemind-gateway-client` key
`ONEMIND_GATEWAY_KEY`) with the `kubeseal` pattern in `CLAUDE.md`, commit, and restart `agentgateway` (env is read at start). Use
`--from-file=<key>=/dev/stdin` fed by `printf`, so the value never appears in argv or the shell history.

## Known behaviors (Phase 0)

- Sessions are encrypted with a key the gateway generates at start; with one replica a restart invalidates open sessions and Hermes
  reconnects. Add `session.secretName` before running a second replica.
- `apiKey` defaults to a mode that lets requests without a key through. Every route is rendered with `mode: strict`; do not remove it.
- With one dead target in a group, `failureMode: failOpen` keeps the group answering; a call routed to the dead target returns an error.

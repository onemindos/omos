# onemind-hermes

Dedicated Mac SSH bootstrap, owner-local enrollment and pending acceptance gates:
[`mac-access/README.md`](mac-access/README.md). No authenticated Mac access is
claimed by the source wiring alone. The existing initial-config prepares owner-only
0400 credentials for the nonroot hermes account in a bounded RAM emptyDir; main
never mounts the raw Secret and no credential is copied to the PVC.

Legacy (`onemind-hermes`): the OneMind Hermes instance. It replaced the old `hermes-legacy`
Deployment, which was removed from Git on 2026-11-03 (see [`RETIREMENT-AUDIT.md`](RETIREMENT-AUDIT.md)), and the
`hermes-fresh` launch of October 26, 2026 (same image, same design, new names and new storage). It does not
import any existing agent data.

**Agent persona: Legacy, a single soul.** The pods carry
`onemind.io/agent: legacy`; the Deployment, Service and storage are named
`onemind-hermes`. The four profile souls from the old installation are not
separate agents here.

- Argo Application: `platform-onemind-hermes`, path `platform/hermes/onemind-hermes`.
- Shared objects it depends on are owned by `platform-hermes-base` (path `platform/hermes/base`): the
  `hermes` Namespace and the Secrets `hermes-dashboard-auth`, `omos-nats-creds` and `hermes-nats-ca`. They
  used to belong to `platform-legacy`, which no longer exists
  (`platform/hermes/base/README.md`). The tool gateway's `onemind` target, `onemind-mcp`, is
  `platform-onemind-mcp` (path `platform/hermes/mcp`).
- Deployment/Service/Ingress/NetworkPolicy: `onemind-hermes`, namespace `hermes`.
- Private dashboard: `https://onemind-hermes.tail717ba0.ts.net` once the
  Tailscale operator provisions it. Use the existing Hermes dashboard username
  and password.
- Image: official Hermes v0.23.0 (`nousresearch/hermes-agent:v0.23.0`), pinned by the multi-architecture index digest in all three containers (see [Upgrading Hermes](#upgrading-hermes)).
- New 25 GiB PVC `onemind-hermes-data` on `do-block-storage-retain`, so deleting
  the claim keeps the volume. No old volume is mounted or cloned; the unused
  `hermes-fresh-data` claim is not carried over.
- Initial config selects the existing preferred model/provider, but **no model
  credentials are imported**. Sign in to the model on the new instance.
- Native bundled skills/default features remain available. No custom plugins,
  MCP servers, channels, schedules, Legacy memory or conversation history are
  copied. `MIGRATION.md` lists what moves over and in which order.

## Tool gateway (MCP)

Legacy gets its tools from the agentgateway tool gateway (`platform/tool-gateway`, runbook
[`docs/operations/tool-gateway.md`](../../../docs/operations/tool-gateway.md)): servers `ops`, `data` and `tak`,
tools named `<target>_<tool>`. The Deployment passes `ONEMIND_GATEWAY_KEY` from Secret `onemind-gateway-client`;
`config.defaults.yaml` references it as `${ONEMIND_GATEWAY_KEY}` in the `Authorization` header. The default
profile reads it from the container environment; role profiles read it from their own `.env` (see
[Worker roles and the tool gateway](#worker-roles-and-the-tool-gateway)).

`config.defaults.yaml` only seeds an empty volume. To change the live `/opt/data/config.yaml` run
`mcp/apply.sh` (add `--dry-run` first): it sets every server in the defaults (`ops`, `data`, `tak`, `dropbox`,
`todoist`, `whop` and the nine Zoho entries) exactly as there, removes retired `clickup`, the old LiteLLM `onemind` entry and the
retired LiteLLM comment block, writes `config.yaml.bak-mcp-<timestamp>` before changing anything and touches no
other key. A second run prints only `unchanged`. Hermes reads `mcp_servers` at start, so restart the gateway
(`nats/apply.sh --restart`) afterwards. Check with `hermes mcp test ops|data|tak` inside the pod.

**Tool filters on the hosted servers.** `hermes mcp test` lists everything the server offers; the filters below
decide what the agent can call (`hermes mcp list` shows `N excluded` / `N selected`). Filter names are the
upstream tool names, exact or globs, and a name that does not match filters nothing: the earlier Dropbox
exclusions (`GetUsageAndQuota` and three more, in CamelCase) matched no tool.

| Server | Filter | Effect |
|---|---|---|
| `dropbox` | `exclude` | Reads and normal writes stay. Excluded: `delete*`, `restore_*`, `sign_cancel_*`, `run_write_dropbox_operation` (a generic writer that could reach a delete), plus the four quota/file-request tools the owner excluded before. |
| `whop` | `include` (200 of 464) | Reads, plus low-risk writes (messages, forum posts, reactions, notifications, course content, leads, notes, experiences, products, files and exports). Not listed: delete/remove, creating or changing payments, payouts, transfers, fees and refunds, cards and bank accounts, account and user changes, tokens and API keys. The reconciler refuses any such name in `mcp/whop/server.json`. |
| `zoho-*` | `include` | Read-only lists, unchanged (checked: no delete, remove, void or cancel tool). |
| `ops`, `data`, `tak` | gateway allowlist | `platform/tool-gateway`. `ops` includes two write tools by owner decision: `kubernetes_resources_create_or_update` and `onemind_tile38_set_point`. |

### Retired ClickUp MCP: bootstrap and recovery

ClickUp is no longer in the defaults. The `initial-config` init container removes
only `mcp_servers.clickup` through `hermes config unset` when present, after home
ownership repair and outside the empty-volume seed guard. Thus an existing or
restored older default-home config cannot reactivate ClickUp on the next start.
The manual MCP applier also removes that retired entry. Todoist remains the
owner's task authority and Dropbox the knowledge authority; neither is removed.

This config-only retirement does not read or delete OAuth files, revoke vendor
grants, change other profiles, migrate tasks, or delete a ClickUp account,
workspace or task. Previously removed local authentication is not recreated by
these manifests. Older encrypted backups may still contain historical OAuth
state; startup removes the config entry, not backup contents. The inactive
`onepassword/` proposal and historical audit records remain source-only/history,
not bootstrap inputs. Configuration exports are snapshots, not startup reconcilers.

Regression: `python -m unittest scripts.tests.test_gitops_retired_clickup_mcp -v`.
The startup test uses the actual installed Hermes CLI in an isolated home; it
mocks ownership operations and unrelated server reconcilers, not ClickUp removal.
These source tests do not establish deployment or a live backup restore.

## Worker roles and the tool gateway

Hermes v0.23.0 always multiplexes profiles: `gateway.multiplex_profiles: false` is retired and is rewritten to
`true` at start (`hermes_cli/gateway_multiplex_mode.py`), and `auto_multiplex_migration` only matters for a
single-profile install. So it cannot be turned off, and under it a role profile resolves `${ONEMIND_GATEWAY_KEY}`
in MCP headers from its own `.env` only (`agent/secret_scope.py`); the container environment is the default
profile's. Without the key every role logged `${ONEMIND_GATEWAY_KEY} ... is not set in this profile's .env` for
`ops`, `data` and `tak` about every ten minutes, and the role children had no gateway tools.

`roles/apply.sh` therefore gives each role (`roles/roles.yaml`, `mcp:` list) only its servers in `mcp_servers` and
the matching `mcp-<server>` entries in `platform_toolsets.cli` (Hermes accepts `mcp-<name>` or the bare name; a
listed name makes the list an allowlist of servers), and copies `ONEMIND_GATEWAY_KEY` from the pod environment
into the role's `.env` (mode 0600, never printed, no backup copy). After a key rotation run `roles/apply.sh`
again. ops-cluster gets ops and data, tak-mission tak and data, builder ops, researcher data, checker all three.

## Permanent Zoho MCP connections

The nine owner-provided Zoho endpoints are managed through strict-scoped
SealedSecrets and environment-only URL references. A Git-owned startup reconciler
uses the supported Hermes CLI to add only these nine entries on existing as well
as clean volumes, preserving other integrations. Initial exposure is an explicit
read-only allowlist, not all upstream actions. See [Zoho setup, verification and
recovery](mcp/zoho/README.md). Source preparation and live acceptance are separate.

## Permanent hosted Whop MCP

The existing hosted Whop MCP definition reconciles through the supported Hermes
config CLI on every start, scoped only to `mcp_servers.whop`. Private OAuth state
remains on the existing data PVC; its encrypted backups, not Git token snapshots,
are the authentication recovery path. No standalone Whop CLI install, startup
configuration or login is included. See [Whop MCP startup and recovery](mcp/whop/README.md).
Source tests do not establish live backup/restore success or authorize financial
actions.

## Permanent hosted Todoist MCP

Todoist is a major operational integration. The supported Hermes config CLI
reconciles only `mcp_servers.todoist` on every default-home startup, including
existing PVCs; empty PVCs also receive matching public defaults. OAuth stays
private on the existing PVC. Sampling is disabled, the working exclusions are
preserved, and no new tool allowlist, task migration or ClickUp replacement is
introduced. See [Todoist startup, recovery and pending backup acceptance](mcp/todoist/README.md).
Source tests are not deployment or restore acceptance.

## Native 1Password app credentials

**Source-only draft — NOT DEPLOYED.** Proposed native-source settings, pinned
`op` installer, reconciliation, runtime correction and tests remain under
[`onepassword/`](onepassword/README.md). Production Deployment, defaults and
Kustomization are byte-identical to the original base: no bootstrap injection,
installer, ConfigMap or native-source activation is added. The existing runtime
image is unchanged. Activation requires a built/tested corrected image pinned by
its real digest, a real rotated strictly scoped bootstrap SealedSecret, reviewed
GitOps wiring and authorized runtime/recovery acceptance. No new MCP/Connect
server is introduced.

## Labels and selectors

`app.kubernetes.io/name: onemind-hermes` is the selector for the Deployment,
Service and NetworkPolicy; it keeps Legacy's selectors
(`app.kubernetes.io/name: hermes` plus `onemind.io/agent: legacy`) separate.
Selecting on `onemind.io/agent: legacy` alone is ambiguous in principle, because the label names the persona, not
the instance. The only selector that did that, the `hermes-legacy` PodDisruptionBudget, was removed with the old
instance. Keep new selectors on `app.kubernetes.io/name`.

## Initial configuration (`config.defaults.yaml`)

The init container copies it once into the empty PVC with exclusive creation.
Later pod starts keep configuration authored through Hermes, so changing this
file after first start has no effect on the running instance. Every key was
checked against the v0.23.0 code in the image and passed its
`validate_config_structure`:

| Key | Value | Why |
|---|---|---|
| `model.provider` / `model.default` | `openai-codex` / `gpt-6.1-sol` | Same route as Legacy |
| `credential_pool_strategies.openai-codex` | `least_used` | Spread turns across signed-in Codex credentials instead of draining the first (`fill_first` default) |
| `fallback_providers` | `ollama-cloud` / `glm-5.3` | Used only if the primary fails. Served by the `ollama-pro` credential in Hermes' auth pool (`hermes auth list`), not by an `OLLAMA_API_KEY` variable; cron `infra-health` runs on it (`glm-5.3`, 2026-11-03) |
| `agent.disabled_toolsets` | `computer_use`, `browser`, `discord`, `discord_admin`, `yuanbao`, `feishu_doc`, `feishu_drive`, `x_search`, `video_gen` | Desktop, browser, chat-platform and media tools a cloud agent does not need |
| `approvals.mode` | `manual` | Dangerous commands wait for a human instead of the LLM judge (`smart` default) |
| `approvals.cron_mode`, `single_query_mode`, `unattended_mode` | `deny` | Cron, one-shot and API-server sessions never run flagged commands (matches the default; pinned) |
| `security.redact_secrets` / `allow_private_urls` | `true` / `false` | Defaults, pinned |
| `plugins.enabled` | `[]` | No custom plugins in the seed. The apply scripts add `legacy-roles` (`roles/apply.sh`) and, at the NATS cutover, `nats-platform` (`nats/apply.sh`) |
| `mcp_servers` | `ops`, `data`, `tak` | The tool gateway routes (`http://agentgateway.tool-gateway.svc.cluster.local:8080/<group>/mcp`), key `Authorization: Bearer ${ONEMIND_GATEWAY_KEY}`. Eager connect (`lazy` unset), `trust` default. `/execute` is not added until the execute gate is deployed |

`agent.disabled_toolsets` is applied last and per tool. It lists `yuanbao` (the
Yuanbao platform tools) rather than `hermes-yuanbao`: the latter is the whole
Yuanbao platform composite and includes terminal, file, web, skills, delegation
and kanban tools, so disabling it would remove toolsets this instance keeps.

The terminal backend stays `local`. The container is the sandbox: no
service-account token, no privilege escalation, `NET_RAW` dropped, and the
NetworkPolicy below. Hermes' Docker, Singularity, SSH, Modal and Daytona backends
either do not run inside a pod or send commands to external hosts.

## Isolation and network

The pod has no Kubernetes service-account token, NATS sidecar, observer
credential or Paperclip credential. The NetworkPolicy allows:

- ingress on 9119 (dashboard) from the `tailscale` namespace;
- egress to cluster DNS;
- egress on TCP 8080 to the tool gateway pod in `tool-gateway`
  (`app.kubernetes.io/name: agentgateway-standalone`), the MCP front door; its own
  NetworkPolicy admits this pod;
- egress on TCP 4222 to the NATS pods in `nats` (a direct client of the bus, see
  [NATS](#nats-legacys-bus-presence)); the `nats-clients` policy in
  `platform/nats/network-policy.yaml` admits this pod;
- egress to the Tailscale operator's proxy pods for the owner's devices, one rule per device and port
  (`omos-mbp` 22, `omos-s24` 8022, `omos-ha-edge` 8123, `omos-dell-zeus` 22), matched by the pod labels
  `tailscale.com/parent-resource` and `-ns: hermes`. The Services are in `platform/hermes/base/egress.yaml` and
  the tailnet policy `platform/core/policy.hujson` allows the same ports. The pod has no `tailscaled`;
- egress on 443 to public addresses, with private, cluster and tailnet ranges
  excluded.

The API is keyed independently and is not published by the Service or Ingress.
Dashboard session signing uses its own secret. Only the existing
`hermes-dashboard-auth` login Secret is shared. Plaintext credentials are not in
Git.

### Sealed secrets

Both SealedSecrets are namespace-wide scoped
(`sealedsecrets.bitnami.com/namespace-wide: "true"`), so they are bound to the
`hermes` namespace but not to a name. They were renamed by changing only
`metadata.name`, and the encrypted values were reused:

- `onemind-hermes-runtime` (was `hermes-fresh-runtime`): `API_SERVER_KEY`,
  `HERMES_DASHBOARD_BASIC_AUTH_SECRET`.
- `onemind-hermes-restic` (was `hermes-fresh-restic`): Restic storage
  credentials and password unchanged; `RESTIC_REPOSITORY` re-sealed with the new
  prefix (below) using `kubeseal --merge-into`, so no other value was decrypted.

Deployment uses one replica and `Recreate`. PVC annotations prevent Argo
prune/delete of its data.

## Scratch volume

Agent work (terminal commands, Go/uv/pip/npm caches, worktrees, throwaway databases) lives on its
own volume, `onemind-hermes-scratch` (50Gi, `do-block-storage`, reclaim `Delete`), mounted at
`/opt/scratch`. Before this, terminal commands started in the gateway's working directory on the data
volume. On 2026-11-12 an agent's build caches (4.7G in `/opt/data/agent-a-20261112`, plus 1.9G of uv
cache) filled the then-10Gi data volume, and `state.db` writes failed with `disk I/O error`.

- `TERMINAL_CWD=/opt/scratch/work` is where terminal commands start.
- `XDG_CACHE_HOME`, `GOPATH`, `GOCACHE`, `UV_CACHE_DIR`, `PIP_CACHE_DIR` and `npm_config_cache` point
  under `/opt/scratch`. They are not on Hermes's child-env blocklist, so terminal commands inherit them.
- `initial-config` creates `/opt/scratch/{work,cache}` owned by `hermes` (top level only).
- **Not backed up.** No ReplicationSource covers it, and deleting the claim deletes the volume. Push work
  to Git; nothing should exist only here.
- If it fills, agent builds get `No space left on device`; the gateway and `state.db` are unaffected.
  Clear `/opt/scratch/cache` freely.

**Memory is not isolated by this.** Agent commands still run inside the `hermes` container and share
its 3Gi limit with the gateway (the 2026-11-12 19:51 UTC OOM kill). Hermes's own per-command cap
(`TERMINAL_LOCAL_MEMORY_MAX_MB`) works through `systemd-run --user --scope`. This container has no
systemd and mounts `/sys/fs/cgroup` read-only, so that cap cannot take effect here. Real memory
isolation needs agent commands to run in a separate pod with its own limits, e.g. Hermes's `ssh`
terminal backend pointed at a sandbox pod. That is a follow-up.

## Backup

`ReplicationSource onemind-hermes-data` snapshots the PVC nightly at 01:25 UTC
using the existing VolSync mover/namespace policy, with explicit
`do-block-storage` snapshot and cache classes (temporary claims). Retention
follows Legacy's 7 daily / 4 weekly / 3 monthly / 1 yearly policy.

Its Restic repository is a new, independently initialized prefix,
`backups/restic/hermes/onemind-hermes-data` in the `onemind-backups` bucket. That
is the `backups/restic/<namespace>/<pvc>` layout every other VolSync source uses
and that `platform/backups/tools/mirror.py` expects. It is not shared with the
`hermes-fresh` repository or with the retired Legacy `hermes-data` repository (its restic data is kept). It shares
the existing storage credentials and recovery password; preserve the existing
off-cluster recovery password. First-backup completion and a restore test must
be observed; this source configuration alone does not establish recovery.

## Launch acceptance and gradual migration

After Argo reconciles, record the actual rollout, gateway health, dashboard auth
and tailnet access in the deployment PR. A ready pod is not proof of working
chat; model sign-in and a real response are a separate acceptance step.

Add selected memory, skill packages and integrations incrementally, following
`MIGRATION.md`. Inspect credential-bearing assets before importing them. Keep
Legacy running until the owner chooses the final cutover.

## Scheduled jobs

Four cron jobs run the daily brief, two cheap watchdogs and the nightly config snapshot.
Everything is in Git under `cron/` and applied by one idempotent script. Nothing is created
by hand in the pod.

| Job | Schedule (UTC) | Mechanism | Model |
|---|---|---|---|
| `legacy-sovereign-brief` | daily 07:00 | script gathers facts, model writes <180 words | `ollama-cloud` / `glm-5.3`, low effort |
| `infra-health` | every 15m | `--monitor-script`: silent unless the output changes | same, only on a change |
| `security-scan` | daily 06:30 | `--monitor-script`, same contract | same, only on a change |
| `config-export` | nightly 02:30 | `--no-agent`: the script is the job | none. **Paused**, see below |

All deliver to Telegram `8297858681`. The model jobs use the already-paid Ollama Cloud
subscription (the `ollama-pro` credential in the pool; verified with a one-token call), not
Codex, and never spend Codex quota. The old `so-infra-health` job (every 30 minutes, full model
run) hit `429 tokens_per_day`; here a healthy cluster costs zero tokens: the monitor script's
output is `OK`, its hash is unchanged, and Hermes skips the model. A problem changes the output,
wakes the model once, and stays silent until the output changes again (a new problem, or
recovery). The very first run after creation also counts as a change; the prompts answer
`[SILENT]` when everything is already `OK`.

Files:

| Path | Role |
|---|---|
| `cron/jobs.yaml` | Declarative job list: schedule, prompt, skills, model, delivery, `paused` |
| `cron/apply_jobs.py` | Applier, runs in the pod as `hermes`, through Hermes' own `cron.jobs` API |
| `cron/apply.sh` | Host wrapper: `cron/apply.sh [--dry-run]`. `--install-local` first copies this checkout's scripts into the pod, for testing before merge |
| `cron/scripts/legacy_common.py` | MCP client and Prometheus/Alertmanager/Argo helpers (stdlib only) |
| `cron/scripts/legacy_brief_facts.py` | Facts for the brief: nodes, bad workloads, volumes, alerts now and in 24h, Argo drift, cost, NATS, Codex usage |
| `cron/scripts/legacy_codex_usage.py` | Codex usage per credential (below) |
| `cron/scripts/legacy_infra_health.py`, `legacy_security_scan.py` | Monitor sources, print `OK` or one stable line per problem |
| `cron/scripts/legacy_config_export.py` | Snapshot, secret scan, commit and push |

How they reach the pod: `kustomization.yaml` builds the ConfigMap `onemind-hermes-cron` from
these files, and the `initial-config` init container copies them on **every** start (unlike
`config.yaml` they are Git-owned): scripts to `/opt/data/scripts` (the only directory the cron
runner accepts, symlinks out of it are refused) and `jobs.yaml` plus `apply_jobs.py` to
`/opt/data/cron-defs`. Changing a script therefore rolls the pod once. Jobs themselves live in
`/opt/data/cron/jobs.json` on the PVC (so VolSync backs them up) and are created only by
`cron/apply.sh`. A second apply prints `unchanged` for every job; fields that drifted are
updated, and `paused` in `jobs.yaml` is enforced both ways. Jobs not in the file are left alone.

Data sources and network: the pod has no Kubernetes token and egress only to DNS, the tool gateway
(`tool-gateway`, TCP 8080), Synapse and public 443. So the scripts call the read-only `onemind`
tools on the gateway's `/ops` route with `ONEMIND_GATEWAY_KEY` (`onemind_prometheus_query`,
`onemind_alertmanager_alerts`, `onemind_argocd_applications`, `onemind_nats_*`,
`onemind_opencost_allocation`; the gateway is stateful, so `legacy_common.py` runs
initialize and notifications/initialized before each call). Verified: Hermes' script runner passes that
variable to scripts (`build_subprocess_env`). The export's SSH goes to `ssh.github.com:443`,
which the public 443 rule already allows. Prometheus has no cert-manager or Falco metrics, so
the security scan reads Falco and certificate findings from Alertmanager only.

Approvals: `approvals.cron_mode: deny` governs commands an agent runs through its terminal tool.
The scripts run in the cron runner, not through that path, and the job prompts tell the model not
to run commands (the facts are already in the prompt), so nothing in these jobs can need a flagged
command. A flagged command would be denied, not asked.

### Codex usage per credential

`hermes usage --json` only reports the active credential. In Hermes v0.23.0
`agent/account_usage.fetch_account_usage(provider, base_url=, api_key=, read_only=)` queries
the ChatGPT usage endpoint for an explicit token, and `hermes_cli.auth.read_credential_pool()`
lists the pool entries (`codex-1`, `codex-2`). `legacy_codex_usage.py` calls the first for each
entry with `read_only=True`, the variant the model picker uses: a 401 is reported, never
repaired, and nothing is refreshed, rotated or written. Tokens are never printed. Output:

```
codex-1 (Plus): Session 0% (resets 2026-10-30 19:11 UTC), Weekly 31% (resets 2026-11-04 10:04 UTC)
codex-2 (Plus): Session 0% (resets 2026-10-30 19:45 UTC), Weekly 3% (resets 2026-11-04 15:00 UTC)
```

`--json` prints the full per-credential record. The brief script embeds it.

### config-export and the Git credential it needs

`legacy_config_export.py` writes an explicit allowlist, not `hermes profile export` (a 12 MB
tarball of the whole home, with sessions, hub caches and bundled skills): `config.yaml` with
secret-looking keys blanked and `${ENV}` references kept, `SOUL.md`, `memories/MEMORY.md` and
`USER.md`, `skills/index.json` (name, category, origin: bundled, user or external), the contents
of user-created skills, and the job definitions without run history. The whole snapshot is then
scanned (and scanned again on the exact tree about to be committed): token shapes, private-key
headers, credential assignments, URLs with embedded credentials or secret query parameters, bearer
tokens, and every secret value this
pod holds (environment and `auth.json`) verbatim. One hit aborts the run, names the file and rule
(never the value) and commits nothing. Tested here with a fake key and token (aborted) and
against a local bare repository (first run commits, second run is silent).

**Blocker found 2026-11-03:** the private repository `onemindos/onemind-hermes-config` exists, but
the GitHub organisation has deploy keys disabled (`deploy_keys_enabled_for_repositories: false`),
so adding a deploy key returns 422. The job stays paused until the owner enables deploy keys
for the organisation or chooses another credential.

onemind-hermes has **no Git write credential** (no token variable, no SSH key, nothing under
`HOME`), so `config-export` is created **paused** with that reason, and the pod mounts an
optional Secret for it. Owner steps to enable it:

1. Create a private repository for the snapshots, for example `onemindos/legacy-snapshots`.
   (Keeping agent commits out of `onemindos/omos` is deliberate; the snapshot contains memory
   text, so it must be private.)
2. Create a **deploy key with write access on that one repository** (preferred over a token:
   it cannot touch any other repository). Generate it locally,
   `ssh-keygen -t ed25519 -N '' -C legacy-config-export -f ./legacy-export-key`, add the `.pub`
   under the repository's Settings, Deploy keys, "Allow write access". Restrict the branch with a
   ruleset if wanted (default branch `main`).
3. Seal it, namespace-wide like the other onemind-hermes secrets, into
   `platform/hermes/onemind-hermes/export.sealedsecret.yaml`:

   ```bash
   kubeseal --fetch-cert --controller-name=sealed-secrets-controller \
     --controller-namespace=kube-system > /tmp/pub-cert.pem
   kubectl create secret generic onemind-hermes-export -n hermes \
     --from-file=ssh-privatekey=./legacy-export-key \
     --from-literal=repo=ssh://git@ssh.github.com:443/onemindos/legacy-snapshots.git \
     --dry-run=client -o json \
     | kubeseal --cert /tmp/pub-cert.pem --format yaml --scope namespace-wide \
     > platform/hermes/onemind-hermes/export.sealedsecret.yaml
   ```

   Add the file to `kustomization.yaml`, then delete the local key. An optional `branch` key
   (default `main`) can be added with `--from-literal=branch=...`. A fine-grained token over
   HTTPS would also work but needs the script changed; the deploy key needs no code change.
4. After Argo syncs (secret volumes are read at pod start, so the pod restarts once if it was
   already running), set `paused: false` for `config-export` in `cron/jobs.yaml`, merge, and
   run `cron/apply.sh`.

Until then the job stays paused; if it is ever run without the Secret it fails loudly rather
than silently skipping a backup.

## How Legacy works: doctrine and worker roles

One agent, one soul. GE, HP, LE and SO are modes of Legacy, not profiles. When work is big or parallel she
dispatches five worker roles on Codex and a `checker` verifies their output before she reports it. The doctrine
is the "How you work" section of `platform/hermes/seed/onemind-hermes/SOUL.md` (read every turn, about 400 more
tokens than before); the role definitions are `roles/roles.yaml`.

| Role | Model (`openai-codex`) | Toolsets | Job |
|---|---|---|---|
| `ops-cluster` | `gpt-6.1-sol` | terminal, file, web, skills, todo, session_search | cluster, Argo, K8s, NATS, Tailscale ops |
| `tak-mission` | `gpt-6.1-sol` | same | TAK, missions, CoT, video |
| `builder` | `gpt-6.1-sol` | + code_execution | code, manifests, branches, PRs (never merges) |
| `researcher` | `gpt-6.1-terra` | web, file, skills, todo, session_search | web and docs research |
| `checker` | `gpt-6.1-terra` | terminal, file, web, skills, todo, session_search | verifies other workers' output |

Models are those `provider_models_cache.json` lists for `openai-codex` in the pod. `terra` is the cheaper tier
(about half of `sol` per token in the models.dev cache); `luna` is cheaper still but too weak for review.

### What Hermes v0.23.0 actually supports (read from `/opt/hermes`; same as v0.22.7 unless noted)

- **`delegate_task` has no named roles.** The bundled `default` backend gives every child one route: the
  `delegation.provider` / `delegation.model` section or the parent's own (`plugins/delegation/default/__init__.py`
  `route_task`). `provider`, `model`, `role` and `toolsets` are reserved task names and not model-facing
  (`agent/delegation_provider.py` `RESERVED_TASK_PARAMETERS`; `role` is ignored, capability is depth-derived,
  `tools/delegate_tool.py:242`). A pluggable backend can add an enum parameter and answer with a route, including
  `profile` = run the child as that Hermes profile (`tools/delegate_tool_routes.py` `_profile_route`: the
  profile's own config, model, SOUL and toolsets). So a per-role model needs a backend plugin, and this PR ships
  one: `roles/plugin/legacy-roles/` adds `worker` (enum of the five) to every task; no `worker` = the default
  route. It is a user plugin (`$HERMES_HOME/plugins/`, `kind: backend`, gated by `plugins.enabled`,
  `hermes_cli/plugins_discovery.py` `gate_manifest`), selected by `delegation.backend: legacy-roles`.
  New in v0.23.0: the `default` backend is a bundled plugin (`plugins/delegation/default`, log line "Plugin 'default'
  registered delegation provider"); `legacy-roles` registers beside it and `worker: checker` still routes to the
  `checker` profile (scratch run: the child got the checker brief and ran a tool). A session with a later-result
  consumer, such as a NATS chat, now runs `delegate_task` in the background: the call returns "dispatched" at once
  and the result re-enters the chat as a new message.
- **`hermes kanban swarm` roles are profiles.** `--worker PROFILE:TITLE[:SKILL,..]`, `--verifier PROFILE`,
  `--synthesizer PROFILE` (`hermes_cli/kanban_parser.py:222`, `kanban_swarm.py`); the card's assignee must satisfy
  `profile_exists` or the dispatcher will not spawn it (`kanban_db_dispatch.py:2077`, `2329`). The worker runs
  `hermes -p <profile>` with that profile's `HERMES_HOME`, model and `platform_toolsets.cli`
  (`_default_spawn`, `_resolve_worker_cli_toolsets`). The dispatcher is embedded in the gateway (60 s tick; gateway
  log "kanban dispatcher: embedded in gateway"). `kanban swarm` is CLI-only (no tool), so Legacy runs it through
  her terminal tool.
- **Profiles share the Codex pool; no second sign-in.** `read_credential_pool` returns the default profile's
  `credential_pool` whenever the profile has zero entries for that provider (`hermes_cli/auth.py:948`, global-root
  `auth.json` read-only fallback `:507`); `_load_provider_state` does the same for OAuth state (`:836`), and a
  refresh that rotates the grant writes it back to the root store (`agent/credential_pool.py:724`). So role
  profiles carry no `auth.json`, and `credential_pool_strategies: {openai-codex: least_used}` in their own
  config spreads turns across both accounts. (`hermes profile create --clone` would copy the grant and fork it;
  `strip_cloned_single_use_oauth_grants` exists to undo that. The applier never clones.)
- **Where the soul comes from.** `/opt/data/SOUL.md`, a plain file. It is identical to
  `platform/hermes/seed/onemind-hermes/SOUL.md` (checked before this change); the init container does not seed it,
  `config.defaults.yaml` only seeds `config.yaml`, and `MIGRATION.md` step 2 copied it by hand. A role profile's
  `SOUL.md` is its whole brief (`roles.yaml`), so a role has no identity of its own.
- **Skills** are shared read-only: each role sets `skills.external_dirs` to `/opt/data/skills` and
  `/opt/data/onemind-skills`, so swarm cards that name a skill (`requesting-code-review`, `humanizer`) resolve
  without copying 56 MB per role.

### Applying it (idempotent)

`config.defaults.yaml` only seeds an empty volume, so a running instance gets this with one script:

```bash
platform/hermes/onemind-hermes/roles/apply.sh --dry-run   # shows create/update/unchanged
platform/hermes/onemind-hermes/roles/apply.sh
```

It pushes `roles/` and the repository `SOUL.md` to `/opt/data/roles-defs`, then runs `apply_roles.py` as `hermes`:
installs `SOUL.md` and the `legacy-roles` plugin; creates the five profiles (`--no-alias --no-skills`, no
credentials) with their model, `credential_pool_strategies`, `platform_toolsets.cli`, `skills.external_dirs`, and
the default profile's `approvals`, `security`, `agent.disabled_toolsets` and `mcp_servers` (copied raw, `${VAR}`
unexpanded, so it follows whatever the default profile has); sets `delegation.backend: legacy-roles` and adds
`legacy-roles` to `plugins.enabled`. Anything it overwrites is copied first to
`/opt/data/backups/legacy-roles-<UTC stamp>/`; a second run prints `unchanged` for everything. It never reads or
writes `auth.json` or any `.env`.

The `delegation` and `plugins` keys are deliberately **not** in `config.defaults.yaml`: on a fresh volume the
backend would be selected before its plugin exists, and an unavailable backend fails `delegate_task` closed.
After applying, the plugin loads at the next gateway start, and a new session (`/new`) is needed for the new soul
and tool schema. Rollback: `hermes config set delegation.backend default`, and restore `SOUL.md` from the backup
directory.

Workers run with `approvals.mode: manual`, which has no human in a worker session: a flagged command times out
and is denied, which is the safe direction. Workers get no authority Legacy lacks.

## NATS: Legacy's bus presence

Legacy is a **direct client of the cluster bus**, like every other cloud workload: no leaf sidecar (the
retired `hermes-legacy` sidecar was the one exception; `CLAUDE.md` "Bus Authority" no longer has one).
Inventory of what Legacy does on NATS, what was ported and the switch order are in
[`MIGRATION.md`](MIGRATION.md#nats-inventory-and-cutover-step-7). Summary:

- Identity `agents.{prompt,status,hb}.hermes.zeus.legacy` through the Synadia agent protocol plugin
  `nats-platform` (vendored in `platform/hermes/onemind-hermes/plugins/nats-platform/`; no start-up hot-patch, because
  Hermes v0.23.0 has no native NATS platform and the one line the patch fixed is now in the vendored file).
- Connection settings are in `deployment.yaml` (URL `tls://nats.nats.svc.cluster.local:4222`, user and password from
  Secret `omos-nats-creds`, CA from Secret `hermes-nats-ca`, identity, `HERMES_NATS_REFUSE_DUPLICATE=true`). They are
  inert until the plugin is listed in `plugins.enabled`.
- Python packages come from the `install-gateway-deps` init container (hash-pinned in `nats/gateway-deps.txt`,
  installed once to `/opt/data/.local/gateway-deps`).
- **Applying** (idempotent, backs up what it overwrites, restarts only the gateway process on request):

  ```sh
  platform/hermes/onemind-hermes/nats/apply.sh --status
  platform/hermes/onemind-hermes/nats/apply.sh --dry-run
  platform/hermes/onemind-hermes/nats/apply.sh --restart          # nats-platform only
  platform/hermes/onemind-hermes/nats/apply.sh --with-context     # also nats-context (off by default)
  platform/hermes/onemind-hermes/nats/apply.sh --disable --restart   # rollback switch
  ```

- `HERMES_NATS_REFUSE_DUPLICATE=true` stays: if another process ever holds the identity, onemind-hermes connects, sees it,
  and waits instead of registering. It held that guard while `hermes-legacy` still ran, which made enabling it first safe.
- `tak-platform` is not ported; PR #585 replaces it with a TAK-native client.

## Hermes Desktop: feature status

Audit of 2026-11-04 for the Mac app (`/Applications/Hermes.app`) connected to this instance. The app's
`connection.json` and `connections.json` point at `https://onemind-hermes.tail717ba0.ts.net` (OAuth mode) and agree
with each other. Method: the endpoints the app calls were extracted from its bundle and called in-process against this
pod's real data (`desktop/apply.sh --check`: status codes and key names only), and the pod logs were read. The UI itself
was not clicked, so "WORKS" means the server side answers and has a provider behind it, not that every screen was seen.
Re-run `desktop/apply.sh --check` after any change.

| Feature | Status | Cause or action |
|---|---|---|
| Sign-in, connection | WORKS | The 502s and ticket errors at 21:18 to 21:20 EST on 2026-11-03 were the pod restart (new pod started 02:20 UTC). |
| Chat, sessions, search, stats | WORKS | 497 sessions stored; every route answers 200. |
| Profiles (roles) | WORKS | `default` plus the five roles. |
| Skills (list, toggle, hub) | WORKS | 98 skills. Editing skills under `/opt/data/onemind-skills` FAILED (`PermissionError`): the tree was copied with Mac ownership (uid 501) and 120 AppleDouble `._*` files. Fixed by `desktop/apply.sh`. |
| Skill lock | WORKS from v0.22.7 | `PUT /api/skills/lock` did not exist in v0.22.5; v0.22.6 added it (`hermes_cli/web_routers/skills.py`, `tools/skill_lock.py`). Route present in the v0.22.7 scratch run; the toggle itself was not clicked. |
| Capabilities, dashboard plugins | WORKS | `kanban` and `hermes-achievements` are bundled dashboard plugins; they load whatever `plugins.enabled` says (that list gates agent plugins). |
| Kanban board | WORKS | Board, stats and workflow answer; the dispatcher runs in the gateway (60 s tick). |
| Mac-side plugins `hermes-project-stewardship`, `home-dashboard` | BROKEN (owner decision) | They live in `~/.hermes/desktop-plugins` on the Mac and call `/api/plugins/<name>/...`, which only exist if the matching server plugin is installed. The pod has none (404; the stewardship popups call is in the Desktop log). See decisions. `hermes-bots` only uses profile and cron routes, which work. |
| Cron | WORKS | 4 jobs; `config-export` stays paused (deploy-key blocker above). |
| Memory | WORKS | Built-in files only. `USER.md` is at 1353 of its 1375-character limit: a further addition will be refused. |
| Models, credentials | WORKS | `openai-codex` / `gpt-6.1-sol`, fallback `ollama-cloud`. |
| MCP management | WORKS | 16 servers listed; `ops`, `data`, `tak` answer `hermes mcp test`. `whop` OAuth re-prompts now and then and recovers by itself; Todoist reported its daily MCP call limit once (2026-11-03). |
| Files, image upload, git panel | WORKS | `git` is installed. `gh` is not, so "create PR" has no GitHub login to use (decision below). |
| Terminal | WORKS | Local backend only; Docker, Modal, Daytona and SSH backends are unavailable by design. |
| Usage, insights, logs, curator | WORKS | |
| Messaging platforms | WORKS | Telegram, Matrix, NATS, API server and webhook connected. Discord and Slack have no token and are not enabled; the `discord` toolsets are in `agent.disabled_toolsets`. |
| Image generation | WORKS (not exercised) | Provider `openai-codex` (gpt-image-2 through the signed-in Codex login); spends Codex quota. |
| Text-to-speech | FIXED, pending apply | Edge TTS needs the `edge_tts` package, which the image lacks and the read-only venv cannot install. Added to `nats/gateway-deps.txt` (hash-pinned; installed in a scratch directory in the pod and streamed audio on 2026-11-04). |
| Speech-to-text (dictation, wake word) | BROKEN (needs a key or a decision) | No provider: no `faster-whisper` and no OpenAI, Groq, Mistral, xAI or ElevenLabs key. Gateway log 2026-11-02: "Wake word needs speech-to-text configured". |
| Realtime voice (`voice-live`) | BROKEN (needs a key) | Wants `OPENAI_API_KEY`; the Codex login does not cover it. |
| Computer use | DISABLED BY DESIGN | `agent.disabled_toolsets: computer_use`. A headless Linux pod has no desktop. |
| Browser automation | DISABLED BY DESIGN | `agent.disabled_toolsets: browser` (no browser backend or key). |
| Video generation, X search, Discord, Feishu, Yuanbao | DISABLED BY DESIGN | Same list. |
| Local models | DISABLED BY DESIGN | A cloud pod has no GPU; `local-models` is off. |
| Approvals | WORKS, differs from the seed | Live `approvals.mode` is `smart` (the seed says `manual`) with a three-entry `command_allowlist`; set from outside Git. See decisions. |
| Webapp (Desktop in a browser tab) | WORKS from v0.23.0, owner to open | `https://onemind-hermes-web.tail717ba0.ts.net`, tailnet only. Scratch run: sign-in page, then the app bundle after login. See [The Webapp](#the-webapp). |
| Updates | MANAGED | The pod updates by image bump, not from the app (`/api/hermes/update/check`). |

**Version skew.** Closed by the v0.22.7 image bump (see [Upgrading Hermes](#upgrading-hermes)). The installed Desktop was
rebuilt on 2026-11-02 from the Mac's `~/.hermes/hermes-agent` at commit `f0d94882e89` (tag `v0.22.6`). Against v0.22.5 only
`PUT /api/skills/lock` was missing among the calls its Electron paths make; v0.22.7 has it. Still answering 404 on
v0.22.7 and in the `desktop/apply.sh --check` in-process run on v0.23.0: `/api/desktop-plugins/*`,
`/api/ui-preferences` and `/api/webapp-access/status`. They belong to the Webapp listener, not the dashboard on 9119. With
the Webapp on (v0.23.0), port 9120 answers `/api/ui-preferences` and `/api/webapp-access/status` with 401 until you sign
in, so they exist there. They are not regressions. The app's `update-check-cache.json`
was last written 2026-09-25 and is stale. The Mac also runs its own old gateway (`ai.hermes.gateway`, platform NATS only,
identity `legacy-edge`) and a stale `~/.hermes` profile; neither touches this instance.

**Owner decisions**

1. Speech-to-text: give the pod a cloud key (Groq or OpenAI; seal it like the other secrets) or approve local
   `faster-whisper` (about 200 MB of packages plus a model download, and CPU/memory on a memory-saturated node).
2. `OPENAI_API_KEY` for realtime voice, only if wanted.
3. Install the server halves of the Mac plugins? `hermes-project-stewardship` adds a `pre_llm_call` hook on every
   turn and writes `kanban.db`; both are third-party community code. Not installed without an audit.
4. Server image bump `v0.22.5` to `v0.22.7`: done. `v0.22.7` to `v0.23.0` with the Webapp: this change (see [Upgrading Hermes](#upgrading-hermes)).
5. `approvals.mode`: `smart` (live) or `manual` (seed). `manual` makes flagged commands wait for a person, which the
   Desktop approval prompt serves but Telegram, NATS and cron sessions cannot answer.
6. `gh` in the image plus a fine-grained GitHub token, only if Desktop "create PR" from the pod is wanted.
7. Mac app update to `v0.22.7-desktop`: report only, nothing was installed.

`desktop/apply.sh` is the only thing here that changes the live volume; `desktop/apply.sh --check` is read-only.

### The Webapp

Hermes v0.23.0 can serve Hermes Desktop in a browser tab. Here it is on, and only the tailnet can reach it.

- **URL:** `https://onemind-hermes-web.tail717ba0.ts.net` (Tailscale Ingress `onemind-hermes-web`, `service.yaml`). The
  dashboard stays at `https://onemind-hermes.tail717ba0.ts.net`. There is no nginx or public route, and the tailnet
  policy needs no change (`platform/core/policy.hujson` grant 1: your devices reach `*`).
- **Sign-in:** the same login as the dashboard (Secret `hermes-dashboard-auth`). A wildcard bind turns the sign-in gate
  on; without a session the listener answers a 302 to `/login`.
- **How it is on:** `webapp:` in `config.yaml` (`enabled: true`, `host: 0.0.0.0`, `port: 9120`, `public_url` as above).
  The gateway gives the Webapp its own listener, which replaces the private loopback listener it opened before. The API
  server (8642, the `/health` probe) and the dashboard (9119, the liveness probe) are separate and did not change in the
  scratch run. `config.defaults.yaml` seeds it on a new volume; a running instance gets it from `desktop/apply.sh`, which
  edits only the `webapp` block, backs `config.yaml` up first and prints `unchanged` on a second run:

  ```sh
  platform/hermes/onemind-hermes/desktop/apply.sh --dry-run
  platform/hermes/onemind-hermes/desktop/apply.sh --restart    # gateway process only, not the pod
  platform/hermes/onemind-hermes/desktop/apply.sh --check      # reports the config and what the listener answers
  ```
- **Turn it off:** `desktop/apply.sh --webapp-off --restart` (sets `webapp.enabled: false`). To remove the route as well,
  delete the `onemind-hermes-web` Ingress, the `webapp` Service port and container port, and the 9120 NetworkPolicy rule.
- The NetworkPolicy lets only the `tailscale` namespace reach 9120 (beside 9119). Memory in the scratch run was about
  620 MiB for the container with the Webapp on, so requests and limits are unchanged.

## Upgrading Hermes

How the image is bumped, what to check, and how to go back. Last done: v0.22.7 to v0.23.0 (2026-11-07).

**Where the version is pinned.** `deployment.yaml` (init containers `initial-config` and `install-gateway-deps`, and the
`hermes` container: all three must match) by the multi-architecture **index** digest, never a bare tag. The inert
`onepassword/` draft carries its own pins (`proposal/deployment.yaml`, `proposal/prepare_image.py`,
`runtime-fix/Dockerfile`, `runtime-fix/apply_source_patch.py`); keep them on the same version. The Hermes version is
also written in this README, `MIGRATION.md`, `config.defaults.yaml` comments, `docs/STACK.md` and the portal and
Backstage catalog descriptions.

**What to compare between versions** (do this against the two images, not the release notes: the v0.22.6 and v0.22.7
notes say only that the full notes ship with v0.23.0). Extract `/opt/hermes` from both images and diff what Legacy
depends on:

- plugin API used by `plugins/nats-platform`, `plugins/nats-context` and `roles/plugin/legacy-roles`
  (`PluginContext.register_platform/register_hook/register_tool/register_delegation_provider`,
  `gateway.run._profile_runtime_scope`, `_resolve_runtime_agent_kwargs`, `tools.approval_context`,
  `agent.delegation_provider`, `gateway.platforms.base`); an AST pass that imports every `hermes*`, `gateway`, `tools`
  and `agent` symbol the repo uses finds renamed or removed ones;
- `hermes_cli.config.atomic_config_write/atomic_config_replace/read_user_config_raw` (the apply scripts), `cron.jobs`
  (`apply_jobs.py`), `hermes mcp ...` (the three MCP reconcilers in the init container);
- the venv's Python (3.14.9 in both; the wheels in `nats/gateway-deps.txt` must still install), `/opt/hermes/tools/uv-*/uv`
  (the `install-gateway-deps` path) and the s6 slot of the gateway (`gateway-default`, used by `nats/apply.sh --restart`
  and `desktop/apply.sh --restart`; its directory is `/run/hermes-service` since v0.23.0);
- one-way migrations on `/opt/data` (below).

**v0.22.7 to v0.23.0: breaking changes found.** One that touches the repo: the gateway's s6 slot moved. Everything else
the repo calls is unchanged. Read from the two images and run in scratch containers on 2026-11-07:

| Area | Change | Effect here |
|---|---|---|
| s6 gateway slot | Per-profile gateway slots moved from `/run/service/` to `/run/hermes-service/`, supervised by a new s6-rc service `profile-gateways` (a nested `s6-svscan` that runs as `hermes`; `hermes_cli/service_manager.py` `S6_DYNAMIC_SCANDIR`). `/run/service/gateway-default` no longer exists. | `nats/apply.sh --restart` used the old path and would have failed. Changed to `/run/hermes-service/gateway-default`, the same for `desktop/apply.sh --restart`. In the scratch run `s6-svc -r` on the new path restarted the gateway as root; the old path printed "unable to control". No other `/run/service` use in the repo. |
| Role profile slots | At boot Hermes registers an s6 slot for every profile on the volume (`container_boot`), but starts only the default one. | The five role profiles show "down (not started yet)"; only `gateway-default` runs, so no extra memory. |
| Config schema | `_config_version` stays 52 (the same number in both images). No migration ran. | None. `state.db` is schema 32 in both. |
| Delegation | The `default` backend is now a bundled plugin; background delegation is the default for NATS chats. | `legacy-roles` still loads and `worker: checker` still routes to the checker profile. The result of a background `delegate_task` still cannot be delivered to a NATS chat: `Watch notification not admitted for nats chat=legacy: not_found` appeared again, so the known problem is **not** fixed in v0.23.0. |
| Security | The approval gate is stricter (per release notes). | `approvals.mode` unchanged. Not exercised beyond the code diff. |
| Webapp | New listener `webapp:` (`gateway/webapp_listener.py`). Off by default. | Turned on by this change, see [The Webapp](#the-webapp). |
| Plugin and config API | Unchanged: `gateway.run._profile_runtime_scope/_resolve_runtime_agent_kwargs/_resolve_gateway_model/_load_gateway_config`, `tools.approval_context`, `agent.delegation_provider`, `register_platform/register_hook/register_tool/register_delegation_provider`, `hermes_cli.config.atomic_config_write/atomic_config_replace/read_user_config_raw`. | The AST pass over 53 imported symbols found one missing in both versions: `tools.approval._get_approval_config` (the known case below). |
| Python, uv | Python 3.14.9, uv 0.13.5 at `tools/uv-0.13.5-linux-<arch>/uv` in both. The six hash-pinned wheels install and import. | None. |
| 1Password draft | The four upstream files `onepassword/runtime-fix` depends on hash the same in both versions. | The patch still applies; the draft's pins moved to v0.23.0. |
| Multiplex, MCP, cron | Source files unchanged for `gateway_multiplex_mode.py`, `secret_scope.py`, `account_usage.py`; `hermes mcp|cron|profile|config` answer. | None. |

**v0.22.5 to v0.22.7: breaking changes found.** None that touches Legacy's code, config keys or the apply scripts.
Specifics, all read from the two images on 2026-11-04:

| Area | Change | Effect here |
|---|---|---|
| Config schema | 51 to 52, run by the image's own boot hook (`docker/stage2-hook.sh` calling `scripts/docker_config_migrate.py` as the `hermes` user) before the gateway starts. It drops `browser.backend` and `browser.extension_control` and moves `AGENT_BROWSER_*` settings from `.env` into `browser.*`. | Legacy sets none of those. In the scratch run the only change to `config.yaml` was `_config_version: 51` to `52`. The hook writes `config.yaml.pre-docker-migrate.<stamp>` and `.env.pre-docker-migrate.<stamp>` to `/opt/data/backups/config/` first. |
| Browser | `agent-browser` and the built-in `browser_*` tools are gone; `browser_exec` is the only browser. The toolset is still called `browser`. | `agent.disabled_toolsets: browser` keeps it off. No key to change. |
| Features moved out of core | Camofox and GitHub Copilot left core in this window (Home Assistant, Spotify, Langfuse, `mcp serve`, Singularity and `import-agent` already had in v0.22.5). At start Hermes installs the plugin from the catalog **only for a home that uses the feature** (`hermes_cli/left_core_migration.py`). | Legacy uses none (Home Assistant stays an edge bridge on the bus), so nothing is fetched. If a boot ever logs a catalog install, a feature is in use that this list did not expect. |
| Delegation | A session now routes on the backend its `delegate_task` schema was built from, pinned to the session row; a config change applies from `/new` (or `/delegation <backend> --now`). The per-session override API (`set_session_backend`) is removed. | `legacy-roles` does not use it. After `roles/apply.sh` changes `delegation.backend`, a running chat keeps its old backend until `/new`. Route resolution for `worker: checker` returned the same profile home and credentials on both versions. |
| Cron | Additive: per-job `timezone`, `grace_minutes`, `catch_up`, `hermes cron cancel`, a `missed` and `cancelled` run status. The executions ledger table is rebuilt once to widen its status check. | `apply_jobs.py` passes only the fields it always did; its dry run and idempotent second run printed the same lines on both versions. |
| Plugin API | `PluginContext.has_capability` moved into a mixin (`PluginCapabilityCtxMixin`); `register_hook` gained a keyword; the `skill_namespace` manifest field became `data_namespace`. | The repo's plugins call only `register_platform`, `register_hook` and `register_tool`. `legacy-roles` logs a `capability_check ... tools.override decision=deny` line at load: expected, it overrides nothing. |
| MCP | Stdio children get `HERMES_HOME` pinned; large SSE events no longer fail; a parked server reports `needs_auth`. OAuth token files are read by a content digest, not only mtime. | No config key changed. `mcp_servers` keys the reconcilers write are unchanged. |
| Multiplex | Unchanged: `gateway.multiplex_profiles: false` is still rewritten to `true`. | None. |
| Security | The approval gate now asks for each content-level finding (credential upload, invisible Unicode) in a command as well as the first match, and redacts approval displays with a new helper. | Stricter prompts for flagged commands; `approvals.mode` unchanged. |
| Launcher, s6, uv | Byte-identical: `bin/hermes`, `hermes_bootstrap.py`, `entrypoint-dispatch.sh`, the s6 service tree, `container_boot.py`, Python 3.14.9, uv 0.13.5 at `tools/uv-0.13.5-linux-x64/uv`. | The init container paths and `nats/apply.sh --restart` (`/run/service/gateway-default`) are unchanged. |
| Python packages | `pydantic` 2.13.4 to 2.14.1, `aiohttp-socks` 0.11.0 to 0.12.0 in the image venv. | The four hash-pinned NATS wheels plus `edge-tts` and `tabulate` install and import unchanged. |
| Bundled skills | The boot hook re-syncs bundled skills into `/opt/data/skills` (about a dozen updated, `microsoft-office` added). | Own skills live in `/opt/data/onemind-skills` and are untouched. |

Not changed by the bump, found along the way: `plugins/nats-platform/adapter.py` imports `_get_approval_config` from
`tools.approval`, which has lived in `tools.approval_context` since before v0.22.5. The import fails inside a
`try`, so the approval round-trip timeout is the 300 s default instead of `approvals.gateway_timeout`. Same on both
versions; left alone here.

Also unchanged and visible only on the NATS surface: a background `delegate_task` finishing after the NATS prompt
returned cannot be delivered back (`Watch notification not admitted for nats chat=legacy: not_found`, then
"exhausted its 8 delivery attempts"). It reproduces on v0.22.5; the result stays queryable (`tools.async_delegation`).

**One-way state on `/opt/data`.** Two things change on the first start of the new image and neither blocks going back:

- `config.yaml` is stamped `_config_version: 52` after the backup above (v0.22.7 to v0.23.0: no change, already 52). v0.22.5 started and ran normally on a volume
  stamped 52 (checked in a scratch rollback run).
- The cron executions table is rebuilt wider; v0.22.5 read it and ran its scheduler fine in the same run.
- `state.db` keeps schema 32 in both versions (`hermes_state_common.SCHEMA_VERSION`), so session history is not migrated.

**Rollback.** Take a fresh snapshot before merging, then:

1. Before merging, take a snapshot that outlives the daily one. `onemind-hermes-data` is backed up by the VolSync source
   `onemind-hermes-data` (`backup.yaml`, daily 01:25 UTC; last sync before this PR: 2026-11-04T01:25:36Z, 36 s), so a restore
   point exists already; the pre-merge snapshot is the fresh one:

   ```sh
   kubectl --context do-nyc3-omos-cloud -n hermes apply -f - <<'EOF'
   apiVersion: snapshot.storage.k8s.io/v1
   kind: VolumeSnapshot
   metadata: {name: onemind-hermes-data-pre-v0-23-0}
   spec:
     volumeSnapshotClassName: onemind-backup-retain
     source: {persistentVolumeClaimName: onemind-hermes-data}
   EOF
   kubectl -n hermes get volumesnapshot onemind-hermes-data-pre-v0-23-0 -o jsonpath='{.status.readyToUse}'   # wait for true
   ```

   (`onemind-backup-retain` keeps the underlying snapshot if the object is deleted.) Delete it after the bump has run clean
   for a few days.
2. Revert: `git revert` the merge commit (it restores the v0.22.7 digest in all three containers and removes the Webapp; the leftover `webapp:` block in `config.yaml` still works on v0.22.7, which serves the sign-in page on 9120 inside the pod, but nothing routes to it; set it off with `desktop/apply.sh --webapp-off` first if you want it closed); Argo rolls the pod.
   The PV is `Retain`, the data is reused as is. No volume restore is needed for an image-only revert.
3. Only if the volume itself is damaged: scale `onemind-hermes` to 0 and restore `onemind-hermes-data` from the snapshot
   taken in step 1 (or the latest VolSync restic snapshot, `docs/portability/REBUILD-FROM-ZERO.md` step 9), then scale back
   up. Restoring is a last resort; it discards sessions written after the snapshot. Never recreate the live claim without
   owner approval (`CLAUDE.md`, deployment pattern item 6).
4. After a revert, `config.yaml` stays stamped 52 (harmless on v0.22.5). To put the pre-migration file back anyway, copy
   `/opt/data/backups/config/config.yaml.good.20261104-140645` (schema 51; the live upgrade did not write a `pre-docker-migrate` file) over `config.yaml` with the gateway stopped.

**Smoke test without the cluster.** The checks above were run in throwaway containers from the pinned images with a scratch
data volume, a local NATS server and a mock OpenAI-compatible model: the shipped init containers, every `*/apply.sh`
Python applier (dry run, apply, second run), real plugin discovery, the s6 start, `/health`, the dashboard status
route, `hermes mcp|cron|profile|config`, a NATS prompt and heartbeat, a model reply, a `delegate_task` with
`worker: checker` that ran a tool call in the child, and then the previous image started again on the migrated volume
(v0.23.0 run on a v0.22.7 volume, then v0.22.7 again: `/health` 200, dashboard 200, NATS connected, `_config_version` 52
both times). The v0.23.0 run also covered the Webapp (sign-in redirect, login, page after login), `desktop/apply.sh` steps
(dry run, apply, second run `unchanged`), the AST import pass and `s6-svc -r` on the new slot path. Each gateway start took
2 to 3 minutes there because the tool-gateway and OAuth MCP servers are unreachable in scratch. Not
covered there: real model providers, the tool gateway (`ops`, `data`, `tak`, `hermes mcp test`), Telegram and Matrix, the
Zoho, Whop and Todoist services, the cluster's own NATS TLS (a scratch CA was used), the Tailscale Ingress and operator proxy, the approval gate, and Desktop's UI. Those are the post-merge checks.

## Notes

- Ollama Cloud fallback: the key lives in the credential pool (`ollama-pro`), so no `OLLAMA_API_KEY` is needed in
  `onemind-hermes-runtime`. Add a second pool entry if a spare key is wanted.
- `ONEMIND_LITELLM_MCP_KEY` was removed from the sealed `onemind-hermes-runtime` Secret (2026-11-06); nothing read it
  since the LiteLLM MCP hub was retired (Phase 6).

## Cleanup after cutover

Removing `platform-hermes-fresh` from Git made the app-of-apps delete that
Application. It had no finalizer and `prune: false`, so its resources were
orphaned and then removed by hand.

Done (2026-10-29), with owner approval:

- Deployment, Service, Ingress and NetworkPolicy `hermes/hermes-fresh`, and the
  tailnet proxy for the `hermes-fresh` node.
- ConfigMap `hermes/hermes-fresh-defaults-*`.
- ReplicationSource `hermes/hermes-fresh-data` and its cache claim
  `volsync-src-hermes-fresh-data-cache`.
- SealedSecrets and Secrets `hermes/hermes-fresh-runtime` and
  `hermes/hermes-fresh-restic`.

Still pending:

- PVC `hermes/hermes-fresh-data` (`do-block-storage`, Delete reclaim policy; no
  sessions were ever recorded on it). Owner decision: keep for now. Do not
  delete it without a fresh snapshot and explicit owner approval.
- Remove the `hermes-fresh` device from the tailnet admin console.
- The `hermes-fresh` Restic prefix in the `onemind-backups` bucket, after the
  above.

## Agent name

The agent is called **Legacy** everywhere she is shown. `onemind-hermes` is only the
instance and infrastructure name. Hermes v0.22.7 has no agent-name config key, so the
name is set in three places: `SOUL.md` (identity), `API_SERVER_MODEL_NAME: Legacy` in
`config.yaml` (the id on `/v1/models`) and the Matrix display name of
`@legacy:onemindos.dev` (set once through the Synapse profile API). The Telegram bot
(`omos_hermes_bot`, first name `omos-hermes`) can only be renamed in BotFather. The old
`hermes-legacy` Deployment was the old Legacy instance, retired on 2026-11-03.

## Langfuse tracing

Active since 2026-11-07. Credentials come from the `hermes-langfuse` SealedSecret via `envFrom` (`HERMES_LANGFUSE_*`, project `hermes` on `langfuse-web.langfuse:3000`) and `HERMES_LANGFUSE_ENV=production`; egress is in `network.yaml`, ingress is `langfuse-web-clients` in `platform/ecosystem/resources/langfuse.yaml`. In v0.23.0 Langfuse is not bundled: the official catalog plugin `langfuse` (2.0.0, pinned to commit e7a218f2) is installed and enabled by the `install-langfuse-plugin` init container (idempotent, fail-open; fetches from GitHub and PyPI). Bumping the plugin means changing `SHA` there. `config.defaults.yaml` keeps `plugins.enabled: []` because the init container enables it on the live config. Langfuse v4 runs in events-only mode, so `/api/public/traces` returns 404: check with `/api/public/v2/observations?limit=1`. The plugin fails open.

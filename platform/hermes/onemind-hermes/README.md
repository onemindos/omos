# onemind-hermes

Legacy (`onemind-hermes`): the new OneMind Hermes instance, running beside the current `hermes-legacy`
Deployment. It replaces the `hermes-fresh` launch of October 26, 2026 (same
image, same design, new names and new storage). It does not import any existing
agent data.

**Agent persona: Legacy, a single soul.** The pods carry
`onemind.io/agent: legacy`; the Deployment, Service and storage are named
`onemind-hermes`. The four profile souls from the old installation are not
separate agents here.

- Argo Application: `platform-onemind-hermes`, path `platform/hermes/onemind-hermes`.
- Deployment/Service/Ingress/NetworkPolicy: `onemind-hermes`, namespace `hermes`.
- Private dashboard: `https://onemind-hermes.tail717ba0.ts.net` once the
  Tailscale operator provisions it. Use the existing Hermes dashboard username
  and password.
- Image: official Hermes v0.22.5, pinned to the verified multi-architecture digest.
- New 25 GiB PVC `onemind-hermes-data` on `do-block-storage-retain`, so deleting
  the claim keeps the volume. No old volume is mounted or cloned; the unused
  `hermes-fresh-data` claim is not carried over.
- Initial config selects the existing preferred model/provider, but **no model
  credentials are imported**. Sign in to the model on the new instance.
- Native bundled skills/default features remain available. No custom plugins,
  MCP servers, channels, schedules, Legacy memory or conversation history are
  copied. `MIGRATION.md` lists what moves over and in which order.

## Labels and selectors

`app.kubernetes.io/name: onemind-hermes` is the selector for the Deployment,
Service and NetworkPolicy; it keeps Legacy's selectors
(`app.kubernetes.io/name: hermes` plus `onemind.io/agent: legacy`) separate.
Because the pods also carry `onemind.io/agent: legacy`, anything that selects on
that label alone matches both Deployments. Today that is only the
`hermes-legacy` PodDisruptionBudget (`maxUnavailable: 1`), which then spans both
pods. Keep new selectors on `app.kubernetes.io/name`.

## Initial configuration (`config.defaults.yaml`)

The init container copies it once into the empty PVC with exclusive creation.
Later pod starts keep configuration authored through Hermes, so changing this
file after first start has no effect on the running instance. Every key was
checked against the v0.22.5 code in the image and passed its
`validate_config_structure`:

| Key | Value | Why |
|---|---|---|
| `model.provider` / `model.default` | `openai-codex` / `gpt-6.1-sol` | Same route as Legacy |
| `credential_pool_strategies.openai-codex` | `least_used` | Spread turns across signed-in Codex credentials instead of draining the first (`fill_first` default) |
| `fallback_providers` | `ollama-cloud` / `glm-5.3` | Used only if the primary fails. Needs `OLLAMA_API_KEY`, not provisioned yet (see TODO) |
| `agent.disabled_toolsets` | `computer_use`, `browser`, `discord`, `discord_admin`, `yuanbao`, `feishu_doc`, `feishu_drive`, `x_search`, `video_gen` | Desktop, browser, chat-platform and media tools a cloud agent does not need |
| `approvals.mode` | `manual` | Dangerous commands wait for a human instead of the LLM judge (`smart` default) |
| `approvals.cron_mode`, `single_query_mode`, `unattended_mode` | `deny` | Cron, one-shot and API-server sessions never run flagged commands (matches the default; pinned) |
| `security.redact_secrets` / `allow_private_urls` | `true` / `false` | Defaults, pinned |
| `plugins.enabled` / `mcp_servers` | `[]` / `{}` | MCP is added later, through LiteLLM |

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
- egress on TCP 4000 to the LiteLLM pods in `onemind-ai`
  (`app.kubernetes.io/name: litellm`), the planned MCP hub. LiteLLM's own
  `litellm-clients` ingress policy does not admit this pod yet; it needs a rule
  for `app.kubernetes.io/name: onemind-hermes` in namespace `hermes` before
  traffic flows;
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

## Backup

`ReplicationSource onemind-hermes-data` snapshots the PVC nightly at 01:25 UTC
using the existing VolSync mover/namespace policy, with explicit
`do-block-storage` snapshot and cache classes (temporary claims). Retention
follows Legacy's 7 daily / 4 weekly / 3 monthly / 1 yearly policy.

Its Restic repository is a new, independently initialized prefix,
`backups/restic/hermes/onemind-hermes-data` in the `onemind-backups` bucket. That
is the `backups/restic/<namespace>/<pvc>` layout every other VolSync source uses
and that `platform/backups/tools/mirror.py` expects. It is not shared with the
`hermes-fresh` repository or with Legacy's `hermes-data` repository. It shares
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

Data sources and network: the pod has no Kubernetes token and egress only to DNS, LiteLLM
(`onemind-ai`, TCP 4000), Synapse and public 443. So the scripts call the read-only `onemind`
MCP tools through the LiteLLM gateway with `ONEMIND_LITELLM_MCP_KEY`
(`onemind-prometheus_query`, `-alertmanager_alerts`, `-argocd_applications`,
`-nats_*`, `-opencost_allocation`). Verified: Hermes' script runner passes that variable to
scripts. No `network.yaml` change was needed; the export's SSH goes to `ssh.github.com:443`,
which the public 443 rule already allows. Prometheus has no cert-manager or Falco metrics, so
the security scan reads Falco and certificate findings from Alertmanager only.

Approvals: `approvals.cron_mode: deny` governs commands an agent runs through its terminal tool.
The scripts run in the cron runner, not through that path, and the job prompts tell the model not
to run commands (the facts are already in the prompt), so nothing in these jobs can need a flagged
command. A flagged command would be denied, not asked.

### Codex usage per credential

`hermes usage --json` only reports the active credential. In Hermes v0.22.5
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
scanned: token shapes, private-key headers, credential assignments, and every secret value this
pod holds (environment and `auth.json`) verbatim. One hit aborts the run, names the file and rule
(never the value) and commits nothing. Tested here with a fake key and token (aborted) and
against a local bare repository (first run commits, second run is silent).

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

## TODO

- Ollama Cloud fallback: seal `OLLAMA_API_KEY` into `onemind-hermes-runtime` (or
  a new namespace-wide SealedSecret) before relying on the fallback.

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
instance and infrastructure name. Hermes v0.22.5 has no agent-name config key, so the
name is set in three places: `SOUL.md` (identity), `API_SERVER_MODEL_NAME: Legacy` in
`config.yaml` (the id on `/v1/models`) and the Matrix display name of
`@legacy:onemindos.dev` (set once through the Synapse profile API). The Telegram bot
(`omos_hermes_bot`, first name `omos-hermes`) can only be renamed in BotFather. The old
`hermes-legacy` Deployment is the old Legacy instance, being retired.

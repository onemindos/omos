# hermes-legacy retirement audit

> **Status (2026-11-03).** This audit blocked deletion on 2026-11-01. Since then: the NATS identity moved to `onemind-hermes`, Paperclip
> was retired (`docs/operations/paperclip-retirement-2026-11-01.md`), the shared Secrets, egress Services and `onemind-mcp` were re-homed
> (`platform/hermes/base/README.md`, phases 1 and 2), and the owner approved removing the old instance without a soak. The removal PR
> deletes the Git content of `platform/hermes/agent` and `bridge`, the Application `platform-legacy` and the other items in section 11.
> The findings below are the 2026-11-01 state and are kept as the record; open items (credential rotation, the PV, the Released
> volumes) are listed in the PR.

Read-only audit, 2026-11-01 (about 05:00 to 05:40 UTC), cluster `do-nyc3-omos-cloud`.
Old instance: `hermes/hermes-legacy` (pod `hermes-legacy-6459d894d-gb74h`, containers `hermes`, `paperclip-api-tls`, `nats-leaf`).
New instance: `hermes/onemind-hermes` (persona Legacy, data `/opt/data`). During the audit its pod was replaced twice by another
agent's cutover work (`platform-onemind-hermes` was OutOfSync), so live facts about it are "as observed at 05:30 UTC".

Sources: live `kubectl` reads and read-only `exec` in both pods; Git state is `origin/main` at 53073089 (PR #613 merged).
The local checkout (`feat/tak-mesh-gateway`) predates `platform/hermes/onemind-hermes/` and `platform/hermes/skills/`, so those were
read with `git show origin/main:...`. Only this file was written. No secret values appear below; the config diff printed key names only.

Legend: **M** = migrated (proof on onemind-hermes), **D** = dropped on purpose (decision cited), **G** = gap (decision or work needed),
**IP** = in progress elsewhere.

## Verdict

Do not delete yet. Four things block deletion, in order of risk:

1. **Argo ownership.** `platform-legacy` owns far more than Legacy: the `hermes` Namespace, the `onemind-mcp` Deployment and Service
   (the tool gateway's `onemind` target), and three Secrets that `onemind-hermes` mounts (`hermes-dashboard-auth`, `omos-nats-creds`,
   `hermes-nats-ca`). It has `resources-finalizer`, `prune: true`, `selfHeal: true`. Deleting the Application cascades into all of them.
2. **Paperclip** has one configured agent ("Hermes") that calls `hermes-paperclip-api`, which is a Service on the old pod only.
   Nothing equivalent exists on `onemind-hermes` (by design today: "no Paperclip credential").
3. **NATS identity** is still held by Legacy (cutover in progress, #613).
4. **`hermes-data` PV has reclaim policy `Delete`** and its PVC is in the Argo app: pruning the PVC destroys the 13 GB volume. The only
   copy after that is the nightly VolSync restic repo.

Counts (by row below): **M 22, D 25, G 23, IP 1.**

---

## 1. Config (`/opt/data/config.yaml`, key by key)

Old file 19 KB, `_config_version` 46; new 4.6 KB, `_config_version` 51 (new relies on upstream defaults and Git-seeded
`config.defaults.yaml`). Secret-looking values were never printed.

| ID | Old key(s) | New | Status | Evidence / decision |
|---|---|---|---|---|
| C01 | `model.provider=openai-codex`, `default=gpt-6.1-sol` | identical | M | both configs |
| C02 | no fallback, single credential | `credential_pool_strategies.openai-codex=least_used` (2 Codex credentials), `fallback_providers[0]=ollama-cloud/glm-5.3` | M (improved) | new config; MIGRATION.md Summary row "Model" |
| C03 | `providers.digitalocean` (about 20 models), `DIGITALOCEAN_AI_KEY` | absent | D | owner memory (Codex x2 + Ollama, no Anthropic); `model_catalog.excluded_providers` on old already listed it. Key must be rotated (section 10) |
| C04 | `providers.bedrock` (about 14 models), `bedrock.region`, `AWS_*` env | absent | D | CLAUDE.md: "remove the AWS keys if confirmed unused"; skills README "no Bedrock-only" |
| C05 | `approvals.mode=smart` plus a 22-entry `command_allowlist` (permanently approved: recursive delete, `execute_code`, shell via `-c`, pipe-to-shell, force push, `reset --hard`, gateway kill/restart, and more) | `approvals.mode=manual`, `timeout=300`, `cron_mode/single_query_mode/unattended_mode=deny`, `security.redact_secrets=true` | D | README "Isolation" and approvals table; deliberate tightening. No explicit note about dropping the allowlist, but it contradicts the SOUL rule that deletes need `execute` |
| C06 | `security.allow_private_urls=true`, `browser.allow_private_urls=true`, `browser.use_real_profile=true` | `allow_private_urls=false`; `browser` in `agent.disabled_toolsets` | D | README line "Desktop, browser, chat-platform and media tools a cloud agent does not need" |
| C07 | `platform_toolsets.*` (hermes-<platform> bundles for 11 platforms) | `agent.disabled_toolsets` = computer_use, browser, discord, discord_admin, yuanbao, feishu_doc, feishu_drive, x_search, video_gen | M | new config |
| C08 | `skills.disabled` (2 entries incl. opencode), `skills.creation_nudge_interval=15` | `skills.disabled` (25 bundled), `skills.external_dirs=[/opt/data/onemind-skills]` | M | skills README "Bundled skills proposed for disabling" |
| C09 | `delegation.max_iterations=250`, `gateway.multiplex_profiles=true` with 4 profile souls | `delegation.backend=legacy-roles`, `plugins.enabled` has `legacy-roles`, profiles `builder, checker, ops-cluster, researcher, tak-mission`, `multiplex_profiles=true` | M (redesigned) | roles plugin in `platform/hermes/onemind-hermes/roles/` |
| C10 | `memory.memory_enabled/user_profile_enabled=true`, char limits 2200/1375, nudge 10, flush 6 | not set (upstream defaults); `memories/MEMORY.md` is 2199 bytes | M (defaults) | not verified that upstream defaults equal 2200/1375; the file size suggests so |
| C11 | `display.*` (skin, pet, tool_progress=all, reactions), `agent.personalities` (surfer, hype), `agent.system_prompt` (6.2 KB: a copy of the old "MOTHER" soul) | absent | D | cosmetic; the system prompt is the superseded soul (section 2) |
| C12 | `dashboard.theme`, `basic_auth.username/password_hash/secret`, `public_url=hermes-legacy...` | `HERMES_DASHBOARD_BASIC_AUTH_SECRET` from env, `envFrom` Secret `hermes-dashboard-auth`, ingress `onemind-hermes.tail717ba0.ts.net` | M (but see W06: Secret is owned by the old Argo app) | live pod env |
| C13 | `tts` (edge), `stt` (local small), `voice.auto_tts`, `wake_word.enabled=true` | absent | D | headless pod; no voice client. Not documented explicitly |
| C14 | `secrets.onepassword` (account, 3 env maps), `OP_SERVICE_ACCOUNT_TOKEN` in `.env` | absent | D | CLAUDE.md "Secrets: SealedSecrets"; skill `onepassword-legacy` is NEEDS-OWNER and not copied |
| C15 | `integrations.wordpress` (url, user) | absent | D | skills README (wordpress-* NEEDS-OWNER, credential to rotate) |
| C16 | `gateway.platforms.telegram.extra` (`fallback_ips`, `rich_messages`, `rich_drafts`, `status_indicator`), `telegram.reactions` | absent | D | cosmetic; Telegram is `connected` on new |
| C17 | `terminal.backend=local` (+container_* keys), `checkpoints`, `compression`, `tool_loop_guardrails`, `code_execution`, `session_reset`, `streaming`, `updates`, `tools.tool_search=on`, `prompt_caching` | not pinned; new env has `TERMINAL_*` and `BROWSER_*` timeouts | M (defaults) | a skim of the diff; none carries OneMind-specific value |
| C18 | `.env` keys `MCP_TODOIST_API_KEY` / `TODOIST_API_TOKEN`, `LMSTUDIO_API_KEY`, `LM_API_KEY`, `TAK_LEGACY_*`, `TAK_CHAT_SUBJECT_*`, `NATS_URL`, `TELEGRAM_HOME_CHANNEL_THREAD_ID` | absent | D | no Todoist MCP or skill exists on either side; TAK keys go with tak-platform (P03); rotate the keys (section 10) |
| C19 | Telegram and Matrix tokens, allowlists | `TELEGRAM_BOT_TOKEN`, `MATRIX_ACCESS_TOKEN` in `onemind-hermes-runtime`; `TELEGRAM_ALLOWED_USERS` and `MATRIX_ALLOWED_USERS` pinned to Zeus only; `GATEWAY_ALLOW_ALL_USERS=false` | M | MIGRATION "Chat channels, moved 2026-10-29" |

Env var names present only on old (process env): `AWS_*`, `DIGITALOCEAN_AI_KEY`, `LMSTUDIO_API_KEY`, `LM_API_KEY`, `ONEMIND_MCP_TOKEN`,
`ONEMIND_MCP_TOOLSET_REVISION`, `TODOIST_API_TOKEN`, `HERMES_DISABLE_LAZY_INSTALLS`, `HERMES_LAZY_INSTALL_TARGET`, `notesPlain`, `password`,
`username` (the last three are stray names, probably an env map from a Secret). Only on new: `API_SERVER_*`, `HERMES_NATS_*`, `NATS_*`,
`ONEMIND_GATEWAY_KEY`, `ONEMIND_LITELLM_MCP_KEY`, `MATRIX_ACCESS_TOKEN`, `TELEGRAM_BOT_TOKEN`, `HERMES_DASHBOARD_BASIC_AUTH_SECRET`.

## 2. Identity and memory

| ID | Item | Status | Evidence |
|---|---|---|---|
| I01 | `SOUL.md`: old 4.2 KB "Legacy, the MOTHER" with edge nodes (`legacy-edge`); new 11.4 KB single-agent soul (mission, four areas SO/HP/LE/GE, principles, rhythm, communication) | M | new SHA256 prefix `d59a1b04b363` equals `platform/hermes/seed/onemind-hermes/SOUL.md`. The mother/edge-node framing is dropped on purpose (MIGRATION: "single Legacy soul") |
| I02 | `memories/MEMORY.md` (old 2201 B, 9 entries; new 2199 B, 10 entries) | M | rewritten and scrubbed (stale TAK public LB, Mac SSH hints, Zeus's WordPress detail removed; tailnet-only TAK, Notion routing, rules added). New file differs from the Git seed because the agent edits it |
| I03 | `memories/USER.md` (old 1378 B, new 1574 B) | M | new SHA256 prefix `1312d69e0855` equals the Git seed; carries every old fact (geo stack, action-first, Notion, OneTAK naming, empire, Matrix) |
| I04 | Profile souls `generational-empire`, `holistic-performance`, `legacy-evolution`, `sovereign-operations` (1.9 to 2.0 KB each, 0 sessions, no skills or cron); 10 empty profile dirs (eden, empire, forge, grid, guardian, haven, heritage, oracle, spartan, trinity) | D | MIGRATION Summary: "No: single Legacy soul; profile souls stay in Git" (`platform/hermes/agent/agents/{ge,hp,le,so}`); the four areas live in the new SOUL and the `onemind-domains` skill |
| I05 | Doctrine and reference files that exist only in Git under `platform/hermes/agent/`: `doctrine/{CREED,PILLARS,PROFILES,EVOLUTION}.md`, `FABRIC-REFERENCE.md`, `LEGACY-ARCHITECTURE.md`, `RUNBOOK.md`, `agents/` (soul per area, `legacy-edge`, 21 under `planned/`), `plugins/{nats-platform,nats-context,onemind-vision}` | **G** | None of it is on the new pod. `nats/apply.sh` and the seed README read `plugins/nats-platform` from `platform/hermes/agent/`, and the MEMORY seed names that path as the home of souls and runbook. Removing the Argo Application does not delete these files, but removing the directory (the natural cleanup) does. **Action:** keep the directory as documentation, or move it (`platform/hermes/onemind-hermes/` or `docs/`) in the same PR that removes the manifests, and fix the references (`nats/apply.sh`, seed README, `.github/workflows/hermes-identity-export.yml`, MIGRATION) |

Old `MEMORY.md` (and its `.pre-edit-20261016`, `.pre-scrub-20261029` backups) contains a plaintext WordPress application password. See C-list in section 10.

## 3. Skills

Old: 202 `SKILL.md` (164 active or stale, 38 under `skills/.archive`). New: 94 (35 OneMind skills in `/opt/data/skills`'s external dir, the rest bundled).
`platform/hermes/skills/` in Git equals the live `/opt/data/onemind-skills` (35 = 35, diff empty).

| ID | Group | Status | Detail |
|---|---|---|---|
| S01 | 35 shipped skills (doks-cluster-ops, nats-operations, tak-server, tak-mission-ops, tailscale-pod-access, composio, onemind-stack, etc.) | M | skills README "Shipped here (35)"; `hermes skills list` on new shows 36 local, 53 builtin |
| S02 | Role skills so-grid, so-guardian, so-oracle, so-trinity, ge-empire, hp-spartan, le-eden, le-haven, le-heritage, so-forge; plus composio-mcp-patterns (folded into composio) and notion-onemind-architecture (renamed `onemind-architecture`) | M | merged into `onemind-domains` and `onemind-code-orchestration` |
| S03 | DROP list (17): omos-ui-build-deploy, legacy-integrations, legacy-model-config, cloudflare-access-patterns, composio-clickup, composio-github, composio-gmail, composio-google, omos-cortex-neurolinked, openclaw-k8s-deploy, repo-absorption, resume-session-on-pod, onemind-xr-platform, onemind-product-strategy, merge-reconciler, hermes-mobile-desktop-sync | D | skills README "DROP" table (reasons: retired things, embedded secrets). **Usage warning:** `onemind-product-strategy` was used 24 times (last 2026-10-25) and `legacy-integrations` 25 times despite being DROP; confirm the owner is fine with it |
| S04 | Upstream optional skills not reinstalled: ascii-art, audiocraft, blogwatcher, comfyui, evaluating-llms-harness, excalidraw, github-auth/-code-review/-issue-to-pr/-issues/-pr-workflow/-repo-management, heartmula, huggingface-hub, jupyter-live-kernel, llama-cpp, lm-evaluation-harness, openhue, petdex, polymarket, pretext, research-paper-writing, segment-anything, serving-llms-vllm, sketch, touchdesigner-mcp, vllm, weights-and-biases, yuanbao | D | MIGRATION "Upstream optional / hub skills: reinstall only if wanted". The `github-*` six are the ones most likely wanted once the GitHub MCP question (M02) is settled |
| S05 | 38 archived skills in `skills/.archive` | D | MIGRATION "stay behind" |
| S06 | NEEDS-OWNER, not copied (13): windows-network-drop, windows-zeus-access, mac-cloud-sync, mac-storage-ops, wordpress-legacy, wordpress-on-digitalocean, zeus-site, onepassword-legacy, notion-workspace, vault-dedup-audit, external-drive-cloud-sync, modesty-ai-daily-briefing, grounddata-ai-daily-news | **G** | The owner call is still pending in the README. Real usage in `.usage.json`: notion-workspace 139 uses (last 2026-10-28), zeus-site 31, vault-dedup-audit 25, onepassword-legacy 18, android-adb-management 52 (imported). **Action:** owner marks each keep or drop; for "keep", copy the SKILL.md (scrubbed of secrets; README lists which hold a WordPress app password, Notion token, embedded tokens in grounddata scripts) into `platform/hermes/skills/` before the volume goes |
| S07 | Skills in neither the README nor MIGRATION: `miro-mcp`, `open-source-frontend-evaluation` | **G** (minor) | tie `miro-mcp` to the Miro decision (M05); classify the other |

## 4. Plugins (6 on old)

| ID | Old plugin | New status | Status |
|---|---|---|---|
| P01 | `nats-platform` | vendored in `/opt/data/plugins/nats-platform`, listed in `plugins.enabled`, `nats` platform state `retrying: identity hermes:zeus:legacy is held by another live instance` | IP (N01 below) |
| P02 | `nats-context` | vendored in Git (`platform/hermes/agent/plugins/nats-context`), **not enabled**, not on the pod | D for now: MIGRATION "Owner decisions": wait for PRs #562/#564 (raw feeds and generic `nats_publish` are removed on purpose). It is what injects the `[LIVE OPERATIONAL PICTURE]` block |
| P03 | `platforms/tak-platform` | not ported | D: MIGRATION "Not ported... PR #585 replaces it". Side effect of retirement: Legacy's 30 s presence beacon to TAK (callsign) stops until #585 |
| P04 | `onemind-maven` (directory `onemind-vision`, 20 KB, vision query, anomaly detection, SITREP, COA) | in Git, not on the pod | **G** (minor): MIGRATION says only "already in Git... not part of this step"; no decision to install or drop. Not used (no tool-call evidence checked) |
| P05 | `memory-wiki` (dashboard tab, MIT), `hermes-project-stewardship` (40 MB with its own venv; `plugin-data/hermes-dockyard/dockyard.db` has schema only, no rows) | not on the pod | **G** (minor): third party, MIGRATION lists them but records no decision. Recommend: drop, record it |
| P06 | (new only) `legacy-roles` | enabled | n/a |

## 5. MCP servers (13 on old)

| ID | Server(s) | On new | Status |
|---|---|---|---|
| M01 | `onemind` (direct to `onemind-mcp:8000`; was failing 401 on old) | via tool gateway `ops` group (the `onemind` + `kubernetes` targets): `hermes mcp test ops` connected, 38 tools. Also `data` (14 tools) and `tak` (25 tools) groups, which old never had | M |
| M02 | `composio`, `github` (both enabled and working on old) | **not configured** | **G**: MIGRATION step 4 says "Add `composio` and `github` MCP"; skill `composio` is shipped and refers to a server the pod does not have. Gateway `dev`/`work` groups exist as routes but are not wired for Hermes. **Action:** decide, then add as `mcp_servers` entries (OAuth for GitHub, Composio consumer key as SealedSecret) or through the gateway |
| M03 | `cloudflare`, `raindrop`, `miro` (enabled, OAuth tokens present in `/opt/data/mcp-tokens`) | not configured | **G**: MIGRATION "Open decisions: whether miro, cloudflare and raindrop MCP are still wanted". Owner memory says agents reach Cloudflare through the official MCP. **Action:** decide; OAuth tokens are on the old volume only and cannot be copied (re-authorise on new if kept) |
| M04 | `notion` (disabled), `context7` (disabled), `unreal-engine` (disabled, localhost) | not configured | D: already disabled on old |
| M05 | `zoho-crm`, `zoho-projects`, `zoho-salesiq`, `zoho-workdrive` (disabled; URLs carry key material) | not configured | D: already disabled; skills README flags the key-in-URL pattern |

New-only servers: `dropbox` (27 tools), `clickup` (67 tools), both OAuth, both connected.

## 6. Cron jobs (10 old, all paused since 2026-10-01; 4 new)

| ID | Jobs | Status | Notes |
|---|---|---|---|
| R01 | `legacy-sovereign-brief` (daily 07:00 UTC), `so-infra-health` -> `infra-health` (every 15 min, monitor script, silent when healthy), `so-security-scan` -> `security-scan` (daily 06:30) | M | all three enabled on new and delivering: last runs 2026-10-31 07:00, 2026-11-01 05:18, 2026-10-31 06:30, status ok; Git-managed in `cron/jobs.yaml` |
| R02 | `so-daily-digest`, `ge-morning-business`, `ge-weekly-ops`, `hp-morning-readiness`, `hp-weekly-review`, `le-home-status`, `le-weekly-estate` | D (de facto) | paused since 2026-10-01 after `429 tokens_per_day`; nobody recreated them; MIGRATION still lists "Which of the 10 paused cron jobs come back" as open |
| R03 | New-only `config-export` (nightly snapshot to Git) | paused until a Git write credential exists | n/a. It is the replacement for the identity-export workflow (W09) |
| R04 | Decision record plus archive of the 7 job prompts (`/opt/data/cron/jobs.json`, 1.1 MB `cron/output`) | **G** (minor) | **Action:** owner confirms the 7 are dead; copy `jobs.json` prompts into Git or Notion before deletion |

## 7. Gateway platforms and channels

Sessions by `source` in each `state.db` (read-only query). Old: 367 sessions, 71,205 messages. New: 170 sessions, 818 messages (since 2026-10-29).

| ID | Channel | Old (all time / last 14 d, latest) | New | Status |
|---|---|---|---|---|
| H01 | Hermes Desktop | 50 sessions (18 in 14 d), latest 2026-10-28 03:57 | 12 sessions, latest today 05:19 | M |
| H02 | Telegram | 3 (latest 2026-08-21); adapter `disconnected` (disabled 2026-10-29) | `connected`; 1 session 2026-10-30; bot token and allowlist on new | M |
| H03 | Matrix | 0 ever; disabled 2026-10-30 01:52 | `connected`, 0 sessions yet | M |
| H04 | Dashboard / web UI | `hermes-agent.tail717ba0.ts.net` (Service `hermes-dashboard`) | `onemind-hermes.tail717ba0.ts.net` | M |
| H05 | API server (`:8642`) | 3 sessions (all Paperclip era, latest 2026-10-24 21:21) | enabled, `connected`, 1 session 2026-10-30 | M for the generic endpoint; Paperclip path is PC1 |
| H06 | `tak` adapter | `connected`, 0 sessions ever | not ported | D (P03) |
| H07 | `nats` | `connected` on the leaf | `retrying` (waiting for Legacy to let go) | IP (N01) |
| H08 | `cron`, `tui`, `oneshot`, `subagent` | cron 309 (last 2026-10-01 08:01), tui 2 (Aug) | cron 147, tui 1, oneshot 6, subagent 2 | n/a |

## 8. Paperclip (the blocker)

What the old pod does for it: sidecar `paperclip-api-tls` (nginx 1.30.6-alpine, digest-pinned) terminates TLS on `:8643` with a leaf from a private
CA (`hermes-paperclip-ca`, Issuer plus Certificate resources in `platform/hermes/agent/paperclip-api.yaml`) and proxies only
`GET /health`, `GET /v1/capabilities`, `POST /v1/runs`, `GET /v1/runs/{id}[/events]`, `POST /v1/runs/{id}/stop` to the Hermes API server on
`127.0.0.1:8642`. Service `hermes-paperclip-api` (ClusterIP, 8643) selects `app.kubernetes.io/name=hermes, onemind.io/agent=legacy`, so it has no
endpoints once the old pod is gone. Secret `paperclip-agent` (key `credentials.json`) is the callback identity mounted in the old pod.

Who calls it (live, read-only Postgres query in `paperclip`): one agent, **"Hermes"**, adapter `hermes_gateway`, status `idle`, created 2026-10-24, whose
adapter config points at `hermes-paperclip-api`; **1 heartbeat run ever, 2026-10-24 21:21, succeeded**. Paperclip logs mention Hermes nowhere in the last 3000 lines; the
sidecar has access logging off and its log holds only startup connection-refused lines from 2026-10-30. Old `state.db` shows 3 `api_server`
sessions, latest 2026-10-24 21:21. So: **used once, as a test, a week ago**; no recurring traffic seen (the Paperclip deployment has a `HEARTBEAT_SCHEDULER_ENABLED` setting; its value was not read).
Paperclip NetworkPolicy allows 8643 egress only to the old pod's labels; the `paperclip` namespace holds `hermes-gateway-ca` (public CA) for trust.

| ID | Item | Status | Action |
|---|---|---|---|
| PC1 | Transport: Service, TLS proxy, private CA and trust Secret, NetworkPolicy rules both sides | **G** | Owner decides. (a) **Retire the integration**: delete the Paperclip agent first (HERMES.md "Recovery" order), then remove the sidecar, Service, Issuers, Certificates via Git. (b) **Move it**: add an equivalent proxy to `onemind-hermes` (or a small standalone Deployment in front of its `:8642`), add the Paperclip-side NP selector `app.kubernetes.io/name: onemind-hermes`, re-issue and re-seal the CA, repoint the Paperclip agent, run one bounded test. Note `onemind-hermes` currently exposes `:8642` only inside the pod's own network (Service lists only 9119) |
| PC2 | Callback credential `paperclip-agent` | **G** | If (b), seal a new one for `onemind-hermes`; the README says the new pod "has no Paperclip credential" on purpose. If (a), rotate/delete it in Paperclip |

If old Hermes disappears with PC1 unresolved, the Paperclip agent fails its next run with a connection error; nothing else in Paperclip breaks (UI, DB and backups are separate).

## 9. NATS (covered by PR #613, not touched)

| ID | Item | Status |
|---|---|---|
| N01 | PR #613 `hermes-nats-cutover` merged (53073089). Live at 05:28 UTC: `onemind-hermes` connects as a direct client (no leaf), `nats-platform` loaded, **identity not taken**: the adapter logs `identity hermes:zeus:legacy is live elsewhere; not starting (HERMES_NATS_REFUSE_DUPLICATE)` and retries with back-off. Legacy still holds `agents.prompt/status.hermes.zeus.legacy` and heartbeats. Remaining: MIGRATION cutover steps 4 to 6 (stop Legacy's NATS, verify one instance, Agent Console and omos-ui show Legacy live). After retirement, delete the Legacy rule in the `nats-leafnodes` NetworkPolicy and the `hermes-leaf-conf` Secret, and stop sharing `leaf_tak`. `omos-nats-creds` is the shared `onemind` user that **the new pod uses**: do not rotate it before PR #562 gives Hermes its own user | IP |

## 10. Data on the volume

`/opt/data` is 13 GB (`du`): repos 2.6G, state.db 1.4G (WAL 10M), home 1.4G, omos-infra 1.0G, omos-infra-upgrade 964M, `.quarantine` 929M,
OpenClaw 906M, takdocs 637M, lazy-packages 515M, `.cache` 464M, omos-atoc 383M, `.npm` 230M, google_venv 218M, zeus-site 185M, lsp 157M, `.local` 120M,
bin 114M, drone_sim 109M, skills 70M, pipecat 66M, `.venv` 61M, kubectl 60M, profiles 56M, runtimes 55M, onemind-hud 55M, logs 52M, webui 46M, plugins 41M,
tak-data-packages 33M. The new volume is 498 MB.

| ID | Item | Status | Detail |
|---|---|---|---|
| V01 | **Git checkouts with work not on a remote** | **G** | `repos/omos-legacy` 2 unpushed commits (repo archived 2026-10-16); `repos/omos-infra` 1 dirty + 1 unpushed; `repos/OMOS` 1 dirty; `repos/omos-voice` 2 dirty; `omos-blueprints` 2 dirty; `OpenClaw` 1 dirty; `ghl-mcp-pro` 1 dirty (third party fork); one repo has credentials embedded in its remote URL (rotate and strip). Everything else is clean and pushed (`omos-atoc` 11,813 commits, `OMOS` 12,314 commits, `omos-infra` 586). **Action:** push or `git bundle` the unpushed ones |
| V02 | **Unversioned working directories** (no `.git`) | **G** | `onemind-hud` 55M, `drone_sim` 109M, `webui` 46M, `pets` 3M, `sites` 624K, `deployments` 204K, `launchops` 32K, `redpanda-connect` 48K, `omos/` (Aug 4), `telegram-business` 32K, `apple-notes-export` 424K (personal notes export, not reproducible), `tak-data-packages` 33M, `tak-certs-rebuild` 328K, `skills/` custom skills not shipped (section 3), `profiles/` 56M (4 soul-only profiles, no sessions). **Action:** owner triages; archive (tar to Dropbox/Notion) anything wanted |
| V03 | **Research notes, reports, loose files** (not reproducible) | **G** | Top level: `3d-tiles-production-research.md` 9K, `COPILOTKIT.md` 8K, `MAC-FILE-AUDIT.md` 11K, `NATS-for-TAK-Transport-Architecture.md` 12K, `ORGANIZATION-PLAN.md` 10K, `OneMind-CloudTAK-Roadmap.md` 7K, `OneMind_Android_APK_Roadmap.md` 13K, `OneMind_TAK_COP_Architecture.md` 8K, `S24_ADB_ELIMINATE_MAC_RESEARCH.md` 12K, `S24_Ultra_Onemind_Status.md` 14K, `TAK-AI-Business-Plan-2026.md` 16K, `TAILSCALE_*` (3 files, 5 to 9K), `ar-glasses-tactical-comparison-2025.md` 7K, `atak-s25-ultra-setup.md` 22K, `cesium-ecosystem-report.md` 18K, `iphone-offload-checklist.md` 21K, `iphone-to-s25-ultra-transition.md` 10K, `opentakserver-research.md` 11K, `tak-server-research.md` 15K, `tak_3d_alternatives_research.md` 12K, `tak_uas_research.md` 13K, `s24-research-analysis.md`, `s24_presence_analysis.md`, `hp_tasks.json`, `OneMind_OS_OpenRouter.pdf`, `onemind_openrouter_onepager.html`. Dirs `reports/` (4.1M: existing-cluster-work-platforms, hermes-first-operations-research, hermes-mobile-client-comparison, onemind-cluster-audit-2026-10-27, onemind-long-term-platform-recommendation, onemind-platform-decision, feed-platform-audit-2026-10-28 + zip), `research/` 120K, `research-s24`, `research_s24`, `s24_research`, `notes/` 8K, `takdocs` 637M (check for a git origin; reproducible if a doc mirror), `cron/output` 1.1M. `fix_*.py`, `test_*.py`, APKs and ten `config.yaml.bak-*` are clutter (D). **Action:** copy the notes and `reports/` somewhere durable (Git `docs/` or Notion/Dropbox); about 1 MB total |
| V04 | **Credential material to rotate** (names only) | **G** | Files: `.env`, `.env.bak-20260811T160842Z`, `.env.bak-20260827T102642Z`, `.env.bak-20260905T104630Z`, `.env.pre-chat-switch-20261029`, `.ha_token`, `.room-link-grant-secret`, `auth.json`, `tokens.txt`, `tak_admin.key`, `tak_admin.p12`, `tak_admin.pem`, `tak_ca.pem` (public CA), `google_client_secret.json`, `google_token.json`, `hermes.cmd`, `mcp-tokens/` (clickup, cloudflare, dropbox, miro, raindrop, webflow: `*.json`, `*.client.json`, `*.meta.json`), `.ssh/` (known_hosts only), `config.yaml` and every `config.yaml.bak-*`, `.before-onemind-mcp`, `.pre-chat-switch-20261029`, `.pre-matrix-dev-202610120843` (Matrix token, Composio consumer key, GitHub bearer, dashboard hash/secret, 1Password env map), `profiles/*/config.yaml` and `auth.json`, one `.git/config` with a URL credential, skills named in the skills README "Secrets found in candidates" list, `home/`, `tmp/`, `scripts/` (untriaged). **Old `memories/MEMORY.md` and its two backups contain a plaintext WordPress application password for `wp.zeusdelacruz.com`.** In-cluster SealedSecrets that carry old credentials: `hermes-env` (`DIGITALOCEAN_AI_KEY`, `AWS_*`, `API_SERVER_KEY`, 1Password and Todoist tokens), `provider-config-backup`, `provider-config-pre-paperclip`, `provider-config-pre-litellm-v1`, `paperclip-agent`, `cloudtak-agent-token`, `hermes-leaf-conf`. Rotate after retirement: TAK admin cert (MIGRATION), Google client and token, WordPress app password, Composio key, GitHub token, Matrix and Telegram bot tokens (they were copied to the new instance; rotating them needs the new SealedSecret resealed at the same time), Cloudflare/Miro/Raindrop/Clickup/Dropbox/Webflow OAuth grants, 1Password service account, DigitalOcean AI key, AWS keys, Todoist, dashboard auth. **Do not rotate** `omos-nats-creds`, `hermes-dashboard-auth` or `hermes-nats-ca` without resealing: `onemind-hermes` mounts all three |
| V05 | Reproducible clutter: `repos` (clean, pushed), `omos-infra`, `omos-infra-upgrade`, `OpenClaw`, `lazy-packages`, `.cache`, `.npm`, `google_venv`, `.venv*`, `lsp`, `bin`, `kubectl`, `runtimes`, `logs`, `.local`, ten `config.yaml.bak-*` | D | MIGRATION "Repository checkouts and caches are reproducible" |
| V06 | **Session history** (`state.db` 1.4 GB, 367 sessions, 71,205 messages): desktop 50 sessions / 46,273 messages, cron 309 / 13,320, telegram 3 / 1,286, tui 2 / 10,298, api_server 3 / 28. About 1.1 GB of the file is the FTS and trigram search indexes (`dbstat`: messages table 240 MB). Also `.quarantine/20261016/state_compact.db` (929 MB, an October compaction copy). New instance holds none of it (starts 2026-10-29) | **G** (decision) | **Action:** owner decides whether the 58 human/API sessions are worth keeping (desktop chats are the only real history). Cheap option: export those 58 sessions to JSONL (about tens of MB) and archive with the VolSync snapshot; discard cron sessions. The VolSync restic repo already holds the full file |
| V07 | **Backup of `hermes-data`** | **G** (verify) | `ReplicationSource volume-c96e219d2314` (sourcePVC `hermes-data`, Snapshot copy, restic to S3, daily 01:00 UTC, retain 7 daily/4 weekly/3 monthly/1 yearly): **last success 2026-11-01T01:02:29Z, 2m30s**, next 2026-11-02T01:00. No restore drill was checked in this audit. The PV `pvc-993e5ef7...` is `do-block-storage` with reclaim policy **Delete**. **Action:** trigger or wait for one final sync after the last writes, run a restore into a scratch claim, then proceed. The ReplicationSource lives in `platform/backups/volsync/sources.yaml` and must be edited when the PVC goes (keep the repo credentials) |

`hermes-fresh-data` (25 GiB, `do-block-storage`, created 2026-10-26) is the predecessor PVC. No pod mounts it. Same cleanup, outside this retirement.

## 11. Other dependents of `hermes-legacy`

| ID | Dependent | Status | Detail / action |
|---|---|---|---|
| W01 | `nats-leaf` sidecar, `hermes-leaf-conf`, `leaf_tak` | D | replaced by direct client (#613). Delete the rule and Secret after cutover |
| W02 | `patch_nats_platform.py` hot-patch (`nats-platform-patch-cm`) | M | now a normal line in the vendored adapter (MIGRATION) |
| W03 | Init containers writing the LiteLLM provider and the `onemind` MCP entry (`hermes-litellm-provider`, `hermes-onemind-mcp-config` ConfigMaps) | D | LiteLLM no longer serves MCP; gateway replaces it |
| W04 | `hermes-observer` ServiceAccount, ClusterRole/Binding `hermes-cluster-observer` (read-only) | D | README: "no Kubernetes service-account token... or observer credential"; cluster reads go through the gateway `kubernetes` target. |
| W05 | `cloudtak-agent-token` SealedSecret (Hermes to CloudTAK API) | **G** (minor) | not mounted on new. Decide whether Legacy still needs the CloudTAK API; if so, reseal into `onemind-hermes-runtime` |
| W06 | Secrets owned by `platform-legacy` (`platform/hermes/bridge`) **that `onemind-hermes` mounts**: `hermes-dashboard-auth` (envFrom), `omos-nats-creds` (`NATS_USER`/`NATS_PASSWORD`), `hermes-nats-ca` (CA volume). Also `ghcr` | **G (blocking)** | move to a path owned by `platform-onemind-hermes` (or a new app) **before** deleting `platform-legacy`, or they are deleted and the new pod loses its dashboard auth and the NATS connection |
| W07 | `onemind-mcp` Deployment, Service, PDB, NetworkPolicy, ConfigMaps (`onemind-mcp-apps`, `onemind-mcp-ca`), SealedSecrets (`onemind-mcp-auth`, `onemind-mcp-registry`) owned by `platform-legacy` source `platform/hermes/mcp` | **G (blocking)** | this is the gateway's `onemind` target: the new `ops` group depends on it. Move to its own Application first. Its NetworkPolicy ingress rule names the old pod (`app.kubernetes.io/name=hermes, onemind.io/agent=legacy`), plus the gateway rule; drop the Legacy rule afterwards |
| W08 | Tailscale egress Services `omos-mbp` (22), `omos-s24` (8022), `omos-ha-edge` (8123) in `hermes`, owned by `platform-legacy` (`egress.yaml`) | **G** (minor) | `onemind-hermes` NetworkPolicy blocks them; MEMORY says restoring is one PR. Decide keep (move to the new app) or drop |
| W09 | `.github/workflows/hermes-identity-export.yml` and `scripts/hermes_identity_export.py` exec into the pod labelled `hermes/legacy`; runner `github-runner-hermes` (label `onemind-hermes-readonly`) is 100 restarts in 13 days; `scripts/onemind_mcp_acceptance.py`, `scripts/onemind_mcp_probe.py` and `scripts/tests/test_gitops_paperclip_hermes_transport.py` name `hermes-legacy` | **G** (minor) | point to `onemind-hermes` or delete; the paused `config-export` cron is the replacement |
| W10 | References to clean: NetworkPolicies `nats-leafnodes`, `litellm-clients` (port 4000), `onemind-mcp`, `paperclip-network`; Backstage and portal catalog entries `hermes-legacy` and `k8s-hermes-hermes-legacy-*` (apps.yaml, cluster.yaml, plugins.yaml in `platform/backstage/catalog/` and `services/onemind-portal/catalog/`); `docs/STACK.md`, `docs/catalog/system-map.md`, `docs/infra/CLUSTER-SERVICE-INVENTORY.md`, `platform/argocd/app-of-apps.yaml`; Tailscale proxy `ts-hermes-dashboard-6thbx` + Ingress `hermes-dashboard` (`hermes-agent.tail717ba0.ts.net`); `agent-console` text only (it discovers by NATS identity, not by pod) | **G** (minor) | one cleanup PR. Nothing else in the cluster calls the old Services: omos-ui/Agent Console, LiteLLM and Backstage reference it only in policy and catalog text |

## 12. Argo CD

| ID | Item | Status |
|---|---|---|
| Z01 | Application **`platform-legacy`** (sync wave 3). Three sources on `main`: `platform/hermes/mcp`, `platform/hermes/agent` (Kustomize), `platform/hermes/bridge` (plain directory, `recurse: false`, excludes `*.secret*`). `syncPolicy.automated`: **`prune: true`, `selfHeal: true`**, `CreateNamespace=true`. **`metadata.finalizers: [resources-finalizer.argocd.argoproj.io]`** (cascading delete). Live: Synced / Healthy, 39 resources: Namespace `hermes`, Deployments `hermes-legacy` and `onemind-mcp`, PVC `hermes-data`, 6 Services (`hermes-dashboard`, `hermes-paperclip-api`, `onemind-mcp`, and the three egress Services), Ingress `hermes-dashboard`, 2 NetworkPolicies, PDB, ServiceAccount, ClusterRole and Binding, 2 Issuers and 2 Certificates, 6 ConfigMaps, 12 SealedSecrets. **Deleting the Application (or a prune after removing a source) cascades into every one of those**, including the `hermes` Namespace (so `onemind-hermes`, `onemind-hermes-data` and the `volume-*` objects) and `hermes-data` (PV reclaim `Delete`). `ignoreDifferences` covers the PVC storageClass. | **G (blocking)** |
| Z02 | Application `platform-onemind-hermes`: source `platform/hermes/onemind-hermes`, automated `prune: false`, `selfHeal: true`, **no finalizer**; currently OutOfSync on ConfigMap hashes (rollout in flight). It does not own the Namespace, `omos-nats-creds`, `hermes-nats-ca` or `hermes-dashboard-auth`. Related: `platform-paperclip`, `platform-agent-fabric` (Agent Console; no pod dependency) | info |

---

## Recommended retirement order

1. **Finish the NATS cutover (#613 steps 4 to 6).** Scale/disable Legacy's NATS so one instance holds the identity; verify heartbeat, Agent Console and a prompt. Legacy is still serving Paperclip and the TAK beacon at this point.
2. **Decide Paperclip** (PC1/PC2): retire the Paperclip "Hermes" agent, or move the TLS proxy to `onemind-hermes` and prove one run. It was used once, a week ago.
3. **Owner decisions in one pass:** the 13 NEEDS-OWNER skills (S06), composio/github/cloudflare/raindrop/miro MCP (M02, M03), the 7 dead cron jobs (R02, R04), CloudTAK token (W05), the two third-party plugins (P04, P05), Todoist and `onemind-product-strategy` usage, session history (V06).
4. **Salvage:** push or bundle unpushed repos (V01); copy notes, `reports/`, `apple-notes-export`, custom skills, anything wanted from the unversioned dirs (V02, V03); export the 58 human/API sessions if wanted (V06); copy the 7 cron prompts.
5. **Safety net, with owner approval:** one final VolSync run after the last writes and a restore into a scratch claim (V07); then patch PV `pvc-993e5ef7...` to `Retain` (CLAUDE.md item 6 allows this with approval) so no Git change can destroy it.
6. **Re-home what the old app owns (one Git PR, merge first):** move `hermes-dashboard-auth`, `omos-nats-creds`, `hermes-nats-ca` into `platform/hermes/onemind-hermes/` (or a shared app); move `platform/hermes/mcp` (onemind-mcp) to its own Application or into `platform-onemind-hermes`; decide the Namespace owner and the egress Services (W06 to W08); move or archive the Git-only `platform/hermes/agent/` content (I05). Use `Prune=false` annotations or remove the finalizer first so moving a resource between Applications never deletes it.
7. **Retire the workload (second PR):** remove the Deployment, Services, Ingress, NetworkPolicy, Issuers/Certificates, observer RBAC, `paperclip-agent`, `provider-config-*`, `hermes-env`, leaf Secrets from `platform/hermes/agent` and `bridge`. Keep the **PVC out of this PR** (or keep it unpruned). Do not `kubectl scale` as a shortcut: `selfHeal` reverts it; use Git. Watch Telegram, Matrix, Desktop, cron, NATS, Paperclip for a cooling period (suggest 7 to 14 days).
8. **Rotate** the credentials in V04 (except the three shared Secrets, which are resealed together with their new owner; `omos-nats-creds` waits for #562). Clean the references (W09, W10).
9. **Last, with explicit owner approval (CLAUDE.md: never delete PVCs without a fresh snapshot):** delete PVC `hermes-data` and `hermes-fresh-data`, remove the `volume-c96e219d2314` ReplicationSource from `platform/backups/volsync/sources.yaml`, keep the restic repository, update docs (STACK, system-map) and the Backstage catalog, then delete or archive `platform-legacy` once it only owns what you intend.

Not verified in this audit: that a VolSync restore actually works; whether `hermes-project-stewardship` or `onemind-vision` ever ran tools; the upstream defaults for memory limits; Matrix inbound traffic on the new instance (0 sessions so far).

## Salvage — 2026-11-01

Owner approved "Push to Git, archive the rest". Read-only on the volume (`kubectl exec` and tar out); every push gated by `gitleaks` (working tree and new commits), no findings pushed.

**Pushed to Git.** Branch `salvage/2026-11-01-hermes-legacy` on `onemindos/omos` (public; commit `707fc7c4`, based on `5ba89dbb`, one commit, not reviewed): `repos/omos-infra` dirty `platform/matrix/base/synapse/manifests/configmap.yaml` (4+/6-) and `repos/redpanda-console-deploy` untracked `platform/argocd/platform-redpanda-console.yaml` and `platform/observe/redpanda-console/{ingress,service}.yaml`. `repos/omos-infra`'s unpushed commit `7b4b59e7` (feed-mta) is already on `origin/feat/feed-mta`.

**Archived** in the private repo `onemindos/hermes-legacy-archive` (670 files, 62 MB, README lists contents). Remotes that no longer exist (404): `omos-legacy` (2 unpushed commits as `.mbox`), `omos-voice` (untracked `deploy/livekit`), `omos-blueprints` (`index.html` diff). `repos/OMOS` (stash plus one untracked file) is kept there because its history is not in `onemindos/omos`. `OpenClaw` overlay file likewise. Also: top-level notes and research, `reports/`, `research*`, `notes`, `cron/output`, `apple-notes-export`, `tak-data-packages`, `onemind-hud` (no APK), `drone_sim` (no venv), `webui`, `pets`, `sites`, `deployments`, `launchops`, `redpanda-connect`, `omos`, `telegram-business` (no db), the two present NEEDS-OWNER skills.

**Removed by gitleaks (names only, 44 files).** 37 under `webui/sessions/` (chat and run journals), `deployments/voice/{save_ar.py,save_ar_to_notion.py,mumble-bridge.yaml}`, `deployments/takml-server/k8s-manifests.yaml`, `tak-data-packages/overlays/live/faa_tfrs.kml`, `reports/feed-platform-audit-2026-10-28/source-catalog-240-records.csv` and that report's `.zip`, `repos/omos-voice/deploy/livekit/README.md`. Those stay only in the VolSync snapshot.

**Not copied.** Credential material per V04, `webui` key and session files, `tak-certs-rebuild`, `profiles/`, `state.db` and other databases. Large or reproducible: `takdocs` (637 MB, public TAK docs), `onemind-hud/OneMindHUD-0.1.0-debug.apk` (57 MB), `drone_sim/.venv` (111 MB), `goodlock.apk` (5.6 MB), `OneMind-Zeus-ATAK.dp`/`.zip` (6 KB, may hold a cert), `OneMind-DimOS-Bridge.tar.gz` (8 KB), `nyc_dot_cameras.zip` (66 KB). Not worth saving: `ghl-mcp-pro` (`package-lock.json`), `plugins/hermes-project-stewardship` (`.hermes-catalog.json`).

**Other findings.** `repos/omos-nats` has a credential embedded in its remote URL (rotate; this is V04's repo). Clean and pushed: `omos-atoc`, `omos-ui`, `atak-plugins`, `hermes-plugin-tak`, `omos-infra-upgrade`, `zeus-site`, `pipecat`, the rest. `repos/omos-rads` is 1 behind; the top-level `omos-infra` is 1204 behind.

## Final snapshot and restore drill (2026-11-01)

**Provider snapshot.** `hermes-data-final-2026-11-01`, id `8f5b4644-d5b6-11f1-889c-fa2c53cb4798`, taken 2026-11-01T05:39:09Z from DO volume `61aa7782-78b9-11f1-9870-0a58ac1463ec` (PV `pvc-993e5ef7-9942-460b-801c-615eab02ae60`, PVC `hermes/hermes-data`), tag `hermes-legacy-retirement`. Size 15.61 GiB (min disk 25 GiB). Crash-consistent; the pod kept running. The PV reclaim policy is `Retain` and the PV is still `Bound`.

**VolSync.** `hermes/volume-c96e219d2314`: last sync 2026-11-01T01:02:29Z (2m30s), idle (`WaitingForSchedule`), next 2026-11-02T01:00Z. The CR was not touched. The 01:02Z run is the latest restic copy at the time of this drill, not a final copy: any change after that time is only in the provider snapshot.

**Restore drill (passed).** Followed `docs/operations/hermes-restore-drill-2026-10-16.md`. It ran in namespace `hermes`, not a throwaway namespace, because the restic repository Secret is namespace-bound and was not copied. All drill resources were named `hermes-restic-drill-20261101*` and labelled `onemind.io/restore-test=drill-20261101`.
- Restore: `ReplicationDestination` (Direct, `do-block-storage` 25 GiB) finished at 2026-11-01T05:46:46Z in 7m24s.
- Verify Job: deny-all NetworkPolicy, no gateway started.
- `state.db` on the restored copy: 1,495,732,224 bytes, 42 tables, `PRAGMA integrity_check` = ok, sessions 367, messages 71,205, session_model_usage 2,069. Live (read-only query): sessions 367, messages 71,205.
- Size and count: restored 13,336,420 KB and 332,303 files, against live 13,336,532 KB and 332,302 files. Top-level directory sizes match to within about 0.1% (repos 2.54 GiB, state.db 1.39 GiB, home 1.31 GiB, omos-infra 0.98 GiB, `.quarantine` 0.91 GiB, OpenClaw 0.88 GiB). The small differences are live writes after the 01:02Z backup.
- Cleanup: Jobs, NetworkPolicy, `ReplicationDestination` and PVC deleted. Drill PV `pvc-9c68f298-c121-4575-988c-5553f530d1b9` and DO volume `97bd3c41-d5b6-11f1-a034-dab28107f2f3` are gone (404). No drill PVC, mover cache PVC or Job remains.

**Not verified.** The restored app config and secrets were not inspected (by design). The first verify Job hit its 900s deadline during `integrity_check`; it was rerun with a longer deadline and passed. `state_compact.db` was not examined.

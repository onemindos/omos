# hermes-legacy → onemind-hermes migration: inventory and order

Read-only inventory of `hermes-legacy` taken 2026-10-29 from the live pod
(`/opt/data`, `state.db`, `config.yaml`, the `nats-leaf` monitor). It records
what Legacy actually uses, so the migration copies that and nothing else. See
`README.md` for how onemind-hermes is deployed. In this document "Legacy" means
the `hermes-legacy` Deployment; onemind-hermes keeps the Legacy persona as a
single soul.

## Summary

| Area | Legacy today | Actually in use | Migrate? |
|---|---|---|---|
| Chat traffic | 367 sessions, 71k messages | Hermes desktop app (28 sessions in last 30 days, latest 2026-10-28); Paperclip API (3, latest 2026-10-24) | Yes: desktop + API first |
| Telegram | enabled | Last chat session 2026-08-21; only used as cron delivery target | With cron, if cron returns |
| Cron | 10 jobs | **All 10 paused since 2026-10-01** (last runs hit `429 tokens_per_day`) | Decide per job; none is live |
| Profiles | 14 dirs | 4 have content (`generational-empire`, `holistic-performance`, `legacy-evolution`, `sovereign-operations`): SOUL.md only, 0 sessions, no skills or cron. 10 are empty | No: single Legacy soul; profile souls stay in Git |
| Skills | 202 (46 bundled, 38 archived) | ~104 not shipped with onemind-hermes (list below) | OneMind-authored set only |
| Plugins | 6 enabled | `nats-platform`, `nats-context`, `platforms/tak-platform`, `onemind-maven` (OneMind); `memory-wiki`, `hermes-project-stewardship` (third party) | OneMind ones with NATS step |
| MCP servers | 13 configured | Enabled: `onemind` (failing 401), `composio`, `github`, `cloudflare`, `raindrop`, `miro` (no OAuth tokens) | `onemind`, `composio`, `github`; others on request |
| Gateway platforms | `nats`, `tak`, `telegram`, `matrix` enabled | NATS: leaf moved ~684k msgs in, ~5.6k out over 2 days (feed consumer, heartbeats). No sessions recorded from `nats`, `tak` or `matrix` | NATS last (cutover) |
| Memory | `MEMORY.md` 2.2 KB, `USER.md` 1.4 KB | Yes | Copy after review |
| Model | `openai-codex` / `gpt-6.1-sol`; providers `digitalocean`, legacy `bedrock` | Yes | onemind-hermes already selects the same; sign in |
| Data volume | 13 GB | Mostly clutter (below) | No bulk copy |

## Disk breakdown (13 GB)

`repos` 2.6G, `state.db` 1.4G, `home` 1.4G, `omos-infra` 1.0G,
`omos-infra-upgrade` 964M, `.quarantine` 929M, `OpenClaw` 906M, `takdocs` 637M,
`lazy-packages` 515M, `.cache` 464M, `omos-atoc` 383M, `google_venv` 218M, plus
about 100 loose files (research notes, `fix_*.py` scripts, APKs, PDFs, ten
`config.yaml.bak-*`).

Repository checkouts and caches are reproducible. Research notes and reports
that matter belong in Git or Notion, not on the agent volume.

**Credential material on the volume, outside SealedSecrets** (names only):
`tak_admin.key`, `tak_admin.p12`, `tak_admin.pem`, `tak_ca.pem`,
`google_client_secret.json`, `google_token.json`, `tokens.txt`, `auth.json`,
`mcp-tokens/`. Do not copy these to onemind-hermes. Rotate the TAK admin cert and the
Google client after cutover; recreate needed ones as SealedSecrets.

## Skills not shipped with onemind-hermes

onemind-hermes (Hermes v0.22.5) bundles 58 skills. Legacy has about 104 non-archived skills that
onemind-hermes lacks. They fall into three groups:

**OneMind / Zeus authored: candidates to bring over (review each).**
Agent roles: `so-grid`, `so-guardian`, `so-oracle`, `so-trinity`, `so-forge`,
`ge-empire`, `hp-spartan`, `le-eden`, `le-haven`, `le-heritage`.
Platform: `doks-cluster-ops`, `onemind-stack`, `onemind-fabric-deployment`,
`onemind-inter-agent-protocol`, `nats-operations`, `ha-nats-bridge`,
`tak-server`, `tak-mission-ops`, `tak-mission-runner`, `mission-replay`,
`omos-feed-builder`, `omos-ui-build-deploy`, `matrix-doks-ops`,
`hermes-mcp-config`, `legacy-integrations`, `legacy-model-config`,
`hermes-deployment-topology`, `hermes-plugin-pattern`,
`cloudflare-access-patterns`, `tailscale-pod-access`,
`tailscale-continuous-mesh-access`, `composio`, `composio-*` (5).
Personal/ops: `mac-*` (4), `s24-*` (2), `android-adb-management`,
`windows-*` (2), `find-zeus`, `zeus-site`, `wordpress-*` (2), `notion-*` (2),
`onepassword-legacy`, `credit-bureau-dispute`, `vault-dedup-audit`,
`external-drive-cloud-sync`, `session-librarian`,
`modesty-ai-daily-briefing`, `grounddata-ai-daily-news`.

**Upstream optional / hub skills: reinstall from the Hermes hub only if wanted.**
`github/*` (6), `mlops/*` (9), `creative/*` (6), `research/{polymarket,blogwatcher,research-paper-writing}`,
`productivity/{nano-pdf,ocr-and-documents,petdex}`, `smart-home/openhue`,
`data-science/jupyter-live-kernel`, `media/heartmula`,
`software-development/plan`, `yuanbao`, `hermes-desktop-plugins`.

**Likely stale: leave behind.** Skills describing retired or speculative work:
`omos-cortex-neurolinked`, `omos-legacy-genome`, `legacy-self-awareness`,
`openclaw-k8s-deploy`, `repo-absorption`, `resume-session-on-pod`,
`onemind-xr-platform`, `onemind-product-strategy`, `merge-reconciler`,
`hermes-mobile-desktop-sync`, and the 38 under `.archive/`.

Skills that change often should live in Git (`platform/hermes/agent/`) and be
seeded, not hand-copied from the volume.

## Legacy-only wiring (in `platform/hermes/agent/deployment.yaml`)

- `nats-leaf` sidecar with its own bus identity (`hermes-leaf-conf`).
- `patch_nats_platform.py` applied to the gateway at startup (code hot-patch).
- Init containers that write the LiteLLM provider and the `onemind` MCP entry
  into `config.yaml`.
- Paperclip HTTPS proxy sidecar (`paperclip-api-tls`) and `paperclip-agent`
  secret.
- Read-only cluster-observer service-account token.
- Main `hermes` container image is a tag without a digest.

## Migration order

Each step ends with an observed check on onemind-hermes. Legacy keeps running until step 7.

1. **Model and chat.** Sign in to the model on onemind-hermes; get a real reply from the
   dashboard and from the Hermes desktop app pointed at onemind-hermes. This replaces
   Legacy's main traffic.
2. **Identity.** Seed the Legacy `SOUL.md` and doctrine from Git
   (`platform/hermes/agent/agents/`, `doctrine/`). onemind-hermes runs a single
   soul, so the four profile souls are not seeded as separate agents. Copy
   `MEMORY.md` and `USER.md` after reading them.
3. **OneMind MCP.** Give onemind-hermes its own `onemind-mcp` token (SealedSecret), allow
   the egress in `network.yaml`, and confirm tools list. Fix Legacy's 401 only if
   cutover is far away.
4. **Skills.** Import the reviewed OneMind set, preferably by moving it into Git
   and seeding it. Add `composio` and `github` MCP.
5. **Paperclip API.** Move the TLS proxy and agent credential; confirm a
   Paperclip call against onemind-hermes.
6. **Cron and Telegram.** Recreate only the jobs the owner wants, with a token
   budget that avoids the October `tokens_per_day` failures; enable Telegram as
   their delivery target.
7. **NATS cutover.** Prepared in Git (this section's PR): `nats-platform` and `nats-context`
   vendored into `platform/hermes/agent/plugins/` without the startup hot-patch, onemind-hermes a
   **direct client** of the cluster bus (no leaf sidecar), identity `agents.*.hermes.zeus.legacy` guarded
   against a double presence. The order of the switch, the checks and the rollback are in
   [NATS inventory and cutover](#nats-inventory-and-cutover-step-7) below. `tak-platform` is not ported.
8. **Retire Legacy.** Keep a VolSync snapshot of `hermes-data`, then, with
   owner approval, remove the Legacy Deployment and Application. Rotate the
   credential files listed above.

## NATS inventory and cutover (step 7)

### What `hermes-legacy` does on NATS (measured 2026-11-01 04:48 UTC, uptime 2d2h)

Read from the sidecar's monitor (`/connz?subs=detail`, `/leafz`, `/varz`), the plugin sources on the pod and the
Synadia SDK in `/gateway-deps`. Nothing was changed on Legacy.

The leaf (`hermes-legacy-leaf`) has 4 client connections, all from the one `hermes` container, plus the uplink
to `nats-2` as `leaf_tak` into the cluster account (local account `OMOS`, hub side `ONEMIND`). Over the 2d2h the
uplink carried 652,395 messages from the cluster and 18,324 to it; the four clients together published those
18.3k (heartbeats, beacons, stats pings, below) and were delivered about 1.32M. So Legacy is almost all
listener: the earlier "684k in, 5.6k out" is the same picture seen from the hub.

| Connection | Plugin | Subscribes | Publishes | Notes |
|---|---|---|---|---|
| cid 9, the agent itself | `nats-platform` (Synadia `hermes-nats-gateway` 0.1.2; SDK `synadia-ai-agents` 0.8.1, `-agent-service` 0.5.0) | `agents.prompt.hermes.zeus.legacy` and `agents.status.hermes.zeus.legacy`, both queue group `agents`; micro-service `$SRV.PING\|INFO\|STATS` (with `.agents` and `.agents.<id>`) | heartbeat on `agents.hb.hermes.zeus.legacy` every 30 s (about 6,100 sent), replies on the caller's inbox | Identity is `{agent}.{owner}.{session}` = `hermes`, `zeus`, `legacy` (`HERMES_NATS_SESSION_NAME=legacy`, `platforms.nats.heartbeat_interval_s: 30`). The heartbeat body is the protocol 0.3 section 8.3 payload: `agent`, `owner`, `session`, `instance_id`, `ts`, `interval_s`, `protocol_version`, `endpoints`. That, and `$SRV.PING.agents`, is what discovery (the Agent Console in `agent-fabric`, omos-ui) reads. No JetStream consumer, no KV |
| cid 7 and 8, the live picture | `nats-context` (OneMind, source never in Git) | core subscriptions, no queue group: `tak.cot.>` (622,475 msgs), `tel.>` (29,312), `intel.>` (576), `agents.hb.>` (6,109, every agent's heartbeat), plus `ent.>`, `det.>`, `evt.>`, `evt.threat.>`, `evt.mission.>`, `evt.task.>`, `cmd.alert.>`, `svc.bridge.>`, `sys.>` (0 so far) | `$SYS.REQ.SERVER.PING` once a minute per connection; on demand the nine tools (`nats_publish`, `nats_request`, `nats_peek`, `nats_jetstream_fetch`, `nats_kv_get`, `nats_kv_put`, `nats_kv_list`, `nats_server_info`, `nats_list_streams`) open a short connection each | This is the "feed consumer": 661k of the 1.34M outbound messages, counted twice because the gateway and the dashboard process each load the plugin. Core subscriptions only. It prepends a `[LIVE OPERATIONAL PICTURE]` block to every model call |
| cid 10 `hermes-tak-adapter` | `platforms/tak-platform` | `tak.cot.out` (0 messages ever) | presence beacon to `tak.cot.inject` every 30 s (about 6,100), and the same beacon plus GeoChat replies straight to `takserver-messaging:8089` over TLS with the TAK admin certificate and `CERT_NONE` | Nothing publishes `tak.cot.out`, so it has never heard a GeoChat (PR #585 section 5.1 explains why). No session was ever recorded from `tak` |

The hot-patch (`patch_nats_platform.py`, `nats-platform-patch-cm.yaml`) changes one thing in `adapter.py`: it moves
`set_current_session_key` and `reset_current_session_key` from `tools.approval` to `tools.approval_context`.
Hermes v0.22.5 has them **only** in `tools.approval_context` and has **no native NATS platform**
(`grep` over `gateway/`, `hermes_cli/`, `plugins/`, `tools/`, `agent/` finds none), so the plugin is still needed
and the fix is still needed. It is now a normal line in the vendored `adapter.py`, not a patch applied at start.

### What was ported, what was not

| Item | Decision |
|---|---|
| `nats-platform` | Vendored to `platform/hermes/agent/plugins/nats-platform/` (runtime files, Apache-2.0 `LICENSE`). Three small changes, listed in its `ONEMIND-CHANGES.md`: the approval import, credentials and TLS from the environment, and refusing a duplicate identity. Loads through the real v0.22.5 plugin discovery in the live onemind-hermes image (checked read-only in `/tmp`, nothing kept) |
| `nats-context` | Vendored to `platform/hermes/agent/plugins/nats-context/`, with `conn.py` so it can use the same credentials. **Not enabled by default** (`apply.sh --with-context`): it takes the raw feeds into the model prompt and gives it a generic `nats_publish`, which PR #562 (no raw feeds, no `cmd.*` for Hermes) and PR #564 (`nats_publish` retired behind a flag, `onemind-bus` tools instead) remove on purpose. Owner decision below |
| `platforms/tak-platform` | **Not ported.** It hears nothing today, it needs `tak_admin.pem` on the volume (which this document says not to copy) and it disables certificate verification. PR #585 replaces it with a TAK-native client using its own certificate. It stays on Legacy until then, so Legacy keeps its beacon and nothing is lost that worked |
| `onemind-maven` | Is `onemind-vision` (renamed), already in Git at `platform/hermes/agent/plugins/onemind-vision/`. It makes no NATS connection (the only "NATS" in it is a string in a prompt). Not part of this step |
| `memory-wiki`, `hermes-project-stewardship` | Third party, no NATS. Not part of this step |

### How onemind-hermes connects

Same shape as every other cloud workload (`CLAUDE.md` Bus Authority): in-cluster `tls://nats.nats.svc.cluster.local:4222`,
user and password from a Secret, certificate verified against the cluster CA, **no sidecar**.

| Piece | Where |
|---|---|
| Credentials | `NATS_USER` and `NATS_PASSWORD` from the existing Secret `omos-nats-creds` in `hermes` (`platform/hermes/bridge/omos-nats-creds.sealedsecret.yaml`; keys `user`, `password`; the shared `onemind` client user, byte-identical to the copy in `omos`). **No new secret and no server change was needed.** PR #562 will give Hermes its own user in account `AGENTS`: reseal these two values from it, nothing else changes |
| CA | Secret `hermes-nats-ca` (`ca.crt`, same certificate as `nats-server-tls`), mounted at `/etc/nats-ca`; `NATS_TLS_SERVER_NAME=nats.nats.svc.cluster.local` |
| Dependencies | init container `install-gateway-deps`: `nats-py`, `nkeys`, the two Synadia packages, pinned with hashes in `nats/gateway-deps.txt`, `--no-deps` so the venv's pydantic and pynacl are not replaced, installed to `/opt/data/.local/gateway-deps` and added to `sys.path` by the plugin itself (the `hermes` launcher unsets `PYTHONPATH` and the venv has no user site, so a variable or a `.pth` does not work). PyPI is needed only when the pin changes. A failed install does not stop the pod |
| Identity | `HERMES_NATS_AGENT=hermes`, `HERMES_NATS_OWNER=zeus`, `HERMES_NATS_SESSION_NAME=legacy`: the subjects above, unchanged |
| Authorisation | `NATS_ALLOW_ALL_USERS=true` for the NATS platform only. The adapter sets every caller's user id to the session name, so a per-user allowlist cannot work; the bus authenticates. `GATEWAY_ALLOW_ALL_USERS` stays `false` |
| Network | egress to the NATS pods on 4222 in `network.yaml`; ingress for this pod in the `nats-clients` policy (`platform/nats/network-policy.yaml`) |
| Switch | the plugin is only loaded when listed in `plugins.enabled`: `nats/apply.sh` adds it, `nats/apply.sh --disable` removes it |

Until the plugin is listed, all of the above is inert. The pod does not touch the bus.

### No double presence

Two bodies behind one identity split the prompts between them and both heartbeat, and upstream only warns about it.
The vendored adapter, with `HERMES_NATS_REFUSE_DUPLICATE=true` (set in the Deployment), asks the bus whether
`agents.status.hermes.zeus.legacy` already answers after it connects. If it does, it closes the connection and does
not register; the gateway retries with a back-off of 30 s, 60 s, ... up to 5 minutes and never gives up. So
onemind-hermes can be switched on **while Legacy still holds the identity**: it connects, authenticates, probes and
waits (that already proves the credentials, the CA, the network policy and the permissions), and it takes the
identity over by itself when Legacy lets go.

### Order

1. **Merge** (Argo rolls the Deployment; the pod restarts once, with a short Telegram and Matrix gap).
2. `nats/apply.sh --status`: every NATS variable `set`, `gateway deps ... installed`.
3. `nats/apply.sh --dry-run`, then `nats/apply.sh --restart` (copies `nats-platform`, lists it in `plugins.enabled`,
   restarts only the gateway process). Check that onemind-hermes is **connected but not registered**: its
   connection is in `/connz` (client name `onemind-hermes-nats-platform`, user `onemind`) and it holds no
   `agents.prompt...` subscription; the gateway log says the identity is live elsewhere.
   Read `/connz` without any new access: `kubectl get --raw "/api/v1/namespaces/nats/pods/nats-0:8222/proxy/connz?subs=detail&auth=true"`
   (a client is on one of the three servers; look at all three).
4. **Owner: let Legacy go.** Either
   - scale `hermes-legacy` to 0 (this also stops its leaf, `nats-context` and `tak-platform`), or
   - if Paperclip stays on Legacy for now: in Legacy's live `config.yaml` set `platforms.nats.enabled: false` and
     `platforms.tak.enabled: false` and restart its gateway, as was done for Telegram and Matrix. (Legacy's
     `nats-context` keeps listening; it publishes nothing but a stats ping.)

   Legacy has no duplicate guard (it runs the upstream adapter), so it must be **disabled or stopped**, not just left
   alone: if it restarted with `nats` enabled while onemind-hermes holds the identity, the prompts would split.
5. onemind-hermes registers by itself within 30 s to 5 minutes; to take it immediately, `nats/apply.sh --restart`.
6. **Verify** (all four, the first before step 4 so there is something to compare):
   - `$SRV.INFO.agents` / `nats req '$SRV.PING.agents'` lists one `hermes` instance with the new `instance_id`;
   - `agents.hb.hermes.zeus.legacy` heartbeats every 30 s from the new instance only;
   - the Agent Console (`agent-fabric`) and omos-ui show Legacy as live;
   - a prompt to `agents.prompt.hermes.zeus.legacy` (Agent Console chat, or the Synadia SDK) streams a reply, and a
     slash command (`/help`) answers.
7. **Rollback**, at any point: `nats/apply.sh --disable --restart` on onemind-hermes, then scale Legacy back up or
   flip `platforms.nats.enabled` back to `true` on Legacy and restart its gateway. Backups of every file and
   `config.yaml` the script changes are under `/opt/data/backups/legacy-nats-<stamp>/`.

A rehearsal under a different identity is possible without any other change: temporarily set
`HERMES_NATS_SESSION_NAME` to `legacy-next` in `deployment.yaml` (a new subject family, so no collision and no
guard), prompt `agents.prompt.hermes.zeus.legacy-next`, then set it back.

### Owner decisions and later work

- **Paperclip** decides whether step 4 can scale Legacy to 0 (this PR does not touch Paperclip).
- **`nats-context`**: `--with-context` brings back the live picture and the generic NATS tools. Recommended to wait for
  PRs #562 and #564, because with the shared `onemind` user nothing on the server side limits what the model can
  publish (the SOUL rule that physical actions need `execute` is the only guard, the same as on Legacy today).
- The shared `onemind` credential is readable by the agent's shell (the terminal backend is `local`). PR #562's
  `hermes` user is the fix; until then this is the same exposure Legacy had through `hermes-env`.
- After Legacy is retired (step 8): delete the Legacy rule in the `nats-leafnodes` policy and the `hermes-leaf-conf`
  Secret, and stop `leaf_tak` being shared.

## Chat channels (Telegram and Matrix), moved 2026-10-29

Same bot identities, switched from Legacy to onemind-hermes. Only one instance may poll a
Telegram bot (a second gets 409 Conflict), so Legacy has `telegram` and `matrix` disabled
(live `config.yaml` only, because Legacy's channel config was never committed) and keeps
`nats` and `tak`.

| Item | Where |
|---|---|
| `TELEGRAM_BOT_TOKEN`, `MATRIX_ACCESS_TOKEN` | `onemind-hermes-runtime` SealedSecret, copied from Legacy's `/opt/data/.env` and `config.yaml` (neither was in a Kubernetes Secret) |
| Telegram allowlist and home channel | `config.defaults.yaml`: only `8297858681` (Zeus) |
| Matrix allowlist | `config.defaults.yaml`: only `@zeus:onemindos.dev`; the bot is `@legacy:onemindos.dev` |
| Matrix homeserver | in-cluster `http://synapse.matrix.svc.cluster.local:8008`; egress in `network.yaml`; the `matrix` namespace has no NetworkPolicy that blocks ingress |
| Approvals | unchanged: `manual`, and the SOUL rule that deletes and physical actions need Zeus to type `execute` or `override` |

Legacy's `GATEWAY_ALLOW_ALL_USERS=true` (in `hermes-env`) is not carried over: onemind-hermes
pins it `false`. E2EE is off, as on Legacy. `config.defaults.yaml` is copied only once into an
empty PVC, so edit the live `/opt/data/config.yaml` as well.

## Open decisions for the owner

- Which of the 10 paused cron jobs come back, if any.
- ~~Whether Telegram and Matrix stay as Hermes channels.~~ Decided 2026-10-29: they stay, moved to onemind-hermes (see "Chat channels" above).
- Which personal/ops skills (Mac, S24, Windows, WordPress, credit dispute) belong
  on the cloud agent at all.
- Whether `miro`, `cloudflare` and `raindrop` MCP are still wanted.

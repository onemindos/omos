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
7. **NATS cutover.** Port the `nats-platform`/`nats-context`/`tak-platform`
   plugins without the startup hot-patch, give onemind-hermes the leaf sidecar and
   `HERMES_NATS_SESSION_NAME`, scale Legacy to 0, and verify omos-ui discovery
   and heartbeats on `agents.*`.
8. **Retire Legacy.** Keep a VolSync snapshot of `hermes-data`, then, with
   owner approval, remove the Legacy Deployment and Application. Rotate the
   credential files listed above.

## Open decisions for the owner

- Which of the 10 paused cron jobs come back, if any.
- Whether Telegram and Matrix stay as Hermes channels.
- Which personal/ops skills (Mac, S24, Windows, WordPress, credit dispute) belong
  on the cloud agent at all.
- Whether `miro`, `cloudflare` and `raindrop` MCP are still wanted.

# OneMind skills for `onemind-hermes`

Curated from the live Legacy pod (`/opt/data/skills`, reviewed 2026-10-29) for the
single-agent instance `onemind-hermes` (persona: Legacy, one soul). The inventory
is "Skills not shipped with the stock image" in `platform/hermes/onemind-hermes/MIGRATION.md`.

Rules used for curation:

- One agent. The ten role skills (Grid, Guardian, Oracle, Trinity, Forge, Empire,
  Spartan, Eden, Haven, Heritage) were separate personas. Their procedures are
  now sections of two skills; no skill tells the agent to delegate to a sibling.
- Facts were checked against `CLAUDE.md` (Argo/Git first, SealedSecrets, tailnet-only
  TAK, GHCR, no Harbor/ESO/OpenClaw/Bedrock-only). Anything that described a
  retired thing was dropped or rewritten.
- Every kept file was scanned for tokens, keys, passwords, private addresses with
  credentials, chat ids, phone numbers and e-mail addresses. None are in this
  directory. Binaries, caches and files over 200 KB were not copied (none of the
  kept files were that large; the largest candidate trees were dropped).
- Personal data (family, home address, personal accounts) is not committed.

## Shipped here (11 skills)

| Skill | Path | Built from | Note |
|---|---|---|---|
| onemind-domains | `devops/onemind-domains` | so-grid, so-guardian, so-oracle, so-trinity, ge-empire, hp-spartan, le-eden, le-haven, le-heritage | One skill, a section per domain. |
| onemind-code-orchestration | `software-development/onemind-code-orchestration` | so-forge | Worker/PR orchestration. |
| doks-cluster-ops | `devops/doks-cluster-ops` | doks-cluster-ops | Rewritten; 3 of 18 references kept (stateful triage, RWO deadlock, GHCR rollout). |
| nats-operations | `devops/nats-operations` | nats-operations | Rewritten to the locked grammar and leaf-sidecar facts. |
| ha-nats-bridge | `devops/ha-nats-bridge` | ha-nats-bridge | Subject contract only; install steps with leaf/HA credentials removed. |
| tak-server | `devops/tak-server` | tak-server | Rewritten (the original was a July 2026 public-LB snapshot with credentials). Keeps `references/tak-cot-format.md`. |
| omos-feed-builder | `devops/omos-feed-builder` | omos-feed-builder | Trimmed: Mac paths, stale feed inventory, ESO text removed. |
| hermes-mcp-config | `devops/hermes-mcp-config` | hermes-mcp-config | Rewritten for v0.22.5 (`hermes mcp ...`, Git as source of truth). |
| composio | `devops/composio` | composio, composio-mcp-patterns | Rewritten without the consumer key or account ids. |
| hermes-plugin-pattern | `software-development/hermes-plugin-pattern` | hermes-plugin-pattern | Copied as is (clean). |
| hermes-deployment-topology | `autonomous-ai-agents/hermes-deployment-topology` | hermes-deployment-topology | Copied as is; one example password replaced by a placeholder. |

## Decision table

Counts (69 candidates): KEEP 20 (nine role skills merged into one, forge into
another, ten platform skills shipped or folded), DROP 26, NEEDS-OWNER 23. The 38 skills under `.archive/` were not individually
reviewed and stay behind.

### KEEP (shipped or merged)

| Candidate | Decision | Reason |
|---|---|---|
| so-grid | KEEP (merged) | Infra procedures kept in `onemind-domains`; node count/paths updated. |
| so-guardian | KEEP (merged) | Security checks kept in `onemind-domains`. |
| so-oracle | KEEP (merged) | Analytics rules kept; ClickHouse port and Composio account ids dropped. |
| so-trinity | KEEP (merged) | Bridge checks kept; the retired public TAK LB address dropped. |
| ge-empire | KEEP (merged) | Business-ops rules kept; personal e-mail and ClickUp/Composio ids dropped. |
| hp-spartan | KEEP (merged) | Kept as a short section; the 07:00 cron is not assumed. |
| le-eden | KEEP (merged) | Kept; old `rad.*` subjects flagged as outside the locked grammar. |
| le-haven | KEEP (merged) | Kept; home IP, SSH user and device ids dropped. |
| le-heritage | KEEP (merged) | Kept as a short section; family names and vault paths dropped. |
| so-forge | KEEP | Becomes `onemind-code-orchestration`; stale repo names fixed. |
| doks-cluster-ops | KEEP (rewritten) | Core ops guidance; stale items (apply-then-sync, token minting, ESO, 1Password, kubectl delete rs as routine) removed. |
| nats-operations | KEEP (rewritten) | Needed; "no JetStream" and benchmark claims replaced by verify-first guidance. |
| ha-nats-bridge | KEEP (trimmed) | Subject contract is current; embedded leaf credentials and HA token flow removed. |
| tak-server | KEEP (rewritten) | Needed; rewritten to tailnet-only topology and the Helm CLI release. |
| omos-feed-builder | KEEP (trimmed) | Mostly design rules; Mac vault path and stale feed list removed. |
| hermes-mcp-config | KEEP (rewritten) | Needed for the "add composio and github MCP" step; Zoho key-in-URL table and pod hacks removed. |
| composio | KEEP (rewritten) | Needed for the migration plan; consumer key and account ids removed. |
| composio-mcp-patterns | KEEP (folded) | Useful header/session facts folded into `composio`; the file held the consumer key. |
| hermes-plugin-pattern | KEEP | Clean, generic, useful for the nats-platform port. |
| hermes-deployment-topology | KEEP | Clean reference to upstream env vars. |

### DROP

| Candidate | Reason |
|---|---|
| onemind-stack | Historical stack decision log (Dropbox/ClickUp/Zoho/Circle); contains a WordPress app password. |
| onemind-fabric-deployment | Multi-agent profile scaffolding (5 profiles, pillar models); superseded by one soul. |
| onemind-inter-agent-protocol | Mother-to-pillar delegation between sibling agents; no siblings now. |
| tak-mission-ops | July 2026 public-TAK material, default passphrases, Mac/WinTAK/ADB flows; superseded by `tak-server` and the runbook. |
| tak-mission-runner | One-off "Gas Station Run" live protocol with a real location and device UID. |
| mission-replay | Depends on a pod script with embedded database credentials and the retired public TAK address. |
| omos-ui-build-deploy | Describes `:latest`, a Cloudflare tunnel and creating a pull secret from a personal token. |
| legacy-integrations | Zoho/GHL key-in-URL setup; contains keys; per-integration setup belongs in SealedSecrets. |
| legacy-model-config | ESO plus 1Password key flow (retired); contains a key pattern. Sealed secrets own `hermes-env`. |
| cloudflare-access-patterns | Global-key fallback and personal e-mail; upstream Cloudflare skills cover this. |
| tailscale-pod-access | SSH from the pod to the Mac, Windows and phone, with a pinned key; wrong model for a cloud agent. |
| tailscale-continuous-mesh-access | Userspace tailscaled lifecycle hacks. |
| composio-clickup, composio-github, composio-gmail, composio-google | Personal account ids baked in; re-derive from `COMPOSIO_SEARCH_TOOLS`. |
| omos-cortex-neurolinked | NeuroLinked simulator on the Mac; speculative. |
| omos-legacy-genome | Superseded by `platform/hermes/agent/` in Git. |
| legacy-self-awareness | Large (26 files) snapshot of a retired Legacy deployment (CloudTAK, ESO, argocd-operations); contains credentials. |
| openclaw-k8s-deploy | OpenClaw is retired. |
| repo-absorption | One-time repo consolidation. |
| resume-session-on-pod | Mac-terminal helper for the old pod. |
| onemind-xr-platform, onemind-product-strategy | Speculative product material. |
| merge-reconciler, hermes-mobile-desktop-sync | Speculative or already stale (inventory list). |

### NEEDS-OWNER (not copied; owner decides)

Personal or machine-specific. Listed with the owner's call needed.

| Candidate | Why it needs a decision |
|---|---|
| credit-bureau-dispute | Personal legal/financial letters. |
| windows-network-drop, windows-zeus-access | Owner's Windows laptops; remote access. |
| s24-app-control, s24-termux-ssh, android-adb-management | Owner's phone; adb/SSH from the pod. |
| mac-cloud-sync, mac-disk-hygiene, mac-repo-hygiene, mac-storage-ops | Owner's Mac housekeeping; not a cloud-agent job. |
| find-zeus | Personal location/contact lookup. |
| wordpress-legacy, wordpress-on-digitalocean, zeus-site | Personal site; held a WordPress app password and a CMS database credential (rotate). |
| onepassword-legacy | 1Password on the old pod; secrets here are SealedSecrets. |
| notion-workspace, notion-onemind-architecture | Notion routing; Notion migration status unknown; one file held a token (rotate). |
| vault-dedup-audit, external-drive-cloud-sync, session-librarian | Personal vault/drive utilities. |
| modesty-ai-daily-briefing, grounddata-ai-daily-news | Daily-news jobs; the grounddata scripts embed secrets and are 248 KB. |
| matrix-doks-ops | Technically current, but whether Matrix stays a Hermes channel is an open owner decision; the template held a credential. |

## Secrets found in candidates (paths only, values not reproduced)

All are on the Legacy pod volume, in skills that were dropped or rewritten here.
Treat each as exposed on that volume and rotate on retirement of Legacy.

- `devops/composio/SKILL.md`, `devops/composio-mcp-patterns/SKILL.md`, `devops/composio-mcp-patterns/references/mcp-http-session.md`, `devops/legacy-integrations/references/composio-integration.md`: Composio consumer key.
- `devops/legacy-integrations/SKILL.md`: integration keys.
- `devops/ha-nats-bridge/SKILL.md`: NATS leaf user/password and a Home Assistant token reference.
- `devops/onemind-stack/SKILL.md`, `devops/wordpress-legacy/SKILL.md`: WordPress application password.
- `devops/tak-server/SKILL.md`, `devops/tak-server/references/*`, `devops/tak-mission-ops/*`, `devops/mission-replay/SKILL.md`: TAK admin web password, certificate passphrases, database password.
- `devops/notion-workspace/SKILL.md`, `devops/notion-workspace/references/notion-mcp-setup.md`: Notion token.
- `devops/legacy-model-config/SKILL.md`: DigitalOcean inference key pattern.
- `devops/doks-cluster-ops/references/tailscale-egress-device-map.md`, `devops/legacy-self-awareness/references/*`, `devops/zeus-site/references/payload-db-ops.md`: credentials in examples.
- `matrix-doks-ops/templates/homeserver-clean.yaml`: Synapse secret.
- `grounddata-ai-daily-news/scripts/lib/engagement_state.py`, `engagement_delivery.py`: embedded tokens.

## How they are installed into `onemind-hermes`

Verified on the predecessor `hermes-fresh` (same image) (Hermes v0.22.5):

- `/opt/hermes/cli-config.yaml.example` documents `skills.external_dirs`: extra
  directories whose skills are loaded without copying them into `~/.hermes/skills/`.
  They are read-only, entries are `~`/`${VAR}` expanded, and a local skill wins
  when names collide. Skills are found by walking the tree (`os.walk`, symlinks
  followed), so the `category/name/SKILL.md` layout here works unchanged.
- `/opt/hermes/tools/skills_sync.py` seeds only the bundled skills into
  `/opt/data/skills`, tracked by `.bundled_manifest`; it also skips names provided
  by `external_dirs`, so they are never shadowed or overwritten.

Recommended mechanism: `external_dirs`.

1. Make this tree available in the pod at a fixed path, for example
   `/opt/onemind-skills`: an init container (or a git-sync sidecar) that copies
   `platform/hermes/skills/` from the repo revision Argo deploys. At about 110 KB
   it also fits a handful of ConfigMaps, but nested directories need explicit
   `items` mappings, so a copy step is simpler.
2. In the instance config (seeded from Git, not edited in the pod):

   ```yaml
   skills:
     external_dirs:
       - /opt/onemind-skills
   ```

3. Restart or `/reset`; check with `hermes skills list` and `skill_view`.

Copying directly into `/opt/data/skills` also works (user-added directories are
left alone) but drifts from Git and is hand-maintained, which is what the migration
plan wants to avoid.

## Bundled skills proposed for disabling on the new instance

Bundled count on the stock image is 58. Keep `devops`, `software-development`, `research`
(except the two noted), `productivity` (document skills) and `autonomous-ai-agents`
(except `computer-use`), as requested.

Mechanism: `hermes skills config` is the interactive picker; the persisted form is
`skills.disabled` (a list of names) in `config.yaml`, with `skills.platform_disabled.<platform>`
for per-channel lists. Both are read by `agent/skill_utils.py`. Essential skills cannot
be disabled. Put the list in the Git-seeded config so it survives rebuilds.

| Group | Disable | Reason |
|---|---|---|
| apple (4) | apple-notes, apple-reminders, findmy, imessage | macOS-only CLIs (`memo`, `remindctl`, FindMy.app, `imsg`); the pod is Linux. |
| creative | ascii-video, baoyu-infographic, manim-video, p5js, popular-web-designs, songwriting-and-ai-music, design-md, claude-design | Generative art, video and design systems; not part of C2/ISR operations. Keep `architecture-diagram` and `humanizer` (diagrams and clean prose are useful). |
| media | gif-search, songsee | Entertainment/audio analysis. `youtube-content` optional: keep if research use is wanted. |
| social-media | xurl | X/Twitter posting needs an account and credentials nobody has provisioned. |
| email | himalaya | Raw IMAP/SMTP CLI; mail goes through Composio. Keep `email-inbox-triage` if a mailbox is connected. |
| note-taking | obsidian | Needs the owner's vault, which is not on this instance. |
| productivity | airtable, box, teams-meeting-pipeline, product-price-monitor | No Airtable, Box or Teams integration is planned; price monitoring is off-mission. |
| productivity (decide) | google-workspace, notion | Overlap Composio (Google) and the Notion decision; enable only if the owner wants the direct path. |
| research | competitor-news-monitor, llm-wiki | Not tied to a current job; re-enable on request. Keep `arxiv`, `grounded-citations`. |
| software-development | inspecting-hermes-desktop-dom | Desktop-app DOM inspection; irrelevant to a pod. |
| autonomous-ai-agents | computer-use | Drives a desktop GUI; a headless pod has none. |

Keep as is: `devops/sdlc-review`, `software-development` (the rest),
`research/arxiv`, `research/grounded-citations`, `productivity` documents
(`docx`, `pdf`, `powerpoint`, `xlsx`, `maps`, `meeting-action-items`,
`document-to-action-items`, `weekly-review-planning`), `autonomous-ai-agents`
(`claude-code`, `codex`, `opencode`, `hermes-agent`), `web/blocked-page-recovery`.
That disables 25 and leaves 33 of the 58 bundled skills enabled (31 if both
"decide" rows are also disabled).

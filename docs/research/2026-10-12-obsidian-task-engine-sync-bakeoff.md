# Obsidian task engine: sync and version-control bake-off

Research date: 2026-10-11/12, America/New_York. Status: **proposal, paused by the user
on 2026-10-12**. Nothing here is deployed. The staged manifests in
[`platform/vault/`](../../platform/vault/README.md) are not referenced by any Argo CD
Application.

This document records a user direction that differs from the earlier recommendation in
[`self-hosted-capability-plan.md`](../self-hosted-capability-plan.md) (Outline + Plane +
Vikunja). It does not silently replace that plan: the decision becomes final only after
the bake-off below passes and the capability-plan gates are re-evaluated.

## 1. Decision context

User requirements stated in the 2026-10-11 session (in addition to
[`MISSION.md`](../MISSION.md)):

- Replace Notion, Todoist and ClickUp for personal **and** business tasks/projects.
- Obsidian becomes the primary human interface on **macOS and Android** (iOS not planned;
  Windows/Linux possible later). `omos-ui` is not used.
- The task system must run **24/7**, sync across devices, keep **version history**, and
  **must not lose data** ("I cannot afford to lose anything").
- Must work **offline at the homestead** (home cluster, local Wi-Fi mesh, no internet),
  later over mesh networks, and reconcile with DOKS when connectivity returns.
- AI-native: agents read/write tasks; webhooks in and out; events usable for a COP.
- Private distribution of OneMind plugins/theme is fine (BRAT / git); no community-store
  release required.

Working architecture agreed in the session (files are the contract, Obsidian is a client):

```
OWNED
 1. Contract  OneMind task schema v1 — Markdown + YAML, TaskNotes-compatible (to be written)
 2. Sync      bake-off winner (Syncthing or Self-hosted LiveSync)
 3. Engine    vault-engine (Go, not yet written): file watcher → NATS events, REST/MCP,
              webhooks in/out, scheduler (reminders, recurrence), index, ICS feed,
              conflict merge, approval gate
 4. Clients we own: web console (tasks + COP), Android capture (HTTP Shortcuts/Tasker → API),
              ntfy push, CloudTAK task layer, agents via MCP, Telegram/Matrix bot
RENTED
 5. Obsidian  editor UI on Mac + Android: TaskNotes, Bases, OneMind theme/plugins
```

Long-term test: if Obsidian disappeared, tasks, reminders, approvals, COP and agents must
keep working through the owned layers. Only the editing experience is lost.

## 2. Verified facts about Obsidian (2026-10-11/12)

| Item | Fact | Source |
|---|---|---|
| Platforms | macOS, Windows, Linux (x64 + ARM64 AppImage, deb, Snap, Flatpak), Android, iOS. Current installer v1.14.4 | [obsidian.md/download](https://obsidian.md/download) |
| Obsidian CLI | Requires installer 1.12.7+; **requires the Obsidian app to be running** (launches it if not) | [Obsidian CLI help source](https://github.com/obsidianmd/obsidian-help/blob/master/en/Extending%20Obsidian/Obsidian%20CLI.md), [help.obsidian.md/cli](https://help.obsidian.md/cli) |
| Obsidian Headless | **Open beta**, `npm install -g obsidian-headless` (Node 22+), command `ob`; talks **only to Obsidian Sync/Publish** and needs an active Sync subscription | [Obsidian Headless](https://github.com/obsidianmd/obsidian-help/blob/master/en/Extending%20Obsidian/Obsidian%20Headless.md), [Headless Sync](https://github.com/obsidianmd/obsidian-help/blob/master/en/Obsidian%20Sync/Headless%20Sync.md) |
| Not provided | No web app, no self-hostable sync server, no server-side plugin runtime (Bases/TaskNotes evaluate only in the app), no official API/webhooks/event stream, no per-user permissions | Derived from the docs above; this is why the owned engine exists |
| Plugins run in-app | Plugin features (including TaskNotes HTTP API/MCP/webhooks) run only while Obsidian runs | [TaskNotes HTTP API doc](https://github.com/callumalpass/tasknotes/blob/main/docs/HTTP_API.md): "Desktop only", "Not available on mobile", `http://localhost:{PORT}` |

### TaskNotes (selected task plugin)

- One Markdown note per task with YAML frontmatter; every view is a Bases query; RRULE
  recurrence with per-instance completion; time tracking; dependencies; calendar views;
  Google/Microsoft OAuth and ICS calendar sync. [README](https://github.com/callumalpass/tasknotes)
- HTTP API + MCP + webhooks (`task.created|updated|deleted|completed|archived|unarchived`,
  time, pomodoro, recurring and reminder events) — **desktop only**.
  [HTTP API](https://github.com/callumalpass/tasknotes/blob/main/docs/HTTP_API.md),
  [Webhooks](https://github.com/callumalpass/tasknotes/blob/main/docs/webhooks.md)
- Registry: 2,051,848 downloads (2026-10-11); repo 2,203 stars, MIT, pushed 2026-10-06.
- Consequence: an edit made on Android never fires a TaskNotes webhook. **The engine,
  watching files server-side, must be the single source of events.**

## 3. Sync plugin review (all user-supplied candidates)

Data pulled 2026-10-12 from the official registry
([community-plugins.json](https://github.com/obsidianmd/obsidian-releases/blob/master/community-plugins.json),
[community-plugin-stats.json](https://github.com/obsidianmd/obsidian-releases/blob/master/community-plugin-stats.json)),
each plugin's `manifest.json`, GitHub repo metadata and README. "Unreviewed" means the
registry description carries "This plugin has not been manually reviewed by Obsidian staff."

Hard filters: (1) Android + macOS, (2) self-hosted on our clusters so it works offline,
(3) plain-file copy on the cluster for git/engine, (4) safe merges, deletes and history.

| Plugin | Downloads | Repo / license | Maintainers | Verified facts | Verdict |
|---|---|---|---|---|---|
| [github-sync](https://community.obsidian.md/plugins/github-sync) | 53,423 | [kevinmkchin/Obsidian-GitHub-Sync](https://github.com/kevinmkchin/Obsidian-GitHub-Sync), MIT | 1 main | `isDesktopOnly: true` in manifest | Rejected: no Android |
| [yaos](https://community.obsidian.md/plugins/yaos) | 38,628 | [kavinsood/yaos](https://github.com/kavinsood/yaos), 0BSD | 1 | Server is a Cloudflare Worker (Yjs CRDT); attachments need R2; unreviewed | Rejected: cannot run on our clusters / offline |
| [powerconnect](https://community.obsidian.md/plugins/powerconnect) | 2,407 | [obsidian-power-connect](https://github.com/obsidian-power-plugins/obsidian-power-connect), MIT | 1, created 2026-08-02 | Dropbox/OneDrive/Google Drive only; has delete guard, dry-run, E2E; unreviewed | Rejected: not self-hosted |
| [vaultbox](https://community.obsidian.md/plugins/vaultbox) | 2,262 | [grumpydev/vaultbox](https://github.com/grumpydev/vaultbox), MIT | 1, 5 stars | Dropbox only; unreviewed | Rejected: not self-hosted |
| [air-sync](https://community.obsidian.md/plugins/air-sync) | 7,960 | [takezoh/obsidian-air-sync](https://github.com/takezoh/obsidian-air-sync), MIT | 1 main | Google Drive/OneDrive/Dropbox; default sign-in uses third-party token server `auth-airsync.takezo.dev`; unreviewed | Rejected |
| [sync-vault-ce](https://community.obsidian.md/plugins/sync-vault-ce) | 14,097 | [abcamus/obsidian-sync-vault-ce](https://github.com/abcamus/obsidian-sync-vault-ce); README badge AGPL-3.0, GitHub detects NOASSERTION | 1 | Cloud drives (examples: Aliyun, Baidu Netdisk); VFS streaming; MCP | Rejected |
| [syncthing-integration](https://community.obsidian.md/plugins/syncthing-integration) | 48,411 | [LBF38/obsidian-syncthing-integration](https://github.com/LBF38/obsidian-syncthing-integration), MIT, 358 stars | 7 | Status + conflict diff for a running Syncthing; does not sync by itself; reviewed | **Optional helper** (preferred) |
| [syncthing-manager](https://community.obsidian.md/plugins/syncthing-manager) | 8,588 | [gustjose/obsidian-syncthing-manager](https://github.com/gustjose/obsidian-syncthing-manager), MIT, 40 stars | 1 main | Conflict resolver, version restore; **Android requires disabling HTTPS in Syncthing GUI**; unreviewed | Helper, not preferred |
| [nextcloud-sync](https://community.obsidian.md/plugins/nextcloud-sync) | 18,081 | [siosig/obsidian-nextcloudsync](https://github.com/siosig/obsidian-nextcloudsync), MIT | **1** (587/588 commits) | Created 2026-06-08; latest `1.0.12-beta.2`; If-Match concurrency, trashbin deletes, server version history, plain files server-side; unreviewed | Watch: good design, too young |
| [fast-note-sync](https://community.obsidian.md/plugins/fast-note-sync) | 163,972 | plugin [haierkeys/obsidian-fast-note-sync](https://github.com/haierkeys/obsidian-fast-note-sync) + server [fast-note-sync-service](https://github.com/haierkeys/fast-note-sync-service), Apache-2.0 | 1 dominant | Server: note history, trash restore, git auto-push, S3/WebDAV backups, OIDC web admin, REST, MCP, SQLite/MySQL/Postgres. **E2E on roadmap**; roadmap also lists "optimize offline note merging and introduce conflict resolution". Last push 2026-08-14; unreviewed | Runner-up / watch |
| [obsidian-livesync](https://community.obsidian.md/plugins/obsidian-livesync) | **1,000,636** | [vrtmrz/obsidian-livesync](https://github.com/vrtmrz/obsidian-livesync), MIT, 12,696 stars | **64 contributors**, 375 versions, pushed 2026-10-05 | CouchDB or S3 object storage or WebRTC P2P; E2E; merges simple conflicts. Docs: "For mobile devices, Obsidian requires a valid SSL certificate." README: back up the vault before installing/upgrading | **Finalist B** |
| Syncthing (no plugin) | core 89,317 stars | [syncthing/syncthing](https://github.com/syncthing/syncthing), MPL-2.0 | large | Android: official app **discontinued** (last release Dec 2024, [repo archived](https://github.com/syncthing/syncthing-android)); maintained fork [researchxxl/syncthing-android](https://github.com/researchxxl/syncthing-android) (3,038 stars, pushed 2026-10-09) via [F-Droid](https://f-droid.org/packages/com.github.catfriend1.syncthingfork), GitHub releases or Obtainium — **not Google Play** | **Finalist A (recommended)** |

Also checked earlier in the session (not sync, recorded for completeness):
[remotely-save](https://github.com/remotely-save/remotely-save) last push 2024-11-10 (avoid);
[obsidian-webhooks](https://github.com/trashhalo/obsidian-webhooks) last push 2022-07-13 (avoid);
[obsidian-git](https://github.com/Vinzent03/obsidian-git) README: mobile "highly unstable",
no SSH on mobile, limited repo size — **desktop-only use**.

### Why the recommendation changed from LiveSync to Syncthing

Earlier in the session LiveSync + `livesync-bridge` was the working choice. Verification found:

1. [`vrtmrz/livesync-bridge`](https://github.com/vrtmrz/livesync-bridge) (the component that
   would give the cluster a plain-file copy) has **no LICENSE file** and **no releases**
   (142 stars, pushed 2026-10-05). Without a license it is all-rights-reserved by default.
2. LiveSync stores the vault as chunks in CouchDB, not files; plain files on the cluster
   depend on the bridge.
3. Mobile requires a valid TLS certificate
   ([setup_own_server.md](https://github.com/vrtmrz/obsidian-livesync/blob/main/docs/setup_own_server.md)),
   which needs a real CA certificate and local DNS for offline homestead operation.
4. LiveSync is a plugin, so it syncs only while Obsidian runs. Syncthing is an OS
   service on Android/macOS and syncs with Obsidian closed.

| | Syncthing (A) | Self-hosted LiveSync (B) |
|---|---|---|
| Cluster copy | Native plain files on the PVC | CouchDB chunks; plain files only via unlicensed bridge |
| Offline / LAN / mesh | Built in: P2P, local discovery | Needs CouchDB reachable with valid TLS |
| Runs with Obsidian closed | Yes | No |
| Concurrent edits | `.sync-conflict-*` copies, no merge | Auto-merge simple conflicts |
| Real time | Seconds | Real time |
| E2E | TLS in transit; untrusted (encrypted) devices supported | Built-in E2E |
| Android distribution | F-Droid / GitHub / Obtainium | Inside Obsidian |
| Independent of Obsidian | Yes | No |

Syncthing's weakness (no merge) is closed by the engine: it holds full git history, so
it can three-way merge `.sync-conflict-*` files against the last committed base, or open
an approval task when it cannot.

Plugins for the base install: **TaskNotes** (required), **Syncthing Integration**
(optional), no sync plugin. Fallback if Syncthing fails the bake-off: **Self-hosted
LiveSync** with a real TLS certificate on CouchDB.

## 4. Never-lose layer (independent of the sync choice)

Sync replicates deletions and mistakes. Recovery comes from separate layers:

```
Mac ⇄ Android ⇄ cluster Syncthing peer (PVC = plain vault)
                      │
 1. vault-git sidecar: commit every 60 s into a git dir OUTSIDE the synced folder
    (GIT_DIR=/var/syncthing/git/vault.git, work tree=/var/syncthing/vault) → push to remote
 2. Syncthing staggered file versioning on the cluster peer
 3. VolSync restic snapshots → DO Spaces (existing platform-volsync / platform-volume-backups)
 4. Mac: Time Machine
 5. Engine delete guard: pause the cluster folder + alert when > N deletions in a window
```

Open choice for the git remote: an in-cluster Forgejo (works offline at the homestead)
versus a private GitHub repository (off-site, internet required). Recommendation: both —
Forgejo as primary, GitHub mirror.

## 5. Bake-off plan

Two isolated test vaults with synthetic data only; no real notes until a winner passes.

Setup: Syncthing peer (A) and CouchDB (B) from `platform/vault/`, one Mac and one Android
device, Syncthing-Fork from F-Droid, LiveSync from the community store.

| # | Test | Pass condition |
|---|---|---|
| 1 | Android offline edit and agent (cluster-side) edit of the same task, then reconnect | No content lost; result merged or both versions preserved and detected |
| 2 | Delete 200 files on one device | Recoverable from git and versioning; delete guard alerts |
| 3 | Rename/move folders on Mac while Android offline | No duplicates or orphaned copies after reconnect |
| 4 | 2,000 TaskNotes task notes | Android opens TaskNotes views acceptably; sync completes |
| 5 | Delete the sync pod / restart the cluster mid-sync | No corruption; sync resumes |
| 6 | Restore a file deleted 3 days earlier | Restored from git within 5 minutes |
| 7 | 72-hour WAN cut (home cluster or simulated) | LAN sync continues; DOKS reconciles after reconnect |
| 8 | Obsidian closed on Android for 24 h while tasks change | Files present when Obsidian opens (A expected yes, B expected no) |
| 9 | Attachment (PDF ≤ 20 MB) added on Android | Arrives on Mac and cluster; committed or excluded by policy |

**Pass = zero data loss and every change recoverable from git.** Record results in this
file before deciding.

Also re-run the personal workspace gates from
[`self-hosted-capability-plan.md`](../self-hosted-capability-plan.md) (calendar, habits,
attachments, AI access, webhooks, mobile, offline) against the TaskNotes + engine design.

## 6. Engine contracts (to design; no new NATS prefixes)

Follow [`NATS-DOCTRINE.md`](../NATS-DOCTRINE.md) §2; domains are the vault pillars
`so`/`hp`/`le`/`ge` (+ `sys`). The doctrine's own example `evt.so.task.created` fits:

```
evt.{pillar}.task.{created|updated|completed|deleted}   facts from file changes
cmd.{pillar}.task.{create|update|complete}              requests the engine may reject
svc.vault.task.{list|get}                               request/reply queries
```

Approval gate: tasks with `approval: required` stop at `awaiting-approval`; only after
human approval does the engine publish the associated `cmd.*`. Webhooks are adapters on
these subjects (HMAC-signed outbound, authenticated inbound via n8n or the engine).

## 7. File policy for the vault

- Keep: `.md`, `.base`, `.canvas`, images used in notes, PDFs that are read/annotated
  (≤ ~20 MB), short voice memos.
- Convert to Markdown, archive the original: `.docx`, `.xlsx`, `.html`.
- Move to object storage and link: video, music/podcasts, `.zip`, `.apk`, installers,
  exports, raw drone imagery.
- Sensitive scans: separate folder excluded from agent access, or outside the ops vault.
- Personal and business: two folders/vaults so business work can be shared later without
  exposing personal notes; the engine indexes both.

## 8. Open questions before resuming

1. Will anyone other than Zeus and agents create/update business tasks? (Decides whether a
   shared mirror such as ClickUp/Plane is needed.)
2. Git remote: Forgejo in-cluster, GitHub private, or both.
3. Home cluster hardware and when it exists (offline tests 7–8 need it or a simulation).
4. Mesh: [`self-hosted-capability-plan.md`](../self-hosted-capability-plan.md) records that
   Headscale is not a drop-in replacement for the Tailscale operator used here; offline
   tailnet coordination needs the NetBird/Nebula pilot decision or LAN-direct addresses.
5. TLS for LiveSync (only if B is chosen): real certificate via DNS-01 vs Tailscale ingress cert.

## 9. Resume checklist

1. Answer the open questions above.
2. Create secrets out of band (see [`platform/vault/README.md`](../../platform/vault/README.md)).
3. Add an Argo CD Application for `platform/vault` (not included on purpose).
4. Run the bake-off; record results here.
5. Write OneMind task schema v1 (TaskNotes-compatible JSON Schema under `schemas/`).
6. Build vault-engine MVP; then web console, Android capture, ntfy, CloudTAK layer.
7. Migrate Notion → Todoist → ClickUp only after the gates pass; keep Todoist 30 days.

# Hermes recoverability plan

**Status (2026-10-16):** decisions made (D1 to D8, section 4). Step 0, Step 1 and most of Steps 2 and 3 are done; see the progress table below. The export is blocked on rotating the credentials found in the live config and skills.
**Written:** 2026-10-16. Live state below was read from the cluster (`do-nyc3-omos-cloud`) and from Zeus's Mac on 2026-10-15/16, read-only.
## Progress

| Step | State |
|------|-------|
| 0 Rescue and scrub | Done (#206) |
| 1 Consolidate into `platform/hermes/agent/`, archive `legacy/` | Done (#207 to #210) |
| 2a Export: script, tests, manual workflow, read-only role, dedicated runner | Built and merged (#211, #212, #213). The runner is online (repo-scoped, label `onemind-hermes-readonly`). The workflow stays off and the first run would refuse until the credentials are rotated |
| 2b Seeding | Script and tests merged (#214), deliberately not wired into the Deployment yet (see `platform/hermes/seed/README.md`) |
| 2c OpenClaw capture | Not started |
| 2d Mac jobs (weekly memory copy, edge identity) | Not started |
| 3 Restore drill | Done: application-level restore passed on 2026-10-16 (`hermes-restore-drill-2026-10-16.md`) |
| 3 Apply the edited soul and the memory correction to the live agent | Done 2026-10-16. `SOUL.md` on the volume now equals `agents/legacy/SOUL.md`; the stale repository entry in `MEMORY.md` was replaced. Safety copies of both old files are on the volume (`*.pre-edit-20261016`) |

Still open: the credential rotation and cleanup (Todoist), removing the WordPress password from memory, the `op` CLI not being on Hermes' PATH, and the dead 929 MiB `state_compact.db` (needs Zeus's go-ahead to remove).

**Goal:** make Hermes' soul, skills, configuration and memory recoverable. If the pod, the volume or the cluster is lost, Hermes (and the agents that run beside it) can be brought back to the state they were in, from Git plus backups, without relying on anyone's memory.

This plan also decides what happens to the `legacy/` directory, which was the first attempt at the same goal.

---

## 1. Where things stand

### 1.1 What Hermes is made of

Hermes runs as `hermes-legacy` in namespace `hermes` on the stock `nousresearch/hermes-agent:v2026.9.24` image, with a `nats-leaf` sidecar. All of its identity and state is **one 25 GiB PVC, `hermes-data`, mounted at `/opt/data`**:

| Part | Location on the volume | Size / shape | In Git today |
|------|------------------------|--------------|--------------|
| Soul | `SOUL.md` | 68 lines | **No.** Matches none of the three copies in Git |
| Config | `config.yaml` | 621 lines, 10+ hand-made `.bak` files | No |
| Skills | `skills/` | 201 `SKILL.md` files, 67 MB | No |
| Declarative memory | `memories/MEMORY.md`, `USER.md` | 3.5 KB total | No (decision D2) |
| Conversation and memory DB | `state.db`, `state_compact.db`, `sessions/` | 1.4 GB + 0.9 GB | No (runtime state) |
| Scheduled jobs | `cron/jobs.json` | 10 jobs, all disabled | No |
| Plugins | `plugins/` (`nats-platform`, `nats-context`, `memory-wiki`, `onemind-vision`, `platforms`) | | Only `onemind-vision` (identical to `plugins/onemind-vision/`) |
| Profiles | `profiles/` | 14 folders: 4 with a soul and config (`sovereign-operations`, `holistic-performance`, `legacy-evolution`, `generational-empire`), 10 empty shells (including `forge`) | Partly |
| Clutter | `repos/` 2.6 GB, `home/` 1.2 GB, two clones of this repo about 2 GB, venvs, caches, `OpenClaw/` 906 MB, `takdocs/`, loose research files | about 15 GiB | n/a |

### 1.2 What backs it up

VolSync `ReplicationSource volume-c96e219d2314` backs up `hermes-data` with restic every night at 01:00 (retain 7 daily, 4 weekly, 3 monthly, 1 yearly; CSI snapshot copy method). This works for disaster recovery. Its limits:

- Restore is all-or-nothing, and about 15 GiB of the volume is clutter.
- A **file-level restore drill has been done**: `platform/backups/volsync/README.md` records that on 2026-10-08 Hermes was restored to a new PVC (319,496 files, matching reference-binary checksum) and the disposable resources removed. The application-level check (booting Hermes from a restored copy and confirming it loads its soul, skills and databases) was done on 2026-10-16 and passed: `hermes-restore-drill-2026-10-16.md`. (An earlier draft of this plan said no restore test existed; that was wrong.)

### 1.3 Why Git cannot rebuild Hermes today

- **The seeding mechanism is dead.** `legacy/seed/01-seed-legacy` copied the persona from `/opt/hermes/seed/`. That path only existed in the old custom image. The live pod runs the stock image and `/opt/hermes/seed` does not exist, so the live `SOUL.md` has been edited by hand on the volume. The three Git copies drifted: `legacy/seed/SOUL.mother.md` (38 differing lines from live, all of them the "daughter" versus "edge nodes" wording), `platform/hermes/agent/agents/legacy/SOUL.md` (80), `legacy/agents/legacy/SOUL.md` (125). These counts measure text difference, not meaning.
- **The live soul carries stale facts.** It says the architecture is AWS Bedrock (the live provider is `digitalocean`, model `nemotron-3-ultra-550b`), lists a NATS account layout that is unverified, and describes `om_presence` as a voice bridge (archived in August). A rebuild from this text would reinstate those beliefs.
- **The documented recovery is weak.** `legacy/RUNBOOKS.md` section 10 says that if the volume is lost, "Hermes memory/skills are gone". Legacy's own `MEMORY.md` still says disaster recovery is "clone omos-legacy and build seed/", a repo that no longer exists and a process that no longer works.
- **Some definitions exist only on the pod.** The pod keeps a clone of `onemindos/omos-legacy`. That GitHub repo no longer exists. The clone has **2 unpushed commits from 2026-09-17** that add the four overseer agents (SO, HP, LE, GE): `agents/{SO,HP,LE,GE}/SOUL.md` and 81 added lines of `AGENTS.yaml`. Apart from those and one `shellcheck` line in `seed/launch.sh`, the clone is identical to `legacy/`. The clone's own git history is damaged (a `git bundle` of it cannot be cloned because of a missing object), so it must not be trusted as a backup.

### 1.4 Other agent runtimes, and what is actually running

| Runtime | Agents | State | In Git |
|---------|--------|-------|--------|
| Hermes (pod) | Legacy plus profiles | `hermes-data`, backed up | Manifests only |
| OpenClaw (`openclaw` ns) | `legacy`, `so`, `hp`, `le`, `ge` | `openclaw-home-pvc` (50 GiB), backed up | Manifests only (`platform/openclaw/`) |
| Legacy-Edge (Zeus's Mac) | one edge agent | `~/.hermes` (8.6 GB, nearly all install and cache), launchd `ai.hermes.gateway` | No |
| kagent | `k8s-agent`, `promql-agent`, `observability-agent`, `cilium-debug-agent` | cluster objects | No Agent manifests found in `platform/` |

Verified on 2026-10-16:

- **On the cluster bus,** only `legacy` registers as an agent service. The overseers have a soul, config and `nats` enabled in their profiles but do not register. The Mac's edge agent does not appear in the cluster listing.
- **Forge has no running profile:** `profiles/forge/` is empty, although `PROFILES.md` and `legacy/AGENTS.yaml` call it active.
- **The nine servants** (Eden, Empire, Grid, Guardian, Haven, Heritage, Oracle, Spartan, Trinity) have empty profile folders and paused cron jobs.
- **The Mac edge identity is small:** `SOUL.md` 5 KB (already written as "edge node", not "daughter"), `config.yaml` 24 KB (provider `digitalocean`, model `deepseek-v4-pro`, NATS session `legacy-edge`), 147 skills (the pod has 201, so they have drifted), `USER.md` 1.3 KB, empty cron, no profiles. Time Machine includes `~/.hermes` (latest backup 2026-10-16 02:11), but it is a local disk only.

### 1.5 The `legacy/` directory (146 files)

| Part | Disposition |
|------|-------------|
| `CREED.md`, `PROFILES.md`, `EVOLUTION.md` | **Keep.** Doctrine: hierarchy law, profile architecture, self-improvement loop and PR conventions. Remove "daughter" wording (D4) |
| `INTEGRATION.md`, `RUNBOOKS.md`, pitfalls in `FABRIC.md` | **Keep, distilled** into one Hermes runbook (startup order, dispatch routing, edge delegation, NATS leaf pitfalls). Strip addresses and credentials |
| `AGENTS.yaml` (three versions) | **Rebuild** as one verified registry with a status per agent (D5) |
| SOULs for Legacy, Legacy-Edge, SO, HP, LE, GE | **Keep.** Legacy's is the live file, captured verbatim, then split (D4) |
| `agents/Forge` | **Keep, marked planned** (its profile is empty) |
| SOULs and CONFIGs for the nine servants | **Keep as `agents/planned/`** with a README on the hierarchy, stale model fields flagged (D6) |
| 14 agent folders in the May 2026 OpenClaw/Codex format (`zenith`, `unity`, `eternal`, `sentinel`, `sovereignty`, `apex`, `creed`, `eagleeye`, `genesis`, `helix`, `librarian`, `nomad`, `overwatch`, `pulse`) | **Archive** to the vault, after extracting the pillar and domain map into `PILLARS.md` (D7) |
| `seed/plugins/nats-platform` (29 files) | **Keep** after comparing with the live plugin |
| `seed/Dockerfile`, `launch.sh`, `01-seed-legacy`, `config.hub.yaml`, `hermes_nats_host.py` | **Archive.** They belong to the retired custom-image path |
| `ARCHITECTURE.md`, `SERVICES.md`, `MACHINES.md`, `FABRIC.md`, `services/*.md`, `apps/`, `nats/`, `cognitive/` | **Archive.** July/August snapshots superseded by `docs/STACK.md` |
| `docs/tak-3d-alternatives.md`, `deploy/takserver/tak-webview-cesium.yaml` (not running), `config/wintak-omos-dell.xml`, `scripts/optimize-dell-win10.ps1`, `scripts/activate-servant.sh` | **Archive** |

### 1.6 Sensitive content found in `legacy/`

The repo is private, has one collaborator and no forks. Gitleaks did not flag these:

- **The default TAK client-certificate password**, in `SERVICES.md` and `MACHINES.md`. It is committed in **21 files across the repo** (TAK portal code, `tak-server/**`, the onetak and takml manifests, and others) and in Hermes' live `MEMORY.md`. Tracked as a separate task (section 5).
- **A truncated Notion token** in `INTEGRATION.md`. Only the first and last few characters are shown, so it cannot be used to authenticate.
- A droplet public IP, the home Home Assistant LAN address, the TAK load-balancer IP, SSH usernames and key paths, and the location of the Home Assistant token file.
- **Hermes' live `memories/MEMORY.md` contains a plaintext WordPress application password** (for the retired `zeusdelacruz.com` WordPress droplet). Memory is loaded into the model prompt every session, so it has been sent to the model provider, and it is on the volume and in the backups. The value is not recorded here. Tracked as a separate task: revoke it, then remove the entry (a live change, done with a fresh backup). The earlier memory scan missed it because the pattern was `app_pass=...`.
- Separately, **the Hermes pod holds a GitHub OAuth token for `Zeus-Delacruz`** whose scopes include `repo`, `workflow`, `write:packages`, `delete_repo`, `admin:org`, `admin:gpg_key` and `admin:ssh_signing_key`. It is on the backed-up volume. Tracked as a separate task (section 5).

---

## 2. Target design

Three layers. Each covers a different failure.

| Layer | Covers | Mechanism | Status |
|-------|--------|-----------|--------|
| A. Restore | Lost volume or cluster, including runtime state (`state.db`, sessions, memory) | VolSync restic backups, plus a tested restore | Exists, untested |
| B. Identity as code | Rebuilding Hermes' soul, skills, config and plugins from Git | Files in `platform/hermes/agent/`, seeded into an empty volume | Missing |
| C. Capture | Hermes edits its own skills, soul and memory, so Git must follow the pod | Scheduled export of an allow-listed set of files into a reviewable PR | Missing |

### 2.1 One canonical home in Git

`platform/hermes/agent/` is already Argo-managed (`platform-legacy`). Its `kustomization.yaml` lists resources explicitly (`namespace`, `pvc`, `deployment`, `dashboard`), so adding files there does not change what is deployed. Identity stays in this repo (D3). Proposed layout:

```
platform/hermes/agent/
  agents/                  one folder per agent: SOUL.md (+ CONFIG.md where it exists)
    legacy/ legacy-edge/ so/ hp/ le/ ge/ _template/
    planned/               Forge and the nine servants, with a README on the hierarchy
  profiles/                per-profile config templates (secret-free)
  config/config.template.yaml
  skills/                  custom skills only (not the upstream bundle)
  plugins/nats-platform/   plugin source
  cron/jobs.json
  AGENTS.yaml              one verified registry, status per agent
  doctrine/                CREED, PROFILES, EVOLUTION, PILLARS
  RUNBOOK.md               distilled from RUNBOOKS, INTEGRATION and FABRIC pitfalls
  FABRIC-REFERENCE.md      stack facts split out of the soul, kept current from docs/STACK.md
```

### 2.2 What is captured, and what is not

| Captured into Git | Never captured into the infra repo |
|-------------------|------------------------------------|
| `SOUL.md` (root and each profile) | `auth.json`, `.env`, tokens, keys, certificates |
| `config.yaml` after redaction (secret values replaced by references) | `state.db`, `state_compact.db`, `sessions/` (runtime state, covered by Layer A) |
| Custom skills | `memories/USER.md` and `MEMORY.md` (personal and business detail; see D2) |
| Plugin source, `cron/jobs.json` | `repos/`, `home/`, venvs, caches, loose research files |

Capture uses an allow-list of named paths, not an exclude-list, so new files on the volume are never swept in by accident.

### 2.3 Why the export is a workflow, not a sidecar

- **Anything on `main` is deployed** (ArgoCD syncs it), and GitHub will not enforce branch protection on this private repo. A write key in the pod would therefore be a path from the pod to the cluster. The pod gets no write credential.
- The repo already has the pattern: `onemind-inventory-refresh.yml` runs on the in-cluster self-hosted runner with `contents: write` and `pull-requests: write` and opens a PR with `peter-evans/create-pull-request`.
- The runner currently has no cluster permissions, so the new workflow needs a narrow read-only role. `pods/exec` cannot be limited to one command, so the workflow uses a fixed copy script from the repo, and the role is scoped to the `hermes` namespace.
- PRs created with `GITHUB_TOKEN` do not trigger the normal CI checks, so the workflow runs gitleaks itself and refuses to open a PR if anything is found.

---

## 3. The steps

Each step is its own PR (or small set of PRs), reviewed before the next begins. Step 2d (the Mac job) is done last.

### Step 0: Rescue and scrub (urgent, small)

**Why first:** the overseer definitions exist only on a pod volume, and sensitive text is sitting in Git.

1. Add the four overseer SOULs and the pod's `AGENTS.yaml` from the pod clone, unchanged. They go into `legacy/` for now so authorship and content are preserved; Step 1 moves them. The files were already copied off the pod and verified identical to the pod's working tree (kept in the session scratchpad; the pod's volume also holds them and VolSync backs it up nightly).
2. Add a one-time snapshot of the live `SOUL.md` files (root plus the four profiles) **verbatim**, so the original wording, including "daughter", is in Git history before any edit. Config and memory are not included here.
3. Remove the sensitive content in 1.6 from `legacy/`: delete the credential lines, replace addresses and SSH details with names. **No Git history rewrite** (D1).

**Done when:** `gitleaks` passes, a grep for the listed items in `legacy/` returns nothing, and the overseer files are byte-identical to the pod copies.
**Risk:** low, docs only.
**Rollback:** revert the PR.

### Step 1: Consolidate into `platform/hermes/agent/` and archive the rest

1a. Move the keep-list items from 1.5 into the layout in 2.1.
1b. **Rebuild `AGENTS.yaml`** (D5): start from the pod's version, give every agent a `status` (`running`: Legacy on the pod and Legacy-Edge on the Mac; `configured`: SO, HP, LE, GE; `planned`: Forge and the nine servants), correct the image and model fields to live values, and delete the comment about the image baking in the soul.
1c. **Edit the soul** (D4): start from the verbatim live file, remove the "daughter" framing in favour of neutral edge-node wording (the seed copy already has this wording), and move the stack facts (the Bedrock line, ports, accounts, live services) into `FABRIC-REFERENCE.md`. Remove "daughter" from `CREED.md`, the registry and the docs too. This changes the soul in Git only; the running agent is not touched (see Step 3).
1d. **Write `PILLARS.md`** (D7): the four pillars (SO 00 to 24, HP 25 to 49, LE 50 to 74, GE 75 to 99), what each owns, and the routing rules (a domain agent acts only inside its range, routes anything else to Legacy, no direct cross-domain contact).
1e. Add `agents/planned/README.md` (D6) with the hierarchy (Zeus, Legacy, overseers, servants) and a note that none of these agents runs and that several CONFIG model fields are stale.
1f. Update every reference to the old paths: `README.md`, `CONTRIBUTING.md`, `renovate.json`, `CLAUDE.md`, `docs/REPO-LAYOUT.md`, `docs/hermes-omos-ui.md`, and the `platform/hermes/agent/` docs that point at `platform/legacy/...`.
1g. Move the archive items from 1.5 (including the 14 May folders) to `~/OneMind/omos-legacy-archive/` in the vault, verified identical first. Git history keeps everything.
1h. Delete `legacy/`.

**Done when:** `legacy/` is gone, CI passes, Argo shows no change for `platform-legacy` (the kustomize resource list is untouched), and no document points at a path that no longer exists.
**Risk:** low to medium (links, and the registry rebuild).
**Rollback:** revert; the archive copy and Git history both remain.

### Step 2: Capture and seed

2a. **Export workflow (Layer C), D3.** Rules for skills and config.yaml are in `hermes-export-rules.md`; the script, workflow (manual and off) and read-only role are built (section 6 there), and a dry run against the live pod refused as designed. **Prerequisite (D11):** the credentials found in the live config and skills are rotated and moved to env or `op://` references first, or the first run is refused. `hermes-identity-export.yml`, modeled on `onemind-inventory-refresh.yml`. It runs nightly at 02:00 (after the 01:00 backup) on the self-hosted runner, copies the allow-listed paths from the pod with a fixed script, redacts `config.yaml`, runs gitleaks, and opens a PR limited to the identity paths. Needs a read-only role for the runner in the `hermes` namespace.

2b. **Seed (Layer B).** On an empty volume only, an init container copies SOULs, the config template, custom skills, plugins and `cron/jobs.json` into `/opt/data`. CI builds a small image from the repo files so the content is pinned by digest (ConfigMaps would hit the 1 MiB limit). It never overwrites an existing `SOUL.md`; a deliberate `HERMES_RESEED` flag forces a reset. This replaces the dead `01-seed-legacy` path.

2c. Apply the same capture to OpenClaw's five agents (their workspace `SOUL.md`/`IDENTITY.md` files); `platform/openclaw/` is already in Git.

**Done when:** a nightly PR appears when a skill or SOUL changes, an empty-volume start comes up with the Git-defined identity, and neither path ever touches secrets or databases.
**Risk:** medium. It adds a runner permission and an init container to the Hermes pod. Roll out the export first, seed second.
**Rollback:** remove the workflow and role, and the init container; the volume is untouched.

### Step 3: Prove it, and apply the soul change

1. Restore the latest VolSync snapshot into a **scratch namespace** (never over the live volume). Check `SOUL.md`, `config.yaml`, skills, memory files, and `PRAGMA integrity_check` on both databases. A file-level restore was already verified on 2026-10-08 (see 1.2); this step adds the application-level check.
2. Boot a second Hermes against the restored volume with the NATS platform disabled, so it cannot collide with the live agent's bus identity. Confirm it loads the soul and skills.
3. Repeat the seed path from an empty volume and compare with the live agent.
4. Record the result in `docs/operations/` and in `docs/STACK.md`. Re-run quarterly.
5. **Apply the edited soul to the live agent** (Step 1c) and correct the stale line in Legacy's `MEMORY.md` about the old disaster-recovery process. Both are live changes: do them deliberately, after the restore drill, with a fresh backup first.
6. Optional: stop backing up the clutter (two repo clones, venvs, caches) to shrink the restic repository.

**Done when:** a written restore drill exists with measured time and result, and the live soul and memory match Git.
**Risk:** low if confined to a scratch namespace; the live edits in item 5 are small and reversible from the backup.

### Step 2d (last): Mac jobs

A weekly launchd job on the Mac (D2 and D8):

- copies `memories/USER.md` and `MEMORY.md` from the pod into the vault (`~/OneMind/05 Hermes (The Runtime)/`, which auto-commits to the private vault repo). VolSync stays the authority; this is the readable copy;
- captures the Mac edge's allow-listed identity (`SOUL.md`, redacted `config.yaml`, custom skills, cron jobs) and opens a PR to this repo with the Mac's existing Git credentials. Time Machine stays the restore layer for `~/.hermes`.

Open item: a rule to tell custom skills from the bundled upstream ones, so the export does not commit all 147 (or 201) skills. To be worked out before 2a as well.

---

## 4. Decisions (recorded 2026-10-16)

| # | Decision | Outcome |
|---|----------|---------|
| D1 | Exposed credentials: rotate, or rewrite history? | **No history rewrite.** Remove the lines in Step 0. The Notion token is truncated and unusable. Replacing the default TAK cert password across all 21 files is a separate task |
| D2 | Where does personal memory live? | **Not in the infra repo.** VolSync is authoritative; operational lines move into the runbook and skills; a weekly Mac job copies the memory files into the vault |
| D3 | Credential for the export | **No credential on the pod.** A scheduled workflow on the in-cluster runner, in this repo, opens the PRs (section 2.3). The over-privileged GitHub token on the pod is replaced as a separate task |
| D4 | Which `SOUL.md` is canonical? | **The live file**, captured verbatim in Step 0, then edited: stale stack facts split out, and the "daughter" framing removed in favour of edge-node wording. The three older copies go to the archive |
| D5 | How to reconcile the three `AGENTS.yaml` files? | **One registry**, rebuilt from the pod's version and verified against live, with a `status` per agent (`running`, `configured`, `planned`) |
| D6 | The nine planned servants | **Keep as `agents/planned/`** with a README on the hierarchy, stale model fields flagged, status `planned` |
| D7 | The 14 May-2026 agent folders | **Archive to the vault**, after writing `PILLARS.md` with the pillars, domain ranges and routing rules |
| D8 | Legacy-Edge on the Mac | **Included, done last:** the weekly Mac job also captures the edge identity as a PR. Time Machine stays the restore layer |

## 5. Tracked outside this plan

- Replace the default TAK cert password in all 21 files and in the TAK server (Todoist: "Replace default TAK cert passwords across omos repo and TAK server").
- Replace the over-privileged GitHub token on the Hermes pod with a narrow one (Todoist: "Replace the over-privileged GitHub token on the Hermes pod with a narrow one").
- Revoke the WordPress application password in Hermes' memory, then remove that memory entry (Todoist: "Revoke the WordPress app password stored in plaintext in Hermes' memory, then remove that memory entry").
- Forge: decide whether it is a planned agent or is dropped. Its profile is empty today.

## 6. Not verified

- I did not read every servant SOUL in full, nor the plugin code line by line. I have not diffed the live `nats-platform` plugin against `legacy/seed/plugins/nats-platform`.
- kagent's four agents: I found no manifests for them in `platform/`; where they are defined is unconfirmed.
- Whether `memory-wiki`, `nats-context` and `platforms` on the volume come from upstream, from a separate repo, or were hand-written.
- Why the overseers and the Mac edge agent are not visible in the cluster's agent listing even though their config enables NATS.
- The NATS account layout named in the live soul (`SYS`, `CORE`, `FABRIC`, `INTERFACE`).
- The Home Assistant, S24 and Dell integrations described in the old docs were not examined.

# Hermes recoverability plan

**Status:** proposal for review. Nothing in this plan has been executed.
**Written:** 2026-10-16. Live state below was read from the cluster (`do-nyc3-omos-cloud`) on 2026-10-15/16, read-only.
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
| Declarative memory | `memories/MEMORY.md`, `USER.md` | small | No (personal) |
| Conversation and memory DB | `state.db`, `state_compact.db`, `sessions/` | 1.4 GB + 0.9 GB | No (runtime state) |
| Scheduled jobs | `cron/jobs.json` | 10 jobs, all disabled | No |
| Plugins | `plugins/` (`nats-platform`, `nats-context`, `memory-wiki`, `onemind-vision`, `platforms`) | | Only `onemind-vision` (identical to `plugins/onemind-vision/`) |
| Profiles | `profiles/` | 14 folders: 4 with real state (`sovereign-operations`, `holistic-performance`, `legacy-evolution`, `generational-empire`), 10 empty shells | Partly |
| Clutter | `repos/` 2.6 GB, `home/` 1.2 GB, two clones of this repo about 2 GB, venvs, caches, `OpenClaw/` 906 MB, `takdocs/`, loose research files | about 15 GiB | n/a |

### 1.2 What backs it up

VolSync `ReplicationSource volume-c96e219d2314` backs up `hermes-data` with restic every night at 01:00 (retain 7 daily, 4 weekly, 3 monthly, 1 yearly; CSI snapshot copy method). This works for disaster recovery. Its limits:

- Restore is all-or-nothing, and about 15 GiB of the volume is clutter.
- No Hermes restore test is recorded in `docs/operations/`.

### 1.3 Why Git cannot rebuild Hermes today

- **The seeding mechanism is dead.** `legacy/seed/01-seed-legacy` copied the persona from `/opt/hermes/seed/`. That path only existed in the old custom image. The live pod runs the stock image and `/opt/hermes/seed` does not exist, so the live `SOUL.md` has been edited by hand on the volume. The three Git copies drifted: `legacy/seed/SOUL.mother.md` (38 differing lines from live), `platform/hermes/agent/agents/legacy/SOUL.md` (80), `legacy/agents/legacy/SOUL.md` (125). These counts measure text difference, not meaning.
- **The documented recovery is weak.** `legacy/RUNBOOKS.md` section 10 says that if the volume is lost, "Hermes memory/skills are gone".
- **Some definitions exist only on the pod.** The pod keeps a clone of `onemindos/omos-legacy`. That GitHub repo no longer exists. The clone has **2 unpushed commits from 2026-09-17** that add the four overseer agents (SO, HP, LE, GE): `agents/{SO,HP,LE,GE}/SOUL.md` and 81 lines of `AGENTS.yaml`. Apart from those and one `shellcheck` line in `seed/launch.sh`, the clone is identical to `legacy/`. The clone's own git history is damaged (a `git bundle` of it cannot be cloned because of a missing object), so it must not be trusted as a backup.

### 1.4 Other agent runtimes

| Runtime | Agents | State | In Git |
|---------|--------|-------|--------|
| Hermes (pod) | Legacy plus profiles | `hermes-data`, backed up | Manifests only |
| OpenClaw (`openclaw` ns) | `legacy`, `so`, `hp`, `le`, `ge` | `openclaw-home-pvc` (50 GiB), backed up | Manifests only (`platform/openclaw/`) |
| Legacy-Edge (Zeus's Mac) | the "daughter" | `~/.hermes`, launchd | No |
| kagent | `k8s-agent`, `promql-agent`, `observability-agent`, `cilium-debug-agent` | cluster objects | No Agent manifests found in `platform/` |

### 1.5 The `legacy/` directory (146 files)

| Part | Disposition |
|------|-------------|
| `CREED.md`, `PROFILES.md`, `EVOLUTION.md` | **Keep.** Doctrine: hierarchy law, profile architecture, self-improvement loop and PR conventions |
| `INTEGRATION.md`, `RUNBOOKS.md`, pitfalls in `FABRIC.md` | **Keep, distilled** into one Hermes runbook (startup order, dispatch routing, edge delegation, NATS leaf pitfalls). Strip addresses and credentials |
| `AGENTS.yaml` | **Keep, reconcile.** Three versions exist (332 lines in `legacy/`, 80 in `platform/hermes/agent/`, and the pod's with the overseers) |
| SOULs for Legacy, Legacy-Edge, Forge, SO, HP, LE, GE | **Keep** |
| SOULs and CONFIGs for the nine planned servants (Eden, Empire, Grid, Guardian, Haven, Heritage, Oracle, Spartan, Trinity) | **Keep as "planned"**: their profiles are empty and their cron jobs are paused |
| `seed/plugins/nats-platform` (29 files) | **Keep** after comparing with the live plugin |
| `seed/Dockerfile`, `launch.sh`, `01-seed-legacy`, `config.hub.yaml`, `hermes_nats_host.py` | **Archive.** They belong to the retired custom-image path |
| `ARCHITECTURE.md`, `SERVICES.md`, `MACHINES.md`, `FABRIC.md`, `services/*.md`, `apps/`, `nats/`, `cognitive/` | **Archive.** July/August snapshots superseded by `docs/STACK.md`. They still describe Bedrock Sonnet 4.6, TAK 5.7, NFS storage, a `cloudtak` namespace and `rad.*` streams |
| 14 agent folders in the May 2026 OpenClaw/Codex format (`zenith`, `unity`, `eternal`, `sentinel`, `sovereignty`, `apex`, `creed`, `eagleeye`, `genesis`, `helix`, `librarian`, `nomad`, `overwatch`, `pulse`) | **Archive.** OpenClaw's agents are now `ge`, `hp`, `le`, `legacy`, `so` |
| `docs/tak-3d-alternatives.md`, `deploy/takserver/tak-webview-cesium.yaml` (not running), `config/wintak-omos-dell.xml`, `scripts/optimize-dell-win10.ps1`, `scripts/activate-servant.sh` | **Archive** (the script is useful only if servants are activated) |

### 1.6 Sensitive content found in `legacy/`

The repo is private, but these should not be in Git and gitleaks did not flag them:

- The default TAK client-certificate password, in `SERVICES.md` and `MACHINES.md`.
- A truncated Notion token in `INTEGRATION.md`. Truncated, but treat it as exposed.
- A droplet public IP, the home Home Assistant LAN address, the TAK load-balancer IP, SSH usernames and key paths, and the location of the Home Assistant token file.

---

## 2. Target design

Three layers. Each covers a different failure.

| Layer | Covers | Mechanism | Status |
|-------|--------|-----------|--------|
| A. Restore | Lost volume or cluster, including runtime state (`state.db`, sessions, memory) | VolSync restic backups, plus a tested restore | Exists, untested |
| B. Identity as code | Rebuilding Hermes' soul, skills, config and plugins from Git | Files in `platform/hermes/agent/`, seeded into an empty volume | Missing |
| C. Capture | Hermes edits its own skills, soul and memory, so Git must follow the pod | Nightly export of an allow-listed set of files into a reviewable PR | Missing |

### 2.1 One canonical home in Git

`platform/hermes/agent/` is already Argo-managed (`platform-legacy`) and already holds `AGENTS.yaml`, `CREED.md`, and the persona folders. Its `kustomization.yaml` lists resources explicitly (`namespace`, `pvc`, `deployment`, `dashboard`), so adding files there does not change what is deployed. Proposed layout:

```
platform/hermes/agent/
  agents/                  one folder per agent: SOUL.md (+ CONFIG.md where it exists)
    legacy/ legacy-edge/ forge/ so/ hp/ le/ ge/ _template/
    planned/               the nine servants, marked planned
  profiles/                per-profile config templates (secret-free)
  config/config.template.yaml
  skills/                  custom skills only (not the upstream bundle)
  plugins/nats-platform/   plugin source
  cron/jobs.json
  AGENTS.yaml             single reconciled registry
  doctrine/                CREED, PROFILES, EVOLUTION
  RUNBOOK.md               distilled from RUNBOOKS, INTEGRATION and FABRIC pitfalls
```

### 2.2 What is captured, and what is not

| Captured into Git | Never captured into the infra repo |
|-------------------|------------------------------------|
| `SOUL.md` (root and each profile) | `auth.json`, `.env`, tokens, keys, certificates |
| `config.yaml` after redaction (secret values replaced by references) | `state.db`, `state_compact.db`, `sessions/` (runtime state, covered by Layer A) |
| Custom skills | `memories/USER.md` and `MEMORY.md` (personal; see decision D2) |
| Plugin source, `cron/jobs.json` | `repos/`, `home/`, venvs, caches, loose research files |

Capture uses an allow-list of named paths, not an exclude-list, so new files on the volume are never swept in by accident.

---

## 3. The steps

Each step is its own PR (or small set of PRs), reviewed before the next begins.

### Step 0: Rescue and scrub (urgent, small)

**Why first:** the overseer definitions exist only on a pod volume, and sensitive text is sitting in Git.

1. Add the four overseer SOULs and the updated `AGENTS.yaml` from the pod clone, unchanged. They go into `legacy/` for now so authorship and content are preserved; Step 1 moves them. The files were already copied off the pod and verified identical to the pod's working tree (kept in the session scratchpad; the pod's own volume also holds them and VolSync backs it up nightly).
2. Add a one-time snapshot of the live `SOUL.md` files (root plus the four profiles), so the live wording is in Git even before reconciliation. Config and memory are not included here.
3. Remove the sensitive content listed in 1.6 from `legacy/`: delete the credential lines, replace addresses and SSH details with names.
4. **Zeus:** rotate the Notion token and change the TAK certificate password if it is still the default. Removing text from the current tree does not remove it from Git history.

**Done when:** `gitleaks` passes, a grep for the listed items returns nothing, and the overseer files are byte-identical to the pod copies.
**Risk:** low, docs only.
**Rollback:** revert the PR.

### Step 1: Consolidate into `platform/hermes/agent/` and archive the rest

1a. Move the keep-list items from 1.5 into the layout in 2.1. Reconcile `AGENTS.yaml` into one file (decision D5) and choose the canonical Legacy `SOUL.md` (decision D4).
1b. Update every reference to the old paths: `README.md`, `CONTRIBUTING.md`, `renovate.json`, `CLAUDE.md`, `docs/REPO-LAYOUT.md`, `docs/hermes-omos-ui.md`, and the `platform/hermes/agent/` docs that point at `platform/legacy/...`.
1c. Move the archive items from 1.5 to `~/OneMind/omos-legacy-archive/` in the vault, verified identical first, as with the other moved folders. Git history keeps everything.
1d. Delete `legacy/`.

**Done when:** `legacy/` is gone, CI passes, Argo shows no change for `platform-legacy` (the kustomize resource list is untouched), and no document points at a path that no longer exists.
**Risk:** low to medium (links, and the `AGENTS.yaml` merge).
**Rollback:** revert; the archive copy and Git history both remain.

### Step 2: Capture and seed

2a. **Export (Layer C).** A small sidecar container in the `hermes` pod mounts `hermes-data` read-only and runs nightly at 02:00 (after the 01:00 backup). It copies the allow-listed paths, passes `config.yaml` through a redaction filter, and pushes a branch and opens a PR. Hermes only changes through review, matching the "OBSERVE, LEARN, BUILD, PR, APPROVE" loop in `EVOLUTION.md`. A sidecar is used because the volume is single-writer (`ReadWriteOnce`), so a separate CronJob cannot reliably mount it. The write credential is decision D3.

2b. **Seed (Layer B).** On an empty volume only, an init container copies SOULs, the config template, custom skills, plugins and `cron/jobs.json` into `/opt/data`. CI builds a small image from the repo files so the content is pinned by digest (ConfigMaps would hit the 1 MiB limit). It never overwrites an existing `SOUL.md`; a deliberate `HERMES_RESEED` flag forces a reset. This replaces the dead `01-seed-legacy` path.

2c. Apply the same capture to OpenClaw's five agents (their workspace `SOUL.md`/`IDENTITY.md` files); `platform/openclaw/` is already in Git.

**Done when:** a nightly PR appears when a skill or SOUL changes, an empty-volume start comes up with the Git-defined identity, and neither path ever touches secrets or databases.
**Risk:** medium. It adds a container to the Hermes pod and a push credential. Roll out to the export first, seed second.
**Rollback:** remove the sidecar and init container; the volume is untouched.

### Step 3: Prove it

1. Restore the latest VolSync snapshot into a **scratch namespace** (never over the live volume). Check `SOUL.md`, `config.yaml`, skills, memory files, and `PRAGMA integrity_check` on both databases.
2. Boot a second Hermes against the restored volume with the NATS platform disabled, so it cannot collide with the live agent's bus identity. Confirm it loads the soul and skills.
3. Repeat the seed path from an empty volume and compare with the live agent.
4. Record the result in `docs/operations/` and in `docs/STACK.md`. Re-run quarterly.
5. Optional: stop backing up the clutter (two repo clones, venvs, caches) to shrink the restic repository.

**Done when:** a written restore drill exists with measured time and result.
**Risk:** low if confined to a scratch namespace.

---

## 4. Decisions needed from Zeus

| # | Decision | Recommendation |
|---|----------|----------------|
| D1 | Rotate the exposed Notion token and TAK cert password; rotate rather than rewrite Git history? | Rotate. The repo is private; rewriting history is disruptive |
| D2 | Where does personal memory (`USER.md`, `MEMORY.md`) live? | In the vault (private, auto-backed-up), not the infra repo. VolSync already covers it |
| D3 | Write credential for the export PRs | A repo-scoped deploy key in a SealedSecret, used by the sidecar. Alternative: Hermes' own GitHub login, which ties recovery to the agent being healthy |
| D4 | Which `SOUL.md` is canonical for Legacy? | The live one (68 lines), since it is what is running. Keep the three older copies in the archive for comparison |
| D5 | How to reconcile the three `AGENTS.yaml` files? | Start from the pod's version (it has the overseers), then bring over entries from the 332-line file that are still true |
| D6 | The nine planned servants | Keep as `agents/planned/` so the intent is not lost, clearly marked not running |
| D7 | The 14 May-2026 OpenClaw/Codex agent folders | Archive to the vault |
| D8 | Legacy-Edge on the Mac (`~/.hermes`) | Separate follow-up; out of scope for this plan |

## 5. Not verified, out of scope

- I did not read every servant SOUL in full, nor the plugin code line by line. I have not diffed the live `nats-platform` plugin against `legacy/seed/plugins/nats-platform`.
- kagent's four agents: I found no manifests for them in `platform/`; where they are defined is unconfirmed.
- Whether `memory-wiki`, `nats-context` and `platforms` on the volume come from upstream, from a separate repo, or were hand-written.
- Legacy-Edge on the Mac, and the Home Assistant, S24 and Dell integrations described in the old docs, were not examined.

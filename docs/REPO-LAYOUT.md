# Repository layout

What each top-level path in `onemindos/omos` is for, whether it is live, and what to consolidate. Status was judged from what CI builds, what Argo CD syncs, git history, and the live cluster on 2026-10-09. Nothing here has been moved or deleted; the consolidation section is a proposal for separate pull requests.

**Status key:** **Active** = deployed, built by CI, or edited regularly. **Reference** = read by people or agents, not deployed. **Legacy** = superseded or duplicated elsewhere. **Future** = planned, not deployed. **Broken** = does not work as committed.

## Top-level paths

| Path | Purpose | Status | Evidence / notes |
|------|---------|--------|------------------|
| `platform/` | Everything Argo CD syncs to the cluster: Application CRs (`argocd/`), per-layer manifests, Helm values (`values/`), ecosystem charts and values (`ecosystem/`), backup jobs (`backups/`), NATS config and 39 stream CRs (`nats/`) | **Active** | 552 files, 465 commits. 25 `platform-*` Applications point into it. Sub-directories with **no Application** are listed below |
| `platform/nats-tools/` | Manifests for five NATS web UIs | **Active but unmanaged** | Pods run, but no Argo Application applies this path (see `docs/STACK.md`, "Not tracked by Argo or Helm") |
| `platform/matrix/` | Matrix/Synapse and Element manifests, deploy script, bridge scripts | **Active but unmanaged** | Pods run in `matrix`; no Application. Mixes manifests with scripts, config and chatops notes |
| `platform/geo/` | Namespace and sealed secrets for a `geo` namespace | **Legacy** | The `geo` namespace no longer exists live; the geo stack runs in `onetak` |
| `platform/install.sh` | Imperative bootstrap that runs `helm upgrade --install` for 10 of the 11 Helm-CLI releases (everything except `takserver`) | **Active, risky** | Creates the GitOps gap in `docs/STACK.md`. Keep as a documented cold-start path until the releases are adopted by Argo |
| `platform/takserver/` | TAK Server deployment: Helm chart (`chart/`, with vendored subcharts) and `values-onemind.yaml` (applied by Helm CLI), the Argo-managed overlays `deploy/` (`platform-takserver`) and `fedhub/` (`platform-fedhub`), and the README and runbook. Moved from the top-level `tak-server/` on 2026-10-17; image builds and the plugin are in `services/takserver/`, the edge TAK deployment in `infra/edge/takserver/` | **Active** | `deploy/` and `fedhub/` are Argo-managed; `chart/` is applied by Helm CLI until the `takserver` adoption |
| `platform/security/` | Zero-trust layer: Kyverno, Falco, SPIRE (+ CRDs) and OpenFGA Applications plus their `onemind.io/layer: security` namespaces, applied by `platform-security` (added 2026-10-17; moved from the ecosystem foundation and applications groups) | **Active** | Each tool keeps its own namespace; OpenFGA's own namespace waits on the OneMind Core decision |
| `infra/` | Terraform for the DO VPC, DOKS cluster and project; `ecosystem-workers/` is a second root for the ecosystem node pool. `edge/mac-leaf/` holds the Mac NATS leaf config and launchd plists (moved from the former top-level `edge/` on 2026-10-17; a reference copy, not deployed) | **Active** | `terraform fmt` and `validate` run in CI. Defaults differ from the live cluster, so do not apply without a clean plan. Only hardware-bound bridges belong in `edge/`; it will be replaced or joined by a home-server equivalent |
| `services/` | All first-party source code that CI builds into images. `tak-portal` (built by `build-tak-portal.yml`, serves `onemindos.dev` and `enroll.onemindos.dev`) and `onemind-portal` (a Backstage app, built by `build-onemind-portal.yml`) moved here from the former `apps/` folder on 2026-10-17. `schemas/` holds the versioned event/command contracts and code generators (moved from the top level on 2026-10-17); `schemas/envelope/` and `schemas/missions/` trigger `build-onemind-core.yml`, and they are the source of truth for bus contracts alongside `docs/NATS-DOCTRINE.md`. `onemind-core` and `onemind-postgres` are built by `build-onemind-core.yml`. The rest (`ctl.py`, `om_*.py`, `subscriber*.py`, `bridge/`, `nex/`, `omd/`, launch scripts) are Python/Go helpers from the Mac-era fabric | **Mixed** | 169 commits. Only `onemind-core` and `onemind-postgres` are wired to CI. `om_presence.py` contains a committed default token (see "Plaintext credentials" in `docs/STACK.md`) |
| `PLATFORM-ROADMAP.md` | Living platform plan: four layers and AI trust model, namespace-to-layer map, OpenChoreo migration phases with a PR tracker (added 2026-10-18) | **Active** | Kept at the root so any agent can find it. Update its tracker in the PR that completes roadmap work |
| `docs/` | Doctrine, architecture, runbooks, operations evidence, research | **Active** | Start at `docs/README.md`. Dated files (`cluster-audit-2026-10-07.md`, `docs/operations/*2026-10-0x*`, `infra/CLUSTER-SERVICE-INVENTORY.md`) are snapshots; `docs/STACK.md` supersedes the inventory for current state |
| `reports/` | Removed. Its one file moved to `docs/operations/cluster-review-2026-10-09.md` | **Moved** | Dated reports now live in `docs/operations/` |
| `scripts/` | CI helpers (`validate_changed_yaml.py`, `validate_manifests.py`, `validate_helm_renders.py`, `validate_ecosystem_charts.py`) and ecosystem helpers | **Active** | `check_secret_patterns.py` is no longer called by CI (replaced by gitleaks); safe to remove in a cleanup PR |
| `plugins/` | Removed. Its one plugin, `onemind-vision`, moved to `platform/hermes/agent/plugins/onemind-vision/` | **Moved** | It is live on the Hermes pod (byte-identical), installed by hand; see `platform/hermes/agent/plugins/README.md` |
| `blueprints/` | Moved out of the repo to the local vault (`~/OneMind/blueprints-html`). History is retained in git | **Moved** | 28 files of HTML planning material plus a model-management script |
| `legacy/` | Removed. Archived to the local vault (`~/OneMind/omos-legacy-archive`); the agent's souls, registry, doctrine, runbook and fabric reference now live in `platform/hermes/agent/`. History is retained in git | **Moved** | The migrated copy of the old `omos-legacy` repo (146 files); that repo no longer exists on GitHub |
| `omos-infra` | Removed. Was a dangling submodule pointer (mode 160000, commit `81dc734e`) with no `.gitmodules` | **Removed** | Old name of this repo; workflow comments may still say "omos-infra main" |
| `OneMind-Platform/` | Moved out of the repo to the local vault (`~/OneMind/OneMind-Platform`). History is retained in git | **Moved** | Five capability notes (ArcCortex TAK VR, Defender AI, Leica BLK, Octave, platform) |
| `rads/` | Moved out of the repo to the local vault (`~/OneMind/rads`). History is retained in git | **Moved** | RADS (Regenerative, Autonomous, Decentralized, Self-sustainable): homestead digital-twin prototype (Python agent, Docker Compose, DimOS bridge, docs). Uses its own `rad.*` subjects, which would need porting to the locked grammar |
| `docs/references/` | Research notes on external projects (38 files; the top-level `references/` folder was merged in on 2026-10-17) | **Reference** | Not code. Contains one false-positive secret hit (`references/ai/orp.md`) |
| `docs/architecture/` | Design documents: `OneMind-OS-Architecture.md` (unified architecture), `ONEMIND-MASTER-PICTURE.md` (system synthesis), `OMOS-DATA-PICTURE.md` (data pipeline), and `decisions/`. Moved from the repo root on 2026-10-17 | **Active (under review)** | Being reviewed and updated; until then `docs/STACK.md` and `CLAUDE.md` win on live facts |
| `CLAUDE.md`, `AGENTS.md` | Instructions for AI coding agents | **Active** | `CLAUDE.md` is the operating context; `AGENTS.md` points agents at `docs/MISSION.md` |
| `.github/` | `quality.yml` (gates), `ecosystem-charts.yml`, image builds, Tailscale ACL sync, inventory refresh | **Active** | See "CI" in `CONTRIBUTING.md` and `quality.yml` header |
| `renovate.json`, `.gitleaks.toml`, `.gitleaks-baseline.json`, `.yamllint.yml`, `.pre-commit-config.yaml`, `Makefile`, `.editorconfig` | Repo-level tooling | **Active** | |

There is no `platform/future/` directory in this tree. Planned work lives in `docs/research/` and `docs/MASTER-PLAN.md` (the HTML blueprints moved to the local vault as `blueprints-html`).

## Workflows that look stale

| Workflow | Why it looks stale | Suggested action |
|----------|--------------------|------------------|
| `build-cloudtak.yml` | Built `cloudtak-api` from a separate fork; no `cloudtak` workload exists any more | **Removed 2026-10-17** (history keeps it) |
| `build-onetak.yml` | Built `onetak-{api,events,tiles}`, replaced by tak-portal and the geo stack | **Removed 2026-10-17** |
| `build-intel-services.yml` | Triggered on `platform/intel/**`, which does not exist; the intel droplet is gone | **Removed 2026-10-17** |
| `build-nats-tools.yml`, `build-takml-vision.yml`, `build-platform-backup.yml`, `build-backstage-backup.yml`, `build-fedhub-ldap.yml`, `build-tak-portal.yml`, `build-onemind-*.yml`, `build-takserver.yml` | Match live workloads | Keep |

## Proposed consolidation (documentation only, for separate PRs)

Each step is independent and reversible. Do them one PR at a time; none changes cluster behavior.

1. ~~**Remove the dangling `omos-infra` gitlink.**~~ Done.
2. **Give every live workload an owner in Git.** Add Applications (or fold into an existing one) for `platform/nats-tools` and `platform/matrix`, with `prune: false` at first. Then adopt the Helm-CLI releases in the order listed in `docs/STACK.md`.
3. ~~**Pick one home for the Legacy agent.**~~ Done: `platform/hermes/agent/` is the home, and `legacy/` was archived to the vault (see `docs/operations/hermes-recoverability-plan.md`).
4. **Retire dead directories.** After confirming nothing reads them: `platform/geo/` (namespace gone), `scripts/check_secret_patterns.py` (superseded by gitleaks).
6. ~~**Decide on `rads/`.**~~ Done: moved to the local vault.
7. **Clean `services/`.** Keep `onemind-core` and `onemind-postgres` (CI-built). Move the Mac-era scripts (`ctl.py`, `om_*.py`, `subscriber*.py`, `launch*.sh`, `create_*.py`) into `services/legacy-scripts/` or the edge repo, and rotate the committed token in `om_presence.py`.
8. **Archive stale workflows** (table above) once CloudTAK/OneTAK retirement is confirmed.
9. **Split `platform/matrix/`**: manifests stay, scripts and chatops notes go to `docs/` or `scripts/`.

## Rules for new files

- Anything that Argo syncs lives under `platform/<layer>/` and has (or extends) an Application in `platform/argocd/`. If no Application points at it, it is not deployed.
- Helm values live in `platform/values/` (cluster add-ons) or `platform/ecosystem/values/` (ecosystem charts). Chart versions are pinned in the Application's `targetRevision`.
- Dated reports go in `docs/operations/`, not the repo root.
- Never commit plaintext credentials, even in examples. CI fails on new findings (`.gitleaks.toml`).

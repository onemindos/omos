# F8 — Add and remove apps: lifecycle templates and stack hygiene

Work order for one implementing agent. Branch: `flagship/f8-app-lifecycle`.
Program index: [README](README.md). Depends on: F2. Absorbs draft #351.

## Outcome Zeus sees

- **Add an app**: one form (Create → "Add app to OneMind"). Pick the hosting
  type (K8s via Argo, DigitalOcean, SaaS, home/device), the pillar, the System,
  the owner, the health check, the data/backup plan and the access. It opens a
  **reviewed GitHub PR** containing the catalog entry and, for Kubernetes, the
  Argo Application and manifests skeleton. Nothing changes live until Zeus
  merges.
- **Retire an app**: Create → "Retire app". Pick the app; it shows everything
  attached (workloads, Argo app, volumes and snapshots, DNS, secrets, Gatus
  check, dashboards, alerts, cost), asks for data retention (export, keep the
  final backup, retention days), and opens a **removal PR**. The PR deletes
  manifests and the catalog entry, and lists the manual follow-ups (volume and
  snapshot deletion after retention). No deletion runs from the portal.
- **Stack hygiene report** (`/stack-hygiene`, weekly):
  - Argo apps without a catalog entry, and catalog apps without a running workload
  - workloads with no owner app
  - orphan DO volumes, old snapshots (F6)
  - apps with no health check, no backup or no runbook
  - images with no digest pin
  This shows what to add and what to remove.

## Evidence

- Templates exist: `services/onemind-portal/templates/{app-onboarding,deployment-change,recovery-preparation,service,service-registration}`.
  Draft #351 extends `app-onboarding/skeleton/README.md` with the per-app
  contract. No retirement template exists.
- 69 Argo apps, 6 mapped. 58 DO snapshots with no retention view.

## Build

1. Fold #351's contract into an **"Add app" template** with
   `publish:github:pull-request` (the existing GitHub scaffolder module) and
   per-hosting branches. Use the existing GitHub App credentials; no new
   secrets.
2. **"Retire app" template**: an `EntityPicker` plus a custom action that
   gathers the attached-resource inventory, read-only, from the catalog
   relations, the F2 mappings, F6 and the Argo app. It writes a removal PR with
   a checklist in its body. Retention defaults: keep the last backup for 30 days.
3. **Hygiene report**: a backend scheduled task (weekly, plus on demand)
   producing JSON plus a page. Each finding links to "Add" or "Retire" with
   pre-filled parameters.
4. **Dry run**: both templates support the scaffolder dry-run in the editor.

## Out of scope

Executing deletes, Argo syncs, or `kubectl` from the portal.

## Files

`services/onemind-portal/templates/{app-add,app-retire}/**` (new; keep the old
templates until Zeus says to remove them), `plugins/plugin-onemind-services-backend/**`
(retire inventory action, hygiene task), `plugins/plugin-onemind-services/**`
(hygiene page), and tests.

## Acceptance

- [ ] Dry run of "Add app" (K8s and SaaS variants) produces the expected files;
      a real run opens a draft PR on a test branch.
- [ ] "Retire app" on `stac-browser` lists its workloads, Argo app, links and
      alerts, and produces a PR that removes only that app's files (the test
      asserts the file list).
- [ ] The hygiene report flags the unmapped Argo apps and at least the snapshot
      age findings.
- [ ] No template step calls `kubectl`, Argo sync or DO write APIs (test).
      Required gates pass.

## Status

(Implementing agent: write progress here.)

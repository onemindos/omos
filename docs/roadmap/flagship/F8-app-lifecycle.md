# F8 — Add and remove apps: lifecycle templates and stack hygiene

Work order for one implementing agent. Branch: `flagship/f8-app-lifecycle`.
Program index: [README](README.md). Depends on: F2. Absorbs closed draft #351.

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
  Closed draft #351 extends `app-onboarding/skeleton/README.md` with the per-app
  contract. No retirement template exists.
- 60 Argo apps and 53 DO snapshots with no retention view (2026-11-04).
- Paperclip, Nextcloud and Hermes Legacy were retired by hand in November. The
  2026-11-03 cluster workup still found leftover volumes and stale alerts, which
  a retirement checklist should catch.

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
- [ ] The hygiene report flags any unmapped Argo app and at least the snapshot
      age findings.
- [ ] No template step calls `kubectl`, Argo sync or DO write APIs (test).
      Required gates pass.

## Status

**Built 2026-11-05 on `feat/backstage-app-lifecycle`; not deployed.** The image builds in CI on merge and
`platform/backstage/deployment.yaml` still pins the old digest, so none of this is visible until a reviewed digest update.
Design, boundaries and the file lists: [App lifecycle (F8)](../../../platform/backstage/README.md#app-lifecycle-f8).

| Item | State | Notes |
|---|---|---|
| "Add app" template (`add-app`) | Built | K8s, DigitalOcean, SaaS and home-device variants. K8s opens a PR with the Argo Application (manual sync, no finalizer), manifests, tailnet or Cloudflare-allowlisted ingress, VolSync plan, secret placeholder, STACK section, catalog inputs and the Argo count bump. Folds in #351's per-app contract as the PR checklist. |
| "Retire app" template (`retire-app`) | Built | `EntityPicker` plus the read-only `onemind:app-retire:plan` action; app-wide or single-workload. Default retention: keep the last backup 30 days. |
| Stack hygiene (`/stack-hygiene`) | Built | Weekly backend task, stored JSON, page, admin-only refresh. Orphan volumes and old snapshots are `unavailable` until F6. |
| Old templates | Kept | `app-onboarding`, `deployment-change`, `recovery-preparation`, `service`, `service-registration` are untouched until Zeus says to remove them. |

**Decisions and deviations**

- No token field and no new secret: the publish step uses the GitHub App in `integrations.github`. The README documents that App
  as read-only, so **it needs Contents and Pull requests write on `onemindos/omos`** before a real run can open a PR.
- No new RBAC. The retire inventory lists what `backstage-observer` already reads. Secrets, VolumeSnapshots, ReplicationSources and
  PV reclaim policy are not readable, and are reported as unavailable with a "check by hand" note, never as none.
- The generated catalog output (`apps.yaml`, `cluster.yaml`, both copies, runbook pages) is not hand-edited. Templates edit its inputs
  (`workload-map.yaml`, `structure.yaml`, `coverage-exceptions.yaml`, `ref-baseline.json`, the two hard-coded Argo counts) and the PR
  checklist starts with `node scripts/curate-structure.cjs`. That sequence was run in a scratch copy of the repository for a new K8s
  app and for retiring `thingsboard`: `test-catalog-structure.cjs` and `test-coverage.cjs` pass after it.
- "Retire on `stac-browser`": the live catalog has no `stac-browser` entity (it is a workload of the `stac` app), so the test adds the
  workload entity the generator would create and runs against the real `platform/stac` files.
- Retirement keeps (and lists) markdown, restic credentials, shared source directories and protected platform directories
  (`platform/argocd`, `security`, `backups`, ...) instead of deleting them.

**Acceptance**

- [x] Dry run of "Add app" (K8s tailnet+VolSync, K8s Cloudflare, SaaS) produces the expected files through the real `fetch:template`
      action, and the publish step opens a draft PR with no token (`templates.test.ts`, with a fake GitHub client). A real run against a
      test branch needs the deployed image and the App permission above. **Open.**
- [x] "Retire app" on `stac-browser` lists its Deployment, the `platform-stac` Argo application (finalizer, link), its ingress and the
      alert/dashboard search, and the PR deletes only `platform/stac/browser.yaml`, edits only `platform/stac/ingress.yaml` and adds the
      retirement record; the test asserts the file list. The NetworkPolicy that selects the browser is flagged, not edited.
- [x] The hygiene report flags every seeded fixture (unmapped Argo app, catalog app with no workload, ownerless workload, unpinned image,
      no probe, no backup, no runbook). Snapshot age and orphan volumes are `unavailable` with a reason because F6 is not built, so the
      acceptance wording "at least the snapshot age findings" is **not met** until F6.
- [x] No template step calls `kubectl`, Argo sync or a DigitalOcean API (steps and lifecycle source are scanned in
      `templates.test.ts` and `test_gitops_portal_lifecycle.py`).
- [x] `yarn install --immutable`, `tsc`, lint of the changed packages, the portal node and gitops tests, and the backend build pass,
      except the pre-existing failures below.

**Pre-existing failures (not from this change).** `repo test`: 640 pass, 8 fail, all in `plugin-onemind-inventory-backend`
(`provider`, `placement`, `entities`, `startup`, `live`; catalog drift, the same eight F9 recorded). `lint` also reports three errors in
`plugin-onemind-services` test files this change does not touch (`AppHeader.test.tsx`, `cards.test.tsx`, `links.test.ts`).
`yarn build:all` stops in the `plugin-onemind-theme` type rollup (`ThemeMode` not exported), a package this change does not touch
(not re-run on `origin/main`); `yarn build:backend` passes (its bundle contains the lifecycle code) and so does the app build (`yarn workspace app build`).

**Check after deployment** As an administrator: Create shows "Add app to OneMind" and "Retire app"; run Retire app on a workload
without ticking the confirmation box and read the inventory, then run it once for real and confirm the draft PR deletes only what the
list says. Open `/stack-hygiene`, press Refresh now (as a non-administrator it says administrators only), and confirm the two
DigitalOcean checks read "Unavailable". Confirm the GitHub App can open a PR before the first real run.

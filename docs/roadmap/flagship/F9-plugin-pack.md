# F9 — Plugin pack: verified additions

Work order for one implementing agent. Branch: `flagship/f9-plugin-pack`.
Program index: [README](README.md). One plugin per commit, so any can be
reverted alone.

## Outcome Zeus sees

- **Cost.** A working Cost page with real OpenCost data per app, namespace and
  pillar. Today the portal has only the `/api/proxy/opencost/allocation` route.
- **Security.** Kyverno policy results per app ("3 policy violations").
- **Shortcuts and bookmarks.** Personal shortcuts in the sidebar; bookmarks
  per entity.
- **Admin tools.** DevTools (installed packages, config, health) and the
  unprocessed-entities view, for admins only.

## Candidates, with reasons (verify each against Backstage 1.55.3 before adding)

| Plugin | Why it fits | Notes |
|---|---|---|
| `@backstage-community/plugin-opencost` | OpenCost already runs (`opencost` namespace) | Check its new-frontend support; otherwise wrap with core-compat. Read through the existing proxy route and reuse `services/onemind-portal/scripts/cost-allocation.cjs` for unknown/partial data. |
| Kyverno Policy Reporter (`kyverno` Policy Reporter plugin or TeraSky Kyverno) | Kyverno installed | Needs Policy Reporter UI/API in-cluster (GitOps, manual sync) |
| Shortcuts (`@backstage/plugin-shortcuts`) | Personal sidebar shortcuts | Core-maintained |
| DevTools (`@backstage/plugin-devtools`) | Admin diagnostics | Admin-only permission; no secret display |
| Catalog unprocessed entities (`@backstage/plugin-catalog-unprocessed-entities`) | Shows why entities fail to ingest | Admin-only |
| Optional: Q&A (`drodil/backstage-plugin-qeta`) | Family/household knowledge Q&A | Only if Zeus wants it |
| Optional: Trivy Operator + Policy Reporter source | Vulnerability reports per workload | Capacity check first; see the pilot in #442 |

Not recommended now: Datadog, PagerDuty, Snyk, SonarQube and other SaaS
products we do not use. The existing Grafana, Prometheus, Loki and Alertmanager
cover monitoring.

## Build

For each plugin: pin the exact version, register it in `App.tsx` and
`backend/index.ts`, add config without secrets, add an entity card or page via
the F5 composition, add any SealedSecret placeholders for Zeus, add a
permission negative test where it has backend routes, and attach a screenshot.

## Files

`services/onemind-portal/packages/{app,backend}/{package.json,src/**}`,
`yarn.lock`, `app-config.yaml`, `platform/**` (only Policy Reporter/Trivy GitOps
if approved; manual sync), and docs.

## Acceptance

- [ ] Cost page shows non-empty namespace allocation for the last 7 days that
      matches OpenCost's own API within 1%.
- [ ] Each added plugin renders on a real entity, or shows a clear
      "not configured" state.
- [ ] Admin-only plugins return 403 for non-admins (test).
- [ ] `yarn install --immutable` passes and required gates pass. Bundle-size
      change is recorded.

## Status

**Built 2026-11-05 on `feat/backstage-plugin-pack`, one commit per plugin; not deployed.** The image builds in CI on merge and
`platform/backstage/deployment.yaml` still pins the old digest, so nothing below is visible until a reviewed digest update.
Details and boundaries: [Plugin pack](../../../platform/backstage/README.md#plugin-pack-f9).

| Plugin | State | Notes |
|---|---|---|
| OpenCost Cost page | Built (internal plugin) | The community OpenCost plugin has no new-frontend export and its own API config; a small page reads the existing proxy route and validates with a TypeScript port of `scripts/cost-allocation.cjs` (parity-tested). Unknown and partial data are labelled, never zero; `__unmounted__` cost is shown separately; the page says these are estimates, not the DigitalOcean bill. |
| Shortcuts | Built | `@backstage-community/plugin-shortcuts` 0.25.1, not `@backstage/plugin-shortcuts` (deprecated since April 2024, not in the 1.55.3 manifest). Legacy plugin converted; desktop sidebar only. |
| DevTools | Built, admin only | 0.1.43 / backend 0.5.21. Info view only: `devtools.config` and the scheduler views are closed for everyone. `devtools.info` and `devtools.external-dependencies` have action `read`, so the policy now lists them as administrator-only. |
| Unprocessed entities | Built, admin only | 0.2.35 / module 0.6.16. Standalone page enabled, DevTools copy disabled. |
| Kyverno policy results | Built without a new workload | Policy Reporter needs another workload; TeraSky needs core packages newer than 1.55.3. A card reads the PolicyReports Kyverno already writes through the backend's existing read-only cluster access. Adds `get`/`list` on `policyreports.wgpolicyk8s.io` to the `backstage-observer` ClusterRole. |
| Q&A, Trivy Operator, bookmarks per entity | Not built | Q&A skipped by owner decision; Trivy Operator already runs in `trivy-system`; bookmarks were not requested. |

**Acceptance**

- [ ] Cost page matches OpenCost's own API within 1%. The page shows OpenCost's `totalCost` per namespace unchanged
      (no rescaling), and the real payload captured on 2026-11-05 is its regression fixture (`takserver` $14.70 over 7 days,
      `__unmounted__` $26.97). A live comparison needs the deployed image. **Open.**
- [x] Each plugin renders or shows a clear state: tests cover Cost (available, partial, unknown, empty, error), the Kyverno card
      (zero, violation, unavailable, stale, no report) and the DevTools and unprocessed-entities tabs in a test app. Not yet seen in
      a browser against live data. **Open: check after deployment.**
- [x] Admin-only plugins return 403 for non-admins: `adminRoutes.test.ts` and `unprocessedRoutes.test.ts` run the real routes
      behind the real policy (non-admin, service credential, and `/config` for an administrator).
- [x] `yarn install --immutable`, `tsc`, lint, the portal's node tests, the app and backend builds, and
      `test_gitops_portal_*.py` pass. `repo test`: 600 pass; 8 fail in `plugin-onemind-inventory-backend`, which fail the same way
      on `origin/main` (catalog drift, unrelated).
- [x] Bundle size: app JavaScript 16,324,215 to 16,816,248 bytes uncompressed (+492 KB, +3.0%); `dist` 89,156 to 91,324 KB.

**Check after deployment** Sign in as an administrator and as a non-administrator. Both: Insights > Cost shows namespaces with a
total and an estimate notice; an app page shows "Policy results (Kyverno)". Administrator only: Admin > DevTools (Info tab only) and
Admin > Unprocessed entities. As a non-administrator, `/api/devtools/info`, `/api/devtools/config` and
`/api/catalog/entities/unprocessed/failed` return 403.

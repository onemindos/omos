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

(Implementing agent: write progress here.)

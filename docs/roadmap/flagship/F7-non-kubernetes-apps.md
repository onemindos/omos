# F7 — Non-Kubernetes apps: model, live health and access

Work order for one implementing agent. Branch: `flagship/f7-non-k8s-apps`.
Program index: [README](README.md). Depends on: F2.

## Outcome Zeus sees

SaaS apps (Zoho, Gmail, Dropbox, Monarch, 1Password …), home and edge services,
devices and DO-hosted non-K8s services are **managed like first-class apps**:

- **Same page.** The F5 app page, with a hosting badge.
- **Live health.** Live health for anything reachable: HTTP(S) or TCP checks
  and certificate expiry, from **Gatus** (recommended addition) running
  in-cluster now and at home later.
- **Facts.** Owner, pillar, plan/cost, renewal date, data export method, backup
  location, recovery steps, and who has access (text, never credentials).
- **Views.** A "Not on Kubernetes" view grouped by hosting type.

## Evidence

- 325 Components. Many are external apps (Apple Music, 1Password, Zoho apps)
  marked `type: website, lifecycle: production` with no health signal. The home
  launcher has hosting chips (Self-hosted / Cloud service / Local / Reference /
  Planned / Unknown), but these are not a catalog model.
- `/today` shows Zoho Projects/Calendar "not connected".

## Build

1. **Annotation contract** (documented in `docs/catalog/app-contract.md`,
   validated by a test):
   - `onemind.io/hosting`: `kubernetes|digitalocean|saas|home|device|desktop|mobile`
   - `onemind.io/health-check`: the Gatus endpoint name, if any
   - `onemind.io/data-export`, `onemind.io/backup-location`,
     `onemind.io/recovery-runbook`
   - `onemind.io/renewal`, `onemind.io/plan`
   - `onemind.io/access`: who has access, as text only
   Backfill the hosting annotation for all curated apps; set the rest to
   "Unknown" visibly.
2. **Gatus**: add a GitOps app (`platform/gatus/**` plus an Argo Application with
   **manual sync**, so Zeus activates it). Endpoints come from config-as-code
   (`platform/gatus/endpoints.yaml`), generated from catalog entries that carry
   `onemind.io/health-check`. Expose Prometheus metrics. Add alert rules for
   down endpoints and for certificates expiring in under 14 days, routed through
   Alertmanager to F3. Check resource requests against current node capacity
   and record them.
3. **Status card** on the F5 page reads Gatus results (Gatus API or Prometheus).
4. **"Not on Kubernetes" view**: a saved catalog filter plus a section on the
   pillar pages.

## Out of scope

Storing app credentials, and SSO changes for SaaS apps.

## Files

`services/onemind-portal/catalog/apps.yaml` (annotations), `platform/gatus/**`,
`platform/argocd/platform-gatus.yaml` (manual sync), `platform/monitoring/**`
(rules), `plugins/plugin-onemind-services/**` (status card), docs and tests.

## Acceptance

- [ ] Every curated app has `onemind.io/hosting`. The test fails on a missing
      or invalid value.
- [ ] Gatus manifests render (`validate_helm_renders`/kubeconform) and sync is
      manual. After Zeus syncs: at least 10 endpoints green, and one deliberate
      failing test endpoint shows red on the page and fires one alert.
- [ ] SaaS app page shows owner, renewal, export and recovery fields.
- [ ] Required gates pass.

## Status

(Implementing agent: write progress here.)

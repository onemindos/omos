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

- Gatus still has no Argo Application (checked 2026-11-04).
- 325 Components (2026-10-21). Many are external apps (Apple Music, 1Password, Zoho apps)
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

2026-11-05, branch `feat/f7-non-k8s-apps` (local, not pushed or merged). Source and tests are done; **nothing is
deployed**: `platform-gatus` has manual sync and has not been synced, and the portal changes are not in an image yet.

- **Contract**: [docs/catalog/app-contract.md](../../catalog/app-contract.md). Canonical hosting values plus the existing
  ones kept as legacy values (`digitalocean-nyc3` and the rest), mapped to a kind; nothing was mass-renamed.
  `scripts/app-contract.cjs` and `test-app-contract.cjs` fail on a missing or invalid value. All 197 curated Components
  have a valid value. New annotation `onemind.io/health-url` (the URL Gatus probes); `onemind.io/health-check` is now the
  endpoint name (the F8 add template wrote a URL there; it now writes name and URL separately and uses canonical
  hosting values).
- **Backfill**: every Component already had a legacy value except Forgejo (added). Seven cluster apps that were labelled
  as external or mesh despite an Argo application were corrected to `digitalocean-nyc3` (Prometheus, Loki, TAK Server,
  OneTAK, ClickHouse, WebODM, NATS Surveyor). Result by kind: kubernetes 57, saas 97, home 4, digitalocean 4,
  desktop 11, web 6, planned 18, unknown 0. Plan, renewal, access, export and backup were **not** invented: they are absent
  and the page says "Not recorded".
- **Gatus**: `platform/gatus/**` and `platform/argocd/platform-gatus.yaml` (manual sync, wave 5). Pinned `v5.37.0` by
  digest, in-memory storage (history resets on restart), requests 10m CPU / 48Mi, limits 200m / 128Mi, hardened pod, Tailscale
  Ingress `status`, ServiceMonitor, alerts for down 5m, certificate under 14 days and Gatus not scraped. 13 generated
  endpoints (10 in-cluster, 3 SaaS) plus the deliberate failure `test-deliberate-failure`. Config and metrics were verified
  against the real image locally. Gmail and Monarch have no catalog Component, so they have no check yet.
- **Portal**: App facts card (all Components), Uptime check card (apps with `health-check`, through the new GET-only
  `/api/proxy/gatus/statuses` route), and the "Not on Kubernetes" page (`/non-kubernetes`, sidebar under Services).
  The Gatus API does not export the exact certificate expiry, so the card shows whether the 14 day condition holds; the
  exact date is the Prometheus series `gatus_results_certificate_expiration_seconds`.
- Not done: the existing combined status dot does not include Gatus yet. `backstage-cli repo build` fails in this
  environment on a rollup error in `plugin-onemind-theme` type output (untouched by this branch); tsc, lint and tests pass.
- **Owner steps after merge**: build and roll the portal image; sync `platform-gatus` in Argo CD; open
  `https://status.tail717ba0.ts.net`; confirm 10+ green and `test-deliberate-failure` red, then wait five minutes for
  `GatusEndpointDown`; then remove the deliberate failure ([how](../../../platform/gatus/README.md#remove-the-deliberate-failure)).
- Acceptance still open: green endpoints, the red page and the alert are **not verified live**.

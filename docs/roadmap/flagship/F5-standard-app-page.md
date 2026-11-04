# F5 — One standard app page for every service

Work order for one implementing agent. Branch: `flagship/f5-standard-app-page`.
Program index: [README](README.md). Depends on: F2 (app-centric entities).

## Outcome Zeus sees

Every app, whether Kubernetes, DigitalOcean, SaaS or home device, opens to the
**same layout**:

- **Header**: name, pillar badge, hosting badge (K8s / DO / SaaS / Home /
  Device), live status dot, owner, and a primary "Open app" button.
- **Overview** tab: live status summary (alerts, Argo sync and health, ready
  replicas, uptime check), key links, dependencies (graph card), data and
  backup status, cost this month.
- **Live** tab: Kubernetes resources and pods (native plugin), Argo history,
  recent events.
- **Metrics** tab: Grafana dashboard plus scoped Prometheus alerts (existing).
- **Logs** tab: prefer an existing plugin (KwirthLog or Kubelog for live pod logs);
  otherwise a Loki query through a Grafana Explore link or embed. No new log store.
- **Cost** tab: OpenCost allocation for the app's namespace and labels (F9).
- **Docs** tab: TechDocs and runbook.
- **Code** tab: repo, PRs and CI. Show "Connect GitHub" as a button there, and
  **never** an automatic pop-up.

Tabs with no data are hidden, not shown empty. Each card shows its data source
and observed time once, compactly. Do not repeat the long "this is not proof"
disclaimers on every card; put one info tooltip per card instead.

## Evidence

- `/catalog/default/component/k8s-backstage-backstage-21b99354` (the best-wired
  page today) has tabs Overview, TechDocs, APIs, Deployment, GitHub Actions, Tech
  Insights, Pull/Merge Requests, Operations. It opens an unrequested **GitHub
  OAuth dialog** and repeats disclaimer text on many cards.
- Composition modules: `packages/app/src/componentOverviewModule.tsx`,
  `catalogCompositionModules.ts`. Existing tests: `ComponentOverview.test.tsx`,
  `componentOverviewModule.test.ts`, `catalogCompositionModules.test.ts`.

## Coverage contract: every service shows everything

Zeus's requirement: **every service** shows its Argo app, Grafana dashboard,
alerts, Kubernetes resources, logs, cost and runbook. Not 6 of 69.

| Capability | Plugin | Annotation(s) each Kubernetes app gets | Source of the value |
|---|---|---|---|
| Argo CD status/history | Argo CD (installed) | `argocd/app-name` | F2 `workload-map.yaml` (all 69 apps) |
| Kubernetes resources/pods | Kubernetes (installed) | `backstage.io/kubernetes-id` or `backstage.io/kubernetes-label-selector` + `backstage.io/kubernetes-namespace` | F2/F11 |
| Grafana dashboard | Grafana (installed) | `grafana/dashboard-selector` (tag based) | **One generic "OneMind workload" dashboard** with `namespace`/`workload` variables, provisioned through GitOps in the existing Grafana folder; every app links it with its own variables. Apps with a bespoke dashboard tag it `app:<name>` instead |
| Alerts | Grafana / Roadie Prometheus (installed) | `grafana/alert-label-selector`, `prometheus.io/alert` | namespace + app labels |
| Logs | KwirthLog/Kubelog, or a Loki link | plugin's annotation, or a generated Explore URL | namespace + selector |
| Cost | OpenCost (F9) | namespace/label selector | namespace |
| Docs/runbook | TechDocs (installed) | `backstage.io/techdocs-ref`, `onemind.io/recovery-runbook` | a default runbook template per app |
| Source/CI | GitHub (installed, connect button) | `github.com/project-slug` | repo map |

Rules:
- Annotations are **generated** from the F2 mapping by the curate scripts, not
  hand-typed per app, so a new app gets full coverage automatically.
- A coverage test (extend `audit-catalog.cjs`, reusing #355's report) **fails the
  build** when a Kubernetes app lacks any required annotation, unless it is
  listed in a reviewed exceptions file with a reason.
- Non-Kubernetes apps (F7) get the subset that applies: health check,
  Grafana/Gatus dashboard, runbook and cost.

## Plugin first

Fill tabs with existing plugins wherever possible:
- Kubernetes (installed), plus Red Hat Application Topology or TeraSky
  Kubernetes resources, for the Live tab
- Argo CD (installed)
- Grafana and Roadie Prometheus (installed) for Metrics
- KwirthLog/Kubelog for Logs
- OpenCost for Cost (F9)
- TechDocs (installed) for Docs
- GitHub Actions/PRs (installed, behind a connect button) for Code

Custom code is limited to the layout, the status dot and the hosting-aware filters.

## Build

1. Define the page once with new-frontend-system `EntityContentBlueprint` /
   `EntityCardBlueprint` extensions and `filter` expressions by
   `onemind.io/hosting` and annotation presence. Do not copy a separate layout
   per kind.
2. **Status dot.** One shared hook combining alerts, Argo, workload readiness
   and uptime into `ok | warn | down | unknown`. Unknown is grey, never green.
3. **GitHub.** Stop auto-triggering OAuth. Render GitHub cards only after the
   user connects, with an explicit button.
4. **Disclaimers.** Move evidence and disclaimer text into per-card tooltips
   driven by one component.
5. Apply the same layout to System and Domain pages: a summary of member app
   statuses.

## Files

`services/onemind-portal/packages/app/src/{componentOverviewModule,catalogCompositionModules,prometheusPlugin}.*`,
`plugins/plugin-onemind-services/**`, `plugins/plugin-onemind-common/**` (shared
status hook), and their tests.

## Acceptance

- [ ] Coverage report: 69/69 Argo-managed apps have Argo, Kubernetes, Grafana dashboard, alerts, logs and runbook links, or a reviewed exception. Spot-check 5 apps signed in.
- [ ] The same tab set appears on a K8s app (backstage), a SaaS app (for example
      Zoho), and a DO resource (after F6). Tabs without data are hidden.
- [ ] Opening any app page triggers no OAuth pop-up.
- [ ] The status dot matches Live Ops (F4) for the same app and minute.
- [ ] Logs tab opens the correct Loki query for the app.
- [ ] Existing layout tests are updated, not deleted. Required gates pass.
      Screenshots attached (desktop/phone, dark/light).

## Status

Built by an implementing agent, completed and reviewed by the coordinator (2026-10-22). Source only: the catalog annotations,
runbooks and Grafana dashboard go live when Argo applies them (annotations and the dashboard need no image); the page layout,
tabs and GitHub gate need an image release.

**Done**
- **Coverage contract.** `scripts/curate-coverage.cjs` generates, from F2's reviewed workload map and never by hand: Argo app, Kubernetes
  selector and namespaces, Grafana dashboard path and alert selector, Loki query, runbook, OpenCost namespaces and source. It only
  fills what is missing and never replaces a reviewed value. `audit-catalog.cjs` now fails when a Kubernetes app lacks a required
  capability unless `catalog/coverage-exceptions.yaml` lists it with a reason (29 reviewed exceptions, each stating why; a stale
  exception also fails). Negative tests cover both directions.
- **Per-app runbooks.** 64 generated runbook pages under `runbooks/docs/apps/` (TechDocs, baked into the image).
- **One Grafana dashboard for every app.** `platform/monitoring/grafana-dashboard-onemind-workload.yaml` (uid `onemind-workload`,
  variables `namespace` and `workload`, eight panels), provisioned like the existing dashboards (sidecar label plus folder
  annotation). Datasource `prometheus` verified against the live datasource ConfigMap.
- **Standard page.** One header for every Component, System and Domain (Backstage UI `Header`: name, kind, pillar and hosting badges,
  status dot, owner, lifecycle, Open app, star, more menu), a Live status card, a Key links card, and Logs, Cost and Runbook tabs that
  appear only when the app has the data. System and Domain pages get an App status summary of member apps. OpenChoreo-managed
  entities keep their own header. The existing Kubernetes, Argo CD, Grafana, Prometheus, TechDocs and GitHub plugins remain the Live,
  Metrics, Docs and Code views.
- **One status model.** `useAppStatus` / `StatusDot` in `plugin-onemind-common`: ok only when something is connected and healthy; grey
  unknown for missing or stale data, never green. The source sits behind one adapter (`appStatusApiRef`) that reads
  `/api/onemind/status` (still served by F4's collector, field names verified) and the scoped Prometheus routes, so it can move to the
  Live Ops aggregator in one line.
- **No automatic GitHub sign-in.** Native `PluginWrapperBlueprint` modules wrap the `github-pull-requests` and `github-actions` plugins:
  their cards and tabs render only after one "Connect GitHub" click (shared across the page, remembered in the browser, works with
  storage blocked), so the OAuth dialog no longer opens on every app page.
- One info tooltip per card replaces the repeated "not proof of health" paragraphs.

**Plugin-first (rule 11)** Evaluated against Backstage 1.55.3 and kept as the per-entity views (installed): Kubernetes, Argo CD, Grafana,
Roadie Prometheus, TechDocs, GitHub Actions and Pull Requests. Not added in this PR: Red Hat Application Topology, TeraSky Kubernetes
resources, Kwirth/Kubelog and OpenCost (OpenCost belongs to F9). Custom code is limited to the generated annotations, the header
and cards, the status hook and the GitHub gate. Dependencies added: `@backstage/catalog-model`, `@backstage/core-components`,
`@backstage/ui` (all pinned to the versions already in the lockfile) and test-only `@backstage/test-utils`.

**Tests** common 37, services 30 (cards, header, GitHub gate, links), app 56 incl. the whole-app render, catalog node tests 60 on the host
and 55 pass + 5 skipped in a simulated image build context (repo-level tests use the `repoOnly` skip guard).

**Not verified** The page has been rendered only in Jest (real extensions and cards inside Backstage's test app), not yet in a browser
against live data; that check is part of the release acceptance. Coverage numbers per capability should be re-read from the audit after
release. The Cost tab links to OpenCost until F9 adds an allocation card.

**Review notes** The agent's branch had no registration in `App.tsx` and two type errors (`type: 'summary'` is not a valid card type in this
Backstage version; now `content`); both fixed. Header, cards and gate gained render tests.

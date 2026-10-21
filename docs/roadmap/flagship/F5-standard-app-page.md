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

- [ ] The same tab set appears on a K8s app (backstage), a SaaS app (for example
      Zoho), and a DO resource (after F6). Tabs without data are hidden.
- [ ] Opening any app page triggers no OAuth pop-up.
- [ ] The status dot matches Live Ops (F4) for the same app and minute.
- [ ] Logs tab opens the correct Loki query for the app.
- [ ] Existing layout tests are updated, not deleted. Required gates pass.
      Screenshots attached (desktop/phone, dark/light).

## Status

(Implementing agent: write progress here.)

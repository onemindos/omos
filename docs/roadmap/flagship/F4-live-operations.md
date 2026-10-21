# F4 — Live Operations: "what is wrong right now"

Work order for one implementing agent. Branch: `flagship/f4-live-operations`.
Program index: [README](README.md). Soft dependency: F3 (signals).

## Outcome Zeus sees

A **Live Ops** page (`/live`) that answers "is anything broken, and where" in
one screen, updating in real time:

- **Health banner.** Green, amber or red, with the reason
  ("2 critical alerts, 1 app degraded").
- **Firing alerts.** Grouped by app and pillar, with severity, age, the affected
  app, a runbook link and a Grafana link.
- **Argo CD.** Applications OutOfSync, Degraded, Progressing for more than 10
  minutes, or whose last operation failed.
- **Workloads.** Crash-looping, Pending or not-ready pods, plus recent restarts
  and OOM kills, each linked to the owning app (via F2 mapping).
- **Nodes.** Ready state, CPU, memory and disk pressure, and allocatable versus
  requested for the 4 DOKS nodes.
- **Storage.** PVCs over 80% full, and volume and snapshot growth.
- **Backups.** Last success per job (CronJobs, VolSync, Velero, CNPG). The status
  must tell **retired/superseded** jobs (for example `backstage-backup`, replaced
  by `backup-postgres-backstage`, and `backup-clickhouse`, replaced by VolSync),
  **pending setup** (`backup-offsite-copy`, waiting for credentials) and
  **failing/overdue** apart. Use an annotation such as
  `onemind.io/backup-state: active|superseded|pending-setup` on the CronJob, with
  `onemind.io/superseded-by`. Today all three show only as "suspended".
- **Uptime.** External and endpoint checks: the DO uptime check now, Gatus
  later (F7).
- **Events.** The last 30 minutes of Kubernetes Warning events, deduplicated.
- **Embedded Grafana panels.** Cluster CPU and memory, network, and the NATS
  and Postgres overview.

## Evidence

- `/today` "Needs attention" already reads Argo, workload readiness and backup
  jobs (`plugins/plugin-onemind-home/src/{LiveOverview,ActiveStack}.tsx`,
  operations backend). It refreshes every minute and lists 176 workloads.
- Per-service Prometheus alerts exist only on entity pages
  (`plugins/plugin-onemind-services-backend/src/prometheusAlerts.ts`), scoped
  per namespace. There is no cluster-wide view.
- Loki, Tempo, Grafana, Alertmanager and kube-state-metrics are installed in
  `monitoring`.

## Plugin first

Compose existing plugins before writing panels. Candidates to evaluate and record:
- native Kubernetes plugin (installed) for workload and pod state
- Red Hat **Application Topology for Kubernetes** for a live graph
- TeraSky **Kubernetes resources** for a CRD-aware resource view
- **KwirthMetrics** / **KwirthLog** or **Kubelog** for live metrics and logs
- **k8sgpt** for "why is this failing" explanations (needs the k8sgpt operator)
- Argo CD and Grafana plugins (installed)
- Roadie Prometheus (installed) for alert tables

Only the cluster-wide aggregation and health banner, which no listed plugin
provides, should be custom.

## Build

1. **Backend aggregator** in the operations backend: one read-only route
   (`/api/onemind-operations/live`). It fans out to Prometheus (instant queries),
   the Alertmanager API, the Argo CD API (existing integration) and the
   Kubernetes API (existing read-only identity), with per-source timeout,
   `observedAt` and `status: ok|stale|error`. Cache for 15 seconds.
2. **Live updates.** Push changes over Signals (F3 channel). Poll fallback is
   30 seconds.
3. **Frontend page** in `plugin-onemind-home` (or a new
   `plugin-onemind-live`): use a dense table and card layout, with each item
   linking to the entity page tab that explains it. Mobile-first stacking.
4. **Reuse.** Move Today's "Needs attention" logic into the aggregator so
   Home (F1), Today and Live Ops share one source.
5. **Queries.** Keep PromQL in a reviewed file
   (`plugins/.../liveQueries.ts`) with unit tests on fixture responses.
6. **Grafana embeds.** Use the existing Grafana plugin or iframe with the
   configured Grafana URL. No new datasource credentials.

## Out of scope

Write actions (restart, sync, scale). They stay out until the #339 approval
gates exist. Links to Argo and Grafana are fine.

## Files

`plugins/plugin-onemind-operations-backend/**`, `plugins/plugin-onemind-home/**`
(or a new `plugins/plugin-onemind-live/**`), `packages/app/src/App.tsx`
(registration), `packages/backend/src/index.ts` (if a new backend module), tests,
`docs/operations/live-ops.md`.

## Acceptance

- [ ] With today's cluster, the page shows `backup-offsite-copy` as pending setup, the two superseded jobs as retired (not attention items), and the
      `platform-matrix` OutOfSync / `platform-openclaw` Progressing states (or
      their current equivalents).
- [ ] Stopping one source (a fixture with a Prometheus timeout) shows that panel
      as `error` with a time; the others still render.
- [ ] A test alert appears on the page within 5 seconds without reload (with F3).
- [ ] Every item links to the owning app (via F2) or to "Unassigned".
- [ ] Route denied to unauthenticated and non-admin users where required
      (negative test). Required gates pass.

## Status

(Implementing agent: write progress here.)

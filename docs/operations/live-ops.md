# Live Ops (flagship F4)

Work order: `docs/roadmap/flagship/F4-live-operations.md`. This describes **source
behavior**. Nothing here is live until the F4 image is released (`pipeline.py publish`,
digest-pin PR, Argo) and the three CronJob annotations are synced by Argo.

## What it is

`/live` answers "is anything broken, and where" on one screen. Home's status strip,
Today's "Needs attention" and Live Ops all read the **same aggregator output**, so the
numbers cannot disagree.

```
Prometheus ─┐
Alertmanager ┤
Argo CD (CRs)├─▶ collector (operations backend) ─▶ GET /api/onemind/live    (Live Ops page)
Kubernetes  ─┤                                  └▶ GET /api/onemind/status  (Home strip, Today)
Service checks┘
```

- Code: `services/onemind-portal/plugins/plugin-onemind-operations-backend/src/`
  (`collect.ts` assembles, `build.ts` holds the pure rules, `liveQueries.ts` holds the
  reviewed PromQL, `backupState.ts` holds the backup rules, `source.ts` the timeouts).
- Page: `services/onemind-portal/plugins/plugin-onemind-live/`.
- Routes are read-only GETs and need a signed-in **user** (`allow: ['user']`): no token is
  401, a service or plugin credential is refused. They read nothing the existing
  `backstage-observer` identity (`platform/backstage/rbac.yaml`) cannot already read, and
  the identity was not widened.

## Sources and honesty

Each source has its own 5 to 8 second timeout and a 15 second cache. The response reports
`status: ok | stale | error` and the time the data was actually read (`observedAt`) for every
panel.

| Status | Meaning | Page shows |
|---|---|---|
| `ok` | read within the cache window | panel |
| `stale` | the latest read failed; data up to 5 minutes old is served | panel, "Stale", original time |
| `error` | no usable data | "Unavailable"; counts become `unknown`, never `0` |

One failing source never blanks the rest. If nothing is wrong but a source is missing,
the banner says it **cannot confirm health** (`unknown`), not "healthy".

| Panel | Sources |
|---|---|
| Firing alerts | Alertmanager `GET /api/v2/alerts` (Watchdog and InfoInhibitor dropped) |
| Argo CD | Applications in the `argocd` namespace (OutOfSync, Degraded, Progressing over 10 minutes, failed last operation) |
| Workloads | Pods, Deployments, StatefulSets (crash loop, image pull, config error, Pending over 2 minutes, NotReady over 5 minutes, restarts and OOM kills in the last hour) |
| Nodes | Nodes + pod requests (Kubernetes), CPU, memory, disk (Prometheus) |
| Storage | `kubelet_volume_stats_*` (over 80 percent, or projected full within 4 days) |
| Backups | CronJobs (Kubernetes), VolSync `volsync_volume_out_of_sync` (Prometheus) |
| Uptime | The existing cluster-side endpoint checks |
| Warning events | `Warning` events, last 30 minutes, repeats merged |
| Grafana | Five `d-solo` panels, embedded only when `onemind.grafanaEmbedding.enabled` is true |

## Backup states

The CronJob annotation (in Git) says what a job is **for**; its status says whether it is
**working**.

| `onemind.io/backup-state` | Meaning | Attention item? | Shown as |
|---|---|---|---|
| `superseded` (+ `onemind.io/superseded-by`) | replaced on purpose | no | Retired |
| `pending-setup` (+ `onemind.io/pending-reason`) | waiting for something only a person can do | no, counted separately | Needs setup |
| absent or `active` | expected to work | yes when suspended, never succeeded, last run failed, or overdue | OK / Last run failed / Overdue / Never succeeded / Suspended |

Overdue is 36 hours without success for jobs that run daily or more often (schedule has `*`
in day-of-month, month and day-of-week). Weekly or monthly jobs have no assumed deadline
unless they carry `onemind.io/backup-max-age-hours`. A mistyped state value counts as
`active`, so the mistake shows up as an attention item rather than hiding a problem.

Current annotations:

| CronJob | State | Detail |
|---|---|---|
| `backstage/backstage-backup` | superseded | by `backup-postgres-backstage` |
| `clickhouse/backup-clickhouse` | superseded | by the daily VolSync snapshot |
| `ops/backup-offsite-copy` | pending-setup | waits for SealedSecret `offsite-copy-credentials` (`platform/backups/OFFSITE-COPY.md` steps 1 to 5) |

All three keep `suspend: true`. Not covered yet: Velero (exports no metrics) and
CloudNativePG (its metrics carry no cluster name and read `0`, meaning no backup recorded);
the page says so.

The F3 observer (`plugin-onemind-events-backend`) only watches the CronJobs named in
`onemind.operationalEvents.backupCronJobs` (default `backup-postgres-backstage`), so the two
superseded jobs never raise inbox items. It does not read these annotations.

## Live updates

The page and the Home strip (and Today) refresh once on load, every 30 seconds, and about
0.4 seconds after any message on the F3 Signals channel `onemind:alerts`. Without Signals
the page says it is polling every 30 seconds.

## Check it

```bash
export PATH="$HOME/.nvm/versions/node/v24.18.0/bin:$PATH"
cd services/onemind-portal
node .yarn/releases/yarn-4.13.0.cjs workspace @internal/plugin-onemind-operations-backend test
node .yarn/releases/yarn-4.13.0.cjs workspace @internal/plugin-onemind-live test
node .yarn/releases/yarn-4.13.0.cjs workspace @internal/plugin-onemind-home test
```

Signed-in acceptance after release: open `/live`, confirm `backup-offsite-copy` reads
"Needs setup", the two superseded jobs read "Retired" and are not in the attention count,
and that a test alert from Alertmanager appears without a reload.

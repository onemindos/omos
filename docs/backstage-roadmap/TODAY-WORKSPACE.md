# OneMind Today workspace

`/today` combines live operational observations, the Notifications inbox, Zoho tasks and
maintained GitHub pull-request review/open-work cards. OneMind's light/dark brand assets
are preserved. The app directory remains on `/`. New/default Home layouts place Today
first; existing saved Home layouts can open Today through sidebar/directory navigation.

## Data boundaries

Cluster checks refresh each minute, initially showing five attention items with an explicit
button for all observations/metrics. Inbox refreshes every 30 seconds; an empty inbox is
not proof of health. GitHub authorization remains separate from Authentik sign-in and both
queries are restricted to onemindos/omos, the read-only App's installed repository. Zoho
remains configuration-required without the US tasks.READ grant and numeric portal ID.
No tasks/dates are fabricated; stale/partial observations stay visible. Calendar, NATS agent
completion events, billing adapters and email channels are not enabled.

## Internal deployment and backup events

`operationalEvents.ts` attaches to the existing onemind plugin. Backstage's global scheduler
observes each minute with existing read-only Kubernetes permissions. No external webhook
credential, additional cluster permission or Kubernetes mutation is introduced. Production
sets `onemind.operationalEvents.enabled`; local development defaults off.

Sources: terminal Argo Application operations in argocd (Succeeded/Failed/Error with
finishedAt), and terminal jobs owned by CronJob backstage-backup in backstage (Complete/
Failed). Standalone manual Jobs and unrelated schedules are outside this source's coverage.

The observer checks the live onemind-admins group and reciprocal User memberOf relations.
Recipients are fixed to that group. The first successful source scan establishes a durable
baseline without historical notification replay. Later new occurrences use stable scopes:
Application UID + operation finish time/phase, or Job UID + terminal condition. Raw Kubernetes
payloads, credentials, arbitrary recipients and arbitrary links are excluded.

Two additive tables in the existing onemind database persist baseline/processing records
and source status; existing all-plugin backups include them. The global scheduler avoids
ordinary overlapping replica runs. A crash after notification processing but before recording
may retry the same scope; exactly-once delivery is not claimed. notificationService.send
acknowledges processing, not delivery. `GET /api/onemind/event-sources` requires a signed-in
user and returns persisted source readiness/failures and scan timestamps. Older successful
scans are marked stale in Today. Each source has independent coverage.

Operation completion does not establish current health. Backup Job success does not separately
prove archive restoration; inspect the actual restoreVerified report.

## Recovery and acceptance

Disable the observer in production config to pause future processing while preserving its
baseline. Roll back the pinned image through Git/Argo if needed. Schema changes are additive;
no catalog, identity or storage deletion occurs.

Acceptance covers compilation/build, CI, pinned rollout, signed-in Today/source rendering,
anonymous denial, intended member inbox delivery, source boundaries, themes and two ready
replicas with Argo Healthy/Synced. Record runtime evidence after release; compilation alone
is not delivery evidence.

References:
- https://backstage.io/docs/notifications/
- https://backstage.io/docs/notifications/usage/
- https://backstage.io/docs/backend-system/core-services/scheduler/
- https://backstage.io/docs/backend-system/core-services/database/

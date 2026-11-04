# Dapr single-scheduler request: blocked on chart and data migration

Task 1 requests one scheduler and one placement server for the single-operator
cluster. This assessment does **not** implement the scheduler reduction. No
deployment, values, Argo policy, backup resources or data were changed.

## Verified baseline

Read-only observations on 2026-10-23 at approximately 18:52–18:56 UTC, using context
`do-nyc3-omos-cloud`:

- Argo Application `argocd/ecosystem-dapr`: Synced/Healthy, chart `dapr` 1.18.5,
  Git revision `c95e2470249e6bfad59a9384430321304ed26ce8`.
- `dapr-system/dapr-scheduler-server`: 3/3; image
  `ghcr.io/dapr/scheduler:1.18.5`; embedded etcd enabled with three initial peers.
- `dapr-system/dapr-placement-server`: 1/1; image
  `ghcr.io/dapr/placement:1.18.5`. Placement already meets the request.
- Scheduler StatefulSet PVC retention: `whenScaled: Retain`,
  `whenDeleted: Retain`. Each scheduler PV has reclaim policy `Retain`.
- The Application automatically syncs with self-heal enabled, pruning disabled.
  A warning in a PR body would not prevent a merged manifest from applying.

## Why a values-only change cannot implement this

The owning Application is
`platform/ecosystem/foundation-applications/ecosystem-dapr.yaml`; its values are
`platform/ecosystem/values/dapr.yaml`. The pinned chart's
`charts/dapr_scheduler/templates/dapr_scheduler_statefulset.yaml` hardcodes
`spec.replicas: 3` and a three-member `--etcd-initial-cluster` argument.
It does not use `dapr_scheduler.replicaCount` or `global.ha.replicaCount` for these.
Setting either to 1 therefore does not produce the requested scheduler.

`global.ha.enabled` is already false. The placement template uses that setting
and its default `ha: false` to render one placement replica. No placement edit
is necessary.

Reproduce without changing the cluster:

```bash
helm template dapr dapr --repo https://dapr.github.io/helm-charts \
  --version 1.18.5 --namespace dapr-system \
  -f platform/ecosystem/values/dapr.yaml \
  --set dapr_scheduler.replicaCount=1 --set global.ha.replicaCount=1 \
  > /tmp/dapr-single-scheduler-render.yaml
```

Inspect both StatefulSets and the scheduler peer list in that render: the
scheduler still has three replicas and placement has one.

The [Dapr 1.18 scheduler documentation](https://v1-18.docs.dapr.io/concepts/dapr-services/scheduler/)
warns that changing scheduler replica counts with embedded etcd risks losing
scheduling data. A replica-only patch does not reconfigure existing etcd
membership. Merely retaining ordinal 0's volume is not a migration plan.
Workflow timers and reminders are scheduler data; a broken quorum could stall
them indefinitely rather than cause a brief restart pause.

## Scheduler volume inventory

All three claims were Bound, 5 GiB, RWO, using `do-block-storage-retain`.
Names are an observation, not permission to delete anything.

| PVC in `dapr-system` | PV | Potential future role |
| --- | --- | --- |
| `dapr-scheduler-data-dir-dapr-scheduler-server-0` | `pvc-3c5d2282-3513-45f8-b45c-16e69edbafb6` | Candidate retained member; requires a proven migration |
| `dapr-scheduler-data-dir-dapr-scheduler-server-1` | `pvc-901288d7-ba97-4b83-8ed7-2cdc99c059d7` | Candidate retirement only after migration and recovery verification |
| `dapr-scheduler-data-dir-dapr-scheduler-server-2` | `pvc-ff670dea-b2c7-4457-a593-4c0f005caf55` | Candidate retirement only after migration and recovery verification |

These volumes are in use now, not orphaned. A scale-down with the observed
retention policy would leave claims 1 and 2 and their PVs intact. Their combined
10 GiB allocation would remain billable. Deleting a claim with PV reclaim policy
Retain does not itself delete the underlying DigitalOcean disk.

VolSync also has three separate 1 GiB cache claims in this namespace. Those are
not the scheduler data volumes. Backup definitions under
`platform/backups/volsync/` are explicitly outside this task's allowed paths.
Any retirement plan that needs to change them must stop for separate review.

## Reviewer steps and acceptance gates

After this documentation-only PR merges, no sync, restart or cleanup is required.
Keep the current healthy scheduler and all volumes intact. Task 1 remains blocked.

A subsequent implementation needs a reviewed Argo-owned rendering solution for
one scheduler **and** one etcd peer, plus a rehearsed data migration. A chart
upgrade, fork or external etcd deployment is a separate design decision; none
is selected or authorized by this assessment.

Before a migration PR can merge, require a fresh, verified scheduler backup and
restore rehearsal, workflow/reminder inventory, a maintenance window, a rollback
procedure and sequencing that accounts for automated Argo reconciliation. Preserve
the existing volumes until restored jobs and an application-level workflow timer
have been verified. Backup/export data must stay outside Git and PR logs.

After the future migration, verify the exact merged Argo revision, Synced/Healthy,
one ready scheduler, one ready placement server, correct single-member etcd
membership, and successful OneMind Core workflow start, timer and resume behavior.
A single scheduler has no failover: its restart briefly pauses workflow timers
when storage and membership are healthy; node/disk recovery can take longer.

Only then may Zeus separately approve removal of the two unused scheduler claims,
retained PVs and corresponding DigitalOcean disks, after checking backups and
pod references again. There are deliberately no deletion commands in this PR.

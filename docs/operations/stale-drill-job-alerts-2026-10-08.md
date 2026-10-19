# Stale drill Jobs and KubeJobFailed alerts, 2026-10-08

Several manual backup drill Jobs from Oct-8 remain in the cluster with `Error` pods, so Prometheus continues to raise `KubeJobFailed` (cluster-review A10). These are stale drill artifacts: later variants, and the scheduled `ops/backup-cluster-recovery`, succeeded.

**Operator action:** audit and delete only the six identified failed drill Jobs by name through the normal live-cluster change process. Do not add deletion to Git-owned workloads.

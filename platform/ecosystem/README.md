# OneMind CPU ecosystem rollout

The foundation parent activates seven pinned Helm releases through Argo CD:
CloudNativePG, Dapr, KEDA, Argo Rollouts, OpenCost, Kyverno and Falco.
Their values are in `values/`; `foundation-lock.json` records release provenance.
Existing Cilium, monitoring, NATS, Authentik and application data remain owned by
existing applications. New workloads use the separate ecosystem worker pool;
Falco agents inspect the participating cluster nodes.

Kyverno initially has no enforcing policies. Falco reports JSON detection events
into the existing container log pipeline. OpenCost uses the existing Prometheus.
Dapr sidecars are opt-in. No GPU is requested.

Argo auto-sync uses self-heal with pruning disabled for the initial additions.
Resources and persistent volumes require an explicit reviewed retirement change.
Subsequent application releases must supply their own database credentials,
identity integration, backups and functional acceptance evidence.

## Acceptance

Render each pinned chart with its checked-in values, then verify Argo sync,
controller readiness, PVC binding and logs. Verify that the existing applications
remain healthy. Record successful functional checks separately from readiness.
Pending application definitions outside `foundation-applications/` are not active.

# OneMind CPU ecosystem worker pool

Standalone Terraform root adding a pool to the existing cloud cluster. It does not
manage the legacy root's VPC, cluster or core node pool. Two NYC3 AMD workers each
provide 8 vCPU / 32 GiB; the configured autoscaler range is 2–4 workers.

The provider API listed this size at $168/month or $0.25/hour per worker on
2026-10-08: two baseline workers are approximately $336/month, excluding volumes,
load balancers, model usage and other services. Actual hourly usage and future
pricing determine billing. No GPU is requested.

The release helper initializes a private Spaces state backend using existing
backup credentials obtained locally from Kubernetes, with credentials passed only
through subprocess environment variables. State and plans must never be committed.
Commit the provider lock file. Review the plan's one pool creation before applying.

Once the pool is Ready, Helm values select its `onemind.io/compute-role=ecosystem`
label for the larger new workloads. Existing application scheduling is preserved.

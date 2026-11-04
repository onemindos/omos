# Existing cluster metrics-server

This adopts only the existing manually installed `kube-system/metrics-server`
Deployment into `platform-metrics-server` to repair its missing memory limit.
The existing Service, APIService, ServiceAccount and RBAC are preserved; there is
no new metrics service, role or namespace. The dedicated Application avoids the
legacy `platform-core` blanket Deployment-resource ignore rule.

Its unchanged100m CPU/200Mi memory requests are retained; memory is capped at512Mi
against an observed seven-day maximum working set of50.6Mi. Server dry-run proved
that this limit is the only workload-spec change. The existing v0.8.1 image, probes,
arguments, selectors, rollout strategy and credentials are preserved. Prune/Delete/
Replace/Force are disabled; the Application has no deletion finalizer and prune is
off. Normal Deployment rolling update keeps the old metrics pod until a new one is
Ready. Verify the exact runtime image and metrics API after adoption.

# Resource requests and limits — 2026-10-23

Task 3 is based on a read-only workload inventory from `do-nyc3-omos-cloud`,
per-container `kubectl top`, and Prometheus range queries. The current inventory
contains 55 incomplete container resource entries, including system, generated,
TAK and explicitly off-limits workloads; the earlier 22-workload audit count is
not recertified. This PR changes eligible source manifests and owning chart values,
never generated live objects. It does not claim all cluster violations are fixed.

CPU history is `avg_over_time(rate(container_cpu_usage_seconds_total[5m])[7d:15m])`;
memory history is `avg_over_time(container_memory_working_set_bytes[7d:15m])` plus
`max_over_time(...[7d:15m])`, filtered by namespace/container and excluding POD and
empty container labels. Below, historical numbers are the maximum of available
per-pod series, avoiding double-counting old replicas. The query spans seven days,
but scrape continuity/full seven-day retention and peak capture are not certified.
Sampled maxima can miss short bursts. CPU is converted from cores to millicores;
memory is converted from bytes to MiB. Init-container measurements are often absent.

## Changed budgets

All new memory limits are 2x their new memory request. CPU limits are unchanged.
Existing complete application budgets are preserved except the ingress memory
request, raised above observed usage so its new limit can have safe headroom.

| Workload/container and owning source | CPU request old → new | Memory request old → new | Memory limit old → new | Evidence and rationale |
| --- | --- | --- | --- | --- |
| ingress-nginx/controller; platform/values/ingress-nginx.yaml | 100m → 100m | 90Mi → 192Mi | absent → 384Mi | Top 2–3m / 156–162Mi; 7d CPU avg up to 2.68m, memory avg 151.61Mi, sampled max 166.23Mi |
| Prometheus and Alertmanager config-reloader, including init-config-reloader; platform/values/kube-prometheus-stack.yaml | absent → 2m | absent → 64Mi | absent → 128Mi | Top 1m / 25–39Mi; 7d avg up to 0.38m / 38.58Mi, sampled max 58.38Mi; shared operator override covers both workloads |
| hermes/hermes-legacy install-gateway-deps; platform/hermes/agent/deployment.yaml | absent → 50m | absent → 128Mi | absent → 256Mi | Estimate, no init series; package download/unpack requires more headroom than a file check |
| matrix/synapse fix-permissions; platform/matrix/argo/synapse.yaml | absent → 10m | absent → 32Mi | absent → 64Mi | Estimate, no init series; one-shot chown |
| matrix/postgres init-data-dir; platform/matrix/argo/postgres.yaml | absent → 10m | absent → 32Mi | absent → 64Mi | Estimate, no init series; directory creation/chown |
| matrix/redis init-data-dir; platform/matrix/argo/redis.yaml | absent → 10m | absent → 32Mi | absent → 64Mi | Estimate, no init series; directory creation/chown |
| onetak/webodm-db fix-perms; platform/onetak/webodm.yaml | absent → 10m | absent → 32Mi | absent → 64Mi | Estimate, no init series; directory permissions |
| takml/takml-server check-jar; platform/takml/deployment.yaml | absent → 10m | absent → 32Mi | absent → 64Mi | Estimate, no init series; shell file-size/existence check; TAK Server resources are excluded |
| onemind-core/openfga migrate-database; platform/ecosystem/values/openfga.yaml | absent → 50m | absent → 64Mi | absent → 128Mi | Estimate, no init series; small database migration client; main server resources unchanged |
| velero/velero plugin copy init; platform/ecosystem/values/velero.yaml | absent → 10m | absent → 32Mi | absent → 64Mi | Estimate, no init series; plugin binary copy, not a backup upload |
| velero/node-agent; same values | absent → 25m | absent → 128Mi | absent → 256Mi | Top 1–3m / 25–43Mi; 7d avg up to 2.52m / 35.20Mi, sampled max 46.65Mi; conservative extra working room for future backup traffic |
| falco/falcoctl-artifact-follow; platform/ecosystem/values/falco.yaml | absent → 10m | absent → 128Mi | absent → 256Mi | Top 0–1m / 25–44Mi; 7d avg up to 0.09m / 44.68Mi, sampled max 99.25Mi; extra headroom for artifact download/unpack bursts missed by sampling |
| falco/falcoctl-artifact-install; same values | absent → 10m | absent → 128Mi | absent → 256Mi | Estimate; align with the updater download/unpack headroom |
| tailscale/operator; platform/values/tailscale-operator.yaml | absent → 25m | absent → 192Mi | absent → 384Mi | Top 5m / 81Mi; 7d avg up to 6.92m / 146.60Mi, sampled max 153.23Mi; correct ignored operator.resources to supported operatorConfig.resources |
| default-class Tailscale proxy sysctler; platform/core/tailscale-proxyclass.yaml | absent → 10m | absent → 16Mi | absent → 32Mi | Estimate, no init series; supported tailscaleInitContainer.resources only; existing proxy main-container budgets stay unchanged |

The Velero plugin retains its existing S3-protocol adapter image and provider
string; its configured endpoint is **DigitalOcean Spaces**, not an AWS service.
No AWS CLI, AWS account, credential or endpoint was used or introduced.

## Owning-chart evidence

Supported values were inspected and rendered at the pinned versions:

- ingress-nginx 4.15.1: controller.resources.
- kube-prometheus-stack 87.16.1: prometheusOperator.prometheusConfigReloader.resources.
- OpenFGA 0.3.16: datastore.migrations.resources, verified on migrate-database.
- Velero 12.2.1: initContainers[].resources and nodeAgent.resources.
- Falco 9.2.0: falcoctl.artifact.install/follow.resources.
- Tailscale operator 1.98.4: operatorConfig.resources; the previous operator.resources
  stanza was ignored, which explains the live operator gap. The installed
  ProxyClass CRD and pinned chart both support tailscaleInitContainer.resources.
  The existing proxyConfig.defaultProxyClass is retained. See the
  [pinned chart values](https://github.com/tailscale/tailscale/blob/v1.98.4/cmd/k8s-operator/deploy/chart/values.yaml).

OpenFGA and Falco Application sources live under platform/security, outside the
current validate_helm_renders.py discovery directories; they were also rendered
explicitly and their target containers checked, rather than trusting ignored keys.
No chart version, image, OAuth credential, auth policy, exposure or backup schedule
is changed. Resource estimates are starting budgets, not load-test results.

## Exclusions and residual findings

- DOKS-managed Cilium/Hubble, CSI, metrics-server, konnectivity, DigitalOcean
  agents and telemetry are excluded: change through the provider, not unmanaged
  workload patches in this repository.
- TAK Server, FedHub, LDAP, TAK NATS/PostgreSQL, voice sidecars and retention resources
  are excluded as requested; no platform/takserver resources changed.
- Ollama, OpenChoreo, VolSync and the forbidden ecosystem OpenSearch/ESO/control-plane
  application paths are unchanged. In particular, OneMind Search init containers
  belong to the off-limits OpenSearch application.
- Generated Tailscale StatefulSets are never patched. The shared default ProxyClass
  change will apply only where that class is selected. ProxyGroup egress has no
  proxyClass and no ProxyGroup manifest was found in Git, so its sysctler gap remains
  a separate ownership/adoption issue. Do not imply every generated proxy is covered.
- No policy-report objects were returned by the PolicyReports API; observed missing
  resource fields, not a reconstructed audit count, drive this PR.

## Exact reviewer steps after merge

1. Recheck node allocatable/request headroom and the final combined diff before
   merging. Active sidecars add scheduling requests; init requests often fit under
   the already larger application budget, but this is not a capacity guarantee.
2. Allow only each owning Argo Application to reconcile. Expect bounded pod rolls
   for changed templates/operator-generated workloads; observe readiness and
   access to critical admin routes, especially ingress and the Tailscale operator.
   No hand scaling, Helm upgrade, kubectl apply, generated-object edit or secret
   operation is required.
3. Keep Matrix manual. Its missing synapse-homeserver Secret must be handled using
   the separately reviewed Matrix drift runbook before any application sync.
   Resource-only intent does not justify syncing unrelated Matrix drift.
4. Inspect rendered/live resources on every changed main/init container, and check
   Prometheus/Alertmanager reloaders and default-class proxy regeneration specifically.
   Watch Pending, OOMKilled, restart counts and actual backup/model gateway behavior.
   Re-measure these initial budgets after a representative workload/backup cycle.
5. Confirm the actual Kyverno findings decrease; explicitly retain provider-managed,
   forbidden and unowned egress exceptions. Do not claim a clean cluster audit.
6. If a budget causes failure, use a reviewed Git adjustment/revert and ordinary
   Argo reconciliation. Keep storage, credentials, routes and backup schedules intact.

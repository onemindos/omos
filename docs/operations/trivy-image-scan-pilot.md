# Trivy Operator image vulnerability pilot

Prepared 2026-10-24 for Zeus; retargeted 2026-11-04 to `auth,nats,onetak,n8n` after
owner approval to deploy (Nextcloud was retired 2026-11-01 and is no longer a target).
**Not live-accepted by this PR.** Merging deploys it: `platform-security` auto-syncs
the RBAC/quota and the `ecosystem-trivy-operator` child (see GitOps rollout).
Implementation PR: [#442](https://github.com/onemindos/omos/pull/442).
This implements the requested bounded scanning pilot alongside Kyverno/Falco.
It does not block deployments, patch apps, install an admission webhook, or run
commands inside application pods. Authentik's login program is separate.

## Scope and ownership

| Item | Configuration |
|---|---|
| Parent | Existing `platform-security`, `platform/security/`, prune off |
| Helm child | `ecosystem-trivy-operator`, release `trivy-operator`, namespace `trivy-system`, sync wave 5, prune off, no cascading finalizer |
| Chart | Aqua `trivy-operator` **0.37.0**, operator **0.35.0**, scanner **0.75.0**; archive checksum in both chart metadata/verified foundation lock as appropriate |
| Image references | Operator and scanner index digests pinned in `platform/ecosystem/values/trivy-operator.yaml`; observed public mirror manifest digests on preparation date |
| Workloads | Current ReplicaSets and StatefulSets in **auth,nats,onetak,n8n**, including ordinary and init containers. Public images are scanned; the private `ghcr.io/onemindos` image (`tak-portal`) is **not** (see Private images) |
| Scan jobs | In `trivy-system`; **one concurrent Job**, 15-minute Job deadline, 12-minute scanner timeout, 5-minute retry delay/cleanup TTL |
| Reports | Namespaced `VulnerabilityReport`; 24-hour TTL triggers refresh, and workload changes trigger new reports |
| Optional scanners | Configuration, RBAC, nodes/host filesystem, compliance, secret and SBOM generation/cache **disabled**; no compliance specs instantiated |
| Monitoring | Chart ServiceMonitor with `release: monitoring`, honorLabels true, 60-second scrape; dashboard UID `onemind-trivy-pilot`, folder OneMind |
| Storage/exposure | Ephemeral scanner cache; **no PVC, LoadBalancer, ingress, additional node or paid service** |

No report is a claim that an image is exploitable, clean or patched. Review CVE,
package, installed/fixed versions, applicability and image digest before opening
a remediation PR. All severities and unfixed findings remain visible. We do not
hide unresolved vulnerabilities just to obtain a green dashboard.

The standalone scanner downloads current vulnerability data and public target
images over HTTPS. This has real registry/network/CPU/ephemeral-storage cost;
one Job is not one container: a workload with two containers runs two scanner
containers. It does not read live application volumes or database contents.
An initial scan is not continuously observing runtime attacks; Falco handles
that separate concern. Private GHCR image coverage is a later scoped change,
not a reason to give this pilot cross-namespace credential access.

## Read-only capacity and integration preflight

Captured **2026-10-24 03:25 EDT / 07:25 UTC**, context
`do-nyc3-omos-cloud`. Desired pod requests include regular containers, native
sidecars, sequential init maxima and pod overhead; completed/failed pods excluded.
These are scheduling requests, not measured consumption.

| Node | Requested CPU / allocatable cores | Requested memory / allocatable MiB |
|---|---|---|
| `onemind-ecosystem-cpu-3rhwrp` | 5.609 / 7.88 | 22329 / 28767 |
| `onemind-hub-core-v2-3rtf47` | 3.734 / 7.88 | 9678 / 13646 |
| `onemind-hub-core-v2-3rtf4m` | 3.532 / 7.88 | 9576 / 13646 |
| `onemind-hub-core-v2-3rtf4q` | 5.567 / 7.88 | 9485 / 13646 |

At the same preparation snapshot, `kubectl top nodes` showed memory use 48–65%
on hub nodes and 55% on the ecosystem node. Both pilot namespaces exist. The
existing Prometheus selects ServiceMonitors labeled `release: monitoring` in all
namespaces. The Kubernetes API Service endpoint port is 443. No Trivy deployment
or Aqua report CRDs were found during the preceding read-only inventory.
Recheck these facts and node ephemeral storage before merge/rollout; they can change.

Workload inventory re-read from the cluster on **2026-11-04** (read-only; current,
non-zero-replica ReplicaSets and StatefulSets): **21 controllers, 24 ordinary
containers and 1 init container**. On 2026-11-07 NUI replaced heynats, nats-dashboard and nats-ui (a public image),
which leaves 19 controllers and 22 ordinary containers; the counts below were not re-read.

| Namespace | Controllers | Containers (init) |
|---|---|---|
| `auth` | 4: authentik server, worker, redis, postgres | 4 |
| `nats` | 5: nats StatefulSet (nats + reloader), nack, nui, nats-surveyor, leafz-exporter | 6 |
| `onetak` | 8: nodeodm, mediamtx, tak-portal, tile38, tileserver-gl, webodm (webapp + worker), webodm-broker, webodm-db | 9 (+1 init: `fix-perms` busybox) |
| `n8n` | 2: n8n (n8n + task-runners), postgres | 3 |

The largest Job is two scanner containers (nats StatefulSet, webodm, n8n, and
webodm-db with its init container). `mediamtx:1-ffmpeg`, `authentik-redis`
(`redis:7-alpine`), `postgres:17-alpine` (auth) and `redis:7.0.10` are mutable
tags; most others carry digests. Acceptance must match each report to the
actual running digest, not just these tags. The largest image (nodeodm) is about
800 MB compressed, so the 4Gi per-container ephemeral limit and 12-minute scanner
timeout are expected to hold; measure it, and raise only if a real OOM, eviction or
timeout is recorded.

### Private images (not scanned)

No registry credential is created or copied (`privateRegistryScanSecretsNames: {}`,
`accessGlobalSecretsAndServiceAccount: false`, no Secret read outside `trivy-system`).
This `ghcr.io/onemindos` image, which pods pull with the `ghcr` image pull
secret, **cannot be scanned**:

| Namespace | Workload | Image |
|---|---|---|
| `onetak` | `tak-portal` | `ghcr.io/onemindos/tak-portal@sha256:eaa14729...` |

The scanner pulls anonymously, is refused, and its container fails. In pinned operator
v0.35.0 a Job with no successfully completed container is deleted and **no
VulnerabilityReport is written**; there is no immediate retry (the 5-minute delay
only applies to concurrency push-back). The gap is therefore visible as: no report for
that workload, no row in the Grafana table (No data), a failed-scan log line in the
operator, and an events/Job failure in `trivy-system`. It must never be read as
clean. The operator may attempt the scan again on a later reconcile (operator
restart, workload change, cache resync); count those attempts during acceptance.
Other `ghcr.io` images (`goauthentik/server`, `nats-nui/nui`) are public.

- Operator request **50m / 128Mi**, limit **500m / 512Mi**; ephemeral storage
  request 64Mi, limit 256Mi.
- Each scanner container requests **100m / 512Mi / 1Gi ephemeral**, limited to
  **1 CPU / 2Gi memory / 4Gi ephemeral**. The widest current workloads have two
  scanner containers, so a Job requests 200m / 1Gi plus the operator: 250m / 1152Mi
  maximum, limits 4Gi plus the operator's 512Mi (4.5Gi, within the 5Gi quota). This
  is capacity planning from the current workload set, not a universal per-Job maximum.
  The 4 target namespaces add controllers but not concurrency: still **one Job at a
  time**, so the 19 workloads are scanned serially. Because the 5-minute push-back
  applies to the 18 waiting workloads, the first full pass can take well over an hour;
  this is accepted for a pilot and not tuned.
- Namespace quota caps total requests at **1 CPU / 2Gi memory / 4Gi ephemeral**
  and limits at **3 CPU / 5Gi memory / 10Gi ephemeral**, with four pods and five
  Job objects; these still suit the larger set (unchanged, see the 2026-11-04
  capacity re-check). This deliberately stops a future larger workload from consuming
  an unreviewed budget. Quota rejection is a failed/missing scan, not clean coverage.
- Operator and scanner jobs prefer the existing ecosystem pool, with scheduling
  fallback. No node-name pin, node scaling or reserved high priority is introduced.
- Standalone mode avoids a new persistent Trivy server; the tradeoff is repeated
  vulnerability DB downloads per Job. Measure duration, throttling, rate limits,
  storage and memory/OOM before changing cache architecture or concurrency.

Capacity re-check **2026-11-04** (read-only; sum of container requests per node,
allocatable memory minus requested): `onemind-ecosystem-cpu-3rhwrp` about **6.5 GiB**
free; hub nodes `3rtf47` / `3rtf4m` / `3rtf4q` about **4.4 / 4.8 / 4.1 GiB** free;
CPU requests 3.1 to 5.5 of 7.88 cores. `kubectl top` memory use was 58 to 61% on all
four nodes. The operator (128Mi) plus one Job (1Gi) fits on any node, and jobs prefer
the ecosystem pool. The first table above is the original 2026-10-24 snapshot.

## Permissions and network boundary

The chart's broad operator RBAC is **disabled**. Explicit roles/bindings are
owned by `platform-security` in `trivy-pilot-rbac.yaml` before the child starts:

- Namespace-bound workload metadata reads and vulnerability-report CRUD in
  `auth`, `nats`, `onetak`, `n8n` and `trivy-system`. No writes to Deployments, StatefulSets,
  pods or application configuration. Listing metadata includes workload specs;
  it is not permission to read application Secret objects.
- Read-only `exposedsecretreports` and `sbomreports` in those namespaces are
  necessary because pinned v0.35.0 unconditionally registers `Owns()` watches
  for both in its workload controller. Their scanners and report writes are off.
- Job create/delete, scanner pod-log reading, ConfigMap initialization and
  transient/empty Secret management only in `trivy-system`. Public pilot images
  need no registry credentials. No application Secret or service-account reads,
  node/proxy access, cluster-admin, cluster-wide Job create or cluster-report writes.
- Cluster-scoped rights are CRD get/list/watch and get of the three namespace
  records. MultiNamespace cache includes both targets plus the operator namespace.
- Operator needs its own Kubernetes API token. Scanner pods explicitly disable
  token automount, run non-root with RuntimeDefault seccomp, drop capabilities,
  disallow privilege escalation and use a read-only root filesystem with scratch
  volumes. Node-collector hostPath defaults are removed as well as disabling scanning.
- Upstream chart **always creates three report-view aggregate ClusterRoles**, even
  when operator RBAC is off. These add vulnerability/config/secret-report reads
  to standard view/edit/admin roles, including existing observers with those roles.
  This intentional visibility is report access only, not Secret object access or
  any new machine write right. Secret/config report generation remains disabled;
  reassess report visibility before enabling those features or broadening scope.
- Native NetworkPolicy restricts ingress to port 8080 from `monitoring`; egress
  permits kube-system DNS and TCP 443 for API/registries/databases. HTTPS egress
  is deliberately destination-unrestricted, so it can reach any HTTPS service;
  it is not a registry FQDN allowlist. Enforce/report this limit honestly. There
  is no node socket, host networking or mount of application storage.

Do not enable global credential access or grant an extra permission merely to
silence a Forbidden error. Diagnose the specific watched resource against the
pinned source and update tests/review first. The operator may create empty internal
configuration Secrets; no actual credentials or plaintext secret values enter Git.

## GitOps rollout and acceptance

All live checks below remain **pending**. Source/render checks do not close them.
Merging the PR makes the child eligible for the existing automated GitOps flow;
there is no imperative Helm installation or `kubectl apply` step.

1. Recheck capacity, API/CRD compatibility and existing Argo operation state.
   `platform-security` applies the restricted namespace at wave -1, then RBAC,
   quota and policy before the child Application at wave 5. The Helm child owns
   its CRDs/operator/metrics ServiceMonitor; no resource has two owners. The
   Grafana ConfigMap is owned separately by existing `platform-monitoring`.
2. Check actual child **operation revision/chart**, values source revision,
   Deployment image digest, 1/1 Ready and logs. The root's compared Git revision
   alone is not proof that the new child was created or synchronized. Existing
   stalled root/parent operations can delay rollout; diagnose them separately.
3. Use read-only authorization checks (the operator SA does not need impersonation):

   ```bash
   kubectl --context do-nyc3-omos-cloud auth can-i create jobs -n trivy-system --as=system:serviceaccount:trivy-system:trivy-operator
   kubectl --context do-nyc3-omos-cloud auth can-i create jobs -n auth --as=system:serviceaccount:trivy-system:trivy-operator
   kubectl --context do-nyc3-omos-cloud auth can-i get secrets -n nats --as=system:serviceaccount:trivy-system:trivy-operator
   kubectl --context do-nyc3-omos-cloud auth can-i patch deployments -n n8n --as=system:serviceaccount:trivy-system:trivy-operator
   ```

   Expected answers: yes, no, no, no (`nextcloud` is no longer a namespace or target). Also verify report writes allowed only in
   the pilot namespaces, read-only SBOM/secret-report watches work, and node/proxy
   access is denied. Check all explicit bindings, not just one Role document.
4. Observe at most one active scan Job, all scanner pods in `trivy-system`, no
   projected service-account token/host mounts, and no optional scanner Jobs.
   Record scan duration, memory/CPU/storage peaks, failures/retries and quota
   events. Confirm Authentik, NATS, OneTAK and n8n remain healthy and their app specs
   unchanged; NATS (3-node JetStream) and Authentik are the sensitive ones.
5. For **each current pilot ReplicaSet/StatefulSet container and init container**,
   enumerate image references/digests and match an actual VulnerabilityReport.
   Cover all four namespaces (19 controllers). Compare report scanner version,
   artifact digest, creationTimestamp, summary and a sample package/fixed version.
   Mutable target image tags may move: confirm report digest against running
   `status.containerStatuses[].imageID`. Inactive ReplicaSets are intentionally
   excluded. A private-image failure, OOM, timeout or absent report stays visible
   as incomplete coverage. Expect exactly one missing report (`tak-portal`, a
   private image) and verify the failures are
   logged, not looping. Do not add credentials or change application images
   automatically to obtain a successful report.
6. Prometheus target must be discovered/up and emit
   `trivy_image_vulnerabilities` for `namespace` in `auth`, `nats`, `onetak` and `n8n`.
   HonorLabels preserves the report's workload namespace against the scrape
   target's namespace; per-CVE ID metrics are off to keep cardinality bounded.
   Open Grafana `/d/onemind-trivy-pilot` using existing login, verify OneMind folder,
   namespace/workload/image/severity table, scrape state and findings versus CRDs.
   Dashboard sums per-container findings, not unique global CVEs. No data is
   never converted to zero; an up scanner is not evidence of complete coverage.
7. After the 24-hour report refresh boundary, confirm recreated reports and
   current digests without duplicate active jobs or accumulating completed pods.
   TTL deletion can produce a temporary no-data gap; record it. Check that reports
   outside the pilot are not generated and cluster state is otherwise unchanged.

Save dated release evidence (merge SHA, source/operation revisions, pinned images,
replicas, coverage inventory, failures, metric samples and screenshots without
tokens) here or in a linked operations receipt. No report DB backup is added:
these reports/cache are reproducible assessments, not authoritative app data.
Export an approved redacted snapshot if a remediation audit needs retention;
the 24-hour TTL is not an archive policy.

## Stop, rollback and follow-up

Stop expansion on Forbidden loops, OOM/eviction, repeated registry failures,
quota rejection, application health regression or misleading dashboard results.
Do not claim the pilot accepted simply because its Deployment is Ready.

First stop further scan creation by a reviewed Git change setting
`operator.replicas: 0`, preserving its Application and values. Existing jobs may
finish until their deadline; verify that separately. Revert the Grafana ConfigMap
addition through its owner if necessary. No PVC, namespace, app credential or
application workload needs deletion to pause the pilot.

Prune is off at both security parent and scanner child. A Git revert/removal of
the Application alone therefore does **not** uninstall the live operator. A full
retirement needs an explicit reviewed cleanup list after scale-to-zero and job
completion. Keep CRDs/reports and namespace until approved; deleting a report
CRD would delete its reports. Do not use Helm CLI uninstall, blanket namespace
deletion or finalizer stripping as rollback.

After acceptance, separately consider private-image credential references (a
SealedSecret in `trivy-system` plus `privateRegistryScanSecretsNames`, needing owner approval),
additional namespaces, configuration/RBAC reports, actionable alerts and OneMind
catalog links. Reassess permissions/capacity/visibility for each. Node collectors,
admission enforcement, automatic remediation and paid services stay out of this
pilot. Retain Kyverno and Falco ownership/behavior.

## Validation and sources

Run `make validate`, changed-YAML/manifest gates, pinned foundation chart checks,
`scripts/validate_helm_renders.py` and the GitOps safety contracts. The new Trivy
contracts render the checksum-verified chart and evaluate effective RBAC across
all explicit bindings, including a deliberately overbroad negative fixture.
The standard Helm gate now discovers `platform/security/` Applications too.
Local `TRIVY_CHART_ARCHIVE` can point to an already downloaded archive, but its
bytes must still match the committed checksum. Record actual pass/failure results
in the PR; unrelated repository lint debt is not a green full validation.

Preparation receipt: all 96 GitOps contracts passed, including six Trivy
render/security tests; changed YAML lint passed for seven files; plain-manifest
schema checks passed for 16 resources. The rendered chart plus new manifests
passed strict schema checks for 24 resources; 12 CRDs, absent from the public
schema catalog, passed separate Kubernetes server-side create dry runs with no
persistence. The eight checksum-locked foundation charts rendered successfully,
and the standard changed Helm gate discovered/rendered the new security child.
The running Prometheus accepted all six dashboard expressions (empty scanner
results are expected before deployment). Secret scan passed. `make validate`
stopped at the existing YAML debt in 40 files, none modified by this PR; later
make targets were not reached. No scanner deployment or finding is claimed.

- [Pinned operator release](https://github.com/aquasecurity/trivy-operator/releases/tag/v0.35.0).
- [Chart 0.37.0](https://github.com/aquasecurity/helm-charts/releases/tag/trivy-operator-0.37.0).
- [Operator configuration](https://aquasecurity.github.io/trivy-operator/latest/getting-started/installation/configuration/).
- [Metrics](https://aquasecurity.github.io/trivy-operator/latest/tutorials/integrations/metrics/).
- [Pinned workload watches](https://github.com/aquasecurity/trivy-operator/blob/v0.35.0/pkg/vulnerabilityreport/controller/workload.go).
- [Pinned centralized job builder](https://github.com/aquasecurity/trivy-operator/blob/v0.35.0/pkg/vulnerabilityreport/builder.go).
- [Pinned client/cache behavior](https://github.com/aquasecurity/trivy-operator/blob/v0.35.0/pkg/operator/cache.go).

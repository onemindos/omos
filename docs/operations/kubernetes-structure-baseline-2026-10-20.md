# Kubernetes / Backstage structure baseline — 2026-10-20

Read-only snapshot using context `do-nyc3-omos-cloud`, captured **16:33–16:37 UTC**.
This supplements older `STACK.md` inventories; it does not recertify all app
functions, permissions or backups. No Secrets were fetched for this inventory.

## Observed state

| Observation | Result | Limit |
|---|---|---|
| Kubernetes | 4 Ready amd64 nodes, kubelet 1.36.5, containerd 2.2.3 | One DOKS cluster, not a verified home/fleet platform |
| Argo Applications | 69 total; **67 Synced/Healthy** | Matrix OutOfSync/Healthy; OpenClaw Synced/Progressing |
| Workload objects | 188 Deployment/StatefulSet/DaemonSet objects | Includes zero-replica/system objects; not 188 validated applications |
| Pods | 252 total; no Pending or unready nonterminal pods in the bounded sweep | Completed/failed pods are excluded from the readiness assertion |
| Namespaces / CRDs | 54 / 196 | Counts are inventory, not security/feature acceptance |
| Backstage | Deployment 2/2; `platform-backstage` Synced/Healthy | Does not prove every catalog link, permission, adapter or underlying app works |
| NATS | StatefulSet 3/3; Legacy and Agent Console/Tap 1/1 | End-to-end inference and permission behavior remain separate gates |
| Authentik | Server 2/2; worker/Postgres/Redis present | Local/home recovery and all application SSO flows not tested in this sweep |
| Observability | Prometheus, Alertmanager, Grafana, Loki, Tempo, OTel agent, OpenSearch, OpenChoreo Observer/adapters present | Existing adapter/auth/data correctness follow-ups still apply |
| Security | Kyverno controllers; Falco 4/4; SPIRE server and node agents; OpenFGA present | No assertion of policy coverage or enforcement for every workload |
| Storage / recovery controllers | CNPG/Barman, Velero/node agents, VolSync, DOKS CSI present | No new application-consistent restore or independent-copy acceptance performed |

At the original 16:33–16:37 UTC snapshot, Backstage's deployed image was:
`ghcr.io/onemindos/onemind-portal@sha256:4a52c3b840fcd2b4e2fc0287929e1704721f7420018f2c09deb8543a6aab74da`.
That earlier image is superseded by #336; see the fresh delta below. Source
publication and actual deployment remain separate evidence.

Backstage MCP at **16:33:15 UTC** returned timestamped HTTP observations around
16:33:04: Identity, Grafana, Headlamp and NATS dashboard all HTTP 200; namespace
CPU/memory samples were available. The tool explicitly scopes these as internal
reachability/metrics, **not** browser login, semantic app health or NATS delivery.
In particular, Grafana's own health endpoint does not prove the portal adapter
has a working authenticated query credential.

## Platforms present versus proposed

Present: Argo CD/ApplicationSet/Rollouts; Backstage; OpenChoreo control/data/
observability planes; Authentik/Tailscale; NATS/JetStream/NACK; Hermes/Legacy and
the existing agent console/tap; CNPG/Barman; Harbor/Trivy; Dapr; KEDA; OpenCost;
metrics/logs/traces; Kyverno/Falco/SPIRE/OpenFGA; Velero/VolSync; Cilium/Hubble.

No matching CRD group was found for Crossplane, KubeVela, Karmada, KubeEdge,
OpenYurt, Rancher or KubeVirt in this cluster. That observation is narrower than
proving there is no standalone service anywhere. They remain candidates, not
installed capabilities. No OpenChoreo workflow-plane namespace was observed;
the workflow-plane pilot remains an explicit plan (#328), not a working CI path.

Legacy's earlier authorized probe on this date reached the actual agent API but
chat ended in an error. Runtime logs identified configured DigitalOcean inference
HTTP 402, account blocked. See the evidence and acceptance gates in
[#332](https://github.com/onemindos/omos/pull/332); no provider/billing/key changes.

## Public ingress lifecycle finding

Observed controller image:
`registry.k8s.io/ingress-nginx/controller:v1.15.1@sha256:594ceea76b01c592858f803f9ff4d2cb40542cae2060410b2c95f75907d659e1`.
The official community ingress-nginx project announced retirement in March 2026
and no further security maintenance afterward. This identifies a lifecycle risk,
not a demonstrated exploit in the observed image.
[Official retirement notice](https://kubernetes.io/blog/2025/11/11/ingress-nginx-retirement/)

27 Ingress resources were listed; **8 explicitly select class `nginx`**:

- `auth/authentik-ingress`
- `matrix/matrix-ingress`
- `matrix/matrix-wellknown`
- `nats/nats-ws-ingress`
- `onetak/geo-services-ingress`
- `onetak/onetak-ingress`
- `onetak/tak-portal-ingress`
- `takserver/fedhub-ingress`

This is a resource list, not eight hostnames or a complete enumeration of TCP/UDP
listeners. Both GatewayClasses `cilium` and `onemind-private` reference
`io.cilium/gateway-controller`; 3 HTTPRoutes were observed. **Private gateway
success does not prove public ingress annotation/feature compatibility.** Do not
upgrade provider-managed Cilium/CRDs or cut over identity based on this snapshot.

## What still blocks trustworthy production expansion

- Broader Phase 6 trust/security/recovery gates remain; the scoped #330 → #335 →
  #336 source/image/GitOps/live chain is verified below. STAC schema and Grafana issues require
  their own current verification/repair; not silently recertified here.
- Legacy inference failure; complete identity, database and recovery gates before
  activating the staged LibreChat surface.
- Earlier documented Helm-CLI/GitOps exceptions and independent backup gaps must
  be rechecked before adoption/cutover; this sweep did not read release Secrets.
- Earlier Mac free space was 13–14 GiB; a fresh check during this review shows
  about 26 GiB. No space has been allocated or external volume approved for a VM.
  External data/Time Machine drives are not approved for repurposing.
  Closed-lid/reboot/FileVault and local outage recovery remain untested.
- GitHub Quality Gate run `39458922274` could not start its jobs: check annotations
  report failed recent account payments or a spending-limit restriction. No billing
  change was attempted. Local `make validate` still fails on existing YAML debt
  beginning `platform/onetak/media-mtx.yaml:86`. Neither is a passing CI result.

## Fresh structural review — 17:40 UTC

Read-only inventory repeated while reviewing #333 after #331/#335/#336 merged:

- Four Ready amd64 nodes, kubelet v1.36.5, containerd 2.2.3; 54 namespaces,
  196 CRDs, 188 Deployment/StatefulSet/DaemonSet objects and 252 pods. No Pending
  or unready nonterminal non-deleting pods in this sweep.
- Argo: 69 Applications, 67 Synced/Healthy; the same Matrix OutOfSync/Healthy
  and OpenClaw Synced/Progressing exceptions. Counts do not certify all app behavior.
- Backstage now runs source `c9a05f8d0b1fc1809b0496945d99593e884f561d`, digest
  `sha256:6c729bd20c0a0a40330be8f0c7bd5213c59c7540bd1ea54c38e05f8c63fc3e7c`.
  Its two new pods are Ready with no restarts in the observation. Scoped browser,
  716-entity catalog and both-pod denial/MCP checks are in the
  [release evidence](stac-operational-slice-2026-10-20.md), not inferred from counts.
- **NATS ownership correction:** StatefulSet `nats/nats` has Argo tracking ID
  `platform-nats:apps/StatefulSet:nats/nats`, managed-by `argocd`, and is listed
  in `platform-nats` status. Do not follow an older report by blindly adopting
  it again or running Helm upgrade. Release remnants/field ownership still need
  K02 audit; metadata alone does not prove absence of every overlap.
- Argo's application-controller StatefulSet has Helm metadata, no observed Argo
  tracking ID/status owner in this inventory. Several TAK core resources likewise
  have no observed Argo tracking; TAK overlays/plugins and federation resources do.
  This is candidate audit evidence, not proof of missing deployment intent or
  permission to replace them. No Helm release Secrets were read.
- ESO deployment and CRD groups `external-secrets.io`/`generators.external-secrets.io`
  are present; `ecosystem-external-secrets` tracks the controller Deployment.
  This does not prove an external store is configured or authorize replacing the
  repository's SealedSecrets credential pattern.
- Eight explicit nginx-class Ingress resources and three HTTPRoutes remain.
  #334 owns route/parity assessment; no public route, DNS, TLS or identity change here.
- No matching Crossplane/KubeVela/Karmada/KubeEdge/OpenYurt/Rancher/KubeVirt CRD
  groups appeared in this repeat. No new platforms were installed by this review.

Private local observation files: `/tmp/onemind-333-workloads.json`,
`/tmp/onemind-333-argo.json`, `/tmp/onemind-333-inventory.json` and
`/tmp/onemind-333-network.json`. These are bounded operator snapshots, not backups
or a live-state API. No Secret objects or credential values were requested.

## Reproduction (read-only)

```bash
kubectl --context do-nyc3-omos-cloud get nodes -o json
kubectl --context do-nyc3-omos-cloud get applications.argoproj.io -A -o json
kubectl --context do-nyc3-omos-cloud get deployment,statefulset,daemonset -A -o json
kubectl --context do-nyc3-omos-cloud get pods -A -o json
kubectl --context do-nyc3-omos-cloud get namespaces,customresourcedefinitions
kubectl --context do-nyc3-omos-cloud get ingress,httproutes -A -o json
kubectl --context do-nyc3-omos-cloud get gatewayclasses
```

See the proposed [operating structure and queue](../architecture/kubernetes-operating-structure.md)
for how these findings become small, gated implementation PRs.

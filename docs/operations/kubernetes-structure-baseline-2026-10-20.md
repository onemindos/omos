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

Backstage's deployed image remains:
`ghcr.io/onemindos/onemind-portal@sha256:4a52c3b840fcd2b4e2fc0287929e1704721f7420018f2c09deb8543a6aab74da`.
Source in #330 is not in that image merely because the PR exists.

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

- Phase 6 portal/API semantic gates and #330 source → image → GitOps → live
  verification. Reported STAC database/schema and Grafana adapter issues require
  their own current verification/repair; not silently recertified here.
- Legacy inference failure; complete identity, database and recovery gates before
  activating the staged LibreChat surface.
- Earlier documented Helm-CLI/GitOps exceptions and independent backup gaps must
  be rechecked before adoption/cutover; this sweep did not read release Secrets.
- Mac production VM cannot responsibly be added with only 13–14 GiB internal
  free space. External data/Time Machine drives are not approved for repurposing.
  Closed-lid/reboot/FileVault and local outage recovery remain untested.
- GitHub Quality Gate run `39458922274` could not start its jobs: check annotations
  report failed recent account payments or a spending-limit restriction. No billing
  change was attempted. Local `make validate` still fails on existing YAML debt
  beginning `platform/onetak/media-mtx.yaml:86`. Neither is a passing CI result.

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

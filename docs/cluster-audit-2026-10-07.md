# Cluster ownership and capacity audit

Observed 2026-10-07 against `do-nyc3-omos-cloud`. This is an initial controller inventory, not a completed audit of every Kubernetes resource, secret, persistent volume, external dependency, or backup. No workloads were removed or resized.

## Findings

- Three Ready nodes; 130 Deployment/StatefulSet/DaemonSet controllers; 19 Argo applications; 11 Helm releases.
- All Argo applications Synced; platform-takserver Degraded. Its owned workload health requires investigation; do not infer an outage solely from the application summary.
- Node memory snapshot: 14350 MiB (105%), 10387 MiB (76%), 11384 MiB (83%) of allocatable memory. CPU 15%, 21%, 12%. These are point-in-time measurements, not a sizing forecast. All nodes report MemoryPressure=False.
- Summed container working sets are materially lower than node working sets. Investigate host/system consumption, accounting, and time alignment before attributing all memory to application pods.
- Largest pod working sets: Hermes 2426 MiB, OpenClaw 1678 MiB, WebODM 1653 MiB, FedHub 1324 MiB. Size alone is not justification to remove them.
- Ten controllers have at least one container without a memory request, primarily system components and reloaders. Assess each owner before changing requests.
- Warning events: DigitalOcean telemetry config reloader cannot retrieve imagePullSecret onemind-hub. Check provider-managed ownership before repairing.
- Terraform default node size is s-4vcpu-8gb; observed node allocatable CPU is 7880m and memory approximately 13.3 GiB. Defaults are not evidence of applied configuration. Locate effective inputs/state and confirm DO resource identities before any Terraform apply.
- Terraform remote state backend is not declared in checked-in configuration. Establish actual state location and backup/locking before infrastructure changes.
- Argo sources used the old omos-infra repository name and HEAD. Both URLs resolved to main commit 86e4fa72 during preflight. Updated Git sources to canonical omos URL and explicit main; chart versions preserved.
- Several Applications ignore first-container resource differences and template annotations. Synced therefore does not establish full configuration equality. Narrow ignores only after reviewing the live differences.

## Controller ownership inventory

Argo ownership below uses live Application status.resources, not labels alone. Operator-generated resources are distinguished from controllers without an observed owner. Helm metadata identifies a release but does not establish that current Git values reproduce it.

| Namespace | Kind | Name | Argo inventory | Helm release | Controller owner |
|---|---|---|---|---|---|
| argocd | Deployment | argocd-applicationset-controller | — | argocd | — |
| argocd | Deployment | argocd-dex-server | — | argocd | — |
| argocd | Deployment | argocd-notifications-controller | — | argocd | — |
| argocd | Deployment | argocd-redis | — | argocd | — |
| argocd | Deployment | argocd-repo-server | — | argocd | — |
| argocd | Deployment | argocd-server | — | argocd | — |
| argocd | StatefulSet | argocd-application-controller | — | argocd | — |
| auth | Deployment | authentik-redis | platform-auth-infra | — | — |
| auth | Deployment | platform-auth-helm-authentik-server | platform-auth-helm | — | — |
| auth | Deployment | platform-auth-helm-authentik-worker | platform-auth-helm | — | — |
| auth | StatefulSet | authentik-postgres | platform-auth-infra | — | — |
| backstage | Deployment | backstage | platform-backstage | — | — |
| backstage | StatefulSet | backstage-postgres | platform-backstage | — | — |
| cert-manager | Deployment | cert-manager | — | cert-manager | — |
| cert-manager | Deployment | cert-manager-cainjector | — | cert-manager | — |
| cert-manager | Deployment | cert-manager-webhook | — | cert-manager | — |
| clickhouse | StatefulSet | chi-onemind-onemind-0-0 | — | — | ClickHouseInstallation/onemind |
| cloudflare | Deployment | cloudflared | platform-core | — | — |
| github-runner | Deployment | github-runner | platform-core | — | — |
| hermes | Deployment | hermes-legacy | platform-legacy | — | — |
| ingress-nginx | Deployment | ingress-nginx-controller | — | ingress-nginx | — |
| kube | Deployment | headlamp | platform-kube | — | — |
| kube | Deployment | kubernetes-dashboard | platform-kube | — | — |
| kube-system | DaemonSet | cilium | — | — | — |
| kube-system | DaemonSet | cpc-bridge-proxy-ebpf | — | — | — |
| kube-system | DaemonSet | csi-do-node | — | — | — |
| kube-system | DaemonSet | do-node-agent | — | — | — |
| kube-system | DaemonSet | do-node-agent-amd-device-metrics-exporter | — | — | — |
| kube-system | DaemonSet | do-node-agent-nvidia-dcgm-exporter | — | — | — |
| kube-system | DaemonSet | doks-telemetry-config-reloader | — | — | — |
| kube-system | Deployment | clickhouse-operator | — | — | — |
| kube-system | Deployment | coredns | — | — | — |
| kube-system | Deployment | hubble-relay | — | — | — |
| kube-system | Deployment | hubble-ui | — | — | — |
| kube-system | Deployment | konnectivity-agent | — | — | — |
| kube-system | Deployment | metrics-server | — | — | — |
| kube-system | Deployment | platform-clickhouse-operator | platform-clickhouse-operator | platform-clickhouse-operator | — |
| kube-system | Deployment | sealed-secrets-controller | platform-core | — | — |
| matrix | Deployment | element-web | — | — | — |
| matrix | Deployment | synapse | — | — | — |
| matrix | StatefulSet | postgres | — | — | — |
| matrix | StatefulSet | redis | — | — | — |
| monitoring | DaemonSet | monitoring-prometheus-node-exporter | — | monitoring | — |
| monitoring | DaemonSet | otel-collector-opentelemetry-collector-agent | — | otel-collector | — |
| monitoring | Deployment | loki-gateway | — | loki | — |
| monitoring | Deployment | monitoring-grafana | — | monitoring | — |
| monitoring | Deployment | monitoring-kube-prometheus-operator | — | monitoring | — |
| monitoring | Deployment | monitoring-kube-state-metrics | — | monitoring | — |
| monitoring | StatefulSet | alertmanager-monitoring-kube-prometheus-alertmanager | — | monitoring | Alertmanager/monitoring-kube-prometheus-alertmanager |
| monitoring | StatefulSet | loki | — | loki | — |
| monitoring | StatefulSet | prometheus-monitoring-kube-prometheus-prometheus | — | monitoring | Prometheus/monitoring-kube-prometheus-prometheus |
| monitoring | StatefulSet | tempo | — | tempo | — |
| nats | Deployment | agent-web-ui | — | — | — |
| nats | Deployment | agenttap | — | — | — |
| nats | Deployment | heynats | — | — | — |
| nats | Deployment | nack | — | nack | — |
| nats | Deployment | nats-dashboard | — | — | — |
| nats | Deployment | nats-metrics-scrape | platform-monitoring | — | — |
| nats | Deployment | nats-surveyor | platform-monitoring | — | — |
| nats | Deployment | nats-ui | — | — | — |
| nats | StatefulSet | nats | — | nats | — |
| omos | Deployment | feed-intel | platform-omos | — | — |
| omos | Deployment | feed-satellites | platform-omos | — | — |
| omos | Deployment | omos-mission-engine | platform-omos | — | — |
| omos | Deployment | omos-tak-projector | platform-omos | — | — |
| omos | Deployment | redpanda-connect | platform-redpanda-connect | — | — |
| omos | Deployment | sink-clickhouse | platform-omos | — | — |
| onetak | Deployment | maputnik | platform-onetak | — | — |
| onetak | Deployment | martin | platform-onetak | — | — |
| onetak | Deployment | nodeodm | platform-onetak | — | — |
| onetak | Deployment | nominatim | platform-onetak | — | — |
| onetak | Deployment | onetak-geo-db | platform-onetak | — | — |
| onetak | Deployment | onetak-media | platform-onetak | — | — |
| onetak | Deployment | onetak-minio | platform-onetak | — | — |
| onetak | Deployment | onetak-postgis | platform-onetak | — | — |
| onetak | Deployment | overpass | platform-onetak | — | — |
| onetak | Deployment | pg-featureserv | platform-onetak | — | — |
| onetak | Deployment | stac-api | platform-onetak | — | — |
| onetak | Deployment | stac-browser | platform-onetak | — | — |
| onetak | Deployment | tak-portal | platform-onetak | — | — |
| onetak | Deployment | tile38 | platform-onetak | — | — |
| onetak | Deployment | tileserver-gl | platform-onetak | — | — |
| onetak | Deployment | valhalla | platform-onetak | — | — |
| onetak | Deployment | webodm | platform-onetak | — | — |
| onetak | Deployment | webodm-broker | platform-onetak | — | — |
| onetak | Deployment | webodm-db | platform-onetak | — | — |
| openclaw | Deployment | openclaw | — | — | — |
| tailscale | Deployment | operator | — | tailscale-operator | — |
| tailscale | StatefulSet | egress | — | — | ProxyGroup/egress |
| tailscale | StatefulSet | ts-agent-web-ui-l74d4 | — | — | — |
| tailscale | StatefulSet | ts-agenttap-8ddc6 | — | — | — |
| tailscale | StatefulSet | ts-argocd-slkzh | — | — | — |
| tailscale | StatefulSet | ts-grafana-ksprn | — | — | — |
| tailscale | StatefulSet | ts-headlamp-4qwkc | — | — | — |
| tailscale | StatefulSet | ts-hermes-dashboard-6thbx | — | — | — |
| tailscale | StatefulSet | ts-heynats-p5dmf | — | — | — |
| tailscale | StatefulSet | ts-hubble-tw85d | — | — | — |
| tailscale | StatefulSet | ts-kubernetes-dashboard-9r9fp | — | — | — |
| tailscale | StatefulSet | ts-nats-dashboard-242wv | — | — | — |
| tailscale | StatefulSet | ts-nats-ui-d6rrg | — | — | — |
| tailscale | StatefulSet | ts-omos-dell-edge-qfrh9 | — | — | — |
| tailscale | StatefulSet | ts-omos-dell-zeus-wsv6m | — | — | — |
| tailscale | StatefulSet | ts-omos-ha-edge-rlv6k | — | — | — |
| tailscale | StatefulSet | ts-omos-intel-zzn2k | — | — | — |
| tailscale | StatefulSet | ts-omos-mbp-hgq9w | — | — | — |
| tailscale | StatefulSet | ts-omos-nats-fs6r4 | — | — | — |
| tailscale | StatefulSet | ts-omos-s24-44l4v | — | — | — |
| tailscale | StatefulSet | ts-onemind-auth-4w5hk | — | — | — |
| tailscale | StatefulSet | ts-onemind-preview-shtd7 | — | — | — |
| tailscale | StatefulSet | ts-openclaw-gateway-s7hcm | — | — | — |
| takml | Deployment | takml-db | platform-takml | — | — |
| takml | Deployment | takml-server | platform-takml | — | — |
| takml | Deployment | takml-vision-worker | platform-takml | — | — |
| takserver | Deployment | fedhub | platform-fedhub | — | — |
| takserver | Deployment | mumble-bridge | platform-takserver | — | — |
| takserver | Deployment | murmur | platform-takserver | — | — |
| takserver | Deployment | tak-lb-proxy | platform-takserver | — | — |
| takserver | Deployment | tak-nats-leaf | platform-takserver | — | — |
| takserver | Deployment | takserver-api | — | takserver | — |
| takserver | Deployment | takserver-config | — | takserver | — |
| takserver | Deployment | takserver-messaging | — | takserver | — |
| takserver | Deployment | takserver-nats-box | — | takserver | — |
| takserver | Deployment | takserver-plugins | platform-takserver | takserver | — |
| takserver | StatefulSet | fedhub-mongodb | platform-fedhub | — | — |
| takserver | StatefulSet | tak-ldap | platform-takserver | — | — |
| takserver | StatefulSet | takserver-ignite | — | takserver | — |
| takserver | StatefulSet | takserver-nats | — | takserver | — |
| takserver | StatefulSet | takserver-postgresql | — | takserver | — |
| zeus-site | Deployment | zeus-site | platform-zeus-site | — | — |
| zeus-site | Deployment | zeus-site-postgres | platform-zeus-site | — | — |

## Helm release inventory

| Namespace | Release | Chart | Revision |
|---|---|---|---|
| argocd | argocd | argo-cd-10.1.2 | 3 |
| cert-manager | cert-manager | cert-manager-v1.20.3 | 5 |
| ingress-nginx | ingress-nginx | ingress-nginx-4.15.1 | 1 |
| monitoring | loki | loki-7.0.0 | 4 |
| monitoring | monitoring | kube-prometheus-stack-87.16.1 | 8 |
| nats | nack | nack-0.34.0 | 4 |
| nats | nats | nats-2.14.2 | 10 |
| monitoring | otel-collector | opentelemetry-collector-0.165.0 | 3 |
| tailscale | tailscale-operator | tailscale-operator-1.98.4 | 3 |
| takserver | takserver | takserver-4.6.0 | 34 |
| monitoring | tempo | tempo-1.24.4 | 4 |

## Next work, in order

1. Collect memory history, node/system memory breakdown, requests, limits, and scheduling distribution. Avoid lowering limits from a single snapshot.
2. Trace Matrix (four controllers), five NATS tool deployments, and OpenClaw. These have no observed direct Argo inventory, Helm release metadata, or controller owner. Matrix/NATS tools have repository manifests; OpenClaw workload source still needs locating.
3. Reproduce every manual Helm release from pinned charts and sanitized Git values before transferring ownership to Argo. Do not adopt resources with automatic pruning until reviewed.
4. Explain TAK application health and inspect any Helm/Argo overlap.
5. Inventory persistent volumes, recovery procedures, sealed-secret recovery keys, network exposure, and service dependencies.
6. Prepare a bounded Netdata pilot after measuring available memory. Keep existing monitoring during evaluation.
7. Present removals individually for user confirmation.

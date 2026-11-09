# Authentik login inventory baseline — 2026-10-23

Read-only observation: 2026-10-23 20:44:45 EDT / 2026-10-24 00:44:45 UTC.
`kubectl` context: `do-nyc3-omos-cloud`.
Git source inspected: `08c65b58b8ade541415559311405c4f329c1ffb1` (`origin/main`).

This is a narrow workload/route/Argo snapshot for the
[shared-login plan](../roadmap/authentik-unified-login.md), not a browser-login,
protocol, licensing, data-integrity or cluster-wide health certification.
Deployment/StatefulSet readiness and ingress configuration do not prove an app
is reachable, usable or signed in. No Secret data, provider responses, application
databases or user sessions were read. No cluster mutation was performed for this
baseline. Earlier troubleshooting in this chat is not acceptance evidence.

## Selected application workloads

| Namespace | Workload | Kind | Ready / desired |
|---|---|---|---|
| `agent-fabric` | `agent-console` | Deployment | 1/1 |
| `agent-fabric` | `agent-tap` | Deployment | 1/1 |
| `argocd` | `argocd-server` | Deployment | 1/1 |
| `auth` | `authentik-redis` | Deployment | 1/1 |
| `auth` | `platform-auth-helm-authentik-server` | Deployment | 2/2 |
| `auth` | `platform-auth-helm-authentik-worker` | Deployment | 1/1 |
| `backstage` | `backstage` | Deployment | 2/2 |
| `dp-default-onetak-development-b188bf98` | `stac-browser-development-0c6b80f7` | Deployment | 1/1 |
| `hermes` | `hermes-legacy` | Deployment | 1/1 |
| `kube` | `headlamp` | Deployment | 1/1 |
| `kube` | `kubernetes-dashboard` | Deployment | 1/1 |
| `matrix` | `element-web` | Deployment | 3/3 |
| `matrix` | `synapse` | Deployment | 1/1 |
| `monitoring` | `monitoring-grafana` | Deployment | 1/1 |
| `n8n` | `n8n` | Deployment | 1/1 |
| `nats` | `heynats` | Deployment | 1/1 |
| `nats` | `nats-dashboard` | Deployment | 1/1 |
| `nats` | `nats-surveyor` | Deployment | 1/1 |
| `nats` | `nats-ui` | Deployment | 1/1 |
| `nextcloud` | `nextcloud` | Deployment | 1/1 |
| `nextcloud` | `nextcloud-redis` | Deployment | 1/1 |
| `onemind-ai` | `litellm` | Deployment | 1/1 |
| `onemind-core` | `onemind-core` | Deployment | 2/2 |
| `onemind-core` | `onemind-ingestion` | Deployment | 1/1 |
| `onemind-core` | `openfga` | Deployment | 1/1 |
| `onetak` | `nodeodm` | Deployment | 1/1 |
| `onetak` | `onetak-media` | Deployment | 1/1 |
| `onetak` | `tak-portal` | Deployment | 1/1 |
| `onetak` | `tile38` | Deployment | 1/1 |
| `onetak` | `tileserver-gl` | Deployment | 1/1 |
| `onetak` | `webodm` | Deployment | 1/1 |
| `openchoreo-control-plane` | `openchoreo-api` | Deployment | 1/1 |
| `opencost` | `opencost` | Deployment | 1/1 |
| `takml` | `takml-server` | Deployment | 1/1 |
| `takserver` | `fedhub` | Deployment | 1/1 |
| `auth` | `authentik-postgres` | StatefulSet | 1/1 |
| `clickhouse` | `chi-onemind-onemind-0-0` | StatefulSet | 1/1 |
| `nats` | `nats` | StatefulSet | 3/3 |
| `onemind-search` | `onemind-search-master` | StatefulSet | 1/1 |

No Paperclip or Langfuse Deployment/StatefulSet was returned in the all-namespace
query. This does not prove absence of every possible Job/external service.
OpenChoreo/STAC and OpenSearch were still returned despite the approved retirement
direction. Do not resurrect them or assume the retirement has completed.

## Configured browser and protocol entry points

Tailnet names below come from Ingress TLS hosts; `.tail717ba0.ts.net` is the
existing tailnet suffix in repository configuration. This is configured routing,
not a DNS/TLS/HTTP probe. Empty public rule lists are normal for these private
Ingress objects. Service-backed listeners, internal UIs and generated HTTPRoutes
are separate from Ingress inventory.

| Namespace | Ingress | Class | Configured host/name |
|---|---|---|---|
| `agent-fabric` | `agent-console` | `tailscale` | `agent-console` |
| `agent-fabric` | `agent-tap` | `tailscale` | `agent-tap` |
| `argocd` | `argocd` | `tailscale` | `argocd` |
| `auth` | `authentik-ingress` | `nginx` | `auth.onemindos.dev` |
| `auth` | `onemind-auth` | `tailscale` | `onemind-auth` |
| `backstage` | `onemind-preview` | `tailscale` | `onemind-preview` |
| `hermes` | `hermes-dashboard` | `tailscale` | `hermes-agent` |
| `kube-system` | `hubble` | `tailscale` | `hubble` |
| `kube` | `headlamp` | `tailscale` | `headlamp` |
| `kube` | `kubernetes-dashboard` | `tailscale` | `k8s-dashboard` |
| `matrix` | `matrix-ingress` | `nginx` | `matrix.onemindos.dev`, `element.onemindos.dev` |
| `matrix` | `matrix-wellknown` | `nginx` | `onemindos.dev` |
| `monitoring` | `grafana` | `tailscale` | `grafana` |
| `n8n` | `n8n` | `tailscale` | `n8n` |
| `nats` | `heynats` | `tailscale` | `heynats` |
| `nats` | `nats-dashboard` | `tailscale` | `nats-dashboard` |
| `nats` | `nats-monitor` | `tailscale` | `nats-monitor` |
| `nats` | `nats-ui` | `tailscale` | `nats-ui` |
| `nats` | `nats-ws` | `tailscale` | `nats-ws` |
| `nats` | `nats-ws-ingress` | `nginx` | `ws.onemindos.dev` |
| `nextcloud` | `nextcloud` | `tailscale` | `onemind-files` |
| `onetak` | `geo-services-ingress` | `nginx` | `tileserver.onemindos.dev` |
| `onetak` | `onetak-ingress` | `nginx` | `onemindos.dev` |
| `onetak` | `tak-portal-ingress` | `tailscale` | `enroll` |
| `opencost` | `opencost` | `tailscale` | `opencost` |

FedHub is a Service-backed tailnet TLS listener, not an Ingress: source declares
`fedhub.tail717ba0.ts.net:9443` for OAuth and `:9100` for client certificates.
TAK/Mumble and NATS L4 listeners are protocol endpoints, not browser SSO targets.
Public root, tiles and Matrix routes must retain their deliberate per-path
behavior; do not protect every path by adding a blanket browser redirect.

### Generated HTTPRoutes still observed

| Namespace | HTTPRoute | Configured hostname |
|---|---|---|
| `dp-default-onetak-development-b188bf98` | `stac-browser-http-02eaf4a9` | `development-default.apps.openchoreo.100-110-185-30.nip.io` |
| `openchoreo-control-plane` | `openchoreo-api` | `onemind-openchoreo-api.tail717ba0.ts.net` |
| `openchoreo-control-plane` | `portal-assistant` | not specified |

## Identity rollout and ownership findings

| Auth bootstrap Job | Succeeded | Failed pods | Conditions |
|---|---|---|---|
| `argocd-authentik-oidc-v2` | 0 | 4 | FailureTarget=True (BackoffLimitExceeded), Failed=True (BackoffLimitExceeded) |
| `onemind-platform-oidc-v1` | 0 | 1 | FailureTarget=True (BackoffLimitExceeded), Failed=True (BackoffLimitExceeded) |

| Argo Application | Sync | Health | Compared revision | Operation phase / revision |
|---|---|---|---|---|
| `app-of-apps` | OutOfSync | Progressing | `08c65b58b8ade541415559311405c4f329c1ffb1` | Running / `22b86049153004ff4d27fffd3a1182137cc21e9d` |
| `ecosystem-external-secrets` | Unknown | Healthy | `not reported` | Succeeded / `not reported` |
| `ecosystem-oc-control` | Unknown | Healthy | `not reported` | Succeeded / `not reported` |
| `ecosystem-opensearch` | Unknown | Healthy | `not reported` | Succeeded / `not reported` |
| `platform-argocd` | OutOfSync | Healthy | `not reported` | not reported / `not reported` |
| `platform-auth-helm` | Synced | Healthy | `2026.8.2` | Succeeded / `2026.8.2` |
| `platform-auth-infra` | Synced | Healthy | `08c65b58b8ade541415559311405c4f329c1ffb1` | Succeeded / `0d8c354f45d189ffa39059ae08847fb7ff8ea88d` |
| `platform-backups` | OutOfSync | Healthy | `08c65b58b8ade541415559311405c4f329c1ffb1` | Succeeded / `c95e2470249e6bfad59a9384430321304ed26ce8` |
| `platform-ecosystem-applications` | OutOfSync | Healthy | `08c65b58b8ade541415559311405c4f329c1ffb1` | Succeeded / `0d8c354f45d189ffa39059ae08847fb7ff8ea88d` |
| `platform-ecosystem-resources` | OutOfSync | Healthy | `08c65b58b8ade541415559311405c4f329c1ffb1` | Running / `f5f4dca037767c28de3f13a8874b4f66815186eb` |
| `platform-matrix` | OutOfSync | Healthy | `08c65b58b8ade541415559311405c4f329c1ffb1` | not reported / `not reported` |
| `platform-openchoreo-defaults` | Unknown | Progressing | `not reported` | Succeeded / `not reported` |
| `platform-openchoreo-projects` | Unknown | Progressing | `main` | Succeeded / `68a5c612916e8caafa4df945d375dbfc3f7aacb2` |
| `platform-volume-backups` | OutOfSync | Healthy | `08c65b58b8ade541415559311405c4f329c1ffb1` | Succeeded / `c95e2470249e6bfad59a9384430321304ed26ce8` |

In particular, `platform-ecosystem-resources` had a newer comparison revision
but an old running operation. No successful application of #440 or completed
browser SSO is certified by this baseline. Repair/verification is package S00.
`platform-argocd` has no operation reported; do not assume Argo self-adoption
has completed. Matrix controller metadata lacked an Argo tracking annotation
while its Application was OutOfSync; reconcile ownership before changing login.
Missing controller annotations alone do not establish that a resource is unowned.

## Reproduce the observation

These commands read metadata/status; inspect selected fields rather than
printing complete app configuration or Secrets. Record a fresh timestamp.

```bash
kubectl --context do-nyc3-omos-cloud get deployments,statefulsets -A \
  -o custom-columns='NS:.metadata.namespace,NAME:.metadata.name,READY:.status.readyReplicas,DESIRED:.spec.replicas'
kubectl --context do-nyc3-omos-cloud get ingress -A
kubectl --context do-nyc3-omos-cloud get httproutes.gateway.networking.k8s.io -A
kubectl --context do-nyc3-omos-cloud -n auth get jobs
kubectl --context do-nyc3-omos-cloud -n argocd get applications
kubectl --context do-nyc3-omos-cloud -n argocd get app platform-ecosystem-resources \
  -o jsonpath='{.status.sync.revision}{" compared; operation "}{.status.operationState.syncResult.revision}{"\n"}'
```

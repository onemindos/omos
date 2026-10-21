# OneMind catalog system map (F2)

Status: **proposed, awaiting Zeus's approval** (flagship work order F2, PR #358). Evidence date: 2026-10-21.
Sources: `catalog/apps.yaml` (154 entities), `catalog/cluster.yaml` (514 entities), and read-only
`kubectl --context do-nyc3-omos-cloud -n argocd get applications` (69 applications) and `get deployments,statefulsets -A` (176 live workloads;
the catalog also retains 10 workloads that are no longer observed, marked `missing-review-required`).

Machine-readable form: `services/onemind-portal/catalog/structure.yaml` (Systems, pillars, placement) and
`services/onemind-portal/catalog/workload-map.yaml` (workloads and Argo applications to apps). This page is generated from the same data.

## Decisions for Zeus

1. **Approve the target Systems and pillar column below.** 22 Systems (limit 25), each placed in a Domain. Pillars are `proposed`. Systems marked Unclassified sit under the parent domain `OneMind` until you decide.
2. **Seven Systems are Unclassified on purpose:** Intelligence & feeds (proposed SO: feeds that serve TAK and the bus), Communication (shared), Productivity & SaaS (Zoho and Google Workspace may be business GE or personal), Personal tools & devices (mixed), Brand & social media (proposed GE per the pillar doctrine "brand", but several accounts are personal), Reference & community, Unassigned workloads. Say which you want assigned.
3. **GE, LE and HP are thin.** Only Business operations (GE: GoHighLevel, Circle), Home & homestead (LE: Home Assistant, FarmBot) and Personal finance (HP) have content today. There are no ventures, holdings, health or family systems in the catalog yet; none are invented here.
4. **Mercury** (bank) sits in Personal finance (HP) but its pillar is left Unclassified: it may be a business account (GE). Credit bureaus and credit-building apps are HP.
5. **Several Argo applications share one app** (for example 12 OpenChoreo applications, 4 Argo CD applications). The Argo CD plugin reads one `argocd/app-name` per Component, so the app carries the primary application and lists the rest in `onemind.io/argocd-apps`. Approve this, or ask for one Component per Argo application.
6. **Removed Systems keep a redirect entry** (`catalog/redirects.yaml`, 74 current Systems, 61 get a redirect). Backstage has no built-in entity redirect, so the file and its test exist now, but making the old URLs forward needs a small route in the app (outside this work order's files). Say if you want it added here or in F1.
7. **Resource type rename deferred.** The consistent set (`database`, `volume`, `bucket`, `cluster`, `network`, `code-repository`, `capability`, `endpoint`) needs edits to `plugin-onemind-common/AccessDeployments.tsx`, `report-service-coverage.cjs` and `scripts/tests/test_gitops_topology.py`, which other work orders touch. Approve doing it in a follow-up.
8. **Duplicate launcher entries** (`external-argocd`, `external-kubernetes-dashboard`, `external-tailscale`) are kept and retitled so titles are unique; they look like merge candidates.
9. **Ollama.** `component:default/ollama` is documented as planned/local, but a cluster Ollama is deployed (`onemind-ai/onemind-cpu-model`, Argo `ecosystem-ollama`). The workload is mapped to that entity; its lifecycle text is not changed.
10. **`onetak-stac-readonly`** pointed at `system:default/onetak`, which never existed. It moves to Geospatial.

## Target Systems

| System | Name (ref) | Domain / pillar | Contents |
|---|---|---|---|
| Platform core | `system:default/infrastructure` | SO | Cluster foundations: networking and ingress, certificates, DNS, autoscaling, operators, Tailscale and Cloudflare, and the DigitalOcean account. |
| Identity & access | `system:default/identity-access` | SO | Sign-in, authorization, workload identity, secrets management and admission/runtime policy. |
| GitOps & delivery | `system:default/gitops-delivery` | SO | Argo CD and Argo Rollouts, the container registry, the GitHub runner and source hosting. |
| Developer platform | `system:default/developer-platform` | SO | The Backstage portal, OpenChoreo control and data planes, cluster dashboards and runbooks. |
| Observability | `system:default/observability` | SO | Metrics, logs, traces, dashboards, alerting, cost and cluster-health tooling. |
| Data & storage | `system:default/data-storage` | SO | Databases and operators, object and block storage, search, backups and recovery. |
| Messaging (NATS) | `system:default/messaging-nats` | SO | The NATS bus, its UIs and tools, and the data pipelines that publish to or consume from it. |
| TAK & field ops | `system:default/tak-ecosystem` | SO | TAK Server and its plugins, OneTAK, FedHub, TAK-ML, voice and the TAK clients and references. |
| Geospatial, drones & robotics | `system:default/geospatial` | SO | Map, routing, geocoding and imagery services, drone software and firmware, and robotics platforms. |
| AI & agents | `system:default/ai-agents` | SO | Hermes (Legacy), agent consoles, model gateways and runtimes, model providers and workflow automation. |
| Code repositories | `system:default/repo-onemind` | SO | Source repositories in the onemindos GitHub organization. Hidden by default. |
| Mesh & devices | `system:default/onemind-mesh` | SO | Tailscale mesh endpoints (devices and service proxies). Discovery snapshot, not live health. Hidden by default. |
| Personal finance | `system:default/personal-finance` | HP | Personal banking, credit monitoring and credit-building tools. |
| Home & homestead | `system:default/home-homestead` | LE | Home automation and homestead/farm automation. |
| Business operations | `system:default/business-operations` | GE | Customer relationship, community and other business-operations platforms. |
| Intelligence & feeds | `system:default/intelligence` | OneMind (**Unclassified**) | Intelligence feeds, monitors and the observational API. Pillar not yet reviewed. |
| Communication | `system:default/communication` | OneMind (**Unclassified**) | Messaging, mail and phone services, self-hosted and external. Shared across missions; pillar not yet reviewed. |
| Productivity & SaaS | `system:default/productivity-saas` | OneMind (**Unclassified**) | Zoho, Google Workspace, Notion, Dropbox and other hosted productivity services. Mixed use; pillar not yet reviewed. |
| Personal tools & devices | `system:default/personal-tools` | OneMind (**Unclassified**) | Desktop applications, hardware vendors and personal devices. Mixed use; pillar not yet reviewed. |
| Brand & social media | `system:default/social-media` | OneMind (**Unclassified**) | Social and brand accounts. Pillar not yet reviewed (GE is the likely home). |
| Reference & community | `system:default/reference-links` | OneMind (**Unclassified**) | Foundations, standards bodies and community reference sites. Hidden by default. |
| Unassigned workloads | `system:default/unassigned-workloads` | OneMind (**Unclassified**) | Kubernetes workloads that no reviewed mapping assigns to an app yet. Nothing is hidden: new workloads land here until mapped. |

Domains: `domain:default/onemind` (parent) with children `onemind-so`, `onemind-le`, `onemind-hp`, `onemind-ge`, owner `group:default/onemind-admins`.
Names are kept where an existing System survives (so its URL and favorites still work); every other old System gets a redirect entry.

## Every current System to its target

| Current System | Members | Target System | Pillar | Reason |
|---|---|---|---|---|
| `infrastructure` (Infrastructure) | 14 | Platform core | SO | Name kept; retitled "Platform core". Its members are split by function (see the app table). |
| `ai-agents` (AI & Agents) | 6 | AI & agents | SO | Name kept. |
| `communication` (Communication) | 5 | Communication | **Unclassified** | Name kept; now also holds mail and phone apps. |
| `tak-ecosystem` (TAK Ecosystem) | 8 | TAK & field ops | SO | Name kept; retitled "TAK & field ops". |
| `intelligence` (Intelligence) | 4 | Intelligence & feeds | **Unclassified** | Name kept; retitled "Intelligence & feeds". Grafana, NATS UI, Agent Monitor and ClickHouse move to the Systems that own them. |
| `data-storage` (Data & Storage) | 3 | Data & storage | SO | Name kept. |
| `observability` (Observability) | 6 | Observability | SO | Name kept. |
| `geospatial` (Geospatial) | 9 | Geospatial, drones & robotics | SO | Name kept; retitled to include drones and robotics. |
| `zoho-apps` (Zoho Apps) | 12 | Productivity & SaaS | **Unclassified** | Zoho is one of several hosted productivity suites; merged into Productivity & SaaS. |
| `personal-tools` (Personal Tools) | 18 | Personal tools & devices | **Unclassified** | Name kept. |
| `workspaces` (Workspaces) | 6 | Productivity & SaaS | **Unclassified** | Mixed SaaS: GitHub moves to GitOps & delivery, Circle to Business operations, Composio to AI & agents, the rest to Productivity & SaaS. |
| `personal-finance` (Personal Finance) | 6 | Personal finance | HP | Name kept. |
| `external` (External) | 4 | Developer platform | SO | Duplicate launcher links ("external-*") and OneMind OS; each now sits with the service it points at. |
| `social-media` (Social Media) | 12 | Brand & social media | **Unclassified** | Name kept; retitled "Brand & social media". |
| `cluster-agent-fabric` (Namespace / agent-fabric) | 2 | AI & agents | SO | Namespace grouping dissolved; agent-fabric workloads attach to their app, whose System is AI & agents. |
| `cluster-argocd` (Namespace / argocd) | 8 | GitOps & delivery | SO | Namespace grouping dissolved; argocd workloads attach to their app, whose System is GitOps & delivery. |
| `cluster-auth` (Namespace / auth) | 4 | Identity & access | SO | Namespace grouping dissolved; auth workloads attach to their app, whose System is Identity & access. |
| `cluster-backstage` (Namespace / backstage) | 3 | Developer platform | SO | Namespace grouping dissolved; backstage workloads attach to their app, whose System is Developer platform. |
| `cluster-cert-manager` (Namespace / cert-manager) | 3 | Platform core | SO | Namespace grouping dissolved; cert-manager workloads attach to their app, whose System is Platform core. |
| `cluster-cloudflare` (Namespace / cloudflare) | 2 | Platform core | SO | Namespace grouping dissolved; cloudflare workloads attach to their app, whose System is Platform core. |
| `cluster-cnpg-system` (Namespace / cnpg-system) | 2 | Data & storage | SO | Namespace grouping dissolved; cnpg-system workloads attach to their app, whose System is Data & storage. |
| `cluster-dapr-system` (Namespace / dapr-system) | 5 | Platform core | SO | Namespace grouping dissolved; dapr-system workloads attach to their app, whose System is Platform core. |
| `cluster-dp-default-onetak-development-b188bf98` (Namespace / dp-default-onetak-development-b188bf98) | 1 | Geospatial, drones & robotics | SO | Namespace grouping dissolved; dp-default-onetak-development-b188bf98 workloads attach to their app, whose System is Geospatial, drones & robotics. |
| `cluster-external-dns` (Namespace / external-dns) | 1 | Platform core | SO | Namespace grouping dissolved; external-dns workloads attach to their app, whose System is Platform core. |
| `cluster-external-secrets` (Namespace / external-secrets) | 3 | Identity & access | SO | Namespace grouping dissolved; external-secrets workloads attach to their app, whose System is Identity & access. |
| `cluster-github-runner` (Namespace / github-runner) | 2 | Platform core | SO | Namespace grouping dissolved; github-runner workloads attach to their app, whose System is Platform core. |
| `cluster-hermes` (Namespace / hermes) | 1 | AI & agents | SO | Namespace grouping dissolved; hermes workloads attach to their app, whose System is AI & agents. |
| `cluster-ingress-nginx` (Namespace / ingress-nginx) | 1 | Platform core | SO | Namespace grouping dissolved; ingress-nginx workloads attach to their app, whose System is Platform core. |
| `cluster-keda` (Namespace / keda) | 3 | Platform core | SO | Namespace grouping dissolved; keda workloads attach to their app, whose System is Platform core. |
| `cluster-kube-system` (Namespace / kube-system) | 8 | Platform core | SO | Namespace grouping dissolved; kube-system workloads attach to their app, whose System is Platform core. |
| `cluster-kube` (Namespace / kube) | 2 | Developer platform | SO | Namespace grouping dissolved; kube workloads attach to their app, whose System is Developer platform. |
| `cluster-kyverno` (Namespace / kyverno) | 4 | Identity & access | SO | Namespace grouping dissolved; kyverno workloads attach to their app, whose System is Identity & access. |
| `cluster-matrix` (Namespace / matrix) | 4 | Communication | **Unclassified** | Namespace grouping dissolved; matrix workloads attach to their app, whose System is Communication. |
| `cluster-monitoring` (Namespace / monitoring) | 9 | Observability | SO | Namespace grouping dissolved; monitoring workloads attach to their app, whose System is Observability. |
| `cluster-n8n` (Namespace / n8n) | 2 | AI & agents | SO | Namespace grouping dissolved; n8n workloads attach to their app, whose System is AI & agents. |
| `cluster-nats` (Namespace / nats) | 12 | Messaging (NATS) | SO | Namespace grouping dissolved; nats workloads attach to their app, whose System is Messaging (NATS). |
| `cluster-omos` (Namespace / omos) | 6 | Messaging (NATS) | SO | Namespace grouping dissolved; omos workloads attach to their app, whose System is Messaging (NATS). |
| `cluster-onemind-ai` (Namespace / onemind-ai) | 2 | AI & agents | SO | Namespace grouping dissolved; onemind-ai workloads attach to their app, whose System is AI & agents. |
| `cluster-onemind-core` (Namespace / onemind-core) | 3 | Platform core | SO | Namespace grouping dissolved; onemind-core workloads attach to their app, whose System is Platform core. |
| `cluster-onemind-registry` (Namespace / onemind-registry) | 8 | GitOps & delivery | SO | Namespace grouping dissolved; onemind-registry workloads attach to their app, whose System is GitOps & delivery. |
| `cluster-onetak` (Namespace / onetak) | 19 | TAK & field ops | SO | Namespace grouping dissolved; onetak workloads attach to their app, whose System is TAK & field ops. |
| `cluster-openchoreo-control-plane` (Namespace / openchoreo-control-plane) | 4 | Developer platform | SO | Namespace grouping dissolved; openchoreo-control-plane workloads attach to their app, whose System is Developer platform. |
| `cluster-openchoreo-data-plane` (Namespace / openchoreo-data-plane) | 1 | Developer platform | SO | Namespace grouping dissolved; openchoreo-data-plane workloads attach to their app, whose System is Developer platform. |
| `cluster-openchoreo-observability-plane` (Namespace / openchoreo-observability-plane) | 8 | Observability | SO | Namespace grouping dissolved; openchoreo-observability-plane workloads attach to their app, whose System is Observability. |
| `cluster-openclaw` (Namespace / openclaw) | 1 | AI & agents | SO | Namespace grouping dissolved; openclaw workloads attach to their app, whose System is AI & agents. |
| `cluster-opencost` (Namespace / opencost) | 1 | Observability | SO | Namespace grouping dissolved; opencost workloads attach to their app, whose System is Observability. |
| `cluster-ops` (Namespace / ops) | 1 | Observability | SO | Namespace grouping dissolved; ops workloads attach to their app, whose System is Observability. |
| `cluster-tailscale` (Namespace / tailscale) | 33 | Platform core | SO | Namespace grouping dissolved; tailscale workloads attach to their app, whose System is Platform core. |
| `cluster-takml` (Namespace / takml) | 3 | TAK & field ops | SO | Namespace grouping dissolved; takml workloads attach to their app, whose System is TAK & field ops. |
| `cluster-takserver` (Namespace / takserver) | 15 | TAK & field ops | SO | Namespace grouping dissolved; takserver workloads attach to their app, whose System is TAK & field ops. |
| `cluster-velero` (Namespace / velero) | 1 | Data & storage | SO | Namespace grouping dissolved; velero workloads attach to their app, whose System is Data & storage. |
| `cluster-clickhouse` (Namespace / clickhouse) | 1 | Data & storage | SO | Namespace grouping dissolved; clickhouse workloads attach to their app, whose System is Data & storage. |
| `cluster-onemind-search` (Namespace / onemind-search) | 1 | Data & storage | SO | Namespace grouping dissolved; onemind-search workloads attach to their app, whose System is Data & storage. |
| `cluster-spire-system` (Namespace / spire-system) | 1 | Identity & access | SO | Namespace grouping dissolved; spire-system workloads attach to their app, whose System is Identity & access. |
| `onemind-mesh` (OneMind mesh) | 47 | Mesh & devices | SO | Name kept; retitled "Mesh & devices". |
| `fabric-geo` (Fabric / 🌍 Geo) | 13 | Geospatial, drones & robotics | SO | Notion registry category; its documented capabilities are Resources placed by topic in Geospatial, drones & robotics. |
| `fabric-data` (Fabric / 💾 Data) | 4 | Data & storage | SO | Notion registry category; its documented capabilities are Resources placed by topic in Data & storage. |
| `fabric-brandmarketing` (Fabric / 🎯 Brand & Marketing) | 9 | Brand & social media | **Unclassified** | Notion registry category; its documented capabilities are Resources placed by topic in Brand & social media. |
| `fabric-googleworkspace` (Fabric / 🔵 Google Workspace) | 7 | Productivity & SaaS | **Unclassified** | Notion registry category; its documented capabilities are Resources placed by topic in Productivity & SaaS. |
| `fabric-infrastructure` (Fabric / 🏗️ Infrastructure) | 8 | Platform core | SO | Notion registry category; its documented capabilities are Resources placed by topic in Platform core. |
| `fabric-tak` (Fabric / 🎯 TAK) | 19 | TAK & field ops | SO | Notion registry category; its documented capabilities are Resources placed by topic in TAK & field ops. |
| `fabric-nats` (Fabric / 🚌 NATS) | 12 | Messaging (NATS) | SO | Notion registry category; its documented capabilities are Resources placed by topic in Messaging (NATS). |
| `fabric-intel` (Fabric / 📡 Intel) | 9 | Intelligence & feeds | **Unclassified** | Notion registry category; its documented capabilities are Resources placed by topic in Intelligence & feeds. |
| `fabric-bridge` (Fabric / 🌉 Bridge) | 9 | Messaging (NATS) | SO | Notion registry category; its documented capabilities are Resources placed by topic in Messaging (NATS). |
| `fabric-finance` (Fabric / 💰 Finance) | 3 | Personal finance | HP | Notion registry category; its documented capabilities are Resources placed by topic in Personal finance. |
| `fabric-comms` (Fabric / 📡 Comms) | 3 | Communication | **Unclassified** | Notion registry category; its documented capabilities are Resources placed by topic in Communication. |
| `fabric-plugins` (Fabric / 🔌 Plugins) | 10 | TAK & field ops | SO | Notion registry category; its documented capabilities are Resources placed by topic in TAK & field ops. |
| `fabric-agents` (Fabric / 🧠 Agents) | 13 | AI & agents | SO | Notion registry category; its documented capabilities are Resources placed by topic in AI & agents. |
| `fabric-businessops` (Fabric / 💼 Business Ops) | 2 | Business operations | GE | Notion registry category; its documented capabilities are Resources placed by topic in Business operations. |
| `fabric-devops` (Fabric / ⚙️ DevOps) | 8 | GitOps & delivery | SO | Notion registry category; its documented capabilities are Resources placed by topic in GitOps & delivery. |
| `fabric-observability` (Fabric / 📊 Observability) | 1 | Observability | SO | Notion registry category; its documented capabilities are Resources placed by topic in Observability. |
| `fabric-productivity` (Fabric / 🛠️ Productivity) | 1 | Productivity & SaaS | **Unclassified** | Notion registry category; its documented capabilities are Resources placed by topic in Productivity & SaaS. |
| `fabric-iotautomation` (Fabric / 🤖 IoT / Automation) | 1 | Home & homestead | LE | Notion registry category; its documented capabilities are Resources placed by topic in Home & homestead. |
| `repo-onemind` (OneMind repositories) | 83 | Code repositories | SO | Name kept; retitled "Code repositories". |

Where a current System splits (Platform core, Intelligence, Workspaces, External, Personal tools), members follow the app table below, not the System-level default.

## Argo CD applications (69) to app

Each application resolves to exactly one app Component. The app carries `argocd/app-name` (the first application listed for it) and `onemind.io/argocd-apps` (all of them).

| Argo application | Namespace | App | System |
|---|---|---|---|
| `app-of-apps` | argocd | ArgoCD (`argocd`) | GitOps & delivery |
| `argo-rollouts` | argocd | Argo Rollouts (`argo-rollouts`) | GitOps & delivery |
| `ecosystem-barman` | cnpg-system | CloudNativePG (`cloudnative-pg`) | Data & storage |
| `ecosystem-cnpg` | cnpg-system | CloudNativePG (`cloudnative-pg`) | Data & storage |
| `ecosystem-dapr` | dapr-system | Dapr (`dapr`) | Platform core |
| `ecosystem-external-secrets` | external-secrets | External Secrets Operator (`external-secrets`) | Identity & access |
| `ecosystem-falco` | falco | Falco (`falco`) | Identity & access |
| `ecosystem-harbor` | onemind-registry | Harbor (`harbor`) | GitOps & delivery |
| `ecosystem-keda` | keda | KEDA (`keda`) | Platform core |
| `ecosystem-kyverno` | (cluster-scoped / multiple) | Kyverno (`kyverno`) | Identity & access |
| `ecosystem-litellm` | onemind-ai | LiteLLM (`litellm`) | AI & agents |
| `ecosystem-oc-control` | openchoreo-control-plane | OpenChoreo (`openchoreo`) | Developer platform |
| `ecosystem-ollama` | onemind-ai | Ollama (`ollama`) | AI & agents |
| `ecosystem-onemind-core` | onemind-core | OneMind core services (`onemind-core`) | Platform core |
| `ecosystem-opencost` | opencost | OpenCost (`opencost`) | Observability |
| `ecosystem-openfga` | onemind-core | OpenFGA (`openfga`) | Identity & access |
| `ecosystem-opensearch` | onemind-search | OpenSearch (`opensearch`) | Data & storage |
| `ecosystem-spire` | (cluster-scoped / multiple) | SPIRE (`spire`) | Identity & access |
| `ecosystem-spire-crds` | spire-system | SPIRE (`spire`) | Identity & access |
| `ecosystem-velero` | velero | Velero (`velero`) | Data & storage |
| `edge-cert-manager` | cert-manager | cert-manager (`cert-manager`) | Platform core |
| `edge-ingress-nginx` | ingress-nginx | ingress-nginx (`ingress-nginx`) | Platform core |
| `edge-nack` | nats | NATS JetStream (`nats`) | Messaging (NATS) |
| `edge-tailscale-operator` | tailscale | Tailscale operator and proxies (`tailscale-operator`) | Platform core |
| `monitoring-kube-prometheus-stack` | monitoring | Prometheus (`prometheus`) | Observability |
| `monitoring-loki` | monitoring | Loki (`loki`) | Observability |
| `monitoring-otel-collector` | monitoring | OpenTelemetry Collector (`opentelemetry-collector`) | Observability |
| `monitoring-tempo` | monitoring | Tempo (`tempo`) | Observability |
| `platform-agent-fabric` | agent-fabric | Agent Web UI (`agent-web-ui`) | AI & agents |
| `platform-auth-helm` | auth | Authentik (`authentik`) | Identity & access |
| `platform-auth-infra` | auth | Authentik (`authentik`) | Identity & access |
| `platform-backstage` | backstage | Backstage portal (`backstage`) | Developer platform |
| `platform-backups` | ops | Platform Backup and Recovery (`platform-backup-recovery`) | Data & storage |
| `platform-clickhouse-operator` | kube-system | ClickHouse (`clickhouse`) | Data & storage |
| `platform-cloudflared` | cloudflare | Cloudflare Tunnel (standby) (`cloudflared`) | Platform core |
| `platform-core` | (cluster-scoped / multiple) | Platform core manifests (`platform-core`) | Platform core |
| `platform-ecosystem-applications` | argocd | ArgoCD (`argocd`) | GitOps & delivery |
| `platform-ecosystem-foundation` | argocd | ArgoCD (`argocd`) | GitOps & delivery |
| `platform-ecosystem-resources` | argocd | ArgoCD (`argocd`) | GitOps & delivery |
| `platform-external-dns` | external-dns | ExternalDNS (`external-dns`) | Platform core |
| `platform-fedhub` | takserver | FedHub (`fedhub`) | TAK & field ops |
| `platform-kube` | kube | Headlamp (`headlamp`) | Developer platform |
| `platform-legacy` | hermes | Hermes (Legacy) (`hermes-legacy`) | AI & agents |
| `platform-matrix` | matrix | Matrix Synapse (`matrix-synapse`) | Communication |
| `platform-monitoring` | monitoring | Grafana (`grafana`) | Observability |
| `platform-n8n` | n8n | n8n (`n8n`) | AI & agents |
| `platform-nats` | nats | NATS JetStream (`nats`) | Messaging (NATS) |
| `platform-nats-tools` | nats | NATS JetStream (`nats`) | Messaging (NATS) |
| `platform-observability` | clickhouse | ClickHouse (`clickhouse`) | Data & storage |
| `platform-omos` | omos | OMOS data pipeline (`omos-data-pipeline`) | Messaging (NATS) |
| `platform-onetak` | onetak | OneTAK (`onetak`) | TAK & field ops |
| `platform-openchoreo-data-plane` | openchoreo-data-plane | OpenChoreo (`openchoreo`) | Developer platform |
| `platform-openchoreo-defaults` | default | OpenChoreo (`openchoreo`) | Developer platform |
| `platform-openchoreo-events` | openchoreo-observability-plane | OpenChoreo observability plane (`openchoreo-observability`) | Observability |
| `platform-openchoreo-finops` | openchoreo-observability-plane | OpenChoreo observability plane (`openchoreo-observability`) | Observability |
| `platform-openchoreo-logs` | openchoreo-observability-plane | OpenChoreo observability plane (`openchoreo-observability`) | Observability |
| `platform-openchoreo-observability` | openchoreo-observability-plane | OpenChoreo observability plane (`openchoreo-observability`) | Observability |
| `platform-openchoreo-platform` | default | OpenChoreo (`openchoreo`) | Developer platform |
| `platform-openchoreo-prereqs` | openchoreo-secrets | OpenChoreo (`openchoreo`) | Developer platform |
| `platform-openchoreo-projects` | default | OpenChoreo (`openchoreo`) | Developer platform |
| `platform-openchoreo-traces` | openchoreo-observability-plane | OpenChoreo observability plane (`openchoreo-observability`) | Observability |
| `platform-openclaw` | openclaw | OpenClaw (`openclaw`) | AI & agents |
| `platform-ops` | ops | Cluster health observer (`health-keeper`) | Observability |
| `platform-redpanda-connect` | omos | Redpanda Connect (`redpanda-connect`) | Messaging (NATS) |
| `platform-security` | argocd | Kyverno (`kyverno`) | Identity & access |
| `platform-takml` | takml | TAK-ML (`takml`) | TAK & field ops |
| `platform-takserver` | takserver | TAK Server (`tak-server`) | TAK & field ops |
| `platform-volsync` | ops | VolSync (`volsync`) | Data & storage |
| `platform-volume-backups` | ops | Platform Backup and Recovery (`platform-backup-recovery`) | Data & storage |

## Workloads to app (186 generated workloads)

Workloads become `spec.subcomponentOf` their app and join the app's System. Full list: `catalog/workload-map.yaml`. Unmapped workloads (none today) land in Unassigned workloads.

| Namespace | Workloads | App(s) |
|---|---|---|
| agent-fabric | 2 | Agent Web UI (1), Agent Monitor (1) |
| argocd | 8 | ArgoCD (7), Argo Rollouts (1) |
| auth | 4 | Authentik (4) |
| backstage | 2 | Backstage portal (2) |
| cert-manager | 3 | cert-manager (3) |
| clickhouse | 1 | ClickHouse (1) |
| cloudflare | 2 | Cloudflare Tunnel (standby) (2) |
| cnpg-system | 2 | CloudNativePG (2) |
| dapr-system | 5 | Dapr (5) |
| dp-default-onetak-development-b188bf98 | 2 | OneTAK geospatial services (2) |
| external-dns | 1 | ExternalDNS (1) |
| external-secrets | 3 | External Secrets Operator (3) |
| github-runner | 2 | Platform core manifests (2) |
| hermes | 1 | Hermes (Legacy) (1) |
| ingress-nginx | 1 | ingress-nginx (1) |
| keda | 3 | KEDA (3) |
| kube | 2 | Headlamp (1), Kubernetes Dashboard (1) |
| kube-system | 8 | Kubernetes system services (3), Hubble UI (2), ClickHouse (2), Platform core manifests (1) |
| kyverno | 4 | Kyverno (4) |
| matrix | 4 | Element Chat (1), Matrix Synapse (3) |
| monitoring | 8 | Loki (2), Grafana (1), Prometheus (4), Tempo (1) |
| n8n | 2 | n8n (2) |
| nats | 9 | HeyNATS (1), NATS JetStream (3), NATS Dashboard (1), NATS Surveyor (1), NATS UI (1), Agent Web UI (1), Agent Monitor (1) |
| omos | 6 | OMOS data pipeline (5), Redpanda Connect (1) |
| onemind-ai | 2 | LiteLLM (1), Ollama (1) |
| onemind-core | 3 | OneMind core services (2), OpenFGA (1) |
| onemind-registry | 8 | Harbor (8) |
| onemind-search | 1 | OpenSearch (1) |
| onetak | 18 | Maputnik (1), OneTAK geospatial services (7), WebODM (4), OneTAK (4), TileServer GL (1), Valhalla (1) |
| openchoreo-control-plane | 4 | OpenChoreo (4) |
| openchoreo-data-plane | 1 | OpenChoreo (1) |
| openchoreo-observability-plane | 8 | OpenChoreo observability plane (8) |
| openclaw | 1 | OpenClaw (1) |
| opencost | 1 | OpenCost (1) |
| ops | 1 | VolSync (1) |
| spire-system | 1 | SPIRE (1) |
| tailscale | 33 | Tailscale operator and proxies (33) |
| takml | 3 | TAK-ML (3) |
| takserver | 15 | FedHub (2), TAK Server (13) |
| velero | 1 | Velero (1) |

## New app Components (36)

These Argo-managed products have no curated app today, so the workloads would stay orphaned. Each is a Component of type `service`, owner `onemind-admins`.

| App | Name | System | Argo application(s) |
|---|---|---|---|
| Platform core manifests | `platform-core` | Platform core | platform-core |
| cert-manager | `cert-manager` | Platform core | edge-cert-manager |
| ingress-nginx | `ingress-nginx` | Platform core | edge-ingress-nginx |
| ExternalDNS | `external-dns` | Platform core | platform-external-dns |
| Cloudflare Tunnel (standby) | `cloudflared` | Platform core | platform-cloudflared |
| Tailscale operator and proxies | `tailscale-operator` | Platform core | edge-tailscale-operator |
| KEDA | `keda` | Platform core | ecosystem-keda |
| Dapr | `dapr` | Platform core | ecosystem-dapr |
| Kubernetes system services | `kubernetes-system` | Platform core | none |
| OneMind core services | `onemind-core` | Platform core | ecosystem-onemind-core |
| Authentik | `authentik` | Identity & access | platform-auth-helm, platform-auth-infra |
| OpenFGA | `openfga` | Identity & access | ecosystem-openfga |
| SPIRE | `spire` | Identity & access | ecosystem-spire, ecosystem-spire-crds |
| Kyverno | `kyverno` | Identity & access | ecosystem-kyverno, platform-security |
| Falco | `falco` | Identity & access | ecosystem-falco |
| External Secrets Operator | `external-secrets` | Identity & access | ecosystem-external-secrets |
| Argo Rollouts | `argo-rollouts` | GitOps & delivery | argo-rollouts |
| Harbor | `harbor` | GitOps & delivery | ecosystem-harbor |
| Backstage portal | `backstage` | Developer platform | platform-backstage |
| OpenChoreo | `openchoreo` | Developer platform | ecosystem-oc-control, platform-openchoreo-data-plane, platform-openchoreo-defaults, platform-openchoreo-platform, platform-openchoreo-prereqs, platform-openchoreo-projects |
| OpenChoreo observability plane | `openchoreo-observability` | Observability | platform-openchoreo-events, platform-openchoreo-finops, platform-openchoreo-logs, platform-openchoreo-observability, platform-openchoreo-traces |
| Tempo | `tempo` | Observability | monitoring-tempo |
| OpenTelemetry Collector | `opentelemetry-collector` | Observability | monitoring-otel-collector |
| OpenCost | `opencost` | Observability | ecosystem-opencost |
| Cluster health observer | `health-keeper` | Observability | platform-ops |
| CloudNativePG | `cloudnative-pg` | Data & storage | ecosystem-barman, ecosystem-cnpg |
| OpenSearch | `opensearch` | Data & storage | ecosystem-opensearch |
| Velero | `velero` | Data & storage | ecosystem-velero |
| VolSync | `volsync` | Data & storage | platform-volsync |
| NATS JetStream | `nats` | Messaging (NATS) | edge-nack, platform-nats, platform-nats-tools |
| Redpanda Connect | `redpanda-connect` | Messaging (NATS) | platform-redpanda-connect |
| OMOS data pipeline | `omos-data-pipeline` | Messaging (NATS) | platform-omos |
| TAK-ML | `takml` | TAK & field ops | platform-takml |
| OneTAK geospatial services | `geo-data-services` | Geospatial, drones & robotics | none |
| n8n | `n8n` | AI & agents | platform-n8n |
| LiteLLM | `litellm` | AI & agents | ecosystem-litellm |

## Existing curated apps to System (140)

| App | Name | Current System | Target System | Pillar |
|---|---|---|---|---|
| ArgoCD | `argocd` | infrastructure | GitOps & delivery | SO |
| Kubernetes Dashboard | `kubernetes-dashboard` | infrastructure | Developer platform | SO |
| Headlamp | `headlamp` | infrastructure | Developer platform | SO |
| Tailscale | `tailscale` | infrastructure | Platform core | SO |
| Cloudflare | `cloudflare` | infrastructure | Platform core | SO |
| DigitalOcean | `digitalocean` | infrastructure | Platform core | SO |
| Hermes (Legacy) | `hermes-legacy` | ai-agents | AI & agents | SO |
| Agent Web UI | `agent-web-ui` | ai-agents | AI & agents | SO |
| OpenClaw | `openclaw` | ai-agents | AI & agents | SO |
| NATS Dashboard | `nats-dashboard` | ai-agents | Messaging (NATS) | SO |
| Element Chat | `element-chat` | communication | Communication | **Unclassified** |
| Matrix Synapse | `matrix-synapse` | communication | Communication | **Unclassified** |
| HeyNATS | `heynats` | communication | Messaging (NATS) | SO |
| Google Messages | `google-messages` | communication | Communication | **Unclassified** |
| CloudTAK | `cloudtak` | tak-ecosystem | TAK & field ops | SO |
| TAK Server | `tak-server` | tak-ecosystem | TAK & field ops | SO |
| OneTAK | `onetak` | tak-ecosystem | TAK & field ops | SO |
| FedHub | `fedhub` | tak-ecosystem | TAK & field ops | SO |
| TAK Product Center | `tak-product-center` | tak-ecosystem | TAK & field ops | SO |
| TAK.gov | `tak-gov` | tak-ecosystem | TAK & field ops | SO |
| EolianVR | `eolianvr` | tak-ecosystem | TAK & field ops | SO |
| PyTAK | `pytak` | tak-ecosystem | TAK & field ops | SO |
| Grafana | `grafana` | intelligence | Observability | SO |
| NATS UI | `nats-ui` | intelligence | Messaging (NATS) | SO |
| Agent Monitor | `agent-monitor` | intelligence | AI & agents | SO |
| ClickHouse | `clickhouse` | intelligence | Data & storage | SO |
| MinIO | `minio` | data-storage | Data & storage | SO |
| PostGIS | `postgis` | data-storage | Data & storage | SO |
| Longhorn | `longhorn` | data-storage | Data & storage | SO |
| Hubble UI | `hubble-ui` | observability | Observability | SO |
| Prometheus | `prometheus` | observability | Observability | SO |
| NATS Surveyor | `nats-surveyor` | observability | Messaging (NATS) | SO |
| TileServer GL | `tileserver-gl` | geospatial | Geospatial, drones & robotics | SO |
| Valhalla | `valhalla` | geospatial | Geospatial, drones & robotics | SO |
| Maputnik | `maputnik` | geospatial | Geospatial, drones & robotics | SO |
| WebODM | `webodm` | geospatial | Geospatial, drones & robotics | SO |
| Zoho One | `zoho-one` | zoho-apps | Productivity & SaaS | **Unclassified** |
| Zoho Mail | `zoho-mail` | zoho-apps | Productivity & SaaS | **Unclassified** |
| Zoho Store | `zoho-store` | zoho-apps | Productivity & SaaS | **Unclassified** |
| Zoho Account | `zoho-account` | zoho-apps | Productivity & SaaS | **Unclassified** |
| Zoho MCP Client | `zoho-mcp-client` | zoho-apps | Productivity & SaaS | **Unclassified** |
| Zoho API Console | `zoho-api-console` | zoho-apps | Productivity & SaaS | **Unclassified** |
| Zoho Social | `zoho-social` | zoho-apps | Productivity & SaaS | **Unclassified** |
| Zoho Workplace | `zoho-workplace` | zoho-apps | Productivity & SaaS | **Unclassified** |
| Zoho Projects | `zoho-projects` | zoho-apps | Productivity & SaaS | **Unclassified** |
| Zoho Flow | `zoho-flow` | zoho-apps | Productivity & SaaS | **Unclassified** |
| Zoho IoT | `zoho-iot` | zoho-apps | Productivity & SaaS | **Unclassified** |
| Zoho Zia Agents | `zoho-zia-agents` | zoho-apps | Productivity & SaaS | **Unclassified** |
| 1Password | `1password` | personal-tools | Identity & access | SO |
| Dropbox | `dropbox` | personal-tools | Productivity & SaaS | **Unclassified** |
| Raindrop.io | `raindrop-io` | personal-tools | Productivity & SaaS | **Unclassified** |
| AT&T Prepaid | `at-t-prepaid` | personal-tools | Communication | **Unclassified** |
| GitHub | `github` | workspaces | GitOps & delivery | SO |
| Composio | `composio` | workspaces | AI & agents | SO |
| Miro | `miro` | workspaces | Productivity & SaaS | **Unclassified** |
| Circle | `circle` | workspaces | Business operations | GE |
| Mercury | `mercury` | personal-finance | Personal finance | **Unclassified** |
| Credit Karma | `credit-karma` | personal-finance | Personal finance | HP |
| Experian | `experian` | personal-finance | Personal finance | HP |
| TransUnion | `transunion` | personal-finance | Personal finance | HP |
| Equifax | `equifax` | personal-finance | Personal finance | HP |
| Innovis | `innovis` | personal-finance | Personal finance | HP |
| OneMind OS | `onemind-os` | external | Developer platform | SO |
| ArgoCD | `external-argocd` | external | GitOps & delivery | SO |
| Kubernetes Dashboard | `external-kubernetes-dashboard` | external | Developer platform | SO |
| Tailscale | `external-tailscale` | external | Platform core | SO |
| Tailscale Admin | `tailscale-admin` | infrastructure | Platform core | SO |
| OneMind Operations Runbooks | `onemind-runbooks` | infrastructure | Developer platform | SO |
| OneMind Observational API | `onemind-observations` | (none) | Intelligence & feeds | **Unclassified** |
| QGroundControl | `qgroundcontrol` | geospatial | Geospatial, drones & robotics | SO |
| LingBot-Map | `lingbot-map` | geospatial | Geospatial, drones & robotics | SO |
| PX4 Autopilot | `px4` | geospatial | Geospatial, drones & robotics | SO |
| ArduPilot | `ardupilot` | geospatial | Geospatial, drones & robotics | SO |
| Platform Backup and Recovery | `platform-backup-recovery` | infrastructure | Data & storage | SO |
| Linux Foundation | `linux-foundation` | (none) | Reference & community | **Unclassified** |
| Cloud Native Computing Foundation (CNCF) | `cncf` | (none) | Reference & community | **Unclassified** |
| CNCF Landscape | `cncf-landscape` | (none) | Reference & community | **Unclassified** |
| LF Edge | `lf-edge` | (none) | Reference & community | **Unclassified** |
| OpenSSF | `openssf` | (none) | Reference & community | **Unclassified** |
| Open Robotics | `open-robotics` | (none) | Reference & community | **Unclassified** |
| Proton Mail | `proton-mail` | (none) | Communication | **Unclassified** |
| Autoware | `autoware` | (none) | Geospatial, drones & robotics | SO |
| ROS 2 | `ros-2` | (none) | Geospatial, drones & robotics | SO |
| Eclipse Zenoh | `eclipse-zenoh` | (none) | Geospatial, drones & robotics | SO |
| Eclipse Ditto | `eclipse-ditto` | (none) | Geospatial, drones & robotics | SO |
| NVIDIA Omniverse | `nvidia-omniverse` | (none) | Geospatial, drones & robotics | SO |
| vLLM | `vllm` | (none) | AI & agents | SO |
| LM Studio | `lm-studio` | (none) | AI & agents | SO |
| oMLX | `omlx` | (none) | AI & agents | SO |
| Jan | `jan` | (none) | AI & agents | SO |
| Ollama | `ollama` | (none) | AI & agents | SO |
| DigitalOcean Inference | `digitalocean-inference` | (none) | AI & agents | SO |
| NVIDIA Build / NIM | `nvidia-build` | (none) | AI & agents | SO |
| Together AI | `together-ai` | (none) | AI & agents | SO |
| Fireworks AI | `fireworks-ai` | (none) | AI & agents | SO |
| Alibaba / Qwen | `qwen-models` | (none) | AI & agents | SO |
| Google / Gemma | `gemma-models` | (none) | AI & agents | SO |
| Mistral / Ministral | `ministral-models` | (none) | AI & agents | SO |
| NVIDIA / Nemotron | `nemotron-models` | (none) | AI & agents | SO |
| Deskflow | `deskflow` | personal-tools | Personal tools & devices | **Unclassified** |
| Facebook — Zeus | `facebook-zeus` | social-media | Brand & social media | **Unclassified** |
| Facebook — Legacy Platform | `facebook-legacy` | social-media | Brand & social media | **Unclassified** |
| X — Zeus | `x-zeus` | social-media | Brand & social media | **Unclassified** |
| LinkedIn — Zeus | `linkedin-zeus` | social-media | Brand & social media | **Unclassified** |
| TikTok — Zeus | `tiktok-zeus` | social-media | Brand & social media | **Unclassified** |
| Instagram — Zeus | `instagram-zeus` | social-media | Brand & social media | **Unclassified** |
| Twitch — Zeus | `twitch-zeus` | social-media | Brand & social media | **Unclassified** |
| Reddit — callsign_zero0 | `reddit-zeus` | social-media | Brand & social media | **Unclassified** |
| Pinterest — Zeus | `pinterest-zeus` | social-media | Brand & social media | **Unclassified** |
| Threads — Zeus | `threads-zeus` | social-media | Brand & social media | **Unclassified** |
| YouTube — profile pending | `youtube-zeus` | social-media | Brand & social media | **Unclassified** |
| LinkedIn — Company Admin | `linkedin-company-admin` | social-media | Brand & social media | **Unclassified** |
| Bambu Studio | `bambu-studio` | personal-tools | Personal tools & devices | **Unclassified** |
| Bambu Lab | `bambu-lab` | personal-tools | Personal tools & devices | **Unclassified** |
| Foxglove Studio | `foxglove-studio` | geospatial | Geospatial, drones & robotics | SO |
| OBS Studio | `obs-studio` | personal-tools | Personal tools & devices | **Unclassified** |
| Blender | `blender` | personal-tools | Personal tools & devices | **Unclassified** |
| MadMapper | `madmapper` | personal-tools | Personal tools & devices | **Unclassified** |
| TouchDesigner | `touchdesigner` | personal-tools | Personal tools & devices | **Unclassified** |
| Termius | `termius` | infrastructure | Personal tools & devices | **Unclassified** |
| Brave Browser | `brave-browser` | personal-tools | Personal tools & devices | **Unclassified** |
| Google Chrome | `google-chrome` | personal-tools | Personal tools & devices | **Unclassified** |
| Google Cloud | `google-cloud` | infrastructure | Platform core | SO |
| Riverside.fm | `riverside-fm` | workspaces | Productivity & SaaS | **Unclassified** |
| Notion Backups | `notion-backups` | infrastructure | Data & storage | SO |
| ProBackup.io | `probackup-io` | infrastructure | Data & storage | SO |
| Strava | `strava` | personal-tools | Personal tools & devices | **Unclassified** |
| Apple Music | `apple-music` | personal-tools | Personal tools & devices | **Unclassified** |
| Langfuse | `langfuse` | observability | Observability | SO |
| Uptime Kuma | `uptime-kuma` | observability | Observability | SO |
| MTA | `mta` | communication | Communication | **Unclassified** |
| Agent Mail | `agent-mail` | ai-agents | AI & agents | SO |
| turbovec | `turbovec` | ai-agents | AI & agents | SO |
| Samsung KNOX | `samsung-knox` | personal-tools | Personal tools & devices | **Unclassified** |
| Samsung One UI | `samsung-one-ui` | personal-tools | Personal tools & devices | **Unclassified** |
| Samsung DeX | `samsung-dex` | personal-tools | Personal tools & devices | **Unclassified** |
| Notion | `notion` | workspaces | Productivity & SaaS | **Unclassified** |
| Home Assistant | `home-assistant` | infrastructure | Home & homestead | LE |
| Loki | `loki` | observability | Observability | SO |
| OneTAK STAC read-only API | `onetak-stac-readonly` | onetak | Geospatial, drones & robotics | SO |

## Resources (Systems only change)

Resources keep their names. Documented capabilities (Notion registry), repositories, mesh endpoints and PVCs move to the successor System of their old System (table above): for example all 83 repositories stay in Code repositories; the 5 PVCs follow their namespace's app (backstage, grafana, nats) via the namespace row. `do-nyc3-omos-cloud` (no System today) joins Platform core.

## Pillar rules applied

- A pillar is set on a System only when the doctrine names the subject (SO: infrastructure, security, development, automation, agents, TAK, drones, robotics; HP: personal finance; LE: home and homestead; GE: business operations). Everything else is Unclassified.
- Entities inherit the pillar of their System as label `onemind.io/pillar`, unless overridden in `structure.yaml` (`pillarOverrides`). No `pillars-served` values are set: no cross-pillar service is approved yet.
- Namespace, vendor, host and launcher category were not used as evidence.

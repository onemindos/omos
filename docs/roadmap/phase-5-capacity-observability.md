# Phase 5 brief: capacity, then the observability plane

Status when written (2026-10-18): not started. Tracker: [`PLATFORM-ROADMAP.md`](../../PLATFORM-ROADMAP.md). Read [`HANDOFF.md`](HANDOFF.md) first.

## Goal

Make room, then install OpenChoreo's observability plane on top of the monitoring stack we already
run, and give Claude the observer MCP server (logs, metrics, traces, alerts, incidents, costs).

## Starting state

- Nodes: 1 `onemind-ecosystem-cpu` (8 CPU, 32 GiB) and 3 `onemind-hub-core-v2` (8 CPU, 16 GiB each,
  recycled by DOKS on 2026-10-18). The ecosystem node is at about 99% CPU requested; the hub-core
  nodes had room after the recycle, but nothing rebalances pods.
- Findings and request cuts: `docs/operations/cluster-review-2026-10-09.md` (H1).
- Existing stack: kube-prometheus-stack, Loki, Tempo, an OTel collector (all in `monitoring`),
  OpenSearch (`onemind-search`), OpenCost, KEDA.

## Part A: capacity

1. Re-measure: requests and real usage per node and pod (`onemind-kubernetes` `nodes_top`,
   `pods_top`, `resources_list`; Prometheus). Update the H1 table.
2. Cut over-requests with Git changes (one PR per layer): the review lists onetak, onemind-data,
   onemind-core, matrix, kyverno, spire and openchoreo; raise under-requests for Prometheus.
3. Add pod affinity for ecosystem workloads to prefer the ecosystem pool, TAK/NATS/auth on hub-core.
   Never name a node.
4. Consider the Kubernetes descheduler (Layer 1, Argo CD) so pods rebalance after node
   replacement; or add a node (Terraform, `infra/`) if cuts are not enough.
5. Target: no node above about 80% of memory or CPU requested.

## Part B: observability plane (install guide v1.3.x, step 7)

Follow `https://openchoreo.dev/docs/getting-started/try-it-out/on-your-environment/` (source:
`openchoreo.github.io` repo, `versioned_docs/version-v1.3.x/getting-started/try-it-out/on-your-environment.mdx`,
section "Step 7: Setup Observability Plane") and the reference
`https://openchoreo.dev/docs/reference/helm/observability-plane/`. Charts (guide versions):
`openchoreo-observability-plane` 1.3.1, `observability-logs-opensearch` 0.6.0,
`observability-metrics-prometheus` 0.7.1, `observability-tracing-opensearch` 0.6.0,
`observability-events-otel-collector` 0.1.2, plus Fluent Bit for logs. Registry: `ghcr.io/openchoreo/helm-charts`.

1. New Argo Applications under `platform/openchoreo/observability/` (same pattern as
   `platform/argocd/platform-openchoreo-data-plane.yaml`; values in Git).
2. Reuse what exists where the modules allow: point the logs and traces modules at the existing
   OpenSearch (credentials through SealedSecrets, the guide's ExternalSecret keys are
   `opensearch-username`, `opensearch-password`, `observer-oauth-client-secret`) and the metrics module
   at the existing Prometheus. Note our logs live in Loki and traces in Tempo; OpenChoreo has no Loki or
   Tempo module, so decide per signal (see the community modules repo
   `https://github.com/openchoreo/community-modules`) and record it.
3. Copy `cluster-gateway-ca` into `openchoreo-observability-plane` (as phase 2c did for the data plane),
   install, then register `ClusterObservabilityPlane` `default` (agent CA from `cluster-agent-tls`, `observerURL`)
   and link the data plane: `observabilityPlaneRef` (kind `ClusterObservabilityPlane`, name `default`).
4. Authentication on the observer: set `security.oidc.issuer` and `jwksUrl` to Authentik and
   `observer.security.subjectTypes` to match the control plane (`user` maps `sub`, `service_account`
   maps `client_id`); the roles in `platform/openchoreo/platform/ai-operator-read.yaml` already contain
   the observability `:view` actions (`logs`, `metrics`, `traces`, `events`, `alerts`, `incidents`,
   `rcareport`, `finops`).
5. Expose the observer on the tailnet (Tailscale ingress or the gateway) under a stable hostname.
6. Add `onemind-openchoreo-obs` (MCP at `<observer>/mcp`) to `scripts/ai-operator/setup.sh`, same
   token helper (`scripts/ai-operator/openchoreo-token.sh`) and docs. Re-enable the console's
   observability and workflow features if the planes exist (`platform/ecosystem/values/oc-control.yaml`).
7. FinOps: the `finops-opencost` module from community-modules, using the existing OpenCost.

## Gate

Nodes under the capacity target; a phase-3 component shows logs and metrics through the
observer API and through the MCP server; the observability Argo apps are Synced and Healthy.

## Risks

- The ecosystem node has no CPU left, and the observability plane is the heaviest part of OpenChoreo.
- Multi-signal duplication (Loki and OpenSearch both store logs): choose one source of truth per signal.
- The workflow plane (build pipelines, guide step 6) is not required while images are prebuilt.

# Tool gateway: Kubernetes backend (Phase 3a)

Design: [`docs/architecture/onemind-tool-gateway.md`](../architecture/onemind-tool-gateway.md) sections 2.1 row 1, 3 and 7. Generic recipe and
renderer: [`tool-gateway.md`](tool-gateway.md).

## What it is

| Thing | Value |
|---|---|
| Server | `containers/kubernetes-mcp-server` v0.1.1 (released 2026-10-29), toolset `core` |
| Image | `ghcr.io/containers/kubernetes-mcp-server:v0.1.1@sha256:d616331a99184d42f7471eac5ae257192e7065bb684ad115cc95c2e76cfbc2d1` (multi-arch index; the same digest is on `quay.io/containers/kubernetes_mcp_server`) |
| Workload | Deployment and Service `mcp-kubernetes` in `tool-gateway`, 1 replica, `onemind-ecosystem-cpu`, requests 25m / 64Mi, limits 500m / 256Mi |
| Identity | ServiceAccount `tool-gateway/k8s-operator` (not the Mac's `kube/ai-operator` token) |
| Gateway target | `kubernetes` in group `ops`; tools appear to Hermes as `mcp_ops_kubernetes_<tool>` |
| Files | `platform/tool-gateway/manifests/backends/kubernetes/`, `manifests/egress-kubernetes.yaml`, `targets/kubernetes.yaml`, `policy/classification/kubernetes.yaml` |

The server has no credential of its own to call it: the NetworkPolicy admits only the gateway pod on 8080, and the gateway route needs the client
key. Its egress is DNS and the API server (Cilium `kube-apiserver` entity, port 443).

## Tools (14 live)

| Class | Tools | Where |
|---|---|---|
| read (10) | `events_list`, `namespaces_list`, `nodes_top`, `pods_get`, `pods_list`, `pods_list_in_namespace`, `pods_log`, `pods_top`, `resources_get`, `resources_list` | `/ops/mcp` |
| write (1) | `resources_create_or_update` (server-side apply on existing objects) | `/ops/mcp` |
| destructive (3) | `pods_delete`, `resources_delete`, `resources_scale` | `/execute/mcp` only (closed until Phase 4) |

Turned off in the server config (`disabled_tools`), so they do not exist: `pods_exec`, `pods_run` (arbitrary code in the cluster; no RBAC is
granted), `nodes_log`, `nodes_stats_summary` (need `nodes/proxy`, which allows exec on every pod of the node) and `projects_list` (OpenShift).
The `config` toolset (kubeconfig views) is not enabled; the design listed it, but it does nothing in a pod. To offer `pods_exec` later: grant
`pods/exec` in `rbac.yaml`, remove it from `disabled_tools`, list it under `destructive`.

## RBAC ceiling (what the API server enforces whatever the gateway does)

- **Read**: built-in `view` plus `tool-gateway-k8s-read` (nodes, PVs, storage classes, CRDs, metrics, and the platform CRD groups of
  `platform/kube/ai-operator.yaml`). **No Secrets**, and the server also denies the `Secret` kind. `clickhouse.altinity.com` and `dapr.io` are
  left out of the CRD list on purpose: their specs hold inline passwords or connection strings. Logs and ConfigMaps are readable (as for
  `ai-operator`); a secret printed to a log is readable.
- **Write**: `patch` on Deployments, StatefulSets, DaemonSets and CronJobs; `get/update/patch` on the scale subresources; `patch/update` on the one
  ConfigMap `tool-gateway/gateway-selftest`. **No `create`** anywhere (new objects come from Git through Argo CD).
- **Delete**: `pods` only. Nothing else.
- Not granted: Secrets, `pods/exec`, `pods/attach`, `pods/portforward`, `pods/eviction`, `nodes/proxy`, RBAC kinds, namespaces.
  `scripts/tests/test_gitops_tool_gateway_kubernetes.py` fails on any of these.

## Known limits (read before extending)

1. **Argument predicates do not work in agentgateway v1.7.0** `allow` rules (tested against the pinned image: a guard on `mcp.tool.arguments`
   never fired, and listing/calling behaved as if the variable were absent). So `resources_scale` is destructive outright, and
   `resources_create_or_update` is write even though a manifest with `replicas: 0` scales a workload to zero. RBAC cannot separate the two. The
   Phase 4 execute gate (AuthZEN, which sees the arguments) is where "apply with replicas 0" and similar rules belong.
2. `patch` on workloads lets a write-class call change an image, env, or mount a Secret into a pod and print it to the logs. The design accepts
   write-class workload edits; flag it if the threat model changes (a Kyverno policy on `tool-gateway`-originated changes would close it).
3. `resources_delete` is classified destructive but the API server refuses it for anything except Pods.

## Verify live (after merge, after Argo syncs and the gateway ConfigMap reloads, 60 to 120 s)

From a workstation, with the key piped to curl on stdin so it is never printed (pattern in `tool-gateway.md`, Operate):

1. `kubectl -n tool-gateway get deploy mcp-kubernetes` Ready; `kubectl auth can-i get secrets --as=system:serviceaccount:tool-gateway:k8s-operator -A` is `no`.
2. `initialize`, `notifications/initialized`, `tools/list` on `/ops/mcp`: 11 tools with the `kubernetes_` prefix plus the `onemind` ones; none of
   `pods_delete`, `resources_delete`, `resources_scale`.
3. Read: `tools/call kubernetes_pods_list_in_namespace {"namespace":"tool-gateway"}`.
4. Safe write: `tools/call kubernetes_resources_create_or_update` with a ConfigMap `gateway-selftest` manifest that adds label
   `gateway-test: "<timestamp>"`; then `kubectl -n tool-gateway get cm gateway-selftest --show-labels`.
5. Destructive: `tools/call kubernetes_pods_delete` returns JSON-RPC `-32602 unknown resource`; `/execute/mcp` `tools/list` is empty.

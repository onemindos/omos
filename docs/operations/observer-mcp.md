# Authenticated Observer MCP

Phase 5C connects the official OpenChoreo Observer at
https://onemind-observer.tail717ba0.ts.net/mcp. Merge this stage only after the plane
and modules are healthy and capacity is within target.

The same Authentik ai-operator identity and ai-operator-read OpenChoreo role used
by the control-plane MCP provide access. Tailnet reachability and JWT authorization
are both required. This does not grant Kubernetes writes or raw Secrets access.

The stage registers ClusterObservabilityPlane/default with the observability
agent's own CA reference, links ClusterDataPlane/default explicitly, and enables
console observability. It wires OTel component logs/traces to the prepared OpenSearch
indices and exposes OpenChoreo labels in the existing kube-state-metrics. This is
ordered after module setup to avoid creating indices with accidental dynamic mappings.
The Alertmanager Config matcher strategy allows the explicitly labeled OpenChoreo
route to receive alerts from component namespaces, while retaining other receivers.
Workflow features stay disabled.

Run from the repository on the operator Mac:

    scripts/ai-operator/setup.sh
    python3 scripts/ai-operator/verify-observer.py --component stac-browser

Setup verifies unauthenticated rejection and authenticated initialization/tool
discovery before registering onemind-openchoreo-obs in Claude Code and Codex.
Claude uses the existing token helper; Codex uses OPENCHOREO_MCP_TOKEN, like the
control-plane MCP. No credential is committed or printed.

Read tools include query_component_logs, query_resource_metrics, query_traces,
query_alerts, query_incidents and query_costs. STAC Browser emits no application
traces by default; no traces differs from a failed adapter. Optional RCA/AI analysis
stays disabled. Alerts, incidents and allocations may be empty; no data is fabricated.

## Live Phase 5 gate

The optional component check verifies transport/tool execution. It does not treat
an empty response as proof telemetry exists. Make normal requests to the pilot,
then confirm recent logs and CPU/memory samples in both console and MCP. Also check
the plane's agent connection, Argo health, per-node requested capacity and alert-store
backup/restore. Record evidence before marking the roadmap Done.

Status: configuration prepared; live rollout, telemetry and console checks pending.

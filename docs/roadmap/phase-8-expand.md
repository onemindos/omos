# Phase 8 brief: expand (home server, Hermes)

Status when written (2026-10-18): not started. Needs Zeus's decisions on hardware and priorities.
Tracker: [`PLATFORM-ROADMAP.md`](../../PLATFORM-ROADMAP.md). Read [`HANDOFF.md`](HANDOFF.md) first.

## Goal

1. A home server as a second OpenChoreo data plane (and a NATS leaf, later a second authority), matching
   `docs/MISSION.md`: home is the primary source of authority, cloud is hosting and recovery.
2. Hermes ("Legacy") gets a model provider and the same read-only tools Claude uses.
3. Revisit where TAK Server runs.

## Open decisions for Zeus (ask before building)

- Home hardware, OS (Talos Linux was recommended), network and outage targets
  (`docs/MISSION.md`, "Existing context and decisions still open").
- Whether the OpenChoreo control plane stays in DOKS or moves home. Check what a home data plane does
  when the cloud control plane is unreachable: running workloads should continue, but no changes apply.
- Which model Hermes uses (through LiteLLM in `onemind-ai`, or local Ollama) and its budget.

## Steps

### A. Home server as a data plane
1. Cluster: install Talos Linux (or chosen OS), join the tailnet. Register it in Argo CD as a second
   cluster (`argocd cluster add`, or a declarative cluster Secret) and deploy the same foundation basics
   (cert-manager, Gateway API CRDs, External Secrets) from Git.
2. OpenChoreo multi-cluster: follow `https://openchoreo.dev/docs/platform-engineer-guide/multi-cluster-connectivity/`.
   The data-plane agent dials out to the control plane's cluster gateway over WebSocket (port 8443, path
   `/ws`), so a home server behind NAT only needs outbound access. The `cluster-gateway-external` Service is
   `ClusterIP` today; expose it on the tailnet (Tailscale LoadBalancer or ingress) and give the agent
   `clusterAgent.serverUrl: wss://<that host>:8443/ws` with a matching certificate SAN. Register a second
   `ClusterDataPlane` (own `planeID`, the agent's CA) and a home `Environment`. Docs note: gateway is a single
   replica, agent CAs are self-signed and copied by hand.
3. NATS: run a leaf node at home (see `infra/edge/mac-leaf/` and `docs/NATS-DOCTRINE.md`).
4. Move one component to the home environment as the test (ReleaseBinding to the home environment).

### B. Hermes
1. Give Hermes a model (`platform/hermes/agent/`, LiteLLM keys in `platform/ecosystem/resources/secrets/`).
2. Decide the tool gateway for in-cluster agents: LiteLLM's MCP gateway or agentgateway (Solo.io, Apache 2.0,
   an OpenChoreo ecosystem module). Both need per-tool read-only rules.
3. Connect Hermes to the Backstage and OpenChoreo MCP servers with its own read-only identities (a
   dedicated Authentik service account per agent, not the human's or Claude's), and to the Kubernetes
   MCP server behind a view-only ServiceAccount. Hermes currently has no Kubernetes RBAC (its runbook).
4. Move the `hermes` namespace into OpenChoreo last (phase 7 order).

### C. TAK Server placement
`takserver` is Layer 1 for now (Helm CLI release plus Argo overlays). Revisit once the home server exists.

## Gate

A component runs on the home data plane and is manageable from the portal; Hermes answers "what is
unhealthy and why?" through MCP with read-only access; the mission's cloud-loss test is written down and run.

## Risks

- Cluster-gateway HA is pre-release, agent certificates are self-signed (rotation breaks registration unless handled).
- Home exposure and Tailscale ACLs: `platform/core/policy.hujson`.
- Do not give any agent write access above Layer 3.

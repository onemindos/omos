# Phase 5: capacity and observability — completed

Live acceptance passed on 2026-10-19. The authoritative tracker and evidence are in
[PLATFORM-ROADMAP.md, Phase 5](../../PLATFORM-ROADMAP.md#phase-5-capacity-then-observability).
This document records the implementation, not an unexecuted deployment proposal.

## Delivered

1. Moved measured requests off hub-core through ecosystem preference; right-sized
   Hermes and added supported Tailscale proxy resources/namespace memory defaults.
   Every node is below 78% of both memory and CPU requested. No extra node or
   descheduler was needed; no immutable Loki storage template was changed.
2. Installed official OpenChoreo observability core 1.3.1 and pinned logs/traces,
   metrics, events and FinOps adapters using existing OpenSearch, Prometheus and
   OpenCost. Existing OTel retains Loki/Tempo and copies only component-tagged signals
   to OpenSearch. There is no duplicate application collector or monitoring stack.
3. Registered the plane with its own client CA, linked the default data plane, and
   enabled observability in the built-in console. Workflow features remain disabled.
4. Registered authenticated Observer MCP in Claude Code/Codex. Claude's real
   connection check reports Connected. Anonymous requests are rejected; existing
   read-only identities/permissions and issuer/audience checks are retained.
5. Added daily encrypted alert-store backup and proved an isolated, read-only
   SQLite restore. The restore PVC is protected and retained, not silently deleted.

## Gate evidence

- Five ordinary pilot requests returned HTTP 200. The signed-in console displayed
  five real STAC Browser logs and CPU/memory graphs for development.
- MCP returned five logs, eight points in each of six resource series and one cost
  allocation. All six required tool categories execute; empty traces/alerts/incidents
  are valid, not fabricated samples or failed adapters.
- All five plane/module Argo Applications are Synced/Healthy; authenticated agent
  connected, both index-setup Jobs successful, existing OTel 4/4 ready.
- Restic snapshot 583bae7f restored into new storage. SQLite integrity_check=ok,
  alert_entries and incident_entries tables present. Daily 04:50 UTC backup restored
  after the one-off manual trigger. Production PVC never overwritten.
- Portal remains 2/2 available, readiness 200. Only pre-existing Matrix/OpenClaw
  Argo exceptions remain; no Pending pods.

## Live issues fixed before completion

- Upstream setup wraps long Base64 Basic-auth headers: Git-owned hooks normalize
  wrapping while retaining the same pinned images, scripts and credentials.
- Event label app vs app.kubernetes.io mapping collision: enrich only OpenChoreo labels.
- Authentik's internal token URL minted an untrusted internal-host issuer: use the
  existing public HTTPS endpoint, not disabled issuer validation.
- Console fetches the registered Observer URL directly: use the authenticated
  tailnet HTTPS URL, not cluster-only DNS. Backend/frontend URL caches last five minutes.

Implementation PRs #300–#310 are individually linked in the roadmap. Local gates
passed except existing whole-repository YAML lint debt. GitHub later stopped starting
CI jobs due its account payment/spending-limit state; no billing or merge-rule bypass.

Phase 6 is next, but was not started as part of Phase 5.

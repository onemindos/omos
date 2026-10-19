# OpenChoreo observability plane

Stage 5B, following the capacity PRs. These are deployment sources; live acceptance
is pending until merge, Argo reconciliation and the checks below.

Core: openchoreo-observability-plane 1.3.1. Logs/tracing adapters: 0.6.0.
Metrics adapter: 0.7.1. Events collector: 0.1.2. OpenCost adapter: 0.2.0.
All images and charts are version-pinned.

## Existing backends and signal ownership

- Prometheus remains the only metrics backend. KSM exposes OpenChoreo identity
  labels. The upstream metrics chart hard-codes port 9091; the rendered adapter
  resources in resources/metrics-adapter.yaml change it to the existing Prometheus
  on port 9090 and add ecosystem affinity. No second monitoring stack is installed.
- OpenSearch remains the only search cluster. A module-local ExternalName Service
  points to the existing search Service. Credentials are copied by a narrowly
  scoped, idempotent Argo hook; no values enter Git.
- The next registration stage wires the existing OTel DaemonSet, retaining Loki/Tempo.
  Additional pipelines copy only signals carrying an OpenChoreo component UID to
  OpenSearch. OpenChoreo queries that component-specific copy; existing Loki/Tempo
  data is retained. No Fluent Bit or second log/trace collector runs.
- OTel bodymap maps component logs to the module's Fluent Bit-compatible schema.
  Labels replace dots with underscores and preserve slash suffixes. Runtime target,
  namespace, pod and container fields are retained. Trace export uses the official
  module's SS4O mapping instead of the log-only bodymap mapping.
- The single events collector supplies enriched Kubernetes events. It does not
  collect application logs. Keep one replica to avoid duplicate event watches.
- FinOps uses the existing OpenCost through its official adapter. Optional AI
  FinOps and RCA agents stay disabled; this does not invent incident records.

The existing OpenSearch endpoint uses its current demo/self-signed TLS setup.
Module clients and OTel retain the upstream internal TLS verification exception.
Replace it with managed backend trust in a separate backend-hardening change.

## Authentication, state and recovery

The Observer checks Authentik JWT issuer/audience and resolves sub, matching the
control plane's read-only AI role. Its UID resolver reuses the existing read-only
catalog account, through a runtime copy of the console's service credential.
No credential value or private key is stored here.

The hook reads only named source Secrets in control-plane, search and monitoring
namespaces. Target get/patch permissions name the generated Secrets; create is
limited to the two target namespaces. It does not rotate credentials. Re-sync
this Application after changing a source credential.

The SQLite alert PVC has a daily VolSync snapshot source and a distinct encrypted
restic repository. The hook reuses the existing backup account/encryption material
with a different repository path. It never writes to Loki's repository. A successful
backup and isolated restore remain rollout acceptance checks.

## Order and acceptance

1. Reconcile capacity/resource PRs; require about 80% or less requested memory per node.
2. Reconcile this stage. Verify core, logs, traces, events and FinOps Argo Applications
   are Synced/Healthy, setup Jobs succeed, and alert-store backup starts succeeding.
3. After index-setup Jobs succeed, merge the registration/MCP stage: enable OTel/KSM,
   pin this plane's own agent CA and link the
   default data plane, enables the console and extends operator setup.
4. Run the authenticated Observer verifier and component log/metric queries. Generate
   normal HTTP traffic to STAC Browser, then inspect real samples in console and MCP.
   An empty result alone does not pass the gate.
5. Record node totals, Argo revision/health, agent connection, log/metric sample counts
   and a human console check in PLATFORM-ROADMAP.md before closing Phase 5.

Local checks: chart renders, Kustomize render, YAML lint, OTel 0.156.0 configuration
validation, and isolated HTTP ingestion/export proving component log schema,
unrelated-log exclusion and preservation of the original Loki branch.
Full make validate still fails on existing YAML debt in onetak, TAK and Matrix;
changed-file gates must stay green.

References:
- [OpenChoreo installation](https://openchoreo.dev/docs/getting-started/try-it-out/on-your-environment/)
- [Community modules](https://github.com/openchoreo/community-modules)
- [OTel 0.156 OpenSearch mappings](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/v0.156.0/exporter/opensearchexporter/README.md)

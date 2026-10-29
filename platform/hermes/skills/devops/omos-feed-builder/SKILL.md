---
name: omos-feed-builder
description: Build OMOS feed pods for the NATS intelligence fabric.
---

# OMOS Feed Builder

Build and assess real-time sources for the OneMind Intelligence Fabric. Feed runtimes can emit typed entity/intel envelopes or CoT XML according to the deployed consumer contract. Do not assume that every source is a CoT producer, that every route uses the same subject, or that every sink archives every route.

## Architecture

Target boundary (verify actual subjects and schemas before migration):

```
External source -> collector/decoder -> validated observations -> durable NATS/JetStream
                                                        -> archive / AI / apps
                                                        -> controlled TAK projection -> TAK ingress -> clients
Tiles/terrain/assets -> map/asset services -> client layers
Audio/video -> media services -> client players (metadata/derived observations on NATS)
```

Feed collectors should not each open their own TAK connection. Use an independently observable TAK publication/ingress boundary. Distinguish broker persistence, TAK acceptance and client rendering; they are different proofs.

## Files Per Feed (2 files)

1. `platform/omos/feed-{name}-script.yaml` - ConfigMap with Python script
2. `platform/omos/feed-{name}.yaml` - Deployment manifest

## Location

Canonical GitOps repo: `onemindos/omos`, branch `main`. Verify the live Argo CD Application's `spec.source` before editing.
Feed manifests: `platform/omos/`.
Redpanda Connect: `platform/data/redpanda-connect/deployment.yaml`, Argo CD application `platform-redpanda-connect`, destination namespace `omos`.
Guide: inspect `platform/omos/FEED-DEVELOPMENT.md` in the canonical revision if present.

## Key Steps

1. Inspect the live Argo source/revision, workload sharing, accepted schemas and available credentials without printing secret values. Distinguish repo intent from running code.
2. Establish each authorized source's executable request/subscription, identity/freshness rules, output policy and representative fixtures before choosing its runtime.
3. Implement and test the collector/decoder and durable sink/control behavior. Pin dependencies; do not rely on unpinned runtime package installation for a reproducible production image.
4. Change the canonical GitOps manifests and use the approved Argo deployment path. Do not use ad-hoc apply/scale as a substitute for persistent feed controls; do not deploy during a read-only assessment.
5. Verify source response, validation, bus persistence, archival commit, TAK receipt and selected-client display with matching event/UID evidence. A log line or existing `cot_router` row alone does not certify current end-to-end delivery.

## CoT Type Reference

- Trains: `a-n-G-E-V-C-R`
- Aircraft: `a-f-A-M-F-Q`
- Ships: `a-f-S-X-M`
- Ground vehicle: `a-n-G-E-V-C`
- Earthquake/hazard: `b-r-f-h-h`
- Fire/thermal: `b-t-f`

## Toggle

Use a persisted, authorized per-source control contract with desired/applied revision reporting, separate TAK publication policy and per-client visibility. Inspect whether multiple adapters share one Deployment before stopping any workload. Argo-managed replicas can reconcile ad-hoc scaling back, and stopping a shared worker pauses unrelated sources. Treat pod scaling as an explicitly approved operational stop, not the operator's durable feed on/off feature.

## Requirements Discovery and Operator Controls

- Before choosing a feed runtime, read the full source catalog and the original architecture intent, not just the running pods. The owner's scope includes reliable TAK visualization plus per-source on/off control; a generic HTTP-workflow comparison is insufficient.
- Source-catalog and design notes may exist in the owner's vault or Notion exports. Discover the requested location instead of assuming a path, and treat them as requirements and historical metadata, not live deployment evidence.
- Parse YAML frontmatter and count actual `kind: record` files. Exclude host/index notes and `.base` view definitions. Compare the actual count with the index's declared count and report discrepancies; do not discard records to force a match. Preserve original Notion IDs, declared subjects, and source paths. Provider homepages, web maps, protocol standards, and paid dashboards are not verified ingestion endpoints.
- Separate live entity/alert projection, raster/vector map delivery, media transport, and non-geospatial enrichment. Do not force basemaps, audio/video bytes, or non-geolocated indicators into CoT points or onto one event subject.
- Define three independent controls: collection enabled, TAK publication enabled, and per-client layer visibility. Turning off TAK routing must not implicitly stop archival collection; hiding a layer must not change another operator's view. Pod scaling alone is not a feed-control API, particularly for shared workers and GitOps-managed replicas.
- Persist desired controls on the fabric and report applied state separately. Reconcile after restart, authorize and audit changes, and define the behavior of queued/expired observations and already-visible markers. A control request accepted is not proof that the worker applied it.
- Require per-stage evidence: last source success, validation/quarantine, durable bus publication, projector forwarding, TAK receipt, and selected-client rendering. A pod readiness check, a JetStream ACK, a TCP write, and map rendering prove different things. Never promise blanket end-to-end exactly-once delivery or immediate remote marker removal.
- Read historical architecture snippets critically. Validate commands, API request formats, dedupe keys, DLQ routing, subject/schema compatibility, TLS verification, and restart persistence against current official docs and actual code before reusing them.

## UI Selection Without an Engine Migration

- Separate pipeline authoring, bus administration, operator feed controls and map rendering. Do not recommend replacing the collector solely to obtain an interface.
- When the owner asks for a UI to see/monitor Redpanda Connect, prioritize provider/third-party Connect-specific interfaces and existing observability—not another NATS UI or a proposed Node-RED/FlowFuse dashboard. He already has NATS UIs. A monitoring request is not permission to deploy or migrate runtimes. Show an actual provider screenshot or accessible UI, clearly distinguished from live OneMind data.
- Redpanda publishes the `Redpanda Connect` Grafana dashboard, catalog ID **22157**, owned by **Redpanda Data, Inc.**; the official Connect Prometheus/Grafana setup guide links it. Inspect/download the real dashboard JSON before claiming compatibility. Its revision-1 template has no stream selector and includes `batch_created`, which some pipelines do not expose; preserve valid panels and adapt missing metrics, units and per-stream filters instead of treating every empty panel as collector failure. Connect streams-mode metrics can carry a `stream` label. Prometheus scrape availability is separate from an HTTP 200 metrics endpoint.
- Check existing Grafana/Prometheus/Loki and current scrape targets before adding infrastructure. Probe Connect's read-only `/metrics`, `/streams` and `/ready` endpoints where authorized. Dashboard counters prove specific Connect stages, not TAK/client receipt. Grafana monitoring does not automatically supply source on/off controls or a visual pipeline editor.
- Redpanda Console exists, and Redpanda Cloud documentation describes synchronized YAML/visual pipeline editing and monitoring. Do not equate Cloud-managed pipeline support, Kafka Connect management or a UI changelog entry with tested control of a standalone self-hosted Connect streams process publishing to NATS. Verify the supported backend/integration and license for the actual deployment model. Redpanda marketing can lag released documentation; check current docs and source rather than repeating a 'coming soon' claim.
- NUI (`nats-nui/nui`) is a desktop/web NATS GUI for Core messages, JetStream streams and KV entries/history. It is an engineering/admin option, not an automatic OneMind source registry or TAK visibility controller. Use scoped credentials; destructive stream/KV operations must not be exposed as general operator toggles.
- Node-RED's editor authors flows; FlowFuse Dashboard (`@flowfuse/node-red-dashboard`, Dashboard 2.0) supplies an open-source operator dashboard with widgets. Use it as a potential control/telemetry frontend without forcing all collectors into Node-RED. Source/TAK/local-view semantics and secure NATS/API integration still require implementation and acceptance tests. Do not use the deprecated original `node-red-dashboard` package for new work.
- When exposing Connect runtime management, inspect the pinned binary's help/docs. Newer Benthos-core releases gate streams mutation endpoints behind `--bind-http` and dynamic-component endpoints behind `bind_http`; do not assume an open HTTP listener exposes them. Keep mutation APIs authenticated and private and define GitOps versus runtime-state ownership.

## Choosing a Feed Runtime

- Separate ingestion/transformation from the NATS broker. Redpanda Connect can publish directly to NATS; replacing Connect is not a broker migration, and its presence is not evidence that a Kafka/Redpanda broker is deployed.
- Audit live Deployments and Argo CD provenance before recommending another workflow platform. Existing n8n may be sufficient for scheduled REST polling; do not add Node-RED solely because historical inventory lists it.
- For n8n/NATS integration, inspect the primary repository `synadia-labs/n8n-nodes-synadia`: it documents Core NATS and JetStream publishing/triggers plus KV/Object Store. It is a community-node package, not a built-in n8n node. Verify installation, deployed-version compatibility, TLS/auth, and actual ACK/retry behavior before promising delivery guarantees.
- Consider Node-RED for a dedicated visual event-driven feed editor. NATS support uses contributed nodes or custom nodes; audit their actual implementation, maintenance, and JetStream semantics rather than treating a palette listing as production proof.
- Keep long-lived/high-volume source connections and specialized codecs in small NATS-native adapters where that reduces operational complexity. GTFS-Realtime Protobuf decoding and satellite propagation are not replaced merely by adding HTTP/JSON workflow nodes.
- When proposing normalized observations followed by shared TAK projection, label this as a target architecture. Inspect the deployed projector's accepted schema/subjects before claiming compatibility; use current contracts for a bounded migration pilot.
- Assess replacements with representative payload tests, source timeouts, NATS disconnect/reconnect, duplicate suppression, restart recovery, and end-to-end sink verification. A ready pod or successful workflow execution is not proof of successful delivery. Preserve the previous manifests and an approved rollback path until the replacement passes.
- Test KV watching beyond the initial snapshot with the exact pinned NATS client version. Initial-replay markers and async-iterator termination semantics vary; changing `break` to `continue` is not by itself a proven fix. Exercise updates, delete/purge, reconnection, startup gating and applied-revision reporting.
- ACK archival messages only after confirmed storage commit under a defined duplicate/idempotency policy. A durable consumer configured with ACK policy none and RAM-only discard-on-error batches is not a loss-resistant archive. A Core subscriber does not gain replay just because a matching stream stores the subject.
- Dedupe by source-appropriate event/revision identity, not only entity UID. Repeated identities can carry valid moving positions or revised/cancelled alerts. Parse timezone offsets correctly, distinguish missing coordinates from valid numeric zero, and do not substitute `(0,0)` for unknown geometry.

## Delta Filtering (Mandatory)

Every feed must implement delta filtering to avoid NATS/ClickHouse bloat. Only republish if moved more than delta_meters OR stale window approaching.

## Mission Routing

Add `<marti><dest mission="Land"/></marti>` to route to specific TAK mission.
Missions: Air, Sea, Land, Hazards, Intel, Drones, Team, OneMind.

## NATS Connection (production prerequisites)

Discover the live account, authorized subjects and stream bindings before connecting. Mount the trusted CA, use a server URL matching the certificate SAN, enable certificate validation, and use secret/credential references without logging their values. Do not copy historical `CERT_NONE`, `check_hostname=False` or `skip_cert_verify: true` examples into a production collector.

Use JetStream publishing with PubAck when persistence is required and a matching authorized stream is configured. Use Core publication only where its transient semantics are intentional; publishing CoT on Core is not proof of TAK receipt or map visibility. Never claim a stream/consumer on one account protects an unrelated leaf/account route.

## Secrets

Feed API keys live in a SealedSecret beside the workload (`platform/omos/`), referenced with `secretKeyRef`. Never print their values.

## Existing feeds

The set of feeds, their replica counts and which are live changes. List them from the cluster and `platform/omos/` (`kubectl -n omos get deploy`) rather than trusting any written inventory.

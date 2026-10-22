# Backstage service coverage workstream

This draft prepares an offline Component coverage report for Phase 6. It adds no
catalog mappings and certifies no all-app integration, live health, authenticated
AI access, recovery, migration or deployment. The coordinator owns the roadmap
update and PR preparation. No production catalog or existing audit file changes.

## Run against a complete export

From `services/onemind-portal`, with Node 22:

```sh
node scripts/report-service-coverage.cjs complete-entities.json
# Or pipe an already obtained complete entity array into stdin:
node scripts/report-service-coverage.cjs - < complete-entities.json
node --test scripts/test-service-coverage.cjs
```

Use the same complete `/entities` JSON array convention as
`scripts/audit-catalog.cjs`. Query envelopes (`items`, cursors), source YAML,
Kubernetes snapshots and malformed entities are rejected. A bare array cannot
prove it contains every authorized entity: the operator must ensure all pages
and kinds were included, with no field projection or permission-filtered gaps.
An empty array is accepted but is not evidence of an empty running catalog.
Run the existing catalog audit separately for freshness, errors and integrity.
The new report neither imports nor modifies that audit.

The report has one sorted row per canonical `component:namespace/name`;
namespace defaults to `default` for catalog identity only. It never defaults the
Kubernetes workload namespace. Noncomponents contribute to target resolution and
counts, not coverage rows. Duplicate canonical IDs are rejected rather than
silently combined. Limits: 32 MiB per CLI input, 10,000 entities/exclusions and
10,000 references per relevant array; no silent truncation. CLI success is 0,
invalid input/read/usage is 2; missing coverage does not fail the process.
There is no clock-dependent output. It does not fetch referenced files or URLs.

## Interpretation

Every field reports `available`, `missing` or `not-applicable`, fixed evidence
keys/modes, and `verification: not-tested`. Available means evidence exists, not
that a configuration is correct. Partial annotation sets remain visible through
the individual evidence keys; a link alone is never a configured native match.
No coverage percentage is a readiness gate.

| Field | Evidence considered |
| --- | --- |
| source | Source-location, GitHub project slug, deployment-source annotations |
| owner | `spec.owner` or `ownedBy` relation; owner identity/access is not verified |
| kubernetesSelector | Kubernetes label-selector or Kubernetes ID annotation |
| kubernetesNamespace | Explicit Kubernetes namespace annotation |
| location | Typed location link or dependency targeting a catalog Kubernetes cluster Resource |
| argo | App-name/app-selector annotation or typed Argo link |
| grafana | Dashboard/overview/alert selector annotations or typed Grafana link |
| prometheus | Service/alert/rule/labels annotations, typed link, and separately the three fixed native scopes in `prometheusApi.ts` |
| operations | Namespace/STAC source predicate from `ServicesModule.tsx`, or typed Operations link |
| dependency | Deduplicated `spec.dependsOn`/`dependsOn` relation counts, with present-target count |
| runbook | Typed runbook link or exact existing Operations/Recovery runbook title; generic TechDocs alone is insufficient |
| aiApi | Provided API refs/relations and typed API/OpenAPI/MCP links; API and definition presence counts |

API metadata never proves agent authentication, authorized reads/writes, contract
conformance or endpoint reachability. `authenticatedAccess` stays `not-tested`.
Consumed API refs are not a provided application interface. Missing targets may
reflect incomplete input or permissions and require review, not entity deletion.
The source predicates are intentionally offline mirrors: changes to the frontend
filters or fixed scope IDs require corresponding report review and tests.

Only validated stable entity refs, fixed keys, states and counts leave the
reporter. Annotation values, selectors, link URLs/titles, API definitions,
metadata descriptions and parser snippets are not copied. Error messages are
fixed and do not echo input paths or JSON. Retain reports under the normal
catalog access policy: stable entity IDs themselves are inventory information.

Applicability is unknown by default; absence does not imply external hosting,
retirement or exemption. An optional second JSON file can record separately
reviewed exclusions without changing the catalog:

```json
[{"ref":"component:default/example","fields":["kubernetesSelector","kubernetesNamespace"]}]
```

```sh
node scripts/report-service-coverage.cjs complete-entities.json reviewed-exclusions.json
```

This is a format illustration, not an approved exclusion. Review the rationale
outside the report with the app owner. Unknown refs/fields, duplicate exclusion
records and exclusions conflicting with present evidence are rejected. No
exclusions are bundled or inferred by this change.

## Scoped rollout queue and remaining gates

Prepare one app per subsequent review; preserve existing IDs and ownership. Each
batch begins with source review and a fresh complete authorized export. These
are candidate review scopes, not authorization to deploy or invent mappings:

1. Portal: `component:default/k8s-backstage-backstage-21b99354`.
2. Portal database: `component:default/k8s-backstage-backstage-postgres-17d2f351`.
3. NATS: `component:default/k8s-nats-nats-0dee54ac`.
4. Existing STAC pilot: `component:default/stac-browser`.
5. Select the next single app from the report with its owner; confirm canonical
   identity and backing manifests before proposing any additional mapping.

For each app, review missing fields, target resolution, exact workload selector,
namespace, deployment location and Argo ownership against source. Review dashboard
and alert scopes separately; a namespace-wide Operations page is not proof of
per-app telemetry. Specify the runbook and provided API/authentication contract
without creating credentials, access grants, storage or controllers.

Coordinator gates before a PR: add the required roadmap entry, review these three
files, run the repository quality gate in an authorized environment, and decide
how the standalone test joins CI (existing pipeline files are outside this
worker's scope). `make validate` is not executed by this offline worker because
it includes whole-tree scans and Terraform initialization.

After any separately approved metadata/release change, record exact source/image
revision, Argo/catalog convergence, signed-in native views, read-only machine
behavior and existing anonymous/write denials. Verify data correctness and
recovery separately when relevant. Failures and unavailable states remain open;
no phase completion or all-app coverage follows from generating this report.

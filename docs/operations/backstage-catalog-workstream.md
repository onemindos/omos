# Backstage catalog launcher review workstream

This implementation prepares offline review tooling, not completion of launcher
repairs. Source inspection in the assigned worktree found both app catalog copies
byte-identical and all 16 historical launcher findings still tagged
`launcher-review-required`. This is source evidence, not a new live observation.
No endpoint was contacted, TLS bypassed, DNS changed, credential read, workload
changed or stable entity identity rewritten.

## Offline report contract

Run from `services/onemind-portal`:

```sh
node scripts/report-launcher-review.cjs complete-entities.json
node scripts/report-launcher-review.cjs - < complete-entities.json
node --test scripts/test-launcher-review.cjs
```

Like `audit-catalog.cjs`, input is the complete authorized Backstage `/entities`
JSON array. Query/paginated envelopes, malformed entity metadata and duplicate
canonical identities are rejected. A bare array cannot prove its own completeness:
the caller must not supply a page or permission-filtered subset as complete.
The report is dependency-free and uses only Node built-ins. It never fetches,
probes, follows redirects, edits metadata or infers a workload failure. Exit 0
means a report was generated, including unresolved findings; exit 2 means invalid
input/arguments or an unreadable input. Errors deliberately omit parser excerpts
and input filenames. There is no automatic closure or deployment gate.

Schema version 1 reports canonical lower-case `kind:namespace/name` references
(default namespace `default`), required-tag findings, four other reviewed launchers
in current source, and a fixed 16-identity historical closure ledger. Evidence
references identify annotation keys and the historical review document; arbitrary
annotation/link values are never used as evidence references. Classifications
come only from the four actual source values: `dns-unresolved-from-operator-mac`,
`certificate-untrusted-from-operator-mac`, `timeout-from-operator-mac`, and
`reachable-from-operator-mac`. Unrecognized observations become `unknown` without
echoing their content. Required-tag rows always remain `unresolved`, even if a
reachable observation conflicts with the tag. An absent historical entity stays
`unresolved-entity-absent`; tag removal yields `closure-review-required`, never
an automatic claim of repair. New required-tag entities also appear in findings.

`onemind.io/launcher-reviewed-at` is the source observation timestamp, not the
report-generation time or inventory timestamp. Timestamp-shaped values retain
their exact original precision and timezone. Missing, invalid, impossible calendar
or clock dates, future, stale and fresh observations are distinguished. Freshness
uses a seven-day default; the exported function accepts a positive finite day
threshold and explicit clock for reproducibility. Age equal to the threshold
is fresh. A malformed arbitrary value is redacted and marked invalid, because
printing it could disclose a URL, secret or personal path. Freshness describes
age of an operator-path observation; it does not resolve a finding or certify
current health. The output intentionally omits launcher URLs, descriptions,
free-form review text, local input paths and credentials.

## Historical 16-finding closure ledger matched to current source

Historical basis: [launcher review and repairs](portal-acceptance-followthrough-2026-10-20.md#launcher-review-and-repairs).
All rows below still have the required tag and source timestamp
`2026-10-20T08:02:51.107324+00:00`. The source inspection does not replace that
timestamp with today's date. All are **unresolved**; no closure is inferred from
old counts, working Grafana integration, portal readiness or current pod health.
Line anchors refer to the identical reviewed app sources at implementation time.

| Canonical identity | Current classification | Closure | Source evidence in both copies |
|---|---|---|---|
| `component:default/cloudtak` | dns-unresolved | Unresolved | [production](../../platform/backstage/catalog/apps.yaml#L658), [development](../../services/onemind-portal/catalog/apps.yaml#L658) |
| `component:default/tak-server` | certificate-untrusted | Unresolved | [production](../../platform/backstage/catalog/apps.yaml#L693), [development](../../services/onemind-portal/catalog/apps.yaml#L693) |
| `component:default/onetak` | dns-unresolved | Unresolved | [production](../../platform/backstage/catalog/apps.yaml#L729), [development](../../services/onemind-portal/catalog/apps.yaml#L729) |
| `component:default/fedhub` | dns-unresolved | Unresolved | [production](../../platform/backstage/catalog/apps.yaml#L762), [development](../../services/onemind-portal/catalog/apps.yaml#L762) |
| `component:default/clickhouse` | dns-unresolved | Unresolved | [production](../../platform/backstage/catalog/apps.yaml#L1048), [development](../../services/onemind-portal/catalog/apps.yaml#L1048) |
| `component:default/minio` | dns-unresolved | Unresolved | [production](../../platform/backstage/catalog/apps.yaml#L1089), [development](../../services/onemind-portal/catalog/apps.yaml#L1089) |
| `component:default/postgis` | dns-unresolved | Unresolved | [production](../../platform/backstage/catalog/apps.yaml#L1113), [development](../../services/onemind-portal/catalog/apps.yaml#L1113) |
| `component:default/longhorn` | dns-unresolved | Unresolved | [production](../../platform/backstage/catalog/apps.yaml#L1137), [development](../../services/onemind-portal/catalog/apps.yaml#L1137) |
| `component:default/prometheus` | dns-unresolved | Unresolved | [production](../../platform/backstage/catalog/apps.yaml#L1193), [development](../../services/onemind-portal/catalog/apps.yaml#L1193) |
| `component:default/nats-surveyor` | dns-unresolved | Unresolved | [production](../../platform/backstage/catalog/apps.yaml#L1233), [development](../../services/onemind-portal/catalog/apps.yaml#L1233) |
| `component:default/tileserver-gl` | dns-unresolved | Unresolved | [production](../../platform/backstage/catalog/apps.yaml#L1299), [development](../../services/onemind-portal/catalog/apps.yaml#L1299) |
| `component:default/valhalla` | dns-unresolved | Unresolved | [production](../../platform/backstage/catalog/apps.yaml#L1361), [development](../../services/onemind-portal/catalog/apps.yaml#L1361) |
| `component:default/maputnik` | dns-unresolved | Unresolved | [production](../../platform/backstage/catalog/apps.yaml#L1423), [development](../../services/onemind-portal/catalog/apps.yaml#L1423) |
| `component:default/webodm` | dns-unresolved | Unresolved | [production](../../platform/backstage/catalog/apps.yaml#L1485), [development](../../services/onemind-portal/catalog/apps.yaml#L1485) |
| `component:default/home-assistant` | timeout | Unresolved | [production](../../platform/backstage/catalog/apps.yaml#L5156), [development](../../services/onemind-portal/catalog/apps.yaml#L5156) |
| `component:default/openclaw` | dns-unresolved | Unresolved | [production](../../platform/backstage/catalog/apps.yaml#L420), [development](../../services/onemind-portal/catalog/apps.yaml#L420) |

The earlier four replacement observations (`agent-web-ui`, `element-chat`,
`matrix-synapse`, `agent-monitor`) are separate from these historical 16. Their
current source observations are reachable from the operator Mac at
`2026-10-20T08:06:39.971995+00:00` with no required tag. Their count does not prove
closure of any of the 16. Current required-tag classification totals are 14 DNS,
one certificate and one timeout, derived from the actual paired source.

## Recommendations to coordinator and remaining gates

- Keep this PR described as review tooling. Reconcile shared roadmap/handoff
  wording in the coordinator's owned changes; this worker does not edit either.
  Preserve Phase 6 In progress and Phase 7 Not started. The 16-row ledger remains
  open pending per-identity reviewed evidence, not a batch count reduction.
- The handoff's “Latest scoped release” #336 digest and its older “Prometheus
  not yet deployed” / “topology next” text are historical. The later local
  receipts record #346/#347 Prometheus acceptance and #349 topology acceptance,
  with the newer image digest and OpenCost next. Recommend labeling historical
  sections or consolidating a dated current summary against those receipts,
  rather than treating the oldest next-work statement as current authorization.
- Historical Grafana 401 references in the acceptance review are superseded for
  the scoped native integration by #340–#342 and folder follow-through #344.
  Preserve the historical evidence and its time; those later receipts do not
  close the Prometheus launcher DNS finding or any other launcher path.
- Before review/merge, coordinator runs appropriate repository gates, including
  `make validate`, and records current blockers honestly. This worker adds a
  standalone Node test file; wiring it into shared CI/package/release gates is
  a coordinator change outside these owned files. Existing whole-tree lint
  debt and hosted CI billing reports are historical until freshly checked.
- A subsequent authorized complete catalog export can exercise the CLI against
  processed entities. The paired-source regression validates today's annotation
  facts through a deliberately narrow fixture extractor, not a runtime YAML
  parser or claim that source YAML is a complete processed catalog. No live
  export or deployment was requested/performed by this worker.
- Any actual repair needs separate path-specific, timestamped evidence and human
  closure review. Ordinary TLS verification and intended access boundaries stay
  intact. Do not probe endpoints, bypass certificates, expose services, delete
  entities or change network policy merely to make this report green.

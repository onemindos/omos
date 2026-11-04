# OpenCost namespace allocation foundation

Draft, offline implementation only. No backend route, proxy, UI wiring, permission,
network, storage or deployment change. No live acceptance was attempted. Phase 6
remains in progress; the coordinator owns the separate roadmap update.

## Source and existing integration

Repository source inspected:

- `platform/ecosystem/foundation-applications/ecosystem-opencost.yaml` pins chart
  2.5.32; `platform/ecosystem/values/opencost.yaml` configures cluster identity
  `onemind-cloud` and the existing Prometheus service. It does not establish a
  verified pricing currency, retention, accuracy or observed running image.
- `platform/argocd/platform-openchoreo-observability-modules.yaml` pins
  `finops-opencost` 0.2.0. `platform/openchoreo/observability/finops-values.yaml`
  disables its bundled OpenCost and points the adapter to the existing
  `http://opencost.opencost.svc.cluster.local:9003` and `http://observer:8080`.
  Adapter implementation and a captured allocation response are not tracked here.
- `services/onemind-portal/packages/app/src/App.tsx` registers the native
  observability plugin; its package manifest pins 1.3.1. Backend `src/index.ts`
  registers the native observability backend. The published 1.3.1 package was
  inspected in memory, without installation: `dist/alpha.esm.js` already defines
  `/cost-insights` and a cost summary entity card. Its `ObservabilityApi.esm.js`
  `getCosts` calls the observer's existing
  `/api/v1alpha1/costs/namespaces/{namespace}/environments/{environment}` contract,
  **not** raw OpenCost AllocationJSON. Its client replaces absent CPU/memory cost
  and efficiency with zero; `useCostInsights` also degrades some failed subsidiary
  requests to empty results. Those semantics need explicit review before presenting
  partial observations as complete. Do not create a second Cost Insights page.
- The existing `plugins/plugin-onemind-common-backend/src/costs.ts` and
  `plugins/plugin-onemind-home/src/CostPanel.tsx` concern DigitalOcean account
  balances/usage. They explicitly have no namespace allocation. Keep this data
  separate; no account token or billing integration is needed for this parser.

Official schema references:

- [Pinned chart metadata](https://raw.githubusercontent.com/opencost/opencost-helm-chart/opencost-2.5.32/charts/opencost/Chart.yaml)
  defaults to exporter 1.121.3; this is source evidence, not image verification.
- [AllocationJSON 1.121.3](https://github.com/opencost/opencost/blob/v1.121.3/core/pkg/opencost/allocation_json.go)
  supplies the wire field names and nullable numeric representation.
- [Allocation calculations 1.121.3](https://github.com/opencost/opencost/blob/v1.121.3/core/pkg/opencost/allocation.go)
  define active minutes from start/end and incorporate resource adjustments in
  total cost. Signed adjustment fields are preserved and are never added twice.
- [Allocation API](https://opencost.io/docs/integrations/api/) and
  [examples](https://opencost.io/docs/integrations/api-examples/) document the
  `code`/`data` envelope, namespace aggregation and allocation sets. These are
  allocation estimates, distinct from the cloud-cost reporting API.
- [Pinned native plugin package metadata](https://registry.npmjs.org/@openchoreo/backstage-plugin-openchoreo-observability/1.3.1)
  identifies the exact published package inspected above.

## Pure contract

`services/onemind-portal/scripts/cost-allocation.cjs` exports:

- `validateWindow({start,end})`: positive explicit interval, strict UTC whole-second
  timestamps matching the pinned exporter. Invalid dates, relative expressions,
  offsets, fractional seconds, open/reversed windows are deliberately unsupported.
- `parseNamespaceAllocation(payload, {window, currency?, namespaces?})`: consumes an
  already decoded raw `/allocation` response for `aggregate=namespace`, one set
  covering the requested window. Multiple sets are rejected rather than summed
  without overlap/coverage guarantees. It performs no requests or authorization.
- `observationState({observedAt, now, maxAgeMs, sourceError?})`: explicit clock/TTL
  freshness classification, with no timers or cache. `now` and `maxAgeMs` are
  finite epoch milliseconds and duration milliseconds, respectively; `observedAt`
  is the successful fetch time in the same strict UTC timestamp format.

Example (offline; `response` is an already decoded value):

```js
const { parseNamespaceAllocation } = require('./cost-allocation.cjs');
const result = parseNamespaceAllocation(response, {
  window: { start: '2026-10-20T00:00:00Z', end: '2026-10-20T02:00:00Z' },
  namespaces: ['backstage', 'nats'],
  // Supply currency only after verifying source pricing configuration.
});
```

Identity requires a Kubernetes namespace key matching both allocation `name` and
`properties.namespace`. Special idle/shared/unallocated/unmounted identities are
excluded with a partial-state reason, never renamed into a namespace. A namespace
filter selects exact names; missing requested namespaces remain explicitly missing.
It is not an authorization boundary. No filter means all returned namespace rows,
not an assurance that every cluster namespace or cluster expense was returned.
Multi-cluster identity separation is outside this contract.

Every accepted row must match the requested allocation `window`; active `start` and
`end` must be contained within it. `minutes` must match the active interval.
`totalCost` is the upstream amount accumulated over that window, in the configured
pricing currency. `averageHourlyCost = totalCost / requestedWindowHours`, including
for short-lived workloads; this is an average over the full query window, not an
instantaneous utilization rate, forecast, monthly price or invoice amount.

Numeric strings, infinity, NaN and negative ordinary costs are invalid. Known
resource adjustment fields alone accept signed finite values because upstream adds
these to resource costs. This conservative slice rejects negative net totals and
negative shared/external costs; it does not clamp them to zero. A future extension
needs an actual source fixture and reviewed semantics. Unknown breakdown fields
remain `null` with `breakdownComplete: false`; they never reconstruct an absent
upstream total. Additional upstream fields are intentionally ignored, not certified
valid. Currency defaults to `null`: AllocationJSON itself carries no currency.
An explicit uppercase three-letter currency is a caller assertion, not verification
or conversion. In particular the USD of DigitalOcean billing is not inherited.

States and presentation rules:

| State | Meaning |
| --- | --- |
| `available` | Selected namespace totals/time metadata validated; inspect breakdown completeness separately |
| `unknown` | Missing response/data; no totals inferred |
| `empty` | Successful empty set/range; no zero assumed |
| `partial` | Missing selection, invalid row/metric, excluded special allocation, source warnings/errors, or overflow |
| `invalid` | Unsupported/malformed envelope or set layout |
| `source-error` | Upstream reported failure; no successful observation |

Only a complete selected result exposes aggregate `totalCost` and average rate;
otherwise they are `null`. Valid individual rows may survive a partial response.
`available` describes validation of this response, not telemetry coverage, pricing
accuracy, source freshness, provider reconciliation or an all-clear.

Freshness is separate: successful historical queries can be current observations.
At TTL expiry (`age >= maxAgeMs`) or a new fetch failure, an existing observation
is `stale`. Failure without an observation is `source-error`; absent/future
observation time is `unknown`. A future backend must retain the last successful
fetch timestamp and payload on failure, show stale values visibly, and never reset
freshness or turn failure into zero. The helper does not manage retention itself.

## Remaining integration and acceptance gates

- [ ] Capture a sanitized real response from the approved existing source, verify
  actual exporter image/schema, currency, pricing basis, retention and telemetry
  coverage. Add reviewed regression fixtures; current fixtures are synthetic,
  shaped from the pinned upstream schema, not captured deployment evidence.
- [ ] Inspect the installed FinOps adapter 0.2.0 and observer translation end to end.
  Compare raw allocation semantics with the existing observer cost item contract
  before placing/reusing this parser. A `.cjs` script is a foundation, not a
  backend runtime import or automatically bundled implementation.
- [ ] Map Kubernetes namespaces to OpenChoreo namespace/environment/project scopes
  through authoritative existing metadata. Never equate those identities by name.
- [ ] Reuse the native page/card and API extension points. Review current null-to-zero
  coercion and swallowed subsidiary-request errors; preserve visible unknown,
  partial, source-error and stale states without manufacturing observer fields.
  Keep signed adjustments, allocation window, currency and average-rate labels.
- [ ] Any separately authorized backend work must preserve existing user/machine
  boundaries, validate scopes, use bounded fixed reads and response limits/timeouts,
  and test cache failure/expiry behavior. No endpoint, proxy, RBAC or networking
  addition is included or authorized by this document.
- [ ] Add meaningful backend/native UI tests and wire this offline regression into
  the coordinator's local/Linux gates when runtime integration is prepared.
- [ ] Coordinator adds roadmap receipt, runs repository validation before a PR,
  records any pre-existing failures and obtains normal review. This worker did
  not run `make validate`: its target includes broad YAML scans and Terraform
  initialization, beyond this bounded offline slice.
- [ ] Release/image/GitOps acceptance remains separate from source review. Later
  approved production checks must prove signed-in native display, real scoped
  source values/window/currency, no-data/outage/stale presentation and preserved
  denial boundaries. No live acceptance, deployment or invoice guarantee here.

Offline regression command:

```sh
/Users/zeus/.nvm/versions/node/v22.23.3/bin/node --test services/onemind-portal/scripts/test-cost-allocation.cjs
```

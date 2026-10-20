# Scoped Prometheus alerts in Backstage

Status: scoped production acceptance passed, 2026-10-20 23:30–23:34 UTC.
User requested this after accepting native Grafana folder links (#344/#345).

## Native integration and scope

Use `@roadiehq/backstage-plugin-prometheus` 3.4.1's native `PrometheusAlertStatus`
through Backstage's public legacy-plugin conversion and API module override.
Compatibility bridge 0.5.15 matches this portal's frontend API 0.18.1; the older
0.5.13 resolves an incompatible frontend API generation and is not used.
Do not patch node_modules or replace OneMind's home/catalog/OpenChoreo layout.

[Native plugin reference](https://github.com/RoadieHQ/roadie-backstage-plugins/blob/main/plugins/frontend/backstage-plugin-prometheus/README.md)
describes its alert table and per-entity instance annotations. This host uses a
bounded backend adapter instead of the general-purpose Prometheus proxy/client:
`GET /api/onemind-services/prometheus/{backstage,postgres,nats}`. Existing HTTP auth
requires a real user, not a machine/service token. No arbitrary URL, namespace,
query parameter, PromQL input or Prometheus admin/write API is accepted. No new
credentials, role bindings, cluster controllers or public endpoints are created.

The adapter performs only a fixed internal `/api/v1/rules?type=alert` GET via the
existing 5-second absolute-deadline/2 MB/no-redirect HTTP reader. Validated source
cache/in-flight requests are shared for 15 seconds; failed reads never fall back
to stale success. Strict protocol validation rejects malformed/partial responses;
raw upstream errors, file paths and undocumented fields are not emitted.

Exact namespace plus Deployment/StatefulSet/current-pattern pod labels scope
Kubernetes alerts. Reviewed custom bus rules explicitly include NATS surveyor,
consumer backlog and TAK CoT flow dependencies even without namespace labels.
Namespace-only, service-only arbitrary rules and cluster-wide alerts do not become
service ownership. Generic Kubernetes rule families are candidate coverage, not
proof of full alert/query/application coverage. No firing or pending match is
explicitly **not proof the service is healthy**. Evaluation errors, evaluations
older than three minutes, future clock skew and missing candidate rules remain
partial/not-configured; source/network failure is unknown/unavailable.

Only stable Backstage, PostgreSQL and NATS catalog identities opt into the card.
Paired catalog annotations are preserved by inventory curation. Display observed
timestamp, 30-second refresh, severity/state/affected target/summary and existing
Grafana/Kubernetes/recovery links. Native inactive green rows are disabled; empty
Grafana-managed rules remain distinct from Prometheus rules.

## Acceptance and release plan

- Backend fixtures: firing/pending, cross-namespace/database exclusions, explicit
  service rules, no match, missing coverage, stale/future timestamps, malformed
  protocol, evaluation errors, coalescing/cache and no stale fallback.
- Frontend fixtures: stable identity opt-in, safe path/API use, coalescing,
  wrong-scope/stale failures, native compatibility extension, no-match disclaimer,
  unavailable/error without leaked details and partial coverage with active alerts.
- These tests run in both portable local/Linux pipeline and production Dockerfile.
- Run required exact-head gates, normal source PR merge, Linux image build/push,
  verified registry source/digest/platform, separate reviewed GitOps pin PR and
  natural Argo rollout. No GitHub billing/check bypass or imperative deployment.
- Live acceptance requires both replicas at the new actual imageID, 2/2 ready,
  DB spec/revision unchanged, retained Grafana/PVC and NATS intact, normal signed-in
  native card behavior for all three services and unchanged machine/anonymous/write
  denial boundaries. Fixture outage behavior is not a production outage drill.

Phase 6 stays In progress; Phase 7 not started. Permissioned action draft #339 is
not deployed or authorized by this observation-only feature.


## Source and image release

Source [#346](https://github.com/onemindos/omos/pull/346) merged `cbf0e637`;
exact implementation head `725c5e4e544ef2f2b20344e52993d88e0df708b2` passed required
portable gates and Linux image gates. Native/API/card tests 13, backend source
fixtures 12, HTTP auth/route fixtures 4; catalog tests 13 including scope curation.
The real 167-rule protocol parsed at preflight: 10 candidate families for
Backstage/DB, 15 for NATS; no scoped active instance then. Fixture firing/pending
and outage checks are tests, not production fault injection or an alert delivery drill.

Published source-pinned linux/amd64 image:
`ghcr.io/onemindos/onemind-portal@sha256:e103bf7c0103b615bef2f4b9dfcb603ef952af7d9eb95d498117d1f48a5dc0e6`.
Registry source/platform/digest checks passed. Promotion
[#347](https://github.com/onemindos/omos/pull/347) uses the portable pipeline;
old image remains the reviewed rollback input, not an automatic data rollback.
Receipts: `/tmp/onemind-prometheus-source-check`, `/tmp/onemind-prometheus-publish`,
`/tmp/onemind-prometheus-promotion`. No real credentials were passed to image build;
Docker's pre-existing build-only AUTH_* placeholder warnings are not runtime keys.
Required gates are separate from existing `make validate` whole-tree YAML failure
beginning at `platform/onetak/media-mtx.yaml:86` and GitHub billing-blocked CI.
Live acceptance is still pending the reviewed pin merge and Argo rollout.


## Production acceptance — 23:30–23:34 UTC

Promotion [#347](https://github.com/onemindos/omos/pull/347) merged
`f27bb770b4a4cf90f1e8b61d9e69438720e5345e`. Portable verification reached
`verify_passed`: owning `platform-backstage` Synced/Healthy at that exact revision,
Deployment generation 76 observed and 2/2 ready/updated/available, with actual
imageIDs on both pods equal to the published index digest above:

- `backstage-77d6cd788f-qmdj9`
- `backstage-77d6cd788f-z2mcq`

Normal Argo polling/reconciliation and zero-unavailable rolling update were used;
no force sync, imperative apply/restart or admin/check bypass. PostgreSQL's full
specification and `backstage-postgres-6dcc65dc8` controller revision are unchanged.
Grafana's Deployment UID/template remain unchanged, retained seeded PVC still
Bound, backup schedule enabled; NATS remains 3/3. No observability rule, datasource,
chart or workload-controller ownership changed.

Signed-in native cards were verified for Backstage, PostgreSQL and NATS. Scopes,
source timestamps, 30-second refresh and explicit no-match/not-health disclaimer
appear correctly. Candidate-family counts are 10, 10 and 15 respectively, not a
complete coverage guarantee. Existing Grafana folder/title links and home/catalog
navigation remain present. Native seven-column table uses horizontal scrolling
at this viewport; this does not certify every plugin's responsive/accessibility UI.
The underlying Backstage status badge hides its text from accessibility trees;
textual accessible state presentation remains a bounded upstream/host follow-up.

During the real rolling update, the Backstage card displayed two scoped
`KubePodNotReady` warning rows for the new/terminating pods. Prometheus historical
`ALERTS` samples independently confirm these were **pending**, not firing. They
cleared naturally on the next observation once rollout readiness settled, and the
card then showed the explicit no-match disclaimer. The rule's annotation says
“more than 15 minutes” even while pending: that is its configured summary, not a
claim that the for-duration elapsed or a notification was sent. No production
failure was injected. Firing/error/stale/partial scenarios are fixture-tested;
Alertmanager inhibition/silencing and notification delivery are not certified.

Redacted HTTP request-log counts show real user GET 200 responses on **both**
replicas and collectively all three scope endpoints. Per-replica localhost tests
show anonymous GET 401, existing machine GET/POST 403, existing machine catalog
read 200 and Grafana proxy 403. Kubernetes service-account Secret/exec/patch/delete
remain denied, Deployment read allowed. No machine access expansion or AI action
executor was introduced; #339 remains draft.

The Linux publish receipt is `publish_passed`; pin exact-head required gates
passed at `6468077203c911ddce8c68e94b16856040ec22d1`. Existing whole-tree YAML
lint debt and hosted GitHub billing-blocked jobs remain **not green**. Local
receipts: `/tmp/onemind-prometheus-source-check`, `/tmp/onemind-prometheus-publish`,
`/tmp/onemind-prometheus-pin-check`, `/tmp/onemind-prometheus-rollout`.
Redacted live evidence in `/tmp/onemind-prometheus-ops/`: `live-acceptance.json`,
`api-log-proof.json`, `rollout-alert-history.json`. Screenshots in
`/Users/zeus/.codex/visualizations/2026/10/20/prometheus-native-alerts/`:
`backstage-alerts.png`, `postgres-alerts.png`, `nats-alerts.png`.
These local paths are operator evidence, not guaranteed permanent remote artifacts.

This accepts the scoped observation feature, not all apps/COP coverage or Phase 6
completion. Next meaningful slice is topology/dependency visibility, then OpenCost;
alert-state accessibility and broader service/cluster alert coverage need explicit
follow-up. No service migration or permissioned mutation is implied.

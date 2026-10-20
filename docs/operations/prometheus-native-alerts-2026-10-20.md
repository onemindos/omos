# Scoped Prometheus alerts in Backstage

Status: source implementation and release gates in progress; not yet deployed.
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

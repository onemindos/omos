# Public ingress lifecycle: assessment and reversible migration plan

Status: **planning only / not deployed**. Requested as part of the full Kubernetes
structure work on 2026-10-20. This is a foundation lifecycle/security follow-up,
not an OpenChoreo app migration and not permission to change DNS, public exposure,
identity routing, the DOKS CNI, credentials or billing.

## Why this is a priority

The official community **ingress-nginx** retirement notice ends maintenance in
March 2026, including security fixes. Our 2026-10-20 read-only inventory still
finds that controller and eight explicitly `nginx`-class Ingress resources.
Existing deployment health does not resolve the maintenance risk. This is not a
claim of a demonstrated exploit or retirement of all products named NGINX.
[Official notice](https://kubernetes.io/blog/2025/11/11/ingress-nginx-retirement/),
[Kubernetes security/steering statement](https://kubernetes.io/blog/2026/01/29/ingress-nginx-statement/)

Observed around 16:33–16:40 UTC, context `do-nyc3-omos-cloud`:

- Argo owner `edge-ingress-nginx`, pinned chart `4.15.1` in Git, automated
  self-heal, prune off. Never `helm upgrade` this adopted release out of band.
- Live controller image `registry.k8s.io/ingress-nginx/controller:v1.15.1`
  at digest `sha256:594ceea76b01c592858f803f9ff4d2cb40542cae2060410b2c95f75907d659e1`.
- DOKS-managed Cilium agent image `ghcr.io/digitalocean-packages/cilium:v1.19.8`.
- Existing GatewayClasses `cilium` and `onemind-private`, controller
  `io.cilium/gateway-controller`. Three HTTPRoutes were observed.
- HTTPRoute CRD annotation reports Gateway API bundle **v1.2.1**, managed by DOKS.
  Current upstream docs describe newer API capabilities: do not copy their
  install/upgrade commands into this managed cluster.

## Route inventory and compatibility work

These are **live resource/host observations**, not automatic deployment targets.
Source paths are inspection entry points; compare effective Argo render to live
before changing them. Certificate/key contents were not fetched.

| Live Ingress | Hosts observed | Owner / inspection source | Required parity tests |
|---|---|---|---|
| `auth/authentik-ingress` | `auth.onemindos.dev` | `platform-auth-infra`; `platform/auth/ingress.yaml` | Cloudflare-only source path, large auth headers/buffers, timeout, cookies, OIDC issuer/redirect/JWKS/logout. **Not first cutover.** |
| `matrix/matrix-ingress` | `matrix.onemindos.dev`, `element.onemindos.dev` | No Argo tracking annotation observed; `platform/matrix/base/ingress/manifests/ingress.yaml` is an older `.com` source | Resolve #192/current ownership/domain drift before using source; redirects, 100 MB upload, long-lived traffic and federation policy. |
| `matrix/matrix-wellknown` | `onemindos.dev` | No Argo tracking annotation observed; confirm actual owner/source | Exact `/.well-known` paths, discovery JSON/content types and coexistence with TAK portal's same-host routes. |
| `nats/nats-ws-ingress` | `ws.onemindos.dev` | `platform-nats`; `platform/nats/ws-ingress.yaml` | Authorized WebSocket upgrade, no buffering, long read/send timeout, Cloudflare restriction, publish/subscribe/reconnect; no anonymous bus exposure. |
| `onetak/geo-services-ingress` | `valhalla`, `nominatim`, `overpass`, `tileserver`, `maputnik` under `onemindos.dev` | `platform-onetak`; `platform/onetak/geo-services-ingress.yaml` | Per-host backend/port/path, body size, timeout and Cloudflare restriction; choose one approved noncritical read-only canary, not all hosts together. |
| `onetak/onetak-ingress` | `onemindos.dev`, `tiles.onemindos.dev` | `platform-onetak`; `platform/onetak/ingress.yaml` | Portal versus tile routing, body size, HTTP/streaming behavior, long timeout, Cloudflare source handling and same-host well-known precedence. |
| `onetak/tak-portal-ingress` | `enroll.onemindos.dev` | `platform-onetak`; `platform/onetak/tak-portal-ingress.yaml` | Enrollment/access policy, uploads, error handling and timeout; existing lack of source allowlist is a separate exposure decision, not silently widened or changed. |
| `takserver/fedhub-ingress` | `fedhub.onemindos.dev` | `platform-fedhub`; `platform/takserver/fedhub/fedhub-ingress.yaml` | Admin authorization, certificate upload, redirects and connect/read/send timeout; independently verify current DNS/public exposure. |

27 total Ingress resources were listed. This table covers the eight explicit
`nginx` selections only: also inspect ingress-class annotations/default-class
behavior, controller TCP/UDP ConfigMaps, LoadBalancers, webhooks, generated
routes, external-dns and certificates before deciding the old controller is unused.
Do **not** confuse FedHub's web route with TAK's independent TCP/UDP listeners.

## Preferred assessment target, not a selected live replacement

First test **Cilium Gateway API** because this cluster already runs it. Reuse
supported APIs and platform skills before installing another proxy. Private
OpenChoreo gateway success does not demonstrate public ingress feature parity.
[Cilium Gateway API documentation](https://docs.cilium.io/en/stable/network/servicemesh/gateway-api/gateway-api/)

If the installed provider-managed version cannot meet a required feature safely,
document the exact gap and compare one maintained alternative with its pinned
API/plugin/release support. Do not install multiple gateways “just in case,” take
over the DOKS CNI, or upgrade managed CRDs as an incidental migration step.

Nginx annotations are not portable policy. The assessment must map each to an
implemented feature/configuration, a tested equivalent, or a **blocking gap**:

- Host/path matching, precedence, request size, timeouts, redirects and backend
  protocols, including WebSocket/SSE behavior where relevant.
- Client source identity and trusted-forwarded-header chain through the actual
  DigitalOcean LB/Cloudflare path. An arbitrary client-supplied `X-Forwarded-For`
  is not an allowlist boundary. Verify IPv4/IPv6 and direct-origin denial.
- TLS termination, SNI, certificate issuance/renewal, origin encryption, secret
  ownership and exact Authentik OIDC issuer/callback semantics. Never disable
  certificate verification or move secrets into documentation.
- NetworkPolicy/Cilium identity handling for gateway-to-backend paths. Existing
  private `onemind-private` must stay private; do not reuse its listener/address
  as a public endpoint. Limit allowed namespaces/routes/reference grants.
- External-dns ownership, DNS TTL, actual origin address and collision handling.
  No unintended public exposure of Backstage, cluster administration or NATS.

## Implementation PR queue and gates

- [x] P0 — Record inventory, maintenance notice, ownership gaps and proposal.
      This checklist tick is **planning complete**, not migration complete.
- [ ] P1 — Capture an effective rendered route/annotation/listener matrix and
      current DNS/TLS/exposure evidence. Reconcile Matrix drift before its move.
      Select one candidate and test strategy; human approves public exposure/cost.
- [ ] P2 — Add a separate Argo-owned public gateway canary configuration in its
      own PR, constrained to the approved test route. No new CNI/CRD ownership.
      Validate schemas against **installed** APIs and deny unauthorized sources.
      No production DNS switch. Confirm target budget/capacity first.
- [ ] P3 — Test an approved noncritical read-only host using a controlled route
      before moving production traffic. Check positive and negative access,
      app-level success, logs/metrics and rollback. Retain the existing route/LB.
- [ ] P4 — Migrate in separate small PRs by feature/criticality. Do not move
      Authentik, enrollment, NATS WebSocket or FedHub until their dedicated tests
      pass. Keep TAK core data/schema/federation/listener migrations out of scope.
- [ ] P5 — After all users and listeners are accounted for, approval + evidence
      authorizes old-controller retirement. Remove dependencies/webhooks/DNS and
      stop reconciliation through Git deliberately. Do not merely delete a pod
      that Argo will recreate. Do not delete certificates, volumes or backup data.

Each implementation PR records versions, ownership, observed UTC timestamps,
requests/limits, exposure changes, route tests, user-impact window, rollback and
the exact CI/local validation results. Only P1 assessment is the immediate next
execution task; later checkboxes are unexecuted work.

## Rollback and emergency access

Keep current LB, old route definitions and certificate ownership until the
cutover/rollback test passes and the approved observation window closes. A revert
must restore the earlier Git revision and any authorized DNS record; record TTL
and propagation expectations before moving traffic. No dual reconciliation of
the same resource. No simultaneous controller retirement and application upgrade.

Keep independent private Authentik/Argo/Headlamp access and a tested recovery
path throughout. Public login breakage must not lock out the operator. Retiring
unmaintained software is important, but not permission for an untested big-bang
identity/network replacement.

## Dependencies / non-goals

Linked to the full-stack proposal [#333](https://github.com/onemindos/omos/pull/333)
and the readiness queue [#331](https://github.com/onemindos/omos/pull/331). Neither
this plan nor an installed gateway makes Phase 6, Phase 7 or all app migrations
complete. No live resource, traffic, provider, billing or credential change is
included in this PR.

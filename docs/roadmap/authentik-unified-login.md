# Authentik shared login across OneMind applications

Owner: Zeus. Requested 2026-10-23 (America/New_York).
Status: implementation plan; no login cutover or production acceptance in this PR.
Master PR: [#441](https://github.com/onemindos/omos/pull/441).

## Outcome and scope

Sign in to Authentik once, open OneMind, and move between permitted applications
without typing another password where the application supports SSO. Keep each
application's existing data, account identity and permissions. Show an honest
login status for every app in the OneMind directory, including exceptions.

This is the master implementation and acceptance contract for that work. It
covers native browser SSO, protected web tools, native clients and machine
access. It does not make every protocol accept a browser cookie or give every
signed-in person administrator access. Each implementation package below gets
its own small PR and release evidence. Merging this planning PR changes docs only.

Read [the mission](../MISSION.md), [platform tracker](../../PLATFORM-ROADMAP.md),
[live baseline](../operations/authentik-login-baseline-2026-10-23.md) and
[repository workflow](../../CONTRIBUTING.md). Coordinate with
[zero-trust PR #399](https://github.com/onemindos/omos/pull/399) for NATS
identities, network isolation and authorization; this plan does not replace it.
Keep the approved OpenChoreo/OpenSearch retirement direction in the platform
decision log. Do not rebuild those integrations to complete this program.

## Evidence and the first blocker

Read-only observation on 2026-10-23 at 20:44 EDT (2026-10-24 00:44 UTC), context
`do-nyc3-omos-cloud`, against Git main `08c65b58`:

- Authentik server 2/2, worker 1/1; both Authentik Applications Synced/Healthy.
- `argocd-authentik-oidc-v2` and `onemind-platform-oidc-v1` report Failed with
  `BackoffLimitExceeded`. `platform-ecosystem-resources` compares to `08c65b58`
  while its running operation still targets `f5f4dca0`.
- [#440](https://github.com/onemindos/omos/pull/440) merged the direct application
  lookup fix. A newer comparison revision is not proof its hook ran successfully.
- The browser-facing workloads in the matrix below are Ready where marked
  observed; readiness does not establish successful SSO, account mapping or access.
- Paperclip is prepared in Git, but no Paperclip Deployment/StatefulSet was
  returned. Langfuse likewise had no returned workload. Do not advertise either
  as a currently available SSO application.

S00 must explain and resolve the old operation through supported Argo operation
controls in a separately reviewed foundation repair. Compare desired hook source,
operation revision and actual Job template. Do not repeatedly erase status,
strip finalizers or force-delete hooks to manufacture a green state. Preserve
the diagnosis and separately record successful bootstrap and interactive login.

## Login contract

1. **Native SSO:** an application trusts an Authentik OIDC or SAML provider and
   creates its own user session. Prefer OIDC when the installed app supports it.
   A brief redirect or initial consent may still occur. Session expiry or
   stronger authentication for a sensitive operation can prompt again.
2. **Protected entrance:** an Authentik proxy authenticates browser access to an
   app. It is full SSO only if the app safely accepts that identity or has no
   separate user login. Otherwise label it **Protected; separate app login**.
   Never remove application authentication simply to make a double login disappear.
3. **Native/device access:** WebDAV app passwords, Matrix client tokens, Mumble
   certificates and TAK enrollment remain protocol-specific. Test actual clients;
   an authenticated browser page is not device-login acceptance.
   The shared session belongs to a browser profile on a device. Another browser,
   a phone or an embedded client browser may need its own initial sign-in.
4. **Machine access:** each AI, workflow and service uses its own scoped
   credentials, audience and API authorization. Do not copy human cookies, expose
   Authentik bootstrap tokens to agents or make an agent a foundation admin.

The intended browser flow is Authentik -> OneMind launcher -> app's supported
login endpoint -> Authentik checks the existing session -> application session.
Backstage stays the main workspace; Authentik's application dashboard is an
alternate launcher and recovery convenience. Neither substitutes for app APIs,
data synchronization, fine-grained authorization or agent tool integration.

## Application integration matrix

Every proposed route below still needs exact installed-version compatibility,
owner/source verification and the acceptance receipt. **Configured** means source
exists; it does not mean a fresh cross-app login was verified in this task.
Registry slugs already in use must be preserved.

### Native SSO and identity-aware applications

| App / namespace | Observed or documented state | Integration and next gate | Owning source / implementation package |
|---|---|---|---|
| OneMind Backstage / `backstage` | 2/2 Ready; Authentik sign-in configured in existing source, not freshly browser-tested here | Keep current signed-subject/catalog-user mapping; use as the login-session control app, verify refresh and denials on both replicas | `platform/backstage/`, `services/onemind-portal/`; S00, S08 |
| Argo CD / `argocd` | Server 1/1; OIDC configured, bootstrap failed | Repair hook convergence; preserve `argocd-admins`, separate read-only mapping and local recovery; browser and CLI SSO tests | `platform/values/argocd.yaml`, `platform/ecosystem/resources/argocd-authentik-provider.yaml`; S00 |
| FedHub / `takserver` | 1/1; Authentik OAuth configured in source | Verify `:9443` OAuth entry, redirects, trust chain and `fedhub-admins`; retain certificate-only recovery path on `:9100` | `platform/takserver/fedhub/`, `platform/ecosystem/resources/fedhub-authentik-provider.yaml`; S00 |
| Nextcloud / `nextcloud` | 1/1 plus Ready Redis; native login and active device sync | Supported `user_oidc` backend; explicitly link Authentik subject to existing non-admin `zeus` account and retained files; preserve device app passwords and direct local recovery | `platform/nextcloud/`; S02 |
| Grafana / `monitoring` | 1/1; persistent database, existing portal Viewer integration | Generic OAuth with explicit role mapping; default Viewer, no automatic server-admin; preserve dashboard UIDs, folders, existing service-account tokens and PVC | `platform/values/kube-prometheus-stack.yaml`; S03 |
| Matrix Synapse + Element / `matrix` | Synapse 1/1, Element 3/3; Application OutOfSync | Configure OIDC on Synapse and advertise SSO to Element; first prove safe mapping to existing Matrix IDs, encrypted rooms and native clients; reconcile the owning manifest/secret gap before rollout | `platform/matrix/argo/`, `platform/argocd/platform-matrix.yaml`; S04 |
| LiteLLM UI / `onemind-ai` | 1/1; internal deployment, no ingress observed | Assess installed generic OIDC and entitlement; upstream currently documents free SSO for up to five users. Verify actual version/license before enabling; private UI route only if required, preserve API keys and budgets | `platform/ecosystem/applications/ecosystem-litellm.yaml`, `platform/ecosystem/values/`; S05 |
| Headlamp / `kube` | 1/1; deployment currently has no OIDC arguments | OIDC is supported upstream; first prove DOKS API-server trust or a separately reviewed supported identity bridge, then per-user Kubernetes RBAC. Keep existing access until proven; no shared admin bearer token | `platform/kube/headlamp.yaml`; S07 |
| n8n / `n8n` | 1/1; own account login documented | Native SSO is edition-gated; verify current license/installed feature. Keep separate login unless an eligible supported route is approved; no purchase or license-bypass patch | `platform/n8n/`; S05 |
| WebODM / `onetak` | 1/1 plus Ready broker/DB and NodeODM | Investigate supported auth/plugin route for installed image; link existing projects/users before any SSO claim. Gateway-only remains separate app login if native identity is unavailable | `platform/onetak/webodm.yaml`; S06 |
| TAK Portal / `onetak` | 1/1; private enrollment and public root share the app | Assess OIDC for operator pages or protected private enrollment; preserve public landing behavior and certificate/package downloads. Test QR/mobile enrollment and avoid protecting the public root by accident | `platform/onetak/tak-portal*.yaml`, `services/tak-portal/`; S06 |

The plan deliberately leaves unsupported or paid integrations conditional. An
explicit exception with a reason is a completed assessment, not verified SSO.

The following are callback candidates for implementation review, not newly
registered endpoints. Verify the pinned app's exact route and current provider
before writing it; preserve existing valid callbacks, and use separate client
registrations for isolated test deployments.

| App | Candidate callback at its verified canonical HTTPS origin | Claim / identity gate |
|---|---|---|
| Argo CD | `/auth/callback`, `/pkce/verify`; existing CLI loopback callback retained as required | `openid profile email groups`; exact `argocd-admins` role mapping |
| Nextcloud `user_oidc` | `/apps/user_oidc/code` | Stable subject with explicitly rehearsed existing `zeus` account mapping; no automatic email-based privilege linking |
| Grafana | `/login/generic_oauth` | `openid profile email` plus explicit group/role claims; Viewer default |
| Synapse | `/_synapse/client/oidc/callback` on the homeserver origin, not the Element frontend | Preserve existing Matrix user ID mapping; verify required scopes against pinned Synapse |
| LiteLLM | `/sso/callback` | Version/license-specific user/team/role claim mapping, separate from API virtual keys |
| Headlamp | `/oidc-callback` | Correct issuer/audience accepted by Kubernetes or reviewed identity bridge; app login alone is insufficient |

FedHub uses the redirect defined in its current configuration on the `:9443`
OAuth listener; inspect and preserve it rather than borrowing another app's
callback. Post-logout callback registration is a separate supported-version check.

### Tools to assess behind an Authentik proxy

| App / namespace | Observed entry / owner | Work and required limit |
|---|---|---|
| Hubble UI / `kube-system` | Tailnet ingress `hubble`; DOKS-managed Cilium components | First read-only proxy pilot; own only the reviewed frontend route/proxy. Do not alter DOKS-managed components or broaden network-view permissions |
| OpenCost / `opencost` | Tailnet ingress `opencost`; `ecosystem-opencost` | Protect human UI; preserve Backstage/backend queries, metric scraping and distinct machine authentication |
| HeyNATS / `nats` | Tailnet `heynats`; `platform-nats-tools` | Pilot after Hubble; check whether it exposes credentials, connection settings or publish/admin actions |
| NATS Dashboard / `nats` | Tailnet `nats-dashboard`; `platform-nats-tools` | Separate provider/policy; preserve required WebSocket/SSE traffic; identify mutation capability before granting viewer access |
| NATS UI / `nats` | Tailnet `nats-ui`; `platform-nats-tools` | Same protocol tests; do not turn shared broker credentials into implicit per-person broker authorization |
| NATS monitor / `nats` | Tailnet `nats-monitor`; resolve exact Service/backend owner at S01 | Protect operator HTTP view only; keep broker monitoring/scrapes internal and separate from browser redirect flows |
| Agent Console / `agent-fabric` | 1/1, tailnet `agent-console`; `platform-agent-fabric` | Classify read and command actions; preserve server-side action authorization and audit identity; gateway access alone cannot grant agent-control rights |
| Agent Tap / `agent-fabric` | 1/1, tailnet `agent-tap`; `platform-agent-fabric` | Protect stream visibility; test long-lived streams, reconnect and authorization on reconnect |
| Hermes dashboard / `hermes` | `hermes-legacy` 1/1, tailnet `hermes-agent`; `platform-legacy` | Keep current dashboard authentication until supported trusted identity is proven; separate chat/control from read-only status; keep probe port private |
| Kubernetes Dashboard / `kube` | 1/1, tailnet `k8s-dashboard`; `platform-kube` | Gateway can protect entrance only. Source includes skip/insecure login flags; assess service-account rights in S07, do not present this as per-user Kubernetes SSO or enable more bypasses |

These are S01/S06/S07 candidates, not an assertion that each currently lacks
authentication. Inventory the actual backend and sensitive operations first.
Prometheus and Alertmanager have internal human interfaces as well as machine
endpoints; use Grafana/OneMind for ordinary visibility. Any optional private human
route needs its own assessment and policy; do not expose them just to add tiles.

### Backend, device, pending and retiring services

| Service family | Disposition |
|---|---|
| TAK Server, ATAK/iTAK and federation | Keep certificate/protocol authentication, enrollment and revocation. Authentik can govern approved human portal entry; it cannot replace field-client certificates by a browser session. TAK core ownership is checked separately from overlays |
| Mumble/Murmur + voice bridge | Keep current client identities/certificates and voice bridge credentials. LDAP/external-authenticator research is a separate compatibility task, not promised passwordless browser SSO |
| NATS TCP/leaf/WebSocket | Keep broker identities and subject permissions. Browser gateways are not broker authenticators. Coordinate per-agent credentials with #399; never wrap the public WebSocket protocol endpoint in an interactive redirect |
| PostgreSQL/CNPG, Redis, MongoDB, ClickHouse, OpenFGA, model runtime | Internal data/service credentials and network rules; no new human web exposure. Database or tool access requires an independent scoped API contract |
| MediaMTX, NodeODM, Tile38, tileserver, TAK-ML, feeds/projectors/ingestion | Identify each human/admin surface versus data plane. Preserve camera/stream/tile/worker protocols; authenticate machine APIs separately. Public tiles may intentionally remain public |
| Paperclip | Source prepared in `platform/paperclip/`, no workload observed; preserve #438 installation/recovery gates. Check pinned version's actual upstream SSO support before adding an Authentik provider; an open SSO proposal is not released functionality |
| Langfuse | No workload observed, bootstrap source still references it. Reconcile intended lifecycle before deciding whether to deploy or remove stale identity; supported upstream SSO is not evidence of a deployed app |
| STAC Browser + OpenChoreo | Pilot browser still Ready, generated HTTPRoute exists; approved retirement is recorded. Track removal to acceptance instead of introducing new auth work or resurrecting generated resources |
| OpenSearch + External Secrets | Workloads still observed while approved retirement is pending. No new user dashboard, provider, controller or identity integration in this program |
| CloudTAK / LibreChat and other proposal-only apps | Reuse their existing draft PRs; add identity requirements at their own installation gate. Not counted as hosted/SSO-ready by this baseline |
| Kubernetes operators, controllers, exporters, backups and Tailscale proxies | Workload identities, Kubernetes RBAC and protocol access; not separate end-user login applications |

## Shared implementation design

### Identity, accounts and permissions

- Reuse the canonical Authentik issuer origin in each existing integration;
  proposed new integrations use `https://auth.onemindos.dev` with their own
  provider issuer path. Verify discovery, `issuer`, signing keys and TLS from
  both browser and pod. Do not alternate public/tailnet issuer aliases casually.
- Preserve existing client IDs, slugs, signing-subject modes and memberships.
  Use stable subject mapping for new accounts. Before changing any existing
  mapping, export identifiers and rehearse account linking on an isolated copy.
  Email alone must not silently link privileged accounts or merge users.
- Retain `onemind-admins`, `argocd-admins` and `fedhub-admins` semantics. Propose
  app-specific member/viewer/editor groups in implementation PRs; ordinary
  family/worker membership must not inherit foundation administration.
- For each app, write the group-to-role map and denial cases before enabling
  auto-provisioning. Default-deny unmapped users; app launch visibility is not
  an authorization control. Existing OpenFGA remains an app authorization
  dependency where already used, not a replacement for every app's native RBAC.
- Pilot with Zeus plus a deliberately unprivileged test identity; test revoked
  and unassigned users too. Never use `akadmin` as a shared family or agent login.
- Draft session policy: normal browser SSO session 12 hours; sensitive admin
  reauthentication within 1 hour; MFA required for administrators. Confirm exact
  supported controls and usability before applying. Document each app's own
  session lifetime, refresh, revocation and maximum stale-access window.
- Treat single logout separately: many apps retain their own sessions/tokens.
  Configure supported logout mechanisms, test them and document remaining
  sessions. Never promise that an Authentik logout instantly revokes every
  Matrix token, Nextcloud app password or API key.

### Gateway topology and access boundaries

Prefer native OIDC/SAML for user-aware apps. For tools without a supported
integration, evaluate the existing Authentik embedded/deployed outpost before
adding another deployment. A private candidate path is:

```text
browser on tailnet -> existing private HTTPS hostname -> Authentik proxy
                   -> allowed app upstream
                   -> Authentik login/session and per-app group policy
```

Use a separate application policy per protected tool. Domain-wide forward auth
cannot supply distinct per-app authorization rules. Tailscale ingress is not
ingress-nginx: adding nginx `auth-url` annotations to `ingressClassName:
tailscale` is not an implementation. Select a supported Authentik proxy path,
keep the ingress frontend Service in its required namespace, and explicitly
prove how it reaches the selected outpost/upstream. Do not assume cross-namespace
Ingress Service references work.

Keep current private hostnames/TLS and access boundaries where possible. Map the
old direct route, alternate ingress, LoadBalancer, NodePort, public root and
in-cluster Service path for each candidate. Enforce trusted headers and upstream
network isolation so a client cannot bypass auth or forge a user. Strip incoming
identity headers and trust only the selected proxy. Keep protocol/API routes on
their own explicit authentication, not an unauthenticated catch-all exception.

Test WebSocket upgrade, SSE streaming, idle reconnect, large upload/download,
cookies, CSRF, callback paths, deep links and redirect loops. Restrict return
URLs and exact callback URIs; no broad wildcard callbacks. Existing public
landing/federation/discovery/media paths require an explicit per-path decision.
Adding Authentik must not make an internal app public or break machine traffic.

### GitOps ownership, secrets and bootstrap reliability

- Authentik stays Layer 1 under `platform-auth-helm` / `platform-auth-infra`.
  Extend the existing provider hooks/blueprints with a clearly named owner per
  object; first inventory embedded outpost state. No second controller or
  manual Helm owner. Application configuration stays with its current owner.
- Each integration records app, namespace, canonical URL, public/private paths,
  provider type/slug, client ID, exact callbacks, Secret *references*, required
  claims, group-to-role map, session/logout behavior, owner and rollback commit.
- SealedSecrets remain the committed pattern. Generate unique per-app client
  credentials through the approved workflow; seal both namespace-bound consumer
  references when required. Never print raw API provider responses, client
  secrets, cookies, bearer tokens, user claims or personal data into logs/PRs.
- Make provider reconciliation idempotent: explicit grant types, non-destructive
  property-mapping/redirect updates, direct slug lookup with expected 404 handling,
  pagination where lists are used, bounded retries and deadlines. Fail on real
  authorization/validation errors with sanitized operation/status diagnostics.
  Avoid broad cleanup, deleting existing providers or changing client secrets
  on every reconciliation. Include create/update/repeat/failure test fixtures.
- Separate foundation bootstrap from optional app configuration where practical;
  a stale optional Langfuse hook must not indefinitely block unrelated foundation
  resources. Design the owning Argo Application/sync-wave boundary explicitly;
  do not move resources into multiple owners or broaden prune policies.
- Proxy/outpost deployments need pinned versions, measured resource requests and
  limits, readiness/liveness, policies, placement and any relevant disruption
  strategy. Record incremental capacity/cost and recovery dependencies first.
  No new node, paid plan or global gateway/controller is assumed by this plan.

### Recovery and the home-first mission

Retain independent operator access to Authentik, Kubernetes, Argo CD and the
existing supported local app administrators until recovery is rehearsed.
Record secure recovery paths without placing passwords or recovery codes in
Git. MFA recovery materials and sealing keys must be recoverable outside the
cluster; check existing encrypted Authentik database/media/config backups.
Before extending the identity foundation or changing account/subject mappings,
require a dated successful encrypted Authentik restore receipt for the current
deployment: database integrity, relevant media/config, signing/encryption material
and independent administrator recovery. Rehearse in an isolated instance with
no outbound production integrations. A backup Job completion alone is insufficient.
If evidence is missing, add the recovery work to S00 before that identity cutover;
read-only diagnosis and preparation can continue.

Before linking Nextcloud/Grafana/Matrix accounts, verify a fresh app-correct
backup and isolated restore; account mapping and database changes are not
reversible by an image rollback alone. Test IdP outage without taking the
production IdP offline: use isolated fixtures or a restricted test client, then
document actual production break-glass verification separately.

Today's Authentik is cloud-hosted. Cloud outage does not magically permit new
home logins. Keep existing device/offline operation; require a separately tested
home identity/recovery design before essential home control depends on this SSO.
Do not introduce a second writable identity authority or claim offline login
from cached browser sessions. Recovery must not require the failing portal.

## Ordered implementation tracker

Only the inventory/planning work is prepared here. All implementation and live
acceptance boxes remain open. Open each package sequentially against current
main, update this table and the root roadmap with its PR and evidence. Resolve
the affected app's blockers rather than treating every optional app as a gate
for the entire program.

| ID | Deliverable / source boundary | Dependency | Status / acceptance |
|---|---|---|---|
| S00 | Repair existing Argo/FedHub/portal identity bootstrap and record existing slugs/roles, supported operation recovery and break-glass | First | [ ] implementation; [ ] correct-revision hooks succeed twice; [ ] portal/Argo browser + CLI tests; [ ] FedHub `:9443` OAuth/admin and `:9100` certificate recovery; [ ] encrypted Authentik restore gate before identity cutover |
| S01 | Versioned per-app connection records; private outpost/proxy pilot for Hubble, including route/bypass map and resource budget | S00 | [ ] reviewed topology; [ ] pilot allowed/denied/spoof/outage tests; [ ] rollback |
| S02 | Nextcloud existing-account OIDC and device compatibility | S00, backup/restore | [ ] isolated account-link rehearsal; [ ] same files/shares; [ ] Mac + Android/WebDAV sync; [ ] recovery |
| S03 | Grafana generic OAuth with Viewer default and explicit admin groups | S00, backup/restore | [ ] isolated existing-account linking rehearsal; [ ] existing DB/dashboards/tokens retained; [ ] viewer/admin/denied tests; [ ] portal plugin regression |
| S04 | Matrix/Synapse OIDC, Element web and mobile SSO, account/device encryption preservation | S00, owner/secret reconciliation, backup/restore | [ ] isolated existing-account linking rehearsal; [ ] existing Matrix IDs preserved; [ ] rooms/history/keys/device tests; [ ] protocol paths |
| S05 | LiteLLM version/license proof and private UI SSO; n8n edition assessment and explicit route decision | S00 | [ ] entitlement recorded; [ ] supported integrations tested or exceptions recorded; [ ] API/webhook preservation |
| S06 | Roll proxy pattern to OpenCost/NATS views, agent UIs and Hermes; assess WebODM and TAK Portal native integration | S01, per-app action/identity assessment | [ ] one app per rollout; [ ] streams and data paths; [ ] role denial; [ ] no bypass; [ ] unknown-native exceptions |
| S07 | Headlamp/DOKS trust proof; Kubernetes Dashboard access review with per-user RBAC and independent recovery | S00, supported API identity design | [ ] no shared admin identity; [ ] read/write/Secret denials; [ ] recovery route. Blocked until trust design proven |
| S08 | OneMind app directory and Authentik launch metadata: canonical URLs, group visibility and honest login labels | Accepted app receipts | [ ] stable existing catalog IDs; [ ] deep links/mobile; [ ] accessible statuses; [ ] no tokens in frontend |
| S09 | Service/AI credential and tool-access inventory; link scoped implementation to #399 and #339 rather than human SSO grants | Per-app API contracts | [ ] caller/audience/scope/expiry/revocation table; [ ] independent machine denials; [ ] no foundation write expansion |
| S10 | Cross-app acceptance, session/role revocation, recovery receipt and updated support runbook | Accepted applicable packages | [ ] cold login + cross-app journey; [ ] expiry/logout limits; [ ] restore/recovery evidence; [ ] all remaining exceptions listed |

Paperclip/Langfuse installation or retirement and STAC/OpenChoreo/OpenSearch
teardown remain their own lifecycle work, not prerequisites that this plan
silently executes. Do not mark the entire program blocked because one
edition-gated app cannot provide SSO.

## Acceptance contract and evidence receipt

For **each** app, save a dated receipt under `docs/operations/` with its package,
source/merge SHA, pinned image/chart, owning Application, operation revision,
actual ready replicas, private/public route, tested role, browser/device and
result. Keep provider IDs and credential references only. Redact HARs/screenshots
before committing; never save raw auth traffic in Git.

- [ ] **Compatibility:** installed version/license and authoritative upstream
  settings verified; discovery/issuer/JWKS reachable from browser and backend.
- [ ] **GitOps:** changed-source checks, appropriate behavioral/security fixtures,
  manifest/schema/Helm renders when applicable, secret scan and `make validate`
  recorded. Existing unrelated failures identified precisely rather than waived.
- [ ] **Reconciliation:** correct revision actually ran, hooks bounded/successful,
  repeat reconciliation creates no duplicate provider, account or role binding.
- [ ] **Cold login:** a fresh browser profile reaches the intended app via
  Authentik; bad callback/state/nonce/issuer/audience is rejected by the supported
  library, and the intended existing account is selected.
- [ ] **Shared session:** from that profile visit OneMind -> Nextcloud -> Grafana
  -> Element -> Argo/FedHub (where assigned); no repeated password challenge within
  policy. Record initial consent and intentional admin reauthentication separately.
- [ ] **Authorization:** unassigned user denied; viewer cannot mutate/admin;
  intended admin can perform an explicitly reviewed harmless action. Direct app
  URL, alternate backend path, spoofed identity headers and new/revoked users tested.
- [ ] **Session lifecycle:** refresh, expiry, group removal, user disable, logout
  and stale-session bounds tested. Record per-app revocation delay and residual
  device/API sessions; do not label single logout passed from one redirect.
- [ ] **Session policy:** administrator MFA and the approved session/reauthentication
  intervals are enforced, including a negative test for a stale admin session.
  If the proposed 12-hour/1-hour intervals change during review, record the
  approved values and test those exact values instead.
- [ ] **Application continuity:** original data/IDs/permissions present; native
  devices, APIs, webhooks, discovery/federation, WebSocket/SSE and upload flows
  appropriate to that app still work. Live deny tests avoid destructive mutations.
- [ ] **Machine separation:** existing scoped backend integrations still work;
  human login does not grant machine tools or foundation writes. Anonymous and
  wrong-audience machine callers are denied.
- [ ] **Recovery:** old independent access verified, IdP-unavailable behavior
  documented, rollback rehearsed and relevant isolated restore accepted.
- [ ] **Directory:** canonical app link, authentication mode, access role and
  actual evidence status visible; conditional apps never presented as verified.

An app can be **Not assessed**, **Prepared**, **Configured**, **SSO verified**,
**Protected; separate app login**, **Protocol credentials**, **Blocked** (with
reason) or **Retiring/not deployed**. These are proposed directory states; no
existing UI schema is changed in this planning PR. Healthy Kubernetes pods or
an Authentik tile cannot by themselves set **SSO verified**.

## Rollback and stop conditions

Stop the affected rollout on duplicate accounts, changed ownership/permissions,
lost file/chat history, failed refresh, widened admin rights, direct-route bypass,
broken enrollment/sync/webhooks or loss of recovery access. Preserve evidence and
restore the affected app's reviewed previous configuration through its owner.
Do not roll back the whole cluster or delete app/identity databases.

Native SSO rollback restores previous auth configuration and tested local access;
it may also require an app-correct data restore if account mapping changed state.
Proxy rollback restores the previous *private* ingress backend and policy only
after its security posture is reviewed; removing protection must not silently
expose an app. Secrets/certificates and devices remain independently revocable.
Keep the pre-change Secret reference and data backup until acceptance. No
namespace/PVC/backup deletion or retirement is part of this login program.

## Primary references

Checked 2026-10-23; recheck against installed versions when implementing.

- [Authentik proxy modes](https://docs.goauthentik.io/add-secure-apps/providers/proxy/)
  and [per-app versus domain forward auth](https://docs.goauthentik.io/add-secure-apps/providers/proxy/forward_auth/).
- [Nextcloud integration and existing-account mapping](https://integrations.goauthentik.io/chat-communication-collaboration/nextcloud/).
- [Grafana integration](https://integrations.goauthentik.io/monitoring/grafana/).
- [Matrix Synapse integration](https://integrations.goauthentik.io/chat-communication-collaboration/matrix-synapse/).
- [Argo CD SSO](https://argo-cd.readthedocs.io/en/stable/operator-manual/user-management/).
- [Headlamp OIDC requirements](https://headlamp.dev/docs/latest/installation/in-cluster/oidc/).
- [LiteLLM SSO entitlement](https://docs.litellm.ai/docs/enterprise).
- [n8n SSO availability](https://docs.n8n.io/user-management/oidc/).
- [Tailscale operator architecture](https://tailscale.com/docs/kubernetes-operator/concepts/architecture).
- [Paperclip authentication](https://github.com/paperclipai/paperclip/blob/master/docs/api/authentication.md)
  and [upstream SSO proposal](https://github.com/paperclipai/paperclip/issues/3028)
  (a proposal, not a supported-version guarantee).
- [Self-hosted Langfuse authentication](https://langfuse.com/self-hosting/security/authentication-and-sso).

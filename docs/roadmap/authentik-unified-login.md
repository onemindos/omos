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

## Refreshed inventory and sort — 2026-11-08 (read-only, supersedes the matrix below where they differ)

**Status of this plan against the cluster today.** The plan written on 2026-10-23 is
out of date in three ways. (1) Several apps it lists as "to do" are done: Argo CD OIDC
(#392), portal/Backstage OIDC (#402), Grafana OIDC (#757), the FedHub provider (#427),
NUI behind a proxy outpost (#720, #721), and two apps it did not know about,
Karakeep and Langfuse, now sign in through Authentik. (2) Apps it lists no longer
exist or are parked: Nextcloud, Matrix/Synapse/Element, LiteLLM, Paperclip, the
OpenChoreo and OpenSearch pieces, and the HeyNATS/NATS Dashboard/NATS monitor tools
(only NUI remains). Skip S02, S04 and S05-LiteLLM. (3) Apps it does not mention are now
live: Forgejo, Vaultwarden, Karakeep, ThingsBoard, OpenRemote, Gatus, CloudTAK,
STAC browser, Hermes web app. The "Argo/platform OIDC Jobs Failed" blocker in
[Evidence](#evidence-and-the-first-blocker) is from 2026-10-23 and was not re-checked
here. Open #565 (NATS role modes) is a different topic.

Snapshot: `kubectl get ingress -A` plus Deployments, ConfigMaps and env names (no values),
2026-11-08. All 26 tailnet ingresses and 3 public nginx ingresses were reviewed. The
Authentik server/worker run; an Authentik database move is in progress elsewhere, so
nothing below was tested by logging in.

Upstream facts were checked against vendor docs on 2026-11-08 unless marked "verify".

### A) Already on Authentik

| App | URL | How it works today | Gap to close |
|---|---|---|---|
| Argo CD | `argocd.tail717ba0.ts.net` | OIDC (`argocd-cm` oidc.config, PKCE, `argocd-admins` -> admin) | Local `admin` account setting should be re-checked; keep one break-glass login |
| Backstage portal | `portal.tail717ba0.ts.net` | Authentik sign-in (#402); preview at `onemind-preview` | None for sign-in |
| Grafana | `grafana.tail717ba0.ts.net` | Generic OAuth, `authentik Admins` -> Admin, else Viewer | Local admin login still exists (keep as break-glass) |
| Langfuse | `onemind-langfuse.tail717ba0.ts.net` | Generic OIDC; live env has `AUTH_DISABLE_USERNAME_PASSWORD=true` | `AUTH_DISABLE_SIGNUP=false`: confirm only the owner can enroll |
| Karakeep | `karakeep.tail717ba0.ts.net` | Native OIDC, password login disabled, auto-redirect | `DISABLE_SIGNUPS` still `false`: set `true` after the owner's first sign-in |
| NUI (NATS UI) | `nui.tail717ba0.ts.net` | Authentik proxy outpost (`nui-auth`, goauthentik proxy 2026.8.2) in front of the app | Broker credentials inside NUI are still shared, not per person |
| FedHub | `fedhub.tail717ba0.ts.net` | Authentik OAuth provider (#427) on the web UI | Verify end to end after the DB move; certificate path stays |
| Authentik itself | `auth.onemindos.dev` (public) and `onemind-auth` (tailnet) | Own accounts | Owner MFA required before wider rollout |

### B) Native OIDC/SAML is available and free (easy)

| App | URL | Native support | Effort | Risk / note |
|---|---|---|---|---|
| Forgejo 15.0.11 | `forgejo.tail717ba0.ts.net` | OAuth2/OIDC auth source, free; Authentik has an official guide | Small: add provider + auth source in Forgejo admin (or `FORGEJO__oauth2_client__*` env), link existing owner account | Account linking by hand; keep local admin for recovery. Git over HTTPS tokens stay separate |
| ThingsBoard 4.3.1.8 | `thingsboard.tail717ba0.ts.net` | OAuth 2.0 login is in the Community docs (supports generic providers). Verify on the installed CE build: some edition comparisons list "SSO" under PE | Medium: client in the tenant admin UI, mapper from email to tenant user | Devices keep access tokens, MQTT and API keys untouched. Tenant admin stays local |
| OpenRemote (manager + own Keycloak) | `openremote.tail717ba0.ts.net` | Keycloak identity brokering: add Authentik as an OIDC identity provider in the `master` or app realm | Medium | OpenRemote needs its Keycloak; this only changes who can log in. Keep the `admin` bootstrap account. Do not replace Keycloak |
| Vaultwarden 1.37.4 | `vaultwarden.tail717ba0.ts.net` | Native SSO (OIDC) since 1.35 | Small config, high care | The master password is still required (SSO does not unlock the vault). Do not set SSO-only until SSO is proven; a broken Authentik must not lock the vault out. Do this last |
| Gatus | `status.tail717ba0.ts.net` | Native OIDC in the config (`security.oidc`) | Small | Low value (read-only status page); can also use C |

### C) Through the Authentik proxy outpost (works, app keeps its own user)

An outpost fronts the app and Authentik decides who gets in. The app behind it
either has no login or keeps its own. Tailscale ingress cannot use nginx forward auth,
so each app gets its own `nui-auth` style outpost Deployment and the ingress points
at it (the pattern already live for NUI). Memory on the hub-core pool is tight, so
budget requests per outpost and share one outpost for several hosts only after
testing.

| App | URL | Today | Effort / risk |
|---|---|---|---|
| Hubble UI | `hubble.tail717ba0.ts.net` | No login (tailnet only) | Pilot app: read-only, no state. DOKS-managed Cilium, so only the ingress route changes |
| OpenCost | `opencost.tail717ba0.ts.net` | No login | Easy. Backstage and Prometheus use the in-cluster service, not the tailnet ingress |
| STAC browser | `stac-browser.tail717ba0.ts.net` (and `stac`) | No login; STAC API on the other ingress | Easy for the browser; the API ingress (`stac`) needs a machine path. Check the retirement status first |
| Agent Console | `agent-console.tail717ba0.ts.net` | Tailnet only; server-side action checks | Medium: has actions, test streaming |
| Agent Tap | `agent-tap.tail717ba0.ts.net` | Tailnet only | Medium: long-lived streams, reconnect behavior |
| Hermes dashboard | `onemind-hermes.tail717ba0.ts.net` | Own dashboard username/password (`hermes-dashboard-auth`) | Medium. Outpost adds a gate; the dashboard login stays unless header auth is added upstream |
| Hermes web app | `onemind-hermes-web.tail717ba0.ts.net` | Own app login | Same. Native clients and API calls must bypass |
| WebODM | none (ClusterIP `webodm`, no ingress today) | Own Django users | No SSO upstream (no header auth). Gate only when an ingress is added |
| Headlamp | `headlamp.tail717ba0.ts.net` | In-cluster service account (everyone with tailnet access has its rights) | Gate entrance now. Native Headlamp OIDC is not useful on DOKS because the API server will not accept Authentik tokens (DOKS has no custom OIDC setting) |
| n8n 2.42.4 | `n8n.tail717ba0.ts.net` | Own owner account | Community edition has no SSO (docs list SSO as Business/Enterprise, verify on the pricing page). Gate the UI with an outpost; webhook paths must be exempt or kept internal |
| TAK Portal | `enroll.tail717ba0.ts.net` | Own login (and public landing at `onemindos.dev`) | Gate `enroll` only if it is private to staff. Do not touch the public root or enrollment downloads |

### D) Cannot or not worth it

| App | Why |
|---|---|
| TAK Server web 8443/8446, ATAK/iTAK/WebTAK | Client certificate authentication; a browser SSO cannot replace it. Authentik may gate the portal only |
| CloudTAK (`cloudtak`, `cloudtak-tiles`) | Signs in against TAK Server credentials/certificates; wrapping it in Authentik would not create a CloudTAK user. Revisit if CloudTAK adds OIDC |
| Mumble/Murmur | Client certificates and bridge credentials; no browser SSO |
| NATS (4222, leaf 7422, WebSocket `ws.onemindos.dev`) | Broker credentials per client; not a browser login |
| OpenRemote Keycloak master admin, Vaultwarden master password, Argo/Grafana/Forgejo local admins | Keep one local break-glass login per app so a bad Authentik does not lock everyone out |
| ClickHouse, Velero, Kyverno, Falco/Falcosidekick, OpenFGA, CNPG databases | No human web UI deployed (APIs and databases only); use Grafana or Argo for visibility |
| Tile services (`tileserver.onemindos.dev`), MediaMTX | Public data plane for tiles/streams, machine access; not a login surface |
| Element/Matrix, Nextcloud, LiteLLM, Paperclip | Not deployed (parked or retired); no work |

### Proposed order

1. Prerequisite: Authentik database move finished and Authentik verified (login, MFA on owner, existing OIDC clients re-tested). Nothing starts before that.
2. Finish A: Karakeep `DISABLE_SIGNUPS=true`, confirm Langfuse signup limit, test Argo/Grafana/FedHub/portal after the move.
3. B in this order: Forgejo, ThingsBoard, OpenRemote (Keycloak broker), Gatus, Vaultwarden last.
4. C in this order: Hubble (pilot), OpenCost, STAC browser, Agent Console, Agent Tap, Hermes dashboard and web app, Headlamp, n8n, TAK Portal. One app per PR, each with a rollback of "point the ingress back at the app service".

### Progress and lessons — 2026-11-09

Done (sign-in goes to Authentik; each app keeps one local break-glass login):

| App | PRs | State |
|---|---|---|
| Karakeep | #795, #796 | New sign-ups closed (`DISABLE_SIGNUPS=true`); owner already linked |
| Forgejo | #797, #798 | Authentik auth source; local admin kept |
| ThingsBoard | #799, #801 | Authentik OAuth2 (CE supports it); sysadmin kept |
| OpenRemote | #803 | Authentik identity provider in its Keycloak; Keycloak admin kept |
| Langfuse | none | Sign-up left open until the owner's first Authentik sign-in (0 users); then set `AUTH_DISABLE_SIGNUP=true` |

Authentik login gates in front of the tailnet admin UIs (2026-11-10; admins group only; each app keeps its own login):
one Authentik outpost `gate-outpost` and one token, with a small replica of it in each app's namespace
(a `gate-outpost.yaml` plus a sealed token beside the app's manifests). PR A adds the replica and its
NetworkPolicy rules; PR B switches only the Ingress backend, in place. The providers come from the Sync hook
`platform/auth/gate-providers.yaml`; a hook-only change does not auto-sync, so sync `platform-auth-infra` once to the merge SHA.

| App (tailnet name) | PRs | State |
|---|---|---|
| Headlamp (`kube`) | #804, #805, #809, #810, #821, #826 | Gated |
| OpenCost | #830, #831 | Gated |
| STAC browser (the `stac` API ingress is untouched) | #832, #833 | Gated |
| Gatus (`status`) | #834, #835 | Gated |
| Hubble UI | #838, #839 | Gated; the hand-applied Ingress is now in `platform/kube/hubble/` with `Prune=false` |
| Agent Console, Agent Tap | #836, #837 | Gated (one replica for both) |
| TAK Portal (`enroll`) | none | Skipped on purpose: phones self-enroll there (`/`, `/enroll`, then the QR `/download/<token>`), and gating it would break enrollment |
| Hermes dashboard and web app | #840 (gate replica, now removed by #846); #844 (Authentik provider); #845 (HOLD) | Not gated by a proxy. Built-in OIDC login chosen: Authentik side is live (#844); the Hermes env change is in draft PR #845 labelled HOLD and turns on at the owner's "restart now" (see below). Gate replica, its token and the two proxy providers are gone; the Hermes pod was not restarted |

#### Hermes: built-in Authentik login (prepared 2026-11-10, waiting for the owner's restart)

The owner chose Hermes's own OIDC login (dashboard_auth `self_hosted` plugin: authorization code + PKCE, public client)
over the proxy gate, so the gate was removed (#846) and the Ingresses were never changed (no certificate burn).

- Authentik (live, #844): provider `OneMind Hermes` (public client `onemind-hermes`, RS256), application `hermes`, bound to
  `authentik Admins` only (Hermes does not check groups). Issuer `https://auth.onemindos.dev/application/o/hermes/`.
  Redirect URIs: `https://onemind-hermes.tail717ba0.ts.net/auth/callback` and `https://onemind-hermes-web.tail717ba0.ts.net/auth/callback`.
- Hermes (draft PR #845, label HOLD): three env vars on the `hermes` container: `HERMES_DASHBOARD_OIDC_ISSUER`,
  `HERMES_DASHBOARD_OIDC_CLIENT_ID`, `HERMES_DASHBOARD_PUBLIC_URL` (the pod sees plain http behind the Tailscale Ingress, so the
  redirect URI must be declared). The web app (9120) uses `webapp.public_url`, already set. Merging changes the pod template, so
  the Deployment (`Recreate`) restarts Hermes; nothing else triggers it. The basic login (`hermes-dashboard-auth`) stays as the fallback.
- Post-restart check: both login pages show an Authentik button; an admin signs in through Authentik on both; username/password
  still works; a non-admin Authentik user is refused by Authentik.
- Rollback: revert #845 (one more restart). Basic login works throughout.

Lessons from the Headlamp pilot:

- A Tailscale ProxyGroup ingress cannot use an ExternalName backend ("invalid ClusterIP"), so its
  backend must be a Service in the ingress's own namespace. One shared outpost Deployment cannot
  sit behind many ingresses. Use one Authentik outpost (one token) with a small replica in each
  app's namespace, and change only the ingress backend.
- **Never delete or recreate an existing tailnet ingress, or change its first TLS host.** Each
  recreate issues a new Let's Encrypt certificate, and the limit is 5 per name per week. Headlamp's
  name used all 5 on 2026-11-08 and was down until it moved to `kube.tail717ba0.ts.net` (#810).
  For the Hermes UIs this would mean days offline. Switch backends in place, one app per PR, and
  let Argo's auto-sync apply each merge (do not race it with manual syncs).

### Owner actions by hand

- Sign in to each B app once through Authentik to link or create the account, and approve linking Forgejo/ThingsBoard to the existing owner users.
- Karakeep: first OIDC sign-in, then approve `DISABLE_SIGNUPS=true`.
- Turn on MFA for the Authentik owner account.
- Vaultwarden: keep the master password; no change to it.
- Decide whether n8n is worth a paid license (otherwise outpost only).

## Application integration matrix

> Note: the matrix below is the original 2026-10-23 draft. Where it conflicts with the refreshed sort above (apps already done, apps no longer deployed), the refreshed sort wins.

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

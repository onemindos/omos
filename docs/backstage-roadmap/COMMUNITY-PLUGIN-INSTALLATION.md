# Community plugin installation

## Included

Eight frontend plugins are pinned in the application: GitHub Actions, Roadie GitHub Pull Requests, Grafana, community Argo CD, Tech Insights, ADR, Tech Radar and Announcements. Their required backend packages, GitHub OAuth provider and ADR/announcement search collators are registered.

Tech Radar uses committed OneMind data instead of upstream sample data. Tech Insights collects the three built-in fact sets every six hours, retains two weeks and checks recorded group ownership. This is metadata completeness, not an operational certification. Announcements uses the existing authenticated user/admin policy; its global banner is disabled to avoid unsolicited notifications. GitHub review/open-PR home widgets are available through home customization rather than forced into existing saved layouts.

The Argo backend's direct mutation routes are blocked by a module. Its action registry also uses the existing global `readOnly:true`, non-destructive filter. It is not added to `backend.actions.pluginSources`. Use a read-only Argo account anyway; never install an administrator token. Existing custom Operations views remain available until equivalent plugin views have been accepted.

## Activation requirements

These packages are installed in source. They are not yet deployed, and installation does not establish external connections. Existing backstage namespace Secrets contain no new plugin integration credentials.

### GitHub

Create or reuse a dedicated OAuth app. Callback:
`https://onemind-preview.tail717ba0.ts.net/api/auth/github/handler/frame`.

Add `auth.providers.github.production.clientId` and `clientSecret` through sealed server-side configuration. Do **not** add a GitHub sign-in resolver or change the Authentik sign-in page: GitHub is supplementary authorization for the plugin APIs. Private repository access may require OAuth `repo` scope (broader than read-only); inspect the actual consent and plugin controls before activation. If that grant is unacceptable, assess plugin support for a constrained GitHub App instead of granting a broad token blindly. Repository entities need `github.com/project-slug: onemindos/omos` or their actual verified repository.

The Actions plugin can expose workflow controls depending on the release. Authorization alone is not evidence of read-only behavior. Inspect the resolved UI/API paths and enforce the desired GitHub privileges before live activation.

### Grafana

Create a Viewer service account token. Configure the proxy endpoint `/grafana/api` with the real internal Grafana Service URL, `credentials: require`, `allowedMethods: [GET, HEAD]`, and a server-side Bearer token. Keep TLS verification enabled for HTTPS upstreams. The public domain is already configured as `https://grafana.tail717ba0.ts.net`.

Set verified `grafana/dashboard-selector` and `grafana/alert-label-selector` annotations on matching entities. The committed NATS dashboard has `nats` and `onemind` tags; do not assign it indiscriminately. Dashboard embedding additionally needs a working Grafana session/frame policy; API tokens do not authenticate iframe viewers.

### Argo CD

Configure `argocd.appLocatorMethods` with type `config`, instance name, real URL and a dedicated token restricted to applications get. The default is an empty list, avoiding upstream fallback username/password attempts. Add `argocd/app-name` (or selector, not both) and instance annotations only after verifying actual application association. Never set `localDevelopment: true`: upstream uses that setting to disable TLS verification globally.

### ADR and search

The repository is private. ADR readers and search collators need a server-side GitHub reader credential limited to necessary repository contents. Set an absolute `backstage.io/adr-location` URL after that reader is connected; mounted file catalog locations cannot infer a GitHub-relative path reliably. Use MADR-formatted Git records. The plugin is installed but has no fabricated architecture records or unresolved private-reader annotations.

### Secrets and rollout

Keep credential values out of chat and Git. Seal credentials using the existing cluster controller and namespace-specific name. Add only the required secret keys/configuration to the deployment, build through the existing GitHub workflow, and pin the resulting image digest through a release PR. Verify PostgreSQL plugin migrations and signed-in behavior before promoting. This PR does not change the running image or create provider grants.

## Deferred because prerequisites are absent

| Candidate | Reason | Next requirement |
| --- | --- | --- |
| Argo Workflows | Cluster has Argo CD Application CRD but no Workflow/WorkflowTemplate CRDs | Deliberately adopt and operate the workflow engine first; it is distinct from Argo CD |
| Matomo analytics | No Matomo integration/service found in reviewed committed configuration | Decide tracking scope/retention, deploy self-hosted Matomo and configure its URL/site ID |
| API linting | No verified API contract set identified in the reviewed platform/catalog configuration | Publish real OpenAPI contracts and select a matching lint engine/plugin |
| Cost Insights | Billing/allocation adapter absent | Connect reliable cost data; do not ship its example client as actual spending |

No custom replacement is being built for these deferred capabilities. Revisit existing plugins when their inputs are ready.

## Build evidence

- TypeScript compilation passed after registration and radar API override correction.
- Frontend production build passed.
- Backend production bundle passed with `--skip-build-dependencies` after the separately built frontend. The initial automatic dependency build could not find a `yarn` executable in the shell PATH.
- No tests were added or run. No live provider calls, browser acceptance, migrations or rollout have been verified for this batch.
- Yarn reports pre-existing React/Material UI and test-library peer warnings, the app now explicitly supplies React 18 types for the Argo peer requirement; legacy Material UI still reports its non-overlapping React/type ranges.

Maintainer sources and selection criteria are in [PLUGIN-FIRST-PLAN.md](PLUGIN-FIRST-PLAN.md).

## Provider activation on 2026-10-07

- Organization-owned `OneMind Portal Readonly` GitHub App installed only on `onemindos/omos`. Repository Actions, Contents, Commit statuses, Metadata and Pull Requests are read-only; user email read supports authorization. Webhooks disabled. Backend reader uses automatically refreshed installation tokens. A private-content read through Backstage's installed credentials provider returned 200. Authentik remains portal login.
- Grafana service account `onemind-backstage`, Viewer, has a 90-day token; 31 dashboards were readable at provisioning. Backend proxy requires portal authentication and allows GET/HEAD only.
- Default Argo project role `onemind-portal` has only `p, proj:default:onemind-portal, applications, get, default/*, allow`. Its 90-day project token was generated using official CLI core mode. Permission probes returned get=yes, sync=no, delete=no. This role is bootstrap state on the existing default AppProject: preserve/recreate this exact policy during recovery without replacing the project's other rules.
- Provider values are sealed in `backstage-plugin-integrations`. Owner-only local copies are under `~/.config/onemind/plugin-activation`; no plaintext appears in source. Argo/Grafana tokens require rotation before expiration. No recurring automation was created.
- Argo's existing HTTPS mesh endpoint uses the established egress ProxyGroup and exact CoreDNS rewrite. TLS validation stays enabled; Argo ignores operator-owned Service ExternalName.
- A Git-owned ADR records the administrator's plugin-first decision. Runbooks catalog points to that private directory. NATS Dashboard uses the verified Grafana `nats` tag.

The personal OAuth App and temporary fine-grained reader token created during setup should be retired after the GitHub App's live acceptance. Release PR records the pinned image and runtime acceptance.

### Argo ingress correction

The existing unmanaged Argo Tailscale ingress forwarded to port 80 while the server enforced HTTPS, causing a redirect loop. The release adopts that existing ingress into GitOps and changes the upstream port to 443. Tailscale's operator uses its built-in `https+insecure` upstream mode for port 443: traffic to the self-signed Argo backend is encrypted, but that proxy hop does not validate the self-signed backend certificate. The portal validates the Tailscale HTTPS certificate, and Argo server TLS remains enabled. No global Node TLS bypass or Argo `localDevelopment` setting is enabled.

Sources: https://tailscale.com/docs/kubernetes-operator/egress/access-tailnet-service and https://github.com/tailscale/tailscale/blob/main/cmd/k8s-operator/ingress.go.

### Database pool budget

During the production rolling update, the second new replica exhausted PostgreSQL's 100-connection limit while two serving replicas remained available. Set `backend.database.knexConfig.pool` to min 0, max 3, idle timeout 10 seconds, with a 30-second acquisition timeout. The installed Backstage PostgreSQL connector reads this exact `knexConfig` path. Idle pools can drain and catalog concurrency is bounded during surge. PostgreSQL's limit/storage and the two-replica availability policy stay intact; review this budget as traffic and plugin count grow.

Sources: https://backstage.io/docs/tutorials/configuring-plugin-databases/ and https://knexjs.org/guide/#pool.

### Browser API allowlist

The GitHub plugins use their documented direct-browser GitHub APIs with user authorization. Permit `https://api.github.com` in `backend.csp.connect-src` alongside self, and GitHub's avatar host in `img-src` alongside self/data. Do not permit arbitrary external connections or relax script policy. Grafana and Argo data continue through authenticated backend integration routes.

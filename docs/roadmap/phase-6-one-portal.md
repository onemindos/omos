# Phase 6 brief: one portal

Status when written (2026-10-18): not started. Depends on phase 4 (plugin packages) and ideally phase 5.
Tracker: [`PLATFORM-ROADMAP.md`](../../PLATFORM-ROADMAP.md). Read [`HANDOFF.md`](HANDOFF.md) first.

## Goal

Our own portal (`services/onemind-portal`) shows its home screens, its full catalog and the
OpenChoreo tabs (Deploy, and with phase 5 Logs, Metrics, Traces, Alerts), and the temporary built-in
console is turned off.

## Starting state

- Our portal: Backstage 1.55.3, new frontend and backend systems, Authentik sign-in, a custom
  permission policy (`packages/backend/src/permissions.ts`), MCP actions with a restricted machine token.
- OpenChoreo plugins 1.3.x are tested against Backstage **1.51.0**
  (`https://openchoreo.dev/docs/platform-engineer-guide/backstage-plugins/compatibility-matrix/`):
  Node 22/24, Yarn 4.13.x, `@backstage/cli` 0.36.x, exact pins for `frontend-plugin-api` 0.17.0,
  `frontend-defaults` 0.5.2, `plugin-app-react` 0.2.3, `plugin-scaffolder` 1.37.0 (we run 0.18.1, 0.5.6,
  0.2.7, 1.39.0). 1.55 is untested; the doc asks to report results to `openchoreo/backstage-plugins`.
- Built-in console: `platform/openchoreo/console/`, values in `platform/ecosystem/values/oc-control.yaml`
  (`backstage.enabled: true`), secrets `openchoreo-control-plane/backstage-secrets`, bindings in
  `platform/openchoreo/platform/console-bindings.yaml`.

## Steps

1. **Decide the version path** on a branch: try the plugins on 1.55.3 (resolutions as the doc lists
   for the NFS path), or pin to 1.51.0 (`yarn backstage-cli versions:bump --release 1.51.0`). Record the
   outcome in the roadmap decision log.
2. **Install** following `https://openchoreo.dev/docs/platform-engineer-guide/backstage-plugins/installing-into-existing-backstage/`
   (reference wiring: `openchoreo/backstage-plugins` `packages/portal-app` and `portal-backend`, tag for 1.3.x):
   packages `@openchoreo/backstage-plugin`, `...-backend`, `...-catalog-backend-module`,
   `...-scaffolder-backend-module`, `...-auth-backend-module-openchoreo-auth`,
   `...-permission-backend-module-openchoreo-policy`, `@openchoreo/openchoreo-auth`; frontend features
   `openchoreoPluginAlpha`, `openChoreoAppModule`, `openChoreoEntityGroupsModule`; optional observability,
   CI and workflows packs once their planes exist. Add `@openchoreo/*` to `npmPreapprovedPackages`
   in `.yarnrc.yml`.
3. **Auth.** The guide wires an `openchoreo-auth` sign-in provider against the OpenChoreo IdP and forwards
   the user's token as `x-openchoreo-token`; it does not document non-default IdPs. Our IdP is Authentik
   (provider `onemind-openchoreo`, which already lists the portal's redirect URI
   `https://onemind-preview.tail717ba0.ts.net/api/auth/openchoreo-auth/handler/frame`, and has grant types
   `authorization_code`, `refresh_token`, `client_credentials` since phase 2). Keep the portal's existing
   user mapping (`onemind.io/oidc-sub`); the built-in console already proves the Authentik flow.
4. **Permissions.** Only one permission policy can be active. Merge the OpenChoreo policy with
   `permissions.ts` (the `onemind-admins` and read-only rules) and re-test the machine token restrictions.
5. **Config.** `openchoreo.baseUrl` (in-cluster `http://openchoreo-api.openchoreo-control-plane.svc.cluster.local:8080/api/v1`),
   `openchoreo.auth` (client credentials: the `backstage-catalog` account, reuse the Secret keys from phase 2f),
   `openchoreo.features.auth.enabled` and `...authz.enabled` mirroring the cluster (both on),
   `openchoreo.schedule` (integers in seconds), `catalog.rules.allow` including `Template` and `Environment`,
   `openchoreo.defaultOwner`.
6. **Events.** The event forwarder already posts to `http://backstage.backstage.svc.cluster.local:7007/api/events/http/openchoreo`.
7. **Build and release**: PR for code, CI image build, separate PR bumping the image digest in
   `platform/backstage/deployment.yaml`.
8. **Switch off the built-in console**: `backstage.enabled: false`, delete `platform/openchoreo/console/`
   and its Argo app, remove the console redirect URI from `auth-provider-bootstrap.yaml`, retire the
   `backstage-catalog` flow only if the portal does not reuse it.
9. **Clean authorization**: remove the inert chart bootstrap bindings (see roadmap phase 2 findings) with a
   `bootstrap` override (`openchoreoApi.config.security.authorization.bootstrap`; the whole array is replaced,
   so keep the roles and bindings you need).

## Gate

One portal: sign-in through Authentik, our home screens and 555-entity catalog, OpenChoreo entities with
a working Deploy tab, anonymous API requests denied, the AI read-only token unchanged, built-in console gone.

## Risks

- Dependency pins and `@mui` / React versions (see the compatibility page); a downgrade changes the whole portal.
- Mixing legacy and new-frontend wiring duplicates tabs; use the new frontend system only.
- Missing `openChoreoAppModule` gives "No QueryClient set" and 401/"No permissions" errors.
- Feature flags must match the cluster or tabs show "No permissions".

## References

`https://openchoreo.dev/docs/platform-engineer-guide/backstage-plugins/overview/`,
`.../catalog-sync/`, `.../permission-policy/`, `.../troubleshooting/`, `https://github.com/openchoreo/backstage-plugins`.

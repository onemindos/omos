# Phase 6 integration verification

Status: implementation in progress, not production compatible/deployed yet.
The existing OneMind Backstage 1.55.3 remains the host. OpenChoreo 1.3.1 adds
core entity views and observability. Existing branding/home/plugins remain wired.
Do not retire the temporary console or mark Phase 6 complete from build evidence.

## Reference and deviations

- Official existing-host guide: https://openchoreo.dev/docs/platform-engineer-guide/backstage-plugins/installing-into-existing-backstage/
- Compatibility matrix: https://openchoreo.dev/docs/platform-engineer-guide/backstage-plugins/compatibility-matrix/
- Permission composition: https://openchoreo.dev/docs/platform-engineer-guide/backstage-plugins/permission-policy/
- Pinned source: openchoreo/backstage-plugins v1.3.1, commit
  `106dfa4d4c4d4418c974e5f941223e3320e0d3af`.

Upstream tests Backstage 1.51.0. This host retains 1.55.3 and pins one compatible
frontend API cohort; local compilation/bundling is evidence only for that gate.

1. Retain the Backstage verified OIDC authenticator, register `openchoreo-auth`
   against the OpenChoreo Authentik provider, and resolve the signed subject to
   `onemind.io/oidc-sub`. No email fallback, guest access or claim-based admin
   enrollment. Capability pre-cache uses the resolved catalog entity reference.
2. Retain exactly one top-level host policy. Non-admin mutations are denied first;
   OpenChoreo/catalog/scaffolder decisions are then delegated. Host admin membership
   does not override an upstream DENY or promote CONDITIONAL to ALLOW. Unknown
   identity and upstream errors fail closed. Auth/authz disable flags fail startup.
3. Upstream 1.3.1 does not export its contextual ABAC router at the package root.
   `upstreamRouter.ts` isolates its shipped `dist/router.cjs.js` API and rejects any
   other package version. This is an explicit private-API maintenance boundary,
   not a copied CEL implementation. Review it on every package update. The cache
   route only admits `plugin:auth`; contextual queries also enforce the host policy.
4. A host module enforces admin-only REST mutations before the OpenChoreo proxy.
   Machine/service callers are denied there. The existing MCP plugin allowlist,
   read-only tool filter and direct machine catalog-write guard remain unchanged.
5. Replace the upstream global FetchApi override with token forwarding scoped to
   the host backend. Direct IDP-token requests only reach explicit approved API URL
   prefixes, never arbitrary catalog-supplied hosts. Scheme, origin and path segment
   boundaries are checked; token-bearing redirects fail. Keep the IDP token out of
   unrelated external requests. Do not expose the backend-only service credential.
6. Workflows/build-plane/assistant/secret-management packs are not enabled. The
   observability pack uses the plane verified in Phase 5. Event sync starts in
   polling-only mode until event auth/ingress is deliberately reviewed.

## Gates, in order

- [x] Local immutable dependency resolution, full TypeScript and frontend bundle.
- [x] Identity, policy intersection, request boundary, pinned router and credential
      forwarding regression tests (37 tests).
- [x] Real Linux release image with all gates inside the Dockerfile.
- [x] Isolated fixture backend ready, network disabled and in-memory database; no production DB contact.
- [ ] Private candidate using real Authentik and OpenChoreo ready.
- [ ] Real Authentik login: same catalog user and intended ownership groups.
- [ ] Existing home/brand/catalog/TechDocs/settings/operational features preserved.
- [ ] Synced OpenChoreo catalog and pilot Cell/Deploy/Definition views show real data.
- [ ] Pilot runtime logs/metrics visible through the observability pack.
- [ ] Anonymous denial; non-admin mutation denial (including direct REST/context
      endpoint); machine read-only MCP works and direct writes remain denied.
- [ ] Human admin receives only intended OpenChoreo rights, no catalog/FGA bypass.
- [ ] Backups and plugin database migration/rollback compatibility reviewed.
- [ ] Separate reviewed GitOps digest/config promotion; Argo synced, 2/2 ready.
- [ ] Repeat browser/API/security gates live before console retirement.

Preview needs the existing user-client credential in `backstage/openchoreo-oidc`
and the restricted `backstage-catalog` service-client credential. Reuse it through
an explicitly managed secret/bootstrap path, not the human client secret for
background sync. Never copy credentials into source or logs. Production manifests
remain unchanged until the preview and migration gates pass.

CI currently cannot start because of account payment/spending limits. Do not mark
checks green or bypass required checks. `make validate` has existing unrelated YAML
lint failures; report them separately from changed-file/local test results.

Candidate deployment in #312: `https://onemind-portal-candidate.tail717ba0.ts.net`.
It uses the existing catalog ConfigMap, in-memory SQLite, no production PostgreSQL
configuration/password and no mounted Kubernetes service-account token. Its user
client is the existing `openchoreo-oidc` Secret; the restricted existing catalog
service credential is resealed as `openchoreo-catalog`, without rotation. One exact
candidate callback is added through the existing Authentik bootstrap Job. Main
portal image/configuration and the temporary console are not changed by this step.

Verified source `4d99acd971ca5510ba2511ff51631d56f2ed7086`, published digest
`sha256:a640cf8746f0d3eb40348e0cb4817b65491955b56453b9d67133050047f60efb`.
In addition to 37 unit/regression tests, isolated HTTP checks verified readiness
200; anonymous catalog/platform/permission/cache/context endpoints 401; guest
refresh 404; machine catalog read 200, catalog mutation and platform/permission
access 403; MCP tools/list returned six explicitly read-only, non-destructive tools.
These fixture checks do not prove real Authentik/browser or production behavior.

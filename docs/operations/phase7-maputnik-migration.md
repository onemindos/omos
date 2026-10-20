# Phase 7a: Maputnik stateless migration

Status: candidate declarations prepared; not a completed migration until the live
checks below are recorded. The existing `onetak/maputnik` Deployment, Service and
ingress remain unchanged. No storage, secrets, shared-service permissions or image
builds are involved.

## Scope and ownership

The user requested quick wins and a Kubernetes common operational picture on
2026-10-20. This intentionally puts one stateless `onetak` app before the otherwise
unchanged Phase 7 project order. It does not mark the `onetak` project complete.

- `platform-openchoreo-projects` reconciles Component, Workload, immutable
  ComponentRelease, ReleaseBinding and a distinct canonical-path redirect.
- OpenChoreo alone owns the new runtime Deployment, Service, configuration and
  generated prefix route in the existing development data-plane namespace.
- Argo continues to own the old application until a separately verified cutover.
- The candidate uses the exact digest of the already running AMD64 image. This
  digest is not evidence of ARM64/Mac compatibility.
- Readiness, startup and liveness HTTP probes check the container on port 8888.
- Resource budget matches the legacy app: requests 50m/64Mi, limits 200m/128Mi.
- The extra candidate is temporary capacity: remove the old workload only after
  functional verification and the agreed observation/cutover period.

## Why the redirect exists

The existing HTML references `app.0c94b3d092b32f56859e.js` and `manifest.json`
relatively. Opening the bare generated path would resolve them outside the app
prefix. The exact-match HTTPRoute temporarily redirects `/maputnik-http` to
`/maputnik-http/`; OpenChoreo's existing prefix route strips the prefix when
forwarding to the application. A 302 avoids caching a permanent redirect during
the parallel migration. No shared gateway or ComponentType is changed.

Candidate URL (private gateway; HTTP transported over the tailnet):
`http://development-default.apps.openchoreo.100-110-185-30.nip.io/maputnik-http/`.

The legacy catalog's public Maputnik hostname did not resolve from the Mac during
preflight. Do not claim the old public hostname works merely because its Ingress
exists. This PR does not change public DNS or exposure.

## Live acceptance

1. Argo `platform-openchoreo-projects` is Synced/Healthy at the candidate commit.
2. `maputnik-development` reports a healthy release and one Ready runtime pod.
3. Candidate imageID matches the pinned digest; readiness/startup/liveness probes exist.
4. Redirect's namespace/parent/host agree with the actual generated route;
   both are accepted/resolved by Gateway API, and the bare prefix returns 302 with
   the trailing-slash Location, without a redirect loop.
5. HTML, referenced JS and manifest return appropriate content types and HTTP 200.
   JS must not be the SPA's HTML fallback. A browser loads the editor controls.
6. The OpenChoreo catalog provider exposes the new component in existing Backstage;
   its Deploy and runtime observations resolve. Read-only MCP lists the component
   and the release. Record timestamps and do not confuse inventory with health.
7. No new Pending pods or unrelated application regressions. Legacy Deployment is
   still available until cutover; do not remove any PVC, namespace or backup.

## Follow-up cutover

After acceptance/observation, preserve or deliberately repair the intended human
URL, update the curated catalog and runtime relationships in both catalog copies,
then retire the old Deployment through a separate PR. Keep one canonical launcher
and an explicit relation to the generated OpenChoreo entity rather than duplicate
unrelated app cards. Verify the URL again after cutover.

Rollback of this parallel candidate is a Git revert of its declarations. It leaves
the legacy service untouched. Never delete the whole shared data-plane namespace
or the `onetak` Project to roll back one app.

## Evidence

- Preflight: legacy Deployment 1/1, no configured volumes or secrets; image digest
  `sha256:35d3074330b01a55b14cd29c7a0168f09728afabc673b7d262eb70546eed7d75`.
- Current OpenChoreo CRDs support managed Resources and StatefulSet ComponentTypes.
- Live rollout, UI, MCP and browser evidence: pending.

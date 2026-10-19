# Phase 3 brief: pilot one app in OpenChoreo

Status (2026-10-18): in progress; parallel deployment under verification. Tracker row and rules: [`PLATFORM-ROADMAP.md`](../../PLATFORM-ROADMAP.md).
Read [`HANDOFF.md`](HANDOFF.md) first.

## Goal

Prove OpenChoreo saves work before anything important moves: run one small, stateless app as an
OpenChoreo Component next to its existing Argo-managed copy, switch over, and delete the old
manifests. If it is not clearly easier, stop and record why before phase 7.

## Starting state

- OpenChoreo has a `default` project, environments `development`, `staging`, `production`, a
  `default` deployment pipeline, the data plane `default`, and the component type
  `onemind-service` (`platform/openchoreo/platform/onemind-service.yaml`).
- No Component, Workload or ReleaseBinding exists. Apps will be reachable at
  `<environment>-<namespace>.apps.openchoreo.100-110-185-30.nip.io`, path prefix
  `/<component>-<endpoint>`, tailnet only, plain HTTP.

## Candidate

Preferred: `stac-browser` (namespace `onetak`, `platform/onetak/stac.yaml`). Image
`ghcr.io/radiantearth/stac-browser:latest`, port 8080, env `CATALOG_URL=http://stac-api:8082`.
Alternative: `agent-console` (namespace `agent-fabric`).

Check first (read-only): what else references `stac-browser` (the Service is `stac-browser:8085`;
grep `platform/` and the portal catalog), and whether the SPA works under a path prefix. STAC
Browser has a path-prefix setting; if the app breaks under `/stac-browser-http/`, pick a different
candidate rather than fighting it.

## Steps (one PR each, update the tracker in each)

1. **Projects Argo app.** New Application `platform-openchoreo-projects` (pattern:
   `platform/argocd/platform-openchoreo-platform.yaml`) for `platform/openchoreo/projects/`,
   destination namespace `default` (OpenChoreo's namespace is `default`).
2. **Project and Component.** Under `platform/openchoreo/projects/onetak/`: a Project `onetak`
   (`ClusterProjectType` `default`, pipeline `default`), a ProjectReleaseBinding per environment as
   in `samples/getting-started/all.yaml`, a Component with
   `componentType: deployment/onemind-service` (check the exact reference form in
   `openchoreo/samples/from-image/react-starter-web-app/react-starter.yaml` at tag v1.3.1) and a
   Workload with the image, port, an external HTTP endpoint and the env. The Component runs in
   another namespace (`dp-onetak-development-*`), so `CATALOG_URL` must be
   `http://stac-api.onetak.svc.cluster.local:8082`. Pin the image by digest, not `:latest`.
3. **Release and bind.** Create the ComponentRelease/ReleaseBinding for `development` (GitOps: the
   binding is a file). Verify the generated Deployment, Service and HTTPRoute with
   `onemind-kubernetes`, and that `get_resource_tree` / `list_release_bindings` via
   `onemind-openchoreo` show it.
4. **Test.** From a tailnet device: `curl -i http://development-default.apps.openchoreo.100-110-185-30.nip.io/stac-browser-<endpoint>/`
   and open it in a browser. Compare with the existing app.
5. **Cutover.** Point whatever used the old Service at the new one (or keep the old Service name via
   a thin Service if needed), then delete the old manifests from `platform/onetak/stac.yaml`
   (only the `stac-browser` Deployment and Service). Argo prunes them.
6. **Write the verdict** in the roadmap: files before/after, what was awkward, whether to continue.

## Gate

The component is Ready in OpenChoreo, serves traffic through the data-plane gateway, the old
manifests are removed, and the verdict is written. Evidence goes in the roadmap.

## Risks

- Path-prefix routing breaks single-page apps (see above).
- The HTTPRoute has no CORS filter (the Gateway API CRDs here are v1.2.1); `onemind-service`
  already omits it.
- ResourceTypes (`postgres`, `valkey`, `nats` from the defaults) are not needed for this app.
- The ecosystem node is at about 99% CPU requested; new pods may land on hub-core (fine).

## References

- OpenChoreo docs: Developer abstractions, Runtime model, Resource relationships
  (`https://openchoreo.dev/docs/category/concepts/`), GitOps
  (`https://openchoreo.dev/docs/platform-engineer-guide/gitops/overview/`).
- `platform/openchoreo/` (existing layout), `docs/operations/ai-operator-access.md` (read-only tools).

## Execution record (2026-10-18)

- [#286](https://github.com/onemindos/omos/pull/286): Argo project owner and OneTAK
  project bindings. Verified live: Synced/Healthy, OneTAK Ready, all three bindings Ready.
- [#287](https://github.com/onemindos/omos/pull/287): parallel STAC Browser Component,
  pinned image, explicit release and development binding. Argo reconciled it, but
  the first pod failed at container creation because file `mountPath` must be the
  parent directory. The generated path included the filename twice.
- [#288](https://github.com/onemindos/omos/pull/288): correct that directory in release
  `stac-browser-v2`. All four CI checks passed; live verification is pending.
- The existing `onetak/stac-browser` Deployment and Service remain in place.
- In-app browser navigation to the pilot's HTTP nip.io URL returned
  `net::ERR_BLOCKED_BY_CLIENT`; an HTTP check cannot substitute for a visual review.
- No cutover, phase completion or phase-4 work is claimed. The maintenance comparison
  must count both authored app contracts and generated immutable snapshots.

The STAC image expects `SB_catalogUrl`, `SB_catalogTitle` and `SB_pathPrefix`; the
old manifest's `CATALOG_URL` and `STAC_BROWSER_TITLE` are ineffective. The pilot
proxies catalog reads through its own nginx because the SPA runs on the user's
machine, where Kubernetes DNS names do not resolve. The gateway retains the app's
path prefix; nginx rewrites only the upstream catalog base URL in JSON responses.

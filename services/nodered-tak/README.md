# Pinned Node-RED + TAK conversion image

This is an image source, not a published image or a live integration. Node-RED's
supported distribution is its [official container](https://nodered.org/docs/getting-started/docker).
The upstream project does not publish a supported Helm chart; OneMind supplies
ordinary Kubernetes manifests through Argo rather than adopting a community chart.

| Input | Pin |
| --- | --- |
| Official Node-RED + Node.js | `nodered/node-red:5.0.9-24@sha256:e085d4a06591fd0aca2ea2a91364d27f1a0e31e885675f77935e10f95412a963` |
| TAK converter | `node-red-contrib-tak@4.2.2`, npm integrity in `package-lock.json` |
| TAK upstream source / npm gitHead | [`490e3801b76f5737a6e021297b5efad9c9ca5042`](https://github.com/snstac/node-red-contrib-tak/tree/490e3801b76f5737a6e021297b5efad9c9ca5042) |
| Proposed publication | `ghcr.io/onemindos/nodered-tak:5.0.9-tak4.2.2-r1` |

Pins checked 2026-10-24 against upstream releases, the anonymous Docker registry
manifest index, npm package metadata and the tagged source. The lock fixes the
converter's transitive dependencies too. npm lifecycle scripts are disabled.
The dev-only Node-RED package supports a real local authentication/palette smoke
test; production retains the official image's runtime and installs only the
converter dependency tree into `/opt/onemind-nodes`.
Production omits dev, peer and optional dependencies so the converter's optional
Node-RED peer does not install a second runtime; the official base provides it.

`bootstrap.cjs` links the reviewed module into `/data/node_modules`, seeds an empty
flow file only if absent, and refuses to replace a conflicting installed module.
It does not overwrite flows, credentials, settings or context. Settings are loaded
from the immutable image, outside `/data`. Every restart uses `--safe`; a deliberate
authenticated Deploy is required to start flows. The disabled example is bundled
under `/opt/onemind/examples` and in Git, never imported automatically.

The TAK package encodes/decodes CoT; it does not itself enroll with or connect to a
TAK Server. The three converter outputs include binary protocol forms. The example
connects only the JSON output to a local debug node, with manual inject, no repeat,
no TCP/MQTT/NATS/HTTP transport, and a disabled tab.

## Build and verify

```bash
cd services/nodered-tak
npm ci --ignore-scripts --no-audit --no-fund
npm test
docker build -t onemind-nodered-tak:smoke .
node test/image-smoke.cjs
```

Use Node.js 24. npm 12 may additionally require `--allow-remote=all` for the
dev-only upstream Node-RED npm dependency's pinned npm tarball; this is a
per-command install option, not a change to provider or credential configuration.
The real-runtime smoke test creates only loopback HTTP traffic and ephemeral test
credentials in process memory, and deletes its temporary user directory afterward.
The image smoke test requires Docker, checks the actual read-only container and
TAK palette with ephemeral memory-only credentials, then removes its own container.
Neither test connects to TAK or an event bus.

`.github/workflows/build-nodered-tak.yml` builds and tests on pull requests without
publishing. After review/merge, a coordinator may explicitly dispatch the workflow
on `main` with `publish=true`. It tests the image before publication, publishes the
version tag and source-SHA tag, and records an immutable GHCR digest in the run
summary. Publication is restricted to that explicit main-branch dispatch. The
workflow uses the repository's `GITHUB_TOKEN`; no new registry credentials belong
in Git. Confirm GHCR pull visibility/permissions before activation; a private image
needs a separately sealed app-namespace pull Secret and Deployment reference.

Before the deployment change is activated, replace its version-tag reference with
the actual `ghcr.io/onemindos/nodered-tak@sha256:...` reported by the successful
build. No application-image digest has been invented or claimed verified here.
The first commit can be integrated independently to publish this image before the
manual deployment commit. [Deployment, enrollment and recovery gates](../../platform/nodered/README.md).

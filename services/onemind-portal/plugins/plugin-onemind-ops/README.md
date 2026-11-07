# plugin-onemind-ops

The Ops Center page (`/ops`): a CesiumJS globe of CoT tracks, the live bus event stream, a namespace health grid,
the Argo strip and Kubernetes Warning events.

Data comes from the `livebridge` backend plugin: `GET /api/livebridge/snapshot`, then deltas on the four Signals
channels in `src/contract.ts` (a verbatim copy of the backend's contract; change both together). `useLiveFeed`
feeds a pure reducer (`liveState.ts`) that locks onto one backend replica (`origin`), ignores old `seq` values and
re-fetches the snapshot after an origin switch and every 5 minutes.

Cesium is lazy-loaded (the only import of `cesium` is in `CesiumGlobe.tsx`, behind `React.lazy`). Its static
assets are copied to `packages/app/public/cesium` by `scripts/copy-cesium-assets.cjs`, which the app `build` and
`start` scripts, the root `start`/`build:*` scripts and the backend Dockerfile run first. No Cesium Ion token is
used; imagery is the OpenStreetMap tile server. The portal CSP must allow `tile.openstreetmap.org` (`img-src`,
`connect-src`) and `blob:` workers (`worker-src`).

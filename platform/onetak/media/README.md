# onemind-media: OneMind video (FEATURES K1, K2)

A second, hardened MediaMTX for OneMind video next to the old `onetak-media`
(`../media-mtx.yaml`), plus `video-registry`, which links each stream to its entity on the bus.
Applied only by `platform/argocd/platform-onetak-media.yaml` (manual sync). Design:
onemind-xr `docs/research/video-and-plugins.md` (Recommendation 1 and "How it applies" 3–6).

| Concern | How it is handled | File |
|---|---|---|
| Image | `bluenviron/mediamtx:1.21.3@sha256:59e4…586c` (Docker Hub digest resolved 2026-10-29); no `-ffmpeg` variant | `deployment.yaml` |
| Auth | `authMethod: internal`. One user per bus identity (BUS-DESIGN section 7), same name, own password: `mavlink-bridge` and `edge-home` publish, `onemind-client` reads and plays back, `tak-viewer` and `takml` read, `video-registry` uses the API. The only password-less user is `metrics` from 127.0.0.1. No anonymous publish or read on any port | `mediamtx.yml` |
| Passwords | Not in Git. `pass: from-secret` placeholders are replaced by `MTX_AUTHINTERNALUSERS_<i>_PASS` from the `onemind-media-users` Secret (proven on v1.21.3; the unit test `TestDeploymentEnvMatchesUsers` keeps the env index and the user order in step) | `deployment.yaml`, `SEALING.md` |
| Path allowlist | Only `drone/sim-*`, `robot/sim-*` (simulated assets in this run) and `cam/<id>` exist; anything else is refused before auth. `overridePublisher: false` | `mediamtx.yml` |
| Recording | fMP4, 60 s segments, `%z` in the path; drones and robots 7 days, cameras 3 days, on `do-block-storage-retain` (20 Gi) | `mediamtx.yml`, `pvc.yaml` |
| Exposure | Tailscale Service `media` only (RTSP 8554, LL-HLS 8888, WebRTC 8889 + ICE 8189/udp, SRT 8890/udp, playback 9996). API and metrics are never on the tailnet. No Ingress, no public LB | `tailnet-service.yaml`, `service.yaml` |
| Placement | Ecosystem pool (`onemind.io/compute-role: ecosystem`), non-root, read-only root FS, requests and limits | both Deployments |
| Registry | `video-registry` polls `/v3/paths/list` every 2 s; KV `OVERWATCH_VIDEO` key `{entity_type}.{entity_id}`; `evt.video.{entity_id}.online\|offline`; `svc.registry.video.get\|list` | `video-registry.yaml`, `services/video-registry` |
| Backup | None by decision: off-site copies of video cost money (decision 0004 #5). Add a VolSync source only if the owner wants recordings backed up | `pvc.yaml` |

## Sync prerequisites (in order)

1. The owner seals the two Secrets (`SEALING.md`), commits the two `*.sealedsecret.yaml` files
   here and adds them to `resources:` in `kustomization.yaml`.
2. The first `build-video-registry` run on `main` pushed `ghcr.io/onemindos/omos-video-registry`;
   pin its digest in `video-registry.yaml`.
3. The owner approves the 20 Gi recordings volume (billed once synced).
4. Sync `platform-onetak-media` by hand and check: both pods Ready on the ecosystem node; the
   operator created the tailnet node `media`; `MEDIA_HOST` matches its MagicDNS name.

## Owner checks after the first sync

```bash
# From a tailnet device. Passwords come from the keychain or a password manager, never typed
# into shell history: `read -rs P` then use "$P".
ffmpeg -re -f lavfi -i testsrc2=size=1280x720:rate=30 -c:v libx264 -preset ultrafast -tune zerolatency \
  -g 30 -f rtsp -rtsp_transport tcp "rtsp://mavlink-bridge:$P@media.tail717ba0.ts.net:8554/drone/sim-x1"
# Expected: publishing starts. Without credentials: "401 Unauthorized".
# In Safari: http://media.tail717ba0.ts.net:8888/drone/sim-x1/ , sign in as onemind-client:
# the test pattern plays (LL-HLS).
nats req svc.registry.video.get '{"entity":"ent.drone.sim-x1.state"}'
# Expected: "live": true, "entity_subject": "ent.drone.sim-x1.state", tak_url on media.<tailnet>.
```

## Cut-over of the old onetak-media (not in this PR)

The old `onetak-media` (`../media-mtx.yaml`, no auth, floating `1-ffmpeg` tag, hub-core pool) is
applied by the auto-synced `platform-onetak`, so removing it here would prune it on merge. After
the owner has checked onemind-media:

1. Point the takml worker `FEEDS` at `rtsp://takml:<pw>@onemind-media.onetak.svc:8554/<class>/<id>`
   (its password from a takml-namespace SealedSecret) and move its detections from
   `fabric.vision.*` to `det.{domain}.{type}.{id}`.
2. Delete `../media-mtx.yaml` in a separate PR.
3. CloudTAK media stays off (its README "Media replacement gate").

## Unverified until the first sync

- UDP through the Tailscale operator proxy (SRT 8890, WebRTC ICE 8189). RTSP over TCP and LL-HLS
  do not depend on it. WebRTC may also need `webrtcAdditionalHosts` set to the tailnet name.
- MediaMTX as UID 65532 with a read-only root filesystem (the image runs as root by default;
  recordings go to the mounted volume, HLS stays in memory).

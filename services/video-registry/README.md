# video-registry (FEATURES K2, K3)

Links every MediaMTX stream to the OneMind entity that owns it. Video never travels on NATS;
only its URLs, codec, latency class and permissions do (onemind-xr `docs/design/BUS-DESIGN.md`
section 2).

```
MediaMTX /v3/paths/list ──poll 2 s──► video-registry ──► KV OVERWATCH_VIDEO   key drone.sim-x1
   (user video-registry)                    │          ──► evt.video.sim-x1.online|offline (core)
   mediamtx.yml (same ConfigMap) ──────────►│          ◄── svc.registry.video.get | list (RPC)
```

| Path | Entity | KV key |
|---|---|---|
| `drone/sim-x1` | `ent.drone.sim-x1.state` | `drone.sim-x1` |
| `robot/sim-r1` | `ent.robot.sim-r1.state` | `robot.sim-r1` |
| `cam/gate-north` | `ent.camera.gate-north.state` | `camera.gate-north` |

Drones and robots must be `sim-*` in this run (run rule 6); other paths are counted as rejected.

## Record (`v: 1`)

`stream_id`, `entity_type`, `entity_id`, `entity_subject`, `live`, `reason` (`publisher_stopped`,
`mediamtx_unreachable`), `source_type`, `tracks`, `codec`, `width`, `height`, `klv`,
`latency_class` (best class for native/TAK readers: `tactical`), `latency_by_transport`
(`whep` control, `rtsp`/`srt` tactical, `hls` situational), `urls` (`rtsp`, `hls`, `whep`, `srt`,
`playback`; never credentials), `tak_url` (only when it passes the TAK hygiene of
robotics-cot-profile.md section 5), `site_url`, `permissions` (`publish`/`read`/`playback` bus
identities, derived from the mounted `mediamtx.yml`), `recording`, `media` (the entry for the
entity's `media[]`, entity.md section 3), `ready_since`, `last_seen`, `updated_at`, `origin`, `seq`.

Only the entity service writes `ent.*.state`; it is meant to merge `media` from this bucket
(owner proposal K-P2, a follow-up to A4). Until then clients get the stream with
`svc.registry.video.get {"entity":"ent.drone.sim-x1.state"}`.

## `internal/takvideo` (K3)

Renders `<__video uid url><ConnectionEntry …/></__video>` and the `/Marti/api/video` feed body,
and enforces the URL rules: no loopback, link-local, pod, `.svc`/`.cluster.local` or (for the hub)
private address; no `@`, query or fragment. The hub projector (A6, not in omos yet) imports it.
`TestVideoURLHygiene` is robotics-cot-profile.md section 9 test 3.

## Configuration

`NATS_URL`, `NATS_USER`/`NATS_PASSWORD` or `NATS_CREDS`; `MEDIAMTX_API`, `MEDIAMTX_USER`,
`MEDIAMTX_PASSWORD`, `MEDIAMTX_CONFIG`; `KV_BUCKET` (`OVERWATCH_VIDEO`, never created by the
service: the NACK CR owns it); `ORIGIN`; `MEDIA_HOST`, `SITE_MEDIA_HOST`, `MEDIA_*_PORT`;
`POLL_INTERVAL` (2 s); `STALE_AFTER` (30 s; MediaMTX unreachable this long marks every live stream
offline); `HTTP_ADDR` (`:8080`: `/healthz`, `/readyz`, `/metrics`).

`video-registry select ent.drone.sim-x1.state` asks a running registry for an entity's stream.

## Tests

```bash
go vet ./... && go test -race -count=1 ./...
```

| Test | Needs | Proves |
|---|---|---|
| `TestParsePath`, `TestParseEntityRef`, `TestBuildRecord` | — | path to entity, URLs without credentials, no `tak_url` for loopback |
| `TestNoAnonymousAccess`, `TestPermissionsAndRecording`, `TestDeploymentEnvMatchesUsers` | — | the cluster `mediamtx.yml` and `deployment.yaml` |
| `TestVideoURLHygiene`, `TestCheck`, `TestFeed` | — | K3 renderer rules |
| `TestRegistryLifecycle`, `TestRestartMarksVanishedStreamsOffline`, `TestMissingBucketFails` | `nats-server` | KV, events, RPC, idempotence, MediaMTX outage, restart |
| `TestE2EMediaMTX` | `nats-server`, `mediamtx` v1.21.3, `ffmpeg`, `ffprobe` | K1 + K2 against the real MediaMTX running the cluster file: refusals, LL-HLS, a decoded frame, recordings roll, online/offline |

Binaries are found on `PATH` or through `NATS_SERVER_BIN`, `MEDIAMTX_BIN`, `FFMPEG_BIN`,
`FFPROBE_BIN`; a test skips when one is missing. CI (`build-video-registry.yml`) installs all of
them, MediaMTX checksum-verified.

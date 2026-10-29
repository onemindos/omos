# Edge MediaMTX (FEATURES K9)

MediaMTX on the owner's Mac or the home server. LAN devices publish and read here with low
latency; every allowed path is forwarded to the cloud `onemind-media` over SRT as the site
identity `edge-home`, through the tailnet. If the cloud or the link is down, LAN viewing and the
local 24 h recording keep working (the same rule as the site NATS leaf).

| File | Purpose |
|---|---|
| `edge-media.sh` | `run`, `selftest`, `keychain-set` |
| `mediamtx-edge.yml.tmpl` | the edge config; rendered at start into a private temp file with the passwords, deleted on exit |

Credentials never touch Git or argv: the environment (`ONEMIND_MEDIA_EDGE_HOME`,
`ONEMIND_MEDIA_EDGE_PUBLISHER`, `ONEMIND_MEDIA_LAN_VIEWER`) or the macOS keychain (service
`onemind-media`, accounts `edge-home`, `edge-publisher`, `lan-viewer`). The script needs only
bash 3.2 (the macOS default), `mediamtx`, and for `selftest` also `ffmpeg`, `ffprobe`, `curl` and
`openssl`.

## Build machine

```bash
infra/edge/mac-media/edge-media.sh selftest
```

Starts a local "cloud" MediaMTX with the cluster `platform/onetak/media/mediamtx.yml` (only
addresses changed), the edge from the template, and an ffmpeg test pattern as `drone/sim-x1`,
then checks: anonymous publish to the edge refused; the cloud path is fed over SRT
(`source.type=srtConn`); the cloud client and the LAN viewer both read H.264 640x360; anonymous
reads of edge and cloud refused. Prints `SELFTEST OK`.

## Owner (Mac)

```bash
brew install mediamtx ffmpeg                       # or the release binaries
infra/edge/mac-media/edge-media.sh keychain-set    # edge-home from SEALING.md; the other two new
infra/edge/mac-media/edge-media.sh run
# On the Mac: publish the camera as cam/maccam
ffmpeg -f avfoundation -framerate 30 -i "0" -c:v libx264 -preset ultrafast -tune zerolatency -g 30 \
  -f rtsp -rtsp_transport tcp "rtsp://edge-publisher:$(security find-generic-password -s onemind-media -a edge-publisher -w)@127.0.0.1:8554/cam/maccam"
```

Expected: on a second LAN device, `rtsp://<mac-ip>:8554/cam/maccam` (user `lan-viewer`) plays in
VLC or ATAK; on a tailnet device, `rtsp://media.tail717ba0.ts.net:8554/cam/maccam` (user
`onemind-client`) plays the same stream; `nats req svc.registry.video.get
'{"entity":"ent.camera.maccam.state"}'` shows it live. This replaces the old push in
`../mac-leaf/ffmpeg-cam-push.plist`, which targets a host no Service provides.

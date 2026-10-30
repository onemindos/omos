# S24 Ultra — Installed ATAK Plugin Inventory (July 2026)

Confirmed via `adb shell pm list packages | grep -Ei "atak|tak|geocam|ar\.|voice|ptt|mumble"` on device <SERIAL> (ATAK CIV 5.7.0.10).

## Core ATAK
- `com.atakmap.app.civ` — ATAK-CIV 5.7.0.10 (build e0d26f0e)
- `com.atakmap.takcam` — TAK Camera
- `com.atakmap.helmcam` — Helm Camera

## Camera / Video / AR
| Package | Name | Notes |
|---------|------|-------|
| `com.atakmap.android.geocam.plugin` | GeoCam | ATAK camera streaming plugin |
| `com.partech.geocamera` | PAR GeoCam | GeoCam engine — streams RTSP to VCM |
| `com.atakmap.android.arvideooverlay.plugin` | AR Video Overlay | Overlays COP onto camera (screen only, no external API) |
| `com.atakmap.android.viscacam.plugin` | VISCA Cam | PTZ camera control |
| `com.atakmap.android.videocollections.plugin` | Video Collections | VCM video feed manager |
| `com.atakmap.android.vns.plugin` | VNS | Video/network streaming |
| `com.atakmap.android.vans.plugin` | VANS | Video/audio network service |
| `com.google.ar.core` | ARCore | Base AR SDK (system) |
| `com.google.ar.lens` | Google Lens AR | — |
| `com.par.imagemarkup` | PAR Image Markup | Image annotation tool |

## Voice / PTT
| Package | Name | Notes |
|---------|------|-------|
| `com.atakmap.android.taktalk.plugin` | **TAK Talk** | ✅ Native ATAK PTT via Mumble protocol — PRIMARY for Legacy PTT |
| `com.atakmap.android.wave.plugin` | **WAVE** | Motorola WAVE PTT (requires WAVE subscription) |
| `com.beartooth.beartoothtakplugin` | Beartooth | Off-grid mesh radio PTT |

## Communications / Chat
| Package | Name | Notes |
|---------|------|-------|
| `com.atakmap.android.takchat.plugin` | TAK Chat | GeoChat — wire to Legacy via NATS |
| `com.atakmap.android.icet.plugin` | ICET | Interoperability comms tool |

## Drone / UAS
| Package | Name | Notes |
|---------|------|-------|
| `com.atakmap.android.uastool.plugin` | UAS Tool | Full MAVLink drone ops |
| `com.atakmap.android.uastool.groundsdk8` | UAS Tool (GroundSDK8) | Parrot/Anafi SDK |
| `com.atakmap.android.uastool.autelmavsdk` | UAS Tool (Autel) | Autel MAVLink |
| `com.atakmap.android.atakgo.dji` | ATAKGo v4 (DJI) | DJI integration |
| `com.atakmap.android.atakgov5.dji` | ATAKGo v5 (DJI) | DJI v5 |
| `com.atakmap.android.dronehone.plugin` | Drone Hone | RF drone detection |
| `com.atakmap.android.wasp.plugin` | WASP | Wide-area surveillance planning |

## Situational Awareness / Data
| Package | Name | Notes |
|---------|------|-------|
| `com.atakmap.android.adsbdirect.plugin` | ADS-B Direct | Live aircraft on map |
| `com.atakmap.android.winddata.plugin` | Wind Data | Wind/weather overlay |
| `com.atakmap.android.compassnav.plugin` | Compass Nav | Navigation overlay |
| `gov.tak.taktracker` | TakTracker | Lightweight blue force tracking |
| `com.somewearlabs.swtak.plugin` | SomewearTAK | Satellite comms bridge |

## Planning / Intelligence
| Package | Name | Notes |
|---------|------|-------|
| `com.atakmap.android.reports.plugin` | Reports | 9-line, SALUTE, SPOT forms |
| `com.atakmap.android.grgbuilder.plugin` | GRG Builder | Range cards / sector sketches |
| `com.atakmap.android.blockgenerator.plugin` | Block Generator | Tactical shape generator |
| `com.atakmap.android.firesurvey.plugin` | Fire Survey | Fire area mapping |
| `com.atakmap.android.datasync.plugin` | Data Sync | Enterprise sync |
| `com.atakmap.android.cot_utility.plugin` | CoT Utility | Raw CoT manipulation |

## Suspicious / Flagged
| Package | Installed | Notes |
|---------|-----------|-------|
| `com.atak.nightvision` | 2026-07-16 | Non-standard namespace (lacks `atakmap.android`). Cannot inspect signing on Android 14 without root. Monitor for unusual network activity. |

## Total: 33 packages (3 core + 30 plugins)

---

## Key Integration Notes for OneMind

### Legacy PTT via TAK Talk (PRIMARY PATH)
- TAK Talk uses **Mumble protocol** — needs a Murmur server pod on DOKS
- Deploy Murmur on DOKS → configure TAK Talk to connect → Legacy joins same channel via `mumble-py` or LiveKit bridge
- TAK server does NOT include a built-in Murmur server
- Alternative: use `om_presence` + LiveKit (already deployed) as the PTT bridge, bypass native TAK Talk

### GeoCam → AI Vision
1. GeoCam plugin streams S24 camera as RTSP to MediaMTX (`rtsp://<mediamtx-ip>:8554/geocam`)
2. Stream registered in TAK VCM via `POST /Marti/api/video`
3. Legacy subscribes to RTSP stream → frames → vision model
4. AR Video Overlay is screen-only — does NOT expose an external API; use the raw RTSP stream

### AR Video Overlay
- Overlays COP data (markers, tracks) onto live camera feed
- Screen-only output — no external frame API
- For AI vision: access underlying camera via GeoCam/RTSP, not AR overlay output

---

## Unused ATAK Capabilities (to activate)

| Feature | How to Activate |
|---------|----------------|
| **ADS-B Direct** | Tools → ADS-B Direct → Enable; configure API source |
| **GeoFencing** | Draw shape → right-click → Set as GeoFence |
| **Video feeds (RTSP)** | S24 camera → RTSP stream → ATAK Video Manager |
| **TAK Missions** | Server → Missions → subscribe and push data packages |
| **TAK Chat → Legacy** | Wire `takchat` channel to NATS `agents.prompt.hermes.zeus.legacy` |
| **Persistent Tracks** | Settings → Device → GPS → Save Track |
| **9-line / SALUTE reports** | Reports plugin → field forms |

---

## OneMind APK Roadmap (July 2026)

Full roadmap: `/opt/data/OneMind_Android_APK_Roadmap.md` and `~/Desktop/OneMind_Android_APK_Roadmap.md` on Mac.

**Phase 1:** OneMind Agent APK — NATS WebSocket client, device telemetry to `fabric.telemetry.zeus-s24.*`, subscribes to `agents.cmd.zeus-s24.>` for Legacy commands. Stack: Kotlin + `io.nats:jnats` + ForegroundService.

**Phase 2:** Legacy TAC APK (PTT→Whisper→NATS→Legacy→TTS), ADS-B bridge, MTA GTFS-RT bridge, NYC DOT camera layer.

**Phase 3:** ATAK Plugin (Legacy Overlay) — native plugin SDK, "Ask Legacy" button, AI annotations, CoT injection from cloud.

**Phase 4:** HomeBase APK (HA via NATS), Family COP, Intelligence Collector.

### NATS transport
- Primary: `wss://<mac-tailnet-host>:8080` (Tailscale WebSocket)
- All telemetry → ClickHouse `onemind` database, 10yr TTL
- Device fabric identity: `zeus-s24` / serial `<SERIAL>`

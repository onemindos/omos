---
title: "TAK Integration Reference"
type: reference
status: historical
created: 2026-06-12
updated: 2026-10-27
owner: zeus
category: infrastructure
domain: SO/03
tags: [tak, atak, wintak, itak, freetakserver, nats, cot, plugin, sdk, integration, reference]
related:
  - "[[TACMAP-COP-STANDARD]]"
  - "[[NATS-DOCTRINE]]"
  - "[[NATS-SETUP]]"
  - "[[CLICKHOUSE-TIMEMACHINE]]"
  - "[[ARCHITECTURE]]"
ai_summary: "Complete TAK ecosystem reference — every product (ATAK/WinTAK/iTAK/WebTAK/TAKServer), access tiers, official and third-party plugins, FreeTAKServer ZMQ architecture and REST API, ATAK/WinTAK plugin SDK structure, NATS integration patterns aligned to the real live OneMind NATS setup (onemind-mac, accounts CORE/FABRIC/INTERFACE, streams EVT/TEL/CMD, subjects per NATS-DOCTRINE §9), NATS-CoT bridge code, and deployment options. Companion to TACMAP-COP-STANDARD (doctrine); this doc is the implementation playbook."
---

# TAK Integration Reference

> **HISTORICAL (marked 2026-10-27). Do not use this as a description of the live system.**
> It was written in June 2026, for a Mac-hosted NATS (`onemind-mac`, accounts CORE/FABRIC/INTERFACE)
> and a planned OpenTAKServer/FreeTAKServer deployment. Neither matches what runs today:
>
> - **Live:** official TAK Server **5.8** on OneMind GHCR images in namespace `takserver` on DOKS,
>   reachable only over the tailnet (`tak.tail717ba0.ts.net`: 8089, 8443, 8446, 64738).
>   See [`platform/takserver/README.md`](../../platform/takserver/README.md).
> - **Bus:** the in-cluster 3-node NATS (`nats.nats`) is the authority. The NatsBridgePlugin
>   reaches it through `tak-nats-leaf`.
> - **Subjects:** the live bridge uses `tak.cot.out` (TAK → bus, stream `COT`) and `tak.cot.inject`
>   (bus → TAK), not `evt.tak.cot.{uid}` as in §5–§6 below. That gap is an open decision:
>   [NATS-DOCTRINE §9.8](../NATS-DOCTRINE.md).
>
> The product catalog, CoT protocol notes and SDK sections remain useful background.

> **Companion to [[TACMAP-COP-STANDARD]]** (doctrine — what/why) and **[[NATS-DOCTRINE]] §9**
> (geo-envelope standard and locked subject taxonomy).
> This doc is the implementation playbook: products, access, plugins, SDK scaffolds,
> bridge code, and deployment — all aligned to the live OneMind NATS setup.

---

## 0. How This Fits the Build Plan

TAK integration maps to **Phase 6 — OM Connect** in [[MASTER-PLAN]] (Fabric Bridges).
The NATS↔CoT bridge is a Wombat-era bridge service, same shelf as HA and ESP32 bridges.

| Phase | TAK milestone |
|---|---|
| **T0** | This doc + [[TACMAP-COP-STANDARD]] ratified. Subjects locked in [[NATS-DOCTRINE]] §9. |
| **T1** | NATS-CoT bridge live → agents/nodes publishing `tel.fleet.>` appear on ATAK map. |
| **T2** | TAK Server (OpenTAKServer) on omos-cloud via Docker, Tailscale-bound. |
| **T3** | Field clients (ATAK/iTAK/WinTAK) connect; CoT from field flows back to bus via bridge. |
| **T4** | `cmd.fleet.{class}.{id}.{verb}` from map → actual asset command. |
| **T5** | Cross-ring correlation + Time Machine scrub on TACMAP. |

---

## 1. TAK Product Catalog

### 1.1 Clients

| Product | Platform | Description | Access |
|---|---|---|---|
| **ATAK-CIV** | Android | Flagship mobile client — full mapping, BFT, comms, video, MIL-STD-2525 symbology, APK plugin SDK | Free — Google Play or tak.gov account |
| **ATAK (military)** | Android | Full ATAK with additional classified capabilities | .mil CAC + tak.gov military account |
| **iTAK** | iOS | Lightweight iOS client; no plugin SDK (iOS sandbox) | Free — App Store; tak.gov for expanded features |
| **WinTAK** | Windows | Desktop command-post client, near-ATAK parity, DLL plugin SDK | tak.gov account (free, no CAC required for basic) |
| **WebTAK** | Browser | In-browser SA picture — no install, reads from TAK Server Marti API | tak.gov account; server-side deployment |
| **TAK Tracker** | Android/iOS | Position reporting only — simplified for non-tactical users | Free — App Store / Google Play |
| **CloudTAK** | Browser | Open-source Vue/TS web COP over Marti API | Self-hosted; AGPL-3.0 — github.com/dfpc-coe/CloudTAK |

### 1.2 Servers

| Product | Description | License | Notes |
|---|---|---|---|
| **TAK Server (TPC official)** | Java/Spring + PostGIS — CoT routing :8089 TLS, Marti REST :8443, WebTAK, federation, PKI | Apache-2.0 | github.com/TAK-Product-Center/Server |
| **OpenTAKServer** | Lighter Python reimplementation — PostGIS, Marti subset, edge-friendly | MIT | docs.opentakserver.io — **recommended for T2 deployment** |
| **FreeTAKServer (FTS)** | Python 3 Flask MVC, ZMQ internal bus, REST API, WebTAK, Node-RED | Apache-2.0 | github.com/FreeTAKTeam/FreeTAKServer — best for integration dev/ZMQ |

### 1.3 Supporting Tools

| Product | Description | License |
|---|---|---|
| **MediaMTX** | RTSP/WebRTC/HLS video relay — the de-facto TAK video server | MIT — github.com/bluenviron/mediamtx |
| **pytak** | Python CoT read/write library — the bridge core | Apache-2.0 — github.com/snstac/pytak |
| **takproto** | Python CoT protobuf support | MIT — github.com/snstac/takproto |
| **node-cot / node-tak** | TypeScript CoT libs | Apache-2.0 — github.com/dfpc-coe |
| **milsymbol** | MIL-STD-2525 symbology rendering in JS | MIT — github.com/spatialillusions/milsymbol |
| **adsbcot / aiscot** | ADS-B aircraft / AIS vessel → CoT via pytak | Apache-2.0 — the exact pattern the NATS bridge copies |

---

## 2. Access Requirements

| Tier | Who | How |
|---|---|---|
| **Public / Civilian** | General public, hobbyists, SAR, emergency management | Free tak.gov account (email signup); ATAK-CIV via Google Play |
| **Public Safety** | Fire, EMS, law enforcement | tak.gov "Public Safety" org account — unlocks additional plugins |
| **Government / Military (NIPR)** | .mil / .gov personnel | CAC/PIV + tak.gov account tied to DoD identity |
| **Classified** | Cleared personnel on classified networks | Separate classified distribution — not via public tak.gov |

**For OneMind:** a free tak.gov email account unlocks WinTAK + official plugins. No CAC needed.
OpenTAKServer and FreeTAKServer require zero tak.gov account.

---

## 3. Plugin Ecosystem

### 3.1 Official / TPC Plugins (require tak.gov account)

| Plugin | Function |
|---|---|
| **HAMMER** | Handheld Active Mapping and Mission Exploitation for Reconnaissance |
| **SNAP** | Situational Network Awareness Plugin |
| **FROGTAC** | Tactical graphics and MIL-STD-2525 symbology overlays |
| **SALUTE** | Structured spot report tool |
| **CasEvac** | Casualty evacuation reporting and tracking |
| **KILSWITCH** | Force protection / blue-force data filtering |
| **missionKEY** | Cryptographic key management for encrypted comms |
| **PACE** | Primary/Alternate/Contingency/Emergency comm planning |
| **DataSync** | Multi-device data package synchronization |
| **Grid Reference Tools** | MGRS, USNG, UTM calculators and conversions |

### 3.2 Third-Party / Community Plugins

| Plugin | Function | Relevance to OneMind |
|---|---|---|
| **Meshtastic ATAK Plugin** | LoRa mesh radio → CoT positions in TAK | High — off-grid comms, homestead mesh nodes |
| **FreeTAKServer Plugin** | Connects ATAK to FTS | Dev/ZMQ integration path |
| **Video Plugin** | RTSP/RTMP stream viewer in-map | Drone/camera feeds via MediaMTX |
| **Draw Plugin** | Collaborative map drawing and annotation | Mission planning |
| **CoT Logger** | Logs all CoT messages to file | Debug and audit |
| **adsbcot** | ADS-B aircraft → CoT via pytak | Outer-ring air traffic feed |
| **aiscot** | AIS vessel → CoT via pytak | Outer-ring maritime feed |
| **DJI/MAVLink → CoT** | Drone telemetry bridge | Drone inner-ring `tel.fleet.drone.` feed |
| **OpenTAK ICU** | Phone camera → live video icon on map | Field camera streams |

Plugins distribute as `.apk` files sideloaded on Android or as TAK Server data packages.

---

## 4. CoT Protocol

CoT (Cursor-on-Target) is the protocol all TAK products speak — a tiny XML event expressing **what / when / where**. Unchanged base class since 2002.

### 4.1 Message Structure

```xml
<event version="2.0"
       uid="onemind-agent-legacy-001"
       type="a-f-G-U-C-I"
       time="2026-06-12T00:00:00Z"
       start="2026-06-12T00:00:00Z"
       stale="2026-06-12T00:05:00Z"
       how="m-g">
  <point lat="32.715736"
         lon="-117.161087"
         hae="0"
         ce="10"
         le="10"/>
  <detail>
    <contact callsign="LEGACY-AGENT"/>
    <remarks>OneMind AI agent — published by NATS bridge</remarks>
    <track course="0" speed="0"/>
  </detail>
</event>
```

Transport: TCP or UDP to TAK Server port 8089 (TLS) or 8087 (plain).
LAN discovery: UDP multicast `239.2.3.1:6969`.

### 4.2 CoT Type Code Structure

`{affiliation}-{battle-dimension}-{function}...` — maps directly to MIL-STD-2525 symbology.

| Code | Description |
|---|---|
| `a-f-G-U-C` | Friendly Ground Unit Combat |
| `a-f-G-U-C-I` | Friendly Ground Unit Combat Infantry → **AI agent** |
| `a-f-G-E-V` | Friendly Ground Equipment Vehicle → **robot / vehicle** |
| `a-f-A-M-F-Q` | Friendly Air Unmanned Fixed-wing → **drone** |
| `a-f-G-I-U-T-C` | Friendly Ground Installation → **node / fixed structure** |
| `a-h-G-U-C` | Hostile Ground Unit Combat |
| `a-u-G-U-C` | Unknown Ground Unit |
| `b-a-g-e-t` | Event — GeoFence breach |

The type vocabulary in [[NATS-DOCTRINE]] §9.3 maps 1:1 to CoT types — the bridge is a pure field rename.

### 4.3 Geo-Envelope → CoT Field Map

Per [[NATS-DOCTRINE]] §9.1–9.2, the OneMind geo-envelope maps directly to CoT fields:

```
lat          → CoT point/@lat
lon          → CoT point/@lon
hae          → CoT point/@hae
course       → CoT detail/track/@course
speed        → CoT detail/track/@speed
ce, le       → CoT point/@ce, @le
type         → CoT event/@type  (use §9.3 type taxonomy)
uid          → CoT event/@uid   (format: onemind-{class}-{id})
callsign     → CoT detail/contact/@callsign
t            → CoT event/@time, @start (epoch ms → ISO-8601 UTC)
stale        → CoT event/@stale = t + stale_s
```

---

## 5. NATS Integration — Live System Context

> **[[NATS-SETUP]] is authoritative for the live config. This section shows how TAK maps onto it.**

### 5.1 Live NATS Server (onemind-mac)

```
Binary:   /opt/homebrew/bin/nats-server  v2.14.1
Config:   ~/.nats/config/onemind.conf
Data:     ~/.nats/data
Ports:    4222 (client)  8222 (monitor)  8080 (WebSocket)  1883 (MQTT)  7422 (leaf)
Domain:   onemind
Service:  LaunchAgent  io.nats.server
Creds:    ~/.nats/config/accounts.creds  (chmod 600, never commit)
```

### 5.2 Account Layout

| Account | Purpose | JetStream | User | TAK bridge uses |
|---|---|---|---|---|
| `CORE` | Agents + meta services, bridge | yes | `core` | Bridge publishes `evt.tak.cot.>`, subscribes `tel.fleet.>` |
| `FABRIC` | Telemetry + IoT (sensors, robots) | yes | `fabric` | Fleet assets publish `tel.fleet.>` here |
| `INTERFACE` | Apps — read-only imports | yes | `interface` | ATAK/WinTAK plugins connect with this user |
| `GLOBAL` | **DEPRECATED** — migrate clients | yes | `onemind` | — |

Cross-account wiring (from [[NATS-SETUP]]):
```
CORE    imports: tel.> from FABRIC   ← bridge reads fleet positions here
INTERFACE imports: evt.> from CORE   ← ATAK plugin reads TAK CoT events
```

### 5.3 Streams

| Stream | Account | Subjects | Retention | TAK relevance |
|---|---|---|---|---|
| **EVT** | CORE | `evt.>` | 30 days | `evt.tak.cot.{uid}` — all field CoT lands here |
| **TEL** | FABRIC | `tel.>` | 7 days | `tel.fleet.{class}.{id}.position` — source for bridge |
| **CMD** | CORE | `cmd.>` | 24h WorkQueue | `cmd.fleet.{class}.{id}.{verb}` — map commands |

### 5.4 TAK Subject Taxonomy (from [[NATS-DOCTRINE]] §9.1 — LOCKED)

```
# Inner ring — fleet (participatory)
tel.fleet.{class}.{id}.position       class ∈ drone|robot|vehicle|node|edge|agent
tel.fleet.{class}.{id}.status         battery, mode, link health
cmd.fleet.{class}.{id}.{verb}         estop | goto | rtl | arm | ...
cmd.mission.{id}.{verb}               create | waypoint | abort | complete

# Outer ring — world intel (observational)
evt.world.air.{id}                    aircraft (ADS-B)
evt.world.air.mil.{id}                military aircraft
evt.world.sea.{id}                    vessels (AIS)
evt.world.conflict.{id}               conflict events (ACLED)
evt.world.fire.{id}                   wildfire (NASA FIRMS)
evt.world.seismic.{id}                earthquakes (USGS)
evt.world.chokepoint.{id}             strategic waterway status

# TAK bridge — field devices (ATAK/iTAK/WinTAK)
evt.tak.cot.{uid}                     CoT events FROM field TAK clients → EVT stream
evt.tak.status.{uid}                  TAK client connection state

# Map control (agent/app → map)
cmd.map.filter.set
cmd.map.highlight.{uid}
cmd.map.viewport.set
cmd.map.geofence.{id}.{verb}
cmd.map.annotation.{id}.{verb}
evt.map.geofence.{id}.{enter|exit}    Watch Job fires
```

### 5.5 Geo-Envelope Message Shape (from [[NATS-DOCTRINE]] §9.2)

```json
{
  "uid":      "agent-legacy",
  "type":     "fleet.agent",
  "lat":      32.715736,
  "lon":      -117.161087,
  "hae":      0,
  "course":   0,
  "speed":    0,
  "callsign": "LEGACY-AGENT",
  "stale":    300,
  "t":        1718150400000,
  "meta": {
    "battery_pct": -1,
    "mode":        "active",
    "node_id":     "omos-mbp"
  }
}
```

Required fields: `uid`, `type`, `lat`, `lon`, `t`. All others default to `-1` / `""`.
`hae` = 0 for ground assets. `t` = epoch milliseconds UTC (never ISO strings in geo hot-path).

---

## 6. NATS ↔ CoT Bridge

The only code OneMind builds — ~150 lines Python using `pytak` + `nats-py`.
Lives at `services/tak-bridge/` in `onemindos-core`.

### 6.1 Architecture

```
omd nodes / agents publish:
  tel.fleet.agent.legacy.position  { geo-envelope JSON }
  tel.fleet.drone.dji1.position
        │
        ▼ [FABRIC account, TEL stream]
  [NATS onemind-mac]
        │ CORE imports tel.> from FABRIC
        ▼ [CORE account, user: core]
  [NATS-CoT Bridge]  services/tak-bridge/nats_cot_bridge.py
        │
        ├── TCP CoT XML ──▶ OpenTAKServer :8089 TLS  (Tailscale: omos-cloud)
        │
        └── CoT FROM TAK ──▶ pub: evt.tak.cot.{uid}  [CORE, EVT stream]
                                   └──▶ ClickHouse onemind.events_v2  (Time Machine)
```

### 6.2 Bridge Implementation

```python
# services/tak-bridge/nats_cot_bridge.py
import asyncio, json, os, time
from datetime import datetime, timezone, timedelta
import nats
import pytak
import xml.etree.ElementTree as ET

TAK_HOST  = "omos-cloud.tail717ba0.ts.net"
TAK_PORT  = 8089
NATS_URL  = "nats://localhost:4222"
NATS_USER = "core"

# Type lookup — NATS-DOCTRINE §9.3 → CoT type codes
COT_TYPES = {
    "fleet.agent":    "a-f-G-U-C-I",
    "fleet.robot":    "a-f-G-E-V",
    "fleet.drone":    "a-f-A-M-F-Q",
    "fleet.vehicle":  "a-f-G-E-V",
    "fleet.node":     "a-f-G-I-U-T-C",
    "fleet.edge":     "a-f-G-U-C",
    "fleet.aircraft": "a-f-A-C-F",
    "fleet.vessel":   "a-f-S-X-M",
}

def geo_to_cot(msg: dict) -> str:
    uid      = f"onemind-{msg['uid']}"
    cot_type = COT_TYPES.get(msg.get("type", ""), "a-f-G-U-C")
    callsign = msg.get("callsign", uid)
    t_ms     = msg.get("t", int(time.time() * 1000))
    stale_s  = msg.get("stale", 300)
    now      = datetime.fromtimestamp(t_ms / 1000, tz=timezone.utc)
    stale    = now + timedelta(seconds=stale_s)
    fmt      = "%Y-%m-%dT%H:%M:%SZ"
    return f"""<event version="2.0" uid="{uid}" type="{cot_type}"
         time="{now.strftime(fmt)}" start="{now.strftime(fmt)}"
         stale="{stale.strftime(fmt)}" how="m-g">
  <point lat="{msg['lat']}" lon="{msg['lon']}"
         hae="{msg.get('hae', 0)}" ce="10" le="10"/>
  <detail>
    <contact callsign="{callsign}"/>
    <track course="{msg.get('course', 0)}" speed="{msg.get('speed', 0)}"/>
    <remarks>OneMind {msg.get('type', '')}</remarks>
  </detail>
</event>"""

def make_envelope(uid: str, raw_xml: str) -> bytes:
    return json.dumps({
        "v": 1, "id": f"tak-{uid}",
        "ts": datetime.now(timezone.utc).isoformat(),
        "type": "tak.atak", "source": "tak-bridge",
        "subject": f"evt.tak.cot.{uid}",
        "data": raw_xml
    }).encode()

async def main(password: str):
    nc  = await nats.connect(NATS_URL, user=NATS_USER, password=password)
    tak = await pytak.TCPClient.create(TAK_HOST, TAK_PORT)

    async def on_position(msg):
        try:
            data = json.loads(msg.data.decode())
            await tak.send(geo_to_cot(data).encode())
        except Exception as e:
            print(f"bridge→tak error: {e}")

    await nc.subscribe("tel.fleet.>", cb=on_position)

    async def from_tak():
        async for raw in tak.receive():
            try:
                text = raw.decode() if isinstance(raw, bytes) else raw
                uid  = ET.fromstring(text).get("uid", "unknown")
                await nc.publish(f"evt.tak.cot.{uid}", make_envelope(uid, text))
            except Exception as e:
                print(f"tak→bridge error: {e}")

    asyncio.create_task(from_tak())
    await asyncio.Future()

if __name__ == "__main__":
    # NATS_CORE_PASSWORD set in LaunchAgent env — never hardcode
    asyncio.run(main(os.environ["NATS_CORE_PASSWORD"]))
```

### 6.3 FreeTAKServer ZMQ → NATS (dev path)

FTS exposes raw ZMQ ports — useful for local integration dev without deploying a TAK Server.

| FTS Port | Direction | What |
|---|---|---|
| **15557** | OUT (FTS publishes) | All incoming CoT — subscribe to consume |
| **15558** | IN (FTS subscribes) | Inject CoT into FTS — publish here to push markers to field clients |

```python
# services/tak-bridge/fts_zmq_bridge.py
import asyncio, zmq, zmq.asyncio, nats, json, os
from datetime import datetime, timezone
import xml.etree.ElementTree as ET

async def main(password: str):
    ctx      = zmq.asyncio.Context()
    sock     = ctx.socket(zmq.SUB)
    sock.connect("tcp://localhost:15557")
    sock.subscribe(b"")
    nc = await nats.connect("nats://localhost:4222", user="core", password=password)

    while True:
        raw = await sock.recv()
        try:
            uid = ET.fromstring(raw.decode()).get("uid", "unknown")
            envelope = json.dumps({
                "v": 1, "id": f"tak-{uid}",
                "ts": datetime.now(timezone.utc).isoformat(),
                "type": "tak.atak", "source": "fts-zmq-bridge",
                "subject": f"evt.tak.cot.{uid}",
                "data": raw.decode()
            }).encode()
            await nc.publish(f"evt.tak.cot.{uid}", envelope)
        except Exception as e:
            print(f"zmq-bridge error: {e}")

if __name__ == "__main__":
    asyncio.run(main(os.environ["NATS_CORE_PASSWORD"]))
```

### 6.4 ClickHouse — TAK Events in the Time Machine

TAK events on `evt.tak.cot.>` auto-land in ClickHouse via the live `nats_evt` → `evt_to_events_v2`
pipeline (see [[CLICKHOUSE-TIMEMACHINE]]).

ClickHouse: local, unauthed, **HTTP :8124**, **TCP :9001**.

```sql
-- All CoT received from field clients (last 24h)
SELECT ts, JSONExtractString(data, 'uid') AS uid,
       JSONExtractFloat(data, 'lat') AS lat,
       JSONExtractFloat(data, 'lon') AS lon
FROM onemind.events_v2
WHERE subject LIKE 'evt.tak.cot.%'
  AND ts > now() - INTERVAL 24 HOUR
ORDER BY ts DESC;

-- Fleet position replay for Time Machine scrub
SELECT ts, subject, data FROM onemind.events_v2
WHERE subject LIKE 'tel.fleet.%.position'
  AND ts BETWEEN {t1} AND {t2}
ORDER BY ts ASC;

-- Full history union (legacy + new, per CLICKHOUSE-TIMEMACHINE §2.1)
SELECT timestamp AS ts, subject, payload AS data FROM onemind.events
UNION ALL
SELECT ts, subject, data FROM onemind.events_v2
ORDER BY ts DESC;
```

---

## 7. ATAK Plugin SDK (Android)

### 7.1 Plugin Structure

```
onemind-atak-plugin/
├── app/src/main/
│   ├── java/com/onemind/takplugin/
│   │   ├── PluginLifecycle.java      ← entry point (AbstractPlugin)
│   │   ├── DropDownReceiver.java     ← slide-in UI panel
│   │   ├── CotEventReceiver.java     ← hooks all incoming CoT
│   │   └── NatsBridgeService.java    ← Android Service, NATS WS client
│   ├── res/layout/plugin_panel.xml
│   └── assets/plugin-description.xml
└── app/build.gradle
```

### 7.2 Plugin Manifest

```xml
<!-- assets/plugin-description.xml -->
<plugin name="OneMind TAK Bridge"
        classname="com.onemind.takplugin.PluginLifecycle"
        description="Live NATS bus integration for OneMind OS"
        version="1.0"
        requires="com.atakmap.app@4.8.0.0"/>
```

### 7.3 Lifecycle Entry Point

```java
public class PluginLifecycle extends AbstractPlugin {

    @Override
    public void onCreate(final Context ctx, final View view) {
        AtakBroadcast.getInstance().registerReceiver(cotReceiver,
            new AtakBroadcast.DocumentedIntentFilter(
                CotMapComponent.INTENT_NEW_COT));
        // NatsBridgeService connects to ws://onemind-mac.tail717ba0.ts.net:8080
        // user=interface, password from Android secure storage
        ctx.startService(new Intent(ctx, NatsBridgeService.class));
    }

    @Override
    public void onDestroy(final Context ctx) {
        AtakBroadcast.getInstance().unregisterReceiver(cotReceiver);
        ctx.stopService(new Intent(ctx, NatsBridgeService.class));
    }

    // All field CoT → NATS bus as evt.tak.cot.{uid} (lands in EVT stream)
    private final BroadcastReceiver cotReceiver = new BroadcastReceiver() {
        @Override
        public void onReceive(Context ctx, Intent intent) {
            CotEvent event = intent.getParcelableExtra("cotEvent");
            if (event != null) {
                NatsBridgeService svc = NatsBridgeService.getInstance();
                if (svc != null)
                    svc.publish("evt.tak.cot." + event.getUID(),
                                event.toString().getBytes());
            }
        }
    };
}
```

The plugin uses the **INTERFACE** account (WebSocket `:8080`) — read-only imports, cannot
publish to `tel.>` or `cmd.>`. Only publishes to `evt.tak.cot.>` which CORE receives.

### 7.4 Key ATAK Java APIs

| API | Purpose |
|---|---|
| `MapView.getRootGroup()` | Access the map; add/remove map items |
| `CotDispatcher.getInstance().dispatchTo(event)` | Inject CoT into ATAK's pipeline |
| `AtakBroadcast.getInstance()` | Android intent bus within ATAK |
| `DropDownReceiver` | Slide-in side panels |
| `MapGroup` | Organize markers into named groups |
| `Marker` | Individual map point — UID, callsign, icon, color |
| `CotEvent` | CoT event object — UID, type, lat/lon, detail XML |

### 7.5 Build Requirements

- `atak-core.jar` — from tak.gov (free account required)
- `atak-javadoc.jar`
- Android SDK API 26+, Gradle 7.x
- Sign with own keystore — ATAK verifies signatures at sideload
- Sideload: Settings → Apps → Install Unknown Apps → Files

---

## 8. WinTAK Plugin SDK (Windows)

```csharp
// OneMindPlugin.cs — C# .NET DLL, implements IWinTAKPlugin
public class OneMindPlugin : IWinTAKPlugin {
    private IConnection _nats;

    public void Initialize(IPluginHost host) {
        var opts = ConnectionFactory.GetDefaultOptions();
        opts.Url      = "nats://onemind-mac.tail717ba0.ts.net:4222";
        opts.User     = "interface";
        opts.Password = SecureStorage.Get("nats_interface_password");
        _nats = new ConnectionFactory().CreateConnection(opts);

        // All incoming CoT → bus as evt.tak.cot.{uid}
        host.CotEventReceived += (_, e) =>
            _nats.Publish($"evt.tak.cot.{e.Uid}",
                          Encoding.UTF8.GetBytes(e.ToXml()));

        // cmd.map.highlight.> → WinTAK map action
        _nats.SubscribeAsync("cmd.map.highlight.>", (_, args) => {
            var uid = args.Message.Subject.Split('.').Last();
            host.HighlightMapItem(uid);
        });

        host.RegisterPanel("OneMind", new OneMindPanel());
    }
}
```

**NuGet:** `NATS.Net` v2.x (preferred) or `NATS.Client` v1.x

---

## 9. FreeTAKServer — Internal Architecture

FTS 2.x uses **ZeroMQ** as its internal event bus, not NATS.

```
[ATAK Client] ─CoT TCP─▶ [FTS Core Controller]
                                  │
                            ZMQ PUB/SUB (internal)
                                  │
             ┌────────────────────┼──────────────────────┐
        [CoT Service]       [REST API :19023]    [ZMQ Bridge]
                                                       │
                                           External ZMQ :15557 (PUB)
                                           External ZMQ :15558 (SUB)
```

### 9.1 FTS REST API

```
Base: http://<host>:19023/

POST /ManageGeoObject/postGeoObject         — inject map geo-object
GET  /ManageGeoObject/getGeoObject/{uid}    — retrieve object
DELETE /ManageGeoObject/deleteGeoObject     — remove marker
POST /ManageEmergency/postEmergency         — create emergency marker
POST /ManageChat/postChatToAll              — broadcast chat to all clients
POST /ManageMission/postMission             — mission management
POST /DataPackageTable/postDP               — upload data package
```

### 9.2 FTS vs OpenTAKServer vs Official TAK Server

| Aspect | FreeTAKServer | OpenTAKServer | Official TPC TAK Server |
|---|---|---|---|
| Language | Python 3 / Flask | Python 3 | Java / Spring |
| License | Apache-2.0 | MIT | Apache-2.0 |
| RAM | ~300 MB | ~500 MB | ~2–4 GB |
| ZMQ bridge | ✅ native | ❌ | ❌ |
| REST API | ✅ :19023 | Marti-compatible | Full Marti :8443 |
| Federation | Limited | Partial | Full |
| Best for | Integration dev, ZMQ path | **T2 deployment on omos-cloud** | Production/federation |

---

## 10. Full Integration Diagram

```
┌────────────────────────────────────────────────────────────────────────┐
│                   MAC — onemind-mac  (local)                            │
│                                                                          │
│   OpenClaw Agents (Legacy, Zenith…)     omd daemon                     │
│   Hermes Gateway (:8642/:8787/:9119)                                    │
│         │                                                                │
│         │  publish: tel.fleet.{class}.{id}.position                    │
│         │           { geo-envelope JSON }  → FABRIC account             │
│         ▼                                                                │
│   ┌─────────────────────────────────────┐                               │
│   │  NATS  onemind-mac  v2.14.1         │                               │
│   │  :4222  :8080 ws  :1883 mqtt        │                               │
│   │  :7422 leaf  :8222 monitor          │                               │
│   │                                     │                               │
│   │  FABRIC  TEL stream  tel.>  7d      │ ← fleet positions             │
│   │  CORE    EVT stream  evt.>  30d     │ ← TAK CoT, world events       │
│   │  CORE    CMD stream  cmd.>  24h wq  │ ← map/fleet commands          │
│   │  INTERFACE  (app read view)         │                               │
│   └──────────────┬──────────────────────┘                               │
│                  │                                                        │
│     ┌────────────┴───────────────┐                                       │
│     ▼                            ▼                                       │
│  [NATS-CoT Bridge]          [FTS ZMQ Bridge]  ← dev path only           │
│  services/tak-bridge/        services/tak-bridge/                        │
│  account: CORE               account: CORE                               │
└──────────────│───────────────────────────────────────────────────────────┘
               │ Tailscale  tail717ba0.ts.net
               ▼
   ┌────────────────────────────────────┐
   │  omos-cloud (VPS)                  │
   │  OpenTAKServer (Docker)            │
   │    :8089 TLS CoT                   │
   │    :8443 Marti REST                │
   │    :8080 WebTAK                    │
   │                                    │
   │  MediaMTX (optional)               │
   │    :8554 RTSP  :8888 WebRTC        │
   └──────────────┬─────────────────────┘
                  │ CoT — all clients on Tailscale
   ┌──────────────┼──────────────┐
   ▼              ▼              ▼
 [ATAK]         [iTAK]       [WinTAK]
 Android          iOS          Windows
   │
   └── OneMind plugin (optional)
       WS → onemind-mac:8080  user=interface
       publishes:  evt.tak.cot.{uid}    → EVT stream → ClickHouse
       subscribes: cmd.map.highlight.>  → map highlight action
```

---

## 11. Deployment Options

| Option | Stack | RAM | Notes | Verdict |
|---|---|---|---|---|
| **OpenTAKServer + Docker (omos-cloud)** | Python, MIT | ~500 MB | Tailscale-bound, no tak.gov needed, Marti API | **Recommended T2** |
| **Official TAK Server + Docker** | Java/Spring, Apache-2.0 | ~2–4 GB | Full Marti, federation, WebTAK | Production / federation |
| **FreeTAKServer (dev)** | Python, Apache-2.0 | ~300 MB | ZMQ bridge native, best for local testing | Dev / ZMQ path |
| **Commercial hosted** | Broadleaf / ParaZero etc. | — | Per-seat, cloud-coupled | Avoid — not sovereign |

**Sovereignty rule (from [[TACMAP-COP-STANDARD]]):** everything Tailscale-bound.
TAK Server public IP stays dark. No COP data on the public internet.

---

## 12. Quick-Start Checklist

```
□ Create free tak.gov account (unlocks WinTAK + official plugins — optional for OpenTAKServer)
□ Install ATAK-CIV on Android (Google Play, free)
□ Verify NATS live: nats account info --user core --password <see ~/.nats/config/accounts.creds>
□ Deploy OpenTAKServer on omos-cloud via Docker (Tailscale-bound, :8089 TLS)
□ pip install pytak nats-py
□ Write services/tak-bridge/nats_cot_bridge.py (§6.2 above)
□ Add NATS_CORE_PASSWORD to bridge LaunchAgent environment
□ Configure ATAK: Server → tak-server.tail717ba0.ts.net :8089 TLS (phone on Tailscale)
□ Test: publish tel.fleet.agent.legacy.position → confirm marker appears in ATAK
□ Verify ClickHouse: clickhouse client --port 9001 \
    --query "SELECT count() FROM onemind.events_v2 WHERE subject LIKE 'evt.tak.%'"
□ (Optional) Build + sideload ATAK OneMind plugin → field CoT flows to evt.tak.cot.> on bus
```

---

## 13. Source Index

| Reference | URL |
|---|---|
| TAK Server (official, Apache-2.0) | https://github.com/TAK-Product-Center/Server |
| ATAK-CIV source | https://github.com/TAK-Product-Center/atak-civ |
| OpenTAKServer | https://docs.opentakserver.io |
| FreeTAKServer | https://github.com/FreeTAKTeam/FreeTAKServer |
| pytak | https://github.com/snstac/pytak |
| takproto | https://github.com/snstac/takproto |
| adsbcot / aiscot | https://github.com/snstac/adsbcot · https://github.com/snstac/aiscot |
| CloudTAK (web COP) | https://github.com/dfpc-coe/CloudTAK |
| node-cot / node-tak | https://github.com/dfpc-coe/node-CoT · https://github.com/dfpc-coe/node-tak |
| MediaMTX | https://github.com/bluenviron/mediamtx |
| Meshtastic ATAK plugin | https://github.com/meshtastic/ATAK-Plugin |
| milsymbol | https://github.com/spatialillusions/milsymbol |
| MapLibre GL JS | https://maplibre.org |
| MITRE CoT Developer's Guide | https://nps.edu/documents/104517539/109705106/COT+Developer+Guide.pdf |
| CoT type codes / MIL-STD-2525 | https://freetakteam.github.io/FreeTAKServer-User-Docs/About/architecture/mil_std_2525 |
| NATS docs | https://docs.nats.io |
| nats.py client | https://github.com/nats-io/nats.py |
| NATS.Net client | https://github.com/nats-io/nats.net.v2 |

---

*Maintained in `onemindos-core/docs/infra/` alongside [[NATS-SETUP]] and [[CLICKHOUSE-TIMEMACHINE]].*
*Update when: TAK Server version changes, new plugins adopted, NATS account layout changes,
or subject taxonomy updated in [[NATS-DOCTRINE]] §9.*

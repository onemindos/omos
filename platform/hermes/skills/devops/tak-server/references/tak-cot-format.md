# Cursor-on-Target (CoT) XML Reference

## Minimal valid event

```xml
<event version="2.0" uid="unique-id" type="a-f-G-U-C"
  time="2026-07-13T06:35:00Z" start="2026-07-13T06:35:00Z"
  stale="2026-07-13T06:40:00Z" how="m-g">
  <point lat="41.0" lon="-73.6" hae="100" ce="10" le="10"/>
  <detail>
    <contact callsign="My Callsign"/>
    <remarks>Optional human-readable note</remarks>
  </detail>
</event>
```

## Required attributes

| Attribute | Description | Example |
|-----------|-------------|---------|
| `version` | Always `"2.0"` | `"2.0"` |
| `uid` | Globally unique ID | `"legacy-test-001"` |
| `type` | CoT type (see below) | `"a-f-G-U-C"` |
| `time` | Event time (ISO 8601) | `"2026-07-13T06:35:00Z"` |
| `start` | Same as time for point events | Same as `time` |
| `stale` | Expiry time (~5 min after time) | `"2026-07-13T06:40:00Z"` |
| `how` | How position was obtained | `"m-g"` (machine-generated) |

## Common CoT types

| Type | Meaning |
|------|---------|
| `a-f-G-U-C` | Friendly ground unit — unknown (default for manual injection) |
| `a-f-G-U` | Friendly ground unit |
| `a-f-A-C` | Friendly aircraft (generic) |
| `a-f-A-M-F-Q` | **Friendly UAS/drone (quadrotor)** — renders drone icon on TAK map ✓ |
| `a-f-A-M-F-F` | Friendly fixed-wing UAS |
| `a-f-A-M-H` | Friendly helicopter / rotary-wing aircraft |
| `a-f-A-M-F` | Friendly military fixed-wing |
| `a-f-S-C` | Friendly sensor |
| `a-h-G` | Hostile ground unit |
| `a-n-G` | Neutral ground unit |
| `b-a-o-tbl` | Emergency beacon (911/PLB) |
| `b-r-f-h` | Fire hazard |
| `b-r-f` | Fire |

Type breakdown: `a` = atom (point), `f/h/n/u` = friend/hostile/neutral/unknown, `G/A/S` = ground/air/sensor, `M` = military, `F` = fixed-wing, `Q` = quadrotor, `H` = helicopter

### Drone/UAS CoT example (`a-f-A-M-F-Q`)

```xml
<event version="2.0" uid="DRONE-CALLSIGN-001" type="a-f-A-M-F-Q"
  time="2026-07-14T08:00:00.0Z" start="2026-07-14T08:00:00.0Z"
  stale="2026-07-14T08:00:10.0Z" how="m-g">
  <point lat="41.0219000" lon="-73.6404000" hae="100.0" ce="5.0" le="5.0"/>
  <detail>
    <track speed="29.2" course="45.0"/>   <!-- speed in knots, course in degrees -->
    <uid Droid="DRONE-CALLSIGN"/>
    <contact callsign="DRONE-CALLSIGN"/>
    <remarks>Mode:patrol Alt:100m Batt:99% Msgs:1234</remarks>
    <status battery="99"/>
    <precisionlocation geopointsrc="gps" altsrc="gps"/>
  </detail>
</event>
```

**Key fields for aircraft SA:**
- `<track speed="..." course="..."/>` — speed in **knots** (not m/s), course in degrees true north. Enables the directional arrow on TAK map.
- `<uid Droid="..."/>` — sets the display name (Android ATAK convention, required for label)
- `stale` — keep 10s for live 1Hz feeds so the icon disappears quickly if the feed drops
- `hae` = height above ellipsoid in meters (same as AGL for low-altitude ops)

## Point element

| Attribute | Description |
|-----------|-------------|
| `lat` | Latitude (-90 to 90) |
| `lon` | Longitude (-180 to 180) |
| `hae` | Height above ellipsoid (meters) |
| `ce` | Circular error (meters, 9999999 = unknown) |
| `le` | Linear error (meters, 9999999 = unknown) |

## Detail elements (all optional)

| Element | Attributes/Children |
|---------|-------------------|
| `<contact callsign="..."/>` | Display name on map |
| `<remarks>...</remarks>` | Text shown in tooltip |
| `<status readiness="true"/>` | Readiness indicator |
| `<track course="..." speed="..."/>` | Course (deg), speed (m/s) |
| `<precisionlocation geopointsrc="..." altsrc="..."/>` | GPS quality |
| `<link uid="..." type="..." relation="..." />` | Link to another event |
| `<_flow-tags_ TAK-Server="..." time="..."/>` | Source tag (added by server) |

## Publishing from Hermes

```bash
# Direct NATS injection (preferred)
nats pub tak.cot.inject '<event ...>...</event>'

# Python helper
python3 -c "
import subprocess, json, datetime
now = datetime.datetime.utcnow()
stale = now + datetime.timedelta(minutes=5)
event = f'''<event version=\"2.0\" uid=\"legacy-{now.timestamp()}\" type=\"a-f-G-U-C\"
  time=\"{now.isoformat()}Z\" start=\"{now.isoformat()}Z\"
  stale=\"{stale.isoformat()}Z\" how=\"m-g\">
  <point lat=\"41.0\" lon=\"-73.6\" hae=\"100\" ce=\"10\" le=\"10\"/>
  <detail><contact callsign=\"Legacy\"/><remarks>Injected by Hermes</remarks></detail>
</event>'''
subprocess.run(['nats','pub','tak.cot.inject', event])
"
```

## Pitfalls

- **UID must be globally unique** — collisions cause TAK clients to silently drop duplicates. Use timestamp + prefix.
- **stale time** — TAK clients remove events after stale. Keep it 5-10 minutes for live data, longer for static markers.
- **XML must be valid** — no unescaped `&`, `<`, or `>` inside element text. Use CDATA sections if needed: `<remarks><![CDATA[text with & special chars]]></remarks>`
- **lat/lon precision** — TAK displays to ~6 decimal places. More precision is fine but unnecessary.
- The `tak.cot.inject` subject routes through the NatsBridgePlugin inside the TAK plugins pod. It connects to `tak-nats-leaf` (ClusterIP), which remotes to the DOKS hub. No auth needed — the leaf handles creds.

# iTAK vs ATAK Feature Gap — Full Comparison

Session confirmed 2026-07-14 during Zeus's iPhone 17 Pro Max → Samsung S26 Ultra ATAK migration.

## Origin

ATAK was built by the Air Force Research Laboratory (AFRL) for Android in 2010. iOS got a port — iTAK — years later. The gap is architectural, not a version lag.

| Factor | Android | iOS |
|--------|---------|-----|
| Side-loading apps | ✅ Native, unrestricted | ⚠️ AltStore / Enterprise only |
| Plugin architecture | ✅ Full APK plugin system | ❌ No equivalent |
| Background processes | ✅ Unrestricted | ⚠️ iOS kills background apps |
| Hardware PTT integration | ✅ Full Bluetooth/USB | ❌ Limited |
| TAK.gov distribution | ✅ Direct APK download | ⚠️ App Store dependency |
| Government MDM | ✅ Knox / Android Enterprise | ⚠️ Limited enrollment |
| Sensor/device integration | ✅ Full USB OTG, Bluetooth serial | ❌ Locked |
| Platform age / maturity | ✅ Since 2010 — 16+ years | ⚠️ Port added years later |

## Feature-by-Feature

| Feature | iTAK (iPhone) | ATAK (Android) |
|---------|--------------|----------------|
| Plugins | Almost none | Hundreds — full ecosystem |
| Mission management | Receive only | Create, manage, log |
| MEDEVAC / 9-Line | ❌ No | ✅ Built-in |
| SALUTE reports | ❌ No | ✅ Built-in |
| PACE planning | ❌ No | ✅ Full |
| Video in-map (RTSP/WebRTC) | ❌ No | ✅ Live PiP on COP |
| PTT hardware (Wave/Zello) | ⚠️ Limited | ✅ Full integration |
| ExCheck creation | ❌ View only | ✅ Create + manage |
| Data packages | ⚠️ Receive only | ✅ Build + push to team |
| Offline map cache | ⚠️ Limited | ✅ Full tile caching |
| MIL-STD-2525 symbols | ⚠️ Partial | ✅ Complete library |
| Legacy AI C2 | ⚠️ Text chat only | ✅ Map markers = commands |
| Drone integration | ❌ No | ✅ Apollo plugin + more |
| Update cycle | ⚠️ Lags ATAK | ✅ Always current |

## OneMind / Legacy AI C2 Impact

On iTAK, Legacy AI reads GeoChat messages and replies in text. On ATAK, the map IS the command interface.

| Zeus Does in ATAK | Legacy Sees on NATS | Legacy Does Back |
|-------------------|---------------------|------------------|
| Drop yellow marker + remarks | CoT `tak.cot.out` | Researches intel, pushes waypoint back to map |
| GeoChat to Legacy room | `b-t-f` CoT message | Responds in chat AND drops map markers simultaneously |
| Blue "INTEL" marker | Mission UID add event | Logs intel point, cross-references, acknowledges |
| Green "OBJ COMPLETE" | Mission change event | Logs completion, updates mission log on TAK server |
| Red "HOLD" marker | CoT event | Logs hold, maintains position watch, sends ACK |
| Emergency button | `a-f-G-U-C-E` CoT | Priority response, pushes support info as map objects |
| Change team color Red→Green | SA update | Status signal — Legacy reads as mission state change |

The iTAK text interface becomes a shared Common Operating Picture with full bidirectional C2 on ATAK.

# ATAK-CIV 5.7.0 Loadout Reference

## Source Location (Mac Desktop)
`/Users/zeus/Desktop/atak /`

## Core APKs

| File | Size | SHA-Prefix | Description |
|------|------|------------|-------------|
| `ATAK-5.7.0.10-e0d26f0e-civ-release.apk` | 351 MB | e0d26f0e | Full CIV with bundled US offline maps |
| `ATAK-5.7.0.10-e0d26f0e-civSmall-release.apk` | 107 MB | e0d26f0e | Small CIV, no bundled maps |
| `TakTracker-8.1.0-5d05c2f-release (1).apk` | 66 MB | 5d05c2f | Lightweight blue force tracking |

## Plugins (30 total, ~2.75 GB)

### Core Ops (Tier 1)
| Plugin | Size | Function |
|--------|------|----------|
| datasync-4.0.4-a858109b | 4.1 MB | Data sync with TAK server |
| takchat-1.12-e54a9cda | 4.6 MB | GeoChat text comms |
| compassnav-1.0-4f89153f | 61 KB | Compass navigation overlay |
| reports-6.7-cff20a38 | 6.4 MB | Mission reports (CAS, medevac) |
| grgbuilder-3.14-24cb3168 | 4.7 MB | Grid Reference Graphic builder |
| hammer-1.2-311f24c3 | 954 KB | Target marking / range tools |
| geotakcam-1.1-33a9d47b | 77 KB | Geo-tagged photos into ATAK |
| imagemarkup-6.2.1-0d162e7 | 4.5 MB | Annotate photos on map |
| networkmonitor-7.8.2-e0ec589 | 4.8 MB | Debug TAK server connectivity |

### Situational Awareness (Tier 2)
| Plugin | Size | Function |
|--------|------|----------|
| adsb-direct-7.0-c2bf4cbf | 4.8 MB | Live aircraft tracking (ADS-B) |
| vns-4.0-7ad24562 | 4.9 MB | Visual Navigation System |
| wave-2.1-248180cd | 4.6 MB | Wave form overlay (comms planning) |
| fireareasurvey-2.1-be65806c | 301 KB | Fire area survey |
| videocollections-1.1-60fafaa3 | 70 KB | Video collections |
| blockgenerator-1.0-ddc6e9df | 25 KB | Block generator |
| winds-data-1.0-14bc4514 | 3.1 MB | Winds data |
| avo-1.2-0cb4a345 | 89 KB | Aviation operations overlay |
| camera-6.2.1-750de46 | 118 MB | Camera integration |

### Drone/UAS (Tier 3 — functional with SITL)
| Plugin | Size | Function |
|--------|------|----------|
| UAS_Tool/uastool-13.0.4-2fcea3d5 | 311 MB | Core drone plugin (MAVLink/MAVSDK) |
| UAS_Tool/UASTool-DIUBLUE-13.0.4-cac8963a | 129 MB | DIU Blue UAS interface |
| UAS_Tool/ATAKGo-4.18.0.190-CIV-dcd3ce6 | 50 MB | DJI control (v4) |
| UAS_Tool/ATAKGo-v5_5.18.0.81-CIV_988da93 | 148 MB | DJI control (v5, with LDM) |
| UAS_Tool/MAVSDK-Anafi-1.3.0-ede45529 | 63 MB | Parrot Anafi SDK |
| UAS_Tool/MAVSDK-Autel-1.2.0-4a767855 | 61 MB | Autel SDK |
| dronehone-3.3-6bf8a41d | 5.7 MB | Drone detection/telemetry |
| UAS_Tool/ESP32C3_crsfToCot.bin | 969 KB | ESP32 firmware (CRSF→CoT relay) — NOT an APK |

### Hardware-Dependent (Tier 4 — dormant without device)
| Plugin | Size | Needs |
|--------|------|-------|
| taktalk-1.0.0-2b537c1d | 1.4 GB | Mic + TAK server (functional but huge) |
| beartooth-mkii-4.2.0-dc8f9815 | 44 MB | Beartooth mesh radio |
| somewear-atak-1.1.23.0 | 75 MB | Somewear satcom device |
| viscacam-1.0-de013a20 | 2.5 MB | Specific camera hardware |
| helmcam-7.0.1-8123dd5 | 24 MB | Helm cam hardware |
| serialmonitor-6.4.1-c22be48 | 4.4 MB | Serial-connected device |
| vans-1.2-5d45b957 | 80 MB | Specialized VANS UI |
| wasp-5.7.0.1.DEV-10991797 | 13 MB | DEV build (may have bugs) |
| icu-7.2.3-a0f4814 | 37 MB | Ice operations (niche) |

### DJI LDM License
File: `UAS_Tool_5.7.0/license_com.atakmap.android.atakgov5.dji_-_New_DJI_Local_Data_Mode_(LDM)_License_for_ATAK_GO_v5_Only_-_Valid_from_01-05-2026_to_01-05-2027.txt`
Valid: Jan 2026 — Jan 2027

## Device Resource Check

| Resource | ATAK Full Loadout Needs | S24 Ultra Has | Headroom |
|----------|------------------------|---------------|----------|
| Storage | ~3.1 GB (CIV + 30 plugins + TakTracker) | 453 GB free | 0.7% |
| RAM | ~1.5-2 GB peak (all plugins active) | 12 GB (6 GB free) | 80% free |
| CPU | Low (I/O + GPS polling) | SD 8 Gen 3 | Overkill |
| GPU | Map rendering, satellite imagery | Adreno 750 | Overkill |

## Notes
- taktalk is 1.4 GB but functional (voice comms via TAK server + mic)
- wasp is a DEV build — may have bugs, install at own risk
- ESP32C3_crsfToCot.bin is firmware for ESP32 board, NOT an APK — do not adb install
- TakTracker is a standalone APK, not a plugin — install separately alongside ATAK-CIV

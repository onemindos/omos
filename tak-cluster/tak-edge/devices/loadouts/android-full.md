# Loadout: android-full

Standard Android ATAK build for a trusted field device, provisioned over USB.

## What gets installed
- **ATAK-CIV** full release APK (`client/ATAK-*-civ-release.apk`).
- Plugins from `client/loadout/Plugins/`. Install only what the device's mission needs; the full set is large (taktalk ~1.5GB, camera ~123MB, vans ~83MB).

## Gotchas learned on real hardware (SM-T320, LineageOS/Android 15)
- **civSmall is arm64-only.** `ATAK-*-civSmall-release.apk` ships only `lib/arm64-v8a/`, so it will NOT install on 32-bit (`armeabi-v7a`) devices — use the **full civ** APK, which is a fat APK (v7a + arm64 + x86). Check with `unzip -l APK 'lib/*'`.
- **Large-APK install over USB on macOS fails with `splice failed: EINVAL`.** Fix: `adb install --no-streaming -g APK` (legacy push install). Applies to the 368MB civ APK and big plugins.
- **Plugins are version-locked to the ATAK core.** A major core upgrade (e.g. 4.6 → 5.7) means old plugin packages won't load — uninstall and reinstall the matching 5.7 plugin builds.
- **Core plugin set used:** takchat, datasync, reports, networkmonitor, imagemarkup (all installed cleanly on v7a).

## Recommended plugin subsets by role
| Role | Plugins |
|---|---|
| Basic SA / comms | takchat, networkmonitor |
| UAV operator | UAS_Tool_5.7.0 (+ matching ATAKGo for DJI generation), MAVSDK-* for Parrot/Autel |
| Recon / imagery | camera, imagemarkup, geotakcam, helmcam |
| Comms-heavy | takchat, taktalk (needs SIP/VoIP backend), somewear (satellite) |

> Install ONE ATAKGo version matching the drone generation. `ESP32C3_crsfToCot.bin` is firmware, flashed to ESP32-C3 hardware — not an APK.

## Provision (hard-wired)
```bash
# 1) issue the cert + TAK user (once, from deploy/)
deploy/scripts/tak-add-device.sh field-android-01 field-android-01 Red "Team Member"

# 2) plug in device (USB debugging on), then:
devices/scripts/provision-android.sh field-android-01 --plugins
```

## Enable USB debugging (device side)
1. Settings → About phone → tap **Build number** 7×.
2. Settings → Developer options → enable **USB debugging**.
3. Authorize this Mac when the RSA prompt appears (check "always allow").

## After provisioning
- ATAK → Settings → Network → Manage Server Connections → Add → cert + `8089` SSL.
- Confirm callsign appears on WebTAK (`https://192.168.1.183:8446`) and update `devices/registry.md` → `live`.

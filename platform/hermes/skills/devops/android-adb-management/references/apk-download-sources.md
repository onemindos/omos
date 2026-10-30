# APK Download Sources for Android Tactical Setup

## Verified URLs (as of July 2026)

### Tailscale (Android)
- **URL:** `https://pkgs.tailscale.com/stable/tailscale-android-universal-1.98.8.apk`
- **Source:** pkgs.tailscale.com (official)
- **Size:** ~96 MB (universal APK)
- **Package:** `com.tailscale.ipn`
- **Note:** Larger than Play Store version (includes all architectures). Requires Android 8+.

### GPSTest
- **URL:** `https://github.com/barbeau/gpstest/releases/download/v3.10.5/googleRelease-v3.10.5.apk`
- **Source:** GitHub releases (official, by Sean Barbeau / USF)
- **Size:** ~4.2 MB
- **Package:** `com.android.gpstest`
- **Purpose:** GNSS diagnostics — satellite lock, signal strength, sky plot, dual-frequency (L1+L5) verification. Critical for ATAK GPS debugging.

### F-Droid
- **URL:** `https://f-droid.org/repo/org.fdroid.fdroid_2000011.apk`
- **Source:** f-droid.org (official)
- **Size:** ~13 MB
- **Package:** `org.fdroid.fdroid`
- **Purpose:** Open-source app store, no Google account needed.

### Termux
- **URL:** `https://f-droid.org/repo/com.termux_1022.apk`
- **Source:** F-Droid (DO NOT use Play Store version — it's outdated and abandoned)
- **Size:** ~110 MB
- **Package:** `com.termux`
- **Purpose:** Full Linux terminal on Android. SSH, Python, bash, proot.

### Good Lock
- **Source:** Galaxy Store only (`com.sec.android.app.samsungapps`)
- **Package:** `com.samsung.android.goodlock`
- **Requirements:** Samsung account login + Galaxy Store installed
- **Note:** No sideload path. Not on Play Store, F-Droid, or APKMirror (Samsung-ecosystem-locked).

### Galaxy Store (if debloated)
- **Source:** System partition `/system/priv-app/GalaxyApps_OPEN/GalaxyApps_OPEN.apk`
- **Package:** `com.sec.android.app.samsungapps`
- **Restore:** `adb pull` from system + `adb install`

## Download Pattern (from Legacy Pod)

```bash
# Download to pod, SCP to Mac, ADB install to device
curl -sL -o /opt/data/<app>.apk "<URL>"
scp -o "ProxyCommand=..." -i <key> /opt/data/<app>.apk zeus@<mac-tailnet-host>:/tmp/
ssh ... zeus@<mac-tailnet-host> 'adb -s <SERIAL> install /tmp/<app>.apk'
```

## F-Droid Recommended Apps (Tactical)

| App | Package | Purpose |
|-----|---------|---------|
| Termux | com.termux | Terminal, SSH, Python |
| Termux:API | com.termux.api | GPS, camera, battery access from Termux |
| Termux:Boot | com.termux.boot | Auto-start scripts on boot |
| NetGuard | eu.faircode.netguard | Per-app firewall (no root) |
| Orbot | org.torproject.android | Tor on Android |
| OsmAnd~ | net.osmand.plus | Offline maps (ATAK can use) |
| OpenCamera | net.sourceforge.opencamera | Manual camera controls |
| LocalSend | com.zaneschepke.wirelessfiletransfer | AirDrop equivalent |
| DAVx⁵ | at.bitfire.davdroid | CalDAV/CardDAV sync |
| K-9 Mail | com.fsck.k9 | Open-source email |

# F-Droid Bulk Install Session — S24 Ultra (July 2026)

> Device: Samsung S24 Ultra (SM-S928U1), ADB serial <SERIAL>
> Bridge: Legacy pod → Tailscale SSH → Mac (Mac @ <mac-tailnet-host>) → USB → S24
> All commands run via `ssh ... zeus@<mac-tailnet-host> 'export PATH="/opt/homebrew/bin:$PATH"; ...'`

## F-Droid Index Structure

F-Droid's repo index is a signed JAR at `https://f-droid.org/repo/index.jar`.
Inside: `index.xml` (14MB XML), NOT `index-v1.json`. The XML uses this structure:

```xml
<application id="net.osmand.plus">
  <name>OsmAnd~</name>
  <package>
    <version>5.3.8</version>
    <versioncode>530801</versioncode>
    <apkname>net.osmand.plus_530801.apk</apkname>
    <nativecode>armeabi-v7a</nativecode>  <!-- ABI-specific! -->
    ...
  </package>
  <package>
    <versioncode>530803</versioncode>
    <apkname>net.osmand.plus_530803.apk</apkname>
    <nativecode>arm64-v8a</nativecode>
    ...
  </package>
</application>
```

Key: `apkname` is the filename on the repo (`https://f-droid.org/repo/<apkname>`).
`nativecode` tells you which ABI variant it is (when split APKs exist).

## Installed Packages (14 total)

### From F-Droid Main Repo (12)

| Short Name | Package ID | Version Code | Version Name | APK Filename |
|------------|-------------|-------------|---------------|--------------|
| aurora | com.aurora.store | 72 | 4.8.0 | com.aurora.store_72.apk |
| breezyweather | org.breezyweather | 60103 | 6.1.3_freenet | org.breezyweather_60103.apk |
| connectbot | org.connectbot | 11007000 | 1.10.7-oss | org.connectbot_11007000.apk |
| davx5 | at.bitfire.davdroid | 405120004 | 4.5.12-ose | at.bitfire.davdroid_405120004.apk |
| mumla | se.lublin.mumla | 3070007 | 3.7.0 | se.lublin.mumla_3070007.apk |
| netguard | eu.faircode.netguard | 2024122801 | 2.332 | eu.faircode.netguard_2024122801.apk |
| newpipe | org.schabi.newpipe | 1009 | 0.28.4 | org.schabi.newpipe_1009.apk |
| openkeychain | org.sufficientlysecure.keychain | 58902 | 5.8.2 | org.sufficientlysecure.keychain_58902.apk |
| osmand (arm64) | net.osmand.plus | 531003 | 5.3.10 | net.osmand.plus_531003.apk |
| shelter | net.typeblog.shelter | 23 | 1.8 | net.typeblog.shelter_23.apk |
| trailsense | com.kylecorry.trail_sense | 143 | 8.0.0 | com.kylecorry.trail_sense_143.apk |
| gpslogger | com.mendhak.gpslogger | 133 | 133 | com.mendhak.gpslogger_133.apk |
| syncthingfork | com.github.catfriend1.syncthingfork | 2010000 | 2.1.0.0 | com.github.catfriend1.syncthingfork_2010000.apk |

### From GitHub Releases (1)

| Short Name | Package ID | Source | Version |
|------------|-------------|--------|---------|
| orbot | org.torproject.android | guardianproject/orbot releases | 17.9.5-RC-4-tor-0.4.9.11 |

URL: `https://github.com/guardianproject/orbot-android/releases/download/17.9.5-RC-4-tor-0.4.9.11/Orbot-17.9.5-RC-4-tor-0.4.9.11-fullperm-arm64-v8a-release.apk`

### Already on Device (pre-installed)

- com.tailscale.ipn (Tailscale)
- com.termux + com.termux.boot (Termux)
- org.fdroid.fdroid (F-Droid)
- com.atakmap.app.civ (ATAK-CIV + 28 plugins)

## Issues Encountered and Resolved

### 1. OsmAnd ABI Mismatch (INSTALL_FAILED_NO_MATCHING_ABI)

The F-Droid index lists OsmAnd as three split APKs per version:
- versionCode 530801 → armeabi-v7a (default, listed first)
- versionCode 530802 → x86/x86_64
- versionCode 530803 → arm64-v8a

The S24 Ultra is arm64-v8a only (no 32-bit compat). The default (first/last depending on
parse order) was armeabi-v7a, which failed. Fix: explicitly search for the arm64-v8a
variant by checking the `nativecode` child element. The latest version (5.3.10) has
versionCode 531003 for arm64.

### 2. Missing Packages from F-Droid Main Repo

- **com.nutomic.syncthingandroid** — NOT in F-Droid index. The original Syncthing Android
  app was removed from F-Droid. Replaced with `com.github.catfriend1.syncthingfork`
  (actively maintained fork, same functionality).
- **org.torproject.android** (Orbot) — NOT in F-Droid main repo. Downloaded from
  GitHub guardianproject/orbot releases (arm64-v8a variant).
- **eu.basicairdata.graziano.gpslogger** — NOT in F-Droid index. Used
  `com.mendhak.gpslogger` instead (GPSLogger by mendhak, functionally equivalent).

### 3. Python String Quoting in SSH

When passing Python scripts through SSH to the Mac, f-strings with escaped quotes
(`print(f"{info[\"name\"]}")`) fail with SyntaxError due to shell quote nesting.
Fix: use string concatenation instead of f-strings:
```python
print(targets[pkg_id] + "|" + pkg_id + "|" + str(vc) + "|" + str(vn) + "|" + str(apk))
```

### From F-Droid Main Repo (2 additional, installed in second pass)

| Short Name | Package ID | Version Code | Version Name | APK Filename |
|------------|-------------|-------------|---------------|--------------|
| ha-companion | io.homeassistant.companion.android.minimal | 60404000 | 2026.4.4-minimal | io.homeassistant.companion.android.minimal_60404000.apk |
| grocy | xyz.zedler.patrick.grocy | 61 | 3.8.1 | xyz.zedler.patrick.grocy_61.apk |

### From GitHub Releases (1 additional, installed in second pass)

| Short Name | Package ID | Source | Version |
|------------|-------------|--------|---------|
| signal | org.thoughtcrime.securesms | signalapp/Signal-Android releases | 8.19.2 |

URL pattern: `https://github.com/signalapp/Signal-Android/releases/download/v8.19.2/Signal-Android-universal-8.19.2.apk`
Note: signal.org/install redirects, do NOT use it — go directly to GitHub releases.

## All 20 Apps — Final Status

ALL 20 OneMind-recommended apps installed on S24 (July 2026):
- 14 from F-Droid main repo
- 3 from GitHub releases (Signal, Orbot)
- 3 pre-installed (Tailscale, Termux, Obsidian via Play Store)

Verified: `adb shell pm list packages -3 | grep -iE "tailscale|osmand|syncthing|termux|connectbot|securesms|torproject|netguard|keychain|gpslogger|trail|breezy|homeassistant|obsidian|aurora|newpipe|davdroid|shelter|fdroid|mumla|grocy"` → 23 results (includes Termux:Boot, F-Droid, ATAK plugins matching)

## Total Download Size

~330 MB across 14 APKs (largest: OsmAnd arm64 at 162 MB, Syncthing Fork at 65 MB).
All downloads completed in under 2 minutes from the Mac.

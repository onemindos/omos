# Samsung S24 Ultra — OneMind Fabric Device Status Report
## Generated: July 16, 2026 by Legacy (Mother)

---

## DEVICE HARDWARE

| Field | Value |
|-------|-------|
| Model | SM-S928U1 (Samsung Galaxy S24 Ultra) |
| Serial | R5CX11QVGGX |
| CPU Architecture | arm64-v8a (64-bit only, no 32-bit compat) |
| Android Version | 16 |
| SDK Level | 36 |
| Build | BP4A.251205.006.S928U1UES6DZF2 |
| Display | 1080x2340 @ 450dpi |
| Storage | 461GB total, 18GB used, 443GB free (5%) |
| RAM | Not queried (S24 Ultra = 12GB) |
| Battery | 100%, Li-ion, USB charging |
| Carrier | Google Fi (currently no SIM active — Emergency calls only) |
| Samsung Account | None (debloated) |
| Knox Status | Warranty bit: 0, Verified boot: green, Flash locked: 1 (Knox intact, no root) |

## NETWORK

| Interface | IP |
|-----------|-----|
| WiFi (wlan0) | 192.168.1.184/24 (LAN) |
| Tailscale | Active (com.tailscale.ipn installed) |
| Cellular | No active SIM data connection |

## DEVICE STRATEGY

- Debloated Samsung S24 Ultra — no Samsung account, no cloud services
- Burner Google account for Play Services only
- Local backups via Smart Switch/ADB — never cloud
- Knox intact (not rooted) — Secure Folder available natively
- Sovereignty-first: F-Droid primary, Aurora Store secondary, Play Store last resort
- ATAK as primary tactical application (C2 platform)

---

## INSTALLED APPLICATIONS (70 third-party packages)

### TIER 1 — CORE FABRIC (5 apps)

| # | App | Package | Version | Source | Status |
|---|-----|---------|---------|--------|--------|
| 1 | Tailscale | com.tailscale.ipn | — | Direct APK | Pre-installed |
| 2 | OsmAnd+ | net.osmand.plus | 5.3.10 (arm64-v8a) | F-Droid | NEW — installed this session |
| 3 | Syncthing Fork | com.github.catfriend1.syncthingfork | 2.1.0.0 | F-Droid | NEW — installed this session |
| 4 | Termux | com.termux | — | F-Droid | Pre-installed |
| 5 | ConnectBot | org.connectbot | 1.10.7-oss | F-Droid | NEW — installed this session |

### TIER 2 — COMMS & SECURITY (5 apps)

| # | App | Package | Version | Source | Status |
|---|-----|---------|---------|--------|--------|
| 6 | Signal | org.thoughtcrime.securesms | 8.19.2 | GitHub Releases | NEW — installed this session |
| 7 | Mumla (Mumble) | se.lublin.mumla | 3.7.0 | F-Droid | NEW — installed this session |
| 8 | Orbot (Tor) | org.torproject.android | 17.9.5-RC-4 | GitHub Releases | NEW — installed this session |
| 9 | NetGuard | eu.faircode.netguard | 2.332 | F-Droid | NEW — installed this session |
| 10 | OpenKeychain | org.sufficientlysecure.keychain | 5.8.2 | F-Droid | NEW — installed this session |

### TIER 3 — TACTICAL & SITUATIONAL AWARENESS (3 apps)

| # | App | Package | Version | Source | Status |
|---|-----|---------|---------|--------|--------|
| 11 | GPSLogger | com.mendhak.gpslogger | 133 | F-Droid | NEW — installed this session |
| 12 | Trail Sense | com.kylecorry.trail_sense | 8.0.0 | F-Droid | NEW — installed this session |
| 13 | Breezy Weather | org.breezyweather | 6.1.3 | F-Droid | NEW — installed this session |

### TIER 4 — HOMESTEAD & SELF-SUFFICIENCY (3 apps)

| # | App | Package | Version | Source | Status |
|---|-----|---------|---------|--------|--------|
| 14 | HA Companion | io.homeassistant.companion.android.minimal | 2026.4.4-minimal | F-Droid | NEW — installed this session |
| 15 | Grocy | xyz.zedler.patrick.grocy | 3.8.1 | F-Droid | NEW — installed this session |
| 16 | Obsidian | md.obsidian | — | Play Store | Pre-installed |

### TIER 5 — PRIVACY & SOVEREIGNTY (4 apps)

| # | App | Package | Version | Source | Status |
|---|-----|---------|---------|--------|--------|
| 17 | Aurora Store | com.aurora.store | 4.8.0 | F-Droid | NEW — installed this session |
| 18 | NewPipe | org.schabi.newpipe | 0.28.4 | F-Droid | NEW — installed this session |
| 19 | DAVx5 | at.bitfire.davdroid | 4.5.12-ose | F-Droid | NEW — installed this session |
| 20 | Shelter | net.typeblog.shelter | 1.8 | F-Droid | NEW — installed this session |

### TIER 6 — TACTICAL / ATAK ECOSYSTEM (33 apps)

| App | Package | Role |
|-----|---------|------|
| ATAK Civ | com.atakmap.app.civ | Main ATAK C2 application |
| TAK Tracker | gov.tak.taktracker | Location tracking |
| ATAK Helm Cam | com.atakmap.helmcam | Helmet camera integration |
| TAK Cam | com.atakmap.takcam | Camera integration |
| ATAK Reports Plugin | com.atakmap.android.reports.plugin | Reporting |
| ATAK Wasp Plugin | com.atakmap.android.wasp.plugin | Wasp integration |
| ATAK Wind Data Plugin | com.atakmap.android.winddata.plugin | Wind data |
| ATAK Autel MAVSDK | com.atakmap.android.uastool.autelmavsdk | Autel drone control |
| ATAK Ground SDK 8 | com.atakmap.android.uastool.groundsdk8 | Drone ground station |
| ATAK UAS Tool Plugin | com.atakmap.android.uastool.plugin | UAS tools |
| ATAK TakTalk Plugin | com.atakmap.android.taktalk.plugin | Voice comms |
| ATAK TakChat Plugin | com.atakmap.android.takchat.plugin | Chat |
| ATAK VANS Plugin | com.atakmap.android.vans.plugin | VANS integration |
| ATAK ADSB Direct | com.atakmap.android.adsbdirect.plugin | ADS-B aircraft tracking |
| ATAK Fire Survey | com.atakmap.android.firesurvey.plugin | Fire survey |
| ATAK ATAKGo DJI | com.atakmap.android.atakgo.dji | DJI drone Go |
| ATAK ATAKGo V5 DJI | com.atakmap.android.atakgov5.dji | DJI V5 SDK |
| ATAK Data Sync | com.atakmap.android.datasync.plugin | Data synchronization |
| ATAK AR Video Overlay | com.atakmap.android.arvideooverlay.plugin | AR video |
| ATAK GRG Builder | com.atakmap.android.grgbuilder.plugin | Grid reference generator |
| ATAK Visca Cam | com.atakmap.android.viscacam.plugin | VISCA camera control |
| ATAK Compass Nav | com.atakmap.android.compassnav.plugin | Compass navigation |
| ATAK Drone Hone | com.atakmap.android.dronehone.plugin | Drone hone |
| ATAK Geo Cam | com.atakmap.android.geocam.plugin | Geo camera |
| ATAK Block Generator | com.atakmap.android.blockgenerator.plugin | Block generation |
| ATAK ICET | com.atakmap.android.icet.plugin | ICET |
| ATAK Wave Plugin | com.atakmap.android.wave.plugin | Wave comms |
| ATAK COT Utility | com.atakmap.android.cot_utility.plugin | CoT utility |
| ATAK Video Collections | com.atakmap.android.videocollections.plugin | Video management |
| ATAK VNS Plugin | com.atakmap.android.vns.plugin | VNS |
| Beartooth TAK | com.beartooth.beartoothtakplugin | Beartooth mesh radio |
| Somewear TAK | com.somewearlabs.swtak.plugin | Somewear satellite comms |
| GPSTest | com.android.gpstest | GPS diagnostics |

### TIER 7 — PRODUCTIVITY & UTILITIES (10 apps)

| App | Package | Source |
|-----|---------|--------|
| F-Droid | org.fdroid.fdroid | F-Droid |
| Brave Browser | com.brave.browser | Play Store |
| Fastmail | com.fastmail.app | Play Store |
| 1Password | com.onepassword.android | Play Store |
| Todoist | com.todoist | Play Store |
| Clean My Phone | com.cleanmyphone.easy.clean | Play Store |
| OpenKuma | app.openkuma | Unknown |
| ParTech Image Markup | com.par.imagemarkup | Unknown |
| ParTech Geo Camera | com.partech.geocamera | Unknown |
| ParTech Network Monitor | com.partech.networkmonitor | Unknown |
| ParTech Serial Manager | com.partech.serialmanager | Unknown |

### SAMSUNG SYSTEM (190 packages, stock + debloated)

Samsung system packages remain on device (Knox intact, not rooted).
Good Lock + Sound Assistant installed for UI customization.
Secure Folder (Knox) available as native alternative to Shelter.

---

## APP SOURCES USED

| Source | Status | Package Count Installed |
|--------|--------|--------------------------|
| F-Droid (f-droid.org/repo) | Primary | 14 apps |
| GitHub Releases (signal.org, guardianproject) | Secondary | 2 apps (Signal, Orbot) |
| Google Play Store | Pre-existing | ~10 apps |
| Aurora Store | Available as alternative to Play Store | 0 new (app itself from F-Droid) |
| Samsung Galaxy Store | Pre-existing on device | Unknown (Samsung system apps) |
| Direct APK sideload | Used for Tailscale | Pre-installed |

---

## INSTALLATION NOTES

1. **OsmAnd+** — F-Droid default APK (versionCode 530801) was armeabi-v7a only. S24 Ultra is arm64-v8a with no 32-bit compatibility. Downloaded arm64-v8a build (versionCode 531003, v5.3.10) instead.
2. **Syncthing** — Original nutomic build (com.nutomic.syncthingandroid) not available on F-Droid. Used actively maintained fork (com.github.catfriend1.syncthingfork, v2.1.0.0).
3. **Orbot** — Not on F-Droid main repo. Pulled from GitHub guardianproject/orbot releases (v17.9.5-RC-4, arm64-v8a, fullperm).
4. **Signal** — Not on F-Droid. signal.org direct URL returned 404. Downloaded from GitHub signalapp/Signal-Android releases (v8.19.2, universal).
5. **GPSLogger** — Package ID on F-Droid is com.mendhak.gpslogger (not eu.basicairdata.graziano.gpslogger as initially assumed).
6. **HA Companion** — F-Droid has the "minimal" variant (io.homeassistant.companion.android.minimal), which is the F-Droid-compliant build without proprietary dependencies.
7. **Grocy** — Package ID is xyz.zedler.patrick.grocy (the Grocy Android client, not the server).
8. All ADB installs show installerPackageName=null (sideloaded), except Obsidian (Play Store) and Play Services framework.
9. Installation performed from Mac (zeus@100.74.182.117) via ADB over USB/Tailscale, with Legacy orchestrating from DOKS cloud pod via SSH through Tailscale.

---

## PACKAGES NOT YET CONFIGURED

These apps are installed but require on-device setup:

- **Syncthing Fork** — Needs folder pairing with Mac/HA
- **HA Companion** — Needs HA Yellow URL (100.64.25.12) and auth token
- **Signal** — Needs phone number verification
- **Mumla/Mumble** — Needs server address (self-hosted Mumble server)
- **Orbot** — Needs Tor circuit configuration
- **DAVx5** — Needs CalDAV/CardDAV server URL
- **Grocy** — Needs Grocy server URL (self-hosted)
- **OpenKeychain** — Needs PGP key import/creation
- **NetGuard** — Needs per-app firewall rules configuration
- **Shelter** — Needs work profile setup
- **GPSLogger** — Needs output path/GPX upload endpoint configuration
- **Trail Sense** — Needs compass/barometer calibration

---

## DEVICE SECURITY POSTURE

| Layer | Status |
|-------|--------|
| Root | No (Knox intact, flash locked) |
| Verified Boot | Green (stock firmware) |
| Samsung Account | None |
| Google Account | Burner account only |
| Cloud Backup | Disabled — local ADB/Smart Switch only |
| Firewall | NetGuard installed (needs configuration) |
| Tor | Orbot installed (needs configuration) |
| Work Profile | Shelter installed (Samsung Secure Folder also native) |
| Encryption | Full disk encryption (Samsung default) |
| ADB | Enabled (USB debugging on) |
| MDM/Knox Guard | Dormant (not enrolled) |

---

## RECOMMENDED NEXT STEPS

1. **Configure NetGuard** — Block all Samsung/Google telemetry apps from phoning home
2. **Set up Syncthing** — Pair S24 ↔ Mac ↔ HA Yellow for sovereign file sync
3. **Configure HA Companion** — Connect to HA Yellow at 100.64.25.12
4. **Set up Mumble server** — Deploy Mumble server on fabric for tactical voice
5. **Calibrate Trail Sense** — Compass, barometer, sun position for field ops
6. **Configure GPSLogger** — Set GPX output and PostGIS/TAK COP feed
7. **Set up DAVx5** — Connect to self-hosted Radicale/Nextcloud for contacts/calendar sync
8. **Verify Signal** — Complete phone number registration
9. **Import PGP keys** — Load existing keys into OpenKeychain
10. **Configure Shelter** — Set up work profile for untrusted apps
11. **Aurora Store** — Log in anonymously for any remaining Play Store-only apps

---

## ADDITIONAL APP SOURCES (beyond Play Store, F-Droid, Aurora, Samsung Store)

1. **Obtainium** (F-Droid) — Auto-updates apps directly from their GitHub/GitLab releases. Eliminates the need to manually check for updates. Each app is fetched from its own source repo, signed and verified.
2. **APKMirror** (apkmirror.com) — Community-maintained APK mirror. Signed APKs from Play Store, verified by checksum. Good for region-locked or older versions. Web-only, no app client.
3. **Aptoide** (F-Droid or direct) — Independent app store with per-developer repositories. Less curated than F-Droid but larger selection. Multiple community stores.
4. **IzzyOnDroid** (F-Droid repo) — Additional F-Droid repository with apps that dont meet F-Droids strict anti-commercial criteria but are still open source. Add to F-Droid as a secondary repo.
5. **GitHub/GitLab Releases** — Many open-source apps publish APKs directly in their releases. We already used this for Signal and Orbot. Obtainium automates this.
6. **F-Droid Custom Repos** — Any developer can host their own F-Droid-compatible repository. Useful for self-hosted or private app distribution across the fabric.
7. **AAB to APK services** — Bundle to APK extraction (e.g., APKMirror, bundletool). For apps that only distribute Android App Bundles.
8. **XDA Labs** (labs.xda-developers.com) — Developer-focused APK distribution for beta builds, mods, and tools. Useful for rooted devices or custom ROMs.
9. **Direct developer APK** — Many apps (Signal, Tailscale, Tor Browser) offer direct APK downloads from their websites. Most sovereign approach — no middleman at all.

### Sovereignty-Ranked Sources (best to worst for privacy):

1. Direct developer website/GitHub releases (no middleman)
2. F-Droid main repo (verified, reproducible builds)
3. IzzyOnDroid F-Droid repo (open source, less strict criteria)
4. Obtainium (automates direct-from-source)
5. Aurora Store (anonymous Play Store proxy)
6. APKMirror (community mirror, signed APKs)
7. Aptoide (independent stores, less curated)
8. Google Play Store (least sovereign, requires Google account)
9. Samsung Galaxy Store (locked to Samsung ecosystem)

---

## DOCUMENT METADATA

- Generated by: Legacy (Mother) — OneMind OS Sovereign Agent
- Location: DOKS Cloud Hub (hermes-legacy pod)
- Orchestration path: Legacy (cloud) → SSH via Tailscale → Mac (zeus@100.74.182.117) → ADB → S24 Ultra (R5CX11QVGGX)
- Date: July 16, 2026
- Document version: 1.0

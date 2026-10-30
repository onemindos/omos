# F-Droid Recommendations — Derived from OneMind Project

> Source: OneMind vault (Convergence docs, SO/Infrastructure, Homestead Stack,
> Infection Stack, Samsung Tactical Platform doc, Security/Shield comms, TAK loadouts).
> Generated July 2026 by analyzing `/Users/zeus/onemind (Before clean)/` and
> `/Users/zeus/Convergence docs/` on the Mac via SSH.

## Method

The recommendations below were derived by:
1. SSH to Mac → read OneMind project structure (4 pillars: SO, HP, LE, GE)
2. Read Samsung_S23_Tactical_Mobile_Platform.md (app ecosystem, NATS integration)
3. Read Full_Self_Sovereign_Homestead_Stack.md (self-hosted services, mesh tools)
4. Read OneMind_Infection_Stack.md (device onboarding: Tailscale + NATS + KubeSolo)
5. Read TAK loadouts (android-full.md, provision-android.sh)
6. Read Security/Shield/comms (encrypted messaging, VPN, privacy)
7. Cross-referenced with OneMind's sovereignty-first, local-first, air-gapped-capable doctrine

## Tier 1: CORE FABRIC (maps directly to NATS/TAK/Mesh architecture)

| App | F-Droid Package | OneMind Alignment |
|-----|----------------|-------------------|
| **Tailscale** | (direct APK or IzzyOnDroid) | Infection Stack Layer 1: every device joins tailnet |
| **OsmAnd~** | net.osmand.plus | Offline vector maps; custom tile overlays; supplements ATAK maps; Homestead Stack terrain/mapping |
| **Syncthing** | com.github.catfriend1.syncthingfork | Homestead Stack explicit entry: P2P mesh-native file sync, no cloud. NOTE: original nutomic package removed from F-Droid — use the catfriend1 fork |
| **Termux** | com.termux | Infection Stack compute layer; SSH/NATS leaf node |
| **ConnectBot** | org.connectbot | Quick SSH GUI when Termux not convenient; fabric access |

## Tier 2: COMMS & SECURITY (Security/Shield + tactical encryption)

| App | F-Droid Package | OneMind Alignment |
|-----|----------------|-------------------|
| **Signal** | (direct APK — not F-Droid) | Samsung Tactical doc: encrypted messaging; E2E family/ops comms |
| **Mumble** | se.lublin.mumla | Convergence doc (Mumble.md): PTT VoIP; pairs with LiveKit PTT architecture |
| **Orbot** | org.torproject.android | Sovereignty layer; "no external leakage" air-gapped doctrine. NOTE: not on F-Droid main repo — download from GitHub guardianproject/orbot releases (arm64-v8a variant) |
| **NetGuard** | eu.faircode.netguard | Per-app firewall (no root); sovereignty-first debloated phone strategy |
| **OpenKeychain** | org.sufficientlysecure.keychain | PGP key management; aligns with X.509/mTLS cert infrastructure for TAK |

## Tier 3: TACTICAL & SITUATIONAL AWARENESS (TAK/COP/ISR doctrine)

| App | F-Droid Package | OneMind Alignment |
|-----|----------------|-------------------|
| **BasicAirData GPS Logger** | com.mendhak.gpslogger | Position track logging → PostGIS/TAK COP; mission replay. NOTE: original eu.basicairdata package not on F-Droid — com.mendhak.gpslogger is the equivalent |
| **Trail Sense** | com.kylecorry.trail_sense | Offline compass/barometer/GPS/weather; air-gapped field ops scenario |
| **Breezy Weather** | org.breezyweather | Weather data; aligns with weewx/emonCMS sensor doctrine |

## Tier 4: HOMESTEAD & SELF-HOSTED (Full Self-Sovereign Homestead Stack)

| App | F-Droid Package | OneMind Alignment |
|-----|----------------|-------------------|
| **Home Assistant** | (direct APK or Play Store) | HA Yellow at <ha-tailnet-host> with 185+ entities |
| **Grocy** | (browser PWA) | Homestead Stack explicit entry: pantry tracking, barcode scanning |
| **Obsidian** | (direct APK) | Knowledge base platform; Convergence docs/OneMind vault; Syncthing sync |

## Tier 5: PRIVACY & SOVEREIGNTY (replace Google services)

| App | F-Droid Package | OneMind Alignment |
|-----|----------------|-------------------|
| **Aurora Store** | com.aurora.store | Anonymous Play Store access; reduces "burner Google account" exposure |
| **NewPipe** | org.schabi.newpipe | YouTube without tracking; no account needed |
| **DAVx⁵** | at.bitfire.davdroid | CalDAV/CardDAV; sovereign alternative to Google Contacts/Calendar |
| **Shelter** | net.typeblog.shelter | Work profile isolation; sandbox untrusted apps |

## Install Priority (highest immediate impact first)

```
1. Syncthing      — immediate fabric value, file sovereignty
2. OsmAnd~        — fixes broken ATAK maps problem (offline tile fallback)
3. NetGuard       — firewall on debloated phone = clean sovereignty
4. Mumble (Mumla) — tactical voice (pairs with LiveKit PTT stack)
5. Aurora Store   — escape Google dependency for remaining apps
6. GPS Logger     — feed position data to PostGIS/TAK COP
7. Breezy Weather — sensor data, offline-capable
```

## Key Insights from OneMind Project

- `/Users/zeus/onemind` does NOT exist as a directory — the correct path is
  `/Users/zeus/onemind (Before clean)/` (with the trailing "(Before clean)")
- The OneMind project has NO explicit "F-Droid list" file — recommendations must
  be synthesized from the architecture docs
- The Full Self-Sovereign Homestead Stack explicitly lists Syncthing, Grocy, HA,
  OpenSprinkler, farmOS, Frigate NVR, and other self-hosted tools
- The Infection Stack defines the 4-layer device onboarding: Tailscale → NATS leaf →
  KubeSolo/NEX → Portainer Edge Agent
- The Samsung Tactical Platform doc references: Signal, GrapheneOS, NATS clients,
  OpenAndroidMaps, DDF catalog clients
- Mumble has its own Convergence doc (Mumble.md) — it's an explicit OneMind component

## Notes

- These recommendations are for the S24 Ultra (SM-S928U1) running Android 14+
- F-Droid version of Termux is the only maintained one (Play Store is dead)
- OsmAnd~ from F-Droid is the FOSS version (equivalent to OsmAnd+ paid)
- Signal is NOT on F-Droid (uses Google libraries) — install from signal.org/android/apk
- Home Assistant Companion not on F-Droid — use direct APK from HA GitHub releases

---
name: TAK client onboarding (iTAK / ATAK)
description: Connect iPhone/iPad (iTAK) and Android (ATAK) tactical devices to the self-hosted OneMind TAK Server over Tailscale. Use this whenever adding, re-enrolling, or troubleshooting a device connection. Covers the ONE reliable route (soft-cert data package) and the dead-ends that wasted hours — read the "Do NOT repeat" section before trying anything else.
---

# TAK client onboarding

## TL;DR — the reliable route
Build a **soft-cert data package** (client cert + CA truststore bundled) and AirDrop / push it to the
device. That's it. No enrollment, no passwords, no trust prompts.

```bash
cd /Users/zeus/ATAK
CERTS=server/standard/takserver-docker-5.7-RELEASE-43/tak/certs/files
# <user> <client.p12> <ca-truststore.p12> <host> <out.zip>
devices/scripts/build-itak-datapackage.sh omos-iphone17 \
  "$CERTS/omos-iphone17.p12" "$CERTS/truststore-root.p12" \
  100.74.182.117 deploy/exports/OneMind-iPhone17.zip
```
Then: AirDrop the `.zip` (macOS) → on device tap the AirDrop notification → **Open in iTAK** (iOS) or
push to `/sdcard/atak/tools/datapackage/` and import via Import Manager (Android/ATAK). iTAK auto-adds
server **OneMind TAK** with the cert loaded and connects.

**Verified working 2026-07-02:** iPhone connected on first import (server log:
`Set client for subscription: tls:3 to iPhone (...) source=192.168.65.1`).

## Why this works (the one fact that matters)
**iTAK does NOT use iOS's system certificate trust.** It keeps its own trust store. So on a
**self-signed** TAK Server, the CA has to travel *inside* the data package. Put the CA there and trust
is solved; the bundled client cert is the login (mutual TLS on port 8089).

## Do NOT repeat these dead-ends (the loop we already burned hours on)
| Tried | Why it failed |
|---|---|
| iTAK **Connect with Credentials** (user/pass enrollment) | iTAK won't complete TLS to a server it doesn't trust; fails "error authenticating with the server" and **never reaches the server** (zero log trace). Only works AFTER the CA is already in iTAK via a package. |
| Installing the CA as an **iOS profile** + Certificate Trust Settings | Fixes **Safari** but not iTAK — iTAK ignores iOS system trust. Harmless but useless for iTAK. |
| Hand-built package with **two manifests** (`MANIFEST/manifest.xml` + root `manifest.xml`) | iTAK: "the format of this configuration could not be read." Use exactly ONE manifest at `MANIFEST/manifest.xml`. |
| Manual CSR to `/Marti/api/tls/signClient` with only `CN=` | Server 500 "CSR validation failed!" — the CSR subject MUST include `O=TAK, OU=TAK-OPS` (the `nameEntries`). Real clients read these from `tls/config` first. This is a server-signer red herring, not a real fault. |

## Exact data-package format (why hand-builds fail)
Matches maintained refs (talentedbrute/takserver-datapackage-builder, Cloud-RF/tak-server certDP.sh),
verified with iTAK 5.x. The build script encodes all of this — prefer the script over hand-editing.
- **One** manifest: `MANIFEST/manifest.xml`, `version="2"`. Never also a root `manifest.xml`.
- Manifest `zipEntry` names are **flat** (`omos-iphone17.p12`) — **no** `cert/` prefix.
- `secure.pref`: cert paths **do** use the `cert/` prefix (`cert/omos-iphone17.p12`,
  `cert/truststore.p12`), as **non-indexed** `caLocation`/`certificateLocation` in the
  `com.atakmap.app_preferences` block. `connectString0` stays indexed in `cot_streams`.
- p12 files sit **flat at the zip root**; iTAK relocates them into its internal `cert/` dir on import.
- p12 password = the keystore password (TAK default until step-2 rotation), **legacy-encrypted** → read with `openssl ... -legacy`.

## Verify BEFORE handing a package to the user
```bash
# 1) does the client cert authenticate to the server over real mTLS? (expect: Verify return code: 0)
d=$(mktemp -d); P=$CERTS/omos-iphone17.p12
openssl pkcs12 -in "$P" -passin env:CAPASS -legacy -nokeys -clcerts -out $d/c.pem 2>/dev/null
openssl pkcs12 -in "$P" -passin env:CAPASS -legacy -nocerts -nodes -out $d/k.pem 2>/dev/null
openssl pkcs12 -in "$CERTS/truststore-root.p12" -passin env:CAPASS -legacy -nokeys -cacerts -out $d/ca.pem 2>/dev/null
echo | openssl s_client -connect 100.74.182.117:8089 -cert $d/c.pem -key $d/k.pem -CAfile $d/ca.pem 2>/dev/null | grep "Verify return code"
```
```bash
# 2) after the user connects, confirm on the server:
docker exec takserver bash -lc 'grep -iE "Set client for subscription" /opt/tak/logs/takserver-messaging.log | tail -5'
# NOTE: Docker Desktop NATs the device's Tailscale IP, so connections show source=192.168.65.1, NOT the
# device's 100.x Tailscale address. Don't be fooled by that — match on timing + callsign instead.
```

## Server facts (OneMind TAK)
- Host: `100.74.182.117` (Tailscale) / `192.168.1.183` (LAN). Container: `takserver`.
- Ports: **8089** CoT / client mTLS · **8446** cert enrollment (HTTPS) · **8450** WebTAK.
- CA: `TAK-CA`, truststore `truststore-root.p12`, all p12 passwords are the keystore password (default until step-2 rotation).
- Server cert SANs include both `100.74.182.117` and `192.168.1.183` (+ `localhost`, DNS `takserver`).
- Credential-enrollment users (fallback only): `omos-iphone17` / `omos-ipad`, password
  set per device; never commit, group Red. Users persist in bind-mounted `UserAuthenticationFile.xml`.

## Device registry
| Device | cert CN / user | package | callsign | team |
|---|---|---|---|---|
| iPhone 17 | `omos-iphone17` | `deploy/exports/OneMind-iPhone17.zip` | `zeus` | Red |
| iPad | `omos-ipad` | `deploy/exports/OneMind-iPad.zip` | `zeus-pad` | Red |
| Android tablet | `omos-tablet` | `deploy/exports/OneMind-Tablet.zip` | `zeus-tab` | Red |

After connecting, set identity in-app: Settings → Identity → callsign + Team **Red** (default callsign
is "iPhone"/"iPad" until changed).

## Related
- [scripts/build-itak-datapackage.sh](scripts/build-itak-datapackage.sh) — the generator.
- [scripts/ios-enroll.md](scripts/ios-enroll.md) — iOS-specific steps + fallbacks.
- [../deploy/docs/enrollment.md](../deploy/docs/enrollment.md) — server-side enrollment config + the
  `signClient` CSR subject-DN gotcha and `JDK_JAVA_OPTIONS` add-opens dependency.

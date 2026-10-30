# iOS enrollment

iOS can't be driven by adb, and iTAK is sandboxed — a Mac can't set its server connection remotely.

**KEY FACT (learned the hard way):** iTAK does **NOT** use iOS's system certificate trust. Installing
the CA as an iOS profile (Settings → Certificate Trust) fixes *Safari* but does nothing for iTAK — iTAK
keeps its own trust store. So with a **self-signed** TAK Server:
- **Connect with Credentials fails** ("error authenticating with the server") — iTAK won't complete TLS
  to a server it doesn't trust, and never even reaches the server (zero log trace).
- The reliable route is a **soft-cert data package** that bundles the client cert **and** the CA
  truststore. iTAK loads the CA from the package into its own store → trust solved, no enrollment, no
  prompt. The cert *is* the login.

Full server-side setup + verification: [deploy/docs/enrollment.md](../../deploy/docs/enrollment.md).

## Path 0 — Soft-cert data package (THE reliable route)
Build it with the dedicated script (format matches the maintained
talentedbrute / Cloud-RF references, verified with iTAK 5.x):
```bash
CERTS=server/standard/takserver-docker-5.7-RELEASE-43/tak/certs/files
devices/scripts/build-itak-datapackage.sh omos-iphone17 \
  "$CERTS/omos-iphone17.p12" "$CERTS/truststore-root.p12" \
  100.74.182.117 deploy/exports/OneMind-iPhone17.zip
```
Then:
1. AirDrop `deploy/exports/OneMind-iPhone17.zip` from the Mac → iPhone.
2. On the iPhone: tap the AirDrop notification (or Files → the zip) → **Open in iTAK** / Share → iTAK.
3. iTAK imports it and auto-adds server **OneMind TAK** (`100.74.182.117:8089:ssl`) with the cert
   already loaded — no username/password, no cert-trust prompt.
4. Open iTAK's server list → ensure **OneMind TAK** is enabled → it connects (green).
5. Settings → Identity → callsign (iPhone `zeus`, iPad `zeus-pad`), Team **Red**.

**Critical format details** (why earlier hand-built packages failed "format could not be read"):
- Exactly **one** manifest at `MANIFEST/manifest.xml` — never also a stray root `manifest.xml`.
- `manifest.xml` `zipEntry` names are **flat** (e.g. `omos-iphone17.p12`), no `cert/` prefix.
- `secure.pref` cert paths **do** use the `cert/` prefix (`cert/omos-iphone17.p12`,
  `cert/truststore.p12`) with **non-indexed** `caLocation`/`certificateLocation` in the
  `com.atakmap.app_preferences` block; `connectString0` stays indexed in `cot_streams`.
- p12 files sit **flat at the zip root**; iTAK relocates them into its internal `cert/` dir on import.
- p12 password is the keystore password (TAK default until step-2 rotation; see Secret tak-keystores in ns takserver) (legacy-encrypted — read with `openssl ... -legacy`).

Verify the cert authenticates *before* handing it to the user:
```bash
# extract from the p12 and do a real mTLS handshake to the CoT port
openssl s_client -connect 100.74.182.117:8089 -cert c.pem -key k.pem -CAfile ca.pem </dev/null 2>/dev/null \
  | grep "Verify return code"      # → 0 (ok) means the server accepts this client cert
```

Same `.zip` works on Android/ATAK: `adb push OneMind-iPhone17.zip /sdcard/atak/tools/datapackage/`
then import via ATAK Import Manager (ATAK also accepts the `MANIFEST/manifest.xml` layout).

## Path C — Connect with Credentials (only works with CA pre-trusted)
Requires the CA already in iTAK's store (i.e. after a Path 0 import). Standalone on a self-signed
server it fails. In iTAK → add server → **Connect with Credentials**:
- **Address:** `100.74.182.117`  **Port:** `8089`  **Protocol:** `SSL`
- **Username:** `omos-iphone17` / `omos-ipad`   **Password:** set per device; never commit

## Path A — HTTPS enrollment (no cables)
1. Device on the same network (LAN or Tailscale).
2. Safari → `https://192.168.1.183:8446/Marti/enrollment/` (or Tailscale IP).
3. Accept the cert warning → log in with the device's TAK user → download the enrollment profile.
4. Settings → General → VPN & Device Management → install the profile.
5. Open iTAK — server auto-configures. Confirm callsign on WebTAK.

## Path B — Cert file import (AirDrop or USB)
1. Issue cert: `deploy/scripts/tak-add-device.sh <cn> <user> Red "Team Member"`.
2. AirDrop `deploy/exports/<cn>.p12` and `server/.../certs/files/ca.pem` to the device.
3. Open each file → iOS installs the profile (Settings → VPN & Device Management).
4. Certificate Trust Settings → enable full trust for **TAK-CA**.
5. iTAK → add server → `8089` SSL/TLS → select the imported cert.

## Optional: scripted over USB with Apple Configurator
```bash
brew install --cask apple-configurator   # provides cfgutil
cfgutil list                             # tether device, read ECID
# Wrap ca.pem (+ client identity) into a .mobileconfig, then:
cfgutil -e <ECID> install-profile tak-enroll.mobileconfig
```
Use this when onboarding several iPhones repeatably. For one or two devices, Path A is faster.

# TAK Server — Certificate Enrollment ("Connect with Credentials")

How the iOS/Android clients get a cert automatically by logging in with a username + password,
instead of hand-importing a `.p12`. This is the **recommended** path for iTAK (its data-package
importer is fussy; enrollment "just works").

Set up and verified 2026-07-02. If something here stops making sense, this doc is the source of truth
for how the server was configured.

---

## What the client enters (iTAK → "Connect with Credentials")

| Field | Value |
|---|---|
| Name | `OneMind TAK` (any label) |
| Address | `100.74.182.117` (Tailscale) or `192.168.1.183` (LAN) |
| Port | `8089`  ← the CoT port, NOT the enrollment port |
| Protocol | `SSL` |
| Username | `omos-iphone17` / `omos-ipad` (per device) |
| Password | `OneMind-TAK-2026` |

The client connects CoT on **8089** and silently enrolls its cert on **8446**. Accept the
"trust server certificate" prompt on first connect (it's fetching our self-signed CA). Set the
callsign afterwards (iPhone `zeus`, iPad `zeus-pad`, Team Red).

---

## How enrollment is wired on the server (4 pieces)

### 1. `certificateSigning` block in CoreConfig
Lets the server sign client CSRs with our CA. **Placement matters** — per `CoreConfig.xsd` the
element order is `... dissemination → certificateSigning → logging → security ...`, so it is a
**direct child of `<Configuration>`, immediately before `<security>`** (NOT inside `<security>`).
Putting it inside `<security>` makes the API log `CertificateSigning element not found in CoreConfig!`
and return HTTP 500.

```xml
<dissemination smartRetry="false"/>
<certificateSigning CA="TAKServer">
    <certificateConfig>
        <nameEntries>
            <nameEntry name="O" value="TAK"/>
            <nameEntry name="OU" value="TAK-OPS"/>
        </nameEntries>
    </certificateConfig>
    <TAKServerCAConfig
        keystore="JKS"
        keystoreFile="/opt/tak/certs/files/ca-signing.jks"
        keystorePass="atakatak"
        validityDays="365"
        signatureAlg="SHA256WithRSA"
        CAkey="/opt/tak/certs/files/ca-plain.key"
        CAcertificate="/opt/tak/certs/files/ca.pem"/>
</certificateSigning>
<security> ... </security>
```
`CAkey` is the **unencrypted** CA private key (`openssl rsa -in ca-do-not-share.key -passin pass:atakatak
-out ca-plain.key`); `CAcertificate` is `ca.pem`. Both are gitignored (plaintext key).

### 2. CA signing keystore — `certs/files/ca-signing.jks`
A JKS holding the CA cert **and its private key** so the server can sign. Built from `ca.pem` +
`ca-do-not-share.key` (CA password `atakatak`):

```bash
# in a throwaway container with certs/files mounted (openssl/keytool live there)
openssl pkcs12 -export -legacy -in ca.pem -inkey ca-do-not-share.key -passin pass:atakatak \
  -out ca-signing.p12 -passout pass:atakatak -name ca
keytool -importkeystore -srckeystore ca-signing.p12 -srcstoretype PKCS12 -srcstorepass atakatak \
  -destkeystore ca-signing.jks -deststoretype JKS -deststorepass atakatak -destkeypass atakatak -noprompt
```
(`certs/files` is bind-mounted read-only into `takserver`; the file just needs to exist on the host.)

### 3. Standard port mapping — enrollment on host **8446**
Clients hard-assume enrollment at `<server>:8446`. The compose was remapped to TAK-standard 1:1:

```yaml
ports:
  - "8089:8089"   # CoT TLS — device connections
  - "8450:8443"   # WebTAK UI + REST API (host 8443 is taken → use 8450)
  - "8444:8444"   # Federation
  - "8446:8446"   # Cert enrollment (standard; do NOT remap)
```
Old broken mapping was `8446→8443` (WebTAK) + `8448→8446` (enrollment) — clients tried 8446, hit
WebTAK (which demands a client cert), and enrollment failed.

### 4. Persistent user file — `tak/UserAuthenticationFile.xml` bind-mount
**Critical gotcha:** `UserAuthenticationFile.xml` lives inside the container and is **wiped on every
`docker compose up` recreate** unless bind-mounted. Users (both password and cert entries) vanished
after the port change until this mount was added:

```yaml
volumes:
  - ./tak/UserAuthenticationFile.xml:/opt/tak/UserAuthenticationFile.xml   # persist users
```

---

## Managing enrollment users

```bash
# password user (for credential enrollment) — creates or updates
docker exec takserver bash -lc \
  "cd /opt/tak && java -jar utils/UserManager.jar usermod -p '<password>' -g Red <username>"

# cert user (for pre-issued .p12 / data-package auth, e.g. the tablet)
docker exec takserver bash -lc \
  "cd /opt/tak && java -jar utils/UserManager.jar certmod -g Red /opt/tak/certs/files/<cn>.pem"

# list users
grep -oE 'identifier="[^"]+"' tak/UserAuthenticationFile.xml
```

Current users (all group Red): `omos-iphone17`, `omos-ipad` (password `OneMind-TAK-2026`),
`omos-tablet` (cert-based, uses its data-package `.p12`).

---

## Verify the endpoint (from the host)

```bash
# anonymous → 401 (endpoint alive, auth-gated)
curl -sk -o /dev/null -w '%{http_code}\n' https://100.74.182.117:8446/Marti/api/tls/config
# with creds → 200 + <certificateConfig> XML
curl -sk -u 'omos-iphone17:OneMind-TAK-2026' https://100.74.182.117:8446/Marti/api/tls/config
```
200 = enrollment ready. 500 = `certificateSigning` misplaced/missing. 401 = bad creds or user file
not loaded. 000 = wrong port (enrollment is 8446, not 8443/8450).

### Proving the signer works end-to-end (not just config)
`tls/config` returning 200 only proves the config is *readable*. To prove the server actually **signs**,
POST a real CSR — but the **CSR subject MUST include the O/OU from `nameEntries`** or TAK's
`validateCSR` returns false → `CSR validation failed!` → HTTP 500. Real clients do this automatically
(they GET `tls/config`, read the nameEntries, then build the CSR). Manual test:
```bash
openssl req -new -newkey rsa:2048 -nodes -keyout t.key -out t.csr \
  -subj "/O=TAK/OU=TAK-OPS/CN=omos-iphone17"          # O + OU MUST match nameEntries
openssl req -in t.csr -outform DER | base64 | tr -d '\n' > t.b64
curl -sk -u 'omos-iphone17:OneMind-TAK-2026' --data-binary @t.b64 -H 'Content-Type: text/plain' \
  "https://100.74.182.117:8446/Marti/api/tls/signClient/v2?clientUid=test&version=2"
# → 200 + JSON {"signedCert":"<base64 DER>"}  = signer fully working
```
Gotcha history: a bare `CN=...` CSR gives a misleading `CSR validation failed!` 500 that looks like a
signer/keystore fault but is really just a subject-DN mismatch. Also note TAK's signing uses JVM-internal
`sun.security.pkcs10/x509` classes — the container's `setenv.sh` supplies the required `--add-opens`
via `JDK_JAVA_OPTIONS` (don't strip them, or signing breaks).

---

## Why not the data package for iTAK?
It can work, but iTAK requires a **flat** zip (certs + `.pref` at root, no `cert/` subfolder) or it
errors "the format of this configuration could not be read." Credentials enrollment sidesteps all
of that. The generator (`build-datapackage.sh`) now produces the flat layout for when you do want it
(and for ATAK, which is less picky).

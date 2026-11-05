# TAK CA signing-key password rotation, step 1 (2026-10-29)

Owner-approved "step 1": rotate only the password that protects the **CA signing key**.
Device-facing passwords (8089/8443/8446 `keystorePass`/`truststorePass`, client `.p12`,
truststore handed to devices) stay at the stock value until **step 2**. No EUD-facing
workload was restarted and no connection changed.

## CA key password consumers (verified 2026-10-29)

| Consumer | Where the CA key lives | How the password is read | Pickup | Status |
|---|---|---|---|---|
| TAK Server enrollment v2 (`takserver-messaging`) | `ca-signing.jks` from SealedSecret `tak-enrollment-signer`, projected at `/certs-configmap/` | `TAKServerCAConfig keystorePass` in `core-config` ConfigMap, copied to an emptyDir by the `copy-coreconfig` init container | **Only at pod start** | **Done in step 2** |
| TAK Portal (`onetak/tak-portal`) | `/certs/ca.key` from Secret `tak-portal-ca-certs` (was an **unencrypted** PEM, `CA_KEY_PASS` empty) | env `CA_KEY_PASS` | Pod restart (not EUD-facing) | **Prepared in this PR**: Secret `tak-portal-ca-certs-v2` holds the key encrypted, `CA_KEY_PASS` comes from the same Secret |
| `cert-migration` ConfigMap (`ca-do-not-share.key`, `root-ca-do-not-share.key`, `config-takserver.cfg`) | Hand-managed ConfigMap, mounted at `/certs-configmap/` in api/config/plugins/messaging | TAK cert tooling (`makeCert.sh`, `revokeCert.sh`) run by an operator inside a pod; TAK Server itself does not read these keys (`CAkey=/certs-configmap/ca-plain.key` is absent) | Operator action | **Done in step 2** (see below) |
| `platform/onetak/provider-patch.yaml` (`CA_PASSPHRASE`) | `/home/etl/tak-ca/ca.key`, not provisioned | Hard-coded in a ConfigMap script | None: the ConfigMap is not mounted by any workload | TODO comment added; fix when/if revived |
| `platform/cloudtak` | `cloudtak-tak-ca` holds only `ca.pem` | n/a | n/a | No CA key, nothing to rotate |

Finding for step 2: the CA private keys sit in the `cert-migration` **ConfigMap** (encrypted
with the stock password). Moving them to a Secret is part of step 2.

## Done in this PR

1. Backup (cluster only, never committed): Secret `tak-ca-backup-20261029` in `takserver`
   holds the previous `ca-signing.jks`, both `*-do-not-share.key` files, `ca.pem`, and the
   portal's previous `ca.key`/`ca.pem`. Delete it after step 2 is verified.
2. New password generated in a pipe and never displayed. The enrollment JKS was re-encrypted
   (`keytool -keypasswd` and `-storepasswd`, alias `tak-ca`) and the portal PEM re-encrypted
   (PKCS#8, AES-256-CBC), both from a private `mktemp` directory that was removed.
3. SealedSecrets (namespace-wide):
   - `platform/takserver/deploy/tak-ca-credentials.sealedsecret.yaml`: `ca-signing.jks`
     (new password) and `ca-key-password`. Not mounted by anything yet.
   - `platform/onetak/tak-portal-ca-v2.sealedsecret.yaml`: encrypted `ca.key`, `ca.pem`,
     `ca-key-password`.
4. `platform/onetak/tak-portal.yaml` reads the key and password from `tak-portal-ca-certs-v2`
   (a new name, so a rollout cannot mount the old plaintext key together with a password).
5. Plaintext scrub of docs and TODO markers on runtime files that wait for step 2.

CA fingerprint is unchanged: `72:93:0A:5B:...:3A:15` (SHA-256).

## Verification (2026-10-29)

- New password opens the JKS; the stock password is rejected; JKS certificate fingerprint
  unchanged; private key recovered with the new key password and its public key equals `ca.pem`.
- Portal encrypted key public key equals the portal `ca.pem`.
- Both SealedSecrets applied to the cluster unsealed with the expected keys.
- A throwaway pod from the pinned portal image mounted `tak-portal-ca-certs-v2` and ran the
  portal's own `_issue_cert()`: a client certificate was signed with the CA (SHA-256
  `72930a5b...fb33a15`). The pod was deleted.
- TAK Server was untouched: messaging pod UID and restart count unchanged; established
  connections on 8089 (counted from `/proc/net/tcp` in the messaging container) unchanged.

## Step 2 (2026-10-29): TAK Server side, server stores and CA keys

Owner-approved window ("Do it now"). Everything below is in the step-2 PR; the live rollout is
Argo-driven (no `helm upgrade`, the release stays at revision 34).

### Consumers found

What actually loads the server stores is the `cert-migration` ConfigMap mounted at `/certs-configmap/`
(live server certificate SHA-256 `B7:2E:4F:01:...:D3:61` on 8089, 8443, 8444 and 8446, identical to the
`takserver.jks` in that ConfigMap). The Secret `tak-keystores` is mounted by api, config and messaging but
is **not** read: its `cp /keystores/* /certs/files/` targets a directory that does not exist in the image.

| Consumer | Uses | Step 2 |
|---|---|---|
| `takserver-config`, `-messaging` | `/certs-configmap/*.jks` + CoreConfig passwords | **Switched** to `tak-server-stores` and the substituted CoreConfig |
| `takserver-api`, `-plugins` | `/certs-configmap/*.jks`; passwords come from the config service, no local CoreConfig | **Switched** to `tak-server-stores` (restarted after config) |
| FedHub (`fedhub-certs` SealedSecret) | its own copies, `keystorePassword: atakatak` in `fedhub-config.yaml` | Unchanged, own files |
| TAK-ML (`tak-keystores-copy` in ns takml) | own client store copy, `atakatak` in `configmap-server.yaml` | Unchanged, own files |
| CloudTAK (`platform/cloudtak`) | `cloudtak-tak-ca` holds `ca.pem` only; admin client `.p12` entered by the operator | Unchanged |
| TAK Portal | `CERT_PASS` for the `.p12` it builds, independent of server stores (`services/tak-portal/main.py`, `_issue_cert()` and the package manifest) | Unchanged; its CA key password was step 1 |
| `services/takserver/provision-channels.sh` | `admin.p12` in `cert-migration` with `atakatak` | Unchanged (client file, still mounted) |
| `platform/onetak/provider-patch.yaml` | dormant ConfigMap, not mounted | Unchanged |
| `tak-keystores` Secret (api/config/messaging mounts) | not read by the server | Left in place, retire in a follow-up |

### What changed

1. `core-config-v2` (Argo, `deploy/configmap-core-config-v2.yaml`) is the live CoreConfig with three
   placeholders instead of passwords. The stale 5.7 `deploy/CoreConfig.xml` is removed from Git. The LDAP
   bind credential, which sat in the old `core-config` ConfigMap in plaintext, is now a placeholder fed from
   `tak-ldap-creds`.
2. Init container `copy-coreconfig` (config and messaging) substitutes the placeholders from
   `tak-ca-credentials`, `tak-server-stores` and `tak-ldap-creds` and fails on a leftover placeholder.
3. SealedSecret `tak-server-stores`: `takserver.jks`, `truststore-root.jks`, `fed-truststore.jks`
   re-encrypted with a new random 32-character password (same certificates; legacy PKCS12 encoding kept),
   plus `keystore-password`.
4. Messaging takes `ca-signing.jks` from `tak-ca-credentials` (new CA password) instead of
   `tak-enrollment-signer`.
5. SealedSecret `tak-ca-keys` holds `ca-do-not-share.key` and `root-ca-do-not-share.key` re-encrypted with the
   CA key password (public keys checked against `ca.pem`). Nothing mounts it. After the rollout the keys,
   the old server stores and `*-do-not-share.key` are removed from the live `cert-migration` ConfigMap.
6. `/certs-configmap/` in the four TAK pods is a projected volume (allowlisted `cert-migration` keys plus the
   two Secrets). Sync waves order the roll: config 0, messaging 1, api 2, plugins 3.

Leaf private keys remain in `cert-migration` (`takserver.key`, `admin.key`, `user.key` are **unencrypted**
PEM; `takserver.p12`, `admin.p12`, `user.p12`, `admin.jks`, `user.jks`, `truststore-root.p12` still use the
stock password). Rotating those means re-issuing or re-wrapping the server/admin identities and is left as
step 3.

### Rollout record (2026-10-29, UTC)

1. 02:41 PR #577 merged. Argo rolled config (ok) then messaging, which crash-looped:
   `JwtUtils.loadKeyPair` "keystore password was incorrect". Cause: the config service restores the CoreConfig
   held in Ignite memory and rewrote the new file back to the old passwords (see the runbook).
2. 02:46 rollback by `kubectl rollout undo` of config and messaging with Argo auto-sync paused (Argo,
   including app-of-apps selfHeal, re-applied the new specs twice until both were paused). Old pods Ready
   about 02:50; fingerprints and connections matched the preflight. Outage of 8089: about 8 minutes.
3. 02:52 full bounce: api, plugins, messaging, config scaled to 0, `takserver-ignite-0` deleted and Ready,
   Argo auto-sync restored. config, messaging, api and plugins came up on the new specs in that order;
   all Ready by 02:59 (second outage about 7 minutes).

### Verification (after the second roll)

- All TAK Server pods Ready, 0 restarts; config pod `CoreConfig.xml` holds the new passwords.
- Server certificate SHA-256 on 8089, 8443, 8444, 8446 identical to the preflight
  (`B7:2E:4F:01:...:D3:61`). `GET /Marti/api/version` with the admin client certificate returned 200.
- Established connections on 8089: 1 before, 1 after (the one peer is the hermes-legacy pod, 10.109.3.95, which dials 8089 every 30 s
  without a client certificate and logs `PEER_DID_NOT_RETURN_A_CERTIFICATE` on the server; unrelated to the
  keystores, and not an EUD, so the EUD count could not be separately confirmed).
- Plugin manager: `pluginCoreConfigApi available: true`. FedHub, tak-ldap, the NATS leaf, CloudTAK and
  TAK-ML: no keystore/SSL errors, no restarts caused by the change. The LDAP `ou=Group` lookup error 32 is
  older than this change (first seen 2026-10-29 22:01 UTC); binds as `cn=ldapadm` succeed.
- Not tested end to end: a certificate signing request through TAK Server (the admin client certificate
  gets 403 on `/Marti/api/tls/*`, and no user password was used). The new CA password was verified to open
  `ca-signing.jks` in the running messaging image with the exact Secret value, and the messaging pod started
  with the substituted `keystorePass`. Run one CloudTAK enrollment as in
  `tak-enrollment-repair-2026-10-25.md` at the next opportunity.
- No lasting objects were created by the checks.

### Cleanup done after verification

- Live `cert-migration` ConfigMap: removed `ca-do-not-share.key`, `root-ca-do-not-share.key`,
  `takserver.jks`, `truststore-root.jks`, `fed-truststore.jks` (checked byte for byte against
  `tak-ca-backup-20261029` and `tak-step2-backup-20261029` first). The projected volumes never mounted them.
- Live `core-config` ConfigMap deleted (it held the stock passwords and the LDAP bind password); its content is
  in `tak-step2-backup-20261029`. A live-only `core-config` mount on `takserver-plugins` (not in Git, unused)
  was removed.
- SealedSecret `tak-enrollment-signer` retired from Git (Argo prunes the Secret); nothing mounts it.
- Reverting #577 alone is no longer a complete rollback: restore `core-config` and the removed
  `cert-migration` keys from `tak-step2-backup-20261029` first, and bounce Ignite as described in the runbook.
- Step 3 is prepared (PR `chore/tak-ca-password-rotation-step3`, runbook section "STEP 3" in
  `platform/takserver/DEPLOY-RUNBOOK.md`, script `scripts/tak/rotate-step3.sh`); the items below that it covers are
  done when that window is.
- Remaining: delete `tak-step2-backup-20261029` and `tak-ca-backup-20261029` after the soak period;
  retire the unused `tak-keystores` Secret and its mounts; step 3 (unencrypted `takserver.key`, `admin.key`,
  `user.key` and the stock-password `*.p12`/`admin.jks`/`user.jks` in `cert-migration`, plus FedHub, TAK-ML
  and portal device passwords).

### Rollback

Cluster backup: Secret `tak-step2-backup-20261029` (`takserver`): `tak-keystores`,
`tak-enrollment-signer`, the old `core-config` and `cert-migration` contents, and the previous Deployment
specs. Before the cleanup above: revert the step-2 PR (Argo re-rolls config, messaging, api, plugins to the old mounts) and bounce Ignite so the cluster drops the new cached CoreConfig.
Delete `tak-step2-backup-20261029` and `tak-ca-backup-20261029` once the follow-up is done.

## Reading the new password

The owner can read it with:

```bash
kubectl --context do-nyc3-omos-cloud -n takserver get secret tak-ca-credentials -o jsonpath='{.data.ca-key-password}' | base64 -d
```

(The Secret exists in the cluster now because the SealedSecret was applied once for
verification; after merge Argo owns it.)

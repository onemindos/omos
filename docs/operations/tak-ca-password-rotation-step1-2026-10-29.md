# TAK CA signing-key password rotation, step 1 (2026-10-29)

Owner-approved "step 1": rotate only the password that protects the **CA signing key**.
Device-facing passwords (8089/8443/8446 `keystorePass`/`truststorePass`, client `.p12`,
truststore handed to devices) stay at the stock value until **step 2**. No EUD-facing
workload was restarted and no connection changed.

## CA key password consumers (verified 2026-10-29)

| Consumer | Where the CA key lives | How the password is read | Pickup | Status |
|---|---|---|---|---|
| TAK Server enrollment v2 (`takserver-messaging`) | `ca-signing.jks` from SealedSecret `tak-enrollment-signer`, projected at `/certs-configmap/` | `TAKServerCAConfig keystorePass` in `core-config` ConfigMap, copied to an emptyDir by the `copy-coreconfig` init container | **Only at pod start** | **Deferred to step 2** |
| TAK Portal (`onetak/tak-portal`) | `/certs/ca.key` from Secret `tak-portal-ca-certs` (was an **unencrypted** PEM, `CA_KEY_PASS` empty) | env `CA_KEY_PASS` | Pod restart (not EUD-facing) | **Prepared in this PR**: Secret `tak-portal-ca-certs-v2` holds the key encrypted, `CA_KEY_PASS` comes from the same Secret |
| `cert-migration` ConfigMap (`ca-do-not-share.key`, `root-ca-do-not-share.key`, `config-takserver.cfg`) | Hand-managed ConfigMap, mounted at `/certs-configmap/` in api/config/plugins/messaging | TAK cert tooling (`makeCert.sh`, `revokeCert.sh`) run by an operator inside a pod; TAK Server itself does not read these keys (`CAkey=/certs-configmap/ca-plain.key` is absent) | Operator action | **Deferred to step 2** (re-encrypt keys, update cfg, then move out of a ConfigMap) |
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

## Why the TAK Server half is deferred

`takserver-messaging` reads `keystorePass` from CoreConfig only when its pod starts, and
`ca-signing.jks` is projected into the same pod. Switching the JKS without changing
`keystorePass` in the same restart breaks enrollment signing; switching only CoreConfig has
the same effect. Both change together, so it needs a messaging restart (Recreate strategy:
a brief reconnect for EUDs). CoreConfig also cannot read a password from the environment,
so the init container must substitute it from the Secret. Do this in the step-2 window.

## Step-2 apply (do not run before the maintenance window)

1. Merge this PR (Argo creates both SealedSecrets and rolls the portal).
2. In `takserver-messaging-deployment.yaml`: project `tak-ca-credentials` (`ca-signing.jks`,
   mode 0400) instead of `tak-enrollment-signer`; add env `CA_KEY_PASS` from
   `tak-ca-credentials/ca-key-password`; change `copy-coreconfig` to a shell that writes
   `CoreConfig.xml` with the placeholder `__CA_KEY_PASS__` replaced from `$CA_KEY_PASS`
   (use a delimiter-safe replacement; the password is alphanumeric); put the placeholder in
   the `TAKServerCAConfig keystorePass` of `core-config`.
3. Sync during the window, wait for messaging Ready, verify CloudTAK enrollment as in
   `tak-enrollment-repair-2026-10-25.md` and that EUD sessions on 8089 return.
4. Rollback: revert the Deployment and `core-config` to the previous revision; the previous
   JKS is still in `tak-enrollment-signer` and `tak-ca-backup-20261029`.
5. Then retire `tak-enrollment-signer`, rotate the device-facing passwords, and move the
   `cert-migration` keys into a Secret.

## Reading the new password

The owner can read it with:

```bash
kubectl --context do-nyc3-omos-cloud -n takserver get secret tak-ca-credentials -o jsonpath='{.data.ca-key-password}' | base64 -d
```

(The Secret exists in the cluster now because the SealedSecret was applied once for
verification; after merge Argo owns it.)

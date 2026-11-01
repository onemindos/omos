# Private IoT pilot credential preparation

No actual credentials are included. Both Applications stay manual-sync and their
`sealed/` directories are intentionally empty. Missing SealedSecrets block Pod
startup; do not create out-of-band live Secrets to bypass that boundary.

Use a private JSON file outside the repository, owner-only mode0600. Preserve it
in the owner's encrypted credential/recovery store. Install `kubeseal`, Python3
with `PyYAML` and `bcrypt` locally for the ThingsBoard encoder; use the existing controller's
public certificate. Do not generate a new sealing controller or rotate its key.

Input is a mapping of the exact Secret names below to mappings of exact key names.
No values or sample credentials are supplied in this guide:

| App | Secret | Required input keys |
|---|---|---|
| OpenRemote | `openremote-env` | `DB_PASSWORD`, `ADMIN_PASSWORD` |
| OpenRemote | `openremote-database-restic` | `RESTIC_REPOSITORY`, `RESTIC_PASSWORD`, `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY` |
| OpenRemote | `openremote-storage-restic` | same four Restic keys |
| ThingsBoard | `thingsboard-env` | `DB_PASSWORD`, `SYSADMIN_EMAIL`, `SYSADMIN_PASSWORD`, `JWT_TOKEN_SIGNING_KEY` |
| ThingsBoard | `thingsboard-database-restic` | same four Restic keys |

Each password must be random, at least24 characters. BCrypt accepts a maximum72
UTF8 bytes for the ThingsBoard owner password. The JWT key must be Base64 encoding
of at least32 random bytes. Each volume needs a distinct HTTPS S3 repository prefix
and its own Restic password; choose backup credentials scoped to those prefixes.
Do not point both volumes at one Restic repository. Do not reuse the app's admin
credential as an API identity or backup encryption key.

Run one app at a time, using your actual private paths:

```sh
python scripts/prepare_iot_secrets.py openremote --private-json "$PRIVATE_IOT_JSON" --cert "$SEALED_SECRETS_CERT"
python scripts/prepare_iot_secrets.py thingsboard --private-json "$PRIVATE_TB_JSON" --cert "$SEALED_SECRETS_CERT"
```

The helper validates input, hashes the ThingsBoard owner password and uses strict
namespace/name-bound sealing. It writes only ciphertext and updates the app's
`sealed/kustomization.yaml`; subprocess diagnostics are suppressed on failure.
It never applies anything to the cluster and refuses to overwrite existing files.
Inspect the encrypted diff, run validation and commit through normal review before
manual Argo sync. Preserve passwords, Restic repositories/passwords, original
app secrets and sealing-key recovery independently of the cluster. Credential
rotation and retained-data adoption need their own reviewed procedure.

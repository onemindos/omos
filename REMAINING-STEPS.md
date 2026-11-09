# Remaining owner steps: offsite backup copy (`ops/backup-offsite-copy`)

Branch: `fix/offsite-copy-credentials-wiring` — every step that can be completed in
Git **without real secret values** is done here. What remains requires the real
destination credentials, which only the owner may produce (never fabricated, never
committed in plaintext).

## Already done on this branch (in Git)

| Step (OFFSITE-COPY.md) | State |
|---|---|
| CronJob, NetworkPolicy, ConfigMap, job script | Committed earlier; script is CI-tested (real rclone engine) |
| `offsite-copy-config` destination defaults | Done: DigitalOcean Spaces bucket `onemind-backups-sfo3`, region `sfo3`, mirror prefix `mirror/` |
| Kustomization wiring for the credentials SealedSecret | Done on this branch: `sealed/ops-offsite-copy-credentials.yaml` is listed in `platform/backups/kustomization.yaml` |
| `sealed/ops-offsite-copy-credentials.yaml` | Present as a blatantly-marked **placeholder** SealedSecret. Its marker-text values cannot be unsealed by the controller (no Secret is ever created from it), cannot be mistaken for real ciphertext (real kubeseal values are ~400+ char base64 starting with `Ag`), and pass CI. Replace it; do not edit it by hand. |
| Alert behaviour documentation | `CriticalBackupScheduleSuspended` already exists; doc notes updated |

Destination decision: **already recorded in Git** (`offsite-copy-config`: sfo3 /
`onemind-backups-sfo3`). Only change it first if you want Backblaze B2 or a home
MinIO instead — see OFFSITE-COPY.md "Moving to Backblaze B2 or a home server".

## What only the owner can do

### 1-3. Destination bucket, scoped key, optional lifecycle (control panel / doctl)

Exactly as documented in `platform/backups/OFFSITE-COPY.md` steps 1-3:

```sh
# step 2: bucket-scoped Spaces key (values shown ONCE — keep them out of chat/Git)
doctl spaces keys create onemind-offsite-copy \
  --grants 'bucket=onemind-backups;permission=read,bucket=onemind-backups-sfo3;permission=readwrite'
```

### 4. Seal the credentials (replaces the placeholder)

Run from the repository root, with `kubeseal` able to reach the cluster:

```sh
umask 077
printf '%s' '<spaces access key id>'     > src-id
printf '%s' '<spaces secret access key>' > src-secret
printf '%s' '<spaces access key id>'     > dst-id      # same pair for DO at both ends
printf '%s' '<spaces secret access key>' > dst-secret

kubeseal --fetch-cert --controller-name=sealed-secrets-controller \
  --controller-namespace=kube-system > /tmp/pub-cert.pem

kubectl create secret generic offsite-copy-credentials -n ops \
  --from-file=RCLONE_CONFIG_SRC_ACCESS_KEY_ID=./src-id \
  --from-file=RCLONE_CONFIG_SRC_SECRET_ACCESS_KEY=./src-secret \
  --from-file=RCLONE_CONFIG_DST_ACCESS_KEY_ID=./dst-id \
  --from-file=RCLONE_CONFIG_DST_SECRET_ACCESS_KEY=./dst-secret \
  --dry-run=client -o json \
  | kubeseal --cert /tmp/pub-cert.pem --format yaml --scope namespace-wide \
  > platform/backups/sealed/ops-offsite-copy-credentials.yaml

rm -f src-id src-secret dst-id dst-secret /tmp/pub-cert.pem
```

Result: the placeholder file is overwritten **at the same path**; every
`encryptedData` value starts with `Ag` and is ~400+ chars. If it doesn't, the
sealing failed — do not commit it.

### 5. Enable (same commit as the sealed file)

In `platform/backups/offsite-copy.yaml`:

- set `spec.suspend: false`
- delete BOTH CronJob metadata annotations:
  `onemind.io/backup-state: pending-setup` and `onemind.io/pending-reason`

Clearing the annotations is what moves the portal (Home / Today / Live Ops) from
"Needs setup" to actively tracking the job, and what enrols it in
`ScheduledBackupOverdue`. No kustomization change is needed — the entry is wired.
Commit, push, merge. Argo applies; the controller unseals the real Secret.

### 6. Seed and verify

```sh
kubectl -n ops get secret offsite-copy-credentials      # must exist after sync
kubectl -n ops create job --from=cronjob/backup-offsite-copy offsite-seed-1
kubectl -n ops logs -f job/offsite-seed-1
```

First run copies everything (an hour or more). Then confirm the heartbeat object
`s3://onemind-backups-sfo3/_offsite/last-success.json` exists, and spot-check with
`restic check` against a mirrored repository and `tools/decrypt.py` on a mirrored
archive (OFFSITE-COPY.md step 6 and "Restoring from the mirror").

## Alert expectations after enabling

- `CriticalBackupScheduleSuspended` (the alert firing now, warning severity):
  clears on its own, but only **~7 days after** `suspend` flips to false — its
  expression requires `min_over_time(kube_cronjob_spec_suspend[7d]) == 1`, so the
  7-day window must fill with unsuspended samples. That lag is expected, not a
  failure of step 5.
- `ScheduledBackupOverdue` takes over immediately after unsuspension (fires if a
  schedule older than 30 h has no successful run).
- `OffsiteBackupCopyFailed` fires 5 minutes after any failed Job.

## Never do

- Never commit the plaintext `kubectl create secret` manifest, the `src-id` /
  `*secret` files, or the Spaces key values — to any branch, document, or chat.
- Never enable the CronJob while the placeholder file is still in place.
- Never delete `mirror/` objects from the destination bucket automatically
  (`rclone copy` semantics are deliberate; see OFFSITE-COPY.md retention policy).

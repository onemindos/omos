# DIAGNOSIS — VolSync sources ntfy-data and home-assistant-config never synced

Date: 2026-11-09. Branch: fix/volsync-unsynced-ntfy-ha.

## Root cause

- `ntfy/ntfy-data`: the ReplicationSource shipped with `spec.paused: true` and
  its restic repository Secret `ntfy-data-restic` was never sealed/applied —
  the `TODO(owner)` comment in `platform/ntfy/backup.yaml` and the follow-up
  note in `platform/ntfy/README.md` documented this as deferred work that was
  never completed.
- `home-assistant/home-assistant-config`: identical condition — `paused: true`
  plus a missing `home-assistant-config-restic` SealedSecret
  (`platform/home-assistant/backup.yaml`, TODO + README note).

A paused source never spawns the restic mover, so `status.lastSyncTime` is
never set and the `PortableVolumeBackupOutOfSync` (role=source) alert fires.
Unpausing alone would still fail every run because the Secret named in
`spec.restic.repository` does not exist. Both defects had to be fixed together.

## Field-by-field diff vs known-good sources (forgejo/forgejo-data,
hermes/onemind-hermes-data, karakeep/*, volsync/sources.yaml set)

| Field | Broken pair | Known-good | Verdict |
|---|---|---|---|
| `spec.paused` | `true` | `false` (all) | ROOT CAUSE 1 |
| restic repo Secret `<pvc>-restic` in namespace | absent | present SealedSecret in kustomization | ROOT CAUSE 2 |
| sourcePVC name | matches PVC `ntfy-data` / `home-assistant-config` (verified in `data-pvc.yaml` / `config-pvc.yaml`) | — | OK |
| `moverServiceAccount` ServiceAccount `platform-volume-backup` | declared in same `backup.yaml` | identical | OK |
| `cacheStorageClassName`, `storageClassName`, `volumeSnapshotClassName` | `do-block-storage` | `do-block-storage` | OK |
| `pruneIntervalDays` 14, retain 7d/4w/3m/1y | identical | identical | OK (untouched) |
| schedule (`50 14 * * *` / `20 14 * * *`) | valid cron, valid trigger | valid | OK (untouched) |
| NetworkPolicy `platform-volume-backup-egress` (DNS + S3 proxy 10.10.255.254:443) | present | present | OK |
| moverSecurityContext | matches app UID (ntfy 1000, HA root) — same practice as per-app good sources | per-app | OK |
| kustomization inclusion of the source | both already referenced (`backup.yaml`) | — | OK (that's why the paused objects exist live) |

Everything else was structurally identical to the good pattern; only the pause
flag and the missing SealedSecret distinguished them.

## What changed (this branch)

- `platform/ntfy/backup.yaml` — `paused: true` → `paused: false`; TODO comment
  replaced by current-state note.
- `platform/ntfy/backup-data.sealedsecret.yaml` — NEW. SealedSecret
  `ntfy-data-restic` (ns ntfy), karakeep pattern: sync-wave `-1`, strict scope,
  keys RESTIC_REPOSITORY, RESTIC_PASSWORD, AWS_ACCESS_KEY_ID,
  AWS_SECRET_ACCESS_KEY, AWS_DEFAULT_REGION. Payloads are PLACEHOLDER base64
  (decode to readable marker text, not credentials).
- `platform/ntfy/kustomization.yaml` — added `backup-data.sealedsecret.yaml`.
- `platform/home-assistant/backup.yaml` — same unpause/comment change.
- `platform/home-assistant/backup-config.sealedsecret.yaml` — NEW. SealedSecret
  `home-assistant-config-restic` (ns home-assistant), same placeholder pattern.
- `platform/home-assistant/kustomization.yaml` — added
  `backup-config.sealedsecret.yaml`.
- `platform/ntfy/README.md`, `platform/home-assistant/README.md` — status text
  updated to match reality.

No changes to the other 24 sources, schedules, retention, or prune semantics.

## REQUIRED OWNER ACTION — reseal the placeholders

The two new SealedSecrets contain placeholder payloads. The sealed-secrets
controller CANNOT decrypt them (it will log an error and no Secret will
materialize) until every value is replaced with real kubeseal output. This is
deliberate: no credentials may be committed to this repo, even sealed by
anyone other than the key holder.

For each of the two (substitute NAME, NAMESPACE, file path accordingly —
`ntfy-data-restic`/ntfy/`backup-data.sealedsecret.yaml` and
`home-assistant-config-restic`/home-assistant/`backup-config.sealedsecret.yaml`):

1. Gather the plaintext values WITHOUT writing them to the repo:
   - S3 endpoint/bucket and access/secret key + region: identical to the
     existing volume backups. Read a live good secret, e.g.
     `kubectl -n forgejo get secret forgejo-data-restic -o jsonpath='{.data.RESTIC_REPOSITORY}' | base64 -d`
     (read-only; keep output out of the repo).
   - RESTIC_REPOSITORY: same scheme/host/bucket as the good secret, new
     per-source path prefix (`ntfy-data` / `home-assistant-config`) — every
     source owns its repository prefix (see platform/backups/volsync/README.md).
   - RESTIC_PASSWORD: the platform restic password preserved at
     `~/.config/onemind/recovery/restic-password` (volsync README). Confirm it
     matches a good secret's password before sealing:
     `diff <(sudo cat ~/.config/onemind/recovery/restic-password) <(kubectl -n forgejo get secret forgejo-data-restic -o jsonpath='{.data.RESTIC_PASSWORD}' | base64 -d)`
2. Seal:
   ```sh
   kubectl create secret generic ntfy-data-restic -n ntfy \
     --from-literal=RESTIC_REPOSITORY="$REPO_URL" \
     --from-literal=RESTIC_PASSWORD="$RESTIC_PASSWORD" \
     --from-literal=AWS_ACCESS_KEY_ID="$AWS_ACCESS_KEY_ID" \
     --from-literal=AWS_SECRET_ACCESS_KEY="$AWS_SECRET_ACCESS_KEY" \
     --from-literal=AWS_DEFAULT_REGION="$AWS_DEFAULT_REGION" \
     --dry-run=client -o yaml | kubeseal --format yaml > /tmp/sealed.yaml
   ```
   then copy ONLY `spec.encryptedData` from /tmp/sealed.yaml into
   `platform/ntfy/backup-data.sealedsecret.yaml`, overwriting the five
   placeholder lines (keep metadata, sync-wave annotation and template as-is).
   Repeat for home-assistant. Shred /tmp/sealed.yaml. Never paste values into
   chat or docs.
3. Commit the resealed files on this branch (or a follow-up branch) and let
   Argo CD sync. Do NOT `kubectl apply` the secrets — GitOps only.

## Post-Argo-sync verification (after reseal + sync)

```sh
# 1. SealedSecrets decrypted into Secrets (both must exist):
kubectl -n ntfy get secret ntfy-data-restic
kubectl -n home-assistant get secret home-assistant-config-restic

# 2. Sources are unpaused and scheduling:
kubectl -n ntfy get replicationsource ntfy-data -o jsonpath='{.spec.paused}{"\n"}{.status.nextSyncTime}{"\n"}'
kubectl -n home-assistant get replicationsource home-assistant-config -o jsonpath='{.spec.paused}{"\n"}{.status.nextSyncTime}{"\n"}'

# 3. After the first scheduled window (ntfy 14:50Z, home-assistant 14:20Z),
#    lastSyncTime must be set and mover Jobs completed:
kubectl -n ntfy get replicationsource ntfy-data -o jsonpath='{.status.lastSyncTime}'
kubectl -n home-assistant get replicationsource home-assistant-config -o jsonpath='{.status.lastSyncTime}'
kubectl -n ntfy get jobs
kubectl -n home-assistant get jobs
kubectl -n ntfy describe replicationsource ntfy-data | tail -20

# 4. Alert clears in Prometheus/Alertmanager (PortableVolumeBackupOutOfSync,
#    role=source, names ntfy-data / home-assistant-config).
```

To force an immediate first sync instead of waiting for the window (optional,
read-mostly but creates one Job):
`kubectl -n ntfy patch replicationsource ntfy-data --type merge -p '{"spec":{"trigger":{"manual":"init-2026-11-09"}}}'`
— then revert the trigger to the schedule afterwards via Git.

## Rollback

`git revert` this branch's commit (or reapply `paused: true` on both sources).
Rollback removes the two SealedSecret CRs from Argo's desired state; the
sealed-secrets controller deletes the managed Secret when its SealedSecret is
removed. No restic repository data is ever deleted by VolSync when a
ReplicationSource changes or disappears, so pausing/unpausing is data-safe.
The other 24 sources are untouched by this change and unaffected by rollback.

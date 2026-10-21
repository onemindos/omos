# Hermes restore drill, 2026-10-16

**Result: passed.** The latest nightly VolSync backup of `hermes-data` was restored to a new volume and a copy of Hermes booted against it at application level. The live volume and the running agent were not touched. This closes the gap left by the 2026-10-08 drill, which only verified files.

## What was done

1. **Restore.** A new 25 GiB PVC `hermes-restic-drill-20261016` and a `ReplicationDestination` using the same mover settings as the 2026-10-08 drill (`copyMethod: Direct`, repository `volume-c96e219d2314-restic`, ServiceAccount `platform-volume-backup`). The latest snapshot (taken 01:02 UTC) came back in about 4 minutes: 11.4 GiB, 368,391 files and directories.
2. **Application-level check.** A Job with a deny-all network policy (no network at all, so nothing could reach NATS, Telegram or Matrix) mounted only the restored copy, set `HERMES_HOME=/data` and ran Hermes' own commands. It never started the gateway.
3. **Cleanup.** The Job, network policy, `ReplicationDestination` and PVC were deleted by name. `hermes-data` and the live pod were checked afterwards and are unchanged.

## What the restored copy showed

| Check | Result |
|-------|--------|
| `hermes --version` | v0.21.5 (2026.9.24), same as the live agent |
| `hermes config path` / `config check` / `config get model.provider` | config found, loads, provider `digitalocean` |
| `SOUL.md` checksum | identical to the live soul at the time |
| `cron/jobs.json` checksum | identical to live |
| Skills on disk | 166, same as live |
| Skills as Hermes sees them (`hermes skills list`) | 170 (2 hub, 40 builtin, 128 local) |
| Profiles | 14, same as live |
| `state.db` | 1,426 MB, `PRAGMA quick_check` ok, 42 tables |
| `state_compact.db` | see below |

## Finding: `state_compact.db` is dead weight

The live `state_compact.db` is 929 MiB, last written 2026-09-26, with a 1 KiB `-journal` beside it. The journal header records an original size of 0 pages and 0 records, and the main file's first page is not a valid SQLite header ("file is not a database"). This is the leftover of a compaction that never finished. Opening it normally lets SQLite roll the journal back, which truncated the restored copy to an empty file in the drill. Nothing in the Hermes code references `state_compact`, and no process has it open. It adds 929 MiB to every nightly backup.

It was **not** touched on the live volume (it was inspected read-only with `immutable=1`). Removing it needs Zeus's explicit go-ahead.

## Other things the drill turned up

- Hermes prints on every command: "1Password: secrets.onepassword.enabled is true but the op CLI was not found on PATH". The `op` binary exists on the volume (`/opt/data/home/.local/bin/op`) but is not on the default PATH, so the `op://` references in `config.yaml` do not resolve. Moving credentials to `op://` references, as the credential cleanup plans, needs `op` on PATH and a service-account token first.

## Manifests used (for the record, not applied by Argo)

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata: {name: hermes-restic-drill-20261016, namespace: hermes, labels: {onemind.io/restore-test: drill-20261016}}
spec: {accessModes: [ReadWriteOnce], storageClassName: do-block-storage, resources: {requests: {storage: 25Gi}}}
---
apiVersion: volsync.backube/v1alpha1
kind: ReplicationDestination
metadata: {name: hermes-restic-drill-20261016, namespace: hermes, labels: {onemind.io/restore-test: drill-20261016}}
spec:
  trigger: {manual: drill-20261016}
  restic:
    repository: volume-c96e219d2314-restic
    copyMethod: Direct
    storageClassName: do-block-storage
    cacheStorageClassName: do-block-storage
    cacheCapacity: 1Gi
    moverServiceAccount: platform-volume-backup
    moverPodLabels: {app: platform-volume-backup}
    moverSecurityContext: {runAsUser: 0, runAsGroup: 0, seccompProfile: {type: RuntimeDefault}}
    destinationPVC: hermes-restic-drill-20261016
    cleanupCachePVC: true
```

The verify Job used the live image `nousresearch/hermes-agent:v2026.9.24` with `HERMES_HOME=/data`, `HERMES_DISABLE_LAZY_INSTALLS=1`, the restored PVC mounted at `/data`, and a `NetworkPolicy` selecting the Job's label with both policy types and no rules (deny all).

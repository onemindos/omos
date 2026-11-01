# Paperclip retirement — 2026-11-01

Owner decision (2026-11-01): "Remove Paperclip entirely", after a backup.

## Final backups (taken 2026-11-01, before this change)

| What | Evidence |
|---|---|
| Logical dump | `pg_dump -Fc` of database `paperclip` (PostgreSQL 17.11), 1,322,573 bytes, 242 tables, `pg_restore --list` verified. SHA-256 `3a56afce0e90cf0a70adbef91c4118fe4373b5ecc1f85032aac5cf60a35b662c`. Held locally by the operator, mode 0600; never committed or uploaded (may contain user data). |
| Volume snapshot, database PV | DO snapshot `d17ca36c-d5b6-11f1-b952-52ac49c6ab7a` (`paperclip-postgres-final-2026-11-01`) of volume `b7e8e08f-cf83-11f1-a7b7-1ea8959b39f9` (PV `pvc-71031808-f86e-4635-b54a-ee50f8a46d1b`) |
| Volume snapshot, app PV | DO snapshot `d1c9b17b-d5b6-11f1-b3d4-5a83ceea2d76` (`paperclip-data-final-2026-11-01`) of volume `b7954c14-cf83-11f1-abe9-26fe4a441cbc` (PV `pvc-ae84233c-1653-45c8-8ab5-18f853664efc`) |
| Existing restic repositories | VolSync (`paperclip-data`, `paperclip-postgres`) and the encrypted PG dump job write to the platform backup bucket under their own prefixes. They are not deleted by this change. |

Snapshots are tagged `paperclip-retirement`. The snapshots report 0.00 GiB at creation
because the DO API shows size lazily.

## What merging the removal PR does

The Application `platform-paperclip` has `prune: false` and **no** `resources-finalizer`. The
parent `app-of-apps` has `prune: true`, so merging deletes only the Application CR. Because
there is no cascading finalizer, **Argo CD does not delete any Paperclip workload, PVC,
Secret or the namespace**. Everything keeps running, unmanaged, until the live teardown below.
This is the safe behaviour: nothing is destroyed by a merge.

Both data PVCs already use `do-block-storage-retain` and their PVs are already
`persistentVolumeReclaimPolicy: Retain` (verified 2026-11-01), so deleting the claims or the
namespace leaves the DO volumes in place (PVs go `Released`). The two VolSync cache PVCs are
`do-block-storage` (reclaim `Delete`) and are intentionally disposable.

## Pre-merge (coordinator, optional because already Retain)

Idempotent re-assertion of the reclaim policy; run it if there is any doubt:

```bash
kubectl --context do-nyc3-omos-cloud patch pv pvc-71031808-f86e-4635-b54a-ee50f8a46d1b \
  -p '{"spec":{"persistentVolumeReclaimPolicy":"Retain"}}'
kubectl --context do-nyc3-omos-cloud patch pv pvc-ae84233c-1653-45c8-8ab5-18f853664efc \
  -p '{"spec":{"persistentVolumeReclaimPolicy":"Retain"}}'
```

## Live teardown after merge (owner-approved; not part of the PR)

Stop traffic and writers first, then remove the namespace. With `selfHeal` gone (the
Application no longer exists) nothing recreates the objects.

```bash
K="kubectl --context do-nyc3-omos-cloud"
$K -n paperclip delete ingress paperclip                       # drops paperclip.tail717ba0.ts.net
$K -n paperclip delete cronjob backup-postgres-paperclip
$K -n paperclip delete replicationsource paperclip-data paperclip-postgres
$K -n paperclip scale deploy/paperclip --replicas=0
$K delete namespace paperclip                                  # deletes PVCs; PVs stay Released (Retain)
$K get pv | grep paperclip                                     # expect two Released PVs
```

Afterwards, remove the stale Tailscale device `paperclip` from the tailnet admin console if
the operator did not clean it up.

## Deferred, needs separate owner approval

- Deleting the two Released PVs and their DO volumes (`b7e8e08f-...`, `b7954c14-...`).
- Deleting the DO snapshots `paperclip-*-final-2026-11-01` and the restic repositories.
- Hermes-side leftovers (`platform/hermes/agent/paperclip-api*`, `paperclip-agent`
  SealedSecret, network-policy rules, the disabled `paperclip` MCP adapter entry in
  `platform/hermes/mcp/apps.yaml`) are retired with the old Hermes, not here.
- `docs/PLATFORM-EXTRACTION.md`, `docs/architecture/onemind-tool-gateway.md` and the archived
  registry page `docs/inventory/fabric/` still mention Paperclip as historical or planned
  content and are left as written.

## Restore, if ever needed

Recreate a PostgreSQL 17 instance and run `pg_restore -d <db> <dump>` from the operator's
copy, or create a volume from the DO snapshot. The removed manifests are in Git history
before this change (`platform/paperclip/`).

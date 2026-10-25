# Vaultwarden recovery

## What is protected

Vaultwarden runs one SQLite writer on `data-pvc`. At **13:20 UTC**, the
`vaultwarden-sqlite-backup` CronJob runs the pinned server image's native
`/vaultwarden backup` command as UID/GID 1000 on the app's node. Vaultwarden 1.37.4
opens the live database read-only and uses SQLite `VACUUM INTO` to finish a separate
`/data/db_YYYYMMDD_HHMMSS.sqlite3` file. The job checks the result, keeps the latest
seven strictly named native copies, and atomically replaces
`/data/.sqlite-backup-complete.json` with the completed filename, SHA256 and UTC
time. It never removes `db.sqlite3`, WAL files, attachments or arbitrary `db_*`
files. Native plaintext copies remain inside the protected source claim; no
plaintext backup staging PVC, upload directory or external bucket is created.

At **13:40 UTC**, the `vaultwarden-data` VolSync source takes a DigitalOcean
whole-volume snapshot and uploads it through Restic to the dedicated
`onemind-backups/backups/restic/vaultwarden/data-pvc` Spaces prefix. Restic encrypts
contents and filenames with a fresh password unique to this repository. The
strict namespace/name-bound `vaultwarden-data-restic` SealedSecret reuses only
the established Spaces storage credentials. No account or vault password is used.
The native job has no Kubernetes API token, workload permissions or network
access. Only the VolSync mover has the existing DNS and HTTPS backup egress.

The full archive includes the completed native SQLite file, `attachments/`,
`sends/`, signing keys (`rsa_key*`), any `config.json`, the completion marker, and
other files on `/data`. Keys and configuration are sensitive even when no vault
account has been created. Restic retention is 7 daily, 4 weekly, 3 monthly and 1 yearly;
pruning runs at 14 day intervals. The source claim, existing repositories, storage
credentials and SealedSecrets controller recovery keys must be retained.

## Consistency and failure limits

The **completed native SQLite copy** is a consistent standalone database. The
whole-volume snapshot also includes the live `db.sqlite3` and its journals/WAL,
which are only crash-consistent. Use the native file selected by the archived
marker for the verified SQLite recovery path. Do not combine it with the live
database's WAL or SHM files. Attachments and sends can change between the native
SQLite copy and the volume snapshot; this is not an atomic application-wide
backup. Pause writes for a coordinated final migration if that is required.

The two schedules are independent. A failed native job leaves its previous
completion marker intact; VolSync can still successfully archive that older
copy. A healthy/successful ReplicationSource alone does not prove a fresh or
nonempty application backup. Check both Job success and the archived marker's
freshness, then test recovery. The 20 minute interval exceeds the native Job's
10 minute deadline. Job overlap is forbidden, retries and runtime are bounded,
and required pod affinity avoids a second-node RWO attach.

## First archive acceptance

Run these checks after reviewed Git/Argo deployment; preserve all production files
and existing archives. The initial results are recorded below. Repeat the checks
when validating later changes or recovery points.

1. Require the application Ready and SQLite initialized. Run one instance of the
   native backup job through the reviewed GitOps rollout, then require that Job
   Successful. Read only its filename/hash/time marker; require the marker to be
   from this rollout and no more than 30 minutes old. Do not create a fake user or
   reset credentials to populate a test database.
2. Confirm the repository SealedSecret Ready and VolSync source present. A source
   may run on first reconciliation before the first native backup. Ignore that
   early archive for acceptance. Add a unique `spec.trigger.manual` tag through
   Git/Argo **after** the fresh native marker, then wait for the source's
   `status.lastManualSync` to equal the tag and its successful completion time to
   follow the native job. Require an actual nonempty Restic snapshot containing
   the marker and its selected native SQLite file.
3. Carry out the isolated archive verification below. Record snapshot/time,
   marker/time, file/hash match, SQLite integrity and schema counts, without
   printing contents or secret values.
4. Remove `spec.trigger.manual` through Git/Argo once accepted. In the installed
   VolSync 0.16 controller, a nonempty manual tag takes priority over the cron
   schedule; leaving it in Git prevents subsequent scheduled runs. Confirm the
   daily source `nextSyncTime` returns. Preserve the 13:40 schedule and repository.

## Isolated encrypted archive verification

Use an official pinned Restic binary verified against its published checksums.
Read the existing `vaultwarden/vaultwarden-data-restic` Secret into the verifier
process's memory and pass values only through the Restic subprocess environment.
Do not print values, put them in command arguments, or write a plaintext env file.
Use `--no-lock snapshots --json`, `--no-lock ls --json SNAPSHOT`, and
`--no-lock dump SNAPSHOT PATH` only. These read the repository without remote
locks, pruning, initialization or writes. Never mount or write the source PVC.

List the accepted immutable snapshot and dump its completion marker into a local
task scratch directory with mode 0700 and files mode 0600. Validate its strict
native filename and freshness against the snapshot. Dump **that same snapshot's**
selected native database and require its SHA256 to equal the archived marker.
Open only the restored file with SQLite, require `PRAGMA integrity_check` to
return `ok`, and require the expected Vaultwarden tables (including `users`,
`ciphers`, `attachments`, `organizations`) even when their contents are empty.
List/hash attachment, key and config files privately without displaying contents.
An empty initial account table is valid before the owner's own enrollment.
Remove local decrypted database/key copies after recording the bounded receipt.

This proves retrieval/decryption and SQLite structure from an actual remote
archive. It does not prove owner login, client sync, attachment retrieval or a
complete running server restore. Those require a separate bounded test after
owner enrollment; do not invent an identity to claim them.

## Full recovery procedure

1. Preserve the current source claim and archive. Recover the repository address,
   Restic password and Spaces access independently from the strict SealedSecret
   plus protected controller recovery keys, or an owner-held secure credential
   escrow. A SealedSecret ciphertext alone cannot decrypt its secret without the
   matching controller private key. No separate owner escrow is claimed here.
2. Restore the chosen full encrypted archive into an **isolated new directory or
   claim**, not the production source. Verify its marker/hash and SQLite integrity
   as above. Retain its attachments, sends, signing keys and configuration.
3. Keep the target Vaultwarden stopped. Remove only the restored target's live
   `db.sqlite3`, `db.sqlite3-wal`, `db.sqlite3-shm` and old journal, then install the
   verified native standalone file as `db.sqlite3`. Never replay a live DB's WAL
   against the native copy. Preserve the original source unchanged. Set target
   ownership to 1000:1000 and restrictive permissions, then start one compatible
   pinned Vaultwarden pod with no public ingress.
4. After explicit owner login, verify vault items, an attachment download, a client
   sync and the intended private origin. Record actual results and retain the
   original claim/archive until the owner accepts recovery. Follow upstream
   version/migration guidance before changing server versions.

## Verified initial recovery — 2026-10-25 UTC

The initial encrypted archive and an isolated server restore passed. Snapshot
`57a7711032f0628ca20003192b1465976fc3820bee8e4437dbe854fb3b526eab`, taken at
**23:09:37.87693726 UTC**, contains the native
`db_20261025_230700.sqlite3` backup completed at **23:07:00 UTC**. Its archived
completion marker identifies the same database and SHA256; the database was
about 158 seconds old at snapshot time. The source acknowledged manual tag
`initial-native-20261025T230700Z` and completed successfully at 23:09:42 UTC.

Read-only Restic retrieval restored all 6 files (591,654 bytes), including the signing
key, into a private local target. The selected standalone native database matched
its marker hash, passed SQLite integrity and foreign-key checks, and contained
29 tables. The test installed that native copy as the target's `db.sqlite3` after
removing only the target's live DB/WAL/SHM/journal; it preserved the production
claim and did not create remote locks or repository writes.

The same pinned Vaultwarden 1.37.4 image then started against the recovered state
in an isolated local container as UID 1000 with a read-only root filesystem and all
Linux capabilities dropped. Its database-aware `/alive` check passed, the web UI
returned 200, and registration remained disabled. The test containers, restored
plaintext database/key directories and temporary runtime environment were removed.

This initial vault contained 0 users and 0 ciphers, with no attachments, sends or
`config.json`. No account was fabricated. The proof covers archive decryption,
SQLite structure, signing-key recovery and server startup; owner login, populated
vault/client sync and attachment recovery still require the owner's own enrollment
and data. It does not claim those workflows passed.

Sources:
[Vaultwarden backup guidance](https://github.com/dani-garcia/vaultwarden/wiki/Backing-up-your-vault),
[1.37.4 native SQLite implementation](https://github.com/dani-garcia/vaultwarden/blob/1.37.4/src/db/mod.rs),
[VolSync Restic source](https://volsync.readthedocs.io/en/stable/usage/restic/index.html),
[Restic restore and dump](https://restic.readthedocs.io/en/stable/050_restore.html).

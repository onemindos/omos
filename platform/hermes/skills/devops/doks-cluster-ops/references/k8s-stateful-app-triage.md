# K8s Stateful App Triage

Rules for any stateful Deployment/StatefulSet on DOKS that crashloops, fails to start, or cannot write config. Apply before deeper debugging.

## 1. Postgres: PGDATA subdirectory (avoid `lost+found`)

When a Postgres container mounts a PVC at `/var/lib/postgresql/data`, the directory is not empty — it contains `lost+found` from the block device mount. initdb refuses to run in a non-empty directory. Fix: set `PGDATA=/var/lib/postgresql/data/pgdata` and create that subdirectory in an initContainer (run as root, chown to 999:999).

## 2. Postgres collation: use `LC_COLLATE 'C'` / `LC_CTYPE 'C'` for apps that mandate it

Synapse (and some other apps) reject databases created with locale collation (e.g. `en_US.utf8`). When an app errors with `Database has incorrect collation`, recreate the database with `LC_COLLATE 'C' LC_CTYPE 'C' TEMPLATE template0`. Do not rely on the cluster default.

## 3. Redis: flag is `--requirepass` (underscore breaks config)

The Redis command-line flag for password auth is `--requirepass` (no underscore). Passing `--require-password` produces `Bad directive or wrong number of arguments`. Check `redis-server --help` for the exact spelling when auth fails at startup.

## 4. Read-only ConfigMap mounts: write generated files to the PVC, not the ConfigMap

ConfigMap volumes are read-only in the pod. If an app generates config files (e.g. Synapse's `log.config`), point the generation path at the mounted PVC (e.g. `/data/log.config`) — not at the ConfigMap mount. Otherwise the app crashes with `Read-only file system`.

## 5. Image registry: verify the image name/tag before assuming a pull failure is auth

`element-web` (and some other images) moved from a vendor-specific registry to a canonical registry (e.g. `vectorim/element-web`). An `ErrImagePull` with a registry that requires auth often means the image name is wrong — search the official docs for the canonical image name before creating image-pull secrets.

## 6. RWO PVC + initContainer permission pattern

For any app running as non-root that needs to write to a fresh PVC: use an initContainer running as root (`securityContext.runAsUser: 0`) that `chown`s the PVC to the app's UID/GID (e.g. 991:991 for Synapse, 999:999 for Postgres). Set `fsGroup` on the pod so the volume is group-writable. Do not rely on the app to fix permissions at runtime.

## 7. Element Web image: use `vectorim/element-web`

The Element Web image is published under the `vectorim` org (e.g. `vectorim/element-web:v1.11.70`), not `elementhq/element-web`. Using the latter produces ErrImagePull.

## 8. Redis command: `--dir` for data subdirectory

When a Redis container mounts a PVC and the default data directory contains `lost+found`, pass `--dir /data/redis-data` (created by the initContainer) so Redis writes to a subdirectory, not the mount root.

## 10. Synapse auto-generates log config; pre-seed it or the file mount breaks
Synapse always tries to generate `log.config` on first run if the path doesn't exist. If `log_config` points at the ConfigMap mount, it crashes with `Read-only file system`. If the ConfigMap pre-seeds a log.config with a FileHandler, it crashes because the image root FS is read-only. Pre-seed with a StreamHandler (`ext://sys.stderr`) and point the path at the writable PVC (`/data/log.config`).

## 11. Drop and re-create database if collation was wrong on first run
Synapse checks DB collation at startup (`IncorrectDatabaseSetup`). If the database was created with the default locale (`en_US.utf8`), no config change fixes it — DROP DATABASE and re-CREATE with `LC_COLLATE 'C' LC_CTYPE 'C' TEMPLATE template0`. Do this BEFORE the first successful Synapse start to avoid re-processing.

## 12. Element Web image: `vectorim/element-web` (elementhq does not exist)
The official Element Web image is `vectorim/element-web:v1.11.70+`. `elementhq/element-web` returns ErrImagePull with auth errors — do not create pull secrets, fix the image name. The old `elementhq` org was retired; docs may still reference it.

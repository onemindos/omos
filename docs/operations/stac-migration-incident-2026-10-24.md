# STAC migration containment — 2026-10-24

The first automated `platform-stac` rollout began after the OpenChoreo deletion blocker cleared at approximately 08:20 UTC. At 08:22 UTC the shared `onemind-data` cluster became unready. Primary logs show `relation "pgstac.migrations" does not exist`, followed by a server process terminated by signal 11 (segmentation fault), `restart_after_crash is off`, and PostgreSQL shutdown. The migration hook repeatedly retried; this affected OpenFGA/core and other shared database clients. The failing SQL/native extension and a causal reproduction remain unverified. Do not retry the migration against the shared database to investigate.

Emergency containment: terminate the STAC Argo operation, delete only the migration Job, and temporarily annotate `platform-stac` with `argocd.argoproj.io/skip-reconcile=true`. No database, PVC, role, credential, disk, recovery point or backup was deleted or restored. The browser can remain running, but API and imagery acceptance are unfinished.

This change makes the Application manual with no sync retry, and the Job single-attempt with a ten-minute deadline. A PostgreSQL migration failure must require an explicit investigation, rather than repeated restarts of a shared service.

After merge, let app-of-apps apply the manual Application policy without pruning other parent Applications. Verify `spec.syncPolicy.automated` is absent, the old operation has stopped and the migration Job is absent. Remove the temporary skip-reconcile annotation, but do not start a sync. Confirm the primary and both replicas recover, and verify database-dependent clients separately. If automatic recovery does not complete, follow the existing CNPG recovery runbook; do not recreate PVCs or reset WAL.

Before enabling STAC: reproduce the exact pgstac image/SQL on a disposable instance of the exact current PostgreSQL image, investigate the native crash, validate migrations/roles/extensions, verify a recoverable backup and review the complete shared-database impact. Roll out once under observation after the reviewed fix. Re-enable automation only in a separate reviewed change after live acceptance.

# Observer recovery gate (2026-10-19)

One-shot Argo-owned restore from the distinct encrypted observer-alerts repository
into a NEW 1 GiB PVC. The production PVC is never mounted here or modified.
The verification Job has a read-only mount, no token, and denied network egress; it opens no listening service.
The mover and verifier share the existing ecosystem pool to avoid RWO attach races.
Success requires a completed ReplicationDestination plus SQLite integrity_check=ok
and the alert_entries/incident_entries tables. No sample alerts are fabricated.
The restore PVC is explicitly protected from Argo prune/deletion. Keep it as evidence;
cleanup requires explicit storage-deletion approval, never deleting the source PVC.

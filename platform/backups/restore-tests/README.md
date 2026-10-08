# One-off snapshot restoration

These manifests are deliberately excluded from the recurring backup Kustomization.
After their reviewed PR merges, apply the test manifest in the explicit production
context. This is a verification operation, not a replacement application deployment.
The Authentik test restores its retained snapshot to a new 5 GiB claim, starts
PostgreSQL 17 on loopback, and queries the restored application's table inventory.
All recovery writes target the clone. Production claim bindings are unchanged.

```sh
kubectl --context do-nyc3-omos-cloud apply -f platform/backups/restore-tests/authentik-snapshot.yaml
kubectl --context do-nyc3-omos-cloud -n auth get job auth-snapshot-verify-20261008
kubectl --context do-nyc3-omos-cloud -n auth logs job/auth-snapshot-verify-20261008
```

Record completion and sanitized result before cleanup. Delete only this temporary
Job and clone claim after successful verification. Its Delete-policy storage class
removes the temporary provider volume; the original claim and Retain-policy
snapshot stay intact. Do not reapply a completed test after cleanup unintentionally.

```sh
kubectl --context do-nyc3-omos-cloud -n auth delete job auth-snapshot-verify-20261008
kubectl --context do-nyc3-omos-cloud -n auth delete pvc auth-snapshot-verify-20261008
```

A table inventory proves that this snapshot can boot PostgreSQL and read schema.
It does not prove login/session recovery, every row, or recovery of every other PVC.

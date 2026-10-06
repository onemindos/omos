# Isolated recovery exercise

These resources are deliberately excluded from the production kustomization. They create an ephemeral PostgreSQL 17 instance with emptyDir storage, not a production PVC.

1. Apply `instance.yaml` and wait for `backstage-recovery` to be Ready.
2. Create ConfigMap `backstage-recovery-script` from `verify.py` in namespace backstage.
3. Apply `job.yaml`; wait for completion and retain its JSON evidence.
4. Delete the Job, Pod, Service, NetworkPolicy and ConfigMap after recording evidence.

The verifier downloads the newest encrypted Spaces archive, validates AES-GCM with the existing key, restores every application database and checks catalog entities. PGHOST must be the dedicated recovery Service. Retry requires a fresh temporary instance; never point this script at production. The exercise proves application database recovery, not Authentik recovery or complete infrastructure disaster recovery.

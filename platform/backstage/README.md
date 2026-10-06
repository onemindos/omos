# OneMind Backstage

Private preview: https://onemind-preview.tail717ba0.ts.net
Identity provider: https://onemind-auth.tail717ba0.ts.net

This application is built from upstream Backstage 1.55.3 using the new frontend and backend systems. Source is `apps/onemind-portal`. Production catalog files are mounted from the generated ConfigMap in this directory; reviewed inventory changes do not require rebuilding the image. Argo CD manages this directory after merge. Argo uses the existing authenticated repository URL (the old omos-infra alias redirects to omos). Container images are pinned by digest; application changes require a successful image build and a reviewed digest update. The previous Homepage application was retired after browser review. Its inactive source remains in platform/homepage for recovery.

## Identity and access

Tailscale provides private HTTPS access. Authentik OIDC supplies Backstage identity and membership. The existing `akadmin` account is registered as the initial catalog user and sole allowed Authentik application user. Backstage maps the signed stable subject to `onemind.io/oidc-sub` in the user catalog; it does not trust unverified email as identity. The `onemind-admins` group permits configured Backstage actions; other identified users receive read permissions. No guest login is enabled. Adding a user requires a catalog annotation `onemind.io/oidc-sub` matching the signed subject, Authentik application access, and deliberate group membership.

The Kubernetes service account has observation permissions, including pod logs, but cannot read Secrets or modify workloads. Portal actions do not receive a GitHub credential. MCP exposes configured catalog actions; machine authentication must be configured separately before connecting an AI client.

Cluster OIDC requests use the existing Tailscale egress ProxyGroup. CoreDNS rewrites only the identity provider hostname to the operator-managed egress Service while TLS still validates the original hostname. The operator owns that Service's ExternalName, which Argo ignores. CoreDNS must reload newly added custom files; verify resolution and HTTPS after a change.

## Database and recovery

PostgreSQL 17 uses a 10 GiB retained DigitalOcean volume. The application role can create plugin databases but is not a superuser. Credentials are SealedSecrets. The single application replica uses Recreate updates to prevent concurrent initial migrations; updates have a brief availability gap. PostgreSQL and the portal are not highly available.

`backstage-backup` runs daily at 07:15 UTC. It dumps all databases owned by the application role, encrypts the compressed archive with AES-256-GCM, and uploads to the existing private Spaces bucket under `backups/backstage/`. Each run downloads the object, checks authenticated decryption, restores the catalog to a uniquely named temporary database, verifies tables, and removes only that temporary database. Job failure means backup or recovery verification failed. Objects are retained; no automatic deletion policy is configured.

Recovery requires the backup Secret's encryption key, Spaces access, database credentials, and access to the SealedSecrets controller's sealing key when rebuilding the cluster. Keep an independent secure copy of recovery keys; Git alone cannot decrypt a SealedSecret on a new controller. Do not rotate or remove the encryption key without preserving access to prior archives.

For a recovery exercise, create a Job from the CronJob and inspect its `restoreVerified` output. For full recovery, stop Backstage, restore the selected archive into a separate PostgreSQL instance using the same application role and each original database name, verify all databases and catalog entities, then point the application at that instance. Do not restore over a running database. The backup includes database state; application code, catalog source and deployment manifests live in Git. No TechDocs content is published yet.

## Catalog and live state

Initial app links were imported from the now-archived Homepage via `scripts/migrate-homepage.cjs`. Maintain new links in the Backstage directory JSON and both apps.yaml catalog copies; rerunning the legacy migration can overwrite later additions. Cluster catalog identities are an explicit reviewed snapshot from `scripts/snapshot-cluster.cjs`, not automatic discovery. Stable hash suffixes prevent link/workload identity collisions. Rerun and review the snapshot when adding or removing workloads. The snapshot script updates both development and production copies. The migration script mirrors app links; user enrollment must update both org.yaml copies. The operations panel queries the live cluster every minute and marks unavailable sources and stale results. Existing unrelated unhealthy workloads are shown, not hidden.

## Release checks

Use the official Backstage Kubernetes, Docker, OIDC, permissions and Kubernetes integration documentation. Require: immutable dependency install, TypeScript check, backend/container build, server-side manifest validation, healthy rollout, verified OIDC callback and identity, authenticated catalog and live status, denied anonymous APIs, read-only cluster permissions, and an encrypted backup with successful restore verification. Verify actual browser rendering before replacing Homepage.

Do not claim offline operation: self-hosted pages still require reachability to the host. Existing cloud app links require their providers. Home hosting can use the same containers, Git catalog and restored PostgreSQL with a suitable storage class and mesh routing.

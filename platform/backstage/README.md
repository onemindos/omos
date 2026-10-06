# OneMind Backstage

Private preview: https://onemind-preview.tail717ba0.ts.net
Identity provider: https://onemind-auth.tail717ba0.ts.net

This application is built from upstream Backstage 1.55.3 using the new frontend and backend systems. Source is `apps/onemind-portal`. Production catalog files are mounted from the generated ConfigMap in this directory; reviewed inventory changes do not require rebuilding the image. Argo CD manages this directory after merge. Argo uses the existing authenticated repository URL (the old omos-infra alias redirects to omos). Container images are pinned by digest; application changes require a successful image build and a reviewed digest update. The previous Homepage application was retired after browser review. Its source was removed; Git history retains the pre-retirement configuration for recovery.

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

Maintain app links in the Backstage directory JSON and both apps.yaml catalog copies. Cluster catalog identities are an explicit reviewed snapshot from `scripts/snapshot-cluster.cjs`, not automatic discovery. Stable hash suffixes prevent link/workload identity collisions. Rerun and review the snapshot when adding or removing workloads. The snapshot script updates both development and production copies. user enrollment must update both org.yaml copies. The operations panel queries the live cluster every minute and marks unavailable sources and stale results. Existing unrelated unhealthy workloads are shown, not hidden.

## Release checks

Use the official Backstage Kubernetes, Docker, OIDC, permissions and Kubernetes integration documentation. Require: immutable dependency install, TypeScript check, backend/container build, server-side manifest validation, healthy rollout, verified OIDC callback and identity, authenticated catalog and live status, denied anonymous APIs, read-only cluster permissions, and an encrypted backup with successful restore verification. Verify actual browser rendering before replacing Homepage.

Do not claim offline operation: self-hosted pages still require reachability to the host. Existing cloud app links require their providers. Home hosting can use the same containers, Git catalog and restored PostgreSQL with a suitable storage class and mesh routing.

## Inventory maintenance

Run `node scripts/enrich-inventory.cjs` from `apps/onemind-portal` on the authorized Mac to refresh both catalog copies. It checks the cluster context, captures workload identities and Tailscale endpoint names/OS/DNS, and connects app links to workloads only through matching ingress hosts and service selectors. It does not collect credentials, user profiles or IP addresses. Review the Git diff before merging; the script is not scheduled automatically. Plain `snapshot-cluster.cjs` replaces cluster.yaml and omits mesh enrichment, so use the enrichment command for the complete inventory.

Hosting annotations distinguish verified cluster routing from unverified mesh endpoints and external providers. Device location and purpose are explicitly unreviewed; an offline discovery snapshot is not a service incident. The catalog graph exposes verified application-to-workload and workload-to-cluster relationships. Cloud API integrations, repository discovery beyond the infrastructure repository, and physical device classification still require reviewed sources.

## Curated Fabric registry

The reviewed Notion snapshot (`catalog/fabric-source.json`, 174 records, reviewed 2026-10-06) enriches existing app and mesh entries and adds documented capabilities, planned items and references. `scripts/reconcile-fabric.cjs` merges duplicate titles and Google Calendar/Drive/Photos aliases while retaining every source page link. Recorded status and hosting are provenance annotations, not live health or verified deployment. Original protocol addresses and old IPs remain in source documentation rather than becoming launch buttons. A source snapshot update requires review. The complete enrichment command reapplies this reconciliation so cluster refresh cannot erase curated records.

## P2 completed inventory scope

The complete inventory refresh also records accessible OneMind repositories and exact runtime identifiers, then verifies reviewed aliases. Tags `deployment-verified`, `repository-verified`, `documented-unverified`, `planned`, and `reference` are available through Catalog's Tags filter. A repository match is not evidence of deployment; a deployment match is not live health. Device roles remain documented unless independently inspected. Unknown SaaS account state stays unknown.

All 174 Fabric source pages were fetched. 165 had no body content; nine contained historical notes. Git-owned Markdown is in `docs/inventory/fabric/` and the completeness manifest is `docs/inventory/migration-manifest.json`. Catalog source links now point to these archives rather than Notion. This removes the catalog's Notion content dependency; it is not a full export of the rest of the Notion workspace, attachments, comments or unrelated pages. Notes retain historical dates and are not validated deployment recipes.

## P3 health scope

The authenticated operations endpoint now checks four fixed internal service targets: Authentik readiness, Grafana health, Headlamp web response and NATS dashboard web response. Checks run on demand with a shared 60-second cache and in-flight deduplication; each GET has a five-second absolute deadline and two-megabyte response cap, no redirects or client-supplied target URLs. This is cluster-side point-in-time availability, not uptime history, browser login, tailnet path verification, NATS message delivery or TAK certificate validation. Non-200 results are down except 401/403 which are explicitly authentication-required. Transport failures are unknown.

The same server reads fixed Prometheus queries for five namespaces: backstage, auth, nats, takserver and onetak. CPU is a five-minute rate in cores; memory is working-set bytes. Empty series are unknown, not zero. No monitoring credentials reach the browser and no cluster write permissions were added. The panel exposes timestamps and investigation/catalog links. External monitoring and continuous alert routing can be added separately when needed.

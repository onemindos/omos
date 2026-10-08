# Recovery and migration to home infrastructure

Git, portable data archives and independently held keys are the recovery inputs.
DigitalOcean snapshots alone cannot supply a home restore. Run this procedure on
an isolated target before changing production DNS, identity callbacks or traffic.
Never attach a test instance to production NATS routes or enable live automations.

1. Record hardware architecture, disks, power protection, network interfaces and
   failure domains. Existing app images include amd64-only custom builds; verify
   every required digest on the target architecture before selecting ARM hardware.
   Three replicas on one physical computer do not survive loss of that computer.
2. Install a supported Kubernetes distribution and a CSI driver appropriate to the
   actual disks. Create an explicit home storage-class overlay; remove DigitalOcean
   load-balancer annotations and snapshot-class dependencies in that overlay.
   Preserve original production manifests until the isolated deployment passes.
3. Recover the off-cluster logical archive key, Restic password and SealedSecrets
   private keys. Keep copies independent of the cluster and this laptop. Validate
   object-store access separately; rotating an access key must preserve encryption
   keys. Do not wholesale replay service-account tokens or provider-specific secrets.
4. Bootstrap namespaces, networking/DNS, ingress, PKI, SealedSecrets and Argo. Apply
   `platform/argocd/app-of-apps.yaml` only after selecting the home overlay/revision:
   current main describes DigitalOcean and must not be blindly applied to new disks.
   Restore the original sealing keys before expecting existing ciphertext to decrypt.
5. Restore data to NEW PVCs. Use verified logical PostgreSQL archives into compatible
   versions/extensions (TAK needs PostGIS). Backstage needs its complete plugin DB
   set, not a single database. Restore NATS into isolated servers/accounts first.
   Restic exports can restore files onto a different CSI filesystem; raw DB copies
   require database recovery and application checks. Never overwrite source PVCs.
6. Recover identity and registry access before dependent services. Preserve issuer
   and client identity, user/group mappings, signing keys and local certificate trust.
   Validate callback URLs and mesh DNS. Ensure local DNS/PKI and operator login work
   without cloud connectivity. Preserve mesh routes; approve any route overlap.
7. Resolve manually managed releases from their recorded chart version, values and
   configuration, including NATS and Matrix. Adopt each release through a focused
   migration with a diff and data-preservation check. Mixed ownership is not a
   reproducible bootstrap contract.
8. Test user login, authenticated AI calls with least privilege, TAK CoT ingress and
   delivery, databases, queue/stream continuity, uploads, scheduled jobs and restore
   markers. Exercise node loss and a disconnected WAN. Record measured recovery
   time, recovered timestamp and application-specific data checks.
9. Plan a final consistent export and explicit cutover window. Disable duplicate
   workers/schedulers, move DNS/mesh endpoints, verify traffic and retain a tested
   rollback source. Cloud infrastructure can remain an outpost after home becomes
   primary. Public tunnels are optional and must not become a local recovery need.

## Repository layout

`platform/argocd` owns application paths; `platform/<service>` owns shared platform
manifests; `tak-server/deploy` owns TAK workloads; `apps/onemind-portal` owns portal
source/catalog/runbooks; `platform/backups` owns backup workloads, ciphertext and
restore templates; `docs/operations` owns observed evidence and recovery procedures.
Keep these stable paths while reorganizing service internals. Update Argo and
Kustomize references in the same PR as any move. Protected local inventories and
cleartext recovery keys must not enter Git.

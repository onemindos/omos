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

## Independent recovery inputs

The protected recovery copy is at `~/OneMindRecovery/2026-10-08`; encryption keys
remain separately under `~/.config/onemind/recovery`. Consult mirror-manifest.json
and offline-integrity.json for the copied repository set and actual full-data
integrity checks. This laptop copy is independent of DigitalOcean, but loss of the
laptop still requires another physical/offsite copy of data AND keys.

`platform/backups/tools/mirror.py` downloads verified encrypted archives and Restic
repositories without deleting source data; capacity and free-space guards are
mandatory. `tools/decrypt.py` handles the current archive and two older formats
without a running cluster. `tools/export_helm.py` extracts deployed Helm release
inputs from the authenticated Kubernetes archive. Eleven release inputs were
recovered locally. Extracted values/manifests can contain credentials: keep them
private and out of Git. Existing Helm manifests can predate later GitOps overrides;
review NATS sizing, TAK quotas, provider storage and ingress before applying them.

No new home storage overlay can be verified before the disks and CSI driver are
chosen. Recovery instructions and portable backups are preparation; they are not
evidence of a successfully migrated or WAN-independent production system.

A private AES-GCM archive of the tracked main-branch repository was also created
and authenticated locally, including mission memory, Argo paths and recovery tools.
repository-recovery.json records its exact commit and checksum; it contains no
Git history. Keep a copy of the archive decryption tool available independently
to avoid needing GitHub merely to open the recovery bundle.

A disconnected-WAN drill also needs required container images and Helm charts
available locally, plus usable DNS, PKI and identity. No air-gapped registry/chart
cache or disconnected application exercise was established by this review.

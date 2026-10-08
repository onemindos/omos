# Cluster readiness: October 8, 2026

This is an evidence record, not a guarantee of reliability. The review targets
`do-nyc3-omos-cloud`. Kubernetes, TAK, authenticated AI access and a future
home-primary deployment remain the mission. Source manifests, observed deployment
state and successful restore exercises are distinct evidence.

## Verified repairs

- Three workers were Ready on Kubernetes 1.36.0. Kubelet heap profiles matched
  upstream issue 139823 / fix 139850. Sequential kubelet service restarts reclaimed
  approximately 10 GB without rebooting workers. This mitigates accumulated memory;
  it does not fix the running release. The fix was backported in PR 140066 and is
  included in 1.36.3. DigitalOcean publicly lists 1.36.3-do.5; account-supported
  upgrade options still require authenticated provider access.
- Health keeper observes unhealthy pods and cannot delete pods (PR 84). Its scheduled
  jobs completed. The original diagnostic job had already been garbage-collected.
- TAK's Argo Degraded status came from an unmanaged registry Secret, despite ready
  workloads. Its working bytes were resealed and adopted; Argo became Healthy.
- NATS snapshot attempts exposed a separate real OOM at the previous 1 GiB limit.
  All three servers now have a 2 GiB limit, GOMEMLIMIT=1536MiB, requests of 100m/512Mi
  and required separation by hostname. Partition 0 rollout completed. The verified
  backup restored all 39 streams into an isolated server; production servers did
  not restart during that successful exercise. One stream is single-replica;
  this backup does not make it highly available.
- Vision and TAK plugin requests were reduced from observed over-reservations while
  preserving limits. TAK plugins gained startup/readiness probes. Ready rollouts
  demonstrate deployment health, not a peak-load capacity test.
- CI now enforces changed YAML on hosted runners and reports the existing full-repo
  YAML debt separately. Full `make validate` still fails on legacy formatting;
  passing changed-file checks must not be represented as a clean repository.

## Backup evidence and limits

Encrypted object-store archives were downloaded, authenticated and restored into
separate disposable database/server processes. Completed production exercises:

| Service | Verified restored content |
| --- | --- |
| Authentik | 230 application tables |
| n8n | 148 tables |
| Matrix | 151 tables |
| TAKML | 1 table |
| Zeus site | 20 tables |
| Backstage | All 23 role-owned plugin databases; 145 tables |
| NATS | 39 streams, including consumers in the archive |
| TAK / PostGIS | 64 tables; isolated restore completed in a 144-minute exercise |

These checks establish archive readability and restore execution. They do not prove
all user workflows or row-by-row equality. TAK's logical COPY data expanded far
beyond its physical database size, so the job deadline is now six hours. Scheduled
recurrence and measured recovery time on home hardware still need observation.
ClickHouse's unverified legacy backup was suspended after its encrypted filesystem
backup booted successfully on a disposable PVC.

All 35 original bound PVCs have Ready CSI snapshots in the dated readiness set,
using `onemind-backup-retain`. A separate Authentik snapshot clone booted PostgreSQL
and exposed 230 tables; that temporary clone and test Job were removed. Snapshots
are crash-consistent and DigitalOcean-specific, not portable offsite archives.
Retention may incur provider charges and requires an explicit lifecycle decision.

All 35 original volume sources completed their initial checks: 32 encrypted Restic
repositories and three confirmed empty volumes (Nominatim, Overpass and tileserver).
All sources now have daily schedules staggered 20 minutes apart. Raw database
volume copies remain crash-consistent. The filesystem inventory excludes temporary
test/backup volumes. The 35 additional PVCs are persistent 1 GiB Restic caches;
no source PVC was removed.

Representative filesystem restore drills used NEW PVCs and were removed afterward:
Authentik booted PostgreSQL with 230 tables; Hermes restored 319,496 files and its
reference binary checksum matched; ClickHouse booted with nine tables and
74,916,246 rows. These are representative checks, not full validation of all apps.

An encrypted recovery archive was downloaded and decrypted with exact equality,
including 241 Secrets, 154 ConfigMaps and two active SealedSecrets private keys.
The Kubernetes archive passed authenticated decryption and payload equality; it
was not applied wholesale to another cluster. Recovery encryption keys are held
outside Git and Kubernetes under `~/.config/onemind/recovery/`, mode 0600. An
independent offline copy is still required; a file on this laptop is not sufficient
protection against loss of the laptop.

## Networking, ownership and retained namespaces

Internal HTTP checks passed for Backstage, Authentik, Matrix/Element, WebODM, STAC,
TAK portal, Martin, Maputnik, pg-featureserv, Valhalla and tileserver. n8n's local
health endpoint passed; cross-namespace access from TAK was intentionally denied
by its policy. OpenClaw returned an access gate. These checks are not authenticated
end-to-end application acceptance tests or a live TAK CoT exchange.

Backup egress is restricted to DNS and HTTPS, including the observed Spaces DNS
address 10.10.255.254/32 where required. Tailscale's intentional routes and the
user's machine remain connected; the disconnection request was withdrawn.

NATS still has split ownership: Argo owns the StatefulSet; its existing Helm release
owns configuration and Services. Matrix and some legacy services also have manual
ownership. Do not run an overlapping Helm upgrade or assert complete GitOps parity.
No namespace was established as safely unused. Standard namespaces and dormant
Git-owned services were preserved.

## Outstanding dependencies and acceptance gates

- Local `doctl` still returns 401. Provider-supported Kubernetes upgrade selection,
  account backup settings and rotation of the previously exposed Spaces access key
  remain blocked. Enter credentials locally with `doctl auth init`, not in chat.
- Observe the first actual scheduled recurrence and retention/pruning cycle. One-off
  successes do not establish a backup service's long-term reliability. Extend
  database-specific recovery checks where raw volume copies are the only coverage.
- Resolve manual ownership, remaining YAML validation debt and load-test capacity
  before treating the cluster as ready for unlimited additional applications.
- Home migration requires target hardware/storage and an isolated recovery drill.
  See [the migration runbook](home-migration.md). No home deployment has been tested.

Upstream evidence: https://github.com/kubernetes/kubernetes/issues/139823,
https://github.com/kubernetes/kubernetes/pull/139850,
https://github.com/kubernetes/kubernetes/pull/140066,
https://docs.digitalocean.com/products/kubernetes/details/changelog/.

## Final operational checks

All 23 Argo applications were Synced and Healthy; all enabled Deployment and
StatefulSet controllers were ready, with no Running/Pending pod exceptions.
All original 35 PVCs remained Bound and all 35 retained baseline snapshots Ready.
The attached sanitized inventory records the observation timestamp and volume map:
[cluster-readiness-evidence](cluster-readiness-2026-10-08.json).

The persistent GitHub runner now survives jobs. Incorrect nonempty "false" flags
were removed, images pinned, Kubernetes token mounts and direct CI rollout writes
removed. DinD's unauthenticated TCP listener binds only 127.0.0.1; ingress is denied
and egress limited to DNS/public HTTP(S). A real self-hosted smoke job passed after
the networking change (run 37780763483). Its privileged DinD container still deserves
isolation from untrusted jobs; this review does not certify arbitrary builds.

Backstage's AI identity read the new backup-service catalog entry (HTTP 200);
unauthenticated access returned 401 and a machine write returned 403. Catalog
integration provides service inventory and operational links, not credential access.

VolSync metrics scrape is UP. Four loaded rules report healthy evaluation:
ScheduledBackupOverdue, PortableVolumeBackupOutOfSync, VolSyncMetricsUnavailable
and KubeletResidentMemoryHigh. No alerts were active at the final observation.
Kubelet RSS remained about 156–171 MiB, versus multi-gigabyte heaps before mitigation.

Changed-file CI, server-side dry-runs, backup image integration tests, restore
exercises and Terraform validation passed. Full-repository YAML validation still
reports 498 errors and four warnings (formatting debt); independent ShellCheck
fails with 36 diagnostics in existing scripts. Neither is a passing full-repo gate.

Recent CSI deletion warnings came from temporary clones still detaching; the
final inventory had no non-Bound PVs and no remaining restore/clone PVCs.
Namespaces were retained because none met evidence-based deletion criteria.
Repository paths remained stable, with service-specific backup/restore/runbook
folders; unrelated changes in the mission checkout were preserved.

## Capacity and remaining resilience limits

At the final sizing check, per-node memory reservations were approximately
89–94% and CPU reservations 66–84%; observed memory working sets ranged 46–89%.
Backstage's consistently observed 279–303 MiB per replica supports reducing its
request from 1 GiB to 512 MiB while retaining a 2 GiB limit. Idle measurements do
not justify reducing batch geospatial/photogrammetry reservations blindly.
The three-worker cluster does not yet demonstrate enough spare capacity to
reschedule all workloads after losing a worker. Several databases remain
single-instance; backups improve recovery but do not provide database failover.
A large app expansion needs workload sizing, load tests and an explicit capacity
decision. HPA alone would not create additional worker capacity.

## Independent recovery copy

The protected laptop copy contains all 32 encrypted Restic repositories and ten
verified logical/configuration archive checksums: 3,103 objects, 49.03 GiB.
It was copied without deleting remote data, with a 40 GiB free-space reserve.
All 32 repositories passed full offline Restic check --read-data; results are
recorded in the accompanying evidence JSON. A download/checksum is distinct from an application restore.

The authenticated Kubernetes archive was decrypted locally without provider
access; 11 deployed Helm release inputs were extracted into protected local
files. These bootstrap files can contain plaintext credentials, so the whole
recovery directory must remain private. The encrypted data copies and independently
held keys still need another physical/offsite copy to survive laptop loss.

A Backstage replica was rebalanced from the most reserved worker using the
Eviction API; the existing disruption budget allowed it and another replica
remained ready. Final resource readings and rollout state are recorded separately.
A tracked-source recovery archive was encrypted and authenticated locally,
including mission memory and the GitOps/recovery source, without Git history.

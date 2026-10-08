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

These checks establish archive readability and restore execution. They do not prove
all user workflows or row-by-row equality. TAK's roughly 10 GB database exercise is
still running at this record's first publication; do not infer success. ClickHouse's
legacy backup reports success without an isolated restore or encryption assurance.

All 35 original bound PVCs have Ready CSI snapshots in the dated readiness set,
using `onemind-backup-retain`. A separate Authentik snapshot clone booted PostgreSQL
and exposed 230 tables; that temporary clone and test Job were removed. Snapshots
are crash-consistent and DigitalOcean-specific, not portable offsite archives.
Retention may incur provider charges and requires an explicit lifecycle decision.

VolSync/Restic filesystem exports are staged under `platform/backups/volsync`.
Initially only Hermes is active; paused sources are not covered by completed
filesystem exports. Raw database volume copies remain crash-consistent. The
filesystem inventory excludes temporary test/backup volumes.

An encrypted recovery archive was downloaded and decrypted with exact equality,
including active SealedSecrets private keys. Recovery encryption keys are held
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
- Self-hosted GitHub runner repeatedly exits and enters restart backoff. Hosted
  quality checks no longer depend on it; its operation still needs repair.
- Establish completed filesystem exports and representative application restores;
  complete TAK and ClickHouse verification and extend remaining database-specific
  recovery where logical dumps are needed.
- Resolve manual ownership, remaining YAML validation debt and load-test capacity
  before treating the cluster as ready for unlimited additional applications.
- Home migration requires target hardware/storage and an isolated recovery drill.
  See [the migration runbook](home-migration.md). No home deployment has been tested.

Upstream evidence: https://github.com/kubernetes/kubernetes/issues/139823,
https://github.com/kubernetes/kubernetes/pull/139850,
https://github.com/kubernetes/kubernetes/pull/140066,
https://docs.digitalocean.com/products/kubernetes/details/changelog/.

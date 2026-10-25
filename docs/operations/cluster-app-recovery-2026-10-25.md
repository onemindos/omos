# Cluster application recovery — 2026-10-25

Zeus requested STAC repair, first deployment of CloudTAK and Karakeep, and resolution of the ten OutOfSync Argo Applications. This is bounded application recovery; it does not change platform phases or reduce node pools.

## Verified starting point

At 15:46–15:48 UTC: four Ready nodes, 186 Running/Ready pods, 76 Bound PVCs, 53/56 Healthy Argo Applications and 46/56 Synced. STAC read/write API Deployments were absent, CloudTAK/Karakeep namespaces absent, and the disposable STAC validation control failed before its migration. Scheduled logical backups and 35 VolSync sources completed today.

## Reconciliation decisions

- Adopt Argo CD with the existing pinned chart, no prune/replace and no pod-template changes. Exact-main diff showed only tracking annotations. The manual sync succeeded; all six serving pod UIDs remained unchanged. Manual policy stays in place.
- Reconcile the app-of-apps STAC-validation child definition: only the explicit directory.recurse=false field differed. This does not retry the diagnostic Jobs.
- Seal Matrix’s current homeserver.yaml directly through memory/pipes into a strict SealedSecret. Preserve bytes and server identity; do not rotate credentials, reset databases, or delete existing claims. First selectively sync the sealed prerequisite and compare the resulting config privately; only then adopt remaining resources.
- Declare the nine existing OpenSearch/OpenChoreo recovery assets explicitly under their current owning Applications. The retired workloads stay removed. Namespaces, the readiness snapshot, backup credentials and the retired volume’s existing encrypted backup remain protected with Prune/Delete/Replace/Force=false. This resolves actual missing desired-state ownership without deleting data or using IgnoreExtraneous. Any disposal is a separate reviewed retirement.
- Correct the isolated STAC fixture’s restricted-setting observation and UTF8 initialization before native reproduction. No production migration is authorized by the fixture alone; use its comparison evidence to choose a remedy.

## Node reduction

Current requests total approximately 47.2 GiB before new launches. Removing a 16 GiB core node leaves about 54.7 GiB allocatable across two core plus one ecosystem worker, before new apps and maintenance headroom. NATS requires three separate hosts, so three total nodes leave no spare host during disruption. Root Terraform defaults do not match the live core pool; do not apply a count-only edit or change the ecosystem pool below its required minimum. Keep four nodes while recovering and measuring these services.

## Acceptance

Source preparation, isolated reproduction, Argo sync, ready workloads and application behavior are separate gates. Record deployment evidence below after each actual rollout. Existing account login and historical Matrix message/media acceptance remain human checks; do not send messages on the user’s behalf.

## Launch implementation

- CloudTAK: all 226 pinned migrations and all four images passed isolated startup against the exact PostgreSQL digest with AGE preloaded. Add only the existing sealed cloudtak managed role, with no superuser/createdb/createrole, before the manual first sync. Dependency waves order database, API, workers and private ingress. TAK first-run integration remains a separate observed gate.
- Karakeep: one web/SQLite writer, one Chrome and one Meilisearch pod; private access and an app-only Authentik provider restricted to the stable existing owner UID. Password authentication is disabled. Initial OIDC-only enrollment requires DISABLE_SIGNUPS=false under that sole-user provider; tighten it after the owner's first login. Daily encrypted whole-volume snapshots preserve crash-consistent recovery; initial export/restore remains to be observed.
- Matrix: actual adoption exposed a pre-existing RollingUpdate/RWO deadlock: the replacement pod was scheduled on another node while the old Ready pod retained its media volume. Use Recreate for the single-instance Synapse Deployment to release that attachment before starting its replacement. Preserve all claims, Service identities, config bytes and databases.

The ten Karakeep launch contracts, all 27 rendered resources, both installed VolSync schemas and five sealed-secret controller checks passed. Independent review found no actionable Karakeep launch issues. CloudTAK's 21 schemas and launch YAML lint passed.

## Observed acceptance progress

PRs [#458](https://github.com/onemindos/omos/pull/458), [#459](https://github.com/onemindos/omos/pull/459) and [#460](https://github.com/onemindos/omos/pull/460) passed all GitHub checks and merged. Matrix is Synced/Healthy after the Recreate correction and ordinary Redis/PostgreSQL init-budget adoption; its original claims and Services will receive a final identity check. Karakeep is Synced/Healthy with all three pods Ready. Five retained-recovery/Argo ownership issues are resolved without deleting storage.

The corrected STAC control and AGE-preloaded Jobs both completed natively on amd64 nodes, each passing full migration plus idempotent rerun with no container restarts, SIGSEGV, PANIC or OOM. Raw failed diagnostic evidence and successful logs are saved privately in this task's work/stac-native-before and work/stac-native-after. Production STAC remains manual pending the new explicit preflight CNPG recovery point.

Zeus authorized discovering CloudTAK credentials from TAK. Existing LDAP had no user accounts and the human File password was a nonrecoverable hash. A dedicated LDAP login for the already-existing nonhuman CloudTAK certificate principal was created with its existing OneMind group semantics and Team Member role, without TAK admin privileges. Initial password change followed the unchanged LDAP password policy. Only the final username/password are strictly sealed in the app namespace; no human password, File user, global policy or TAK CoreConfig was changed. Preserve restricted offline recovery material and rotate through normal LDAP self-service before the existing 30-day policy expiry.

## Production STAC and initialized Karakeep recovery verified

On 2026-10-25, preflight CNPG backup `stac-repair-preflight-20261025` completed at 17:53:28 UTC. The one observed production STAC migration completed in eight seconds and reported PgSTAC 0.9.12. Both API pods are Ready. Private landing/collections/search return 200, Browser returns 200, and the exposed reader rejects a POST to collections with 405. All three shared PostgreSQL pod UIDs and container restart counts match the original task baseline; CNPG stayed 3/3 healthy.

Karakeep's initialized encrypted data and Meilisearch backups completed at 18:01:23 UTC. The data snapshot contains six files/807624 bytes; a local read-only Restic restore of both SQLite databases and their matching WAL files passed integrity checks (36 app tables, three queue tables). Meilisearch's five-file/110598-byte snapshot is nonempty. No production PVC or remote archive was modified by the restore check. Remove the accepted manual tags to restore daily cron and verify nextSyncTime. Full owner/bookmark/attachment/search application restore remains a broader human acceptance check.

Matrix's original three PVC UIDs/PV bindings and four Service UIDs/clusterIPs are unchanged. The obsolete plaintext synapse-config ConfigMap was removed only after checking no current pod references it, the replacement is byte-identical, and the new Synapse is Ready. The client versions endpoint still responds successfully.

PR [#463](https://github.com/onemindos/omos/pull/463) passed all four checks and
merged. VolSync manual tags are absent: the next data backup is scheduled for
2026-10-26 12:40 UTC and the next Meilisearch backup for 13:00 UTC.

After launch, node descriptions total 51,916 MiB (approximately 50.7 GiB) of
reserved memory. All four nodes remain Ready and all 80 claims Bound. The
three-node target would have little remaining reservation headroom before
maintenance, and NATS still requires three distinct hosts. No Terraform node
count or live pool size was changed. OpenProject, LibreChat and Vaultwarden
were checked for deployed workloads/services and are absent.

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

# Matrix drift and gated reconciliation — 2026-10-23

Task 4: inspected read-only on `do-nyc3-omos-cloud`. Matrix is currently Healthy,
with 11 Argo resources OutOfSync and no automated sync. This PR preserves that
manual policy. It corrects public domain metadata, adopts the existing ingress
routes, and protects the data claim. **Do not fully sync yet:** the existing
Secret prerequisite has not been met.

## Findings and per-resource decisions

The live Argo source is `platform/matrix/argo`, main, with ServerSideApply=true.
Git was staged for adoption on 2026-10-10 but intentionally did not finish the
secret conversion. `Secret/synapse-homeserver` is absent; the running Synapse
still mounts `ConfigMap/synapse-config`. The missing `synapse-log-config` is a
planned non-secret split from that existing ConfigMap. No secret values were
printed, copied into Git, or changed during inspection. Only resource identities,
public Element config, non-secret log config and allowlisted server names were
examined for reporting.

| Argo resource | Observed functional difference | Decision |
| --- | --- | --- |
| ConfigMap/element-web-config | Git points at matrix.onemindos.com / server onemindos.com; live uses matrix.onemindos.dev / onemindos.dev | Live wins; update Git to the working identity |
| ConfigMap/synapse-log-config | Missing live; log.config is currently in synapse-config, semantically equal to Git | Git split is correct; create only after sealed-config prerequisite is met |
| PVC/synapse-data | Same 50 GiB request, storage class and access mode; no functional drift | Keep existing claim; add explicit no-prune/no-delete/no-replace/no-force annotations |
| Service/element-web | Desired selectors and ports match; live adds allocated/default service fields | Git correct; adopt without replacing the Service |
| Service/postgres | Desired selectors and ports match; allocated/default fields are live | Git correct; preserve clusterIP and identity |
| Service/redis | Desired selectors and ports match; allocated/default fields are live | Git correct; preserve clusterIP and identity |
| Service/synapse | Desired selectors and ports match; allocated/default fields are live | Git correct; preserve clusterIP and identity |
| Deployment/element-web | Image, resources and config mount match; API supplies termination defaults | Git correct; no template change required by this PR |
| Deployment/synapse | Git projects absent Secret/synapse-homeserver plus absent synapse-log-config; live mounts synapse-config | Git secret-backed design is correct **only after** sealing the current configuration; block workload sync until then |
| StatefulSet/postgres | Image, requests/limits, init command, storage template match; live adds API defaults/template status | Git correct; no replace, storage edit, database reset or forced restart |
| StatefulSet/redis | Image, requests/limits, init command, storage template match; live adds API defaults/template status | Git correct; no replace, storage edit or forced restart |

Functional comparisons checked Git-specified fields recursively against the live
objects, distinguishing API default fields from meaningful differences. Matching
objects still await Argo ownership/adoption metadata; OutOfSync alone does not
mean every one of the 11 needs an application behavior change.

Two existing, currently untracked Ingresses are added to the Argo directory:
`matrix-ingress` (matrix.onemindos.dev and element.onemindos.dev) and
`matrix-wellknown` (onemindos.dev/.well-known/matrix). Their live selectors,
routes and annotations are preserved. This does not reset server_name, change
DNS/TLS, enable federation or introduce a new public route.

### Existing PR #192

[#192](https://github.com/onemindos/omos/pull/192) is directly related and remains
open at `fb397afe09d5d5f2d48a9bea470f97a43e245217`. Its description records that the
.dev move, server identity change and a database reset were performed on
2026-10-12. Those are historical claims, not operations to repeat. Current live
Element/Ingress/allowlisted Synapse server_name confirm the .dev identity.
This PR carries only the live public-config/route parity and safe adoption work;
#192 also changes catalog and other docs/config files. Do not merge both blindly:
Claude must reconcile overlapping files and preserve useful remaining changes.
No remote comment, PR closure or merge was performed.

## Data guard

Observed claim `matrix/synapse-data`:

- UID: `69d8e9ca-cd95-4bcc-91ac-13f7b85c55b5`.
- Bound PV: `pvc-69d8e9ca-cd95-4bcc-91ac-13f7b85c55b5`.
- Request: 50 GiB, ReadWriteOnce, `do-block-storage`.
- PostgreSQL and Redis claims are separate: `data-postgres-0` (20 GiB) and
  `data-redis-0` (10 GiB); preserve them as well.

PVC Argo sync options explicitly disable Prune, Delete, Replace and Force.
These guards reduce Argo lifecycle risk; they cannot prevent an administrator
from deleting resources outside the reviewed procedure. Never use force/replace
or delete/create to resolve PVC drift. If a preview proposes recreation, stop.
A current verified backup/restore record is required before changing workloads;
this review did not create a backup or certify recovery.

## Exact reviewer procedure after merge

1. Keep `platform-matrix` manual; no full app sync, no automated prune, no
   namespace or PVC deletion. Record current claim UIDs/PV names, resource UIDs,
   images, pod readiness and the Argo Git revision with read-only kubectl.
2. Read `platform/matrix/SECRETS-TODO.md` step 1. Zeus seals the **current live**
   homeserver.yaml through a pipe directly into a SealedSecret; no rotation,
   plaintext file or terminal output. Use context `do-nyc3-omos-cloud`, name
   `synapse-homeserver`, namespace `matrix`, controller
   `kube-system/sealed-secrets-controller`. Commit only
   `platform/matrix/argo/synapse-homeserver.sealedsecret.yaml` in a separate
   reviewed PR. No placeholder Secret or speculative ciphertext is included here.
3. After that credential PR merges, review Argo's exact-main diff and selectively
   sync **only** the SealedSecret. Do not select Synapse, its PVC, StatefulSets or
   any other resource. Confirm the controller reports Synced and the resulting
   Secret exists with key homeserver.yaml, without displaying its value. Confirm
   a secure local equality check of the sealed source and resulting config; do
   not put either in logs. If absent or mismatched, stop while old Synapse serves.
4. Review the new exact-main diff. No server_name/domain rollback, immutable PVC
   change, StatefulSet claim-template change, Force, Replace or prune is permitted.
   If #430 (image pins) or a resource-budget PR has merged, include its changes in
   a fresh review; this recorded diff is no longer the whole proposal.
5. Selectively sync the public ConfigMap, both existing Ingresses, and
   synapse-log-config. Ensure the current log.config is semantically identical
   before creating the split log ConfigMap. Capture and preserve Service UIDs and
   clusterIPs while adopting the four existing Services. Apply the PVC protection
   annotation only after checking its UID/PV/size against the recorded baseline.
6. Once the verified backup and secret gates pass, selectively sync Synapse with
   ordinary server-side apply. The existing Recreate/rolling strategy is not
   modified here. Observe startup and readiness before considering any remaining
   adoption; do not sync/restart Postgres or Redis just to clear a status badge.
7. Confirm the original claim UIDs/PVs and 50/20/10 GiB allocations are unchanged.
   Check `https://matrix.onemindos.dev/_matrix/client/versions` and Element's
   config.json. Zeus logs in with the existing account and sends a unique test
   message to Hermes from Matrix, confirming the reply and pre-existing room/media
   history. The agent did not send a message or spend model tokens in this review.
8. Confirm Argo Synced/Healthy at the exact reviewed revision only after all
   remaining adoption steps are reviewed. Leave automated sync disabled. Do not
   delete the old synapse-config, PVCs or backups as part of this PR.

On failed Synapse readiness, stop further syncs. Preserve the old ConfigMap and
all storage; restore the prior config-volume reference through a reviewed Git
revert and manual ordinary Argo sync. Do not reset the database, change
server_name, reseed users or roll back disk data to fix configuration drift.

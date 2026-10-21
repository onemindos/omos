# Catalog inventory reconciliation — 2026-10-20

Follow-up to the authorized portal stabilization, not an application migration.
The existing transactional refresh ran on a clean branch from main against the
expected DOKS context, current Tailscale peers, and the authorized GitHub organization.
All three discovery sources completed. This changes catalog metadata, not workloads,
permissions, secrets, storage, workflow planes or application deployment ownership.

## Reviewed invariants

- 176 Deployment/StatefulSet identities across 40 namespaces; 35 observed mesh
  endpoints; 83 repository records. These are timestamped observations, not uptime,
  device ownership, application health or proof every endpoint is reachable.
- Catalog source grows from 555 to 664 entities: **109 added, zero removed**
  (20 Systems, 68 Components, 21 Resources). Existing stable references are retained.
- 22 previously discovered identities are absent from this observation and marked
  `missing-review-required`, retaining their last successful observation time.
  They were not deleted or silently marked current/healthy. This includes the
  retired original STAC Browser workload; the OpenChoreo pilot has its own identity.
- No new unresolved owner, system or dependency references; no duplicate source
  identities. Both production/source copies of apps.yaml and cluster.yaml match.
- Argo CD/Tailscale overlapping entries are retained for reference/favorites
  compatibility; the two external entries now explicitly annotate their canonical
  service reference and depend on it. The launcher already groups matching URLs.
  These are not two
  distinct live installations. Bambu Lab and Bambu Studio remain distinct concepts.
  No canonical identity or personal favorites were deleted during this refresh.
- Historical archived Fabric documents and their migration manifest are preserved;
  regenerated historical-document changes were excluded from this inventory PR.
- No proposal branch `origin/automation/onemind-inventory-review` exists locally.
  The refresh therefore compares retention with main's committed inventory; it
  does not claim accumulated absence counts from unavailable prior proposals.

## Validation and live acceptance

Source-reference/retention/pair checks, changed YAML lint and Kustomize rendering
passed. Secret and rendered-manifest gates are run before PR publication.
Repository-wide make validate still encounters existing unrelated YAML lint debt;
GitHub jobs still cannot execute due the existing account billing restriction.

After merge, Argo regenerates the catalog ConfigMap and performs its normal portal
rollout. Verify two ready replicas on the approved stabilization image, catalog
processing, source origins/counts, owner/dependency resolution, home/favorites and
machine read-only permissions. Do not equate the source file change with processed
live acceptance. The four inactive default builders still reference a genuinely
absent workflow plane; workflows remain disabled, and no plane is fabricated.


### Visible absence qualifiers

The follow-up release annotates descriptions of the same 22 missing-review records
with a human-visible qualifier, preserving every ID, URL and observation timestamp.
The refresh operator applies it idempotently in future; rediscovery reconstructs
the entity from the new observation. This avoids search/About implying that a
retained historical workload is still a current deployment. No entity is retired.

# Portable release and lifecycle pipeline

Implemented for the **portal profile first**; additional apps require a reviewed
profile/build/acceptance adapter. This is an operator-run local/Linux pipeline,
not a new in-cluster CI controller, cron, billing workaround or unrestricted AI
deployment account. “Agno” in the request is interpreted as **Argo CD**.

```text
committed source → required local gates → clean Linux build → GHCR digest proof
       → GitOps image-pin PR → reviewed normal merge → Argo reconciliation
       → exact revision/image/replica checks → browser/app/security acceptance
```

Zeus explicitly approved the verified local/Linux executor while GitHub-hosted
jobs cannot start because of account payment/spending restrictions. This does not
make those jobs green, waive required gates, permit administrator merge bypasses
or change billing. The existing hosted workflow remains intact and runs the same
pipeline safety contracts when execution is available.

## Invariants

- **Git is intent; Argo is the workload writer.** No imperative workload apply,
  patch, rollout restart, forced sync, automatic rollback, Helm install/upgrade,
  prune approval or deletion is executed by this pipeline.
- A checked-in profile pins the repo, context, Argo Application/source/destination,
  namespace, workload/container, registry and replica count. Unknown apps fail
  closed. First profile is `portal`; do not use it as a general cluster-admin API.
- Refuse a dirty checkout, unapproved origin, stale image/source, ambiguous image
  replacement, wrong registry/platform, another pipeline lock, or reused run output.
- The repository-wide advisory lock serializes these pipeline commands across
  linked worktrees. Other tools/operators must still coordinate; it is not a
  distributed GitHub/cluster lock. Later superseding desired state fails review
  rather than being silently rolled back.
- Build only tracked committed source; existing `release-image.sh` performs the
  Linux gates. Registry labels, immutable digest and actual architecture are checked;
  labels are not a claim of cryptographically signed provenance.
- Never read live Secrets or emit credentials. Logs/receipts use private temporary
  directories outside Git. Receipts record stage outcomes, source, image, old image
  and PR/revision evidence. Do not commit raw run files or paste log contents blindly.
- Separate a healthy rollout from app correctness. `verify` does **not** claim SSO,
  token expiry, MCP isolation, catalog topology, STAC query semantics or restoration
  passed. Keep their explicit live/browser evidence and phase gates.

## Prerequisites

Use Python 3.9+ with PyYAML, Git, gh, Docker/BuildKit, kubectl, Helm, kubeconform,
gitleaks, ShellCheck, Terraform and the existing approved authentication. Portal
source checks need Node **22 or 24**, committed Yarn 4.13.0 and network/cache access.
The pipeline does not install runtimes, change kube context defaults, copy/rotate
tokens, start a cloud VM or repair billing. Keep builds off production K8s nodes.
`publish` requires at least 8 GiB host free before starting; this is a minimum
guard, not a guarantee that every Docker/VM storage configuration has enough space.

The old TechDocs `npx` command was version-pinned at its root but re-resolved its
transitive graph every build. An actual Linux rebuild failed with npm `ETARGET`
for `@backstage/plugin-search-common@^1.2.26`. `@techdocs/cli` 1.12.0 is now a
workspace dev dependency in the immutable Yarn lock; the build executes that
installed CLI and still requires TechDocs metadata and search artifacts. This
removes that floating resolution, not every external apt/registry availability risk.

## Commands

Run from `/Users/zeus/omos` (or the approved Linux checkout). Each invocation
creates fresh private evidence and prints its path. `--output` may select a new
directory **outside** the repo. No daemon is installed.

```bash
# Validate the committed PR range; no cloud/workload mutations.
python3 scripts/gitops/pipeline.py check --profile portal --base origin/main

# Validate, build Linux/AMD64, push and verify; no deployment yet.
python3 scripts/gitops/pipeline.py publish --profile portal --base origin/main

# Revalidate a previously emitted real release.json; never accept a digest alone.
python3 scripts/gitops/pipeline.py artifact --profile portal --release /tmp/RUN/image/release.json

# Only after accepted source is merged: create a separate pin PR from origin/main.
# The source/build-workflow bytes must still match the image's source commit.
python3 scripts/gitops/pipeline.py promote --profile portal \
  --release /tmp/RUN/image/release.json --branch codex/portal-verified-promotion

# Review/normal merge is a separate approved operation. No --admin/force bypass.
# Then verify the exact merged pin revision, digest, generation and all replicas:
python3 scripts/gitops/pipeline.py verify --profile portal \
  --release /tmp/RUN/image/release.json --revision MERGED_40_CHARACTER_COMMIT --timeout 600

# Read-only retirement preflight: inventories resources/PVCs, never removes them.
python3 scripts/gitops/pipeline.py retirement-plan --profile portal
```

After `promote`, the agent must attach the created PR to its chat with the normal
Codex artifact tool. If opening the PR was interrupted, inspect the branch/remote
before retrying; do not overwrite a branch or create a duplicate promotion.

### Required validation versus known debt

`check` requires changed-YAML lint, manifest schemas, Helm renders, tracked shell
errors, Terraform formatting/init-without-backend/validation, lifecycle tests,
image-command contracts, commit-range and exact tracked-tree secret scans. Portal
source changes also require immutable install, TypeScript, lint, identity/permission/
STAC tests, layout tests and catalog audit tests. `publish` additionally runs the
actual Dockerfile's Linux tests, TechDocs/backend build and final runtime artifacts.

`make validate` is also recorded as a **diagnostic**. Existing whole-tree YAML debt
starting at `platform/onetak/media-mtx.yaml:86` is not called passing; new changed
YAML and all required checks still stop the pipeline. An unrelated diagnostic
failure is not permission to skip a new required failure. Read each run receipt.

### Promotion and reconciliation

`promote` verifies the release record against the live registry and requires its
source to be merged into main with unchanged app/build-workflow bytes. It creates
one new branch/commit changing **only the unique profile container image**, runs
required gates, pushes and opens a PR. It does not merge it or mutate the roadmap;
link the promotion/live evidence in the owning phase tracker before closing it.

Existing `platform-backstage` automated sync detects merged desired state. The
pipeline does not add another writer or change its prune/self-heal options.
`verify` requires Synced/Healthy at the pin revision, no running/failed Argo
operation, current Deployment generation, the approved replica count, all updated/
available/ready replicas, matching nonterminating pods and actual runtime image IDs.
For an OCI index, the remotely verified target-platform manifest is also accepted.
A later main revision is accepted only if it descends from the requested revision
and both the entire owned GitOps path and its Application are unchanged.

Timeout/unknown/partial state is failure, not a reason to force-sync or restart.
Inspect Argo diff/operation and pod events. A failed deployment must not trigger
an automatic image rollback over newer desired state or a database schema change.
Prepare a reviewed Git revert with the app's compatibility/recovery evidence.

## Safe removals: deliberate retirement, never a blanket delete button

The implemented `retirement-plan` is **read-only planning**, not a removal executor.
Portal is a protected foundation profile. Even an unprotected future app does not
receive automatic deletion from this command. It records Argo's actual resource
inventory and namespace PVCs to retain; unavailable inventory is not an empty app.

Before implementing/executing an app-specific retirement PR:

1. Confirm exact resources, owning controller, dependencies, catalog/launcher links,
   external DNS/LBs, traffic and data export/retention policy.
2. Obtain a fresh app-consistent backup and isolated restore for stateful resources;
   preserve recovery/encryption keys independently. State age/RPO/RTO explicitly.
3. Get explicit resource-level removal approval. Remove routes/traffic and quiesce
   writers in reversible steps before removing workload intent.
4. Review Argo prune, Application/ApplicationSet finalizers, parent app-of-apps,
   hooks and retention annotations. Deleting an Application may orphan resources
   or cascade-delete them depending on configuration; neither is assumed safe.
5. Preserve Namespace/PVC/PV/Secret/SealedSecret resources and backup data unless
   separately approved with recovery evidence. Never toggle global prune/allowEmpty
   to make a retirement succeed. Argo, not the operator script, executes approved
   desired-state changes; verify absence/retention after the approved sync.

No retirement was executed during the first portal release. Adding an app profile,
Helm/OpenChoreo lifecycle adapter, semantic checks, signing, notification routing
or scheduling is a separate tested PR, not an implied universal capability.

## References

- [Argo automatic sync](https://argo-cd.readthedocs.io/en/stable/user-guide/auto_sync/)
- [Prune/delete confirmation and sync options](https://argo-cd.readthedocs.io/en/stable/user-guide/sync-options/)
- [Application deletion and finalizers](https://argo-cd.readthedocs.io/en/stable/user-guide/app_deletion/)
- [Existing portable image release](../../services/onemind-portal/scripts/RELEASE.md)

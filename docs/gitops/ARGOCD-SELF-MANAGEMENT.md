# Argo CD managing itself

Argo CD is the one adopted-Helm candidate with no Argo Application yet (see
`HELM-ADOPTION.md`). `platform/argocd/platform-argocd.yaml` defines it. This page is the
plan, the risks and the exact steps. **Merging the file applies nothing**: the Application
has no automated sync and is only synced by a person.

## What the Application is

| Setting | Value | Why |
|---|---|---|
| Chart | `argo-cd` 10.1.2 from `https://argoproj.github.io/argo-helm`, release name `argocd` | Same chart and version as the live Helm release (revision 3, app v3.4.4) |
| Values | `$values/platform/values/argocd.yaml` | The file Git already holds |
| Inline overrides | `crds.install=false`, `configs.secret.createSecret=false`, `notifications.secret.create=false` | Safety switches, see below |
| Sync | **manual** (no `automated`) | The first sync must be a reviewed step. Intended end state: `automated: {prune: false, selfHeal: true}` in a later PR |
| Prune | never (do not pass `--prune`, do not tick Prune in the UI) | Argo must never delete part of itself |
| Finalizer | none | Deleting the Application can never cascade to Argo CD |
| Sync options | `ServerSideApply=true`, `SkipDryRunOnMissingResource=true`, no `Replace=true` | Adopt by field ownership; never replace |

Inline overrides, and what they leave alone:

- `crds.install=false`: the `applications`, `applicationsets` and `appprojects` CRDs stay
  out of Argo's sync scope. They are chart templates with `helm.sh/resource-policy: keep`.
  Change them only by hand on a chart bump, before the bump:
  `helm template argocd argo/argo-cd --version <new> -n argocd --set crds.install=true | kubectl apply --server-side --force-conflicts -f -`
  limited to the `CustomResourceDefinition` documents.
- `argocd-secret` and `argocd-notifications-secret` are not rendered. Argo CD writes the
  live ones itself (server signing key, admin password hash, webhook secrets); the chart's
  render holds no data and would add nothing but risk. The live Secrets stay as they are.

## Read-only comparison done for this PR

`helm template` of the pinned chart with `platform/values/argocd.yaml`, compared object by
object with `helm get manifest argocd -n argocd` (labels `helm.sh/chart` and
`app.kubernetes.io/version`, and `checksum/*` annotations ignored):

- 49 objects exist in both. 44 are identical. Two are the chart's Secrets (see the inline overrides).
  The other three are ConfigMaps that differ only in the keys that PR #392
  changed and the live release does not have yet: `argocd-cm` `url` (live is
  `https://argocd.example.com`) and `oidc.config`, `argocd-cmd-params-cm`
  `server.dex.server` and `server.dex.server.strict.tls`, `argocd-rbac-cm` `policy.csv`.
- 6 objects are live-only: the Dex Deployment, Service, ServiceAccount, Role, RoleBinding
  and NetworkPolicy (all `argocd-dex-server`). Git disables Dex.
- 4 objects are render-only: the `argocd-redis-secret-init` Job, ServiceAccount, Role and
  RoleBinding. They are Helm hooks, which `helm get manifest` omits.
- Live values differ from Git only in the `configs.cm`, `configs.rbac` and `dex` keys that
  #392 added (the pending upgrade); nothing else.

## Relationship to the pending `helm upgrade` from #392

#392 (merged) changed `platform/values/argocd.yaml` (Authentik sign-in, `url`, RBAC, Dex off)
but its `helm upgrade` has **not** been run. **Run it first, then adopt.**

If the Application were synced first, Argo CD would itself roll out those changes (restarting
`argocd-server` while someone may be signing in with a not-yet-verified `oidc.config`) and,
because prune is off, would leave the six Dex objects behind as orphans. Helm removes them
cleanly and gives `helm rollback argocd 3` as the undo. After the upgrade, the first Argo
sync is a no-op, which is what makes adoption verifiable.

## Adoption risks

1. **Helm ownership metadata.** Live objects carry `meta.helm.sh/release-name: argocd`,
   `meta.helm.sh/release-namespace: argocd` and `app.kubernetes.io/managed-by: Helm`. Argo
   leaves those in place and adds its tracking annotation. Consequence: Helm still believes it
   owns the release. After adoption never run `helm upgrade`, `helm rollback` or
   `helm uninstall argocd` (uninstall deletes every object in the release record, including
   Argo CD itself). The `sh.helm.release.v1.argocd.v*` Secrets stay as history.
2. **Server-side apply field ownership.** The first sync makes the `argocd-controller`
   field manager co-own fields with `helm`. With `Replace=true` or client-side apply this
   could drop live-only fields; it is not used.
3. **Self-restart.** Syncing a change to the controller StatefulSet restarts the
   controller that is running the sync. This is normal; the operation resumes. Keep a
   second way in (below).
4. **Bad `argocd-cm`/RBAC.** A wrong `oidc.config` or `policy.csv` can lock users out.
   `admin.enabled: "true"` is kept as break-glass until SSO is proven.
5. **Hooks.** `argocd-redis-secret-init` runs as a PreSync hook on every sync. It only
   creates the `argocd-redis` Secret when missing.
6. **Tracking collisions.** Argo CD 3 tracks by annotation (default). Do not set
   `application.resourceTrackingMethod: label`: the chart's `app.kubernetes.io/instance: argocd`
   label would make the Application own every object that carries it.
7. **Bootstrap.** `platform/install.sh` still runs `helm upgrade --install argocd`. On a rebuild
   that is the bootstrap step; run it only on an empty cluster, and never on this one.
8. **CRD upgrades are manual** (see above) and must precede any chart bump.

## Exact safe adoption steps

Nothing here is run by this PR. Every step below is for a person; stop at the first
unexpected result.

1. **Pre-checks (read-only).**
   - `helm -n argocd list` shows `argocd` revision 3, chart `argo-cd-10.1.2`, `deployed`.
   - `kubectl -n argocd get secret argocd-oidc-authentik argocd-secret` both exist
     (names only; do not print them). The Authentik application `argocd` exists
     (`platform/ecosystem/resources/argocd-authentik-provider.yaml`).
   - `kubectl -n argocd get applications,appprojects -o yaml > <file outside the repo>` as an export.
   - You can sign in with the local `admin` account (break-glass).
2. **Apply #392 with Helm** (the pending upgrade), in a quiet window:
   `helm upgrade argocd argo/argo-cd --version 10.1.2 -n argocd -f platform/values/argocd.yaml --atomic --timeout 10m`.
   Expect the Dex objects to be deleted and `argocd-server` and `argocd-repo-server` to
   restart. Verify: SSO sign-in through Authentik works in a private window, an
   `argocd-admins` member is admin, `admin` still works, `kubectl -n argocd get applications`
   still shows everything Synced. Undo with `helm rollback argocd 3 -n argocd`.
3. **Merge this PR.** `app-of-apps` creates `platform-argocd`; it does not sync.
4. **Review the diff before any sync.** `argocd app diff platform-argocd` must show **no change**
   to any Deployment, StatefulSet, Service, ConfigMap, Role or ClusterRole. Expected output is
   empty or limited to metadata (tracking annotation, managed fields) and the PreSync hook.
   Any pod-template change means stop and fix Git or the values, not the cluster.
5. **First sync, by hand, without prune or replace.** `argocd app sync platform-argocd`
   (or Sync in the UI with Prune and Replace unticked). Check that no Argo CD pod restarted
   (`kubectl -n argocd get pods`: ages and restart counts unchanged), `argocd-redis-secret-init`
   completed, and the Application is `Synced` and `Healthy`.
6. **Soak** a few days, including a change to one `platform/values/argocd.yaml` resource
   value through a PR to see Argo apply it to itself.
7. **Enable automation** in a follow-up PR: add
   `automated: {prune: false, selfHeal: true}` to `platform-argocd.yaml`, mark `argocd` as
   Adopted in `HELM-ADOPTION.md`, and update `docs/STACK.md`.

## Break-glass

If Argo CD is down or mis-synced: `kubectl -n argocd rollout undo` of the broken Deployment or
StatefulSet, or `kubectl -n argocd rollout restart`; sign in with the local `admin` account
(`admin.enabled: "true"`); as a last resort re-run `helm upgrade` with the last good values
from Git history (this reverts to the Helm-owned state described above). Applications keep
running without Argo CD; nothing in the cluster depends on it being up.

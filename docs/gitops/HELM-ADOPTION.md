# Adopting Helm-CLI releases into Argo CD

Eleven releases on `do-nyc3-omos-cloud` were installed with `helm upgrade --install`
(see `platform/install.sh` and `docs/STACK.md`). This document is the pattern for moving
each one under Argo CD with zero functional change, and the status of every release.

## The rule

**Never run `helm upgrade`, `helm install` or `helm uninstall` on an adopted release.**
After adoption Argo CD is the only writer. A CLI upgrade creates a second Helm revision
whose manifest drifts from Git, and Argo (self-heal) will then fight it. To change an
adopted release, edit its values file or chart version in Git and merge. Read-only
commands (`helm get values|manifest|history`, `helm list`, `helm template`) are fine.

The existing Helm release Secrets (`sh.helm.release.v1.*`) stay in place as history.
Argo renders with `helm template` and applies the output itself; it does not touch them.

## The pattern

One Application per release, `platform/argocd/<group>-<release>.yaml`, picked up by
`app-of-apps`. It is a multi-source Application: source 1 is the chart (pinned repo URL
and exact version, same `releaseName` and namespace as the live release), source 2 is
this repository with `ref: values`, so the values file is read from
`$values/platform/values/<release>.yaml`.

| Setting | Why |
|---------|-----|
| `releaseName` = live release name, `destination.namespace` = live namespace | Names and labels render exactly as the live objects, so the first sync is a no-op adoption |
| `targetRevision` = exact live chart version | Zero functional change; upgrades are a later, reviewed PR |
| `prune: false`, `selfHeal: true` | Argo never deletes anything; drift is corrected |
| no `resources-finalizer` | Deleting the Application can never cascade to workloads, PVCs or CRDs |
| `ServerSideApply=true`, `compare-options: ServerSideDiff=true` | Adopts existing objects by field ownership; no giant last-applied annotation (matters for CRDs) |
| `SkipDryRunOnMissingResource=true` | Allows a release and its CRs to sync in one pass |
| `retry` (10, 15s x2, max 3m) | Rides out webhook/API start-up races |
| `ignoreDifferences` + `RespectIgnoreDifferences=true` | For fields filled at runtime, mainly webhook `caBundle` |

### Adoption checklist per release (all read-only against the cluster)

1. `helm get values <rel> -n <ns> -o yaml` must equal the repo values file semantically.
   Fix the repo file to match live, not the other way round.
2. `helm template` the pinned chart with the repo values, then compare to
   `helm get manifest <rel> -n <ns>` per resource (kind, namespace, name), ignoring
   `helm.sh/chart`, `chart`, `app.kubernetes.io/version` labels and `checksum/*` annotations.
   The only expected additions are Helm hook resources, which `helm get manifest` omits.
3. Search the chart for `genCA`, `genSelfSignedCert`, `randAlphaNum`, `lookup`. Anything that
   renders differently on every sync must be neutralised (existing-secret option,
   `ignoreDifferences`, or a value pinned in a Secret) before adoption.
4. CRDs: note whether they come from the chart's `crds/` dir (`skipCrds: true` excludes
   those) or from `templates/` (`skipCrds` does **not** exclude them). Never allow prune or
   replace on CRDs; prune is off and ServerSideApply is used, never `Replace=true`.
5. Secrets: reference existing Secrets by name only. If a chart would render a Secret whose
   real value lives elsewhere (SealedSecret, ExternalSecret, out-of-band), unset the
   chart input so the Secret is not rendered at all. Never commit a secret value.
6. Run `python3 scripts/validate_helm_renders.py --all`, `python3 scripts/validate_manifests.py`
   and `yamllint -c .yamllint.yml` on the touched files.

### Helm hooks under Argo

Argo maps `helm.sh/hook` pre-* to PreSync and post-* to PostSync and runs them on each sync
that has changes. ingress-nginx (`ingress-nginx-admission-create` / `-patch` Jobs) and
cert-manager (`cert-manager-startupapicheck` Job) have such hooks. They are idempotent and
self-deleting (`hook-delete-policy`), so they appear as short-lived Jobs during a sync.

## Release status

**8 adopted, 3 remaining** (2026-10-17; the four monitoring releases were adopted after the 4/7 check that morning). All 11 still appear in `helm list -A`, because adoption leaves the Helm release record in place. Argo renders the chart itself and Helm records stop updating.

`platform/install.sh` still runs `helm upgrade --install` for the adopted releases too. Do not run it: it would fight Argo.

| Release | Namespace | Chart | Status | Application | Values |
|---------|-----------|-------|--------|-------------|--------|
| `monitoring` | monitoring | kube-prometheus-stack-87.16.1 | Adopted 2026-10-17 (PR #171 was closed unmerged; redone as adoption only, no alerting changes) | `monitoring-kube-prometheus-stack` | `platform/values/kube-prometheus-stack.yaml` |
| `loki` | monitoring | loki-7.0.0 | Adopted 2026-10-17 (PR #171 was closed unmerged; redone as adoption only, no alerting changes) | `monitoring-loki` | `platform/values/loki.yaml` |
| `tempo` | monitoring | tempo-1.24.4 | Adopted 2026-10-17 (PR #171 was closed unmerged; redone as adoption only, no alerting changes) | `monitoring-tempo` | `platform/values/tempo.yaml` |
| `otel-collector` | monitoring | opentelemetry-collector-0.165.0 | Adopted 2026-10-17 (PR #171 was closed unmerged; redone as adoption only, no alerting changes) | `monitoring-otel-collector` | `platform/values/otel-collector.yaml` |
| `ingress-nginx` | ingress-nginx | ingress-nginx-4.15.1 | Adopted in #176 | `edge-ingress-nginx` | `platform/values/ingress-nginx.yaml` |
| `cert-manager` | cert-manager | cert-manager-v1.20.3 | Adopted in #176 | `edge-cert-manager` | `platform/values/cert-manager.yaml` |
| `tailscale-operator` | tailscale | tailscale-operator-1.98.4 | Adopted in #176 | `edge-tailscale-operator` | `platform/values/tailscale-operator.yaml` |
| `nack` | nats | nack-0.34.0 | Adopted in #176 | `edge-nack` | `platform/nats/nack-values.yaml` |
| `nats` | nats | nats-2.14.2 | Remaining | none | `platform/nats/values.template.yaml` plus gitignored `values.secret.yaml` |
| `takserver` | takserver | takserver-4.6.0 | Remaining | none | `platform/takserver/values-onemind.yaml`, chart `platform/takserver/chart` |
| `argocd` | argocd | argo-cd-10.1.2 | Defined, not yet adopted (manual sync only; see `ARGOCD-SELF-MANAGEMENT.md`) | `platform-argocd` | `platform/values/argocd.yaml` |

### Notes for this PR

- **ingress-nginx**: live values equal the repo file. `crds.enabled` in the file is not a
  chart key for ingress-nginx and is a no-op; it is kept so the file stays identical to the
  live values. Admission webhook caBundle is ignored.
- **cert-manager**: live values equal the repo file. The six CRDs are rendered from chart
  templates (`crds.enabled: true`, each with `helm.sh/resource-policy: keep`) and match
  live; they stay chart-managed, server-side applied, prune off. Webhook caBundles (both
  webhook configurations) are injected by cainjector and ignored.
- **tailscale-operator**: the repo file intentionally omits `oauth.clientId`, which the live
  release has. With it set, the chart renders `Secret/operator-oauth` with an **empty**
  `client_secret`; the real Secret is owned by a SealedSecret and Argo would have blanked
  the credential. With it unset the Secret is not rendered and the operator Deployment
  renders identically (it always mounts `operator-oauth`). The Secret is referenced by name
  only. The seven tailscale CRDs come from templates and match live.
- **nack**: live values equal `platform/nats/nack-values.yaml`, including
  `jetstream.nats.url: null`. The controller therefore has no server URL and logs
  reconcile timeouts today (also noted in `docs/STACK.md`). Adoption keeps that as is; the
  follow-up fix must supply the connection through an existing Secret rather than a URL with
  an embedded password. The nack chart has no CRDs (they live in `platform/nats/`).

### Why the last three are last

- **nats** (10 revisions): JetStream StatefulSet with PVCs and an out-of-band
  `values.secret.yaml` holding passwords. `platform-nats` already also manages manifests in
  the same namespace, and the note in `docs/STACK.md` says Argo and Helm both own the
  StatefulSet. Needs: secrets moved to existing Secrets or ESO/SealedSecrets first, a
  VolumeSnapshot of every JetStream PVC before the first sync, and a render-vs-live diff that
  proves no StatefulSet template change (a change rolls the NATS cluster).
- **takserver** (34 revisions, local chart): the most important stateful workload, with a
  CA, Postgres data and certificates. Needs volume snapshots (and a restore drill) first, a
  render diff against the live manifest, and a freeze on `helm upgrade` from the existing
  deploy scripts before Argo takes over.
- **argocd**: Argo cannot safely adopt itself in one step. The Application and the step-by-step plan are in `ARGOCD-SELF-MANAGEMENT.md`; the pending #392 `helm upgrade` goes first. Needs a self-management approach
  (an Application for the `argocd` release with `prune: false`, manual sync only at first,
  `Replace` forbidden, and CRDs kept out of sync scope), a break-glass path that restores
  Argo with `helm` if a bad sync takes it down, and a maintenance window. Do it last so a
  failure cannot strand the other adopted releases.

## Expected behaviour on first sync

Each Application should show `Synced` after the first sync, with no pod restarts: the
render equals the live manifest apart from Helm hook Jobs (which run once) and server-side
apply ownership metadata. If a workload restarts on the first sync, treat it as a defect and
stop: it means the render differs from live.


## TAK partial ownership — 2026-10-26

The existing `platform-takserver` Application now owns the exact existing
`takserver-postgresql` and `takserver-nats` StatefulSets and `cot-retention` CronJob.
The first sync must preserve all three server-defaulted specs byte-for-byte; resource
budgets are a separate follow-up. Namespace, source claims, PVC templates, images,
selectors, credentials, schedules and workload names are preserved. Prune/Delete/
Replace/Force are explicitly disabled for the adopted resources. Existing Helm
revision34 and ownership labels are retained as history.

Do not run whole-release Helm upgrades, rollbacks or uninstalls on `takserver`;
they would compete with these resources and the already owned overlay Deployments.
No repository shell script upgrades that release; legacy commands in the upstream
chart README are not instructions for this live cluster. Full remaining release
adoption is a different task. The old bootstrap script must not be run on this
Argo-managed cluster.

See [recovery preflight and evidence](../operations/2026-10-26-tak-policy-preflight.md)
for fresh encrypted CSI-backed snapshots, PostgreSQL isolated logical restoration
and the bounded internal NATS archive/server test. No source claim is deleted,
recreated or shrunk; the live PostgreSQL claim remains20Gi even though its existing
immutable StatefulSet claim template is8Gi.

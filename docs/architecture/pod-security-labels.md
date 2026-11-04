# Pod Security labels

Pod Security Admission reads three labels on each namespace. `enforce` rejects a Pod that fails
the named profile. `warn` and `audit` never reject anything: `warn` returns a warning to the
client (`kubectl apply`, `helm upgrade`, Argo CD or an operator) and `audit` adds an annotation
to the API server audit event. Both also check workload templates (Deployments, StatefulSets,
Jobs and so on), not only Pods. This repo does not collect the control-plane audit log, so in
practice `warn` is the signal you see.

| Label | Value | Namespaces |
|---|---|---|
| `warn`, `audit` | `restricted` | Every application namespace whose Namespace manifest Argo CD applies, except those listed below |
| `enforce` | `restricted` | `karakeep`, `n8n`, `stac`, `stac-validation`, `vaultwarden` |
| `enforce` | `baseline` | `nextcloud` |

No `*-version` label is set, so each mode uses the profile of the cluster's Kubernetes version.

## Left without warn and audit labels

- `falco`, `velero`, `tailscale`, `monitoring`, `github-runner`: their workloads need privileged
  containers, host namespaces or hostPath volumes by design, so the warnings would never clear.
- `default`: no long-running workloads.
- `kube-system`, `kube-public`, `kube-node-lease`: managed by DOKS, not defined in Git.

## VolSync privileged movers

When a Namespace manifest carries `volsync.backube/privileged-movers: "true"`, VolSync runs that
namespace's backup movers as UID 0 with the `CHOWN`, `DAC_OVERRIDE` and `FOWNER` capabilities.
That passes `baseline` and fails `restricted`. The VolSync controller creates each mover Job at
backup time, so the warning appears only in the controller log, repeated on each reconcile while
the mover runs. It never shows up in Argo CD or kubectl output.

```bash
kubectl --context do-nyc3-omos-cloud get ns -o jsonpath='{range .items[?(@.metadata.annotations.volsync\.backube/privileged-movers=="true")]}{.metadata.name}{"\n"}{end}'
kubectl --context do-nyc3-omos-cloud -n ops logs deploy/platform-volsync | grep 'would violate PodSecurity'
```

## Rules

- A new application namespace gets `pod-security.kubernetes.io/warn: restricted` and
  `pod-security.kubernetes.io/audit: restricted` in its Namespace manifest.
- Raising `enforce` is a separate change for one namespace at a time, made after that
  namespace stops producing warnings. Check the VolSync controller log as well as Argo CD and
  kubectl output, because backup mover warnings appear only there.
- A namespace with the privileged-movers annotation can go no higher than `enforce: baseline`,
  as `nextcloud` does. Before it can take `enforce: restricted`, each of its ReplicationSources
  needs a non-root `moverSecurityContext` and the annotation must be removed, as in `karakeep`,
  `n8n` and `vaultwarden`. Otherwise admission rejects the mover pods and the
  backups stop with no error at sync time; the `PortableVolumeBackupOutOfSync` alert is the
  first sign.
- Any namespace label change, these included, gives every pod in that namespace a new Cilium
  security identity, because Cilium keeps namespace labels in the identity. Pods do not restart,
  but make the change in a quiet window and watch Hubble for policy drops afterwards.

Check the live labels (read-only):

```bash
kubectl --context do-nyc3-omos-cloud get ns -L pod-security.kubernetes.io/enforce,pod-security.kubernetes.io/warn,pod-security.kubernetes.io/audit
```

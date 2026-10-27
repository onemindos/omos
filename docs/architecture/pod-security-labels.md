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
| `enforce` | `restricted` | `karakeep`, `n8n`, `paperclip`, `stac`, `stac-validation`, `vaultwarden` |
| `enforce` | `baseline` | `nextcloud` |

No `*-version` label is set, so each mode uses the profile of the cluster's Kubernetes version.

## Left without warn and audit labels

- `falco`, `velero`, `tailscale`, `monitoring`, `github-runner`: their workloads need privileged
  containers, host namespaces or hostPath volumes by design, so the warnings would never clear.
- `default`, `onemind-search`, `openchoreo-observability-plane`: no workloads run there.
- `kube-system`, `kube-public`, `kube-node-lease`: managed by DOKS, not defined in Git.

## Rules

- A new application namespace gets `pod-security.kubernetes.io/warn: restricted` and
  `pod-security.kubernetes.io/audit: restricted` in its Namespace manifest.
- Raising `enforce` is a separate change for one namespace at a time, made after that
  namespace's workloads stop producing warnings.

Check the live labels (read-only):

```bash
kubectl --context do-nyc3-omos-cloud get ns -L pod-security.kubernetes.io/enforce,pod-security.kubernetes.io/warn,pod-security.kubernetes.io/audit
```

# Gatus status checks

> Generated from the reviewed catalog (workload-map.yaml, structure.yaml). It lists what is recorded, not what is currently healthy; check the Live tab for current state.

Live HTTP(S) checks and certificate expiry for apps that run outside Kubernetes or have no probe of their own. Deployed by Argo application platform-gatus; Deployment gatus/gatus observed 1/1 running on 2026-11-07.

## Where it runs

- Catalog entity: `component:default/gatus` in System `observability`
- Argo CD application: `platform-gatus`
- Namespace: `gatus`

## Workloads

| Workload | Kind | Namespace |
| --- | --- | --- |
| `gatus` | Deployment | `gatus` |

## When it looks wrong

1. Open the app page. The status dot combines Argo CD, workload readiness, alerts and uptime; grey means unknown, not healthy.
2. **Live** tab: confirm pods and replicas. **Logs** tab: open the generated Loki query. **Metrics** tab: open the workload dashboard with this app's variables.
3. Find the cause before acting. A restart that hides the cause is not a fix.

## Safe recovery

- Change the cluster only through Git and Argo CD (pull request, merge, sync). Do not `kubectl apply`, patch or delete by hand, and do not run Helm directly.
- A Git change is reversible with `git revert`; Argo CD then restores the previous state.
- Data stores: confirm the latest backup in the [backups runbook](../backups.md) before any change that touches storage.

## Generated links

- Loki query: `{k8s_namespace_name="gatus"}`
- Grafana dashboard path: `d/onemind-workload/onemind-workload?var-namespace=gatus`
- Source: `onemindos/omos`

## Escalate

If this app is unreachable and the cause is not clear after the Live and Logs tabs, treat it as an incident: record what was observed (time, tab, message) and follow the [resilience runbook](../resilience.md).

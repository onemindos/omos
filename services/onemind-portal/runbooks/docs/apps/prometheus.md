# Prometheus

> Generated from the reviewed catalog (workload-map.yaml, structure.yaml). It lists what is recorded, not what is currently healthy; check the Live tab for current state.

Metrics scraper

## Where it runs

- Catalog entity: `component:default/prometheus` in System `observability`
- Argo CD application: `monitoring-kube-prometheus-stack`
- Namespace: `monitoring`

## Workloads

| Workload | Kind | Namespace |
| --- | --- | --- |
| `alertmanager-monitoring-kube-prometheus-alertmanager` | StatefulSet | `monitoring` |
| `monitoring-kube-prometheus-operator` | Deployment | `monitoring` |
| `monitoring-kube-state-metrics` | Deployment | `monitoring` |
| `prometheus-monitoring-kube-prometheus-prometheus` | StatefulSet | `monitoring` |

## When it looks wrong

1. Open the app page. The status dot combines Argo CD, workload readiness, alerts and uptime; grey means unknown, not healthy.
2. **Live** tab: confirm pods and replicas. **Logs** tab: open the generated Loki query. **Metrics** tab: open the workload dashboard with this app's variables.
3. Find the cause before acting. A restart that hides the cause is not a fix.

## Safe recovery

- Change the cluster only through Git and Argo CD (pull request, merge, sync). Do not `kubectl apply`, patch or delete by hand, and do not run Helm directly.
- A Git change is reversible with `git revert`; Argo CD then restores the previous state.
- Data stores: confirm the latest backup in the [backups runbook](../backups.md) before any change that touches storage.

## Generated links

- Loki query: `{k8s_namespace_name="monitoring", k8s_pod_name=~"(alertmanager-monitoring-kube-prometheus-alertmanager|monitoring-kube-prometheus-operator|monitoring-kube-state-metrics|prometheus-monitoring-kube-prometheus-prometheus)-.*"}`
- Grafana dashboard path: `d/onemind-workload/onemind-workload?var-namespace=monitoring&var-workload=alertmanager-monitoring-kube-prometheus-alertmanager&var-workload=monitoring-kube-prometheus-operator&var-workload=monitoring-kube-state-metrics&var-workload=prometheus-monitoring-kube-prometheus-prometheus`
- Source: `onemindos/omos`

## Escalate

If this app is unreachable and the cause is not clear after the Live and Logs tabs, treat it as an incident: record what was observed (time, tab, message) and follow the [resilience runbook](../resilience.md).

# NATS JetStream

> Generated from the reviewed catalog (workload-map.yaml, structure.yaml). It lists what is recorded, not what is currently healthy; check the Live tab for current state.

The authority message bus (3-node JetStream), the NACK controller and NATS tools.

## Where it runs

- Catalog entity: `component:default/nats` in System `messaging-nats`
- Argo CD applications: `edge-nack`, `platform-nats`, `platform-nats-tools`
- Namespace: `nats`

## Workloads

| Workload | Kind | Namespace |
| --- | --- | --- |
| `nack` | Deployment | `nats` |
| `nats` | StatefulSet | `nats` |
| `nats-metrics-scrape` | Deployment | `nats` |

## When it looks wrong

1. Open the app page. The status dot combines Argo CD, workload readiness, alerts and uptime; grey means unknown, not healthy.
2. **Live** tab: confirm pods and replicas. **Logs** tab: open the generated Loki query. **Metrics** tab: open the workload dashboard with this app's variables.
3. Find the cause before acting. A restart that hides the cause is not a fix.

## Safe recovery

- Change the cluster only through Git and Argo CD (pull request, merge, sync). Do not `kubectl apply`, patch or delete by hand, and do not run Helm directly.
- A Git change is reversible with `git revert`; Argo CD then restores the previous state.
- Data stores: confirm the latest backup in the [backups runbook](../backups.md) before any change that touches storage.

## Generated links

- Loki query: `{k8s_namespace_name="nats", k8s_pod_name=~"(nack|nats|nats-metrics-scrape)-.*"}`
- Grafana dashboard path: `d/onemind-workload/onemind-workload?var-namespace=nats&var-workload=nack&var-workload=nats&var-workload=nats-metrics-scrape`
- Source: `onemindos/omos`

## Escalate

If this app is unreachable and the cause is not clear after the Live and Logs tabs, treat it as an incident: record what was observed (time, tab, message) and follow the [resilience runbook](../resilience.md).

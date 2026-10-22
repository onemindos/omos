# TAK Server

> Generated from the reviewed catalog (workload-map.yaml, structure.yaml). It lists what is recorded, not what is currently healthy; check the Live tab for current state.

TAK 5.8 hardened

## Where it runs

- Catalog entity: `component:default/tak-server` in System `tak-ecosystem`
- Argo CD application: `platform-takserver`
- Namespace: `takserver`

## Workloads

| Workload | Kind | Namespace |
| --- | --- | --- |
| `mumble-bridge` | Deployment | `takserver` |
| `murmur` | Deployment | `takserver` |
| `tak-lb-proxy` | Deployment | `takserver` |
| `tak-ldap` | StatefulSet | `takserver` |
| `tak-nats-leaf` | Deployment | `takserver` |
| `takserver-api` | Deployment | `takserver` |
| `takserver-config` | Deployment | `takserver` |
| `takserver-ignite` | StatefulSet | `takserver` |
| `takserver-messaging` | Deployment | `takserver` |
| `takserver-nats` | StatefulSet | `takserver` |
| `takserver-nats-box` | Deployment | `takserver` |
| `takserver-plugins` | Deployment | `takserver` |
| `takserver-postgresql` | StatefulSet | `takserver` |

## When it looks wrong

1. Open the app page. The status dot combines Argo CD, workload readiness, alerts and uptime; grey means unknown, not healthy.
2. **Live** tab: confirm pods and replicas. **Logs** tab: open the generated Loki query. **Metrics** tab: open the workload dashboard with this app's variables.
3. Find the cause before acting. A restart that hides the cause is not a fix.

## Safe recovery

- Change the cluster only through Git and Argo CD (pull request, merge, sync). Do not `kubectl apply`, patch or delete by hand, and do not run Helm directly.
- A Git change is reversible with `git revert`; Argo CD then restores the previous state.
- Data stores: confirm the latest backup in the [backups runbook](../backups.md) before any change that touches storage.

## Generated links

- Loki query: `{k8s_namespace_name="takserver", k8s_pod_name=~"(mumble-bridge|murmur|tak-lb-proxy|tak-ldap|tak-nats-leaf|takserver-api|takserver-config|takserver-ignite|takserver-messaging|takserver-nats|takserver-nats-box|takserver-plugins|takserver-postgresql)-.*"}`
- Grafana dashboard path: `d/onemind-workload/onemind-workload?var-namespace=takserver&var-workload=mumble-bridge&var-workload=murmur&var-workload=tak-lb-proxy&var-workload=tak-ldap&var-workload=tak-nats-leaf&var-workload=takserver-api&var-workload=takserver-config&var-workload=takserver-ignite&var-workload=takserver-messaging&var-workload=takserver-nats&var-workload=takserver-nats-box&var-workload=takserver-plugins&var-workload=takserver-postgresql`
- Source: `onemindos/omos`

## Escalate

If this app is unreachable and the cause is not clear after the Live and Logs tabs, treat it as an incident: record what was observed (time, tab, message) and follow the [resilience runbook](../resilience.md).

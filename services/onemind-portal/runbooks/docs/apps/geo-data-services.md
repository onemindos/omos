# OneTAK geospatial services

> Generated from the reviewed catalog (workload-map.yaml, structure.yaml). It lists what is recorded, not what is currently healthy; check the Live tab for current state.

Martin, Overpass, pg_featureserv, Tile38, STAC API/browser and the geo database.

## Where it runs

- Catalog entity: `component:default/geo-data-services` in System `geospatial`
- Argo CD application: none recorded for this app
- Namespaces: `onetak`, `dp-default-onetak-development-b188bf98`

## Workloads

| Workload | Kind | Namespace |
| --- | --- | --- |
| `martin` | Deployment | `onetak` |
| `onetak-geo-db` | Deployment | `onetak` |
| `overpass` | Deployment | `onetak` |
| `pg-featureserv` | Deployment | `onetak` |
| `stac-api` | Deployment | `onetak` |
| `stac-browser-development-0c6b80f7` | Deployment | `dp-default-onetak-development-b188bf98` |
| `tile38` | Deployment | `onetak` |

## When it looks wrong

1. Open the app page. The status dot combines Argo CD, workload readiness, alerts and uptime; grey means unknown, not healthy.
2. **Live** tab: confirm pods and replicas. **Logs** tab: open the generated Loki query. **Metrics** tab: open the workload dashboard with this app's variables.
3. Find the cause before acting. A restart that hides the cause is not a fix.

## Safe recovery

- Change the cluster only through Git and Argo CD (pull request, merge, sync). Do not `kubectl apply`, patch or delete by hand, and do not run Helm directly.
- A Git change is reversible with `git revert`; Argo CD then restores the previous state.
- Data stores: confirm the latest backup in the [backups runbook](../backups.md) before any change that touches storage.

## Generated links

- Loki query: `{k8s_namespace_name=~"onetak|dp-default-onetak-development-b188bf98", k8s_pod_name=~"(martin|onetak-geo-db|overpass|pg-featureserv|stac-api|stac-browser-development-0c6b80f7|tile38)-.*"}`
- Grafana dashboard path: `d/onemind-workload/onemind-workload?var-namespace=onetak&var-namespace=dp-default-onetak-development-b188bf98&var-workload=martin&var-workload=onetak-geo-db&var-workload=overpass&var-workload=pg-featureserv&var-workload=stac-api&var-workload=stac-browser-development-0c6b80f7&var-workload=tile38`
- Source: `onemindos/omos`

## Escalate

If this app is unreachable and the cause is not clear after the Live and Logs tabs, treat it as an incident: record what was observed (time, tab, message) and follow the [resilience runbook](../resilience.md).

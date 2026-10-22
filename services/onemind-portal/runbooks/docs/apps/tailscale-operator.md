# Tailscale operator and proxies

> Generated from the reviewed catalog (workload-map.yaml, structure.yaml). It lists what is recorded, not what is currently healthy; check the Live tab for current state.

Tailscale Kubernetes operator and the per-service tailnet proxies it creates.

## Where it runs

- Catalog entity: `component:default/tailscale-operator` in System `infrastructure`
- Argo CD application: `edge-tailscale-operator`
- Namespace: `tailscale`

## Workloads

| Workload | Kind | Namespace |
| --- | --- | --- |
| `egress` | StatefulSet | `tailscale` |
| `operator` | Deployment | `tailscale` |
| `ts-agent-console-8khxt` | StatefulSet | `tailscale` |
| `ts-agent-tap-v7ktr` | StatefulSet | `tailscale` |
| `ts-argocd-slkzh` | StatefulSet | `tailscale` |
| `ts-cilium-gateway-gateway-default-bl4tw` | StatefulSet | `tailscale` |
| `ts-cilium-gateway-gateway-default-hl5fh` | StatefulSet | `tailscale` |
| `ts-grafana-ksprn` | StatefulSet | `tailscale` |
| `ts-headlamp-4qwkc` | StatefulSet | `tailscale` |
| `ts-hermes-dashboard-6thbx` | StatefulSet | `tailscale` |
| `ts-heynats-p5dmf` | StatefulSet | `tailscale` |
| `ts-hubble-tw85d` | StatefulSet | `tailscale` |
| `ts-kubernetes-dashboard-9r9fp` | StatefulSet | `tailscale` |
| `ts-n8n-pmtds` | StatefulSet | `tailscale` |
| `ts-nats-dashboard-242wv` | StatefulSet | `tailscale` |
| `ts-nats-monitor-zh589` | StatefulSet | `tailscale` |
| `ts-nats-ui-d6rrg` | StatefulSet | `tailscale` |
| `ts-nats-ws-4q6bc` | StatefulSet | `tailscale` |
| `ts-observer-pftsp` | StatefulSet | `tailscale` |
| `ts-omos-ha-edge-rlv6k` | StatefulSet | `tailscale` |
| `ts-omos-mbp-hgq9w` | StatefulSet | `tailscale` |
| `ts-omos-nats-fs6r4` | StatefulSet | `tailscale` |
| `ts-omos-s24-44l4v` | StatefulSet | `tailscale` |
| `ts-onemind-auth-4w5hk` | StatefulSet | `tailscale` |
| `ts-onemind-preview-shtd7` | StatefulSet | `tailscale` |
| `ts-opencost-rqpmq` | StatefulSet | `tailscale` |
| `ts-tak-tailnet-jtgn4` | StatefulSet | `tailscale` |

## When it looks wrong

1. Open the app page. The status dot combines Argo CD, workload readiness, alerts and uptime; grey means unknown, not healthy.
2. **Live** tab: confirm pods and replicas. **Logs** tab: open the generated Loki query. **Metrics** tab: open the workload dashboard with this app's variables.
3. Find the cause before acting. A restart that hides the cause is not a fix.

## Safe recovery

- Change the cluster only through Git and Argo CD (pull request, merge, sync). Do not `kubectl apply`, patch or delete by hand, and do not run Helm directly.
- A Git change is reversible with `git revert`; Argo CD then restores the previous state.
- Data stores: confirm the latest backup in the [backups runbook](../backups.md) before any change that touches storage.

## Generated links

- Loki query: `{k8s_namespace_name="tailscale"}`
- Grafana dashboard path: `d/onemind-workload/onemind-workload?var-namespace=tailscale`
- Source: `onemindos/omos`

## Escalate

If this app is unreachable and the cause is not clear after the Live and Logs tabs, treat it as an incident: record what was observed (time, tab, message) and follow the [resilience runbook](../resilience.md).

# OneMind live bridge (backend, pluginId `livebridge`)

Streams real-time data to browsers through Backstage Signals and serves a snapshot. It is
**read-only** everywhere: NATS core subscriptions only (never publish, never JetStream
consumers, streams or KV), and Kubernetes `get`/`list`/`watch` only.

The wire contract is `src/contract.ts`. It is shared with the Ops Center page; change both
copies together.

## Sources

| Source     | What                                                                                   | Notes                                                                                                                                           |
| ---------- | -------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------- |
| NATS       | `tak.cot.>`, `intel.>`, `evt.>`, `det.>`, `ent.>`, `tel.>`, `cmd.>`, `agents.status.>` | `_INBOX.` and `$` subjects are ignored. `disabled` without `NATS_USER` and `NATS_PASSWORD`. Reconnects forever.                                 |
| Kubernetes | pods and Warning events, cluster-wide                                                  | In-cluster config, else the default kubeconfig. Re-lists and re-watches on any end or error (410 Gone included) with backoff of 1 s up to 60 s. |
| Argo CD    | `applications.argoproj.io` in `argocd`                                                 | Same loop.                                                                                                                                      |

The ServiceAccount `backstage` already has get/list/watch on pods, events (`backstage-observer`
ClusterRole) and Argo applications (`backstage-deployment-observer` Role in `argocd`), so no
RBAC change is needed (`platform/backstage/rbac.yaml`).

## Output

| Channel                | Content                                             | Rate                                                                                 |
| ---------------------- | --------------------------------------------------- | ------------------------------------------------------------------------------------ |
| `onemind:live:tracks`  | CoT positions (upserts, removes)                    | at most 1 per second                                                                 |
| `onemind:live:events`  | sampled bus events                                  | at most 1 per second, 50 events, 2 per subject per second; `dropped` counts the rest |
| `onemind:live:cluster` | changed namespace cells, Argo cells, Warning events | at most 1 per second, only on change                                                 |
| `onemind:live:status`  | source status                                       | every 5 s                                                                            |

Every message carries `origin` (hostname plus random suffix) and an increasing `seq`. No
message is empty and each is below 56 KB of JSON (the Signals limit is 60 KB); larger deltas are
split. `GET /api/livebridge/snapshot` returns the full `LiveSnapshot` to signed-in users.

## Caps (memory is bounded)

- 5000 tracks, least recently updated evicted; kept until the CoT `stale` time but at least
  2 minutes and at most 30 minutes.
- Ring of the last 200 bus events; at most 50 warnings; one slim record per pod.
- CoT larger than 64 KB is skipped; JSON summaries are parsed only for sampled events and only
  below 16 KB.

## Namespace health

`down`: any Failed pod, or a container waiting in CrashLoopBackOff, ImagePullBackOff or
ErrImagePull. `warn`: otherwise any Pending pod, running pod not ready, or restart in the last
hour. `ok` otherwise. Restarts in the last hour are approximated from `restartCount` increases
seen since the replica started plus containers whose `lastState.terminated.finishedAt` is within
the hour. Completed Job pods are not counted.

## MCP action

`livebridge_snapshot_summary` (read-only, non-destructive, idempotent) returns namespace health
counts, unhealthy namespaces, Argo applications that are not Synced and Healthy, track counts by
CoT affiliation, bus rates and the last 20 events. `livebridge` is in
`backend.actions.pluginSources`. The data is what this one replica currently holds.

## Config

```yaml
onemind:
  livebridge:
    nats:
      url: nats://nats.nats.svc.cluster.local:4222 # default
      user: ${NATS_USER} # secret portal-nats-creds
      password: ${NATS_PASSWORD}
      # caFile: /path/ca.crt                          # only for a TLS listener
    kubernetes:
      enabled: true
      argoNamespace: argocd
```

The two replicas each keep their own picture and publish under their own `origin`; clients lock
onto one origin.

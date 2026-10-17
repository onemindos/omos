# Monitoring and health

## Access and dependencies

Open `https://grafana.tail717ba0.ts.net/`. Prometheus, Grafana, Loki and Tempo run in namespace `monitoring`. P3 uses fixed Prometheus queries for backstage, auth, nats, takserver and onetak.

## Interpret the panel

CPU is a five-minute rate in cores; memory is working-set usage. Missing series are unknown, not zero. Four internal service checks report response status and latency. Transport errors are unknown; 401/403 are authentication-required. Checks occur on demand with a shared minute cache, not continuous uptime history. They do not validate Tailscale access or full user workflows.

## Investigate

```sh
kubectl -n monitoring get deployments,statefulsets,pods,services,pvc
kubectl -n backstage logs deployment/backstage --tail=100
```

If metrics disappear, distinguish an unavailable Prometheus API, an empty query result, a scrape failure and an unavailable dashboard. Inspect target health and query the intended namespace in Grafana. Check scrape labels and sample freshness before changing metric names or assuming zero usage.

For a failed HTTP check, use the catalog/investigation link and verify the service from the correct network. A successful internal check with an inaccessible tailnet URL points to routing, DNS, certificates or identity rather than proving the user path is healthy.

## Recovery

Preserve time-series storage and credentials. Correct scrape/routing configuration through a reviewed GitOps PR. Do not delete monitoring volumes to resolve an empty chart. If a data restore is needed, first verify that monitoring backups and a compatible restore process exist; the portal's backup does not back up Prometheus, Loki or Tempo.

## Verify completion

Confirm current queries return expected namespaces and meaningful samples, the portal labels missing sources clearly, and Grafana opens for the authorized user. Review the affected application's real workflow separately from monitoring availability.

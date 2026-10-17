# NATS fabric

## Access and dependencies

NATS runs in namespace `nats`; historical `nats.hermes.svc` addresses are outdated. Use the current catalog and service definitions. Browser investigation tools include `https://nats-dashboard-1.tail717ba0.ts.net/`, `https://nats-ui-1.tail717ba0.ts.net/`, and Grafana. Client traffic uses the configured TLS/authentication settings; a web UI response does not prove broker delivery.

## Investigate

```sh
kubectl -n nats get statefulsets,deployments,pods,services,pvc
kubectl -n nats get events --sort-by=.metadata.creationTimestamp
kubectl -n nats get cronjob backup-nats-streams
```

Inspect broker readiness, recent restarts, storage capacity, consumer lag and pending acknowledgements through the authorized dashboards. Determine whether the failure affects publishing, consuming, a specific subject, or a bridge. Check producer/consumer configuration against the current namespace and catalog. Inspect targeted logs without copying credentials or message payloads into public records.

## Recovery and backups

The backup schedule previously had no successful run recorded in the attention panel. That is not proof of a recoverable backup. Confirm the actual backup destination, retention and restore procedure from the current GitOps configuration before relying on it.

Do not delete JetStream volumes, reset streams, purge messages or recreate consumers as a generic recovery step. For an application/configuration regression, use a reviewed known-good GitOps change. A data recovery requires a verified snapshot, a separate target and a scoped restore procedure; preserve the original data until verification.

## Verify completion

Verify broker and storage readiness, then perform an authorized publish/consume round trip on a dedicated test subject with non-sensitive data. Confirm the affected consumer resumes and lag decreases. Confirm dashboard availability separately from message delivery.

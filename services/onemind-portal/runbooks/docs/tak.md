# TAK services

## Access and dependencies

Core TAK runs in namespace `takserver`; geospatial services and associated data services run in `onetak`; ML services run in `takml`. Use catalog access links and current certificates. Archived notes mentioning `cloudtak` or an old public IP are historical. Do not disable TLS validation to make a client connect.

## Investigate

```sh
kubectl -n takserver get deployments,statefulsets,pods,services,pvc
kubectl -n onetak get deployments,pods,services,pvc
kubectl -n takml get deployments,pods,services
kubectl -n argocd get application platform-takserver
kubectl -n takserver get cronjob backup-takserver-postgres
```

Separate certificate/enrollment issues, message routing issues, database/storage issues and map/geo service problems. For a degraded application, inspect the specific workload's readiness and events. Confirm client endpoint and certificate validity against the deployed configuration. The database is `takserver-postgresql`. Do not infer connectivity from a website or repository being present.

## Recovery and backups

The TAK database backup schedule previously showed no successful run recorded. Verify backup objects and a restore exercise before claiming database recovery is available. Preserve certificates, signing material and database volumes. Roll back configuration through reviewed GitOps changes when compatible. Restore data to a separate target with a service-specific procedure before cutover; avoid unreviewed database consolidation or certificate replacement during incident response.

## Verify completion

Verify the affected workload and dependency are ready. Use an authorized test client with a valid certificate to confirm the intended CoT exchange or map workflow with non-sensitive data. Check downstream bridge/consumer behavior if the incident involved NATS. Record remaining degraded components rather than declaring the whole TAK stack healthy.

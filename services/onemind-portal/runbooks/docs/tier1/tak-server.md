# TAK Server

TAK Server 5.8 on OneMind hardened images, in namespace `takserver`: messaging, API, config, plugins (with the NatsBridgePlugin that carries CoT to and from the bus), Ignite, PostgreSQL/PostGIS and TAK's own internal NATS. End-user devices reach it only over the tailnet.

Sources: `platform/takserver/README.md`, `platform/takserver/DEPLOY-RUNBOOK.md` (section 8, verification only; the rest is historical bring-up), `docs/gitops/HELM-ADOPTION.md` (TAK sections), `docs/operations/2026-10-26-tak-policy-preflight.md`, `docs/operations/tak-enrollment-repair-2026-10-25.md`, `docs/networking/NETWORKING.md`.

## Ownership and access

- `platform-takserver` applies `platform/takserver/deploy`: the LB proxy, NATS leaf, Mumble, LDAP and sealed secrets, plus the adopted messaging, plugins, API and config Deployments, the Ignite, PostgreSQL and internal NATS StatefulSets and the `cot-retention` CronJob.
- The historical Helm release `takserver` (revision 34) still records these objects, and parts of it are not Git-owned yet. Whole-release `helm upgrade`, `helm rollback` and `helm uninstall` are frozen: they would create a second writer for workloads Argo CD now owns. Legacy commands in the chart's upstream README are not instructions for this cluster.
- The legacy ConfigMaps `cert-migration`, `core-config` and `tak-userauth` and some runtime Secrets are not Git-owned, and some of those ConfigMaps hold credential, private-key or keystore data. Never copy their contents into Git or logs.
- Access is `tak.tail717ba0.ts.net` through the Tailscale Service `tak-tailnet` (8089, 8443, 8446 and 64738/TCP). The public load balancer was removed on October 23, 2026, and federation ports 9101/9102 are not exposed. FedHub (`platform-fedhub`) is a separate tier 2 app.
- Bus subjects follow `docs/NATS-DOCTRINE.md`; TAK's internal NATS never carries OneMind traffic.

## Investigate (read-only)

```sh
kubectl --context do-nyc3-omos-cloud -n argocd get application platform-takserver
kubectl --context do-nyc3-omos-cloud -n takserver get deployments,statefulsets,pods,pvc
kubectl --context do-nyc3-omos-cloud -n takserver get service tak-tailnet
kubectl --context do-nyc3-omos-cloud -n takserver logs deploy/takserver-config | grep "Topology snapshot" | tail -1
kubectl --context do-nyc3-omos-cloud -n takserver logs deploy/takserver-api | grep "Tomcat started"
```

A joined Ignite grid reports `ACTIVE` in the config tier's topology snapshot. Separate certificate or enrollment failures, message routing, database or storage problems and the NATS bridge before acting; more symptoms are in [TAK services](../tak.md). Do not disable TLS validation to make a client connect.

Certificate enrollment V2 signs with `/certs-configmap/ca-signing.jks`, projected into messaging only from the SealedSecret `tak-enrollment-signer`. When that file was missing, enrollment returned HTTP 500 while existing clients kept working; the repair and its rollback are in `docs/operations/tak-enrollment-repair-2026-10-25.md`.

## Backups and recovery

- `backup-takserver-postgres` (namespace `takserver`, daily at 02:00 UTC, about 1.5 to 2.5 hours) uploads an encrypted dump of the `cot` database, restores it into an isolated PostGIS sidecar on loopback, and only then writes its verification marker.
- VolSync copies the TAK volumes daily, including PostgreSQL and the internal NATS store. Retained CSI checkpoints of both claims from October 17, 2026 also exist.
- On October 26, 2026 the preflight confirmed a verified logical restore (64 tables) and booted the internal NATS archive with the exact running server version. It did not test a raw PostgreSQL volume restore or a full running TAK recovery.
- Restore to a separate target, never over the source claims. Preserve the CA, certificates, signing material and database volumes. Roll back configuration through a reviewed Git change.

## Verify completion

`platform-takserver` is Synced and Healthy and the workloads are Ready. An authorized test client with a valid certificate connects over the tailnet and completes the expected CoT exchange or enrollment with non-sensitive data. Record any component that is still degraded rather than declaring the whole TAK stack healthy.

Generated workload list: [TAK Server app runbook](../apps/tak-server.md).

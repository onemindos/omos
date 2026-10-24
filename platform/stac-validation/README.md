# Disposable STAC migration validation

Source-only diagnostic fixture for the 2026-10-24 shared-PostgreSQL incident.
It has not been executed, and does not establish the crash cause or authorize a
production migration. Keep `platform-stac` manual and its migration Job absent.

`platform-stac-validation` is a manual Argo CD Application without automated sync,
pruning, retries or a deletion finalizer. Its root Application may register it;
only an explicitly requested sync runs these diagnostics. No Helm or imperative
workload installation is required. Coordinator owns review, merge, manual sync,
evidence collection and cleanup. No production manifests or credentials change.

The source adds **one** Argo Application to inventory (+1). Its reviewed canonical
mapping points to `component:default/stac`; catalog guards now expect 55 source
Applications, including the CloudTAK Application merged separately.

## Cases and boundaries

Control (wave 1) has no AGE preload. AGE (wave 2) sets
`shared_preload_libraries=age`. Argo waits for the control Job to complete before
starting AGE, so only one case runs during a normal full sync. If control fails,
stop and inspect it; AGE does not run automatically. Each case receives a separate
fresh `emptyDir`; the only intended database difference is AGE preload. Both run
PostgreSQL 16 with the production max_connections/shared_buffers values.

Each native PostgreSQL sidecar starts from initdb. The client creates synthetic
`stac` LOGIN CREATEROLE NOSUPERUSER, owned database `stac`, and PostGIS, btree_gist
and unaccent as bootstrap superuser, matching production's documented prerequisites.
AGE is deliberately absent from the synthetic STAC database. Migration uses
`PGOPTIONS=-c createrole_self_grant=set,inherit` and runs unchanged `pypgstac migrate`
with its original transaction behavior. Only after success does it set collection
extent updates, check version, rerun the idempotent migration and read schema state.
No remote source records or inference are involved.

There is no Service, ingress, PVC, host mount, Secret, image pull credential,
production service address or Kubernetes API access. PostgreSQL binds only
`127.0.0.1`. Passwordless trust exists only inside that isolated pod; Unix sockets
are on the PostgreSQL container's private `/tmp`. The namespace denies all ingress
and egress, including DNS. Kubelet image pulls happen outside pod networking;
the pinned GHCR images must be anonymously pullable, otherwise stop at
ImagePullBackOff rather than injecting production credentials.

The Jobs use Kubernetes native sidecars (`initContainers[].restartPolicy: Always`)
and pod `restartPolicy: Never`. Kubernetes 1.36 supports this lifecycle: the client
finishes without a database container keeping its Job running indefinitely.
There are zero Job retries and a 600-second deadline. PostgreSQL's
`restart_after_crash=off` and an emptyDir attempt marker prevent a second database
start if its native sidecar restarts; the original server log is reprinted.

Per case: PostgreSQL request 50m CPU/384Mi memory, limit 1 CPU/1Gi; client request
50m/128Mi, limit 1 CPU/512Mi. Scratch disk limit is 1Gi plus 96Mi temporary space.
Normal sync serializes cases; requests total 100m/512Mi and memory limits 1.5Gi.
Selecting both Jobs manually can bypass wave ordering and is outside this runbook.
All containers run nonroot with read-only roots, dropped capabilities and no
service-account token. No core dump support, debug package or binary is installed.

## Evidence and cleanup

After the coordinator's reviewed manual sync, collect each Job's `migrate` and
`postgres` container logs, Pod status (including previous state/restart count),
Job conditions, and the deployed imageIDs. Logs contain only this fixture's
synthetic SQL. PostgreSQL logs every statement with timestamp/PID/user/database,
verbose errors and `log_min_error_statement=error`; its wrapper prints the full
server log when PostgreSQL exits, including on client completion/sidecar shutdown.
If PostgreSQL crashes, a restarted sidecar reprints the original log and waits;
it does not retry SQL. Collect `--previous` logs too when available. The wrapper
buffers logs on scratch disk until exit; inspect that file only within this
synthetic pod if live progress is needed.

Both Jobs must complete successfully, report version 0.9.12 and pass the second
migration, with no PostgreSQL crash/restart in logs or Pod status. A failure is
an observation to investigate, not permission to try the shared database. A/B
results alone are a hypothesis test; identify the failing statement before
proposing a runtime fix. If server logging cannot be collected, retain the pod
and gather its synthetic scratch log before deleting it.

Jobs have no TTL or hook deletion policy so evidence survives until review.
After saving evidence, coordinator deletes only `stac-validation-control` and
`stac-validation-age` Jobs (their emptyDir data disappears with the pods). Do not
sync again to rerun them without a fresh review. Then retire this Application
and source in a cleanup PR. Its lack of finalizer means deleting the Application
does not delete its namespace or resources; remove remaining diagnostic
ConfigMap, policy and empty namespace explicitly only after verifying no other
resources were introduced. Never delete production resources or retained PVCs.

## Source proof and local checks

- PostgreSQL image equals `platform/ecosystem/resources/postgres-image-catalog.yaml`:
  digest `35302d70df121bc86b370ebad800d85401edcdb890f8146cd3db690826a9925d`.
- Client equals `platform/stac/migrate.yaml`: pgstac v0.9.12 digest
  `672faad776b55afe15fadb27589c390c4fc93de5bc73e8889d2c29a600c5dba4`.
- `services/onemind-postgres/Dockerfile` pins AGE commit
  `2db2f060c4c9265a14d40f007eb8c56febf31e4c`. Its globally installed native hooks
  are a candidate interaction, not a demonstrated cause.
- [Kubernetes native sidecar lifecycle](https://kubernetes.io/docs/concepts/workloads/pods/sidecar-containers/)
  documents startup probes and Job completion behavior.

Local YAML/schema, embedded shell and safety-contract checks validate source;
they do not execute PostgreSQL or prove migrations passed. The local Docker daemon
was unavailable when this fixture was authored. Record actual run results and
imageIDs separately after execution; never substitute manifest pins for imageIDs.

Authoring verification (source only): four safety-contract tests passed; embedded
Bash syntax and ShellCheck passed; all six manifests passed strict kubeconform
1.36.0 with zero skipped schemas; changed YAML lint passed. `make validate` was
run and stopped at pre-existing unrelated YAML lint errors under onetak, TAK,
NATS, monitoring, data and matrix. No unrelated files were reformatted.

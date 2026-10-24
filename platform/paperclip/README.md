# Paperclip business agent workspace

Prepared 2026-10-23 at Zeus's request in [PR #438](https://github.com/onemindos/omos/pull/438).
**Deployed 2026-10-24:** app and database Ready; authenticated/private health
verified, and Zeus created the first owner and OneMind Enterprise workspace.
Hermes gateway transport is deployed; its Paperclip task round trip is under
verification. Other acceptance items below remain open until evidenced.
This is a user-directed installation alongside the existing portal work; Phase 6
remains In progress and Phases 7/8 remain Not started. It introduces no new
controller or transfer of existing workload ownership.

## Deployment contract

- Argo Application `platform-paperclip`, namespace `paperclip`, sync wave 4.
  Source is this directory on reviewed `main`; automatic self-heal, prune off.
- Paperclip `2026.1012.0`, release commit
  `bba6a45f0d364cfd962ac59b3bdf2426a8222adb`. The public GHCR multi-architecture
  index was resolved to the digest in `deployment.yaml`; no mutable image tag.
- Proposed private URL: <https://paperclip.tail717ba0.ts.net/>. Existing Tailscale
  operator terminates HTTPS. Paperclip's own authenticated mode is enabled;
  tailnet membership is not a Paperclip account or board permission.
- One non-root app replica and a dedicated PostgreSQL 17 instance. A non-superuser
  `paperclip` database role owns only the new database; existing shared databases
  and application credentials are unchanged.
- Two retained 10 GiB source PVCs plus two 1 GiB VolSync cache volumes and temporary
  CSI clones/snapshots. This consumes provider storage and existing compute;
  no node-pool expansion or model billing change is configured. Initial steady
  requests: app 250m/512 MiB, database 50m/128 MiB, plus the Tailscale proxy.
- No Kubernetes API token, RBAC, host mounts, container socket, public ingress,
  NATS integration or new permissions for existing agents. Restricted PodSecurity,
  default-deny network policy, and public HTTP(S) egress with private ranges blocked.
  Outbound HTTPS still permits sending data to public services; it is not a DLP boundary.

## Identity and agent onboarding

The release image includes Claude Code and Codex CLIs. Having an adapter does not
mean an agent is connected or authenticated. The image does not install the Hermes
CLI: Hermes needs a separately prepared runtime image or a compatible remote
Hermes gateway; do not install tools imperatively into a running production pod.
Do not mount/copy the Mac's agent homes, plugins, chats or provider credentials.
Their complete desktop session state is not implicitly inherited.

`HEARTBEAT_SCHEDULER_ENABLED=false` deliberately leaves the initial installation in
setup mode (the release disables the heartbeat service, including task execution).
After the server is healthy, the operator must:

1. Sign in/create an account at the private URL and claim the fresh instance.
   `authenticated/private` supports a browser board claim on this release.
2. Create the company and roles, then connect the intended runtime accounts through
   the authenticated UI. Verify account identity, credential scope and actual billing
   for each runtime. No model/provider credentials are supplied in Git.
3. Prepare/verify Hermes execution independently. Prove a bounded task with each
   runtime, session continuity, handoff, pause and review behavior before business use.
4. Review per-agent tool permissions. Upstream coding adapters can bypass interactive
   tool approval by default; board approvals do not sandbox tools. Local adapters run
   inside this container and can reach its own app data. For stronger separation,
   use dedicated execution environments; do not grant cluster administration.
5. Set budgets and approve a Git change enabling the heartbeat service. Confirm
   actual cost reporting and the limits that apply to API versus subscription runs.

The phone browser or an independently evaluated mobile client can connect over
Tailscale. Phone notifications, native client compatibility, offline behavior and
actual business workflows have not been tested by preparing these manifests.
The server must be reachable; cloud model execution still needs its provider.

## Secrets and recovery

Only cluster-bound SealedSecrets are committed. App session/tool-action signing
keys, local-encryption master key and database passwords are independently generated.
The app receives its database URL and app keys, not the PostgreSQL admin password.
The database pod receives only its own bootstrap passwords. Existing object-store
and registry credentials are resealed into this namespace for the existing backup
pipeline, with new Paperclip archive and Restic encryption keys and repository paths.

New recovery material is outside Git at
`~/.config/onemind/recovery/paperclip-credentials.json` (mode 0600). Preserve an
independent offline copy plus the existing cluster sealing-key recovery export.
A local recovery file is not itself an independent/offsite backup. Do not rotate
keys during redeploys; restoring the database without the matching Paperclip master
key cannot decrypt saved provider secrets.

- `backup-postgres-paperclip`: 06:00 UTC daily logical dump to
  `backups/postgres-paperclip`, using the existing pinned PostgreSQL-17 backup image.
  Every run downloads/decrypts its archive and restores into a disposable loopback-only
  PostgreSQL sidecar. A verified archive does not prove all application workflows.
- `paperclip-data` and `paperclip-postgres`: daily encrypted VolSync/Restic exports
  at 12:00 and 12:20 UTC, including app workspaces/assets and the database volume.
  They have separate repository prefixes and the existing retention pattern.
  Raw database snapshots are crash-consistent; use logical archives for DB recovery.
- These backups use existing storage in the same provider/region as the cluster.
  No geographically independent backup or coordinated app/DB restore is claimed.
- Source volumes and the StatefulSet are protected against Argo pruning/deletion;
  storage uses Retain. Namespace deletion can still affect claims. Never delete the
  namespace, PVCs, PVs, archives or encryption keys as a rollback shortcut.

## Acceptance after approved merge

Argo, rather than `kubectl apply` or Helm CLI, must create every resource.
Do not report installation from PR creation alone.

- [ ] `platform-paperclip` Synced/Healthy at the intended revision.
- [ ] SealedSecrets Synced; source PVCs Bound; app and PostgreSQL ready without restarts.
- [ ] Private HTTPS works on Mac and Samsung over Tailscale; public access is absent.
- [ ] `/api/health` returns status `ok`; anonymous protected API requests are denied.
- [ ] Owner account/claim completed; other tailnet users do not acquire board access.
- [ ] Restart/reconcile preserves database rows, assets and encryption-key usability.
- [ ] Logical archive uploaded and isolated restore verification passed; both VolSync
      exports completed. Restore app and DB together in a separate namespace and verify
      login, company/task state, assets and decryption without contacting providers.
- [ ] Each intended runtime passes its authenticated execution and review checks;
      scheduler enablement and permissions separately approved.

After merge, inspect with the explicit `do-nyc3-omos-cloud` context:

```sh
kubectl --context do-nyc3-omos-cloud get application platform-paperclip -n argocd
kubectl --context do-nyc3-omos-cloud get pods,pvc,sealedsecrets,replicationsources -n paperclip
```

Rollback before any user data: disable sync through a reviewed Git change and
retain all storage/secrets. After data exists, take a verified fresh backup before
version changes. The server runs release migrations on its dedicated database at
startup; reverting the image does not undo migrations. Restore a compatible
app/DB pair rather than guessing backward schema compatibility.

## Validation and sources

The local Docker daemon is unavailable, so no container boot test was performed.
Validation performed: Kustomize rendered 22 workload resources; strict kubeconform
validated all 23 resources including the Argo Application, with none skipped.
Changed YAML lint and the repository-configured, baseline-aware whole-tree Gitleaks
scan passed. Volume references, immutable image pins and sealed-secret key bindings
were independently checked. Terraform format/validate passed.

`make validate` was run and failed on pre-existing whole-tree YAML debt (including
`platform/onetak/media-mtx.yaml` and legacy Matrix manifests). The separate
`make lint-shell` also fails in existing TAK scripts. No unrelated lint cleanup
was made. These preparation checks do not replace CI or live acceptance.

- [Pinned upstream Dockerfile](https://github.com/paperclipai/paperclip/blob/v2026.1012.0/Dockerfile)
- [Pinned Docker entrypoint](https://github.com/paperclipai/paperclip/blob/v2026.1012.0/scripts/docker-entrypoint.sh)
- [Pinned configuration](https://github.com/paperclipai/paperclip/blob/v2026.1012.0/server/src/config.ts)
- [Agent adapters](https://docs.paperclip.ing/reference/adapters/overview/)
- [Existing backup contract](../backups/image/README.md)
- [Existing VolSync recovery notes](../backups/volsync/README.md)

## Hermes connection

[Private native Hermes adapter transport and activation](HERMES.md) separates
GitOps transport, public CA trust, agent credentials and functional acceptance.

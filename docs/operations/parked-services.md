# Parked services register

The record of every service removed from the cluster during the 2026 cost and cleanup review, so nothing is lost: where its configuration went, where its data is, and how to bring it back. The full inventory and the reasons are in the "Cluster Keep or Cut" review (2026-10-16).

## How a service is removed (one at a time)

1. **Discuss.** What it does, who uses it (evidence from logs, traffic and dependencies), what data it holds, what would be lost. Zeus decides: keep, park or cut.
2. **Preserve.**
   - Trigger a fresh VolSync backup of each of its volumes and wait for it to succeed.
   - Note the restic repository name for each volume below. The repository stays in the backup bucket after the service is gone, and nothing prunes it once its backup source is removed.
3. **Remove through Git (one PR).**
   - Move its Application (and values) into `platform/parked/<service>/`.
   - Remove its VolSync backup source from `platform/backups/volsync/sources.yaml`, and its namespace from `platform/ecosystem/resources/namespaces.yaml` if it has its own.
   - CI must pass; Zeus approves the merge.
4. **Apply.** The ecosystem parent Applications run with pruning off, so removing a file does not delete anything by itself. After the merge, delete the service's Application with cascade, once, by name. This is the only manual step, and the PR states the exact command.
5. **Verify.**
   - Only that service's objects are gone.
   - Every other Application is still Synced and Healthy.
   - Hermes is still on the bus, and the public sites still answer.
6. **Data.** Volumes created by a StatefulSet survive the cascade. They are deleted only after 30 days and only with Zeus's explicit go-ahead; until then the service can be restored in minutes.
7. **Record.** Fill in the row below and update the review page.

## Register

| Service | Decision | Date | Config now at | Data backup (restic repository, last snapshot) | Volumes deleted | Notes |
|---------|----------|------|---------------|-----------------------------------------------|-----------------|-------|
| omos-mission-engine | Parked | 2026-10-17 | `platform/parked/omos-mission-engine/` (was `platform/omos/`) | None: it ran at 0 replicas with no volume | n/a | Image `ghcr.io/onemindos/omos-mission-engine:latest`; its source code exists nowhere, so it cannot be rebuilt. Bringing it back needs the code first |
| kagent + agent-sandbox | Removed (parked) | 2026-10-17 | `platform/parked/kagent/` (was `platform/ecosystem/`) | No volumes. `kagent` database dropped 2026-10-17 after a verified `pg_dump` at `~/Backups/omos/kagent/kagent-2026-10-17.dump` (SHA-256 `d1ed3ba3…4221964`) | Database dropped 2026-10-17 | 0 use in 7 days. Hermes is the agent harness. Vault profile: `Software/Self-Hosted/AI and Agents/Agent Platforms/kagent - Kubernetes AI Agent Runtime.md`. LiteLLM key and Authentik provider to revoke by hand |

## Restoring a parked service

1. Move `platform/parked/<service>/` back to where the register says it came from, and add its namespace and backup source back.
2. If its volumes were deleted, restore them first: create a PVC and a VolSync `ReplicationDestination` pointing at the restic repository in the register (the 2026-10-16 Hermes restore drill shows the exact pattern, `hermes-restore-drill-2026-10-16.md`).
3. Merge; ArgoCD recreates the service.

## Known drift not handled here

- `platform-matrix` stays on manual sync. Git does not match the running Matrix server: Git still points Element at `onemindos.com` while the live server is on `onemindos.dev` (open PR 192), and Git expects a `synapse-homeserver` Secret that does not exist (the live config is in the `synapse-config` ConfigMap). Syncing it as it is would break Matrix. It needs its own careful step: finish PR 192, turn the live homeserver config into a SealedSecret, diff to zero, then enable sync.

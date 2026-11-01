# Self-hosted Beeper bridges

Owner-approved deployment queue, 2026-10-31. This PR prepares WhatsApp, Telegram, Google Messages, Instagram and X together. Source preparation is not a running deployment or authenticated account acceptance.

The Argo Application `platform-beeper-bridges` owns this directory and stays on manual sync for the owner's deployment batch. Do not apply these workloads individually. No existing Matrix, Hermes, messaging account or Docker app is changed by preparing this source.

## Runtime

Each connector has one non-root Deployment with a Recreate rollout and its own 5 GiB `do-block-storage-retain` claim. The upstream manager image is digest-pinned. `artifacts-lock.json` records immutable public CI job URLs, exact source commits and SHA256 checksums for the five Linux amd64 connector binaries. Those downloads were checked as ELF64 x86-64; their authentication and messaging have not been tested.

The init container verifies the cached or downloaded binary. The manager uses a custom startup command so restarts cannot replace it with an unreviewed upstream build. The generated configuration and session/SQLite state persist; runtime configuration applies Telegram application credentials privately. Connector upgrades require updating the lock and checking database compatibility. CI artifact retention is upstream-controlled; the retained claim and encrypted backups preserve the downloaded binary, and a durable internal artifact mirror is a follow-up for rebuild independence.

The pods initiate network connections to Beeper and the external messenger services. They have no Kubernetes API token, no public Service/Ingress, and no direct dependency on the independent Synapse server. Bridge Manager connects to a Beeper account/backend. This is self-hosted transport, not a fully independent Beeper service. Self-hosted accounts are exempt from Beeper's connected-account quota; premium app features are separate. [Upstream container](https://github.com/beeper/bridge-manager/tree/main/docker), [self-hosting and account exemption](https://developers.beeper.com/bridges/self-hosting/).

## Initial credentials

Required resources are `beeper-bridge-auth` (Beeper Bridge Manager authorization), `beeper-telegram-auth` (Telegram API ID/hash), and five separate encrypted backup repository credentials. There are no dummy credentials or plaintext overrides in Git. `sealed/kustomization.yaml` intentionally starts empty until real authorization exists.

Complete a production `bbctl login` on the owner's machine. Telegram application credentials come from [my.telegram.org/apps](https://my.telegram.org/apps); keep the ID and hash in private local files or protected mounts, not chat, shell arguments or repository files.

Fetch the public sealing certificate from the verified cluster, then prepare the encrypted files:

```bash
kubeseal --context do-nyc3-omos-cloud --fetch-cert \
  --controller-name sealed-secrets-controller \
  --controller-namespace kube-system > /tmp/beeper-sealing-cert.pem

python3 scripts/prepare_beeper_secrets.py \
  --context do-nyc3-omos-cloud \
  --cert /tmp/beeper-sealing-cert.pem \
  --bbctl-config '/Users/zeus/Library/Application Support/bbctl/config.json' \
  --telegram-api-id-file /path/to/private/telegram-api-id \
  --telegram-api-hash-file /path/to/private/telegram-api-hash
```

The helper reads existing `ops/platform-backup-credentials` into process memory, reuses its established storage access for dedicated repository prefixes, and generates a distinct Restic password for each bridge. It seals all resources before saving any output, refuses to overwrite previous encrypted resources, and writes ciphertext plus an updated sealed-resource Kustomization. It never applies objects or logs credential values. Re-running it is not a rotation procedure.

Review/commit the encrypted outputs and merge the deployment PR in the intended batch. Then sync `platform-beeper-bridges` through Argo. Do not substitute the independent Synapse access token or a Beeper Desktop MCP token. Check SealedSecrets readiness before diagnosing workload startup.

## Connect all five accounts

Bridge names are `sh-onemind-whatsapp`, `sh-onemind-telegram`, `sh-onemind-gmessages`, `sh-onemind-instagram` and `sh-onemind-twitter`. Use the corresponding self-hosted bridge bot in Beeper to pair/log in. Pairing is owner-assisted and does not authorize unsolicited test messages to contacts.

All five are included in the same deployment batch. Verify each connection on Android and Mac: inbound/outbound messages using an owner-approved test conversation, media, expected history and reconnect behavior. Existing hosted/on-device connections should be assessed before removal because re-adding a connection may not recover all history. Do not delete a Beeper bridge or reset recovery keys as a restart step.

Google Messages still requires an online Android phone. Telegram requires its own application API credentials. Native apps remain needed for unsupported calls/account-management features. Pods becoming Ready only proves the runtime responds; account login and useful sync are separate acceptance checks.

## Backup and recovery

`backup.yaml` provides one daily VolSync snapshot-to-Restic source for every bridge claim, on staggered schedules. Each uses its own encrypted repository prefix/password. Source claims use Retain storage; temporary snapshot-copy/cache claims explicitly use `do-block-storage`, so routine backup cleanup does not retain billed throwaway volumes.

The whole-volume copy contains SQLite databases, WAL/session files, generated credentials, encryption keys and the locked connector binary. It is crash-consistent, not a complete export of every upstream conversation. Verify the archived files and SQLite recovery in an isolated new claim, then launch the same pinned connector with only one active writer. Remote services may require reauthentication after restore. Never replace the production PVC or delete existing snapshots to perform a drill.

Preserve the repository passwords and SealedSecrets controller recovery keys independently. The sealed ciphertext alone is insufficient if the controller's private keys are lost. A completed backup and a successful restore are distinct live acceptance checks; neither has been performed for this newly prepared app.

Rollback is a reviewed Git change to stop the five Deployments or restore the previous lock. Keep PVCs, encrypted repository credentials and remote archives. Do not delete the namespace or bridge registrations as rollback.

## Human and AI access

Human access uses Beeper Android/macOS. Transport pods are not an MCP endpoint. The `beeper` entry in Hermes's application inventory is explicitly disabled until a separately authorized Desktop or Android MCP connection passes private reachability and authentication checks. No message content is granted to agents by this PR.

The existing DigitalOcean cluster is an on-grid deployment. The user's future off-grid Matrix service needs a locally reachable server and local inference for offline Hermes; this PR does not claim or implement that site.

## Validation

Required local checks include `make validate`, scoped changed-YAML lint, `kubectl kustomize platform/beeper-bridges`, manifest schema validation, script ShellCheck and the Beeper GitOps contracts. Existing whole-repository validation debt must be reported separately from new-file results. No Terraform apply, Helm install or imperative workload apply is part of validation.

Source validation on October 31: 7 offline runtime tests and the repository's
218 GitOps tests passed (6 existing skips), with changed-file YAML lint,
ShellCheck, Kustomize/schema validation, installed VolSync schema validation,
MCP coverage, paired catalog generation, Terraform validation and secret scan.
`make validate` stops on existing unrelated YAML formatting debt; the scoped
new-file checks pass. GitHub CI runs on [PR #605](https://github.com/onemindos/omos/pull/605).

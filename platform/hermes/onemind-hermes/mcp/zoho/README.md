# Permanent Zoho MCP connections

These nine connections belong to the official `onemind-hermes` default profile.
They do not create a local Desktop installation or a second agent.

## Durable configuration

- Each `*.sealedsecret.yaml` is strict-scoped to its exact `hermes/<name>-mcp`
  Secret and contains an endpoint encrypted with official `kubeseal`. Do not rename
  a Secret by editing metadata alone: strict-scoped ciphertext must be re-sealed.
- Argo applies these Secrets in sync wave -2 before the Deployment. The main
  container receives each `ZOHO_*_MCP_URL` through `envFrom`. No plaintext endpoint
  appears in a ConfigMap, startup argument, manifest, or this document.
- `servers.json` holds only `${ENV}` references and exact read-tool allowlists.
  Hermes v0.22.5's `tools/mcp_tool_config.py` recursively interpolates URL values;
  this is not an assumed shell expansion. Missing variables fail the connection.
- The existing `initial-config` container runs `reconcile.py` on every pod start,
  including an existing PVC. It calls the supported Hermes `config set` command
  for only the nine owned entries. It leaves `ops`, `data`, `tak`, `dropbox`,
  `clickup`, provider sign-ins and all unrelated configuration alone.
- `config.defaults.yaml` has identical entries for a clean-volume rebuild.
  Kustomize hashes the reconciliation ConfigMap, so changing the definitions
  changes the pod template and triggers Argo's rollout. No new workload or PVC.
- Only the default profile is changed. Role-profile MCP synchronization is a
  separate workflow; this migration does not silently modify other profiles.

## Initial scope

All nine servers are configured, with exact allowlists for identity, navigation,
content reads and selected reporting. This deliberately does **not** expose every
upstream tool. Mutations, user administration, signing, payment initiation and
removal are not included. Expand an allowlist through a reviewed PR after checking
its tool schema and effects. Remote LLM sampling, resources and prompts are off.
Provider-side grants may remain broader than the client allowlist; the allowlist
is not a revocation of the underlying Zoho service authorization.

The canonical Projects connection is `zoho-projects`. The old combined
`zoho-projects-workdrive` pending artifact is not installed and is not deleted.
WorkDrive is its own connection.

## Validation and deployment gates

Run:

```sh
python platform/hermes/onemind-hermes/mcp/zoho/test_reconcile.py -v
kubectl kustomize platform/hermes/onemind-hermes
make validate
```

Before merging, check remote CI. After owner-authorized merge:

1. Observe Argo's `platform-onemind-hermes` revision and health. Inspect all nine
   SealedSecret conditions (`Synced=True`) without displaying Secret data.
2. Observe `initial-config` and pod readiness. Nine CLI writes are a bounded
   startup operation; the init limit permits 192 MiB without increasing requests.
3. In the restarted official instance, run `hermes mcp test <name>` for each
   connection. **Capture and redact its output**: the CLI transport line can print
   the credential-bearing endpoint. Report only status and tool counts.
4. Execute minimal allowed identity/organization reads. Discovery does not prove
   each Zoho application's connection is authorized.
5. Confirm the raw config still contains environment references, not resolved
   endpoints, and that unrelated MCP entries remain present.

Source tests and a valid rendered manifest do not prove the active controller can
unseal a credential. The public certificate is tracked at
`scripts/data/onemind-mcp-sealing.crt`; validate live unsealing after deployment.
If unavailable or obsolete, keep deployment blocked and re-seal with the active
controller's public certificate; never commit a private sealing key.

## Recovery and credential rotation

A pod replacement reloads Secrets and reconciles all nine definitions. A clean
PVC receives the same public configuration. Cluster rebuild also needs the
existing SealedSecrets controller private-key recovery procedure; ciphertext
alone cannot reconstruct that key. Existing PVC backups remain the source for
Hermes conversations and OAuth grants for Dropbox/ClickUp; this Zoho change does
not make those other providers' authentication independently rebuildable.

Zoho Connection-mode endpoints do not use a Hermes browser OAuth callback. The
underlying Zoho service connection must remain authorized, and revoked/rotated
endpoints still need replacement. Do not promise access will never need renewal.

To rotate, seal the replacement endpoint for the same Secret name/key and submit
a PR. A Secret-only change does not itself restart a pod whose credentials use
`envFrom`: also update a non-secret pod-template revision annotation to trigger a
GitOps rollout. Never paste credentials into chat, Git, docs or config. Do not
use chat history as a credential store.

Rollback requires a reviewed Git change that explicitly disables these entries
through the startup reconciler; reverting the seed alone does not remove entries
already on an existing PVC. Preserve credentials and pending artifacts unless the
owner explicitly authorizes the specific deletion with `execute` or `override`.

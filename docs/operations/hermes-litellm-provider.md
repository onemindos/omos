# Hermes through LiteLLM — credential-gated rollout

Task 6 prepares Hermes' OpenAI-compatible route at
`http://litellm.onemind-ai.svc.cluster.local:4000/v1`. It does not deploy or claim a
successful live model response. The existing NATS leaf, persona, tools and Matrix
identity remain in place.

## Model and runtime contract

Default alias `hermes-agent` routes to OpenRouter's
[Claude Sonnet 4.6](https://openrouter.ai/anthropic/claude-sonnet-4.6), restricted
to the Anthropic upstream with provider fallback disabled. No AWS endpoint,
credential or service is configured. The [OpenRouter provider restriction](https://openrouter.ai/docs/guides/routing/provider-selection)
is explicit so that it does not select a different hosting provider.
LiteLLM's model-group fallback is `hermes-agent -> local-tools` (in-cluster
qwen3:4b); `local-embed` (bge-m3) remains a separate embedding route. The small
local model has weaker reasoning/tool reliability and an 8192 context budget;
this is a degraded fallback, not cloud-equivalent agent performance. The local
service still depends on the current DigitalOcean cluster; this does not prove
home/WAN-independent operation. Fallback can take time while upstream retries
expire; it is not a zero-delay failover.

Read-only metadata inspection found `OPENROUTER_API_KEY` in `litellm-providers`,
but did not read its value, verify its validity, spend inference tokens, or prove
available credit. Existing Git aliases onemind-chat/fast still route locally;
the new cloud route is scoped to Hermes. The key allowlist is precisely
`[hermes-agent, local-tools, local-embed]`.

The provider merge was checked against Hermes tag `v2026.9.24`, source commit
`e3dd27ee2d8b011737a4eea8e3eb3d711ab78690`, particularly
[custom provider resolution](https://github.com/NousResearch/hermes-agent/blob/v2026.9.24/hermes_cli/runtime_provider_custom.py).
It supports a `providers.litellm` entry with `key_env`, `default_model` and
chat_completions. The init container atomically changes only provider/model
routing in the existing PVC config, clears stale explicit model endpoint/key
overrides, preserves other settings and file ownership/mode, and never writes
the virtual key into config. It requires an existing config when activating.
The deployed tag's filesystem/runtime dependency parity still needs post-merge
verification; no container was executed during this review.

Without the consumer Secret, the first-run merge leaves current config unchanged.
This is only a staged activation guard: after activation, removing the key will
not undo the persisted provider config. Keep the Secret until rollback restores
previous provider settings. No persona, skill, plugin, NATS or database file is
replaced. Snapshot/recovery verification must precede activation.

## Secret preparation by Zeus

There are two namespace-bound encrypted copies of the same dedicated virtual key:
`onemind-ai/litellm-hermes-client` (key `key`) for the existing bootstrap Job and
`hermes/hermes-litellm` (key `HERMES_LITELLM_API_KEY`) for the agent. This avoids
editing or replacing the existing core/openchoreo client-key Secret. Missing
Hermes credentials are skipped by bootstrap, so existing clients still provision.
No ServiceAccount access to Secrets or new RBAC is introduced.

The `.yaml.example` is deliberately not a Kustomize resource. Never apply a dummy
key or placeholder ciphertext. On Zeus' secure workstation, prepare a fresh key
without displaying it, then seal both copies (these are instructions, not actions
performed by this PR):

```bash
umask 077
keydir=$(mktemp -d)
python3 - "$keydir/key" <<'PYKEY'
import pathlib, secrets, sys
pathlib.Path(sys.argv[1]).write_text("sk-" + secrets.token_urlsafe(32))
PYKEY
kubectl --context do-nyc3-omos-cloud -n onemind-ai create secret generic \
  litellm-hermes-client --from-file=key="$keydir/key" --dry-run=client -o yaml \
| kubeseal --context do-nyc3-omos-cloud \
  --controller-name sealed-secrets-controller --controller-namespace kube-system \
  --format json \
  > platform/ecosystem/resources/secrets/onemind-ai--litellm-hermes-client.json
kubectl --context do-nyc3-omos-cloud -n hermes create secret generic \
  hermes-litellm --from-file=HERMES_LITELLM_API_KEY="$keydir/key" \
  --dry-run=client -o yaml \
| kubeseal --context do-nyc3-omos-cloud \
  --controller-name sealed-secrets-controller --controller-namespace kube-system \
  --format yaml > platform/hermes/agent/hermes-litellm.sealedsecret.yaml
rm -f "$keydir/key"
rmdir "$keydir"
```

Use strict default sealing scope (namespace and name). Review and commit only the
sealed outputs in a credential PR. Add the Hermes sealed file to
`platform/hermes/agent/kustomization.yaml`; do not add the example. No change is
needed to the ecosystem directory source to discover its JSON SealedSecret.
Do not change existing provider credentials, client keys or billing. Remove
plaintext temporary key material on any failed preparation as well.

## Exact reviewer rollout and verification

1. Before merging activation credentials, record original Hermes readiness,
   persistent claim identity and a verified recoverable config backup privately.
   Review this PR plus any separately merged init-resource or image-pin changes.
   This PR's optional credential references allow source to merge without
   activating the new provider, but Argo may still roll the pod template.
2. Let `ecosystem-litellm` reconcile the new cloud alias and local fallback. Verify
   config and readiness; do not make a direct cloud request or change providers
   outside the gateway.
3. Merge/reconcile only the **bootstrap** sealed output first through
   `platform-ecosystem-resources`. Its existing Sync hook
   `onemind-model-clients-v1` creates/updates the supplied dedicated Hermes key.
   Verify Job success and its restricted model list in the authorized LiteLLM UI;
   do not display the key. A missing provider credential or unauthorized route
   must be resolved before consumer activation. No hand-applied Kubernetes secret.
4. After the credential is registered, merge/reconcile the consumer sealed output
   and Kustomize registration through `platform-legacy`. Creating a Secret alone
   does not update an existing process environment. The config commit changes the
   Pod template via a non-secret activation annotation, for example
   `onemind.io/model-credential-revision: "hermes-v1"`, so Argo rolls Hermes once.
   Review that annotation in the credential PR; no manual rollout restart.
5. Check `kubectl --context do-nyc3-omos-cloud -n hermes rollout status deploy/hermes-legacy`
   and inspect bounded `-c merge-litellm-provider` and `-c hermes` logs. No full
   config/env dump: confirm selected alias/endpoint with a redacted read only.
   `hermes-data` must retain its UID, mounts, persona and NATS leaf identity.
6. From Matrix or Agent Console, Zeus sends a unique short prompt and verifies a
   reply from Legacy. Confirm LiteLLM observed alias/provider without logging
   prompt contents or keys. Request a harmless tool query to verify tool calling.
7. Verify `local-tools` explicitly through the gateway with the restricted key.
   In an isolated test, exercise the configured cloud-error fallback; do not
   deliberately disable the shared production OpenRouter credential. Record
   response latency, actual model/provider, context behavior and degraded tool
   quality. No claim of offline fallback acceptance until this passes.
8. Verify Hermes can reach TCP/4000 and that a pod without Legacy's exact labels
   cannot use this newly added ingress rule. Existing broader namespace grants
   are unchanged. Read-only live policy inspection found no Hermes-wide egress
   restriction; the added rule requires both namespace and pod labels.

Rollback is a reviewed Git change restoring prior provider routing on the PVC
through the same merge mechanism (or a verified secure config restore), followed
by ordinary Argo reconciliation. Removing the init container alone does not undo
its earlier persistent config edit. Keep the virtual key and gateway route until
rollback is verified; no PVC deletion, database operation or credential rotation
is part of this source PR.

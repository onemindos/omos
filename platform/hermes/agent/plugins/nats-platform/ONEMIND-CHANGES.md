# OneMind changes to nats-platform

Upstream: `synadia-ai/hermes-nats-gateway` 0.1.2 (Apache-2.0, `LICENSE` kept). This directory holds the
runtime files only (no tests, lockfile or CI). The base is the copy running in `hermes-legacy` on
2026-11-01, which already carried the approval-import fix below. See the Git history of `adapter.py`
for the exact diff.

| # | Change | Why |
|---|---|---|
| 0 | `from tools.approval_context import set_current_session_key, reset_current_session_key` | Hermes moved these out of `tools.approval`. Legacy applied this as a startup hot-patch (`patch_nats_platform.py`, `nats-platform-patch-cm.yaml`). v0.22.5 defines them only in `tools.approval_context`, so the fix is still needed; it is now an ordinary reviewed line, not a patch |
| 1 | `_connect_options()`: `NATS_USER`, `NATS_PASSWORD`, `NATS_CA_FILE`, `NATS_TLS_SERVER_NAME`, `NATS_CLIENT_NAME` are passed to `nats.connect` | Upstream assumes a local leaf. Legacy now connects directly to the cluster bus over TLS with a password. With none of the variables set the behaviour is upstream's |
| 2 | `_refuse_duplicate_identity()` (`HERMES_NATS_REFUSE_DUPLICATE`), and `_warn_if_identity_already_live` returns a bool | Upstream starts a second body behind a live identity and splits the prompts between them. With the flag the adapter stays down until the identity is free |
| 3 | `__init__.py` appends `/opt/data/.local/gateway-deps` to `sys.path` | The `hermes` launcher unsets `PYTHONPATH` and the venv has no user site, so the init container's packages cannot be reached any other way (found at the live cutover: with a `PYTHONPATH` variable the gateway silently skipped the platform) |

When upstream changes, re-vendor and re-apply these four, and keep the diff small.

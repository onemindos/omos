# Hermes plugins

The plugins installed on the Hermes volume (`/opt/data/plugins/`) and where each one's source lives. Checked 2026-10-16, NATS rows updated 2026-11-01.

| Plugin on the volume | Source | Notes |
|----------------------|--------|-------|
| `nats-platform` | This repo: `platform/hermes/agent/plugins/nats-platform/` (upstream `synadia-ai/hermes-nats-gateway` 0.1.2, Apache-2.0, vendored; changes in its `ONEMIND-CHANGES.md`) | Vendored from the running Legacy pod. The approval-import line the hot-patch (`nats-platform-patch-cm.yaml`) fixed is now in `adapter.py`; the patch is only still used by `hermes-legacy`. Installed on onemind-hermes by `platform/hermes/onemind-hermes/nats/apply.sh`; hermes-legacy still has its own hand-installed copy |
| `onemind-vision` | This repo: `platform/hermes/agent/plugins/onemind-vision/` (next to this file) | Byte-identical to the live copy. Installed by hand, so it is a GitOps gap. Formerly `onemind-maven` |
| `nats-context` | This repo: `platform/hermes/agent/plugins/nats-context/` (OneMind-authored; no upstream repo was found) | Copy from the running Legacy pod plus `conn.py` (credentials and CA from the environment). Installed on onemind-hermes only with `nats/apply.sh --with-context` |
| `platforms/` (the TAK platform) | Not in Git. `tak-platform` 0.1.0, 469 lines, on the Legacy volume | Deliberately not ported to onemind-hermes: PR #585 replaces it with a TAK-native client. It needs the TAK admin certificate on the volume and disables certificate verification |
| `memory-wiki` | Upstream `NousResearch/hermes-memory-wiki` (MIT), per `.install-metadata.json` | Dashboard only, no NATS |

To rebuild a plugin on hermes-legacy, reinstall the upstream or repo version and let the init container apply the patch. Do not edit plugins in place on a volume without exporting the change.

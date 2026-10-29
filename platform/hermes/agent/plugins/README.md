# Hermes plugins

The plugins installed on the Hermes volume (`/opt/data/plugins/`) and where each one's source lives. Checked 2026-10-16.

| Plugin on the volume | Source | Notes |
|----------------------|--------|-------|
| `nats-platform` | Upstream: `synadia-ai/hermes-nats-gateway` (Apache-2.0), version 0.1.2 | Not vendored here. The live copy differs from the upstream code the old `legacy/` copy held by two lines in `adapter.py`, which are the approval-import patch applied by the init container in `deployment.yaml` (`nats-platform-patch-cm.yaml`). How the plugin is first installed is not recorded; plan step 2 (seeding an empty volume) will capture it |
| `onemind-vision` | This repo: `platform/hermes/agent/plugins/onemind-vision/` (next to this file) | Byte-identical to the live copy. Installed by hand, so it is a GitOps gap |
| `onemind-bus` | This repo: `platform/hermes/agent/plugins/onemind-bus/` | Not installed yet. Bus tools for I2/D5/D11 (entity and layer reads, policy proposals, `cmd.map.*`); retires `nats_publish`. Install steps in its README |
| `nats-context` | Probably `onemindos/hermes-plugin-nats-context` | Not compared with the live copy |
| `platforms/` (the TAK platform) | Probably `onemindos/hermes-plugin-tak` | Not compared with the live copy |
| `memory-wiki` | Unknown (upstream, a separate repo, or hand-written) | |

To rebuild a plugin, reinstall the upstream or repo version and let the init container apply the patch. Do not edit plugins in place on the volume without exporting the change.

# Hermes plugins

The Hermes plugins whose source is in Git, and what `onemind-hermes` does with each. The old `hermes-legacy` volume, which
held hand-installed copies and two third-party plugins, was retired on 2026-11-03. Facts about `onemind-hermes` come from
the retirement audit (2026-11-01) and `nats/apply.sh`.

| Plugin | Source | State on onemind-hermes |
|--------|--------|-------------------------|
| `nats-platform` | `plugins/nats-platform/` (upstream `synadia-ai/hermes-nats-gateway` 0.1.2, Apache-2.0, vendored; changes in its `ONEMIND-CHANGES.md`). The approval-import line that the old start-up hot-patch fixed is now in `adapter.py` | Installed by `platform/hermes/onemind-hermes/nats/apply.sh` and listed in `plugins.enabled`. Gives Legacy its identity on the bus |
| `nats-context` | `plugins/nats-context/` (OneMind-authored, with `conn.py` so it takes credentials and CA from the environment) | Vendored, **not enabled**. Install only with `nats/apply.sh --with-context`, after the raw-feed and generic `nats_publish` decisions (PRs #562 and #564) |
| `onemind-vision` | `plugins/onemind-vision/` (formerly `onemind-maven`) | Source only, not installed, no install step |
| `tak-platform` (`platforms/`) | Not in Git | Not ported. PR #585 replaces it with a TAK-native client |
| `memory-wiki`, `hermes-project-stewardship` | Upstream, third party | Not installed; dropped with the old volume |

Do not edit plugins in place on the volume without exporting the change here.

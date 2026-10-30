# Plugin ecosystem catalog

Shows every plugin, adapter and bridge OneMind builds or runs, what it connects to and where it runs.
System: `onemind-plugin-ecosystem` (Domain `onemind-so`). 137 entities as of 2026-10-29.

## Where to look (tailnet Backstage, https://onemind-preview.tail717ba0.ts.net)

- `/catalog/default/system/onemind-plugin-ecosystem` : the system page; its Diagram tab draws every plugin, subject and host with their relations.
- `/catalog?filters[kind]=component&filters[type]=hermes-plugin` : the plugin list by type. Types: `tak-server-plugin`, `atak-plugin`, `hermes-plugin`, `mcp-adapter`, `bridge`, `backstage-plugin`. The Lifecycle filter separates `production` (live or installed), `parked` and `planned`; the Tag filter has `plugin`, `nats`, `third-party`, `flagged`, `device`, `cluster`, `edge`.
- A plugin's Relations tab (catalog-graph, already in the deployed image) shows what it depends on: host workload, NATS subject Resources, APIs.
- Resources of type `nats-subject` (flagged `locked` or `off-grammar`), `nats-stream` and `plugin-host` show the other side: open a subject to see every plugin that publishes or subscribes to it.

## Entity model

| Thing | Kind / type | Notes |
|---|---|---|
| Plugin, adapter, bridge | Component, `spec.type` = plugin type | `lifecycle`: live and installed = `production`, `parked`, `planned`. `subcomponentOf` the app it extends (`tak-server`, `hermes-legacy`, `onemind-mcp`, `backstage`) |
| Where it runs | `dependsOn` the existing workload Component (`k8s-...`) or a `plugin-host` Resource (ATAK on Android, iTAK, Mac edge leaf, Hermes plugin volume) | Hosts depend on `tak-server` or `nats`, so the chain plugin, host, system is visible |
| What it connects to | `onemind.io/connects-to` annotation (comma list) plus `dependsOn` for catalogued systems | |
| Bus subjects | Resource `nats-subject` (`nats-subject-*`), depends on Resource `nats-stream` (`nats-stream-*`), which depends on Component `nats` | Plugins `dependsOn` every subject they publish or subscribe to; direction is in `onemind.io/nats-publishes` / `onemind.io/nats-subscribes` |
| Bus contracts | Existing APIs `tak-cot-bridge`, `nats-agent-protocol`, `nats-events`, `nats-commands`, `nats-telemetry`, `nats-intel` via `providesApis` / `consumesApis` | No new API entities |

Annotations (all `onemind.io/`): `plugin-type`, `plugin-subtype`, `plugin-status` (live, parked, planned, installed), `plugin-origin` (onemind, upstream, third-party), `runs-on-kind` (cluster, device, edge), `runs-on`, `connects-to`, `nats-publishes`, `nats-subscribes`, `plugin-version`, `plugin-source`, `plugin-build`, `plugin-verified`. Labels repeat type, status, origin, runs-on-kind and (on subjects) `subject-grammar` for API queries such as `metadata.labels.onemind.io/plugin-status=live`.

Names: `<type>-<slug>`, for example `tak-plugin-nats-bridge`, `hermes-plugin-nats-context`, `mcp-adapter-prometheus`, `atak-plugin-taktalk`, `bridge-mumble`, `backstage-plugin-onemind-services`.

## Files

| File | Role |
|---|---|
| `inventory.yaml` | Reviewed source: hosts, streams, subjects, plugin records, ATAK package list. Hand-edited. |
| `../plugins.yaml` and `services/onemind-portal/catalog/plugins.yaml` | Generated, identical. Never edit by hand. |
| `scripts/generate_plugin_catalog.py` | Renders both. MCP adapters come from `platform/hermes/mcp/apps.yaml` and Backstage plugins from `services/onemind-portal/plugins/*/package.json`, so they cannot drift. `--check` (run by `scripts/tests/test_gitops_plugin_catalog.py`) fails on stale output, dangling references, bad names, labels or tags. |

Rollout: `platform/backstage` generates ConfigMap `backstage-catalog-plugins` (separate from `backstage-catalog`, which is near the 1 MiB limit), mounted at `/app/catalog-plugins`, and `app-config.cluster.yaml` registers `/app/catalog-plugins/plugins.yaml`. No image change. To add or change a plugin: edit `inventory.yaml`, run the generator, commit all three files.

## Future: generating from plugin.yaml

The per-plugin `plugin.yaml` contract is being defined in `docs/architecture/hermes-connection-plugins.md`. The records in `inventory.yaml` already use the field names below, so a generator step can read `**/plugin.yaml` and append records instead of (or merged over) the hand-written ones. Mapping, with the inventory field on the left:

| plugin.yaml / inventory field | Catalog output |
|---|---|
| `name`, `title`, `description` | `metadata.name` (prefixed with the type), `title`, `description` |
| `type` | `spec.type`, label and annotation `plugin-type`, tag |
| `subtype` (Hermes `kind`: platform, standalone, backend, dashboard) | annotation `plugin-subtype` |
| `status` live / parked / planned | `spec.lifecycle` production / parked / planned, label and tag |
| `origin` onemind / upstream / third-party | label `plugin-origin` |
| `version` | annotation `plugin-version` |
| `source` (repo path) | `plugin-source`, `backstage.io/source-location` |
| `build` | `plugin-build` |
| `runsOn.kind`, `runsOn.location`, `runsOn.host` | `runs-on-kind`, `runs-on`, and `dependsOn` the host entity ref |
| `connectsTo` | `connects-to` |
| `publishes`, `subscribes` (subject slugs) | `nats-publishes` / `nats-subscribes` and `dependsOn` `resource:default/nats-subject-<slug>` |
| `providesApis`, `consumesApis` | `spec.providesApis` / `consumesApis` |
| `dependsOn`, `subcomponentOf` | same fields |
| `verified` | `plugin-verified` |

A plugin that uses a subject not yet in `subjects:` needs one new entry (slug, subject, grammar, streams, description) so the Resource exists. Keep `status` honest: `live` only when a pod, device or volume was read, otherwise `parked` or `planned`, with the evidence in `verified`.

## Findings recorded here (2026-10-29)

- Not in the repo, so not catalogued: any OneMind-built ATAK plugin, an Obsidian plugin (Obsidian is a Nextcloud-synced vault, plus a Hermes `obsidian` skill), VR-TAK (only docs mention a distribution server that is not in the tree), a Home Assistant plugin beyond the omos-bridge Go adapter. iTAK has no plugin API; its host entry exists so it shows what feeds it.
- `onemind-maven` was renamed `onemind-vision`; the live volume holds only the latter.
- The Hermes volume also holds `hermes-project-stewardship`, which `platform/hermes/agent/plugins/README.md` does not list. `nats-context` and `tak-platform` have no source in this repo.
- Off-grammar subjects (`voice.transcript`, the Matrix bridge's `alerts.*` / `infra.health.*` / `legacy.*`, the omos-bridge `omos.*` family) are tagged `off-grammar`; the omos-bridge Home Assistant adapter and the `ha-nats-bridge` skill disagree on subjects.
- `cot-tap` has no Argo manifest or pod; the NATS-to-Matrix bridge is in `platform/matrix/base` but not in the Argo overlay and has no pod.
- Third-party ATAK plugins (32) come from the July 2026 adb inventory and were not re-checked; `com.atak.nightvision` is tagged `flagged`.

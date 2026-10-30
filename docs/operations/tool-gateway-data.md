# Tool gateway: keyless data backends (Phase 3d)

Design rows 8, 13, 16, 20, 21, 23 in [`docs/architecture/onemind-tool-gateway.md`](../architecture/onemind-tool-gateway.md). Recipe:
[`tool-gateway.md`](tool-gateway.md). Every credential used here already existed in the cluster or is generated and sealed; the owner creates nothing.

| App | Route | Target | Group | Tools (class) | Credential |
|---|---|---|---|---|---|
| STAC | OpenAPI, spec trimmed to GET in Git | `stac` | `data` | 5 read: `list_collections`, `get_collection`, `list_items`, `get_item`, `search` | none (read API only; the writer is not a target) |
| onemind-core | OpenAPI, hand-written read subset | `core` | `data` | 6 read: `list_capabilities`, `list_objects`, `list_links`, `list_agents`, `list_missions`, `get_mission` | `READ_API_KEY` resealed as `gateway-target-core` |
| ClickHouse | official `ClickHouse/mcp-clickhouse` 0.7.0, Deployment `mcp-clickhouse` | `clickhouse` | `data` | 3 read: `list_databases`, `list_tables`, `run_query` | user `legacy_ro` (`readonly=1`), password sealed into `clickhouse/clickhouse-legacy-ro` and `tool-gateway/mcp-clickhouse-credentials`; MCP bearer in `gateway-target-clickhouse` |
| Ollama | existing `onemind-mcp` adapter | `onemind` | `ops` | 2 read: `ollama_models`, `ollama_running_models` | the adapter's own bearer |
| OpenCost | existing `onemind-mcp` adapter | `onemind` | `ops` | 1 read: `opencost_allocation` | same |
| Tile38 | new adapter `tile38` in `services/onemind-mcp` | `onemind` | `ops` | 6 read, 1 write (`tile38_set_point`), 2 destructive (`tile38_delete_object`, `tile38_drop_collection`) | none (Tile38 has no auth) |

Why Ollama, OpenCost and Tile38 are in `ops`: classification and group membership are per target, and these tools are served by the
`onemind` target, which the design keeps in `ops`. Giving them their own target would need a second `onemind-mcp` deployment.

## ClickHouse user

`legacy_ro` is declared in `platform/data/clickhouse/clickhouse.yaml` (the Altinity CHI, applied by `platform-observability`): password from
Secret `clickhouse-legacy-ro` through `valueFrom.secretKeyRef`, profile `legacy_ro` with `readonly=1`, 30 s execution, 1 GB memory and
100,000 result rows (`result_overflow_mode: break`), network `10.0.0.0/8`. The ClickHouse namespace has no ingress policy, and adding one
would cut off Grafana, Redpanda Connect and the backups, so only the gateway-side egress and the `mcp-clickhouse` policies exist.

## Tile38 rollout (second step)

The adapter is merged disabled: the live `onemind-mcp` image does not know the `tile38` adapter and would fail at start if `apps.yaml`
named it. After this PR merges, CI pushes `ghcr.io/onemindos/onemind-mcp:sha-<commit>`; a follow-up PR pins that digest in
`platform/hermes/mcp/deployment.yaml`, adds `TILE38_URL=http://tile38.onetak.svc.cluster.local:9851` and enables the `tile38` app
(`allow_writes: true`, `allow_destructive` stays false, so the two destructive tools are not registered until the execute gate exists).

## OpenAPI specs

`scripts/gateway/trim_openapi_get.py` turns `platform/tool-gateway/openapi/<name>.recipe.yaml` plus its source spec into ConfigMap
`gateway-openapi-<name>` (`manifests/openapi-<name>.configmap.yaml`), mounted by `values.yaml` at `/specs/<name>/`. The
gateway's parser accepts OpenAPI 3.0 only (it rejects the 3.1 `exclusiveMinimum` numbers of the STAC spec), so the script emits plain 3.0
parameter schemas and a generic response. A test fails when the committed ConfigMaps are stale.

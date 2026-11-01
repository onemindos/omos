# nats-context in Git

OneMind-authored (`author: OneMind OS`); the source repo was never confirmed, so this directory is now the
source. The base is the copy running in `hermes-legacy` on 2026-11-01 (`state.py`, `formatter.py`,
`plugin.yaml` and `__init__.py` are unchanged).

Changes:

| File | Change | Why |
|---|---|---|
| `conn.py` (new) | `nats_url()` and `connect_options()` read `NATS_URL`, `NATS_USER`, `NATS_PASSWORD`, `NATS_CA_FILE`, `NATS_TLS_SERVER_NAME`, `NATS_CONTEXT_CLIENT_NAME` | The plugin dialled `nats://localhost:4222` (the leaf sidecar) with no credentials and no TLS. onemind-hermes has no sidecar |
| `subscriptions.py` | `nats.connect(..., **connect_options())` | same |
| `tools.py` | the hard-coded `NATS_URL = "nats://localhost:4222"` and 9 `nats.connect(NATS_URL)` calls use `nats_url()` and `connect_options()` | same |
| `__init__.py` | appends `/opt/data/.local/gateway-deps` to `sys.path` before the imports | The `hermes` launcher unsets `PYTHONPATH`; same reason as `nats-platform` |

With the environment unset the behaviour is what it was.

What it does on the bus (measured on Legacy's leaf, 2026-11-01): one core subscription per pattern in
`subscriptions.SUBJECTS` (`tak.cot.>`, `ent.>`, `tel.>`, `intel.>`, `det.>`, `evt.>`, `evt.threat.>`,
`evt.mission.>`, `evt.task.>`, `cmd.alert.>`, `agents.hb.>`, `svc.bridge.>`, `sys.>`), one `$SYS.REQ.SERVER.PING`
request a minute, and nine tools (`nats_publish`, `nats_request`, `nats_peek`, `nats_jetstream_fetch`,
`nats_kv_get`, `nats_kv_put`, `nats_kv_list`, `nats_server_info`, `nats_list_streams`). Hermes loads it in two
processes (gateway and dashboard), so Legacy shows two such connections. It owns no JetStream consumer.

It is **not enabled by default** on onemind-hermes: it receives the raw feeds (622k of the 684k messages the
Legacy leaf carried in two days were `tak.cot.>`) and gives the model a generic `nats_publish`, both of which
PR #562 and PR #564 take away from Hermes on purpose. See `platform/hermes/onemind-hermes/MIGRATION.md`.

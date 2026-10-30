# execute-gate

AuthZEN decision point for the OneMind tool gateway `/execute/mcp` route (design:
[`docs/architecture/onemind-tool-gateway.md`](../../docs/architecture/onemind-tool-gateway.md) section 5.3, Phase 4).
Python 3.12, stdlib HTTP server and Matrix client, `nats-py` for the KV bucket and bus events, `PyYAML` for the policy file.
**Not deployed**: code and image build only; manifests, KV CR, Matrix bot and `extAuthz` wiring come with the gateway core.

## Behaviour

- `POST /access/v1/evaluation` (AuthZEN): subject `hermes`, action `tools/call`, resource id `<target>/<tool>`,
  `context.arguments` (JSON value or JSON text). Response `{"decision": bool, "context": {"reason", "request_id"}}`. Default and every
  error path is `false`.
- Approval key = sha256 of resource + canonical argument JSON (sorted keys, no whitespace, integral floats as ints).
- No approval: create a pending record (4-hex id, TTL 600 s), post it to the approvals room, deny. Same call while pending: deny, no
  second post (re-posts only if the first post failed).
- Approval: only a plain `m.text` message in the approvals room from exactly `OWNER_MXID`, matching `^(execute|override) <id>$`.
  `execute` approves `destructive`, `override` approves `physical`; the wrong word, wrong id, expired id, edits, other senders, the
  bot and other rooms never approve. Event ids are de-duplicated.
- Next matching call: allowed once; the record is deleted by compare-and-set, so a replay is denied and opens a new request.
- Classification: anything not matching `physical` globs in the policy file is `destructive` (default deny).
- Audit: JSON line on stdout for every decision (argument values never logged), and `evt.toolgw.approval.{requested,granted,denied,consumed}`
  on NATS (fits `evt.{domain}.{entity}.{event}`). The `EVT` stream does not yet cover `evt.toolgw.>`; core publish works without it,
  persistence needs that subject added in `platform/nats/streams/evt.yaml`.
- `GET /healthz`, `GET /readyz` (store connected and Matrix sync thread alive).

## Configuration (environment)

| Variable | Default | Notes |
|---|---|---|
| `APPROVALS_ROOM_ID` | required | internal room id (`!...:onemindos.dev`) of the private, **unencrypted** approvals room |
| `MATRIX_ACCESS_TOKEN` | required | bot `@gate:onemindos.dev` (Secret) |
| `MATRIX_HOMESERVER` | `http://synapse.matrix.svc.cluster.local:8008` | confirm service name at deploy |
| `OWNER_MXID` | `@zeus:onemindos.dev` | only sender whose approval counts |
| `BOT_MXID` | whoami | must differ from owner |
| `NATS_URL`, `NATS_USER`, `NATS_PASSWORD` | in-cluster NATS | KV bucket `exec-approvals` (`KV_BUCKET`) |
| `STORE` | `nats` | `memory` for local runs only |
| `POLICY_FILE` | none | see `policy.example.yaml` |
| `ALLOWED_SUBJECTS` | `hermes` | comma list |
| `APPROVAL_TTL_SECONDS` | `600` | |
| `GATE_AUTH_TOKEN` | none | optional bearer the gateway sends |
| `PORT` | `8080` | |

## Tests

```
pip install -r requirements-dev.txt && python -m pytest -q
```

The NATS adapter (`NatsKvStore`) needs a live server and is covered by the Phase 4 live verification (a) to (h) in the design.

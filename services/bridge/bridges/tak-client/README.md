# tak-client: TAK Server streaming client bridge

`tak-client` joins TAK Server the way an ATAK device does, on the TLS streaming input
(8089) with its own client certificate, and connects that stream to the authority bus:

```
TAK Server ──► tak.cot.out     CoT XML, acked JetStream publish (stream COT), Nats-Msg-Id dedup
TAK Server ◄── tak.cot.inject  queue group tak-cot-inject, shared with NatsBridgePlugin 1.1.0+
```

It is feature A21 of the OneMind build (the bus-to-TAK Server leg) and replaces
`services/takserver/cot-tap`. Both subjects already exist in the locked grammar
(`tak.cot.inject` "CoT → TAK inject", `tak.cot.>` "CoT from TAK") and in stream `COT`
(`platform/nats/streams/cot.yaml`); no new subject is introduced.

## Why it exists

Nothing has published `tak.cot.out` in the 30-day retention of `COT`, and nothing consumes
`tak.cot.inject` into TAK (onemind-xr `docs/research/omos-live-state.md` §6):

- The live `takserver-plugins` JVM starts with
  `-Dloader.path=WEB-INF/lib-provided,WEB-INF/lib,WEB-INF/classes` (no `lib/`) and logs no
  NatsBridgePlugin line, so the plugin JAR is not loaded (details in
  `services/takserver/plugin-src/README.md`).
- Even when loaded, TAK Server's cluster-mode plugin manager does not deliver CoT to a
  plugin's `onMessage` (`services/takserver/cot-tap/README.md`), so the plugin cannot be the
  TAK → bus leg.
- cot-tap, the stopgap for that, is not deployed (`kubectl get deploy cot-tap` → NotFound).

## Protocol

From the published spec (atak-civ `takproto/README.txt`), not from TAK Server source:

| Step | What tak-client does |
|---|---|
| Connect | TCP + TLS 1.2+ with the client certificate. The server certificate is verified against `TAK_CA_FILE`; the hostname too when `TAK_SERVER_NAME` is set (Common Commo checks the chain only, which is the default here because TAK Server certificates rarely name the in-cluster Service) |
| XML streaming | Splits the stream on `</event>`, drops anything before `<event` (declaration, whitespace). Each event it writes starts with an XML declaration and a newline. One event is at most 1 MiB |
| Negotiation | On `t-x-takp-v` offering `<TakProtocolSupport version="1"/>` (and `TAK_PROTOCOL=auto`) it sends `t-x-takp-q` `<TakRequest version="1"/>` and stops writing XML. On `t-x-takp-r` `status="true"` both directions switch to v1 streaming frames (`0xbf`, varint length, `TakMessage`). On `false` it stays on XML. With no answer in 60 s it disconnects, and later sessions stay on XML |
| Keep-alive | After Commo: ping (`t-x-c-t`, uid `<id>-ping`) after 15 s without receiving, every 4.5 s, and drop the connection after 25 s of silence |
| Reconnect | Exponential backoff with jitter, 1 s to 30 s |
| Received CoT | Stream plumbing (`t-x-takp-*`, `t-x-c-t`, `t-x-c-t-r`) is dropped. Everything else goes to `tak.cot.out`: the original bytes on an XML stream, re-rendered XML on a v1 stream. Headers `Nats-Msg-Id: takc-<hash of uid, type, time, lat, lon>` (same fields as the mesh gateway's `cot-` id, different prefix), `Tak-Client`, `Tak-Proto` (0 or 1), `Tak-Server`, and `Tak-Undecoded-Extensions` when set |
| Echo filter | Every `tak.cot.inject` message is remembered for `ECHO_TTL` (uid, type, time). When TAK relays such an event back, it is dropped instead of being republished as if a TAK user had sent it |
| Inject | Joins queue group `INJECT_QUEUE` on `tak.cot.inject` only while the TAK stream is up, so NatsBridgePlugin (same group) takes every message while tak-client is down. Each message must parse as one CoT event; it is re-rendered, so a second event appended to a message is never sent. Stream-control types are refused. The queue is bounded (1024, oldest dropped) and items older than 30 s are dropped, not sent late |

The TAK Protocol v1 codec (`takproto.go`, `cotxml.go`) is copied unchanged from the mesh
gateway (`services/bridge/bridges/tak`, PR #551, Apache-2.0); move both copies into one
shared module after both PRs merge.

## Configuration

| Variable | Default | Meaning |
|---|---|---|
| `NATS_URL`, `NATS_USER`, `NATS_PASSWORD` | `nats://127.0.0.1:4222` | Bus connection. In the cluster `nats://nats.nats.svc.cluster.local:4222` with credentials from a Secret |
| `TAK_ADDR` | `takserver-messaging-service.takserver.svc.cluster.local:8089` | TAK Server streaming input |
| `TAK_CERT_FILE`, `TAK_KEY_FILE` | required | PEM client certificate and key (TAK user `omos-tak-client`) |
| `TAK_CA_FILE` | required | PEM CA that signed the TAK Server certificate |
| `TAK_SERVER_NAME` | empty | Also verify the server hostname |
| `TAK_INSECURE_SKIP_VERIFY` | `false` | Do not verify the server (logged as a warning; local tests only) |
| `TAK_CLIENT_ID` | `omos-tak-client` | CoT uid for pings and negotiation; `Tak-Client` header |
| `TAK_PROTOCOL` | `auto` | `auto` accepts a v1 offer; `xml` never negotiates |
| `INJECT_ENABLED`, `INJECT_QUEUE` | `true`, `tak-cot-inject` | Inject leg |
| `ECHO_SUPPRESS`, `ECHO_TTL` | `true`, `2m` | Echo filter |
| `OUT_JETSTREAM` | `true` | Acked publish (needs stream `COT`); `false` = core publish |
| `HTTP_ADDR` | `:8080` | `/healthz`, `/readyz` (TAK and NATS connected), `/metrics` (Prometheus text) |

## Tests

```sh
go vet ./... && go test -race -count=1 ./...      # NATS_SERVER_BIN or nats-server on PATH for integration
PYTHON=/path/to/venv/bin/python ./test/e2e.sh     # binary against test/peer.py (needs takproto, nats CLI, openssl)
```

- Unit: v1 stream framing, XML splitting across reads, oversize and malformed frames,
  negotiation parsing, control types, echo filter, dedup id, config validation, bounded
  inject queue.
- Scripted TAK-stream peer (`peer_test.go`, mutual TLS with a per-run test CA): XML stream,
  v1 negotiation (accept, decline, no answer), `TAK_PROTOCOL=xml`, pings and reset,
  reconnect after a server drop, server and client certificate checks, hostname check.
- Integration with nats-server: TAK → `COT` in XML and v1 with headers, replay dedup,
  inject, echo suppression, the queue group shared with a plugin stand-in (and the plugin
  taking everything while TAK is down), smuggled second event, NATS restart mid-stream
  (no loss, no duplicates).
- `test/e2e.sh`: the built binary against `test/peer.py`, an independent Python peer that
  encodes v1 with the `takproto` package (MIT), in XML and v1 mode.

No test needs a real TAK Server. The live check is in `platform/takserver/tak-client/README.md`.

## Retiring cot-tap

1. Merge this bridge; seal its certificate and NATS credentials; sync
   `platform-takserver-tak-client` (`platform/takserver/tak-client/README.md`).
2. Confirm `tak.cot.out` messages with `Tak-Client: omos-tak-client` in `COT`.
3. Set `natsProducerEnabled: false` in the NatsBridgePlugin config once the 1.1.0 JAR runs,
   so tak-client is the only `tak.cot.out` producer.
4. After a week without gaps, delete `services/takserver/cot-tap/` and the `cot-tap-certs`
   Secret if one exists, and revoke the `omos-cot-tap` TAK user and certificate.

Not in this version: Data Sync missions (FEATURES C7) and video feed registration
(`POST /Marti/api/video`), both planned for this client over 8443.

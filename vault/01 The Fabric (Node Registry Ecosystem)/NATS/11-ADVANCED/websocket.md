---
title: "WebSocket"
section: "11-ADVANCED"
type: note
tags: [nats, websocket, browser, ws, wss, dashboard, nats.ws]
source: https://docs.nats.io/running-a-nats-service/configuration/websocket
created: 2026-06-14
---

# WebSocket

NATS speaks **WebSocket** natively — a browser connects directly to `nats-server` over WS/WSS and gets full NATS (pub/sub, request-reply, JetStream, KV). No gateway, no REST shim, no Socket.IO layer. This is what makes real-time web dashboards on NATS trivial.

See also: [[NATS/08-CLIENTS-SDKS/javascript|JavaScript Client]] · [[NATS/11-ADVANCED/mqtt|MQTT]] · [[NATS/05-SECURITY/tls-encryption|TLS]]

---

## Why it exists

Browsers can't open raw TCP sockets — so a normal NATS client (which needs TCP) can't run in a browser. WebSocket is the browser-native, firewall-friendly transport. With WS enabled, the **same `nats-server`** serves browser clients alongside TCP clients, all on one bus. A dashboard subscribes to `telemetry.>` and updates live as services publish.

---

## Enabling it

```text
websocket {
  listen: "0.0.0.0:8080"
  # Production: terminate TLS here (browsers want WSS)
  tls {
    cert_file: "/etc/nats/server-cert.pem"
    key_file:  "/etc/nats/server-key.pem"
  }
  # no_tls: true            # only for localhost/dev behind a proxy
  # same_origin: false
  # allowed_origins: ["https://command.onemind.local"]
  compression: true         # permessage-deflate
  handshake_timeout: "2s"
}
```

- **TLS strongly recommended** — browsers prefer `wss://`, and mixed-content rules block `ws://` from an HTTPS page.
- **`allowed_origins`** — restrict which web origins may connect (CORS-style defense).
- **`compression`** — permessage-deflate, useful for chatty text payloads.

---

## Connecting from the browser

The JavaScript client (`nats.ws`, now folded into the unified `@nats-io/nats` / `nats.js` family) auto-selects WebSocket in a browser:

```javascript
import { connect } from "@nats-io/nats";   // or nats.ws build

const nc = await connect({
  servers: "wss://command.onemind.local:8080",
  // user/pass, token, or creds as configured
});

// Subscribe — live updates straight into the UI
const sub = nc.subscribe("telemetry.barn.>");
for await (const m of sub) {
  render(m.subject, m.json());
}

// Request-reply works too
const res = await nc.request("svc.status", undefined, { timeout: 1000 });

// JetStream / KV also available over WS
```

Same API as the Node/Deno client — the transport is the only difference. See [[NATS/08-CLIENTS-SDKS/javascript|JavaScript Client]].

---

## What works over WebSocket
- **Everything the client protocol supports**: Core pub/sub, queue groups, request-reply, headers.
- **JetStream**: publish, consumers, **KV**, **Object Store** — all usable from the browser.
- **Auth**: token, user/pass, NKEY/JWT creds.

The WS transport simply frames the same [[NATS/10-PROTOCOL-INTERNALS/wire-protocol|wire protocol]] inside WebSocket messages — so feature parity is automatic.

---

## Gotchas
- **Mixed content** — an HTTPS page may only open `wss://` (TLS). Plan certs accordingly.
- **Origins** — lock `allowed_origins` in production; an open WS endpoint is reachable from any web page.
- **Auth in the browser** — never ship an operator/account seed to a browser. Use scoped **User** creds or short-lived tokens, ideally minted server-side. See [[NATS/05-SECURITY/jwt-nkeys-decentralized|JWT & NKEYS]].
- **Load balancers** — must support WebSocket upgrade and sticky enough behavior; set sane idle timeouts (NATS PING/PONG keeps it alive).

---

## OneMind Relevance (high)
- **OneMind Command** dashboards subscribe to live fabric subjects (telemetry, agent status, topology) directly over **WSS** — no backend polling layer. This is the standard the app spec assumes.
- Browser clients get the **same** JetStream/KV access as services, so a dashboard can read config from a KV bucket or tail a stream natively.
- Lock `allowed_origins` to the Command app origin; mint scoped browser creds — never expose broad fabric access to a web page.

## Sources
- WebSocket config — https://docs.nats.io/running-a-nats-service/configuration/websocket
- WebSocket (concepts) — https://docs.nats.io/running-a-nats-service/configuration/websocket/websocket_conf
- nats.ws — https://github.com/nats-io/nats.ws
- JS client — https://github.com/nats-io/nats.js

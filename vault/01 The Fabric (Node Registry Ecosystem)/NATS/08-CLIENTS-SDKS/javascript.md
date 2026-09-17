---
title: "JavaScript / TypeScript Client (nats.js)"
section: "08-CLIENTS-SDKS"
type: note
tags: [nats, javascript, typescript, nodejs, deno, browser, websocket]
source: https://github.com/nats-io/nats.js
created: 2026-06-14
---

# JavaScript / TypeScript Client (`nats.js`)

Runs in **Node.js**, **Deno**, **Bun**, and the **browser** (via WebSocket). The browser capability is unique — apps connect *directly* to the bus. Monorepo: `nats-io/nats.js` (packages: `@nats-io/nats-core`, `@nats-io/jetstream`, `@nats-io/kv`, `@nats-io/obj`, `@nats-io/services`, transports `@nats-io/transport-node` / `nats.ws`).

See also: [[NATS/08-CLIENTS-SDKS/README|Clients Index]] · [[NATS/11-ADVANCED/websocket|WebSocket]]

---

## Install & Connect

```shell
npm install @nats-io/nats-core @nats-io/transport-node
```

```typescript
import { connect } from "@nats-io/transport-node";  // Node
// browser: import { wsconnect } from "@nats-io/nats-core"; (WebSocket)

const nc = await connect({
  servers: "nats://localhost:4222",
  name: "my-service",
  // creds / token / user+pass / tls options
  maxReconnectAttempts: -1,
  reconnectTimeWait: 2000,
});
```

Browser (WebSocket to `:8080`):
```typescript
import { wsconnect } from "@nats-io/nats-core";
const nc = await wsconnect({ servers: "wss://nats.example.com:8080" });
```

---

## Core Pub/Sub & Request-Reply

```typescript
import { StringCodec } from "@nats-io/nats-core";
const sc = StringCodec();

nc.publish("orders.new", sc.encode(JSON.stringify({ id: 1 })));

const sub = nc.subscribe("orders.*");
(async () => { for await (const m of sub) console.log(m.subject, sc.decode(m.data)); })();

// queue group
nc.subscribe("orders.new", { queue: "workers", callback: (err, m) => {/*...*/} });

// request-reply
const r = await nc.request("time.now", sc.encode(""), { timeout: 1000 });
console.log(sc.decode(r.data));
```

> Subscriptions are **async iterables** — `for await (const m of sub)` is the idiomatic consume loop.

---

## JetStream

```typescript
import { jetstream, jetstreamManager } from "@nats-io/jetstream";

const jsm = await jetstreamManager(nc);
await jsm.streams.add({ name: "ORDERS", subjects: ["orders.>"] });

const js = jetstream(nc);
await js.publish("orders.new", sc.encode("..."), { msgID: "order-1" });  // dedup

// pull consumer
const c = await js.consumers.get("ORDERS", "dispatch");
const msgs = await c.fetch({ max_messages: 10 });
for await (const m of msgs) { /* process */ m.ack(); }
```

KV / Object Store / Services live in `@nats-io/kv`, `@nats-io/obj`, `@nats-io/services` with the same concepts.

---

## Browser Use (the differentiator)

With **WebSocket** enabled on the server (`websocket { port: 8080 }`), a browser app connects straight to NATS — no backend proxy:
- Live dashboards subscribing to `metrics.>`
- Collaborative UIs over pub/sub
- Real-time feeds without polling

OneMind's **live bus viewer** and the Flutter app's `NatsMonitorService` use this pattern (WebSocket subscriber to bus subjects). Secure browser connections with **TLS (`wss://`)** and scoped **account/permissions** — the browser is untrusted, so give it a least-privilege user (read-only subscribe to specific subjects).

---

## Notes

- **Codecs:** `StringCodec` / `JSONCodec` encode/decode payloads (which are `Uint8Array`).
- **Deno/Bun:** use the matching transport import; core API is identical.
- The JS client moved to a **scoped-package monorepo** (`@nats-io/*`) — older docs reference a single `nats` package; the split packages are current.

## Sources
- nats.js — https://github.com/nats-io/nats.js
- nats.ws (browser/WebSocket) — https://github.com/nats-io/nats.ws
- JetStream (js) — https://github.com/nats-io/nats.js/tree/main/jetstream
- NATS by Example (JS) — https://natsbyexample.com

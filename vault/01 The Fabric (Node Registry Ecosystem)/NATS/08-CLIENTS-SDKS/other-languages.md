---
title: "Other Languages — Java, C#, C, Ruby, Elixir & more"
section: "08-CLIENTS-SDKS"
type: note
tags: [nats, java, csharp, dotnet, c, ruby, elixir, clients]
source: https://docs.nats.io/using-nats/developer
created: 2026-06-14
---

# Other Official Clients

NATS maintains first-party clients across the ecosystem. All share the same concepts (connect → pub/sub → request-reply → JetStream → KV/OS → Micro); only syntax differs.

See also: [[NATS/08-CLIENTS-SDKS/README|Clients Index]]

---

## Java — `nats.java`

```java
// repo: nats-io/nats.java  (Maven: io.nats:jnats)
Connection nc = Nats.connect("nats://localhost:4222");
nc.publish("orders.new", "{\"id\":1}".getBytes());

Subscription sub = nc.subscribe("orders.*");
Message m = sub.nextMessage(Duration.ofSeconds(1));

// request-reply
Message reply = nc.request("time.now", null, Duration.ofSeconds(1));

// JetStream
JetStream js = nc.jetStream();
js.publish("orders.new", data);
JetStreamSubscription jsub = js.subscribe("orders.new",
    PushSubscribeOptions.builder().durable("dispatch").build());

nc.close();
```
JVM-native, mature, used widely in enterprise. Full JetStream/KV/OS/Micro support.

---

## C# / .NET — `nats.net`

```csharp
// repo: nats-io/nats.net  (NuGet: NATS.Net)
await using var nc = new NatsConnection();
await nc.PublishAsync("orders.new", new byte[]{...});

await foreach (var msg in nc.SubscribeAsync<string>("orders.*"))
    Console.WriteLine(msg.Data);

// request-reply
var reply = await nc.RequestAsync<string,string>("time.now", "");

// JetStream
var js = new NatsJSContext(nc);
await js.CreateStreamAsync(new StreamConfig("ORDERS", new[]{"orders.>"}));
```
Modern async/await .NET client (the current `NATS.Net` v2 line). Full feature parity.

---

## C — `nats.c`

```c
// repo: nats-io/nats.c
natsConnection *nc = NULL;
natsConnection_ConnectTo(&nc, "nats://localhost:4222");
natsConnection_PublishString(nc, "orders.new", "{\"id\":1}");

natsSubscription *sub = NULL;
natsConnection_Subscribe(&sub, nc, "orders.*", onMsg, NULL);

// request-reply
natsMsg *reply = NULL;
natsConnection_RequestString(&reply, nc, "time.now", "", 1000);
```
Native C, **embeddable** in C/C++ apps and other-language FFI bindings. JetStream supported. Good for firmware/embedded contexts.

---

## Ruby — `nats-pure.rb`

```ruby
# repo: nats-io/nats-pure.rb  (gem: nats-pure)
require 'nats/client'
nats = NATS.connect("nats://localhost:4222")
nats.publish("orders.new", '{"id":1}')
nats.subscribe("orders.*") { |msg| puts msg.data }
resp = nats.request("time.now", "", timeout: 1)
js = nats.jetstream
```
Pure-Ruby implementation (no C extension). JetStream supported.

---

## Elixir / Others

- **Elixir** — `gnat` (community, widely used) and others.
- **Zig, Crystal, PHP, Dart, Swift, Haskell, Kotlin, …** — official or strong community clients exist for most languages; check the org and the docs client list.
- **Dart/Flutter** — relevant to OneMind Command; community `dart_nats` / WebSocket bridge patterns (the app currently uses a WebSocket `NatsMonitorService`).

---

## Choosing & Consistency

| Need | Reach for |
|------|-----------|
| Reference/most features first | **Go** |
| Async services, agents | **Python**, **Rust**, **C#** |
| Browser / Node / Deno | **JS/TS** |
| JVM shop | **Java** |
| Embedded/firmware/FFI | **C** |

Because all clients speak the same **wire protocol** ([[NATS/10-PROTOCOL-INTERNALS/wire-protocol|Protocol]]), a Go service, a Python agent, a browser dashboard, and a C firmware node interoperate seamlessly on one bus — polyglot by design.

## Sources
- Developer guide (all langs) — https://docs.nats.io/using-nats/developer
- Clients list — https://nats.io/download/#nats-clients
- GitHub org — https://github.com/nats-io
- NATS by Example — https://natsbyexample.com

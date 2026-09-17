---
title: "Go Client (nats.go)"
section: "08-CLIENTS-SDKS"
type: note
tags: [nats, go, golang, nats.go, client, embedded]
source: https://github.com/nats-io/nats.go
created: 2026-06-14
---

# Go Client (`nats.go`)

The **reference** NATS client — most complete, and the one that can **embed** the server in-process. Module: `github.com/nats-io/nats.go`.

See also: [[NATS/08-CLIENTS-SDKS/README|Clients Index]] · [[NATS/11-ADVANCED/embedded-nats|Embedded NATS]]

---

## Connect

```go
import "github.com/nats-io/nats.go"

nc, err := nats.Connect("nats://localhost:4222",
    nats.Name("my-service"),
    nats.UserCredentials("zeus.creds"),   // JWT creds
    nats.MaxReconnects(-1),
    nats.ReconnectWait(2*time.Second),
    nats.ReconnectJitter(100*time.Millisecond, time.Second),
)
defer nc.Drain()                          // graceful close
```

---

## Core Pub/Sub & Request-Reply

```go
// Publish
nc.Publish("orders.new", []byte(`{"id":1}`))

// Async subscribe
sub, _ := nc.Subscribe("orders.*", func(m *nats.Msg) {
    fmt.Printf("got %s on %s\n", m.Data, m.Subject)
})

// Queue subscribe (load balanced)
nc.QueueSubscribe("orders.new", "workers", handler)

// Request-reply
msg, err := nc.Request("time.now", nil, time.Second)

// Responder
nc.Subscribe("time.now", func(m *nats.Msg) {
    m.Respond([]byte(time.Now().String()))
})
nc.Flush()
```

---

## JetStream (new API: `jetstream` package)

```go
import "github.com/nats-io/nats.go/jetstream"

js, _ := jetstream.New(nc)

// Create/Update a stream
js.CreateStream(ctx, jetstream.StreamConfig{
    Name:     "ORDERS",
    Subjects: []string{"orders.>"},
    Storage:  jetstream.FileStorage,
    Retention: jetstream.LimitsPolicy,
})

// Publish (with dedup id)
js.Publish(ctx, "orders.new", data, jetstream.WithMsgID("order-1"))

// Pull consumer (preferred)
cons, _ := js.CreateOrUpdateConsumer(ctx, "ORDERS", jetstream.ConsumerConfig{
    Durable:   "dispatch",
    AckPolicy: jetstream.AckExplicitPolicy,
    FilterSubject: "orders.new",
})
msgs, _ := cons.Fetch(10)
for m := range msgs.Messages() {
    process(m.Data())
    m.Ack()
}

// Or consume continuously
cons.Consume(func(m jetstream.Msg) { process(m.Data()); m.Ack() })
```

> Two JetStream APIs exist: the **legacy** `nc.JetStream()` (on `nats.go`) and the **newer** `jetstream` subpackage (cleaner, context-aware). Prefer the **new `jetstream` package** for new code.

---

## KV & Object Store

```go
kv, _ := js.CreateKeyValue(ctx, jetstream.KeyValueConfig{Bucket: "CONFIG", History: 5})
kv.Put(ctx, "version", []byte("2.14.2"))
e, _ := kv.Get(ctx, "version")
w, _ := kv.Watch(ctx, "version")
for u := range w.Updates() { /* live changes */ }

obs, _ := js.CreateObjectStore(ctx, jetstream.ObjectStoreConfig{Bucket: "ASSETS"})
obs.PutBytes(ctx, "model.bin", data)
```

---

## Micro Service

```go
import "github.com/nats-io/nats.go/micro"

srv, _ := micro.AddService(nc, micro.Config{Name: "calc", Version: "1.0.0"})
grp := srv.AddGroup("math")
grp.AddEndpoint("add", micro.HandlerFunc(func(r micro.Request) {
    r.Respond([]byte("result"))
}))
```

---

## Notes & Gotchas

- **Always `Drain()`** (not just `Close()`) for graceful shutdown so in-flight messages finish.
- Use **context** with the new JetStream API for timeouts/cancellation.
- The Go client is where **embedded NATS** lives (`server.NewServer(...)`) — run a full broker inside your Go binary; see [[NATS/11-ADVANCED/embedded-nats|Embedded NATS]].
- `nats.go` tracks server features fastest — new JetStream capabilities land here first.

## Sources
- nats.go — https://github.com/nats-io/nats.go
- jetstream package — https://pkg.go.dev/github.com/nats-io/nats.go/jetstream
- micro package — https://pkg.go.dev/github.com/nats-io/nats.go/micro
- NATS by Example (Go) — https://natsbyexample.com

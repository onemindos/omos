---
title: "Rust Client (async-nats)"
section: "08-CLIENTS-SDKS"
type: note
tags: [nats, rust, async-nats, tokio, client]
source: https://github.com/nats-io/nats.rs
created: 2026-06-14
---

# Rust Client (`async-nats`)

High-performance async client built on **Tokio**. Crate: **`async-nats`** (repo `nats-io/nats.rs`). Ideal for performance-critical edge agents and embedded-class services.

See also: [[NATS/08-CLIENTS-SDKS/README|Clients Index]] · [[NATS/02-JETSTREAM/README|JetStream]]

---

## Add & Connect

```toml
# Cargo.toml
[dependencies]
async-nats = "0.35"
tokio = { version = "1", features = ["full"] }
futures = "0.3"
```

```rust
use futures::StreamExt;

#[tokio::main]
async fn main() -> Result<(), async_nats::Error> {
    let client = async_nats::ConnectOptions::new()
        .name("my-service")
        .credentials_file("zeus.creds").await?
        .connect("nats://localhost:4222")
        .await?;
    // ... use client ...
    client.drain().await?;
    Ok(())
}
```

---

## Core Pub/Sub & Request-Reply

```rust
use bytes::Bytes;

// publish
client.publish("orders.new", Bytes::from("{\"id\":1}")).await?;

// subscribe
let mut sub = client.subscribe("orders.*").await?;
while let Some(msg) = sub.next().await {
    println!("{}: {:?}", msg.subject, msg.payload);
}

// queue subscribe
let mut q = client.queue_subscribe("orders.new", "workers".into()).await?;

// request-reply
let resp = client.request("time.now", "".into()).await?;
```

---

## JetStream

```rust
let js = async_nats::jetstream::new(client);

// stream
js.create_stream(async_nats::jetstream::stream::Config {
    name: "ORDERS".into(),
    subjects: vec!["orders.>".into()],
    ..Default::default()
}).await?;

// publish with dedup
js.publish("orders.new", "...".into())
  .await?;   // builder supports headers / msg id

// pull consumer
let consumer = js.get_stream("ORDERS").await?
    .create_consumer(async_nats::jetstream::consumer::pull::Config {
        durable_name: Some("dispatch".into()),
        ..Default::default()
    }).await?;

let mut messages = consumer.messages().await?;
while let Some(Ok(msg)) = messages.next().await {
    // process
    msg.ack().await?;
}
```

KV and Object Store are under `async_nats::jetstream::kv` / `::object_store` with the same semantics.

---

## Why Rust for the Edge

- **Tiny, fast, no GC** — great for resource-constrained edge nodes (robots, sensors) where a Rust agent + `async-nats` sips memory.
- **Tokio async** — handles many concurrent subscriptions efficiently.
- Pairs well with **embedded NATS** (run the Go server in-process elsewhere, connect Rust services to it) or as leaf-node-local clients.

## Sources
- nats.rs — https://github.com/nats-io/nats.rs
- async-nats crate — https://crates.io/crates/async-nats
- docs.rs — https://docs.rs/async-nats
- NATS by Example (Rust) — https://natsbyexample.com

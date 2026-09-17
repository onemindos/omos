---
title: "KV & Object Store — Section Index"
section: "03-KV-OBJECT-STORE"
type: index
tags: [nats, jetstream, key-value, object-store]
created: 2026-06-14
---

# 03 — Key-Value & Object Store

Two higher-level abstractions **built on JetStream streams**. KV is a distributed key-value store with watch/history/TTL/atomic ops; Object Store is chunked blob storage. Both inherit JetStream's replication, persistence, and clustering.

## Notes
- [[NATS/03-KV-OBJECT-STORE/key-value-store|Key-Value Store]] — buckets, get/put/update, watch, history, revisions, TTL, atomic CAS
- [[NATS/03-KV-OBJECT-STORE/object-store|Object Store]] — chunked large blobs, metadata, watch, links

## Mental model
> A **KV bucket** is a JetStream stream named `KV_<bucket>` with subjects `$KV.<bucket>.>`, `MaxMsgsPerSubject` controlling history depth, and `update` implemented via expected-last-sequence (atomic compare-and-set). An **Object Store bucket** is a stream named `OBJ_<bucket>` that chunks large payloads across messages. Everything you know about [[NATS/02-JETSTREAM/streams|streams]] applies underneath.

## See also
- [[NATS/02-JETSTREAM/dedup-and-ordering|Dedup & Ordering]] (expected-last-sequence powers KV atomic updates)
- [[NATS/02-JETSTREAM/retention-and-limits|Retention & Limits]] (history = MaxMsgsPerSubject)

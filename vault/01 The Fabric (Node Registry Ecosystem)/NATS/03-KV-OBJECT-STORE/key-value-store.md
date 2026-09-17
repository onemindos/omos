---
title: "JetStream Key-Value Store"
section: "03-KV-OBJECT-STORE"
type: note
tags: [nats, jetstream, key-value, kv, watch, history, ttl, atomic]
source: https://docs.nats.io/nats-concepts/jetstream/key-value-store
created: 2026-06-14
---

# JetStream Key-Value Store

A distributed, persistent **key-value store** with live **watch**, **history**, **revisions**, **TTL**, and **atomic compare-and-set** — built directly on a JetStream stream.

See also: [[NATS/03-KV-OBJECT-STORE/object-store|Object Store]] · [[NATS/02-JETSTREAM/dedup-and-ordering|Dedup & Ordering]]

---

## What It Is

KV gives you `get` / `put` / `update` / `delete` / `purge` over named **buckets**, plus things a plain KV store usually can't do:
- **Watch** — subscribe to live changes on a key or key pattern (push notifications on write).
- **History** — keep the last N revisions of each key (not just the current value).
- **Revisions** — every write returns a monotonically increasing revision number.
- **TTL** — keys (or the whole bucket) can expire automatically.
- **Atomic create/update** — compare-and-set via expected revision (optimistic concurrency).

### Under the hood
A bucket `CONFIG` is a JetStream stream **`KV_CONFIG`** with subjects **`$KV.CONFIG.>`**. A key `foo` maps to subject `$KV.CONFIG.foo`. `History` = the stream's **`MaxMsgsPerSubject`** (how many revisions per key). `update` uses **`Nats-Expected-Last-Subject-Sequence`** so two writers can't silently clobber each other. This is why everything in [[NATS/02-JETSTREAM/dedup-and-ordering|Dedup & Ordering]] matters here.

---

## CLI

```shell
# Create a bucket with 5 revisions of history and a 1h TTL
nats kv add CONFIG --history=5 --ttl=1h

nats kv put CONFIG nats.version 2.14.2     # write
nats kv get CONFIG nats.version            # read
nats kv history CONFIG nats.version        # all revisions
nats kv watch CONFIG                        # live changes (Ctrl-C to stop)
nats kv del CONFIG nats.version            # delete (tombstone)
nats kv purge CONFIG nats.version          # remove key + all history
nats kv ls CONFIG                           # list keys
```

---

## Client API (Python / Go)

```python
# Python
js = nc.jetstream()
kv = await js.create_key_value(bucket="CONFIG", history=5, ttl=3600)
await kv.put("nats.version", b"2.14.2")
entry = await kv.get("nats.version")
print(entry.value, entry.revision)

# atomic compare-and-set: only succeeds if current revision == entry.revision
await kv.update("nats.version", b"2.14.3", last=entry.revision)

# watch
watcher = await kv.watch("nats.version")
async for e in watcher:
    print("changed:", e.key, e.value, e.revision)
```

```go
// Go
kv, _ := js.CreateKeyValue(&nats.KeyValueConfig{Bucket: "CONFIG", History: 5})
kv.Put("nats.version", []byte("2.14.2"))
e, _ := kv.Get("nats.version")
kv.Update("nats.version", []byte("2.14.3"), e.Revision()) // CAS
w, _ := kv.Watch("nats.version")
for entry := range w.Updates() { /* ... */ }
```

---

## Operations Reference

| Op | Meaning |
|----|---------|
| **Put** | Write/overwrite a value, returns new revision |
| **Get** | Read current value (+ revision); `GetRevision` reads a specific past revision |
| **Create** | Put only if the key doesn't exist (atomic) |
| **Update** | Put only if current revision == expected (atomic CAS) |
| **Delete** | Soft-delete (writes a delete marker; history preserved up to limits) |
| **Purge** | Hard-remove a key and **all** its history |
| **Watch / WatchAll** | Stream live changes for a key/pattern/whole bucket |
| **History** | All retained revisions of a key |
| **Keys** | List all keys in the bucket |
| **Status** | Bucket metadata (values, history, bytes, ttl) |

---

## Bucket Config Knobs

| Config | Effect |
|--------|--------|
| **History** | Revisions kept per key (1–64). `MaxMsgsPerSubject` under the hood |
| **TTL** | Auto-expire entries after this duration |
| **MaxValueSize** | Cap per-value bytes |
| **MaxBytes** | Cap total bucket storage |
| **Storage** | `file` or `memory` |
| **Replicas** | R1/R3/R5 (inherits JetStream HA) |
| **Compression** | Optional S2 compression of stored values |

---

## Watch — the killer feature

`watch` turns KV into a **live configuration/coordination bus**: a service `put`s config, and every watcher is notified instantly (it's a JetStream consumer under the hood). No polling. This is exactly how OneMind's `zeus-context` works — a background **KV Watch** on a separate connection sets `_context_changed=True`, and the next LLM turn injects the new context. Presence, focus-mode, feature flags, and live config all ride KV watch.

---

## Pitfalls (verified)

- **No dots in bucket names** — KV bucket names cannot contain `.` (dots are subject separators). Use `coral-cache`, not `coral.cache`. Same rule as any identity destined for a subject.
- **History = revisions, not a log** — `History=5` keeps the last 5 writes per key. Reading full write history at the stream level: `nats stream view KV_<bucket>` or `nats_stream_read stream="KV_<bucket>"`.
- **Changing the JetStream `domain` orphans buckets** — buckets created under the old domain prefix become invisible. Recreate after any domain change.
- **Atomic update needs the right `last`** — pass the revision you read; a stale revision makes `update` fail (that's the point — it prevents lost updates).

## Sources
- Key-Value Store — https://docs.nats.io/nats-concepts/jetstream/key-value-store
- KV walkthrough — https://docs.nats.io/using-nats/developer/develop_jetstream/kv
- KV CLI — https://docs.nats.io/using-nats/jetstream/kv
- Expected-last-sequence (CAS) — https://docs.nats.io/nats-concepts/jetstream/streams

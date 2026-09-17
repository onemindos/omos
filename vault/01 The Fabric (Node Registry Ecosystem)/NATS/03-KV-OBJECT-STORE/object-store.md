---
title: "JetStream Object Store"
section: "03-KV-OBJECT-STORE"
type: note
tags: [nats, jetstream, object-store, blobs, chunking]
source: https://docs.nats.io/nats-concepts/jetstream/obj_store
created: 2026-06-14
---

# JetStream Object Store

Store and distribute **large binary objects** (files, firmware, models, media) over NATS, with metadata, watch, and chunked transfer — built on a JetStream stream.

See also: [[NATS/03-KV-OBJECT-STORE/key-value-store|Key-Value Store]] · [[NATS/02-JETSTREAM/streams|Streams]]

---

## What It Is

The Object Store is the "S3-lite over NATS" abstraction. Where KV holds small values, Object Store holds **arbitrarily large blobs** by **chunking** them across many stream messages and reassembling on read. Each object has a name, metadata (size, digest, headers, description), and a content digest (SHA-256) for integrity.

### Under the hood
An object bucket `ASSETS` is a JetStream stream **`OBJ_ASSETS`** with subjects:
- `$O.ASSETS.C.>` — **chunks** (the object payload, split into pieces)
- `$O.ASSETS.M.>` — **meta** (object info: name, size, digest, headers)

Reading streams chunks back in order, verifying the digest. Because it's a stream, you get replication (R1/R3/R5), persistence, and watch for free.

---

## CLI

```shell
nats object add ASSETS                       # create bucket
nats object put ASSETS ./firmware-v2.bin     # upload (chunked)
nats object get ASSETS firmware-v2.bin       # download
nats object ls ASSETS                          # list objects
nats object info ASSETS firmware-v2.bin       # metadata + digest
nats object del ASSETS firmware-v2.bin        # delete
nats object watch ASSETS                       # live changes
```

---

## Client API

```python
# Python
js = nc.jetstream()
obs = await js.create_object_store(bucket="ASSETS")
# put from bytes or a file
await obs.put("model.gguf", open("model.gguf","rb").read())
info = await obs.get_info("model.gguf")
print(info.size, info.digest)
res = await obs.get("model.gguf")            # res.data is the reassembled blob
```

```go
// Go
obs, _ := js.CreateObjectStore(&nats.ObjectStoreConfig{Bucket: "ASSETS"})
obs.PutBytes("model.gguf", data)
b, _ := obs.GetBytes("model.gguf")
```

---

## Operations

| Op | Meaning |
|----|---------|
| **Put** | Upload an object (auto-chunked); returns ObjectInfo |
| **Get** | Download + reassemble; verifies digest |
| **GetInfo** | Metadata only (size, digest, headers, mtime) |
| **Delete** | Remove an object (chunks + meta) |
| **List** | Enumerate objects in the bucket |
| **Watch** | Live notifications on object add/update/delete |
| **AddLink / AddBucketLink** | A named link pointing to another object or bucket |
| **Seal** | Make a bucket read-only (no further writes) |

---

## Config Knobs

| Config | Effect |
|--------|--------|
| **Storage** | `file` or `memory` |
| **MaxBytes** | Total bucket size cap |
| **TTL** | Auto-expire objects |
| **Replicas** | R1/R3/R5 |
| **Compression** | Optional S2 compression |
| **Metadata** | Arbitrary key-value metadata on the bucket |

---

## Chunking & Size

Objects are split into chunks (default chunk size in the hundreds of KB, tunable). This is why Object Store handles blobs far larger than the server's `max_payload` (~1MB default) — no single message holds the whole object. Practical for firmware, ML model files, images, documents, backups distributed across the fabric. For *truly* massive archival (many GB–TB each, infinite retention), a dedicated object store (S3/MinIO) is still the better tool — Object Store shines for **distributing moderate blobs over the same bus** you already run.

---

## Use Cases

- **Firmware / OTA updates** — push a binary to an Object Store bucket; edge nodes `watch` and pull on change.
- **ML model distribution** — ship `.gguf`/weights to inference nodes over the bus.
- **Media / document handoff** — agents exchange images, PDFs, recordings.
- **Config bundles** — larger-than-KV config artifacts with integrity checking.

---

## OneMind Relevance

For the fleet (firmware to barn-pi, models to inference nodes, media between agents), Object Store lets the **existing NATS bus** double as the artifact-distribution channel — no separate object server at homestead scale. Pair with **watch** for pull-on-change OTA, and with **leaf-node domains** so an edge node mirrors only the objects it needs.

## Sources
- Object Store — https://docs.nats.io/nats-concepts/jetstream/obj_store
- Object Store walkthrough — https://docs.nats.io/using-nats/developer/develop_jetstream/object
- Object Store CLI — https://docs.nats.io/using-nats/jetstream/obj_store

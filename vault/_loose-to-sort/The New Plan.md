# The New Plan

## A) The Repo Split Timing

P1-A (the NATS foundation) just completed in command for bootstrap convenience — but per this decision, NATS config and canon docs belong in **core**. Recommendation: create `onemindos-core` now, relocate the canon docs + P1-A's NATS output there, and keep `command` as pure UI. Clean break while the codebase is still small.

## B) The `omd` Daemon Language

Still open from §8: **Go** (single static binary, cross-compiles clean to the barn Pi and ARM robots, matches the NATS/NEX/Wombat ecosystem) with a Python lane for agent glue — that's the recommendation.

And separately: **P1-A is sitting in review-blocked state** — the worker reports all 7 verifications pass and pushed its branch. Should review that diff against the doctrine before building on it.

Sequence options:
1. Create `onemindos-core` and do the split first
2. Review/merge P1-A first so the foundation is locked
3. Both — review P1-A, then immediately relocate it into the new core repo

---

## `omd` is NOT the Leaf Node

They're different *kinds* of thing:

- **Leaf node** = a NATS **server** running in leaf mode. It's infrastructure — the local switchboard that bridges this machine's subjects up to core and holds local JetStream (store-and-forward when the link drops). Still the "telephone network," just physically on that machine.
- **omd** = a NATS **client**. It connects *to* a server and has hands on the machine.

So on any given node you'd naively have **two processes**: the leaf server (bus) + `omd` (worker). That's the thing your instinct is rebelling against — and correctly.

---

## Your Instinct: "Do it all in one universe" → Yes, and Go makes it possible

Because `omd` is in Go, it can **embed the NATS leaf server inside itself**. The NATS server is a Go library (`nats-server/v2/server`) — you instantiate it in-process. So `omd` becomes **one binary, one process, one universe** per node:

```
omd (single Go binary on a node)
├── embedded NATS leaf server → connects up to core, holds local JetStream
└── worker logic → registers node actor, capabilities, runs commands
```

That's exactly your "same universe." One download, one process, the node *is* one thing. This is a proven pattern — it's how Synadia's NEX runs.

---

## This is Another Reason Go Was Right

**You can only embed the NATS server if you're in Go.** It's a Go library. Python literally cannot do this — it would *force* you into two separate processes forever. So "do it all in one universe" isn't just compatible with the Go decision — it's an argument *for* it. Had we picked Python, your instinct here would be impossible to satisfy.

---

## The One Exception: Core/Hub Stays Standalone

Don't embed the **core** into a worker daemon. The core holds the durable streams (EVT/TEL/CMD), the system account, the permanent record — it's bedrock. It must outlive any worker and be independently upgradeable.

- **Core/hub node** → standalone nats-server (shared infrastructure). *This is what your Mac is running right now.*
- **Edge nodes** (barn Pi, office box, robots) → `omd` with an **embedded leaf**. One universe each.

So on *this* machine today, `omd` is correctly just a plain client connecting to the standalone server built in P1-A. When the barn Pi gets added, *that* node's `omd` embeds a leaf and dials home.

---

## The Honest Trade-off (Why Not Embed Everywhere)

**Lifecycle coupling.** If the bus lives inside `omd`, then every time you restart or upgrade the worker logic, the local bus bounces too — dropping consumers and buffered messages.

- On **edge nodes** that rarely change, simplicity wins → embed. One universe.
- On the **core and active dev machine**, you want the bus rock-stable while iterating on worker code → keep them separate.

**The rule: embed where simplicity matters, separate where stability matters.**

---

## What This Means for the Running Card

OMD-1 is still correct as-is. Even an embedded-leaf `omd` needs all the client-side skeleton it's building (register actor, advertise, handle ping). Embedding the leaf is a *later additive* card for edge nodes — it doesn't change the skeleton. No need to touch the worker.

**Resolution:** `omd` embeds a leaf on edge nodes (one universe), core stays standalone. Record in `ARCHITECTURE.md §8` decision log — add an "OMD-2: embed leaf node" card to the backlog for when the first edge node comes up.

Obsessed

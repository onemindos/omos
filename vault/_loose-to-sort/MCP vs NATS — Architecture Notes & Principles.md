# MCP vs NATS — Architecture Notes & Principles

## MCP vs NATS: The Case for Direct

Yes — 100% agree. The MCP adds a whole Node.js process + stdio overhead for what is just HTTP calls to HA's REST API. A direct Python tool like the NATS one hits `http://100.64.25.12:8123/api/` directly over Tailscale, zero intermediary.

---

## Synadia Agent Protocol for NATS

Previously called "NATS Agent Protocol" — renamed to "Synadia Agent Protocol for NATS."

The trajectory: Synadia becomes the AWS of agent infrastructure. Not the models, not the reasoning — the **discovery, routing, identity, delivery guarantees, and observability layer** that sits underneath any heterogeneous agent fleet. NATS already powers 450M+ downloads of enterprise messaging. They're extending that to agents.

---

## The Unification Layer — One Mind, Literally

The thing you're reaching for is the unification layer: every entry in 05 and 08 becomes an addressable node on the bus, with one identity, one subject, live state, and an event stream. The bus is the membrane that makes the digital catalog and the physical catalog the same fabric.

That's One Mind — not metaphorically, literally: **one addressable namespace for your entire existence.**

---

## Hard Rules & Principles

### 1. Robots and Drones — Bridge the Seam, Do NOT Replace Their Language

This is a hard real-time engineering line. A drone's inner control loop (PX4/ArduPilot/MAVLink at 100–400Hz, ROS2 over DDS) cannot tolerate general-bus latency/jitter.

**The native pattern:**
- The robot runs its own local NATS leaf
- Inner loop stays onboard/native
- Only telemetry + high-level commands cross to the hub
- Bridge MAVLink↔NATS and ROS2↔NATS at the command/telemetry layer
- **Never put flight control on the shared bus**

### 2. Build the Spine First — Sprawl Without a Spine is the Risk

"A lot of the services we have is very disorganized even though up and running." Adding more nodes onto a foundation with no identity, no canonical schema, no live registry = more chaos, faster.

The entire bus is one `$G` account with one shared password right now.

**Before a single new node:** identity + subject schema + registry-as-KV must land first. This is the one non-negotiable. Build the spine, then hang everything off it.

### 3. "Wrap Everything" Silently Becomes "Maintain Everything"

Every bespoke bridge is code you own forever. The native discipline (one wrapper pattern, micro-services, import/export) is the only thing standing between you and a 200-daemon swamp. The temptation to hand-roll "just one quick bridge" per service is the trap that kills this.

### 4. Not Everything Deserves to Be on the Bus

Honest filter for each node: **does it stream/react in real time, or do you just call it occasionally?**

- **Bus-native** — things that emit and react continuously
- **On-demand wrap** — things you touch monthly

Putting a once-a-month API on the bus is pure overhead. "Anything could be a node" is true. "Everything should be a streaming node" is not.

### 5. One MBP is a Single Point of Failure

The moment the bus is load-bearing for physical things (locks, security, presence), the MBP going down = your digital life goes dark. Leaf nodes give edge survival (local buffering), but the hub still matters.

**Before this carries real-world weight:** the hub needs redundancy — cluster, or hot standby on the spare machine. That idle box is a perfect use for this.

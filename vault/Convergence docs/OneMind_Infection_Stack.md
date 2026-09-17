# The "OneMind Infection" Stack — Device Onboarding Architecture

**Core Principle:** Every device that joins the mesh becomes an organ, not a router. It publishes what it can sense/do, subscribes to commands it can execute. The bus is the nervous system; every device is a nerve ending with a fixed contract.

---

## The 4-Layer Mesh Kit (Every Device)

Every device that joins OneMind receives the same 4-layer infrastructure:

| Layer | Tool | Role |
|-------|------|------|
| **Network** | Tailscale | Private mesh — device joins tailnet, no public exposure |
| **Bus** | NATS leaf node | Nervous system — connects to hub, carries subjects/events |
| **Compute** | KubeSolo or NEX | Workload execution |
| **Management** | Portainer Edge Agent | Central pane — Portainer already manages everything |

### Critical Piece: KubeSolo + Portainer Integration

KubeSolo is a Portainer project. The Edge Agent connects your KubeSolo fleet into one Portainer operator plane:
- GitOps deployments
- RBAC
- Remote lifecycle across every node
- Fits your existing Portainer doctrine exactly

---

## KubeSolo vs NEX — How to Decide

They're not competing. They answer different questions:

### NEX (NATS Execution Engine)

**Best for:** Dedicated edge processors, robots, sensors

**Characteristics:**
- Device runs NATS-native workloads only
- Lightweight functions triggered by bus events
- Zero container orchestration overhead
- True sovereignty moat
- Minimal resource footprint

### KubeSolo

**Best for:** Raspberry Pis, field laptops, branch sites, any device that previously ran Docker Compose

**Characteristics:**
- Device needs full containerized services with K8s semantics
- Helm support
- Operators
- Health probes
- RBAC
- <200MB RAM
- ARM64 native
- Full Kubernetes API

### Coexistence

Both can coexist on the same device:
- NEX handles bus-native execution
- KubeSolo manages long-running containerized services
- NATS leaf node underlies both

---

## The Infection Playbook (IaC)

Device joins OneMind via 4-layer bootstrap:

```
Device joins OneMind:
  1. curl | sh → Tailscale
     (joins tailnet)
     
  2. nats-server leaf config → joins bus
     (ONEMIND account, leaf_tak creds)
     
  3. curl | sh → KubeSolo   (if capacity)
     OR
     nex agent (if ultra-light)
     
  4. Portainer Edge Agent
     (registers with your existing Portainer pane)
```

### IaC Expression

Fully GitOps-expressible:
- 4 shell scripts + 2 manifests
- Committed to omos-infra
- Every device becomes identically structured
- Centrally managed
- Bus-connected
- Privately networked

---

## Verdict: The Right Pattern

**Yes.** NATS + leaf + Tailscale + KubeSolo/NEX + Portainer Edge is the IaC mesh architecture.

**The unlock:** Portainer-KubeSolo native integration means you don't manage individual KubeSolo nodes directly. Portainer does it at fleet scale. Single management pane doctrine holds.

---

## The Core Reframe: Organs, Not Routers

**You already have this understanding.**

A device joins OneMind as an organ, not a router.

**What it doesn't do:**
- Move subjects for others
- Route traffic

**What it does:**
- Serves subjects
- Publishes what it can sense/do
- Subscribes to commands it can execute

**The bus is the nervous system.**
**Every device is a nerve ending with a fixed contract.**

---

## The Lingua Franca — Three Universal Things

Identical on every device in the world:

### 1. NKey/JWT Identity

Who am I? What subjects am I allowed to touch?

### 2. NATS Subjects

The address space (the same everywhere)

### 3. Envelope Schema

The message shape (the same everywhere)

### Language Agnostic

Any language with a NATS client can speak this:
- Go robot
- Swift iPhone
- Rust vehicle ECU
- Python sensor

**Same subjects, same envelope. That's why NATS is the universal layer, not Flutter.**

```
Flutter = the face (human UI)
NATS = the nerve (device capability bridge)
```

---

## The Layered Pattern — What Every Device Carries

Bottom-up. The spine (L0–L3) is universal and identical everywhere. Everything above varies by device class.

### Layer Definition Table

| Layer | What | Robot | Vehicle | iPhone | Pi/Gateway | Sensor |
|-------|------|-------|---------|--------|-----------|--------|
| **L0** | Identity | NKey/JWT, scoped creds, subject namespace | ✓ | ✓ | ✓ | ✓ | ✓ |
| **L1** | Network | Tailscale (reach across NAT/internet) | ✓ | ✓ | ✓ | ✓ | ~ |
| **L2** | Bus | NATS connection | leaf | leaf | client | leaf | client |
| **L3** | Contract | Envelope + subject taxonomy | ✓ | ✓ | ✓ | ✓ | ✓ |
| **L4** | Capability | Native bridge (sensors/actuators → bus) | ROS2 | CAN/ECU | Swift | GPIO | firmware |
| **L5** | Compute | NEX / KubeSolo (workloads) | both | both | — | KubeSolo | — |
| **L6** | Management | Portainer Edge / control plane | ✓ | ✓ | via bus | ✓ | via bus |
| **L7** | Console | Flutter (human face) | opt | opt | ✓ | — | — |

### The Spine (L0–L3)

```
L0–L3 = the spine

If a device has these four, it's already OneMind:
  ├─ Identity (knows who it is)
  ├─ Reachability (can be reached across networks)
  ├─ Bus connection (speaks NATS)
  └─ Language (common envelope schema)

Everything above is just how much it can do.
```

---

## The Organ Contract — What Every Node Implements

Regardless of tier, every organ implements the same five methods.

**This is the reusable code/SDK pattern.**

### Five-Method Contract

Every device publishes/subscribes on these patterns:

#### 1. Announce — Presence Heartbeat + Last-Will

**Purpose:** Bus knows when device dies

**Subject:** `device.<class>.<id>.status`

**Payload:** Heartbeat + Last-Will (JetStream-backed)

**Frequency:** Every N seconds (configurable)

#### 2. Self-Describe — Capability Manifest

**Purpose:** Mesh discovers what device can do

**Subject:** `device.<class>.<id>.capabilities`

**Payload:** JSON manifest
```json
{
  "class": "robot",
  "id": "r001",
  "capabilities": [
    { "action": "move_forward", "args": ["speed"] },
    { "sensor": "gps", "rate": 1 },
    { "sensor": "camera", "rate": 30 }
  ],
  "compute": "NEX+KubeSolo",
  "network": "Tailscale"
}
```

#### 3. Serve Telemetry — Stream What It Senses

**Purpose:** Continuous sensor data

**Subject:** `device.<class>.<id>.<sensor>`

**Payload:** Sensor reading (time-series)

**Examples:**
- `device.robot.r001.gps` → latitude, longitude, accuracy
- `device.vehicle.v001.can_obd2` → RPM, fuel, temperature
- `device.sensor.s001.temperature` → degrees + timestamp

#### 4. Accept Commands — Request/Reply

**Purpose:** Device executes commands

**Subject:** `device.<class>.<id>.cmd.<action>`

**Pattern:** Request/reply via NATS (synchronous)

**Examples:**
- `device.robot.r001.cmd.move_forward` → request: {"speed": 0.5} → reply: {"status": "ok"}
- `device.vehicle.v001.cmd.unlock` → request: {} → reply: {"status": "unlocked"}
- `device.sensor.s001.cmd.calibrate` → request: {} → reply: {"status": "calibration_ok"}

#### 5. Survive Disconnection — JetStream Store-and-Forward

**Purpose:** Sovereignty — works offline, cloud-disposable

**Mechanism:**
- Telemetry published to JetStream (persisted on edge hub)
- Device buffers locally when bus unavailable
- Resyncs when reconnected
- No data loss (eventual consistency)

**Result:** Device is autonomous. Cloud is optional.

---

## The OneMind Node SDK

**Write the five-method contract once per language.**

### Language Targets

- [ ] Go (robotics, edge processors)
- [ ] Swift (iOS organ process)
- [ ] Kotlin (Android organ process)
- [ ] Rust (vehicle ECU, high-safety)
- [ ] Python (sensors, Raspberry Pi)
- [ ] TypeScript (browser/desktop)

### What Each SDK Provides

```
OneMindNode (language-specific):
  ├─ announce()           → heartbeat + LW
  ├─ capabilities()       → self-describe
  ├─ telemetry(channel)   → stream data
  ├─ command(action)      → request/reply handler
  ├─ offline_buffer()     → JetStream resilience
  ├─ identity_from_jwt()  → NKey/JWT validation
  └─ connect(nats_url)    → leaf connection
```

Every device that imports OneMindNode becomes an organ.

---

## Device Tiers (Which Layers You Carry)

### Tier 1 — Capability Node (iPhone, Sensor)

**Examples:** iPhone, smart sensor, wearable

**Carries:**
- L0–L3 (spine)
- L4 (native bridge)
- Native NATS client only

**Compute:** None (native platform)

**Management:** Via bus commands

**Characteristics:**
- Embedded NATS client
- Native bridge to hardware
- No orchestration
- Keep-alive is native platform's problem
  - iOS: location entitlement
  - Android: foreground service
  - Sensor: firmware timer

---

### Tier 2 — Workload Node (Pi, Dell, Gateway)

**Examples:** Raspberry Pi, Dell OMOS, field gateway

**Carries:**
- L0–L6 (full spine + compute + management)
- NATS leaf node
- NEX and/or KubeSolo

**Compute:** Containerized services (Helm, operators)

**Management:** Portainer Edge Agent + central pane

**Characteristics:**
- Runs long-lived services
- Can orchestrate other workloads
- Fully managed by Portainer
- Offline-capable (local JetStream)
- GitOps deployable

---

### Tier 3 — Autonomous Edge (Robot, Vehicle)

**Examples:** Robot, autonomous vehicle, tactical system

**Carries:**
- L0–L6 (full spine)
- L5 (both NEX + KubeSolo)
- NATS leaf node + NEX + KubeSolo

**Compute:** Bus-native (NEX) + containerized (KubeSolo)

**Management:** Portainer Edge + local decision loop

**Characteristics:**
- Operates fully offline
- Local decision-making (no cloud required)
- Syncs when bus reconnects
- Full sovereignty
- Can act as edge hub for other devices
- Highest autonomy tier

---

## On Flutter — The Adjustment You're Sensing

**Decouple face from nerve.**

### Current Assumption (Incorrect)

Flutter must be the always-on "organ" that serves capabilities to the mesh.

### The Problem

iOS background rules prevent Flutter from being your always-on node:
- Can't hold location-update entitlements in background
- Can't maintain consistent network presence
- Dies when app backgrounded

Same on robots: ROS2 cannot run as Flutter widget. Can't speak native CAN bus. Can't hold real-time guarantees.

### The Correction

**Flutter = the console (L7).** Cross-platform human UI.
- Great for atoc (tactical console)
- Great for dashboards
- Great for family views
- Keep it

**Native = the organ (L4).** The always-on capability bridge.
- iOS: Swift holding location-updates entitlement
- Android: Kotlin with foreground service
- Robot: ROS2 native process
- Vehicle: CAN bus native driver
- Stays alive even when user isn't looking

### The Pattern: Layered App

One app can contain both:

```
App Architecture:
├─ Native layer (always-on organ)
│   ├─ Swift process (L4: location, sensors)
│   ├─ NATS SDK integration
│   └─ JetStream buffering
│
└─ Flutter layer (human UI console)
    ├─ Cross-platform UI
    ├─ Subscribes to telemetry from native layer
    ├─ Sends commands via native layer
    └─ Displays mesh status
```

### Don't Force Flutter to Be the Nerve

Flutter is the face. It's where humans look. It's not where the device listens (that's Swift/Kotlin/native).

---

## The Doctrine-Level Question

**What's the first device type you want to target?**

Options:
- Raspberry Pi (workload node, gateway)
- The OMOS Dell (Tier 2/3 hybrid)
- Something else (specify)

This is the doctrine-level decision — the one that makes OneMind actually one mind.

---

## Verdict: The Pattern and the Route

**Yes — this is the pattern.**

**The components:**
1. ✅ NATS-SDK-per-language is the lingua franca
2. ✅ Five-method Organ Contract is the reusable spine
3. ✅ Tailscale is reachability
4. ✅ NEX/KubeSolo is compute-when-needed
5. ✅ Portainer is the central management pane

**The outcome:**
Every device in the world that imports the OneMind Node SDK becomes an organ.

**That's literally one mind.**

---

## Recommendation: Codify as Doctrine

This is foundational enough that it should be version-controlled in omos-infra.

### Suggested Artifact

**File:** `reference/device-node-contract.md`

**Content:**
1. Layered model (L0–L7)
2. Organ Contract (five methods)
3. Tier definitions (Tier 1–3)
4. Per-language SDK plan
5. Bootstrap playbook (IaC)
6. Flutter decoupling pattern
7. Device class taxonomy

**Why:**
This becomes the canonical reference every team member builds devices against. It's the pattern language. It's the contract. It's how every device becomes an organ, not a one-off implementation.

---

## Summary: The OneMind Infection

### How It Works

```
Device physical connection
    ↓
Tailscale joins (network layer)
    ↓
NATS leaf connects (bus layer)
    ↓
KubeSolo or NEX bootstraps (compute layer)
    ↓
Portainer Edge registers (management layer)
    ↓
Device runs OneMind Node SDK (organ layer)
    ↓
Device announces itself + publishes capabilities
    ↓
Bus discovers device
    ↓
Device serves telemetry + accepts commands
    ↓
Portainer manages workloads centrally
    ↓
Device is now an organ in OneMind
    ↓
Fully offline-capable, bus-aware, centrally managed
```

### The Infection Is Complete

Every device is:
- ✅ Identifiable (L0)
- ✅ Reachable (L1)
- ✅ On the bus (L2)
- ✅ Speaking the language (L3)
- ✅ Bridging sensors/actuators (L4)
- ✅ Running workloads (L5)
- ✅ Managed centrally (L6)
- ✅ Visible to humans (L7)

---

**Captured:** 2026-07-11

**Status:** Foundational architecture for OneMind device onboarding

**Relevance:** Complete pattern for converting physical devices into mesh organs

**Key Innovation:** Five-method Organ Contract as the universal reusable pattern

**Doctrine:** NATS-SDK per language, Tailscale for networking, Portainer for management, native code for always-on sensing, Flutter for human UI

**First Priority:** Choose device tier (Tier 1/2/3) to target first, then codify device-node-contract.md in omos-infra

**Philosophy:** Every device becomes an organ, not a router. The nervous system (NATS) reaches every nerve ending (device). That's how one mind works.

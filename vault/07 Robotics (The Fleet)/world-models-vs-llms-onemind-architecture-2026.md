# World Models vs. LLMs — Deep Brief for OneMind Omnipresence Architecture

**Filed:** 2026-06-07 · **For:** Zeus · **By:** Legacy
**Trigger:** Introduced to world models (WMs); friend argues they're "the future," do things LLMs can't, run locally, and give "real understanding of the world." Wants an always-on (24/7) agent to watch his back, learn, predict, and act across robotics/drones/wearables/infra.
**Primary sources:** LeWorldModel paper (Maes, Le Lidec, Scieur, **LeCun**, Balestriero — arXiv 2603.19312v3, Jun 2026); AbdelStark/worldforge repo; Meta V-JEPA 2 release (Jun 2025); HF V-JEPA 2 / V-JEPA 2-AC docs.

---

## TL;DR (the honest version)

1. **Your friend is directionally right but technically imprecise.** World models *are* a real, distinct, and important paradigm — and they genuinely do things LLMs can't (predict physical consequences, plan in latent space, run cheaply/locally, detect "surprise"). But a world model is **NOT a replacement for an LLM and CANNOT be the "higher meta-agent" that manages other agents.** That's a category error.

2. **A world model is a *cost oracle / physics simulator*, not a reasoner or an orchestrator.** It answers *"if I take action A, what does the world look like next, and is that closer to my goal?"* — extremely fast, in latent space, no language. It does **not** do symbolic reasoning, tool-calling, multi-step planning over APIs, or natural-language communication. Those remain LLM jobs.

3. **The correct mental model is a STACK, not a winner.** WMs and LLMs are complementary layers. The frontier (and what worldforge literally implements) is: **LLM/policy proposes actions → world model scores/rolls them out → controller executes → repeat.** Neither replaces the other.

4. **For OneMind specifically:** The world model belongs at the **bottom (physical/robotics) layer**, not the top. Your "omnipresent meta-agent" stays an LLM-class system (GPT-realtime + Hermes). World models get bolted UNDER it, per-robot/per-sensor, to give the physical layer fast predictive control and "always-on perception that doesn't burn tokens."

5. **The "always-on without burning tokens" instinct is your single best insight here.** That's a legitimately correct reason to want a WM (or a small always-on perception model) in the loop. More below.

---

## Part 1 — What a world model actually IS (precisely)

A world model learns to **predict the consequences of actions** in a compact **latent space** (not pixels, not words). Canonical structure (confirmed in both LeWM and V-JEPA 2):

| Component | Job |
|---|---|
| **Encoder** | Compresses raw observation (camera frame, sensor stream) → low-dim latent state `z` |
| **Predictor (dynamics)** | Given `z_t` + action `a_t` → predicts `z_t+1` ("imagine the next state") |
| **No decoder needed** | JEPA-class WMs predict in *latent* space, never reconstruct pixels — that's why they're fast |

What this buys you (all evidenced in the papers):

- **Imagination/planning:** roll out hundreds of hypothetical action sequences in <1s and pick the best — *before* moving a real motor. LeWM plans full trajectories in **0.98s vs 47s** for a foundation-model WM (DINO-WM) — **~48× faster**.
- **Physical understanding:** LeWM's latent space linearly encodes real physical quantities (object location, angle) — provable by probing.
- **Surprise / anomaly detection:** WMs spike on physically impossible events (object teleports) — a built-in "that's wrong" signal. **This is the "watch my back" primitive.**
- **Sample efficiency:** V-JEPA 2 learns robot control from only **62 hours** of robot data (vs thousands typically), because it pre-trained on 1M+ hours of *passive* video first.
- **Local / cheap:** LeWM is **15M params, trains on a single GPU in hours.** This is the basis of your friend's "runs locally easier" claim — and it's TRUE for this class.

### What a world model is NOT (this is the crux)
- ❌ It does not reason in language or chains of thought.
- ❌ It does not call tools/APIs or orchestrate software.
- ❌ It does not communicate with you.
- ❌ It does not "manage other agents." It has no concept of an agent, a task queue, or a tool.
- ❌ It is **task/domain-bound** to what it was trained to perceive and act within (a robot arm WM knows arm+table physics, not your calendar).

> A world model is to a robot what the cerebellum + visual cortex are to you: fast, subconscious, predictive motor/perceptual control. **It is not the prefrontal cortex.** Your LLM layer is the prefrontal cortex.

---

## Part 2 — The model family your friend named (decoded)

| Name he said | What it actually is | Relevance to you |
|---|---|---|
| **I-JEPA** | Image JEPA (Meta, 2023). Self-supervised *image* representation learning by predicting masked patches in latent space. The original JEPA proof. | Foundational, not directly deployable for control. |
| **V-JEPA** | Video JEPA (Meta, 2024). Same idea on video. Learns physical dynamics from passive video. | Foundation for V-JEPA 2. |
| **V-JEPA 2** | **The big one.** Meta, Jun 2025. 1.2B-param video world model. SOTA physical understanding + prediction + **zero-shot robot planning.** Open weights, commercial-OK. | ⭐ **Most directly usable for your robotics layer.** |
| **V-JEPA 2-AC** | The **action-conditioned** post-trained checkpoint of V-JEPA 2 (the "-AC"). Trained on ~62h robot trajectories so the predictor takes *actions* as input. Deploys zero-shot on a Franka arm via goal-image planning. | ⭐ This is the *control-ready* version. "v-jepa2-ac" = what you actually wire to a robot. |
| **LeWM (LeWorldModel)** | The PDF. Maes et al. **w/ Yann LeCun**. First JEPA that trains *stably end-to-end from pixels* with one hyperparameter, 15M params, single GPU, ~48× faster planning. | ⭐ The "runs locally easy" poster child. Research-grade but the efficiency direction you want. |
| **"lewm"** | = LeWM above (same thing, the repo short-name `le-wm`). | Same. |
| **"nexarbadas 2.0"** | Almost certainly a voice-transcription mash-up. Two real things it could mean: **(a) Nexar** (dashcam/driving data + driving world models) and **(b) BADAS** — "Context-Aware Collision Prediction Using Real-World Dashcam Data" (arXiv 2510.14876, Oct 2025), a Nexar-affiliated driving-safety model. "2.0" likely a version. | Niche (driving/collision). Park it unless we go deep on the autonomous-vehicle node. **Worth re-asking him to confirm the exact name.** |

**Adjacent heavy hitters worth knowing (not from him, but you should):**
- **DreamerV4 / Dreamer line (Hafner, DeepMind)** — train agents *inside* the world model ("in imagination"). The RL flagship.
- **NVIDIA Cosmos** — "world foundation models" explicitly for robotics/AV; pairs with their policy models. (worldforge has a `cosmos-policy` provider.)
- **Genie 3 (DeepMind, 2025)** — real-time *generative* interactive worlds from text. Different sub-branch (generative WM, not JEPA predictive-control WM).
- **GR00T (NVIDIA)** — humanoid robot policy foundation model. (worldforge has a `gr00t` provider.) Relevant to your **G1** plan.
- **MuZero** — the ancestor: learned a model "sufficient for planning," mastered Go/chess/Atari without rules.

---

## Part 3 — WorldForge (the repo he sent) — what it tells us

WorldForge (AbdelStark) is **not a model** — it's a **harness** that composes world-model workflows. It is the single most useful artifact he sent, because it shows *the actual architecture pattern the whole field has converged on.* The backbone loop:

```
observe state
  → policy proposes candidate actions        (policy provider: GR00T, LeRobot, Cosmos-Policy)
  → world model scores / rolls out futures   (score/predict provider: LeWM, V-JEPA)
  → latent MPC controller picks lowest-cost   (CEM / receding-horizon optimizer)
  → execute through a provider               (robot)
  → persist, evaluate, observe again
```

Five capabilities: `predict`, `score`, `policy`, `embed`, `plan`. Critical insight baked into its design:

> **"LeWorldModel is a score provider, not a video generator. Cosmos-Policy, GR00T, and LeRobot are policy providers, not predictive world models."**

The framework *forcibly separates* "thing that proposes actions" from "thing that predicts/scores consequences." **That separation is the lesson.** Your LLM/agent layer is a fancy policy proposer + goal-setter. The world model is the consequence scorer. They are different boxes on purpose.

---

## Part 4 — Grading your friend's specific claims

| Claim | Verdict | Nuance |
|---|---|---|
| "World models are the future" | ✅ **True** (for physical/embodied AI) | They're *a* future, the embodied/robotics one. Not a replacement for language agents. |
| "There are things LLMs cannot do" | ✅ **True** | Fast physical prediction, latent planning, sub-second motor control, token-free always-on perception, calibrated physical "surprise." |
| "Have to use text-based LLMs to do more real-world things" | ⚠️ **Half-true / muddled** | LLMs are *reactive symbolic* layers. For *physical* real-world action, a WM is more direct. But for *digital* real-world action (tools, comms, orchestration) the LLM is exactly what you need. |
| "Run world models locally easier" | ✅ **True for the efficient class** (LeWM 15M, single GPU) | NOT true for all (V-JEPA 2 is 1.2B; generative WMs like Genie are huge). The *predictive-control* WMs are genuinely light. |
| "More interference [inference] for real understanding of the world" | ✅ **True, with an asterisk** | They encode *physical* structure better. They understand *physics*, not *meaning/intent/language*. "Understanding" is domain-specific. |
| WM as "higher meta-agent managing other agents" | ❌ **False — category error** | A WM has no concept of agents, tasks, or tools. The meta-agent/orchestrator MUST be the LLM-class system. WMs sit UNDER it. |

**Bottom line on the friend:** He's a credible signal — he's pointing you at exactly the right *new layer* (and LeCun himself is behind this thesis). But he's overshooting on "WM replaces/sits-above the LLM stack." The reality is **WM is a new bottom layer, not a new top.**

---

## Part 5 — The "always-on / don't burn tokens" question (your best instinct)

You said the sharpest thing in the whole message: *"a robot is on 24/7 reasoning and perceiving... without having to worry about it being activated based on tokens."* This is the **real, correct reason** to bring world models in. Here's the precise truth:

- **LLM agents are fundamentally reactive and token-metered.** Every "thought" costs tokens and latency. You cannot run GPT-realtime as a literal 24/7 continuous perceptual loop economically — it's the wrong tool for "always watching."
- **A world model (or small perception model) IS the right substrate for continuous perception.** It runs locally, cheaply, forever, on a sensor/camera/robot, encoding the world into latents every frame and flagging *surprise* (anomalies) without any token cost.
- **The pattern:** WM/perception layer watches 24/7 → only when it detects something salient (surprise spike, goal-relevant event) does it **wake the LLM layer** to reason, decide, communicate, and call tools. This is exactly **interrupt-driven architecture**: cheap always-on sensor, expensive reasoning on-demand.

> This maps PERFECTLY onto your existing autonomy ladder (L0–L3) and your Hermes cron/webhook instinct — except the trigger graduates from "time/webhook" to "learned perceptual salience." That's the upgrade.

---

## Part 6 — Recommended OneMind architecture (the synthesis)

**Do NOT collapse everything into one model. Build the stack. Five layers:**

```
┌─────────────────────────────────────────────────────────────────┐
│ L5  PRESENCE / INTERFACE                                        │
│     GPT-realtime (voice, live reasoning) + wearables            │
│     (glasses HUD, watch telemetry, laser projector output)      │
│     → How you and OneMind talk. The "face."                     │
├─────────────────────────────────────────────────────────────────┤
│ L4  ORCHESTRATOR / META-AGENT  ← THIS is your "manager"         │
│     Hermes (Legacy) — memory, tool-calling, agent routing,      │
│     autonomy ladder L0–L3, cron/webhooks, NATS fabric.          │
│     LLM-class. This is the prefrontal cortex. NOT a WM.         │
├─────────────────────────────────────────────────────────────────┤
│ L3  SPECIALIST AGENTS  (domain LLM agents)                      │
│     Infra, security, comms, business, research — your fleet.    │
├─────────────────────────────────────────────────────────────────┤
│ L2  WORLD-MODEL CONTROL LAYER  ← NEW. Where WMs live.          │
│     Per-robot/per-domain world models scoring actions in        │
│     latent space. worldforge-style loop. V-JEPA 2-AC / LeWM /   │
│     Cosmos. Fast latent MPC. Sub-second physical planning.      │
├─────────────────────────────────────────────────────────────────┤
│ L1  ALWAYS-ON PERCEPTION  ← NEW. The "24/7 watch."            │
│     Local WM encoders + small perception nets running           │
│     continuously on cameras/sensors/wearables. Token-free.      │
│     Emits latents + SURPRISE signals → interrupts L4.           │
├─────────────────────────────────────────────────────────────────┤
│ L0  PHYSICAL FABRIC                                             │
│     Robots (DIMOS, OM1, Go2, G1), drones, vehicles, 3D         │
│     printers, cameras, projectors, IoT — via NATS.             │
└─────────────────────────────────────────────────────────────────┘
```

**Why this is right:**
- The **omnipresent "manager"** you want = **L4 Hermes/Legacy**, unchanged. It already does the orchestration, memory, tools, autonomy ladder. The WM does NOT take this job.
- **World models slot in at L1–L2**, giving you the two things LLMs genuinely can't do: (a) token-free 24/7 perception/anomaly-watch, (b) fast physical action planning for robots/drones/vehicles.
- **The interrupt loop** (L1 surprise → wake L4) is how you get "always-on watching my back" *without* paying tokens to think continuously.
- Everything rides your **existing NATS fabric** — WMs become just more nodes/subjects on the mesh.

---

## Part 7 — Concrete next steps (when you're ready, no action taken yet)

**Phase 0 — Learn the loop (cheap, today-grade):**
- `uv add worldforge-ai` → run the mock provider loop. Zero GPU/robot. Lets you *feel* the policy→score→plan→execute cycle. This is the single best way to internalize the architecture.

**Phase 1 — First real WM (robotics node):**
- Pull **V-JEPA 2-AC** (open weights) onto the box with a GPU. Wire it as a worldforge `score`/`predict` provider behind one robot (Go2 is the natural first target given your roadmap).
- Goal-image planning: give it a target frame, let it plan the action chunk, DIMOS/OM1 executes.

**Phase 2 — Always-on perception (the "watch my back" layer):**
- Stand up a small local WM encoder on a fixed camera (House/Barn/Office node). Stream latents + surprise score onto NATS. Wire a Hermes webhook to the surprise signal → Legacy gets interrupted only when something physically notable happens.
- ⚠️ **Subject naming superseded:** the provisional `om.perception.*` names below are corrected to the spine-compliant `tel/evt/gov.sensor.*` contract in **`01 The Fabric/ONEMIND-LAYER-STACK.md §3.2`** (the canonical L1 definition). The *reasoning* here stands; use that doc for the actual subjects when building.

**Phase 3 — Wearables as sensors+interface:**
- Glasses/watch telemetry → L1 perception. Glasses HUD / laser projector → L5 output surface. The WM consumes egocentric video (V-JEPA 2 is *trained* on egocentric data — strong fit).

**Open item:** Re-confirm with your friend the exact name behind "nexarbadas 2.0" (likely Nexar/BADAS driving model). Park unless AV node goes active.

---

## One-line answer to "what do you think / what do you recommend"

**Don't chase one model that does everything — that model doesn't exist, and a world model is specifically *not* it.** Build the stack: keep Hermes/Legacy as your omnipresent orchestrator (L4), add world models as a NEW bottom layer (L1 always-on perception + L2 physical control) that finally gives OneMind the token-free 24/7 watching and fast physical action your friend is right to be excited about. The world model is the cerebellum; Legacy stays the mind.

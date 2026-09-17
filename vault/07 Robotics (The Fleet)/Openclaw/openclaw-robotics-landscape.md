---
id: openclaw-robotics-landscape
title: "OpenClaw × Robotics — Landscape Overview"
type: research
status: active
created: 2026-05-04
owner: zeus
tags: [type/research, domain/robotics, domain/agents, domain/ai]
---

# OpenClaw × Robotics — Landscape Overview

> OpenClaw is jumping off the screen and into physical hardware. 2026 is the year AI agents step out of the chatbox and into the physical world.

**Source:** [When AI Gets Physical Hands — EvoAI Labs (Mar 2026)](https://evoailabs.medium.com/when-ai-gets-physical-hands-a-review-of-openclaw-on-the-unitree-g1-and-other-robots-0fbf06a1d4c8)

---

## What OpenClaw Is (Physical Context)

OpenClaw (formerly Clawdbot/Moltbot) is the open-source AI agent framework created by Peter Steinberger (now at OpenAI). It connects LLMs (Claude, DeepSeek, OpenAI GPTs) to tools, desktops, and now — physical robots. You interact via Telegram, Discord, Signal, or WhatsApp.

**Key insight:** OpenClaw is not the brain. It's the hands. It provides tool execution, multi-step task automation, and physical hardware interfaces. The brain problem is what QwestorClaw and SingularityNET are solving.

---

## Hardware Integrations

### Unitree G1 (Humanoid Robot)
- **Price:** ~$16,000
- **Specs:** Bipedal, 3D LiDAR, depth cameras, advanced joint mobility
- **Integration:** `unitree-robot` OpenClaw skill — control via text commands ("forward 1m", "turn left 45")
- **Camera:** Insight9 AI stereo camera — RGB-D frames streamed back to Telegram
- **Navigation:** Optional TinyNav for autonomous path planning + obstacle avoidance
- **Folder:** `unitree (Unitree Research)/`
- **Links:**
  - [Unitree G1](https://www.unitree.com/g1)
  - [RosClaw Plugin (GitHub)](https://github.com/PlaiPin/rosclaw)

### Aero Hand Open (3D-Printed Robotic Hand)
- **Type:** Fully 3D-printed, 16-joint robotic hand
- **Integration:** OpenClaw agent (Claude) + USB camera for visual self-verification
- **Capability:** Autonomous firmware calibration, gesture formation (fist, peace sign), live Telegram narration
- **Source:** Adafruit feature project
- **Folder:** `aero-hand (Aero Open Hand)/`

### NERO 7-Axis Robotic Arm
- **Type:** Industrial-grade 7-axis robotic arm
- **SDK:** pyAgxArm
- **Integration:** Natural language → OpenClaw → auto-generates + executes Python kinematics scripts
- **Key value:** No manual kinematics coding — describe the movement, OpenClaw writes the code
- **Folder:** `nero-arm (NERO 7-Axis Arm)/`

---

## Platforms & Frameworks

### Peaq Robotics SDK
- **What it is:** Blockchain network making robots "OpenClaw-ready" at fleet scale
- **Capabilities:**
  - Robots connect to shared OpenClaw agents
  - Download reusable skills
  - Orchestrate multi-machine workflows
  - Machine-to-machine payments (on-chain financial entities)
- **Relevance to One Mind:** DePIN (Decentralized Physical Infrastructure Network) — robots as autonomous economic agents
- **Folder (Research):** `peaq-sdk (Peaq Robotics SDK)/`
- **Folder (Investment):** `GE 79 Investments/Crypto/peaq (Peaq DePIN Network)/`
- **Links:**
  - [Peaq Robotics SDK Blog](https://www.peaq.xyz/blog/make-your-robot-openclaw-ready-with-robotics-sdk)
  - [OpenClaw Robot Skill](https://playbooks.com/skills/openclaw/skills/openclaw-robotics)

### QwestorClaw
- **What it is:** Architecture by Ben Goertzel (SingularityNET) pairing OpenClaw's execution layer with SingularityNET's persistent cognitive memory
- **Problem it solves:** OpenClaw has no persistent memory or self-model — it forgets between sessions. Current LLMs lack high-level abstraction, persistent memory, and proactive self-understanding needed for AGI.
- **Solution:** Secure integration of OpenClaw (hands) + SingularityNET AtomSpace (brain) = cognitive flywheel
- **Folder:** `SO 01 Intelligence/Agents/agent-docs/qwestorclaw/`
- **Links:**
  - [Ben Goertzel Substack on OpenClaw](https://bengoertzel.substack.com/p/openclaw-amazing-hands-for-a-brain)

### SingularityNET Cognitive Framework
- **What it is:** Ben Goertzel's AGI platform — provides the persistent memory layer QwestorClaw pairs with OpenClaw
- **Core component:** **AtomSpace** — a hypergraph knowledge store where facts, patterns, and relationships persist permanently across sessions
- **Architecture:** OpenCog — multiple cognitive processes running in parallel (logical reasoning, pattern mining, evolutionary learning) all sharing the same AtomSpace
- **Cognitive synergy:** Different cognitive algorithms mutually assist each other through shared memory
- **Goal:** Cumulative world model that grows with every interaction — not just context window, but a persistent knowledge graph
- **Folder:** `SO 01 Intelligence/Agents/agent-docs/singularitynet/`

---

## The Cognitive Flywheel

> "A cognitive flywheel where AI can learn cumulatively, expanding its abilities to eventually bootstrap genuine AGI." — Ben Goertzel

The architecture Goertzel describes for QwestorClaw:

```
Agent executes task (OpenClaw)
    ↓
Writes learnings to AtomSpace (SingularityNET)
    ↓
Next task starts with richer context
    ↓
Better execution → more learnings
    ↓
Repeat → compounding intelligence
```

**Hard to start. Once spinning — self-accelerating.**

**One Mind mapping:** This is exactly what CODE×'s **× Multiply** step is. Every Execute cycle feeds back into Capture → richer Organize → smarter Direct → better Execute. QwestorClaw is building this for AI agents. One Mind is building it for a human life. Same architecture, different substrate.

- **Folder:** `SO 00 Framework/doctrine/cognitive-flywheel/`

---

## Security Warning

> Cybersecurity experts warn that giving an AI agent autonomous control over physical hardware — combined with human identity credentials — creates massive governance and security challenges.

Physical AI agents inherit all OpenClaw security considerations plus:
- Physical harm potential (robots can cause injury)
- Credential exposure (human identity + physical access)
- Firmware attack surface
- Network bridge between digital and physical

**Relevant vault:** `SO 04 Security/Protocols/` — machine-level ROE applies here

---

## Vault Map

| Topic | Location |
|-------|----------|
| Unitree G1 research | `SO 06 Research/Robotics/unitree/` |
| Aero Hand Open | `SO 06 Research/Robotics/aero-hand/` |
| NERO 7-Axis Arm | `SO 06 Research/Robotics/nero-arm/` |
| Peaq SDK research | `SO 06 Research/Robotics/peaq-sdk/` |
| Peaq as investment | `GE 79 Investments/Crypto/peaq/` |
| QwestorClaw | `SO 01 Intelligence/Agents/agent-docs/qwestorclaw/` |
| SingularityNET | `SO 01 Intelligence/Agents/agent-docs/singularitynet/` |
| Cognitive Flywheel | `SO 00 Framework/doctrine/cognitive-flywheel/` |


---

## Related

- [[00-24 SO (Sovereign Operations)/06 Research (R&D) (The Lab)/06 Research|06 Research]]

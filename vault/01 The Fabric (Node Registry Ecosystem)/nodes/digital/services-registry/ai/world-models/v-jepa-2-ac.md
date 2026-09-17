---
title: V-JEPA 2-AC — Action-Conditioned World Model for Robot Control
type: research-brief
domain: SO
tags: [world-models, jepa, meta, robotics, action-conditioned, control, mpc, franka, embodied-ai, onemind-ready]
date: 2026-06-08
status: reference
links:
  paper: https://arxiv.org/abs/2506.09985
  github: https://github.com/facebookresearch/vjepa2
  weights: https://huggingface.co/facebook/vjepa2-ac-vith-fpc16-256
  blog: https://ai.meta.com/blog/v-jepa-2-world-model-physical-reasoning-video/
---

# V-JEPA 2-AC — Action-Conditioned World Model for Robot Control

## One-Line

The **control-ready checkpoint** of V-JEPA 2. Post-trained on ~62 hours of robot trajectories, it takes robot actions as inputs and predicts future world states — enabling zero-shot robot planning via goal-image latent MPC.

## What "-AC" Means

`-AC` = **Action-Conditioned**. This is a **post-trained variant** of the V-JEPA 2 base model (see `v-jepa-2.md`). The base model predicts video purely from passive observation. The -AC checkpoint additionally takes a **robot action vector** as input to the predictor, so it can answer:

> *"If I move the arm to position X with velocity V, what does the scene look like next?"*

Without -AC: "what will probably happen next?"  
With -AC: "what happens **if I do action A**?" ← this is what control requires

## Provenance

| Field | Value |
|---|---|
| **Authors** | Same as V-JEPA 2 — Meta AI / FAIR |
| **Institution** | Meta AI Research |
| **Paper** | "V-JEPA 2: Self-Supervised Video Models Enable Understanding, Prediction and Planning" |
| **Released** | June 2025 (same release as V-JEPA 2 base) |
| **License** | Open weights, commercially permitted |
| **Base** | Post-trained from V-JEPA 2 (1.2B ViT-H) |

## How It's Made

### Stage 1 — Passive Pretraining (inherited from V-JEPA 2 base)

```
1M+ hours of internet video
    → masked spatio-temporal latent prediction
    → learns world physics without any labels
```

### Stage 2 — Action-Conditioned Post-Training

```
~62 hours of Franka arm robot trajectories
    (proprioceptive joint states + gripper + end-effector actions)
    → predictor conditioned on action vector a_t
    → learns: p(z_{t+1} | z_t, a_t)
    → model predicts next latent state given current state + action
```

**Critical insight:** Only **62 hours** of robot data are needed because Stage 1 already encoded the physics of the world. Stage 2 only teaches the model *how to relate actions to those already-understood dynamics*.

Compare: a model starting from scratch would need orders of magnitude more robot data.

## Architecture

| Component | Description |
|---|---|
| **Backbone** | ViT-H (inherited from V-JEPA 2 base), ~1.2B params |
| **Predictor modification** | Action vector `a_t` concatenated or cross-attended into predictor input |
| **Action space** | Franka arm proprioception (7-DoF joints + gripper) — but architecture generalizes |
| **Input** | Video frames + action sequence |
| **Output** | Predicted latent states z_{t+1}, z_{t+2}, ... z_{t+k} |
| **Inference** | Single GPU; 64 frames at 256×256 |

## How It Plans (Goal-Image Latent MPC)

V-JEPA 2-AC enables **goal-image planning** via latent Model Predictive Control (MPC):

```
GOAL: reach a target state (given as a goal image G)

1. Encode goal image G → target latent z_goal
2. Encode current state (video) → current latent z_0
3. POLICY proposes N candidate action sequences {a_0, ..., a_k}^i  (i = 1..N)
4. V-JEPA 2-AC SCORES each sequence:
       for each action sequence: roll out z_1, z_2, ..., z_k
       score = ||z_k - z_goal||  (latent distance to goal)
5. Pick lowest-cost sequence → execute first action a_0
6. Observe new state → repeat (receding horizon)
```

This is **zero-shot**: the model was trained on Franka data but the planning algorithm generalizes to any robot by swapping the action space.

### Zero-Shot Generalization

V-JEPA 2-AC deployed zero-shot on a **Franka Emika arm** for tasks including:
- Pick and place
- Block stacking
- Goal-image reaching

Performance: **zero-shot, no task-specific fine-tuning** — the physics from Stage 1 + action dynamics from Stage 2 are sufficient.

## What It Needs at Runtime

| Requirement | Detail |
|---|---|
| **GPU** | ~1 A100 or equivalent for real-time inference |
| **Goal specification** | A goal image (photo of desired end state) |
| **Robot interface** | Action publisher + observation subscriber (ROS2/NATS/direct) |
| **Policy proposer** | CEM, random shooting, or an LLM goal-setter (any proposer works) |
| **Frame rate** | 15–30 fps recommended for real-time control |

## OneMind Integration

### NATS Message Subjects

```
# Publish: camera observations (from robot node)
tel.robot.{id}.camera

# Publish: proprioceptive state (from robot node)
tel.robot.{id}.joints

# Subscribe: planned action command (from WM controller)
cmd.robot.{id}.action

# Publish: world model anomaly / surprise
evt.robot.{id}.wm_surprise
```

### Control Loop Integration

```
                     ┌─────────────────────────────┐
                     │     Legacy/Hermes (L4)       │
                     │  sets goal, monitors alerts  │
                     └──────────────┬──────────────┘
                                    │ goal image
                     ┌──────────────▼──────────────┐
                     │    V-JEPA 2-AC Controller   │
                     │  latent MPC, action planning │
                     └──────────────┬──────────────┘
                                    │ cmd.robot.{id}.action
                     ┌──────────────▼──────────────┐
                     │   Robot Node (L0 Fabric)    │
                     │  DIMOS / OM1 / Go2 / G1    │
                     └─────────────────────────────┘
```

### Fleet Deployment Targets

| Device | Priority | Notes |
|---|---|---|
| **DIMOS** | High | Arm manipulation tasks, pick-place, reach |
| **OM1** | High | Mobile manipulation, goal-image navigation |
| **Go2** | Medium | Whole-body action conditioning (different action space) |
| **G1** (planned) | High | Humanoid — exactly the Franka-class application |
| **Drones** | Medium | Action space: throttle/pitch/roll/yaw — requires domain adaptation |

## Pitfalls

- **-AC ≠ base model** — the base V-JEPA 2 checkpoint CANNOT do control. Always use the `-ac` checkpoint for planning.
- **Action space matters** — the published -AC is trained on Franka proprioceptive actions. For Go2 or drones, you need to either: (a) fine-tune Stage 2 with your robot's action data, or (b) map/normalize the action space carefully.
- **Goal image quality** — goal-image planning requires a clear, unambiguous image of the desired end state. Blurry or ambiguous goal images degrade planning quality.
- **Horizon length** — longer planning horizons (more steps) accumulate latent prediction error. Use receding-horizon (short horizon + re-plan) for reliability.
- **Not a replacement for SLAM** — this is a *control* world model, not a *mapping* system. Pair with explicit mapping for navigation.

## Relation to Other Models

```
V-JEPA 2 (base, Jun 2025)
    └── V-JEPA 2-AC (action-conditioned) ← THIS FILE
            ↕ pairs with
        Policy proposer (GR00T / LeRobot / Cosmos-Policy / LLM-generated goals)
```

In the canonical control loop:
- **V-JEPA 2-AC = WORLD MODEL / SCORER** (scores action consequences)
- **GR00T / LeRobot = POLICY PROPOSER** (generates candidate action sequences)
- **Legacy/Hermes = GOAL SETTER** (sets what end state to reach)

## Resources

- **Paper:** https://arxiv.org/abs/2506.09985
- **GitHub:** https://github.com/facebookresearch/vjepa2
- **Weights (HuggingFace):** `facebook/vjepa2-ac-vith-fpc16-256`
- **Blog:** https://ai.meta.com/blog/v-jepa-2-world-model-physical-reasoning-video/
- **WorldForge provider:** `vjepa2-ac` (`score`, `plan` capabilities)
- **Related OneMind doc:** `v-jepa-2.md` (base model), `world-models-vs-llms-onemind-architecture-2026.md`

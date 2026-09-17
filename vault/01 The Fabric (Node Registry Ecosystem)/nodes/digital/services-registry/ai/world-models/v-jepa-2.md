---
title: V-JEPA 2 — Video World Model with Zero-Shot Robot Planning
type: research-brief
domain: SO
tags: [world-models, jepa, meta, video, robotics, zero-shot, embodied-ai, onemind-ready]
date: 2026-06-08
status: reference
links:
  paper: https://arxiv.org/abs/2506.09985
  github: https://github.com/facebookresearch/vjepa2
  blog: https://ai.meta.com/blog/v-jepa-2-world-model-physical-reasoning-video/
  weights: https://huggingface.co/facebook/vjepa2-vith-fpc64-256
---

# V-JEPA 2 — Video World Model with Zero-Shot Robot Planning

## One-Line

Meta's June 2025 video world model — 1.2B params, SOTA physical reasoning, zero-shot robot planning from goal images. The primary usable robotics-layer world model for OneMind's fleet.

## Provenance

| Field | Value |
|---|---|
| **Authors** | Mido Assran, Adrien Bardes, David Fan, Quentin Duval, Abhimanyu Dubey, Piotr Bojanowski, Pascal Vincent, Yann LeCun, Nicolas Ballas et al. |
| **Institution** | Meta AI Research / FAIR |
| **Paper** | "V-JEPA 2: Self-Supervised Video Models Enable Understanding, Prediction and Planning" |
| **arXiv** | 2506.09985 (est.; see official link) |
| **Released** | June 2025 |
| **License** | Open weights, commercially permitted (Meta Research License) |
| **Lineage** | JEPA family — direct successor to V-JEPA (2024) |

## What It Does

V-JEPA 2 is a **video world model** with two modes:

1. **Base model (V-JEPA 2)** — passive video understanding, physical reasoning, video QA, zero-shot prediction
2. **Action-conditioned model (V-JEPA 2-AC)** — post-trained with robot action data → full robot planning/control

This file covers the **base model**. See `v-jepa-2-ac.md` for the control-ready variant.

### What "World Model" Means Here

V-JEPA 2 encodes the world into a **compact latent representation** and can predict forward in that latent space: *"given what I've seen, what will happen next?"*

This enables:
- **Goal-image planning:** "I want to reach state Z — which action sequence gets me there?"
- **Physical reasoning:** Understands object permanence, causality, collision, gravity — without language
- **Zero-shot generalization:** Never saw a Franka arm during pretraining, yet plans for it

## Training: Two Stages

### Stage 1 — Actionless Pretraining (V-JEPA 2 base)

```
1M+ hours of internet video (passive, no labels, no actions)
    → masked spatio-temporal prediction in latent space
    → learns rich physical world representation
```

Identical mechanism to V-JEPA: predict masked video regions in latent space. Scale: 1M+ hours of video. Result: a model that *understands physics*.

### Stage 2 — Action-Conditioned Post-Training (→ V-JEPA 2-AC)

```
~62 hours of robot trajectory data (Franka arm)
    → predictor conditioned on robot actions
    → learns: "if I do action A, what happens?"
```

See `v-jepa-2-ac.md` for details on the action-conditioned checkpoint.

## Architecture

| Component | Description |
|---|---|
| **Backbone** | Video Vision Transformer (ViT-H) |
| **Parameters** | ~1.2B |
| **Context encoder** | Encodes observed video frames into latent embeddings |
| **Target encoder** | EMA copy — produces stable prediction targets |
| **Predictor** | Spatio-temporal transformer — predicts future latent states |
| **Input resolution** | 256×256, 64 frames per clip |
| **Patch size** | 16×16 spatial, 2-frame temporal |

No pixel decoder. No generative component. Pure latent prediction.

## Key Properties

| Property | Detail |
|---|---|
| **Latent-space only** | Never reconstructs pixels — efficient, semantic |
| **Action-free pretraining** | Stage 1 requires zero robot data |
| **Zero-shot robot planning** | Plans Franka arm tasks without robot-specific pretraining |
| **Open weights** | Commercially permitted |
| **Physical benchmarks** | SOTA on IntPhys 2, MVPBench, CausalVQA |

## Benchmarks

| Benchmark | What It Tests | V-JEPA 2 Result |
|---|---|---|
| **IntPhys 2** | Physically plausible vs. implausible (violation-of-expectation) | SOTA |
| **MVPBench** | Minimal Video Pairs — physical understanding, shortcut-resistant | SOTA |
| **CausalVQA** | Causal cause-and-effect (counterfactual, anticipation, planning) | SOTA |
| **Something-Something V2** | Fine-grained action recognition | Competitive |

Compared to video-language models (GPT-4o, Gemini): V-JEPA 2 **without language** outperforms language-enabled models on *physical reasoning*. Key result: language is not needed for physical understanding.

## OneMind Relevance

**Role: Primary Robotics-Layer World Model (Base)**

V-JEPA 2 is the world model you'd use for:

- **Passive perception** at L1 — encode sensor video 24/7, token-free, emit surprise signal when physics is violated
- **Pre-training foundation** for downstream robot adaptation (via V-JEPA 2-AC)
- **Physical scene understanding** — does the world state make sense? Is something anomalous?

### Fleet Applications

| Device | Use |
|---|---|
| **DIMOS** | Perception backbone, anomaly detection, pre-training for arm control |
| **OM1** | Video understanding, environment modeling |
| **Go2** | Navigation world model (with action conditioning) |
| **G1** (planned) | Humanoid whole-body planning |
| **Drones** | Aerial scene prediction, obstacle anticipation |
| **Greenhouse** | Anomaly detection, plant state prediction |

### Hardware Requirements

- **Inference:** Single GPU (A100 or equivalent) for real-time video encoding
- **Fine-tuning (Stage 2):** ~62h robot data + GPU cluster (but Meta did this — you use the V-JEPA 2-AC checkpoint)
- **Edge inference:** Feasible on a workstation; not yet optimized for embedded hardware (use ONNX/TRT export)

## Stack Placement

```
L4 ORCHESTRATOR  Legacy/Hermes — sets goals, reasons in language
──────────────────────────────────────────────────────
L2 WM CONTROL    V-JEPA 2-AC (action-conditioned) ← deploy for robot control
L1 ALWAYS-ON     V-JEPA 2 base ← deploy for 24/7 passive perception + anomaly detection
L0 FABRIC        cameras, robots, drones, sensors
```

The base model runs at **L1** (passive, always-on). The -AC checkpoint runs at **L2** (active control).

## Relation to Other Models

```
I-JEPA (Jan 2023) — image, foundational
    └── V-JEPA (Jan/Apr 2024) — video, temporal dynamics
            └── V-JEPA 2 (Jun 2025) — 1.2B, zero-shot robot planning ← THIS FILE
                    └── V-JEPA 2-AC — action-conditioned checkpoint ← for control
```

## Resources

- **Paper:** https://arxiv.org/abs/2506.09985
- **GitHub:** https://github.com/facebookresearch/vjepa2
- **Blog:** https://ai.meta.com/blog/v-jepa-2-world-model-physical-reasoning-video/
- **Weights (HuggingFace):** `facebook/vjepa2-vith-fpc64-256`
- **WorldForge provider:** `vjepa2` (`score`, `embed`, `predict` capabilities)

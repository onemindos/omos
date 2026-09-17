---
title: V-JEPA — Video Joint-Embedding Predictive Architecture
type: research-brief
domain: SO
tags: [world-models, jepa, meta, video, self-supervised, temporal-dynamics, embodied-ai]
date: 2026-06-08
status: reference
links:
  paper: https://arxiv.org/abs/2404.08471
  github: https://github.com/facebookresearch/jepa
  blog: https://ai.meta.com/blog/v-jepa-yann-lecun-self-supervised-learning/
---

# V-JEPA — Video Joint-Embedding Predictive Architecture

## One-Line

Meta 2024 extension of I-JEPA to video — learns physical dynamics by predicting masked spatiotemporal regions in latent space from passive (unlabeled) video. The foundation for V-JEPA 2.

## Provenance

| Field | Value |
|---|---|
| **Authors** | Adrien Bardes, Quentin Duval, Jean Ponce, Yann LeCun et al. |
| **Institution** | Meta AI Research / FAIR |
| **Paper** | "Revisiting Feature Prediction for Learning Visual Representations from Video" |
| **arXiv** | 2404.08471 |
| **Published** | January 2024 (blog), April 2024 (arXiv) |
| **Lineage** | JEPA family (LeCun) — direct successor to I-JEPA |

## What It Does

V-JEPA extends the I-JEPA framework from still images to **video** — learning spatio-temporal representations by predicting masked volumes of a video in latent space, without reconstructing pixels.

### The Core Idea

```
video clip → context encoder → context embedding
                                       ↓
masked spatio-temporal tubes ─→ PREDICTOR → predicted latent z_t1..tn
                                       ↓
                             compare to: target encoder → target latent z_t1..tn
                                       ↓
                              minimize prediction error over time
```

- **Masking strategy:** Large spatiotemporal blocks — forces prediction of what happens *next* in time
- **No pixel reconstruction:** The model never generates frames — it only predicts in latent space
- **Passive video only:** No action labels, no robot data — just watching YouTube/video

### Why This Matters for Physical Dynamics

By predicting *what happens next* in a video (in latent space), the model must internalize:
- Object permanence — objects continue existing when occluded
- Motion continuity — objects follow physical trajectories
- Causal relations — one event leads to another
- Gravity, collision, inertia (implicitly)

This is what "physical world model" means — the representations encode physics without any explicit physics labels.

## What It Does NOT Do

- ❌ Not action-conditioned — cannot take actions as inputs
- ❌ Not suitable for robot control directly (that requires V-JEPA 2-AC)
- ❌ No planning capability in V-JEPA 1 — predicts but cannot be queried "what happens if I do X"
- ❌ Not a generative model — cannot produce video frames

## Architecture

| Component | Description |
|---|---|
| **Context encoder** | Video ViT (ViT-L/16×2 or ViT-H/16×2) — processes unmasked video patches |
| **Target encoder** | EMA copy of context encoder — stable target representations |
| **Predictor** | Spatio-temporal transformer — conditioned on mask positions, predicts latent targets |
| **Masking** | Multi-block masking of large spatio-temporal regions (not random individual patches) |

**Input:** Video clips (16–32 frames)  
**Output:** Latent embedding of the full video (context + predicted regions)  
**No decoder, no pixel output**

## Parameters

Primary reported checkpoints:
- **ViT-L (V-JEPA):** ~307M params
- **ViT-H (V-JEPA):** ~632M params

## Benchmarks

Evaluated on:

| Benchmark | Task | Result |
|---|---|---|
| Something-Something V2 | Action recognition (frozen probe) | Competitive with prior SOTA |
| Kinetics-400 | Action recognition (frozen probe) | Strong linear probe performance |
| AVA | Spatiotemporal action detection | Outperforms pixel-reconstruction baselines |
| **Physical reasoning (zero-shot)** | Tested in follow-up work | Establishes basis for V-JEPA 2 improvements |

Key finding: V-JEPA representations generalize to downstream tasks *without fine-tuning* better than masked pixel reconstruction (VideoMAE, etc.).

## Relation to Other Models

```
I-JEPA (2023) — image proof-of-concept
    └── V-JEPA (2024) — extends to video, temporal dynamics ← YOU ARE HERE
            └── V-JEPA 2 (Jun 2025) — 1.2B params, zero-shot robot planning
                    └── V-JEPA 2-AC — action-conditioned, Franka arm control
```

V-JEPA is the **temporal dynamics bridge** between I-JEPA (static images) and V-JEPA 2 (full world model with robot control).

## OneMind Relevance

**Role: Intermediate Reference / Precursor to V-JEPA 2**

V-JEPA itself is NOT deployable for robot control — it has no action-conditioning. However:

1. Understanding V-JEPA explains *why* V-JEPA 2 works — it's the pretraining stage
2. V-JEPA's passive-video pretraining on 1M+ hours of video IS the first stage of V-JEPA 2 training
3. V-JEPA's architecture (ViT + spatio-temporal masking + latent prediction) is directly inherited by V-JEPA 2

**For Zeus's fleet:** V-JEPA is a precursor. The model you'd actually wire to DIMOS, OM1, or Go2 is V-JEPA 2-AC, not V-JEPA.

## Stack Placement

```
L2 WM CONTROL    → V-JEPA 2-AC (control-ready), not V-JEPA
L1 ALWAYS-ON     → could use V-JEPA for passive scene understanding (perception only)
```

V-JEPA could theoretically run at L1 as a pure perception backbone (encode video, emit embeddings), but V-JEPA 2 is strictly better. The main value of understanding V-JEPA is intellectual lineage.

## Resources

- **Paper:** https://arxiv.org/abs/2404.08471
- **GitHub:** https://github.com/facebookresearch/jepa
- **Blog:** https://ai.meta.com/blog/v-jepa-yann-lecun-self-supervised-learning/
- **Weights:** Available on HuggingFace (`facebook/vjepa-vith16-*`)

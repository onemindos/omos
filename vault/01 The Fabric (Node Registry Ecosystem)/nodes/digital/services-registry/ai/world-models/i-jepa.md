---
title: I-JEPA — Image Joint-Embedding Predictive Architecture
type: research-brief
domain: SO
tags: [world-models, jepa, meta, self-supervised, computer-vision, embodied-ai]
date: 2026-06-08
status: reference
links:
  paper: https://arxiv.org/abs/2301.08243
  github: https://github.com/facebookresearch/ijepa
  blog: https://ai.meta.com/blog/yann-lecun-ai-model-i-jepa/
---

# I-JEPA — Image Joint-Embedding Predictive Architecture

## One-Line

Foundational Meta 2023 model that proved self-supervised image representations can be learned by predicting masked patches **in latent space** — never reconstructing pixels.

## Provenance

| Field | Value |
|---|---|
| **Authors** | Mahmoud Assran, Quentin Duval, Ishan Misra, Piotr Bojanowski, Pascal Vincent, Michael Rabbat, Yann LeCun, Nicolas Ballas |
| **Institution** | Meta AI Research |
| **Paper** | "Self-Supervised Learning from Images with a Joint-Embedding Predictive Architecture" |
| **arXiv** | 2301.08243 |
| **Published** | January 2023 |
| **Lineage** | JEPA family (LeCun) |

## What It Does

I-JEPA is a **self-supervised image representation learning** model. Given a masked image, it predicts the latent representations of the missing patches — without ever reconstructing the actual pixels.

### The Core Idea

```
input image → context encoder → context embedding
                                        ↓
masked regions ──────────────→ PREDICTOR → predicted latent z
                                        ↓
                              compare to: target encoder → target latent z
                                        ↓
                               minimize prediction error
```

- The **target encoder** produces stable latent targets (uses EMA / stop-gradient of context encoder)
- The **predictor** learns to fill in masked regions in latent space, not pixel space
- This forces the model to learn **semantic, physical representations** rather than low-level textures

### Why Latent Prediction Matters

- Pixel-space reconstruction (MAE, BEiT, etc.) forces the model to memorize texture details — wasted capacity
- Latent-space prediction forces the model to learn *meaning* — what is there, not what it looks like
- This is LeCun's central thesis: intelligence requires internal models of the world, not generative pixel models

## What It Does NOT Do

- ❌ Not action-conditioned — cannot plan or control anything
- ❌ Not temporal — single images only (V-JEPA extends this to video)
- ❌ Not deployable for robotics control directly
- ❌ Not a world model in the full sense — no dynamics, no forward prediction over time

## Architecture

| Component | Description |
|---|---|
| **Context encoder** | ViT (Vision Transformer) variant — encodes unmasked patches |
| **Target encoder** | EMA copy of context encoder — provides stable latent targets |
| **Predictor** | Narrow transformer — maps context to predictions for masked regions |
| **Masking** | Large contiguous block masks (not random patches) — forces long-range reasoning |

**Key design choice:** The predictor is intentionally *narrow* (small) to prevent it from learning to copy the encoder — it must actually predict.

## Parameters

Trained in multiple sizes mirroring ViT variants (ViT-H/14 is the largest published checkpoint). Exact sizes not the focus — the architecture, not scale, is the contribution.

## Benchmarks

I-JEPA was evaluated on standard downstream classification benchmarks with frozen representations:

- **ImageNet-1K linear probe** — competitive with or exceeds pixel-reconstruction baselines (MAE, data2vec)
- **Low-shot classification** — strong performance with limited labeled data
- **Semantic segmentation** — representations generalize well to dense prediction
- **Ablations** confirm: latent prediction > pixel prediction for semantic quality

## OneMind Relevance

**Role: Foundational / Historical Reference**

I-JEPA is not deployable for control or robotics — it's the proof-of-concept that established the JEPA paradigm. Its contributions:

1. Validated that latent-space prediction produces higher-quality representations than pixel reconstruction
2. Established the context encoder + target encoder + predictor architecture that V-JEPA and V-JEPA 2 inherit
3. Is the intellectual foundation for the entire JEPA lineage including LeWM

**Relevance to Zeus's fleet:** Indirect — it's the ancestor of V-JEPA 2 which IS relevant to DIMOS, OM1, Go2, G1, drones. Understanding I-JEPA gives intuition for why V-JEPA 2 works.

## Relation to Other Models

```
I-JEPA (2023) — image, foundational proof
    └── V-JEPA (2024) — extends to video, temporal dynamics
            └── V-JEPA 2 (Jun 2025) — adds zero-shot robot planning
                    └── V-JEPA 2-AC — action-conditioned, control-ready
```

## Stack Placement

```
L4 ORCHESTRATOR  Legacy/Hermes
L3 SPECIALISTS   domain agents
──────────────────────────────
L2 WM CONTROL    V-JEPA 2-AC, LeWM (NOT I-JEPA — it's a repr. model, not a controller)
L1 ALWAYS-ON     perception + anomaly detection
L0 FABRIC        sensors, robots, drones
```

I-JEPA itself does not slot into L2. It is the **upstream pretraining paradigm** that informs L2 model design.

## Resources

- **Paper:** https://arxiv.org/abs/2301.08243
- **GitHub:** https://github.com/facebookresearch/ijepa
- **Blog:** https://ai.meta.com/blog/yann-lecun-ai-model-i-jepa/
- **Weights:** Available on HuggingFace (`facebook/ijepa_vith14_1percent`, etc.)

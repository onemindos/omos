---
id: so-01-ai-models-nvidia
title: "AI Models — NVIDIA"
type: reference
status: active
created: 2026-04-12
modified: 2026-04-12
owner: zeus
tags: [type/reference, status/active, domain/so, ai/models, provider/nvidia]
ai_summary: "NVIDIA AI model reference — 6 vertical AI platforms: Nemotron (agentic), Cosmos (physical AI), GR00T (robotics), DRIVE (autonomous vehicles), Earth-2 (physics/climate), BioNeMo (biomedical). Access via NVIDIA NIM API (OpenAI-compatible)."
---

# NVIDIA AI

> STACK / SPECIALIST provider. NVIDIA operates **6 vertical AI platforms**, each targeting a different physical-world domain.
> All accessible via NVIDIA NIM (Inference Microservices) — a unified OpenAI-compatible API, or as self-hosted containers on GPU hardware.

---

## Overview

NVIDIA is not just a GPU company. They've built complete AI model stacks for:
agentic reasoning, physical world simulation, robotics, autonomous vehicles, earth
science, and biomedicine. Each platform is a separate product line — not just one model family.

- **API Base URL:** `https://integrate.api.nvidia.com/v1`
- **SDK:** OpenAI-compatible — use OpenAI SDK with custom `baseURL`
- **Key Location:** 1Password → SO vault → `NVIDIA NIM API Key`
- **Console / Key:** https://build.nvidia.com (get free API key)
- **Self-host:** NIM containers via NGC (`nvcr.io/nim/...`) — run on any NVIDIA GPU

---

## The 6 NVIDIA AI Platforms

---

### 1. Nemotron — Agentic AI

> **Purpose:** LLMs and reward models optimized for agentic workflows — complex reasoning, tool use, multi-step planning, and coding. Used to build and train AI agents.

| Model | Params | Context | Specialty |
|-------|--------|---------|----------|
| `nvidia/nemotron-3-super-120b-a12b` | 120B (MoE) | 1M | Agentic reasoning, coding, tool calling |
| `nvidia/nemotron-4-340b-instruct` | 340B | 4K | High-capacity reasoning |
| `nvidia/nemotron-mini-4b-instruct` | 4B | 4K | Lightweight, on-device |
| `nvidia/llama-3.1-nemotron-ultra-253b` | 253B | 128K | Ultra-scale reasoning |
| `nvidia/nemotron-4-340b-reward` | 340B | 4K | Reward model for RLHF / agent evaluation |

**What it's for:**
- Building autonomous agents that plan, reason, and use tools
- Training + fine-tuning other models via reward modeling
- Complex multi-step task execution
- Powers NVIDIA's NemoClaw personal AI agent system

**Stack use:** Agent evaluation/reward scoring, complex reasoning fallback, agentic pipeline backbone

---

### 2. Cosmos — Physical AI

> **Purpose:** World Foundation Models (WFMs) that understand and generate **physical world video** — the "language model for the real world." Generates realistic synthetic data of physical environments for training robots, AV systems, and simulations.

| Model | Type | Input → Output |
|-------|------|----------------|
| `nvidia/cosmos-1.0-diffusion-7b-text2world` | World gen | Text → video of physical world |
| `nvidia/cosmos-1.0-diffusion-14b-video2world` | World continuation | Video → extended physical world video |
| `nvidia/cosmos-1.0-guardrail` | Safety | Content filtering for Cosmos outputs |
| `nvidia/cosmos-tokenizer-cv8x8x8` | Tokenizer | Compress video for WFM inputs |
| `nvidia/cosmos-nemotron-34b` | Language | Physical-world-aware LLM |

**What it's for:**
- Generate **synthetic training data** for robotics and AV (massively reduces real-world data collection)
- Simulate physical environments for robot training (manipulation, navigation)
- AV scenario generation (edge cases, rare events)
- Digital twin asset generation

**Stack use:** Synthetic data generation for robot/AV training pipelines

---

### 3. GR00T — Humanoid Robotics AI

> **Purpose:** Foundation model for **humanoid robots**. Gives robots the ability to understand instructions, reason about manipulation tasks, and learn new skills from few demonstrations. Part of the Isaac robotics platform.

| Model | Description |
|-------|-------------|
| `nvidia/gr00t-n1-2b` | 2B parameter humanoid robot action model |
| `nvidia/gr00t-n1.5` | Latest version — improved dexterity + cross-embodiment |
| **Isaac GR00T** | Full SDK wrapping GR00T model (sim-to-real pipeline, RL training) |

**What it's for:**
- Humanoid robot skill learning — imitation learning from video demonstrations
- Cross-embodiment: one model works across different robot hardware designs
- Manipulation: pick/place, assembly, dexterous tasks
- Robot policy training via NVIDIA Isaac Sim (Omniverse-based simulator)

**Components:**
- **GR00T-Mimic:** Generate training data from human video demos
- **GR00T-Control:** Low-level motor control policies
- **Isaac RL:** Reinforcement learning environment for robot skill refinement

**Stack use:** Robotics integration layer — runs on Jetson Orin edge hardware

---

### 4. Alpamayo — Autonomous Vehicles

> **Purpose:** NVIDIA's end-to-end autonomous vehicle AI platform. Alpamayo is the full-stack AV solution — covering perception, prediction, planning, training data generation, and in-car inference. Built on DRIVE Thor compute + Cosmos synthetic data pipelines.

| Component | Description |
|-----------|-------------|
| **Alpamayo** | End-to-end AV foundation model — full stack: perception → prediction → planning |
| **Cosmos AV Blueprints** | Synthetic AV scenario generation via Cosmos WFMs |
| **NVIDIA Hydra-MDP** | Multi-policy end-to-end driving model — generates diverse driving behaviors |
| **DRIVE Orin / Thor** | In-car compute platforms running AV inference |
| **NVIDIAmaps / DriveSim** | HD map AI + photorealistic AV simulation |
| `nvidia/cosmos-dataset-search` | Vector search over large AV video datasets for post-training |

**What it's for:**
- **Alpamayo:** Full-stack AV platform unifying all components into an end-to-end trainable system
- **Hydra-MDP:** Trains AV models to handle diverse driving policies (cautious → aggressive)
- **Cosmos for AV:** Generate rare/dangerous scenarios (construction zones, adverse weather) synthetically
- **DriveSim:** Photorealistic simulation for AV validation without real-world miles
- Full perception → prediction → planning pipeline running on DRIVE Thor hardware

**Stack use:** AV training data pipeline, simulation, policy training

---

### 5. Earth-2 — AI Physics / Climate

> **Purpose:** AI-native digital twin of Earth. Uses neural network surrogates to simulate atmospheric physics, weather, and climate **10,000× faster** than traditional numerical models. Runs on GPU instead of supercomputers.

| Model | Type | Description |
|-------|------|-------------|
| `nvidia/fourcastnet` | Weather | Global data-driven weather forecasting |
| **CorrDiff** | Downscaling | Diffusion model for km-scale climate detail |
| **StormCast** | Storm | High-res storm prediction (NOAA partnership) |
| **Earth-2 Studio** | Platform | Interactive climate/weather AI visualization |
| **PhysicsNeMo** | Physics AI | Neural operator framework for physical simulations |

**What it's for:**
- 10-day global weather forecasts in seconds (vs hours for traditional NWP)
- High-resolution regional climate downscaling (global → city-block scale)
- Renewable energy site assessment (solar/wind yield prediction)
- Flood, hurricane, wildfire risk modeling
- Building simulation surrogate models for engineering (CFD, FEA)

**Stack use:** Environmental data feeds (LE/GE climate analysis), physics simulation research

---

### 6. BioNeMo — Biomedical AI

> **Purpose:** AI for drug discovery, protein structure prediction, molecular design, and genomics. NVIDIA's platform for accelerating life science research with GPU-optimized foundation models.

| Model | Domain | Description |
|-------|--------|-------------|
| `nvidia/esm2nv` | Proteins | ESM-2 protein language model (650M–3B) |
| `nvidia/molmim` | Molecules | Molecular design via masked image modeling |
| `nvidia/diffdock` | Drug docking | Diffusion model for protein–ligand docking |
| `nvidia/equidock` | Docking | Protein–protein docking prediction |
| `nvidia/alphafold2-multimer` | Structure | Protein complex structure prediction |
| `nvidia/generative-chemistry` | Drug design | De novo molecule generation |
| `nvidia/dnabert` | Genomics | DNA sequence language model |

**What it's for:**
- **Drug discovery:** Design novel drug-like molecules, predict binding affinity
- **Protein structure:** Predict 3D protein folding + protein–protein interaction
- **Molecular docking:** Screen millions of drug candidates against target proteins
- **Genomics:** Analyze DNA/RNA sequences for variant effects, gene regulation

**Stack use:** Research pipeline for health/longevity domain (HP/25–49), biomedical IP

---

## API Setup

All NVIDIA NIM models use the OpenAI-compatible format:

```bash
export NVIDIA_API_KEY="nvapi-..."
```

```typescript
import OpenAI from 'openai';

const client = new OpenAI({
  apiKey: process.env.NVIDIA_API_KEY,
  baseURL: 'https://integrate.api.nvidia.com/v1',
});

// Nemotron — agentic reasoning
const response = await client.chat.completions.create({
  model: 'nvidia/llama-3.1-nemotron-ultra-253b',
  messages: [{ role: 'user', content: 'Plan a 3-step research task' }],
  max_tokens: 1024,
});

// Cosmos — text to physical world video
const cosmos = await client.chat.completions.create({
  model: 'nvidia/cosmos-1.0-diffusion-7b-text2world',
  messages: [{ role: 'user', content: 'A robot arm picking up a red cup on a table' }],
});
```

## Self-Hosted via NIM (Docker)

```bash
# Pull NIM container for a specific model
docker run --gpus all -p 8000:8000 \
  -e NVIDIA_API_KEY=$NVIDIA_API_KEY \
  nvcr.io/nim/nvidia/nemotron-mini-4b-instruct:latest

# Then hit localhost:8000/v1 with standard OpenAI format
```

---

## Usage in One Mind Stack

| Use Case | Platform | Model |
|----------|----------|-------|
| **Agentic reasoning / reward eval** | Nemotron | `nemotron-3-super-120b-a12b` |
| **Robot training data generation** | Cosmos + GR00T | Cosmos WFMs + GR00T-N1.5 |
| **AV scenario synthesis** | Alpamayo | Cosmos AV blueprints, Hydra-MDP, DRIVE Thor |
| **Weather / climate modeling** | Earth-2 | FourCastNet + CorrDiff |
| **Drug/molecule research** | BioNeMo | DiffDock + MolMIM |
| **Local agentic inference** | Nemotron | `nemotron-mini-4b-instruct` (Jetson) |

---

## Notes

- **Free tier:** NVIDIA Build (`build.nvidia.com`) offers free serverless API access for development — 1,000 API credits on signup
- **Enterprise:** NIM microservices licensed per GPU node for self-hosting (DGX, Jetson, cloud VMs)
- **Jetson:** Nemotron-Mini and GR00T models run on Jetson Orin edge hardware — enables on-device robot AI
- **Cosmos license:** Open-weights with NVIDIA Open Model License — commercial use permitted with restrictions on competing WFM development
- **BioNeMo Framework:** Separate from the API — open-source training framework at `github.com/NVIDIA/BioNeMo`
- **PhysicsNeMo:** Neural operator framework for engineering simulations (CFD, heat transfer) — distinct from Earth-2's climate models
- **Pricing:** NIM API is credits-based; BioNeMo cloud is usage-based per inference. Self-hosted is hardware-cost only.


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]

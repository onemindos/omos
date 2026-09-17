# Local Inference & World Models — OneMind OS
**Filed:** 2026-06-11 · **For:** Zeus · **By:** Legacy  
**Hardware baseline:** MBP M2 Max · 32 GB unified memory · 38 GPU cores · Metal 4 · macOS 26.5

---

## TL;DR

| Question | Answer |
|----------|--------|
| Best local LLM stack for your Mac? | **MLX** (fastest) + **Ollama** (daily driver, NATS-ready) |
| vLLM on your Mac? | Not recommended — Linux/NVIDIA first, experimental on Apple Silicon |
| LM Studio? | Good GUI browser/tester; already installed |
| Biggest model you can run well? | **Qwen 2.5 32B Q4_K_M** or **Gemma 3 27B Q4_K_M** (~18-19 GB) |
| Can you run world models now? | **Yes** — WorldForge mock today; LeWM this week; V-JEPA 2-AC next sprint |
| World model harness? | **WorldForge** (`uv add worldforge-ai`) — the canonical harness |

---

## Part 1 — Your Hardware, Translated for Inference

```
Apple M2 Max · 32 GB Unified Memory
┌────────────────────────────────────────────────────────────┐
│  CPU: 12 cores (8 perf + 4 eff)                           │
│  GPU: 38 cores · Metal 4                                   │
│  Memory: 32 GB shared (CPU + GPU — same pool)              │
│  Memory bandwidth: ~400 GB/s                               │
│  Available for models: ~26-28 GB (after OS + apps)         │
└────────────────────────────────────────────────────────────┘
```

**The key insight:** Unified memory means your GPU "VRAM" IS your RAM. 32 GB is available to the model. No PCIe bottleneck. This makes Apple Silicon unusually capable for local inference despite "no discrete GPU."

**Bottleneck for LLMs:** Memory bandwidth (400 GB/s), not compute. Token generation speed ∝ bandwidth ÷ model bytes loaded per token.

---

## Part 2 — Serving Stack Comparison

| Tool | Apple Silicon | Speed on M2 Max | API | Status | Best For |
|------|--------------|----------------|-----|--------|---------|
| **MLX** (Apple) | ⭐ Native Metal | Fastest (80-90% of theoretical BW) | OpenAI-compat | Not installed | Speed-critical, research |
| **Ollama** | ✅ llama.cpp/Metal | Good (65-75% of BW) | OpenAI-compat REST | **Installed** (no models) | Daily driver, NATS routing |
| **LM Studio** | ✅ llama.cpp/Metal | Good (65-75% of BW) | OpenAI-compat | **Installed** (no models) | GUI browsing/testing |
| **llama.cpp** (direct) | ✅ Metal native | Good, max control | CLI + server | Not installed | Fine-tuned control |
| **vLLM** | ⚠️ Experimental | Broken/limited on Mac | OpenAI-compat | Not installed | Linux+NVIDIA only |
| **Jan.ai** | ✅ llama.cpp | Same as LM Studio | OpenAI-compat | Not installed | Optional GUI alt |

### Winner for OneMind OS: **MLX + Ollama**

```
MLX  → fastest raw speed, use when you need peak tokens/s
Ollama → already installed, OpenAI API, NATS-routeable, best for
         Hermes → route LLM calls to local model via base_url override
```

**vLLM verdict:** Skip it on Mac. vLLM's PagedAttention KV cache is CUDA-first — the Apple Silicon backend is community-maintained, often broken, and gives you none of the production advantages. vLLM belongs on a Linux GPU node (future fleet).

---

## Part 3 — Full Model Table (M2 Max 32 GB)

> **Key:** ✅ runs well | ⚠️ tight/slow | ❌ won't fit | ⭐ recommended

### GGUF Quantization Size Reference
`Q4_K_M ≈ 4.5 bits/param ÷ 8 = 0.5625 bytes/param`  
Token speed formula: `≈ 400 GB/s ÷ model_bytes`

| Model | Params | Quant | VRAM Est. | Tokens/s Est. | Fits? | Recommend? | Best Use |
|-------|--------|-------|-----------|--------------|-------|-----------|---------|
| Gemma 3 2B | 2B | Q4_K_M | ~1.4 GB | 120–150 | ✅ | Fast tasks | Edge/tool routing |
| Llama 3.2 3B | 3B | Q4_K_M | ~1.9 GB | 100–130 | ✅ | Background agents |
| Phi-4 Mini | 3.8B | Q4_K_M | ~2.4 GB | 85–110 | ✅ | Coding assistant |
| Llama 3.1 8B | 8B | Q4_K_M | ~4.9 GB | 60–80 | ✅ | ⭐ General daily use |
| Gemma 3 9B | 9B | Q4_K_M | ~5.6 GB | 55–72 | ✅ | ⭐ General + vision |
| DeepSeek R1 8B | 8B | Q4_K_M | ~4.9 GB | 60–80 | ✅ | ⭐ Reasoning chain |
| Llama 3.3 70B | 70B | Q2_K | ~22 GB | 12–18 | ⚠️ | Avoid (quality loss) | — |
| Qwen 2.5 14B | 14B | Q4_K_M | ~8.7 GB | 35–50 | ✅ | ⭐ Strong reasoning |
| Phi-4 14.7B | 14.7B | Q4_K_M | ~9.2 GB | 32–46 | ✅ | ⭐ Coding/STEM |
| Mistral Small 22B | 22B | Q4_K_M | ~13.7 GB | 22–32 | ✅ | ⭐ Balance speed/quality |
| Gemma 3 27B | 27B | Q4_K_M | ~16.9 GB | 18–26 | ✅ | ⭐ Multimodal/vision |
| **Qwen 2.5 32B** | 32B | Q4_K_M | ~19.6 GB | 15–22 | ✅ | ⭐⭐ **Best quality fits** |
| Mistral Large 2411 | 123B | — | ~75 GB | — | ❌ | Too large | — |
| Llama 3.3 70B | 70B | Q4_K_M | ~43 GB | — | ❌ | Too large | — |
| Qwen 2.5 72B | 72B | Q4_K_M | ~45 GB | — | ❌ | Too large | — |
| DeepSeek V3 671B | 671B | — | >>32 GB | — | ❌ | Cloud only | — |

### Recommended Load-out for OneMind OS
```
Slot 1 (fast agent) :  Llama 3.1 8B Q4_K_M   → 4.9 GB  → Hermes tool router
Slot 2 (reasoning)  :  Qwen 2.5 14B Q4_K_M   → 8.7 GB  → Planning, analysis
Slot 3 (max quality):  Qwen 2.5 32B Q4_K_M   → 19.6 GB → Best offline frontier
```
All three together = 33.2 GB — only run one at a time in Ollama unless you want to swap.

### MLX-Specific Models (faster than Ollama/llama.cpp on M2)
```bash
mlx_lm.generate --model mlx-community/Qwen2.5-32B-Instruct-4bit --prompt "..."
mlx_lm.generate --model mlx-community/Llama-3.1-8B-Instruct-4bit --prompt "..."
mlx_lm.server --model mlx-community/Qwen2.5-14B-Instruct-4bit  # OpenAI API on :8080
```
MLX community hub: `huggingface.co/mlx-community` — pre-quantized versions of all major models.

---

## Part 4 — Quick Start Commands

### Install MLX (takes 2 min)
```bash
pip3 install mlx-lm
# First model pull
mlx_lm.generate --model mlx-community/Llama-3.1-8B-Instruct-4bit \
  --prompt "Hello, test inference"
```

### Load Models into Ollama (already installed)
```bash
ollama pull qwen2.5:14b          # 8.7 GB
ollama pull llama3.1:8b           # 4.9 GB
ollama pull deepseek-r1:8b        # 4.9 GB
ollama pull qwen2.5:32b           # 19.6 GB — max quality

# Serve with OpenAI-compat API (for Hermes routing)
ollama serve  # runs on localhost:11434
```

### Point Hermes at Ollama
```bash
hermes config set provider.base_url http://localhost:11434/v1
hermes config set provider.model qwen2.5:14b
```

---

## Part 5 — World Models: Can We Run Them NOW?

**Short answer: YES for the right class. Here's the breakdown:**

### World Model Readiness on M2 Max 32 GB

| Model | Params | Size | Runs on Mac? | How | Effort |
|-------|--------|------|-------------|-----|--------|
| **WorldForge mock** | — | 0 | ✅ TODAY | `uv add worldforge-ai` | 5 min |
| **LeWM** | 15M | ~34 MB | ✅ This week | PyTorch CPU/Metal | 2-4 hrs setup |
| **V-JEPA 2** | 1.2B | ~2.4 GB | ✅ This week | PyTorch or MLX | 2-4 hrs |
| **V-JEPA 2-AC** | 1.2B (ft) | ~2.4 GB | ✅ Sprint 1 | HuggingFace + Metal | 4-8 hrs |
| **DreamerV4** | 200M+ | ~400 MB+ | ⚠️ Partial | PyTorch CPU only | Research mode |
| **NVIDIA Cosmos** | 4B–14B | 8–28 GB | ❌ No | NVIDIA-only runtime | Skip for Mac |
| **Genie 3** | Not released | — | ❌ Not available | — | Future |
| **I-JEPA** | ViT-H/14 | ~630 MB | ✅ | PyTorch + Metal | 1 hr |
| **Nexar BADAS 2.0** | VLM-class | — | ⚠️ Limited | Inference only | Research only |

### The Harness: WorldForge

WorldForge (`github.com/AbdelStark/worldforge`) is the canonical harness. It separates:
```
POLICY  provider  →  proposes candidate actions  (GR00T, LeRobot, LLM goal-setter)
WM      provider  →  scores/rolls out futures     (V-JEPA 2-AC, LeWM, DreamerV4)
MPC     optimizer →  picks lowest-cost action     (CEM / receding-horizon)
ROBOT   executor  →  runs the action              (DIMOS, OM1, Go2, G1)
```

5 capabilities: `predict`, `score`, `policy`, `embed`, `plan`

### Phase 0 — Right Now (0 GPU required)
```bash
uv add worldforge-ai
uv run worldforge doctor --registered-only
uv run python examples/latent_mpc_planning.py   # mock provider, zero GPU
```
This teaches the loop without touching torch or downloading checkpoints.

### Phase 1 — LeWM on Mac (This Week)
```bash
pip3 install torch torchvision  # ~2 GB, CPU+Metal
# Clone and install LeWM
git clone https://github.com/YOUR_LEWM_REPO
# 15M params → 34 MB, trains in hours on M2 Max Metal
# Plug in as `score` provider in worldforge
```
**Why start with LeWM:** 15M params. Fits on literally any machine. Gives you the surprise/anomaly primitive immediately. Proof of concept for the whole loop.

### Phase 2 — V-JEPA 2-AC (Sprint 1)
```bash
# Download action-conditioned checkpoint (~2.4 GB)
huggingface-cli download facebook/vjepa2-ac-vitl-256 --local-dir ./vjepa2-ac
# Wire as WM provider in worldforge config
# This is the checkpoint you'll eventually run on DIMOS/OM1/Go2
```
**Why V-JEPA 2-AC matters for OneMind:** Trained on 62h robot trajectories. Zero-shot on Franka arm via goal-image planning. Open weights, commercial-OK. This is your physical layer for DIMOS, OM1, eventually G1.

---

## Part 6 — Vault Intel on World Models

What's already filed in the vault:

### Service Registry Cards (Fabric — 6 filed, 2026-06-08)
`SO/01 The Fabric/nodes/digital/services-registry/ai/world-models/`
- `i-jepa.md` — Meta Jan 2023, foundational JEPA proof
- `v-jepa.md` — Meta Jan 2024, video extension
- `v-jepa-2.md` — Meta Jun 2025, 1.2B, SOTA physical reasoning, commercial-OK
- `v-jepa-2-ac.md` — Action-conditioned checkpoint, ~62h Franka data, zero-shot planning
- `lewm.md` — 15M params, single GPU, 48× faster than DINO-WM, SIGReg
- `nexar-badas-2.md` — VLM driving WM, SOTA vs GPT-4o + Gemini 2.0 Flash, fleet dashcam

### Architecture Doc
`SO/08 Research (The Lab)/robotics/world-models-vs-llms-onemind-architecture-2026.md`
- Full deep brief filed 2026-06-07
- Covers: WM vs LLM roles, the canonical control loop, OneMind stack placement, always-on pattern

### Unitree/Robot Integration
`SO/08 Research (The Lab)/robotics/unitree/unitree-openmind-integration.md`
- References shared world model across fleet (line 41)
- FABRIC network WM node architecture (line 392)

### OneMind Stack Placement (from vault)
```
L5 PRESENCE      GPT-realtime + wearables
L4 ORCHESTRATOR  Hermes/Legacy — memory, tools, routing         ← LLM layer
L3 SPECIALISTS   domain agents (infra, security, comms, biz)
─────────────────────────────────────────────────────────────
L2 WM CONTROL    per-robot world models, latent MPC             ← NEW (V-JEPA 2-AC, LeWM)
L1 ALWAYS-ON     local WM perception, 24/7, token-free          ← NEW (surprise → interrupt L4)
L0 FABRIC        robots/drones/vehicles/cameras — via NATS
```

---

## Part 7 — Recommended Action Plan

### This Session
```bash
# 1. Install MLX (fastest for M2 Max)
pip3 install mlx-lm

# 2. Pull your first model via Ollama
ollama pull llama3.1:8b
ollama pull qwen2.5:14b

# 3. Run WorldForge Phase 0
pip3 install uv
uv add worldforge-ai
uv run worldforge doctor --registered-only
```

### This Week
- Pull Qwen 2.5 32B Q4_K_M (max quality local model)
- Install LeWM → wire as score provider in worldforge
- Test the propose/score/execute loop with mock policy + LeWM score

### Sprint 1
- Download V-JEPA 2-AC checkpoint
- Wire to worldforge as WM provider
- Test on DIMOS sensor stream → surprise signal → NATS interrupt

### Later (when Linux GPU node exists)
- Deploy vLLM on a Linux VPS/fleet node for production batching
- Cosmos + GR00T belong there, not on the Mac

---

## Summary

| Verdict | Detail |
|---------|--------|
| **Best stack** | MLX (fastest) + Ollama (NATS/Hermes integration) |
| **Skip vLLM on Mac** | Linux/NVIDIA only — add to VPS node later |
| **Max quality model** | Qwen 2.5 32B Q4_K_M · ~19.6 GB · 15-22 t/s |
| **Sweet spot model** | Qwen 2.5 14B Q4_K_M · ~8.7 GB · 35-50 t/s |
| **World models** | YES — WorldForge today, LeWM this week, V-JEPA 2-AC next sprint |
| **Cosmos/Genie** | No — NVIDIA-only or unreleased, skip for Mac |
| **Your WM stack** | WorldForge harness · LeWM (score) · V-JEPA 2-AC (robot WM) |
| **NATS integration** | WM nodes → NATS subjects like `wm.dimos.surprise` → interrupt L4 |

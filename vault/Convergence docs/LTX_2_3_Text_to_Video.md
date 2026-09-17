# LTX 2.3 — Open-Source Text-to-Video (T2V) AI Model with Flexible VRAM Requirements

**Official Links:**
- **Hugging Face Model:** https://huggingface.co/Lightricks/LTX-2.3
- **GitHub Repository:** https://github.com/Lightricks/LTX-Video
- **Civitai Workflow (Low-VRAM):** https://civitai.com/articles/30959/ltx-23-t2v-simple-workflow-for-low-vram-gpus-8gb-or-10gb-vram
- **Paper (arxiv):** https://arxiv.org/abs/2601.03233
- **LTX.io Official:** https://ltx.io/

**Type:** Open-source DiT-based audio-video foundation model

**Release:** 2026 (recent advancement to LTX-2)

**Model Size:** 22B parameters

**License:** LTX-2 Community License Agreement (free for non-commercial + some commercial use)

---

## What Is LTX 2.3?

**LTX 2.3** is an **open-source text-to-video (T2V) generative AI model** that:
- ✅ Generates synchronized video + audio in a SINGLE model
- ✅ Works on consumer GPUs (6GB to 32GB+ VRAM)
- ✅ Scalable from low-end laptops to high-end workstations
- ✅ Open weights (Lightricks, permissive license)
- ✅ Optimizable via GGUF quantization for local deployment
- ✅ Integrates with ComfyUI (open-source video generation platform)

**Core Innovation:** Joint audio-video generation (not separate models)

---

## VRAM Requirements by Use Case

### Official Desktop Minimum

**32GB VRAM:** RTX 5090 or equivalent
- Full model operation (bf16 precision)
- No optimization/quantization needed
- Best quality output
- Fastest generation speed

### Realistic Deployment Scenarios

#### 1. **6GB VRAM (Budget Laptop)**
**Hardware:** RTX 3060 (6GB GDDR6)

**Configuration:**
- GGUF quantization (Q4 or Q5)
- System RAM offloading (requires 16-32GB system RAM)
- ComfyUI with memory optimization
- Batch size: 1 only
- Inference: 1080p 5-10 second clips

**Performance:**
- Generation time: 10-15 minutes per 5-second clip
- Quality: Good (quantization tradeoff)
- Practical for: Experimenting, learning, infrequent generation

**Setup:** Requires careful tuning of ComfyUI memory settings

#### 2. **8-10GB VRAM (Entry-Level Desktop)**
**Hardware:** RTX 4070 (12GB), RTX 3060 Ti (8GB), RTX 4060 Ti (8GB)

**Configuration:**
- GGUF quantization (Q4/Q5) OR light precision reduction (fp8)
- System RAM offloading (16GB+ system RAM recommended)
- ComfyUI with standard optimizations
- Batch size: 1
- Inference: 1080p 5-15 second clips

**Performance:**
- Generation time: 5-10 minutes per 10-second clip
- Quality: Good-to-excellent
- Practical for: Daily video generation, lightweight workflows

**Tested & Confirmed:** Civitai article confirmed 8-10GB hardware works reliably

#### 3. **12-16GB VRAM (Mid-Range GPU)**
**Hardware:** RTX 4070 Ti (12GB), RTX 3080 (10GB), RTX 4070 (12GB)

**Configuration:**
- fp8 precision (mixed-precision)
- Minimal offloading needed
- ComfyUI standard settings
- Batch size: 1-2
- Inference: 1080p 10-20 second clips

**Performance:**
- Generation time: 3-8 minutes per 10-second clip
- Quality: Excellent
- Practical for: Professional daily workflow, multiple generation passes

**Status:** Recommended "daily driver" VRAM tier

#### 4. **24GB VRAM (High-End Consumer)**
**Hardware:** RTX 4090 (24GB), RTX 6000 (48GB), A6000 (48GB)

**Configuration:**
- bf16 precision (full quality)
- No offloading needed
- Standard ComfyUI settings
- Batch size: 2-4 possible
- Inference: 1080p 12-20 seconds, short 4K via tiling

**Performance:**
- Generation time: 2-5 minutes per 10-second clip
- Quality: Maximum (full model precision)
- Practical for: Professional production, long-form content, 4K experiments

**Extra Capability:** Tiling-based upscaling for higher resolution

#### 5. **32GB+ VRAM (Data Center / Professional)**
**Hardware:** RTX 5090 (32GB), A100 (40GB), H100 (80GB)

**Configuration:**
- bf16 precision (full model)
- No optimizations needed
- High batch sizes (4-8)
- Inference: 1080p 20-30 seconds, 4K 10-15 seconds

**Performance:**
- Generation time: 1-3 minutes per 10-second 1080p
- Quality: Maximum
- Practical for: Batch production, research, commercial deployment

---

## Model Variants

### Main Checkpoints (Hugging Face)

| Name | Use Case | VRAM | Quality | Speed |
|------|----------|------|---------|-------|
| **ltx-2.3-22b-dev** | Full model, flexible, trainable | 32GB (bf16) | Maximum | Slowest |
| **ltx-2.3-22b-distilled** | Distilled version, 8 steps, CFG=1 | 16GB (fp8) | High | Medium |
| **ltx-2.3-22b-distilled-1.1** | Distilled v1.1, improved audio | 16GB (fp8) | High | Medium |
| **ltx-2.3-22b-distilled-lora-384** | LoRA version (fine-tune friendly) | 12GB (q8) | High | Medium |
| **ltx-2.3-22b-distilled-lora-384-1.1** | LoRA v1.1 | 12GB (q8) | High | Medium |

### Upscaler Models

| Name | Purpose | Application |
|------|---------|------------|
| **ltx-2.3-spatial-upscaler-x2-1.1** | 2x resolution upscaler | Multi-stage for 4K |
| **ltx-2.3-spatial-upscaler-x1.5-1.0** | 1.5x resolution upscaler | Intermediate resolution |
| **ltx-2.3-temporal-upscaler-x2-1.0** | 2x temporal upscaler (FPS) | Smoother motion |

---

## How LTX 2.3 Works

### Architecture

**Type:** DiT (Diffusion Transformer)

**Key Features:**
- Diffusion-based generation (iterative refinement)
- Transformer backbone (attention-based)
- Joint audio-video processing
- Unified latent space (audio + video share representation)

### Generation Pipeline

**Step 1: Text Encoding**
- Input text → CLIP text encoder → semantic embedding
- Dual CLIP models (Gemma-3, LTX text projection)
- Embedding contains prompt semantics

**Step 2: Diffusion Process**
- Start with noise (latent space)
- Iteratively denoise guided by text embedding
- Typically 8-50 steps (configurable)
- Distilled models: 8 steps (faster)
- Full models: 50+ steps (higher quality)

**Step 3: Audio Synthesis**
- Generate audio latents simultaneously
- Text embedding guides audio characteristics
- Synchronized with video timing

**Step 4: Decoding**
- Latent video → VAE decoder → 1080p video
- Latent audio → VAE decoder → audio waveform
- Output: MP4 with embedded audio

### ComfyUI Integration

**Workflow:**
```
Text Prompt
  ↓
LTX Loader (GGUF quantized model)
  ↓
CLIP Text Encode (dual CLIP models)
  ↓
LTX Sampler (diffusion iterations)
  ↓
VAE Decode (latent → video)
  ↓
Audio Processor (extract/sync audio)
  ↓
Output: MP4 video + audio
```

---

## Quantization Options for Low-VRAM

### GGUF Quantization (Recommended for Low-VRAM)

**What is GGUF?**
- Quantization format for LLMs + generative models
- Reduces model size (32GB bf16 → 8GB Q4)
- Minimal quality loss with Q4/Q5
- Supported by ComfyUI + LM Studio + Ollama

**Quantization Levels:**

| Level | Size Reduction | Quality Loss | VRAM Requirement | Speed |
|-------|---|---|---|---|
| **Q2_K** | 90% | High | 2-3GB | Fastest |
| **Q3_K** | 85% | Medium | 3-4GB | Fast |
| **Q4_K_M** | 75% | Minimal | 6-8GB | Medium |
| **Q5_K_M** | 65% | Very minimal | 10-12GB | Slower |
| **Q6_K** | 55% | Negligible | 14-16GB | Slowest |
| **bf16** | None (baseline) | None | 32GB | Baseline |

**Recommended for Low-VRAM:** Q4_K_M (best quality/VRAM tradeoff)

**Civitai Article Resources:**
- Unet (GGUF): https://huggingface.co/unsloth/LTX-2.3-GGUF/blob/main/ltx-2.3-22b-dev-Q4_K_M.gguf
- Video VAE: https://huggingface.co/Kijai/LTX2.3_comfy/blob/main/vae/LTX23_video_vae_bf16.safetensors
- Audio VAE: https://huggingface.co/Kijai/LTX2.3_comfy/blob/main/vae/LTX23_audio_vae_bf16.safetensors
- Dual CLIP Models: https://huggingface.co/Comfy-Org/ltx-2/blob/main/split_files/text_encoders/

### Other Optimization Techniques

**1. Flash Attention**
- Reduces memory during transformer operations
- 10-20% VRAM reduction with minimal quality loss
- Requires CUDA compute capability 8.0+

**2. Mixed Precision (fp8)**
- Combines fp32 + fp8 selectively
- 30-40% memory reduction
- Minimal quality impact
- Faster computation on compatible hardware

**3. Gradient Checkpointing**
- Saves activation memory during generation
- Trade: Slower generation (extra compute)
- Useful for training/fine-tuning on limited VRAM

**4. Model Offloading**
- Move model layers to system RAM when not in use
- Trade: Slower (PCIe bandwidth limited)
- Enables generation on smaller GPUs
- Tested: 6GB GPU + 32GB system RAM viable

---

## Generation Parameters

### Core Prompts

**Text Prompt (required):**
- Natural language description of desired video
- Example: "A futuristic city at night with flying cars and neon lights"
- Length: 10-150 words (optimal)
- More detailed = better adherence but slower generation

**Negative Prompt (optional):**
- What NOT to generate
- Example: "blurry, low quality, distorted"

### Technical Parameters

| Parameter | Range | Default | Effect |
|-----------|-------|---------|--------|
| **Steps** | 8-50 | 8 (distilled) / 50 (full) | More steps = higher quality but slower |
| **CFG Scale** | 1.0-20.0 | 7.5 | Higher = stronger prompt adherence |
| **Seed** | 0-2^32 | Random | Fixed seed = reproducible output |
| **Resolution** | 256-2048 | 1024 | 1080p = 1024x576 or 1024x1024 |
| **Duration** | 5-120 frames | 61 (2.5s @ 24fps) | Longer = more VRAM required |
| **FPS** | 8-60 | 24 | Frame rate of output |

---

## Real-World Performance (Civitai Test)

**Hardware:** 10GB VRAM (RTX 3060/4070 Ti) + 32GB system RAM

**Test Conditions:**
- Windows 11 dual-monitor setup
- Other apps running (Chrome, video player)
- ComfyUI VRAM fluctuation: 8-10GB

**Results:**
- **Output:** 10-second 1080p video
- **Generation Time:** ~7 minutes
- **Quality:** Good (no visible degradation from quantization)
- **Audio:** Synchronized, natural

**Key Insight:** Real-world performance depends heavily on:
- Quantization (Q4 vs Q5 = 2-3 min difference)
- System load (background apps impact)
- GPU architecture (newer = more efficient)

---

## Practical Deployment Guide

### Setup 1: Low-VRAM Laptop (6GB GPU)

**Hardware:**
- RTX 3060 or RTX 4060 (6GB)
- 16GB+ system RAM
- SSD for model caching

**Installation:**
```bash
# Install ComfyUI
git clone https://github.com/comfyanonymous/ComfyUI.git
cd ComfyUI
python -m venv venv
source venv/bin/activate  # Linux/Mac
venv\Scripts\activate     # Windows

# Install dependencies
pip install -r requirements.txt

# Download LTX 2.3 GGUF quantized model
# From Civitai or HuggingFace (Q4_K_M recommended)

# Configure ComfyUI for low-memory mode:
# Edit launch_colab.py or use --lowvram flag
```

**Workflow:**
- Load Q4_K_M GGUF model
- Use system RAM offloading
- Generate 5-10 second clips
- Expect 10-15 minutes per clip

### Setup 2: Mid-Range Desktop (12-16GB GPU)

**Hardware:**
- RTX 4070 / RTX 3080 (12-16GB)
- 16GB+ system RAM
- NVMe SSD

**Installation:**
Same as above, but:
- Can use distilled model directly (no quantization)
- Use fp8 precision mixed mode
- Skip aggressive memory optimizations

**Workflow:**
- Load ltx-2.3-22b-distilled model
- Generate 10-20 second clips
- Expect 3-8 minutes per clip

### Setup 3: Professional Workstation (24GB+ GPU)

**Hardware:**
- RTX 4090 / RTX 5090 (24GB+)
- 32GB+ system RAM
- Fast SSD/NVMe

**Installation:**
- Use full model (bf16 precision)
- Enable batching (generate multiple videos)
- Optional: Enable upscaling models

**Workflow:**
- Load full ltx-2.3-22b-dev model
- Generate multiple clips in parallel
- Expect 2-5 minutes per clip
- Can attempt 4K with upscaling

---

## Integration with OneMind/Convergence

### Creative Content Generation at Scale

**Pattern:**
```
OneMind Agent (Eve/Claude Code)
  ↓ (text prompts)
LTX 2.3 Generator (ComfyUI)
  ↓ (video + audio)
Content Archive (ClickHouse + S3)
  ↓ (indexing + retrieval)
Operator Interface (TAK Server, Unreal Engine)
  ↓ (display)
Tactical/Training Deployment
```

### Applications

1. **Dojo Training Content**
   - Generate tactical scenario videos automatically
   - Text-to-scenario: "Urban compound assault, daylight, three entrances"
   - Quick iteration on training content
   - Scale training across multiple scenarios

2. **Automated Intelligence Briefing Videos**
   - Text summary → AI-generated explainer video
   - Geographic narration + synchronized visuals
   - Faster intelligence dissemination

3. **Simulation Content**
   - Generate diverse scenario variations
   - Fill gaps in training dataset
   - Cost-effective vs. filming

4. **Homestead Documentation**
   - Auto-generate instructional videos
   - "How to maintain solar panel array"
   - Accessible locally (no internet needed)

---

## Comparison to Commercial Alternatives

| Feature | LTX 2.3 | Runway ML | Sora (OpenAI) | Pika |
|---------|---------|-----------|---|---|
| **Open Source** | ✅ Yes | ❌ No | ❌ No | ❌ No |
| **Local Deployment** | ✅ Yes (ComfyUI) | ❌ Cloud only | ❌ Cloud only | ❌ Cloud only |
| **Cost** | Free | $10-50/mo | $20/mo (waitlist) | Free-$10/mo |
| **Audio Generation** | ✅ Native | ❌ Separate | ❌ Separate | ❌ Separate |
| **VRAM Flexibility** | ✅ 6GB-32GB | N/A | N/A | N/A |
| **Quantization Support** | ✅ GGUF | N/A | N/A | N/A |
| **Quality (1080p 10s)** | High | Excellent | Excellent | Very Good |
| **Speed** | 2-10 min | Real-time | Unknown | Real-time |
| **Training Data Restrictions** | Yes (LTX license) | Yes | Yes | Yes |

**LTX 2.3 Advantages:**
- Open source (can modify/improve)
- Runs locally (privacy, offline capability)
- Scalable VRAM (works on laptops)
- Free forever (no subscription)
- Audio included (single model)

---

## Future Directions & Roadmap

### Potential Improvements (2026+)

1. **Longer Duration**
   - Currently: 5-10 seconds typical
   - Goal: 30-60 second videos feasible

2. **Higher Resolution**
   - Currently: 1080p comfortable
   - Goal: 4K native (not just upscaling)

3. **Better Audio Quality**
   - Current: Acceptable, synchronized
   - Goal: Studio-quality audio generation

4. **Faster Generation**
   - Current: 2-10 minutes
   - Goal: 30 seconds (real-time or near-real-time)

5. **LoRA Fine-tuning**
   - Enable user-specific style training
   - Domain-specific models (tactical, medical, etc.)

### Research Directions

- Multimodal guidance (image + text + audio prompts)
- Longer temporal coherence (maintain consistency across scenes)
- Physical accuracy (realistic physics simulation)
- Controllable generation (precise placement of objects)

---

## Limitations

### Quality Constraints

1. **Temporal Coherence** — Motion can become jittery at edges (scene transitions)
2. **Object Consistency** — Objects may change appearance across clip
3. **Prompt Adherence** — Complex prompts sometimes misinterpreted
4. **Realistic Physics** — May not follow realistic physics in all cases

### Performance Constraints

1. **VRAM** — Still requires significant GPU memory (6GB minimum)
2. **Speed** — Minutes per clip, not real-time
3. **Resolution** — 1080p typical, 4K requires upscaling
4. **Duration** — 5-20 seconds practical, longer clips require tiling

### Data Constraints

1. **Training Data** — Trained on public internet video (bias toward popular content)
2. **Watermarks** — May occasionally reproduce watermarks from training data
3. **Copyright** — Model can inadvertently generate copyrighted content

---

## Resource Links

### Official

- **HuggingFace Model:** https://huggingface.co/Lightricks/LTX-2.3
- **GitHub Repository:** https://github.com/Lightricks/LTX-Video (10.7k stars)
- **Research Paper:** https://arxiv.org/abs/2601.03233
- **Official Site:** https://ltx.io/

### Community Resources

- **Civitai Workflows:** https://civitai.com/articles/30959/ltx-23-t2v-simple-workflow-for-low-vram-gpus-8gb-or-10gb-vram
- **ComfyUI GitHub:** https://github.com/comfyanonymous/ComfyUI
- **Kijai LTX2.3 ComfyUI Support:** https://huggingface.co/Kijai/LTX2.3_comfy
- **GGUF Quantization Guide:** https://huggingface.co/unsloth/LTX-2.3-GGUF

### Model Checkpoints

- **Unsloth GGUF Variants:** https://huggingface.co/unsloth/LTX-2.3-GGUF
- **Kijai VAE/Encoders:** https://huggingface.co/Kijai/LTX2.3_comfy
- **Comfy-Org Split Files:** https://huggingface.co/Comfy-Org/ltx-2

---

## Summary

**LTX 2.3 = Open-source text-to-video generation for everyone**

✅ **Key Strengths:**
- Open weights (permissive license)
- Flexible VRAM (6GB to 32GB+)
- Audio + video in one model
- Local deployment (privacy, offline)
- Active community support

✅ **Practical Tiers:**
- 6GB laptop: Learning, experiments
- 12-16GB desktop: Daily professional use
- 24GB+ workstation: Production batch generation

✅ **Integration Ready:**
- ComfyUI (standard tool)
- HuggingFace ecosystem
- GGUF quantization (proven method)

**For OneMind Convergence:** Video generation layer for training content, tactical briefing automation, scenario diversity for Dojo training

---

**Captured:** 2026-07-11

**Primary Sources:**
- HuggingFace Model Card: https://huggingface.co/Lightricks/LTX-2.3
- Civitai Low-VRAM Workflow: https://civitai.com/articles/30959/ltx-23-t2v-simple-workflow-for-low-vram-gpus-8gb-or-10gb-vram
- User-provided VRAM specifications

**Status:** Production-ready, active development (2026)

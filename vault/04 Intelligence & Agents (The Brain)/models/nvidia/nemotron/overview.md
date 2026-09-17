# Nemotron

NVIDIA's Nemotron family — foundation and reward models for alignment and agent use.

| Model | Type | Params | Active | Context | Notes |
|-------|------|--------|--------|---------|-------|
| `nvidia/NVIDIA-Nemotron-3-Ultra-550B-A55B` | LLM MoE | 550B | 55B | 1M | **Trainable base model** — best open foundation for custom training; weights H1 2026 |
| `nvidia/NVIDIA-Nemotron-3-Super-120B-A12B-BF16` | LLM MoE | 120B | 12B | 1M | Released Mar 11 2026 — 2.2× throughput vs GPT-OSS-120B |
| `nvidia/NVIDIA-Nemotron-3-Nano-30B-A3B-BF16` | LLM MoE | 30B | 3B | — | Nano — lightweight MoE |
| `nvidia/NVIDIA-Nemotron-3-Nano-4B-BF16` | LLM | 4B | 4B | — | Edge-deployable |
| `nvidia/llama-3.1-nemotron-ultra-253b-v1` | LLM | 253B | 253B | — | Previous gen Nemotron Ultra (Llama 3.1-based) |
| `nvidia/nemotron-4-340b-instruct` | LLM | 340B | 340B | — | Nemotron-4 instruction-tuned |
| `nvidia/nemotron-4-340b-reward` | Reward Model | 340B | 340B | — | RLHF scoring — used for training pipelines |

## Nemotron 3 Ultra — Key Details

Announced at **GTC San Jose 2026** · Best Open Base Model (announced)

| Spec | Value |
|------|-------|
| Total params | 550B |
| Active params / token | up to 55B |
| Architecture | Hybrid Mamba-Transformer MoE |
| Context window | 1M tokens |
| Pre-training precision | NVFP4 |
| Status | **Base only / Trainable** — weights releasing H1 2026; designed for fine-tuning, not out-of-box use |

**Architecture innovations (shared with Nemotron 3 Super):**
- **LatentMoE** — compresses tokens into low-rank latent space before routing; activates 4× more expert specialists at same inference cost
- **Multi-Token Prediction (MTP)** — predicts multiple future tokens per forward pass; enables built-in speculative decoding
- **Mamba-2 layers** — linear-time complexity over sequence length; makes 1M-token context practical

**Benchmarks vs GLM-4.5-355B-A32B and Kimi-K2-1026B-A33B (on GB200 NVL72):**

| Benchmark | Nemotron 3 Ultra | GLM-4.5-355B | Kimi-K2 |
|-----------|-----------------|-------------|---------|
| MMLU Pro | **79.0** | 65.6 | 69.3 |
| MMLU | **89.1** | 86.3 | 88.0 |
| Code | **85.3** | 76.2 | 75.3 |
| Math | **85.4** | 72.1 | 79.5 |
| Multilingual | **89.0** | 83.3 | 84.2 |
| Peak Throughput | **5× vs GLM** | 1× | ~2.5× |

> 🏋️ **This is a trainable base model.** Nemotron 3 Ultra is explicitly designed as the best open starting point for custom training — domain-specific fine-tuning, RL post-training, and custom instruction tuning pipelines. It is NOT meant to be used out of the box as an assistant. If you need a deployable model now, use Nemotron 3 Super.

## API Access

All Nemotron models are available via NVIDIA NIM (NVIDIA Inference Microservice):

```
Base: https://integrate.api.nvidia.com/v1
Key: NVIDIA_API_KEY (build.nvidia.com)
```

## Stack Use Cases

- **Nemotron 3 Ultra (when released):** best open base for domain fine-tuning, RL post-training, custom instruction pipelines
- **Nemotron 3 Super:** production agentic workloads — high throughput, 1M context, open weights available now
- **Nemotron-4-340b-reward:** synthetic data generation and RLHF annotation pipelines
- **Nemotron-3-Nano:** edge inference where full GPU unavailable


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]

# GPU Analysis — Does Legacy Need GPU Power?

> Decision date: 2026-03-16
> Status: **NOT NEEDED NOW** — revisit when local inference or OM1 robotics become active

---

## Current Architecture (No GPU Required)

Legacy's entire stack is CPU/memory-bound:

| Component | Runtime | GPU needed? |
|-----------|---------|-------------|
| OpenClaw Gateway | Node.js 22 | No |
| 9 AI agents (Legacy + 8 domain) | Azure OpenAI API (GPT-5.2 / GPT-4o) | No — GPU is on Azure's side |
| Embeddings | OpenAI text-embedding-3-small (API) | No — remote |
| NATS JetStream | Native binary | No |
| PostgreSQL / TimescaleDB | Docker | No |
| legacy-bridge | Python | No |
| codex-watch / codex-pull | Git + systemd | No |

**All LLM inference runs on Azure/OpenAI GPUs via API.** The VPS only orchestrates, routes, and stores.

### Current Droplet

| Field | Value |
|-------|-------|
| Name | `onemind-main` |
| Spec | 4 vCPU / 16 GB RAM / 200 GB SSD (AMD) |
| Size slug | `s-4vcpu-16gb-amd` |
| Cost | $84/mo |
| Next scale target | 8 vCPU / 32 GB RAM (when agent load increases) |

---

## When GPU Would Be Justified

### Use Case 1: Local LLM Inference
- Run Llama 3.3 70B, Mistral, Qwen, etc. locally via Ollama / vLLM / TGI
- **Why**: Zero API costs at scale, no rate limits, full data privacy, lower latency
- **When**: If Azure OpenAI costs exceed ~$500/mo or privacy requirements change
- **Minimum GPU**: NVIDIA A100 40GB or H100 80GB

### Use Case 2: Local Embeddings
- Replace OpenAI text-embedding-3-small with local model (e.g., BGE, E5)
- **Why**: Faster re-indexing, no API dependency, unlimited throughput
- **When**: Codex vault grows very large or real-time semantic search becomes critical
- **Minimum GPU**: NVIDIA A10 or T4 (lighter workload)

### Use Case 3: Computer Vision / OM1 Robotics
- Run vision models locally for camera feeds, object detection, scene understanding
- **Why**: Real-time inference required for physical robots, can't tolerate API round-trips
- **When**: OM1 / LeRobot workloads go live
- **Minimum GPU**: NVIDIA T4 or L4

### Use Case 4: Fine-Tuning
- Fine-tune models on Zeus's data (codex, memory, communication patterns)
- **Why**: Personalized model behavior, domain-specific knowledge baked in
- **When**: After collecting enough training data and defining fine-tune objectives
- **Minimum GPU**: NVIDIA A100 80GB (or H100 for larger models)

---

## DigitalOcean GPU Droplets

### Availability
- DO offers **GPU Droplets** with NVIDIA H100 GPUs
- Region: Limited availability (check `doctl compute size list --output json | grep gpu`)
- Pre-configured with NVIDIA drivers + CUDA toolkit

### Pricing

| GPU | vCPU | RAM | Storage | Cost |
|-----|------|-----|---------|------|
| 1x H100 80GB | 20 | 240 GB | 450 GB NVMe | ~$2.50-3.50/hr ($1,800-2,500/mo) |

### How It Would Work

```
┌─────────────────────────────────┐    ┌──────────────────────────────────┐
│  onemind-main (current)         │    │  onemind-gpu (new, separate)     │
│  4 vCPU / 16 GB / $84/mo       │    │  H100 GPU / ~$2,500/mo          │
│                                 │    │                                  │
│  OpenClaw Gateway               │    │  Ollama / vLLM / TGI            │
│  9 Agents                       │    │  Local LLM serving              │
│  NATS, PostgreSQL               │◄──►│  Embeddings server              │
│  All orchestration              │    │  Vision models (OM1)            │
│                                 │    │                                  │
│  Calls GPU droplet via          │    │  Exposes API on Tailscale:      │
│  Tailscale private network      │    │  :11434 (Ollama)                │
│                                 │    │  :8000 (vLLM)                   │
└─────────────────────────────────┘    └──────────────────────────────────┘
```

**Key principle**: Keep GPU workloads on a **separate droplet**. Don't mix inference with orchestration. This lets you spin the GPU droplet up/down based on need.

### Setup Steps (When Ready)

```bash
# 1. Create GPU droplet
doctl compute droplet create onemind-gpu \
  --region nyc3 \
  --size gpu-h100x1-80gb \
  --image ubuntu-24-04-x64 \
  --vpc-uuid 2e617f7b-2721-4371-b34c-b4f5d91aec76 \
  --ssh-keys <key-id>

# 2. Install Tailscale (private network access)
curl -fsSL https://tailscale.com/install.sh | sh
tailscale up --authkey=<tskey>

# 3. Install Ollama
curl -fsSL https://ollama.com/install.sh | sh
ollama pull llama3.3:70b

# 4. Update OpenClaw config to route to local model
# In openclaw.json, change model provider to point to Tailscale IP:11434
```

---

## Cost Comparison

| Scenario | Monthly Cost | Model Quality |
|----------|-------------|---------------|
| **Current** (Azure OpenAI API) | $84 droplet + ~$50-200 API usage | GPT-5.2 (best) |
| **GPU Droplet** (always-on) | $84 + ~$2,500 GPU | Llama 3.3 70B (good, not GPT-5.2) |
| **GPU Droplet** (on-demand, 4h/day) | $84 + ~$300 GPU | Llama 3.3 70B (intermittent) |
| **Hybrid** (API + local for embeddings) | $84 + ~$150 small GPU + reduced API | Best of both |

---

## Recommendation

### Now (2026-03)
- **No GPU needed.** Stay on cloud LLMs via Azure OpenAI.
- Use DO credit toward the planned **8 vCPU / 32 GB RAM** upgrade for `onemind-main`.

### Future Triggers to Revisit
- [ ] Azure OpenAI monthly spend exceeds $500/mo consistently
- [ ] OM1 / LeRobot workloads go live and need real-time vision inference
- [ ] Privacy requirements mandate no external API calls for certain data
- [ ] Fine-tuning experiments begin
- [ ] DO credit approaching expiration with significant balance remaining

### If GPU Is Added Later
1. **Separate droplet** — never mix GPU inference with Legacy's orchestration runtime
2. **Tailscale networking** — private connectivity, no public exposure
3. **On-demand scheduling** — spin up only when needed to control costs
4. **NATS integration** — GPU droplet subscribes to inference request subjects, publishes results

---

## Related Docs

- [infra/README.md](infra/README.md) — Current VPS architecture
- [infra/VPS-MAIN.md](infra/VPS-MAIN.md) — Main droplet details
- [infra/deploy/DROPLET-RESIZE.md](infra/deploy/DROPLET-RESIZE.md) — Resize history and future scale target
- [infra/ops/SCALING.md](infra/ops/SCALING.md) — Scaling guide

# Open-Source Embedding Models — 2026 Reference

> **Context**: Google launched [Gemini Embedding 2](https://blog.google/innovation-and-ai/models-and-research/gemini-models/gemini-embedding-2/) as a proprietary, API-only embedding model. This document covers the open-source ecosystem that matches or exceeds it — models you can self-host, run locally, and integrate without per-token costs or vendor lock-in.

---

## What Is Gemini Embedding 2?

Gemini Embedding 2 is Google's latest text embedding model. Key specs:

- **3072-dimensional** output vectors
- **2048-token** context window
- **Multilingual** — supports 100+ languages
- **State-of-the-art** MTEB benchmark scores (top proprietary tier)
- **API-only** — cloud-hosted, pay-per-use, no self-hosting option
- Designed for: semantic search, RAG, classification, clustering, reranking

**Why it matters**: Long context + high dimensionality + multilingual = very strong for knowledge base retrieval (e.g., codex/vault RAG systems).

---

## Open-Source Alternatives

These models match or exceed Gemini Embedding 2 in specific domains and are fully self-hostable.

---

### 1. `nomic-embed-text-v2` — Best Overall Open-Source

| Property | Value |
|----------|-------|
| Developer | Nomic AI |
| License | **Apache 2.0** |
| Dimensions | 768 |
| Context | 8,192 tokens |
| Languages | Multilingual |
| Size | ~280M params |
| MTEB Score | Top open-source tier |

**Architecture**: Mixture-of-Experts (MoE) — novel for embedding models. Significantly outperforms older models at its parameter count.

**Why it replaces Gemini**: Longer context (8k vs 2k), completely free, runs locally via Ollama, Apache 2.0 license. For vault/codex RAG, the 8k window means whole notes can be embedded without chunking.

```bash
# Run locally
ollama pull nomic-embed-text
```

**Best for**: RAG pipelines, semantic search, local OneMind codex indexing.

---

### 2. `BGE-M3` — Best Multilingual + Multi-Function

| Property | Value |
|----------|-------|
| Developer | BAAI (Beijing Academy of AI) |
| License | **Apache 2.0** |
| Dimensions | 1024 |
| Context | 8,192 tokens |
| Languages | 100+ |
| Size | ~570M params |
| MTEB Score | Top multilingual tier |

**Architecture**: Three retrieval modes in one model:
- **Dense retrieval** — standard embedding similarity
- **Sparse retrieval** — BM25-style keyword matching
- **ColBERT retrieval** — token-level late interaction (highest precision)

**Why it replaces Gemini**: BGE-M3 with sparse+dense hybrid retrieval outperforms Gemini Embedding 2 on multilingual retrieval benchmarks. The three-mode architecture matches what proprietary models offer.

```bash
# Via HuggingFace
pip install sentence-transformers
from sentence_transformers import SentenceTransformer
model = SentenceTransformer("BAAI/bge-m3")
```

**Best for**: Production RAG with hybrid retrieval, multilingual knowledge bases, pgvector + sparse index combos.

---

### 3. `gte-Qwen2-7B-instruct` — Best Long Context

| Property | Value |
|----------|-------|
| Developer | Alibaba (Tongyi) |
| License | **Apache 2.0** |
| Dimensions | 3584 |
| Context | **131,072 tokens** |
| Languages | Multilingual |
| Size | 7B params |
| MTEB Score | Near-proprietary top tier |

**Architecture**: Qwen2 LLM backbone fine-tuned for embedding. The 131k context window is unmatched.

**Why it replaces Gemini**: Gemini Embedding 2 tops out at 2048 tokens. `gte-Qwen2-7B` handles entire documents, codebases, or long conversation histories as a single embedding. Also 3584 dims vs 3072 — denser representation.

**Best for**: Document-level search, embedding entire files, long-form content retrieval. Needs a GPU (7B params).

---

### 4. `e5-mistral-7b-instruct` — Best Instruction-Tuned

| Property | Value |
|----------|-------|
| Developer | Microsoft Research |
| License | **MIT** |
| Dimensions | 4096 |
| Context | 32,768 tokens |
| Languages | Multilingual |
| Size | 7B params |
| MTEB Score | Top open-source on many tasks |

**Architecture**: Mistral 7B backbone with instruction tuning for embeddings. Prepend task instructions like `Represent this sentence for searching relevant passages:` to guide embedding generation.

**Why it replaces Gemini**: MIT license, 32k context, 4096 dims (more expressive than Gemini's 3072), and instruction-following means you can optimize embeddings per task type (search vs clustering vs classification).

**Best for**: Task-specific embedding pipelines, research, high-quality retrieval with GPU resources.

---

### 5. `mxbai-embed-large-v1` — Best Small/Fast

| Property | Value |
|----------|-------|
| Developer | Mixedbread AI |
| License | **Apache 2.0** |
| Dimensions | 1024 |
| Context | 512 tokens |
| Languages | English-primary |
| Size | ~335M params |
| MTEB Score | #1 on MTEB at its size class (as of early 2026) |

**Why it replaces Gemini**: For English-only tasks with short text (product descriptions, support tickets, short notes), `mxbai-embed-large` beats Gemini on MTEB while being tiny, fast, and free.

```bash
ollama pull mxbai-embed-large
```

**Best for**: Fast local inference, high-throughput pipelines, English-only use cases.

---

### 6. `Jina Embeddings v3` — Best Web/Multilingual

| Property | Value |
|----------|-------|
| Developer | Jina AI |
| License | **Apache 2.0** |
| Dimensions | 1024 |
| Context | 8,192 tokens |
| Languages | 89 languages |
| Size | ~570M params |

**Why it replaces Gemini**: Apache 2.0, 89 languages, 8k context, task-specific adapters (retrieval, STS, classification) all in one model. API available but fully self-hostable.

**Best for**: Multilingual RAG, web content indexing, production-grade open-source pipelines.

---

### 7. `all-MiniLM-L6-v2` — Best Ultra-Lightweight

| Property | Value |
|----------|-------|
| Developer | Microsoft / SBERT team |
| License | **Apache 2.0** |
| Dimensions | 384 |
| Context | 256 tokens |
| Size | ~22M params |

**Why it still matters**: Tiny, runs on CPU in milliseconds. For quick fuzzy search, autocomplete, or edge/IoT applications where Gemini's API latency is unacceptable.

```bash
ollama pull all-minilm
```

**Best for**: Real-time search suggestions, on-device embedding, latency-critical paths.

---

## Direct Comparison: Open-Source vs Gemini Embedding 2

| Model | Context | Dims | Multilingual | Self-Host | License | Speed | Quality |
|-------|---------|------|-------------|-----------|---------|-------|---------|
| **Gemini Embedding 2** | 2048 | 3072 | ✅ 100+ | ❌ API only | Proprietary | Fast (API) | ⭐⭐⭐⭐⭐ |
| `nomic-embed-text-v2` | 8192 | 768 | ✅ | ✅ | Apache 2.0 | Fast | ⭐⭐⭐⭐½ |
| `BGE-M3` | 8192 | 1024 | ✅ 100+ | ✅ | Apache 2.0 | Medium | ⭐⭐⭐⭐⭐ |
| `gte-Qwen2-7B` | 131072 | 3584 | ✅ | ✅ (GPU) | Apache 2.0 | Slow | ⭐⭐⭐⭐⭐ |
| `e5-mistral-7b` | 32768 | 4096 | ✅ | ✅ (GPU) | MIT | Slow | ⭐⭐⭐⭐⭐ |
| `mxbai-embed-large` | 512 | 1024 | Partial | ✅ | Apache 2.0 | Very Fast | ⭐⭐⭐⭐ |
| `Jina v3` | 8192 | 1024 | ✅ 89 | ✅ | Apache 2.0 | Fast | ⭐⭐⭐⭐½ |
| `all-MiniLM-L6` | 256 | 384 | Partial | ✅ | Apache 2.0 | Instant | ⭐⭐⭐ |

**Key wins over Gemini**:
- `nomic-embed-text-v2` → 4× longer context, no API costs, local
- `BGE-M3` → hybrid retrieval modes Gemini doesn't offer
- `gte-Qwen2-7B` → 64× longer context (131k vs 2k)
- `e5-mistral-7b` → higher dims, 16× longer context, MIT license
- All → no per-token cost, no network dependency, no data leaving your infrastructure

---

## OneMind Stack Recommendations

### For Codex/Vault RAG (Legacy's memory search)

**Recommended**: `BGE-M3` via pgvector + hybrid index
**Reason**: 8k context fits full codex notes, sparse+dense hybrid matches how Legacy searches (keywords + semantic)

```python
# Already compatible with your pgvector setup
from sentence_transformers import SentenceTransformer
model = SentenceTransformer("BAAI/bge-m3")
embedding = model.encode("query text", normalize_embeddings=True)
```

### For Local/Fast Inference (Ollama)

```bash
ollama pull nomic-embed-text   # best quality/speed ratio
ollama pull mxbai-embed-large  # fastest
```

### For Long Documents

**Recommended**: `gte-Qwen2-7B-instruct`
**Reason**: 131k context means entire skill files, workspace docs, or memory logs can be embedded as-is

### For Production API Replacement

Drop-in Gemini Embedding 2 replacement via **Jina v3 API** (same interface, open-source model, self-hostable) or run `nomic-embed-text-v2` behind a FastAPI endpoint.

---

## Resources

- [MTEB Leaderboard](https://huggingface.co/spaces/mteb/leaderboard) — Live benchmark rankings
- [Nomic embed-text-v2](https://huggingface.co/nomic-ai/nomic-embed-text-v2)
- [BGE-M3](https://huggingface.co/BAAI/bge-m3)
- [gte-Qwen2-7B](https://huggingface.co/Alibaba-NLP/gte-Qwen2-7B-instruct)
- [e5-mistral-7b](https://huggingface.co/intfloat/e5-mistral-7b-instruct)
- [mxbai-embed-large](https://huggingface.co/mixedbread-ai/mxbai-embed-large-v1)
- [Jina Embeddings v3](https://huggingface.co/jinaai/jina-embeddings-v3)
- [Gemini Embedding 2 announcement](https://blog.google/innovation-and-ai/models-and-research/gemini-models/gemini-embedding-2/)

---

*Created: 2026-03-16 | Category: AI Models / Embeddings / Research*

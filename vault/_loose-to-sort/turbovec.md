# turbovec — One Mind RAG Engine

> Rust vector index with Python bindings. Faster than FAISS, 8x more memory efficient. 100% local. No data leaves the machine.

## Source
- GitHub: https://github.com/RyanCodrai/turbovec
- Built on Google Research's TurboQuant algorithm
- Benchmarked on Apple M ARM — our exact hardware

## Status
- [ ] Installed ✅
- [ ] Ingest pipeline built
- [ ] Embedded into Legacy query loop

## Planned Stack

```
one-mind/              ← single source of truth
     ↓
Ingest pipeline        ← parse by file type
  .md   → direct text
  .pdf  → pymupdf / marker-pdf
  .mp3  → Whisper → transcript
  .mp4  → Whisper → transcript
  .csv  → pandas → text rows
     ↓
Embedding model        ← Ollama + nomic-embed-text (local)
     ↓
turbovec index         ← stored at _system/rag/
     ↓
Legacy queries         ← natural language over full vault
```

## Framework Integration
```bash
pip install turbovec[langchain]
# or
pip install turbovec[llama-index]
```

## Basic Usage
```python
from turbovec import TurboQuantIndex

index = TurboQuantIndex(dim=1536, bit_width=4)
index.add(vectors)
scores, indices = index.search(query, k=10)
index.write("_system/rag/one-mind.tv")
```

## Notes
- turbovec = vector search only — does NOT parse files
- Ingestion pipeline needs to be built per file type
- One folder (one-mind/) = one RAG brain = one index
- Exclude heavy media from indexing — text/md only for speed
- PDFs added selectively when needed searchable

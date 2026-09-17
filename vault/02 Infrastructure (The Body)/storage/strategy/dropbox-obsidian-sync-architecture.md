---
id: so-02-20260514-dropbox-obsidian-sync
title: "Dropbox + Obsidian Sync — Architecture Decision & RAG Strategy"
type: reference
status: active
created: 2026-04-01
modified: 2026-05-14
owner: zeus
tags: [path/so, domain/02-infrastructure, type/reference, topic/dropbox, topic/obsidian-sync, topic/rag, topic/vps, status/active]
ai_summary: "Research from Scan 1: Dropbox CLI (dbxcli), pgvector RAG pipeline, Obsidian Sync Standard recommendation. Final: $16/mo for 2TB Dropbox + Obsidian Sync Standard."
source: "Apple Notes — Scan 1 (2026-04-01)"
---

## Conclusion / Recommendation

**Option A (Dropbox + Obsidian Sync Standard) is the right move.**

| Tool | Purpose | Cost |
|------|---------|------|
| Dropbox 2TB | Binary files (PDFs, images), VPS agent access, RAG, USB-portable | ~$12/mo |
| Obsidian Sync Standard | Mobile read/write for markdown vault (618 MB fits in 1 GB limit) | $4/mo |
| **Total** | | **~$16/mo** |

- Git = version control backbone (markdown only, committed + pushed to GitHub)
- Dropbox = everything Git doesn't track (PDFs, images, binaries)
- Obsidian Sync = phone access (markdown only — don't sync PDFs/images to phone)
- Dual-sync conflict risk is low (one person, <5s window to create conflict; Dropbox creates a "conflicted copy" — no data loss)

---

## Dropbox CLI — dbxcli

Two tools available:
- **dbxcli** — Dropbox's official CLI (Go binary)
- Dropbox HTTP API v2

### Install on Mac / VPS

```bash
brew install dbxcli
# Or: github.com/dropbox/dbxcli
```

### Key Commands

```bash
dbxcli ls /onemind-codex/           # List files
dbxcli get /onemind-codex/file.md   # Download file
dbxcli put local.md /onemind-codex/ # Upload file
dbxcli search "credit repair"       # Search filenames
dbxcli du                           # Show space usage
```

### VPS Use Case

Install dbxcli on VPS → agents can read/write binary files (PDFs, images) that Git doesn't track:
```bash
dbxcli get /onemind-codex/_docs/some-contract.pdf
# Agent pulls specific PDF for RAG processing without syncing full vault
```

---

## Dropbox HTTP API v2

Base URL: `api.dropboxapi.com`

| Endpoint | Purpose |
|----------|---------|
| `files/list_folder` | Browse directory |
| `files/download` | Download a file |
| `files/search` | Search by filename |
| `files/upload` | Upload <150 MB |
| `files/list_folder/longpoll` | Watch for changes (event-driven) |

> **Gold:** `longpoll` endpoint — OpenClaw agent subscribes to Dropbox change events and reacts to new files automatically.

---

## RAG / PDF Indexing Options

### Option A — OpenClaw Native Search (Already Running)

- Hybrid BM25 + vector search over the codex
- `memorySearch.extraPaths → /home/onemind/onemind-os/codex` — already indexing vault
- OpenAI embeddings, MMR diversity reranking, temporal decay
- **Limitation:** Only searches git-tracked files (~3,601 markdown files). Does NOT search PDFs/binaries.

### Option B — pgvector + PDF Processing (Best for this stack)

`onemind-postgres` (PostgreSQL + pgvector) already running on VPS.

**Pipeline:**
1. `dbxcli` or Dropbox API pulls PDFs from Dropbox to VPS
2. Python extracts text (`PyPDF2`, `pdfplumber`, or `unstructured`)
3. Embed with OpenAI or local model (Nomic Embed)
4. Store chunks + vectors in pgvector
5. OpenClaw agents query via custom skill

### Option C — ChromaDB / Milvus Lite (Simpler, Self-Contained)

- **MemSearch** (Zilliz/memsearch) — ONNX local embeddings, no API key needed, BM25+vector
- Mount Dropbox folder → MemSearch watches it → auto-indexes

---

## File Types Indexable for RAG

(Screenshots in original note — general categories: PDF, DOCX, images with OCR, markdown, HTML, CSV, audio with transcription)

---

## Next Steps

- [ ] Purchase Obsidian Sync Standard ($4/mo)
- [ ] Install dbxcli on VPS
- [ ] Authenticate Dropbox API on VPS
- [ ] Decide on RAG pipeline: pgvector (Option B) vs MemSearch (Option C)
- [ ] Configure longpoll webhook for agent-triggered file reactions

---

## Related

- [[3-2-1-backup-strategy]] — Backup architecture
- [[Holistic Storage Location Management]] — Overall storage map
- Domain: SO 02 Infrastructure → Storage → Strategy

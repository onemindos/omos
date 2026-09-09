# Kreuzberg Cloud

**Repo:** github.com/kreuzberg-dev/kreuzberg-cloud  
**License:** BUSL-1.1 (source available, self-hostable)  
**Language:** Rust  
**Role in OneMind:** Document extraction and RAG pipeline for agent document ingestion

REST API + webhook service for extracting text, metadata, tables, and code intelligence from documents. Multi-tenant, presigned uploads, HMAC-signed webhook delivery. RAG-ready: document collections, ingest pipelines, hybrid retrieval with reranking, embedding management.

## Key capabilities

- Extract text, metadata, tables from any document format
- Code-aware extraction across 306 languages (tree-sitter)
- RAG collections with hybrid retrieval + reranking
- Presigned upload URLs, bulk job submission, polling
- Agent plugin via marketplace (`kreuzberg-dev/plugins`) — works with Claude Code, Cursor, Codex, Gemini CLI
- Self-host via Helm on K8s or use managed SaaS (kreuzberg.dev — 10k free pages)

## Why it matters for OneMind

Agents (Hermes, Oracle, Eden) need to ingest field documents — SOPs, maps, PDFs, intel reports. Kreuzberg handles extraction and chunking; the vectors land in the RAG store that agents query.

## Integration path

```
Document uploaded (PDF, DOCX, image, code)
    → Kreuzberg extract API
    → chunks + embeddings
    → RAG collection
    → agent hybrid retrieval on query
```

## Webhook subjects (proposed)

```
evt.docs.{collection}.ingested   — document chunked and indexed
evt.docs.{collection}.failed     — extraction failure
```

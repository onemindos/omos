# LanceDB Maintenance Guide

> OneMind OS Knowledge Base Storage

---

## Overview

LanceDB is an embedded vector database used for RAG (Retrieval-Augmented Generation) in OneMind OS. It stores document embeddings for semantic search.

| Property | Value |
|----------|-------|
| **Type** | Embedded (serverless) |
| **Location** | `/app/data/lancedb` (in container) |
| **Volume** | `agno_data:/app/data` |
| **Search** | Hybrid (vector + full-text BM25) |
| **Embeddings** | AWS Bedrock `amazon.titan-embed-text-v2:0` |
| **Dimensions** | 1024 |
| **Cost** | $0.02/1M tokens |

---

## Why LanceDB?

| Feature | pgvector | LanceDB |
|---------|----------|---------|
| Setup | Extension | Zero config |
| Search | Vector only | Hybrid (vector + BM25) |
| Speed | Good | Optimized for RAG |
| Ops | Manage with Postgres | Self-contained files |

---

## Maintenance Tasks

### Daily: None Required
LanceDB is self-maintaining. No vacuuming, no index rebuilds.

### Weekly: Backup

```bash
# SSH to server
ssh ubuntu@100.102.21.44

# Backup LanceDB data
docker exec onemind-agno tar czf /tmp/lancedb-backup.tar.gz -C /app/data lancedb
docker cp onemind-agno:/tmp/lancedb-backup.tar.gz ~/backups/lancedb-$(date +%Y%m%d).tar.gz
```

### Monthly: Check Storage

```bash
# Check size
docker exec onemind-agno du -sh /app/data/lancedb

# Expected: ~1KB per document chunk
# 1000 documents ≈ 10-50MB depending on chunk size
```

---

## Common Operations

### View Statistics

```bash
curl http://100.102.21.44:8081/api/knowledge/stats
```

### Search Documents

```bash
curl -X POST http://100.102.21.44:8081/api/knowledge/search \
  -H "Content-Type: application/json" \
  -d '{"query": "your search query", "num_results": 5}'
```

### Upload Document

```bash
# File upload
curl -X POST http://100.102.21.44:8081/api/knowledge/upload \
  -F "file=@document.pdf" \
  -F "name=My Document"

# URL indexing
curl -X POST http://100.102.21.44:8081/api/knowledge/url \
  -H "Content-Type: application/json" \
  -d '{"url": "https://example.com/page", "name": "Example Page"}'
```

### Delete Document

```bash
curl -X DELETE http://100.102.21.44:8081/api/knowledge/{document_id}
```

### Reset Database (Delete All)

```bash
# Stop container
docker stop onemind-agno

# Remove LanceDB data
docker run --rm -v agno_data:/data alpine rm -rf /data/lancedb

# Restart container (recreates empty DB)
docker start onemind-agno
```

---

## Restore from Backup

```bash
# Stop container
docker stop onemind-agno

# Remove existing data
docker run --rm -v agno_data:/data alpine rm -rf /data/lancedb

# Restore backup
docker run --rm -v agno_data:/data -v ~/backups:/backups alpine \
  tar xzf /backups/lancedb-20260125.tar.gz -C /data

# Start container
docker start onemind-agno
```

---

## Storage Estimates

| Documents | Chunks (avg 5/doc) | Estimated Size |
|-----------|-------------------|----------------|
| 100 | 500 | 5-10 MB |
| 1,000 | 5,000 | 50-100 MB |
| 10,000 | 50,000 | 500 MB - 1 GB |

Storage grows linearly. Embeddings are 1536 dimensions × 4 bytes = 6KB per chunk plus metadata.

---

## Troubleshooting

### "LanceDB dependencies not installed"

```bash
# In container or during build
pip install lancedb tantivy pyarrow
```

### Search returns empty results

1. Check if documents were uploaded: `GET /api/knowledge/stats`
2. Verify LanceDB folder exists: `docker exec onemind-agno ls -la /app/data/lancedb`
3. Check logs: `docker logs onemind-agno | grep -i lance`

### Slow search performance

LanceDB auto-creates indexes. If slow:
```bash
# Restart to rebuild indexes
docker restart onemind-agno
```

### Corrupted database

Reset and re-upload documents:
```bash
docker exec onemind-agno rm -rf /app/data/lancedb
docker restart onemind-agno
# Re-upload documents via API
```

---

## Configuration

Located in `core/knowledge/lancedb.py`:

```python
# Defaults
uri = "./data/lancedb"           # Storage path
table_name = "onemind_knowledge" # Table name
search_type = SearchType.HYBRID  # Vector + full-text
embedding_model = "text-embedding-3-small"  # OpenAI embeddings
```

Environment variable override:
```bash
LANCEDB_URI=/custom/path/lancedb
```

---

## Cost

| Component | Cost |
|-----------|------|
| LanceDB | Free (open source) |
| Storage | ~$0.10/GB/month (EBS) |
| Embeddings | ~$0.02 per 1M tokens (OpenAI) |

**Example:** 1,000 documents ≈ $0.10 storage + $0.50 one-time embedding cost

---

## Files Reference

| File | Purpose |
|------|---------|
| `core/knowledge/lancedb.py` | LanceDB client wrapper |
| `core/knowledge/api.py` | REST API endpoints |
| `core/knowledge/pipeline.py` | Document processing |
| `core/knowledge/embeddings.py` | OpenAI embedding service |
| `core/knowledge/chunking.py` | Text chunking logic |

---

## See Also

- [ACTIVATION-CHECKLIST.md](../ACTIVATION-CHECKLIST.md) - Deployment steps
- [KNOWLEDGE-RAG-COMPLETE.md](../KNOWLEDGE-RAG-COMPLETE.md) - Full RAG documentation
- [Agno LanceDB Docs](https://docs.agno.com/vectordb/lancedb) - Upstream documentation


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]

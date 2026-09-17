# Linear Issue: Test RAG System with Bedrock Titan

## Issue Details

**Title:** Test Knowledge RAG System with Bedrock Titan Embeddings

**Team:** OneMind OS

**Priority:** High

**Labels:** `testing`, `knowledge`, `rag`, `bedrock`

**Estimate:** 2 points

---

## Description

The Knowledge RAG system has been deployed with AWS Bedrock Titan embeddings (replacing OpenAI). Need comprehensive testing to verify:

1. Document upload works across all supported file types
2. Hybrid search returns relevant results
3. LanceDB storage persists correctly
4. Cost tracking is accurate

### Background

- **Embedding Model:** `amazon.titan-embed-text-v2:0`
- **Dimensions:** 1024
- **Cost:** $0.02/1M tokens (same as OpenAI text-embedding-3-small)
- **Benefit:** No extra API key needed - uses existing AWS credentials

---

## Acceptance Criteria

### 1. File Upload Testing
- [ ] Upload TXT file - verify chunking and embedding
- [ ] Upload PDF file - verify text extraction
- [ ] Upload MD file - verify markdown parsing
- [ ] Upload DOCX file - verify Word extraction
- [ ] Upload code file (.py, .js) - verify syntax preservation
- [ ] Verify error handling for unsupported file types

### 2. URL Indexing Testing
- [ ] Index a simple HTML page
- [ ] Index documentation page (e.g., Agno docs)
- [ ] Verify metadata is stored correctly
- [ ] Test URL that returns non-HTML content

### 3. Search Testing
- [ ] Semantic search - query matches meaning not keywords
- [ ] Keyword search - exact term matching
- [ ] Hybrid search - combined results
- [ ] Test with 0 results query
- [ ] Test relevance ranking

### 4. Storage Testing
- [ ] Verify LanceDB folder created at `/app/data/lancedb`
- [ ] Restart container - verify data persists
- [ ] Delete document - verify removal from index
- [ ] Check stats endpoint after multiple uploads

### 5. Cost Tracking
- [ ] Verify token counting accuracy
- [ ] Verify cost calculation ($0.02/1M tokens)
- [ ] Check stats endpoint returns cumulative totals

---

## Test Commands

### Upload File
```bash
# Create test file
echo "OneMind OS is a personal AI operating system with four pillars: HP, LE, GE, IT." > test.txt

# Upload
curl -X POST http://100.102.21.44:8081/api/knowledge/upload \
  -F "file=@test.txt" \
  -F "name=Test Document"
```

### Index URL
```bash
curl -X POST http://100.102.21.44:8081/api/knowledge/url \
  -H "Content-Type: application/json" \
  -d '{"url": "https://docs.agno.com/introduction", "name": "Agno Introduction"}'
```

### Search
```bash
curl -X POST http://100.102.21.44:8081/api/knowledge/search \
  -H "Content-Type: application/json" \
  -d '{"query": "What are the four pillars?", "num_results": 5}'
```

### Check Stats
```bash
curl http://100.102.21.44:8081/api/knowledge/stats
```

### Delete Document
```bash
curl -X DELETE http://100.102.21.44:8081/api/knowledge/{document_id}
```

---

## Expected Results

### Successful Upload Response
```json
{
  "success": true,
  "name": "Test Document",
  "source": "file",
  "num_chunks": 1,
  "text_length": 85,
  "embedding_tokens": 21,
  "embedding_cost": 4.2e-07,
  "total_time": 0.25
}
```

### Successful Search Response
```json
{
  "query": "What are the four pillars?",
  "results": [
    {
      "rank": 1,
      "text": "OneMind OS is a personal AI operating system with four pillars: HP, LE, GE, IT.",
      "score": 0.89,
      "metadata": {"source": "upload", "filename": "test.txt"}
    }
  ],
  "num_results": 1
}
```

---

## Files to Review

| File | Purpose |
|------|---------|
| `core/knowledge/api.py` | REST API endpoints |
| `core/knowledge/embeddings.py` | Bedrock Titan integration |
| `core/knowledge/pipeline.py` | Document processing |
| `core/knowledge/lancedb.py` | Vector storage |
| `docs/operations/LANCEDB-MAINTENANCE.md` | Maintenance guide |

---

## Related Issues

- OMOS-81: Knowledge/RAG with LanceDb (parent feature)
- Deployment: Bedrock Titan embeddings (completed 2026-01-25)

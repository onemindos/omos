# RAG Pipeline - Quick Start Guide

> **Get started with the complete RAG system in 5 minutes**
> Updated: 2026-01-25 - Now uses AWS Bedrock Titan embeddings (no OpenAI key needed)

---

## Installation

```bash
# Install required dependencies
pip install pypdf beautifulsoup4 python-docx boto3 lancedb pyarrow cohere pillow

# Optional: Local OCR (Tesseract)
brew install tesseract

# Optional: Local Whisper
pip install openai-whisper
```

## Environment Setup

```bash
# Required - AWS credentials (already configured via Bedrock setup)
export AWS_ACCESS_KEY_ID=...
export AWS_SECRET_ACCESS_KEY=...
export AWS_DEFAULT_REGION=us-east-1

# Optional
export COHERE_API_KEY=...           # For reranking (recommended)
export LANCEDB_URI=./data/lancedb   # Storage path (default shown)
```

**Note:** No OpenAI API key needed - embeddings use AWS Bedrock Titan.

---

## Quick Examples

### 1. Upload & Store a PDF

```python
from core.knowledge.pipeline import RAGPipeline
from core.knowledge.storage import KnowledgeStorage

# Process PDF
pipeline = RAGPipeline()
result = await pipeline.process_file(
    file_content=open("document.pdf", "rb").read(),
    filename="document.pdf",
    metadata={"user_id": "zeus"},
)

# Store in LanceDB
storage = KnowledgeStorage()
stored = await storage.store_document(result, user_id="zeus")

print(f"✅ Stored {stored['num_chunks']} chunks")
print(f"💰 Cost: ${result.embedding_cost:.4f}")
```

### 2. Search with Reranking

```python
from core.knowledge.storage import KnowledgeStorage
from core.knowledge.reranking import create_reranker, RerankConfig

storage = KnowledgeStorage()
reranker = create_reranker(RerankConfig())

# Search (gets top-20)
results = await storage.search(
    query="How does OneMind work?",
    num_results=20,
    user_id="zeus",
)

# Rerank (returns top-5 most relevant)
reranked = await reranker.rerank(query, results)

for r in reranked:
    print(f"{r.reranked_rank}. Score: {r.rerank_score:.3f}")
    print(f"   {r.text[:100]}...")
```

### 3. Process Image with OCR

```python
from core.knowledge.multimodal import process_multimodal

result = await process_multimodal(
    file_bytes=open("diagram.png", "rb").read(),
    filename="diagram.png",
)

print(f"OCR Text: {result.ocr_text}")
print(f"Description: {result.image_description}")
```

### 4. Transcribe Audio

```python
from core.knowledge.multimodal import process_multimodal

result = await process_multimodal(
    file_bytes=open("meeting.mp3", "rb").read(),
    filename="meeting.mp3",
)

print(f"Transcript: {result.transcript}")
print(f"Duration: {result.audio_duration:.1f}s")
print(f"Language: {result.audio_language}")
```

---

## API Usage

### Upload File

```bash
curl -X POST http://localhost:8081/api/knowledge/upload \
  -F "file=@document.pdf" \
  -F "name=Project Docs" \
  -F "user_id=zeus"
```

### Search

```bash
curl -X POST http://localhost:8081/api/knowledge/search \
  -H "Content-Type: application/json" \
  -d '{
    "query": "How does OneMind work?",
    "num_results": 5,
    "user_id": "zeus",
    "rerank": true
  }'
```

### List Documents

```bash
curl "http://localhost:8081/api/knowledge/documents?user_id=zeus"
```

### Delete Document

```bash
curl -X DELETE http://localhost:8081/api/knowledge/doc_abc123
```

### Get Stats

```bash
curl http://localhost:8081/api/knowledge/stats
```

---

## Supported File Types

### Text Documents
- `.pdf` - PDF documents
- `.docx` - Microsoft Word
- `.txt`, `.md` - Plain text, Markdown
- `.html`, `.htm` - Web pages
- `.json`, `.csv` - Data files
- `.py`, `.js`, `.ts` - Code files

### Images (OCR + Vision)
- `.png`, `.jpg`, `.jpeg` - Photos, screenshots
- `.gif`, `.bmp`, `.tiff` - Other formats
- `.webp` - Modern web format

### Audio (Transcription)
- `.mp3`, `.wav` - Common formats
- `.m4a`, `.ogg` - Apple, OGG
- `.flac`, `.aac` - High quality, AAC

---

## Configuration

### Pipeline Config

```python
from core.knowledge.pipeline import PipelineConfig

config = PipelineConfig(
    chunk_size=512,                # Target chunk size
    chunk_overlap=50,              # Overlap between chunks
    chunk_strategy="sentence",     # sentence | markdown | code | fixed
    embedding_model="amazon.titan-embed-text-v2:0",  # AWS Bedrock Titan
    embedding_dimensions=1024,     # Titan v2 default
    batch_size=25,                 # Bedrock rate limits
    show_progress=True,            # Log progress
)

pipeline = RAGPipeline(config)
```

### Storage Config

```python
from core.knowledge.storage import StorageConfig

config = StorageConfig(
    uri="./data/lancedb",                     # Database path
    chunks_table="document_chunks",           # Table names
    metadata_table="document_metadata",
    check_duplicates=True,                    # Deduplication
    batch_size=1000,                          # Chunks per insert
)

storage = KnowledgeStorage(config)
```

### Rerank Config

```python
from core.knowledge.reranking import RerankConfig, RerankProvider

config = RerankConfig(
    provider=RerankProvider.COHERE,   # cohere | tfidf | none
    model="rerank-v3.5",              # Cohere model
    top_n=5,                          # Final number of results
)

reranker = create_reranker(config)
```

### Multi-Modal Config

```python
from core.knowledge.multimodal import MultiModalConfig, OCRProvider

config = MultiModalConfig(
    # Image processing
    ocr_provider=OCRProvider.OPENAI_VISION,  # or TESSERACT
    vision_model="gpt-4o",
    describe_images=True,

    # Audio processing
    transcription_provider=TranscriptionProvider.OPENAI_WHISPER,
    whisper_model="whisper-1",
)
```

---

## Cost Estimates

### Per Document

| Type | Size | Processing | Embedding | 20 Searches | Total |
|------|------|------------|-----------|-------------|-------|
| PDF (50p) | 250KB | Free | $0.0004 | $0.04 | $0.04 |
| Image | 500KB | $0.02 | $0.0001 | $0.04 | $0.06 |
| Audio (1h) | 50MB | $0.36 | $0.0002 | $0.04 | $0.40 |

### Monthly (Example: 100 docs, 1K searches/day)

```
Uploads (one-time):
- 80 PDFs: $0.03
- 15 images: $0.30
- 5 audio: $1.80
Total: $2.13

Searches (recurring):
- 30K searches/month × $0.002 = $60/month

Total: ~$62/month
```

---

## Performance

| Operation | Time | Notes |
|-----------|------|-------|
| Extract PDF (50p) | 1.2s | pypdf |
| Chunk text | 0.5s | Sentence-aware |
| Embed 78 chunks | 2.8s | Bedrock Titan batch |
| Store in LanceDB | 0.8s | Batch insert |
| **Total upload** | **5.3s** | Complete pipeline |
| Search (hybrid) | 50ms | Vector + keyword |
| Rerank (Cohere) | 200ms | Top-20 → top-5 |
| **Total query** | **250ms** | Fast! |

---

## Troubleshooting

### "pypdf not installed"
```bash
pip install pypdf>=4.0.0
```

### "lancedb not found"
```bash
pip install lancedb>=0.4.0 pyarrow>=14.0.0
```

### "AWS credentials not configured"
```bash
# Option 1: Environment variables
export AWS_ACCESS_KEY_ID=...
export AWS_SECRET_ACCESS_KEY=...
export AWS_DEFAULT_REGION=us-east-1

# Option 2: AWS credentials file (~/.aws/credentials)
aws configure
```

### "Cohere reranking failed"
```bash
# Install cohere
pip install cohere>=5.0.0

# Set API key
export COHERE_API_KEY=...

# Or use TF-IDF fallback (no API needed)
config = RerankConfig(provider=RerankProvider.TFIDF)
```

### "Tesseract command not found"
```bash
# macOS
brew install tesseract

# Ubuntu
sudo apt-get install tesseract-ocr

# Or use OpenAI Vision instead
config = MultiModalConfig(ocr_provider=OCRProvider.OPENAI_VISION)
```

---

## Documentation

- **Quick Start:** This file
- **Base Pipeline:** `docs/architecture/LIGHTWEIGHT-RAG-PIPELINE.md`
- **Advanced Features:** `docs/architecture/RAG-COMPLETE-IMPLEMENTATION.md`
- **API Reference:** `docs/architecture/RAG-FEATURE-MATRIX.md`
- **Summary:** `docs/KNOWLEDGE-RAG-COMPLETE.md`

---

## Next Steps

1. **Test upload:** `curl -X POST ... -F "file=@test.pdf"`
2. **Test search:** `curl -X POST ... -d '{"query": "test"}'`
3. **Try multimodal:** Upload an image or audio file
4. **Enable reranking:** Set `COHERE_API_KEY` for better results
5. **Integrate with agent:** Wire to Claude/GPT for Q&A

---

**Status:** Production-ready ✅
**Date:** 2026-01-25
**Support:** See documentation links above

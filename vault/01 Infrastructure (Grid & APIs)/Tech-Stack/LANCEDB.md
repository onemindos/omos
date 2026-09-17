# LanceDB

> Serverless vector database for AI applications with native multimodal support.

---

## Overview

| Attribute         | Value                                                           |
| ----------------- | --------------------------------------------------------------- |
| **Type**          | Embedded Vector Database                                        |
| **Storage**       | Local files, S3, GCS, Azure Blob                                |
| **License**       | Apache 2.0                                                      |
| **Languages**     | Python, JavaScript, Rust                                        |
| **Documentation** | [lancedb.github.io/lancedb](https://lancedb.github.io/lancedb/) |

---

## What is LanceDB?

LanceDB is a developer-friendly vector database that provides:

- **Serverless** — No server to manage, runs embedded
- **Multimodal** — Store text, images, audio, video together
- **Columnar format** — Built on Lance format, efficient for ML
- **ANN search** — Fast approximate nearest neighbor search
- **Cloud storage** — Works with S3, GCS, local disk
- **Zero-copy** — Memory-mapped for performance

---

## OneMind Usage

### Primary Use Cases

```
┌─────────────────────────────────────────────────────────────────┐
│                       OneMind LanceDB                            │
│                                                                  │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │                  Document Embeddings                      │   │
│  │  Knowledge chunks • Notes • Research • References        │   │
│  └──────────────────────────────────────────────────────────┘   │
│                                                                  │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │                  Multimodal Search                        │   │
│  │  Images • Screenshots • Voice transcripts • Videos       │   │
│  └──────────────────────────────────────────────────────────┘   │
│                                                                  │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │                  Agent Memory                             │   │
│  │  Episodic memory • Context retrieval • Semantic search   │   │
│  └──────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
```

| Use Case           | Description                         |
| ------------------ | ----------------------------------- |
| **Knowledge Base** | Semantic search over documents      |
| **Image Search**   | Find similar screenshots/photos     |
| **Voice Memory**   | Search transcribed conversations    |
| **Agent Context**  | Retrieve relevant past interactions |
| **RAG Pipeline**   | Fast retrieval for LLM context      |
| **Offline First**  | Works without network access        |

### Why LanceDB for OneMind?

| Advantage           | Benefit                               |
| ------------------- | ------------------------------------- |
| **No server**       | Simpler deployment, embedded in agent |
| **Multimodal**      | Text + images in same DB              |
| **S3 support**      | Cloud backup without complexity       |
| **Fast cold start** | No connection overhead                |
| **Python-native**   | Integrates with Agno easily           |

---

## Installation

```bash
pip install lancedb

# With sentence transformers for embeddings
pip install lancedb sentence-transformers

# With OpenAI embeddings
pip install lancedb openai
```

---

## Basic Usage

### Connect to Database

```python
import lancedb

# Local storage
db = lancedb.connect("./data/lancedb")

# S3 storage (production)
db = lancedb.connect("s3://onemind-data/lancedb")

# Cloud storage with credentials
import os
os.environ["AWS_ACCESS_KEY_ID"] = "..."
os.environ["AWS_SECRET_ACCESS_KEY"] = "..."
db = lancedb.connect("s3://onemind-data/lancedb")
```

### Create Table with Data

```python
import lancedb
from lancedb.pydantic import LanceModel, Vector
from lancedb.embeddings import get_registry

# Define schema with auto-embedding
embed_model = get_registry().get("sentence-transformers").create(
    name="BAAI/bge-small-en-v1.5"
)

class KnowledgeChunk(LanceModel):
    text: str = embed_model.SourceField()
    vector: Vector(384) = embed_model.VectorField()
    source: str
    category: str
    created_at: str

# Connect and create table
db = lancedb.connect("./data/lancedb")
table = db.create_table("knowledge", schema=KnowledgeChunk)

# Add data (embeddings generated automatically)
table.add([
    {
        "text": "The Legacy agent uses Agno framework for tool calling.",
        "source": "docs/architecture.md",
        "category": "architecture",
        "created_at": "2026-01-25"
    },
    {
        "text": "NATS JetStream handles event streaming between components.",
        "source": "docs/infrastructure.md",
        "category": "infrastructure",
        "created_at": "2026-01-25"
    }
])
```

### Search

```python
# Semantic search (auto-embeds query)
results = table.search("how does the agent communicate?").limit(5).to_pandas()

# With metadata filter
results = table.search("agent architecture") \
    .where("category = 'architecture'") \
    .limit(10) \
    .to_pandas()

# Vector search (raw vector)
query_vector = embed_model.compute_query_embeddings("agent tools")[0]
results = table.search(query_vector).limit(5).to_pandas()
```

---

## OneMind Integration

### Knowledge Base Table

```python
from datetime import datetime
from lancedb.pydantic import LanceModel, Vector
from lancedb.embeddings import get_registry

# Use same embedding model as elsewhere
embed_model = get_registry().get("openai").create(
    name="text-embedding-ada-002"
)

class OnemindKnowledge(LanceModel):
    id: str
    content: str = embed_model.SourceField()
    vector: Vector(1536) = embed_model.VectorField()  # ada-002 dimension
    source_type: str  # 'obsidian', 'web', 'manual', 'conversation'
    source_path: str
    domain: str  # 'ui', 'hp', 'le', 'ge'
    tags: list[str]
    created_at: datetime
    updated_at: datetime

# Initialize
db = lancedb.connect("./data/onemind_knowledge")
knowledge_table = db.create_table("knowledge", schema=OnemindKnowledge)
```

### RAG Search Function

```python
async def search_knowledge(
    query: str,
    limit: int = 5,
    domain: str = None,
    source_type: str = None
) -> list[dict]:
    """Search knowledge base for relevant context."""

    search = knowledge_table.search(query)

    # Apply filters
    filters = []
    if domain:
        filters.append(f"domain = '{domain}'")
    if source_type:
        filters.append(f"source_type = '{source_type}'")

    if filters:
        search = search.where(" AND ".join(filters))

    results = search.limit(limit).to_pandas()

    return results.to_dict('records')

# Usage in agent
context = await search_knowledge(
    "how to set up home automation",
    domain="ui",
    limit=10
)
```

### Multimodal Table (Images + Text)

```python
import clip
import torch
from PIL import Image

# CLIP model for images
device = "cuda" if torch.cuda.is_available() else "cpu"
clip_model, preprocess = clip.load("ViT-B/32", device=device)

class MultimodalMemory(LanceModel):
    id: str
    content_type: str  # 'text', 'image', 'audio_transcript'
    text_content: str
    image_path: str = None
    vector: Vector(512)  # CLIP dimension
    created_at: datetime

def embed_content(content_type: str, text: str = None, image_path: str = None):
    """Generate CLIP embedding for text or image."""
    with torch.no_grad():
        if content_type == "image" and image_path:
            image = preprocess(Image.open(image_path)).unsqueeze(0).to(device)
            return clip_model.encode_image(image).cpu().numpy()[0]
        else:
            text_tokens = clip.tokenize([text]).to(device)
            return clip_model.encode_text(text_tokens).cpu().numpy()[0]

# Add image to memory
multimodal_table.add([{
    "id": "img_001",
    "content_type": "image",
    "text_content": "Screenshot of dashboard",
    "image_path": "/screenshots/dashboard.png",
    "vector": embed_content("image", image_path="/screenshots/dashboard.png"),
    "created_at": datetime.now()
}])

# Search by text, find images
results = multimodal_table.search(
    embed_content("text", text="dashboard metrics")
).where("content_type = 'image'").limit(5).to_pandas()
```

### Conversation Memory

```python
class ConversationMemory(LanceModel):
    id: str
    session_id: str
    role: str  # 'user', 'assistant'
    content: str = embed_model.SourceField()
    vector: Vector(1536) = embed_model.VectorField()
    timestamp: datetime
    metadata: dict

async def store_conversation(
    session_id: str,
    role: str,
    content: str,
    metadata: dict = None
):
    """Store conversation turn in memory."""
    conversation_table.add([{
        "id": f"{session_id}_{datetime.now().timestamp()}",
        "session_id": session_id,
        "role": role,
        "content": content,
        "timestamp": datetime.now(),
        "metadata": metadata or {}
    }])

async def get_relevant_history(query: str, limit: int = 10) -> list[dict]:
    """Retrieve semantically relevant conversation history."""
    results = conversation_table.search(query) \
        .limit(limit) \
        .to_pandas()
    return results.to_dict('records')
```

---

## Index Configuration

### IVF-PQ Index (For Large Datasets)

```python
# Create index for faster search on large tables
table.create_index(
    metric="cosine",
    num_partitions=256,
    num_sub_vectors=96,
    accelerator="cuda"  # GPU acceleration if available
)
```

### Search Parameters

```python
# Adjust search quality vs speed
results = table.search(query) \
    .limit(10) \
    .nprobes(20) \  # More partitions = better quality, slower
    .refine_factor(10) \  # Re-rank top results
    .to_pandas()
```

---

## Data Management

### Update Records

```python
# Update existing records
table.update(
    where="source_path = 'docs/old.md'",
    values={"source_path": "docs/new.md", "updated_at": datetime.now()}
)
```

### Delete Records

```python
# Delete by filter
table.delete("source_type = 'deprecated'")

# Delete by IDs
table.delete("id IN ('chunk_001', 'chunk_002')")
```

### Incremental Sync

```python
async def sync_obsidian_notes(vault_path: str):
    """Sync Obsidian notes to LanceDB."""
    import os
    from pathlib import Path

    for md_file in Path(vault_path).rglob("*.md"):
        # Check if already indexed
        existing = table.search(f"source_path = '{md_file}'").limit(1).to_pandas()

        if len(existing) == 0:
            # New file - add
            content = md_file.read_text()
            chunks = chunk_document(content)

            for i, chunk in enumerate(chunks):
                table.add([{
                    "id": f"{md_file}_{i}",
                    "content": chunk,
                    "source_type": "obsidian",
                    "source_path": str(md_file),
                    "domain": detect_domain(md_file),
                    "tags": extract_tags(content),
                    "created_at": datetime.now(),
                    "updated_at": datetime.now()
                }])
```

---

## Production Deployment

### S3 Backend

```python
import os

# Configure AWS credentials
os.environ["AWS_ACCESS_KEY_ID"] = "your_key"
os.environ["AWS_SECRET_ACCESS_KEY"] = "your_secret"
os.environ["AWS_REGION"] = "us-east-1"

# Connect to S3
db = lancedb.connect("s3://onemind-data/lancedb/")
```

### Backup and Restore

```python
# LanceDB uses Lance format - backup is just file copy
import shutil

# Backup
shutil.copytree("./data/lancedb", "./backups/lancedb_20260125")

# For S3, use aws cli
# aws s3 sync s3://onemind-data/lancedb s3://onemind-backups/lancedb_20260125
```

### Concurrent Access

```python
# LanceDB supports concurrent readers
# For writes, use connection per operation

async def add_with_lock(data: list[dict]):
    # Open fresh connection for write
    db = lancedb.connect("./data/lancedb")
    table = db.open_table("knowledge")
    table.add(data)
```

---

## Hybrid Search

Combine vector search with full-text:

```python
# Enable full-text search
table.create_fts_index("content")

# Hybrid search
results = table.search(
    query,
    query_type="hybrid"  # Combines vector + FTS
).limit(10).to_pandas()

# With reranking
from lancedb.rerankers import ColbertReranker

reranker = ColbertReranker()
results = table.search(query) \
    .limit(20) \
    .rerank(reranker) \
    .limit(5) \
    .to_pandas()
```

---

## Monitoring

### Table Statistics

```python
# Get table info
table = db.open_table("knowledge")
print(f"Rows: {table.count_rows()}")
print(f"Schema: {table.schema}")

# List all tables
print(db.table_names())
```

### Query Performance

```python
import time

start = time.time()
results = table.search(query).limit(10).to_pandas()
elapsed = time.time() - start

print(f"Search took {elapsed*1000:.2f}ms")
print(f"Results: {len(results)}")
```

---

## Best Practices

### Chunking Strategy

```python
def chunk_document(text: str, chunk_size: int = 500, overlap: int = 50) -> list[str]:
    """Split document into overlapping chunks."""
    chunks = []
    start = 0

    while start < len(text):
        end = start + chunk_size
        chunk = text[start:end]

        # Try to break at sentence boundary
        if end < len(text):
            last_period = chunk.rfind('.')
            if last_period > chunk_size * 0.7:
                end = start + last_period + 1
                chunk = text[start:end]

        chunks.append(chunk.strip())
        start = end - overlap

    return chunks
```

### Embedding Model Selection

| Model                    | Dimension | Speed  | Quality    |
| ------------------------ | --------- | ------ | ---------- |
| `BAAI/bge-small-en-v1.5` | 384       | Fast   | Good       |
| `BAAI/bge-base-en-v1.5`  | 768       | Medium | Better     |
| `text-embedding-ada-002` | 1536      | API    | Best       |
| `CLIP ViT-B/32`          | 512       | Medium | Multimodal |

---

## Troubleshooting

### Import Errors

```bash
# Install with all dependencies
pip install lancedb[embeddings]

# For GPU support
pip install lancedb[cuda]
```

### S3 Access Issues

```python
# Verify credentials
import boto3
s3 = boto3.client('s3')
s3.list_buckets()  # Should work

# Check bucket permissions
```

### Slow Search

```python
# Create index for large tables (>100k rows)
table.create_index(num_partitions=256)

# Use nprobes for quality control
results = table.search(query).nprobes(10).limit(5).to_pandas()
```

---

## Resources

- [LanceDB Documentation](https://lancedb.github.io/lancedb/)
- [LanceDB GitHub](https://github.com/lancedb/lancedb)
- [Lance Format](https://github.com/lancedb/lance)
- [Embedding Models](https://lancedb.github.io/lancedb/embeddings/)

---

_Last updated: January 2026_

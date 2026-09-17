# LanceDB

> **Serverless Vector Database**

## Overview

LanceDB is a serverless vector database for AI applications. It's used for fast similarity search, document retrieval, and AI memory in OneMind OS.

## Status

- **Active**: ✅ Yes
- **Mode**: Embedded / Local
- **Last Updated**: January 25, 2026

## Use Cases

1. **Semantic Search** - Natural language document search
2. **AI Memory** - Long-term agent memory
3. **RAG** - Retrieval-augmented generation
4. **Document Retrieval** - Find relevant content

## Key Features

- Serverless (no separate server needed)
- Embedded in application
- Disk-based (handles large datasets)
- Native Python integration
- Multi-modal (text, images)
- Full-text + vector hybrid search

## Installation

```bash
pip install lancedb
```

## Basic Usage

### Create Database & Table

```python
import lancedb
from lancedb.embeddings import get_registry

# Connect (creates directory if not exists)
db = lancedb.connect("./data/lancedb")

# Create table with data
data = [
    {"text": "Meeting with team at 3pm", "category": "work"},
    {"text": "Buy groceries after work", "category": "personal"},
    {"text": "Review PR #234", "category": "work"},
]

# Using built-in embedding function
model = get_registry().get("sentence-transformers").create()

table = db.create_table("memories", data=data, embedding=model)
```

### Search

```python
# Vector search
results = table.search("work tasks for today").limit(5).to_list()

# With metadata filter
results = (
    table
    .search("tasks")
    .where("category = 'work'")
    .limit(5)
    .to_list()
)
```

## Embedding Models

### Built-in Models

```python
from lancedb.embeddings import get_registry

# Sentence Transformers (local)
model = get_registry().get("sentence-transformers").create(
    name="all-MiniLM-L6-v2"
)

# OpenAI
model = get_registry().get("openai").create(
    name="text-embedding-3-small"
)

# Cohere
model = get_registry().get("cohere").create()
```

### Custom Embedding Function

```python
from lancedb.embeddings import EmbeddingFunction
import openai

class CustomEmbedding(EmbeddingFunction):
    def embed_documents(self, documents):
        response = openai.embeddings.create(
            model="text-embedding-3-small",
            input=documents
        )
        return [e.embedding for e in response.data]

    def embed_query(self, query):
        return self.embed_documents([query])[0]
```

## Schema Definition

### With Pydantic

```python
from lancedb.pydantic import LanceModel, Vector
from datetime import datetime

class Memory(LanceModel):
    text: str
    category: str
    embedding: Vector(384)  # Dimension of your embedding model
    created_at: datetime = datetime.now()

# Create table with schema
table = db.create_table("memories", schema=Memory)

# Add data
memory = Memory(
    text="Important meeting tomorrow",
    category="work",
    embedding=model.embed_query("Important meeting tomorrow")
)
table.add([memory])
```

## Agent Memory Integration

### Memory Store

```python
import lancedb
from datetime import datetime
from lancedb.embeddings import get_registry

class AgentMemory:
    def __init__(self, db_path: str = "./data/agent_memory"):
        self.db = lancedb.connect(db_path)
        self.model = get_registry().get("sentence-transformers").create()

        if "memories" not in self.db.table_names():
            self.table = self.db.create_table("memories", schema={
                "id": "string",
                "content": "string",
                "role": "string",
                "session_id": "string",
                "timestamp": "timestamp",
                "vector": f"vector[{self.model.ndims()}]"
            })
        else:
            self.table = self.db.open_table("memories")

    def add(self, content: str, role: str, session_id: str):
        embedding = self.model.compute_source_embeddings([content])[0]
        self.table.add([{
            "id": str(uuid.uuid4()),
            "content": content,
            "role": role,
            "session_id": session_id,
            "timestamp": datetime.now(),
            "vector": embedding
        }])

    def search(self, query: str, session_id: str = None, limit: int = 10):
        search = self.table.search(query, vector_column_name="vector")
        if session_id:
            search = search.where(f"session_id = '{session_id}'")
        return search.limit(limit).to_list()

    def get_session_history(self, session_id: str, limit: int = 50):
        return (
            self.table
            .search()
            .where(f"session_id = '{session_id}'")
            .limit(limit)
            .to_pandas()
            .sort_values("timestamp")
        )
```

### Usage

```python
memory = AgentMemory()

# Store conversation
memory.add("What's on my calendar?", "user", "session-123")
memory.add("You have 3 meetings today...", "assistant", "session-123")

# Search memories
relevant = memory.search("calendar meetings", session_id="session-123")
```

## Hybrid Search

### Full-Text + Vector

```python
from lancedb import FTS

# Create full-text index
table.create_fts_index("text")

# Hybrid search
results = (
    table
    .search("meeting schedule", query_type="hybrid")
    .limit(10)
    .to_list()
)
```

## Data Operations

### Update Records

```python
table.update(
    where="category = 'work'",
    values={"priority": "high"}
)
```

### Delete Records

```python
table.delete("category = 'archive'")
```

### Upsert

```python
table.merge_insert(
    "id"  # Primary key column
).when_matched_update_all().when_not_matched_insert_all().execute([
    {"id": "1", "text": "Updated text", ...}
])
```

## Index Configuration

### IVF Index

```python
table.create_index(
    metric="cosine",
    num_partitions=256,
    num_sub_vectors=96,
    index_type="IVF_PQ"
)
```

## Storage Options

### Local File System

```python
db = lancedb.connect("./data/lancedb")
```

### S3

```python
db = lancedb.connect("s3://bucket/path/to/lancedb")
```

### Cloud (LanceDB Cloud)

```python
db = lancedb.connect(
    "db://onemind-db",
    api_key="your-api-key"
)
```

## Comparison with Alternatives

| Feature     | LanceDB | pgvector | Pinecone | Chroma |
| ----------- | ------- | -------- | -------- | ------ |
| Serverless  | ✅      | ❌       | ✅       | ✅     |
| Self-hosted | ✅      | ✅       | ❌       | ✅     |
| Disk-based  | ✅      | ✅       | ❌       | ⚠️     |
| Full-text   | ✅      | ✅       | ❌       | ❌     |
| Multi-modal | ✅      | ❌       | ❌       | ❌     |
| Cost        | Free    | Free     | Paid     | Free   |

## Best Practices

- ✅ Use appropriate embedding dimensions
- ✅ Create indexes for large tables
- ✅ Use metadata filtering for scoped searches
- ✅ Batch insert for bulk operations
- ✅ Choose the right distance metric (cosine for text)
- ✅ Use Pydantic schemas for type safety

## Resources

- [LanceDB Documentation](https://lancedb.github.io/lancedb/)
- [GitHub Repository](https://github.com/lancedb/lancedb)
- [Python Examples](https://github.com/lancedb/lancedb/tree/main/examples)

---

**Related**: [PostgreSQL](./PostgreSQL.md) | [Agno](./Agno.md)

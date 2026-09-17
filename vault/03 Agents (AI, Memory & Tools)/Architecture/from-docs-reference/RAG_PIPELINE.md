# RAG Pipeline Architecture

> Retrieval-Augmented Generation: The architecture for AI systems that use external knowledge.

---

## The RAG Pipeline

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           RAG PIPELINE                                       │
│                                                                              │
│   ┌─────────┐   ┌────────┐   ┌───────┐                                      │
│   │ EXTRACT │──▶│ ENCODE │──▶│ INDEX │                                      │
│   └─────────┘   └────────┘   └───────┘                                      │
│        │             │            │                                          │
│   Documents     Embeddings    Vector DB                                      │
│   (PDFs, web,   (dense        (searchable                                   │
│    APIs, DB)    vectors)      knowledge)                                    │
│                                                                              │
│   ════════════════════════════════════════════════════════════════════════  │
│                         INGESTION PHASE ↑  │  ↓ QUERY PHASE                 │
│   ════════════════════════════════════════════════════════════════════════  │
│                                                                              │
│        Query                                                                 │
│          │                                                                   │
│          ▼                                                                   │
│   ┌──────────┐   ┌────────┐   ┌──────────┐   ┌──────────┐                   │
│   │ RETRIEVE │──▶│ RERANK │──▶│ GENERATE │──▶│ EVALUATE │                   │
│   └──────────┘   └────────┘   └──────────┘   └──────────┘                   │
│        │              │             │              │                         │
│   Semantic       Re-order       LLM with       Quality                      │
│   search         by relevance   context        metrics                      │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## The Seven Stages

### 1. EXTRACT

**Purpose:** Get raw data from various sources and prepare it for processing.

| Source Type | Examples | Extraction Method |
|-------------|----------|-------------------|
| Documents | PDF, DOCX, TXT | PyPDF, python-docx, Unstructured |
| Web | HTML pages, blogs | BeautifulSoup, Scrapy, Puppeteer |
| APIs | REST, GraphQL | httpx, requests |
| Databases | SQL, NoSQL | SQLAlchemy, pymongo |
| Media | Audio, Video | Whisper (transcription) |
| Code | Repositories | Tree-sitter, AST parsing |

**Key Considerations:**
```python
# Document Loading Example
from langchain_community.document_loaders import (
    PyPDFLoader,
    WebBaseLoader,
    UnstructuredMarkdownLoader
)

# Load PDF
pdf_loader = PyPDFLoader("document.pdf")
pdf_docs = pdf_loader.load()

# Load Web Page
web_loader = WebBaseLoader("https://example.com/article")
web_docs = web_loader.load()

# Load Markdown
md_loader = UnstructuredMarkdownLoader("README.md")
md_docs = md_loader.load()
```

**Chunking Strategies:**

| Strategy | Best For | Chunk Size |
|----------|----------|------------|
| Fixed size | General text | 512-1024 tokens |
| Sentence-based | Articles, blogs | 3-5 sentences |
| Paragraph-based | Structured docs | 1-3 paragraphs |
| Semantic | Technical docs | By topic/section |
| Recursive | Code, mixed content | Hierarchical |

```python
from langchain.text_splitter import RecursiveCharacterTextSplitter

splitter = RecursiveCharacterTextSplitter(
    chunk_size=1000,
    chunk_overlap=200,  # Overlap prevents context loss at boundaries
    separators=["\n\n", "\n", ". ", " ", ""]
)

chunks = splitter.split_documents(documents)
```

---

### 2. ENCODE

**Purpose:** Convert text into dense vector representations (embeddings).

```
┌─────────────────────────────────────────────────────────────────┐
│                                                                  │
│   "The quick brown fox"  ──▶  [0.023, -0.041, 0.089, ... ]      │
│                                                                  │
│   Text (human readable)       Embedding (machine readable)       │
│                               1536 dimensions (OpenAI)           │
│                               768 dimensions (sentence-transformers)
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

**Embedding Models:**

| Model | Dimensions | Speed | Quality | Cost |
|-------|------------|-------|---------|------|
| OpenAI text-embedding-3-large | 3072 | Fast | Excellent | $0.13/1M tokens |
| OpenAI text-embedding-3-small | 1536 | Faster | Great | $0.02/1M tokens |
| Cohere embed-v3 | 1024 | Fast | Excellent | $0.10/1M tokens |
| Voyage AI | 1024 | Fast | Excellent | $0.10/1M tokens |
| BGE-large-en | 1024 | Medium | Great | Free (local) |
| all-MiniLM-L6-v2 | 384 | Very Fast | Good | Free (local) |
| nomic-embed-text | 768 | Fast | Great | Free (local) |

**Implementation:**

```python
from openai import OpenAI

client = OpenAI()

def get_embedding(text: str, model: str = "text-embedding-3-small") -> list[float]:
    """Convert text to embedding vector."""
    response = client.embeddings.create(
        input=text,
        model=model
    )
    return response.data[0].embedding

# Batch encoding for efficiency
def batch_encode(texts: list[str], model: str = "text-embedding-3-small") -> list[list[float]]:
    """Encode multiple texts in one API call."""
    response = client.embeddings.create(
        input=texts,
        model=model
    )
    return [item.embedding for item in response.data]
```

**Local Embedding (Free):**

```python
from sentence_transformers import SentenceTransformer

model = SentenceTransformer('all-MiniLM-L6-v2')

# Single text
embedding = model.encode("The quick brown fox")

# Batch (much faster)
embeddings = model.encode([
    "First document",
    "Second document",
    "Third document"
])
```

---

### 3. INDEX

**Purpose:** Store embeddings in a vector database for fast similarity search.

**Vector Database Options:**

| Database | Type | Best For | Scaling |
|----------|------|----------|---------|
| **pgvector** | PostgreSQL extension | Existing Postgres users | Millions |
| **Pinecone** | Managed SaaS | Production, zero-ops | Billions |
| **Weaviate** | Self-hosted/Cloud | Hybrid search | Billions |
| **Qdrant** | Self-hosted/Cloud | Performance critical | Billions |
| **Chroma** | Embedded | Prototyping, local dev | Millions |
| **Milvus** | Self-hosted | Enterprise scale | Billions |
| **FAISS** | Library (Meta) | Research, benchmarking | Billions |
| **LanceDB** | Embedded | Serverless, cheap | Millions |

**pgvector Implementation (OneMindOS Stack):**

```sql
-- Enable extension
CREATE EXTENSION IF NOT EXISTS vector;

-- Create table with embedding column
CREATE TABLE documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    content TEXT NOT NULL,
    embedding vector(1536),  -- Match your model dimensions
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Create HNSW index for fast approximate search
CREATE INDEX ON documents
USING hnsw (embedding vector_cosine_ops)
WITH (m = 16, ef_construction = 64);

-- Insert document with embedding
INSERT INTO documents (content, embedding, metadata)
VALUES (
    'The quick brown fox jumps over the lazy dog',
    '[0.023, -0.041, 0.089, ...]'::vector,
    '{"source": "example.pdf", "page": 1}'
);
```

**Python with pgvector:**

```python
from pgvector.psycopg import register_vector
import psycopg

conn = psycopg.connect("postgresql://...")
register_vector(conn)

# Insert
cur = conn.cursor()
cur.execute(
    "INSERT INTO documents (content, embedding) VALUES (%s, %s)",
    ("Document text", embedding_vector)
)

# Search (returns nearest neighbors)
cur.execute("""
    SELECT content, 1 - (embedding <=> %s) as similarity
    FROM documents
    ORDER BY embedding <=> %s
    LIMIT 10
""", (query_embedding, query_embedding))
```

**Index Types:**

| Index | Speed | Accuracy | Memory |
|-------|-------|----------|--------|
| **IVF** | Fast | ~95% | Low |
| **HNSW** | Faster | ~98% | High |
| **Flat** | Slow (exact) | 100% | Low |

---

### 4. RETRIEVE

**Purpose:** Find the most relevant documents for a given query.

```
┌─────────────────────────────────────────────────────────────────┐
│                                                                  │
│   Query: "How do I reset my password?"                          │
│                                 │                                │
│                                 ▼                                │
│                        ┌──────────────┐                         │
│                        │ Query Encode │                         │
│                        └──────────────┘                         │
│                                 │                                │
│                                 ▼                                │
│                        ┌──────────────┐                         │
│                        │ Vector Search│ ◀── k=10 nearest        │
│                        └──────────────┘                         │
│                                 │                                │
│                                 ▼                                │
│   ┌─────────────────────────────────────────────────────────┐   │
│   │  Doc 1: "Password reset instructions..." (0.92)         │   │
│   │  Doc 2: "Account security settings..." (0.87)           │   │
│   │  Doc 3: "Login troubleshooting guide..." (0.84)         │   │
│   │  ...                                                     │   │
│   └─────────────────────────────────────────────────────────┘   │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

**Retrieval Strategies:**

| Strategy | Description | Use Case |
|----------|-------------|----------|
| **Dense** | Vector similarity only | General semantic search |
| **Sparse** | Keyword matching (BM25) | Exact term matching |
| **Hybrid** | Dense + Sparse combined | Best of both worlds |
| **MMR** | Maximal Marginal Relevance | Diverse results |
| **Multi-query** | Generate query variations | Broader coverage |
| **HyDE** | Hypothetical Document Embeddings | Better query matching |

**Hybrid Search Implementation:**

```python
def hybrid_search(query: str, k: int = 10, alpha: float = 0.5):
    """
    Combine dense (semantic) and sparse (keyword) search.
    alpha: 0 = all sparse, 1 = all dense
    """
    # Dense search (vector similarity)
    query_embedding = get_embedding(query)
    dense_results = vector_db.search(query_embedding, k=k*2)

    # Sparse search (BM25/keyword)
    sparse_results = bm25_search(query, k=k*2)

    # Combine with Reciprocal Rank Fusion
    combined = reciprocal_rank_fusion(
        dense_results,
        sparse_results,
        weights=[alpha, 1-alpha]
    )

    return combined[:k]
```

**Multi-Query Retrieval:**

```python
def multi_query_retrieve(question: str, k: int = 10):
    """Generate multiple query variations for broader retrieval."""

    # Use LLM to generate query variations
    prompt = f"""Generate 3 different versions of this question for search:
    Original: {question}

    Variations:"""

    variations = llm.generate(prompt)  # Returns 3 variations

    # Retrieve for each variation
    all_results = []
    for query in [question] + variations:
        results = retrieve(query, k=k)
        all_results.extend(results)

    # Deduplicate and rank
    return deduplicate_and_rank(all_results)[:k]
```

---

### 5. RERANK

**Purpose:** Re-order retrieved results for better relevance using a more powerful model.

```
┌─────────────────────────────────────────────────────────────────┐
│                     RETRIEVAL vs RERANKING                       │
│                                                                  │
│   Retrieval (Fast, Approximate):                                │
│   - Embedding similarity                                         │
│   - Sub-second for millions of docs                             │
│   - Good recall, okay precision                                 │
│                                                                  │
│   Reranking (Slow, Precise):                                    │
│   - Cross-encoder (sees query + doc together)                   │
│   - Seconds for dozens of docs                                  │
│   - Excellent precision                                          │
│                                                                  │
│   Pipeline: Retrieve 100 → Rerank top 100 → Return top 10       │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

**Reranking Models:**

| Model | Provider | Quality | Speed |
|-------|----------|---------|-------|
| Cohere Rerank v3 | Cohere | Excellent | Fast |
| bge-reranker-large | BAAI | Great | Medium |
| cross-encoder/ms-marco | HuggingFace | Good | Medium |
| Jina Reranker | Jina AI | Great | Fast |
| Voyage Rerank | Voyage AI | Excellent | Fast |

**Implementation:**

```python
import cohere

co = cohere.Client("YOUR_API_KEY")

def rerank(query: str, documents: list[str], top_n: int = 10) -> list[dict]:
    """Rerank documents using Cohere."""

    response = co.rerank(
        query=query,
        documents=documents,
        top_n=top_n,
        model="rerank-english-v3.0"
    )

    return [
        {
            "index": result.index,
            "document": documents[result.index],
            "relevance_score": result.relevance_score
        }
        for result in response.results
    ]

# Usage
retrieved_docs = retrieve(query, k=100)  # Get 100 candidates
reranked = rerank(query, retrieved_docs, top_n=10)  # Keep best 10
```

**Local Reranking (Free):**

```python
from sentence_transformers import CrossEncoder

reranker = CrossEncoder('cross-encoder/ms-marco-MiniLM-L-6-v2')

def local_rerank(query: str, documents: list[str], top_n: int = 10):
    """Rerank using local cross-encoder model."""

    # Create query-document pairs
    pairs = [[query, doc] for doc in documents]

    # Score all pairs
    scores = reranker.predict(pairs)

    # Sort by score
    ranked = sorted(
        zip(documents, scores),
        key=lambda x: x[1],
        reverse=True
    )

    return ranked[:top_n]
```

---

### 6. GENERATE

**Purpose:** Use an LLM to generate a response using retrieved context.

```
┌─────────────────────────────────────────────────────────────────┐
│                                                                  │
│   PROMPT STRUCTURE                                               │
│                                                                  │
│   ┌───────────────────────────────────────────────────────────┐ │
│   │ SYSTEM: You are a helpful assistant. Answer based on the  │ │
│   │ provided context. If the answer isn't in the context,     │ │
│   │ say "I don't have enough information."                    │ │
│   └───────────────────────────────────────────────────────────┘ │
│                              +                                   │
│   ┌───────────────────────────────────────────────────────────┐ │
│   │ CONTEXT:                                                   │ │
│   │ [Document 1]: Password reset can be done via Settings...  │ │
│   │ [Document 2]: For security, passwords must be 12+ chars...│ │
│   │ [Document 3]: Two-factor authentication is recommended... │ │
│   └───────────────────────────────────────────────────────────┘ │
│                              +                                   │
│   ┌───────────────────────────────────────────────────────────┐ │
│   │ QUESTION: How do I reset my password?                     │ │
│   └───────────────────────────────────────────────────────────┘ │
│                              │                                   │
│                              ▼                                   │
│   ┌───────────────────────────────────────────────────────────┐ │
│   │ ANSWER: To reset your password, go to Settings > Account  │ │
│   │ > Security > Reset Password. Your new password must be    │ │
│   │ at least 12 characters. I recommend enabling two-factor   │ │
│   │ authentication for additional security.                   │ │
│   └───────────────────────────────────────────────────────────┘ │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

**Implementation:**

```python
def generate_response(query: str, context_docs: list[str]) -> str:
    """Generate response using retrieved context."""

    # Format context
    context = "\n\n".join([
        f"[Document {i+1}]: {doc}"
        for i, doc in enumerate(context_docs)
    ])

    # Build prompt
    messages = [
        {
            "role": "system",
            "content": """You are a helpful assistant. Answer the user's question
            based ONLY on the provided context. If the answer isn't in the context,
            say "I don't have enough information to answer that."

            Always cite which document(s) you used."""
        },
        {
            "role": "user",
            "content": f"""Context:
{context}

Question: {query}

Answer:"""
        }
    ]

    response = client.chat.completions.create(
        model="gpt-4o",
        messages=messages,
        temperature=0.1  # Low temperature for factual responses
    )

    return response.choices[0].message.content
```

**Advanced: Citation Extraction:**

```python
def generate_with_citations(query: str, context_docs: list[dict]) -> dict:
    """Generate response with extracted citations."""

    # Format context with IDs
    context = "\n\n".join([
        f"[{doc['id']}]: {doc['content']}"
        for doc in context_docs
    ])

    messages = [
        {
            "role": "system",
            "content": """Answer based on the context. Format your response as JSON:
            {
                "answer": "Your answer here",
                "citations": ["doc_id_1", "doc_id_2"],
                "confidence": 0.0-1.0
            }"""
        },
        {
            "role": "user",
            "content": f"Context:\n{context}\n\nQuestion: {query}"
        }
    ]

    response = client.chat.completions.create(
        model="gpt-4o",
        messages=messages,
        response_format={"type": "json_object"}
    )

    return json.loads(response.choices[0].message.content)
```

---

### 7. EVALUATE

**Purpose:** Measure the quality of your RAG system.

**Evaluation Dimensions:**

```
┌─────────────────────────────────────────────────────────────────┐
│                     RAG EVALUATION METRICS                       │
│                                                                  │
│   ┌─────────────┐    ┌─────────────┐    ┌─────────────┐        │
│   │  RETRIEVAL  │    │ GENERATION  │    │   END-TO-END │        │
│   │   QUALITY   │    │   QUALITY   │    │    QUALITY   │        │
│   └─────────────┘    └─────────────┘    └─────────────┘        │
│         │                  │                   │                 │
│   • Precision@K      • Faithfulness      • Answer Relevance    │
│   • Recall@K         • Hallucination     • User Satisfaction   │
│   • MRR              • Coherence         • Task Completion     │
│   • NDCG             • Completeness      • Latency             │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

**Key Metrics:**

| Metric | Measures | Target |
|--------|----------|--------|
| **Context Precision** | Are retrieved docs relevant? | > 0.8 |
| **Context Recall** | Did we find all relevant docs? | > 0.7 |
| **Faithfulness** | Is answer grounded in context? | > 0.9 |
| **Answer Relevance** | Does answer address the question? | > 0.8 |
| **Hallucination Rate** | % of claims not in context | < 5% |

**Evaluation with RAGAS:**

```python
from ragas import evaluate
from ragas.metrics import (
    context_precision,
    context_recall,
    faithfulness,
    answer_relevancy
)
from datasets import Dataset

# Prepare evaluation dataset
eval_data = {
    "question": ["How do I reset my password?", ...],
    "answer": ["To reset your password...", ...],
    "contexts": [["Doc 1 content", "Doc 2 content"], ...],
    "ground_truth": ["The correct answer is...", ...]
}

dataset = Dataset.from_dict(eval_data)

# Run evaluation
results = evaluate(
    dataset,
    metrics=[
        context_precision,
        context_recall,
        faithfulness,
        answer_relevancy
    ]
)

print(results)
# {
#     'context_precision': 0.85,
#     'context_recall': 0.78,
#     'faithfulness': 0.92,
#     'answer_relevancy': 0.88
# }
```

**LLM-as-Judge Evaluation:**

```python
def evaluate_faithfulness(question: str, answer: str, context: str) -> float:
    """Use LLM to evaluate if answer is faithful to context."""

    prompt = f"""Evaluate if the answer is faithful to the given context.

Context: {context}

Question: {question}

Answer: {answer}

Score from 0.0 to 1.0 where:
- 1.0 = Answer is completely supported by context
- 0.5 = Answer is partially supported
- 0.0 = Answer contains claims not in context

Return only the numeric score."""

    response = client.chat.completions.create(
        model="gpt-4o",
        messages=[{"role": "user", "content": prompt}],
        temperature=0
    )

    return float(response.choices[0].message.content.strip())
```

**A/B Testing Framework:**

```python
class RAGExperiment:
    """Track and compare RAG configurations."""

    def __init__(self):
        self.experiments = {}

    def run_experiment(
        self,
        name: str,
        retriever,
        generator,
        test_questions: list[str]
    ):
        results = []

        for question in test_questions:
            # Time the full pipeline
            start = time.time()

            # Retrieve
            docs = retriever.retrieve(question)

            # Generate
            answer = generator.generate(question, docs)

            latency = time.time() - start

            results.append({
                "question": question,
                "answer": answer,
                "docs": docs,
                "latency": latency
            })

        self.experiments[name] = results
        return results

    def compare(self, exp_a: str, exp_b: str):
        """Compare two experiments."""
        # Evaluate both and return comparison metrics
        pass
```

---

## Complete RAG Pipeline Example

```python
from typing import List, Dict, Any
import psycopg
from pgvector.psycopg import register_vector
from openai import OpenAI

class RAGPipeline:
    """Complete RAG pipeline implementation."""

    def __init__(self, db_url: str):
        self.client = OpenAI()
        self.conn = psycopg.connect(db_url)
        register_vector(self.conn)

    # 1. EXTRACT
    def extract(self, source: str) -> List[str]:
        """Extract and chunk documents."""
        from langchain_community.document_loaders import UnstructuredFileLoader
        from langchain.text_splitter import RecursiveCharacterTextSplitter

        loader = UnstructuredFileLoader(source)
        docs = loader.load()

        splitter = RecursiveCharacterTextSplitter(
            chunk_size=1000,
            chunk_overlap=200
        )

        return [chunk.page_content for chunk in splitter.split_documents(docs)]

    # 2. ENCODE
    def encode(self, texts: List[str]) -> List[List[float]]:
        """Convert texts to embeddings."""
        response = self.client.embeddings.create(
            input=texts,
            model="text-embedding-3-small"
        )
        return [item.embedding for item in response.data]

    # 3. INDEX
    def index(self, chunks: List[str], embeddings: List[List[float]], metadata: Dict = None):
        """Store in vector database."""
        cur = self.conn.cursor()

        for chunk, embedding in zip(chunks, embeddings):
            cur.execute(
                """INSERT INTO documents (content, embedding, metadata)
                   VALUES (%s, %s, %s)""",
                (chunk, embedding, metadata or {})
            )

        self.conn.commit()

    # 4. RETRIEVE
    def retrieve(self, query: str, k: int = 10) -> List[Dict]:
        """Find relevant documents."""
        query_embedding = self.encode([query])[0]

        cur = self.conn.cursor()
        cur.execute("""
            SELECT content, metadata, 1 - (embedding <=> %s) as similarity
            FROM documents
            ORDER BY embedding <=> %s
            LIMIT %s
        """, (query_embedding, query_embedding, k))

        return [
            {"content": row[0], "metadata": row[1], "score": row[2]}
            for row in cur.fetchall()
        ]

    # 5. RERANK
    def rerank(self, query: str, docs: List[Dict], top_n: int = 5) -> List[Dict]:
        """Rerank documents for better relevance."""
        import cohere
        co = cohere.Client()

        response = co.rerank(
            query=query,
            documents=[d["content"] for d in docs],
            top_n=top_n,
            model="rerank-english-v3.0"
        )

        return [
            {**docs[r.index], "rerank_score": r.relevance_score}
            for r in response.results
        ]

    # 6. GENERATE
    def generate(self, query: str, context_docs: List[Dict]) -> str:
        """Generate response using context."""
        context = "\n\n".join([
            f"[{i+1}]: {doc['content']}"
            for i, doc in enumerate(context_docs)
        ])

        response = self.client.chat.completions.create(
            model="gpt-4o",
            messages=[
                {
                    "role": "system",
                    "content": "Answer based on the provided context. Cite sources."
                },
                {
                    "role": "user",
                    "content": f"Context:\n{context}\n\nQuestion: {query}"
                }
            ],
            temperature=0.1
        )

        return response.choices[0].message.content

    # 7. EVALUATE
    def evaluate(self, query: str, answer: str, context: List[str]) -> Dict[str, float]:
        """Evaluate response quality."""
        # Faithfulness check
        faithfulness_prompt = f"""
        Context: {' '.join(context)}
        Answer: {answer}

        Is the answer fully supported by the context? Score 0-1:"""

        faith_response = self.client.chat.completions.create(
            model="gpt-4o",
            messages=[{"role": "user", "content": faithfulness_prompt}],
            temperature=0
        )

        return {
            "faithfulness": float(faith_response.choices[0].message.content.strip())
        }

    # FULL PIPELINE
    def query(self, question: str) -> Dict[str, Any]:
        """Run the complete RAG pipeline."""
        # Retrieve
        docs = self.retrieve(question, k=20)

        # Rerank
        reranked = self.rerank(question, docs, top_n=5)

        # Generate
        answer = self.generate(question, reranked)

        # Evaluate
        metrics = self.evaluate(
            question,
            answer,
            [d["content"] for d in reranked]
        )

        return {
            "question": question,
            "answer": answer,
            "sources": reranked,
            "metrics": metrics
        }


# Usage
if __name__ == "__main__":
    rag = RAGPipeline("postgresql://...")

    # Ingest documents
    chunks = rag.extract("knowledge_base.pdf")
    embeddings = rag.encode(chunks)
    rag.index(chunks, embeddings, {"source": "knowledge_base.pdf"})

    # Query
    result = rag.query("How do I reset my password?")
    print(result["answer"])
```

---

## Tool Stack for OneMindOS

| Stage | Tool | Notes |
|-------|------|-------|
| **Extract** | Unstructured, LangChain loaders | Multi-format support |
| **Encode** | OpenAI text-embedding-3-small | Cost-effective, high quality |
| **Index** | pgvector (PostgreSQL) | Already in stack |
| **Retrieve** | pgvector + BM25 hybrid | Best accuracy |
| **Rerank** | Cohere Rerank v3 | Fast, accurate |
| **Generate** | GPT-4o / Claude | Main LLM |
| **Evaluate** | RAGAS + custom | Comprehensive metrics |

---

## Common Pitfalls

| Problem | Cause | Solution |
|---------|-------|----------|
| Irrelevant results | Poor chunking | Tune chunk size, add overlap |
| Missing information | Low recall | Increase k, use hybrid search |
| Hallucinations | Weak grounding | Better prompts, lower temperature |
| Slow queries | Large index | Add HNSW index, reduce k |
| Stale answers | Outdated index | Implement incremental updates |

---

## Resources

- [LangChain RAG Tutorial](https://python.langchain.com/docs/tutorials/rag/)
- [LlamaIndex Documentation](https://docs.llamaindex.ai/)
- [RAGAS Evaluation Framework](https://docs.ragas.io/)
- [pgvector Documentation](https://github.com/pgvector/pgvector)
- [Cohere RAG Guide](https://docs.cohere.com/docs/retrieval-augmented-generation-rag)

---

*Last Updated: 2026-02-15*

---
id: ui-02-20260111053925
title: "Data Flow"
type: concept
pillar: UI
domain: 00-24
folder: "00-24 UI (Unified Intelligence)"
status: active
created: 2026-01-11
modified: 2026-01-11
owner: zeus
tags: [domain/ui, folder/02]
todoist_project: "00-24 UI (Unified Intelligence)"
ai_summary: "Data Flow"
---

![[the actual answer Hybrid with Smart Boundaries.png]]**Data Routing Rules:**

**→ Supabase (pgvector):**

- ClickUp tasks, comments, relationships
- Raindrop bookmarks (metadata + tags)
- Fitness/health metrics (time series)
- Calendar events (structured)
- Email metadata (sender, subject, thread relationships)
- Contact information
- **Anything with relationships or requiring SQL queries**

**→ R2 + AutoRAG/Vectorize:**

- PDF documents
- Word documents, spreadsheets
- Images with text (OCR'd)
- Large text files
- Meeting transcripts
- **Anything that's primarily full-text search**

# Letta Architecture (Simplified)

# 1. Core Memory (Always in context - small)

# - Current conversation

# - User preferences

# - Active goals

  

# 2. Archival Memory (Vector database - searched as needed)

# THIS IS WHERE YOUR DATA LIVES

  

# 3. External Tools (API calls - use sparingly)

# - Real-time data that changes constantly

# - Actions (create task, send email)  
  
from letta import LettaClient

  

client = LettaClient()

  

# Create Letta agent with dual archival sources

agent = client.create_agent(

    archival_storage_connectors=[

        {

            "type": "postgres", 

            "connection_string": "postgresql://supabase...",

            "embedding_config": {

                "model": "text-embedding-3-small",

                "dimensions": 1536

            }

        },

        {

            "type": "vectorize",  # If Letta supports Cloudflare

            "endpoint": "https://api.cloudflare.com/...",

            "index_name": "legacy-documents"

        }

    ]

)

  

# Now when Letta needs context, it queries BOTH automatically  
  
——

**Why This Architecture Wins:**

**1. Natural Data Modeling**

- Tasks are rows in a table (not flattened JSON in R2)
- Documents are blobs in object storage (not database records)

**2. Query Efficiency**

- "Show me all high-priority ClickUp tasks assigned to me" → SQL query (instant)
- "What did we discuss about Q4 strategy?" → Vector search across meeting notes

**3. Cost Optimization**

- Supabase: $25/month for 8GB database (structured data is small)
- R2: $0.015/GB for documents (2TB = $30/month)
- Total: ~$55/month for EVERYTHING

**4. Letta Context Efficiency**

- Letta doesn't load entire database into context
- Semantic search retrieves only relevant chunks
- You can query "tasks from last week" without loading all tasks

**5. Data Sovereignty**

- Supabase can be self-hosted (if needed later)
- R2 is Cloudflare (aligned with your infrastructure)
- No vendor lock-in from AWS/GCP  
      
    ———  
    

**FAQ / Your Specific Questions:**

**Q: Should I use AutoRAG or build custom pipeline?** **A:** Use AutoRAG for documents. It's fine for that. Use Supabase+pgvector for structured data.

**Q: Should everything go to S3/R2 or Postgres?** **A:** BOTH. Different data types need different storage. Don't force documents into Postgres or relational data into object storage.

**Q: Should Letta use tools or archival memory?** **A:** Archival memory for historical data. Tools only for real-time actions (create task, send email).

**Q: Is MindsDB relevant?** **A:** No. MindsDB is for ML predictions on data, not storage/retrieval.

**Q: What about MCP/context windows?** **A:** Exactly - which is why you use Letta's archival memory (vector search) instead of loading everything into context.

**Q: Can Letta access both Supabase and Vectorize?** **A:** Yes. Letta supports custom archival storage connectors. You can write a simple wrapper that queries both.  
  

**The TL;DR Decision:**

**Use Hybrid Architecture:**

- **Supabase + pgvector** for structured/relational data (ClickUp, Raindrop, fitness)
- **R2 + AutoRAG/Vectorize** for documents/files (PDFs, notes, transcripts)
- **Airbyte** for ETL from SaaS apps
- **Letta** queries both via archival memory

**Total Cost: ~$55/month** **Complexity: Medium (but maintainable)** **Scalability: Excellent** **Sovereignty: High**
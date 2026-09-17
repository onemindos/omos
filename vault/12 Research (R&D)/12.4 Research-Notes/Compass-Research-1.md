---
id: ui-15-20260111053926
title: "Self-Hosted AI Agent Infrastructure: Technical Architecture Analysis"
type: concept
pillar: UI
domain: 00-24
folder: "00-24 UI (Unified Intelligence)"
status: active
created: 2026-01-11
modified: 2026-01-11
owner: zeus
tags: [domain/ui, folder/15]
todoist_project: "00-24 UI (Unified Intelligence)"
ai_summary: "Self-Hosted AI Agent Infrastructure: Technical Architecture Analysis"
---

# Self-Hosted AI Agent Infrastructure: Technical Architecture Analysis

**Your farm operations demand offline capability, data sovereignty, and zero vendor lock-in - this architecture delivers all three while eliminating redundant complexity.** Deploy R2 + OCIS, Supabase self-hosted, Airbyte, and Letta with Ollama as your core stack for under $100/month with complete offline functionality. This analysis cuts through 12 critical decision points to eliminate at least 8 unnecessary tools while preserving full capability for your Legacy AI system.

## Critical finding: Most "competitors" solve different problems

**Cloudflare's R2 Auto RAG fundamentally changes your vector database landscape.** This managed service eliminates the need for Pinecone, Milvus, ChromaDB, Qdrant, and Docling entirely - they become redundant when R2 automatically ingests documents, generates embeddings via Workers AI, and stores vectors in Vectorize. Your Obsidian vault gets semantic search capabilities without managing separate vector databases.

**Several platforms on your list aren't ETL tools at all.** Prefect, Dagster, and Kestra are orchestration platforms that schedule when Airbyte runs - they're the conductor, not the musicians. unstructured.io preprocesses documents for LLMs rather than moving data between systems. InfluxDB is a time-series database. Confluent provides Kafka streaming. Understanding these distinctions prevents architectural bloat.

**Your existing n8n and NodeRED deployment already handles workflow orchestration** - adding Kestra creates unnecessary overlap unless you develop complex data pipeline requirements. The modern data stack uses complementary tools for different layers, not competing tools for the same purpose.

## Decision 1: Cloud storage architecture wins with OCIS

**Recommendation: Cloudflare R2 + OCIS (ownCloud Infinite Scale)**

OCIS provides simpler deployment than Nextcloud - single Go binary versus PHP/Apache/MySQL stack, lighter resource requirements (runs on Raspberry Pi to Kubernetes), and Apache 2.0 licensing rather than AGPLv3 copyleft. Both platforms handle Obsidian sync identically via WebDAV using the "Remotely Save" plugin since neither has native Obsidian support.

**Technical configuration for R2 integration:**
- Native S3 support via `s3ng` storage driver
- Stores blobs on R2, metadata on POSIX filesystem
- Region set to "auto" for R2 compatibility
- Handles files up to 5TB with multipart uploads
- Zero egress costs (Cloudflare's killer advantage)

**Docker deployment simplicity:**
```yaml
services:
  ocis:
    image: owncloud/ocis:latest
    environment:
      OCIS_INSECURE: false
      STORAGE_USERS_S3NG_ENDPOINT: https://your-account.r2.cloudflarestorage.com
      STORAGE_USERS_S3NG_REGION: auto
    volumes:
      - ocis-data:/var/lib/ocis
    ports:
      - "9200:9200"
```

**For Nextcloud comparison:** Nextcloud requires separate web server, PHP, Redis, and MySQL/PostgreSQL containers with higher resource consumption and more configuration complexity. Choose Nextcloud only if you need its mature app ecosystem and don't mind AGPLv3 requirements.

**Cost structure:** R2 provides 10GB free storage, then $0.015/GB-month with zero egress fees. For a typical Obsidian vault under 5GB plus additional documents: approximately $0-5/month.

## Decision 2: R2 Auto RAG transforms vector strategy

**Cloudflare R2 Auto RAG (now "AI Search") delivers fully managed RAG** that automatically ingests documents from R2 buckets, converts to vectors via Workers AI, and enables semantic search without managing vector databases. Currently in open beta (free), with production pricing at $0.05 per 100M stored dimensions and $0.01 per million queried dimensions after generous free tiers.

**Document parsing capabilities include:**
- PDF, DOCX, PPTX, XLSX (rich formats, 1MB max)
- HTML, CSV, Markdown (text formats, 4MB max)
- Images with vision-to-text conversion
- Automatic chunking and embedding generation
- Continuous re-indexing on file changes

**How Letta agents access this:**

**Method 1 - Custom Python Tool (Recommended):**
```python
def search_obsidian_vault(query: str) -> str:
    """
    Search the Obsidian vault using R2 AI Search.
    
    Args:
        query (str): The search query
        
    Returns:
        str: JSON string with search results
    """
    import requests
    response = requests.post(
        f"https://{CLOUDFLARE_ACCOUNT}.r2.cloudflarestorage.com/ai-search",
        json={"query": query},
        headers={"Authorization": f"Bearer {API_KEY}"}
    )
    return response.json()["answer"]
```

**Method 2 - Letta Filesystem Tools:**
Mount R2 bucket locally, use built-in `open_file`, `grep_file`, `search_file` tools that agents can call directly without custom code.

**Architecture separation of concerns:**
- **R2 AI Search**: Unstructured knowledge retrieval (Obsidian markdown, PDFs, documents)
- **Postgres pgvector**: Structured data embeddings (ETL from SaaS apps, time-series data)
- **mem0**: Agent memory and user state (conversation context, preferences)

This three-layer approach provides semantic search on documents, vector similarity on structured data, and stateful agent memory - completely different use cases with no redundancy.

### Tools eliminated by R2 Auto RAG + pgvector

**Fully redundant - eliminate these:**
- **Pinecone**: Managed vector DB fully replaced by Vectorize (Cloudflare's vector service)
- **Milvus**: Open-source vector DB for billion-scale deployments - overkill for your use case
- **ChromaDB**: Developer-friendly vector DB - pgvector provides equivalent functionality
- **Qdrant**: High-performance Rust vector engine - covered by Vectorize/pgvector
- **Docling**: IBM's document parser - R2 AI Search handles parsing via Workers AI
- **Standalone Supabase pgvector**: Integrate into your main Postgres deployment instead

**Partially redundant - use only for advanced patterns:**
- **LlamaIndex**: Basic RAG covered by R2 AI Search; useful only for complex retrieval strategies (reranking, multi-hop queries, custom document hierarchies)

**NOT redundant - keep these:**
- **mem0**: Manages agent memory and user state (orthogonal to document retrieval)
- **Meilisearch**: Fast keyword search with typo tolerance (different from semantic search)

**Meilisearch complements rather than competes** with R2 AI Search. Use Meilisearch for instant search-as-you-type UI experiences (sub-50ms responses) and keyword matching, while R2 AI Search handles complex question-answering with LLM-generated comprehensive responses. The Community Edition uses MIT licensing with zero commercial restrictions.

## Decision 3: Supabase self-hosted dominates database choice

**Cloudflare D1 cannot replace Postgres for your requirements** - this is critical. D1 is SQLite-based with a hard 10GB limit, no pgvector support, and incompatible with standard ETL tools expecting Postgres connections. D1 requires separate Vectorize for vectors, creating architectural fragmentation rather than simplification.

**Supabase self-hosted provides the optimal database solution:**

**Technical capabilities:**
- Full PostgreSQL with pgvector pre-installed
- Auto-generated REST API via PostgREST
- GraphQL API via pg_graphql extension
- Realtime subscriptions over WebSocket
- Built-in authentication with GoTrue
- Studio web interface for management

**Deployment via Docker Compose:**
```bash
git clone --depth 1 https://github.com/supabase/supabase
cd supabase/docker
cp .env.example .env
# Edit .env: change POSTGRES_PASSWORD, JWT_SECRET, ANON_KEY, SERVICE_KEY
docker compose up -d
```

Setup time: 30-60 minutes including configuration. The deployment includes everything needed: Postgres with pgvector, Kong gateway, PostgREST API layer, Realtime server, GoTrue auth, and Studio UI.

**Integration architecture for Letta agents:**
```
SaaS Apps (Strava, Todoist, etc.)
    ↓
Airbyte ETL
    ↓
Supabase Postgres (structured data + pgvector embeddings)
    ↓
REST API (auto-generated by PostgREST)
    ↓
Letta Agents (query via custom Python tools)
```

**Cost comparison over 5 years:**
- Supabase Cloud: $1,500 - $7,200 (managed, $25-120/month)
- Supabase Self-Hosted: $1,200 - $2,400 (VPS at $20-40/month)
- Plain Postgres + pgvector: $600 - $1,440 (VPS at $10-30/month)

**Why self-hosted Supabase beats plain Postgres:** The $10-20/month premium over bare Postgres buys you auto-generated REST/GraphQL APIs (eliminating weeks of development), realtime subscriptions, built-in auth, and Studio UI. Letta explicitly recommends Supabase in their documentation, confirming this architectural choice.

**Licensing:** Apache 2.0 for the platform, PostgreSQL License for the database, MIT for various components - all permissive with zero commercial restrictions.

## Decision 4: Airbyte dominates ETL with single-platform strategy

**Airbyte emerges as the clear winner** after eliminating platforms that aren't ETL tools (Prefect, Dagster, Kestra are orchestrators), don't self-host (Fivetran, Rivery), or solve different problems (unstructured.io for document preprocessing, Confluent for streaming, InfluxDB as a database).

**Why Airbyte wins for your use case:**

**1. Custom connector development** - Industry-leading CDK allows building connectors in under 30 minutes with low-code framework. Connectors run as Docker containers (language-agnostic). Essential for your health/fitness apps without pre-built connectors.

**2. Dual destination support** - Native capability to sync single source to both Postgres/Supabase AND S3/R2 simultaneously. Map different data streams to different destinations in one configuration.

**3. Pre-built connector coverage:**
- ✅ Strava: Available and documented
- ✅ Todoist: Available and documented  
- ✅ Common business tools: Excellent (Salesforce, HubSpot, Google Analytics, Slack, etc.)
- ⚠️ Health/fitness apps: Partial (Strava main one, custom connectors needed for others)
- ✅ S3/R2: Full support via S3 connector
- ✅ Postgres/Supabase: Native connector

**4. Active development** - 25,000+ community members, 1,000+ contributors, version 1.0 launched September 2024 (major stability milestone), $181M funding demonstrating commercial viability.

**Critical deployment change in 2024:** Airbyte deprecated Docker Compose, now uses `abctl` (Airbyte Control) which manages Kubernetes automatically. Still runs in Docker but with production-ready architecture:

```bash
# New installation method
curl -LsfS https://get.airbyte.com | bash
abctl local install --host airbyte.yourfarm.com

# Requires 4+ CPUs, 8GB RAM for optimal performance
```

**Portainer compatibility:** Deploy as Kubernetes stack in Portainer rather than traditional Docker Compose. More complex setup but significantly more robust for production.

**Licensing clarity:**
- **Core Platform**: Elastic License 2.0 (source-available, cannot resell as managed service)
- **Connectors & CDK**: MIT License (truly open source)
- **Protocol**: MIT License

**Meltano alternative evaluation:** Meltano offers pure MIT licensing and CLI-first infrastructure-as-code workflows but has no web UI (YAML/command-line only), relies on inconsistently-maintained Singer taps with variable quality, and requires substantially more engineering effort. Choose Meltano only if your team strongly prefers GitOps and CLI tools over graphical interfaces.

**Single-platform recommendation:** Use Airbyte for all ETL rather than mixing platforms. Alternative multi-platform strategies (different tools for SaaS, databases, documents) create unnecessary operational overhead. If pipeline complexity grows, add Dagster for orchestration rather than replacing Airbyte.

**Modern data stack architecture:**
```
Airbyte (Extract & Load)
    ↓
dbt (Transform) - optional, add when needed
    ↓
Dagster/Prefect (Orchestrate) - optional, add when pipeline complexity grows
    ↓
Supabase (Store & Serve)
```

**Cost:** $50-200/month for infrastructure (server resources), $0 for software licensing.

## Decision 5: Multi-tool backup strategy eliminates single-point failures

**No single tool handles all backup requirements effectively** - specialized tools for specialized tasks provides the pragmatic solution.

**Tool 1: rsync for cloud storage → NAS backup**

**Why:** Fast delta transfers over Tailscale VPN, minimal bandwidth consumption, preserves file attributes and permissions, mature and reliable since 1996.

**Implementation for OCIS/Nextcloud:**
```bash
#!/bin/bash
# Daily backup at 2 AM via cron: 0 2 * * *

# For Nextcloud: enable maintenance mode
sudo -u www-data php /var/www/nextcloud/occ maintenance:mode --on

# Backup data, config, theme
rsync -Aavx --delete /var/www/nextcloud/ \
  nas.tailscale:/backup/nextcloud-$(date +%Y%m%d)/

# Backup database
mysqldump --single-transaction -h localhost -u nextcloud -p[password] \
  nextcloud_db > /backup/nextcloud-db-$(date +%Y%m%d).sql

# For Nextcloud: disable maintenance mode
sudo -u www-data php /var/www/nextcloud/occ maintenance:mode --off

# For OCIS: Stop service before backup, restart after
```

**Cost:** $0 (built-in Linux tool). **License:** GPL-3.0.

**Tool 2: rclone for cloud storage → S3 backup**

**Why:** Native S3 API support without mounting, built-in encryption via crypt remotes, multi-threaded transfers, 70+ cloud provider support including all S3-compatible services.

**Implementation with Backblaze B2:**
```bash
# Configure encrypted remote
rclone config  # Set up B2 provider, then create crypt remote on top

# Daily encrypted backup at 3 AM via cron: 0 3 * * *
rclone sync /var/www/nextcloud/data/ \
  b2-encrypted:nextcloud-backup \
  --fast-list --checksum --transfers 32 \
  --log-file=/var/log/rclone-backup.log
```

**Performance tuning:**
- `--fast-list`: 25-30% improvement for large directories
- `--transfers 32`: Parallel uploads for S3
- `--checksum`: Verify integrity rather than size/time only

**Docker Compose deployment:**
```yaml
version: '3'
services:
  rclone-backup:
    image: rclone/rclone:latest
    volumes:
      - /data:/data:ro
      - ./rclone.conf:/config/rclone/rclone.conf
    command: >
      sync /data/ b2-encrypted:nextcloud-backup
      --fast-list --checksum --log-file=/logs/rclone.log
    restart: "no"
```

**Cost:** $0-10/month (Backblaze B2 storage only, no egress fees). **License:** MIT.

**Tool 3: Syncthing for real-time bidirectional sync (optional)**

**Why:** Continuous synchronization with automatic conflict resolution, mesh or hub-spoke topology, useful for active-active scenarios.

**When NOT to use:** Syncthing synchronizes changes immediately including deletions and corruptions - this is NOT backup. Deploy only if you need real-time bidirectional file sharing between locations, not for backup purposes.

**If needed for sync (not backup):**
```yaml
version: "3"
services:
  syncthing:
    image: syncthing/syncthing:latest
    hostname: cloud-vm
    environment:
      - PUID=1000
      - PGID=1000
    volumes:
      - /data/sync:/var/syncthing/data
      - ./config:/var/syncthing/config
    network_mode: host
    restart: unless-stopped
```

**Cost:** $0 (self-hosted). **License:** MPL v2.

**Recommended backup architecture:**

```
Cloud Storage (OCIS/Nextcloud)
    ├─ rsync → NAS (fast, on-site, 2 AM daily)
    └─ rclone → Backblaze B2 (encrypted, off-site, 3 AM daily)

Implements 3-2-1 Rule:
- 3 copies: production + NAS + B2
- 2 media types: local SSD + remote object storage  
- 1 off-site: Backblaze B2
```

**Critical success factors:**
- Test restores monthly (verify backup integrity)
- Monitor backup jobs (set up alerts for failures)
- Document restore procedures
- Automate everything via cron
- Keep rclone.conf backed up separately

**Total cost:** $0-10/month. **Setup time:** 2-3 hours for full implementation.

## Decision 6: Letta self-hosted with Ollama enables offline operation

**For farm operations requiring offline capability, self-hosted Letta with local LLMs is the only viable option.** Letta Cloud requires continuous internet connectivity and stores data on Letta's infrastructure - complete dealbreakers for your requirements regardless of other advantages.

**Self-hosted Letta architecture:**

**Docker deployment:**
```bash
# Single command deployment
docker run \
  -v ~/.letta/.persist/pgdata:/var/lib/postgresql/data \
  -p 8283:8283 \
  --network host \
  -e OLLAMA_BASE_URL="http://localhost:11434" \
  letta/letta:latest

# Includes bundled PostgreSQL for agent state storage
# REST API exposed on port 8283
```

**Offline capability with local LLMs:**

**Step 1 - Install Ollama:**
```bash
curl -fsSL https://ollama.com/install.sh | sh
```

**Step 2 - Download recommended models:**
```bash
ollama pull qwen3-14b        # 14B params - best balance (16GB VRAM)
ollama pull llama3.1-8b      # 8B params - lighter (8GB VRAM)
```

**Step 3 - Configure Letta to use local models:**
Letta automatically detects Ollama at `http://localhost:11434` when running with `--network host`. No additional configuration needed.

**100% offline operation achieved:** Letta server + Ollama + downloaded models + local Postgres = complete independence from internet connectivity after initial setup.

**Resource requirements:**
- **Letta Server**: 2-4GB RAM, 1-2 CPU cores (lightweight)
- **7B Models**: 8GB+ VRAM (GPU) or 16GB+ RAM (CPU-only, slower)
- **14B Models**: 16GB+ VRAM recommended
- **Recommended hardware**: Workstation with 32GB RAM + NVIDIA GPU (12-16GB VRAM)

**Integration architecture for your stack:**

**1. R2 Auto RAG Integration:**
```python
# Custom tool for Obsidian vault search
def search_knowledge_base(query: str) -> str:
    """Search the Obsidian vault using R2 AI Search."""
    import requests
    response = requests.post(
        f"https://{ACCOUNT}.r2.cloudflarestorage.com/ai-search",
        json={"query": query},
        headers={"Authorization": f"Bearer {API_KEY}"}
    )
    return response.json()["answer"]
```

**2. Postgres/pgvector Integration:**
```python
# Custom tool for structured data queries
def search_structured_data(query: str) -> str:
    """Search pgvector database for structured data."""
    import psycopg2
    from pgvector.psycopg2 import register_vector
    conn = psycopg2.connect(SUPABASE_CONNECTION_STRING)
    register_vector(conn)
    cur = conn.cursor()
    # Perform vector similarity search
    cur.execute(
        "SELECT content FROM documents ORDER BY embedding <=> %s LIMIT 5",
        (query_embedding,)
    )
    return cur.fetchall()
```

**3. Home Assistant Integration:**
- Method A: REST API calls via custom Python tools
- Method B: Community-built MCP servers for HA automation
- Method C: Via n8n workflows triggering Letta agents

**4. n8n Integration:**
Official Letta node available (`letta-ai/n8n-nodes-letta`). Send messages to agents from workflows, receive responses. Perfect for triggering agents from automated workflows.

**5. NodeRED Integration:**
HTTP Request node to Letta REST API:
```
POST http://localhost:8283/v1/agents/{agent_id}/messages/create
Authorization: Bearer {token}
Body: {"message": "Your message here"}
```

**6. Farm-specific integrations:**
```python
# Custom tool for farm sensors
def query_farm_sensor(sensor_id: str, metric: str) -> str:
    """Query farm sensor data from IoT devices."""
    import requests
    response = requests.get(
        f"http://farmbot-api/sensors/{sensor_id}/{metric}"
    )
    return response.json()

# Custom tool for EMQX MQTT
def publish_mqtt_message(topic: str, payload: str) -> str:
    """Publish message to EMQX MQTT broker."""
    import paho.mqtt.client as mqtt
    client = mqtt.Client()
    client.connect("localhost", 1883)
    client.publish(topic, payload)
    return "Message published successfully"
```

**Cost comparison over 5 years:**

| Option | Monthly Cost | 5-Year Total | Notes |
|--------|--------------|--------------|-------|
| **Letta Cloud** | $50-100 | $3,000-6,000 | Requires internet, no data sovereignty |
| **Self-hosted + Cloud APIs** | $20-50 | $1,200-3,000 | OpenAI/Anthropic token costs |
| **Self-hosted + Local Models** | $10-20 | $600-1,200 | Electricity only, full offline |

**Savings with self-hosted + local models: $2,400-5,400 over 5 years** while gaining offline capability and complete data control.

**Licensing:** Apache 2.0 (Letta core), PostgreSQL License (database), MIT (Ollama) - all permissive with zero commercial restrictions for farm operations.

## Decision 7: Letta Python SDK dominates API development

**The Letta Python SDK provides the most powerful and flexible tool development approach** for building custom integrations with your farm infrastructure.

**What "building custom tools" means:** Python functions that Letta agents can call to take actions - querying sensors, controlling equipment, accessing data sources, sending notifications. Defined with proper Google Style docstrings, automatically converted to OpenAI function calling schema.

**Letta Python SDK approach:**
```python
from letta_client import Letta

# Initialize client
client = Letta(base_url="http://localhost:8283")

# Define custom tool
def control_irrigation(zone: str, duration: int) -> str:
    """
    Control farm irrigation system.
    
    Args:
        zone (str): Irrigation zone identifier (A1-A4, B1-B4)
        duration (int): Duration in minutes
        
    Returns:
        str: Confirmation message
    """
    import requests
    response = requests.post(
        "http://farmbot-api/irrigation/activate",
        json={"zone": zone, "duration": duration}
    )
    return f"Activated zone {zone} for {duration} minutes"

# Register tool
tool = client.tools.create_from_function(control_irrigation)

# Create agent with tools
agent = client.agents.create(
    model="ollama/qwen3-14b",
    tools=["control_irrigation", "query_farm_sensor", "send_message"],
    memory_blocks=[
        {"label": "farm_context", "value": "Farm operations AI assistant"},
        {"label": "current_season", "value": "Winter wheat growing season"}
    ]
)
```

**Best practices for tool development:**
- All imports inside function scope (sandbox requirement)
- Proper Google Style docstrings (generates function schema)
- Handle errors gracefully with try/except
- Return JSON strings for complex data
- Test tools individually before agent deployment

**Alternative approaches evaluated:**

**Postman:** Useful for testing Letta's REST API endpoints but doesn't build tools. Use for debugging API authentication and exploring endpoint responses, not primary development.

**unified.io / Merge.dev:** Unified API platforms aggregating multiple SaaS APIs (CRM, HR, accounting). Not Letta-specific tools. Could serve as data sources that Letta tools query, but not relevant for building Letta tools themselves.

**MCP (Model Context Protocol):** Pre-built tool servers for common integrations. Use when available (Postgres, GitHub, Slack MCP servers exist), build custom tools when needed for farm-specific logic.

**Composio:** Library of 7,000+ pre-built tools with Letta integration. Useful for standard API integrations but not farm-specific operations.

**Recommendation:** Primary development with Letta Python SDK for custom farm tools, supplementary use of MCP servers for common needs (database queries, file operations), Postman for API testing only.

## Decision 8: Meilisearch complements rather than competes

**Meilisearch serves fundamentally different purposes than R2 Auto RAG + pgvector** - it's not redundant, it's complementary for specific use cases.

**Meilisearch design goals:**
- Lightning-fast keyword search (sub-50ms response times)
- Search-as-you-type autocomplete experiences  
- Typo tolerance and fuzzy matching
- Faceted filtering (e-commerce style)
- Hybrid search combining keywords + vectors

**R2 AI Search design goals:**
- Semantic document understanding
- LLM-generated comprehensive answers
- Vector similarity across unstructured data
- Slower but more intelligent responses

**When to use each:**

**Meilisearch scenarios:**
- Search bar in farm management dashboard
- Product catalog browsing
- Quick lookup by exact identifiers (SKU, sensor ID)
- Navigation and filtering interfaces

**R2 AI Search scenarios:**
- "What did I note about wheat fertilization last March?"
- Complex questions requiring document understanding
- Agent memory retrieval across knowledge base
- Semantic search where keywords aren't known

**Deployment:** Single Rust binary, minimal configuration, ~2GB RAM for production. Docker Compose setup takes 10 minutes.

**Licensing:** MIT (Community Edition) - fully open source, zero commercial restrictions. Enterprise Edition uses BUSL for advanced features like sharding.

**Verdict:** Deploy Meilisearch for fast UI search experiences, use R2 AI Search for agent knowledge retrieval. They solve different problems in your architecture.

## Decision 9: Homarr simplifies dashboard management

**Homarr over Dashy** for your distributed cloud VM and home server deployment because GUI-based configuration eliminates YAML editing friction and native Docker auto-discovery reduces manual maintenance.

**Homarr advantages:**
- Native Docker integration via socket mount (automatic container detection)
- Drag-and-drop UI configuration (zero YAML/JSON required)
- Modern widget system with 30+ integrations
- Multi-user with RBAC support
- Real-time WebSocket updates

**Configuration syncing limitation:** Both Homarr and Dashy require manual file copying to sync configs between servers. Dashy offers encrypted cloud backup feature, but manual sync remains pragmatic for two-server deployment.

**Deployment:**
```yaml
version: '3'
services:
  homarr:
    image: ghcr.io/ajnart/homarr:latest
    container_name: homarr
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock  # Auto-discovery
      - ./homarr/configs:/app/data/configs
      - ./homarr/icons:/app/public/icons
    ports:
      - '7575:7575'
    restart: unless-stopped
```

**Alternative - Dashy if:** You prefer YAML-as-code workflows, want built-in cloud backup for config syncing, or need extensive CSS customization options.

**Licensing:** Homarr (Apache 2.0), Dashy (MIT) - both permissive and actively maintained.

## Decision 10: Add Dockge alongside Portainer for complementary strengths

**Don't replace Portainer - supplement it with Dockge** because they serve complementary purposes with minimal overlap.

**Dockge excels at:**
- Stack-centric management (one folder per stack)
- Readable file structure (`/opt/stacks/[name]/compose.yml`)
- Full Docker Compose CLI compatibility
- Single-page interface (logs, editor, terminal unified)
- Convert docker run to compose.yml
- No vendor lock-in (standard compose files)

**Portainer excels at:**
- Advanced network management
- Individual container operations (not just stacks)
- Volume management UI
- Event logs and audit trails
- Registry management
- Kubernetes support
- Teams and RBAC (Enterprise)

**Recommended workflow:**
- **Daily stack management:** Use Dockge (faster, cleaner interface)
- **Network/volume operations:** Use Portainer (handles advanced features)
- **Individual containers:** Use Portainer (Dockge is stack-only)

**Deployment:**
```yaml
version: '3'
services:
  dockge:
    image: louislam/dockge:1
    container_name: dockge
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
      - ./dockge/data:/app/data
      - /opt/stacks:/opt/stacks  # Stack storage location
    environment:
      - DOCKGE_STACKS_DIR=/opt/stacks
    ports:
      - 5001:5001
    restart: unless-stopped
```

**Both use MIT licensing** with active development. Total setup time: 10 minutes for Dockge addition.

## Decision 11: Beszel with external monitoring solves "who monitors the monitor"

**Beszel provides lightweight infrastructure monitoring** (approximately 50MB RAM versus 500MB+ for Grafana + Prometheus stack) with built-in Docker container stats, simple alerting, and 5-minute setup versus several hours for traditional observability stacks.

**Beszel architecture:**
- Central hub on home server
- Agents on each monitored server (cloud VM, home server)
- WebSocket connections with SSH tunnel fallback
- Monitors: CPU, memory, disk, network, temperature, GPU, per-container stats
- Alerts: SMTP, Discord, Slack, Telegram via Shoutrrr

**Deployment:**
```yaml
# Hub on home server
version: '3'
services:
  beszel:
    image: henrygd/beszel:latest
    container_name: beszel
    volumes:
      - ./beszel/data:/beszel/data
    ports:
      - 8090:8090
    restart: unless-stopped

# Agent on each server
services:
  beszel-agent:
    image: henrygd/beszel-agent:latest
    container_name: beszel-agent
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro
    environment:
      - KEY=your-agent-key
      - HUB=https://beszel.yourdomain.com
    restart: unless-stopped
```

**"Who monitors the monitor" solution:**

**The problem:** If cloud VM goes down, Beszel hub goes down too. Need external monitoring independent of your infrastructure.

**Multi-layer monitoring strategy:**

**Tier 1 - Internal monitoring (Beszel):**
- Beszel hub on home server
- Agents on cloud VM and home server
- Monitors resources and containers
- Alerts via SMTP/Discord/Telegram

**Tier 2 - External availability monitoring (Healthchecks.io):**
- Services ping Healthchecks every 5-10 minutes (dead man's switch)
- If services stop reporting, alerts sent
- Self-host on separate cheap VPS or use cloud ($7/month)
- Independent of your infrastructure

**Tier 3 - Public endpoint monitoring (UptimeRobot):**
- Free tier: 50 monitors at 5-minute intervals
- Monitors HTTP endpoints from external network
- Public status pages
- Redundant alert channel

**Recommended architecture:**
```
Home Server (Beszel Hub)
    ↓ monitors
Cloud VM (Beszel Agent)
    ↓
Both ping Healthchecks.io (external VPS)
    ↓
UptimeRobot monitors public endpoints

Alerts to: Email + Telegram + Discord
```

**When to use Grafana instead:** Deploy full Grafana + Prometheus when you need complex PromQL queries, custom application metrics, multi-year retention, or compliance requirements. For infrastructure and container monitoring at your scale, Beszel suffices.

**Grafana Cloud free tier option:** 10K metrics, 50GB logs, 50GB traces/month hosted by Grafana Labs. Solves the "monitoring goes down" problem by pushing metrics to external service. Good compromise if you want Grafana without self-hosting.

**Cost:** $0-15/month (Healthchecks cloud or small monitoring VPS). **Licensing:** Beszel (MIT), Healthchecks.io (BSD), Uptime Kuma (MIT).

## Decision 12: Skip Kestra given existing n8n and NodeRED

**Kestra provides data pipeline orchestration for complex ETL workflows** - a use case you don't currently have with existing n8n and NodeRED handling your automation needs.

**Understanding the distinctions:**

**Kestra:** Infrastructure-as-Code workflow orchestrator (YAML-based) for data engineering tasks - ETL/ELT pipelines, scheduled batch processing, data transformations, infrastructure automation. Target users: data engineers, DevOps engineers.

**n8n:** Visual workflow automation (low-code) for API integrations and business process automation. Target users: business users, citizen developers, non-technical teams.

**NodeRED:** Flow-based programming for IoT and event-driven systems - MQTT, hardware devices, real-time data flows. Target users: IoT developers, makers, embedded systems.

**Current coverage analysis:**
- API integrations between services: **n8n handles this**
- Home automation flows: **NodeRED handles this**  
- IoT device coordination: **NodeRED handles this**
- MQTT message handling: **NodeRED handles this**
- Complex data transformations: **Not needed yet**
- Scheduled batch processing of large datasets: **Not needed yet**

**When to add Kestra:**
- Building data warehousing workflows
- Complex multi-step ETL pipelines requiring versioning
- Infrastructure automation (Ansible, Terraform orchestration)
- Scheduled batch jobs processing large datasets
- Need for declarative, Git-versioned workflow definitions

**Recommendation:** Skip Kestra initially. Your n8n + NodeRED combination covers current needs. Add Kestra later if you develop data engineering workloads that outgrow n8n's capabilities.

**If you later need orchestration for Airbyte:** Use Dagster or Prefect instead of Kestra - they integrate more naturally with modern data stacks. All three use Apache 2.0 licensing.

## Complete recommended architecture and deployment plan

### Core infrastructure stack

**Storage and knowledge layer:**
- **Cloudflare R2** - Object storage ($0-10/month)
- **OCIS** - WebDAV frontend (self-hosted, $0)
- **R2 AI Search** - Unstructured data RAG (free beta, ~$5/month projected)
- **Meilisearch** - Fast keyword search (self-hosted, $0)

**Database and structured data:**
- **Supabase self-hosted** - Postgres + pgvector + REST API ($20-40/month VPS)

**ETL and data movement:**
- **Airbyte** - All data ingestion ($50-200/month server infrastructure, $0 software)
- **Optional Dagster** - Orchestration when complexity grows (self-hosted, $0)

**Backup and resilience:**
- **rsync** - Cloud-to-NAS backup ($0)
- **rclone** - Cloud-to-S3 encrypted backup ($0-10/month storage)

**AI agent platform:**
- **Letta self-hosted** - Agent framework (self-hosted, $0)
- **Ollama** - Local LLM inference ($10-20/month electricity)
- **mem0** - Agent memory layer (self-hosted, $0)

**Infrastructure management:**
- **Homarr** - Service portal (self-hosted, $0)
- **Portainer** - Container management (self-hosted, $0)
- **Dockge** - Stack management (self-hosted, $0)
- **Beszel** - Internal monitoring (self-hosted, $0)
- **Healthchecks.io** - External monitoring ($0-7/month)
- **n8n** - Workflow automation (existing)
- **NodeRED** - IoT automation (existing)

**Total monthly cost: $90-270** depending on server sizing, $0 software licensing

### What you're definitively eliminating

**Vector databases (R2 Auto RAG + pgvector cover all needs):**
- Pinecone - Fully redundant with Vectorize
- Milvus - Overkill for your scale
- ChromaDB - Redundant with pgvector
- Qdrant - Covered by Vectorize/pgvector

**Document processing (R2 AI Search handles natively):**
- Docling - Parsing covered by Workers AI

**Database platforms (wrong architecture for your needs):**
- Cloudflare D1 - SQLite-based, insufficient capabilities
- Standalone Supabase pgvector - Integrate into main Postgres instead

**ETL platforms (not self-hosted or wrong category):**
- Fivetran - Cloud-only, no self-hosting
- Rivery - Cloud-only, no self-hosting
- Matillion - Limited self-hosting (cloud VMs only)
- Talend/Informatica - Enterprise-focused, sales required
- unstructured.io - Document preprocessing, not general ETL
- Confluent - Streaming platform, not batch ETL
- InfluxDB - Time-series database, not ETL tool

**Orchestration (existing tools sufficient):**
- Kestra - Redundant with n8n + NodeRED for current needs

**Advanced RAG (basic needs covered):**
- LlamaIndex - Only needed for complex retrieval patterns

### Deployment timeline and effort

**Week 1 - Infrastructure foundation (8 hours):**
- Deploy Portainer on both servers (1 hour)
- Add Dockge alongside Portainer (1 hour)
- Deploy Homarr on cloud VM (1 hour)
- Configure Tailscale networking (2 hours)
- Set up monitoring: Beszel hub + agents (2 hours)
- Configure external monitoring: Healthchecks.io (1 hour)

**Week 2 - Storage and database (10 hours):**
- Configure Cloudflare R2 bucket (1 hour)
- Deploy OCIS with R2 backend (2 hours)
- Configure Obsidian sync via WebDAV (1 hour)
- Deploy Supabase self-hosted (3 hours)
- Enable pgvector and test queries (1 hour)
- Set up R2 AI Search indexing (2 hours)

**Week 3 - ETL and backups (12 hours):**
- Deploy Airbyte via abctl (2 hours)
- Configure first connectors (Strava, Todoist) (3 hours)
- Set up dual destinations (Postgres + R2) (2 hours)
- Configure rsync backup to NAS (2 hours)
- Configure rclone backup to Backblaze B2 (2 hours)
- Test restore procedures (1 hour)

**Week 4 - AI agents (12 hours):**
- Install Ollama and download models (1 hour)
- Deploy Letta self-hosted (2 hours)
- Create first custom tools (farm sensors) (3 hours)
- Configure n8n integration (2 hours)
- Test agent interactions with R2 AI Search (2 hours)
- Test agent queries to pgvector (2 hours)

**Week 5 - Integration testing (8 hours):**
- Test Home Assistant integration (2 hours)
- Test NodeRED workflows with Letta (2 hours)
- End-to-end testing of ETL pipeline (2 hours)
- Disaster recovery testing (2 hours)

**Week 6 - Optimization and documentation (6 hours):**
- Performance tuning (2 hours)
- Security hardening (2 hours)
- Documentation of procedures (2 hours)

**Total deployment effort: 56 hours over 6 weeks**

### Service distribution across locations

**Cloud VM (primary public-facing services):**
- OCIS (WebDAV access to Obsidian vault)
- Letta server (agent platform)
- Airbyte (ETL processing)
- Supabase (database and API)
- n8n (workflow automation)
- Homarr (dashboard)
- Portainer + Dockge (container management)
- Beszel agent (monitoring)

**Home Server (backup and monitoring):**
- Beszel hub (monitors everything)
- Beszel agent (self-monitoring)
- NodeRED (IoT automation for home/farm)
- NAS backup target (rsync destination)
- Portainer + Dockge (local container management)
- Optional: Homarr (separate dashboard for local services)

**External/VPS (redundant monitoring):**
- Healthchecks.io (dead man's switch)
- OR Uptime Kuma (uptime monitoring)
- OR use SaaS: Healthchecks.io cloud ($7/month)

### Future scalability path

**Phase 1 (Immediate): Foundation**
- Single Letta server with Ollama
- 1-10 agents for different farm operations
- Basic integration with Home Assistant, n8n, NodeRED
- Obsidian vault as primary knowledge base

**Phase 2 (3-6 months): Expansion**
- Multi-agent system (weather, crops, equipment agents)
- Agents communicate via Letta's multi-agent API
- Scale to 50-100 agents for different farm zones
- Add more custom tools for equipment control

**Phase 3 (6-12 months): Robotics integration**
- Custom Python tools interfacing with ROS topics
- Agents control farm robots (Farmbot, autonomous equipment)
- Real-time sensor integration via EMQX MQTT
- Vision systems for crop monitoring

**Phase 4 (12+ months): AR and advanced AI**
- Brilliant Labs AR glasses integration
- Visual overlays on field views
- Real-time AI decision support
- Multi-modal agents (vision + language + sensor data)

**Phase 5 (Future): Distributed AI**
- Multiple Letta servers across farm infrastructure
- Load balancing with Docker Swarm or Kubernetes
- Multiple local GPUs for parallel inference
- Edge computing on farm equipment

### Critical success factors

**1. Offline operation validated:**
- After initial setup, system runs completely offline
- Local models (Ollama) eliminate internet dependency
- All data stored locally (Supabase, R2 mounted locally if needed)
- Critical for farm operations during connectivity outages

**2. Data sovereignty maintained:**
- Zero data leaves your infrastructure
- All services self-hosted
- No vendor lock-in (open source licenses)
- Complete audit trail of data access

**3. Cost predictability achieved:**
- Fixed infrastructure costs ($90-270/month)
- No usage-based charges
- No per-token pricing
- No surprise bills from cloud providers

**4. Integration flexibility preserved:**
- Standard protocols (REST, WebDAV, PostgreSQL, S3)
- Open source allows custom modifications
- Can replace any component without system rewrite
- Future-proof architecture for robotics and AR

**5. Minimal operational complexity:**
- Docker-based deployment (Portainer stacks)
- Automated backups (rsync, rclone on cron)
- External monitoring (Healthchecks alerts on failures)
- Self-updating containers (Watchtower)

### Licensing summary for all components

| Component | License | Commercial Use | Notes |
|-----------|---------|----------------|-------|
| **OCIS** | Apache 2.0 | ✅ Yes | Pre-built binaries have EULA |
| **R2 AI Search** | Proprietary managed | ✅ Yes | Usage-based pricing |
| **Meilisearch CE** | MIT | ✅ Yes | Enterprise features are BUSL |
| **Supabase** | Apache 2.0 | ✅ Yes | Includes multiple components |
| **PostgreSQL** | PostgreSQL License | ✅ Yes | Very permissive BSD-style |
| **pgvector** | PostgreSQL License | ✅ Yes | Same as Postgres |
| **Airbyte Core** | ELv2 | ✅ Yes | Cannot resell as SaaS |
| **Airbyte Connectors** | MIT | ✅ Yes | Fully permissive |
| **rsync** | GPL-3.0 | ✅ Yes | Copyleft, modifications must be open |
| **rclone** | MIT | ✅ Yes | Fully permissive |
| **Syncthing** | MPL v2 | ✅ Yes | File-level copyleft |
| **Letta** | Apache 2.0 | ✅ Yes | Fully permissive |
| **Ollama** | MIT | ✅ Yes | Fully permissive |
| **mem0** | Apache 2.0 | ✅ Yes | Fully permissive |
| **Homarr** | Apache 2.0 | ✅ Yes | Fully permissive |
| **Dashy** | MIT | ✅ Yes | Fully permissive |
| **Portainer CE** | Zlib | ✅ Yes | Permissive |
| **Dockge** | MIT | ✅ Yes | Fully permissive |
| **Beszel** | MIT | ✅ Yes | Fully permissive |
| **Grafana OSS** | AGPL 3.0 | ⚠️ Restrictions | Network copyleft applies |
| **Prometheus** | Apache 2.0 | ✅ Yes | Fully permissive |
| **n8n** | Sustainable Use | ⚠️ Fair-code | Self-hosting permitted |
| **NodeRED** | Apache 2.0 | ✅ Yes | Fully permissive |
| **Healthchecks.io** | BSD 3-Clause | ✅ Yes | Permissive |
| **Uptime Kuma** | MIT | ✅ Yes | Fully permissive |

**All components are free for commercial farm operations with zero licensing costs.**

This architecture delivers maximum capability with minimal complexity - eliminating redundant tools, maintaining data sovereignty, enabling complete offline operation, and providing predictable costs under $300/month while scaling from immediate deployment to future robotics and AR integration.
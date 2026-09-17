# TwelveLabs Video Intelligence Platform — Full Intelligence Profile

**Date Created:** 2025-06-15  
**Classification:** Research / Technical Intelligence  
**Domain:** 08 Research (The Lab) → ai-video  
**Status:** Active  
**Relevance to OneMind:** High — potential integration for security, robotics, and automated video analysis across the homestead fleet.

---

## Executive Summary

TwelveLabs is a video-first AI platform providing multimodal video understanding through two core models: **Marengo** (embedding/encoder) and **Pegasus** (generative/LLM). The platform treats video as a native input modality rather than a sequence of frames — ingesting visual, audio, speech, and on-screen text simultaneously for holistic scene comprehension.

**Key Differentiator:** While most vision models process frames in isolation, TwelveLabs models understand temporal and cross-modal relationships — akin to how humans watch, listen, and read simultaneously to comprehend context.

---

## 1. Core Models

### 1.1 Marengo — The Embedding Model

**Current Version:** Marengo 3.0

**Purpose:** Multimodal embeddings and semantic search across video, audio, image, and text.

| Capability | Description |
|------------|-------------|
| **Multimodal Processing** | Combines visual, audio, and text elements |
| **Fine-grained Search** | Detects logos, text, objects as small as 10% of frame |
| **Motion Search** | Identifies and analyzes movement within videos |
| **Counting** | Accurately counts objects in frames |
| **Audio Comprehension** | Analyzes music, lyrics, ambient sounds, silence |
| **Composed Search** | Query with text + image combined |
| **Cinematography** | Understands zoom, pan, tracking shots |
| **Sports Intelligence** | Soccer, basketball, baseball, ice hockey, football |

**Input Specifications:**

| Parameter | Specification |
|-----------|---------------|
| Video Duration | 4 seconds to 4 hours |
| Video File Size | ≤ 4 GB |
| Video Resolution | 360x360 to 5184x2160 |
| Aspect Ratio | 1:1 to 1:2.4 or 2.4:1 to 1:1 |
| Image Formats | JPEG, PNG (min 128x128, max 5 MB) |
| Audio Formats | WAV, MP3, FLAC (max 4 GB, 4 hours) |
| Max Text Length | 500 tokens |

**Supported Languages (36 + English):**  
Arabic, Bengali, Chinese (Simplified), Croatian, Cusco, Czech, Danish, Dutch, English, Farsi, Filipino, Finnish, French, German, Greek, Hebrew, Hindi, Hungarian, Indonesian, Italian, Japanese, Korean, Maori, Norwegian, Polish, Portuguese, Romanian, Russian, Spanish, Swahili, Swedish, Telugu, Thai, Turkish, Ukrainian, Vietnamese

**Embedding Dimensions:** 512 (optimized for speed and storage)

**Use Cases:**
- Cross-modal semantic search (text → video, image → video, video → video)
- Similarity matching and recommendation systems
- Integration with vector databases (Pinecone, Weaviate, Milvus, etc.)

---

### 1.2 Pegasus — The Generative Model

**Current Version:** Pegasus 1.5 (Pegasus 1.2 remains available)

**Purpose:** Video-to-text generation, analysis, and structured extraction.

**New in Pegasus 1.5:**
- **Video Segmentation:** Transform raw video into structured, timestamped data
- **Multimodal Prompting:** Include up to 4 reference images in prompts
- **Video Clipping:** Analyze specific time ranges
- **Context Window:** 261,120 tokens (shared between input/output)
- **Max Response Length:** 98,304 tokens (up from 4,096)
- **No Pre-indexing Required:** Analyze directly from URL/asset/base64

**Input Specifications:**

| Parameter | Specification |
|-----------|---------------|
| Video Duration | 4 seconds to 2 hours |
| Video File Size | ≤ 2 GB |
| Resolution/Format | Same as Marengo |

**Supported Languages:**
- Full Support: English
- Partial Support: Arabic, Chinese, French, German, Italian, Japanese, Korean, Portuguese, Russian, Spanish, Thai, Vietnamese

**Use Cases:**
- Video summarization and chapter generation
- Content analysis and action item extraction
- Timestamp identification for specific events
- Image-guided analysis ("find moments like this reference image")
- Structured metadata extraction via video segmentation

---

## 2. Platform Capabilities

### 2.1 Search API

**Cross-Modal Capabilities:**
- Text → Video: Natural language queries against indexed videos
- Image → Video: Visual similarity search
- Text + Image → Video: Composed queries combining both modalities
- Audio → Video: Search by sound/music patterns

**Search Options:**
- `visual`, `audio` modalities
- Threshold tuning (low/medium/high)
- Grouping by video or flat results
- Filtering by metadata

### 2.2 Analyze API

**Synchronous:** Up to 1 hour video, immediate response  
**Asynchronous:** Up to 2 hours video, webhook/poll for completion

**Analysis Modes:**
- General analysis (prompt-based text generation)
- Structured JSON responses (schema-defined)
- **Video Segmentation (Pegasus 1.5):** Define custom segment types and extract timestamped data

### 2.3 Embed API

**Video/Audio Embeddings:**
- Sync: Up to 10 minutes
- Async: 10 minutes to 4 hours

**Image Embeddings:**
- Single images up to 5 MB
- Batch text + image combined embeddings (up to 10 images)

**Text Embeddings:**
- Up to 500 tokens

---

## 3. Access Methods

### 3.1 TwelveLabs Cloud (SaaS)

**API Endpoint:** `https://api.twelvelabs.io/v1.3`

**Authentication:** API Key (x-api-key header)

**SDKs:**
- Python: `pip install twelvelabs`
- Node.js: `npm install twelvelabs-js`

**Sign Up:** https://playground.twelvelabs.io/

### 3.2 Amazon Bedrock (Enterprise)

Both Marengo and Pegasus available on AWS Bedrock for enterprise deployment.

**Benefits:**
- No separate TwelveLabs account needed
- AWS-native security and compliance
- Private VPC deployment options

### 3.3 MCP Server Integration

**MCP Endpoint:** `https://docs.twelvelabs.io/_mcp/server`

**Compatible Clients:**
- Claude Code
- Cursor
- VS Code with MCP
- Gemini CLI
- Goose

**Claude Code Plugin:**
```bash
/plugin marketplace add twelvelabs-io/twelve-labs-claude-code-plugin
/plugin install twelvelabs@twelvelabs-plugins
```

---

## 4. Pricing & Rate Limits

### Plans

| Plan | Cost | Use Case |
|------|------|----------|
| **Free** | $0 | Evaluation, prototyping |
| **Developer Tier 1** | Default with payment method | Production workloads |
| **Developer Tier 2** | $200/month | Higher throughput |
| **Developer Tier 3** | $400/month | Maximum scale |
| **Enterprise** | Custom | Dedicated resources, SLA, custom models |

### Free Plan Limits (Representative)

| Category | DPD | DPH | RPD | RPM | TPD (k) | TPM (k) |
|----------|-----|-----|-----|-----|---------|---------|
| Index | 3,000 min | 600 min | 3,000 | 60 | — | — |
| Search | — | — | 3,000 | 600 | — | — |
| Analyze | 3,000 min | 600 min | 1,000 | 60 | 500 | 30 |
| Embed (Video) | 3,000 min | 600 min | 3,000 | 25 | — | — |

**DPD** = Duration per day; **DPH** = Duration per hour  
**RPD** = Requests per day; **RPM** = Requests per minute  
**TPD** = Tokens per day (thousands); **TPM** = Tokens per minute (thousands)

---

## 5. Partner Integrations

### VideoDB (Real-Time Video Understanding)

**Integration Type:** Streaming video monitoring for live feeds

**How it Works:**
1. Connect RTSP/IP camera via VideoDB RTStream
2. TwelveLabs Pegasus indexes scenes in real-time
3. Define surveillance events in natural language
4. Receive webhook alerts on event detection
5. Access flagged segments via HLS manifests

**Use Cases:**
- Security breach detection
- Safety hazard alerts
- Emergency response triggering
- Flash flood detection
- Baby crib monitoring

**Reference:** https://github.com/twelvelabs-io/twelvelabs-developer-experience/blob/main/integrations/VideoDB/

### Other Notable Integrations

| Partner | Integration |
|---------|-------------|
| **Adobe Premiere Pro** | Search videos, auto-generate chapters in Premiere |
| **Pinecone** | Semantic video search with vector storage |
| **Weaviate** | Multimodal RAG for video Q&A |
| **Chroma** | Video embeddings for retrieval systems |
| **ApertureDB** | Video + metadata unified storage |
| **Backblaze B2** | Media asset management |
| **Databricks** | Enterprise video analytics pipelines |
| **LanceDB** | Edge-friendly vector database |
| **Langflow** | Visual workflow builder for video AI |
| **MindsDB** | SQL-native video analysis |
| **Milvus** | Large-scale semantic video search |
| **MongoDB Atlas** | Vector search integration |
| **n8n** | No-code automation workflows |

---

## 6. Fine-Tuning

**Status:** Available for selected customers only (contact sales@twelvelabs.io)

**Process:**
1. Prepare domain-specific training dataset
2. TwelveLabs runs automated fine-tuning pipeline
3. Deploy custom model (days not weeks)

**Best Suited For:**
- Domain-specific object/action detection
- Custom taxonomies and vocabularies
- Industry-specific terminology

**Limitations:**
- Small/background objects may not improve significantly
- Long-term temporal actions (e.g., 5-minute run) less effective
- Generalization not guaranteed across contexts

---

## 7. Webhooks & Automation

**Supported Events:**
- Video indexing task completion/failure
- Async analysis task completion/failure

**Workflow:**
1. Create HTTP endpoint in your application
2. Register webhook in TwelveLabs Dashboard
3. Receive real-time JSON POST notifications
4. Process payload for downstream automation

**Note:** Webhooks unavailable for Embed API.

**Integration with OneMind:** Could trigger NATS messages for event-driven architecture across the fleet.

---

## 8. Relevance to OneMind OS

### Immediate Applications

| Use Case | TwelveLabs Capability |
|----------|----------------------|
| **Security Camera Analysis** | Real-time event detection via VideoDB integration |
| **Robot Vision (DIMOS, OM1)** | Video understanding for navigation and task completion |
| **Drone Footage Analysis** | Automated indexing and search of aerial captures |
| **Training Content** | Summarization and chapter extraction for education platform |
| **Meeting Analysis** | Automatic action item extraction from recorded calls |
| **Media Organization** | Semantic search across video library |

### Technical Fit

| Requirement | TwelveLabs Support |
|-------------|-------------------|
| Local/On-Prem Deployment | Available via Enterprise |
| AWS Integration | Native Bedrock support |
| Multi-modal Input | Video, Audio, Image, Text |
| Real-Time Processing | Via VideoDB partnership |
| Webhook Integration | Native support |
| MCP Server | Native support for Claude/Cursor |

### Potential Integration Patterns

```
Camera Feed → VideoDB RTStream → TwelveLabs Pegasus → Webhook → OneMind NATS → Agent Action
                                                                                ↓
                                                                      Alert/Dashboard/Automation
```

---

## 9. Resources

### Documentation
- **Main Docs:** https://docs.twelvelabs.io/
- **v1.3 LLMs.txt:** https://docs.twelvelabs.io/v1.3/llms.txt
- **Full Docs Corpus:** https://docs.twelvelabs.io/llms-full.txt
- **API Reference:** https://docs.twelvelabs.io/reference/api-reference

### Developer Resources
- **Python SDK:** https://pypi.org/project/twelvelabs/
- **Node.js SDK:** https://www.npmjs.com/package/twelvelabs-js
- **GitHub Samples:** https://github.com/twelvelabs-io/twelvelabs-python
- **Playground:** https://playground.twelvelabs.io/

### Support
- **Email:** support@twelvelabs.io
- **Sales:** james@twelvelabs.io, sales@twelvelabs.io
- **Discord:** https://discord.com/invite/Sh6BRfakJa

---

## 10. Quick Reference: Model Selection

| Task | Model | Version | Notes |
|------|-------|---------|-------|
| Semantic search | Marengo | 3.0 | Pre-index videos first |
| Embeddings | Marengo | 3.0 | 512-dim vectors |
| Video summarization | Pegasus | 1.5 | Async for >1hr |
| Timestamp extraction | Pegasus | 1.5 | Temporal grounding |
| Structured segmentation | Pegasus | 1.5 | Define schema |
| Image-guided analysis | Pegasus | 1.5 | Include reference imgs |
| Real-time monitoring | Pegasus | 1.2+ | Via VideoDB |

---

## 11. Assessment

**Strengths:**
- Purpose-built for video (not bolted-on vision)
- True multimodal understanding (not siloed pipelines)
- Strong temporal reasoning
- Mature SDK and documentation
- Amazon Bedrock availability for enterprise
- MCP server for AI assistant integration

**Considerations:**
- Cloud-first (enterprise required for on-prem)
- Pricing scales with video duration processed
- Fine-tuning requires sales engagement
- Limited language support for generation (vs. search)

**Verdict:** **High-value candidate** for OneMind security, robotics, and content workflows. Recommend prototyping with Free tier, especially VideoDB integration for homestead camera feeds.

---

**Last Updated:** 2025-06-15  
**Next Review:** As needed based on integration progress

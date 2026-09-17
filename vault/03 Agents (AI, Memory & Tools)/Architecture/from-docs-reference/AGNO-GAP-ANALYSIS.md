# Agno Gap Analysis - OneMind OS

> **Purpose:** Detailed comparison of Agno features vs OneMind OS implementation
> **Date:** 2026-01-26
> **Reference:** [AGNO-COMPLETE-FEATURE-INVENTORY.md](./AGNO-COMPLETE-FEATURE-INVENTORY.md)

---

## Executive Summary

| Metric | Value |
|--------|-------|
| **Total Agno Features** | 186 |
| **Implemented in OneMind** | 118 (63%) |
| **Fully Working** | 85 (46%) |
| **Scaffolded/Partial** | 33 (18%) |
| **Not Implemented** | 68 (37%) |
| **Critical Gaps (P0)** | 8 |
| **High Priority Gaps (P1)** | 22 |
| **Medium Priority Gaps (P2)** | 25 |
| **Low Priority Gaps (P3)** | 13 |

---

## 1. Agent Features Gap Analysis

### 1.1 Agent Core Parameters (50 total)

| Parameter | Agno | OneMind | Status | Priority |
|-----------|:----:|:-------:|--------|----------|
| `id` | ✅ | ✅ | **Working** | - |
| `name` | ✅ | ✅ | **Working** | - |
| `model` | ✅ | ✅ | **Working** | - |
| `description` | ✅ | ✅ | **Working** | - |
| `role` | ✅ | ✅ | **Working** | - |
| `system_message` | ✅ | ✅ | **Working** | - |
| `system_message_role` | ✅ | ⚠️ | Not used | P3 |
| `introduction` | ✅ | ⚠️ | Not used | P3 |
| `instructions` | ✅ | ✅ | **Working** | - |
| `expected_output` | ✅ | ⚠️ | Not used | P2 |
| `additional_context` | ✅ | ⚠️ | Not used | P2 |
| `build_context` | ✅ | ✅ | **Working** | - |
| `build_user_context` | ✅ | ⚠️ | Not used | P3 |
| `add_name_to_context` | ✅ | ⚠️ | Not used | P3 |
| `add_datetime_to_context` | ✅ | ✅ | **Working** | - |
| `add_location_to_context` | ✅ | ✅ | **Working** | - |
| `timezone_identifier` | ✅ | ⚠️ | Not used | P2 |
| `resolve_in_context` | ✅ | ⚠️ | Not used | P2 |
| `tools` | ✅ | ✅ | **Working** | - |
| `show_tool_calls` | ✅ | ⚠️ | Not used | P3 |
| `tool_choice` | ✅ | ⚠️ | Not used | P2 |
| `tool_call_limit` | ✅ | ⚠️ | Not used | P2 |
| `knowledge` | ✅ | ⚠️ | Scaffolded | **P0** |
| `search_knowledge` | ✅ | ⚠️ | Not wired | **P0** |
| `num_documents` | ✅ | ⚠️ | Not used | P1 |
| `db` | ✅ | ✅ | **Working** | - |
| `enable_agentic_memory` | ✅ | ✅ | **Working** | - |
| `enable_user_memories` | ✅ | ✅ | **Working** | - |
| `enable_session_summaries` | ✅ | ✅ | **Working** | - |
| `user_id` | ✅ | ✅ | **Working** | - |
| `session_id` | ✅ | ✅ | **Working** | - |
| `cache_session` | ✅ | ❌ | Not implemented | P2 |
| `session_state` | ✅ | ⚠️ | Not used | P1 |
| `add_session_state_to_context` | ✅ | ⚠️ | Not used | P1 |
| `enable_agentic_state` | ✅ | ⚠️ | Not used | P1 |
| `overwrite_db_session_state` | ✅ | ⚠️ | Not used | P2 |
| `add_history_to_context` | ✅ | ✅ | **Working** | - |
| `num_history_runs` | ✅ | ✅ | **Working** | - |
| `read_chat_history` | ✅ | ✅ | **Working** | - |
| `search_session_history` | ✅ | ⚠️ | Not used | P2 |
| `num_history_sessions` | ✅ | ⚠️ | Not used | P2 |
| `max_tool_calls_from_history` | ✅ | ⚠️ | Not used | P3 |
| `output_schema` | ✅ | ✅ | **Working** | - |
| `input_schema` | ✅ | ⚠️ | Not used | P2 |
| `structured_outputs` | ✅ | ✅ | **Working** | - |
| `parse_response` | ✅ | ✅ | **Working** | - |
| `markdown` | ✅ | ✅ | **Working** | - |
| `parser_model` | ✅ | ❌ | Not implemented | P3 |
| `parser_model_prompt` | ✅ | ❌ | Not implemented | P3 |
| `output_model` | ✅ | ❌ | Not implemented | P3 |
| `output_model_prompt` | ✅ | ❌ | Not implemented | P3 |

### 1.2 Agent Reasoning Parameters (5 total)

| Parameter | Agno | OneMind | Status | Priority |
|-----------|:----:|:-------:|--------|----------|
| `reasoning` | ✅ | ✅ | **Working** | - |
| `reasoning_model` | ✅ | ⚠️ | Not used | P1 |
| `reasoning_agent` | ✅ | ❌ | Not implemented | P2 |
| `reasoning_min_steps` | ✅ | ❌ | Not implemented | P2 |
| `reasoning_max_steps` | ✅ | ❌ | Not implemented | P2 |

### 1.3 Agent Retry Parameters (3 total)

| Parameter | Agno | OneMind | Status | Priority |
|-----------|:----:|:-------:|--------|----------|
| `retries` | ✅ | ⚠️ | Not used | P1 |
| `delay_between_retries` | ✅ | ⚠️ | Not used | P1 |
| `exponential_backoff` | ✅ | ⚠️ | Not used | P1 |

### 1.4 Agent Hooks & Dependencies (4 total)

| Parameter | Agno | OneMind | Status | Priority |
|-----------|:----:|:-------:|--------|----------|
| `pre_hooks` | ✅ | ✅ | **Working** | - |
| `post_hooks` | ✅ | ✅ | **Working** | - |
| `dependencies` | ✅ | ✅ | **Working** | - |
| `add_dependencies_to_context` | ✅ | ✅ | **Working** | - |

### 1.5 Agent Debug Parameters (2 total)

| Parameter | Agno | OneMind | Status | Priority |
|-----------|:----:|:-------:|--------|----------|
| `debug_mode` | ✅ | ✅ | **Working** | - |
| `telemetry` | ✅ | ✅ | **Working** | - |

**Agent Summary:**
- **Working:** 30/64 (47%)
- **Scaffolded/Not Used:** 26/64 (41%)
- **Not Implemented:** 8/64 (12%)

---

## 2. Team Features Gap Analysis

### 2.1 Team Parameters (15 total)

| Parameter | Agno | OneMind | Status | Priority |
|-----------|:----:|:-------:|--------|----------|
| `name` | ✅ | ✅ | **Working** | - |
| `id` | ✅ | ✅ | **Working** | - |
| `model` | ✅ | ⚠️ | Not configured | P1 |
| `members` | ✅ | ✅ | **Working** | - |
| `instructions` | ✅ | ✅ | **Working** | - |
| `description` | ✅ | ⚠️ | Not used | P2 |
| `tools` | ✅ | ⚠️ | Not used | P1 |
| `db` | ✅ | ✅ | **Working** | - |
| `user_id` | ✅ | ✅ | **Working** | - |
| `session_id` | ✅ | ✅ | **Working** | - |
| `show_members_responses` | ✅ | ⚠️ | Not used | P1 |
| `delegate_to_all_members` | ✅ | ❌ | Not implemented | **P0** |
| `enable_agentic_memory` | ✅ | ⚠️ | Not used | P2 |
| `add_history_to_context` | ✅ | ⚠️ | Not used | P2 |
| `markdown` | ✅ | ✅ | **Working** | - |

### 2.2 Team Coordination Modes (5 total)

| Mode | Agno | OneMind | Status | Priority |
|------|:----:|:-------:|--------|----------|
| Route (default) | ✅ | ⚠️ | Basic only | **P0** |
| Coordinate | ✅ | ❌ | Not implemented | **P0** |
| Collaborate | ✅ | ❌ | Not implemented | **P0** |
| Sequential | ✅ | ❌ | Not implemented | P1 |
| Parallel | ✅ | ❌ | Not implemented | P1 |

**Team Summary:**
- **Working:** 8/20 (40%)
- **Scaffolded/Not Used:** 6/20 (30%)
- **Not Implemented:** 6/20 (30%)

---

## 3. Workflow Features Gap Analysis

### 3.1 Workflow CRUD (7 total)

| Feature | Agno | OneMind | Status | Priority |
|---------|:----:|:-------:|--------|----------|
| `GET /workflows` | ✅ | ✅ | **Working** | - |
| `GET /workflows/{id}` | ✅ | ✅ | **Working** | - |
| `GET /workflows/step-types` | ✅ | ✅ | **Working** | - |
| `POST /workflows` | ✅ | ✅ | **Working** | - |
| `PUT /workflows/{id}` | ✅ | ✅ | **Working** | - |
| `DELETE /workflows/{id}` | ✅ | ✅ | **Working** | - |
| `POST /workflows/{id}/run` | ✅ | ✅ | **Working** | - |

### 3.2 Workflow Step Types (6 total)

| Step Type | Agno | OneMind | Status | Priority |
|-----------|:----:|:-------:|--------|----------|
| Step (basic) | ✅ | ✅ | **Working** | - |
| Router | ✅ | ⚠️ | Scaffolded | **P0** |
| Condition | ✅ | ⚠️ | Scaffolded | **P0** |
| Parallel | ✅ | ⚠️ | Scaffolded | **P0** |
| Loop | ✅ | ❌ | Not implemented | P1 |
| Function | ✅ | ⚠️ | Not used | P2 |

### 3.3 Workflow Features (6 total)

| Feature | Agno | OneMind | Status | Priority |
|---------|:----:|:-------:|--------|----------|
| Step input piping | ✅ | ⚠️ | Basic only | P1 |
| Step output handling | ✅ | ⚠️ | Basic only | P1 |
| Input caching | ✅ | ❌ | Not implemented | P2 |
| Error handling/fallback | ✅ | ⚠️ | Basic only | P1 |
| Session state persistence | ✅ | ⚠️ | Not used | P2 |
| HITL at workflow steps | ✅ | ⚠️ | Not tested | P1 |

**Workflow Summary:**
- **Working:** 8/19 (42%)
- **Scaffolded/Basic:** 8/19 (42%)
- **Not Implemented:** 3/19 (16%)

---

## 4. Built-in Tools Gap Analysis

### 4.1 Search & Research Tools (8 total)

| Tool | Agno | OneMind | Status | Priority |
|------|:----:|:-------:|--------|----------|
| DuckDuckGoTools | ✅ | ✅ | **Enabled** | - |
| WikipediaTools | ✅ | ✅ | **Enabled** | - |
| ArxivTools | ✅ | ✅ | **Enabled** | - |
| HackerNewsTools | ✅ | ✅ | **Enabled** | - |
| YouTubeTools | ✅ | ✅ | **Enabled** | - |
| ExaTools | ✅ | ❌ | Not added | P2 |
| TavilyTools | ✅ | ❌ | Not added | P2 |
| BraveSearchTools | ✅ | ❌ | Not added | P3 |

### 4.2 Finance Tools (1 total)

| Tool | Agno | OneMind | Status | Priority |
|------|:----:|:-------:|--------|----------|
| YFinanceTools | ✅ | ✅ | **Enabled** | - |

### 4.3 Code Execution Tools (5 total)

| Tool | Agno | OneMind | Status | Priority |
|------|:----:|:-------:|--------|----------|
| CalculatorTools | ✅ | ✅ | **Enabled** | - |
| ShellTools | ✅ | ❌ | Disabled (needs HITL) | P1 |
| PythonTools | ✅ | ❌ | Disabled (needs HITL) | P1 |
| E2BTools | ✅ | ❌ | Not added | P2 |
| SQLTools | ✅ | ❌ | Not added (needs HITL) | P1 |

### 4.4 Web/Crawling Tools (3 total)

| Tool | Agno | OneMind | Status | Priority |
|------|:----:|:-------:|--------|----------|
| FirecrawlTools | ✅ | ❌ | Not added | P2 |
| Newspaper4kTools | ✅ | ❌ | Not added | P2 |
| CrawlbaseTools | ✅ | ❌ | Not added | P3 |

### 4.5 Communication Tools (5 total)

| Tool | Agno | OneMind | Status | Priority |
|------|:----:|:-------:|--------|----------|
| SlackTools | ✅ | ❌ | Not added | P1 |
| DiscordTools | ✅ | ⚠️ | Custom impl | P2 |
| TelegramTools | ✅ | ⚠️ | Custom impl | P2 |
| EmailTools | ✅ | ❌ | Not added | P1 |
| TwilioTools | ✅ | ❌ | Not added | P2 |

### 4.6 File & Data Tools (5 total)

| Tool | Agno | OneMind | Status | Priority |
|------|:----:|:-------:|--------|----------|
| FileTools | ✅ | ❌ | Not added (needs HITL) | P1 |
| FileGenerationTools | ✅ | ❌ | Not added | P1 |
| CSVTools | ✅ | ❌ | Not added | P2 |
| JSONTools | ✅ | ❌ | Not added | P2 |
| PDFTools | ✅ | ❌ | Not added | P2 |

### 4.7 Integration Tools (5 total)

| Tool | Agno | OneMind | Status | Priority |
|------|:----:|:-------:|--------|----------|
| GitHubTools | ✅ | ✅ | Custom impl | - |
| LinearTools | ✅ | ✅ | Custom impl | - |
| TodoistTools | ✅ | ✅ | Custom impl | - |
| GoogleCalendarTools | ✅ | ❌ | Not added | P2 |
| NotionTools | ✅ | ❌ | Not added | P3 |

### 4.8 Reasoning & Memory Tools (2 total)

| Tool | Agno | OneMind | Status | Priority |
|------|:----:|:-------:|--------|----------|
| ReasoningTools | ✅ | ❌ | Not added | P1 |
| MemoryTools | ✅ | ✅ | Via agentic_memory | - |

### 4.9 MCP Tools (2 total)

| Tool | Agno | OneMind | Status | Priority |
|------|:----:|:-------:|--------|----------|
| MCPTools | ✅ | ⚠️ | Defined, not auto-attached | P1 |
| MultiMCPTools | ✅ | ❌ | Not used | P2 |

**Tools Summary:**
- **Enabled/Working:** 12/36 (33%)
- **Custom Implementation:** 5/36 (14%)
- **Not Added:** 19/36 (53%)

---

## 5. Memory System Gap Analysis

### 5.1 Memory CRUD (6 total)

| Feature | Agno | OneMind | Status | Priority |
|---------|:----:|:-------:|--------|----------|
| GET /memories | ✅ | ✅ | **Working** | - |
| POST /memories | ✅ | ✅ | **Working** | - |
| PUT /memories/{id} | ✅ | ✅ | **Working** | - |
| DELETE /memories/{id} | ✅ | ✅ | **Working** | - |
| Memory search (full-text) | ✅ | ✅ | **Working** | - |
| Memory semantic search | ✅ | ✅ | **Working** | - |

### 5.2 Memory Modes (4 total)

| Mode | Agno | OneMind | Status | Priority |
|------|:----:|:-------:|--------|----------|
| enable_agentic_memory | ✅ | ✅ | **Working** | - |
| enable_user_memories | ✅ | ✅ | **Working** | - |
| enable_session_summaries | ✅ | ✅ | **Working** | - |
| Topic classification | ✅ | ⚠️ | Not used | P2 |

**Memory Summary:**
- **Working:** 9/10 (90%)
- **Not Used:** 1/10 (10%)

---

## 6. Knowledge/RAG Gap Analysis

### 6.1 Knowledge CRUD (5 total)

| Feature | Agno | OneMind | Status | Priority |
|---------|:----:|:-------:|--------|----------|
| POST /knowledge/upload | ✅ | ✅ | **Working** | - |
| POST /knowledge/url | ✅ | ✅ | **Working** | - |
| GET /knowledge/search | ✅ | ✅ | **Working** | - |
| GET /knowledge/stats | ✅ | ✅ | **Working** | - |
| DELETE /knowledge/{id} | ✅ | ❌ | Not implemented | P2 |

### 6.2 Vector Databases (5 total)

| Database | Agno | OneMind | Status | Priority |
|----------|:----:|:-------:|--------|----------|
| LanceDb | ✅ | ✅ | **Working** | - |
| PgVector | ✅ | ✅ | **Working** | - |
| Pinecone | ✅ | ❌ | Not added | P3 |
| Qdrant | ✅ | ❌ | Not added | P3 |
| Chroma | ✅ | ❌ | Not added | P3 |

### 6.3 Search Types (4 total)

| Type | Agno | OneMind | Status | Priority |
|------|:----:|:-------:|--------|----------|
| Vector search | ✅ | ✅ | **Working** | - |
| Keyword search | ✅ | ✅ | **Working** | - |
| Hybrid search | ✅ | ✅ | **Working** | - |
| Reranking | ✅ | ✅ | **Working** | - |

### 6.4 Document Processing (8 total)

| Feature | Agno | OneMind | Status | Priority |
|---------|:----:|:-------:|--------|----------|
| PDF ingestion | ✅ | ✅ | **Working** | - |
| TXT/MD ingestion | ✅ | ✅ | **Working** | - |
| HTML ingestion | ✅ | ✅ | **Working** | - |
| DOCX ingestion | ✅ | ⚠️ | Not tested | P2 |
| CSV ingestion | ✅ | ⚠️ | Not tested | P2 |
| JSON ingestion | ✅ | ⚠️ | Not tested | P2 |
| Semantic chunking | ✅ | ⚠️ | Not used | P1 |
| Agentic chunking | ✅ | ❌ | Not implemented | P2 |

### 6.5 Agent Integration (2 total)

| Feature | Agno | OneMind | Status | Priority |
|---------|:----:|:-------:|--------|----------|
| Auto-search in agent | ✅ | ❌ | **Not wired** | **P0** |
| num_documents config | ✅ | ⚠️ | Not used | P1 |

**Knowledge Summary:**
- **Working:** 14/24 (58%)
- **Not Tested/Used:** 6/24 (25%)
- **Not Implemented:** 4/24 (17%)

---

## 7. MCP Gap Analysis

### 7.1 MCP Server Management (3 total)

| Feature | Agno | OneMind | Status | Priority |
|---------|:----:|:-------:|--------|----------|
| GET /mcp/servers | ✅ | ✅ | **Working** | - |
| POST /mcp/servers | ✅ | ✅ | **Working** | - |
| DELETE /mcp/servers/{id} | ✅ | ✅ | **Working** | - |

### 7.2 MCP Transports (3 total)

| Transport | Agno | OneMind | Status | Priority |
|-----------|:----:|:-------:|--------|----------|
| stdio | ✅ | ✅ | **Supported** | - |
| streamable-http | ✅ | ✅ | **Supported** | - |
| sse | ✅ | ⚠️ | Not tested | P2 |

### 7.3 Pre-configured Servers (10 total)

| Server | Agno | OneMind | Status | Priority |
|--------|:----:|:-------:|--------|----------|
| filesystem | ✅ | ✅ | **Defined** | - |
| memory | ✅ | ✅ | **Defined** | - |
| brave-search | ✅ | ✅ | **Defined** | - |
| github | ✅ | ✅ | **Defined** | - |
| postgres | ✅ | ✅ | **Defined** | - |
| puppeteer | ✅ | ✅ | **Defined** | - |
| context7 | ✅ | ✅ | **Defined** | - |
| slack | ✅ | ❌ | Not defined | P2 |
| google-maps | ✅ | ❌ | Not defined | P2 |
| sequential-thinking | ✅ | ❌ | Not defined | P2 |

### 7.4 MCP Agent Integration (3 total)

| Feature | Agno | OneMind | Status | Priority |
|---------|:----:|:-------:|--------|----------|
| Auto-attach to agent | ✅ | ❌ | **Manual only** | **P0** |
| MCPTools in agent | ✅ | ⚠️ | Available, not default | P1 |
| MultiMCPTools | ✅ | ❌ | Not used | P2 |

**MCP Summary:**
- **Working/Defined:** 13/19 (68%)
- **Not Tested/Used:** 2/19 (11%)
- **Not Implemented:** 4/19 (21%)

---

## 8. HITL Gap Analysis

### 8.1 Tool Confirmation (4 total)

| Feature | Agno | OneMind | Status | Priority |
|---------|:----:|:-------:|--------|----------|
| @tool(requires_confirmation) | ✅ | ✅ | **Working** | - |
| requirement.confirm() | ✅ | ✅ | **Working** | - |
| requirement.reject() | ✅ | ✅ | **Working** | - |
| agent.continue_run() | ✅ | ✅ | **Working** | - |

### 8.2 Approval Queue (5 total)

| Feature | Agno | OneMind | Status | Priority |
|---------|:----:|:-------:|--------|----------|
| GET /approvals | ✅ | ✅ | **Working** | - |
| POST /approvals/{id}/approve | ✅ | ✅ | **Working** | - |
| POST /approvals/{id}/reject | ✅ | ✅ | **Working** | - |
| Approval webhooks | ✅ | ❌ | Not implemented | P2 |
| Approval timeout | ✅ | ❌ | Not implemented | P2 |

**HITL Summary:**
- **Working:** 7/9 (78%)
- **Not Implemented:** 2/9 (22%)

---

## 9. Evaluation System Gap Analysis

### 9.1 Evaluation Types (4 total)

| Type | Agno | OneMind | Status | Priority |
|------|:----:|:-------:|--------|----------|
| AccuracyEval | ✅ | ⚠️ | Scaffolded | P1 |
| PerformanceEval | ✅ | ⚠️ | Scaffolded | P1 |
| ReliabilityEval | ✅ | ⚠️ | Scaffolded | P1 |
| Agent-as-judge | ✅ | ❌ | Not implemented | P2 |

### 9.2 Eval API (3 total)

| Feature | Agno | OneMind | Status | Priority |
|---------|:----:|:-------:|--------|----------|
| POST /eval-runs | ✅ | ⚠️ | Scaffolded | P1 |
| GET /eval-runs | ✅ | ⚠️ | Scaffolded | P1 |
| GET /eval-runs/{id} | ✅ | ⚠️ | Scaffolded | P1 |

**Evals Summary:**
- **Scaffolded:** 6/7 (86%)
- **Not Implemented:** 1/7 (14%)

---

## 10. Multimodal Gap Analysis

### 10.1 Input Modalities (6 total)

| Modality | Agno | OneMind | Status | Priority |
|----------|:----:|:-------:|--------|----------|
| Text input | ✅ | ✅ | **Working** | - |
| Image input | ✅ | ✅ | **Working** | - |
| Audio input | ✅ | ✅ | **Working** | - |
| Video input | ✅ | ⚠️ | Basic only | P2 |
| PDF input | ✅ | ✅ | **Working** | - |
| Document input | ✅ | ✅ | **Working** | - |

### 10.2 Output Modalities (4 total)

| Modality | Agno | OneMind | Status | Priority |
|----------|:----:|:-------:|--------|----------|
| Text output | ✅ | ✅ | **Working** | - |
| Audio output | ✅ | ❌ | Not implemented | P1 |
| Image generation | ✅ | ❌ | Not implemented | P2 |
| File generation | ✅ | ❌ | Not implemented | P1 |

**Multimodal Summary:**
- **Working:** 6/10 (60%)
- **Basic/Not Implemented:** 4/10 (40%)

---

## 11. Observability Gap Analysis

### 11.1 Tracing Backends (3 total)

| Backend | Agno | OneMind | Status | Priority |
|---------|:----:|:-------:|--------|----------|
| Arize Phoenix | ✅ | ✅ | Implemented, not deployed | P1 |
| Langfuse | ✅ | ✅ | Implemented, not deployed | P1 |
| OpenTelemetry | ✅ | ⚠️ | Basic only | P2 |

### 11.2 Debug/Telemetry (5 total)

| Feature | Agno | OneMind | Status | Priority |
|---------|:----:|:-------:|--------|----------|
| debug_mode toggle | ✅ | ✅ | **Working** | - |
| telemetry toggle | ✅ | ✅ | **Working** | - |
| Enhanced metrics | ✅ | ⚠️ | Basic only | P2 |
| Token metrics | ✅ | ✅ | **Working** | - |
| Time to first token | ✅ | ⚠️ | Not used | P3 |

**Observability Summary:**
- **Working:** 4/8 (50%)
- **Basic/Not Deployed:** 4/8 (50%)

---

## 12. Model Providers Gap Analysis

### 12.1 Native Providers (10 total)

| Provider | Agno | OneMind | Status | Priority |
|----------|:----:|:-------:|--------|----------|
| OpenAI | ✅ | ✅ | Via LiteLLM | - |
| Anthropic | ✅ | ✅ | Via Bedrock | - |
| Google/Gemini | ✅ | ✅ | Via LiteLLM | - |
| Groq | ✅ | ✅ | Via LiteLLM | - |
| DeepSeek | ✅ | ✅ | Via LiteLLM | - |
| Mistral | ✅ | ❌ | Not configured | P2 |
| Cohere | ✅ | ❌ | Not configured | P3 |
| Together | ✅ | ❌ | Not configured | P2 |
| Ollama | ✅ | ✅ | **Configured** | - |
| AWS Bedrock | ✅ | ✅ | **Working** | - |

### 12.2 Gateway Providers (3 total)

| Provider | Agno | OneMind | Status | Priority |
|----------|:----:|:-------:|--------|----------|
| Azure | ✅ | ❌ | Not configured | P2 |
| OpenRouter | ✅ | ❌ | Not configured | P2 |
| LiteLLM | ✅ | ✅ | **Working** | - |

**Models Summary:**
- **Working:** 8/13 (62%)
- **Not Configured:** 5/13 (38%)

---

## 13. Storage Gap Analysis

### 13.1 Database Options (3 total)

| Database | Agno | OneMind | Status | Priority |
|----------|:----:|:-------:|--------|----------|
| PostgresDb | ✅ | ✅ | **Working** | - |
| SqliteDb | ✅ | ⚠️ | Available, not used | P3 |
| InMemoryDb | ✅ | ⚠️ | Available, not used | P3 |

### 13.2 Custom Table Config (5 total)

| Feature | Agno | OneMind | Status | Priority |
|---------|:----:|:-------:|--------|----------|
| Custom session table | ✅ | ⚠️ | Not used | P3 |
| Custom memory table | ✅ | ⚠️ | Not used | P3 |
| Custom metrics table | ✅ | ⚠️ | Not used | P3 |
| Custom eval table | ✅ | ⚠️ | Not used | P3 |
| Custom knowledge table | ✅ | ⚠️ | Not used | P3 |

**Storage Summary:**
- **Working:** 1/8 (12%)
- **Available/Not Used:** 7/8 (88%)

---

## 14. Priority Summary

### P0 - Critical (8 features)

| # | Feature | Category | Current State |
|---|---------|----------|---------------|
| 1 | Knowledge auto-search in agent | Knowledge | Not wired |
| 2 | search_knowledge parameter | Agent | Not wired |
| 3 | Team Route mode | Teams | Basic only |
| 4 | Team Coordinate mode | Teams | Not implemented |
| 5 | Team Collaborate mode | Teams | Not implemented |
| 6 | Workflow Router step | Workflows | Scaffolded |
| 7 | Workflow Parallel step | Workflows | Scaffolded |
| 8 | MCP auto-attach to agents | MCP | Manual only |

### P1 - High Priority (22 features)

| # | Feature | Category |
|---|---------|----------|
| 1 | session_state | Agent |
| 2 | enable_agentic_state | Agent |
| 3 | retries/backoff | Agent |
| 4 | reasoning_model | Agent |
| 5 | Team model config | Teams |
| 6 | Team tools | Teams |
| 7 | show_members_responses | Teams |
| 8 | Sequential/Parallel modes | Teams |
| 9 | Loop step | Workflows |
| 10 | Step input/output piping | Workflows |
| 11 | HITL at workflow steps | Workflows |
| 12 | ShellTools (with HITL) | Tools |
| 13 | PythonTools (with HITL) | Tools |
| 14 | SQLTools (with HITL) | Tools |
| 15 | SlackTools | Tools |
| 16 | EmailTools | Tools |
| 17 | FileTools | Tools |
| 18 | FileGenerationTools | Tools |
| 19 | ReasoningTools | Tools |
| 20 | Semantic chunking | Knowledge |
| 21 | MCPTools in agent | MCP |
| 22 | AccuracyEval | Evals |

### P2 - Medium Priority (25 features)

| # | Feature | Category |
|---|---------|----------|
| 1 | expected_output | Agent |
| 2 | timezone_identifier | Agent |
| 3 | tool_choice/tool_call_limit | Agent |
| 4 | input_schema | Agent |
| 5 | reasoning_agent | Agent |
| 6 | reasoning_min/max_steps | Agent |
| 7 | search_session_history | Agent |
| 8 | cache_session | Agent |
| 9 | Team enable_agentic_memory | Teams |
| 10 | Function step | Workflows |
| 11 | Input caching | Workflows |
| 12 | ExaTools | Tools |
| 13 | TavilyTools | Tools |
| 14 | FirecrawlTools | Tools |
| 15 | Newspaper4kTools | Tools |
| 16 | E2BTools | Tools |
| 17 | TwilioTools | Tools |
| 18 | GoogleCalendarTools | Tools |
| 19 | Topic classification | Memory |
| 20 | DELETE /knowledge/{id} | Knowledge |
| 21 | Agentic chunking | Knowledge |
| 22 | MCP SSE transport | MCP |
| 23 | Additional MCP servers | MCP |
| 24 | Approval webhooks/timeout | HITL |
| 25 | Agent-as-judge eval | Evals |

### P3 - Low Priority (13 features)

| # | Feature | Category |
|---|---------|----------|
| 1 | system_message_role | Agent |
| 2 | introduction | Agent |
| 3 | parser_model | Agent |
| 4 | output_model | Agent |
| 5 | show_tool_calls | Agent |
| 6 | BraveSearchTools | Tools |
| 7 | CrawlbaseTools | Tools |
| 8 | NotionTools | Tools |
| 9 | Pinecone/Qdrant/Chroma | Knowledge |
| 10 | Mistral provider | Models |
| 11 | Cohere provider | Models |
| 12 | SqliteDb/InMemoryDb | Storage |
| 13 | Custom table names | Storage |

---

## 15. Implementation Roadmap

### Phase 1: Critical Gaps (P0) - Week 1-2

| Task | Files | Effort |
|------|-------|--------|
| Wire knowledge auto-search to agent | `core/agent/legacy.py` | 4h |
| Implement Team coordinate mode | `core/teams/coordinator.py` | 8h |
| Implement Team collaborate mode | `core/teams/coordinator.py` | 8h |
| Complete Router step execution | `core/workflows/engine.py` | 6h |
| Complete Parallel step execution | `core/workflows/engine.py` | 6h |
| Enable MCP auto-attachment | `core/mcp/tools.py`, `core/agent/factory.py` | 4h |

**Total Phase 1: ~36 hours**

### Phase 2: High Priority (P1) - Week 3-4

| Task | Files | Effort |
|------|-------|--------|
| Add session state features | `core/agent/factory.py` | 4h |
| Add retry/backoff config | `core/agent/factory.py` | 2h |
| Add reasoning_model config | `core/agent/legacy.py` | 2h |
| Add dangerous tools with HITL | `core/tools/builtin/` | 8h |
| Add SlackTools, EmailTools | `core/tools/builtin/` | 4h |
| Add FileGenerationTools | `core/tools/builtin/` | 4h |
| Add ReasoningTools | `core/tools/builtin/` | 4h |
| Implement Loop step | `core/workflows/engine.py` | 6h |
| Implement AccuracyEval | `core/evals/accuracy.py` | 8h |

**Total Phase 2: ~42 hours**

### Phase 3: Medium Priority (P2) - Week 5-6

| Task | Effort |
|------|--------|
| Add remaining agent parameters | 8h |
| Add ExaTools, TavilyTools, etc. | 8h |
| Add FirecrawlTools, Newspaper4k | 6h |
| Add missing knowledge features | 6h |
| Add MCP server definitions | 4h |
| Add approval webhooks/timeout | 6h |
| Deploy Phoenix/Langfuse | 4h |

**Total Phase 3: ~42 hours**

### Phase 4: Low Priority (P3) - Future

| Task | Effort |
|------|--------|
| Add remaining agent parameters | 4h |
| Add optional tools | 6h |
| Add optional vector DBs | 4h |
| Add additional model providers | 4h |

**Total Phase 4: ~18 hours**

---

## 16. Verification Checklist

After implementation, verify each feature:

### Agents
- [ ] Create agent with all core parameters
- [ ] Test knowledge auto-search
- [ ] Test session state persistence
- [ ] Test reasoning with dedicated model
- [ ] Test retry/backoff on failures

### Teams
- [ ] Create team with Route mode
- [ ] Test Coordinate mode
- [ ] Test Collaborate mode
- [ ] Test Sequential mode
- [ ] Test Parallel mode (delegate_to_all_members)

### Workflows
- [ ] Create workflow with Router step
- [ ] Test Condition step
- [ ] Test Parallel step execution
- [ ] Test Loop step
- [ ] Test HITL at workflow steps

### Tools
- [ ] Test each built-in tool individually
- [ ] Test ShellTools with HITL confirmation
- [ ] Test PythonTools with HITL confirmation
- [ ] Test FileGenerationTools

### Knowledge
- [ ] Upload document and verify search works IN agent context
- [ ] Test semantic chunking
- [ ] Test delete endpoint

### MCP
- [ ] Verify MCP auto-attaches to agents
- [ ] Test all defined servers
- [ ] Test MultiMCPTools

### HITL
- [ ] Test tool confirmation flow
- [ ] Test approval queue API
- [ ] Test continue_run

### Evals
- [ ] Run AccuracyEval and verify results stored
- [ ] Run PerformanceEval
- [ ] Run ReliabilityEval

---

*Gap analysis generated 2026-01-26*
*Reference: [AGNO-COMPLETE-FEATURE-INVENTORY.md](./AGNO-COMPLETE-FEATURE-INVENTORY.md)*

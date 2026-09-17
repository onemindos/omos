# OneMind OS API Documentation

Complete API reference for OneMind OS - AI-Native Human Enhancement Framework.

**Base URL**: `http://localhost:8080`

**Interactive Docs**:
- Swagger UI: http://localhost:8080/docs
- ReDoc: http://localhost:8080/redoc

**Version**: 1.0.0

---

## Table of Contents

1. [Overview](#overview)
2. [Authentication](#authentication)
3. [Agents API](#agents-api)
4. [Teams API](#teams-api)
5. [Workflows API](#workflows-api)
6. [Models API](#models-api)
7. [Tools API](#tools-api)
8. [Skills API](#skills-api)
9. [Voice API](#voice-api)
10. [Activity Feed API](#activity-feed-api)
11. [Approvals API (HITL)](#approvals-api-hitl)
12. [Awareness API](#awareness-api)
13. [Sessions API](#sessions-api)
14. [Memory API](#memory-api)
15. [Knowledge API (RAG)](#knowledge-api-rag)
16. [WebSocket Protocols](#websocket-protocols)

---

## Overview

OneMind OS provides a comprehensive API built on AgentOS with custom extensions for:
- Multi-agent orchestration
- Human-in-the-loop workflows
- Native voice AI (Gemini 2.0)
- Unified activity feed
- RAG-powered knowledge management
- Awareness system for context

### Key Features

- **53+ AgentOS Native Endpoints** - Full agent lifecycle management
- **Real-time Streaming** - SSE for agent runs, WebSocket for voice
- **Multi-modal** - Text, voice, images, video
- **60+ Models** - Across 11 providers (AWS, Google, OpenAI, etc.)
- **40+ Tools** - Web search, code execution, smart home, etc.
- **HITL Approvals** - Human confirmation for critical operations

### Response Format

All endpoints return JSON unless otherwise specified.

**Success Response**:
```json
{
  "data": {...},
  "status": "success"
}
```

**Error Response**:
```json
{
  "detail": "Error message",
  "status_code": 400
}
```

---

## Authentication

Currently, authentication is not enforced. All endpoints accept an optional `user_id` query parameter (defaults to `"zeus"`).

**Future**: OAuth2 + JWT tokens planned.

---

## Agents API

Manage and run agents. Agents are specialized AI assistants with tools, memory, and instructions.

### List Agents

```http
GET /agents
```

Returns all registered agents (30+ presets + remote agents).

**Response**:
```json
{
  "agents": [
    {
      "id": "legacy-mother-ai",
      "name": "Legacy Mother AI",
      "description": "Primary coordinator with full toolkit",
      "model": "sonnet",
      "tools": ["search", "code", "home"],
      "category": "core"
    }
  ]
}
```

**Categories**: core, specialist, workflow, multimodal, domain

---

### Get Agent

```http
GET /agents/{agent_id}
```

Get detailed agent configuration.

**Path Parameters**:
- `agent_id` (string, required) - Agent identifier

**Response**:
```json
{
  "id": "coder",
  "name": "Coder",
  "description": "Software development specialist",
  "model": "sonnet",
  "tools": ["github", "shell", "python"],
  "instructions": "You are an expert software engineer...",
  "enable_user_memories": true,
  "markdown": true
}
```

---

### Create Agent

```http
POST /agents
```

Create a new agent.

**Request Body**:
```json
{
  "id": "custom-agent",
  "name": "Custom Agent",
  "description": "My custom agent",
  "model": "sonnet",
  "instructions": "You are a helpful assistant...",
  "tools": ["search", "calculator"],
  "enable_user_memories": true,
  "enable_session_summaries": true
}
```

**Response**: Agent object

---

### Update Agent

```http
PUT /agents/{agent_id}
```

Update agent configuration.

**Request Body**: Partial agent object (only fields to update)

---

### Delete Agent

```http
DELETE /agents/{agent_id}
```

Delete an agent.

---

### Run Agent (Streaming)

```http
POST /agents/{agent_id}/runs
Content-Type: application/x-www-form-urlencoded
```

Execute agent with SSE streaming.

**Form Parameters**:
- `message` (string, required) - User message
- `stream` (boolean, default: true) - Enable streaming
- `session_id` (string, optional) - Session ID for conversation continuity
- `user_id` (string, default: "zeus") - User identifier

**Response**: SSE stream (text/event-stream)

**SSE Event Types**:
- `run_started` - Run initiated
- `text_message_content` - Text chunk (streaming)
- `tool_call_start` - Tool execution started
- `tool_call_end` - Tool execution completed
- `step_started` - Reasoning step started
- `step_completed` - Reasoning step completed
- `run_finished` - Run completed
- `run_error` - Run failed
- `member_response` - Team member response (teams only)

**Example**:
```bash
curl -N http://localhost:8080/agents/coder/runs \
  -d "message=Review this code" \
  -d "stream=true" \
  -d "session_id=session-123"
```

**SSE Event Example**:
```
event: run_started
data: {"run_id": "run-abc123", "agent_id": "coder"}

event: text_message_content
data: {"content": "I'll review the code...", "delta": "I'll"}

event: tool_call_start
data: {"tool_name": "github", "arguments": {"action": "list_prs"}}

event: tool_call_end
data: {"tool_name": "github", "result": "Found 3 PRs"}

event: run_finished
data: {"run_id": "run-abc123", "status": "complete"}
```

---

### Continue Paused Run

```http
POST /agents/{agent_id}/runs/{run_id}/continue
```

Continue a paused run (after HITL approval).

**Request Body**:
```json
{
  "tool_confirmations": [
    {
      "tool_call_id": "call_123",
      "confirmed": true,
      "reason": "Approved by user"
    }
  ]
}
```

---

## Teams API

Orchestrate multiple agents working together.

### Team Modes

- **COORDINATE**: Leader decomposes task → assigns → collects results
- **COLLABORATE**: All agents share context and work together
- **DELEGATE**: Leader picks single best agent

---

### List Teams

```http
GET /teams
```

Returns all registered teams.

**Response**:
```json
{
  "teams": [
    {
      "id": "dev-team",
      "name": "Development Team",
      "description": "Coder + Researcher",
      "mode": "coordinate",
      "members": ["coder", "researcher"]
    }
  ]
}
```

---

### Run Team

```http
POST /teams/{team_id}/runs
```

Execute team with streaming.

**Form Parameters**: Same as agent runs

**Response**: SSE stream with `member_response` events

**Example**:
```bash
curl -N http://localhost:8080/teams/dev-team/runs \
  -d "message=Build a web scraper" \
  -d "stream=true"
```

---

## Workflows API

Create and execute multi-step workflows.

### List Workflows

```http
GET /api/workflows
```

Returns all workflows.

**Response**:
```json
{
  "workflows": [
    {
      "id": "daily-briefing",
      "name": "Daily Briefing",
      "description": "Morning news and tasks",
      "steps": [
        {
          "id": "step-1",
          "name": "Get News",
          "type": "agent",
          "config": {
            "agent_id": "researcher",
            "message": "Get top news"
          }
        }
      ],
      "schedule": "0 8 * * *",
      "is_enabled": true,
      "run_count": 42
    }
  ],
  "total": 12
}
```

---

### Create Workflow

```http
POST /api/workflows
```

**Request Body**:
```json
{
  "name": "My Workflow",
  "description": "Automated task flow",
  "steps": [
    {
      "name": "Research",
      "type": "agent",
      "config": {
        "agent_id": "researcher",
        "message": "Research {input.topic}"
      }
    },
    {
      "name": "Summarize",
      "type": "agent",
      "config": {
        "agent_id": "writer",
        "message": "Summarize: {step_0}"
      }
    }
  ],
  "schedule": null,
  "is_enabled": true,
  "tags": ["automation"]
}
```

---

### Run Workflow

```http
POST /api/workflows/{workflow_id}/run?stream=true
```

Execute workflow with optional streaming.

**Query Parameters**:
- `stream` (boolean, default: false) - Enable SSE streaming

**Request Body**:
```json
{
  "input": {
    "topic": "AI agents"
  },
  "stream": true
}
```

**Response (Streaming)**:
```
event: workflow_started
data: {"run_id": "run-123", "workflow_id": "daily-briefing", "total_steps": 3}

event: step_started
data: {"step_index": 0, "step_id": "step-1", "step_name": "Get News"}

event: step_completed
data: {"step_index": 0, "output": "Top news: ...", "duration_ms": 2500}

event: workflow_completed
data: {"run_id": "run-123", "steps_completed": 3, "output": "...", "duration_ms": 8000}
```

**Response (Non-Streaming)**:
```json
{
  "run_id": "run-123",
  "workflow_id": "daily-briefing",
  "status": "completed",
  "steps_completed": 3,
  "total_steps": 3,
  "output": "Final result...",
  "duration_ms": 8000
}
```

---

### Update Workflow

```http
PUT /api/workflows/{workflow_id}
```

**Request Body**: Partial workflow object

---

### Delete Workflow

```http
DELETE /api/workflows/{workflow_id}
```

---

## Models API

Query available LLM models (60+ across 11 providers).

### List Models

```http
GET /models?available_only=false&provider=google
```

**Query Parameters**:
- `available_only` (boolean, default: false) - Only return models with configured API keys
- `provider` (string, optional) - Filter by provider

**Response**:
```json
{
  "models": [
    {
      "id": "sonnet",
      "name": "Claude 3.5 Sonnet",
      "provider": "bedrock",
      "family": "claude",
      "modality": ["text", "vision"],
      "context_length": 200000,
      "cost_per_1k_tokens": {
        "input": 0.003,
        "output": 0.015
      },
      "available": true,
      "reasoning": false
    },
    {
      "id": "gemini-2.0-flash",
      "name": "Gemini 2.0 Flash",
      "provider": "google",
      "modality": ["text", "vision", "audio"],
      "context_length": 1048576,
      "available": true,
      "native_audio": true
    }
  ]
}
```

**Providers**: bedrock, google, openai, groq, together, deepseek, mistral, cohere, azure, openrouter

---

### List Providers

```http
GET /models/providers
```

Returns metadata about each provider.

**Response**:
```json
{
  "providers": [
    {
      "id": "bedrock",
      "name": "AWS Bedrock",
      "models_count": 15,
      "configured": true,
      "capabilities": ["text", "vision"]
    }
  ]
}
```

---

## Tools API

Query available tools (40+ across 9 categories).

### List Tools

```http
GET /tools
```

**Response**:
```json
{
  "tools": [
    {
      "id": "duckduckgo_search",
      "name": "DuckDuckGo Search",
      "category": "search",
      "description": "Web search using DuckDuckGo",
      "parameters": {
        "query": {
          "type": "string",
          "required": true,
          "description": "Search query"
        },
        "max_results": {
          "type": "integer",
          "default": 5
        }
      },
      "available": true
    }
  ],
  "total": 45
}
```

**Categories**: search, finance, code, dev, comms, home, custom, activity, reasoning, builtin

---

## Skills API

Manage reusable instruction templates for agents.

### List Skills

```http
GET /api/skills
```

Returns all available skills.

**Response**:
```json
{
  "skills": [
    {
      "id": "web-research",
      "name": "Web Research",
      "description": "Advanced web research techniques",
      "version": "1.0.0",
      "author": "OneMind",
      "tags": ["research", "web"],
      "has_scripts": true,
      "has_references": true
    }
  ],
  "total": 8
}
```

---

### Get Skill

```http
GET /api/skills/{skill_id}
```

Get full skill content.

**Response**:
```json
{
  "id": "web-research",
  "name": "Web Research",
  "description": "Advanced web research techniques",
  "version": "1.0.0",
  "author": "OneMind",
  "tags": ["research", "web"],
  "content": "# Web Research Skill\n\nYou are an expert web researcher...",
  "scripts": ["analyze.py"],
  "references": ["search-patterns.md"]
}
```

---

### Get Skill Reference

```http
GET /api/skills/{skill_id}/reference/{filename}
```

Returns plain text reference document.

---

### Get Skill Script

```http
GET /api/skills/{skill_id}/script/{filename}
```

Returns plain text script content.

---

### Clear Skills Cache

```http
POST /api/skills/cache/clear
```

Force reload of skills registry.

---

## Voice API

Native Gemini audio for voice interactions (single-hop, ~300-500ms latency).

### Voice Run (HTTP)

```http
POST /api/voice/run
```

Execute voice agent run with native audio.

**Request Body**:
```json
{
  "audio_base64": "UklGRiQAAABXQVZFZm10...",
  "agent_id": "legacy-mother-ai",
  "session_id": "session-123",
  "user_id": "zeus",
  "voice_id": "kore"
}
```

**Request Fields**:
- `audio_base64` (string, required) - Base64 encoded audio (wav, mp3, webm, ogg)
- `agent_id` (string, default: "legacy-mother-ai") - Agent to use
- `session_id` (string, optional) - For conversation continuity
- `user_id` (string, default: "zeus") - User identifier
- `voice_id` (string, default: "kore") - Voice selection

**Response**:
```json
{
  "response_text": "Hello! How can I help you?",
  "audio_base64": "UklGRiQAAABXQVZFZm10...",
  "voice_id": "kore",
  "agent_id": "legacy-mother-ai",
  "session_id": "session-123"
}
```

**Example**:
```bash
# Encode audio file
AUDIO_BASE64=$(base64 -i recording.wav)

# Send request
curl -X POST http://localhost:8080/api/voice/run \
  -H "Content-Type: application/json" \
  -d "{
    \"audio_base64\": \"$AUDIO_BASE64\",
    \"agent_id\": \"legacy-mother-ai\",
    \"voice_id\": \"kore\"
  }"
```

---

### List Voices

```http
GET /api/voice/voices
```

Returns available Gemini voices.

**Response**:
```json
{
  "voices": [
    {
      "id": "kore",
      "name": "Kore",
      "description": "Female, warm"
    },
    {
      "id": "charon",
      "name": "Charon",
      "description": "Male, deep"
    },
    {
      "id": "aoede",
      "name": "Aoede",
      "description": "Female, expressive"
    },
    {
      "id": "fenrir",
      "name": "Fenrir",
      "description": "Male, confident"
    },
    {
      "id": "puck",
      "name": "Puck",
      "description": "Neutral, friendly"
    }
  ],
  "default": "kore"
}
```

---

### Voice Health Check

```http
GET /api/voice/health
```

Check voice system status.

**Response**:
```json
{
  "status": "ok",
  "gemini_enabled": true,
  "native_audio": true,
  "pipeline": "single-hop"
}
```

---

## Activity Feed API

Unified inbox for messages, notifications, tasks, and approvals.

### List Activity Items

```http
GET /api/activity?status=new,reading&type=task&priority=high&limit=50
```

**Query Parameters**:
- `status` (string, optional) - Comma-separated statuses (new, reading, done, archived, snoozed)
- `type` (string, optional) - Filter by type (message, notification, event, approval, task, alert, reminder, suggestion)
- `priority` (string, optional) - Filter by priority (critical, high, medium, low)
- `assignee` (string, optional) - Filter by assignee type (human, agent, system)
- `is_read` (boolean, optional) - Filter by read status
- `is_starred` (boolean, optional) - Filter by starred
- `limit` (integer, default: 50, max: 200) - Number of items
- `offset` (integer, default: 0) - Pagination offset
- `user_id` (string, default: "zeus") - User identifier

**Response**:
```json
{
  "items": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "type": "task",
      "priority": "high",
      "status": "new",
      "title": "Review PR #123",
      "body": "Code review needed for authentication changes",
      "source": "github",
      "metadata": {
        "repo": "onemind-os/backend",
        "pr_number": 123
      },
      "assignee_type": "human",
      "assignee_id": "zeus",
      "user_id": "zeus",
      "is_read": false,
      "is_starred": false,
      "created_at": "2026-02-05T10:00:00Z",
      "updated_at": "2026-02-05T10:00:00Z"
    }
  ],
  "total": 42,
  "unread_count": 12
}
```

---

### Get Activity Item

```http
GET /api/activity/{item_id}?user_id=zeus
```

Get single activity item.

---

### Create Activity Item

```http
POST /api/activity?awareness_mode=present
```

**Query Parameters**:
- `awareness_mode` (string, optional) - Current awareness mode for triage

**Request Body**:
```json
{
  "type": "task",
  "priority": "high",
  "title": "Review PR #123",
  "body": "Code review needed",
  "source": "github",
  "metadata": {
    "repo": "onemind-os/backend"
  },
  "assignee_type": "human",
  "assignee_id": "zeus",
  "user_id": "zeus"
}
```

**Response**: Created activity item

---

### Update Activity Item

```http
PATCH /api/activity/{item_id}?user_id=zeus
```

**Request Body**: Partial activity item (only fields to update)

---

### Mark as Read

```http
POST /api/activity/{item_id}/read?user_id=zeus
```

---

### Mark All as Read

```http
POST /api/activity/read-all?user_id=zeus
```

Returns:
```json
{
  "updated": 15
}
```

---

### Star/Unstar Item

```http
POST /api/activity/{item_id}/star?user_id=zeus
POST /api/activity/{item_id}/unstar?user_id=zeus
```

---

### Archive Item

```http
POST /api/activity/{item_id}/archive?user_id=zeus
```

---

### Delete Item

```http
DELETE /api/activity/{item_id}?user_id=zeus
```

---

### Snooze Item

```http
POST /api/activity/{item_id}/snooze?user_id=zeus
```

**Request Body**:
```json
{
  "until": "2026-02-06T09:00:00Z"
}
```

---

### Assign Item

```http
POST /api/activity/{item_id}/assign?user_id=zeus
```

**Request Body**:
```json
{
  "assignee_type": "agent",
  "assignee_id": "coder"
}
```

---

### Approve/Deny (Approval Items)

```http
POST /api/activity/{item_id}/approve?user_id=zeus
POST /api/activity/{item_id}/deny?user_id=zeus
```

**Deny Request Body**:
```json
{
  "reason": "Too risky"
}
```

---

### Complete/Fail Task

```http
POST /api/activity/{item_id}/complete?actor=zeus&user_id=zeus
POST /api/activity/{item_id}/fail?actor=legacy&user_id=zeus
```

**Complete Request Body**:
```json
{
  "result": "Task completed successfully"
}
```

**Fail Request Body**:
```json
{
  "reason": "Unable to access resource"
}
```

---

### Get Statistics

```http
GET /api/activity/stats?user_id=zeus
```

**Response**:
```json
{
  "total": 156,
  "unread": 12,
  "starred": 8,
  "by_type": {
    "task": 45,
    "message": 60,
    "notification": 40,
    "approval": 5,
    "alert": 6
  },
  "by_status": {
    "new": 12,
    "reading": 3,
    "done": 130,
    "archived": 11
  },
  "by_priority": {
    "critical": 2,
    "high": 15,
    "medium": 80,
    "low": 59
  }
}
```

---

### List Approvals

```http
GET /api/activity/approvals?limit=50&user_id=zeus
```

Returns activity items of type "approval" with status "new".

---

### List Tasks

```http
GET /api/activity/tasks?status=new,reading&assignee=human&limit=50&user_id=zeus
```

Returns activity items of type "task".

---

### List Agent Tasks

```http
GET /api/activity/agent-tasks?limit=5
```

Returns pending tasks assigned to agents (used by activity worker).

---

### Get Item History

```http
GET /api/activity/{item_id}/history?limit=50
```

Returns history of changes to an activity item.

---

## Approvals API (HITL)

Human-in-the-loop approval system for agent actions.

### List Approvals

```http
GET /api/hitl?user_id=zeus&status=pending
```

**Query Parameters**:
- `user_id` (string, default: "zeus")
- `status` (string, optional) - pending, approved, rejected, expired

**Response**:
```json
{
  "approvals": [
    {
      "id": "run-abc123",
      "run_id": "run-abc123",
      "agent_id": "coder",
      "tool_name": "delete_file",
      "description": "Delete file: important.txt",
      "parameters": {
        "filename": "important.txt"
      },
      "status": "pending",
      "priority": "critical",
      "created_at": "2026-02-05T10:00:00Z",
      "expires_at": "2026-02-05T11:00:00Z"
    }
  ],
  "total": 1,
  "pending_count": 1
}
```

---

### Approve Request

```http
POST /api/hitl/{approval_id}/approve?user_id=zeus
```

**Request Body**:
```json
{
  "reason": "Approved after review"
}
```

**Response**:
```json
{
  "success": true,
  "id": "run-abc123",
  "status": "approved",
  "resolved_by": "zeus",
  "resolved_at": "2026-02-05T10:05:00Z",
  "result": {
    "content": "File deleted successfully",
    "run_id": "run-abc123"
  }
}
```

---

### Reject Request

```http
POST /api/hitl/{approval_id}/reject?user_id=zeus
```

**Request Body**:
```json
{
  "reason": "Too risky"
}
```

---

### List Paused Runs

```http
GET /api/hitl/paused-runs?agent_id=coder
```

Returns paused runs with requirements.

---

### Get Paused Run

```http
GET /api/hitl/paused-runs/{run_id}
```

---

### Continue Paused Run

```http
POST /api/hitl/paused-runs/{run_id}/continue
```

**Request Body**:
```json
{
  "tool_confirmations": [
    {
      "tool_call_id": "call_123",
      "confirmed": true,
      "reason": "Approved"
    },
    {
      "tool_call_id": "call_456",
      "confirmed": false,
      "reason": "Too risky"
    }
  ]
}
```

---

### Cancel Paused Run

```http
DELETE /api/hitl/paused-runs/{run_id}
```

---

### Get Approval Stats

```http
GET /api/hitl/stats
```

**Response**:
```json
{
  "total_paused": 3,
  "active_paused": 2,
  "expired_paused": 1,
  "by_agent": {
    "coder": 2,
    "ops": 1
  }
}
```

---

## Awareness API

System awareness modes for context-aware behavior.

### Get Awareness State

```http
GET /awareness
```

**Response**:
```json
{
  "mode": "present",
  "heartbeat_interval": 5.0,
  "last_update": "2026-02-05T10:00:00Z",
  "active_sessions": 3,
  "pending_approvals": 2
}
```

**Modes**:
- `dormant`: Sleeping (60s heartbeat)
- `aware`: Watching (20s heartbeat)
- `present`: Engaged (5s heartbeat)
- `omnipresent`: Full presence (1s heartbeat)

---

### Set Awareness Mode

```http
POST /awareness/mode/{mode}
```

**Path Parameters**:
- `mode` (string) - dormant, aware, present, omnipresent

**Response**:
```json
{
  "success": true,
  "mode": "present",
  "state": {
    "mode": "present",
    "heartbeat_interval": 5.0
  }
}
```

---

## Sessions API

Conversation session management.

### List Sessions

```http
GET /sessions?agent_id=coder&user_id=zeus&limit=50
```

**Response**:
```json
{
  "sessions": [
    {
      "id": "session-123",
      "agent_id": "coder",
      "user_id": "zeus",
      "created_at": "2026-02-05T09:00:00Z",
      "updated_at": "2026-02-05T10:00:00Z",
      "message_count": 15,
      "summary": "Code review discussion"
    }
  ]
}
```

---

### Get Session

```http
GET /sessions/{session_id}
```

---

### Create Session

```http
POST /sessions
```

**Request Body**:
```json
{
  "agent_id": "coder",
  "user_id": "zeus"
}
```

---

### Delete Session

```http
DELETE /sessions/{session_id}
```

---

## Memory API

User and agent memory management.

### List Memories

```http
GET /memories?user_id=zeus&agent_id=coder&limit=50
```

---

### Create Memory

```http
POST /memories
```

**Request Body**:
```json
{
  "user_id": "zeus",
  "agent_id": "coder",
  "content": "User prefers Python over JavaScript",
  "type": "preference"
}
```

---

### Get Memory

```http
GET /memories/{memory_id}
```

---

### Delete Memory

```http
DELETE /memories/{memory_id}
```

---

## Knowledge API (RAG)

Vector-based knowledge management for RAG.

### List Documents

```http
GET /knowledge?limit=50
```

---

### Add Document

```http
POST /knowledge
```

**Request Body**:
```json
{
  "content": "Text content or URL",
  "metadata": {
    "source": "documentation",
    "tags": ["api", "reference"]
  }
}
```

---

### Search Knowledge

```http
GET /knowledge/search?query=agent+tools&limit=10
```

**Response**:
```json
{
  "results": [
    {
      "content": "Agent tools enable...",
      "score": 0.95,
      "metadata": {
        "source": "documentation"
      }
    }
  ]
}
```

---

### Delete Document

```http
DELETE /knowledge/{document_id}
```

---

## WebSocket Protocols

### Activity Feed WebSocket

```
ws://localhost:8080/api/activity/ws?user_id=zeus
```

**Client → Server**:
```json
{
  "type": "ping"
}
```

**Server → Client**:
```json
{
  "type": "pong"
}
```

```json
{
  "type": "activity.stats.updated",
  "data": {
    "total": 156,
    "unread": 12
  }
}
```

```json
{
  "type": "activity.item.created",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "type": "task",
    "title": "New task"
  }
}
```

```json
{
  "type": "activity.item.updated",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "status": "done"
  }
}
```

---

### Voice Gateway WebSocket

```
ws://localhost:8080/voice/ws/{device_id}
```

**Client → Server**:

Audio chunk:
```json
{
  "type": "audio_chunk",
  "audio_base64": "UklGRiQAAABXQVZFZm10...",
  "agent_id": "legacy-mother-ai",
  "session_id": "session-123",
  "voice_id": "kore"
}
```

Wake word:
```json
{
  "type": "wake_word_detected",
  "timestamp": "2026-02-05T12:00:00Z"
}
```

Stop listening:
```json
{
  "type": "stop_listening"
}
```

Ping:
```json
{
  "type": "ping",
  "timestamp": "2026-02-05T12:00:00Z"
}
```

**Server → Client**:

Audio response:
```json
{
  "type": "audio_chunk",
  "audio_base64": "UklGRiQAAABXQVZFZm10...",
  "transcript": "Hello, how can I help?"
}
```

Text chunk:
```json
{
  "type": "text_chunk",
  "content": "Processing your request..."
}
```

Tool call:
```json
{
  "type": "tool_call",
  "tool_name": "search",
  "arguments": {
    "query": "AI agents"
  }
}
```

Completed:
```json
{
  "type": "completed",
  "message": "Response complete"
}
```

Error:
```json
{
  "type": "error",
  "error": "Failed to process audio"
}
```

Pong:
```json
{
  "type": "pong",
  "timestamp": "2026-02-05T12:00:00Z"
}
```

---

## Health Check

```http
GET /health
```

System health check (AgentOS native).

**Response**:
```json
{
  "status": "healthy",
  "version": "1.0.0",
  "agents": 30,
  "teams": 12,
  "workflows": 5
}
```

---

## Error Codes

| Code | Description |
|------|-------------|
| 200 | Success |
| 201 | Created |
| 204 | No Content (successful deletion) |
| 400 | Bad Request (validation error) |
| 404 | Not Found (resource doesn't exist) |
| 500 | Internal Server Error |
| 503 | Service Unavailable (dependency not initialized) |

---

## Rate Limiting

Currently no rate limiting. Future: 1000 requests/hour per user.

---

## Changelog

### v1.0.0 (2026-02-05)
- Initial API documentation
- Activity Feed API (NEW)
- Voice Gateway WebSocket (NEW)
- Native Gemini audio (single-hop)
- HITL Approvals API
- Workflows API with streaming
- 60+ models, 40+ tools

---

## Support

- GitHub: https://github.com/onemind-os/onemind-os
- Docs: http://localhost:8080/docs
- Email: support@onemind.ai

---

**Generated**: 2026-02-05
**Version**: 1.0.0

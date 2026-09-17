---
title: "Agent Memory (agentmemory)"
type: service-card
status: active
created: 2026-06-08
updated: 2026-06-08
owner: zeus
category: tool
tags: [tool, memory, mcp, ai, persistence, local]
---

# Agent Memory (agentmemory)

> Persistent long-term memory layer for AI agents — local, self-hosted, MCP-connected.

## Role in OneMind OS
- Stores and retrieves memories across agent sessions
- Layer 3 of the OneMind memory stack (agentmemory → future Honcho at Layer 4)
- Used by Legacy, OpenClaw, and all profiles

## Config
- Location: `~/.agentmemory/`
- Service: LaunchAgent, `localhost:3111`
- Viewer: `localhost:3113`
- Version: v0.9.24
- MCP tools: 53 tools registered
- Upgrade: `npm update -g @agentmemory/agentmemory`

## Use Cases
- Cross-session memory persistence
- Entity recall (people, projects, preferences)
- Agent self-improvement via memory feedback

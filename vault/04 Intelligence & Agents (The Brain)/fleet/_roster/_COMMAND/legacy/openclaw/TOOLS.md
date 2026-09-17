---
type: agent-tools
agent: legacy
status: active
updated: 2026-05-03
---
# TOOLS.md — Legacy

Legacy has access to ALL skills as the master orchestrator. Other agents inherit subsets.

## Skills (All)

| Skill | Category | Purpose |
|-------|----------|---------|
| `coding-agent` | Development | Code generation, debugging, PR creation |
| `github` | Development | Repo management, issues, PRs |
| `diffs` | Development | Code diff analysis |
| `clawhub` | Platform | OpenClaw skill marketplace |
| `session-logs` | Platform | Session history and logging |
| `gh-issues` | Development | GitHub issue management |
| `1password` | Security | Credential queries, vault management |
| `oracle` | Intelligence | Reasoning and analysis |
| `healthcheck` | Operations | Service health monitoring |
| `taskflow` | Productivity | Task workflow management |
| `taskflow-inbox-triage` | Productivity | Inbox triage and routing |
| `tavily` | Research | Web search and research |
| `summarize` | Intelligence | Content summarization |
| `openai-whisper` | Media | Audio transcription |
| `video-frames` | Media | Video frame extraction |
| `nano-pdf` | Media | PDF processing and extraction |
| `camsnap` | Media | Camera/image capture |
| `obsidian` | Knowledge | Codex vault read/write |
| `obsidian-vault-maintainer` | Knowledge | Vault structure maintenance |
| `wiki-maintainer` | Knowledge | Wiki content management |
| `weather` | Data | Weather forecasts and alerts |
| `voice-call` | Communication | Voice call capabilities |
| `telegram` | Communication | Telegram messaging |

## MCP Servers

| Server | Purpose |
|--------|---------|
| `context7` | Up-to-date library documentation |
| `github` | GitHub API (repos, issues, PRs) |

## Tool Profile

`full` — all tools, execution, browser, filesystem, voice

## Sandbox

`off` — Legacy runs unsandboxed as the main orchestrator

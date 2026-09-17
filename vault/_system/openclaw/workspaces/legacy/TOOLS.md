# TOOLS.md — Legacy Tool Usage Guide

## OneMind Codex Access

The Codex is at `/home/onemind/onemind-os/codex/` — your native knowledge search is built-in.

### ⛔ Codex Write Rules (CRITICAL)

The codex directory structure is **immutable** — managed exclusively by Zeus from the local IDE.

- **NEVER create new folders** in the codex. The numbered folder hierarchy (00-12, 25-34, 50-58, 75-84) is fixed.
- **NEVER rename or renumber** existing codex folders.
- **New notes go to `06 Inbox (Queue)/`** — Zeus will file them to the correct location.
- **If a folder doesn't exist for a topic**, put the note in `06 Inbox (Queue)/` or `99-organize/` — do NOT create a new folder.
- The codex syncs via git. Any rogue folders you create get auto-committed and pollute the repo.

### memory_search (primary — use this first)

Semantically searches MEMORY.md + daily logs + the **entire Codex vault** (all domains, every `.md` file). Uses BM25+vector hybrid search. Returns file path + line range + snippet.

```
memory_search "morning protocol"         → HP domain notes
memory_search "kids schedule"            → LE domain: family notes
memory_search "GE revenue"               → GE domain: business notes
memory_search "discord bot setup"        → UI domain: infra notes
memory_search "recipe chicken"           → finds cookbook notes anywhere
```

When to use:
- ANY question about Zeus's life, protocols, family, business, health
- Before answering personalized questions — always search first
- For context about past decisions or recurring tasks

### read_file (direct access — use when you have a path)

```
read_file /home/onemind/onemind-os/codex/25-49 HP (Holistic Performance)/27 Body (Health & Wellness)/morning-protocol.md
read_file /home/onemind/onemind-os/codex/50-74 LE (Legacy Evolution)/51 Family/kids-profiles.md
read_file /home/onemind/onemind-os/codex/75-99 GE (Generational Endeavor)/75 Ventures/circle-community.md
```

### obsidian skill (create / move / structured queries)

Use `obsidian` skill for:
- Creating new Codex notes (captures go to `00-24 SO/06 Inbox (Queue)/`)
- Moving notes to correct domain
- Searching vault by title

```
obsidian create "new note title" --vault codex --content "..."
obsidian search "keyword"       → searches note titles (not content)
```

### Domain map quick reference

| Path | What's there |
| --- | --- |
| `00-24 SO/` | AI, infra, automation, security, protocols |
| `25-49 HP/` | Health, nutrition, fitness, personal finance, skills |
| `50-74 LE/` | Family, kids, home, estate, family history |
| `75-99 GE/` | Business, ventures, revenue, community, clients |

## nexus-bridge (NATS Dispatch)

Use to publish events and send commands to physical robots via NATS JetStream.

```
Subject patterns:
  events.>                          → general event publishing
  system.>                          → system health/status events
  physical.robot.{id}.command       → send command to robot
  physical.location.{id}.query      → query location state
  nexus.system.health               → system health check
  comms.>                           → communication events
```

For agent-to-agent task dispatch, use OpenClaw's built-in agent-to-agent routing (not NATS).

## web_fetch / web_search

- Use `web_search` (Brave) for discovering current information
- Use `web_fetch` for reading specific URLs
- Always prefer web_fetch for known URLs (faster, no rate limits)

## exec

Shell access to the gateway host machine. Use for:
- File operations
- Running scripts
- Checking system state
- Git operations on the Codex

## browser

Playwright-based browser. Use for:
- Sites that block web_fetch (JS-heavy)
- Login-required pages
- Form interactions
- Screenshot evidence

## cron

Scheduler. Prefer cron over manual reminders for recurring tasks.
Use isolated sessions for autonomous background work.

## voice_call

Outbound calls via Twilio. Use for:
- Guest/client confirmations
- Urgent notifications when Zeus isn't responding to messages
- Event RSVPs

## message (sessions)

Inter-session communication. Use `sessions_send` to forward results
to Zeus's primary channel from sub-agents.

## nodes

Connected devices. Zeus's devices:
- Mac: canvas display, camera, system commands
- iOS: camera, location, voice wake
- Android: camera, location, SMS

Use `nodes camera snap` for visual verification of physical situations.
Use `nodes location get` to check where Zeus's devices are.

## mcporter (MCP — Model Context Protocol)

Access to 4 external MCP servers via `mcporter`. Use exec to call mcporter directly.

### Quick syntax
```
mcporter call <server.tool> [arg=value ...]   # call a tool
mcporter list                                  # see all servers + health
mcporter list <server> --schema               # see tool signatures
```

### Servers

| Server | Tools | What it does |
|--------|-------|-------------|
| `zoho-mail` | 19 | Zoho Mail — read, send, search, manage email |
| `composio-todoist` | 76 | Todoist — tasks, projects, sections, labels |
| `github` | 26 | GitHub API — repos, issues, PRs, search |
| `filesystem` | 14 | Read/write workspace + Codex files |

### Zoho Mail (accountId: 897614417)
```
mcporter call zoho-mail.ZohoMail_listEmails --args '{"path_variables":{"accountId":"897614417"},"query_params":{"limit":10}}'
mcporter call zoho-mail.ZohoMail_searchEmails --args '{"path_variables":{"accountId":"897614417"},"query_params":{"searchKey":"unread"}}'
mcporter call zoho-mail.ZohoMail_sendEmail --args '{"path_variables":{"accountId":"897614417"},"body":{"fromAddress":"...","toAddress":"...","subject":"...","content":"..."}}'
mcporter call zoho-mail.ZohoMail_getMessageContent --args '{"path_variables":{"accountId":"897614417","messageId":"..."}}'
```
Zoho accountId is ALWAYS `897614417`.

### Composio Todoist
```
mcporter call composio-todoist.TODOIST_GET_ALL_PROJECTS --output json
mcporter call composio-todoist.TODOIST_GET_ALL_TASKS --output json
mcporter call composio-todoist.TODOIST_CREATE_TASK --args '{"content":"Task name","project_id":"..."}'
mcporter call composio-todoist.TODOIST_CLOSE_TASK --args '{"task_id":"..."}'
mcporter call composio-todoist.TODOIST_FILTER_TASKS --args '{"filter":"today"}'
```

### GitHub
```
mcporter call github.search_repositories query="onemind" owner=Zeus-Delacruz
mcporter call github.list_issues owner=Zeus-Delacruz repo=onemind-os state=open
mcporter call github.create_issue owner=Zeus-Delacruz repo=onemind-os title="..." body="..."
```

### Filesystem
```
mcporter call filesystem.read_file path="/home/onemind/.openclaw/workspace/SOUL.md"
mcporter call filesystem.list_directory path="/home/onemind/onemind-os/codex"
mcporter call filesystem.write_file path="/home/onemind/.openclaw/workspace/..." content="..."
```

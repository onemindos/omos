---
title: "CLI — composio"
type: reference
status: active
created: 2026-04-10
domain: SO/01
tags: [type/reference, domain/so, topic/cli, topic/composio, topic/integrations]
ai_summary: "Composio CLI deep reference — connect apps, execute tools across Todoist/Strava/GitHub/Slack from terminal. Extends the composio-cli skill."
---

# composio — Composio CLI

> 1000+ app integrations via CLI. Execute Todoist, Strava, GitHub, Slack, and more without writing API code.
> Full skill: `_codex/skills/composio-cli/SKILL.md`
> Install: `pip install composio-cli` or `npm install -g composio-core`

---

## Auth

```bash
# Login (opens browser)
composio login

# Check auth status
composio whoami

# Use API key directly
export COMPOSIO_API_KEY=<key>
```

---

## Connect Apps

```bash
# List all available apps
composio apps

# Connect an app
composio add todoist
composio add github
composio add strava
composio add slack
composio add notion

# List connected integrations
composio integrations list

# Test a connection
composio tools --app todoist list
```

---

## Execute Tools

```bash
# List tools for an app
composio tools --app todoist list
composio tools --app github list

# Execute a tool
composio tools execute TODOIST_CREATE_TASK \
  --params '{"content": "Task name", "project_id": "123"}'

composio tools execute GITHUB_CREATE_ISSUE \
  --params '{"owner": "OneMind-OS", "repo": "onemind-codex", "title": "Issue title"}'
```

---

## Common Todoist Operations

```bash
# Create task
composio tools execute TODOIST_CREATE_TASK \
  --params '{"content": "Buy groceries", "due_string": "today"}'

# Get tasks
composio tools execute TODOIST_GET_TASKS \
  --params '{"project_id": "2349250282"}'

# Complete task
composio tools execute TODOIST_CLOSE_TASK \
  --params '{"task_id": "12345"}'

# Move to delete project (vault rule — never delete directly)
composio tools execute TODOIST_MOVE_TASK \
  --params '{"task_id": "12345", "project_id": "6gM4gJcq9HmXjmpg"}'
```

---

## Common GitHub Operations

```bash
# Create issue
composio tools execute GITHUB_CREATE_ISSUE \
  --params '{"owner": "OneMind-OS", "repo": "onemind-codex-live", "title": "Bug: ...", "body": "Details..."}'

# List open PRs
composio tools execute GITHUB_LIST_PULL_REQUESTS \
  --params '{"owner": "OneMind-OS", "repo": "onemind-codex-live", "state": "open"}'
```

---

## Common Strava Operations

```bash
# Get recent activities
composio tools execute STRAVA_LIST_ATHLETE_ACTIVITIES \
  --params '{"per_page": 10}'

# Get athlete stats
composio tools execute STRAVA_GET_ATHLETE_STATS \
  --params '{}'
```

---

## In Python Agents

```python
from composio import ComposioToolSet, App

toolset = ComposioToolSet()

# Get tools for specific apps
tools = toolset.get_tools(apps=[App.TODOIST, App.GITHUB])

# Execute directly
result = toolset.execute_action(
    action="TODOIST_CREATE_TASK",
    params={"content": "Task name", "project_id": "123"}
)
```

---

## Trigger Webhooks

```bash
# Set up trigger (listen for events)
composio triggers list --app github
composio triggers enable GITHUB_PULL_REQUEST_EVENT --connection-id <id>
```


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]

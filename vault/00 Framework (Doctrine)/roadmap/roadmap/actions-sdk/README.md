---
id: ui-08-20260113225954
title: "README"
type: index
path: UI
domain: 00-24
folder: "08"
status: active
created: 2026-01-13
modified: 2026-01-13
owner: zeus
tags: [path/ui, path/ui/08, type/index, status/active]
ai_summary: "README"
---
# OneMind Core

Unified API layer for the OneMind system. Provides:

- **Connectors**: Full Python API classes for each platform
- **Tools**: Letta-compatible atomic functions for parallel execution
- **CLI**: Human interface for testing and operations

## Installation

```bash
pip install -e .
```

## Quick Start

### CLI Usage

```bash
# Check connection status
onemind status

# Get today's tasks
onemind todoist tasks

# Get overdue tasks
onemind todoist tasks --filter overdue

# List projects
onemind todoist projects

# Create a task
onemind todoist create "Review PR" --due tomorrow --priority 2
```

### Python Usage

```python
# Using connectors (full API)
from onemind.connectors.todoist import TodoistConnector

connector = TodoistConnector()
tasks = connector.get_tasks(filter="today")

# Using Letta-compatible tools
from onemind.tools.todoist import get_todoist_tasks, create_todoist_task

result = get_todoist_tasks(filter_query="today")
```

## Architecture

```
onemind-core/
├── onemind/
│   ├── config.py           # Configuration management
│   ├── cli.py              # Typer CLI
│   ├── connectors/         # Full API classes
│   │   ├── base.py         # BaseConnector ABC
│   │   ├── todoist.py
│   │   └── ...
│   └── tools/              # Letta-compatible tools
│       ├── todoist/
│       │   ├── atomic.py   # Single-action tools (parallel)
│       │   └── mega.py     # Multi-action tool (complex ops)
│       └── ...
```

## Environment Variables

Set these in `~/.onemind/.env`:

```bash
TODOIST_API_TOKEN=your_token
HA_URL=https://your-ha-instance
HA_TOKEN=your_ha_token
# ... etc
```

## License

MIT

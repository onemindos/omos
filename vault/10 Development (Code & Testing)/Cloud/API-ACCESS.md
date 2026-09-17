---
id: ui-12-20260111053926
title: "One Mind API Access Guide"
type: concept
domain: SO
domain: 00-24
folder: "00-24 SO (Sovereign Operations)"
status: active
created: 2026-01-11
modified: 2026-01-11
owner: zeus
tags: [domain/so, folder/12]
todoist_project: "00-24 SO (Sovereign Operations)"
ai_summary: "One Mind API Access Guide"
---

# One Mind API Access Guide

This document explains how to programmatically interact with all One Mind external services via their APIs. This enables AI-driven continuous maintenance, bulk operations, and automated system management.

## Table of Contents

- [Overview](#overview)
- [Authentication](#authentication)
- [Letta Cloud API](#letta-cloud-api)
- [Todoist API](#todoist-api)
- [N8N API](#n8n-api)
- [GitHub API](#github-api)
- [Common Patterns](#common-patterns)
- [Error Handling](#error-handling)

## Overview

One Mind integrates with multiple external services:

| Service | Purpose | API Documentation |
|---------|---------|-------------------|
| **Letta Cloud** | 41 AI agents | https://docs.letta.com/api |
| **Todoist** | Task management | https://developer.todoist.com/rest/v2 |
| **N8N** | Workflow automation | https://docs.n8n.io/api/ |
| **GitHub** | Version control | https://docs.github.com/en/rest |

All API endpoints are configured in:
```
infrastructure/config/api_endpoints.yaml
```

## Authentication

### Environment Variables

All API keys are stored in environment variables (NEVER commit these):

```bash
# Required
export LETTA_API_KEY="letta_key_xxx"
export TODOIST_API_KEY="todoist_token_xxx"

# Optional
export N8N_API_KEY="n8n_api_key_xxx"
export GITHUB_TOKEN="ghp_xxx"
```

Load from `.env` file:
```bash
cd infrastructure
source .env  # or use dotenv in Python
```

### Loading API Config

Python example:
```python
import os
import yaml
from pathlib import Path

# Load API endpoints config
config_path = Path(__file__).parent.parent / "config" / "api_endpoints.yaml"
with open(config_path) as f:
    api_config = yaml.safe_load(f)

# Get Letta API details
letta_base_url = api_config['letta']['base_url']
letta_api_key = os.getenv(api_config['letta']['api_key_env'])
```

## Letta Cloud API

Letta Cloud hosts all 41 One Mind agents. Use the API to:
- Create/update/delete agents
- Modify agent memory and system prompts
- Send messages to agents
- Manage agent tools and blocks

### Base Configuration

```python
import requests

LETTA_BASE_URL = "https://api.letta.com"
LETTA_API_KEY = os.getenv("LETTA_API_KEY")

headers = {
    "Authorization": f"Bearer {LETTA_API_KEY}",
    "Content-Type": "application/json"
}
```

### List All Agents

```python
response = requests.get(
    f"{LETTA_BASE_URL}/agents",
    headers=headers
)

agents = response.json()
print(f"Total agents: {len(agents)}")
```

### Get Agent Details

```python
agent_id = "agent-xxx"

response = requests.get(
    f"{LETTA_BASE_URL}/agents/{agent_id}",
    headers=headers
)

agent = response.json()
print(f"Agent: {agent['name']}")
print(f"Model: {agent['llm_config']['model']}")
print(f"System: {agent['system'][:100]}...")
```

### Update Agent System Prompt

```python
agent_id = "agent-xxx"
new_system_prompt = """
You are Mother AI, the central coordinator of the One Mind system.
Your responsibilities include...
"""

payload = {
    "system": new_system_prompt
}

response = requests.patch(
    f"{LETTA_BASE_URL}/agents/{agent_id}",
    headers=headers,
    json=payload
)

if response.status_code in [200, 204]:
    print("✓ Agent updated successfully")
else:
    print(f"✗ Update failed: {response.status_code}")
```

### Update Agent Model

```python
payload = {
    "model": "anthropic/claude-sonnet-4-5"  # or claude-haiku-3
}

response = requests.patch(
    f"{LETTA_BASE_URL}/agents/{agent_id}",
    headers=headers,
    json=payload
)
```

### Send Message to Agent

```python
message = "Please perform your morning health check."

payload = {
    "messages": [
        {
            "role": "user",
            "text": message
        }
    ]
}

response = requests.post(
    f"{LETTA_BASE_URL}/agents/{agent_id}/messages",
    headers=headers,
    json=payload
)

messages = response.json()
for msg in messages:
    print(f"{msg['role']}: {msg['text']}")
```

### Bulk Update Pattern

Update all 41 agents at once:

```python
import yaml

# Load agent configs
config_path = Path("infrastructure/config/agents/command.yaml")
with open(config_path) as f:
    config = yaml.safe_load(f)

# Load deployed state to get agent IDs
state_path = Path("infrastructure/config/deployed_state.yaml")
with open(state_path) as f:
    state = yaml.safe_load(f)

# Update each agent
for agent_config in config['agents']:
    agent_name = agent_config['name']
    agent_id = state['deployed_agents'][agent_name]['letta_id']

    if not agent_id:
        print(f"⚠ {agent_name}: Not deployed")
        continue

    payload = {
        "system": agent_config['system_prompt'],
        "model": agent_config.get('model', 'anthropic/claude-sonnet-4-5')
    }

    response = requests.patch(
        f"{LETTA_BASE_URL}/agents/{agent_id}",
        headers=headers,
        json=payload
    )

    if response.status_code in [200, 204]:
        print(f"✓ {agent_name}: Updated")
    else:
        print(f"✗ {agent_name}: Failed ({response.status_code})")
```

## Todoist API

Todoist manages all One Mind tasks and projects across the four pillars (HP, LE, GE, IT).

### Base Configuration

```python
TODOIST_BASE_URL = "https://api.todoist.com/rest/v2"
TODOIST_API_KEY = os.getenv("TODOIST_API_KEY")

headers = {
    "Authorization": f"Bearer {TODOIST_API_KEY}",
    "Content-Type": "application/json"
}
```

### List All Projects

```python
response = requests.get(
    f"{TODOIST_BASE_URL}/projects",
    headers=headers
)

projects = response.json()
for project in projects:
    print(f"Project: {project['name']} (ID: {project['id']})")
```

### Create Task

```python
payload = {
    "content": "Review agent performance metrics",
    "project_id": "2345678901",  # HP project ID
    "due_string": "tomorrow at 9am",
    "priority": 3,  # 1=normal, 2=high, 3=urgent, 4=very urgent
    "labels": ["one-mind", "automation"]
}

response = requests.post(
    f"{TODOIST_BASE_URL}/tasks",
    headers=headers,
    json=payload
)

task = response.json()
print(f"Created task: {task['id']}")
```

### Get All Tasks for Project

```python
project_id = "2345678901"

response = requests.get(
    f"{TODOIST_BASE_URL}/tasks",
    headers=headers,
    params={"project_id": project_id}
)

tasks = response.json()
print(f"Found {len(tasks)} tasks in project")
```

### Complete Task

```python
task_id = "7654321098"

response = requests.post(
    f"{TODOIST_BASE_URL}/tasks/{task_id}/close",
    headers=headers
)

if response.status_code == 204:
    print("✓ Task completed")
```

### Bulk Task Creation

```python
# Create daily review tasks for all 4 pillars
pillars = [
    {"name": "HP", "project_id": "hp_project_id"},
    {"name": "LE", "project_id": "le_project_id"},
    {"name": "GE", "project_id": "ge_project_id"},
    {"name": "IT", "project_id": "it_project_id"},
]

for pillar in pillars:
    payload = {
        "content": f"{pillar['name']} Daily Review",
        "project_id": pillar['project_id'],
        "due_string": "every day at 9pm",
        "labels": ["daily-review", "one-mind"]
    }

    response = requests.post(
        f"{TODOIST_BASE_URL}/tasks",
        headers=headers,
        json=payload
    )

    if response.status_code == 200:
        print(f"✓ Created {pillar['name']} daily review task")
```

## N8N API

N8N orchestrates workflows that connect Letta agents, Todoist, and other services.

### Base Configuration

```python
N8N_BASE_URL = "https://n8n.zeusops.dev/api/v1"
N8N_API_KEY = os.getenv("N8N_API_KEY")

headers = {
    "X-N8N-API-KEY": N8N_API_KEY,  # Note: Different header format
    "Content-Type": "application/json"
}
```

### List All Workflows

```python
response = requests.get(
    f"{N8N_BASE_URL}/workflows",
    headers=headers
)

workflows = response.json()['data']
for workflow in workflows:
    print(f"Workflow: {workflow['name']} (ID: {workflow['id']})")
```

### Get Workflow Details

```python
workflow_id = "123"

response = requests.get(
    f"{N8N_BASE_URL}/workflows/{workflow_id}",
    headers=headers
)

workflow = response.json()
print(f"Workflow: {workflow['name']}")
print(f"Active: {workflow['active']}")
print(f"Nodes: {len(workflow['nodes'])}")
```

### Activate/Deactivate Workflow

```python
workflow_id = "123"

# Activate
response = requests.patch(
    f"{N8N_BASE_URL}/workflows/{workflow_id}",
    headers=headers,
    json={"active": True}
)

# Deactivate
response = requests.patch(
    f"{N8N_BASE_URL}/workflows/{workflow_id}",
    headers=headers,
    json={"active": False}
)
```

### Delete Workflow

```python
workflow_id = "123"

response = requests.delete(
    f"{N8N_BASE_URL}/workflows/{workflow_id}",
    headers=headers
)

if response.status_code == 200:
    print("✓ Workflow deleted")
```

### Import Workflow (Idempotent Deployment)

```python
import json

# Load workflow JSON
workflow_path = Path("infrastructure/workflows/production/one-mind-morning-briefing.json")
with open(workflow_path) as f:
    workflow_data = json.load(f)

workflow_name = workflow_data['name']

# Check if workflow exists
response = requests.get(f"{N8N_BASE_URL}/workflows", headers=headers)
existing_workflows = response.json()['data']

existing = next((w for w in existing_workflows if w['name'] == workflow_name), None)

if existing:
    # Delete existing workflow
    print(f"Deleting existing workflow: {existing['id']}")
    requests.delete(f"{N8N_BASE_URL}/workflows/{existing['id']}", headers=headers)

# Import new workflow
response = requests.post(
    f"{N8N_BASE_URL}/workflows",
    headers=headers,
    json=workflow_data
)

new_workflow = response.json()
print(f"✓ Imported workflow: {new_workflow['id']}")
```

## GitHub API

GitHub stores all infrastructure configurations and enables GitOps workflows.

### Base Configuration

```python
GITHUB_BASE_URL = "https://api.github.com"
GITHUB_TOKEN = os.getenv("GITHUB_TOKEN")

headers = {
    "Authorization": f"token {GITHUB_TOKEN}",
    "Accept": "application/vnd.github.v3+json"
}
```

### Get File Contents

```python
owner = "your-username"
repo = "one-mind"
path = "infrastructure/config/agents/command.yaml"

response = requests.get(
    f"{GITHUB_BASE_URL}/repos/{owner}/{repo}/contents/{path}",
    headers=headers
)

file_data = response.json()
import base64
content = base64.b64decode(file_data['content']).decode('utf-8')
print(content)
```

### Update File Contents

```python
import base64

new_content = """
agents:
  - name: "cmd-mother-ai"
    system_prompt: "Updated prompt..."
"""

# Get current file SHA (required for updates)
response = requests.get(
    f"{GITHUB_BASE_URL}/repos/{owner}/{repo}/contents/{path}",
    headers=headers
)
current_sha = response.json()['sha']

# Update file
payload = {
    "message": "Update agent config via API",
    "content": base64.b64encode(new_content.encode()).decode(),
    "sha": current_sha
}

response = requests.put(
    f"{GITHUB_BASE_URL}/repos/{owner}/{repo}/contents/{path}",
    headers=headers,
    json=payload
)

if response.status_code == 200:
    print("✓ File updated on GitHub")
```

## Common Patterns

### Pattern 1: GitOps Sync Loop

Agents pull their configs from GitHub, detect drift, and update themselves:

```python
# 1. Fetch config from GitHub
response = requests.get(
    f"{GITHUB_BASE_URL}/repos/{owner}/{repo}/contents/infrastructure/config/agents/command.yaml",
    headers=github_headers
)
config_content = base64.b64decode(response.json()['content']).decode('utf-8')
config = yaml.safe_load(config_content)

# 2. For each agent, compare with Letta Cloud state
for agent_config in config['agents']:
    agent_id = get_agent_id(agent_config['name'])

    letta_response = requests.get(
        f"{LETTA_BASE_URL}/agents/{agent_id}",
        headers=letta_headers
    )
    current_agent = letta_response.json()

    # 3. Detect drift
    if agent_config['system_prompt'] != current_agent['system']:
        print(f"Drift detected for {agent_config['name']}")

        # 4. Update agent
        requests.patch(
            f"{LETTA_BASE_URL}/agents/{agent_id}",
            headers=letta_headers,
            json={"system": agent_config['system_prompt']}
        )
```

### Pattern 2: Cross-Service Automation

Create Todoist task from Letta agent message:

```python
# 1. Agent sends message with task request
agent_message = "Create task: Review HP metrics"

# 2. Parse task from message
task_content = agent_message.replace("Create task: ", "")

# 3. Create in Todoist
todoist_response = requests.post(
    f"{TODOIST_BASE_URL}/tasks",
    headers=todoist_headers,
    json={
        "content": task_content,
        "project_id": hp_project_id,
        "labels": ["agent-created"]
    }
)

# 4. Confirm back to agent
task_id = todoist_response.json()['id']
print(f"✓ Created task {task_id} from agent request")
```

### Pattern 3: Health Check Workflow

N8N workflow calls all APIs to verify system health:

```python
def health_check():
    results = {}

    # Check Letta API
    try:
        response = requests.get(f"{LETTA_BASE_URL}/agents", headers=letta_headers, timeout=5)
        results['letta'] = {'status': 'healthy', 'agent_count': len(response.json())}
    except Exception as e:
        results['letta'] = {'status': 'unhealthy', 'error': str(e)}

    # Check Todoist API
    try:
        response = requests.get(f"{TODOIST_BASE_URL}/projects", headers=todoist_headers, timeout=5)
        results['todoist'] = {'status': 'healthy', 'project_count': len(response.json())}
    except Exception as e:
        results['todoist'] = {'status': 'unhealthy', 'error': str(e)}

    # Check N8N API
    try:
        response = requests.get(f"{N8N_BASE_URL}/workflows", headers=n8n_headers, timeout=5)
        results['n8n'] = {'status': 'healthy', 'workflow_count': len(response.json()['data'])}
    except Exception as e:
        results['n8n'] = {'status': 'unhealthy', 'error': str(e)}

    return results
```

## Error Handling

### Retry with Exponential Backoff

```python
import time

def api_call_with_retry(url, headers, max_retries=3):
    for attempt in range(max_retries):
        try:
            response = requests.get(url, headers=headers, timeout=10)

            if response.status_code == 429:  # Rate limit
                wait_time = 2 ** attempt
                print(f"Rate limited, waiting {wait_time}s...")
                time.sleep(wait_time)
                continue

            response.raise_for_status()
            return response.json()

        except requests.exceptions.Timeout:
            if attempt == max_retries - 1:
                raise
            print(f"Timeout, retrying ({attempt + 1}/{max_retries})...")
            time.sleep(2 ** attempt)

        except requests.exceptions.RequestException as e:
            if attempt == max_retries - 1:
                raise
            print(f"Request failed: {e}, retrying...")
            time.sleep(2 ** attempt)

    raise Exception(f"Failed after {max_retries} retries")
```

### Rate Limit Handling

```python
def check_rate_limit(response):
    if 'X-RateLimit-Remaining' in response.headers:
        remaining = int(response.headers['X-RateLimit-Remaining'])
        if remaining < 10:
            reset_time = int(response.headers['X-RateLimit-Reset'])
            wait_seconds = reset_time - time.time()
            print(f"⚠ Low rate limit ({remaining} remaining), waiting {wait_seconds}s")
            time.sleep(wait_seconds)
```

## Best Practices

1. **Always use environment variables for API keys** - Never hardcode credentials
2. **Implement retries with exponential backoff** - Handle transient failures
3. **Check rate limits** - Respect API quotas
4. **Use timeouts** - Prevent hanging requests (default: 10s)
5. **Log all API calls** - For debugging and audit trail
6. **Validate responses** - Check status codes and response structure
7. **Use the deployed_state.yaml** - Track agent IDs and workflow IDs
8. **Version control configs** - Commit all changes to GitHub
9. **Test in dry-run mode** - Before making bulk changes
10. **Monitor API usage** - Track costs and quota consumption

## Next Steps

- See [CONTINUOUS-MAINTENANCE.md](CONTINUOUS-MAINTENANCE.md) for automated maintenance workflows
- See [BULK-OPERATIONS.md](BULK-OPERATIONS.md) for system-wide change procedures
- See [DEPLOYMENT.md](DEPLOYMENT.md) for full deployment guide

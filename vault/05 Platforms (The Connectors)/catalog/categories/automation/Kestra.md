# Kestra

> **Data Orchestration Platform**

## Overview

Kestra is an open-source orchestration platform for building, scheduling, and monitoring complex data pipelines and workflows.

## Status

- **Active**: 🔄 Planned
- **Deployment**: Self-hosted
- **Last Updated**: January 25, 2026

## Use Cases

1. **ETL Pipelines** - Data extraction and transformation
2. **Scheduled Jobs** - Cron-like task scheduling
3. **Data Sync** - Cross-service data synchronization
4. **Batch Processing** - Large data processing jobs
5. **Workflow Orchestration** - Complex multi-step workflows

## Key Features

- YAML-based workflows
- Real-time monitoring
- Built-in triggers (schedule, event, webhook)
- 400+ plugins
- Namespace organization
- Variables and secrets
- Error handling and retries

## Installation

### Docker

```bash
docker run -d \
  --name kestra \
  -p 8080:8080 \
  -v kestra-data:/app/storage \
  kestra/kestra:latest server standalone
```

### Docker Compose

```yaml
services:
  kestra:
    image: kestra/kestra:latest
    command: server standalone
    ports:
      - "8080:8080"
    volumes:
      - kestra-data:/app/storage
      - /var/run/docker.sock:/var/run/docker.sock
    environment:
      KESTRA_CONFIGURATION: |
        datasources:
          postgres:
            url: jdbc:postgresql://postgres:5432/kestra
            driverClassName: org.postgresql.Driver
            username: kestra
            password: ${POSTGRES_PASSWORD}

  postgres:
    image: postgres:16
    environment:
      POSTGRES_USER: kestra
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
      POSTGRES_DB: kestra
    volumes:
      - postgres-data:/var/lib/postgresql/data

volumes:
  kestra-data:
  postgres-data:
```

## Workflow Definition

### Basic Flow

```yaml
id: hello-world
namespace: onemind
description: A simple hello world flow

tasks:
  - id: hello
    type: io.kestra.core.tasks.log.Log
    message: "Hello, World!"
```

### With Multiple Tasks

```yaml
id: daily-sync
namespace: onemind.data
description: Daily data synchronization

inputs:
  - id: date
    type: DATE
    defaults: "{{ now() | date('yyyy-MM-dd') }}"

tasks:
  - id: extract
    type: io.kestra.plugin.scripts.python.Script
    runner: DOCKER
    docker:
      image: python:3.12
    script: |
      import json
      data = extract_data("{{ inputs.date }}")
      print(json.dumps(data))

  - id: transform
    type: io.kestra.plugin.scripts.python.Script
    runner: DOCKER
    docker:
      image: python:3.12
    script: |
      data = '{{ outputs.extract.vars.data }}'
      transformed = transform_data(data)
      print(json.dumps(transformed))

  - id: load
    type: io.kestra.plugin.jdbc.postgresql.Query
    url: jdbc:postgresql://localhost:5432/onemind
    username: "{{ secret('DB_USER') }}"
    password: "{{ secret('DB_PASSWORD') }}"
    sql: |
      INSERT INTO sync_data (date, data)
      VALUES ('{{ inputs.date }}', '{{ outputs.transform.vars.data }}')

triggers:
  - id: daily
    type: io.kestra.core.models.triggers.types.Schedule
    cron: "0 6 * * *" # Every day at 6 AM
```

## Triggers

### Schedule Trigger

```yaml
triggers:
  - id: hourly
    type: io.kestra.core.models.triggers.types.Schedule
    cron: "0 * * * *"
```

### Webhook Trigger

```yaml
triggers:
  - id: webhook
    type: io.kestra.core.models.triggers.types.Webhook
    key: unique-webhook-key
```

### Flow Trigger

```yaml
triggers:
  - id: after_flow
    type: io.kestra.core.models.triggers.types.Flow
    inputs:
      data: "{{ trigger.executionId }}"
    conditions:
      - type: io.kestra.core.models.conditions.types.ExecutionFlowCondition
        namespace: onemind
        flowId: upstream-flow
      - type: io.kestra.core.models.conditions.types.ExecutionStatusCondition
        in:
          - SUCCESS
```

## Task Types

### Script Tasks

```yaml
# Python
- id: python-task
  type: io.kestra.plugin.scripts.python.Script
  script: |
    print("Hello from Python!")

# Shell
- id: shell-task
  type: io.kestra.plugin.scripts.shell.Commands
  commands:
    - echo "Hello from Shell!"

# Node.js
- id: node-task
  type: io.kestra.plugin.scripts.node.Script
  script: |
    console.log("Hello from Node!");
```

### HTTP Tasks

```yaml
- id: api-call
  type: io.kestra.plugin.fs.http.Request
  uri: https://api.example.com/data
  method: GET
  headers:
    Authorization: "Bearer {{ secret('API_TOKEN') }}"
```

### Database Tasks

```yaml
- id: query
  type: io.kestra.plugin.jdbc.postgresql.Query
  url: jdbc:postgresql://localhost:5432/db
  username: user
  password: "{{ secret('DB_PASSWORD') }}"
  sql: SELECT * FROM users
  fetch: true
```

## Error Handling

```yaml
tasks:
  - id: risky-task
    type: io.kestra.plugin.scripts.python.Script
    script: |
      # May fail
      risky_operation()
    retry:
      maxAttempt: 3
      type: constant
      interval: PT5M # 5 minutes

errors:
  - id: on-error
    type: io.kestra.plugin.notifications.slack.SlackIncomingWebhook
    url: "{{ secret('SLACK_WEBHOOK') }}"
    payload: |
      {"text": "Flow failed: {{ flow.id }}"}
```

## Variables and Secrets

### Namespace Variables

```yaml
# In namespace configuration
variables:
  environment: production
  api_base_url: https://api.onemind.io
```

### Using Secrets

```yaml
- id: secure-task
  type: io.kestra.plugin.scripts.python.Script
  env:
    API_KEY: "{{ secret('ANTHROPIC_API_KEY') }}"
  script: |
    import os
    api_key = os.environ['API_KEY']
```

## API Usage

### Trigger Flow via API

```bash
curl -X POST "http://localhost:8080/api/v1/executions/onemind/daily-sync" \
  -H "Content-Type: application/json" \
  -d '{"date": "2026-01-25"}'
```

### Get Execution Status

```bash
curl "http://localhost:8080/api/v1/executions/onemind/daily-sync/exec-id"
```

## Monitoring

### Dashboard

- View all executions
- Execution timeline
- Task logs
- Trigger history
- Namespace metrics

### Alerts

```yaml
listeners:
  - type: io.kestra.core.models.listeners.Listener
    conditions:
      - type: io.kestra.core.models.conditions.types.ExecutionStatusCondition
        in:
          - FAILED
    tasks:
      - id: alert
        type: io.kestra.plugin.notifications.slack.SlackIncomingWebhook
        url: "{{ secret('SLACK_WEBHOOK') }}"
        payload: |
          {"text": "Execution failed!"}
```

## vs N8N

| Feature      | Kestra         | N8N                 |
| ------------ | -------------- | ------------------- |
| Focus        | Data pipelines | Business automation |
| Definition   | YAML           | Visual canvas       |
| Scheduling   | Built-in cron  | Built-in cron       |
| Code support | Native         | Limited             |
| Scale        | Large data     | Medium workflows    |

## Best Practices

- ✅ Organize flows by namespace
- ✅ Use secrets for sensitive data
- ✅ Implement error handling
- ✅ Add retry logic for external calls
- ✅ Use inputs for parameterization
- ✅ Document flows with descriptions
- ✅ Monitor execution metrics

## Resources

- [Kestra Documentation](https://kestra.io/docs/)
- [Plugin Library](https://kestra.io/plugins)
- [GitHub Repository](https://github.com/kestra-io/kestra)

---

**Related**: [N8N](N8N.md) | [PostgreSQL](./PostgreSQL.md)


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]

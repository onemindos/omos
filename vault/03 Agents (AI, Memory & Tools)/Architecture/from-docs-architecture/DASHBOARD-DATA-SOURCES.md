# OneMind OS - Dashboard Data Sources

> How to access unified data for custom dashboards and widgets

---

## Overview

All integrated data flows through NATS and lands in a **unified data store**. Your dashboards can access this data through multiple methods depending on your needs.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                     DASHBOARD DATA ACCESS                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                    Your Dashboard / Widget                          │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                     │           │           │           │                   │
│                     ▼           ▼           ▼           ▼                   │
│              ┌───────────┐ ┌───────────┐ ┌───────────┐ ┌───────────┐       │
│              │  REST API │ │ WebSocket │ │  Direct   │ │  Nango    │       │
│              │ /dashboard│ │   NATS    │ │    SQL    │ │  Proxy    │       │
│              └───────────┘ └───────────┘ └───────────┘ └───────────┘       │
│                     │           │           │           │                   │
│                     └───────────┼───────────┴───────────┘                   │
│                                 │                                           │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                      Unified Data Store                              │    │
│  │  ┌───────────┐  ┌───────────┐  ┌───────────────────────────────────┐│    │
│  │  │PostgreSQL │  │   Redis   │  │        NATS JetStream             ││    │
│  │  │ (cold)    │  │  (hot)    │  │         (events)                  ││    │
│  │  └───────────┘  └───────────┘  └───────────────────────────────────┘│    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Data Access Methods

### 1. REST API (Recommended for UI)

```typescript
// Fetch all tasks from all sources
const response = await fetch('/api/dashboard/tasks');
const tasks = await response.json();

// Response structure
{
  "items": [
    {
      "id": "todoist_123",
      "source": "todoist",
      "type": "task",
      "data": {
        "content": "Review PR #42",
        "due": "2026-01-21",
        "priority": 2
      },
      "synced_at": "2026-01-20T10:00:00Z"
    },
    {
      "id": "github_456",
      "source": "github",
      "type": "task",
      "data": {
        "title": "Fix bug in auth",
        "repo": "onemind-os",
        "labels": ["bug"]
      },
      "synced_at": "2026-01-20T10:05:00Z"
    }
  ],
  "total": 42,
  "synced_at": "2026-01-20T10:05:00Z"
}
```

#### Available Endpoints

| Endpoint | Description | Parameters |
|----------|-------------|------------|
| `GET /dashboard/unified` | All items | `?source=&type=&limit=` |
| `GET /dashboard/tasks` | Tasks across sources | `?due=today&priority=` |
| `GET /dashboard/events` | Calendar events | `?start=&end=` |
| `GET /dashboard/emails` | Recent emails | `?unread=true&limit=` |
| `GET /dashboard/activity` | Activity timeline | `?limit=50` |
| `GET /dashboard/stats` | Aggregated stats | - |

### 2. WebSocket (Real-time Updates)

```typescript
// Connect to NATS via WebSocket
const ws = new WebSocket('ws://localhost:9222');

// Subscribe to events
ws.onopen = () => {
  ws.send(JSON.stringify({
    action: 'subscribe',
    subjects: [
      'events.todoist.*',
      'events.github.*',
      'events.sync.*',
    ],
  }));
};

ws.onmessage = (event) => {
  const { subject, data } = JSON.parse(event.data);

  switch (subject) {
    case 'events.todoist.task.created':
      addTaskToUI(data);
      break;
    case 'events.sync.caldav.completed':
      refreshCalendar();
      break;
  }
};
```

### 3. Direct SQL (Advanced)

For complex queries or custom analytics, connect directly to PostgreSQL.

```sql
-- All tasks due today across all sources
SELECT
  id,
  source,
  data->>'content' as title,
  data->>'due' as due_date,
  data->>'priority' as priority
FROM unified_items
WHERE type = 'task'
  AND (data->>'due')::date = CURRENT_DATE
ORDER BY (data->>'priority')::int DESC;

-- Activity in last 24 hours
SELECT
  source,
  type,
  data,
  synced_at
FROM unified_items
WHERE synced_at > NOW() - INTERVAL '24 hours'
ORDER BY synced_at DESC
LIMIT 100;

-- Items by source
SELECT
  source,
  COUNT(*) as count,
  MAX(synced_at) as last_sync
FROM unified_items
GROUP BY source;

-- Semantic search (with pgvector)
SELECT
  data,
  1 - (embedding <=> $1::vector) as similarity
FROM unified_items
WHERE embedding IS NOT NULL
ORDER BY embedding <=> $1::vector
LIMIT 10;
```

### 4. Nango Proxy (Fresh Data)

For data that needs to be fetched fresh from the source API:

```typescript
// Use Nango to proxy API calls with managed credentials
const response = await fetch('/api/nango/proxy', {
  method: 'POST',
  body: JSON.stringify({
    integration: 'todoist',
    endpoint: '/tasks',
    method: 'GET',
  }),
});

// Nango handles OAuth token refresh automatically
const tasks = await response.json();
```

---

## Unified Data Schema

All data is stored in the `unified_items` table:

```sql
CREATE TABLE unified_items (
  id TEXT PRIMARY KEY,           -- Unique ID (source_externalId)
  source TEXT NOT NULL,          -- Integration source (todoist, github, etc.)
  type TEXT NOT NULL,            -- Item type (task, event, email, etc.)
  external_id TEXT NOT NULL,     -- ID in the source system
  data JSONB NOT NULL,           -- Full payload from source
  embedding VECTOR(1536),        -- Optional: for semantic search
  synced_at TIMESTAMPTZ,         -- Last sync time
  created_at TIMESTAMPTZ,        -- First seen time
  UNIQUE(source, external_id)
);

-- Indexes for fast queries
CREATE INDEX idx_unified_source ON unified_items(source);
CREATE INDEX idx_unified_type ON unified_items(type);
CREATE INDEX idx_unified_synced ON unified_items(synced_at DESC);
CREATE INDEX idx_unified_data ON unified_items USING GIN(data);
```

---

## Data Sources by Type

### Tasks

| Source | Type | Key Fields |
|--------|------|------------|
| Todoist | `task` | content, due, priority, project |
| GitHub Issues | `task` | title, repo, labels, assignee |
| GitHub PRs | `task` | title, repo, reviewers, status |
| Linear | `task` | title, project, status, priority |
| Notion Tasks | `task` | title, database, properties |

### Events

| Source | Type | Key Fields |
|--------|------|------------|
| CalDAV (Google) | `event` | summary, start, end, location |
| CalDAV (iCloud) | `event` | summary, start, end, location |
| Home Assistant | `event` | entity_id, state, attributes |

### Messages

| Source | Type | Key Fields |
|--------|------|------------|
| IMAP (Gmail) | `email` | from, subject, date, preview |
| Slack | `message` | channel, author, text, ts |
| Discord | `message` | guild, channel, author, content |

### Content

| Source | Type | Key Fields |
|--------|------|------------|
| RSS | `article` | title, link, summary, published |
| Extension | `page` | url, title, content, meta |
| Extension | `receipt` | merchant, total, items, date |

---

## React Hooks for Dashboards

```typescript
// hooks/useDashboardData.ts

import { useState, useEffect } from 'react';
import { useCurrentEndpoint } from '@/lib/useCurrentEndpoint';

export function useTasks(filters?: { due?: string; priority?: number }) {
  const [tasks, setTasks] = useState([]);
  const [loading, setLoading] = useState(true);
  const endpoint = useCurrentEndpoint();

  useEffect(() => {
    const params = new URLSearchParams(filters);
    fetch(`${endpoint}/dashboard/tasks?${params}`)
      .then(res => res.json())
      .then(data => {
        setTasks(data.items);
        setLoading(false);
      });
  }, [endpoint, filters]);

  return { tasks, loading };
}

export function useActivityStream(limit = 50) {
  const [activity, setActivity] = useState([]);
  const endpoint = useCurrentEndpoint();

  useEffect(() => {
    // Initial fetch
    fetch(`${endpoint}/dashboard/activity?limit=${limit}`)
      .then(res => res.json())
      .then(data => setActivity(data.items));

    // WebSocket for real-time updates
    const ws = new WebSocket(`${endpoint.replace('http', 'ws')}/dashboard/stream`);
    ws.onmessage = (event) => {
      const newItem = JSON.parse(event.data);
      setActivity(prev => [newItem, ...prev.slice(0, limit - 1)]);
    };

    return () => ws.close();
  }, [endpoint, limit]);

  return activity;
}

export function useCalendarEvents(start: Date, end: Date) {
  const [events, setEvents] = useState([]);
  const endpoint = useCurrentEndpoint();

  useEffect(() => {
    fetch(`${endpoint}/dashboard/events?start=${start.toISOString()}&end=${end.toISOString()}`)
      .then(res => res.json())
      .then(data => setEvents(data.items));
  }, [endpoint, start, end]);

  return events;
}

export function useDashboardStats() {
  const [stats, setStats] = useState(null);
  const endpoint = useCurrentEndpoint();

  useEffect(() => {
    fetch(`${endpoint}/dashboard/stats`)
      .then(res => res.json())
      .then(setStats);

    // Refresh every minute
    const interval = setInterval(() => {
      fetch(`${endpoint}/dashboard/stats`)
        .then(res => res.json())
        .then(setStats);
    }, 60000);

    return () => clearInterval(interval);
  }, [endpoint]);

  return stats;
}
```

---

## Example Dashboard Widgets

### Task Widget

```tsx
function TaskWidget() {
  const { tasks, loading } = useTasks({ due: 'today' });

  if (loading) return <Skeleton />;

  return (
    <Card>
      <CardHeader>
        <CardTitle>Today's Tasks</CardTitle>
      </CardHeader>
      <CardContent>
        {tasks.map(task => (
          <TaskItem
            key={task.id}
            title={task.data.content || task.data.title}
            source={task.source}
            priority={task.data.priority}
            deepLink={getDeepLink(task)}
          />
        ))}
      </CardContent>
    </Card>
  );
}

function getDeepLink(task) {
  switch (task.source) {
    case 'todoist':
      return `todoist://task?id=${task.external_id}`;
    case 'github':
      return `https://github.com/${task.data.repo}/issues/${task.data.number}`;
    default:
      return null;
  }
}
```

### Activity Stream Widget

```tsx
function ActivityWidget() {
  const activity = useActivityStream(20);

  return (
    <Card>
      <CardHeader>
        <CardTitle>Recent Activity</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="space-y-2">
          {activity.map(item => (
            <ActivityItem
              key={item.id}
              icon={getSourceIcon(item.source)}
              title={formatActivity(item)}
              time={item.synced_at}
            />
          ))}
        </div>
      </CardContent>
    </Card>
  );
}

function formatActivity(item) {
  switch (item.type) {
    case 'task':
      return `Task: ${item.data.content || item.data.title}`;
    case 'email':
      return `Email from ${item.data.from}: ${item.data.subject}`;
    case 'event':
      return `Event: ${item.data.summary}`;
    default:
      return `${item.source}: ${item.type}`;
  }
}
```

---

## Deep Links

All items support deep links to open in their native apps:

```typescript
const DEEP_LINKS = {
  todoist: {
    task: (id) => ({
      mobile: `todoist://task?id=${id}`,
      web: `https://todoist.com/app/task/${id}`,
    }),
  },
  github: {
    issue: (repo, number) => ({
      web: `https://github.com/${repo}/issues/${number}`,
    }),
    pr: (repo, number) => ({
      web: `https://github.com/${repo}/pull/${number}`,
    }),
  },
  calendar: {
    event: (id, source) => {
      if (source === 'google') {
        return {
          web: `https://calendar.google.com/calendar/event?eid=${id}`,
        };
      }
      return null;
    },
  },
};
```

---

## Summary

| Need | Method | Latency | Real-time |
|------|--------|---------|-----------|
| UI widgets | REST API | ~50ms | No |
| Live updates | WebSocket | ~10ms | Yes |
| Complex queries | Direct SQL | ~100ms | No |
| Fresh source data | Nango Proxy | ~500ms | No |
| Semantic search | PostgreSQL + pgvector | ~200ms | No |

Use REST API for most dashboard needs, add WebSocket for real-time updates, and fall back to direct SQL for complex analytics.

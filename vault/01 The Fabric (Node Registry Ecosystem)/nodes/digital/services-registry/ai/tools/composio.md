---
name: composio-cli
description: Execute tools, build automations, and run scripts across 1000+ connected apps via Composio CLI. Use when you need to call Todoist, Strava, or any other connected app; build a new automation on the fly; create a script for the VPS agents to run; execute parallel tool calls; or discover new tool slugs. Prefer this over raw API calls — one command replaces pages of curl/Python. This is the integration layer for all One Mind agents (Legacy, Spartan, Empire).
allowed-tools:
  - Bash(composio:*)
  - Bash(python3:*)
  - Bash(bash:*)
  - Read
  - Write
---

# Composio CLI Skill

The integration layer for One Mind OS. One command reaches 1000+ apps without writing API code.

---

## Auth & Setup

```bash
# Check who is logged in
composio whoami

# Re-authenticate if needed
composio login

# Check CLI version
composio version
```

**Active account:** `zeus.vaults@pm.me` | org: `onemind` (ok_On-n6DvXI39s)
**Binary location:** `~/.composio/composio`
**Tool schemas cache:** `~/.composio/tool_definitions/`

---

## Core Command Patterns

### Execute a single tool

```bash
composio execute SLUG -d '{"param": "value"}'
```

### Execute in parallel (fires simultaneously — use for morning briefings, multi-fetch)

```bash
composio execute --parallel \
  STRAVA_LIST_ATHLETE_ACTIVITIES -d '{}' \
  TODOIST_FILTER_TASKS -d '{"filter": "today"}'
```

### Inspect schema before running (don't guess params)

```bash
composio execute SLUG --get-schema
```

### Preview without executing

```bash
composio execute SLUG --dry-run -d '{"param": "value"}'
```

### Discover slugs when you don't know the exact name

```bash
composio search "send slack message"
composio search "github create issue"
```

### List all tools for a connected app

```bash
composio tools list todoist --limit 200
composio tools list strava --limit 200
```

### Connect a new app

```bash
composio link <toolkit>   # e.g. composio link github
```

---

## Connected Apps & Tools

### TODOIST — 84 tools

**Key rules:**

- Always use `task_id` (NOT `id`) for task operations
- Use `project_id` (NOT `project`) for project operations
- To delete: **move to 🗑️ To Delete project `6gM4gJcq9HmXjmpg`** — never hard delete
- Always use v1/v2 variants (not DEPRECATED: prefix ones)

#### Core Task Operations

| Slug                                             | What it does                                        | Key params                                                                                       |
| ------------------------------------------------ | --------------------------------------------------- | ------------------------------------------------------------------------------------------------ |
| `TODOIST_CREATE_TASK`                            | Create a task                                       | `content`, `project_id`, `section_id`, `priority` (1-4), `due_string`, `description`, `labels`   |
| `TODOIST_GET_ALL_TASKS`                          | Fetch all incomplete tasks                          | _(none required)_                                                                                |
| `TODOIST_GET_TASK2`                              | Get single task by ID                               | `task_id`                                                                                        |
| `TODOIST_UPDATE_TASK`                            | Update task fields                                  | `task_id`, then any: `content`, `description`, `due_string`, `priority`, `labels`, `assignee_id` |
| `TODOIST_CLOSE_TASK_V1`                          | Complete/close a task                               | `task_id`                                                                                        |
| `TODOIST_REOPEN_TASK2`                           | Re-open a completed task                            | `task_id`                                                                                        |
| `TODOIST_MOVE_TASK`                              | Move task to different project/section              | `task_id`, `project_id` (or `section_id`)                                                        |
| `TODOIST_FILTER_TASKS`                           | Query tasks with Todoist filter syntax              | `filter` (e.g. `"today"`, `"p1"`, `"@urgent"`)                                                   |
| `TODOIST_BULK_CREATE_TASKS`                      | Create many tasks at once                           | `tasks` array                                                                                    |
| `TODOIST_QUICK_ADD_TASK`                         | Natural language task ("Call Zeus tomorrow at 3pm") | `text`                                                                                           |
| `TODOIST_DELETE_TASK`                            | Hard-delete a task                                  | `task_id` — **use move to 🗑️ instead**                                                           |
| `TODOIST_LIST_COMPLETED_TASKS`                   | See completed tasks                                 | `project_id` (optional)                                                                          |
| `TODOIST_GET_COMPLETED_TASKS_BY_COMPLETION_DATE` | Completed tasks by date range                       | `since`, `until`                                                                                 |
| `TODOIST_REORDER_TASKS`                          | Reorder tasks in a section                          | `items` array                                                                                    |

#### Project Operations

| Slug                                     | What it does                  | Key params                              |
| ---------------------------------------- | ----------------------------- | --------------------------------------- |
| `TODOIST_GET_ALL_PROJECTS`               | List all projects             | _(none)_ — returns `data.projects`      |
| `TODOIST_GET_PROJECT`                    | Get single project            | `project_id`                            |
| `TODOIST_GET_PROJECT_FULL`               | Project with sections/tasks   | `project_id`                            |
| `TODOIST_CREATE_PROJECT2`                | Create a project              | `name`, `parent_id` (optional), `color` |
| `TODOIST_UPDATE_PROJECT2`                | Update project                | `project_id`, `name`, `color`           |
| `TODOIST_ARCHIVE_PROJECT2`               | Archive a project             | `project_id`                            |
| `TODOIST_DELETE_PROJECT2`                | Delete project (hard)         | `project_id`                            |
| `TODOIST_SEARCH_PROJECTS`                | Search projects by name       | `query`                                 |
| `TODOIST_LIST_WORKSPACE_ACTIVE_PROJECTS` | All active workspace projects | _(none)_                                |

#### Section Operations

| Slug                        | What it does               | Key params                               |
| --------------------------- | -------------------------- | ---------------------------------------- |
| `TODOIST_LIST_SECTIONS`     | List sections in a project | `project_id`                             |
| `TODOIST_CREATE_SECTION_V1` | Create a section           | `name`, `project_id`, `order` (optional) |
| `TODOIST_UPDATE_SECTION`    | Rename a section           | `section_id`, `name`                     |
| `TODOIST_DELETE_SECTION2`   | Delete section + tasks     | `section_id`                             |

#### Label Operations

| Slug                      | What it does            | Key params               |
| ------------------------- | ----------------------- | ------------------------ |
| `TODOIST_LIST_LABELS`     | Get all personal labels | _(none)_                 |
| `TODOIST_CREATE_LABEL_V1` | Create a label          | `name`, `color`, `order` |
| `TODOIST_DELETE_LABEL_V1` | Delete a label          | `label_id`               |

#### Comments & Files

| Slug                        | What it does                | Key params                           |
| --------------------------- | --------------------------- | ------------------------------------ |
| `TODOIST_CREATE_COMMENT_V1` | Add comment to task/project | `content`, `task_id` or `project_id` |
| `TODOIST_GET_ALL_COMMENTS`  | Get comments                | `task_id` or `project_id`            |
| `TODOIST_UPLOAD_FILE`       | Upload file attachment      | `filename`, `file`                   |

#### Stats & Productivity

| Slug                             | What it does                     | Key params                     |
| -------------------------------- | -------------------------------- | ------------------------------ |
| `TODOIST_GET_PRODUCTIVITY_STATS` | Karma, completed counts, streaks | _(none)_                       |
| `TODOIST_GET_USER`               | Account info                     | _(none)_                       |
| `TODOIST_LIST_ACTIVITIES`        | Activity log                     | `page`, `limit`, `object_type` |
| `TODOIST_SYNC`                   | Full sync                        | _(see schema)_                 |

#### One Mind Key Project IDs

```
🗑️ To Delete:  6gM4gJcq9HmXjmpg   ← move here instead of deleting
```

#### Todoist Filter Syntax

```
today             # due today
p1                # priority 1 (urgent)
p2                # priority 2 (high)
@label-name       # by label
#Project Name     # by project
overdue           # past due
7 days            # due in next 7 days
assigned to: me   # assigned tasks
```

---

### STRAVA — 34 tools

**Key rules:**

- Most tools work with `{}` defaults (uses authenticated athlete)
- Use `activity_id` not `id` for activity operations
- `get_per_page` defaults to 30; use `before`/`after` (Unix timestamps) for date filtering
- HP domain: all fitness data routes here — update `LAST-WORKOUT.md` at codex root after pulls

#### Activity Operations

| Slug                             | What it does                                 | Key params                                                     |
| -------------------------------- | -------------------------------------------- | -------------------------------------------------------------- |
| `STRAVA_LIST_ATHLETE_ACTIVITIES` | List recent activities (paginated)           | `before`, `after` (Unix ts), `per_page`, `page`                |
| `STRAVA_GET_ACTIVITY`            | Full detail on one activity                  | `id` (activity_id)                                             |
| `STRAVA_GET_ACTIVITY_STREAMS`    | Time-series data (HR, power, cadence, speed) | `id`, `keys` (e.g. `"heartrate,cadence"`)                      |
| `STRAVA_GET_ACTIVITY_ZONES`      | HR and power zones for activity              | `id`                                                           |
| `STRAVA_LIST_ACTIVITY_LAPS`      | Lap-by-lap breakdown                         | `id`                                                           |
| `STRAVA_LIST_ACTIVITY_COMMENTS`  | Comments on activity                         | `activity_id`                                                  |
| `STRAVA_LIST_ACTIVITY_KUDOERS`   | Who kudoed it                                | `activity_id`                                                  |
| `STRAVA_CREATE_AN_ACTIVITY`      | Log a manual activity                        | `name`, `type`, `start_date_local`, `elapsed_time`, `distance` |
| `STRAVA_UPLOAD_ACTIVITY`         | Upload GPX/FIT file                          | `file`, `data_type` (`gpx`/`fit`/`tcx`)                        |
| `STRAVA_GET_UPLOAD`              | Check upload processing status               | `upload_id`                                                    |

#### Athlete Data

| Slug                               | What it does                   | Key params                       |
| ---------------------------------- | ------------------------------ | -------------------------------- |
| `STRAVA_GET_AUTHENTICATED_ATHLETE` | Profile, zones, weight         | _(none)_                         |
| `STRAVA_GET_ATHLETE_STATS`         | Totals (ytd, all-time, recent) | `id` (athlete ID: `125143734`)   |
| `STRAVA_GET_ZONES`                 | HR and power zone definitions  | _(none)_                         |
| `STRAVA_UPDATE_ATHLETE`            | Update weight, city, etc.      | `weight` (kg), `city`, `country` |
| `STRAVA_GET_EQUIPMENT`             | Gear detail (shoes, bike)      | `id` (equipment ID)              |

#### Routes & Segments

| Slug                                | What it does                     | Key params                                         |
| ----------------------------------- | -------------------------------- | -------------------------------------------------- |
| `STRAVA_LIST_ATHLETE_ROUTES`        | Your saved routes                | `id` (athlete ID)                                  |
| `STRAVA_GET_ROUTE`                  | Route detail                     | `id`                                               |
| `STRAVA_GET_ROUTE_STREAMS`          | Route GPS/elevation streams      | `id`                                               |
| `STRAVA_EXPORT_ROUTE_GPX`           | Export route as GPX file         | `id`                                               |
| `STRAVA_EXPORT_ROUTE_TCX`           | Export route as TCX              | `id`                                               |
| `STRAVA_EXPLORE_SEGMENTS`           | Find segments in a bounding box  | `bounds` (lat/lng box), `activity_type`            |
| `STRAVA_GET_SEGMENT`                | Segment detail                   | `id`                                               |
| `STRAVA_GET_SEGMENT_STREAMS`        | Segment GPS/elevation            | `id`                                               |
| `STRAVA_LIST_STARRED_SEGMENTS`      | Your starred segments            | _(none)_                                           |
| `STRAVA_STAR_SEGMENT`               | Star/unstar a segment            | `id`, `starred` (bool)                             |
| `STRAVA_LIST_SEGMENT_EFFORTS`       | Your efforts on a segment        | `segment_id`, `start_date_local`, `end_date_local` |
| `STRAVA_GET_SEGMENT_EFFORT`         | Detail on a segment effort       | `id`                                               |
| `STRAVA_GET_SEGMENT_EFFORT_STREAMS` | Stream data for a segment effort | `id`                                               |

#### Clubs

| Slug                              | What it does           | Key params |
| --------------------------------- | ---------------------- | ---------- |
| `STRAVA_LIST_ATHLETE_CLUBS`       | Clubs you're in        | _(none)_   |
| `STRAVA_GET_CLUB`                 | Club details           | `id`       |
| `STRAVA_LIST_CLUB_MEMBERS`        | Club member list       | `id`       |
| `STRAVA_LIST_CLUB_ACTIVITIES`     | Recent club activities | `id`       |
| `STRAVA_LIST_CLUB_ADMINISTRATORS` | Club admins            | `id`       |

**Athlete ID:** `125143734`

---

## Common One Mind Automation Patterns

### Morning Briefing (parallel pull — Strava + Todoist today)

```bash
composio execute --parallel \
  STRAVA_LIST_ATHLETE_ACTIVITIES -d '{"per_page": 1}' \
  TODOIST_FILTER_TASKS -d '{"filter": "today"}'
```

### Refresh LAST-WORKOUT.md at codex root

```bash
python3 ~/.claude/skills/composio-cli/scripts/update-last-workout.py
# OR pipe composio output directly:
composio execute STRAVA_LIST_ATHLETE_ACTIVITIES -d '{"per_page": 10}' | python3 ~/.claude/skills/composio-cli/scripts/update-last-workout.py
```

### Create a task from a workout (high-HR → rest day)

```bash
# Check last HR then conditionally create recovery task
composio execute STRAVA_LIST_ATHLETE_ACTIVITIES -d '{"per_page": 1}'
composio execute TODOIST_CREATE_TASK -d '{"content": "Rest Day - Recovery", "due_string": "tomorrow", "priority": 2, "labels": ["hp"]}'
```

### Move task to 🗑️ instead of deleting

```bash
composio execute TODOIST_MOVE_TASK -d '{"task_id": "TASK_ID", "project_id": "6gM4gJcq9HmXjmpg"}'
```

### Close a task when done

```bash
composio execute TODOIST_CLOSE_TASK_V1 -d '{"task_id": "TASK_ID"}'
```

### Get today + overdue tasks

```bash
composio execute TODOIST_FILTER_TASKS -d '{"filter": "today | overdue"}'
```

### Get all HP-pillar tasks

```bash
composio execute TODOIST_FILTER_TASKS -d '{"filter": "#25-49 HP"}'
```

---

## Ready-to-Run Scripts

All scripts live in `.claude/skills/composio-cli/scripts/`. They call `composio execute` internally — just run them directly.

| Script                     | What it does                                                                      | Run                                                       |
| -------------------------- | --------------------------------------------------------------------------------- | --------------------------------------------------------- |
| `update-last-workout.py`   | Fetches last Strava workout → writes `LAST-WORKOUT.md` at codex root              | `python3 update-last-workout.py`                          |
| `morning-briefing.sh`      | Parallel Strava + Todoist pull → formatted morning summary                        | `bash morning-briefing.sh`                                |
| `high-hr-recovery-task.py` | If last workout avg HR > threshold → creates Todoist recovery task for tomorrow   | `python3 high-hr-recovery-task.py`                        |
| `strava-weekly-summary.py` | Last 7 days of activities → stats summary, optionally creates Todoist review task | `python3 strava-weekly-summary.py`                        |
| `todoist-project-dump.py`  | Dumps all tasks from a project (text/json/markdown)                               | `python3 todoist-project-dump.py --project-id PROJECT_ID` |

### Script quick-start

```bash
# Change to scripts directory
cd /Users/zeuslegacy/Documents/onemind-codex/.claude/skills/composio-cli/scripts

# Update LAST-WORKOUT.md
python3 update-last-workout.py

# Morning briefing
bash morning-briefing.sh

# Recovery task if yesterday was intense
python3 high-hr-recovery-task.py --threshold 155

# Weekly training stats (and create a Todoist review task)
python3 strava-weekly-summary.py --create-task

# Dump all tasks from a project to markdown
python3 todoist-project-dump.py --project-id 6g6P3vvjx9643JrW --format markdown
```

---

## Scripted Automation (`composio run`)

Use for multi-step workflows that need logic between tool calls:```javascript
// File: /path/to/script.js → run with: composio run script.js

const activities = await execute("STRAVA_LIST_ATHLETE_ACTIVITIES", { per_page: 1 });
const last = activities.data.details[0];

if (last.average_heartrate > 160) {
await execute("TODOIST_CREATE_TASK", {
content: `Recovery needed — avg HR was ${last.average_heartrate}`,
due_string: "tomorrow",
priority: 2
});
}

````

### `composio run` injected globals:
| Function | Use |
|----------|-----|
| `execute(slug, params)` | Run a tool |
| `search(query)` | Find tool slugs |
| `proxy(app, endpoint, method, body)` | Raw API call to connected app |
| `experimental_subAgent(prompt)` | Inline AI reasoning step |
| `z` | Zod schema validation |

---

## Discovering New Integrations

When a new automation is needed and you don't know the slug:

```bash
# Step 1: Find the app and tools
composio search "slack send message"
composio search "google calendar create event"
composio search "notion create page"

# Step 2: Inspect what params it needs
composio execute SLACK_SEND_MESSAGE_AS_USER --get-schema

# Step 3: Connect the app if not yet linked
composio link slack
composio link google-calendar
composio link notion

# Step 4: Execute for real
composio execute SLACK_SEND_MESSAGE_AS_USER -d '{"channel": "#c-general", "text": "Hello from One Mind"}'
````

---

## Key Corrections (Known Gotchas)

| Tool                     | Wrong param                | Correct param |
| ------------------------ | -------------------------- | ------------- |
| `TODOIST_UPDATE_TASK`    | `id`                       | `task_id`     |
| `TODOIST_CLOSE_TASK_V1`  | `id`                       | `task_id`     |
| `TODOIST_MOVE_TASK`      | `project`                  | `project_id`  |
| Any Strava activity tool | `activity_id` in URL style | `id` as param |

---

## VPS Agent Usage (Legacy, Spartan, Empire)

Agents on the VPS can call Composio via shell. Useful patterns for VPS scripts:

```bash
# Run from VPS shell
~/.composio/composio execute TODOIST_CREATE_TASK -d '{"content": "Agent-generated task", "project_id": "PROJECT_ID"}'

# Chain with jq for parsing
~/.composio/composio execute STRAVA_LIST_ATHLETE_ACTIVITIES -d '{"per_page": 1}' | \
  python3 -c "import sys,json; d=json.load(sys.stdin); print(d['data']['details'][0]['average_heartrate'])"
```

---

## Strava Response Parsing

The Strava response is nested — watch for this:

```python
import json, subprocess

result = subprocess.run(
    ["composio", "execute", "STRAVA_LIST_ATHLETE_ACTIVITIES", "-d", '{"per_page": 5}'],
    capture_output=True, text=True
)
data = json.loads(result.stdout)
activities = data["data"]["details"]  # ← nested under data.details
```

---

## Add a New App Connection

To wire up a new app for on-the-fly use:

```bash
composio link <toolkit-slug>   # opens browser OAuth flow
composio search "<app name>"   # find the toolkit slug if unsure
```

Popular toolkits available: `github`, `slack`, `gmail`, `google-calendar`, `notion`, `airtable`, `linear`, `jira`, `discord`, `telegram`, `hubspot`, `stripe`, `shopify`, `youtube`, `twitter`, `spotify`

After linking, confirm with:

```bash
composio tools list <toolkit> --limit 50
```


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]

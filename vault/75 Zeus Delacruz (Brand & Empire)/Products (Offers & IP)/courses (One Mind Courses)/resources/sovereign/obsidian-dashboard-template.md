# Obsidian Dashboard Template

> *Copy this into any domain's `README.md` or a standalone `_DASHBOARD.md` file.*
> *Requires the Dataview plugin installed in Obsidian.*

---

## Life Dashboard

This dashboard uses Dataview queries to pull live data from across your vault. Replace file paths and tags to match your vault structure.

---

## Active Tasks by Priority

```dataview
TASK
FROM "99 Inbox (Queue)"
WHERE !completed
SORT priority ASC
LIMIT 20
```

---

## Goals — All Pillars

```dataview
TABLE
  pillar AS "Pillar",
  status AS "Status",
  file.mtime AS "Updated"
FROM #goal
WHERE status != "completed"
SORT pillar ASC
```

---

## Domain Health Check

```dataview
TABLE
  status AS "Status",
  agent AS "Agent",
  review AS "Review",
  file.mtime AS "Last Updated"
FROM "00-24 SO (Sovereign Operations)" OR "25-49 HP (Holistic Performance)" OR "50-74 LE (Legacy Evolution)" OR "75-99 GE (Generational Empire)"
WHERE file.name = "README"
SORT file.folder ASC
```

---

## Recent Captures (Last 7 Days)

```dataview
LIST
FROM "99 Inbox (Queue)"
WHERE file.ctime >= date(today) - dur(7 days)
SORT file.ctime DESC
LIMIT 15
```

---

## Weekly Reviews — Recent

```dataview
LIST
FROM "25-49 HP (Holistic Performance)/32 Rhythm (Systems & Productivity)/reviews"
WHERE file.name CONTAINS "weekly"
SORT file.ctime DESC
LIMIT 8
```

---

## GE Active Projects

```dataview
TABLE
  status AS "Status",
  file.mtime AS "Updated"
FROM "75-99 GE (Generational Empire)"
WHERE status = "active" AND file.name != "README"
SORT file.mtime DESC
LIMIT 10
```

---

## Customizing This Dashboard

**To add a new section:**

```dataview
TABLE
  [field1] AS "Label 1",
  [field2] AS "Label 2"
FROM "[path to folder]"
WHERE [condition]
SORT [field] ASC/DESC
```

**Useful FROM sources:**

- `FROM #tag` — all files with a tag
- `FROM "folder path"` — all files in a folder
- `FROM [[linked note]]` — all files linking to a note

**Useful WHERE conditions:**

- `WHERE status = "active"` — filter by frontmatter field
- `WHERE file.ctime >= date(today) - dur(7 days)` — recent files
- `WHERE !completed` — for task queries
- `WHERE contains(file.name, "keyword")` — filename search

**Useful fields:**

- `file.name` — filename without extension
- `file.mtime` — last modified date
- `file.ctime` — created date
- `file.folder` — containing folder path
- Any frontmatter field by its key name

---

*Part of the One Mind OS Command Center — Course 6*

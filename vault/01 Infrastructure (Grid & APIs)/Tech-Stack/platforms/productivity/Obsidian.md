# Obsidian

> **Knowledge Management & Personal Wiki**

## Overview

Obsidian is a powerful markdown-based knowledge management tool built on local files. It's used as the foundation for the OneMind Codex knowledge base.

## Status

- **Active**: ✅ Yes
- **Sync**: Git-based (no Obsidian Sync)
- **Last Updated**: January 25, 2026

## Use Cases

1. **OneMind Codex** - Central knowledge repository
2. **Daily Notes** - Journal and daily planning
3. **Project Documentation** - Technical docs
4. **Research Notes** - Learning and reference
5. **Meeting Notes** - Capture discussions

## Key Features

- Local-first (your files, your control)
- Markdown-based (portable, future-proof)
- Bidirectional linking ([[links]])
- Graph view (visualize connections)
- Plugins ecosystem
- Git-friendly

## Vault Structure

```
OneMind-Codex/
├── 00-24 SO (Sovereign Operations)/
│   ├── 00 Framework (Doctrine)/
│   ├── 01 Command (Dashboard)/
│   ├── 02 Agents (AI, Memory)/
│   ├── 06 Protocols (SOPs)/
│   ├── 08 Infrastructure/
│   └── ...
├── 25-49 HP (Holistic Performance)/
├── 50-74 LE (Legacy Evolution)/
├── 75-99 GE (Generational Endeavor)/
└── Templates/
```

## Configuration

### .obsidian/app.json

```json
{
  "theme": "obsidian",
  "baseFontSize": 16,
  "defaultViewMode": "preview",
  "livePreview": true,
  "spellcheck": true
}
```

### Core Plugins Enabled

- [x] Backlinks
- [x] Tags
- [x] Daily notes
- [x] Templates
- [x] Search
- [x] Graph view
- [x] Outline
- [x] Page preview
- [x] Slash commands

## Essential Plugins

### Community Plugins

| Plugin         | Purpose                       | Status       |
| -------------- | ----------------------------- | ------------ |
| **Templater**  | Advanced templates            | ✅ Installed |
| **Dataview**   | Query notes as database       | ✅ Installed |
| **Calendar**   | Calendar view for daily notes | ✅ Installed |
| **Kanban**     | Kanban boards                 | ✅ Installed |
| **Excalidraw** | Diagrams & sketches           | ✅ Installed |
| **Git**        | Version control               | ✅ Installed |
| **Tasks**      | Task management               | ✅ Installed |

## Templates

### Daily Note Template

```markdown
# {{date:YYYY-MM-DD}}

## 🎯 Focus

-

## 📋 Tasks

- [ ]

## 📝 Notes

## 🔗 Links

- [[Previous Day]]
- [[Next Day]]

---

Created: {{date:YYYY-MM-DD HH:mm}}
```

### Meeting Note Template

```markdown
# {{title}}

**Date**: {{date:YYYY-MM-DD}}
**Attendees**:
**Type**:

## Agenda

1.

## Notes

## Action Items

- [ ] @{{name}}:

## Follow-up

---

Tags: #meeting
```

## Dataview Queries

### All Tasks Due Today

```dataview
TASK
WHERE due = date(today)
AND !completed
SORT priority DESC
```

### Recent Notes

```dataview
TABLE file.mtime as "Modified"
FROM ""
SORT file.mtime DESC
LIMIT 10
```

### Project Status

```dataview
TABLE status, priority, due
FROM #project
WHERE status != "Done"
SORT priority DESC
```

## Linking Conventions

### Wikilinks

```markdown
[[Note Title]]
[[Note Title|Display Text]]
[[Folder/Note Title]]
[[Note Title#Heading]]
```

### Tags

```markdown
#project
#meeting
#idea
#reference
#status/active
#priority/high
```

## Git Workflow

### .gitignore

```
.obsidian/workspace.json
.obsidian/workspace-mobile.json
.obsidian/plugins/
.trash/
```

### Sync Setup

```bash
# In vault directory
git init
git remote add origin git@github.com:Zeus-Delacruz/OneMind-Codex.git
git add .
git commit -m "Initial commit"
git push -u origin main
```

### Auto-Commit (via Git plugin)

- Auto-pull on startup
- Auto-commit every 5 minutes
- Auto-push after commit

## Hotkeys

| Shortcut      | Action                   |
| ------------- | ------------------------ |
| `Cmd+O`       | Quick switcher           |
| `Cmd+P`       | Command palette          |
| `Cmd+E`       | Toggle edit/preview      |
| `Cmd+\`       | Toggle sidebar           |
| `Cmd+Shift+F` | Search in all files      |
| `Cmd+Click`   | Open link in new pane    |
| `[[`          | Create wikilink          |
| `Alt+Enter`   | Follow link under cursor |

## Mobile Setup

### Obsidian Mobile

- iOS and Android apps
- Same vault synced via Git
- Working Copy app (iOS) for Git

### Quick Capture

- iOS Shortcuts → Append to daily note
- Share extension → Save links

## Integration with Legacy AI

### Potential Use Cases

1. **Knowledge Search** - Legacy queries Obsidian vault
2. **Note Creation** - "Create a meeting note for..."
3. **Daily Summary** - Pull from daily notes
4. **Research Agent** - Add notes from web research
5. **Memory Sync** - Sync agent memory to Obsidian

### API Access (Local REST API plugin)

```python
# Read note
GET http://localhost:27124/vault/path/to/note.md

# Create note
PUT http://localhost:27124/vault/path/to/new-note.md
Content-Type: text/markdown

# Your note content
```

## Best Practices

- ✅ Use atomic notes (one idea per note)
- ✅ Link liberally (create connections)
- ✅ Use templates for consistency
- ✅ Tag for discoverability
- ✅ Review graph view periodically
- ✅ Daily notes for continuity
- ✅ Weekly reviews
- ✅ Back up vault (Git + cloud backup)

## Resources

- [Obsidian Help](https://help.obsidian.md)
- [Obsidian Forum](https://forum.obsidian.md)
- [Plugin Directory](https://obsidian.md/plugins)

---

**Related**: [Raindrop](./Raindrop.md) | [Todoist](./Todoist.md)

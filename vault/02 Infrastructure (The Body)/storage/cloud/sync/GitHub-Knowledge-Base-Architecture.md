---
id: ui-10-20260111053926
title: "GitHub as Agent Knowledge Base"
type: concept
domain: SO
domain: 00-24
folder: "00-24 SO (Sovereign Operations)"
status: active
created: 2026-01-11
modified: 2026-01-11
owner: zeus
tags: [domain/so, folder/10]
todoist_project: "00-24 SO (Sovereign Operations)"
ai_summary: "GitHub as Agent Knowledge Base"
---

# GitHub as Agent Knowledge Base

How Letta agents access and edit the One Mind vault using GitHub as the single source of truth.

---

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         THE FLOW                                 │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│   You: "What's my morning routine?"                              │
│        │                                                         │
│        ▼                                                         │
│   Legacy-AI calls search_vault("morning routine", pillar="HP")   │
│        │                                                         │
│        ▼                                                         │
│   GitHub Search API:                                             │
│   GET /search/code?q=morning+routine+repo:Zeus-Delacruz/One-Mind-OS+path:10-19
│        │                                                         │
│        ▼                                                         │
│   Returns matching files:                                        │
│   - 10-19 HP/12 Habits/Morning-Routine.md                        │
│   - 10-19 HP/14 Protocols/Daily-Protocol.md                      │
│        │                                                         │
│        ▼                                                         │
│   GitHub Contents API fetches file content                       │
│        │                                                         │
│        ▼                                                         │
│   Agent reads content, answers your question                     │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## Pillar Path Mapping

The numbering system enables targeted searches:

| Pillar | Path Prefix | Scope |
|--------|-------------|-------|
| Command | `path:00-09` | Meta, framework, governance, protocols |
| HP (Holistic Performance) | `path:10-19` | Health, habits, fitness, wellness |
| LE (Legacy & Estate) | `path:20-29` | Family, home, legacy planning |
| GE (Generational Empire) | `path:30-39` | Business, finance, career, investments |
| IT (Innovative Technology) | `path:40-49` | Tech stack, code, infrastructure |

**Example searches:**
```
"budget"           + path:30-39  → Only GE financial docs
"workout"          + path:10-19  → Only HP fitness docs
"deployment"       + path:40-49  → Only IT infrastructure docs
"family trust"     + path:20-29  → Only LE estate docs
```

---

## Letta Tool: search_vault

```python
def search_vault(query: str, pillar: str = None, limit: int = 5):
    """
    Search the One Mind vault via GitHub API.

    Args:
        query: Keywords to search for
        pillar: Optional filter (HP, LE, GE, IT, Command)
        limit: Max results to return (default 5)

    Returns:
        List of matching documents with path and content
    """

    # Map pillar to path prefix
    pillar_paths = {
        "Command": "00-09",
        "HP": "10-19",
        "LE": "20-29",
        "GE": "30-39",
        "IT": "40-49"
    }

    # Build GitHub search query
    search_query = f"{query} repo:Zeus-Delacruz/One-Mind-OS extension:md"
    if pillar and pillar in pillar_paths:
        search_query += f" path:{pillar_paths[pillar]}"

    # Search via GitHub API
    headers = {"Authorization": f"token {GITHUB_TOKEN}"}
    response = requests.get(
        "https://api.github.com/search/code",
        params={"q": search_query, "per_page": limit},
        headers=headers
    )

    results = response.json()["items"]

    # Fetch content for each match
    documents = []
    for item in results:
        content_response = requests.get(
            f"https://api.github.com/repos/Zeus-Delacruz/One-Mind-OS/contents/{item['path']}",
            headers=headers
        )
        content = base64.b64decode(content_response.json()["content"]).decode()

        documents.append({
            "path": item["path"],
            "content": content
        })

    return documents
```

---

## Letta Tool: edit_vault

```python
def edit_vault(file_path: str, new_content: str, commit_message: str):
    """
    Edit a file in the One Mind vault via GitHub API.

    Args:
        file_path: Path to file (e.g., "10-19 HP/12 Habits/Morning-Routine.md")
        new_content: New file content
        commit_message: Descriptive commit message

    Returns:
        Commit SHA and URL
    """

    headers = {"Authorization": f"token {GITHUB_TOKEN}"}

    # Get current file SHA (required for update)
    current = requests.get(
        f"https://api.github.com/repos/Zeus-Delacruz/One-Mind-OS/contents/{file_path}",
        headers=headers
    ).json()

    # Update file
    response = requests.put(
        f"https://api.github.com/repos/Zeus-Delacruz/One-Mind-OS/contents/{file_path}",
        headers=headers,
        json={
            "message": commit_message,
            "content": base64.b64encode(new_content.encode()).decode(),
            "sha": current["sha"],
            "committer": {
                "name": "Legacy-AI",
                "email": "legacy-ai@onemind.os"
            }
        }
    )

    return {
        "sha": response.json()["commit"]["sha"],
        "url": response.json()["commit"]["html_url"]
    }
```

---

## Version Control Benefits

When agents edit documents, Git provides:

| Feature | What You Get |
|---------|--------------|
| **Full History** | Every change ever made, forever |
| **Attribution** | Who changed it (agent commits as "Legacy-AI") |
| **Timestamps** | Exactly when each change occurred |
| **Diffs** | Line-by-line view of what changed |
| **Rollback** | One command to undo any change |
| **Blame** | See who wrote each line |

**Useful commands:**
```bash
# See all agent changes
git log --author="Legacy-AI" --oneline

# See what changed in last commit
git diff HEAD~1

# Undo agent's last edit
git revert HEAD

# See history of specific file
git log --follow -p "10-19 HP/12 Habits/Morning-Routine.md"

# See who changed each line
git blame "10-19 HP/12 Habits/Morning-Routine.md"
```

---

## Rate Limits

| API | Limit | Practical Meaning |
|-----|-------|-------------------|
| Search API | 30 requests/min | 30 searches per minute |
| Contents API | 5,000 requests/hour | ~80 file reads per minute |

For personal use, these limits are **more than enough**. You'd have to run constant automated searches to hit them.

---

## Search Type: Keyword Only

GitHub provides **keyword/exact match** search, not semantic search.

**What works:**
```
Query: "morning routine"
Finds: Files containing "morning" AND "routine"
```

**What doesn't work:**
```
Query: "how do I start my day"
Won't find: Morning-Routine.md (different words, same concept)
```

**Mitigation:** Your numbered folder structure + consistent naming conventions compensate for lack of semantic search. The pillar filter narrows results significantly.

---

## Why GitHub-Only Works for One Mind

| Factor | Benefit |
|--------|---------|
| Numbered folders | Built-in categorization replaces database tags |
| Markdown files | Text-searchable, version-controllable |
| Obsidian sync | Local editing + git push = always current |
| Agent edits | Full audit trail via commits |
| Zero cost | No database to maintain or pay for |
| Zero infra | No vector DB, no embeddings, no sync jobs |

---

## Upgrade Path

If you later need semantic search:

```
Phase 1 (Now):     GitHub Search API — keyword search
Phase 2 (Future):  Add Supabase pgvector — semantic search
```

The tool interface stays the same — only the backend changes. Agents don't know the difference.

---

## Two-System Architecture

| System | Contents | Storage | Agent Access |
|--------|----------|---------|--------------|
| **One-Mind-OS** | Markdown, code, configs | GitHub | Full read/write via API |
| **One-Mind-Files** | PDFs, images, media | Cloud storage (TBD) | Read-only or separate tools |

Git handles text well. Binary files (images, PDFs) live separately to avoid repo bloat.

---

## Summary

- **Search**: GitHub Search API with pillar path filtering
- **Read**: GitHub Contents API
- **Write**: GitHub Contents API with commits
- **Version Control**: Full Git history, diffs, rollback
- **Cost**: $0
- **Maintenance**: None
- **Trade-off**: Keyword search only (no semantic), mitigated by folder structure


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]

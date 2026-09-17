# Codex Skill: Search — Tiered Loading Protocol

How to find information in the codex efficiently. Uses progressive disclosure to minimize context consumption.

> Inspired by OpenViking's L0/L1/L2 protocol and Letta Context Repositories' `system/` directory pattern.

---

## The Four Levels

### System: Pinned Context (always loaded — ~500 tokens)

Files in `_codex/system/` are **pre-loaded** at the start of every interaction. You don't need to search for these — they're already in your context.

**Contains**: Active goals, owner identity, critical reference data.
**Rule**: If you find yourself searching for the same file in >50% of sessions, it should be in `system/`.

### Level 0: Orientation (cheapest — ~200 tokens)

Start broad. Read the domain group or domain README to understand what's available.

**When**: You don't know which domain the answer is in.

**Steps**:
1. If you know the domain group (UI/HP/LE/GE), read its README
2. Scan the domain folder names within the domain group
3. Identify the 1-2 most likely domains

**Example**: "Where's the VPS documentation?"
→ Read `00-24 SO/README.md`
→ See "08 Infrastructure" exists
→ That's your target domain

### Level 1: Discovery (medium — ~500 tokens)

Narrow down. Read the domain README or scan file names.

**When**: You know the domain but not the exact file.

**Steps**:
1. Read the domain folder's README (e.g., `08 Infrastructure/README.md`)
2. List file names in the folder
3. Read `ai_summary` from frontmatter of promising files — this is the primary progressive disclosure mechanism. Every file should have one.
4. Identify the 1-3 most relevant files

**Example**: "What's the VPS setup process?"
→ List files in `08 Infrastructure/`
→ See `vps-setup.md`, `docker-compose-reference.md`, etc.
→ Target: `vps-setup.md`

### Level 2: Deep Read (most expensive — full file)

Read the actual content. Only do this when you've confirmed relevance.

**When**: You've identified the specific file(s) you need.

**Steps**:
1. Read the target file fully
2. Follow `related:` links in frontmatter if you need broader context
3. Follow `[[wiki-links]]` within the content if they lead to needed detail

---

## Search Strategy by Question Type

| Question Type | Strategy |
| --- | --- |
| "Where is X?" | L0 → L1 (folder scan) |
| "How do I do X?" | L0 → L1 → L2 (find the SOP) |
| "What did we decide about X?" | Search for `type: decision` + keyword |
| "What's the status of X?" | L1 in the relevant domain → L2 on the project file |
| "Tell me about person X" | Search for `type: person` + name |
| Cross-domain question | Search across all domain groups with keyword, then L2 on best matches |

## Using Runtime Search (if available)

When the runtime provides search tools (e.g., OpenClaw's `memory_search`):

1. **Filter by domain first** when possible — reduces noise
2. **Use hybrid search** (keyword + semantic) for best results
3. **Check tags** — `tags:` in frontmatter helps cross-domain discovery
4. **Prefer recent over old** — temporal relevance matters

## Caching Strategy

If you discover a file is highly relevant and likely to be referenced again:
- Note its path for re-use in the conversation
- Quote key facts rather than re-reading the whole file
- Build a "working set" of 3-5 most relevant files for the current task

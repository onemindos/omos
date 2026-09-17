---
name: onemind-consolidate
description: Merge scattered daily notes, agent memory logs, and temporary captures into curated permanent knowledge. Use during weekly or monthly reviews to promote ephemeral notes to stable codex entries.
---

# Codex Skill: Consolidate

How to merge scattered daily notes, agent memory logs, and temporary captures into curated, permanent knowledge. Inspired by ClawVault's fact extraction and Letta's memory defragmentation.

---

## When to Consolidate

- Weekly (during weekly review)
- When a topic has 5+ scattered notes that should be one curated document
- When daily logs accumulate insights that deserve permanent homes
- When an agent flags knowledge that keeps being re-discovered

## The Consolidation Process

### Step 1: Identify Candidates

Look for:
- Daily logs older than 7 days that contain facts not yet captured elsewhere
- Multiple notes in the same domain that overlap significantly
- Notes tagged `status/draft` that have matured
- Inbox items that are actually updates to existing notes

### Step 2: Extract Facts

For each candidate, pull out:
- **Permanent facts**: Things that are true indefinitely ("VPS runs Ubuntu 24.04")
- **Decisions**: Choices that were made with rationale
- **Preferences**: "Zeus prefers X over Y"
- **Procedures**: Steps that solved a problem (potential SOPs)
- **Contacts**: People mentioned with their context

### Step 3: Merge or Create

For each extracted fact:
- **If a note already exists on this topic** → update that note with the new information
- **If no note exists** → create a new note in the correct domain using the appropriate template
- **If multiple notes overlap** → merge them into one comprehensive note, archive the originals

### Step 4: Clean Up

- Archive or delete the source daily logs (keep the last 30 days)
- Remove `status/draft` tags from notes that are now curated
- Update `modified:` dates on edited notes
- Verify `[[wiki-links]]` and `related:` fields

---

## Consolidation Patterns

### Pattern: Daily Log → Domain Notes

```
memory/2026-03-15.md: "Discovered NATS needs explicit stream creation for JetStream"
memory/2026-03-16.md: "Created DIGITAL stream with subjects digital.agent.>"
memory/2026-03-17.md: "Added retention policy of 7 days to DIGITAL stream"

→ Consolidate into:
   08 Infrastructure/nats-jetstream-setup.md
   (comprehensive note with all stream configurations)
```

### Pattern: Scattered Notes → Single Reference

```
06 Inbox/vps-ip-note.md
08 Infrastructure/server-setup-draft.md
08 Infrastructure/ssh-access-notes.md

→ Merge into:
   08 Infrastructure/vps-reference.md
   (single authoritative VPS reference document)
```

### Pattern: Repeated Questions → SOP

```
Agent keeps looking up "how to restart OpenClaw" in different logs

→ Create:
   06 Protocols/restart-openclaw.md
   (SOP with clear steps, so it's found on first search next time)
```

---

## Quality Standards for Curated Notes

A consolidated note should:
- Have complete frontmatter
- Be self-contained (reader doesn't need the source logs)
- Use clear section headers
- Link to related notes with `[[wiki-links]]`
- Answer the question someone would search for
- Be in the correct domain folder (not inbox)

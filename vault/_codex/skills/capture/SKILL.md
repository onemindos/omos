# Codex Skill: Capture

How to add new information to the codex. Any incoming data — ideas, notes, decisions, observations, agent outputs — enters through this process.

---

## The Capture Flow

```
Input (any source)
  → Determine type (note, decision, SOP, log, person, etc.)
  → Apply template from _codex/templates/
  → Add frontmatter
  → Save to 06 Inbox (Queue)/ OR to the known domain folder
```

## Steps

1. **Classify the input** — What type of information is this?
   - Idea or observation → `note`
   - A choice that was made → `decision`
   - A process to follow → `sop`
   - Information about a person → `person`
   - A project or initiative → `project`
   - A time-stamped event → `log`

2. **Choose the template** — Use the matching template from `_codex/templates/`:
   - Read the template
   - Fill in the frontmatter fields
   - Write the content in the template's structure

3. **Generate the filename**:
   - Use lowercase with hyphens: `my-new-note.md`
   - For dated entries: `YYYY-MM-DD-slug.md`
   - Keep it descriptive but concise

4. **Determine the destination**:
   - If the domain is obvious → save directly to that domain folder
   - If uncertain → save to `06 Inbox (Queue)/`
   - Never create new top-level folders

5. **Add connections**:
   - Add `related:` links in frontmatter to connected notes
   - Add `[[wiki-links]]` in the body to reference existing knowledge
   - Add appropriate tags

## Capture Sources

| Source | How to Handle |
| --- | --- |
| Human dictation / voice | Transcribe → classify → template → save |
| Agent observation | Create note with `agent/{name}` tag |
| External data (API, sensor) | Create log entry with source metadata |
| Screenshot / image | Save image to `assets/`, create note linking to it |
| Email / message | Extract key info, create note, link to source if available |
| Daily reflection | Use daily template, save to agent's memory or `29 Life Systems/` |

## Quality Checklist

Before saving, verify:

- [ ] Frontmatter has at minimum: `title`, `type`, `created`, `tags`
- [ ] Content is clear enough that someone reading it in 6 months would understand
- [ ] Tags use the `namespace/value` pattern
- [ ] Any referenced notes are wiki-linked
- [ ] File is saved in the right location (or Inbox)

# Organize Workflows — CODE× Method

## Before You Touch Anything — Read Every Time

1. `CLAUDE.md` at vault root — naming rules, what not to touch
2. `_codex/CONVENTIONS.md` — frontmatter spec, tagging system
3. The domain's existing `_README.md` — understand what already exists

---

## Golden Rules

- **NEVER DELETE ANYTHING WITHOUT ZEUS'S EXPLICIT PERMISSION.** Not files. Not folders. Not duplicates. Not empty shells. Nothing. Flag it, tell Zeus, wait for yes.
- Never rename the 4 domain group folders without Zeus's explicit approval
- Never move files between domain groups without telling Zeus first
- One file, one location — never duplicate content across folders
- `_README.md` at every domain root is required — always update it after structural changes

---

## Folder Naming Convention

Format: `Name (Subtitle)/`

```
✅ Content (Content Engine)/
✅ Finance (Real Estate Finance)/
❌ Content/
❌ Finance/
❌ Misc/
```

Every folder must have a subtitle in parentheses. No underscores. No vague names (`Misc/Other/Stuff`). No duplicates at the same level.

**No loose files rule:** Every `.md` file must live inside a named subfolder. A file directly at a domain root (other than `_README.md`) is loose and must be routed immediately.

**Depth rule:** Go 2-3 levels deep. Not more. If you need level 4, reconsider the structure.

---

## Loose File Protocol

1. Read the file — check frontmatter `id:` prefix (e.g., `hp-32` = HP domain 32)
2. Find the most specific existing subfolder it belongs in
3. Move it — if no suitable subfolder exists, create one
4. Update the domain's `_README.md` if the subfolder is new

```bash
# Find loose files at domain root (should only be _README.md)
find "DOMAIN_PATH" -maxdepth 1 -type f -name "*.md" ! -name "_README.md"
```

---

## Full Domain Organization Pass

**Phase 1 — Audit**
```bash
find "DOMAIN" -maxdepth 3 -type d | sort
find "DOMAIN" -name "*.md" | grep -v "_README" | sort
find "DOMAIN" -maxdepth 2 -type f -name "*.md" ! -name "_README.md"
```

**Phase 2 — Assess:** Map each subfolder to its purpose. Identify: duplicates, misplaced files, empty shells, missing subfolders. Tell Zeus what you plan to do before moving anything structural.

**Phase 3 — Clean:** Fix duplicate folders first (merge, don't delete until merged). Route loose files. Rename folders that don't follow convention.

**Phase 4 — Build:** Add missing subfolders as placeholders. Go 2-3 levels deep.

**Phase 5 — Document:** Update `_README.md`. Verify no loose files remain. Run final audit.

---

## Duplicate Folder Protocol

1. List contents of both folders
2. Tell Zeus: "I found duplicate folders [A] and [B]. I plan to merge everything into [A] and remove the empty shell of [B]. Approve?"
3. Wait for approval. Once approved: move unique files → surviving folder → `rmdir` the empty shell
4. **Never use `rm -rf` on any folder. Use `rmdir` only** — it refuses to delete non-empty folders.

---

## What to Never Touch Without Zeus's Approval

| Thing | Why |
|-------|-----|
| Domain group folder names | Locked after naming sessions |
| `ONEMIND-CODEX.md` at vault root | Canonical master document |
| `CLAUDE.md`, `AGENTS.md` | Agent instruction files |
| `_codex/` folder structure | Changes affect all agents |
| Any file with `status: locked` | Explicitly protected |
| Moving content between SO/HP/LE/GE | Cross-domain routing needs Zeus sign-off |

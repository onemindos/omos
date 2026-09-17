# Codex Skill: Maintain

How to keep the vault healthy. Run this during monthly reviews or when the vault feels disorganized.

---

## Health Audit Checklist

### 1. Inbox Hygiene

- Count items in `06 Inbox (Queue)/`
- If > 20 items, run `_codex/skills/organize.md` to process them
- Items older than 30 days in inbox need immediate triage: organize or archive

### 2. Orphan Detection

Find notes that have no incoming links and no tags:
- These are "orphaned" — they exist but nothing references them
- For each orphan, either:
  - Add `[[wiki-links]]` from related notes
  - Add appropriate tags
  - Move to the correct domain if it's misplaced
  - Archive if it's no longer relevant

### 3. Frontmatter Completeness

Scan for notes missing required frontmatter fields:
- Every file needs at minimum: `title`, `type`, `created`, `tags`
- Fix any files missing these fields
- Add `ai_summary` to frequently-accessed files (helps L1 scanning)

### 4. Broken Links

Check for `[[wiki-links]]` that point to non-existent notes:
- If the target was renamed → update the link
- If the target was deleted → remove the link or recreate the note
- If it's a planned note → create a stub with `status/planned`

### 5. Stale Content

Find notes with `modified:` dates older than 6 months:
- Review: is the content still accurate?
- Update if outdated
- Archive if no longer relevant
- Mark as `status/verified` if still good

### 6. Duplicate Detection

Look for notes covering the same topic in different locations:
- Merge them using `_codex/skills/consolidate.md`
- Keep one authoritative version, archive or delete duplicates
- Update all links to point to the surviving note

### 7. Tag Consistency

- Ensure tags follow the `namespace/value` pattern
- Look for variant spellings (`domain/UI` vs `domain/so`)
- Merge or standardize inconsistent tags

### 8. README Coverage

Every domain folder should have a `README.md`:
- Check each folder under the 3 domains + 1 operating layer
- Create READMEs for any folders missing them
- Update existing READMEs if the folder contents have changed significantly

---

## Maintenance Schedule

| Task | Frequency | Typical Time |
| --- | --- | --- |
| Inbox hygiene | Weekly | 5 min |
| Orphan detection | Monthly | 10 min |
| Frontmatter completeness | Monthly | 10 min |
| Broken links | Monthly | 5 min |
| Stale content review | Quarterly | 30 min |
| Duplicate detection | Quarterly | 15 min |
| Tag consistency | Quarterly | 10 min |
| README coverage | After adding new domains | 5 min |

## Severity Levels

| Issue | Severity | Action |
| --- | --- | --- |
| Missing frontmatter | Medium | Fix immediately |
| Broken links | Medium | Fix or create stubs |
| Orphaned notes | Low | Add context when convenient |
| Stale content | Low | Review during quarterly audit |
| Inbox overflow (>20 items) | High | Process immediately |
| Duplicate notes | Medium | Consolidate during next review |

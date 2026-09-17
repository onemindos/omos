# Verifying Report Pack

## Required
- [ ] SKILL.md readable and loaded
- [ ] All 4 templates present in `templates/`
- [ ] Agent can name the 4 report types
- [ ] Agent can list data sources for reports

## Functional Test
- [ ] Generate a mini domain health check for any one domain (just counts + last modified)
- [ ] Output follows the template structure

## Dependencies Check
- [ ] `_codex/system/active-goals.md` accessible (needed for quarterly review)
- [ ] `_codex/memory/learning/signals.jsonl` accessible (needed for pattern extraction)

## Result
✅ All checks pass = report generation ready
⚠️ Memory not deployed = reports work but without learning data
❌ Templates missing = cannot generate formatted reports


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]

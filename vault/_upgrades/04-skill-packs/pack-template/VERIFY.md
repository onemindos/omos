# Pack Template: VERIFY.md

> Post-installation checklist. Run after INSTALL.md to confirm the skill is working.

---

# Verifying {SKILL_NAME} Pack

## Required Checks

### 1. File Structure ✓
- [ ] `README.md` exists in skill directory
- [ ] `SKILL.md` exists and is readable
- [ ] `INSTALL.md` exists
- [ ] `VERIFY.md` exists (this file)

### 2. Dependencies ✓
- [ ] `_codex/system/identity.md` is accessible
- [ ] `_codex/system/active-goals.md` is accessible
- [ ] `_codex/system/preferences.md` is accessible

### 3. Skill Comprehension ✓
- [ ] Agent can describe what this skill does in one sentence
- [ ] Agent can list the trigger conditions
- [ ] Agent can name the workflows available

### 4. Functional Test ✓
- [ ] Execute a simple test case: {describe test}
- [ ] Confirm output matches expected format
- [ ] Confirm output respects preferences from `preferences.md`

## Optional Checks

### 5. Memory Integration
- [ ] `_codex/memory/` exists (if using memory features)
- [ ] Agent can write to `memory/learning/signals.jsonl`

### 6. Cross-Skill Compatibility
- [ ] Skill doesn't conflict with other installed skills
- [ ] Skill's outputs can be consumed by downstream skills

---

## Result

| Status | Meaning |
|--------|---------|
| ✅ All required checks pass | Skill is fully operational |
| ⚠️ Some optional checks fail | Skill works but with reduced capability |
| ❌ Any required check fails | Skill is NOT ready — troubleshoot |

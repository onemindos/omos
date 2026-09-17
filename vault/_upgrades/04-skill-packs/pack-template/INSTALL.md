# Pack Template: INSTALL.md

> AI agents read this file to self-install the skill. Replace {PLACEHOLDERS}.

---

# Installing {SKILL_NAME} Pack

**Version**: 1.0.0
**For**: All OneMind agents (SO, HP, LE, GE)

---

## Phase 1: System Check

Before installing, verify:

- [ ] `_codex/system/` exists and contains identity files
- [ ] `_codex/skills/` directory is accessible
- [ ] Agent has read/write access to the vault

## Phase 2: Context Loading

Load these files into your context:

1. Read `_codex/system/identity.md` — Know who you're serving
2. Read `_codex/system/active-goals.md` — Know current priorities
3. Read `_codex/system/preferences.md` — Know how to communicate

## Phase 3: Skill Activation

1. Read `SKILL.md` in this directory — these are your operating instructions
2. Internalize the workflows and decision trees
3. Note the trigger conditions — when should you activate this skill?

## Phase 4: Verification

Run through `VERIFY.md` to confirm everything is working:

1. Execute each test case
2. Confirm expected outputs
3. Report any failures

## Phase 5: Confirmation

If all checks pass, respond:
```
✅ {SKILL_NAME} pack installed successfully.
- Skill loaded: {SKILL_NAME} v1.0.0
- Workflows available: {list workflows}
- Ready for: {what the agent can now do}
```

If any checks fail, respond:
```
⚠️ {SKILL_NAME} pack installation incomplete.
- Failed checks: {list failures}
- Suggested fix: {what to try}
```

# Hook: SecurityValidate

**Event**: `execute-pre`
**Blocking**: Yes (must complete before command runs)
**Fail Behavior**: Abort command

## Instructions

Before executing any shell command or file modification:

1. Read `_system/security/patterns.yaml`
2. Follow the validation flow in `_system/security/validator.md`
3. Check against bash_patterns, path_protection, vault_rules
4. Take appropriate action: BLOCK / CONFIRM / ALERT / ALLOW
5. Log event to `_system/security/events/YYYY-MM.jsonl`

See `_system/security/validator.md` for full instructions.


---

## Related

- [[00-24 SO (Sovereign Operations)/03 Command (Operations & Navigation) (The Controls)/03 Command|03 Command]]

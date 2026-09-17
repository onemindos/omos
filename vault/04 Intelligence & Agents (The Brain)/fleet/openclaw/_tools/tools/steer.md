---
name: OpenClaw — Steer
updated: 2026-05-11
---

# Steer

`/steer` sends guidance to an **already-active run** without starting a new turn or killing the session. The most powerful real-time control primitive in OpenClaw.

**Official docs:** [docs.openclaw.ai/tools/steer](https://docs.openclaw.ai/tools/steer.md)

---

## Basic Usage

```bash
/steer prefer the smaller patch and keep the tests focused
/tell summarize before making the next tool call
```

- Targets only the current session's active run
- Works independently of the session's `/queue` mode
- Does **not** start a new run if the session is idle
- Injects guidance at the next supported runtime boundary

---

## Steer vs Queue

| Command | When to use |
|---|---|
| `/steer <message>` | Guide the active run right now |
| `/queue steer` | Future normal messages should steer active runs by default |
| `/queue collect` | New messages should wait for a later turn instead |
| `/queue followup` | New messages queue as follow-ups, not interruptions |

---

## Sub-Agents

To steer a child run (not the current session):

```bash
/subagents steer 2 focus only on the API surface
```

Top-level `/steer` always targets the current session's active run only.

---

## ACP Sessions

```bash
/acp steer --session agent:main:acp:session123 tighten the repro
```

---

## One Mind Use Case

Zeus can watch Legacy working on a complex task and nudge it mid-stream:

```
/steer actually prioritize the Telegram channel first, Discord can wait
/steer before you write the file, show me the diff of what you're planning
/steer stop — the approach is wrong, let's use X instead
```

This is how you maintain control over long-running agent tasks without restarting the session and losing all the context Legacy has built up.

---

## Related

- [concepts/queue](../concepts/README.md)
- [tools/subagents](README.md)

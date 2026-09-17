---
name: OpenClaw — BTW (Side Questions)
updated: 2026-05-11
---

# BTW — Side Questions

`/btw` lets you ask a quick side question against the current session context without that question or answer ever entering the transcript. `/side` is an alias.

**Official docs:** [docs.openclaw.ai/tools/btw](https://docs.openclaw.ai/tools/btw.md)

---

## What It Does

```bash
/btw what file are we editing?
/side what changed while the main run continued?
/btw what does this error mean?
/btw summarize the current task in one sentence
```

When you send `/btw <question>`, OpenClaw:
1. Snapshots the current session context
2. Runs a **separate, tool-less, one-shot** model call
3. Answers only the side question
4. Leaves the main run alone
5. Does **not** write the BTW question or answer to session history
6. Emits the answer as a **live side result** (not a normal assistant message)

---

## What It Does NOT Do

- Does not create a new durable session
- Does not continue or resume the main task
- Does not run tools or agent tool loops
- Does not write to transcript history
- Does not appear in `chat.history`
- Does **not** survive a reload — it is intentionally **ephemeral**

---

## Context Behavior

BTW uses the current session as **background context only**. If a main run is active, OpenClaw snapshots the current message state and includes the in-flight main prompt as context, while telling the model:
- Answer only the side question
- Do not resume or complete the unfinished main task
- Do not emit tool calls

---

## Delivery

BTW uses a separate protocol event (`chat.side_result`) instead of the normal `chat` event. This is why it doesn't appear in history.

- **TUI:** Rendered inline, visually distinct from normal replies, dismissible with `Enter`/`Esc`, not replayed on reload
- **Telegram/Discord/WhatsApp:** Delivered as a clearly labeled one-off reply
- **Web/Control UI:** Gateway emits correctly but web client UI for `side_result` is pending

---

## When to Use BTW

Use `/btw` when you want:
- A quick clarification about the current work
- A factual side answer while a long run is still in progress
- A temporary answer that should NOT become part of future session context

## When NOT to Use BTW

If you want the answer to become part of Legacy's future working context — ask normally in the main session instead.

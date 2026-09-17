# OpenClaw Concurrency — How Multiple Requests Are Handled

> Source: OpenClaw 2026.3.7 gateway internals. Last updated: 2026-03-08.

## TL;DR

OpenClaw serializes runs **per session** but allows **parallel sessions** up to a configurable cap. If 5 messages arrive at the same time, the behavior depends on whether they hit the same session or different sessions.

---

## Core Concepts

### Lanes

OpenClaw uses a **lane-aware FIFO queue** to manage concurrency:

| Lane | Default Concurrency | Purpose |
|---|---|---|
| `session:<key>` | **1** (always) | Guarantees one agent run per session at a time |
| `main` | **4** | Global cap for inbound + heartbeat runs |
| `subagent` | **8** | Subagent runs (don't block inbound) |
| `cron` | **3** | Cron/scheduled jobs (independent lane) |

- Every run is queued into **two lanes**: its session lane + the global lane.
- The session lane prevents transcript/tool races within a single conversation.
- The global lane prevents overwhelming the LLM provider with too many parallel requests.

### Queue Modes (Within a Session)

When the agent is busy on a session and more messages arrive for the **same session**:

| Mode | Behavior |
|---|---|
| **collect** (default) | Coalesce all queued messages into a single followup turn after the current run finishes |
| **steer** | Inject message into the current run mid-stream (cancels pending tool calls at next boundary) |
| **followup** | Queue individually, process one at a time after current run ends |
| **steer-backlog** | Steer now AND queue for a followup turn |
| **interrupt** | Abort current run, process newest message |

Legacy uses `collect` — the most efficient mode for multi-message bursts.

---

## Scenario: 5 Messages Arrive Simultaneously

### Case A: Same Session (5 WhatsApp messages from Zeus)

```
Messages: "hey" "also check my email" "and the calendar" "dont forget todoist" "thanks"
                    ↓ all same session key
          ┌─────────────────────────┐
          │  Session: agent:legacy:main  │
          │  Mode: collect          │
          │  Debounce: 1000ms       │
          └─────────────────────────┘
                    ↓
          Coalesced into ONE agent turn:
          "Zeus said: hey / also check my email / and the calendar /
           dont forget todoist / thanks"
                    ↓
          Legacy processes all 5 in a single response
```

- Only **1 agent turn** runs (most efficient)
- Debounce waits 1 second for burst to settle
- If queue cap (20) is exceeded, overflow is summarized as a bullet list

### Case B: Different Sessions (5 different sources)

```
├── WhatsApp (Zeus)        → session: agent:legacy:main
├── Telegram (Zeus)        → session: agent:legacy:main  (same DM scope!)
├── Discord (#general)     → session: agent:legacy:discord:channel:123
├── Voice Call (inbound)   → session: agent:legacy:voice:call:456
└── Webhook (NATS event)   → session: hook:ingress
```

Note: WhatsApp + Telegram **share the same session** when `dmScope: "main"` (default) — they'd be collected together. The other 3 are independent.

```
Global lane: main (cap=4)
├── [RUNNING] agent:legacy:main (WhatsApp + Telegram collected)
├── [RUNNING] agent:legacy:discord:channel:123
├── [RUNNING] agent:legacy:voice:call:456
├── [RUNNING] hook:ingress
└── [QUEUED]  (nothing — all 4 fit within the cap)
```

All 4 unique sessions run **in parallel**. Each gets its own agent turn with its own context.

### Case C: True 5-Way Parallel (all different sessions)

If all 5 hit separate sessions:

```
Global lane: main (cap=4)
├── [RUNNING] Session 1
├── [RUNNING] Session 2
├── [RUNNING] Session 3
├── [RUNNING] Session 4
└── [QUEUED]  Session 5 → waits for any of the 4 to finish
```

The 5th request is queued in the global lane and runs as soon as one slot opens.

---

## Typing Indicators

Typing indicators fire **immediately on enqueue** (before the run starts). So even if a message is queued behind other work, the user sees "Legacy is typing..." right away.

---

## Configuration

```json5
{
  "agents": {
    "defaults": {
      "maxConcurrent": 4,  // global lane cap (increase for high-volume)
      "subagents": {
        "maxConcurrent": 8  // subagent lane cap
      }
    }
  },
  "messages": {
    "queue": {
      "mode": "collect",      // default queue mode
      "debounceMs": 1000,     // wait for burst to settle
      "cap": 20,              // max queued messages per session
      "drop": "summarize",    // overflow policy
      "byChannel": {
        "whatsapp": "collect",
        "telegram": "collect"
      }
    }
  },
  "cron": {
    "maxConcurrentRuns": 3  // cron lane cap
  }
}
```

---

## Key Guarantees

1. **One run per session at a time** — no transcript races, no tool conflicts
2. **Global cap prevents provider rate limits** — defaults to 4 parallel LLM calls
3. **Subagents don't block inbound** — separate lane with higher concurrency
4. **Cron doesn't block chat** — separate lane entirely
5. **Queue modes are per-channel** — WhatsApp can `collect` while Discord uses `steer`
6. **Stop command clears the queue** — `/stop` aborts current run AND clears queued followups for that session


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]

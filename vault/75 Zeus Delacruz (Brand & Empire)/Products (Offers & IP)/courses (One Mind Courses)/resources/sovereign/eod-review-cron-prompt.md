# EOD Review — Cron Prompt

> *Paste this into your OpenClaw cron job configuration for the end-of-day review.*
> *Runs daily at 8:00 PM. Output goes to your Discord #legacy channel.*

---

## Cron Configuration

In your `openclaw.json` cron section:

```json
{
  "cron": [
    {
      "id": "eod-review",
      "schedule": "0 20 * * *",
      "agent": "main",
      "prompt": "[PROMPT BELOW]",
      "channel": "discord",
      "account": "sentinel",
      "enabled": true
    }
  ]
}
```

---

## The Prompt

```
You are Legacy ⚡, Zeus's AI co-pilot. This is the end-of-day review for [CURRENT_DATE].

Generate a concise EOD review using this structure:

**🌙 EOD REVIEW — [DATE]**

**✅ WINS**
What got done today? Focus on meaningful progress — not just task completion. What actually moved the needle?

**🔄 CARRY FORWARD**
What didn't get done that was planned? For each item: should it be rescheduled, delegated, or dropped? Be direct.

**📊 DAY QUALITY**
Rate today 1–10 based on: focus quality, MITs completed, pillar balance. One sentence of honest assessment.

**🔔 TOMORROW'S SETUP**
The single most important thing to tackle first tomorrow. What does Zeus need to know or prepare to execute it?

**💡 PATTERN ALERT**
If you notice a recurring issue from the last 3-7 days (based on your memory of daily reviews), flag it. Otherwise: "No pattern to flag."

Keep it sharp. No preamble. Total length: under 200 words.
```

---

## Customization Notes

**To add task completion tracking:** Prepend with `Check Todoist for tasks completed today using the todoist-api-rest skill.`

**On tone:** This is the end of day — Zeus should be able to read this in 90 seconds and close the day with full context.

**Memory note:** OpenClaw maintains conversation memory. Legacy can reference previous days' EOD reviews when generating the pattern alert — this becomes more valuable over time.

---

*Part of the One Mind OS OpenClaw Deep Dive — Course 3*

# Weekly Report — Cron Prompt

> *Paste this into your OpenClaw cron job configuration for the weekly intelligence report.*
> *Runs every Sunday at 5:00 PM — before you do your weekly review.*

---

## Cron Configuration

In your `openclaw.json` cron section:

```json
{
  "cron": [
    {
      "id": "weekly-report",
      "schedule": "0 17 * * 0",
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
You are Legacy ⚡, Zeus's AI co-pilot. Today is Sunday [CURRENT_DATE]. This is the weekly intelligence report for Week [WEEK_NUMBER].

Generate a comprehensive weekly report. This is the source document Zeus reads before his Sunday weekly review. Make it a complete picture of the week.

**📊 WEEKLY REPORT — Week [WEEK_NUMBER] — [DATE_RANGE]**

**🏆 TOP WINS (3-5)**
The most meaningful things that happened this week. Be specific. Prioritize: completions, breakthroughs, investments in systems.

**📋 COMPLETED TASKS**
Summary by pillar: how many tasks completed in SO / HP / LE / GE. Highlight any notable completions.

**🎯 GOAL PROGRESS**
For each active goal Zeus has: what moved this week? Green/Yellow/Red status. If Yellow or Red: why, and what does Zeus need to decide?

**⚠️ OPEN LOOPS**
Anything that was on the radar this week that didn't close. Not just tasks — conversations, decisions, captures that need routing.

**📈 PATTERNS THIS WEEK**
What themes showed up repeatedly? What was the quality of focus? Did pillars get appropriate attention?

**🔭 NEXT WEEK SETUP**
The 3 most important things next week based on current goal state and open loops. What needs Zeus's direct attention vs. what agents can handle?

**🤖 AGENT ACTIVITY**
Brief note on what specialist agents did this week (if any notable activity). Any agent that needs attention or reconfiguration?

Length: thorough. This replaces needing to scroll through a week of messages. Zeus is about to do his weekly review — give him everything he needs.
```

---

## Customization Notes

**This is the premium cron.** The morning brief is the daily driver. The EOD is the close. The weekly report is the strategic layer.

**When it fires:** Sunday at 5pm — 2 hours before the recommended weekly review time (7pm). Zeus reads the report, then does the review.

**Memory compaction:** If Legacy's memory gets compacted mid-week, some detail may be lost. The weekly report will reflect what's in working memory + the vault. This is a known limitation — resolved by agents writing key events to the vault throughout the week.

---

*Part of the One Mind OS OpenClaw Deep Dive — Course 3*

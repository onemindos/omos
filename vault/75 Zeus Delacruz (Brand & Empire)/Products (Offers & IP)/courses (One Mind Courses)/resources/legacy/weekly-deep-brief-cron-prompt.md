# Weekly Deep Intelligence Brief — Cron Prompt

> *Trinity's weekly synthesis. Runs Saturday 8:00 AM.*
> *Deeper than the daily brief — synthesis, patterns, strategic intelligence.*

---

## Cron Configuration

```json
{
  "id": "weekly-deep-brief",
  "schedule": "0 8 * * 6",
  "agent": "trinity",
  "channel": "discord",
  "account": "sentinel",
  "prompt": "[PROMPT BELOW]",
  "enabled": true
}
```

Runs Saturday — before Zeus's Sunday weekly review prep.

---

## The Prompt

```
You are Trinity 🔺, Zeus's intelligence and research agent. Generate the weekly deep intelligence brief for the week ending [DATE].

**🔺 DEEP INTEL — Week [WEEK_NUMBER]**

**🌊 SIGNAL PATTERNS**
Review the signals from this week's daily briefs. What patterns emerged? What themes showed up more than once? What's building momentum in the tech/business landscape that affects Zeus's domains?

Format as 3-5 pattern observations: [Pattern] → [Evidence] → [Implication for Zeus]

**🔬 RESEARCH PROGRESS**
Active sprint status: [Current sprint name] — [Progress summary]. Is it on track? Key finding so far?

Top 3 items in R&D backlog by priority: list with one-sentence rationale for each.

**🧠 KNOWLEDGE SYNTHESIS**
The most valuable knowledge generated this week across all intelligence gathering. 1-3 items worth permanent vault notes.

Format: [Topic] → [Key insight] → [Recommended note location in vault]

**🎯 STRATEGIC INTELLIGENCE**
One strategic insight this week that is actionable for Zeus's current goals. Not a task — a shift in understanding that should inform his next decisions.

**🔭 NEXT WEEK RADAR**
What to watch next week: 2-3 specific things Trinity will monitor. This sets the intelligence agenda.

**📚 RECOMMENDED READING**
1-2 specific pieces of long-form content (articles, papers, posts) that Trinity recommends Zeus read this week. Include why.

Total length: 400-600 words. This is the synthesis layer — go deeper than daily signals.
```

---

## How This Pairs with Other Automations

```
Monday 9am     → Empire Weekly (business focus)
Daily 6:30am   → Intel Brief (pulse)
Saturday 8am   → Deep Intel Brief (synthesis) ← This file
Sunday 5pm     → Weekly Report (full life view)
Sunday 7pm     → Zeus does weekly review
```

The deep brief feeds the weekly review prep. Zeus reads both (Deep Intel + Weekly Report) before sitting down to review Sunday evening.

---

## Making It More Powerful Over Time

The deep brief improves as Trinity builds memory of:
- What signals matter to Zeus (from feedback over time)
- What Zeus's active research interests are (from sprint log)
- What's already in the vault (from memory indexing)

After 4–8 weeks of running this, ask Trinity: *"What intelligence patterns have you noticed over the past month?"* — the answer becomes a meta-brief that sharpens future intelligence gathering.

---

*Part of the One Mind OS Intelligence & Research — Course 8*

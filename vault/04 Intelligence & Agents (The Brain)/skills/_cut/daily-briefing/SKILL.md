---
name: daily-briefing
description: "Generates personalized morning, evening, and weekly briefings by pulling from calendar, task list, weather, and user context. Use when Zeus says 'what's on today', 'daily briefing', 'what's on my plate', 'weekly review', or at morning/evening startup."
license: Apache-2.0
metadata:
  version: "1.0.0"
  author: zeus
  tags: ["productivity", "briefing", "daily", "summary"]
---

# Daily Briefing Skill

You create personalized daily briefings that help users start their day focused and informed.

## When to Use

- Morning startup routine
- User asks "What's on my plate today?"
- Weekly/monthly review requests
- End of day wrap-up

## Briefing Types

### Morning Briefing
- Today's top priorities (max 3)
- Calendar overview
- Important deadlines
- Quick wins available
- Mindset/motivation note

### Evening Review
- What was accomplished
- What's carrying over
- Tomorrow's preview
- Wins to celebrate

### Weekly Review
- Week's accomplishments
- Goals progress
- Lessons learned
- Next week's focus

## Information Sources

Pull from available context:
- Calendar/schedule
- Task list
- Recent conversations
- User preferences/memories
- Weather (if relevant)

## Briefing Structure

```
## Good [Morning/Afternoon], Zeus!

**[Day], [Date]** | [Weather if available]

### Top Priority
[The ONE thing that matters most today]

### Today's Focus (3 items max)
1. [Priority task with context]
2. [Second priority]
3. [Third priority]

### Calendar
- [Time] - [Event]
- [Time] - [Event]

### Quick Wins (< 15 min each)
- [ ] [Quick task]
- [ ] [Quick task]

### Heads Up
[Anything important to be aware of]

---
*"[Relevant motivational quote or insight]"*
```

## Tone Guidelines

- Energizing but not overwhelming
- Focused on action, not just information
- Personalized to user's style
- Honest about heavy days

## Customization

Adapt based on:
- User's energy level mentioned
- Day of week (Mon = fresh start, Fri = wrap up)
- Known preferences from memory
- Recent context from conversations


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]

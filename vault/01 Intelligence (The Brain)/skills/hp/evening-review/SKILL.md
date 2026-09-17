---
name: evening-review
description: "Runs a structured end-of-day reflection capturing wins, lessons, energy, gratitude, and tomorrow's intention, then saves the entry to the daily note. Use when Zeus says 'evening review', 'end of day', 'daily review', 'reflect', 'wrapping up', or at scheduled evening trigger. Feeds into weekly scorecard and life scores."
allowed-tools:
  - Bash(curl:*)
  - Bash(jq:*)
  - Bash(python3:*)
  - Read
  - Write
requires: [read_file, list_directory, write_file]
offline: true
agent: any
---

# Evening Review — Daily Reflection Skill

> "The unexamined day is not worth repeating."

Runs a structured end-of-day reflection that captures the day's signal and
stores it for compounding intelligence via the × (Multiply) loop.

---

## When to Activate

| Signal | Examples |
|--------|----------|
| Direct request | "evening review", "daily review", "let's reflect" |
| Time-based | "end of day", "wrapping up" |
| Scheduled trigger | N8N fires at configured wind-down time |

---

## Review Flow

### Step 1: Day Scan (automated)

Before asking questions, pull context:

```bash
# What was completed today
curl -s "https://api.todoist.com/rest/v2/tasks?filter=completed today" \
  -H "Authorization: Bearer $TODOIST_API_TOKEN"

# What was due today but not done
curl -s "https://api.todoist.com/rest/v2/tasks?filter=today & !completed" \
  -H "Authorization: Bearer $TODOIST_API_TOKEN"
```

Present: "Here's what I see from today: {completed count} tasks done, {incomplete count} still open. Let me ask you a few questions."

### Step 2: Reflection Questions (conversational)

Ask 4-6 questions, adapting based on responses:

1. **Wins**: "What went well today? What are you proud of?"
2. **Lessons**: "What didn't go as planned? What would you do differently?"
3. **Energy**: "How was your energy today? (1-10) Any patterns you noticed?"
4. **Gratitude**: "What are you grateful for today?"
5. **Tomorrow**: "What's the ONE thing that would make tomorrow great?"
6. **Signal** (optional): "Anything you want to capture or remember from today?"

### Step 3: Synthesis

Summarize the review and propose:
- Tomorrow's top priority
- Any tasks to reschedule
- Any captures to process

---

## Output Format

Append to daily note: `00-24 SO (Sovereign Operations)/02 Command (Dashboard & Navigation)/daily/{YYYY-MM-DD}.md`

```markdown
## 🌙 Evening Review

### ✅ Completed ({count})
- {task 1}
- {task 2}

### ⏳ Carried Forward ({count})
- {incomplete task}

### 🏆 Wins
- {win from conversation}

### 📝 Lessons
- {lesson from conversation}

### ⚡ Energy: {score}/10
{pattern note if any}

### 🙏 Gratitude
- {gratitude item}

### 🎯 Tomorrow's #1
{the one thing}

### 💡 Signals Captured
- {any new captures flagged}
```

---

## Integration Points

- **Weekly Scorecard**: Evening reviews feed into weekly aggregation
- **Life Score**: Energy ratings feed into PI (Performance Intelligence)
- **Morning Briefing**: Tomorrow's #1 becomes next morning's top priority
- **Memory**: Key lessons stored in `_codex/memory/learning/` for pattern detection

---

## Agent Behavior

- Start by showing what was accomplished before asking questions (positive framing)
- Keep it to 5 minutes — this isn't therapy, it's a quick debrief
- If energy is consistently low (3 days in a row), flag it proactively
- Don't ask all 6 questions every time — adapt to mood and energy
- If Zeus says "quick review", do steps 1 + auto-synthesis only (no questions)
- Always end on a forward note (tomorrow's intention)


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]

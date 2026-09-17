---
name: debugging
description: "Isolates and resolves software bugs using a reproduce-isolate-hypothesize-test-fix methodology. Use when Zeus says 'something is broken', 'this isn't working', 'I'm getting an error', or needs root-cause analysis on unexpected behavior."
license: Apache-2.0
metadata:
  version: "1.0.0"
  author: zeus
  tags: ["debugging", "troubleshooting", "problem-solving"]
---

# Debugging Skill

You are a systematic debugger with expertise in isolating and resolving software issues.

## When to Use

- User reports something isn't working
- Error messages need interpretation
- Code is producing unexpected results
- Performance issues need investigation

## Debugging Process

### 1. Reproduce
- Can you consistently reproduce the issue?
- What are the exact steps?
- What environment? (OS, versions, etc.)

### 2. Isolate
- When did it last work?
- What changed recently?
- Minimum reproduction case?

### 3. Hypothesize
- What could cause this behavior?
- List 3-5 possible causes
- Rank by likelihood

### 4. Test
- Test one hypothesis at a time
- Use binary search to narrow down
- Add logging/print statements strategically

### 5. Fix & Verify
- Make the minimal fix
- Verify the fix works
- Check for regressions
- Document what was wrong

## Common Bug Categories

| Category | Symptoms | Check |
|----------|----------|-------|
| **Syntax** | Won't compile/run | Linter, error line |
| **Logic** | Wrong output | Trace execution |
| **State** | Intermittent | Race conditions, init |
| **Data** | Corrupt results | Input validation |
| **Environment** | Works elsewhere | Versions, config |

## Debugging Questions

1. What did you expect to happen?
2. What actually happened?
3. When did it start?
4. Does it happen every time?
5. Any recent changes?

## Output Format

```
## Bug Analysis

**Symptom:** [What's happening]
**Expected:** [What should happen]
**Environment:** [Relevant details]

### Hypotheses
1. [Most likely cause] - [How to test]
2. [Second possibility] - [How to test]
3. [Third possibility] - [How to test]

### Investigation Steps
1. [ ] First thing to check
2. [ ] Second thing to check

### Resolution
[Once found: What was wrong and how to fix it]
```


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]

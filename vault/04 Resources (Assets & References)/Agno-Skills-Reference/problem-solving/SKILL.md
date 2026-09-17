---
name: problem-solving
description: Structured approach to solving complex problems
license: Apache-2.0
metadata:
  version: "1.0.0"
  author: zeus
  tags: ["problem-solving", "analysis", "decision-making"]
---

# Problem Solving Skill

You help users work through complex problems systematically.

## When to Use

- User is stuck on a decision
- Complex problem needs untangling
- Multiple options to evaluate
- Root cause analysis needed

## Problem-Solving Framework

### 1. Define the Problem
- What exactly is the problem?
- Who is affected?
- What's the impact?
- Is this the real problem or a symptom?

### 2. Gather Information
- What do we know?
- What don't we know?
- What assumptions are we making?
- Who else has context?

### 3. Generate Options
- Brainstorm without judging
- Consider obvious AND creative solutions
- What would [expert] do?
- What if we had unlimited resources?

### 4. Evaluate Options
- Pros and cons of each
- Feasibility (can we do it?)
- Impact (will it solve the problem?)
- Risk (what could go wrong?)

### 5. Decide & Act
- Choose the best option
- Define next steps
- Identify who does what
- Set review point

### 6. Review & Learn
- Did it work?
- What would we do differently?
- What did we learn?

## Decision Matrix

For comparing options:

| Criteria | Weight | Option A | Option B | Option C |
|----------|--------|----------|----------|----------|
| Cost | 3 | 4 (12) | 3 (9) | 5 (15) |
| Time | 2 | 3 (6) | 5 (10) | 2 (4) |
| Quality | 3 | 5 (15) | 4 (12) | 3 (9) |
| Risk | 2 | 4 (8) | 3 (6) | 4 (8) |
| **Total** | | **41** | **37** | **36** |

## Five Whys Technique

Keep asking "Why?" to find root cause:

1. Problem: Server is slow
2. Why? → Database queries are slow
3. Why? → Missing indexes
4. Why? → Index was dropped during migration
5. Why? → Migration script had a bug
→ **Root cause: Migration script bug**

## Output Format

```
## Problem Analysis

### Problem Statement
[Clear, specific description]

### Impact
[Who/what is affected and how]

### Root Cause
[Underlying cause, not just symptoms]

---

### Options Considered

**Option A: [Name]**
- Pros: [List]
- Cons: [List]
- Feasibility: [High/Medium/Low]
- Impact: [High/Medium/Low]

**Option B: [Name]**
- Pros: [List]
- Cons: [List]
- Feasibility: [High/Medium/Low]
- Impact: [High/Medium/Low]

---

### Recommendation
[Which option and why]

### Next Steps
1. [Immediate action]
2. [Follow-up action]
3. [Review point]
```

## When to Escalate

- Problem is beyond your expertise
- High stakes decision
- Need buy-in from others
- Blocked by external factors

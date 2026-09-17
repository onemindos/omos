# Hook: RouteToDomain

**Event**: `organize`
**Blocking**: No
**Fail Behavior**: Skip (send to Inbox)

## Instructions

When content needs a home, determine the correct domain:

1. **Check content subject**:
   - AI, tech, infrastructure, security → **SO (00-24)**
   - Personal growth, health, fitness, habits → **HP (25-49)**
   - Family, home, parenting, heritage → **LE (50-74)**
   - Business, revenue, products, marketing → **GE (75-99)**

2. **Narrow to specific domain** by scanning domain folder names.

3. **If uncertain**: Route to `00-24 SO/24 Inbox (Queue)/` — it gets sorted during weekly review.

4. **Never create new top-level folders** — only use existing domain structure.

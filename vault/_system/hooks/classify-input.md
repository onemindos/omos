# Hook: ClassifyInput

**Event**: `capture`
**Blocking**: No
**Fail Behavior**: Skip (agent proceeds with manual classification)

## Instructions

When new information enters the system, classify it:

1. **Determine type**:
   - Idea or observation → `note`
   - A choice made → `decision`
   - A process to follow → `sop`
   - Information about a person → `person`
   - A project or initiative → `project`
   - A time-stamped event → `log`

2. **Suggest template** from `_codex/templates/`:
   - `daily.md` for logs
   - `decision.md` for decisions
   - `person.md` for people
   - `project.md` for projects
   - `review.md` for reviews
   - `sop.md` for SOPs

3. **Tag suggestion**: Recommend `type/{type}` and relevant domain tags.

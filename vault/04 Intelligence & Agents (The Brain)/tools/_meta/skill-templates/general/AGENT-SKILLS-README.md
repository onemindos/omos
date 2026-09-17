# OneMind Skills

Skills are specialized instruction sets that enhance agent capabilities for specific tasks. They use Agno's native `LocalSkills` loader.

## Directory Structure

```
skills/
├── README.md                 # This file
├── code-review/
│   ├── SKILL.md             # Skill definition (required)
│   ├── scripts/             # Optional helper scripts
│   │   └── analyze_complexity.py
│   └── references/          # Optional reference docs
│       └── style-guide.md
├── research/
│   ├── SKILL.md
│   ├── scripts/
│   └── references/
├── task-management/
│   └── SKILL.md
├── debugging/
│   └── SKILL.md
├── writing/
│   └── SKILL.md
├── daily-briefing/
│   ├── SKILL.md
│   └── scripts/
├── project-planning/
│   └── SKILL.md
└── problem-solving/
    └── SKILL.md
```

## Creating a New Skill

### 1. Create Skill Directory

```bash
mkdir backend/agno/skills/my-skill
```

### 2. Create SKILL.md

The `SKILL.md` file is required. It uses YAML frontmatter + markdown content:

```markdown
---
name: my-skill
description: Brief description of what this skill does
license: Apache-2.0
metadata:
  version: "1.0.0"
  author: your-name
  tags: ["tag1", "tag2", "tag3"]
---

# My Skill

Detailed instructions for the agent when using this skill.

## When to Use
- Scenario 1
- Scenario 2

## How to Use
1. Step 1
2. Step 2

## Output Format
```example output```

## Tips
- Tip 1
- Tip 2
```

### 3. Add Scripts (Optional)

Place helper scripts in the `scripts/` subdirectory:

```bash
mkdir backend/agno/skills/my-skill/scripts
```

Example script (`scripts/helper.py`):
```python
#!/usr/bin/env python3
"""
Helper Script

Brief description of what this script does.

Usage:
    python helper.py [args]
"""

def main():
    # Script logic
    pass

if __name__ == "__main__":
    main()
```

### 4. Add References (Optional)

Place reference documents in the `references/` subdirectory:

```bash
mkdir backend/agno/skills/my-skill/references
```

Reference files are markdown documents with templates, guides, or additional context.

## Loading Skills in Agents

Skills are loaded using Agno's native `LocalSkills`:

```python
from agno.agent import Agent
from agno.skills import Skills
from agno.skills.local import LocalSkills
from pathlib import Path

SKILLS_DIR = Path(__file__).parent / "skills"

# Create agent with skills
agent = Agent(
    name="My Agent",
    skills=Skills(
        loaders=[
            LocalSkills(path=SKILLS_DIR)
        ]
    ),
    # Optional: Give agent tools to access skill content
    tools=[
        get_skill_instructions,
        get_skill_reference,
        get_skill_script,
    ]
)
```

### Agent Tools for Skills

These tools allow agents to access skill content dynamically:

```python
from agno.tools import tool

@tool
def get_skill_instructions(skill_name: str) -> str:
    """Get the full instructions from a skill's SKILL.md file."""
    skill_path = SKILLS_DIR / skill_name / "SKILL.md"
    if skill_path.exists():
        return skill_path.read_text()
    return f"Skill '{skill_name}' not found"

@tool
def get_skill_reference(skill_name: str, filename: str) -> str:
    """Get a reference document from a skill."""
    ref_path = SKILLS_DIR / skill_name / "references" / filename
    if ref_path.exists():
        return ref_path.read_text()
    return f"Reference '{filename}' not found in skill '{skill_name}'"

@tool
def get_skill_script(skill_name: str, filename: str) -> str:
    """Get a script from a skill."""
    script_path = SKILLS_DIR / skill_name / "scripts" / filename
    if script_path.exists():
        return script_path.read_text()
    return f"Script '{filename}' not found in skill '{skill_name}'"
```

## API Endpoints

The skills are exposed via REST API for frontend integration:

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/skills` | GET | List all available skills |
| `/api/skills/{id}` | GET | Get skill details with full content |
| `/api/skills/{id}/reference/{filename}` | GET | Get a specific reference file |
| `/api/skills/{id}/script/{filename}` | GET | Get a specific script file |

### Response Format

```json
{
  "skills": [
    {
      "id": "code-review",
      "name": "code-review",
      "description": "Comprehensive code review and quality analysis",
      "metadata": {
        "version": "1.0.0",
        "author": "zeus",
        "tags": ["code", "review", "quality"]
      },
      "has_scripts": true,
      "has_references": true
    }
  ],
  "total": 8
}
```

## Frontend Integration

Skills can be selected when creating/editing agents in the Flutter app:

1. **Skills Tab** - The agent form includes a Skills tab for selection
2. **Skills Provider** - `frontend/lib/agno/skills/providers/skills_provider.dart`
3. **API Client** - Fetches skills from `/api/skills`

### Provider Usage

```dart
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:onemind/agno/skills/providers/skills_provider.dart';

// Watch skills state
final skillsState = ref.watch(skillsProvider);

// Access skills list
final skills = skillsState.skills;

// Filter by search
final filtered = skillsState.filteredSkills;

// Get skills count
final count = ref.watch(skillsCountProvider);
```

## Available Skills

| Skill | Description | Scripts | References |
|-------|-------------|---------|------------|
| `code-review` | Code review and quality analysis | Yes | Yes |
| `research` | Research and information synthesis | Yes | Yes |
| `task-management` | Task organization with GTD | Yes | Yes |
| `debugging` | Systematic debugging approach | No | Yes |
| `writing` | Professional writing assistance | No | Yes |
| `daily-briefing` | Daily briefing generation | Yes | No |
| `project-planning` | Project breakdown and planning | No | No |
| `problem-solving` | Structured problem solving | No | No |

## Best Practices

1. **Keep SKILL.md focused** - Each skill should have a single, clear purpose
2. **Use tags wisely** - Tags help with discovery and categorization
3. **Include examples** - Show sample outputs in the skill instructions
4. **Document scripts** - Add docstrings and usage info to helper scripts
5. **Version your skills** - Use semantic versioning in metadata
6. **Test with agents** - Verify skills work as expected with real agent runs

## Agno LocalSkills Documentation

For more details on the Agno LocalSkills system, see:
- [Agno Skills Documentation](https://docs.agno.com/skills)
- [LocalSkills API Reference](https://docs.agno.com/skills/local)


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]

# Skills System (OMOS-242)

The Skills System enables agents to have attachable, reusable knowledge chunks with **progressive disclosure** - agents see minimal skill descriptions in their system prompt and load full instructions on-demand.

## Overview

Skills define **what agents KNOW** - domain expertise, procedures, guidelines, and specialized knowledge. Unlike tools (which are functions agents can call), skills provide contextual information and instructions.

### Key Concepts

| Concept | Description |
|---------|-------------|
| **Skill** | A reusable knowledge chunk defined in a `SKILL.md` file |
| **Skill ID** | Format: `level/name` (e.g., `shared/three-paths`, `domain/homestead`) |
| **Progressive Disclosure** | Agents see only skill summaries initially, load full content on-demand |
| **Skill Level** | Hierarchy: shared → specialist → team-leader → domain |

## Directory Structure

```
backend/agno/skills/
├── __init__.py           # Exports all skill functions
├── loader.py             # Skills loading and progressive disclosure
├── tracking.py           # Usage tracking and instruction tools
├── api.py                # REST API endpoints
├── shared/               # Skills available to ALL agents
│   └── three-paths/
│       └── SKILL.md
├── specialist/           # Role-specific skills
│   └── code-review/
│       └── SKILL.md
├── team-leader/          # Management skills
│   └── farm-operations/
│       └── SKILL.md
└── domain/               # Domain expertise
    └── homestead/
        └── SKILL.md
```

## SKILL.md Format

Each skill is defined in a `SKILL.md` file with YAML frontmatter:

```markdown
---
name: three-paths
description: Framework for organizing life across HP (self), LE (family), and GE (business)
version: 1.0.0
author: zeus
tags:
  - productivity
  - life-management
  - framework
domain: wisdom
---

# Three Paths Framework

## Overview
The Three Paths organize all aspects of life...

## HP - Holistic Performance
- Health, fitness, habits
- Personal development
- Self-optimization

## LE - Legacy Evolution
- Family, home, heritage
- Relationships

## GE - Generational Endeavor
- Business, wealth creation
- Financial independence
```

### Frontmatter Fields

| Field | Required | Description |
|-------|----------|-------------|
| `name` | Yes | Skill display name |
| `description` | Yes | Short description (shown in progressive disclosure) |
| `version` | No | Semantic version (default: 1.0.0) |
| `author` | No | Skill author |
| `tags` | No | List of tags for filtering |
| `pillar` | No | Associated life pillar |
| `allowed-tools` | No | Tools the skill may use |

## Skill Levels

### 1. Shared (`shared/`)
Available to **all agents**. Common knowledge, frameworks, guidelines.

```
shared/
├── three-paths/       # Life organization framework
├── communication/      # Communication guidelines
└── safety-protocols/   # Safety procedures
```

### 2. Specialist (`specialist/`)
For **specialist agents** with specific roles.

```
specialist/
├── code-review/        # Code review procedures
├── research-synthesis/ # Research compilation methods
└── data-analysis/      # Data analysis techniques
```

### 3. Team Leader (`team-leader/`)
For **managers and coordinators**.

```
team-leader/
├── farm-operations/    # Farm management procedures
├── house-management/   # Home operations
├── security-ops/       # Security protocols
└── executive-strategy/ # Business strategy
```

### 4. Domain (`domain/`)
**Topic-specific expertise**.

```
domain/
├── homestead/          # Homesteading knowledge
├── smart-home/         # Home automation
├── cybersecurity/      # Security expertise
└── business/           # Business domain
```

## Progressive Disclosure

### How It Works

1. **Agent Creation**: Skills are attached to agent via `skills` field
2. **System Prompt**: Agent sees minimal skill summaries
3. **On-Demand Loading**: Agent calls `get_skill_instructions(skill_id)` when needed
4. **Usage Tracking**: Every instruction load is tracked for learning

### System Prompt Context

When an agent has skills configured, this is added to their context:

```markdown
## Available Skills

You have access to the following skills. Use `get_skill_instructions(skill_id)` to load full details when needed:

- **shared/three-paths**: Framework for organizing work across Health, Wealth, Wisdom, and Relationships
- **domain/homestead**: Homesteading procedures and knowledge

Only load skill instructions when you need them for the current task.
```

### Progressive Disclosure Tools

Agents with skills automatically get these tools:

#### `get_skill_instructions(skill_id)`
Load full instructions for a skill.

```python
# Agent tool call
result = await get_skill_instructions("domain/homestead")
# Returns full SKILL.md content with metadata header
```

#### `get_available_skills()`
List all available skills.

```python
# Agent tool call
result = await get_available_skills()
# Returns organized list of all skills with descriptions
```

## API Endpoints

### List Skills
```http
GET /api/skills
```

Response:
```json
{
  "skills": [
    {
      "id": "shared/three-paths",
      "name": "four-pillars",
      "description": "Framework for organizing work...",
      "level": "shared",
      "category": "FOUR PILLARS",
      "has_scripts": false,
      "uses": 42,
      "success_rate": 95.2,
      "avg_time_ms": 1250
    }
  ],
  "by_level": {"shared": 3, "domain": 2},
  "by_category": {"FOUR PILLARS": 1, "HOMESTEAD": 1}
}
```

### Get Single Skill
```http
GET /api/skills/{skill_id}
```

### Create Agent with Skills
```http
POST /api/agents
Content-Type: application/json

{
  "name": "Farm Manager",
  "model_id": "claude-sonnet",
  "skills": [
    "shared/three-paths",
    "team-leader/farm-operations",
    "domain/homestead"
  ],
  "tools": ["weather", "calendar"],
  "enable_agentic_memory": true
}
```

### Update Agent Skills
```http
PATCH /api/agents/{agent_id}
Content-Type: application/json

{
  "skills": ["shared/three-paths", "domain/smart-home"]
}
```

## Backend Usage

### Loading Skills in Agent Factory

```python
from backend.agno.skills import (
    create_skills_loader_from_ids,
    build_skill_context,
    create_skill_instruction_tools,
)

# In create_agent_from_config():
skill_ids = config.get("skills", [])
if skill_ids:
    # Load Agno Skills loader
    skills = create_skills_loader_from_ids(skill_ids)

    # Build context for system prompt
    skill_context = build_skill_context(skill_ids)

    # Add progressive disclosure tools
    skill_tools = create_skill_instruction_tools()
    tools.extend(skill_tools)
```

### Using Preset Skills

```python
from backend.agno.skills import (
    skills_for_farm_manager,
    skills_for_code_agent,
    skills_for_ceo_agent,
)

# Pre-configured skill sets
agent = Agent(
    name="Farm Manager",
    skills=skills_for_farm_manager(),
)
```

### Tracking Skill Usage

```python
from backend.agno.skills import log_skill_usage, get_skill_usage_summary

# Log usage
await log_skill_usage(
    agent_id="farm-manager",
    skill_id="domain/homestead",
    task_type="crop_planning",
    success=True,
    duration_ms=1500,
)

# Get summary
summary = await get_skill_usage_summary(agent_id="farm-manager")
# Returns: {total_uses, success_rate, avg_duration_ms, by_skill, by_task_type}
```

## Flutter Integration

### Skills Provider

```dart
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../providers/skills_provider.dart';

// Watch all skills
final skillsState = ref.watch(skillsProvider);

// Watch skills by level
final skillsByLevel = ref.watch(skillsByLevelProvider);

// Track selected skills for agent creation
final selectedIds = ref.watch(selectedSkillIdsProvider);
```

### Skill Picker Widget

```dart
import '../../skills/widgets/skill_picker.dart';

// In agent form dialog
SkillPicker(
  selectedSkillIds: _selectedSkills,
  onSelectionChanged: (newSelection) {
    setState(() => _selectedSkills = newSelection);
  },
)
```

### Compact Skill Display

```dart
SkillPickerCompact(
  selectedSkillIds: agent.skills.toSet(),
  onEditPressed: () => _showSkillSelector(),
)
```

### Skill Model

```dart
class Skill {
  final String id;           // e.g., "domain/homestead"
  final String name;
  final String description;
  final SkillLevel level;    // shared, specialist, teamLeader, domain
  final SkillCategory category;
  final SkillMetadata metadata;
  final int uses;
  final double successRate;
  final double avgTimeMs;
  final bool isActive;
}
```

## Creating New Skills

### 1. Create Directory

```bash
mkdir -p backend/agno/skills/domain/my-skill
```

### 2. Create SKILL.md

```markdown
---
name: my-skill
description: Short description of what this skill provides
version: 1.0.0
author: your-name
tags:
  - relevant-tag
---

# My Skill

## Overview
Detailed explanation of the skill...

## Procedures
Step-by-step instructions...

## Best Practices
Guidelines and recommendations...
```

### 3. Verify Discovery

```bash
curl http://localhost:8081/api/skills | jq '.skills[] | select(.id == "domain/my-skill")'
```

### 4. Attach to Agent

```json
{
  "name": "My Agent",
  "skills": ["domain/my-skill"]
}
```

## Memory Integration

Skills integrate with agent memory to track usage patterns:

```python
from backend.agno.skills import build_skill_memory_context

# Build context from skill history for injection
context = build_skill_memory_context(agent_id="farm-manager")
```

This adds to agent context:
```markdown
## Your Skill Usage History

You have used skills 42 times with 95% success rate.

Most used skills:
- domain/homestead: 25 uses
- shared/three-paths: 12 uses
- team-leader/farm-operations: 5 uses
```

## Best Practices

### Skill Design

1. **Keep descriptions concise** - Used in progressive disclosure summaries
2. **Structure content clearly** - Use headers, lists, code blocks
3. **Include examples** - Show practical applications
4. **Version appropriately** - Update version on significant changes

### Agent Configuration

1. **Start minimal** - Only attach skills the agent actually needs
2. **Use appropriate levels** - Match skill level to agent role
3. **Combine strategically** - Mix shared + specialist + domain as needed

### Progressive Disclosure

1. **Let agents decide** - Don't pre-load all instructions
2. **Track usage** - Monitor which skills are loaded most
3. **Optimize summaries** - Good descriptions reduce unnecessary loads

## Troubleshooting

### Skill Not Found

```
Error: Skill 'domain/my-skill' not found
```

**Check:**
- Directory exists: `backend/agno/skills/domain/my-skill/`
- SKILL.md exists in directory
- Skill ID format is `level/name`

### Skills Not Loading

```python
# Debug skill loading
from backend.agno.skills.loader import list_available_skills
print(list_available_skills())
```

### Progressive Disclosure Tools Missing

Ensure agent has skills configured:
```python
config = {
    "skills": ["shared/three-paths"],  # Must have at least one skill
    ...
}
```

## Related Files

| File | Purpose |
|------|---------|
| [backend/agno/skills/__init__.py](backend/agno/skills/__init__.py) | Package exports |
| [backend/agno/skills/loader.py](backend/agno/skills/loader.py) | Skill loading functions |
| [backend/agno/skills/tracking.py](backend/agno/skills/tracking.py) | Usage tracking & tools |
| [backend/agno/skills/api.py](backend/agno/skills/api.py) | REST API endpoints |
| [backend/agno/agents/factory.py](backend/agno/agents/factory.py) | Agent creation with skills |
| [backend/agno/agents/models.py](backend/agno/agents/models.py) | API request/response models |
| [backend/agno/agents/persistence.py](backend/agno/agents/persistence.py) | Database persistence |
| [frontend/lib/agno/skills/providers/skills_provider.dart](frontend/lib/agno/skills/providers/skills_provider.dart) | Flutter state management |
| [frontend/lib/agno/skills/widgets/skill_picker.dart](frontend/lib/agno/skills/widgets/skill_picker.dart) | Skill selection widget |
| [frontend/lib/agno/agents/widgets/agent_form_dialog.dart](frontend/lib/agno/agents/widgets/agent_form_dialog.dart) | Agent builder with Skills tab |

# Circle.so AI Agents

> **AI-Powered Community Management and Engagement**

## Overview

| Attribute   | Value                              |
| ----------- | ---------------------------------- |
| Type        | Community AI Platform              |
| Provider    | Circle.so                          |
| Integration | Native + API                       |
| Features    | AI Moderation, Content, Engagement |
| Status      | ✅ Active                          |

## What is Circle AI?

Circle.so has integrated AI capabilities for community management, content creation, and member engagement. These AI agents help automate moderation, generate content suggestions, and enhance member experiences.

## AI Agent Capabilities

### Content AI

| Feature                 | Description                   |
| ----------------------- | ----------------------------- |
| Post Drafting           | AI-assisted post creation     |
| Title Suggestions       | Generate engaging post titles |
| Content Summarization   | TL;DR for long discussions    |
| Hashtag Recommendations | Suggest relevant tags         |

### Moderation AI

| Feature                | Description                    |
| ---------------------- | ------------------------------ |
| Spam Detection         | Auto-flag suspicious content   |
| Toxicity Filtering     | Identify harmful language      |
| Content Classification | Categorize posts automatically |
| User Risk Scoring      | Identify potential bad actors  |

### Engagement AI

| Feature                 | Description                         |
| ----------------------- | ----------------------------------- |
| Member Matching         | Connect similar members             |
| Content Recommendations | Personalized feed curation          |
| Activity Prompts        | Nudge inactive members              |
| Event Suggestions       | Recommend events based on interests |

## OneMind Community Integration

### AI Moderation Setup

```python
import requests

CIRCLE_API = "https://api.circle.so/v1"
CIRCLE_TOKEN = os.getenv("CIRCLE_API_KEY")

# Configure AI moderation
moderation_config = {
    "community_id": COMMUNITY_ID,
    "settings": {
        "spam_detection": True,
        "toxicity_threshold": 0.7,
        "auto_flag": True,
        "auto_remove": False,  # Manual review
        "notify_admins": True
    }
}
```

### AI Content Assistant

```python
# Generate post suggestions
def get_ai_suggestions(topic: str, space_id: str):
    response = requests.post(
        f"{CIRCLE_API}/ai/content/suggest",
        headers={"Authorization": f"Token {CIRCLE_TOKEN}"},
        json={
            "topic": topic,
            "space_id": space_id,
            "type": "discussion",
            "tone": "professional"
        }
    )
    return response.json()["suggestions"]
```

## AI Agent Workflows

### New Member Onboarding Agent

```yaml
trigger: member_joined
actions:
  - send_welcome_dm:
      message: "personalized based on profile"
  - suggest_spaces:
      based_on: "interests"
  - connect_with_members:
      criteria: "similar_interests"
  - schedule_checkin:
      delay: "7_days"
```

### Content Curation Agent

```yaml
trigger: daily_schedule
actions:
  - analyze_trending_topics
  - generate_discussion_prompts
  - identify_top_contributors
  - create_weekly_digest
  - notify_admins_insights
```

### Engagement Recovery Agent

```yaml
trigger: member_inactive_14_days
actions:
  - analyze_past_interests
  - find_relevant_new_content
  - send_personalized_nudge
  - track_reengagement
```

## Circle + OneMind Synergy

| Circle AI Feature | OneMind Integration           |
| ----------------- | ----------------------------- |
| Member Insights   | Feed into CRM personas        |
| Content Analytics | Inform content strategy       |
| Engagement Data   | Train community health models |
| Moderation Logs   | Security monitoring           |

## Best Practices

1. **Review AI Actions** - Monitor AI moderation decisions
2. **Train on Context** - Provide community guidelines to AI
3. **Gradual Automation** - Start with suggestions, then auto-actions
4. **Feedback Loop** - Correct AI mistakes to improve accuracy
5. **Human Oversight** - Keep humans in the loop for sensitive actions

## Resources

- [Circle AI Features](https://circle.so/features/ai)
- [Circle API Documentation](https://api.circle.so)
- [Community Best Practices](https://circle.so/resources)

---

Last Updated: January 25, 2026


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]

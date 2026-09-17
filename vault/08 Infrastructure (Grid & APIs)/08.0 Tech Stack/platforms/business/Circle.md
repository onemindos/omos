# Circle.so

> **Community Platform**

## Overview

Circle.so is a community platform for hosting the OneMind community, courses, and member discussions.

## Status

- **Active**: ✅ Yes
- **Plan**: Pro
- **Last Updated**: January 25, 2026

## Use Cases

1. **Community** - Member discussions
2. **Courses** - Educational content
3. **Events** - Live sessions
4. **Spaces** - Topic organization
5. **Member Directory** - Networking

## Key Features

- Spaces (forums/channels)
- Events (live & recorded)
- Courses with drip content
- Member profiles & directory
- Gamification (points, badges)
- Native mobile apps
- Custom branding

## Community Structure

### Spaces

| Space                   | Type   | Access          |
| ----------------------- | ------ | --------------- |
| **Welcome**             | Basic  | All             |
| **Announcements**       | Basic  | All (read only) |
| **General Discussion**  | Basic  | All             |
| **Show & Tell**         | Basic  | All             |
| **Feature Requests**    | Basic  | All             |
| **Technical Support**   | Basic  | All             |
| **Pro Members**         | Basic  | Pro tier        |
| **Course: OneMind 101** | Course | Pro tier        |

## API Integration

### Authentication

```python
import requests

CIRCLE_API_KEY = os.getenv("CIRCLE_API_KEY")
CIRCLE_COMMUNITY_ID = os.getenv("CIRCLE_COMMUNITY_ID")

headers = {
    "Authorization": f"Token {CIRCLE_API_KEY}",
    "Content-Type": "application/json"
}

base_url = "https://app.circle.so/api/v1"
```

### Get Members

```python
def get_members(limit: int = 100):
    response = requests.get(
        f"{base_url}/community_members",
        headers=headers,
        params={
            "community_id": CIRCLE_COMMUNITY_ID,
            "per_page": limit
        }
    )
    return response.json()
```

### Create Post

```python
def create_post(space_id: str, name: str, body: str, user_email: str):
    response = requests.post(
        f"{base_url}/posts",
        headers=headers,
        json={
            "community_id": CIRCLE_COMMUNITY_ID,
            "space_id": space_id,
            "name": name,
            "body": body,
            "user_email": user_email,
            "is_pinned": False,
            "skip_notifications": False
        }
    )
    return response.json()
```

### Add Member

```python
def add_member(email: str, name: str, space_group_ids: list = None):
    response = requests.post(
        f"{base_url}/community_members",
        headers=headers,
        json={
            "community_id": CIRCLE_COMMUNITY_ID,
            "email": email,
            "name": name,
            "space_group_ids": space_group_ids or [],
            "skip_invitation": False
        }
    )
    return response.json()
```

### Update Member

```python
def update_member(member_id: str, space_group_ids: list):
    response = requests.put(
        f"{base_url}/community_members/{member_id}",
        headers=headers,
        json={
            "community_id": CIRCLE_COMMUNITY_ID,
            "space_group_ids": space_group_ids
        }
    )
    return response.json()
```

## Webhooks

### Webhook Events

| Event             | Description            |
| ----------------- | ---------------------- |
| `member.created`  | New member joined      |
| `member.deleted`  | Member removed         |
| `post.created`    | New post created       |
| `comment.created` | New comment            |
| `event.rsvp`      | Member RSVPed to event |

### Handle Webhook

```python
from fastapi import FastAPI, Request

@app.post("/webhooks/circle")
async def circle_webhook(request: Request):
    payload = await request.json()
    event_type = payload.get("event")
    data = payload.get("data")

    if event_type == "member.created":
        await handle_new_member(data)
    elif event_type == "post.created":
        await handle_new_post(data)

    return {"status": "ok"}

async def handle_new_member(data):
    # Send welcome message
    # Add to email list
    # Notify in Discord
    pass
```

## Stripe Integration

### On Subscription Created

```python
async def handle_stripe_subscription(customer_email: str, plan: str):
    # Find member in Circle
    members = get_members()
    member = next((m for m in members if m['email'] == customer_email), None)

    if member:
        # Upgrade to pro space group
        space_groups = get_pro_space_groups(plan)
        update_member(member['id'], space_groups)
    else:
        # Invite as pro member
        add_member(customer_email, customer_name, get_pro_space_groups(plan))
```

## Legacy AI Integration

### Community Insights

```python
@Tool
def get_community_stats() -> dict:
    """Get community statistics."""
    members = get_members()
    return {
        "total_members": len(members),
        "active_last_week": count_active(members, days=7),
        "new_this_month": count_new(members, days=30)
    }

@Tool
def create_announcement(title: str, content: str) -> str:
    """Create an announcement post in the community."""
    result = create_post(
        space_id=ANNOUNCEMENTS_SPACE_ID,
        name=title,
        body=content,
        user_email=ADMIN_EMAIL
    )
    return f"Announcement created: {result['url']}"
```

## Member Tiers

| Tier           | Access                    | Price  |
| -------------- | ------------------------- | ------ |
| **Free**       | Welcome, General, Support | $0     |
| **Pro**        | All spaces + courses      | $29/mo |
| **Enterprise** | Pro + priority support    | Custom |

## Space Groups

| Group           | Spaces Included                  |
| --------------- | -------------------------------- |
| **Default**     | Welcome, Announcements, General  |
| **Pro Members** | All Default + Pro space, Courses |
| **Enterprise**  | All Pro + Enterprise channel     |

## Events

### Event Types

- **Live Sessions**: Zoom/Meet integration
- **AMAs**: Q&A sessions
- **Workshops**: Interactive learning
- **Office Hours**: 1:1 support

### API: Get Events

```python
def get_upcoming_events():
    response = requests.get(
        f"{base_url}/events",
        headers=headers,
        params={
            "community_id": CIRCLE_COMMUNITY_ID,
            "status": "upcoming"
        }
    )
    return response.json()
```

## Gamification

### Points System

| Action          | Points |
| --------------- | ------ |
| Create post     | 10     |
| Comment         | 5      |
| Like            | 1      |
| Complete course | 100    |
| Attend event    | 25     |

### Badges

- **Founding Member**: Early adopter
- **Contributor**: 100+ posts
- **Helper**: 50+ helpful comments
- **Course Graduate**: Completed all courses

## Best Practices

- ✅ Create welcoming onboarding flow
- ✅ Post regularly (2-3x per week)
- ✅ Respond to posts within 24h
- ✅ Host live events monthly
- ✅ Celebrate member wins
- ✅ Use gamification to encourage engagement
- ✅ Segment content by member tier

## Resources

- [Circle API Documentation](https://api.circle.so/)
- [Circle Help Center](https://help.circle.so/)

---

**Related**: [Stripe](./Stripe.md) | [Whop](./Whop.md)

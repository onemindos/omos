# Zoho Cliq

> **Team Chat & Collaboration Platform**

## Overview

Zoho Cliq is a team communication and collaboration platform with channels, direct messages, audio/video calls, bots, and integrations.

## Status

- **Active**: ✅ Yes
- **Integration**: API + Webhooks + Bots
- **Last Updated**: January 25, 2026

## Use Cases

1. **Team Communication** - Real-time chat
2. **Project Channels** - Organized discussions
3. **Notifications** - Alerts from other systems
4. **Bot Automation** - Automated workflows
5. **Audio/Video Calls** - Team meetings

## Key Features Used

- [ ] Channels (public/private)
- [ ] Direct messages (1:1 and group)
- [ ] Audio/video calls
- [ ] Screen sharing
- [ ] Bots & slash commands
- [ ] File sharing
- [ ] Message threads
- [ ] Webhooks for notifications

## Channel Structure

```
OneMind Cliq Channels:
├── #general (everyone)
├── #announcements (important updates)
├── #development (dev team)
├── #marketing (marketing team)
├── #support (customer support)
├── #sales (sales team)
├── #random (off-topic)
├── #alerts (system notifications)
└── Private Channels
    ├── #leadership
    └── #projects-alpha
```

## Message Features

### Basic Messaging

- Rich text formatting (bold, italic, lists)
- Code blocks with syntax highlighting
- Mentions (@user, @channel, @everyone)
- Emojis & reactions
- Message editing & deletion
- Message pinning
- Threads (reply to specific messages)

### Advanced Features

- **Slash Commands**: `/remind`, `/poll`, `/status`
- **Message Scheduling**: Send message later
- **Message Formatting**: Markdown support
- **Snippets**: Code/text snippets

## Audio/Video Calls

### Call Features

- 1:1 calls
- Group calls (up to 100 participants)
- Screen sharing
- Call recording
- Background blur
- Virtual backgrounds
- Picture-in-picture mode

### Call Quality

- HD video (720p)
- Noise cancellation
- Echo reduction
- Network adaptation

## Bots & Automation

### Custom Bots

**Legacy AI Bot**

```
Purpose: OneMind assistant bot
Commands:
  /status - System status
  /create [task] - Create Todoist task
  /deploy [env] - Trigger deployment
  /help - Show available commands
```

**Notification Bot**

```
Purpose: System alerts & notifications
Sends:
  - CRM: New leads
  - Desk: High-priority tickets
  - GitHub: PR reviews needed
  - Linear: Issue assignments
  - Monitoring: Server alerts
```

### Built-in Bots

- **Poll Bot**: Create polls/surveys
- **Reminder Bot**: Set reminders
- **ToDo Bot**: Personal task list
- **Calendar Bot**: Meeting scheduling

## Slash Commands

| Command   | Purpose        | Example                                                   |
| --------- | -------------- | --------------------------------------------------------- |
| `/remind` | Set reminder   | `/remind me in 1 hour to check email`                     |
| `/poll`   | Create poll    | `/poll "Pizza toppings?" "Pepperoni" "Veggie" "Hawaiian"` |
| `/status` | Set status     | `/status "In a meeting" 30m`                              |
| `/giphy`  | Send GIF       | `/giphy celebration`                                      |
| `/away`   | Set away       | `/away "Lunch break"`                                     |
| `/dnd`    | Do not disturb | `/dnd 2h`                                                 |

## Webhooks Integration

### Incoming Webhooks

Send messages to Cliq from external systems:

```python
import requests

webhook_url = "https://cliq.zoho.com/api/v2/channelsbyname/alerts/message"

payload = {
    "text": "🚨 **High Priority Alert**",
    "card": {
        "title": "Server CPU at 90%",
        "theme": "modern-inline"
    },
    "bot": {
        "name": "Monitoring Bot",
        "image": "https://..."
    }
}

requests.post(webhook_url, json=payload)
```

### Active Webhook Integrations

1. **GitHub** → PR activity to #development
2. **Linear** → Issue updates to #development
3. **Todoist** → Task reminders to #alerts
4. **Zoho CRM** → Hot leads to #sales
5. **Zoho Desk** → Urgent tickets to #support
6. **Server Monitoring** → Alerts to #alerts

## Message Cards

### Card Types

**Notification Card**

```json
{
  "text": "New Lead: Tesla Inc.",
  "card": {
    "title": "CRM Lead Notification",
    "subtitle": "High-value opportunity",
    "theme": "modern-inline",
    "thumbnail": "https://...",
    "buttons": [
      {
        "label": "View in CRM",
        "action": { "type": "open.url", "url": "https://..." }
      }
    ]
  }
}
```

**Poll Card**

```json
{
  "text": "Quick team poll",
  "card": {
    "title": "Best time for standup?",
    "buttons": [
      {
        "label": "9 AM",
        "action": { "type": "invoke.function", "data": { "vote": "9am" } }
      },
      {
        "label": "10 AM",
        "action": { "type": "invoke.function", "data": { "vote": "10am" } }
      },
      {
        "label": "11 AM",
        "action": { "type": "invoke.function", "data": { "vote": "11am" } }
      }
    ]
  }
}
```

## API Integration

### Send Message

```python
POST https://cliq.zoho.com/api/v2/channelsbyname/{channel_name}/message
Headers:
  Authorization: Zoho-oauthtoken {access_token}
Body:
  {
    "text": "Hello from API!"
  }
```

### Get Messages

```python
GET https://cliq.zoho.com/api/v2/chats/{chat_id}/messages
```

### Create Channel

```python
POST https://cliq.zoho.com/api/v2/channels
{
  "name": "New Project",
  "description": "Project discussion channel",
  "type": "public"
}
```

## Mobile & Desktop Apps

### Desktop App (Mac, Windows, Linux)

- Native application
- Always-on notifications
- Better performance than browser
- Keyboard shortcuts
- Minimize to system tray

### Mobile App (iOS, Android)

- Push notifications
- Audio/video calls
- File uploads from camera/gallery
- Offline message drafts

## File Sharing

### Supported Files

- Documents (PDF, DOCX, etc.)
- Images (JPG, PNG, GIF)
- Videos (MP4, MOV)
- Archives (ZIP, RAR)
- Code files

### File Features

- Preview files in-chat
- Search files across channels
- File version history (if from WorkDrive)
- External sharing links

## Search & Organization

### Search Features

- Search messages, files, users
- Filter by: Channel, User, Date, File Type
- Search in specific channel or all
- Advanced operators (`from:user`, `has:file`)

### Organization

- Star important messages
- Pin messages to channel
- Bookmark channels
- Custom channel folders

## Integrations

- **Zoho CRM**: Link contacts, create tasks
- **Zoho Desk**: Create tickets from messages
- **Zoho Projects**: Task/milestone updates
- **Zoho Calendar**: Meeting reminders
- **Zoho Mail**: Email notifications
- **WorkDrive**: File sharing
- **Google Calendar**: Meeting sync
- **GitHub**: PR/issue notifications

## Notification Settings

### Notification Levels

- **All Messages**: Every message
- **@Mentions**: Only when mentioned
- **Nothing**: Muted

### Per-Channel Settings

- Different notification settings per channel
- Mute channels temporarily (1h, 4h, 24h, forever)
- Desktop notifications
- Mobile push notifications
- Email notifications (digest)

## Status & Presence

### Status Options

- 🟢 **Available** - Online and active
- 🟡 **Away** - Idle for 10+ minutes
- 🔴 **Busy** - In a meeting/focused work
- ⚫ **Offline** - Not logged in

### Custom Status

- Set custom status message
- Add emoji
- Set duration (30m, 1h, 2h, end of day)
- Example: "☕ Coffee break - back in 15"

## Legacy AI Integration

### Bot Commands for Legacy AI

```
/status
  → Check: API health, database status, queue depth

/deploy [environment]
  → Trigger: GitHub Actions deployment

/create task [description]
  → Add: New Todoist task

/search [query]
  → Search: Codebase or documentation

/issue [title]
  → Create: Linear issue
```

### Automated Notifications to Cliq

1. **High-value lead** (CRM) → Post to #sales
2. **P0 support ticket** (Desk) → Post to #support + @mention on-call
3. **GitHub PR review** → Post to #development
4. **Server alert** (monitoring) → Post to #alerts
5. **Daily standup reminder** → Post to #general

## Best Practices

- ✅ Use threads for long conversations (keep channels clean)
- ✅ @mention specific people (not @everyone unless urgent)
- ✅ Use appropriate channels (#development for dev, not #general)
- ✅ Pin important messages
- ✅ Set status when busy/away
- ✅ Use reactions instead of "+1" messages
- ✅ Keep files organized (use WorkDrive for long-term storage)
- ✅ Respect "Do Not Disturb" status

## Notes

- Free tier: Unlimited users, 100 GB storage/org
- Paid plans: $2.50/user/month (Professional), includes advanced features
- Video calls: Up to 100 participants (free tier: 25)
- Message history: Unlimited on all plans
- Bots & webhooks: Unlimited on paid plans
- Desktop/mobile apps free for all users

---

**Related**: [WorkDrive](./WorkDrive.md) | [Mail](./Mail.md) | [Desk](./Desk.md)

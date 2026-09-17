# Zoho Mail

> **Business Email & Calendar**

## Overview

Zoho Mail is an ad-free business email hosting service with custom domains, calendar, contacts, tasks, and notes in one integrated platform.

## Status

- **Active**: ✅ Yes
- **Integration**: IMAP/SMTP + API
- **Last Updated**: January 25, 2026

## Use Cases

1. **Business Email** - Professional email with custom domain
2. **Calendar Management** - Schedule meetings and events
3. **Contact Management** - Organize contacts
4. **Email Automation** - Filters, rules, templates
5. **Team Collaboration** - Shared mailboxes & calendars

## Email Accounts

| Email Address           | Purpose                        | Auto-Reply |
| ----------------------- | ------------------------------ | ---------- |
| **zeus@onemind.com**    | Personal inbox                 | No         |
| **support@onemind.com** | Customer support (→ Zoho Desk) | Yes        |
| **info@onemind.com**    | General inquiries              | No         |
| **noreply@onemind.com** | Automated emails only          | Yes        |
| **hello@onemind.com**   | Sales & new business           | No         |

## Email Features

### Composition

- Rich text editor
- Email templates
- Scheduled send
- Delivery tracking (read receipts)
- Priority inbox
- Undo send (30 seconds)
- Inline attachments
- Email signatures (per account/identity)

### Organization

- Folders & labels
- Smart filters (auto-sort)
- Search with operators
- Tags for emails
- Starred/important emails
- Archive (hide from inbox)

### Automation

- Filter rules (if/then)
- Auto-forward
- Auto-reply/vacation responder
- Scheduled emails
- Email reminders

## Email Filters & Rules

### Example Filters

**Filter 1: CRM Leads**

```
IF: From contains "@bigcompany.com"
THEN:
  - Label: "Hot Leads"
  - Star email
  - Forward to: sales@onemind.com
  - Create CRM lead (via Zoho CRM integration)
```

**Filter 2: Support Tickets**

```
IF: To is "support@onemind.com"
THEN:
  - Create Desk ticket
  - Auto-reply: "We received your message..."
  - Move to "Support" folder
```

**Filter 3: Newsletters**

```
IF: Has word "unsubscribe" in body
THEN:
  - Label: "Newsletters"
  - Skip inbox (auto-archive)
```

## Email Templates

### Template Categories

**Sales Templates**

1. Cold Outreach
2. Follow-up
3. Proposal Sent
4. Meeting Request
5. Thank You

**Support Templates**

1. Welcome Email
2. Password Reset
3. Issue Acknowledged
4. Resolution Provided
5. Feedback Request

**Internal Templates**

1. Weekly Update
2. Meeting Invite
3. Project Status
4. Policy Announcement

### Example Template: Meeting Request

```
Subject: Quick Call about [Topic]?

Hi [Name],

I'd love to schedule a quick 15-minute call to discuss [specific topic].

Would any of these times work for you?
- [Date/Time Option 1]
- [Date/Time Option 2]
- [Date/Time Option 3]

If none of these work, feel free to suggest a time that's convenient.

Best,
Zeus
```

## Calendar Integration

### Calendar Features

- Multiple calendars (Work, Personal, Projects)
- Shared team calendars
- Meeting scheduling (Zoho Bookings integration)
- Email meeting invites
- Calendar sync (Google Calendar, Apple Calendar)
- Recurring events
- Reminders (email, push, SMS)
- Time zone support

### Meeting Templates

- 1:1 Meeting (30 min)
- Team Standup (15 min)
- Client Call (1 hour)
- Workshop (2 hours)

## IMAP/SMTP Integration

### IMAP Settings (Incoming)

```
Server: imap.zoho.com
Port: 993 (SSL)
Username: zeus@onemind.com
Password: [app-specific password]
```

### SMTP Settings (Outgoing)

```
Server: smtp.zoho.com
Port: 465 (SSL) or 587 (TLS)
Username: zeus@onemind.com
Password: [app-specific password]
```

### Use Cases for IMAP/SMTP

- Connect desktop email clients (Apple Mail, Outlook, Thunderbird)
- Send emails via Python/Node.js scripts
- Integrate with Legacy AI for automated emails

## API Integration

### Send Email

```python
import requests

url = "https://mail.zoho.com/api/accounts/{accountId}/messages"
headers = {
    "Authorization": f"Bearer {access_token}",
    "Content-Type": "application/json"
}

payload = {
    "fromAddress": "zeus@onemind.com",
    "toAddress": "customer@example.com",
    "subject": "Thanks for signing up!",
    "content": "<h1>Welcome to OneMind</h1>"
}

response = requests.post(url, json=payload, headers=headers)
```

### List Emails

```python
GET https://mail.zoho.com/api/accounts/{accountId}/messages/view
```

### Search Emails

```python
GET https://mail.zoho.com/api/accounts/{accountId}/messages/search?searchKey=from:customer@example.com
```

## Email Signatures

### Primary Signature

```
Zeus Delacruz
Founder, OneMind OS
📧 zeus@onemind.com
🌐 onemind.zeusdelacruz.com
🐙 github.com/Zeus-Delacruz

🤖 Powered by AI | Building the future of personal operating systems
```

### Support Signature

```
OneMind Support Team
📧 support@onemind.com
📚 Help Center: docs.onemind.com

We typically respond within 24 hours.
```

## Spam & Security

### Anti-Spam Features

- Bayesian spam filter
- DNSBL (DNS-based blacklist)
- SPF/DKIM/DMARC verification
- Greylisting
- User-defined blacklist/whitelist

### Security Features

- Two-factor authentication (2FA)
- App-specific passwords
- IP restrictions
- Email encryption (TLS)
- S/MIME support (enterprise)
- Audit logs

## Email Workflow with Other Zoho Apps

### CRM Integration

- Save emails to CRM contact/lead
- Create tasks from emails
- Associate emails with deals
- Email tracking (opens, clicks)

### Desk Integration

- Emails to support@ → Auto-create ticket
- Reply to ticket from Mail
- Email ticket updates to customer

### Campaigns Integration

- Import contacts from Mail
- Send mass emails via Campaigns
- Track campaign opens/clicks

### WorkDrive Integration

- Send large files via WorkDrive links
- Save email attachments to WorkDrive
- Attach files from WorkDrive in emails

## Mobile & Desktop Apps

### Desktop App (Mac, Windows, Linux)

- Native application
- Offline access
- Faster than web interface
- Desktop notifications
- System tray integration

### Mobile App (iOS, Android)

- Push notifications
- Swipe actions (archive, delete, reply)
- Smart inbox
- Quick replies
- Touch ID/Face ID

## Email Migration

### Import from Other Providers

- Gmail (IMAP import)
- Outlook/Office 365
- Yahoo Mail
- Custom IMAP servers

### Migration Process

1. Connect source account
2. Select folders to import
3. Run migration (automatic)
4. Verify imported emails
5. Update DNS records (MX)

## Custom Domain Setup

### DNS Records Required

**MX Records** (Mail delivery)

```
Priority 10: mx.zoho.com
Priority 20: mx2.zoho.com
Priority 50: mx3.zoho.com
```

**SPF Record** (Sender verification)

```
v=spf1 include:zoho.com ~all
```

**DKIM Record** (Email authentication)

```
[Generated by Zoho upon setup]
```

**DMARC Record** (Policy enforcement)

```
v=DMARC1; p=quarantine; rua=mailto:dmarc@onemind.com
```

## Storage & Limits

| Plan             | Storage | Aliases | Max Attachment |
| ---------------- | ------- | ------- | -------------- |
| **Free**         | 5 GB    | 5       | 25 MB          |
| **Mail Lite**    | 10 GB   | 5       | 100 MB         |
| **Mail Premium** | 50 GB   | 10      | 1 GB           |
| **Current**      | 10 GB   | 5       | 100 MB         |

## Legacy AI Integration

### Automated Email Actions

1. **New Lead Email** → Create CRM lead, add to nurture campaign
2. **Support Email** → Create Desk ticket, analyze sentiment
3. **Meeting Request** → Check calendar, suggest times
4. **Invoice Received** → Extract data, log in accounting
5. **Newsletter Received** → Summarize key points, save to notes

### AI Email Assistant Features

- Draft email responses
- Summarize long email threads
- Suggest follow-up actions
- Detect sentiment (urgent, happy, frustrated)
- Auto-categorize emails

## Best Practices

- ✅ Use descriptive subject lines
- ✅ Keep emails concise (3-5 sentences)
- ✅ Use templates for common responses
- ✅ Set up filters to auto-organize inbox
- ✅ Use scheduled send for off-hours drafts
- ✅ Enable 2FA for account security
- ✅ Regularly clean up inbox (archive old emails)
- ✅ Use labels/folders instead of leaving in inbox
- ✅ Unsubscribe from unnecessary newsletters

## Notes

- Free tier: 5 users, 5 GB each, custom domain
- Paid plans: $1/user/month (Lite), $4/user/month (Premium)
- No ads on any plan (unlike Gmail)
- GDPR compliant with EU data centers
- Supports multiple domains per organization
- Mobile apps free for all users

---

**Related**: [Cliq](Cliq.md) | [CRM](CRM.md) | [Campaigns](Campaigns.md)


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]

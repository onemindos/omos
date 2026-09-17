# 🗂️ Legacy AI - Department Auto-Routing Flowchart (Email → Desk)


### 🗂️ Legacy AI - Department Auto-Routing Flowchart (Email → Desk)

This flowchart outlines the automated logic for routing all support and brand emails into **Legacy AI Support** via **Zoho Mail + Zoho Desk**, with clear tagging, department assignment, and notification paths.

This structure ensures:
- ✅ One single inbox: support@legacyhq.ai
- ✅ All brand emails routed and tagged by source
- ✅ Clean department segmentation in Zoho Desk
- ✅ AI/agent-ready ticket management

### 🌐 Inbound Flow: Overview

Customer Sends Email (e.g., [restore@legacyhq.ai](mailto:restore@legacyhq.ai))

↓

Zoho Mail Group (restore@) forwards → [support@legacyhq.ai](mailto:support@legacyhq.ai)

↓

Zoho Mail Filter matches source email

↓

Auto-tags email + moves to 1 main folder

↓

Post notification to relevant Cliq Channel (#Restore, #Acadacom, etc.)

↓

Zoho Desk receives via support@ → auto-tags + assigns Department

### ✅ Department Routing Table
| Source Email | Zoho Mail Filter | Zoho Desk Tag | Assigned Department | Post to Cliq |
| --- | --- | --- | --- | --- |
| restore@legacyhq.ai | Filter: From/To = restore@ | Tag: restore | Dept: Legacy Restore | #Restore |
| acadacom@legacyhq.ai | Filter: From/To = acadacom@ | Tag: acadacom | Dept: Legacy Acadacom | #Acadacom |
| capital@legacyhq.ai | Filter: From/To = capital@ | Tag: capital | Dept: Legacy Capital | #Capital |
| support@legacyhq.ai | Catch-all filter | Tag: general | Dept: Legacy AI Support | #Support-General |
| help@legacyhq.ai (if used) | Alias → support@ | Tag: general | Dept: Legacy AI Support | #Support-General |

### 🧠 Zoho Mail Filter Actions Per Domain
- ✅ Move to Folder: Legacy AI > Zoho Desk - Legacy AI Support
- ✅ Tag as: [restore, acadacom, capital, general]
- ✅ Post to Cliq Channel (based on tag)
- ✅ Forward to: support@legacyplatformholdingsllc.zohodesk.com (Zoho Desk System Address)

### 🧩 Zoho Desk Auto-Tagging + Department Logic
| Trigger | Action |
| --- | --- |
| Email arrives via support@ with tag: restore | Assign to Dept: Legacy Restore |
| Email arrives via support@ with tag: acadacom | Assign to Dept: Legacy Acadacom |
| Email arrives via support@ with tag: capital | Assign to Dept: Legacy Capital |
| All else | Assign to Dept: Legacy AI Support (General) |

### 🧠 Internal Ticket View for Agents
| View | Filtered By |
| --- | --- |
| Legacy Restore Tickets | Dept = Legacy Restore |
| Acadacom Queue | Dept = Legacy Acadacom |
| Capital Tickets | Dept = Legacy Capital |
| General Support | Dept = Legacy AI Support |

### 🔁 Summary of Routing Logic

**All brand/product-specific support messages**:
- Start at their domain inbox (acadacom@, etc.)
- Are auto-forwarded into support@legacyhq.ai
- Are tagged by their source email in Zoho Mail
- Are routed to the correct Desk department automatically
- Trigger a matching Zoho Cliq notification for real-time alerts

### 🛡️ You now operate:
- One inbox, many departments
- Tag-based routing
- AI-agent and self-help ready helpdesk infrastructure
- Circle.so compatible (all forms route to support@ and trigger correct flows)

Let’s deploy this routing now. You’re ready to scale support for thousands, with total clarity and zero inbox chaos.

### Legacy HQ: Founder-Class Systems. 🧠🚀

---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]

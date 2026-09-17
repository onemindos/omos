---
id: ui-08-20260113225954
title: "PLATFORM CURATION"
type: reference
path: UI
domain: 00-24
folder: "08"
status: active
created: 2026-01-13
modified: 2026-01-13
owner: zeus
tags: [path/ui, path/ui/08, type/reference, status/active]
ai_summary: "PLATFORM CURATION"
---
h# Platform Curation - Tech Stack Selection

**Status**: 🔍 Under Review
**Date**: 2026-01-12
**Total Platforms**: 255 items in 1Password
**Purpose**: Select which platforms to integrate into OneMind Cognitive Operating Layer

---

## Instructions

Review each platform and mark with:
- ✅ **KEEP** - Integrate into OneMind (build connector + tools)
- 🔍 **EVALUATE** - Research more before deciding
- ❌ **DELETE** - Remove from 1Password (no longer used)
- 💤 **ARCHIVE** - Keep credentials but no active integration needed

---

## Summary by Category

| Category | Count | Already Integrated |
|----------|-------|-------------------|
| AI/ML & Development | 45+ | Letta ✅, GitHub ✅ |
| Productivity & Tools | 20+ | Todoist ✅, Obsidian ✅ |
| Business & CRM | 15+ | Zoho One ✅ |
| Communication | 10+ | Discord ✅, Twilio ✅ |
| Learning & Courses | 30+ | None |
| E-commerce & Shopping | 25+ | None |
| Health & Wellness | 15+ | None |
| Financial & Business Ops | 20+ | None |
| Entertainment & Gaming | 10+ | None |
| Home & Local Services | 15+ | None |
| Affiliates & Marketing | 20+ | None |
| Miscellaneous | 30+ | None |

---

## 1. AI/ML & Development Platforms

### Core Development (Current Integration)
| Platform | Status | Current Tools | Decision |
|----------|--------|---------------|----------|
| **Letta** | ✅ INTEGRATED | Full API access, 4 agents | ✅ KEEP |
| **GitHub** | ✅ INTEGRATED | API credential | ✅ KEEP |
| **Todoist** | ✅ INTEGRATED | 7 atomic tools | ✅ KEEP |
| **Discord** | ✅ INTEGRATED | Notification tools | ✅ KEEP |
| **Twilio** | ✅ INTEGRATED | SMS/Voice tools | ✅ KEEP |
| **Obsidian** | ✅ CORE | Vault platform | ✅ KEEP |

### AI Platforms (Evaluate for Integration)
| Platform | Description | Recommendation | Decision |
|----------|-------------|----------------|----------|
| **OpenRouter** | Unified LLM API access | 🔍 EVALUATE - May replace multiple AI services | |
| **Voice Flow** | Conversational AI builder | 🔍 EVALUATE - Good for voice interfaces | |
| **Descript** | AI video/audio editing | 🔍 EVALUATE - Content creation automation | |
| **Heygen** | AI video generation | 🔍 EVALUATE - Video content automation | |
| **Agora** | Real-time engagement platform | 🔍 EVALUATE - For live interactions | |
| **Deep Learning AI** | AI learning platform | 💤 ARCHIVE - Educational resource | |

### Development Tools
| Platform | Description | Recommendation | Decision |
|----------|-------------|----------------|----------|
| **Replit** | Online IDE | 🔍 EVALUATE - For remote development | |
| **Cursor** | AI code editor | 💤 ARCHIVE - Local tool, no API needed | |
| **Windsurf** | Code collaboration | ❌ DELETE? - Redundant with GitHub? | |
| **Vercel** | Deployment platform | 🔍 EVALUATE - For web deployments | |
| **AWS** | Cloud infrastructure | 🔍 EVALUATE - May need for scaling | |
| **Elestio** | Managed app hosting | 🔍 EVALUATE - Self-hosted alternatives | |

### Design & Collaboration
| Platform | Description | Recommendation | Decision |
|----------|-------------|----------------|----------|
| **Coda** | Docs + apps | ❌ DELETE? - Have Obsidian + Notion | |
| **Notion** | Knowledge management | 🔍 EVALUATE - Backup to Obsidian? | |
| **Miro** | Whiteboarding | 💤 ARCHIVE - Occasional use | |
| **Linear** | Issue tracking | ❌ DELETE? - Use GitHub Issues | |
| **Webflow** | Website builder | 🔍 EVALUATE - For landing pages | |

---

## 2. Productivity & Organization

### Task Management
| Platform | Description | Recommendation | Decision |
|----------|-------------|----------------|----------|
| **Todoist** | Task manager | ✅ KEEP - Core system | |
| **TickTick** | Task manager | ❌ DELETE - Redundant with Todoist | |
| **Habitica** | Gamified tasks | ❌ DELETE - Not aligned with system | |

### Note-Taking & PKM
| Platform | Description | Recommendation | Decision |
|----------|-------------|----------------|----------|
| **Obsidian** | Markdown vault | ✅ KEEP - Core vault | |
| **Notion** | All-in-one workspace | 🔍 EVALUATE - Backup/migration? | |
| **Raindrop** / **Raindrop.io** | Bookmark manager | ❌ DELETE - Duplicate entries, use Obsidian | |
| **Substack** | Newsletter platform | 💤 ARCHIVE - Publishing only | |

### File Management
| Platform | Description | Recommendation | Decision |
|----------|-------------|----------------|----------|
| **Setapp** | Mac app subscription | 💤 ARCHIVE - Local tool bundle | |

---

## 3. Business & CRM

### Business Operations (Current)
| Platform | Description | Recommendation | Decision |
|----------|-------------|----------------|----------|
| **Zoho One** | CRM/Desk/Calendar | ✅ KEEP - Business suite | |
| **Circle.so** | Community platform | 🔍 EVALUATE - Community management | |

### E-commerce & Fulfillment
| Platform | Description | Recommendation | Decision |
|----------|-------------|----------------|----------|
| **Shopify** | E-commerce | 🔍 EVALUATE - If running stores | |
| **Pirateship.com** | Shipping | 🔍 EVALUATE - With e-commerce | |
| **Sedo** | Domain marketplace | 💤 ARCHIVE - Occasional domain sales | |
| **Alibaba** | Wholesale sourcing | 💤 ARCHIVE - Occasional use | |
| **Envato** | Digital assets | 💤 ARCHIVE - Occasional purchases | |

### Marketing & Design
| Platform | Description | Recommendation | Decision |
|----------|-------------|----------------|----------|
| **Fiverr - Legacy AI** | Freelancer platform | 💤 ARCHIVE - Hiring only | |
| **Vista Print** | Printing services | 💤 ARCHIVE - Occasional use | |
| **Divi WP Theme** | WordPress theme | 💤 ARCHIVE - Legacy websites | |

---

## 4. Communication Platforms

### Integrated
| Platform | Description | Recommendation | Decision |
|----------|-------------|----------------|----------|
| **Discord** | Team communication | ✅ KEEP - Notifications/approvals | |
| **Twilio** | SMS/Voice API | ✅ KEEP - Communication automation | |

### Email
| Platform | Description | Recommendation | Decision |
|----------|-------------|----------------|----------|
| **Zeus Ops - Gmail** | Operations email | ✅ KEEP - May need email tools | |

---

## 5. Learning & Education

### Online Courses
| Platform | Description | Recommendation | Decision |
|----------|-------------|----------------|----------|
| **Udemy** | Course platform | 💤 ARCHIVE - Learning resource | |
| **Skool - All** | Community learning | 💤 ARCHIVE - Learning resource | |
| **Deep Learning AI** | AI courses | 💤 ARCHIVE - Learning resource | |
| **We Are Nocode** | No-code training | 💤 ARCHIVE - Learning resource | |
| **Divit Society** | Community | 💤 ARCHIVE - Learning resource | |
| **Off Grid Dream** | Homesteading | 💤 ARCHIVE - Learning resource | |
| **The Airbnbae Academy** | Real estate training | 💤 ARCHIVE - Learning resource | |

### Credit & Finance Education
| Platform | Description | Recommendation | Decision |
|----------|-------------|----------------|----------|
| **creditsecret.org** | Credit education | 💤 ARCHIVE - Learning resource | |
| **creditlaunchsecrets.com** (2 entries) | Credit training | ❌ DELETE - Duplicate, consolidate | |
| **secretstrilogy.com** | Finance training | 💤 ARCHIVE - Learning resource | |
| **videoschallenge.com** | Marketing training | 💤 ARCHIVE - Learning resource | |
| **milliondollarads.com** | Ads training | 💤 ARCHIVE - Learning resource | |
| **socialmediacodex.com** | Social media training | 💤 ARCHIVE - Learning resource | |

### Sustainability & Homesteading
| Platform | Description | Recommendation | Decision |
|----------|-------------|----------------|----------|
| **thegrownetwork.com** | Homesteading | 💤 ARCHIVE - Learning resource | |
| **growyourownvegetables.org** | Gardening | 💤 ARCHIVE - Learning resource | |
| **takebackyourpower.net** | Off-grid living | 💤 ARCHIVE - Learning resource | |
| **thesovereignsway.com** | Self-sufficiency | 💤 ARCHIVE - Learning resource | |

---

## 6. E-commerce & Shopping

### Major Retailers
| Platform | Description | Recommendation | Decision |
|----------|-------------|----------------|----------|
| **Amazon.com** | Online shopping | 💤 ARCHIVE - Personal shopping | |
| **eBay** / **Ebay - Zeus** | Marketplace | ❌ DELETE - Duplicate, keep one | |
| **Staples** | Office supplies | 💤 ARCHIVE - Personal shopping | |
| **CVS** | Pharmacy | 💤 ARCHIVE - Personal shopping | |
| **Instacart.com** | Grocery delivery | 💤 ARCHIVE - Personal shopping | |

### Specialty Retailers
| Platform | Description | Recommendation | Decision |
|----------|-------------|----------------|----------|
| **Starbucks** | Coffee | 💤 ARCHIVE - Personal account | |
| **Shakeshack** | Restaurant | 💤 ARCHIVE - Personal account | |
| **Regmovies** | Movie theater | 💤 ARCHIVE - Personal account | |
| **Cooper Hawks Membership** | Restaurant | 💤 ARCHIVE - Personal membership | |
| **Airsoft GI** | Airsoft gear | 💤 ARCHIVE - Personal shopping | |
| **Offer up** | Local marketplace | 💤 ARCHIVE - Personal selling | |

---

## 7. Health & Wellness

### Cannabis Dispensaries
| Platform | Description | Recommendation | Decision |
|----------|-------------|----------------|----------|
| **Fine Fettle Dispensary** | Cannabis | 💤 ARCHIVE - Personal account | |
| **Curaleaf** | Cannabis | 💤 ARCHIVE - Personal account | |
| **Mighty Kind** | Cannabis | 💤 ARCHIVE - Personal account | |

### Health Services
| Platform | Description | Recommendation | Decision |
|----------|-------------|----------------|----------|
| **plushcare.com** | Telemedicine | 💤 ARCHIVE - Personal healthcare | |
| **Letsgetchecked** | Lab testing | 💤 ARCHIVE - Personal healthcare | |
| **Access Health CT** | Health insurance | 💤 ARCHIVE - Personal insurance | |
| **doterra.com** | Essential oils | 💤 ARCHIVE - Personal wellness | |
| **Schedule35** | Wellness | 💤 ARCHIVE - Personal wellness | |

### Fitness & Tracking
| Platform | Description | Recommendation | Decision |
|----------|-------------|----------------|----------|
| **Strava** | Fitness tracking | 💤 ARCHIVE - Personal fitness | |
| **AllTrails** | Hiking | 💤 ARCHIVE - Personal recreation | |
| **Equinox plus** | Gym | 💤 ARCHIVE - Personal fitness | |
| **Rise Above - Floatation** | Float therapy | 💤 ARCHIVE - Personal wellness | |

### Astrology
| Platform | Description | Recommendation | Decision |
|----------|-------------|----------------|----------|
| **Co-Star** | Astrology | 💤 ARCHIVE - Personal interest | |

---

## 8. Financial & Business Operations

### Business Management
| Platform | Description | Recommendation | Decision |
|----------|-------------|----------------|----------|
| **Zoho one** | Business suite | ✅ KEEP - Core CRM/business | |
| **Wyoming LLC Attorney Portal** | Legal | 💤 ARCHIVE - Legal access | |
| **Greenwich Resident Account** | Municipal | 💤 ARCHIVE - Local government | |

### Credit & Financial Tools
| Platform | Description | Recommendation | Decision |
|----------|-------------|----------------|----------|
| **Metro 2 Compliance** | Credit reporting | 🔍 EVALUATE - If managing credit | |
| **Dispute fox** | Credit disputes | 🔍 EVALUATE - If managing credit | |

### Tracking & Monitoring
| Platform | Description | Recommendation | Decision |
|----------|-------------|----------------|----------|
| **Bouncie** | Vehicle tracking | 💤 ARCHIVE - Car monitoring | |

---

## 9. Entertainment & Gaming

### Gaming Platforms
| Platform | Description | Recommendation | Decision |
|----------|-------------|----------------|----------|
| **Steam** | Game platform | 💤 ARCHIVE - Personal gaming | |
| **Ubisoft** | Game platform | 💤 ARCHIVE - Personal gaming | |
| **Instant Gaming** | Game keys | 💤 ARCHIVE - Personal gaming | |

### Media
| Platform | Description | Recommendation | Decision |
|----------|-------------|----------------|----------|
| **getwsodott.com** / **getwsodo.com** | Torrent files | ❌ DELETE - Duplicate entries | |
| **Pateron** | Creator support | 💤 ARCHIVE - Supporting creators | |

---

## 10. Affiliates & Marketing Programs

### Affiliate Networks
| Platform | Description | Recommendation | Decision |
|----------|-------------|----------------|----------|
| **Partnerstack** | Affiliate platform | 🔍 EVALUATE - If doing affiliate marketing | |
| **GoHighlevel Affiliate** (2 entries) | Marketing automation | ❌ DELETE - Duplicate, keep one | |

### Brand Affiliates
| Platform | Description | Recommendation | Decision |
|----------|-------------|----------------|----------|
| **Creditcardbroker.com Affiliate** | Credit affiliate | 💤 ARCHIVE - Passive income? | |
| **Swanwick Affiliate** | Blue light glasses | 💤 ARCHIVE - Passive income? | |
| **IDIQ Affiliate** | Health affiliate | 💤 ARCHIVE - Passive income? |
| **Joinsecret** | Finance affiliate | 💤 ARCHIVE - Passive income? | |
| **Monarch For Professionals** | Financial planning | 💤 ARCHIVE - Financial tool | |

---

## 11. MLM & Network Marketing

### MLM Platforms
| Platform | Description | Recommendation | Decision |
|----------|-------------|----------------|----------|
| **Enagic Distributor (Support Portal)** | Kangen water | 💤 ARCHIVE - MLM business | |
| **enagicwebsystem.com** | Enagic system | ❌ DELETE - Duplicate? | |
| **multipure.com - (backoffice)** | Water filters | 💤 ARCHIVE - MLM business | |

---

## 12. Miscellaneous & Uncategorized

### Project Management
| Platform | Description | Recommendation | Decision |
|----------|-------------|----------------|----------|
| **infinity - Project Manager** | Project management | ❌ DELETE - Use Todoist/GitHub | |
| **geru.com** | Unknown | ❌ DELETE - Investigate, likely unused | |

### Exotic Niche
| Platform | Description | Recommendation | Decision |
|----------|-------------|----------------|----------|
| **exoticcarhacks.com** | Car training | 💤 ARCHIVE - Learning resource | |
| **Exotic Car hacks** (duplicate) | Car training | ❌ DELETE - Duplicate | |
| **Fund Launch** | Funding platform | 💤 ARCHIVE - Research resource | |

### Food & Catering
| Platform | Description | Recommendation | Decision |
|----------|-------------|----------------|----------|
| **Garden Catering** | Local catering | 💤 ARCHIVE - Personal/business events | |

---

## Next Steps

1. **Review each section** and mark your decisions in the "Decision" column
2. **Identify priorities** - Which platforms need integration ASAP?
3. **Delete duplicates** - Mark duplicate entries for cleanup
4. **Archive personal** - Personal shopping/gaming accounts don't need API integration
5. **Evaluate business tools** - Which business platforms align with GE (75-99) goals?

## Integration Priority Framework

### Tier 1: Core Infrastructure (Already Done)
- Letta (AI brain)
- Todoist (task management)
- Discord (notifications)
- GitHub (code/docs)
- Obsidian (vault)

### Tier 2: High Priority (Next)
- Zoho One (CRM/business operations)
- Home Assistant (smart home - in progress)
- Email integration (Gmail API)
- OpenRouter (unified LLM access)

### Tier 3: Medium Priority
- Notion (backup/migration)
- Vercel (deployments)
- Twilio (enhanced communication)
- VoiceFlow (voice interfaces)

### Tier 4: Future Evaluation
- E-commerce stack (if running stores)
- Marketing automation
- Content creation tools
- Analytics platforms

---

## Command to Delete Platform

```bash
# Delete a platform from 1Password
op item delete "Platform Name"

# Or by ID
op item delete ITEM_ID
```

## Command to Export Credentials

```bash
# Get credentials for keeping
op item get "Platform Name" --format json

# Export API key/token
op item get "Platform Name" --fields label=password --reveal
```

---
id: ui-15-20260111053926
title: "The definitive bookmark management platform guide for Legacy AI (2025)"
type: concept
pillar: UI
domain: 00-24
folder: "00-24 UI (Unified Intelligence)"
status: active
created: 2026-01-11
modified: 2026-01-11
owner: zeus
tags: [domain/ui, folder/15]
todoist_project: "00-24 UI (Unified Intelligence)"
ai_summary: "The definitive bookmark management platform guide for Legacy AI (2025)"
---

# The definitive bookmark management platform guide for Legacy AI (2025)

**Two major competitors just died**: Omnivore shut down November 30, 2024 and Pocket closed July 8, 2025. This dramatically reshapes the landscape and makes your platform choice more critical than ever.

## Top 3 recommendations for your use case

After evaluating 15+ platforms against your specific requirements (8.5k bookmarks, 50-100 daily additions, AI features, n8n integration, Obsidian workflow), here are the clear winners:

### 1. Readwise Reader — Best overall for knowledge workers
**The comprehensive champion with mature AI and integrations**

Readwise Reader dominates for serious knowledge management in 2025. With both Pocket and Omnivore dead, Reader becomes the undisputed leader for content-first workflows. **Reader replaces 5+ tools** in one platform: bookmark manager, RSS reader, PDF annotator, YouTube transcript tool, and newsletter processor.

**Why it wins for you**: Native Obsidian plugin with Jinja2 templating exports perfect markdown with frontmatter. Ghostreader AI (GPT-5 Mini included) handles summarization and Q&A. The full API (v3) integrates cleanly with n8n via HTTP nodes. Proven scale — handles 10k+ items without performance degradation. The iOS app captures content instantly via share sheet with offline support.

**AI capabilities shine**: Auto-summarization, semantic full-text search, custom AI prompts, and Chat with Highlights across your entire library. Import from Raindrop via CSV upload works seamlessly. Can process YouTube with time-synced transcripts and extract full articles better than Pocket or Instapaper ever did.

**The catch**: Premium pricing at $96/year (billed annually). No native n8n node requires HTTP Request setup. Instagram and Facebook imports need workarounds. But you're replacing $200+ worth of separate tools.

**Score**: 9.2/10 — Production-ready, bootstrap-funded (sustainable), active weekly updates, 30-day trial

### 2. LinkWarden — Best self-hosted option with AI
**Open-source powerhouse for privacy-focused technical users**

LinkWarden emerged as the go-to self-hosted solution after Omnivore's shutdown. It's purpose-built for web archiving with AI auto-tagging across multiple providers (OpenAI, Claude, Ollama locally). **If you run Portainer already, this is your answer**.

**Why it's compelling**: Full-page preservation in 4 formats (screenshot, PDF, HTML, readable text) protects against link rot. Meilisearch powers semantic full-text search across archived content. REST API with OpenAPI spec enables n8n workflows. Self-hosted means unlimited storage and complete control — crucial for 8.5k+ existing items.

**AI implementation**: Configure OpenAI, Claude, or local Ollama models for auto-tagging. The system can tag existing links in bulk or auto-tag new captures. Advanced search with operators (`title:`, `tag:`, `before:`, `after:`) rivals commercial tools.

**Trade-offs**: No official mobile apps yet (announced "coming very soon"), currently PWA plus third-party apps. No Obsidian plugin requires building custom API-based workflows. Import limited to HTML bookmarks (works with Raindrop/Brave exports) but no direct social media imports. Large bulk imports (25MB+ HTML files) need splitting into batches.

**Infrastructure**: Docker-friendly with PostgreSQL + Meilisearch. Cloud option at $3/month or self-host on VPS ($5-10/month). Resource usage moderate at ~200MB RAM but plan for S3 storage as archives grow rapidly.

**Score**: 8.7/10 — Best value, active development (v2.13.1 January 2025), 15k+ GitHub stars

### 3. Raindrop.io Pro + n8n AI layer — Best pragmatic upgrade
**Enhance your current baseline with automation rather than replacing it**

Sometimes the best move is keeping what works and augmenting it strategically. Raindrop Pro ($28/year) already handles the basics well — you've proven it scales to 8.5k items. Rather than a risky migration, **layer n8n workflows on top** to add the AI features Raindrop lacks.

**The hybrid architecture**: Raindrop for capture and storage (mobile apps, browser extensions, permanent copies via Wayback Machine). n8n workflows triggered via Raindrop API webhooks process new bookmarks through OpenAI for tagging, send to review queue in Supabase, then export approved items to Obsidian markdown. You get enterprise-grade AI features at consumer pricing.

**Why this works**: Built-in n8n node for Raindrop makes integration trivial. The API (120 req/min) handles your volume easily. Raindrop's AI semantic search and multilingual capabilities (added April 2024) provide baseline intelligence while your n8n layer adds custom AI prompts and processing. Migration risk: zero.

**Implementation path**: Keep Raindrop for daily capture. Build n8n workflows for content extraction (Jina.ai), AI tagging (GPT-4o-mini at $0.15/1M tokens), human review queue (Supabase), and Obsidian sync. Total added cost: $10-15/month for API services plus n8n hosting.

**Limitations**: Still dependent on Raindrop's roadmap for core features. Solo developer sustainability concerns remain. No custom AI prompts in Raindrop itself. But proven reliability and lowest migration effort.

**Score**: 8.3/10 — Lowest risk, fastest implementation, familiar UX

## Feature comparison matrix

| **Feature** | **Readwise Reader** | **LinkWarden** | **Raindrop + n8n** | **KaraKeep** | **Emerging (Rabrain)** |
|------------|---------------------|----------------|-------------------|--------------|----------------------|
| **AI auto-tagging** | ✅ Excellent (GPT-5) | ✅ Multi-provider | ⚠️ Via n8n custom | ✅ OpenAI/Ollama | ✅ RAG-powered |
| **Semantic search** | ✅ Full-text | ✅ Meilisearch | ✅ Raindrop built-in | ✅ Full-text | ✅ Context-aware |
| **iOS app quality** | ⭐⭐⭐⭐⭐ Native | ⭐⭐⭐ PWA + 3rd party | ⭐⭐⭐⭐⭐ Native | ⭐⭐⭐⭐ Native | ⭐⭐⭐⭐ Native |
| **Raindrop import** | ✅ CSV/highlights | ✅ HTML export | N/A (is Raindrop) | ✅ HTML | ⚠️ Manual |
| **Social media import** | ❌ YouTube only | ❌ Limited | ⚠️ Manual/scripts | ❌ No | ❌ No |
| **Bulk capacity** | ✅ 10k+ proven | ⚠️ Needs robust VPS | ✅ 20k+ reports | ⚠️ OOM issues noted | ⚠️ Unproven |
| **API/n8n** | ✅ REST v3 (HTTP) | ✅ OpenAPI spec | ✅ Built-in node | ⚠️ Minimal docs | ⚠️ Limited |
| **Content archiving** | ✅ Full article | ✅✅✅ 4 formats | ✅ Permanent copies | ✅ Full page + video | ✅ Full text |
| **Obsidian plugin** | ✅✅ Official Jinja2 | ❌ API only | ⚠️ Community plugins | ❌ API only | ⚠️ Planned |
| **Cost (annual)** | $96-156 | $36 cloud / $60 VPS | $28 + $120 n8n | $0 (self-host) | Free beta |
| **Scale performance** | ✅ Excellent | ✅ Good w/ tuning | ✅ Proven 20k | ⚠️ Monitor memory | ⚠️ Unknown |
| **Active development** | ✅✅✅ Weekly | ✅✅✅ Weekly | ✅✅ Established | ✅✅ Very active | ✅ New but active |
| **Production ready** | ✅ Mature | ✅ Stable v2.13 | ✅ 10+ years | ⚠️ Heavy dev | ⚠️ Beta |

**Legend**: ✅ Full support | ⚠️ Partial/workaround needed | ❌ Not available | ✅✅ Superior implementation

## Critical findings that change everything

**Omnivore and Pocket are dead**. Omnivore's cloud shut down November 2024 after ElevenLabs acquisition. Pocket closed completely July 2025 when Mozilla killed the service. These were two of the most recommended alternatives just months ago. This proves the danger of depending on unprofitable consumer services.

**Notion is terrible for 8.5k bookmarks**. Performance degrades significantly above 5,000 database rows. Users report "seconds to tens of seconds" load times at 10k items. The May 2025 pricing change moved Notion AI to Business plan only ($240/year), making it expensive for personal use. Notion works for projects but fails as a dedicated bookmark manager.

**Capacities is not a bookmark manager**. Despite appearing in searches, it's a structured note-taking app with no web archiving, no reader mode, and a philosophy against "importing large amounts of data." Completely wrong tool for this job.

**Self-hosting is viable but demanding**. LinkWarden and KaraKeep require Docker expertise, ongoing maintenance, and careful resource allocation. KaraKeep has a known OOM (out of memory) bug with large smart lists. Budget 4-8GB RAM for serious scale, not the documented minimums.

## Migration strategy from your 6-platform chaos

Your current setup: 8.5k links scattered across Raindrop, Apple Notes, Brave bookmarks, YouTube saves, Instagram saves, Facebook saves. Here's the proven path to consolidation.

### Phase 1: Export everything (Week 1)

**Raindrop.io**: Settings → Export → JSON format. This captures all metadata, collections, tags, highlights. The export is instantaneous regardless of size.

**Brave browser**: Settings → Bookmarks → Export bookmarks → saves as standard HTML. This Netscape format imports into Raindrop, LinkWarden, and most platforms.

**Apple Notes**: Mac required. File → Export as → Markdown (macOS 13+). For older systems, use Exporter app. Warning: complex embedded content may not export cleanly. Extract URLs manually if needed via SQLite database at `~/Library/Group Containers/group.com.apple.notes/`.

**YouTube saves**: Google Takeout → YouTube data → playlists.json contains liked videos and saved playlists. Parse JSON to extract video URLs with titles.

**Instagram/Facebook saves**: No official export exists. Options: browser history export filtered by domain, third-party scrapers (use cautiously), or manual collection via share sheet during review phase. **Recommendation**: Prioritize future captures over historical social saves — these are low-value nostalgia typically.

**Consolidation step**: Merge all exports into single CSV with columns: URL, title, source, date, notes. Use Python/Node script or spreadsheet for deduplication by normalized URL (remove www, tracking params, trailing slashes).

### Phase 2: Clean and prioritize (Days 3-5)

**Deduplication**: Normalize URLs first — remove `?utm_source`, `www.`, trailing slashes. Hash the cleaned URL for comparison. Expect 10-15% duplicates across 6 sources.

**Dead link removal**: Run bulk URL checker (tools: Dead Link Checker, Broken Link Checker browser extension, or n8n workflow with HTTP status checks). Archive.org Wayback Machine can recover some dead links automatically in LinkWarden and Raindrop Pro.

**Priority tiers**: 
- High (500 items): Active projects, recent saves, starred/favorited
- Medium (2000 items): Last 12 months, specific topic areas  
- Low (6000 items): Archive, old saves, duplicate content

This staged approach prevents overwhelming yourself and tests the workflow at scale before full migration.

### Phase 3: Incremental import (Weeks 2-4)

**Week 2 — High priority wave (500 items)**:
- Import to chosen platform (Readwise CSV, LinkWarden HTML, or keep in Raindrop)
- Run through AI tagging workflow
- Test Obsidian export with 50 items
- Refine frontmatter template based on real data
- Fix any workflow errors before larger batches

**Week 3 — Medium priority (2000 items)**:
- Batch process 200-500 items at a time
- Monitor API rate limits (Jina.ai: 20 req/min free, OpenAI: 500 req/day tier 1)
- Review samples from each batch to verify quality
- Adjust AI prompts if categorization is poor

**Week 4 — Low priority archive (6000 items)**:
- Consider: Do you need all of these in active system?
- Alternative: Keep in cold storage (Raindrop archive, separate Obsidian folder)
- Import only if system proves stable with 2.5k items first
- This prevents database bloat and maintains performance

**Total timeline**: 3-4 weeks including testing and adjustments. Do NOT rush — a failed bulk import is painful to reverse.

### Phase 4: Establish daily workflow

**Capture layer** (multiple entry points):
- iOS Shortcuts → n8n webhook → staging database
- Browser bookmarklet (one-click JavaScript)
- Email forwarding (read@readwise.io pattern)
- Raindrop browser extension (if keeping as capture tool)
- Share sheet from any iOS app

**Processing automation** (every 30 minutes):
1. Fetch new items from staging queue
2. Extract content with Jina.ai Reader API
3. AI tagging with GPT-4o-mini (3-5 tags, category, summary)
4. Insert into review queue with status: pending

**Human review** (15-30 min daily):
- Review dashboard shows pending items with AI suggestions
- Approve (confirm tags, add notes) or reject (archive/delete)
- Bulk approval for high-confidence items
- Approved items flow automatically to Obsidian

**Obsidian sync** (hourly):
- Generate markdown from approved queue
- Write files to vault via Google Drive sync or Obsidian Local REST API
- Update database with sync status
- Verify no duplicates against existing vault

**Maintenance** (weekly):
- Review failed processing items
- Check dead links (automated scan)
- Audit AI tagging accuracy
- Prune old staging records

## Integration architecture: n8n + Obsidian in detail

### Core system design

**Database: Supabase** (recommended over Airtable)

Free tier provides 500MB PostgreSQL database with real-time subscriptions. Create three tables:

```sql
-- Bookmark staging and processing
CREATE TABLE bookmarks (
  id UUID PRIMARY KEY,
  url TEXT UNIQUE NOT NULL,
  title TEXT,
  content TEXT,
  ai_summary TEXT,
  ai_tags TEXT[],
  ai_category TEXT,
  review_state TEXT, -- pending/approved/rejected
  obsidian_synced BOOLEAN DEFAULT FALSE,
  source_platform TEXT,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Processing queue for async work
CREATE TABLE processing_queue (
  id UUID PRIMARY KEY,
  bookmark_id UUID REFERENCES bookmarks(id),
  task_type TEXT, -- 'extract' / 'ai_tag' / 'sync_obsidian'
  status TEXT, -- pending/processing/completed/failed
  attempts INT DEFAULT 0
);
```

This architecture prevents workflow timeouts and allows crash-resistant processing.

**n8n workflows** (5 separate workflows, not one monolith):

**Workflow 1: Capture**
- Trigger: Webhook (iOS Shortcut, bookmarklet)
- Normalize URL, check duplicates
- Insert to Supabase bookmarks table
- Send confirmation notification

**Workflow 2: Content Extraction**
- Trigger: Schedule every 30 min
- Fetch pending extraction jobs from queue
- Call Jina.ai Reader API: `https://r.jina.ai/{url}`
- Store markdown content in bookmarks table
- Handle failures with retry logic

**Workflow 3: AI Processing**  
- Trigger: Schedule every 30 min
- Fetch items with content but no AI tags
- Send to OpenAI with prompt:
  ```
  Analyze this bookmark and return JSON:
  {
    "category": "Technology/Finance/Health/etc",
    "tags": ["tag1", "tag2", "tag3"],
    "summary": "2-sentence summary",
    "priority": "high/medium/low"
  }
  
  URL: {url}
  Title: {title}
  Content: {content_excerpt}
  ```
- Store results in bookmarks table
- Update state to pending_review

**Workflow 4: Review Handler**
- Trigger: Webhook (approval link from email)
- Parse approval token
- Update review_state in Supabase
- Trigger next workflow via Supabase onRowChange

**Workflow 5: Obsidian Sync**
- Trigger: Supabase realtime (onRowChange, review_state='approved')
- Generate markdown with frontmatter
- Write to Google Drive (synced with Obsidian vault)
- Update obsidian_synced flag

**Why separate workflows?**: Each completes cleanly without hanging state. If any workflow crashes, Supabase maintains queue state. You can restart individual components without disrupting the entire system.

### Obsidian frontmatter schema

This YAML template provides rich queryable metadata for Dataview plugin and future AI features:

```yaml
---
title: "{{title}}"
url: "{{url}}"
captured: "{{timestamp}}"
source: "{{platform}}"
type: article  # article/video/tweet/product

# AI-generated
category: "{{ai_category}}"
tags:
  - "{{tag1}}"
  - "{{tag2}}"
  - "{{tag3}}"
ai_summary: "{{summary}}"

# Workflow
status: to-read  # to-read/reading/read/archived
priority: "{{ai_priority}}"
reviewed: "{{review_date}}"

# Optional rich metadata
author: "{{author}}"
domain: "{{domain}}"
read_time: {{minutes}}
---

# {{title}}

## AI Summary
{{ai_summary}}

## Content
{{extracted_markdown}}

## Notes
<!-- Your manual notes here -->

## Metadata
- Captured: [[{{capture_date}}]]
- Source: {{source_platform}}
- Tags: {{tags_inline}}
```

**Key frontmatter decisions**:
- **Date format**: ISO 8601 (`YYYY-MM-DD`) for Dataview compatibility
- **Tags**: Both frontmatter array and inline `#tag` for flexibility
- **Status tracking**: Enable processing queue workflows
- **AI fields**: Separate from human edits for transparency

### Obsidian sync methods compared

**Method 1: Google Drive Sync** (recommended for simplicity)
- n8n writes markdown files to Google Drive folder
- Create symbolic link: `ln -s ~/Google\ Drive/ObsidianInbox ~/Obsidian/Vault/Inbox`
- Obsidian detects new files automatically
- Pros: No Obsidian plugins, works offline after sync, handles batch imports
- Cons: Requires Google Drive desktop app, sync delay 1-5 minutes

**Method 2: Obsidian Post Webhook Plugin** (best for real-time)
- Install plugin: obsidian-post-webhook by Masterb1234
- Configure webhook endpoint in plugin settings
- n8n POSTs markdown directly to Obsidian
- Pros: Instant sync, bidirectional capability, no external dependencies
- Cons: Requires Obsidian running, plugin maintenance, local network only (unless Tailscale)

**Method 3: Git-based sync** (best for version control)
- n8n commits markdown to GitHub repo
- Obsidian Git plugin pulls changes
- Pros: Full version history, works remotely, multiple device sync
- Cons: More complex setup, requires Git knowledge, potential merge conflicts

**Recommendation**: Start with Google Drive for reliability during migration. Switch to Post Webhook after system is stable if you want real-time performance.

### Rate limits and scaling considerations

**API costs for your volume** (8.5k migration + 50-100 daily):

**One-time migration**:
- Content extraction (Jina.ai): 8500 × $0.002 = $17
- AI tagging (GPT-4o-mini): 8500 × $0.001 = $8.50
- Total: ~$25-30 for full backlog

**Monthly operations** (75 items/day average):
- Jina.ai: 2250/month × $0.002 = $4.50
- OpenAI: 2250/month × $0.001 = $2.25
- n8n executions: Included in Pro plan ($20/month) or free self-hosted
- Supabase: Free tier sufficient
- Total: $7-27/month depending on n8n hosting choice

**Rate limit management**:
- Jina.ai free: 20 req/min → ~1200/hour. Batch daily captures for free tier.
- OpenAI tier 1: 500 req/day. Upgrade to tier 2 ($50+ spend) unlocks higher limits.
- Batch processing: Process 100 items, wait 5 minutes, repeat. Prevents hitting limits.

**Deduplication at scale**: Before processing 8.5k items, normalize URLs and check against existing vault. Use n8n Code node with JavaScript:

```javascript
const url = new URL($json.url);
const normalized = `${url.protocol}//${url.hostname}${url.pathname}`
  .replace(/\/$/, '')
  .toLowerCase();
return { normalized_url: normalized };
```

Then query Supabase for existing normalized URLs. Skip processing if found. This prevents wasting API calls on duplicates.

## Emerging platforms worth watching (but not ready yet)

**Rabrain** shows promise with RAG-powered conversations and voice memo integration. The CODE methodology alignment (Capture-Organize-Distill-Express) fits knowledge management workflows. But API documentation is sparse and import capabilities are unproven at your scale. Monitor for 6-12 months.

**Screvi** positions as the modern Readwise alternative with Instagram-style highlight review. At $35/year it undercuts Reader significantly while offering Anki-style spaced repetition and OCR for physical books. The iOS app has good reviews. However, launched recently (v1.7.0 October 2025) means it lacks Reader's maturity and integration ecosystem. Consider if budget is primary constraint.

**bundleIQ** offers the most sophisticated API among emerging platforms. The proprietary RAG system and semantic search API appeal for developer integrations. Knowledge graph visualization is unique. But $120/year pricing approaches Reader without the proven Obsidian plugin and scale performance. Better for enterprise teams than solo knowledge work.

**Markwise** emphasizes zero-friction auto-organization — appealing concept but light on specifics about AI model, accuracy metrics, or scale handling. "Trusted by 1,000+ users" is tiny compared to Readwise's hundreds of thousands.

**Reseek** nails Twitter/X archiving with full metadata preservation. If you're heavy on social media research, this fills a real gap. OCR for academic PDFs also strong. Currently free during testing phase makes it risk-free to try alongside primary platform. But API is undocumented and production pricing unknown.

**General assessment**: These emerging platforms demonstrate exciting innovation (RAG conversations, social media archiving, visual organization) but lack the battle-tested stability, integration ecosystems, and scale proof you need for 8.5k+ items. Readwise Reader's 10+ year history (counting legacy Readwise) and bootstrap funding provides sustainability none of these new entrants have proven yet.

## The platform you should actually choose

**For maximum power and minimum risk: Readwise Reader**

Despite the $96/year cost, Reader provides the most complete solution with lowest migration risk. The official Obsidian plugin eliminates custom workflow building. Proven scale performance removes the "will this work with 8.5k items?" anxiety. The 30-day trial lets you test with 500 high-priority bookmarks before committing.

Import path: Export Raindrop as CSV → Upload to Reader → AI processes automatically → Sync to Obsidian with Jinja2 templates. Done. Brave bookmarks import via HTML. YouTube saves via URL list upload. The painful social media imports (Instagram/Facebook) are low-priority nostalgia anyway.

Your n8n workflows become simpler: iOS capture → webhook → Reader API. Review queue optional since Reader's inbox and feed views provide natural review flow. The $8/month effective cost is less than your current scattered system's cognitive overhead.

**For technical users wanting control: LinkWarden + n8n**

If you're already running Portainer and Supabase, self-hosting LinkWarden demonstrates commitment to data sovereignty. The four-format archiving protects against link rot better than any cloud service. Your Docker expertise makes the ~200MB RAM overhead and S3 storage management trivial.

Build custom Obsidian integration via LinkWarden's REST API. The OpenAPI spec makes n8n workflows straightforward. Choose local Ollama models for AI tagging if privacy is paramount, or use OpenAI API for better accuracy.

Import: Raindrop HTML export → LinkWarden (works well), Brave bookmarks → HTML → LinkWarden (native support). Content extraction runs as background job. Plan 2-3 weeks for migration including custom workflow building.

Total cost: $60-120/year VPS or $36/year managed cloud. But invest significant time in setup and maintenance.

**For pragmatic upgraders: Enhanced Raindrop**

Keep your current Raindrop Pro as capture and storage layer. It already works, your mobile muscle memory is established, and the $28/year cost is proven value. Layer n8n intelligence on top rather than migrating.

Build the missing pieces: n8n workflow queries new Raindrop items hourly → Jina.ai content extraction → OpenAI tagging with custom prompts → Supabase review queue → approved items to Obsidian markdown. You get enterprise AI capabilities at consumer pricing.

This hybrid approach provides 80% of Reader's value at 40% of the cost. The trade-off: you're building and maintaining custom workflows instead of using supported integrations. But you've chosen this research because you're technical enough for n8n, so this is viable.

Import: Already complete — your 8.5k items are in Raindrop. Focus on the n8n processing pipeline and Obsidian export workflow. Implementation time: 1-2 weeks.

## What to do Monday morning

**Week 1 action plan** (choose your path):

**If going Readwise Reader**:
1. Start 30-day trial immediately
2. Export high-priority batch from Raindrop (500 items) as CSV
3. Import CSV to Reader, test AI processing
4. Install Obsidian plugin, configure Jinja2 template with frontmatter schema
5. Process first 50 approved items to Obsidian
6. Evaluate: Does this solve the problem? If yes, subscribe and continue migration.

**If going LinkWarden self-hosted**:
1. Deploy LinkWarden via Docker Compose on your Portainer stack
2. Configure environment variables for Meilisearch and PostgreSQL
3. Set up OpenAI API key for AI tagging (or Ollama locally)
4. Import Brave bookmarks HTML (test with 100 items first)
5. Build n8n workflow: LinkWarden API → Fetch new bookmarks → Process → Obsidian markdown
6. Test end-to-end with 50 items before bulk migration

**If enhancing Raindrop**:
1. Keep Raindrop Pro unchanged (your capture layer works)
2. Set up Supabase free tier with bookmark queue tables
3. Create first n8n workflow: Raindrop API → New items → Supabase staging
4. Add Jina.ai content extraction workflow
5. Build AI tagging workflow with OpenAI
6. Test review queue with 20 items before scaling

**Critical success factors** regardless of choice:
- Start with 50-100 items, not full 8.5k
- Test the complete pipeline before bulk importing
- Document your workflows and credentials immediately
- Set up monitoring for failed processing items
- Budget 2-4 weeks for migration, not 2-4 days

The platform matters less than the workflow discipline. Your "Legacy AI" system succeeds when bookmarks flow reliably through capture → AI organization → human review → Obsidian knowledge base. Choose the platform matching your technical comfort level and maintenance appetite, then execute the migration methodically.

The 8.5k bookmark chaos ends now. Start Monday with one high-priority batch and prove the system works. Scale from evidence, not hope.
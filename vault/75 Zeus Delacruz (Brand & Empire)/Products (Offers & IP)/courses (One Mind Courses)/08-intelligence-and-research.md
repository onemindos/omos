# Course 8: INTELLIGENCE & RESEARCH

> **Build an intelligence operation that keeps you ahead of everything.**

**Tier:** Legacy ($97/mo)
**Circle Space:** 🤖 AI & Agents (ID: 2617142)
**Prerequisite:** Courses 1–3 — vault live, agents running, OpenClaw configured
**Target Student:** Someone who wants their system to feed them relevant intelligence — market signals, AI developments, research synthesis, competitive awareness — without spending hours scrolling and reading
**Outcome:** A fully automated intelligence pipeline: signals collected, research synthesized, briefings delivered, knowledge captured in the vault — all without you doing the information gathering manually

---

## What Intelligence Operations Mean for One Mind OS

Domain 06 (Research / R&D) is the sense organ of the system. While you execute, it watches.

Three functions:
1. **Signals** — what's happening in the world that matters to your domains
2. **Research** — deep investigation of specific questions that require synthesis
3. **Knowledge capture** — turning raw intelligence into permanent vault knowledge

Most people consume information reactively. They scroll. They get distracted. They forget 90% of what they read. This course builds the opposite: a system that collects and synthesizes intelligence on your behalf, then surfaces only what matters.

---

## What You'll Build

By the end of this course:
- An RSS/feed pipeline that monitors your key information domains automatically
- A Trinity agent briefing that delivers curated intelligence to Discord daily
- A research protocol for deep investigation — from question to synthesized knowledge in one session
- A "what's new in AI" pipeline that keeps you current on model releases, tools, and developments
- A competitive intelligence feed for your GE domains
- A vault knowledge capture workflow that converts research into permanent notes
- An R&D backlog — a prioritized list of what to investigate next

---

## Course Structure — 6 Modules, 33 Lessons

---

### MODULE 1: THE INTELLIGENCE MINDSET

**Lesson 1.1 — From Information Consumer to Intelligence Operator**
The difference between passively consuming information and running an active intelligence operation.
- The information overload trap: more input does not mean better decisions
- The intelligence operator posture: you define what matters, the system finds it, agents synthesize it
- The signal-to-noise problem: 99% of what the internet produces is irrelevant to your goals
- The pull vs push distinction: information that finds you (signal) vs information you go find (research)
- How One Mind OS handles both: automated feeds for pull, Trinity for push research
- The intelligence question: "What do I need to know to make better decisions in my domains?"

**Lesson 1.2 — The Intelligence Architecture**
The three-layer system: collect → synthesize → deliver.
- Layer 1 — Collection: RSS feeds, newsletters, GitHub watch, Twitter/X lists, Reddit monitoring
- Layer 2 — Synthesis: Trinity processes raw content, extracts signal, discards noise
- Layer 3 — Delivery: structured briefing in Discord, relevant items routed to vault
- The Trinity agent role: Domain 06 operator — R&D, trend monitoring, research execution
- The intelligence loop: signal → synthesis → briefing → action → knowledge capture
- What makes this better than manually reading: Trinity reads 100 sources so you read 10 highlights

**Lesson 1.3 — Defining Your Intelligence Domains**
You can't monitor everything. Define exactly what matters for your goals.
- The domain intelligence map: one intelligence question per active domain
- SO: What new AI tools, models, or infrastructure patterns should I know about?
- HP: What's new in longevity, performance, and cognitive enhancement research?
- LE: What's relevant to my family, home, and community?
- GE: What market signals, competitor moves, or opportunities exist in my ventures?
- Prioritizing domains: where does better intelligence have the highest leverage?
- The intelligence brief format: one paragraph per domain, most important signal first

**Lesson 1.4 — The R&D Backlog**
A prioritized list of questions worth deep investigation.
- What belongs in the R&D backlog: questions that require synthesis, not just Googling
- The R&D backlog format: question + why it matters + priority + estimated depth
- Populating the backlog: every time you want to understand something deeply, add it here
- Prioritization criteria: impact on decisions + time sensitivity + current relevance
- The backlog review: weekly, Trinity checks the backlog and flags high-priority items
- The research sprint: dedicating focused time to one deep research question
- Where the backlog lives: `00-24 SO/06 Research/backlog.md`

---

### MODULE 2: FEED ARCHITECTURE AND SIGNAL COLLECTION

**Lesson 2.1 — RSS Is Not Dead: Building Your Feed Stack**
RSS gives you a direct, algorithm-free feed from any source. This is how professionals monitor information.
- Why RSS over social media: no algorithm, no ads, no engagement manipulation — pure signal
- RSS reader options: Miniflux (self-hosted, open source), Feedly, Inoreader
- Miniflux on your VPS: deploy alongside your other services
- The feed inventory: one RSS URL per source you want to monitor
- Feed categories: AI/ML, infrastructure, your industry, longevity, finance, competitors
- Feed quality filtering: how to evaluate whether a source is worth monitoring
- Hands-on: set up Miniflux on Titan, subscribe to 20 high-signal feeds

**Lesson 2.2 — AI and Technology Feed Stack**
The exact feeds that keep you current on AI, agents, and infrastructure.
- Model releases: Anthropic, OpenAI, Google DeepMind, Mistral — official blogs
- Research papers: arXiv cs.AI, Hugging Face daily papers, Papers With Code
- Infrastructure: DigitalOcean blog, Tailscale blog, Docker Hub
- AI tools and frameworks: LangChain, OpenClaw changelog, clawhub announcements
- Community signals: Hacker News AI tag, r/LocalLLaMA, AI Twitter lists
- The curation principle: 20 high-quality sources beats 200 mediocre ones
- Hands-on: build your personal AI feed stack — add to Miniflux

**Lesson 2.3 — Domain-Specific Intelligence Sources**
Sources for each of your active life domains.
- GE / Business: competitor blogs, industry publications, earnings releases, ProductHunt
- HP / Personal: longevity researchers (Huberman, Attia, Rhonda Patrick), PubMed alerts, fitness publications
- SO / Security: Krebs on Security, CISA alerts, CVE feeds, security mailing lists
- LE / Home: local community feeds, property market signals, relevant newsletters
- The custom alert stack: Google Alerts, GitHub Watch, Twitter/X advanced search saved
- Newsletter vs RSS: when email newsletters are worth subscribing to
- Hands-on: add 5 domain-specific feeds per active domain

**Lesson 2.4 — The Miniflux + Trinity Integration**
Connecting your RSS reader to Trinity for automated synthesis.
- How Trinity reads from Miniflux: the Miniflux API endpoint
- The synthesis prompt: "Here are today's unread articles. Extract the 5 most important signals."
- Filtering by category: Trinity focuses on one category at a time
- The synthesis output format: headline + one-paragraph summary + why it matters
- Routing to Discord: Trinity posts the synthesis to `#intelligence` channel
- The read/unread loop: articles Trinity processes are marked read in Miniflux
- Hands-on: wire Trinity to Miniflux, run first synthesis

---

### MODULE 3: TRINITY — YOUR RESEARCH AGENT

**Lesson 3.1 — Trinity's Role in the Agent Stack**
Trinity is Domain 06 — the Lab. Its job is to investigate, synthesize, and brief.
- Trinity's IDENTITY: research analyst, trend monitor, knowledge synthesizer
- What Trinity does vs what Legacy does: Trinity goes deep, Legacy goes wide
- Trinity's domain: SO 06 Research / R&D — all research tasks route here
- The researcher's workflow: question → sources → synthesis → brief → vault capture
- Trinity vs Google: Trinity reads, synthesizes, and gives you a structured answer — Google gives you 10 blue links
- When to use Trinity: questions that require synthesis from multiple sources

**Lesson 3.2 — Research Protocols**
Structured methodology for different types of research tasks.
- **Quick brief** (5 min): "Trinity, give me a quick overview of [topic]" — summary + 3 key points
- **Deep dive** (30 min): "Trinity, investigate [topic] thoroughly — sources, current state, implications for my stack"
- **Comparative analysis** (20 min): "Trinity, compare [option A] vs [option B] for [use case]"
- **Literature synthesis** (45 min): "Trinity, synthesize recent research on [topic] into a brief I can act on"
- **Competitive analysis** (30 min): "Trinity, what are [competitor] doing that I should know about?"
- Choosing the right protocol: match depth to the importance of the decision
- The research brief template: structure every output the same way for easy consumption

**Lesson 3.3 — Staying Current on AI**
The AI landscape changes every week. Here's how to stay informed without spending hours on it.
- The weekly AI briefing: Trinity synthesizes the week's AI developments every Friday
- Model monitoring: Trinity watches for new model releases and benchmarks
- The "is this relevant to my stack?" filter: Trinity evaluates new tools against your current setup
- OpenClaw changelog monitoring: Trinity watches for updates that affect your agents
- The AI opportunity scan: monthly, Trinity looks for new tools that could upgrade your system
- Knowledge capture: important AI developments become permanent vault notes in SO 06
- Hands-on: set up the weekly AI briefing cron — Trinity to Discord every Friday at 8am

**Lesson 3.4 — Competitive Intelligence**
For your GE domains, knowing what competitors and the market are doing is leverage.
- Defining your competitive landscape: who are the key players in each venture?
- The competitive monitoring stack: RSS + Google Alerts + ProductHunt + LinkedIn (manual)
- The competitive brief: monthly, Trinity synthesizes what's changed in your competitive landscape
- Opportunity signals: new competitors entering, market shifts, pricing changes
- The competitive intelligence vault note: one note per key competitor, updated monthly
- What not to do: don't obsess over competitors — use intelligence for strategy, then focus on your work

**Lesson 3.5 — Research-to-Vault: Knowledge Capture**
Research is only valuable if it compounds. Every significant finding becomes a vault note.
- The research output format: findings → implications → action items → source links
- The knowledge note template: title, domain, date, summary, key points, implications, sources
- Routing research to the right domain: AI research → SO 01, longevity → HP 27, market research → GE 81
- The evergreen principle: write notes that will still be useful in 12 months
- The link graph: connecting new research to existing vault knowledge via backlinks
- The knowledge decay policy: research notes get a review date — Trinity flags stale research quarterly
- Hands-on: convert one research session into 3 permanent vault knowledge notes

---

### MODULE 4: THE INTELLIGENCE BRIEFING SYSTEM

**Lesson 4.1 — The Daily Intelligence Brief**
What Trinity delivers to your Discord every morning before you start work.
- The brief format: 5 items maximum, ranked by relevance to your active goals
- Each item: headline + one sentence summary + why it matters to you specifically
- Sources: Miniflux unread (last 24 hours) + any flagged research items
- The brief cron: Trinity runs at 6:30am (before your morning ritual at 7am)
- Brief length discipline: if it takes more than 3 minutes to read, it's too long
- The feedback loop: after 2 weeks, tell Trinity what types of items are most valuable — update the prompt
- Hands-on: set up the daily intelligence brief cron

**Lesson 4.2 — The Weekly Deep Brief**
One longer synthesis every week — 30 minutes of Trinity's attention on your behalf.
- The weekly brief scope: what changed in your 4–6 key intelligence domains this week?
- Format: one paragraph per domain, most important signal, one implication
- Delivered: Sunday evening so you can factor it into Monday planning
- The AI section: always in the weekly brief — what dropped in AI this week?
- The opportunity section: anything the weekly scan surfaced that's worth acting on
- Trinity's sources: Miniflux weekly digest + R&D backlog review + any flagged items
- The brief as strategy input: the weekly brief feeds directly into your weekly review

**Lesson 4.3 — Research Sprints**
Dedicated blocks of time for deep investigation of your highest-priority R&D questions.
- The research sprint format: 90 minutes, one question, Trinity + you in a focused session
- Sprint structure: question definition (5 min) → Trinity initial synthesis (20 min) → gaps identified (10 min) → deep dive (45 min) → vault capture (10 min)
- The sprint log: a vault record of every sprint, question, and findings
- Monthly sprint planning: review the R&D backlog, pick 2–4 questions to sprint on this month
- Sprint output: minimum one permanent vault note per sprint
- The compounding research library: 12 months of sprints = 96+ deep research notes in your vault

**Lesson 4.4 — Intelligence for Decision Making**
The ultimate purpose of intelligence: better decisions.
- The decision brief: when you have a major decision pending, Trinity prepares a research brief
- Decision types that benefit from intelligence: technology choices, business moves, personal investments, health protocols
- The decision research template: question + options + what I need to know + Trinity brief
- Time-boxing research: set a research deadline before the decision, not an open-ended investigation
- Avoiding analysis paralysis: intelligence informs decisions, it doesn't make them
- The post-decision review: 90 days after a decision, did the intelligence hold up?

---

### MODULE 5: EVALUATING AND ADOPTING NEW AI TOOLS

**Lesson 5.1 — The AI Tool Evaluation Framework**
New AI tools launch daily. This framework tells you what's worth your attention.
- The three evaluation questions: Does it do something I need? Does it do it better than what I have? Can I integrate it into my stack?
- The "shiny object" filter: most new tools solve problems you don't have
- Evaluation criteria: capability, integration path, cost, privacy, longevity (will it exist in 12 months?)
- The 30-minute evaluation: how to assess a new tool quickly before investing more time
- The trial protocol: 14-day trial, specific use case, measure vs current solution
- The adoption decision: only adopt if it replaces something or adds something genuinely missing

**Lesson 5.2 — Model Evaluation**
New models release constantly. How to know when to switch.
- The model benchmarks that matter: MMLU, HumanEval, LMSYS Chatbot Arena — what they test
- Task-specific evaluation: benchmark on YOUR actual use cases, not general benchmarks
- Cost/performance tradeoff: a 10% better model at 5× the cost may not be worth it
- The model swap protocol: test on 20 representative tasks → compare outputs → decide
- Latency and reliability: benchmark accuracy is useless if the model is slow or unreliable
- The current stack decision: when does Haiku vs Sonnet vs extended thinking — for which tasks
- Hands-on: run a 10-task evaluation comparing two models on your real use cases

**Lesson 5.3 — The Emerging Technology Radar**
A structured view of what's coming and what to watch.
- The Technology Radar concept (ThoughtWorks) adapted for One Mind OS
- Four categories: Adopt (use now), Trial (experiment), Assess (watch carefully), Hold (not yet)
- Your current radar: where does each major AI technology sit?
- Quarterly radar update: Trinity reviews and updates the radar based on developments
- The radar as strategy input: informs what you invest time learning vs what you ignore
- Publishing your radar: sharing your technology assessment with the community
- Hands-on: build your first personal technology radar with Trinity's help

**Lesson 5.4 — Learning New Technologies Efficiently**
When you decide to adopt a new technology, how to go from zero to operational fast.
- The 80/20 learning rule: 20% of the technology's features handle 80% of your use cases — learn those first
- The documentation-first approach: read the official docs once, deeply, before watching tutorials
- The context7 MCP server: up-to-date documentation for any library, directly in your agent context
- Learning-by-building: set a specific integration target, learn what you need to achieve it
- The knowledge capture habit: as you learn, write vault notes — you won't remember it otherwise
- The test-in-isolation principle: don't learn a new technology in production
- Hands-on: use the learning protocol to get operational with one tool from your "Trial" radar category

---

### MODULE 6: THE RESEARCH VAULT

**Lesson 6.1 — Domain 06 Architecture**
How to structure `00-24 SO/06 Research` for a working intelligence operation.
- The folder structure: `feeds/`, `briefs/`, `sprints/`, `radar/`, `backlog.md`, `README.md`
- `feeds/`: feed inventory, categories, quality scores
- `briefs/`: dated daily and weekly briefs (auto-generated by Trinity, archived)
- `sprints/`: one note per research sprint with question, findings, vault captures
- `radar/`: the technology radar, updated quarterly
- `backlog.md`: the prioritized R&D question queue
- The domain README: the index Trinity reads to orient itself in this domain

**Lesson 6.2 — The Knowledge Note Standard**
Research notes that age well and retrieve cleanly.
- Atomic vs compound notes: one idea per note vs comprehensive reference notes
- The knowledge note template for research: metadata, summary (3 sentences), key findings, implications, sources, related notes
- Writing for future retrieval: titles and first sentences that are searchable
- The date/version discipline: when something is updated, note the date — knowledge has a shelf life
- Cross-pillar knowledge: a longevity research note that affects both HP and GE
- The contradiction pattern: when new research contradicts old vault knowledge — update explicitly

**Lesson 6.3 — Maintaining the Intelligence Operation**
The system requires maintenance to stay relevant.
- Feed quality audit (monthly): remove low-signal feeds, add new high-signal sources
- The brief quality review (weekly): is Trinity surfacing the right items? Update the prompt if not.
- R&D backlog grooming (weekly): reprioritize based on what's moved in your life
- Knowledge note staleness audit (quarterly): Trinity flags research notes older than 6 months for review
- The intelligence metrics: brief items acted on per week, research sprints completed per month
- The ROI question: is the intelligence operation producing better decisions? Honest quarterly assessment

**Lesson 6.4 — Graduation: Your Intelligence Operation Is Live**
You've built a professional-grade intelligence function for your personal operating system.
- What you've built: feed stack, Trinity briefings, research protocols, knowledge vault, technology radar
- The full pipeline: signal → synthesis → brief → action → vault
- The compounding effect: 6 months from now, your vault has hundreds of research notes — Trinity uses them all
- Intelligence as competitive advantage: you make decisions with better information than everyone around you
- The intelligence handoff: when you spawn new domain agents, they inherit the intelligence architecture
- Graduation challenge: complete 3 research sprints this month using the full protocol

---

## Resources Included

| Resource | Format | Access |
|----------|--------|--------|
| Miniflux setup guide (VPS) | Written | Legacy+ |
| AI feed stack (20 curated RSS URLs) | OPML | Legacy+ |
| Trinity IDENTITY.md + SOUL.md (research persona) | Markdown | Legacy+ |
| Daily intelligence brief cron prompt | Markdown | Legacy+ |
| Weekly deep brief cron prompt | Markdown | Legacy+ |
| Research sprint template | Markdown | Legacy+ |
| R&D backlog template | Markdown | Legacy+ |
| Knowledge note template | Markdown | Legacy+ |
| Technology radar template (Obsidian) | Markdown | Legacy+ |
| AI tool evaluation framework | PDF | Legacy+ |
| Model evaluation worksheet | Spreadsheet | Legacy+ |
| Domain 06 folder structure | Vault template | Legacy+ |

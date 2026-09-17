# Course 3: CLAUDE CODE & AI TOOLS

> **Build on top of your stack. Write skills. Engineer context. Go multimodal.**

**Tier:** Legacy ($97/mo)
**Circle Space:** ⚙️ SO — Sovereign Operations (ID: 2617120)
**Prerequisite:** Course 1 (vault + first agent). Course 2 recommended but not required.
**Target Student:** Someone who wants to extend their system — write custom skills, build MCP servers, use Claude Code as a development partner, and go beyond text
**Outcome:** Custom OpenClaw skills built and deployed, Claude Code running in your workflow, MCP servers configured, multimodal inputs flowing into your agents, context engineered for maximum intelligence

---

## What You'll Build

By the end of this course:
- 2+ custom OpenClaw skills installed and working
- Claude Code set up and integrated into your daily workflow
- MCP servers configured for your stack (filesystem, GitHub, custom)
- A multimodal pipeline: voice, images, and video feeding into your agents
- Context engineering principles applied across all your agent files
- A CLAUDE.md / AGENTS.md that makes AI actually understand your project

---

## Course Structure — 6 Modules, 35 Lessons

---

### MODULE 1: CONTEXT ENGINEERING
*"Garbage in, garbage out. Context engineering is the highest-leverage skill in AI."*

**Lesson 1.1 — What Context Engineering Actually Is**
Prompt engineering is asking better questions. Context engineering is building the environment where the right answer is inevitable.
- The difference: prompt engineering (one message) vs context engineering (the whole system)
- Why this matters more than any other AI skill
- The 4 layers of context: identity, knowledge, instructions, state
- The One Mind OS advantage: your Codex IS your context — structured, persistent, AI-readable
- The cost of bad context: hallucination, drift, repeated mistakes, wasted tokens

**Lesson 1.2 — The Context Budget**
Every AI interaction has a token limit. How you spend it determines what your agent can do.
- What a context window is: the agent's working memory for one session
- Claude Sonnet 4.6: 200K token context window
- What eats tokens: system prompt, memory files, conversation history, tools schema
- Context budgeting: allocate tokens by priority (identity > goals > recent context > history)
- The `#always-loaded` tag: how to mark files that always load into context
- Practical exercise: audit your current agent's context — what's taking up space?

**Lesson 1.3 — Writing Agent Instruction Files That Actually Work**
SOUL.md, USER.md, CODEX.md, AGENTS.md — the files that define agent behavior. How to write them so AI reads them correctly.
- Specificity over vagueness: "respond in under 3 sentences for quick questions" vs "be concise"
- Structured markdown: headers, tables, and bullet points load more reliably than prose
- Behavioral contracts: "When I say X, you always do Y"
- Negative space: what you tell the agent NOT to do is as important as what to do
- The priority stack: which instructions win when they conflict?
- Hands-on: rewrite your current SOUL.md using context engineering principles

**Lesson 1.4 — The System Prompt Architecture**
What goes in the system prompt vs memory files vs runtime messages.
- System prompt: identity + core instructions (loaded every session)
- Memory files: knowledge that changes over time (loaded on demand)
- Runtime messages: dynamic context (today's date, current task, recent events)
- The pinned context pattern: `_codex/system/` files that always load
- Token-efficient identity: how to say everything about yourself in 500 tokens
- The "AI boot file" concept: what compiles before every agent response

**Lesson 1.5 — Context Drift and How to Fix It**
Agents get confused over long conversations. This lesson teaches you to recognize and prevent it.
- What drift looks like: agent forgets instructions, changes personality, ignores guardrails
- Causes: context window filling up, conflicting instructions, missing anchors
- Anchoring techniques: periodic identity reminders, structured checkpoints
- The conversation reset pattern: when to start fresh vs continue
- Compaction as context management: summarizing before the window fills

**Lesson 1.6 — Frontmatter as Context Interface**
YAML frontmatter is not metadata — it's the interface between your files and your agents.
- What agents do with frontmatter: filter, route, prioritize, act
- The standard One Mind frontmatter spec: `id`, `type`, `status`, `tags`, `owner`, `domain`, `agent`
- The `ai_summary` field: pre-computed summaries that save tokens
- The `autonomy` field: L0–L3 directly in the file
- Hands-on: audit 10 files in your vault — upgrade their frontmatter

---

### MODULE 2: CLAUDE CODE
*"Your AI pair programmer. It reads your codebase, writes code, runs commands, and commits — while you direct."*

**Lesson 2.1 — What Claude Code Is (And Isn't)**
Claude Code is not GitHub Copilot. It's an agentic coding tool that operates in your terminal with read/write access to your entire project.
- Claude Code vs Copilot vs Cursor vs Windsurf — honest comparison
- The agentic difference: Claude Code can read files, run commands, write tests, commit code
- When to use Claude Code: complex refactors, new features, debugging, building skills
- When NOT to use it: quick autocomplete (use Copilot), simple edits (just write them)
- The `--dangerously-skip-permissions` flag: what it does, when it's appropriate, the tradeoffs

**Lesson 2.2 — Install and Configure Claude Code**
Step-by-step setup. Running in 10 minutes.
- Install: `npm install -g @anthropic-ai/claude-code`
- Authentication: Anthropic API key or GitHub Copilot (which you already have)
- First run in your vault: `cd ~/Documents/onemind && claude`
- The interactive loop: ask, review, approve, repeat
- VS Code integration: using Claude Code from the terminal panel
- The conversation history: Claude Code remembers your project across sessions

**Lesson 2.3 — CLAUDE.md: Teaching Claude Your Project**
CLAUDE.md is to Claude Code what SOUL.md is to OpenClaw agents. It's the project instruction file.
- What CLAUDE.md contains: project overview, conventions, commands, what NOT to do
- File location: root of your project (or vault)
- The One Mind CLAUDE.md template: download and customize
- How Claude Code loads it: automatically on every session
- Writing for agentic behavior: "When adding a new domain, always create README.md first"
- Hands-on: write CLAUDE.md for your vault and for one code project

**Lesson 2.4 — The Agent Loop in Practice**
How to work with Claude Code effectively — the rhythm that gets things done.
- Plan → Edit → Verify → Commit: the 4-step loop
- Giving good task descriptions: outcome + constraints + context
- Reviewing diffs before approving: what to look for
- Iterating: "that's close but change X" vs starting over
- The multi-step task: breaking complex work into Claude Code-sized chunks
- Hands-on: use Claude Code to build something in your vault (create a new domain folder, add templates, write README)

**Lesson 2.5 — Claude Code for Building OpenClaw Skills**
The most powerful use case: using Claude Code to build and test skills for your agents.
- Scaffold a new skill: `claude "create an OpenClaw skill for Todoist API"`
- Iterative development: run, test, fix, repeat — with Claude Code doing the typing
- The skill test harness: how to verify a skill works before deploying
- Publishing to your local skills library
- Hands-on: build one custom skill using Claude Code

**Lesson 2.6 — GitHub Copilot in VS Code**
Your always-on inline AI. Different tool, different job.
- Copilot Chat vs Claude Code: when to use each
- Custom instructions: `.github/copilot-instructions.md`
- Inline completions: how to guide them effectively
- The agent mode in VS Code: Copilot doing multi-step tasks inside the editor
- Skills and instructions files: building domain knowledge into Copilot
- Hands-on: configure Copilot custom instructions for your One Mind project

---

### MODULE 3: SKILLS FRAMEWORK
*"Skills are what your agents can DO. Build them, install them, own them."*

**Lesson 3.1 — What a Skill Actually Is**
An OpenClaw skill is a Node.js package that exports tools. It's the bridge between your agent's language and the real world.
- The anatomy of a skill: `index.js`, `package.json`, `SKILL.md`, `README.md`
- What skills enable: API calls, file operations, database queries, web requests, anything
- Skills vs MCP servers: when to use which
- The clawhub marketplace: browsing and installing community skills
- Skills you should install on day one: todoist, github, home-assistant, weather

**Lesson 3.2 — Installing and Managing Skills**
The full skill lifecycle: find → install → configure → use → update.
- `clawhub search <query>`: find skills by capability
- `clawhub install <skill-name>`: install to your workspace
- Skill configuration: API keys, endpoints, preferences
- Skill profiles: `full` (all tools) vs `messaging` (read-only) vs custom allow-lists
- Skill versioning: pinning vs auto-update
- Hands-on: install 3 skills and test each one with your agent

**Lesson 3.3 — Anatomy of a Skill (Deep Dive)**
Read a real skill's source code. Understand how it works.
- The `tools` export: an array of tool definitions
- Tool definition structure: name, description, input_schema, execute function
- The execute function: receives parameters, calls an API, returns a result
- Error handling: what happens when an API call fails
- The SKILL.md: what to include so the agent uses the tool correctly
- Walk-through: read the Todoist skill source code end-to-end

**Lesson 3.4 — Build Your First Custom Skill**
Write a skill from scratch. Deploy it. Watch your agent use it.
- Scaffold: `clawhub create my-skill-name`
- Choose your target: a SaaS API you use (Notion, Airtable, your homestead sensor, anything)
- Write the tool definition: name, description, input schema
- Write the execute function: fetch, transform, return
- Install locally: `npm link` or workspace path
- Test: ask your agent to use the skill
- Hands-on: build a skill for one thing in your life

**Lesson 3.5 — The Skill Design Principles**
What makes a skill that agents actually use correctly.
- Tool names: specific and verb-first (`get_todoist_tasks` not `todoist`)
- Descriptions: written for AI, not humans — explain when to use it, not just what it does
- Input schemas: strict typing + clear descriptions for every parameter
- Return values: structured JSON that agents can reason about
- The single responsibility principle: one skill, one job
- Idempotency: why skills should be safe to call multiple times

**Lesson 3.6 — Publishing to clawhub**
Share your skill with the One Mind community and the broader OpenClaw ecosystem.
- `clawhub publish`: the publish workflow
- Versioning with semver
- Writing a great README and SKILL.md
- The community value: your skill → 200,000+ OpenClaw users
- Skill monetization: how creators earn on clawhub (preview of Course 4 business track)

---

### MODULE 4: MCP AND TOOLS
*"MCP is the universal protocol that connects AI to everything. Master it and your agents can reach anything."*

**Lesson 4.1 — What MCP Is (Model Context Protocol)**
MCP is Anthropic's open standard for connecting AI models to external tools and data sources.
- The problem MCP solves: every AI tool had its own custom integration format
- MCP as a universal interface: write once, works with Claude, GPT, Gemini, any MCP-compatible model
- Server vs client: MCP servers expose tools, clients (your agents) call them
- The 3 things an MCP server can expose: tools, resources, prompts
- Why MCP matters for One Mind OS: your vault, your databases, your APIs — all accessible via MCP

**Lesson 4.2 — MCP Servers Already in Your Stack**
What you already have configured — and how to verify it works.
- `@modelcontextprotocol/server-filesystem`: read/write your vault
- `@modelcontextprotocol/server-github`: GitHub repos, issues, PRs
- `@upstash/context7-mcp`: up-to-date library documentation
- Configuration in `openclaw.json`: `mcpServers` block
- Testing each server: ask your agent to use each one
- The MCP tool call: what it looks like in logs, how to trace it

**Lesson 4.3 — The Filesystem MCP (Your Vault Interface)**
The most important MCP server: the one that gives agents eyes and hands on your Codex.
- `read_file`, `write_file`, `list_directory`, `search_files` — the 4 tools you'll use most
- Scoping access: which paths each agent can reach
- The security model: `allowed_directories` in server config
- Agent write workflows: how agents create notes, update frontmatter, log activity
- The auto-commit loop: agent writes → onemind-watch picks it up → GitHub updated in 3 seconds
- Hands-on: ask your agent to create a new capture note in your Inbox

**Lesson 4.4 — Build a Custom MCP Server**
If you have data an agent needs to reach, and there's no existing MCP server for it — build one.
- The MCP SDK: `@modelcontextprotocol/sdk`
- Scaffold a server: `npx @modelcontextprotocol/create-server my-server`
- Define tools: name, description, input schema, handler function
- Define resources: static or dynamic data sources agents can read
- Test with Claude Desktop or OpenClaw
- Hands-on: build an MCP server for something in your stack (your n8n instance, a local database, a custom API)

**Lesson 4.5 — Tools Schema and Tool Calling**
How agents decide which tool to call — and how to make them choose correctly.
- The tool selection problem: agents pick tools based on descriptions
- Writing tool descriptions that get used correctly: specificity, use-cases, when NOT to use
- Input schema design: required vs optional, types, constraints, examples
- The tool call flow: agent reasons → selects tool → sends parameters → receives result → continues
- Debugging tool calls: enable verbose logging, inspect the raw tool call in logs
- The `allowedTools` filter: restricting which tools an agent can use per context

---

### MODULE 5: MULTIMODAL SUPPORT
*"Text is one channel. Your agents can now see, hear, and process anything."*

**Lesson 5.1 — The Multimodal Landscape**
What your agents can process today — and what's coming.
- Text: the foundation (where we've been)
- Images: screenshots, diagrams, photos, documents — Claude can see them
- Voice: speech-to-text in, text-to-speech out
- Video: frame extraction, transcription, analysis
- The One Mind OS multimodal vision: every input type flows into the right agent

**Lesson 5.2 — Voice Input and Output**
Talk to your agents. Hear them respond. The most natural interface for life management.
- Voice-to-text: Whisper API integration, Apple Dictation, system-level STT
- The voice capture workflow: speak → transcribe → agent processes → responds
- Text-to-speech output: agent responses read aloud
- OpenClaw Talk Mode: the voice interface built into the gateway
- Wake word configuration: say "Legacy" → agent activates
- Mobile voice: iPhone Shortcuts → OpenClaw webhook → agent responds via iMessage
- Hands-on: set up voice input for your primary agent

**Lesson 5.3 — Image Processing**
Your agents can analyze images. This lesson builds the pipeline.
- What Claude can do with images: describe, extract text, analyze diagrams, compare
- Sending images to agents: Discord attachment → agent analyzes
- Use cases in One Mind OS:
  - Screenshot a receipt → Oracle categorizes the expense
  - Photo of whiteboard → Legacy creates structured notes
  - Screenshot of error → Grid diagnoses the issue
  - Photo of plant → Eden identifies and advises
- The image pipeline: capture → base64 or URL → agent call → structured output
- Hands-on: set up at least one image-to-agent workflow

**Lesson 5.4 — Document Processing**
PDFs, spreadsheets, Word docs — getting structured intelligence out of unstructured files.
- The text-extractor plugin (already in your Obsidian stack): extracts text from PDFs
- Claude for document analysis: summarize, extract tables, answer questions
- The document workflow: upload → extract → agent processes → writes to vault
- Use cases: contracts to Guardian, financial statements to Oracle, research papers to Trinity
- Hands-on: process one document through an agent and capture the output to your vault

**Lesson 5.5 — Video Processing**
Videos are the hardest format — this lesson shows the practical approach.
- What's possible today: transcript extraction, key frame analysis, summary generation
- Tools: Whisper for transcription, Claude for analysis
- The video-to-knowledge workflow: YouTube video → transcript → agent summarizes → vault knowledge note
- Use cases: meeting recordings, course videos, client calls
- The n8n video pipeline: trigger → download → transcribe → analyze → store
- Limits: what AI can't reliably do with video yet (be realistic)

---

### MODULE 6: MODELS AND INTEGRATIONS
*"Not every task needs the same model. This module teaches you to route intelligently."*

**Lesson 6.1 — The Model Landscape (May 2026)**
Honest overview of what's available and what each model is good at.
- Claude family: Sonnet (balanced), Haiku (fast/cheap), Opus (deep reasoning)
- GPT family: GPT-4o (multimodal), GPT-4o-mini (cheap), o3 (reasoning)
- Gemini family: 2.0 Flash (speed), 2.5 Pro (long context)
- Open source: Llama 3, Mistral, Qwen — when self-hosting makes sense
- The One Mind OS model routing strategy: match model to job, optimize cost

**Lesson 6.2 — Model Interactions and Routing**
Different agents, different models. This lesson teaches intelligent routing.
- The principle: use the most capable model only when you need it
- Orchestrator (Legacy): needs full reasoning — Claude Sonnet or better
- Domain agents: Haiku for most tasks, Sonnet for complex analysis
- Subagents and workers: cheapest model that can do the job
- Routing by task type: classification → Haiku, analysis → Sonnet, code → Sonnet, creative → Sonnet
- Per-request model override: how to escalate dynamically
- Cost tracking: what each model costs per 1M tokens, your estimated monthly bill

**Lesson 6.3 — GitHub Copilot API as Model Provider**
You already have GitHub Copilot. Use it as your model provider — zero extra cost.
- The `github-copilot/` model prefix in OpenClaw
- Available models: `claude-sonnet-4.6`, `gpt-4o`, `o3-mini`
- How authentication works: your existing Copilot subscription
- Rate limits: what to expect, how to handle them
- Fallback: when to add a dedicated Anthropic or OpenAI key

**Lesson 6.4 — Adding New Model Providers**
Expand beyond the default. Add Anthropic direct, OpenAI, or local models.
- Anthropic API: `openclaw models add anthropic --key YOUR_KEY`
- OpenAI API: `openclaw models add openai --key YOUR_KEY`
- Ollama (local): self-hosted models on your VPS — when it makes sense
- Model aliases: shorthand names for your most-used models
- The fallback chain: `primary → fallback → emergency`

**Lesson 6.5 — Course 3 Graduation Challenge**
Build something real.
1. A custom OpenClaw skill deployed and working
2. Claude Code used for at least one real task in your project
3. At least one MCP server configured and tested
4. One multimodal pipeline working (voice OR image OR document)
5. Post a 60-second demo video to the community

---

## Resources Included

| Resource | Format | Access |
|----------|--------|--------|
| CLAUDE.md template (vault + code project) | Download | Legacy+ |
| Custom skill scaffold (`clawhub create` starter) | Download | Legacy+ |
| MCP server scaffold (custom server starter) | Download | Legacy+ |
| Context engineering worksheet | PDF | Legacy+ |
| Model routing decision matrix | PDF | Legacy+ |
| Voice pipeline setup guide | Written | Legacy+ |
| "Why Git-Based PKM is Superior" | Guide | Legacy+ |

---

## Up Next → Course 4: AI AGENT ARCHITECTURE

You can build and extend the system. Course 4 teaches you to design it at scale — evals, guardrails, human-in-the-loop workflows, and reasoning engines that think before they act.

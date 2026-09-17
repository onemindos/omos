# 10 HERMES AGENT HACKS THAT TURNED MY CHAT AGENT INTO A 24/7 SYSTEM

> **Title:** 10 HERMES AGENT HACKS THAT TURNED MY CHAT AGENT INTO A 24/7 SYSTEM
> **Author:** YanXbt ([@IBuzovskyi](https://x.com/IBuzovskyi))
> **Source:** [x.com/IBuzovskyi/status/2062101068842975409](https://x.com/IBuzovskyi/status/2062101068842975409?s=12)
> **Type:** X Article | **Posted:** 2026-06-03 | **Ingested:** 2026-06-07

---

These 10 Hermes Agent hacks saved me 15+ hours every week - and they work for any workflow you run repeatedly. Content, software development, business operations, client management, research, sales. If you do it more than once, Hermes can run it.
My examples are from content and social media automation because that's what I run daily. The setups themselves are domain-agnostic. A cron job that scans X for trending posts uses the same mechanic as one that checks GitHub for open PRs or monitors a CRM for new leads.
Most people use Hermes Agent like a chat app. They open it, type a prompt, get a response, close it.
That leaves 90% of what Hermes can do on the table.
Here are 10 setups that turn Hermes from a chat window into a 24/7 system that works while you sleep, reacts when your workflow changes, and gets sharper with every run.
If you only have time for 3, start with: Cron Jobs (#3), /goal Structure (#4), and Skills (#8). These three alone change how Hermes feels overnight.
Table of Contents
Mission Control — setup: 30 min, saves: 2 hrs/week
Event Triggers — setup: 20 min, saves: 3 hrs/week
Cron Jobs — setup: 10 min, saves: 5 hrs/week ⭐
/goal Structure — setup: 5 min, saves: 4 hrs/week ⭐
Sub-Agents — setup: 5 min, saves: 3 hrs/week
Telegram Workspaces — setup: 10 min, saves: 1 hr/week
Kanban Board — setup: 5 min, saves: 2 hrs/week
Skills as SOPs — setup: 15 min per skill, saves: 5 hrs/week ⭐
Webhooks — setup: 30 min, saves: 3 hrs/week
Separate Agents — setup: 20 min per profile, saves: 4 hrs/week
1. MISSION CONTROL
The first setup and the biggest one: build a dashboard where everything is visible.
When Hermes is doing real work, you don't want that work buried inside a chat thread. You want to see what is running, what is waiting on you, what is blocked, what needs approval, and what changed since yesterday.
Ask Hermes to build it:
bash
Build me a mission control dashboard. Start with:
- A kanban board showing all active agent tasks
- A content pipeline where I can add ideas and track progress
- A memory wiki showing everything we've worked on
- A performance section showing my X and content metrics
Hermes also has a built-in dashboard out of the box:
hermes dashboard
Opens at localhost:9119. Skills, models, cron jobs, profiles, kanban board. Start here, then customize when you need more.
Hermes also just launched a native Desktop app for macOS, Windows, and Linux. Side-by-side preview, file browser, integrated voice. Same data directory as CLI and Telegram. Work from Desktop at your machine, switch to Telegram on the go. One agent, every surface.
Download: hermes-agent.nousresearch.com/desktop
Fastest path to a working agent on any surface:
hermes setup --portal
One OAuth covers the model + web search + image generation + TTS + cloud browser. No separate API keys needed.
Once you add a Kanban board directly into the dashboard, Hermes stops being something you message and starts becoming part of your operating layer. You can see every task, every status, every blocked item in one place instead of digging through chat history.
Mission control shows you the system. The next setup is where Hermes starts waking up without you prompting it.
2. EVENT TRIGGERS
Think about where you already work. Notion. Linear. Google Sheets. Slack. You move a task, update a status, add an idea. Right now, nothing happens after that. You have to remember to tell Hermes about it later.
The fix: make Hermes watch for changes and react automatically.
Example workflow: when you move a video idea to your "To Film" list in Notion, Hermes detects the change and sends you a filming brief on Telegram within minutes.
That brief includes:
Should you film this now, later, or kill it?
The strongest title angle
A 30-second hook for the opening
Proof assets you'll need
A pre-filming checklist
You didn't prompt Hermes. You moved a card. Hermes saw the state change and did the work.
How to set this up:
Option A — Hermes cron job watches for changes (simplest):
Schedule a cron job every 10 minutes:
check my Notion board [board URL].
if any card moved to "To Film" in the last 10 minutes,
research the topic, write a filming brief,
send it to Telegram.
Option B — Webhook trigger (instant):Use Notion automations, Make, or Zapier to send a webhook to Hermes when a card moves. The response is instant instead of polling every 10 minutes.
The principle: when your workflow changes state, Hermes should know what to do next.
3. CRON JOBS
Event triggers react to changes. Cron jobs react to time.
Every morning I get useful information before I ask for it. That shift made Hermes feel like an employee who starts work before I wake up.
Cron jobs that run in my setup:
Every morning at 8am:
send me one AI story worth reacting to on X.

Every 3 hours:
scan X for fresh posts in my niche I should quote tweet.

Every day at 9pm:
check if competitors posted any outlier content today.

Every Monday at 9am:
audit my content board. flag ideas stuck for more than 7 days.

Every Friday at 6pm:
summarize what content shipped this week,
what performed, what didn't, and why.
Setting these up in Hermes is plain English. No crontab syntax. Just tell the agent what you want and when.
The difference: instead of you remembering to ask Hermes things, useful information arrives before you even think about it. That is what makes it feel like a 24/7 assistant.
Scheduled work is half the system. The other half is giving Hermes better objectives than normal prompts.
4. /GOAL WITH STRUCTURE
A normal prompt asks Hermes for one response. /goal gives Hermes an objective to work toward across multiple turns until it's done.
Most people use /goal like a prompt. They type something vague and get vague results back. The difference between a useless /goal and one that ships real work is structure.
The template that works:
/goal [OUTCOME]
using [SOURCES]
with constraints: [CONSTRAINTS]
deliverable: [DELIVERABLE]
Example:
/goal [OUTCOME]
using [SOURCES]
with constraints: [CONSTRAINTS]
deliverable: [DELIVERABLE]
Every part does a job:
Outcome tells Hermes when the goal is achieved
Sources tells it where to look
Constraints tell it what to avoid
Deliverable tells it what "done" looks like
Fuzzy prompts create fuzzy output. Structured goals create useful work.
The interview hack:
If you don't know how to structure your goal, make Hermes do it for you:
I want to use /goal but I don't want a vague goal.
Interview me with only the questions you need.
Then turn my answers into the strongest possible
/goal command. Include the exact outcome, context,
sources, constraints, deliverable,
and when you should stop.
Hermes asks you 5-8 questions, then writes its own /goal command from your answers. The resulting goal is sharper than anything you'd write from scratch.
5. SUB-AGENTS AS A RESEARCH TEAM
One agent gives you an answer. Sub-agents give you a team.
For any research task worth doing, I split it across multiple sub-agents running in parallel. Each one has a different source. The results merge into one recommendation.
Example:
/goal research the best content angle for this week.
spawn 3 sub-agents:

1. scan X for trending posts in AI agents niche,
   pull engagement numbers and hooks that worked

2. analyze my last 30 days of posts,
   find patterns in what performed vs what didn't

3. check competitor accounts,
   flag any outlier content from the last 7 days

combine all three into one recommendation
with the strongest angle, a draft hook,
and proof assets I'll need.
Each sub-agent gets its own context window. Only the final summary returns to the main session. Your main context stays light.
This is the difference between asking one person to research everything (slow, shallow) and having three specialists each go deep on their part (fast, thorough).
Best use cases for sub-agents:
Research across multiple sources
Competitive analysis (one sub-agent per competitor)
Content creation (one for research, one for drafting, one for editing)
Code reviews (one for logic, one for security, one for performance)
6. TELEGRAM TOPICS AS WORKSPACES
Telegram topics turn one chat into separate workspaces. Each topic becomes a different context with a different job.
My setup:
YouTube topic — content planning, scripts, filming briefs
React topic — trending posts on X worth reacting to
Coding topic — technical work, debugging, PRs
Research topic — deep dives, competitor analysis
General topic — smaller tasks, random questions
Different work needs different rooms. When everything runs in one chat, context bleeds. A coding question gets mixed with a content brief. A research task interrupts a draft review.
Topics fix that. Each workspace keeps its own conversation thread. Hermes knows what you're talking about based on which topic you're in.
Set up topics in any Telegram group:
Create a group with your Hermes bot
Enable Topics in group settings
Create a topic per workspace
Start messaging Hermes in each topic separately
Cron jobs per topic:
YouTube topic cron, every morning at 8am:
research trending AI content topics on X and YouTube.
send me 3 content ideas with engagement data.
post results in this topic only.

React topic cron, every 3 hours:
scan X for posts in AI agents niche
with 500+ likes in the last 3 hours.
if any are worth reacting to, draft a quote tweet
and send it here for approval.

Research topic cron, every Monday at 9am:
run a competitor audit.
check what @competitor1 @competitor2 @competitor3
posted last week. flag outliers. send report here.
Each topic gets its own cron job. Research stays in Research. Content stays in Content. No cross-contamination.
7. KANBAN FOR TASK MANAGEMENT
Once Hermes works on more than one thing, you need a board. Otherwise tasks disappear into chat and you forget what you even asked for.
Hermes has a built-in Kanban board. Durable SQLite storage. Shared across all profiles.
hermes kanban list
Shows your board. Drop tasks into triage. The dispatcher auto-assigns them to workers every 60 seconds.
Statuses: Triage → To-Do → Ready → Running → Blocked → Done
What this changes:
You see what is ready, running, and done
You see which agent owns which task
You see what is blocked and why
Crashed tasks get auto-reclaimed (zombie detection)
Heartbeats track worker health
Every /goal you set also becomes a Kanban card automatically. Multiple goals across multiple profiles, all visible on one board.
/goal research competitors → kanban card
/goal draft weekly report → kanban card
/goal triage inbox → kanban card
All three run in parallel. All three tracked in one place.
Kanban is how agent work stops disappearing into the chat.
Morning workflow with Kanban:
/goal research competitors → kanban card
/goal draft weekly report → kanban card
/goal triage inbox → kanban card
Drop 5 tasks at breakfast. By lunch, half are done. You didn't manage any of them.
8. SKILLS AS SOPs
A skill is a standard operating procedure for Hermes. You encode a process once and the agent uses it forever.
Hermes already creates skills on its own after every task. It reviews what worked, saves the workflow as a markdown file in ~/.hermes/skills/, and reuses it next time.
Writing skills intentionally for your key workflows is where the leverage multiplies.
Example — a content creation skill:
Save this as a skill called "content-post":

# Content Post Workflow

1. Check trending topics in AI agents niche via X search
2. Cross-reference with my last 14 days of posts (avoid repeats)
3. Pick the strongest angle based on engagement patterns
4. Write a draft in my voice:
   - ALL CAPS hook with HERMES AGENT as subject
   - arrows → for feature lists
   - No em-dashes, no adverbs, no throat-clearing
   - End with "full setup guide in the article 👇"
5. Score the draft:
   - Hook: does it stop the scroll? (1-10)
   - Bookmark fuel: would someone save this? (1-10)
   - Proof: is every claim backed by a number? (1-10)
6. If any score below 7, rewrite that section
7. Send final draft to Telegram for approval
Now whenever you say "use content-post for today's draft", Hermes runs the entire SOP without you explaining the process again.
Any workflow you explain twice should become a skill. Over time your agent accumulates dozens of skills, each one encoding a process that used to take manual instructions.
Skills are transparent. They live as markdown files. You can read them, edit them, delete them. No black box.
To see all skills:
hermes skills
Or open the dashboard → Skills tab:
hermes dashboard
Hermes ships with 60+ built-in tools across terminal, web, browser, vision, image generation, TTS, and code execution. Skills layer on top of those tools to create full workflows.
9. WEBHOOKS AND EVENT-BASED AGENTS
Cron jobs run because the clock changed. Webhooks run because the world changed.
A few examples of event-based triggers:
A new lead comes in → Hermes researches the company immediately
A GitHub PR opens → Hermes summarizes the changes and flags risks
A competitor posts content → Hermes checks if it's worth reacting to
A meeting transcript drops → Hermes extracts action items and adds tasks to your board
A keyword starts trending → Hermes drafts a content angle
The point is removing dead zones where something important happens and nobody acts on it until later.
Hermes receives webhooks through its gateway. Configure the webhook URL in your automation tool (Make, Zapier, n8n) and point it at your Hermes gateway endpoint.
Example: competitor monitoring via n8n + Hermes:
n8n workflow:
1. RSS trigger watches competitor blog (every 30 min)
2. if new post detected → send webhook to Hermes

Hermes /goal on webhook receive:
/goal a competitor just published: [title] [url].
read the full article via web search.
summarize the key points in 3 lines.
assess: should I react to this on X?
if yes, draft a reaction post in my voice.
send everything to Telegram for approval.
Example: new lead notification:
Zapier trigger: new form submission on website
→ webhook to Hermes

Hermes /goal on webhook receive:
/goal new lead: [name], [company], [email].
research their company via web search.
find their LinkedIn and recent posts.
draft a personalized follow-up email.
send draft to Telegram for approval.
The principle: cron jobs handle time. Webhooks handle events. Together they cover every scenario where Hermes should wake up without you touching it.
10. SEPARATE AGENTS BY JOB
You don't want your accountant to be your dentist. You don't want one agent doing every job with the same model, same tools, same memory, and same permissions.
Hermes profiles let you create separate agents for separate roles. Each profile has its own:
soul.md (personality and rules)
memory (what it knows)
skills (what it can do)
model (how smart it needs to be)
MCP connections (what tools it has access to)
permissions (what it's allowed to do)
Example setup:
hermes profile create content-lead
→ soul.md: you produce content. match my voice.
   use trending data. avoid repeated angles.
→ model: claude-sonnet-4 (needs strong writing)
→ tools: X search, web search, analytics

hermes profile create researcher
→ soul.md: you find information. deep research only.
   no opinions. facts and numbers.
→ model: gpt-5.5 (cheaper, high volume)
→ tools: web search, firecrawl, browser-use

hermes profile create ops
→ soul.md: you handle admin. calendar, email triage,
   reminders. ask for approval before sending anything.
→ model: gpt-5.5 (routine tasks, cost-efficient)
→ tools: email, calendar, notion

hermes profile create code-reviewer
→ soul.md: you review PRs. flag security issues,
   logic errors, performance problems.
→ model: claude-opus-4-8 (needs deep reasoning)
→ tools: github, terminal
Some agents need the smartest model you can afford. Some just need to check a page every hour. Some should have write access. Some should never have write access.
First /goal for each profile:
content-lead:
/goal research trending topics in AI agents niche.
cross-reference with my last 14 days of posts.
draft 2 posts in my voice.
score each on hook strength (1-10) and bookmark value (1-10).
send drafts to Telegram for approval.

researcher:
/goal run a deep research session on [topic].
use web search and firecrawl to pull data from
at least 5 sources. compile findings into
a structured markdown report with numbers,
sources, and key takeaways. save to file.

ops:
/goal check my inbox. summarize what needs a response.
pull today's calendar events.
flag anything urgent.
send morning briefing to Telegram by 9am.

code-reviewer:
/goal review the latest PR on [repo].
check for security issues, logic errors,
and performance problems.
write a review summary with specific line references.
send to Telegram.
Each profile runs its first /goal, learns from the result, and saves the workflow as a skill. Second run is faster. Fifth run is automatic.
Share any profile with one command:
Built a research agent that works? Share it via git:
cd ~/.hermes/profiles/researcher
git init && git add . && git commit -m "initial"
git push origin main
Anyone can install your agent:
hermes profile install github.com/you/researcher
They fill in their own API keys, and the agent runs with your skills, your soul.md, your workflows. Their memories and sessions stay separate.
The unlock is a team of smaller agents with the right brain, tools, memory, and permissions for their specific job.
Run them simultaneously. Each specialised. Each improving independently in their own lane.
HOW THEY CHAIN TOGETHER
These 10 setups compound when you stack them. Here is one chain running in my system right now:
8:00 AM — cron job (#3) fires.

the content-lead profile (#10) wakes up
and starts a /goal (#4) with structure:

"find the 3 strongest content angles for today
using X trending data and my last 14 days of posts."

it spawns 3 sub-agents (#5):
→ sub-agent 1 scans X for trending posts
→ sub-agent 2 pulls my recent post performance
→ sub-agent 3 checks competitor accounts

all three become kanban cards (#7).
dispatcher tracks them in parallel.

sub-agents finish. content-lead runs
the content-post skill (#8) to draft 2 posts.

drafts land in my Content topic
in Telegram (#6) for approval.

I tap approve on one. reject the other.
the approved post publishes via xurl.

10 minutes later a competitor publishes
a reaction to the same topic.
a webhook (#9) fires.
Hermes drafts a follow-up angle
and sends it to my React topic (#6).

I see everything on mission control (#1).
what ran, what shipped, what's pending.
One morning. Seven hacks fired. Two posts ready. Zero manual research.
That is the system.
THE REAL INSIGHT
If Hermes still feels like another chat app, look at the system around it.
Give it a mission control so you can see what's happening. Set up event triggers so it reacts when your workflow changes. Set up cron jobs so useful information arrives before you ask. Use /goal with structure instead of vague prompts. Split research across sub-agents. Separate workspaces with Telegram topics. Track tasks on the Kanban board. Turn repeatable processes into skills. Connect outside events via webhooks. Stop making one agent do every job.
Ten setups. Each one saves hours per week. Stack all ten and Hermes runs your operations while you focus on the work that moves the needle.
The agent is ready. The stack is ready. Wire the system and let it work.
→ Follow for AI agent blueprints every week → Save this. The /goal template and skill examples are worth bookmarking. → Comment "SYSTEM" and I'll send you the full skill templates for each setup.
Want 7 more hacks? The extended version with Session Recall, Token Optimization, Voice Mode, Browser Automation, Overnight /goal Runs, Memory Providers, and Tool Search is on Substack: Read the full 17-hack guide on Substack →
Resources
Official:
Hermes Agent Docs — installation, configuration, full CLI reference
Skills Hub — community skills, browse and install
GitHub — source, issues, PRs
My articles on Hermes:
HERMES AGENT: THE COMPLETE GUIDE — installation, models, dashboard, use cases, security
The Complete Hermes /goal Playbook — 21 Workflows
Hermes /goal — The Full Guide
How to Make Hermes + xurl Actually Work as a System
Hermes x Bitwarden — The Security Stack

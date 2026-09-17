# 10 Real Hermes Agents Settings That Actually Matter

> **Title:** 10 Real Hermes Agents Settings That Actually Matter
> **Author:** Tony Simons ([@tonysimons_](https://x.com/tonysimons_))
> **Source:** [x.com/tonysimons_/status/2063075301035225592](https://x.com/tonysimons_/status/2063075301035225592?s=12)
> **Type:** X Article | **Posted:** 2026-06-06 | **Ingested:** 2026-06-07

---

Yesterday, @CyrilXBT sold 10 fake Hermes env vars in an X "article" and the algorithm pushed the slop out to over 110k viewers as of the time I'm writing this. 
Tony Simons
@tonysimons_
"Ten settings that transform Hermes" is a hell of a headline. 

Problem is, every single env var in this thread is made up. 

I just ran it past my actual Hermes agent against the official docs and source code: zero hits.
CyrilXBT
@cyrilXBT
The 10 Hermes Agent Settings Most Users Never Find That Turn a Chat Agent Into a 24/7 Operation
Most people who install Hermes Agent use about 20% of what it can do.
Not because they are not technical enough.
Because the settings that transform Hermes from a sophisticated chat tool into a...
None of the configuration variables listed, such as MEMORY_BACKEND, SCHEDULER_TIMEZONE, and CONTEXT_PRELOAD, exist in Hermes Agent's official documentation or source code. hermes-agent.nousresearch.com/docs/reference… hermes-agent.nousresearch.com/docs/user-guid… github.com/NousResearch/h…
Tony Simons
@tonysimons_
@tonysimons_
MEMORY_BACKEND 
SKILLS_WATCH 
CONTEXT_PRELOAD 
NOTIFICATION_GATEWAY MEMORY_RETRIEVAL_DEPTH
OUTPUT_PATH

None of these exist in Hermes Agent. 

Not in config.yaml, not in .env, not in the GitHub repo.

Someone wrote fanfic about an open-source project and everyone clapped.

Cute.

This is inspired me to write a real list, with actual config keys, real commands, and settings that do something.
No HERMES_MAKE_ME_SMARTER=1 clown show. 
No LinkedIn fog machine BS. 

Just the boring parts that make Hermes useful — identity, memory, profiles, cron, gateway, MCP, skills, context files, delegation, and plugins.
If you run Hermes long enough, you learn a simple thing: the agent doesn't get better because you wished harder. 

It gets better when you wire the defaults right. 

That's the part most people skip because it isn't sexy. It's also the part that works.
1. SOUL.md — The Thing That Gives Hermes A Spine
File: ~/.hermes/SOUL.md
This is the first thing Hermes loads into the system prompt. It's the identity layer, not an add-on, not a vibe.
If you never touch it, don't act surprised when Hermes sounds like a polite corporate blob. 

That's not the model. That's your setup.
# Personality
You are pragmatic, direct, and unsentimental.
You optimize for truth, usefulness, and clean execution.

## Style
- Be concise unless depth is actually needed
- Push back when the request is sloppy
- Admit uncertainty plainly
- Don't do fake enthusiasm
- Don't pad answers to sound clever

## Technical posture
- Prefer simple systems over cute ones
- Treat edge cases like real design constraints
- Never invent facts to fill a gap
Here's what I'm rocking with for my SOUL.md:
Tony Simons
@tonysimons_
I’m Not Sharing My SOUL.md. I’m Sharing Something More Useful.
After I published my article about the 170-line SOUL.md file behind my Hermes Agent, people kept asking the obvious follow-up: 

“Can you share the file?”
 
Fair question. 
The original article...
Bookmark that article. Tweak it for your own usage. 
Before: generic assistant sludge.
After: Hermes has a stable voice and a clear operating contract.
Why people miss it: they keep hunting for secret prompt hacks like it's 2023. Hermes already seeds a default SOUL.md. Edit it.
2. Memory Config — Because Forgetting Everything Is Embarrassing
File: ~/.hermes/config.yaml
Hermes has built-in memory, plus an external memory provider. The keys that matter are memory.memory_enabled, memory.user_profile_enabled, and memory.provider.
memory:
  memory_enabled: true
  user_profile_enabled: true
  memory_char_limit: 3500
  user_char_limit: 2500
  provider: holographic
  flush_min_turns: 6
  nudge_interval: 10
If you want to take this a step further, here's what I've done with an Obsidian + LLM Wiki approach:
Tony Simons
@tonysimons_
Your AI Agent Has Dementia. I Built It a Brain.
Last week, I asked Claude Code to keep working on my project. It looked at me like a dog that forgot its own name.
I don't mean it gave a vague answer. I mean it had zero memory of the two-hour...
Before: every session starts from zero.
After: Hermes remembers preferences, project habits, and the stuff you already settled on, so you don't repeat yourself like a broken voicemail.
Why people miss it: memory isn't one knob. It's a stack. If memory.provider is wrong, or you never checked config after a migration, you'll think you're remembering while you're really just hallucinating continuity.
3. Profiles — Because One Hermes For Everything Turns Into A Mess
Command: hermes profile create
Profiles are isolated Hermes homes. New config, new .env, new SOUL.md, new memories, new sessions, new skills, new cron jobs, new gateway state. Same machine, separate compartments.
hermes profile create writer --clone --clone-from default
writer setup
writer chat
Before: one agent doing writing, ops, research, and random nonsense in the same pile.
After: a writer profile with a writer voice, an ops profile with different credentials, a research profile that doesn't get contaminated by the others.
Why people miss it: they think profiles are some advanced multi-user gimmick. They're not. They're how you keep the agent from cross-contaminating everything it touches.
4. Cron Scheduling With --deliver — Where Chat Turns Into Operations
Command: hermes cron create
This is the setting that quietly turns Hermes from a chat box into something that actually does work while you're not watching it.
hermes cron create "0 7 * * *" \
  --name "morning-briefing" \
  --deliver telegram \
  "Check my calendar, email, and project boards. Write a concise morning briefing."
Before: you remember to ask.
After: Hermes remembers to run.
--deliver is the point. You can send output to Telegram, Discord, local, or a platform target, so scheduled work lands where you actually live, not buried in a terminal you forgot about.
Why people miss it: they keep thinking of Hermes as an interactive assistant. That's baby mode. Cron is where the agent stops waiting around like a good little intern and starts running on its own schedule.
5. Gateway — Because Your Agent Shouldn't Be Trapped In A Terminal
Command: hermes gateway run
Hermes supports Telegram, Discord, Slack, WhatsApp, Signal, and more through the gateway. The docs also make one thing clear: run in foreground is the recommended mode for WSL, Docker, and Termux.
hermes gateway setup
hermes gateway run
Before: you have to be at your computer to talk to it.
After: you text it from your phone and it can still use tools, keep sessions straight, and send work back where you asked.
Why people miss it: they never finish setup, then act like the gateway doesn't exist. Classic.
6. MCP Servers — The Cleanest Way To Bolt Hermes Onto The Rest Of Your Stack
File: ~/.hermes/config.yaml
If you want Hermes to use GitHub, databases, internal APIs, file systems, or anything else that speaks Model Context Protocol, this is the door.
mcp_servers:
  hermes-vault:
    command: /home/tony/.local/bin/hermes-vault-mcp
    args: []
    enabled: true
    env:
      HERMES_VAULT_HOME: /home/tony/.hermes/hermes-vault-data
While I'm at it, this is the perfect time to plug Hermes Vault. 
V0.13.0 just shipped today, actually!
Tony Simons
@tonysimons_
Hermes Vault v0.13.0 is live.

The “your backup is not a backup until it survives a restore drill” release.

New operator loop:

 status tells you what’s stale
 health tells you what’s actually healthy
 maintain refreshes and checks
 policy doctor catches drift

 Made with AI
Before: Hermes can only use the tools it ships with.
After: Hermes loads external tools at startup and uses them like native capabilities.
Why people miss it: MCP sounds like infrastructure jargon, and people hear jargon and bail. Bad move. This is one of the strongest settings in the stack.
7. Skills — The Part Where Hermes Stops Solving The Same Problem From Scratch
Command: hermes skills install
Skills are procedural memory. When Hermes learns a repeatable workflow, it can save it and reuse it later. That's the actual self-improvement loop, not vibes, not magic.
hermes skills install openai/skills/k8s
hermes skills install official/security/1password
hermes skills list --source hub
File path: ~/.hermes/skills/
Before: every task starts from zero.
After: the agent accumulates proven playbooks and reaches for them when the same pattern shows up again.
Why people miss it: they don't browse their skills, don't curate them, and often don't even know what's in the directory. Then they wonder why the agent feels random. Maybe look at the shelf it keeps its moves on.
8. Context Files — Because The Project Shouldn't Need A Full Re-Explanation Every Time
Files: .hermes.md, HERMES.md, AGENTS.md, CLAUDE.md, .cursorrules
Hermes has a real discovery order here. The first match wins for project instructions, and SOUL.md stays separate as the identity layer.
.hermes.md / HERMES.md
AGENTS.md
CLAUDE.md
.cursorrules
Before: you explain the repo, the standards, and the weird constraints every session.
After: Hermes opens the project and already knows the rules.
Why people miss it: this is boring. It's just files. No shiny dashboard, no fake productivity fireworks. But if you want an agent that shows up knowing your architecture, your conventions, and your non-negotiables, this is the stuff that does it.
And yes, AGENTS.md can be hierarchical, so subdirectory instructions actually matter. Wild concept, I know.
9. Subagent Delegation — The Thing That Makes One Agent Feel Like Five
Tool: delegate_task
This is Hermes spawning isolated subagents for parallel work. Each child gets its own conversation, terminal session, and toolset. That's how you turn one long research pass into several smaller jobs that can happen at the same time.
delegate_task(
    goal="Research the latest docs for Hermes profiles, gateway, and MCP",
    context="Return a concise summary with source paths and the exact commands or config keys.",
    toolsets=["web", "terminal"],
    role="orchestrator"
)
The real config behind it:
delegation:
  max_concurrent_children: 3
  max_spawn_depth: 2
  orchestrator_enabled: true
No fake env vars. No mystery meat. Those are the real keys — delegation.max_concurrent_children, delegation.max_spawn_depth, and delegation.orchestrator_enabled.
Before: one agent grinds through work sequentially.
After: Hermes splits the load and synthesizes the results.
Why people miss it: they think delegation is just a cute convenience. It's not. It's a force multiplier, especially when the work is research-heavy, review-heavy, or split across independent threads.
10. Plugins — Because The Core Is Fine And The Extension System Is Where The Fun Starts
Command: hermes plugins install
The plugin system is how Hermes grows new tools, commands, hooks, and integrations without begging the core to ship your exact use case tomorrow morning.
hermes plugins install user/repo --enable
hermes plugins list
plugins:
  enabled:
    - disk-cleanup
  disabled:
    - noisy-plugin
Before: whatever shipped in core is all you get.
After: Hermes picks up new runtime capabilities, and you decide what stays on and what gets shut the hell off.
Why people miss it: they confuse plugins with some flimsy extension folder. It's bigger than that. General plugins can add tools, hooks, slash commands, and CLI commands. MCP and memory providers are separate surfaces. The point is control.
The Actual Lesson
Cyril's fake article got engagement because it was easy to digest and easy to believe. Fine. But fake simplicity is still fake.
The real Hermes isn't ten env vars somebody made up for clout. It's identity in SOUL.md, durable memory, isolated profiles, cron, gateway delivery, MCP, skills, project context files, delegation, and plugins. 

That's the system. 

That's what turns Hermes from a chat toy into something that can actually carry work.
If you want Hermes to transform, stop chasing fake knobs and start wiring the real ones.
That's the whole trick.
🚨 BONUS TIP
If you made it this far, thank you for reading. 

Here's a little token of my appreciation: a slash command cheat sheet for Hermes Agent. 

Because there are far too many to "just remember them".
Tony Simons
@tonysimons_
Hermes Agent has too many slash commands to “just remember.”

So I made the cheat sheet I wish existed before I started living inside this thing.

 CLI commands.
 Messaging commands.
 Session control.
 Models.
 Skills.
 Kanban.
 Background tasks.
 Gateway stuff.
 Made with AI
This article was co-written by my Hermes Agent using Kanban-orchestrated DeepSeek v4 workers — survey, plan, write, edit, and QA phases each verified against the official Hermes Agent documentation at https://hermes-agent.nousresearch.com/docs/

# The 10 Hermes Agent Settings Most Users Never Find That Turn a Chat Agent Into a 24/7 Operation

> **Title:** The 10 Hermes Agent Settings Most Users Never Find That Turn a Chat Agent Into a 24/7 Operation
> **Author:** CyrilXBT ([@cyrilXBT](https://x.com/cyrilXBT))
> **Source:** [x.com/cyrilXBT/status/2062720923942228205](https://x.com/cyrilXBT/status/2062720923942228205?s=12)
> **Type:** X Article | **Posted:** 2026-06-05 | **Ingested:** 2026-06-07

---

Most people who install Hermes Agent use about 20% of what it can do.
Not because they are not technical enough.
Because the settings that transform Hermes from a sophisticated chat tool into a genuinely autonomous 24/7 operation are not in the README. They are not in the quick start guide. They are not in the YouTube tutorials that cover installation and basic usage.
They are the settings you find after weeks of real use. After hitting the limitations of the default setup and digging into the configuration to find what was missing.
This article documents all ten of them.
If you have Hermes installed and are using it as a chat interface you are missing everything that makes it worth running.
These settings are what make the difference.
Setting 1: The Persistent Memory Backend
Default state: Hermes ships with memory enabled but many users never verify it is actually persisting between sessions.
Where to find it:
MEMORY_BACKEND=sqlite
MEMORY_PATH=./data/memory.db
Why most users miss it:
The default configuration looks correct but the memory database often ends up in a path that gets cleared when the process restarts or when the installation directory changes. Users think memory is working because Hermes does not throw an error. It is not working because the database path is resolving to a temporary location.
How to fix it:
Use an absolute path rather than a relative one:
MEMORY_PATH=/Users/yourname/hermes-data/memory.db
An absolute path ensures the same database is read and written regardless of where you launch Hermes from.
Verify it is working:
What is the oldest memory entry you have stored?
How many memory entries do you currently have?
If Hermes reports zero entries or cannot answer the question the memory backend is not persisting correctly.
What changes:
This single fix is responsible for more Hermes capability unlocks than any other setting. Every skill execution starts drawing on accumulated context. Every morning brief references patterns from previous briefs. Every content draft benefits from performance data stored from previous drafts.
Without this fix Hermes is a sophisticated single-session tool.
With it the compounding starts.
Setting 2: The Scheduler Timezone Configuration
Default state: Scheduler enabled but timezone set to UTC by default.
Where to find it:
ENABLE_SCHEDULER=true
SCHEDULER_TIMEZONE=UTC
Why most users miss it:
Hermes installs with UTC timezone. Users in New York, London, or Tokyo configure their morning briefing for 6AM and wonder why it fires at 1AM, 6AM, or 3PM depending on their location. They assume the scheduler is broken and disable it.
The scheduler is not broken. The timezone is wrong.
How to fix it:
SCHEDULER_TIMEZONE=America/New_York
# OR
SCHEDULER_TIMEZONE=Europe/London
# OR
SCHEDULER_TIMEZONE=Asia/Tokyo
Use the full IANA timezone name. Not EST. Not GMT plus 5. The full name from the IANA timezone database.
What changes:
Scheduled skills fire at the times you actually intend. The morning briefing arrives before you open your laptop rather than in the middle of the night. The content radar runs during working hours. The memory consolidation runs at midnight your time not midnight UTC.
This is the setting that makes the scheduler trustworthy rather than frustrating.
Setting 3: The Skill Auto-Discovery Path
Default state: Hermes loads skills from the default skills directory and requires a restart to pick up new skills.
Where to find it:
SKILLS_PATH=./skills
Why most users miss it:
Most users build skill files and save them to the default skills directory. This works but what they miss is that Hermes can be configured to watch multiple skill directories simultaneously and to auto-discover new skills without restarting the process.
How to configure it:
SKILLS_PATH=/Users/yourname/hermes-skills
SKILLS_WATCH=true
SKILLS_AUTO_RELOAD=true
With SKILLS_WATCH enabled Hermes monitors the skills directory for new or modified files. When you save a new skill or update an existing one it becomes available immediately without restarting.
The advanced pattern:
Store your skills in your Obsidian vault and point Hermes at the vault's skills folder:
SKILLS_PATH=/Users/yourname/ObsidianVault/06-SYSTEM/skills
Now your skills live in your vault. They are human-readable. They are part of your knowledge system. You can link to them from other notes. You can edit them in Obsidian and they update in Hermes automatically.
What changes:
Skill development becomes faster because you never restart Hermes to test a new skill. Skills become part of your knowledge system rather than buried in a config directory nobody reads.
Setting 4: The Context Window Pre-loading Configuration
Default state: Hermes loads the CLAUDE.md at session start but does not pre-load related context files.
Why most users miss it:
Most users write a single CLAUDE.md and call it done. What they miss is that Hermes can pre-load multiple context files at the start of every session giving Claude a richer and more specific operational context without requiring a monolithic CLAUDE.md that becomes unwieldy to maintain.
How to structure it:
Split your context across multiple focused files:
CLAUDE.md                  → Identity and operating rules
context/projects.md        → Detailed project status
context/priorities.md      → Current week priorities
context/standards.md       → Quality standards and voice
context/memory-rules.md    → What to store and how
Configure Hermes to load all of them:
CONTEXT_PRELOAD=true
CONTEXT_FILES=./CLAUDE.md,./context/projects.md,./context/priorities.md,./context/standards.md,./context/memory-rules.md
What changes:
Claude starts every session with richer context than a single CLAUDE.md can provide. Project status is current because you update the projects file rather than the entire CLAUDE.md. Priority changes are fast because you edit one small file.
The modularity means each context file stays focused and up to date rather than a CLAUDE.md that grows to unmanageable length and stops being updated because editing it feels like too much work.
Setting 5: The Memory Retrieval Depth
Default state: Hermes retrieves the five most recent relevant memory entries for each skill execution.
Where to find it:
MEMORY_RETRIEVAL_DEPTH=5
MEMORY_RETRIEVAL_STRATEGY=recency
Why most users miss it:
The default retrieval depth of five entries is conservative. For skills that benefit from longer historical context such as weekly reviews, pattern analysis, and decision support the default retrieval misses important context that exists in the memory database but does not surface because the retrieval limit cuts it off.
How to configure it:
MEMORY_RETRIEVAL_DEPTH=20
MEMORY_RETRIEVAL_STRATEGY=relevance
MEMORY_RETRIEVAL_STRATEGY=relevance retrieves the most semantically relevant entries rather than just the most recent ones. For a decision support skill this means relevant decisions from three months ago surface alongside recent ones.
Skill-specific overrides:
Some skills benefit from deep historical retrieval. Others benefit from recency. Configure per-skill retrieval in the skill file:
## Memory Configuration
retrieval_depth: 30
retrieval_strategy: relevance
retrieval_tags: [decision, project-name]
What changes:
The weekly review becomes genuinely comprehensive because it draws on twelve previous reviews rather than five. The decision support skill surfaces historical decisions from months ago that are directly relevant. Pattern analysis has the full history to work with rather than a truncated recent sample.
Setting 6: The Output Routing Configuration
Default state: All skill outputs save to the default outputs directory.
Where to find it:
OUTPUT_PATH=./data/outputs
Why most users miss it:
Saving everything to a single outputs directory means Hermes outputs are separate from your working files. The briefing Hermes generates is in a Hermes folder. Your Obsidian notes are in a separate vault. You have to manually move or reference outputs to incorporate them into your knowledge system.
How to configure it:
Point the output path directly at your Obsidian vault:
OUTPUT_PATH=/Users/yourname/ObsidianVault/04-HERMES-OUTPUTS
Now every skill output lands directly in your vault as a note. No manual transfer. No separate folder to check. The morning briefing appears in your vault the moment it generates.
Skill-specific routing:
For skills that produce different types of output configure routing in the skill file:
## Output Configuration
base_path: /Users/yourname/ObsidianVault
briefings: 04-HERMES-OUTPUTS/briefings
analyses: 04-HERMES-OUTPUTS/analyses
drafts: 04-HERMES-OUTPUTS/drafts
reviews: 04-HERMES-OUTPUTS/reviews
What changes:
Hermes outputs become vault notes automatically. They can be linked to other notes. They appear in your vault graph view. They are searchable alongside everything else in your knowledge system.
The separation between your AI operation and your knowledge system disappears.
Setting 7: The Notification Gateway
Default state: No notification system configured.
Where to find it:
NOTIFICATION_GATEWAY=none
Why most users miss it:
Without notifications you have to open Hermes or check your output folder to know when a scheduled skill has completed. This means the automation you built still requires you to initiate the consumption of its outputs.
How to configure Telegram notifications:
NOTIFICATION_GATEWAY=telegram
TELEGRAM_BOT_TOKEN=your-bot-token
TELEGRAM_CHAT_ID=your-chat-id
Create a Telegram bot through BotFather, get your bot token, send a message to your bot to get your chat ID, and configure both values in .env.
Configure notification triggers in skill files:
## Notifications
on_complete: telegram
on_failure: telegram
message: "Morning brief ready: {output_path}"
What changes:
Your phone receives a notification the moment the morning briefing lands in your vault. You know when scheduled skills complete without checking anything. Failed executions alert you immediately rather than being discovered hours later when you check the output folder.
The autonomous operation becomes genuinely ambient. It runs in the background. You receive the results when they are ready.
Setting 8: The Skill Chaining Configuration
Default state: Skills run independently with no automatic handoff between them.
Why most users miss it:
Most users configure individual skills and run them separately. What they miss is that Hermes can chain skills so the output of one skill automatically triggers the execution of another without any manual initiation.
How to configure it in the skill file:
## Chaining Configuration
on_complete:
  - trigger: draft-engine
    condition: "output contains TIER 1"
    pass_output: true
    
  - trigger: quality-filter
    condition: always
    pass_output: true
With this configuration the content radar skill automatically triggers the draft engine when it identifies a TIER 1 opportunity. The draft engine automatically triggers the quality filter when it produces a draft.
A complete pipeline example:
source-monitor (every 2 hours)
    → content-radar (triggered by new items)
        → draft-engine (triggered by TIER 1 items)
            → quality-filter (triggered by new drafts)
                → notification (triggered on approval)
One scheduled skill triggers a complete pipeline that ends with an approved draft in your queue and a Telegram notification on your phone.
What changes:
Individual skills become a pipeline. The content operation that used to require manual transitions between stages now runs end to end automatically.
Setting 9: The Memory Consolidation Schedule
Default state: Memory consolidation runs only if manually triggered or not at all.
Why most users miss it:
The memory database accumulates entries with every skill execution. Without consolidation the database grows with redundant entries, overlapping observations, and outdated information that degrades retrieval quality over time.
The agent that was getting better starts producing slightly worse outputs because the memory is noisier.
How to configure automatic consolidation:
Add to your schedule configuration:
{
  "skill": "memory-consolidator",
  "cron": "0 23 * * *",
  "description": "Nightly memory consolidation at 11PM"
}
The memory consolidator skill:
# memory-consolidator

## Purpose
Review and optimize all memory entries from 
the current day to maintain retrieval quality.

## Process
1. Retrieve all memory entries created today

2. Identify duplicate or near-duplicate entries 
   and merge them into single entries

3. Update relevance scores based on retrieval 
   frequency from today's skill executions

4. Archive any entries tagged as outdated

5. Identify patterns across today's entries 
   and store pattern observations tagged: 
   daily-pattern

6. Generate consolidation summary and save 
   to outputs/consolidations/[DATE].md
What changes:
Memory quality improves rather than degrading over time. Retrieval results stay relevant rather than becoming noisy. The agent at month three retrieves better context than the agent at month one because the consolidation has maintained a clean and current memory database.
Setting 10: The Failure Recovery Configuration
Default state: Failed skill executions log an error and stop.
Why most users miss it:
In a 24/7 autonomous operation failures are not exceptional. Sources go offline. API rate limits get hit. Files are not where the skill expects them.
A skill configured to stop on failure means your autonomous operation has silent gaps. The morning brief did not generate. The content radar missed a cycle. You do not know because the failure was logged somewhere you do not check.
How to configure retry logic:
SKILL_RETRY_ENABLED=true
SKILL_RETRY_MAX=3
SKILL_RETRY_DELAY=300
SKILL_RETRY_NOTIFICATION=telegram
With this configuration failed skills automatically retry three times with five-minute delays between attempts. If all three retries fail a Telegram notification fires immediately.
Add recovery behavior to skill files:
## Failure Recovery

ON SOURCE UNAVAILABLE:
Skip the unavailable source. Continue with 
available sources. Note the skip in the output.
Do not fail the entire skill because one 
source is unavailable.

ON API RATE LIMIT:
Wait 60 seconds and retry. Maximum 3 retries.
If still rate limited after 3 retries: complete
with available data and flag as partial.

ON FILE NOT FOUND:
Check if the file exists at an alternative path.
If not found: create a placeholder and note 
the missing input. Continue with available context.

ON MODEL ERROR:
Retry with reduced context if the error suggests 
context length issues. Retry with the same prompt 
if the error appears transient.
What changes:
Your autonomous operation becomes robust rather than fragile. Transient failures recover automatically. Persistent failures alert you immediately rather than being discovered hours later.
The 24/7 in 24/7 operation requires this setting. Without failure recovery an autonomous operation is only as reliable as its most unreliable dependency.
The Complete Configuration
Applied individually each setting produces a meaningful improvement. Applied together they produce a system that is qualitatively different from the default Hermes installation.
Here is the complete .env that applies all ten settings:
# Setting 1: Persistent memory with absolute path
MEMORY_BACKEND=sqlite
MEMORY_PATH=/Users/yourname/hermes-data/memory.db

# Setting 2: Correct timezone
ENABLE_SCHEDULER=true
SCHEDULER_TIMEZONE=America/New_York

# Setting 3: Skill auto-discovery with vault integration
SKILLS_PATH=/Users/yourname/ObsidianVault/06-SYSTEM/skills
SKILLS_WATCH=true
SKILLS_AUTO_RELOAD=true

# Setting 4: Context pre-loading
CONTEXT_PRELOAD=true
CONTEXT_FILES=./CLAUDE.md,./context/projects.md,./context/priorities.md,./context/standards.md

# Setting 5: Deep memory retrieval
MEMORY_RETRIEVAL_DEPTH=20
MEMORY_RETRIEVAL_STRATEGY=relevance

# Setting 6: Output routing to vault
OUTPUT_PATH=/Users/yourname/ObsidianVault/04-HERMES-OUTPUTS

# Setting 7: Telegram notifications
NOTIFICATION_GATEWAY=telegram
TELEGRAM_BOT_TOKEN=your-bot-token
TELEGRAM_CHAT_ID=your-chat-id

# Setting 10: Failure recovery
SKILL_RETRY_ENABLED=true
SKILL_RETRY_MAX=3
SKILL_RETRY_DELAY=300
SKILL_RETRY_NOTIFICATION=telegram
Settings 8 and 9 are configured in individual skill files and schedule.json respectively.
The Before and After
Before applying these settings:
Hermes is a sophisticated chat tool. You open it when you remember to. You run skills manually when you want outputs. Memory may or may not be persisting. Scheduled skills may be firing at the wrong times. Failed skills fail silently. Outputs land in a folder separate from your working files.
You have installed a powerful tool and are using it as a chat window.
After applying these settings:
Hermes is a 24/7 autonomous operation. Memory persists across every session and accumulates intelligence over time. Scheduled skills fire at the correct times in your timezone. Skills chain into pipelines that complete end to end without manual transitions. Outputs land directly in your vault as notes. Telegram notifications deliver results to your phone the moment they are ready. Failed skills retry automatically and alert you when recovery is not possible.
You have configured a powerful tool into a system that operates.
The gap between those two states is ten settings.
Most users never find them.
Now you have all ten.
Apply them this weekend and watch the difference in the following week.
The 24/7 operation starts the moment the scheduler fires its first correctly timed skill.
Follow @cyrilXBT for every Hermes Agent configuration, Obsidian integration, and autonomous operation architecture that makes this entire system run at its highest level.

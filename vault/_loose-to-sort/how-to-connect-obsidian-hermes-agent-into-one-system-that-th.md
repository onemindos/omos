# How to Connect Obsidian + Hermes Agent Into One System That Thinks, Remembers, and Runs Your Life

> **Title:** How to Connect Obsidian + Hermes Agent Into One System That Thinks, Remembers, and Runs Your Life
> **Author:** CyrilXBT ([@cyrilXBT](https://x.com/cyrilXBT))
> **Source:** [x.com/cyrilXBT/status/2061601347271577630](https://x.com/cyrilXBT/status/2061601347271577630?s=12)
> **Type:** X Article | **Posted:** 2026-06-02 | **Ingested:** 2026-06-07

---

Most productivity systems have a memory problem.
Your Obsidian vault stores everything you know but cannot act on it without you initiating something.
Your AI agent can act on anything but forgets everything the moment the session closes.
Connect them and the problem disappears entirely.
Obsidian becomes the permanent knowledge layer. Every note, every decision, every captured idea, every project update lives there as plain text that never disappears.
Hermes Agent becomes the intelligence layer that reads that knowledge, reasons across it, executes workflows against it, and writes new knowledge back into it automatically.
Together they form a system that thinks with your accumulated knowledge, remembers everything across every session, and runs your most important recurring workflows while you sleep.
This guide is the complete build.
Why This Combination Works
Before the architecture understand why these two tools belong together.
The problem with Obsidian alone:
Obsidian is the best knowledge storage system available. Plain text files. Local storage. Permanent. Searchable. Linkable. But Obsidian is passive. It stores what you put in and returns it when you search for it. It does not read your notes and surface insights. It does not run workflows. It does not generate outputs from your accumulated knowledge without you manually asking for them.
The problem with Hermes Agent alone:
Hermes is the best autonomous agent framework available. Persistent memory. Skill execution. Scheduled automation. MCP connections to real tools. But Hermes memory lives in a SQLite database. It is not human-readable. You cannot open it, browse it, or edit it directly. You cannot build a rich knowledge network inside it the way you can in Obsidian.
What the combination produces:
Obsidian handles human-readable permanent storage. You read it, write in it, link notes together, and build a knowledge network that compounds over years.
Hermes reads your Obsidian vault via the Filesystem MCP, reasons across everything in it using Claude, executes automated workflows that produce new knowledge, and writes the outputs back into your vault as new notes.
The vault is not just where you store things. It is the operating environment that Hermes runs against.
Every note you write becomes context Hermes uses.
Every output Hermes generates becomes a note you can read, edit, and link.
The two systems compound together rather than operating in parallel.
The Architecture
The system has four layers that work together.
Layer 1: The Vault Layer — Obsidian
Plain text Markdown files organized in a structure designed for both human navigation and automated processing. Every note you write, every idea you capture, every decision you make lives here.
Layer 2: The Connection Layer — Filesystem MCP
The bridge between Obsidian and Hermes. The Filesystem MCP gives Hermes direct read and write access to your Obsidian vault. Hermes can read any note, search across all notes, create new notes, and update existing ones.
Layer 3: The Intelligence Layer — Hermes + Claude
Hermes executes skill files that call Claude via the Anthropic API. Claude reads your vault, reasons across your accumulated knowledge, and produces outputs that go back into the vault as new notes.
Layer 4: The Automation Layer — Hermes Scheduler
Scheduled workflows that fire automatically at configured times. Morning briefings generated from your vault before you wake up. Weekly reviews synthesized every Sunday evening. Project status updates produced every Monday morning. All of it running without you initiating anything.
Setting Up the Connection
Step 1: Install Hermes Agent
git clone https://github.com/hermes-agent/hermes
cd hermes
npm install
Step 2: Configure Claude as the model
Open .env and configure:
MODEL_PROVIDER=anthropic
MODEL_NAME=claude-opus-4-8
ANTHROPIC_API_KEY=your-anthropic-api-key

MEMORY_BACKEND=sqlite
MEMORY_PATH=./data/memory.db

ENABLE_SCHEDULER=true
SCHEDULER_TIMEZONE=America/New_York
Step 3: Configure the Filesystem MCP to point at your Obsidian vault
MCP_FILESYSTEM_PATH=/Users/yourname/Documents/YourVaultName
This single configuration is what connects everything. Hermes now has direct access to every note in your vault.
Step 4: Install and configure the MCP server
npm install -g @modelcontextprotocol/server-filesystem
npm install -g @modelcontextprotocol/server-brave-search
Step 5: Start Hermes and verify the connection
npm run start
Test the connection:
List the top-level folders in my Obsidian vault 
and describe the structure you see.
Hermes should accurately describe your vault's folder structure. The connection is live.


The Vault Structure for the Combined System
The vault structure needs to accommodate both human navigation and automated processing. Hermes needs to know where to read inputs and where to write outputs. You need to know where to find everything Hermes produces.
VAULT/
    00 - INBOX/
        [raw captures that need processing]
    
    01 - NOTES/
        permanent/
            [atomic notes in your own words]
        daily/
            [YYYY-MM-DD].md
        meetings/
            [YYYY-MM-DD]-[meeting-name].md
    
    02 - PROJECTS/
        [project-name]/
            overview.md
            tasks.md
            notes/
    
    03 - RESOURCES/
        [topic-based reference notes]
    
    04 - HERMES-OUTPUTS/
        briefings/
            [YYYY-MM-DD]-morning-briefing.md
        analyses/
            [YYYY-MM-DD]-[topic]-analysis.md
        syntheses/
            [YYYY-MM-DD]-[topic]-synthesis.md
        reviews/
            [YYYY-MM-DD]-weekly-review.md
    
    05 - ARCHIVE/
        [completed projects and outdated material]
    
    06 - SYSTEM/
        CLAUDE.md
        skills/
        templates/
The HERMES-OUTPUTS folder is the critical addition. Every automated output Hermes generates lands here. You always know where to find what Hermes produced. You can read it, edit it, link it to other notes, and build on it.
The SYSTEM folder contains your CLAUDE.md and optionally copies of your Hermes skill files for reference. Keeping skill documentation in the vault means your notes about how the system works live alongside the knowledge the system produces.


The CLAUDE.md That Powers the System
The CLAUDE.md lives in your Obsidian vault at 06 - SYSTEM/CLAUDE.md. Hermes reads it via the Filesystem MCP at the start of every skill execution.
# Obsidian + Hermes System — CLAUDE.md

## Vault Location
This file is at: 06 - SYSTEM/CLAUDE.md
All outputs save to: 04 - HERMES-OUTPUTS/

## Who I Am
Name: [YOUR NAME]
What I do: [YOUR PRIMARY WORK]
Primary focus: [YOUR MOST IMPORTANT AREA]

## How My Vault Is Organized
00 - INBOX: Raw unprocessed captures
01 - NOTES: Processed permanent notes and daily notes
02 - PROJECTS: Active project folders
03 - RESOURCES: Reference material by topic
04 - HERMES-OUTPUTS: All automated outputs
05 - ARCHIVE: Completed and outdated material
06 - SYSTEM: CLAUDE.md and system files

## Active Projects
[PROJECT 1]: [ONE SENTENCE STATUS]
Next action: [SPECIFIC NEXT STEP]

[PROJECT 2]: [ONE SENTENCE STATUS]
Next action: [SPECIFIC NEXT STEP]

## Current Priorities
1. [MOST IMPORTANT THING THIS WEEK]
2. [SECOND MOST IMPORTANT]
3. [THIRD MOST IMPORTANT]

## Topics I Think About Most
[LIST 5-8 RECURRING TOPICS]

## Writing Voice
[HOW YOU WRITE AND COMMUNICATE]
[Specific things you never say or do]

## Output Instructions
When generating any output:
- Read the relevant section of my vault first
- Reference specific notes by their file path
- Save outputs to the correct subfolder 
  in 04 - HERMES-OUTPUTS
- Date-stamp every file: YYYY-MM-DD-[type]-[topic].md
- Create proper Obsidian frontmatter with type and date

## Memory Instructions
Before any skill execution:
- Read relevant memory from Hermes database
- Read relevant notes from vault via Filesystem MCP
- Combine both sources for full context

After any skill execution:
- Store summary in Hermes memory with appropriate tags
- Create or update the relevant vault note

## Weekly Update
This file is reviewed and updated every Monday.
Current week priorities updated: [DATE]
The Seven Skills That Make the System Work
Skill 1: The Vault-Aware Morning Brief
This skill reads your actual Obsidian vault to generate a morning brief that is specific to your projects, your notes, and your current priorities rather than generic news.
Create skills/vault-morning-brief.md:
# vault-morning-brief

## Purpose
Generate a morning brief using my actual Obsidian 
vault as the primary context, supplemented by 
real-time web search for external developments.

## Trigger
Scheduled daily at 6:00 AM.
Manual: "Generate my morning brief"

## Process
1. Read CLAUDE.md from my vault at 
   06 - SYSTEM/CLAUDE.md

2. Read today's daily note if it exists at
   01 - NOTES/daily/[TODAY'S DATE].md

3. Read all project overview.md files in 
   02 - PROJECTS/ to understand current project status

4. For each active project in CLAUDE.md:
   Check if any tasks are overdue by reading 
   the project's tasks.md file

5. Read the most recent weekly review from 
   04 - HERMES-OUTPUTS/reviews/ to understand 
   last week's priorities and outcomes

6. Use Brave Search to find external developments 
   relevant to my current priorities from CLAUDE.md

7. Synthesize everything into the morning brief:

---
# Morning Brief — [DATE]

## MOST IMPORTANT TODAY
[Single highest-leverage action based on 
vault context and current priorities]

## PROJECT PULSE
For each active project:
[PROJECT NAME]: [One sentence status] 
Next: [Specific next action from tasks.md]

## OPEN LOOPS
[Any items from yesterday's daily note 
or recent notes marked as unfinished]

## EXTERNAL SIGNAL
[1-3 external developments relevant to 
my current priorities]

## FROM MY VAULT
[One connection between something in my 
notes and something happening today]
---

8. Save to: 04 - HERMES-OUTPUTS/briefings/[DATE]-morning-brief.md

9. Create today's daily note at 
   01 - NOTES/daily/[DATE].md if it does not exist,
   with a link to the morning brief

## Output Format
The brief should reference specific note paths 
in my vault. Not generic observations. 
Specific connections to actual notes I have written.
Skill 2: The Inbox Processor
This skill empties your INBOX folder every evening, filing notes into the correct location and creating links to related permanent notes.
Create skills/inbox-processor.md:
# inbox-processor

## Purpose
Process all files in 00 - INBOX, file them 
in the correct vault location, and create 
connections to related existing notes.

## Trigger
Scheduled daily at 8:00 PM.
Manual: "Process my inbox"

## Process
1. Read CLAUDE.md for vault structure and context.

2. List all files currently in 00 - INBOX

3. For each file in INBOX:

   a. READ the file content fully
   
   b. CLASSIFY the note type:
      - PERMANENT IDEA: deserves its own atomic note
      - PROJECT NOTE: belongs in a specific project
      - REFERENCE: belongs in 03 - RESOURCES
      - DAILY CAPTURE: belongs in today's daily note
      - TASK: belongs in a project's tasks.md
   
   c. FIND CONNECTIONS: Search the vault for 
      existing notes that relate to this capture.
      Look in 01 - NOTES/permanent/ primarily.
   
   d. PROCESS the note:
      
      For PERMANENT IDEA:
      Create a new note at 
      01 - NOTES/permanent/[concept-name].md
      with proper frontmatter and links to 
      related notes found in step c.
      
      For PROJECT NOTE:
      Append to the relevant project's notes 
      or create a new note in 
      02 - PROJECTS/[project-name]/notes/
      
      For REFERENCE:
      Create or append to the relevant topic 
      file in 03 - RESOURCES/
      
      For DAILY CAPTURE:
      Append to 01 - NOTES/daily/[TODAY].md
      
      For TASK:
      Append to the relevant project's tasks.md
   
   e. ARCHIVE the original inbox file to 
      05 - ARCHIVE/inbox-processed/

4. Generate a processing summary:
   [N] notes processed
   [N] permanent notes created
   [N] project notes filed
   [N] connections made

5. Save summary to: 
   04 - HERMES-OUTPUTS/[DATE]-inbox-summary.md

## Quality Standard
Every permanent note created should have 
at least one link to an existing note. 
A permanent note with no connections is 
not ready for the permanent folder.
Skill 3: The Project Health Monitor
This skill reads every active project folder in your vault and produces a health assessment every Monday morning.
Create skills/project-health.md:
# project-health

## Purpose
Read all active project folders and produce 
a comprehensive project health report.

## Trigger
Scheduled every Monday at 7:00 AM.
Manual: "Check project health"

## Process
1. Read CLAUDE.md for list of active projects.

2. For each active project in 02 - PROJECTS/:
   
   a. Read overview.md
   b. Read tasks.md
   c. Read all notes in the project's notes/ folder
   d. Check date of last modification on any file 
      in the project folder

3. For each project Claude assesses:

   STATUS: On Track / At Risk / Stalled / Blocked
   
   EVIDENCE: What specifically in the project files 
   indicates this status. Reference actual file content.
   
   NEXT ACTION: Single most important thing this 
   project needs in the next 48 hours.
   
   STALL FLAG: If no file in the project folder 
   has been modified in 7 or more days flag as stalled.

4. Generate the health report:

---
# Project Health Report — [DATE]

## Summary
[N] active projects: [N] on track, [N] at risk, [N] stalled

## Project Status

### [PROJECT NAME]
Status: [STATUS]
Evidence: [SPECIFIC REFERENCE TO VAULT CONTENT]
Next action: [SPECIFIC ACTION]
[Repeat for each project]

## Attention Required
[Any project flagged as stalled or blocked 
with specific recommendation]
---

5. Save to: 04 - HERMES-OUTPUTS/reviews/[DATE]-project-health.md

6. For any stalled project: add a REVIEW NEEDED 
   note to the project's overview.md with today's date

## Note
Reference specific file paths and note content 
in the assessment. Generic observations about 
project status are not useful. The report should 
only say things that can be traced to specific 
notes in the vault.
Skill 4: The Note Connection Finder
Once per week this skill reads all notes created or modified in the past seven days and finds connections to older notes that were not made at creation time.
Create skills/connection-finder.md:
# connection-finder

## Purpose
Find connections between recent notes and 
older permanent notes that were not explicitly 
linked at creation time.

## Trigger
Scheduled every Sunday at 5:00 PM.
Manual: "Find note connections"

## Process
1. List all files in 01 - NOTES/ modified 
   in the past 7 days.

2. For each recently modified note:
   
   a. Read the full note content
   
   b. Search all permanent notes in 
      01 - NOTES/permanent/ for semantic 
      connections to this note's content
   
   c. Identify connections that are NOT 
      already present as wikilinks in the note
   
   d. For each new connection found:
      - Name both notes
      - Describe specifically how they connect
      - Assess connection strength: 
        Strong / Moderate / Weak
      - Suggest the link text to use

3. Filter to only Strong and Moderate connections.

4. Generate the connection report:

---
# Connection Report — [DATE]

## New Connections Found
[N] connections across [N] notes

### [NOTE NAME]
Path: [FILE PATH]
Connects to: [[OTHER NOTE NAME]]
Reason: [SPECIFIC DESCRIPTION OF CONNECTION]
Suggested link text: [HOW TO PHRASE THE LINK]
[Repeat for each connection]

## Recommended Actions
Add these links this week to strengthen 
the knowledge network.
---

5. Save to: 04 - HERMES-OUTPUTS/analyses/[DATE]-connections.md

## Quality Standard
Only surface non-obvious connections. If the 
link is already present or the connection is 
trivially obvious from the note titles skip it. 
The value is in finding the connections a human 
reader would miss.
Skill 5: The Weekly Vault Synthesis
Every Sunday evening this skill reads your entire week of activity in the vault and produces a synthesis that reveals insights visible only across seven days of accumulated notes.
Create skills/weekly-synthesis.md:
# weekly-synthesis

## Purpose
Synthesize the full week's vault activity 
into a comprehensive review that surfaces 
patterns invisible in any individual day.

## Trigger
Scheduled every Sunday at 7:00 PM.
Manual: "Generate weekly synthesis"

## Process
1. Read CLAUDE.md for current priorities 
   and active projects.

2. Read all daily notes from the past 7 days 
   in 01 - NOTES/daily/

3. Read all permanent notes created or 
   modified this week in 01 - NOTES/permanent/

4. Read all Hermes outputs from this week 
   in 04 - HERMES-OUTPUTS/

5. Read the project health report from Monday 
   in 04 - HERMES-OUTPUTS/reviews/

6. Claude synthesizes across all of it:

---
# Weekly Synthesis — [WEEK OF DATE]

## THE WEEK IN ONE LINE
[Single sentence capturing the most 
important thing this week revealed]

## WHAT MOVED FORWARD
[Specific progress on active projects 
with references to vault evidence]

## WHAT DID NOT MOVE
[Honest assessment of stalled items 
and the most likely reason]

## IDEAS THAT EMERGED
[New permanent notes created this week 
that deserve attention]

## THE PATTERN
[One theme that appeared repeatedly 
across the week's notes]

## CONNECTIONS MADE
[Most significant connections found 
between notes this week]

## NEXT WEEK PRIORITIES
[Three specific priorities for next week 
ranked by impact, each with a next action]

## CLAUDE.MD UPDATE NEEDED
[Any priority or project status changes 
that should be reflected in CLAUDE.md]
---

7. Save to: 04 - HERMES-OUTPUTS/reviews/[DATE]-weekly-synthesis.md

8. Update CLAUDE.md Current Priorities section 
   with next week's top 3.

## Note
The synthesis should be the document you read 
on Monday morning before anything else. 
It should orient you to the week ahead using 
the full context of the week that just ended.
Skill 6: The Research-to-Notes Converter
When you drop research material into INBOX, this skill converts it into properly structured permanent notes with connections to your existing knowledge.
Create skills/research-converter.md:
# research-converter

## Purpose
Convert research material dropped in INBOX 
into properly structured permanent notes 
connected to existing vault knowledge.

## Trigger
Manual: "Convert research in inbox" or 
"Process research on [TOPIC]"

## Process
1. Read CLAUDE.md for vault context and 
   current priorities.

2. Identify research material in 00 - INBOX 
   by reading file contents and looking for 
   source URLs, article text, or research notes.

3. For each research item:

   a. READ the full content
   
   b. EXTRACT the key ideas:
      - What is the core claim?
      - What is the supporting evidence?
      - What is genuinely new versus 
        what I likely already know?
   
   c. SEARCH my existing vault for 
      01 - NOTES/permanent/ notes on related topics
   
   d. CREATE a literature note at:
      01 - NOTES/permanent/lit-[source-topic].md
      
      Literature note format:
      ---
      type: literature
      source: [SOURCE NAME]
      date_read: [TODAY]
      topic: [TOPIC]
      ---
      
      # [SOURCE TITLE]
      
      ## Core Argument
      [What this source claims in my own words]
      
      ## Key Evidence  
      [Specific evidence worth remembering]
      
      ## My Reaction
      [What I agree with, doubt, or want to explore]
      
      ## Connections
      [[RELATED PERMANENT NOTE 1]]
      [[RELATED PERMANENT NOTE 2]]
      
      ## Ideas to Develop
      [Specific ideas worth turning into 
      their own permanent notes]
   
   e. For each idea flagged in Ideas to Develop:
      Create a new permanent note at
      01 - NOTES/permanent/[concept-name].md
   
   f. Move original inbox item to archive

4. Generate conversion summary:
   [N] literature notes created
   [N] permanent notes created
   [N] connections to existing notes

5. Save summary to:
   04 - HERMES-OUTPUTS/analyses/[DATE]-research-converted.md

## Quality Standard
Every literature note must produce at least 
one new permanent note. If the research 
produced no new ideas worth capturing as 
permanent notes the research was not worth 
processing in the first place.
Skill 7: The Thinking Partner
This skill activates Claude as an active thinking partner that reads your recent notes and asks questions, surfaces contradictions, and pushes your thinking forward.
Create skills/thinking-partner.md:
# thinking-partner

## Purpose
Activate Claude as an active thinking partner 
that reads my recent vault activity and 
engages with my ideas rather than just 
storing or summarizing them.

## Trigger
Manual: "Thinking partner session" or 
"Challenge my thinking on [TOPIC]"

## Process
1. Read CLAUDE.md for current intellectual context.

2. Read all permanent notes modified in 
   the past 14 days.

3. Read recent daily notes from the past 7 days.

4. Claude identifies:

   INTERESTING TENSIONS:
   Notes where I appear to hold positions 
   that are in tension with each other.
   Not contradictions necessarily. 
   Tensions worth examining.
   
   UNDERDEVELOPED IDEAS:
   Permanent notes that make a claim without 
   fully developing the reasoning or evidence.
   
   MISSING CONNECTIONS:
   Ideas in recent notes that connect to 
   older permanent notes in non-obvious ways.
   
   OPEN QUESTIONS:
   Questions raised in recent notes that 
   have not been answered anywhere in the vault.

5. Claude engages with each finding:
   Not just naming them. Pushing on them.
   
   For TENSIONS: "You wrote X in [NOTE A] 
   and Y in [NOTE B]. These seem to pull 
   in different directions. Which position 
   do you actually hold and why?"
   
   For UNDERDEVELOPED IDEAS: "The note on 
   [TOPIC] makes the claim that [CLAIM]. 
   What is the strongest argument against 
   this? Does your vault contain evidence 
   that challenges it?"
   
   For OPEN QUESTIONS: "You asked [QUESTION] 
   in your note from [DATE]. You have not 
   answered it anywhere. Is this worth 
   investigating or has it become irrelevant?"

6. Save the session to:
   04 - HERMES-OUTPUTS/analyses/[DATE]-thinking-partner.md

## Tone
This skill should feel like talking to a 
smart colleague who has read everything in 
your vault and is genuinely engaging with 
your ideas. Not validating. Not summarizing. 
Pushing. Questioning. Connecting.
The Daily Experience
With the full system running the daily experience of your Obsidian vault changes completely.
6:00 AM: Hermes reads your vault and generates the morning brief. You wake up to a note in HERMES-OUTPUTS/briefings/ that knows your project status, your open loops, your current priorities, and the external developments relevant to your work. Five minutes to read. You know exactly what matters before you see a single notification.
During the day: You capture everything to 00 - INBOX. Ideas, tasks, references, meeting notes. No filing decisions. No thinking about where things go. Just capture.
8:00 PM: Hermes processes your inbox automatically. Every capture gets filed, linked, and connected to your existing knowledge. Your INBOX empties itself. Nothing sits unprocessed.
Monday 7:00 AM: Project health report waiting in HERMES-OUTPUTS/reviews. Every active project assessed. Stalled projects flagged. Next actions surfaced. Your week starts with a complete picture of where everything stands.
Sunday 5:00 PM: Connection report lands in HERMES-OUTPUTS. Links between recent notes and older permanent notes that you did not consciously make surface automatically. Your knowledge network gets denser every week without additional effort.
Sunday 7:00 PM: Weekly synthesis ready. The full week's vault activity synthesized into a document that shows what moved, what did not, what patterns emerged, and what next week's priorities should be. Your CLAUDE.md updates automatically with next week's top three priorities.
What Happens After 90 Days
The power of this combination is not visible on day one.
It becomes visible over time as two separate compounding effects interact.
Obsidian compounds because your permanent notes accumulate connections over months. A note written in January gets linked to notes written in March, April, and May. The network gets denser. The connection finder surfaces increasingly non-obvious relationships. The vault at month three is qualitatively richer than the vault at month one.
Hermes compounds because memory accumulates across every session. The morning brief at day ninety has read three months of your daily notes and knows your patterns. The weekly synthesis at week twelve has twelve weeks of previous syntheses to draw on. The project health monitor at month three knows how your projects typically stall and surfaces that pattern before it becomes a problem.
The interaction compounds because the more Obsidian knows the better Hermes can reason. The more Hermes produces the richer the vault becomes. The richer the vault the better the next Hermes output.
By month three you have a system that no longer feels like a tool you use.
It feels like a thinking partner who has been working alongside you for months, knows everything you have written, and surfaces the connections and insights your own working memory cannot hold.
That is the system this combination is capable of producing.
Build it this weekend.
The compounding starts from the first skill that reads your vault.
Follow @cyrilXBT for every Obsidian system, Hermes Agent skill, and autonomous operation architecture that makes this combination run at its highest level.

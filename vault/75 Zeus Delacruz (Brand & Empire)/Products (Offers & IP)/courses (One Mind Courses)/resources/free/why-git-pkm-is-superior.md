# Why Git-Based PKM Is Superior

> The argument for storing your knowledge in plain text files under version control — and why every alternative eventually fails you.

---

## The Core Argument

Your knowledge is the most valuable thing you own. Where you store it determines whether you actually own it.

Most people store their knowledge in tools they don't control: Notion, Roam, Obsidian Sync, Confluence, Apple Notes. When those tools change their pricing, shut down, change their API, or get acquired — your knowledge is hostage.

Git-based PKM puts your knowledge in:
- **Plain text files** (markdown) — readable by any tool, now or in 20 years
- **Your own filesystem** — on your hardware, not in a SaaS database
- **A version-controlled repository** — every change tracked, every version recoverable

This is the architecture One Mind OS is built on. Here's why it's the right choice.

---

## The 8 Advantages

### 1. Portability — Your Files Work Everywhere

Plain markdown files open in any text editor on any operating system. Obsidian, VS Code, Neovim, iA Writer, Typora, the terminal — they all read the same files.

In 10 years, Obsidian might not exist. Your markdown files will still open fine.

Notion files require Notion. Roam files require Roam. Obsidian vault files in plain markdown require nothing except a computer.

**The test:** Can you read your notes if your PKM app is down? With git-based PKM: yes, always.

### 2. Version History — Nothing Is Ever Lost

Every change to every file is tracked by git. You can see:
- What a note looked like 6 months ago
- When a specific line was added and why
- Who made a change (if you collaborate)
- How your thinking on a topic evolved over years

You can `git diff` any two versions. You can `git checkout` a file back to any previous state. You can `git log` to see the entire history of a note.

No other PKM system gives you this. Notion has "version history" — locked behind a paywall, limited to 30 days on most plans, and you can't programmatically query it.

**The test:** What did your note on [topic] say 3 months ago? With git: `git log --follow notes/topic.md`. With Notion: pay more, maybe.

### 3. Branching — Experiment Without Risk

Want to restructure your entire vault? Create a branch. Experiment freely. If it works, merge. If it doesn't, delete the branch and your main vault is untouched.

Try that in Notion. There is no branch. There is no undo for reorganizations. There is only hoping you remember what it looked like.

**The test:** Can you test a major vault restructure without risking your live notes? With git: yes, trivially. With any SaaS PKM: no.

### 4. Diff and Audit — See Exactly What Changed

`git diff` shows you precisely what changed between any two states. This is invaluable when:
- You want to review a week of writing
- You want to understand why a note evolved the way it did
- You want to audit what your AI agents wrote to your vault
- You're debugging why a linked note broke

In One Mind OS, every agent write to the vault is automatically committed with a message. The git log is a complete audit trail of everything your agents ever changed.

**The test:** What did your agents write to your vault this week? With git: `git log --since="7 days ago" --oneline`. With any other system: hope your agent kept a separate log.

### 5. AI-Native — Every LLM Reads Plain Text

Plain markdown files are what AI agents read best. No API. No export. No conversion. Your agents directly read and write your vault because it's just files on a filesystem.

Notion requires API calls with rate limits. Roam has a proprietary format. Confluence requires authentication and a custom integration. Your git-based vault requires `read_file`.

This is why One Mind OS agents have a `codex-vault` MCP server that maps directly to the file system — because the vault is just files. No abstraction layer. No middleware. No rate limits.

**The test:** Can your AI agent read any file in your knowledge base without special API access? With git-based PKM: yes. With SaaS PKM: only with API integration, auth tokens, and rate limit handling.

### 6. Sovereignty — You Own Your Data

Your vault lives on your hardware. It syncs through GitHub (a repo you control). No company has unilateral access to your knowledge.

Notion reads your data. Obsidian Sync reads your data. Every SaaS system that hosts your data has the technical ability to read it, train on it, or shut you out.

With git-based PKM, your data lives:
- On your Mac (the working copy)
- On your VPS (the agent-accessible copy)
- On GitHub (the sync layer you control)

You can delete any of these. You can move to a self-hosted git server. You are never locked in.

**The test:** Can you fully migrate your knowledge base without the vendor's help? With git: `git clone` anywhere. With Notion: export once (limited format), hope nothing breaks.

### 7. Collaboration Without Platform Lock-In

Git is the universal collaboration protocol. The entire software world already knows how to use it. Pull requests, issues, reviews — all available.

For One Mind OS multi-agent setups: multiple agents can write to the vault simultaneously because git handles merge conflicts. `codex-watch` commits and pushes. `codex-pull` syncs. Conflicts are resolved systematically, not by hoping two simultaneous writes don't corrupt each other.

**The test:** Can two agents write to your knowledge base simultaneously without data corruption? With git: yes, merge conflict resolution is solved. With Notion: last write wins, or the API serializes requests.

### 8. Offline First — Works Without Internet

Your vault works when the internet is down. You can read, write, search, and restructure your notes with zero connectivity. When you're back online, `git push` syncs everything.

SaaS PKMs require internet for basic functionality. Notion doesn't load offline. Roam doesn't load offline. Your plain text files always load.

**The test:** Can you work in your PKM on a 10-hour flight with no wifi? With git-based PKM: yes, fully. With SaaS PKM: limited or impossible.

---

## The Counterarguments (And Why They Don't Hold)

**"Notion has better databases and tables."**
Yes. And you can add dataview queries to markdown that replicate 90% of Notion's database functionality. For the remaining 10%, you have an actual database running on your VPS.

**"Git is too technical for most people."**
For setup: yes, there's a learning curve. After setup: `git pull`, `git push`, `git log` — three commands cover 99% of daily use. The One Mind OS sync is automated — you never type a git command to keep your vault in sync.

**"Obsidian Sync is easier than git."**
Easier to start, yes. But you're paying $10/mo for file sync you can replicate for free with git, and you're trusting Obsidian to hold your data. Obsidian Sync is also a binary sync — you can't branch, diff, or audit it.

**"What about images and PDFs?"**
Git LFS (Large File Storage) handles binary files. Or keep your media in a separate folder that syncs via cloud storage. Your knowledge (text) goes in git. Your assets go in cloud storage. This is a feature, not a bug — clean separation of concerns.

---

## The Practical Setup

One Mind OS uses a two-path architecture:

```
Your Mac (Obsidian / VS Code)
    ↓ git push (on save or manual)
GitHub (your-username/your-vault) ← SOURCE OF TRUTH
    ↓ auto-pull every 10 seconds
Your VPS (AI agents read/write here)
```

Three commands that run the entire sync:
- `codex-watch.service` — detects file changes on VPS, auto-commits and pushes
- `codex-pull.timer` — pulls from GitHub every 10 seconds on VPS
- VS Code git integration or Obsidian Git plugin — push from Mac

Cost: $0 for the git infrastructure. Free GitHub repos are unlimited.

---

## The Decision Framework

**Use git-based PKM if:**
- You want to own your data permanently
- You're building an AI-native system where agents read/write your vault
- You think in years, not weeks — you want your knowledge to compound
- You value auditability and version history
- You operate across multiple devices and possibly a VPS

**Use SaaS PKM if:**
- You need rich collaboration with a team that won't set up git
- You need embedded databases with complex views (Notion-style)
- You want zero setup friction and don't mind the tradeoffs

For One Mind OS: git-based is the only architecture that works. Agents need direct filesystem access. Version history of vault changes is essential. Sovereignty is non-negotiable.

---

*One Mind — weareonemind.com*

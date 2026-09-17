---
id: ui-10-20260111053926
title: "GitHub for personal file storage and Git LFS: a technical deep dive"
type: concept
domain: SO
domain: 00-24
folder: "00-24 SO (Sovereign Operations)"
status: active
created: 2026-01-11
modified: 2026-01-11
owner: zeus
tags: [domain/so, folder/10]
todoist_project: "00-24 SO (Sovereign Operations)"
ai_summary: "GitHub for personal file storage and Git LFS: a technical deep dive"
---

# GitHub for personal file storage and Git LFS: a technical deep dive

**GitHub can technically store personal files, but with significant caveats.** While the Terms of Service don't prohibit non-code content, GitHub explicitly states "Git is not designed to serve as a backup tool." The platform imposes a **100 MiB hard file size limit** (50 MiB warning), recommends repos stay under **1 GB** (5 GB strongly recommended, 10 GB hard limit), and the LFS free tier provides only **10 GiB storage and 10 GiB/month bandwidth**. For Obsidian vaults synced via GitHub, mobile sync remains problematic—the official Obsidian Git plugin developer warns mobile support is "very unstable."

## GitHub's ToS permits non-code repos, but discourages large storage

GitHub's Acceptable Use Policies focus on prohibited **behaviors** (malware, spam, harassment) rather than content types. There's no explicit prohibition against storing documents, notes, or PDFs. However, the documentation carries this critical warning:

> "Git is not designed to serve as a backup tool. However, there are many solutions specifically designed for performing backups, such as Arq, Carbonite, and CrashPlan."

The bandwidth policy adds risk for heavy users: "If we determine your bandwidth usage to be significantly excessive in relation to other users of similar features, we reserve the right to suspend your Account." GitHub reserves the right to delete repositories "placing undue strain on infrastructure" after notice.

**Security posture** is enterprise-grade with SOC 2 Type II, ISO 27001/27017/27018 certifications, TLS encryption in transit, and encryption at rest. However, GitHub does **not** offer zero-knowledge encryption—employees could theoretically access your data, and legal requests apply. For truly sensitive documents, local encryption (GPG, age) before upload is recommended.

| Limit Type | Value | Enforcement |
|------------|-------|-------------|
| Individual file (CLI) | **100 MiB** | Hard block |
| Browser upload | **25 MiB** | Hard block |
| Repository size | **1 GB** recommended, **5 GB** strongly recommended | Email warning |
| Push size | **2 GB** | Enforced |
| Max repositories | **100,000** | Warning at 50,000 |

## Git LFS pricing hits harder than expected at scale

Git LFS free tier provides **10 GiB storage and 10 GiB bandwidth/month** for Free/Pro accounts. GitHub Team/Enterprise gets **250 GiB** of each. The metered billing rates beyond free tier are **$0.07/GiB/month** for storage and **$0.0875/GiB** for bandwidth.

A critical gotcha: **bandwidth counts against the repository owner when anyone clones**. Public repositories with LFS files can have their quota depleted by random users. One game studio documented spending **$3,000/month** on LFS bandwidth before building a caching proxy (25GB game × 1,300 CI runs/month = 32.5 TB bandwidth).

| Plan | Free Storage | Free Bandwidth | Max File Size |
|------|-------------|----------------|---------------|
| GitHub Free/Pro | 10 GiB | 10 GiB/month | **2 GB** |
| GitHub Team | 250 GiB | 250 GiB/month | **4 GB** |
| GitHub Enterprise | 250 GiB | 250 GiB/month | **5 GB** |

**Performance degrades significantly with many small LFS files.** Real-world tests show a 5.3 GB repo with 2,266 LFS files takes **17+ minutes** to clone versus ~7 minutes without LFS. Windows checkout operations with 10,000+ LFS files are "horribly slow" due to per-file process spawning. The ideal pattern is fewer large files rather than many small ones.

Track with LFS: `.pdf`, `.psd`, `.ai`, `.mp3`, `.mp4`, `.fbx`, images, archives, binaries. Do NOT track text files—Git handles these efficiently.

## Sparse checkout enables partial clones, but LFS has limitations

Modern Git (2.25+) supports partial clone combined with sparse checkout for working with repository subsets. Microsoft's **Scalar** tool achieved **3-4 second `git status`** times on the 300GB Windows repository using these techniques.

```bash
# Optimal clone for monorepo work
git clone --filter=blob:none --sparse https://github.com/org/repo.git
cd repo
git sparse-checkout init --cone
git sparse-checkout set src/my-component docs

# For LFS repos—critical: skip smudge, then selectively pull
GIT_LFS_SKIP_SMUDGE=1 git clone --filter=blob:none --sparse <repo-url>
git lfs pull --include="src/my-component/**"
```

**Key limitation:** `git lfs pull` does NOT respect sparse-checkout patterns by default—it will download ALL LFS files. You must use `--include` explicitly to limit downloads. This bug (GitHub Issue #3803) remains open.

Google uses **Piper** (custom VCS on Bigtable/Spanner) for their **86 TB monorepo** with 25,000+ engineers. Microsoft moved from VFS for Git to **Scalar** (now the recommended approach). Meta uses **Sapling** (Rust-based, Git-compatible). For most teams, standard Git with partial clone + sparse checkout suffices.

## Obsidian mobile GitHub sync works, but reliability varies dramatically

### iOS: Working Copy is the primary option

**Working Copy** ($25-36 one-time, free for GitHub Student Pack) is the most reliable iOS solution. Setup: create vault with iCloud disabled → clone repo in Working Copy → "Link Repository to Folder" to Obsidian vault → automate via iOS Shortcuts.

Real user experiences are mixed:
- ✅ "Works well for small-to-medium vaults under 2GB"
- ✅ "Manageable and reliable when trained to pull before working and push after"
- ❌ "Phone starts to burn" with large vaults (1GB+) and attachments
- ❌ "Working Copy is just too touchy... Large numbers of new files cause timeouts"
- ❌ One user concludes: "I don't think I can recommend using Obsidian + Git LFS + iOS"

### Android: Termux provides full git functionality

Termux with native git is the standard Android approach:
```bash
pkg install git openssh
termux-setup-storage
git clone <repo> ~/storage/shared/obsidian-vault
```

Automation via Tasker integration or Termux Widget works, but **disable battery optimization** or Android will kill background Termux processes.

### The Obsidian Git plugin is explicitly not recommended for mobile

From the plugin's README: "The Git implementation on mobile is **very unstable**! I would not recommend using this plugin on mobile." It uses isomorphic-git (JavaScript reimplementation) with severe limitations: no SSH (HTTPS-only), no rebase, checking working directory takes **3+ minutes** for large vaults.

**Cost comparison:**
- Obsidian Sync: $4-8/month (most reliable)
- GitHub + Working Copy: ~$25-36 one-time
- GitHub + Termux: Free (technical setup required)
- iCloud: Free but "not reliable, to say the least" on Windows

## Multi-remote Git backup works, except for LFS

Configure multiple push URLs:
```bash
git remote set-url --add --push origin git@github.com:user/repo.git
git remote set-url --add --push origin git@gitlab.com:user/repo.git
git remote set-url --add --push origin git@bitbucket.org:user/repo.git
```

**Critical LFS limitation:** Git LFS does NOT properly support multiple push URLs. From GitHub Issue #1823: "LFS objects are only pushed to one of the two push repos (the last one configured)." You must manually push LFS to each remote:

```bash
git push origin main && git lfs push --all origin main
git push gitlab main && git lfs push --all gitlab main
```

Each platform bills LFS separately—pushing 10GB to GitHub, GitLab, and Bitbucket costs you storage on all three. GitHub and GitLab offer 10GB free; Bitbucket offers only 1GB on free tier.

## GitHub API supports AI agents with 5,000 requests/hour authenticated

| Authentication | Rate Limit |
|----------------|------------|
| Unauthenticated | 60/hour per IP |
| Personal Access Token | **5,000/hour** |
| GitHub App | 5,000-12,500/hour (scales with repos) |
| Enterprise Cloud | **15,000/hour** |
| GITHUB_TOKEN (Actions) | 1,000/hour per repo |

The **GitHub MCP Server** (https://github.com/github/github-mcp-server) provides Model Context Protocol access for AI agents. Configure with:

```json
{
  "mcpServers": {
    "github": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-github"],
      "env": { "GITHUB_PERSONAL_ACCESS_TOKEN": "<TOKEN>" }
    }
  }
}
```

Letta/MemGPT doesn't have built-in GitHub integration—you must create custom tools wrapping the GitHub API. Key considerations: implement exponential backoff, use fine-grained PATs with minimal scopes, cache responses, and monitor `x-ratelimit-remaining` headers.

## Conclusion

GitHub works for personal file storage with caveats: keep repositories under **5 GB**, use LFS for files **over 50 MB**, and expect to pay if exceeding **10 GiB bandwidth/month**. For Obsidian vaults, the mobile experience requires **manual sync discipline** (Working Copy on iOS, Termux on Android)—Obsidian Sync remains the most reliable option for cross-platform users prioritizing convenience over version control. Multi-remote backup is straightforward for regular Git but requires manual intervention for LFS. AI agents can reliably access GitHub via API or MCP server with 5,000 requests/hour on personal access tokens.

---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]

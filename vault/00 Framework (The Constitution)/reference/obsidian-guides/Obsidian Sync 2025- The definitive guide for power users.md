---
id: ui-06-20260111053926
title: "Obsidian Sync 2025: The definitive guide for power users"
type: reference
domain: SO
domain: 00-24
folder: "00-24 SO (Sovereign Operations)"
status: active
created: 2026-01-11
modified: 2026-01-11
owner: zeus
tags: [domain/so, folder/06]
todoist_project: "00-24 SO (Sovereign Operations)"
ai_summary: "Obsidian Sync 2025: The definitive guide for power users"
---

# Obsidian Sync 2025: The definitive guide for power users

**Official Obsidian Sync is the most reliable option for cross-platform mobile sync, but comes with significant storage and performance limitations for large vaults.** While selective sync capabilities exist and work well, mobile performance degrades substantially with vaults over 10GB, and the 100GB storage ceiling with no enterprise expansion options may constrain long-term growth. Self-hosted alternatives offer unlimited storage but sacrifice the mobile reliability you require.

For your specific use case—primarily markdown with some PDFs, requiring rock-solid iOS and Android support plus selective sync—official Obsidian Sync delivers the best experience despite its premium price. However, if your vault exceeds 20GB or contains 10,000+ files, expect mobile performance compromises that no amount of money can currently solve.

## Current pricing and the hard 100GB ceiling

Obsidian overhauled its sync pricing in April 2024, introducing two tiers. The **Standard Plan at $4/month** (annual billing) provides 1GB storage, 1 remote vault, 5MB per-file limit, and 1 month version history. The **Plus Plan at $8/month** (annual) includes 10GB base storage expandable to 100GB for $16/month, 10 remote vaults, 200MB per-file limit, and 12 months version history.

**There are no options beyond 100GB.** Multiple forum discussions confirm users requesting 500GB-1TB tiers receive no path forward. The official FAQ explicitly states "The maximum storage amount is 100 GB." No enterprise plans exist, and the development team hasn't indicated plans to raise this ceiling. Storage is account-wide rather than per-vault, and critically, **version history counts against your total quota**—users report local vaults of 93MB showing as 6.32GB on Sync due to version accumulation.

Educational institutions, faculty, students, and nonprofits receive 40% discounts across all tiers. Legacy users who subscribed before January 2024 were grandfathered into 50GB plans at $8-10/month, but these plans vanish permanently if canceled.

The 200MB per-file limit (doubled from 100MB in January 2024) accommodates most PDFs and documents but excludes large video files. More importantly, **attachments receive only 2 weeks of version history** regardless of plan tier, compared to 1-12 months for markdown notes—a critical limitation if you rely on recovering previous versions of PDFs or images.

## Selective sync works exactly as you need

**Yes, you can exclude your "07-Resources" folder from sync while syncing markdown.** Obsidian Sync's selective sync operates on a per-device basis through Settings → Sync → Selective Sync → Excluded Folders. You configure which folders to exclude independently on each device, allowing a desktop to maintain the full vault while mobile devices exclude large media folders.

The implementation has important behaviors: excluding a folder prevents sending OR receiving changes on that device, but doesn't automatically delete already-synced content. If a folder already exists on a device, you must exclude it on that specific device to prevent future sync operations. When you exclude a folder on mobile, it stops syncing but doesn't remove existing files unless you manually delete them.

By default, Obsidian Sync handles markdown notes and images automatically. PDFs, audio, and video require explicit enablement in settings. You can implement a strategy where markdown syncs universally, images sync to all devices, but PDFs and other attachments sync only to desktop machines by disabling "Sync PDFs" on mobile and/or excluding specific resource folders.

**Limitations exist:** You can only exclude entire folders, not individual files. The system doesn't support regex patterns or glob matching. There's no "include only" mode—just exclusion. Renaming an excluded folder doesn't update the exclusion setting automatically, requiring manual adjustment.

For your use case, the recommended configuration would: enable markdown sync (default), add "07-Resources" to excluded folders on mobile devices, optionally disable "Sync PDFs" globally on mobile while keeping it enabled on desktop, and use images sparingly on mobile or exclude image-heavy folders.

## Attachment handling: straightforward but quota-hungry

All files—notes, images, PDFs, audio, video—count against your storage quota without compression or optimization. The Plus Plan's 200MB per-file limit accommodates most documents, but version history inflates storage dramatically. Since **attachments retain only 2 weeks of version history** (versus 12 months for notes on Plus), rapidly changing media files quickly consume space through versioning.

Users report storage discrepancies that highlight version accumulation. One documented case showed a 93MB local vault consuming 6.32GB on Sync. Another user's two vaults under 250MB locally totaled 1.13GB remotely. These multipliers stem from version history—older versions auto-delete after retention periods expire, but users approaching limits must manually prune through the interface.

Third-party cloud services alongside Obsidian Sync create conflicts. Official documentation explicitly warns against combining Sync with iCloud, Dropbox, or OneDrive, as "files-on-demand" features can be interpreted as deletions, potentially removing content from your remote vault. If you currently use these services, you must fully commit to Obsidian Sync or a self-hosted alternative—hybrid approaches invite data loss.

## Cross-platform performance: desktop excellence, mobile struggles

Desktop platforms (macOS, Windows, Linux) deliver excellent Obsidian Sync performance. Apple Silicon Macs particularly excel, with M1/M2/M3 chips handling even 12,000-file vaults smoothly after initial cache building. Windows performs reliably, though some users report AMD Threadripper systems paradoxically underperforming compared to M1 MacBooks for identical vault operations. Linux works well via AppImage or .deb packages, though **avoid the Flatpak version**—sandboxing restrictions break sync features.

**Mobile platforms reveal Obsidian Sync's critical weakness.** iOS suffers the most severe issues. Multiple detailed user tests document the problem: a 13GB vault with 1,200 markdown files plus images required over 1 minute to load on iPhone every time the app opened, with attachments causing complete freezes and crashes. Even text-only vaults took 44 seconds to load. One user with 21,000 files reported 45-second load times on iOS. Testing from October 2024 concluded: "Obsidian apps do not seem to be ready for vaults with a large amount of notes or large amount of attachments."

Android fares better but shares challenges. Initial sync for large vaults takes hours. Users report vaults requiring "about a minute or more to load" with "one-time indexing" every app launch. The fundamental issue: **mobile apps reload the entire vault on every open**—there's no incremental loading. One power user with 50,000 notes on desktop syncs only 5,000 to mobile due to performance, noting that even this "introduces friction" for quick capture.

Recent updates provide incremental improvements. Version 1.10.1 in October 2024 made mobile launch "significantly faster" by opening the last file editor while the rest loads in background. Version 1.7.1 improved app load times and responsiveness with iCloud vaults. These optimizations help but don't solve the fundamental architecture—mobile devices still load the full vault structure.

Battery drain compounds mobile issues. iOS users consistently report Obsidian as their "most battery hungry app," exceeding Safari and Firefox. Surface Pro 8 users note high battery usage with Sync enabled on Windows. The filesystem watching required for real-time sync causes continuous background activity.

**Settings and plugin sync has persistent reliability issues across all platforms.** Users report "1-2 years" of plugins and hotkeys failing to sync reliably. Settings upload to the remote vault but don't automatically apply—they require Obsidian restarts or manual intervention. Community plugins particularly suffer; you must manually enable both "Active community plugin list" and "Installed community plugin list" as they're not synced by default. Sync settings themselves never sync between devices, requiring configuration on each machine.

Platform-specific quirks include: iOS restricts vaults to a specific "Obsidian" folder in iCloud Drive due to sandboxing; Linux has a 255-character filename limit causing ENAMETOOLONG errors in dual-boot scenarios where Windows allows longer paths; macOS users experience iCloud Drive delays and battery drain from aggressive file watching.

## Version history: robust for notes, minimal for media

The Plus Plan's 12-month version history for markdown notes provides excellent protection, with unlimited versions within that timeframe—only time-based limits exist, not version counts. You can view diffs between revisions and restore previous versions through the interface easily.

**The 2-week limit for attachments is the critical constraint.** If you modify PDFs frequently or iterate on visual assets, you have only 14 days to recover older versions regardless of plan tier. This asymmetry between note and attachment retention reflects a design choice prioritizing text version control over media versioning.

After subscription expiration, Obsidian retains data for 1 month. Renewing within this window preserves everything, including version history. After 30 days, remote vault data deletes permanently—though local vault data on your devices remains untouched, providing recovery insurance.

Version history contributes to the storage multiplication effect mentioned earlier. Since all versions count against quota, actively edited vaults accumulate history rapidly. Manual pruning becomes available when approaching limits, but no automatic cleanup exists before hitting this threshold.

## Obsidian Publish: completely separate, no integration

Obsidian Publish and Obsidian Sync are **entirely independent services with zero integration.** Publish creates public websites from selected notes ($8-10/month), while Sync handles cross-device synchronization. Neither requires the other; purchasing both provides no bundle discount or enhanced features.

You might use both if managing a published site from multiple devices (Sync keeps your vault consistent while Publish pushes selected notes public) or running a team that collaborates privately via Sync while publishing a subset publicly. However, these are complementary workflows, not integrated features—the services don't communicate.

For your use case focused on personal vault synchronization, Obsidian Publish is irrelevant unless you plan to create a public knowledge base or digital garden.

## Recent updates strengthen security, barely improve mobile

August 2024 brought a minor security upgrade with AES-SIV encryption for file paths and hashes, replacing the less secure IV-derived method. File contents remain encrypted with AES-GCM. A migration assistant in v1.9.11 allows existing vaults to upgrade encryption, appearing as "Upgrade vault encryption" in Settings → Sync.

September-October 2024 updates focused on mobile UX refinements: faster app launch by loading the last file editor while background processes run, improved responsiveness with iCloud vaults, and fixes for iOS adaptive themes and Android toolbar positioning. October's v1.10.1 represents the most significant mobile improvement in months, though it addresses launch speed rather than vault size handling.

The January 2024 pricing overhaul introduced flexible account-wide storage allocation (versus per-vault limits), increased file size limits to 200MB, and expanded vault counts to 10. Regional server selection (North America, Europe, Asia, Oceania) became available, though you can only choose region during initial vault setup—changing requires recreating the remote vault and losing version history.

**No updates address the fundamental mobile architecture issues** with large vaults. The community continues requesting solutions for 10,000+ note performance, better incremental loading, and improved background sync, but these remain unsolved as of November 2025.

## Community reality: desktop praised, mobile criticized

Community feedback bifurcates sharply by vault size and platform. Users with under 5GB and under 5,000 notes overwhelmingly praise Obsidian Sync as "seamless," "unobtrusive," and working "like a charm." A MakeUseOf review from August 2024 noted: "Unlike my previous solution with Google Drive, I don't have much more to say about Obsidian Sync—and that's precisely why it's great. It quietly works in the background without drawing attention to itself."

**Power users with 50GB+ vaults tell a different story.** Documented experiences include:

A detailed January 2023 test of a 13GB vault with 1,200 notes found desktop performance acceptable but iPhone load times over 1 minute, with attachments causing complete unusability. A 2024 user with 21,000 files experienced 45-second iOS load times after enabling iCloud Advanced Data Protection. An 800MB vault with 2,400 files caused mobile app freezing and crashes until the user created a fresh vault with careful plugin management.

Sync conflicts occur infrequently but cause significant issues when they happen. Obsidian uses the diff-match-patch algorithm for automatic conflict resolution, attempting to merge changes. However, users report the algorithm sometimes "overwrote newer data for old" or "duplicates both old and new data within the same note." One user documented their iPhone overwriting newer data when opening the app without simultaneous desktop sync, leaving them "terrified of losing data."

The community requests manual conflict resolution (like Git) rather than automatic merging, but this remains unavailable. Version history provides recovery options, but users must monitor sync logs and check conflicts reactively rather than receiving clear warnings.

Common complaints cluster around: high pricing ($8-10/month seen as "absolutely absurd" by some Redditors comparing to $1 for 50GB iCloud storage), persistent settings/plugin sync failures requiring manual device configuration, slow mobile sync even for small changes, and storage quota inflation from version history without automatic pruning.

The gotcha most frequently surprises users: regional server lock-in. Once you select a server region during initial vault setup, changing requires creating a new remote vault and permanently losing all version history. Users discover this only when experiencing slow international sync and seeking region changes.

## Self-hosted alternatives: unlimited storage, limited mobile reliability

For your requirement of "rock-solid mobile support (iOS and Android)," self-hosted alternatives fall short despite advantages in storage capacity, privacy, and cost.

**Self-Sync/LiveSync** plugin emerges as the most popular self-hosted option, using CouchDB or object storage (MinIO, S3, Cloudflare R2) as backend. It offers end-to-end encryption, active development through 2025, and experimental WebRTC peer-to-peer sync without servers. However, iOS now requires HTTPS connections due to recent Apple policy changes—HTTP fails with "TypeError: Load failed." Setup demands Docker, CouchDB, reverse proxy, certificates, and CORS configuration. Users report configuration taking 20+ minutes to propagate to mobile devices.

Android experiences "Not Ready" errors where devices sync once then fail on subsequent opens, one-way sync problems (download works but upload fails), and confusing non-intuitive configuration UI. GitHub issues reflect user frustration: "Why make it so complex?" with complaints about lack of documentation. Users report sync randomly stopping, requiring manual database rebuilds with unclear recovery options. The setup complexity is "unfathomable" to non-technical users.

**Git-based solutions** work excellently on desktop with automatic commits and full version history but fail on mobile. The Obsidian Git plugin is officially labeled "unreliable on iOS" by developers. iOS users require external paid apps like Working Copy ($20) or complex shortcut automations. Android requires Termux (only via F-Droid, not Play Store) and command-line git operations. Merge conflicts on mobile are "pretty much unsolvable." Multiple sources recommend using Git as backup only, not primary sync: "Use LiveSync/Remotely Save for sync + Git for backup."

**Syncthing** provides the best self-hosted experience but only for Android. The peer-to-peer synchronization works "excellently" on Android with Syncthing-Fork from F-Droid, offering always-on background sync and reliable plugin/theme synchronization. However, battery drain from continuous background operation concerns users—one reported 10% battery consumption in 10 hours without active sync. Syncthing creates .sync-conflict files rather than auto-merging, requiring manual cleanup. Most critically, **Syncthing does not work on iOS**—iOS sandboxing prevents background folder sync, making it Android-only.

Community experiences from 2024-2025 show a pattern: users attempt self-hosting, encounter mobile reliability issues, then switch to official Obsidian Sync. One user summarized: "Such is the life of maintaining your own tools; if I just used iCloud this would just work." Technical users with always-on home servers, Tailscale VPNs, and willingness to troubleshoot report success with LiveSync, but this represents a minority willing to invest 10-20 hours initial setup plus 1-2 hours monthly maintenance.

The mobile reliability hierarchy for self-hosted solutions ranks: Syncthing on Android (7/10 reliability), LiveSync (5/10 - works when working but requires frequent troubleshooting), Remotely Save with cloud storage (6/10 - simpler but periodic sync only), and Git-based (3/10 - not recommended for mobile).

## The verdict: official Sync wins for your use case, with caveats

Given your requirements—rock-solid iOS and Android support, version control, selective sync, markdown plus some PDFs—**official Obsidian Sync at $8-16/month (Plus Plan) is the clear choice.** Self-hosted alternatives' mobile reliability problems, especially on iOS, disqualify them for your "rock-solid" requirement. Cost isn't your concern, and Obsidian Sync eliminates the 10-20 hour setup plus ongoing maintenance burden of self-hosting.

**However, critical caveats apply:**

If your vault grows beyond 20GB or 10,000 files, mobile performance will degrade substantially regardless of hardware. The architecture doesn't support incremental loading—mobile devices load the entire vault on every open. Users with 50,000 notes maintain separate minimal mobile vaults (under 5GB, under 5,000 files) for quick capture while keeping the full vault desktop-only. The 100GB storage ceiling with no expansion path may eventually constrain growth.

The 2-week version history for PDFs and attachments (versus 12 months for markdown) limits recovery windows for media files. Plugin and settings sync remains unreliable across platforms, requiring manual configuration on each device. You'll likely need to configure your plugin setup independently on each mobile device and desktop machine.

**Recommended implementation strategy:**

Start with the Plus Plan ($8/month annual, 10GB base). Enable selective sync immediately, excluding your "07-Resources" folder from mobile devices while syncing it to desktop. Monitor storage consumption through Settings → Sync for the first month—version history will inflate actual usage 2-6x beyond local vault size. If approaching 10GB, expand to 100GB ($16/month) or aggressively prune version history.

Maintain separate .obsidian configurations for mobile versus desktop if needed, disabling heavy plugins like Dataview on mobile to improve performance. Use Git on desktop as a backup (don't sync the .git folder) to provide version control insurance beyond Obsidian Sync's 12-month limit. Test mobile performance extensively before committing your full vault—if you experience 30+ second load times, consider a two-vault strategy (mobile minimal vault + desktop full vault).

Alternative approach only if iOS isn't essential: Syncthing on Android provides excellent reliability for free, with Obsidian Sync handling iOS devices only. This hybrid reduces cost ($4/month Standard Plan) while maintaining Android reliability. However, this adds complexity managing two sync systems.

The disappointing reality: **no solution currently provides rock-solid mobile performance for 50GB+ multimedia vaults in late 2025.** Obsidian Sync offers the best available mobile experience but hits architectural limits around 10-20GB. If your vault will grow beyond this with significant media files, accept mobile as a constrained quick-capture interface rather than full vault access, or wait for future Obsidian updates addressing mobile architecture—though no such updates have been announced.

For now, official Obsidian Sync delivers the most reliable solution for your stated requirements, with selective sync working exactly as you need and cross-platform support that self-hosted alternatives can't match on iOS. Just understand the mobile performance ceiling and plan your vault structure accordingly.

---

## Related

- [[00-24 SO (Sovereign Operations)/00 Framework (Doctrine) (The Constitution)/00 Framework|00 Framework]]

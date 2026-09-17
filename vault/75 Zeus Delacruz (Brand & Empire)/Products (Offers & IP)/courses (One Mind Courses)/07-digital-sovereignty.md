# Course 7: DIGITAL SOVEREIGNTY

> **Your life runs on software. Own it, secure it, control it.**

**Tier:** Legacy ($97/mo)
**Circle Space:** 🤖 AI & Agents (ID: 2617142)
**Prerequisite:** Course 2 — VPS live, SSH access working
**Target Student:** Someone running a production AI stack with real data, real credentials, and real exposure — who has never done a security audit of their own digital life
**Outcome:** A hardened, sovereignty-first digital infrastructure — secrets in 1Password, network on Tailscale, VPS secured, agent stack access-controlled, incident response plan ready

---

## What Digital Sovereignty Means

Digital sovereignty is not paranoia — it's the principle that **you own and control your digital infrastructure**, rather than being owned and controlled by it.

Three levels:
1. **Security** — your system cannot be compromised by external actors
2. **Privacy** — your data stays yours, not harvested by platforms
3. **Resilience** — your system survives failure, breach, or loss without catastrophe

This course teaches all three. In that order.

---

## What You'll Build

By the end of this course:
- All credentials out of files, browsers, and notes — into 1Password with proper vault structure
- Tailscale running as your private network — every device connected, zero public exposure where possible
- VPS hardened: firewall, fail2ban, SSH key-only, no root login
- Agent stack access-controlled: each agent has only the permissions it needs
- A personal threat model that tells you exactly what you're defending and from what
- An incident response plan: what you do in the first 30 minutes of a breach
- A secrets rotation schedule that runs automatically

---

## Course Structure — 7 Modules, 40 Lessons

---

### MODULE 1: THREAT MODELING YOUR DIGITAL LIFE

**Lesson 1.1 — What Is Threat Modeling?**
Security is about tradeoffs. Threat modeling helps you make them intelligently.
- The security mindset: thinking like an attacker to defend like a professional
- The 4 questions every threat model answers: what are you protecting? From whom? How likely? What's the cost?
- Personal threat modeling vs enterprise: your risks are different, your defenses should be too
- The STRIDE model simplified: Spoofing, Tampering, Repudiation, Information Disclosure, Denial of Service, Elevation of Privilege
- Over-engineering the threat model: the trap of defending against nation-states when your risk is account takeover

**Lesson 1.2 — Your Attack Surface**
Everything that can be compromised is attack surface. Map yours.
- Attack surface categories: credentials, devices, accounts, services, code, agents
- The credential inventory: every username/password you have is a potential entry point
- Device exposure: Mac, iPhone, VPS — each is an attack surface
- Account exposure: Google, GitHub, Discord, Circle, Postiz — each can be targeted
- The AI agent surface: your agents have API keys, vault access, shell access — this is new attack surface most people ignore
- The most likely attack vectors for someone like you: credential stuffing, phishing, supply chain, compromised dependency

**Lesson 1.3 — Risk Prioritization**
Not all risks are equal. Spend your security effort where it matters most.
- The risk matrix: likelihood × impact = priority
- High likelihood, high impact: credential theft, account takeover → address immediately
- Low likelihood, high impact: VPS root compromise, vault exfiltration → address second
- High likelihood, low impact: spam, bot traffic → mitigate passively
- Low likelihood, low impact: don't spend time here
- Your top 5 risks right now: an honest audit of where you're most exposed
- The 80/20 of personal security: these 5 things cover 80% of real risk

**Lesson 1.4 — The Sovereignty Principles**
The non-negotiables that guide every security decision.
- Principle 1: Assume breach — design as if it will happen, not if
- Principle 2: Least privilege — every account, agent, and key gets only what it needs
- Principle 3: Defense in depth — multiple layers, no single point of failure
- Principle 4: Sovereignty first — your data on your hardware, not SaaS databases when avoidable
- Principle 5: Audit trail — everything logged, everything reviewable
- Applying these to One Mind OS: how each principle maps to your current stack

---

### MODULE 2: SECRETS MANAGEMENT WITH 1PASSWORD

**Lesson 2.1 — Why Secrets in Files Is Never Acceptable**
Every API key in a file, `.env`, or config is a breach waiting to happen.
- The most common breach pattern: API key committed to git (it's scanned automatically)
- Even private repos get leaked: former employees, supply chain, account compromise
- Environment variables are not safe either: accessible to any process that runs
- The rule: no secret ever touches a file — ever. Reference it, don't store it.
- What counts as a secret: API keys, tokens, passwords, private SSH keys, webhook URLs, database credentials

**Lesson 2.2 — 1Password Architecture for One Mind OS**
How to structure 1Password so every secret is findable and properly scoped.
- The vault structure: one vault per pillar (SO, HP, LE, GE)
- The SO vault: all infrastructure secrets — VPS, GitHub, agents, APIs
- The reference format: `op://vault-name/item-name/field-name`
- Item naming convention: consistent, searchable, unambiguous
- Tags and categories: organizing by service type, not just pillar
- Shared vaults vs personal vaults: when to use each
- The SO vault ID: `aaeyris5z44pdqspk277u2oppi` — your production vault

**Lesson 2.3 — 1Password CLI Integration**
`op` CLI lets scripts and agents access secrets without storing them.
- Installing the `op` CLI: Mac and VPS
- Authentication: `op signin` — how the session works
- Fetching a secret: `op read "op://SO/github-token/credential"`
- Using in shell scripts: `API_KEY=$(op read "op://SO/api-key/credential")`
- Using in configs: referencing `op://` paths in OpenClaw config
- The service account: a non-human 1Password account for automated access (VPS, agents)
- Audit log: every `op read` is logged — you know what accessed what

**Lesson 2.4 — Secrets Rotation**
Credentials don't expire on their own. You have to rotate them.
- Why rotation matters: limits the window of damage if a credential is compromised
- The rotation schedule: critical secrets (30 days), important (90 days), standard (180 days)
- High-priority rotation targets: GitHub tokens, VPS SSH keys, bot tokens, webhook secrets
- The rotation workflow: generate new → update in 1Password → update in services → verify → revoke old
- Agent-assisted rotation: Guardian tracks rotation due dates and prompts you
- The rotation audit: checking that nothing still references a rotated credential
- Hands-on: rotate one credential end-to-end using this workflow

**Lesson 2.5 — Emergency Access and Recovery**
What happens to your system if you lose access to 1Password?
- The Emergency Kit: 1Password's account recovery document — where it lives, who has it
- The trusted contact: someone who can help you recover in a catastrophe
- The backup decryption key: local encrypted copy of critical secrets
- Recovery test: simulate losing your phone and recovering 1Password access
- The "hit by a bus" principle: your system should survive your absence
- What to document for recovery: account IDs, recovery codes, backup contacts

---

### MODULE 3: TAILSCALE — YOUR PRIVATE NETWORK

**Lesson 3.1 — Why Tailscale Changes Everything**
Tailscale is a zero-config, encrypted mesh VPN built on WireGuard. Once you use it, you can't go back.
- The problem it solves: exposing services to the public internet is dangerous and unnecessary
- What Tailscale does: every device on your Tailnet gets a private IP, traffic stays encrypted
- Your Tailnet: `tail717ba0.ts.net` — the private network that connects your devices
- The three nodes: `zeus-macbook-pro` (100.114.142.55), `onemindos-core` (100.98.25.10), `zeus-iphone-17` (100.122.203.68)
- Accessing services privately: OpenClaw dashboard, n8n, Postiz — all via Tailscale, not public internet
- Tailscale vs traditional VPN: no server to manage, no configuration, works everywhere

**Lesson 3.2 — Tailscale Setup and Device Enrollment**
Every device in your life should be on your Tailnet.
- Installing Tailscale: Mac, VPS (Linux), iPhone
- Enrolling a device: `tailscale up` — the 30-second setup
- Verifying the connection: `tailscale status` — see all nodes
- Tailscale DNS: MagicDNS resolves `onemindos-core` to its Tailscale IP automatically
- Device approval: enabling device approval so you authorize each addition
- Key expiry: disabling key expiry for server nodes (VPS should never need reauthentication)
- The Tailnet admin console: managing devices, access rules, DNS

**Lesson 3.3 — Tailscale Serve and Funnel**
Exposing services securely — only to your Tailnet (Serve) or publicly (Funnel).
- `tailscale serve`: expose a local service to your Tailnet only — no public internet exposure
- `tailscale funnel`: expose a service publicly via Tailscale's edge — with your control
- Your current setup: OpenClaw dashboard on `omos-core.tail717ba0.ts.net` via Serve
- n8n on Titan via Serve: `omos-titan.tail717ba0.ts.net`
- When to use Funnel vs a public reverse proxy: Funnel for temporary public exposure, reverse proxy for permanent
- The serve config: `tailscale serve https:443 / proxy http://localhost:PORT`
- Hands-on: expose one new service via Tailscale Serve

**Lesson 3.4 — Access Control with ACLs**
Tailscale ACLs let you define exactly who can reach what on your Tailnet.
- The default ACL: everyone on your Tailnet can reach everything — fine for personal use, not for shared networks
- Writing ACL rules: allow `src` to reach `dst` on `port`
- Tagging devices: `tag:servers` for VPS nodes, `tag:clients` for personal devices
- The principle of least privilege in Tailscale: your iPhone shouldn't be able to SSH to the VPS directly
- Service-level restrictions: only `tag:clients` can reach port 18789 (OpenClaw)
- Testing ACLs: `tailscale ping` and direct service tests after each rule change

**Lesson 3.5 — Tailscale SSH**
Replace traditional SSH with Tailscale SSH — no keys to manage, access control built in.
- How Tailscale SSH works: Tailscale authenticates the connection, no separate SSH key required
- Enabling Tailscale SSH on your VPS: `tailscale up --ssh`
- Connecting: `ssh onemindos-core` without specifying keys
- ACL-based access: define exactly which devices can SSH to which nodes
- Session recording: log SSH sessions for audit trail
- When to keep traditional SSH keys: compatibility with tools that don't support Tailscale SSH

---

### MODULE 4: VPS HARDENING

**Lesson 4.1 — The Hardening Checklist**
A production VPS needs to be locked down before it runs your life.
- What "hardened" means: minimal attack surface, monitored access, resistant to common attacks
- The hardening order: users → SSH → firewall → services → monitoring
- Current hardening state of `onemindos-core`: what was done in April 2026
- What a fresh VPS looks like before hardening: everything open, root login enabled, no firewall
- The risk of skipping hardening: internet scanners find open ports within minutes of provisioning

**Lesson 4.2 — User and SSH Configuration**
The first line of defense is controlling who can access the system at all.
- Create a non-root sudo user: never run services as root
- Disable root SSH login: `PermitRootLogin no` in `sshd_config`
- SSH key-only authentication: `PasswordAuthentication no` — password brute force becomes impossible
- The authorized_keys file: which keys have access, how to audit and remove
- SSH port change (optional): moving from 22 to a non-standard port (security through obscurity — limited value)
- The SSH hardening config: `AllowUsers`, `MaxAuthTries`, `ClientAliveInterval`
- Hands-on: audit your current SSH config against this checklist

**Lesson 4.3 — Firewall with UFW**
A firewall defines exactly what traffic is allowed. Everything else is dropped.
- UFW (Uncomplicated Firewall): the right tool for a personal VPS
- The default policy: deny incoming, allow outgoing — start here
- Allow only what you need: SSH (22 or custom), Tailscale (41641/udp), OpenClaw (18789 — Tailscale only)
- Never expose to public: databases, n8n, agent dashboards — these go behind Tailscale
- The firewall audit: `ufw status verbose` — review every open port
- DigitalOcean Cloud Firewall: the second layer at the network level (before traffic reaches your VPS)
- Hands-on: audit current UFW rules, close anything that shouldn't be public

**Lesson 4.4 — Fail2ban: Automatic Attack Response**
Fail2ban watches your logs and bans IPs that show attack patterns.
- What fail2ban does: detects repeated failed logins, bans the source IP automatically
- Installing and configuring fail2ban on Ubuntu
- The SSH jail: ban IPs after 3 failed SSH attempts for 24 hours
- The OpenClaw jail: protect the gateway from brute force attempts
- Reviewing bans: `fail2ban-client status sshd` — see what's being blocked
- The ban log: how many attacks are being stopped automatically
- Whitelist your own IPs: never ban yourself (add Tailscale IPs to ignore list)

**Lesson 4.5 — Monitoring and Alerting**
Know when something is wrong before it becomes a crisis.
- DigitalOcean monitoring: CPU, RAM, disk, network — built-in alerts
- Alert thresholds: CPU > 80%, disk > 85%, RAM > 90%
- Log monitoring: `journalctl` and what to watch for
- The audit log: who logged in, when, from where — `last` and `who`
- Automated security scan: weekly cron that checks for open ports and unusual processes
- The Guardian agent: how Guardian monitors your security posture and alerts to Discord
- The incident timeline: when an alert fires, what's the first 10 steps?

---

### MODULE 5: AGENT STACK SECURITY

**Lesson 5.1 — The New Attack Surface: AI Agents**
Your agents have API keys, shell access, vault write access, and channel bindings. This is a real attack surface.
- What an agent compromise looks like: attacker gets your bot token → DMs the agent → agent executes commands
- Prompt injection: malicious content in the vault or from a channel that tricks the agent into acting
- The supply chain risk: a clawhub skill that has malicious instructions
- Credential exposure via agents: agent logs that contain API keys
- The principle of least privilege applied to agents: each agent does exactly what it needs to, nothing more

**Lesson 5.2 — Agent Access Control**
Locking down what each agent can access.
- The `allowedTools` config: explicit tool allowlist per agent — if it's not listed, the agent can't use it
- File system boundaries: agents should only write to their designated vault paths
- The domain isolation rule: Spartan cannot write to GE domains, Empire cannot touch HP
- API key scoping: each agent gets its own API key with minimum permissions
- Bot token isolation: each Discord bot is its own bot with its own token
- The audit question: "If this agent was compromised, what could an attacker do?"

**Lesson 5.3 — Prompt Injection Defense**
Prompt injection is the SQL injection of the AI era.
- What prompt injection is: malicious instructions embedded in content the agent reads
- The vault injection risk: someone puts instructions in a vault file an agent will read
- The channel injection risk: user sends crafted message that overrides agent behavior
- Defense principles: never trust external content as instructions, validate source
- The safe reading pattern: agents summarize external content rather than execute it
- Input sanitization in channel bindings: filter patterns before they reach the agent
- Detecting injection attempts: what they look like, how to log and alert on them

**Lesson 5.4 — Secrets in the Agent Stack**
Agents need API keys. Here's how to give them access without exposing credentials.
- The wrong way: hardcoding keys in `openclaw.json` (they end up in logs and git)
- The right way: `op://` references that get resolved at runtime
- The 1Password service account on VPS: automated access without human authentication
- Environment variable injection: how OpenClaw resolves `op://` paths before agents start
- What gets logged: ensure no secret values appear in gateway logs
- The log audit: scan your current logs for credential patterns
- Hands-on: move all credentials from openclaw.json to 1Password references

**Lesson 5.5 — The Clawhub Security Audit**
Every skill you install is code that runs in your agent. Audit before you trust.
- What skills can do: execute shell commands, call external APIs, read/write files
- The skill review checklist: read the SKILL.md in full, check for shell execution, verify API calls
- Trusted vs untrusted skills: verified clawhub skills vs community-contributed
- The sandbox principle: test new skills in isolation before enabling in production
- Keeping skills updated: `clawhub update --all` — outdated skills may have patched vulnerabilities
- Removing unused skills: attack surface reduction — if you don't use it, remove it
- Hands-on: audit every installed skill against the security checklist

---

### MODULE 6: INCIDENT RESPONSE

**Lesson 6.1 — The Incident Response Mindset**
Breaches happen. Your response in the first 30 minutes determines the outcome.
- The two failure modes: panic (makes it worse) vs denial (also makes it worse)
- The IRP mindset: calm, systematic, document everything
- Incident categories: account compromise, credential leak, VPS breach, data exfiltration, agent hijack
- The first rule: don't make changes until you understand what happened
- The incident log: start writing immediately — timestamp everything
- The containment principle: stop the bleeding before you start diagnosing

**Lesson 6.2 — Your Incident Response Plan**
A written plan you can follow when your hands are shaking.
- The One Mind OS IRP: 6 steps that cover any incident category
- Step 1: Detect — how you know something happened (alert, discovery, external notification)
- Step 2: Contain — isolate the affected system/account immediately
- Step 3: Assess — what was accessed? What was compromised? What's the scope?
- Step 4: Eradicate — remove the attacker's access and the vulnerability they used
- Step 5: Recover — restore normal operation from a known-clean state
- Step 6: Learn — post-incident review, update defenses, document what happened
- Hands-on: write your personal IRP using this template — do it before you need it

**Lesson 6.3 — Account Compromise Response**
Someone has your credentials. Here's what you do in the first 30 minutes.
- Immediate actions: change the password, revoke all sessions, check login history
- Token and API key revocation: GitHub, Discord, OpenClaw, any connected service
- Check for persistence: did they create a secondary access method?
- Notify relevant parties: GitHub, Discord, 1Password if they had master access
- The "what did they see?" audit: review logs to understand what was accessed
- Post-compromise hardening: why this happened and how to prevent recurrence

**Lesson 6.4 — VPS Breach Response**
Your server is compromised. This is the most serious incident on your stack.
- Detection signs: unusual processes, unexpected outbound connections, unfamiliar cron jobs
- Immediate containment: snapshot the VPS (evidence), then lock down — revoke SSH keys
- The nuclear option: if in doubt, destroy and rebuild from backup
- Evidence preservation: before you change anything, capture logs and running process list
- The rebuild procedure: fresh VPS from snapshot, rotate all credentials, verify before re-enabling services
- Why regular backups make this recoverable: a clean snapshot is the difference between 2 hours and 2 days

**Lesson 6.5 — The Security Review Cadence**
Security is not a one-time setup. It's an ongoing practice.
- Weekly: check fail2ban logs, review Guardian's Discord reports, verify all services running
- Monthly: rotate scheduled credentials, audit active API keys, review 1Password vault
- Quarterly: full security audit against the hardening checklist, update threat model
- Annual: rotate all SSH keys, review all service accounts, assess new attack surface
- Guardian as your security agent: automated weekly scans with Discord alerts
- The security debt audit: what you've been meaning to secure but haven't — work through it systematically

---

### MODULE 7: DATA SOVEREIGNTY AND PRIVACY

**Lesson 7.1 — Owning Your Data**
The vault is yours. The infrastructure is yours. Here's how to keep it that way.
- The sovereignty audit: where does your data actually live? Map every service.
- Data categories: personal (HP/LE), operational (SO), business (GE) — each has different exposure risk
- The local-first principle: prefer data on your hardware over data in SaaS
- The git-based PKM advantage: your knowledge in plain text files you control
- Cloud backup without cloud dependency: encrypted offsite backup that you own the keys to
- The data minimization principle: don't store data you don't need — less data = less exposure

**Lesson 7.2 — Privacy in the AI Stack**
Your agents process sensitive information. Here's how to keep it private.
- What your agents see: vault content, tasks, calendar, messages — potentially very personal data
- The model API boundary: what Anthropic/OpenAI receives when your agent calls their API
- The GitHub Copilot model: what Microsoft receives from your requests
- Privacy-preserving practices: don't send sensitive financial/medical data to external APIs
- The local model option: Ollama on your VPS — inference stays on your hardware
- Data retention policies: how long do the cloud providers keep your requests?
- The informed choice: knowing the tradeoffs, making intentional decisions

**Lesson 7.3 — Encrypted Backups**
Your backup is useless if it's unencrypted or inaccessible.
- What to back up: vault (git handles this), agent configs, 1Password Emergency Kit, SSH keys
- Encryption before backup: `gpg` symmetric encryption for sensitive archives
- The backup destinations: DigitalOcean Spaces, external drive, iCloud (encrypted)
- The 3-2-1 rule: 3 copies, 2 different media, 1 offsite
- Backup testing: restore from backup quarterly — untested backups are not backups
- The backup inventory: a document that lists what's backed up, where, and when last tested

**Lesson 7.4 — Graduation: Your Sovereign Digital Infrastructure**
You've built a security posture worthy of a professional operator.
- The Digital Sovereignty Scorecard: rate yourself across all 7 modules
- The security one-pager: your personal threat model, key defenses, incident contacts
- What you've built: 1Password architecture, Tailscale private network, hardened VPS, access-controlled agents, incident plan, encrypted backups
- The ongoing commitment: security is a practice, not a project
- The compounding security advantage: each improvement reduces your attack surface and builds the habit
- Graduation challenge: run the full security audit and score at least 80% before completing the course

---

## Resources Included

| Resource | Format | Access |
|----------|--------|--------|
| Personal threat model template | Markdown | Legacy+ |
| 1Password vault structure guide | PDF | Legacy+ |
| VPS hardening checklist (full) | Markdown | Legacy+ |
| UFW rules template | Bash | Legacy+ |
| fail2ban config template | Config | Legacy+ |
| Agent access control matrix | Spreadsheet | Legacy+ |
| Incident response plan template | Markdown | Legacy+ |
| Security review SOP | Markdown | Legacy+ |
| Clawhub skill audit checklist | Markdown | Legacy+ |
| Guardian agent weekly scan prompt | Markdown | Legacy+ |
| Encrypted backup script | Bash | Legacy+ |

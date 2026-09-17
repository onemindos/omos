# OpenClaw Security — CVEs & Hardening

> Extracted from [awesome-openclaw](https://github.com/rohitg00/awesome-openclaw) — March 2026

## Known CVEs (All Patched)

| CVE | Severity | Description | Fixed In |
|-----|----------|-------------|----------|
| **CVE-2026-24763** | High | Docker PATH Command Injection | v2026.1.29 |
| **GHSA-g8p2-7wf7-98mq** | Critical (8.8) | gatewayUrl Token Exfiltration (1-click RCE) | v2026.1.29 |
| **GHSA-q284-4pvr-m585** | High | sshNodeCommand Injection | v2026.1.29 |
| **CVE-2026-25593** | Critical | Unauthenticated Local RCE via WebSocket | v2026.1.30 |
| **CVE-2026-25475** | High | Local File Inclusion via MEDIA: Path | v2026.1.30 |
| **CVE-2026-25253** | Critical (8.8) | 1-Click RCE via Cross-Site WebSocket Hijacking | v2026.1.29 |

> **Update immediately to v2026.2.12+** to patch all known vulnerabilities (40+ security fixes).

## Security Best Practices

1. **Bind to `127.0.0.1` only** (never `0.0.0.0`) — over 135,000 exposed instances found online
2. Use **SSH tunnels** or [Tailscale Serve](https://tailscale.com/kb/1242/tailscale-serve) for remote access
3. Use **auth profiles** (system keychain): `openclaw auth set ANTHROPIC_API_KEY`
4. Configure **firewall**: `ufw default deny incoming && ufw allow ssh && ufw allow 443/tcp && ufw enable`
5. Use **Docker** for container isolation
6. Run regular audits: `openclaw security audit --deep` and `openclaw doctor`
7. Use **Composio** for managed OAuth ([security guide](https://composio.dev/blog/secure-openclaw-moltbot-clawdbot-setup))

## Skills Marketplace Risks

> **~15% of community skills contain malicious instructions**

- **341 malicious skills** discovered in "ClawHavoc" campaign (Atomic Stealer malware)
- **280+ skills leak API keys and PII** ([Snyk](https://snyk.io/blog/openclaw-skills-credential-leaks-research/))
- VirusTotal scanning now integrated
- Always audit skills before installation

## Security Tools

| Tool | Description | Source |
|------|-------------|--------|
| **ClawSec** | Complete security skill suite by Prompt Security | [GitHub](https://github.com/prompt-security/clawsec) |
| **ClawBands** | Security middleware — intercepts tool execution, human-in-the-loop approval | [GitHub](https://github.com/SeyZ/clawbands) |
| **ClawGuard** | Permission manifests, runtime enforcement, sandboxing, audit logging | [GitHub](https://github.com/newtro/ClawGuard) |
| **Aquaman** | Credential isolation proxy — API keys never enter the agent process. Stores in Keychain/1Password/Vault | [GitHub](https://github.com/tech4242/aquaman) |
| **Clawhatch** | Pre-install security scanner — 128 checks, scores skills 0-100 | [GitHub](https://github.com/AISafetyLab/clawhatch) |
| **Claw-Hunter** | MDM-ready scripts to detect shadow OpenClaw agents across endpoints | [GitHub](https://github.com/backslash-security/Claw-Hunter) |

## Security Resources

- [centminmod/explain-openclaw](https://github.com/centminmod/explain-openclaw) — Deep security analysis (73 files, multi-AI)
- [Cisco - AI Agents Are a Security Nightmare](https://blogs.cisco.com/ai/personal-ai-agents-like-openclaw-are-a-security-nightmare)
- [Bitdefender - Malicious Skill Trap](https://www.bitdefender.com/en-us/blog/labs/helpful-skills-or-hidden-payloads-bitdefender-labs-dives-deep-into-the-openclaw-malicious-skill-trap)
- [Snyk - 280+ Leaky Skills](https://snyk.io/blog/openclaw-skills-credential-leaks-research/)
- [VirusTotal - From Automation to Infection](https://blog.virustotal.com/2026/02/from-automation-to-infection-how.html)
- [Semgrep - Security Engineer's Cheat Sheet](https://semgrep.dev/blog/openclaw-security-cheat-sheet)
- [n8n + OpenClaw Security White Paper](https://gripsy.com.au/research/n8n-openclaw-security-whitepaper.html)

## Diagnostic Commands

```bash
openclaw doctor                    # Common issues
openclaw security audit --deep     # Deep security scan
openclaw security audit --fix      # Auto-fix safe issues
node --version                     # Must be 22+
```

## Common Issues

| Issue | Cause | Fix |
|-------|-------|-----|
| Crash during onboarding | < 2 GB RAM | Upgrade to 4 GB+ RAM |
| API key not working | Incorrect format or expired | Re-enter via `openclaw auth set` |
| Port conflict on 18789 | Another service using port | `openclaw gateway --port 18790` |
| Docker sandbox issues | Permissions or config | Run `openclaw doctor` |
| High token consumption | ~35,600 tokens per message overhead | Reduce workspace files |
| Rate limiting (#5159) | Aggressive retry loops | Update (exponential backoff fix) |


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]

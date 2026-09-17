# TOOLS.md — Guardian Tool Usage Guide

## Codex Access

The Codex is at `/home/onemind/onemind-os/codex/` — your native knowledge search is built-in.

### ⛔ Codex Write Rules (CRITICAL)

The codex directory structure is **immutable** — managed exclusively by Zeus from the local IDE.

- **NEVER create new folders** in the codex. The numbered folder hierarchy (00-12, 25-34, 50-58, 75-84) is fixed.
- **NEVER rename or renumber** existing codex folders.
- **New notes go to `00-24 SO/06 Inbox (Queue)/`** — Zeus will file them to the correct location.
- **If a folder doesn't exist for a topic**, put the note in `06 Inbox (Queue)/` or `99-organize/` — do NOT create a new folder.
- The codex syncs via git. Any rogue folders you create get auto-committed and pollute the repo.

### memory_search (primary — use this first)

Semantically searches your memory + the Codex vault. Uses BM25+vector hybrid search.

```
memory_search "credential rotation schedule"   → credential management
memory_search "insurance policy renewal"        → insurance tracking
memory_search "legal deadline court"            → legal cases
memory_search "domain registration expiry"      → domain management
memory_search "trust estate beneficiary"        → estate planning
memory_search "passport SSN vital documents"    → identity documents
```

When to use: ANY question about security, legal, insurance, credentials, or compliance. Always search before answering.

### read_file (direct access)

```
read_file /home/onemind/onemind-os/codex/00-24 SO (Sovereign Operations)/11 Security/...
read_file /home/onemind/onemind-os/codex/25-49 HP (Holistic Performance)/31 Legal/...
read_file /home/onemind/onemind-os/codex/50-74 LE (Legacy Evolution)/57 Legal/...
```

## 1Password (Credential Management)

Use for querying credential status:
- Check vault items and their metadata
- Verify credential existence
- Monitor expiration dates
- NEVER output actual credential values — only metadata and status

## Todoist (Task Management)

Use for creating and tracking security/legal tasks:
- Legal deadline reminders
- Insurance renewal tasks
- Credential rotation schedules
- Compliance filing deadlines

## File Tools (read / write / edit)

You have direct file access to the codex:

- `read` — read any codex file by path
- `write` — create new notes (ONLY in `06 Inbox (Queue)/` or `99-organize/`, or update existing files in your domain)
- `edit` — make precise edits to existing codex documents

To create a new note:
```
write /home/onemind/onemind-os/codex/00-24 SO (Sovereign Operations)/06 Inbox (Queue)/My-New-Note.md "# Title\n\nContent here"
```

To update an existing file in your domain:
```
edit /home/onemind/onemind-os/codex/00-24 SO (Sovereign Operations)/11 Security (Cyber)/existing-note.md
```

## Domain Quick Reference

| Path | What's There |
|------|-------------|
| `00-24 SO/11 Security/11.1 Accounts Recovery/` | Account access & recovery |
| `00-24 SO/11 Security/11.2 Business Security/` | Business cybersecurity |
| `00-24 SO/11 Security/11.3 Domain Management/` | DNS, domains, renewals |
| `00-24 SO/11 Security/11.4 Personal Security/` | Personal threat model |
| `00-24 SO/11 Security/11.5 Credentials/` | Credential inventory |
| `25-49 HP/31 Legal/31.4 Identity & Vitals/` | Passport, SSN, birth certs |
| `25-49 HP/31 Legal/31.11 Asset Protection/` | Privacy, trusts, entities |
| `50-74 LE/57 Legal/57.2 Family Legal Cases/` | Active cases |
| `50-74 LE/57 Legal/57.4 Trusts & Estate/` | Trust management |
| `50-74 LE/57 Legal/57.5 Insurance/` | All insurance policies |

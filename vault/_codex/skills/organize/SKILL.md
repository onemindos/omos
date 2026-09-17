# Codex Skill: Organize

How to move items from the Inbox into the correct 00-99 domain folder. This is step "O" of the CODEX method.

---

## The Organize Decision Tree

For each item in `06 Inbox (Queue)/`, ask:

1. **Is it about technology, AI, infrastructure, or the system itself?**
   → `00-24 SO/` — pick the specific subdomain:
   - Architecture & doctrine → `00 Framework/`
   - Dashboards & UI → `01 Command/`
   - AI agents & memory → `02 Agents/`
   - Procedures & SOPs → `06 Protocols/`
   - Reference material → `04 Resources/`
   - Location info → `05 Environments/`
   - Server/API/infra → `08 Infrastructure/`
   - Automation & workflows → `09 Automation/`
   - Code & development → `10 Development/`
   - Security → `11 Security/`
   - Research & analysis → `12 Research/`

2. **Is it about personal health, skills, finance, or relationships?**
   → `25-49 HP/` — pick the specific subdomain:
   - Identity & purpose → `25 Identity/`
   - Learning & education → `26 Intelligence/`
   - Health & body → `27 Body/`
   - Skills & career → `28 Mastery/`
   - Daily routines & ops → `29 Life Systems/`
   - Money & investments → `30 Finance/`
   - Legal matters → `31 Sovereignty/`
   - Social & relationships → `32 Relationships/`
   - Fun & travel → `33 Joy/`
   - Spirituality → `34 Spirituality/`

3. **Is it about family, home, children, or legacy?**
   → `50-74 LE/` — pick the specific subdomain:
   - Home & property → `50 Home/`
   - Marriage & partnership → `51 Marriage/`
   - Children & parenting → `52 Children/`
   - Parents & elder care → `53 Parents/`
   - Extended family → `54 Family/`
   - Estate planning → `57 Estate/`

4. **Is it about business, ventures, or wealth creation?**
   → `75-99 GE/` — pick the specific subdomain by venture or topic

5. **Still unsure?** Leave it in `06 Inbox (Queue)/` and tag it `status/unprocessed`

## Moving a File

1. Move the file to the target domain folder
2. Update the `tags:` in frontmatter to include the domain tag (e.g., `domain/so`)
3. Check for any `[[wiki-links]]` that may need updating
4. If the file was referenced elsewhere, ensure links still resolve

## Batch Processing

When processing multiple inbox items:

1. Read all items in `06 Inbox (Queue)/`
2. Sort them by likely domain (group related items)
3. Move each one, updating frontmatter as you go
4. After moving all, scan for any broken links

## When an Item Spans Domains

Some items touch multiple domains. Rules:

- **Primary home wins** — put it in the domain it most strongly belongs to
- **Cross-reference with tags** — add tags for secondary domains
- **Add `related:` links** — link to notes in the other domain
- **Never duplicate** — one file, one location, with links

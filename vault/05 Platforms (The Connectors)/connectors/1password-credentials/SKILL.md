# 1Password CLI Skill

Manage 1Password vaults, items, and credentials via the `op` CLI (v2.x).

## Prerequisites

- **Binary**: `/opt/homebrew/bin/op` (macOS via Homebrew)
- **Version**: 2.32.0+
- **Account**: `zeus.vaults@pm.me` at `https://my.1password.com/`
- **User ID**: `2HJ5DYRWZFEMLHYP5UUFL23R34`
- **Auth**: Biometric (Touch ID) — requires `eval $(op signin)` at session start

## Authentication

```bash
# Sign in (biometric prompt)
eval $(op signin)

# Verify session
op whoami
```

Session tokens expire after 30 minutes of inactivity. Re-run `eval $(op signin)` if you get auth errors.

## Vault Architecture

| Vault                                    | Domain | Items | Vault ID                       |
| ---------------------------------------- | ------ | ----- | ------------------------------ |
| 00-24 SO (Sovereign Operations)          | SO     | 80    | _(use `op vault list` to get)_ |
| 25-49 HP (Holistic Performance)          | HP     | 76    |                                |
| 50-74 LE (Legacy Evolution)              | LE     | 9     |                                |
| 75-99 GE (Generational Entrepreneurship) | GE     | 34    |                                |
| Personal                                 | —      | 6     |                                |
| Legacy AI                                | —      | 0     |                                |
| Review                                   | —      | 34    |                                |

**Secrets reference format**: `op://VAULT_NAME/Item Name/field-name`
Never store actual secrets in markdown — only reference them via `op://` URIs.

## Core Commands

### List & Query

```bash
# List all vaults
op vault list --format=json

# List all items (all vaults)
op item list --format=json

# List items in a specific vault
op item list --vault="00-24 SO (Sovereign Operations)" --format=json

# List items by category
op item list --categories=API_CREDENTIAL --format=json
op item list --categories=LOGIN --format=json

# Get item details (includes fields, passwords, URLs)
op item get "Item Title" --format=json
op item get "Item Title" --vault="Vault Name" --format=json

# Get a specific field value
op item get "Item Title" --fields label=password
op item get "Item Title" --fields label=username

# Read a secret by reference URI
op read "op://00-24 SO (Sovereign Operations)/OpenAI/credential"
```

### Create Items

```bash
# Create a login
op item create --category=LOGIN \
  --title="Service Name" \
  --vault="00-24 SO (Sovereign Operations)" \
  --url="https://example.com" \
  --generate-password='letters,digits,symbols,32'

# Create an API credential
op item create --category=API_CREDENTIAL \
  --title="Service API Key" \
  --vault="00-24 SO (Sovereign Operations)"

# Create with specific fields
op item create --category=LOGIN \
  --title="My Service" \
  --vault="25-49 HP (Holistic Performance)" \
  'username=me@zeusdelacruz.com' \
  'password=MySecurePass'
```

### Edit Items

```bash
# Update a field
op item edit "Item Title" username="new@email.com"
op item edit "Item Title" --vault="Vault Name" username="new@email.com"

# Update URL
op item edit "Item Title" --url="https://new-url.com"

# Update tags
op item edit "Item Title" --tags="tag1,tag2"
```

### Move Items Between Vaults

```bash
# Move item to a different vault
op item move "Item Title" --current-vault="Review" --destination-vault="25-49 HP (Holistic Performance)"
```

### Delete Items

```bash
# Archive (soft delete)
op item delete "Item Title" --archive

# Permanent delete
op item delete "Item Title"
```

### Search

```bash
# Search by title (partial match)
op item list --format=json | python3 -c "
import json, sys
items = json.load(sys.stdin)
for i in items:
    if 'search_term' in i['title'].lower():
        print(f\"{i['title']} | {i['vault']['name']} | {i['category']}\")
"

# Find items by email/username
op item list --format=json | python3 -c "
import json, sys
items = json.load(sys.stdin)
for i in items:
    info = i.get('additional_information', '')
    if 'target@email.com' in (info or '').lower():
        print(f\"{i['title']} | {i['vault']['name']} | {info}\")
"
```

## Inventory Generation

Pull a full inventory snapshot:

```bash
# Save raw JSON
op item list --format=json > /tmp/1p_raw.json

# Generate markdown inventory tables
python3 << 'PYEOF'
import json, collections

with open('/tmp/1p_raw.json') as f:
    items = json.load(f)

vault_items = collections.defaultdict(list)
for item in items:
    vault = item.get('vault', {}).get('name', 'Unknown')
    vault_items[vault].append(item)

lines = []
for vault in sorted(vault_items.keys()):
    vitems = sorted(vault_items[vault], key=lambda x: x.get('title', '').lower())
    lines.append(f"### {vault} ({len(vitems)} items)\n")
    lines.append("| # | Title | Type | URL | Updated |")
    lines.append("|---|-------|------|-----|---------|")
    for i, item in enumerate(vitems, 1):
        title = item.get('title', '')
        cat = item.get('category', '')
        urls = item.get('urls', [])
        url = urls[0].get('href', '-') if urls else '-'
        if len(url) > 40:
            url = url[:40]
        updated = item.get('updated_at', '')[:10]
        lines.append(f"| {i} | {title} | {cat} | {url} | {updated} |")
    lines.append("")

with open('/tmp/1p_inventory.md', 'w') as f:
    f.write('\n'.join(lines))
print(f"Written {len(lines)} lines to /tmp/1p_inventory.md")
PYEOF
```

## Email Audit

Find accounts by email domain:

```bash
op item list --format=json | python3 -c "
import json, sys, collections
items = json.load(sys.stdin)
domain_counts = collections.Counter()
for item in items:
    info = item.get('additional_information', '')
    if info and '@' in info:
        domain = info.split('@')[1].lower()
        domain_counts[domain] += 1
for domain, count in domain_counts.most_common():
    print(f'{domain}: {count}')
"
```

## Bulk Operations

### Bulk move items from one vault to another

```bash
op item list --vault="Source Vault" --format=json | python3 -c "
import json, sys, subprocess
items = json.load(sys.stdin)
for item in items:
    title = item['title']
    print(f'Moving: {title}')
    subprocess.run([
        'op', 'item', 'move', title,
        '--current-vault', 'Source Vault',
        '--destination-vault', 'Destination Vault'
    ])
"
```

### Bulk update email for items using a specific old email

```bash
op item list --format=json | python3 -c "
import json, sys, subprocess
items = json.load(sys.stdin)
OLD_EMAIL = 'old@example.com'
NEW_EMAIL = 'new@example.com'
for item in items:
    info = item.get('additional_information', '')
    if info and OLD_EMAIL.lower() in info.lower():
        title = item['title']
        vault = item['vault']['name']
        print(f'Updating: {title} in {vault}')
        subprocess.run([
            'op', 'item', 'edit', title,
            '--vault', vault,
            f'username={NEW_EMAIL}'
        ])
"
```

## Environment Variable Injection

Use `op run` to inject secrets into commands without exposing them:

```bash
# Run a command with secrets from 1Password
op run --env-file=.env.1p -- npm run dev

# .env.1p format:
# OPENAI_API_KEY=op://00-24 SO (Sovereign Operations)/OpenAI/credential
# GITHUB_TOKEN=op://00-24 SO (Sovereign Operations)/GitHub/credential
```

## Service Account Tokens

For headless environments (VPS, CI/CD):

```bash
# Create a service account token (via 1Password web UI)
# Then use on the server:
export OP_SERVICE_ACCOUNT_TOKEN="<token>"
op item list  # no signin needed
```

Existing service account: `1password service token (legacy-onemindos)` in SO vault.

## Codex Integration

- **Inventory doc**: `00-24 SO/Security/Credentials/1PASSWORD-INVENTORY.md`
- **Email migration**: `00-24 SO/Security/Credentials/EMAIL-MIGRATION.md`
- **Unified plan**: `00-24 SO/Security/Credentials/EMAIL-1PASSWORD-UNIFIED.md`

When adding new accounts to 1Password, route to the correct vault:

- Tech/AI/infrastructure → **SO** vault (00-24)
- Personal/health/finance/lifestyle → **HP** vault (25-49)
- Government/family/community → **LE** vault (50-74)
- Business/brand/revenue → **GE** vault (75-99)

## Troubleshooting

| Error                        | Fix                                                        |
| ---------------------------- | ---------------------------------------------------------- |
| `not currently signed in`    | Run `eval $(op signin)`                                    |
| `session expired`            | Re-run `eval $(op signin)`                                 |
| `more than one item matches` | Add `--vault="Vault Name"` to disambiguate                 |
| `item not found`             | Check exact title with `op item list \| grep -i "term"`    |
| `biometric unlock disabled`  | Enable in 1Password Settings → Developer → CLI Integration |


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]

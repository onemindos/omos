---
name: github
description: Manage GitHub repositories, issues, pull requests, and workflows. Use the gh CLI for most operations. Use when you need to create repos, manage issues, review PRs, or work with GitHub Actions.
allowed-tools:
  - Bash(gh:*)
  - Bash(git:*)
  - Bash(curl:*)
  - Read
  - Write
---

# GitHub Skill

Interact with GitHub using the `gh` CLI (preferred) or REST API.

## Environment Variables

```bash
GITHUB_TOKEN        # Personal access token (for API)
GH_TOKEN            # Used by gh CLI (auto-configured if logged in)
```

## gh CLI (Preferred)

### Authentication
```bash
gh auth status                    # Check auth status
gh auth login                     # Interactive login
```

### Repositories
```bash
gh repo list                      # List your repos
gh repo view owner/repo           # View repo details
gh repo clone owner/repo          # Clone repo
gh repo create name --public      # Create new repo
gh repo fork owner/repo           # Fork a repo
```

### Issues
```bash
gh issue list                     # List issues
gh issue list --state open        # Open issues only
gh issue list --label "bug"       # Filter by label
gh issue view 123                 # View issue #123
gh issue create --title "Title" --body "Description"
gh issue close 123                # Close issue
gh issue reopen 123               # Reopen issue
gh issue comment 123 --body "Comment text"
```

### Pull Requests
```bash
gh pr list                        # List PRs
gh pr list --state open           # Open PRs only
gh pr view 123                    # View PR #123
gh pr create --title "Title" --body "Description"
gh pr checkout 123                # Checkout PR locally
gh pr merge 123                   # Merge PR
gh pr close 123                   # Close PR
gh pr review 123 --approve        # Approve PR
gh pr review 123 --request-changes --body "Changes needed"
gh pr comment 123 --body "Comment"
```

### Workflows (Actions)
```bash
gh workflow list                  # List workflows
gh workflow view name.yml         # View workflow
gh run list                       # List recent runs
gh run view 12345                 # View run details
gh run watch 12345                # Watch run in progress
gh workflow run name.yml          # Trigger workflow manually
```

### Releases
```bash
gh release list                   # List releases
gh release view v1.0.0            # View release
gh release create v1.0.0 --title "Release" --notes "Notes"
gh release download v1.0.0        # Download release assets
```

### Gists
```bash
gh gist list                      # List gists
gh gist view ID                   # View gist
gh gist create file.txt           # Create gist from file
gh gist create file.txt --public  # Public gist
```

## REST API (Alternative)

### Base URL
```
https://api.github.com
```

### Authentication Header
```bash
Authorization: Bearer $GITHUB_TOKEN
# or
Authorization: token $GITHUB_TOKEN
```

### Common Endpoints
```bash
# List repos
curl -s -H "Authorization: Bearer $GITHUB_TOKEN" \
  "https://api.github.com/user/repos"

# Get repo
curl -s -H "Authorization: Bearer $GITHUB_TOKEN" \
  "https://api.github.com/repos/owner/repo"

# List issues
curl -s -H "Authorization: Bearer $GITHUB_TOKEN" \
  "https://api.github.com/repos/owner/repo/issues"

# Create issue
curl -s -X POST -H "Authorization: Bearer $GITHUB_TOKEN" \
  -H "Content-Type: application/json" \
  "https://api.github.com/repos/owner/repo/issues" \
  -d '{"title": "Issue title", "body": "Issue description"}'

# Get file contents
curl -s -H "Authorization: Bearer $GITHUB_TOKEN" \
  "https://api.github.com/repos/owner/repo/contents/path/to/file"

# Create/Update file
curl -s -X PUT -H "Authorization: Bearer $GITHUB_TOKEN" \
  -H "Content-Type: application/json" \
  "https://api.github.com/repos/owner/repo/contents/path/to/file" \
  -d '{
    "message": "commit message",
    "content": "BASE64_ENCODED_CONTENT",
    "sha": "FILE_SHA_IF_UPDATING"
  }'
```

## One Mind GitHub Repos

| Repo | Purpose |
|------|---------|
| `OneMind-Codex` | Main vault (this repo) |
| `One-Mind` | Public framework |

## Tips

- Prefer `gh` CLI over raw API calls - it handles auth and pagination
- Use `--json` flag with `gh` for machine-readable output: `gh issue list --json number,title`
- Use `--jq` for filtering: `gh issue list --json number,title --jq '.[].title'`
- `gh api` can call any GitHub API endpoint: `gh api repos/owner/repo`


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]

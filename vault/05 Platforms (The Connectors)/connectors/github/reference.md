# GitHub API Reference

## gh CLI Quick Reference

```bash
# Auth
gh auth status
gh auth login

# Repos
gh repo list [owner]
gh repo view [repo]
gh repo clone [repo]
gh repo create [name] --public|--private

# Issues
gh issue list [--state open|closed|all] [--label name]
gh issue view [number]
gh issue create --title "..." --body "..."
gh issue close [number]
gh issue comment [number] --body "..."

# PRs
gh pr list [--state open|closed|merged|all]
gh pr view [number]
gh pr create --title "..." --body "..."
gh pr checkout [number]
gh pr merge [number] [--merge|--squash|--rebase]
gh pr review [number] --approve|--request-changes

# Actions
gh workflow list
gh run list
gh run view [run-id]
gh workflow run [name]

# JSON output
gh issue list --json number,title,state
gh pr list --json number,title,author --jq '.[].title'
```

## REST API Endpoints

| Resource | Endpoint |
|----------|----------|
| User repos | `GET /user/repos` |
| Repo | `GET /repos/{owner}/{repo}` |
| Issues | `GET /repos/{owner}/{repo}/issues` |
| Create issue | `POST /repos/{owner}/{repo}/issues` |
| PRs | `GET /repos/{owner}/{repo}/pulls` |
| Contents | `GET /repos/{owner}/{repo}/contents/{path}` |
| Commits | `GET /repos/{owner}/{repo}/commits` |

## Rate Limits

- **Authenticated**: 5,000 requests/hour
- **Unauthenticated**: 60 requests/hour


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]

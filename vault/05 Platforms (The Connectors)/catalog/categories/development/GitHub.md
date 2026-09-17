# GitHub

> **Version Control & Code Hosting**

## Overview

GitHub is the primary platform for hosting code, managing issues, running CI/CD, and collaborating on OneMind OS development.

## Status

- **Active**: ✅ Yes
- **Organization**: Zeus-Delacruz
- **Last Updated**: January 25, 2026

## Use Cases

1. **Code Hosting** - All OneMind repositories
2. **Version Control** - Git-based history
3. **CI/CD** - GitHub Actions
4. **Project Management** - Issues & Projects
5. **Documentation** - README, wikis

## Key Repositories

| Repository           | Description                     | Visibility |
| -------------------- | ------------------------------- | ---------- |
| `onemind-os`         | Main OneMind OS project         | Private    |
| `OneMind-Codex`      | Knowledge base (Obsidian vault) | Private    |
| `OneMind-Codex-Live` | Public documentation            | Public     |
| `onemind-ui`         | Flutter mobile app              | Private    |
| `legacy-agent`       | Agent-only repo                 | Private    |

## Git Workflow

### Branch Strategy

```
main
├── develop (integration branch)
│   ├── feature/new-feature
│   ├── feature/another-feature
│   └── fix/bug-fix
└── release/v1.0.0
```

### Commit Convention

```
type(scope): description

Types:
- feat: New feature
- fix: Bug fix
- docs: Documentation
- style: Formatting
- refactor: Code restructure
- test: Tests
- chore: Maintenance

Examples:
feat(agent): add calendar integration tool
fix(ui): resolve dark mode toggle issue
docs: update API documentation
```

## GitHub Actions

### CI Pipeline (.github/workflows/ci.yml)

```yaml
name: CI

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Python
        uses: actions/setup-python@v5
        with:
          python-version: "3.12"

      - name: Install dependencies
        run: pip install -r requirements.txt

      - name: Run tests
        run: pytest -v

      - name: Lint
        run: ruff check .
```

### Deploy Pipeline

```yaml
name: Deploy

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Build and push Docker image
        run: |
          docker build -t legacy-agent .
          docker push ${{ secrets.REGISTRY }}/legacy-agent

      - name: Deploy to server
        run: |
          ssh ${{ secrets.SERVER }} "docker pull && docker-compose up -d"
```

## Issue Management

### Issue Templates

**.github/ISSUE_TEMPLATE/bug_report.md**

```markdown
---
name: Bug report
about: Create a report to help us improve
---

**Describe the bug**
A clear description of what the bug is.

**To Reproduce**
Steps to reproduce:

1. Go to '...'
2. Click on '...'
3. See error

**Expected behavior**
What you expected to happen.

**Screenshots**
If applicable, add screenshots.

**Environment**

- OS: [e.g., macOS 14.0]
- Version: [e.g., v1.0.0]
```

### Labels

| Label              | Purpose                 | Color  |
| ------------------ | ----------------------- | ------ |
| `bug`              | Something isn't working | red    |
| `feature`          | New feature request     | green  |
| `enhancement`      | Improvement to existing | blue   |
| `documentation`    | Docs update needed      | purple |
| `good first issue` | Good for newcomers      | teal   |
| `priority: high`   | Urgent issue            | orange |

## Pull Request Process

### PR Template

```markdown
## Description

Brief description of changes.

## Type of Change

- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Checklist

- [ ] My code follows the project style
- [ ] I have added tests
- [ ] All tests pass
- [ ] I have updated documentation
- [ ] I have added comments where needed

## Related Issues

Closes #123
```

### Review Process

1. Create PR from feature branch
2. Request review from @Zeus-Delacruz
3. Address review comments
4. Squash and merge when approved
5. Delete feature branch

## GitHub CLI (gh)

### Common Commands

```bash
# Clone
gh repo clone Zeus-Delacruz/onemind-os

# Create issue
gh issue create --title "Bug: X not working" --body "Description"

# List issues
gh issue list

# Create PR
gh pr create --title "Add new feature" --body "Description"

# View PR
gh pr view 123

# Merge PR
gh pr merge 123 --squash
```

## Secrets Management

### Repository Secrets

| Secret                  | Purpose        |
| ----------------------- | -------------- |
| `AWS_ACCESS_KEY_ID`     | AWS deployment |
| `AWS_SECRET_ACCESS_KEY` | AWS deployment |
| `DOCKER_USERNAME`       | Docker Hub     |
| `DOCKER_PASSWORD`       | Docker Hub     |
| `SERVER_SSH_KEY`        | Deployment SSH |

### Using Secrets in Actions

```yaml
env:
  AWS_ACCESS_KEY_ID: ${{ secrets.AWS_ACCESS_KEY_ID }}
  AWS_SECRET_ACCESS_KEY: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
```

## GitHub Projects (v2)

### Project Board Setup

- **Backlog**: Ideas and future work
- **To Do**: Ready to work on
- **In Progress**: Currently being worked on
- **In Review**: PR open, awaiting review
- **Done**: Completed and merged

### Automation

- Issue opened → Added to Backlog
- PR opened → Moved to In Review
- PR merged → Moved to Done

## Copilot Integration

### MCP Tools Available

- `search_issues` - Search issues
- `create_issue` - Create new issue
- `create_pull_request` - Open PR
- `get_file_contents` - Read files
- `push_files` - Commit files
- `assign_copilot_to_issue` - Assign coding agent

## Best Practices

- ✅ Write clear commit messages
- ✅ Keep PRs focused (one feature/fix per PR)
- ✅ Use issues to track all work
- ✅ Link PRs to issues
- ✅ Review code before merging
- ✅ Keep main branch stable
- ✅ Use branch protection rules
- ✅ Tag releases with semantic versioning

## Resources

- [GitHub Documentation](https://docs.github.com)
- [GitHub Actions Documentation](https://docs.github.com/actions)
- [GitHub CLI Manual](https://cli.github.com/manual/)

---

**Related**: [VS Code](VS-Code.md) | [Docker](00-24%20SO%20(Sovereign%20Operations)/02%20Infrastructure%20(VPS,%20Servers%20&%20Networking)%20(The%20Body)/Tech-Stack/platforms/development/Docker.md)

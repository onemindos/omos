---
title: "GitLab"
type: reference
status: active
tags: [platform, development, source-control, ci-cd, self-hosted]
---
# GitLab

> Git hosting + CI/CD + DevOps platform. Self-hostable alternative to GitHub.

## Use Case in OneMind OS
- Self-hosted option for private/sovereign repos
- Built-in CI/CD (no GitHub Actions dependency)
- Container registry, package registry
- Issue tracking + boards

## CLI
```bash
glab auth login
glab repo list
glab mr create --fill
glab mr merge
glab ci status
glab issue create --title "..." --description "..."
```

## Self-Hosting
```bash
# Docker (recommended)
docker run --detach \
  --hostname gitlab.local \
  --publish 443:443 --publish 80:80 --publish 22:22 \
  --name gitlab \
  gitlab/gitlab-ee:latest
```

## Key Features
| Feature | Notes |
|---------|-------|
| CI/CD | Built-in runners, YAML pipelines |
| Container Registry | Docker images per-project |
| Pages | Static site hosting |
| Wiki | Per-project documentation |
| DevSecOps | SAST, DAST, dependency scanning |

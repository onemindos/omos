# P5 — Guided GitOps deployment templates

Status: planned; this draft PR contains the implementation plan, not a deployed feature.

## Outcome

Create consistent new services through reviewed GitOps pull requests.

## Implementation

- [ ] Start with one service template covering catalog ownership, deployment and observability.
- [ ] Use official Backstage Scaffolder schemas and narrowly scoped repository credentials.
- [ ] Generate a PR instead of immediately mutating the cluster.
- [ ] Include resource limits, networking, secret references, monitoring and documentation.
- [ ] Apply existing identity and permission rules to allowed actions and repositories.

## Completion criteria

- [ ] A template produces a reviewable PR in an authorized repository.
- [ ] Generated manifests and catalog entries pass repository validation.
- [ ] Credentials and secret values never appear in generated source or logs.
- [ ] Failure and retry behavior avoid duplicate deployments.

## Dependencies and sequencing

Use P3 monitoring conventions and P4 documentation conventions. Deployment remains subject to the existing GitOps review flow.

## Release evidence

Before marking ready: summarize actual changes, relevant verification, image digest when needed, rollout health, limitations, and recovery path. Keep this PR draft until implementation is reviewable.

# P5 — Guided GitOps deployment templates

Status: template implemented; authenticated end-to-end publication requires a user-supplied repository-restricted token.

## Outcome

Create consistent new services through reviewed GitOps pull requests.

## Implementation

- [x] Start with one service template covering catalog ownership, deployment and observability.
- [x] Use official Backstage Scaffolder schemas and narrowly scoped repository credentials.
- [x] Generate a PR instead of immediately mutating the cluster.
- [x] Include resource limits, networking, secret references, monitoring and documentation.
- [x] Apply existing identity and permission rules to allowed actions and repositories.

## Completion criteria

- [ ] A template produces a reviewable PR in an authorized repository.
- [x] Generated manifests and catalog entries pass repository validation.
- [ ] Credentials and secret values never appear in generated source or logs.
- [ ] Failure and retry behavior avoid duplicate deployments.

## Dependencies and sequencing

Use P3 monitoring conventions and P4 documentation conventions. Deployment remains subject to the existing GitOps review flow.

## Release evidence

Before marking ready: summarize actual changes, relevant verification, image digest when needed, rollout health, limitations, and recovery path. Keep this PR draft until implementation is reviewable.

The template creates a draft proposal PR under proposals/services, with a stable branch per service and update=true for retries. It does not activate an Argo Application. GitHub token is supplied through the built-in Secret field; no shared cloud mutation credential is installed. Container and health assumptions must be reviewed before promotion.

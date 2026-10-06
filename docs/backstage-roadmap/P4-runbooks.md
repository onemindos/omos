# P4 — Runbooks and Git-owned documentation

Status: planned; this draft PR contains the implementation plan, not a deployed feature.

## Outcome

Publish useful service access, investigation, recovery and backup procedures through TechDocs, independent of Notion.

## Implementation

- [ ] Review the Git archive of all 174 Fabric records; historical notes are not verified deployment recipes.
- [ ] Start with Backstage, identity, NATS, TAK and monitoring runbooks.
- [ ] Build TechDocs in CI and publish to a self-hosted supported storage backend.
- [ ] Connect catalog entities to their documentation and search.
- [ ] Document prerequisites, evidence, recovery steps and verification without embedding secrets.

## Completion criteria

- [ ] Selected critical services have usable access and recovery runbooks.
- [ ] Documentation renders in Backstage and is searchable.
- [ ] Build and publishing failures are visible.
- [ ] Published documentation does not depend on Notion.

## Dependencies and sequencing

Build alongside P3. Validate procedures against actual deployment; coordinate storage/recovery with P8.

## Release evidence

Before marking ready: summarize actual changes, relevant verification, image digest when needed, rollout health, limitations, and recovery path. Keep this PR draft until implementation is reviewable.

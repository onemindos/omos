# Hermes doctrine and the old Legacy fabric (historical)

Git-only reference content for the OneMind agent fabric. It used to sit next to the `hermes-legacy` manifests in
`platform/hermes/agent/`; that instance and its Argo CD Application (`platform-legacy`) were removed on 2026-11-03,
and these files were moved here so the history survives. Nothing in this directory is deployed.

**Read it as history, not as current instructions.** The running agent is `onemind-hermes`
(`platform/hermes/onemind-hermes/`), a single Legacy soul that is a direct client of the NATS bus. The mother/edge-node
framing, the NATS leaf sidecar, the 10 specialist profiles and the pod-level runbook below describe the retired design.

| File | What it is |
|------|------------|
| `CREED.md`, `PILLARS.md`, `PROFILES.md`, `EVOLUTION.md` | Doctrine of the fabric (the pillars SO, LE, HP, GE are still the catalog's structure) |
| `AGENTS.yaml` | Registry of the agents as of the old design (identities, subjects, locations) |
| `agents/` | One `SOUL.md` per area, `legacy`, `legacy-edge`, a `_template` and `planned/` servants |
| `FABRIC-REFERENCE.md`, `LEGACY-ARCHITECTURE.md`, `RUNBOOK.md` | Fabric facts, the architecture overview and the pod runbook of the old instance |

The single Legacy soul that actually runs is seeded from `platform/hermes/seed/onemind-hermes/SOUL.md`; skills are in
`platform/hermes/skills/`; the plugins `onemind-hermes` installs (`nats-platform`, `nats-context`, `onemind-vision`) are in
`platform/hermes/onemind-hermes/plugins/`. Retirement evidence: `platform/hermes/onemind-hermes/RETIREMENT-AUDIT.md`.

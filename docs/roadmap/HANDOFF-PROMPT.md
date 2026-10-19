# Prompt to give the next AI agent

Copy everything below the line.

---

You are taking over an in-progress platform migration in the repository `onemindos/omos`
(local checkout: `/Users/zeus/omos`). You are working for Zeus Delacruz, a one-person team running a
4-node DigitalOcean Kubernetes cluster (Argo CD GitOps, NATS, TAK Server, Authentik, Tailscale,
Backstage) and adopting OpenChoreo 1.3.1 as an app platform on top of it.

Read in this order before doing anything:
1. `PLATFORM-ROADMAP.md` (repo root): the four layers, the AI trust model, the namespace map, the
   phase tracker, the decision log and the evidence for finished phases.
2. `docs/roadmap/HANDOFF.md`: current state, access, rules learned the hard way.
3. The first unfinished phase in the tracker: Phase 6, brief on draft PR #282.
4. `CLAUDE.md` and `AGENTS.md` (operating rules), `docs/MISSION.md` (why), `docs/STACK.md` (live inventory).

State: phases 0–5 are complete with live evidence in the roadmap. STAC Browser is
managed by OpenChoreo; the portal refactor is deployed. Authenticated observability
reuses existing backends, console/MCP logs and metrics work, and backup recovery
passed. Phase 6 (one portal) is next; proceed only with the user’s authorization.

Working rules:
- Everything reaches the cluster through Git and Argo CD. Open a PR, wait for CI, merge, then verify
  live with read-only access (the `onemind-kubernetes`, `onemind-backstage` and
  `onemind-openchoreo` and `onemind-openchoreo-obs` MCP servers, set up by `scripts/ai-operator/setup.sh`; run it if they are
  missing). Mutating `kubectl` (admin context `do-nyc3-omos-cloud`) is for incidents and
  one-off operations only.
- Update the tracker in `PLATFORM-ROADMAP.md` in the same PR that completes work, with evidence.
  A phase is done only when its gate passes against the live cluster.
- Follow the official OpenChoreo v1.3.x docs and note every deviation in the roadmap.
- Layer 1 (Authentik, Tailscale, cert-manager, Velero, Kyverno, Falco, SPIRE, Argo CD, OpenChoreo itself,
  TAK Server) never moves into OpenChoreo. Never pin a workload to a node name. Use `Recreate` for
  workloads with ReadWriteOnce volumes. Never put a secret in Git in plaintext.
- Do not delete, rotate or overwrite credentials, and do not touch TAK Server or NATS without
  asking Zeus. Ask before anything destructive or hard to reverse.
- Zeus has signed in to the console; Phase 5 logs/metrics browser verification passed.
- GitHub CI later could not start due account billing/spending-limit state. Never
  change billing or bypass required checks. Full make validate still has existing
  unrelated YAML debt; record precisely which local checks pass.

Start by rechecking the live cluster and current tracker. Phase 6 is the next phase,
not phase 3. Report verified changes and remaining work in plain language after each PR.

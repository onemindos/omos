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
3. The first unfinished phase is Phase 6 stabilization; use the latest acceptance
   receipts, readiness checklist and backstage-parallel-run.md before selecting work.
4. `CLAUDE.md` and `AGENTS.md` (operating rules), `docs/MISSION.md` (why), `docs/STACK.md` (live inventory).

State: phases 0–5 are complete. **Phase 6 remains In progress**; Phase 7 and
Phase 8 are Not started. Historical Phase 6 integration acceptance was reopened
for portal trust/coverage stabilization. Source/image/live Grafana, Prometheus
and native topology slices are accepted through #349/#350, not all-app coverage
or authorization to migrate. The existing OneMind portal is preserved; the
temporary console is retired. Read the newest evidence, not old image counts.

The parallel PR run prepares eight bounded workstreams. Agents' offline tests and
draft PRs are not deployed features, complete recovery or AI write permission.
Read docs/roadmap/backstage-parallel-run.md and reconcile any shared tracker
conflicts before reviewing/merging one PR at a time. Keep #339/#332/#284 activation
gates intact; existing drafts are not blanket deployment authorization.

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
- Zeus signed in through Authentik; scoped browser/denial checks passed, but
  broader live expiry/second-role/recovery coverage is not universally certified.
  The temporary console is retired; use onemind-preview.tail717ba0.ts.net.
- GitHub CI later could not start due account billing/spending-limit state. Never
  change billing or bypass required checks. Full make validate still has existing
  unrelated YAML debt; record precisely which local checks pass.

Start with the authoritative current tracker and bounded Phase 6 gaps, not an
old claim that Phase 7 is automatically next. Report verified changes and remaining work in plain language after each PR.

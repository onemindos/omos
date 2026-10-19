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
3. `docs/roadmap/phase-3-pilot-app.md`: your first task. Its draft PR is the branch to work on.
4. `CLAUDE.md` and `AGENTS.md` (operating rules), `docs/MISSION.md` (why), `docs/STACK.md` (live inventory).

State: phases 0 to 2 are done and merged. OpenChoreo's control plane, data plane, default
resources, the `onemind-service` component type, a read-only MCP identity for AI agents and a
built-in console exist. No application is managed by OpenChoreo yet. Your job is phase 3, then
the later phases in order, one pull request per step.

Working rules:
- Everything reaches the cluster through Git and Argo CD. Open a PR, wait for CI, merge, then verify
  live with read-only access (the `onemind-kubernetes`, `onemind-backstage` and
  `onemind-openchoreo` MCP servers, set up by `scripts/ai-operator/setup.sh`; run it if they are
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
- Open question from the last session: Zeus has not yet signed in to the OpenChoreo console as
  `akadmin` (`https://onemind-console.tail717ba0.ts.net`). Ask whether it worked.

Start by confirming the current state (Argo applications Healthy, no pending pods) with the
read-only tools, then begin phase 3 on the draft PR "Roadmap phase 3: pilot one app". Report to Zeus
after each merged PR in plain language: what changed, what you verified, what is next.

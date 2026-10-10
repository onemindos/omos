# OneMind project memory

Read `docs/MISSION.md` before making architecture, application, AI, infrastructure, or
product decisions. It records the user's current mission and durable requirements.
Read `CONTRIBUTING.md` for the repository workflow and `docs/NATS-DOCTRINE.md` before
changing event subjects or bus contracts.

OneMind serves a sovereign, home-centered homestead, family, founder workspace, and
distributed workforce of AI agents. Kubernetes is the application substrate;
TAK and real-time operations are central. Local operation and recovery must not
depend on the cloud. Apps must provide useful human and authenticated AI access.

Distinguish user requirements from proposed designs, source code from verified
deployment, and inventory from live state. The platform assessment in
`docs/research/2026-10-08-workspace-platform-review.md` is a recommendation with
unexecuted acceptance tests, not authorization to deploy or replace services.
Preserve unrelated local changes. Later explicit user instructions take precedence.

Current cluster state, tiers, exposure and GitOps gaps are in `docs/STACK.md`; the
purpose of each top-level directory is in `docs/REPO-LAYOUT.md`. The root `CLAUDE.md`
holds the operating rules (deployment pattern, secrets, NATS grammar). CI is
`.github/workflows/quality.yml`; run `make validate` before opening a pull request.
Never commit plaintext credentials, and do not add Helm releases or `kubectl apply`
workloads that Argo CD does not own.

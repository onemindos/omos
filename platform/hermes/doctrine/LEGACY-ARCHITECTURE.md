# Legacy Architecture

> **Historical (2026-11-03).** Describes the retired `hermes-legacy` design (mother and edge nodes). The running agent is `onemind-hermes`.

See the primary docs:

- **`platform/hermes/doctrine/AGENTS.yaml`** — canonical registry of every agent on the fabric (identities, subjects, locations)
- **`platform/hermes/doctrine/CREED.md`** — fabric hierarchy law (Legacy is the mother; all servants defer to her)
- **`README.md`** (this directory) — what is kept here and why; the old deployment README was deleted with the instance
- **`platform/hermes/doctrine/agents/legacy/SOUL.md`** — Legacy (Mother) system prompt / persona
- **`platform/hermes/doctrine/agents/legacy-edge/SOUL.md`** — Legacy Edge system prompt / persona
- **`platform/hermes/doctrine/agents/_template/SOUL.md`** — template for new servant agents
- **`docs/legacy-agent-fabric.md`** — design doc: omnipresence, self-healing, meta-agent, nauth identity
- **`docs/hermes-omos-ui.md`** — integration contract: what omos-ui expects from Hermes on the NATS bus (discovery, heartbeats, wire protocol)

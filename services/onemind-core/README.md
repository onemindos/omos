# OneMind operational core

App-owned object relationships, audited mission requests and versioned Dapr
workflows. PostgreSQL retains requests, idempotency keys and results. Kubernetes
and Synadia Agent Protocol discovery provide fresh observations; the agent
harness remains the owner of agent configuration.

Named initial capabilities inspect infrastructure, summarize source evidence
through LiteLLM, delegate to a currently discovered operator-selected agent, and
request a restart of specifically designated processing workers. External agent
work with an ambiguous outcome is not automatically dispatched again.

The API requires bearer credentials. Read-service credentials cannot mutate.
The UI gateway credential is backend-only and accepts an asserted authenticated
user identity. The trusted portal must enforce operator authorization before
using that credential; network policy must restrict gateway access accordingly.
No new machine credential grants unrestricted cluster writes.

Local test: `PYTHONPATH=services/onemind-core pytest services/onemind-core/tests`.
GitHub Actions tests Python 3.12 and publishes a digest artifact for deployment.
This source/build PR alone does not activate the core in the live cluster.

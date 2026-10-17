# OneMind observational assistant

This PR adds `/assistant` and the user-authenticated `POST /api/onemind-assistant/ask`.
Administrator membership is required on the backend. Machine credentials are rejected.
The assistant reads the same configured health observations used by the portal and a
bounded, permission-aware catalog query. It does not receive tokens, tools, arbitrary
URLs or mutation access. It cannot deploy, repair or grant access.

Without a configured model, the page is a source explorer and explicitly says that no
language model is connected. This is a usable evidence view, not a simulated AI answer.

## Optional sovereign model connection

Configure server-side values in the deployment configuration only after choosing an
operator-controlled Ollama endpoint and a model already downloaded on that server:

```yaml
onemindAssistant:
  ollama:
    baseUrl: http://ollama.YOUR_NAMESPACE.svc:11434
    model: YOUR_INSTALLED_MODEL
    # Optional server-only bearer credential for a protected gateway:
    # token: ${ONEMIND_MODEL_TOKEN}
```

These are illustrative placeholders, not an installed model service. No endpoint or
model is provisioned by this PR. The configured endpoint receives the question and
bounded source content; treat that as an explicit data-sharing decision. Prefer an
internal service protected by network policy, or an authenticated private gateway.
Do not expose an unauthenticated model service to the internet.

The implementation follows [Ollama's chat API](https://docs.ollama.com/api/chat), uses
non-streamed JSON responses, a generation limit, a 30-second model deadline and a
40-second overall request deadline. Disconnecting the client cancels in-flight fetches. Source text
is treated as untrusted context. Only supplied citation IDs are accepted. Answers are
displayed as plain text and cannot execute tools or inject HTML. No prompt, provider
response body or conversation is logged or persisted by this plugin.

## Limits and release

- Ten requests per user per minute and four concurrent requests per backend replica.
  Limits are local to each replica, not a distributed account quota.
- Catalog matching considers at most 200 entities and returns eight matches; incomplete
  coverage is disclosed. This is not semantic search or exhaustive dependency analysis.
- Health probes have their existing observation cache and limited scope; an HTTP result
  does not prove application correctness. Selected inventory source/state/observation
  annotations are preserved; arbitrary annotations are excluded from model context.
- Generative answers can still be wrong despite valid citation IDs. Check source content
  before acting. No changes are executed and no approval bypass is introduced.
- Build and release through the normal portal image PR/Argo process. Add model egress to
  network policy deliberately if required; this PR does not broaden cluster privileges.

TypeScript compilation was performed. No model credential/compute or live model call
was available for verification; this PR must not be presented as a deployed AI service.

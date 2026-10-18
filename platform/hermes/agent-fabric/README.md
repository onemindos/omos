# Agent fabric

Tools that watch and talk to the agent protocol on NATS (`agents.*`). They belong with
Hermes in Git and run in their own namespace, `agent-fabric` (label
`onemind.io/layer: agents`), applied by the `platform-agent-fabric` Application.

| App | Tailnet URL | What it does |
|---|---|---|
| Agent Console | `https://agent-console.tail717ba0.ts.net` | Lists agents on the bus; chat with one (Hermes `legacy`) |
| agent-tap | `https://agent-tap.tail717ba0.ts.net` | Visualizes agent traffic; needs an image rebuild to stream live |

Security:
- Network: default deny; ingress only from Tailscale proxies; egress only DNS and NATS 4222.
- NATS: own SealedSecret `agent-fabric-nats`. It holds the shared `onemind` login for now;
  a dedicated `agent_fabric` user limited to `agents.>` and `_INBOX.>` comes with the NATS
  adoption (the live NATS config is edited outside Helm, so it cannot be changed safely
  with `helm upgrade` first).
- TLS: the NATS server certificate is verified against the NATS CA (`nats-ca` ConfigMap).
- Images pinned by digest.

The Dockerfiles used by `build-nats-tools.yml` stay in `platform/nats-tools/agent-web-ui/`
and `platform/nats-tools/agenttap/`.

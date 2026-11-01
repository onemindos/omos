# LiveKit pilot

Prepared November 1, 2026 at Zeus's request; **not deployed or a proven Hermes integration**.
The manual-sync `platform-livekit` Application stages one stateless LiveKit SFU.
No Helm-CLI release, live secret, firewall, node label or Hermes workload is changed.

## Scope and placement

One `Recreate` replica on the ecosystem pool, pinned upstream v1.13.11 (amd64/arm64
index verified against Docker Hub). Resources are an initial two-client pilot
budget, not a benchmark or a capacity guarantee. Rooms/participants exist in memory;
there is no Redis/HA, retained conversation storage, media recording or PVC. An SFU
restart loses sessions; clients must reconnect. Git reconstructs the service, and
the SealedSecrets recovery key plus independently backed-up API credentials recover
authentication. LiveKit does not acquire Hermes memory by being installed.

The private signaling URL will be `wss://livekit.tail717ba0.ts.net` **after** hostname,
TLS, tailnet authorization and rollout acceptance. It is not a live launcher yet.
There is no bundled chat UI or public token generator. Use an SDK client or locally
built upstream Meet frontend with a scoped short-lived participant JWT. Never send
API secrets to a browser or a hosted third-party playground.

## Networking decisions and activation prerequisites

The upstream Kubernetes deployment model uses host networking, one SFU per node.
This pilot also uses one replica/Recreate; it cannot do zero-interruption updates.
`--bind $(NODE_INTERNAL_IP)` limits **HTTP signaling** to the host's InternalIP
(Kubernetes expands the environment argument; no shell is used). The ClusterIP
Service and Tailscale Ingress terminate private HTTPS/WSS. The root bind setting is
separate from RTC listeners in upstream v1.13.11 `pkg/service/server.go`.

RTC still needs **direct node connectivity**: UDP 7882 mux and TCP 7881 fallback.
`rtc.use_external_ip` discovers the public mapping through STUN. An HTTPS ingress
or Tailscale HTTPS proxy does not carry this media, and no TCP-only proxy promises
working WebRTC. TURN is disabled, so restrictive networks may fail; enabling TURN,
choosing its TLS/DNS/certificates and opening more ports is separate follow-up work.

Before any Argo sync:

1. Confirm the selected ecosystem node's current InternalIP is actually private,
   and HTTP 7880 is unreachable on its public interface. The downward field follows
   node replacement without hardcoding the observed node address.
2. Verify STUN mapping equals that node's reachable external IP. Review node and
   DigitalOcean firewalls: permit required RTC 7881/TCP and 7882/UDP only for the
   chosen clients/networks. Do not expose 7880 publicly, use Cloudflare proxying for
   media, or copy broad firewall allowances from a tutorial. No firewall rule is
   created by this PR. If external-IP discovery is inappropriate, commit a reviewed
   node/edge-specific `rtc.node_ip` design with `use_external_ip: false` instead.
3. Reserve the host ports and verify no host daemon or other SFU claims them. Keep
   at most one instance per node if extending replicas later. Check available CPU,
   memory and public media bandwidth; node-pool autoscaling is not configured here.
4. Confirm tailnet permissions for `livekit` and its managed ingress backend. JWT
   room grants still authorize participation: tailnet membership is not room access.
5. Prepare the real strict-bound SealedSecret below. Empty `sealed/` intentionally
   makes first start wait for credentials, rather than installing development keys.

The namespace's **PSA enforce=privileged is required for hostNetwork**, not a grant
of privileged containers: the SFU runs UID/GID 10001, read-only root, all capabilities
dropped, no service-account token and RuntimeDefault seccomp. PSA warn/audit remain
restricted to expose the exception. No cluster-wide policy exemption is added.
Existing Kyverno policy does not deny hostNetwork; recheck the installed admission
policy before sync. Kubernetes NetworkPolicy does not reliably isolate host-network
traffic; this design does not claim default-deny isolation. Host firewall, private
HTTP binding, tailnet policy and JWTs are the actual network/access boundaries.

## Credentials and client connection

Provision one strong API key/secret mapping in the existing private credential
store. Keep a permission 0600 YAML map exported **outside the repository**, and an
independent encrypted recovery copy. Do not paste either value into tickets or CLI
arguments. Fetch only the controller's public certificate and run:

```bash
python scripts/prepare_livekit_secret.py --keys-file "$LIVEKIT_PRIVATE_KEY_MAP" --cert "$SEALED_SECRETS_PUBLIC_CERT"
```

The helper passes private JSON through stdin to kubeseal; it writes only the
strict-bound `livekit-api` SealedSecret. It refuses repo-local/group-readable
private input, existing output and unexpected metadata. Add
`livekit-api.sealedsecret.yaml` to `sealed/kustomization.yaml`, review/commit it, then
activate through Argo. The helper never applies anything and does not create keys.

Use upstream `lk token create`/Server SDK **on an authorized backend** with secrets
loaded from the credential store to mint a short-lived room-join token, named
participant identity and only required publish/subscribe grants. Give the client
only URL+JWT; keep administrative room/API credentials server-side. An authenticated
web token service/finished family voice interface is a separate deliverable.

Acceptance: render and schema-check; observe Argo/pod readiness and non-root UID;
check that private HTTPS works and public 7880 fails; validate a two-client room on
Mac and Android with video/audio, actual ICE path, reconnect and expired/incorrect
JWT rejection. Do not infer media readiness from HTTP200. No such live acceptance
has been performed for this PR.

## Hermes path worth qualifying

The community [kortexa-ai/hermes-realtime](https://github.com/kortexa-ai/hermes-realtime)
(previously `hermes-livekit`) documents a LiveKit Conference gateway into Hermes's
normal agent loop. Release v0.5.1 / commit
`6cbda7eb3ef10a8a43a04d6ae94debc1db05c442` is the reviewed candidate. Its docs require
Hermes >=0.20.0, Python 3.11–3.14 and FFmpeg, LiveKit server credentials, and explicit
allowed participants. It uses the selected Hermes profile/STT/TTS; it is more
relevant than starting a second independent agent. Upstream
[NousResearch/hermes-agent#3894](https://github.com/NousResearch/hermes-agent/pull/3894)
closed unmerged; this is **not an official upstream connector** or proof on Legacy.

Qualification should use an isolated copy of the intended Hermes version/profile,
review a pinned plugin's dependencies, check current STT/TTS/provider and memory
behavior, and prove room/participant isolation, tool authorization, cancellation,
interruptions, reconnection and no duplicate tool execution. Community streaming
fixes and its adjacent-fork test dependencies still require compatibility review.
Do not install a moving plugin branch or copy this plugin into the production pod
by hand. No plugin installation, provider change, model download, conversation,
physical command or Legacy restart is part of this PR.

## Follow-up boundaries

Bring one MediaMTX camera through a supported republisher/Ingress only after voice
works; RTSP is not an automatic direct LiveKit Ingress source. Egress, SIP, Portal/
ROS2, recording and robotics actuation are not enabled. TAK remains the map, NATS
remains the durable event bus, and onboard controllers own immediate safety/control.
A DO-hosted server is not off-grid. A later local deployment also needs reachable
local auth/DNS, Hermes, STT/LLM/TTS and any model assets, storage and sufficient compute.

## Official sources

- [Kubernetes/direct networking](https://docs.livekit.io/transport/self-hosting/kubernetes/)
- [RTC ports](https://docs.livekit.io/transport/self-hosting/ports-firewall/)
- [Authentication](https://docs.livekit.io/frontends/build/authentication/)
- [Server configuration](https://github.com/livekit/livekit/blob/v1.13.11/config-sample.yaml)
- [HTTP binding implementation](https://github.com/livekit/livekit/blob/v1.13.11/pkg/service/server.go)
- [Meet frontend](https://github.com/livekit-examples/meet)
- [Portal robotics](https://docs.livekit.io/robotics/integrations/portal/)

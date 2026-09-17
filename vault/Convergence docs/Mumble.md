# Mumble: Open Source, Low-Latency Voice Communication

**Link:** https://www.mumble.info/

**License:** Open Source  
**Latest Version:** 1.5.901  
**Repository:** GitHub (mumble-voip/mumble)  
**Status:** Active & maintained  
**Architecture:** Client-server VoIP  

---

## What It Is

Low-latency, high-quality voice chat platform designed for real-time communication. Open-source alternative to proprietary VoIP systems like Skype, Discord, or TeamSpeak. Emphasizes privacy, security, and self-hosting.

**Core Use:** Tactical/operational comms, gaming, communities, enterprise voice — anywhere low-latency encrypted voice is needed without cloud dependency.

---

## What It Can Do

### Voice Communication
- **Ultra-low latency** — <50ms typical (competitive with proprietary systems)
- **High-quality codec** — Opus codec (adaptive bitrate, crystal-clear audio)
- **Voice activity detection** — Noise suppression, automatic transmission
- **Positional audio** — 3D sound (tactical use: know where teammates are by audio)
- **Echo cancellation** — Professional echo/feedback removal

### Channel Architecture
- **Hierarchical channels** — Organize by team, squad, operation
- **Temporary channels** — Ad-hoc team formations
- **Channel groups** — Organize large operations
- **ACL-based permissions** — Fine-grained access control per channel
- **Subdued mode** — Listen-only channels (broadcast/command)

### User Features
- **Encrypted comms** — TLS for transport, optional per-channel encryption
- **Overlay UI** — In-game/app overlay (who's talking, status)
- **Voice macros** — Quick-bind tactical callouts
- **Recording** — Server-side or client-side audio logs
- **Text chat** — Fallback messaging (when voice unavailable)
- **User groups** — Team management, permission templates
- **Friend lists** — Quick access to team roster
- **Priority speaker** — Command channel authority (interrupt for orders)

### Administration
- **Web-based admin** — Manage users, channels, permissions
- **Fine-grained ACL** — Role-based access control
- **Certificate-based auth** — X.509 client certificates (military-grade)
- **Bandwidth control** — QoS per user/channel
- **Logging & audit** — Track who spoke when (command accountability)
- **Multiple servers** — Federate large operations across regions

---

## Technical Specifications

### Client Support
- **Windows** — Full client + overlay
- **macOS** — Full client + overlay
- **Linux** — Full client (CLI or Qt GUI)
- **iOS** — Mumla app (3rd-party, full-featured)
- **Android** — Plumble app (3rd-party, tested/stable)
- **Web** — Webmumble (browser-based, limited features)

### Server Architecture
- **murmurd** — Lightweight server daemon
- **Resource-efficient** — Runs on minimal hardware (Raspberry Pi capable)
- **Scalable** — Handle 1,000+ concurrent users per instance
- **Memory footprint** — ~50 MB base + ~1 MB per user
- **CPU:** Minimal (codec operations lightweight)

### Network & Codec
- **Codec:** Opus (VBR 16-128 kbps, adapts to network)
- **Protocol:** Custom binary (UDP for speed, TCP fallback)
- **Latency:** <50ms typical (UDP path)
- **Bandwidth:** 24-64 kbps per user (adapts to loss)
- **Encryption:** TLS 1.3 for control, optional Speex encryption for audio
- **Ports:** 64738 UDP (primary), 64738 TCP (fallback)

### Quality of Service (QoS)
- **Voice Activity Detection (VAD)** — Only transmit when speaking
- **Automatic Gain Control (AGC)** — Normalize volume
- **Noise Suppression** — Remove background noise
- **Echo Cancellation** — Clear feedback removal
- **Jitter Buffer** — Smooth playback over variable latency

---

## Key Capabilities for Tactical Operations

### Command & Control
- **Push-to-talk (PTT)** — Press key to transmit (hands-free mode available)
- **Priority speaker** — Override channels for command authority
- **Hierarchical channels** — Squad → Platoon → Company structure
- **Broadcast channels** — One-way command distribution
- **Recorded comms** — Full audio logging for after-action review (AAR)

### Integration Points
- **NATS events** — Mumble can publish presence/comms events to NATS
- **TAK Server integration** — Position updates + Mumble comms over same fabric
- **ATAK plugins** — Custom ATAK app integrations via Mumble API
- **Webhooks** — Trigger actions on user join/leave/channel activity
- **Custom applications** — Mumble SDK for custom tactical apps

### Security
- **End-to-end encryption** — Optional per-channel encryption
- **Server encryption** — TLS transport (no cleartext comms)
- **Certificate auth** — X.509-based mutual authentication
- **Access control** — ACL-based channel permissions (who can access what)
- **Audit trail** — All comms actions logged server-side
- **Air-gapped capable** — Runs entirely on local network (no cloud required)

### Resilience & Redundancy
- **Offline-capable** — Record/replay audio locally if server down
- **Fallback codecs** — Degrades gracefully to lower quality if bandwidth-constrained
- **Multi-server federation** — Link servers for large-scale operations
- **Auto-reconnect** — Resume session after network interruption
- **TCP fallback** — Maintains connectivity even when UDP blocked

---

## Deployment Scenarios

### Scenario 1: Tactical Squad Operations (Air-Gapped)
```
Field Server: Mumble daemon (murmurd) on laptop / KubeSolo container
  ├── Low latency voice (<50ms)
  ├── Encrypted TLS comms
  ├── Hierarchical channels (Squad, Teams)
  └── Server runs locally (no internet required)

Field Team (S23 Tactical phones + tablets):
  ├── Plumble app (Android)
  ├── Connect to local Mumble server via WiFi
  ├── PTT for tactical callouts
  ├── Overlay shows who's speaking
  └── Optional position sync to TAK Server

Integration:
  Mumble comms ← NATS event stream → TAK Server
  → Squad leader sees team location + hears comms
  → Unified tactical picture (audio + position)
```

### Scenario 2: Multi-Site Command Network
```
Central TOC: Mumble server + command channel
  ├── Site A (Forward Operating Base)
  ├── Site B (Tactical Operations Center)
  ├── Site C (Air Support / Drone)
  └── All linked via encrypted channels

Hierarchy:
  Command Channel (exec officers only)
    ├── Squad A (operational)
    ├── Squad B (operational)
    └── Air Support (separate channel, listening to command)

All audio encrypted, logged server-side.
Real-time position updates via NATS (Mumble publishes to fabric).
```

### Scenario 3: Large Exercise / Joint Operations
```
Multiple servers federated:
  - US Army Mumble server
  - US Air Force Mumble server
  - Coalition partner Mumble server
  
Bridge channels:
  - Joint Operations Channel (all 3 servers)
  - Translated text comms for non-Mumble users

Scalability: Thousands of concurrent users across servers.
```

---

## Installation & Configuration

### Server Installation (Linux)
```bash
# Install murmurd (Debian/Ubuntu)
sudo apt-get install mumble-server

# Start server
sudo systemctl start mumble-server

# Default config: /etc/mumble-server.ini
# Edit to customize channels, SSL certs, bandwidth, etc.

# Access admin UI
# https://localhost:443/admin (requires authentication)
```

### Docker / KubeSolo Deployment
```dockerfile
# Dockerfile
FROM ubuntu:22.04
RUN apt-get update && apt-get install -y mumble-server
COPY mumble-server.ini /etc/mumble-server.ini
EXPOSE 64738/udp 64738/tcp
CMD ["murmurd", "-c", "/etc/mumble-server.ini"]
```

```yaml
# Deploy on KubeSolo
apiVersion: apps/v1
kind: Deployment
metadata:
  name: mumble-server
  namespace: tactical
spec:
  replicas: 1
  template:
    spec:
      containers:
      - name: mumble
        image: mumble-server:latest
        ports:
        - containerPort: 64738
          protocol: UDP
        - containerPort: 64738
          protocol: TCP
        volumeMounts:
        - name: config
          mountPath: /etc/mumble-server.ini
      volumes:
      - name: config
        configMap:
          name: mumble-config
```

### Client Configuration (Android / Plumble)
```
Server: 192.168.1.50:64738
Username: [Tactical call sign]
Certificate: [Military X.509 cert]
Channels:
  - Command (priority speaker enabled)
  - Squad A
  - Squad B
  - Air Support (listen-only)
```

---

## Configuration Examples

### Mumble Server Config (mumble-server.ini)
```ini
# Listen ports
port=64738

# Server name (appears in client)
serverName=Tactical Ops Mumble

# Max users
users=500

# Bandwidth per user
bandwidth=128000

# Welcome message (on login)
welcometext=Welcome to Tactical Ops. Maintain comms discipline.

# SSL certificate (for encryption)
sslCert=/etc/mumble/server.crt
sslKey=/etc/mumble/server.key

# Logging (audit trail)
logfile=/var/log/mumble-server/mumble-server.log
```

### Channel ACL (Tactical Hierarchy)
```
Channel: Command
  Owner: Exec Officer
  Inherit: Yes
  Write: @"Exec Officers", @"Executive"
  Read: @"All Users"

Channel: Squad A
  Owner: Squad Lead
  Inherit: Yes
  Write: @"Squad A Members"
  Read: @"Squad A Members"
```

---

## Integration with OneMind OS

### Mumble + NATS
```
Mumble server publishes events to NATS:
  - agents.tactical.comms.user.joined
  - agents.tactical.comms.user.speaking
  - agents.tactical.comms.channel.activity

Subscribers:
  - TAK Server (updates COP when user joins/leaves)
  - Alerting system (urgent comms alerts)
  - Analytics (comms usage metrics)
  - Archival (store audio in ClickHouse)
```

### Mumble + TAK Server
```
User speaks on Mumble "Squad A" channel
  ↓
Mumble publishes to NATS: agents.tactical.comms.speaking
  ↓
TAK Server listener subscribes
  ↓
TAK Server highlights user on map (who's active)
  ↓
Operator hears position + sees speaker on map (unified picture)
```

### Mumble + S23 Tactical Mobile
```
S23 runs Plumble app
  ↓ Connects to Mumble server (local WiFi or 5G)
  ↓ PTT for tactical callouts
  ↓ Overlay UI shows speaker
  ↓ Mumble publishes to NATS
  ↓ TAK Server displays speaker position
  ↓ Squad leader sees full tactical picture
```

---

## Use Cases

### Military/Tactical
- **Squad-level comms** — Low-latency voice for real-time coordination
- **Command & Control** — Hierarchical channels, priority speaker for orders
- **Multi-site ops** — Federated servers for joint operations
- **Encrypted comms** — Secure, auditable, classified-network capable
- **Recording/AAR** — Full audio logging for after-action review

### Emergency Response
- **Disaster response** — Coordinate rescue teams with low-latency voice
- **Crisis management** — Hierarchical channels for command structure
- **Interoperability** — Works across agencies (federated servers)

### Enterprise/Team Coordination
- **Remote teams** — Better than video for voice-only efficiency
- **Low bandwidth** — 24-64 kbps per user (works on poor networks)
- **Private hosting** — No cloud data leakage

### Gaming/Community
- **Competitive gaming** — Ultra-low latency for esports
- **Community events** — Large-scale voice chat (1000+ users)
- **Open-source development** — Team coordination during sprints

---

## Comparison with Alternatives

| Aspect | Mumble | Discord | Skype | TeamSpeak | Signal |
|--------|--------|---------|-------|-----------|--------|
| **Latency** | <50ms | 100-200ms | 150-300ms | <100ms | 100-200ms |
| **Encryption** | TLS+optional | TLS only | Proprietary | TLS+optional | E2E |
| **Open Source** | ✓ Yes | ✗ No | ✗ No | ✗ No | ✓ Yes |
| **Self-hostable** | ✓ Yes | ✗ No | ✗ No | ✓ Yes | ✓ Yes |
| **Bandwidth** | 24-64 kbps | 50-100 kbps | 64-128 kbps | 40-80 kbps | 50-100 kbps |
| **Positional audio** | ✓ Yes | ✗ No | ✗ No | ✓ Yes | ✗ No |
| **Priority speaker** | ✓ Yes | ✗ No | ✗ No | ✓ Yes | ✗ No |
| **Classified capable** | ✓ Yes* | ✗ No | ✗ No | ✓ Yes* | ✓ Yes* |

*With military-grade certificate auth & encryption

---

## Performance & Scalability

### Single Server Capacity
- **Concurrent users:** 500-1000 typical (more with tuning)
- **Memory:** ~50 MB base + ~1 MB per user (1GB for 1000 users)
- **Bandwidth:** ~50 kbps per user (100 users = 5 Mbps aggregate)
- **CPU:** Minimal (codec ops lightweight, negligible on modern CPU)
- **Latency:** <50ms (UDP) or <100ms (TCP fallback)

### Scaling Strategies
1. **Vertical** — Larger server hardware
2. **Horizontal** — Multiple Mumble servers (federated)
3. **Load balancing** — Round-robin DNS or proxy
4. **Regional servers** — Local server per region, federate for inter-region comms

---

## Security Considerations

### Transport Security
- **TLS 1.3** — Encrypted control channel
- **Certificate pinning** — Prevent MITM attacks
- **Forward secrecy** — Session keys not derivable from cert

### Application Security
- **Audio encryption** — Optional per-channel encryption (Speex)
- **Access control** — ACL-based channel permissions
- **Authentication** — Certificate-based (X.509) or password
- **Audit trail** — Log all comms events server-side

### Operational Security (OPSEC)
- **Voice discipline** — No call signs on open channels
- **Comms schedule** — Predictable transmission patterns
- **Emission control (EMCON)** — Limit radio transmissions (signal intelligence threat)
- **Alternate frequencies** — Multiple Mumble servers for redundancy

---

## Troubleshooting & Common Issues

### High Latency
**Cause:** Network path congestion, codec too high bitrate  
**Fix:** Lower bandwidth cap, switch to lower bitrate codec, use UDP (not TCP)

### Audio Drops / Jitter
**Cause:** Packet loss, network congestion  
**Fix:** Lower bitrate, switch to lower quality codec, check network stability

### Server Unreachable
**Cause:** Firewall blocking UDP 64738, server down  
**Fix:** Open firewall port, verify murmurd is running, check DNS resolution

### Echo / Feedback
**Cause:** AGC/echo cancellation disabled, poor quality microphone  
**Fix:** Enable echo cancellation in client, use noise-suppressing mic

---

## Getting Started

### Download
- **Windows/macOS/Linux:** https://www.mumble.info/downloads
- **Android (Plumble):** Google Play Store
- **iOS (Mumla):** App Store

### Setup (First-Time)
1. Install Mumble client
2. Add server: `Add New...` → Server address + port
3. Create user certificate (one-time)
4. Connect to server
5. Join channels
6. Test microphone levels (Audio menu)

### Server Setup
1. Install murmurd (server daemon)
2. Generate SSL certificate (`murmur-user-wrapper`)
3. Edit `/etc/mumble-server.ini` (channels, bandwidth, SSL paths)
4. Start service: `sudo systemctl start mumble-server`
5. Connect client, log in as SuperUser (initial config)

---

## Advanced Topics

### Custom ATAK Integration
```java
// ATAK plugin that publishes Mumble events to TAK
public class MumblePlugin extends MapComponent {
  void onMumbleUserSpeaking(MumbleEvent event) {
    CotEvent cot = new CotEvent("b-m-p-s-p-loc", event.userId);
    cot.setPoint(event.userPosition);
    cot.setTime(System.currentTimeMillis());
    broadcast(cot); // Send to TAK Server
  }
}
```

### NATS Integration
```python
# Mumble server publishes to NATS
import nats

async def mumble_listener():
    nc = await nats.connect("nats://localhost:4222")
    async def handler(msg):
        mumble_event = json.loads(msg.data)
        # Process: user speaking, joined, left, etc.
        print(f"Mumble event: {mumble_event}")
    
    await nc.subscribe("agents.tactical.comms.>", handler)
```

---

## Related Resources in Convergence

- **Samsung_S23_Tactical_Mobile_Platform.md** — Mobile client for Mumble (Plumble app)
- **TAK Server** (memory ref) — TAK Server integration for COP
- **KubeSolo.md** — Deploy Mumble server on Kubernetes
- **NATS** (memory ref) — Event bus for Mumble event streaming

---

## Learning Resources

### Official
- **Docs:** https://www.mumble.info/documentation
- **Server Guide:** https://wiki.mumble.info/wiki/murmurguide
- **SDK:** https://github.com/mumble-voip/mumble/tree/master/src/mumble

### Community
- **Wiki:** https://wiki.mumble.info
- **Issue Tracker:** https://github.com/mumble-voip/mumble/issues
- **Forums:** (linked from main site)

---

**Captured:** 2026-07-09  
**Status:** Active & maintained, production-ready  
**Relevance:** Tactical voice comms, low-latency encrypted, self-hostable, NATS-integrable  
**Integration Priority:** High — primary comms layer for distributed tactical operations

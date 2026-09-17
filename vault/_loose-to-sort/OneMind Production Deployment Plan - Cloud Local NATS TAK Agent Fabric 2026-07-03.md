# OneMind Production Deployment Plan

**Date:** 2026-07-03  
**Owner:** Zeus / OneMind  
**Scope:** Cloud + local production skeleton for OneMind as a NATS-native physical-world command, TAK, robotics, agent, and data fabric.

---

# Executive Summary

The current OneMind stack proved the concept, but it is still in a prototype/lab deployment shape:

```text
macOS host
  ├── native NATS / NEX / some Go/Rust/Python services
  ├── Docker Desktop stack: TAK, Postgres, GeoServer, Murmur, Grafana, etc.
  ├── local files and repos
  └── some services launched manually / via NEX / via scripts
```

This is not the long-term operational form.

The production direction is:

```text
Cloud Regional Hub + Local Physical Hub + Future Edge/Vehicle/Office Leaf Nodes
```

Core principle:

```text
NATS is the fabric.
Kubernetes is the cloud deployment substrate.
Proxmox/Linux is the local production substrate.
Leaf nodes connect sites.
TAK is a tactical interface.
OneMind is the source of truth.
Hermes/Legacy/Letta are reasoning agents.
Watchers/ML/rules are always-on.
Agents act through policy, never directly.
```

The goal is to move from development/prototype into **small-scale production that is massive-ready**.

Budget available:

```text
$5,000 between DigitalOcean and AWS
```

Recommendation:

```text
Launch a DigitalOcean cloud regional hub now.
Run Kubernetes/k3s there.
Deploy NATS cluster with JetStream using Helm.
Use Synadia Insights/Protect trial immediately.
Connect current Mac/local NATS as a leaf node.
Deploy OneMind control-plane, command-gateway, agent-gateway, TAK bridge, mission/feed/scheduler.
Deploy simulation as first-class assets.
Use TAK as the map/field interface.
Use atoc as the command/agent/control interface.
Use Postgres/ClickHouse/Object Storage for durable state/history/media.
Keep all physical-world actions behind command policy and NATS permissions.
Then migrate physical local ops to a Linux/Proxmox hub when hardware arrives.
```

---

# 1. Target Architecture

## 1.1 Production Topology

```text
                          ┌──────────────────────────────┐
                          │    Cloud Regional Hub         │
                          │ DO/AWS Linux/Kubernetes       │
                          │ NATS super/coordination node  │
                          │ observability + control plane │
                          └──────────────┬───────────────┘
                                         │
                                  NATS leaf / TLS
                                         │
          ┌──────────────────────────────┼──────────────────────────────┐
          │                              │                              │
          ▼                              ▼                              ▼
┌──────────────────┐          ┌──────────────────┐          ┌──────────────────┐
│ Local Home Hub   │          │ Office Hub       │          │ Mobile/Field Hub │
│ Linux/Proxmox    │          │ future           │          │ vehicle/kit      │
│ NATS + TAK + AI  │          │ NATS leaf        │          │ NATS leaf        │
└───────┬──────────┘          └───────┬──────────┘          └───────┬──────────┘
        │                             │                             │
        ▼                             ▼                             ▼
Home Assistant                cameras/sensors                drones/robots
TAK/iTAK                      local team                     Meshtastic/sat
cameras                       offices                        field devices
robots                        buildings                      edge AI
```

## 1.2 Operating Principle

```text
Every site can operate locally if disconnected.
Cloud coordinates, observes, stores, deploys, and relays.
```

The cloud is **not** the only brain. The cloud is a regional hub.

Local hubs own low-latency physical-world authority:

```text
robot control
drone control
Home Assistant actions
TAK local awareness
camera ingestion
emergency reflexes
local safety actions
```

Cloud owns:

```text
regional coordination
remote access
global observability
deployment control
backup/archive
cross-site mission state
agent coordination
long-term analytics
```

---

# 2. Deployment Decision

## 2.1 Should Kubernetes be used now?

With the cloud budget available, the recommendation is:

```text
Yes — use Kubernetes/k3s for the cloud regional hub now.
No — do not require Kubernetes for every local/edge node yet.
```

Meaning:

```text
Cloud hub: Kubernetes now
Local Mac: dev lab, not production
Future Linux/Proxmox hub: Docker/Podman + NATS + optional k3s
Field/vehicle nodes: NATS leaf + lightweight services, not heavy K8s initially
```

Why use Kubernetes now:

```text
Helm
GitOps
Terraform
secrets management
rolling deploys
persistent volumes
observability
health checks
ingress
cert-manager
container registry
backup strategy
```

Kubernetes gives the correct production operating model.

But edge/local physical operations should remain simple enough to run in degraded networks.

---

# 3. Cloud Provider Strategy

## 3.1 Recommended split

Use **DigitalOcean first** for the first production skeleton.

Use **AWS selectively**.

## 3.2 Why DigitalOcean first

DigitalOcean is recommended for initial deployment because it is:

```text
simpler
faster
less AWS complexity
predictable cost
DOKS Kubernetes is easy
Droplets are easy
volumes are easy
Spaces object storage is simple
good enough for the first regional hub
```

## 3.3 AWS role

Use AWS for:

```text
Bedrock / model access
S3 archival storage
Route53 DNS if desired
ECR if standardizing container registry there
future GPU / specialized workloads
future GovCloud-style architecture if needed
```

## 3.4 Provider verdict

```text
DigitalOcean = operational OneMind hub v1
AWS = model/cloud capability + archive + future scale
```

If choosing only one platform long-term, AWS has more enterprise capability.  
If choosing fastest production skeleton, DigitalOcean wins.

Recommendation:

```text
DigitalOcean first.
AWS selectively.
```

---

# 4. OneMind Production v1 Stack

## 4.1 Layer 0 — Infrastructure as Code

Use:

```text
Terraform / OpenTofu
Ansible
GitHub Actions
SOPS or 1Password CLI for secrets
```

Purpose:

```text
reproducible cloud
reproducible local hubs
no mystery servers
disaster recovery
repeatable deployment
```

Recommended repo structure:

```text
infra/
  terraform/
    digitalocean/
    aws/
  ansible/
    roles/
      nats/
      k3s/
      docker/
      tak/
      observability/
      one-mind-services/
  profiles/
    mac-dev/
    cloud-regional/
    linux-single-node/
    proxmox-hub/
    mobile-edge/
```

Each deployment profile must define:

```text
what runs where
ports
storage
backup
NATS config
security model
service start order
health checks
secrets
volumes
DNS
TLS
operators/users
```

---

# 5. Layer 1 — Networking

## 5.1 Use

```text
Tailscale or WireGuard
Cloudflare DNS or Route53
cert-manager + Let’s Encrypt
NATS TLS
NATS leaf node TLS
```

## 5.2 Network model

Public internet:

```text
HTTPS ingress
TAK TLS if needed
NATS leaf node if needed
```

Private overlay:

```text
admin
NATS routes
database access
observability
SSH
internal service traffic
```

## 5.3 Expose only necessary ports

Expose publicly only:

```text
443 HTTPS ingress
8089 TAK TLS if needed
7422 NATS leafnode port, locked down
4222 NATS client only if required, ideally private/VPN
```

Avoid exposing:

```text
databases
internal dashboards without auth
NATS monitoring without auth
raw service ports
agent endpoints
```

---

# 6. Layer 2 — NATS Fabric

This is the center.

## 6.1 Cloud hub NATS

Deploy:

```text
NATS 3-node cluster
JetStream enabled
NATS Helm chart on Kubernetes
NATS Box for CLI/admin
persistent volumes
TLS
accounts
nkeys/JWT users
system account
monitoring
```

Official NATS docs confirm Kubernetes deployment via NATS Helm chart.

## 6.2 Core subject domains

```text
ui.>
evt.>
ent.>
cmd.>
tak.>
mission.>
task.>
agent.>
source.>
integration.>
audit.>
risk.>
prediction.>
media.>
sim.>
mesh.>
sat.>
```

## 6.3 NATS accounts

Use accounts to isolate domains:

```text
SYS
CORE
AGENTS
TAK
HOME
ROBOTICS
PUBLIC_DATA
MEDIA
OBS
SIM
SECURITY
```

## 6.4 Account purpose

```text
SYS
  system account and server events

CORE
  om-control-plane
  om-router
  om-command-gateway
  om-feed
  om-mission
  om-scheduler

AGENTS
  Hermes
  Legacy
  Letta
  specialist agents
  agent workers

TAK
  TAK bridge
  TAK projector
  TAK chat
  TAK mission sync

HOME
  Home Assistant
  home sensors
  cameras
  locks/lights/automations

ROBOTICS
  drones
  robots
  MAVLink bridges
  ROS2 bridges
  Dimos bridge
  OM1 bridge

PUBLIC_DATA
  AIS
  ADS-B
  weather
  APRS
  satellite
  RF/public feeds

MEDIA
  video stream registry
  video analytics
  media recorder

OBS
  observability/monitoring

SIM
  simulated drones
  simulated robots
  simulated people
  simulated sensors

SECURITY
  audit
  policy
  alerts
  security monitors
```

## 6.5 Critical security rule

```text
Agents do not get direct publish permission to dangerous cmd.>
Agents publish cmd.proposal.>
om-command-gateway upgrades proposal -> real command after policy
```

---

# 7. Layer 3 — Workload Runtime

## 7.1 Cloud runtime

Use:

```text
DOKS Kubernetes or k3s on DO droplets
Helm
Argo CD or Flux GitOps
container registry: GHCR or AWS ECR
```

## 7.2 Local Linux hub later

Use:

```text
Proxmox
VMs
Docker/Podman Compose
optional k3s
NEX for NATS-native organs
```

## 7.3 Edge/mobile runtime

Use:

```text
NATS leaf node
Docker Compose or systemd
few targeted services
local autonomy
store-and-forward
```

---

# 8. Recommended Cloud Deployment Layout

## 8.1 DigitalOcean option A — Managed DOKS

Good for fastest production.

```text
DOKS cluster:
  3 worker nodes
  4 vCPU / 8GB RAM each to start
  autoscale later

Volumes:
  NATS JetStream
  Postgres
  ClickHouse
  Grafana/Loki if needed

Managed DB:
  optional Postgres managed database
```

Services:

```text
nats
nats-box
om-control-plane
om-router
om-command-gateway
om-connector-manager
om-feed
om-mission
om-scheduler
om-agent-runtime
Hermes gateway
TAK bridge
media metadata service
observability
simulation services
```

## 8.2 DigitalOcean option B — k3s on droplets

More control, more ops work.

```text
3 Ubuntu droplets
  k3s server nodes
  Longhorn or local-path storage
  NATS Helm
  Traefik or nginx ingress
```

Good if the goal is learning full production mechanics.

## 8.3 Recommended option

Start with:

```text
DOKS managed Kubernetes
```

unless you specifically want to own the whole k3s layer.

---

# 9. AWS Usage

Use AWS for:

```text
S3 backup/archive
Bedrock model calls
Route53 DNS
CloudWatch optional
future EKS if we outgrow DO
future EC2 GPU if needed
```

Do not start with too many AWS services. Avoid slowing down deployment with AWS complexity.

Recommended starting AWS services:

```text
S3 bucket for backup/archive
Bedrock for agent/model calls
Route53 only if DNS is standardized there
ECR only if we choose AWS registry
```

---

# 10. OneMind Cloud Regional Hub v0

Deploy this first:

```text
cloud.onemind.local / cloud.onemindos.com
```

Initial services:

```text
1. NATS cluster + JetStream
2. NATS accounts/security
3. Synadia trial integration
4. om-control-plane
5. om-command-gateway
6. om-router
7. om-feed/inbox
8. om-mission
9. om-scheduler
10. Hermes/Legacy NATS gateway
11. TAK bridge/projector
12. media registry
13. observability
14. GitOps/deployment pipeline
15. simulation bridge/services
```

Then connect Mac/local stack as a leaf:

```text
Mac/local NATS -> cloud NATS via leafnode
```

Subject flow local to cloud:

```text
local publishes:
  ent.local.>
  evt.local.>
  tak.cot.out
  ui.summary
  mission.status
  audit.>

cloud publishes back:
  cmd.approved.>
  config.>
  mission.assignments
  agent.tasks
```

---

# 11. Full Technology Map

## 11.1 NATS core

Use:

```text
NATS Server
JetStream
KV
Object Store
NATS Services
NATS leaf nodes
NATS accounts
NKEY/JWT auth
NATS Helm chart
nats-box
```

Why:

```text
core nervous system
realtime pub/sub
command bus
distributed service discovery
durable streams
edge/cloud connectivity
security boundary
```

---

# 12. Synadia Platform / Protect / Insights

Use the 14-day trial immediately.

Purpose:

```text
observe NATS interactions
understand subject flows
detect abnormal usage
enforce visibility/control philosophy
learn production NATS operations
```

Map to OneMind:

```text
Synadia Insights = see the fabric
Synadia Protect = control/risk posture
Synadia Platform = managed/operational NATS layer
```

Use the trial to answer:

```text
What subjects are hot?
What clients are noisy?
Who is publishing what?
What permissions should exist?
What agents/services are risky?
Where is the topology broken?
What flows should be blocked?
What should be imported/exported between accounts?
```

---

# 13. Database / State

NATS is the nervous system.  
Databases are memory/index/query.

Use:

```text
Postgres
  missions
  tasks
  schedules
  users
  teams
  permissions
  connector configs
  command approvals
  session metadata

JetStream
  immutable event streams
  command logs
  replayable events

NATS KV
  latest entity state
  latest source health
  latest service health
  latest mission status
  config flags

ClickHouse
  high-volume telemetry/event analytics
  observability
  historical time-series search

Object storage
  video clips
  images
  data packages
  mission files
  model artifacts
  logs archives
```

Object storage options:

```text
DigitalOcean Spaces
AWS S3
MinIO for local/on-prem
```

---

# 14. Control Plane

Services:

```text
om-control-plane
om-entity-registry
om-source-registry
om-service-registry
om-integration-registry
om-workflow-registry
```

Responsibilities:

```text
knows what exists
publishes ui.* projections
tracks health/freshness
ties sources to entities
keeps atoc lightweight
summarizes systems
connects agents to bounded context
```

Subjects:

```text
ui.now.summary
ui.assets.index
ui.systems.index
ui.feed.recent
ui.missions.index
ui.inbox.index
ui.agents.index
ui.sources.index
ui.integrations.index
```

---

# 15. Command / Action Layer

Use:

```text
om-command-gateway
OPA or Cedar policy engine
NATS permissions
audit stream
human approval workflow
```

Command path:

```text
agent/user/tool
  -> cmd.proposal.*
  -> om-command-gateway
  -> policy
  -> approval if needed
  -> cmd.approved.*
  -> executor connector
  -> evt.command.result
  -> audit.command.*
```

Command classes:

```text
Class 0: read-only
Class 1: reversible low-risk
Class 2: physical-world action
Class 3: safety/security/high-risk
Class 4: emergency reflex
```

Examples:

```text
read system health              Class 0
turn on light                   Class 1
send TAK chat                   Class 1/2
unlock door                     Class 2/3
launch drone                    Class 3
robot navigation                Class 2/3
intercept drone                 Class 3/4
emergency abort                 Class 4
```

---

# 16. Agent Layer

Use:

```text
Hermes
Legacy
Letta
specialist NATS agents
Synadia/NATS agent services pattern
NATS service discovery
NATS request/reply
JetStream task queues
```

Agent model:

```text
always-on watchers
event-triggered Hermes
specialist workers
memory-backed sessions
audited commands
```

Do not make one LLM subscribe to everything.

Use:

```text
om-agent-router
om-agent-session
om-agent-memory
om-agent-tools
om-agent-audit
```

Subjects:

```text
agent.hb.>
agent.session.started
agent.session.updated
agent.task.assigned
agent.tool.request
agent.tool.result
agent.intent.*
cmd.proposal.*
```

Roles:

```text
Hermes
  primary reasoning agent / Zeus operator agent

Legacy
  personality/interface/intelligence gateway

Letta
  memory/project/coding/ops agent

Specialist agents
  domain workers: TAK, robotics, security, media, home, public feeds

Watchers
  always-on deterministic monitoring
```

---

# 17. ML / Anomaly Detection

Use:

```text
Python services
River for online ML
scikit-learn / PyOD for anomaly detection
ONNX Runtime for deployable models
YOLO/Ultralytics for vision
OpenCV
Torch only where needed
ClickHouse for historical features
NATS for model input/output
```

Services:

```text
om-anomaly
om-risk
om-predict
om-vision
om-patterns
om-correlator
```

Consumes:

```text
tel.>
evt.>
ent.*.state
tak.cot.out
homeassistant state
robot telemetry
video detection events
network events
agent actions
```

Publishes:

```text
evt.anomaly.detected
risk.entity.{id}
risk.site.{id}
prediction.{domain}.{id}
cmd.proposal.*
ui.feed.recent
```

Start simple:

```text
rules + thresholds + statistical anomaly
then online ML
then GPU vision
then predictive models
```

---

# 18. TAK Layer

GoTAK/TAK_OS confirms the feature target:

```text
Overwatch / COP
Tactical VMS
TAK Feeds
Team Control
mission sync
markers
chat
video
biometrics
RF planning
MDM
federation
```

OneMind version:

```text
TAK is field/tactical interface.
OneMind is source of truth + AI fabric.
```

Deploy:

```text
TAK Server
Postgres for TAK
Certificate authority
om-tak-bridge
om-tak-projector
om-tak-chat
om-tak-mission
om-tak-video
```

Subjects:

```text
tak.cot.in
tak.cot.out
tak.chat.in
tak.chat.out
tak.mission.in
tak.mission.out
tak.overlay.in
tak.overlay.out
tak.video.available
tak.user.location
```

TAK projection:

```text
ent.*.state
  -> om-tak-projector
  -> visibility policy
  -> CoT
  -> TAK Server
```

TAK input:

```text
TAK user location/chat/marker
  -> om-tak-bridge
  -> NATS event/entity state
  -> agents/control plane
```

---

# 19. Tactical VMS / Video

Use:

```text
MediaMTX or SRS
WebRTC for low latency
RTSP for cameras
SRT for unreliable networks
HLS for browser playback
object storage for clips
NATS for metadata/control only
```

Services:

```text
om-media-registry
om-video-router
om-video-recorder
om-video-analytics
```

Subjects:

```text
media.source.{id}.status
media.stream.{id}.available
media.stream.{id}.bookmark
media.stream.{id}.recording.started
det.visual.{camera_id}.{object_id}
```

Rule:

```text
NATS carries metadata/events/control.
Video bytes travel through media protocols/storage.
```

---

# 20. Home Assistant

Service:

```text
om-homeassistant
```

Technology:

```text
Home Assistant WebSocket API
REST API
MQTT if needed
NATS connector
```

Subjects:

```text
source.homeassistant.status
ent.home.device.{id}.state
evt.homeassistant.state_changed
cmd.homeassistant.{entity_id}.set
cmd.homeassistant.scene.activate
evt.command.result.homeassistant
```

Purpose:

```text
prove bidirectional physical-world integration
```

---

# 21. Composio / SaaS Tools

Service:

```text
om-composio
```

Purpose:

```text
tool capability gateway
```

Subjects:

```text
capability.composio.gmail.send
capability.composio.calendar.create
cmd.composio.gmail.send
cmd.composio.calendar.create
evt.composio.action.result
```

Security:

```text
agents request tool actions
command gateway authorizes
connector executes
audit always
```

---

# 22. Robotics / Drones

Use two tracks:

## 22.1 Simulation

```text
ArduPilot SITL
PX4 SITL
Gazebo/Ignition
Webots if needed
mock entity generators
```

## 22.2 Hardware later

```text
MAVLink
ROS2
DDS
vendor SDKs
Dimos bridge
OM1 robotics bridge
Unitree/Go2 bridge
```

Services:

```text
om-drone-bridge-mavlink
om-robot-bridge-ros2
om-sim-bridge
om-robot-command
om-safety-monitor
```

Subjects:

```text
tel.drone.{id}.mavlink
tel.drone.{id}.pose
tel.drone.{id}.battery
ent.drone.{id}.state
cmd.drone.{id}.takeoff
cmd.drone.{id}.goto
cmd.drone.{id}.land
cmd.drone.{id}.rtl

tel.robot.{id}.pose
tel.robot.{id}.battery
ent.robot.{id}.state
cmd.robot.{id}.navigate
cmd.robot.{id}.stop
cmd.robot.{id}.dock
```

Simulation and hardware use the same subject contract.

```text
simulated drone -> same NATS subjects -> TAK/atoc/agents
real drone      -> same NATS subjects -> TAK/atoc/agents
```

Swap source, keep interface.

---

# 23. Meshtastic / Satellite / Degraded Comms

Use:

```text
om-meshtastic
om-satcom-bridge
store-and-forward queues
NATS leaf when IP exists
compressed low-bandwidth event summaries when not
```

Subjects:

```text
mesh.message.in
mesh.message.out
mesh.node.{id}.state
mesh.position.{id}
sat.link.status
```

For Meshtastic, do not send rich streams. Send:

```text
position
short alerts
mission status
emergency
compressed text
```

---

# 24. Public Feeds / TAK Feeds

GoTAK TAK_OS has TAK Feeds:

```text
ADS-B
AIS
weather
satellite imagery
RF
APRS
```

OneMind version:

```text
om-feed-adsb
om-feed-ais
om-feed-weather
om-feed-aprs
om-feed-satellite
om-feed-rf
om-feed-public-safety
```

Subjects:

```text
source.adsb.status
ent.aircraft.{id}.state
source.ais.status
ent.vessel.{id}.state
source.weather.status
evt.weather.alert
source.aprs.status
ent.tracker.{id}.state
```

TAK projection:

```text
selected feed entities -> om-tak-projector -> CoT -> TAK
```

Atoc projection:

```text
feed state -> ui.feed.* / ui.assets.*
```

---

# 25. Observability

Use:

```text
Prometheus
Grafana
Loki
Tempo or OpenTelemetry
ClickHouse for high-volume events
NATS system events
Synadia Insights
Uptime Kuma optional
```

Observability targets:

```text
NATS health
JetStream lag
leaf node connectivity
connector freshness
agent actions
command approvals
TAK bridge health
video stream status
robot/drone status
ML anomaly rate
security violations
```

Subjects:

```text
sys.health.*
service.hb.*
service.error.*
audit.security.*
```

---

# 26. Security System

Use:

```text
NATS accounts
NKEY/JWT users
NATS permissions
Synadia Protect
OPA/Cedar policy
Vault or 1Password
SOPS
mTLS
audit streams
approval workflow
```

Security model:

```text
identity at connection
least-privilege subjects
agents cannot directly perform dangerous actions
all commands audited
all external bridges isolated
all leaf nodes permissioned
all secrets managed
```

Compromise example:

```text
Legacy compromised
  -> cannot publish cmd.robot.* directly
  -> can only publish cmd.proposal.*
  -> command gateway blocks or requires approval
  -> audit.security.violation emitted
  -> Synadia Protect/Insights shows abnormal behavior
```

---

# 27. Concrete Initial Deployment

## Phase A — Cloud Regional Hub

Deploy first:

```text
DigitalOcean DOKS or k3s
3 nodes
NATS Helm
JetStream persistent volumes
NATS accounts
Synadia trial
Grafana/Prometheus/Loki
Postgres
ClickHouse
MinIO/Spaces/S3
om-control-plane
om-command-gateway
om-router
om-feed
om-mission
om-scheduler
Hermes gateway
agent session registry
audit stream
TAK bridge skeleton
simulation services
```

Estimated resources:

```text
3 medium nodes
1 managed Postgres or self-hosted Postgres
object storage
block volumes
load balancer
domain/certs
```

This fits within the starting $5k budget.

---

## Phase B — Connect Local Mac as Leaf

Current Mac stays as physical lab.

```text
Mac NATS -> cloud NATS leafnode
TAK local -> bridge into NATS
Home Assistant local -> connector -> NATS
local demos -> cloud visible
```

Flow:

```text
local ent./evt./tak. -> cloud summaries
cloud mission/cmd approved -> local execution
```

---

## Phase C — TAK Command Center

Deploy cloud TAK or keep local TAK but bridge both.

Recommended:

```text
cloud TAK test server
local TAK server remains
om-tak-bridge supports both
```

Reason:

```text
cloud TAK tests remote team/federation
local TAK tests physical/low-latency
```

Do not make TAK the source of truth. Make TAK a projection/field interface.

---

## Phase D — Simulation Layer

Deploy simulation immediately.

```text
om-sim-drone
om-sim-robot
om-sim-personnel
om-sim-sensor
```

These publish the same subjects real hardware will publish.

Example:

```text
sim.drone.reaper1 -> ent.drone.reaper1.state
sim.robot.om1     -> ent.robot.om1.state
sim.person.zeus   -> ent.person.zeus.location
```

Then TAK/atoc/agents see live simulated assets.

Later:

```text
replace sim source with real Dimos/OM1/MAVLink/ROS2 bridge
```

Same subject contract.

---

## Phase E — Agent Action Layer

Deploy:

```text
Hermes as NATS service
agent task queue
agent session registry
command proposal flow
approval UI/feed
audit logs
```

Agent cannot directly control:

```text
TAK
robots
drones
Home Assistant locks
notifications
```

Agent proposes; command gateway executes.

---

# 28. First Production Services

## First 10 production services

```text
1. nats-core
2. om-control-plane
3. om-command-gateway
4. om-router
5. om-feed-inbox
6. om-mission
7. om-scheduler
8. om-tak-bridge
9. om-sim-bridge
10. om-agent-gateway-hermes
```

## First 5 connectors

```text
1. TAK
2. Home Assistant
3. simulation
4. Composio
5. public feeds: AIS/ADS-B/weather
```

## First 5 UI pages

```text
1. Overwatch
2. Tactical VMS / Video
3. TAK Feeds / Sources
4. Team Control / People / Devices
5. Agents / Sessions / Inbox
```

---

# 29. Deployment Diagram

```text
                           ┌──────────────────────────────┐
                           │  DigitalOcean/AWS Cloud Hub   │
                           │  Kubernetes/k3s               │
                           │                              │
                           │  NATS cluster + JetStream     │
                           │  Postgres + ClickHouse        │
                           │  Object storage               │
                           │  Synadia Insights/Protect     │
                           │  om-control-plane             │
                           │  om-command-gateway           │
                           │  Hermes gateway               │
                           │  TAK bridge                   │
                           │  simulation services          │
                           └───────────────┬──────────────┘
                                           │
                                  TLS NATS Leaf Node
                                           │
                           ┌───────────────▼──────────────┐
                           │  Local Mac / future Linux Hub │
                           │  NATS local                   │
                           │  TAK Server                   │
                           │  Home Assistant               │
                           │  cameras / local robots       │
                           │  iTAK / ATAK                  │
                           └───────────────┬──────────────┘
                                           │
                         ┌─────────────────┼──────────────────┐
                         ▼                 ▼                  ▼
                  Home devices        TAK operators      robots/drones
                  cameras/sensors     chat/CoT/video     sim/hardware
```

---

# 30. Strong Recommendation

Use this model:

```text
Cloud Kubernetes now.
NATS-first everywhere.
Local leaf nodes for physical operations.
Docker/Podman for normal services.
NEX for NATS-native organs.
Synadia Insights/Protect during trial.
TAK as the field interface.
atoc as command/control interface.
Postgres/ClickHouse/Object storage for durable state.
Simulation now, hardware later behind same subjects.
Agents act only through command gateway.
```

This gets OneMind out of pure development without overbuilding in the wrong direction.

---

# 31. Step-by-Step Execution Plan

## Week 1 — Cloud Hub Skeleton

### 1.1 Create cloud accounts/projects

Tasks:

```text
Create DigitalOcean project: onemind-production
Create AWS project/billing tags: onemind-production
Create DNS zone or choose Cloudflare/Route53
Create container registry: GHCR, DOCR, or ECR
Create S3/Spaces bucket for backups/artifacts
Create secret management plan: 1Password/SOPS/Vault
```

Outputs:

```text
Cloud project exists
DNS exists
Registry exists
Object storage exists
Secrets approach selected
```

### 1.2 Create infra repo structure

Tasks:

```text
infra/terraform/digitalocean
infra/terraform/aws
infra/ansible
infra/profiles/cloud-regional
infra/profiles/mac-dev
infra/profiles/linux-single-node
infra/profiles/proxmox-hub
infra/profiles/mobile-edge
```

Outputs:

```text
versioned infrastructure skeleton
```

### 1.3 Provision Kubernetes

Preferred:

```text
DigitalOcean DOKS
3 worker nodes
4 vCPU / 8GB RAM each to start
```

Alternative:

```text
3 Ubuntu droplets with k3s
```

Outputs:

```text
kubectl access
nodes healthy
load balancer available
storage class available
```

### 1.4 Install core cluster services

Tasks:

```text
Install ingress controller
Install cert-manager
Install external-dns if desired
Install sealed-secrets or SOPS integration
Install metrics-server
Install Argo CD or Flux
```

Outputs:

```text
HTTPS ingress works
certs issue successfully
GitOps ready
```

---

## Week 1 continued — NATS Cloud Fabric

### 1.5 Deploy NATS via Helm

Tasks:

```text
Add NATS Helm repo
Create nats namespace
Deploy 3-node NATS cluster
Enable JetStream
Attach persistent volumes
Enable monitoring
Deploy nats-box
```

Outputs:

```text
NATS cluster up
JetStream enabled
nats-box can publish/subscribe
```

### 1.6 Configure NATS accounts/security

Create accounts:

```text
SYS
CORE
AGENTS
TAK
HOME
ROBOTICS
PUBLIC_DATA
MEDIA
OBS
SIM
SECURITY
```

Create users/creds:

```text
om-control-plane
om-router
om-command-gateway
om-feed
om-mission
om-scheduler
hermes-gateway
legacy-gateway
letta-agent
om-tak-bridge
om-tak-projector
om-homeassistant
om-sim-bridge
om-media-registry
om-observability
```

Outputs:

```text
least-privilege credentials
no broad unauthenticated access
system account active
```

### 1.7 Attach Synadia trial

Tasks:

```text
Connect Synadia Insights/Protect trial
Observe subject traffic
Validate topology
Start mapping permissions
```

Outputs:

```text
Synadia dashboard shows NATS fabric
hot subjects visible
clients visible
risks visible
```

---

## Week 2 — OneMind Core Services

### 2.1 Deploy Postgres

Options:

```text
DO managed Postgres
or in-cluster Postgres with persistent volume
```

Purpose:

```text
missions
tasks
schedules
users
teams
permissions
connector configs
command approvals
agent sessions
```

Outputs:

```text
Postgres reachable privately
backup enabled
schema migration path ready
```

### 2.2 Deploy ClickHouse

Purpose:

```text
high-volume event analytics
observability
historical telemetry
queryable event history
```

Outputs:

```text
ClickHouse running
basic ingest path ready
```

### 2.3 Deploy object storage

Use:

```text
DigitalOcean Spaces
AWS S3
or MinIO for local/on-prem
```

Purpose:

```text
video clips
images
mission packages
data packages
model artifacts
logs
backups
```

Outputs:

```text
bucket exists
credentials secured
lifecycle policy defined
```

### 2.4 Deploy om-control-plane

Responsibilities:

```text
entity registry
source registry
service registry
integration registry
workflow registry
ui.* projections
health/freshness state
```

Outputs:

```text
ui.now.summary
ui.assets.index
ui.systems.index
ui.feed.recent
```

### 2.5 Deploy om-command-gateway

Responsibilities:

```text
accept cmd.proposal.*
classify command risk
apply policy
request approval if needed
publish cmd.approved.*
write audit.command.*
```

Outputs:

```text
agent/user actions go through policy
no dangerous direct agent commands
```

### 2.6 Deploy om-router

Responsibilities:

```text
rules engine
when X happens, do Y
dry-run mode
approval mode
replay mode
```

Outputs:

```text
basic event-to-action rules work
```

### 2.7 Deploy om-feed-inbox

Responsibilities:

```text
human-readable feed
inbox
notifications
read/unread
priority
correlation
```

Outputs:

```text
ui.feed.recent
ui.inbox.index
```

### 2.8 Deploy om-mission

Responsibilities:

```text
missions
assignments
assets
status
timeline
TAK mission projection
```

Outputs:

```text
mission.create
mission.assign
mission.status
ui.missions.index
```

### 2.9 Deploy om-scheduler

Responsibilities:

```text
cron jobs
recurring tasks
delayed actions
scheduled missions
scheduled checks
```

Outputs:

```text
task.schedule.*
mission.schedule.*
evt.schedule.triggered
```

### 2.10 Deploy Hermes gateway

Responsibilities:

```text
Hermes as NATS service
wake on important events
read bounded context
publish cmd.proposal.*
write agent session events
```

Outputs:

```text
agent.hermes available
agent.session.* events
agent proposals audited
```

---

## Week 3 — TAK + Simulation

### 3.1 Deploy cloud TAK test stack

Components:

```text
TAK Server
Postgres
certificate authority
WebTAK if needed
```

Outputs:

```text
cloud TAK server reachable
EUD enrollment possible
certificates managed
```

### 3.2 Keep local TAK stack

Purpose:

```text
local physical/field testing
current iTAK/ATAK workflow
local latency
```

Outputs:

```text
local TAK remains operational
```

### 3.3 Deploy om-tak-bridge

Responsibilities:

```text
CoT in/out
chat in/out
mission in/out
overlay in/out
video metadata
device/user location
```

Subjects:

```text
tak.cot.in
tak.cot.out
tak.chat.in
tak.chat.out
tak.mission.in
tak.mission.out
tak.overlay.in
tak.overlay.out
tak.video.available
tak.user.location
```

Outputs:

```text
TAK events enter NATS
NATS entities/actions can project into TAK
```

### 3.4 Deploy om-tak-projector

Responsibilities:

```text
read ent.*.state
apply visibility policy
convert to CoT
inject into TAK
```

Outputs:

```text
selected OneMind entities appear in TAK
```

### 3.5 Deploy simulation bridge

Services:

```text
om-sim-drone
om-sim-robot
om-sim-personnel
om-sim-sensor
```

Subjects:

```text
ent.drone.reaper1.state
ent.robot.om1.state
ent.person.zeus.state
ent.sensor.driveway.state
```

Outputs:

```text
simulated assets visible in NATS
simulated assets visible in TAK
simulated assets visible in atoc
agents can reason over simulated assets
```

---

## Week 4 — First Real Connectors

### 4.1 Home Assistant connector

Build/deploy:

```text
om-homeassistant
```

Needs:

```text
HA websocket token
REST token
entity mapping
NATS creds
command policy class mapping
```

Subjects:

```text
source.homeassistant.status
ent.home.device.{id}.state
evt.homeassistant.state_changed
cmd.homeassistant.{entity_id}.set
cmd.homeassistant.scene.activate
evt.command.result.homeassistant
```

Outputs:

```text
HA state appears in OneMind
OneMind can request HA actions through command gateway
```

### 4.2 Composio connector

Build/deploy:

```text
om-composio
```

Responsibilities:

```text
discover tools
register capabilities
execute approved actions
publish results
```

Outputs:

```text
Composio tools available as governed capabilities
```

### 4.3 Public feeds

Build/deploy first three:

```text
om-feed-ais
om-feed-adsb
om-feed-weather
```

Outputs:

```text
vessels / aircraft / weather alerts enter OneMind
selected feeds appear in TAK
feed status visible in atoc
```

### 4.4 Media registry

Deploy:

```text
om-media-registry
MediaMTX or SRS
```

Outputs:

```text
video stream metadata in NATS
video bytes handled outside NATS
TAK/atoc can discover streams
```

### 4.5 Basic anomaly detector

Deploy:

```text
om-anomaly v0
```

Start with:

```text
rules
thresholds
freshness checks
statistical deviations
```

Outputs:

```text
evt.anomaly.detected
risk.entity.*
ui.feed.recent anomaly events
```

---

## Week 5 — Security Hardening

### 5.1 Lock down NATS permissions

Tasks:

```text
remove broad wildcard permissions
restrict agents to cmd.proposal.*
restrict connectors to their domains
restrict TAK to tak.* and projected entities
restrict HOME to ent.home.* / evt.homeassistant.* / approved commands
restrict ROBOTICS to robot/drone domains
```

Outputs:

```text
least-privilege bus
```

### 5.2 Implement command classes

Tasks:

```text
Class 0 read-only
Class 1 reversible low-risk
Class 2 physical-world action
Class 3 safety/security/high-risk
Class 4 emergency reflex
```

Outputs:

```text
policy decides approval requirement
```

### 5.3 Approval gates

Tasks:

```text
approval UI event model
TAK notification option
atoc inbox approval
voice approval later
expire/revoke approval
```

Outputs:

```text
high-risk actions require explicit approval
```

### 5.4 Synadia Protect/Insights review

Tasks:

```text
review subject activity
identify noisy clients
identify overbroad permissions
identify abnormal flows
create security baseline
```

Outputs:

```text
NATS fabric security posture documented
```

### 5.5 Credential rotation

Tasks:

```text
NKEY/JWT per service
store in SOPS/1Password/Vault
rotate test creds
remove plaintext creds
```

Outputs:

```text
production credential hygiene
```

---

## Week 6 — Local Production Hub Path

### 6.1 Build/buy Linux/Proxmox hub

Target:

```text
Linux/Proxmox machine
NVMe storage
wired Ethernet
UPS
Tailscale/WireGuard
```

### 6.2 Move local physical authority off Mac

Move:

```text
local NATS
local TAK
Home Assistant bridge
camera/media ingress
robot/drone bridges
local om-command executors
```

Mac becomes:

```text
dev workstation
operator workstation
atoc UI host
coding lab
```

Linux hub becomes:

```text
local physical authority
```

### 6.3 Configure leaf to cloud

```text
Linux hub NATS -> Cloud NATS leaf
```

### 6.4 Validate degraded mode

Tests:

```text
cloud disconnected -> local Home/TAK/robot commands still work
local disconnected -> cloud shows stale/freshness correctly
reconnect -> state syncs
```

---

# 32. Final Production Pattern

The final OneMind pattern is:

```text
Sense
  connectors ingest physical/digital world

Normalize
  Envelope + schemas + entity model

Fuse
  om-unity/entity graph resolves truth

Project
  ui.*, tak.*, feed.*, mission.* views

Reason
  agents consume bounded context/read models

Decide
  workflow/policy/approval

Act
  commands to connectors/robots/TAK/Home Assistant/Composio

Audit
  immutable event log + activity feed

Learn
  agents and operators improve rules, missions, context
```

---

# 33. Immediate Next Action Checklist

If we start execution, do this in order:

```text
1. Create production infra repo/profile structure.
2. Choose DigitalOcean DOKS vs k3s droplets.
3. Create DigitalOcean project.
4. Create DNS/subdomain plan.
5. Create object storage bucket.
6. Create container registry.
7. Provision Kubernetes.
8. Install ingress/cert-manager/GitOps.
9. Deploy NATS Helm with JetStream.
10. Create NATS accounts and users.
11. Connect Synadia Insights/Protect trial.
12. Deploy Postgres.
13. Deploy ClickHouse.
14. Deploy observability.
15. Deploy om-control-plane.
16. Deploy om-command-gateway.
17. Deploy om-router.
18. Deploy om-feed-inbox.
19. Deploy om-mission.
20. Deploy om-scheduler.
21. Deploy Hermes gateway.
22. Deploy TAK bridge/projector.
23. Deploy simulation services.
24. Connect Mac/local NATS as leaf.
25. Validate local-to-cloud event flow.
26. Validate cloud-to-local command approval flow.
27. Deploy Home Assistant connector.
28. Deploy Composio connector.
29. Deploy AIS/ADS-B/weather feeds.
30. Deploy media registry.
31. Deploy basic anomaly detector.
32. Lock down NATS permissions.
33. Implement approval classes.
34. Review Synadia Protect/Insights findings.
35. Plan Linux/Proxmox hub migration.
```

---

# 34. Bottom Line

OneMind should now move into a real production skeleton:

```text
Cloud regional hub first.
Local physical hub next.
NATS leaf nodes between all sites.
Simulation now, real hardware later.
TAK as tactical map and field interface.
atoc as command center.
Agents as reasoners.
Watchers/ML as always-on sentries.
Command gateway as safety boundary.
NATS accounts/security as the transport-level control plane.
```

This is the OneMind version of TAK_OS:

```text
more personal
more agentic
more NATS-native
more local-first
more suitable for one operator now
scalable to multi-site/global operations later
```


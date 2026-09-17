---
id: ui-15-20260111053926
title: "Home Assistant Multi-Location Edge Deployment: The 2025 Reality Check"
type: concept
domain: SO
domain: 00-24
folder: "00-24 SO (Sovereign Operations)"
status: active
created: 2026-01-11
modified: 2026-01-11
owner: zeus
tags: [domain/so, folder/15]
todoist_project: "00-24 SO (Sovereign Operations)"
ai_summary: "Home Assistant Multi-Location Edge Deployment: The 2025 Reality Check"
---

# Home Assistant Multi-Location Edge Deployment: The 2025 Reality Check

**The central HA instance controlling multiple edge locations you're envisioning doesn't exist as an official Home Assistant feature**. Instead, the community has developed several architectural patterns that achieve similar goals through different approaches. Here's what actually works in 2024-2025.

## The fundamental truth about HA and multi-location deployments

Home Assistant was designed for single-location smart home management, and as of 2025, **no official multi-location or edge deployment features exist**. The "Satellite" feature introduced in October 2024 is exclusively for voice assistant hardware (like ESP32 devices with microphones), not for multi-location infrastructure deployment. This is a critical distinction that affects all architectural decisions.

The community consensus has coalesced around **running separate HA instances at each location** linked via the **Remote Home Assistant integration**, rather than attempting to orchestrate edge devices from a single central instance. This approach better aligns with network realities, provides local autonomy, and avoids single points of failure.

## Architecture Pattern 1: Remote Home Assistant integration (recommended)

This is the **most mature and widely adopted solution** for multi-location deployments in 2025, and it best matches your goal of centralized control with distributed execution.

**How it works**: You run a full Home Assistant instance at each physical location (office, home, farm, warehouse). Your primary/central instance connects to remote instances via WebSocket API. Remote entities populate into the central instance as if they were native, and service calls automatically route to the appropriate remote location. This gives you your "single pane of glass" while maintaining autonomous operation at each site.

**Installation**: Available via HACS (Home Assistant Community Store) on both master and remote instances. The remote instance needs minimal configuration (essentially just enabling the integration), while the central instance requires the remote instance's IP address and a long-lived access token.

**Key capabilities**: Bidirectional control works out of the box. You can filter which entities sync (include/exclude by domain or specific entities), add entity prefixes to avoid naming conflicts, and even proxy services for components without entities. State updates happen in real-time via WebSocket with automatic reconnection handling. If a remote site loses connectivity, its entities show as unavailable but local automation continues functioning.

**Network requirements**: Requires IP connectivity between sites, typically via WireGuard VPN (most common in 2025), Tailscale mesh VPN (increasingly popular), or Nabu Casa Cloud ($6.50/month per instance). The WebSocket approach means the central instance initiates connections, so remote sites can be behind NAT without issues if you're using VPN.

**Example configuration for the central instance**:
```yaml
# Installed via UI, but can be configured in YAML:
remote_homeassistant:
  instances:
    - host: farm.vpn.local
      port: 8123
      access_token: !secret farm_token
      entity_prefix: "farm_"
      include:
        domains: [sensor, switch, light, climate, binary_sensor]
      exclude:
        entities: [persistent_notification.*]
      load_components:
        - zwave  # Enables Z-Wave services from remote
```

This architecture directly addresses your desire for simplicity and single pane of glass management while avoiding the "manage 5 separate configurations" problem. You configure automations on the central instance that control devices across all locations, and the integration handles routing commands to the appropriate site.

## Docker orchestration for HA: What works and what doesn't

Your Docker/Portainer experience is valuable, but the reality is more nuanced than expected for multi-location HA deployments.

### Portainer Edge Agent for multi-site management

Portainer Edge Agent offers **centralized monitoring and deployment** but not true multi-site orchestration. The architecture works via polling: Edge Agents at remote locations poll the central Portainer server every 5 seconds (configurable), agents initiate all connections outbound (no inbound firewall rules needed), and management occurs through reverse tunnels that auto-close after 5 minutes of inactivity.

**Best use case for HA**: Deploy a Raspberry Pi or similar device at each remote location running the Edge Agent. This device hosts USB-dependent services (Zigbee2MQTT with local coordinator, Z-Wave JS with USB stick) that communicate back to your central HA via MQTT or WebSocket. Portainer gives you centralized management of these edge containers without requiring VPN.

**Real-world implementation pattern** from the community: Central HA instance (on Docker Swarm or standalone) handles automations. Remote RockPro64 nodes connected via Edge Agent run Zigbee2MQTT and Z-Wave JS containers with USB device passthrough. Services expose MQTT/WebSocket endpoints that central HA connects to. You manage stack deployments across all sites from a single Portainer UI.

**Limitations**: This gives you deployment and monitoring centralization but each location still runs independent services. There's no automatic failover between sites, and it's not truly orchestrating HA across locations. You're essentially using Portainer to manage the infrastructure that your central HA connects to.

### Docker Swarm for HA: Single-site HA only

Docker Swarm delivers **high availability within a single location**, not across WAN-separated sites. The community has developed sophisticated patterns using IPVLAN networking (to maintain device discovery via mDNS), GlusterFS for real-time config synchronization, and MariaDB Galera clusters for database replication.

The **HAHA (Highly Available Home Assistant) project** on GitHub provides Ansible playbooks for automated 3-node Raspberry Pi clusters with automatic failover, shared storage, and database replication. This is excellent for ensuring your main home stays online if a Pi fails, but **cannot span multiple physical locations** like your office, farm, and warehouse.

**Why Swarm doesn't work for multi-site**: Overlay networks require low-latency same-subnet connectivity. GlusterFS has quorum and split-brain issues over WAN. Swarm has no multi-site awareness in scheduling. The community consensus from 2024 discussions is clear: multiple users attempted multi-location swarm and concluded it's impractical.

**When to use Swarm**: If your central location needs 99.9% uptime (3+ Raspberry Pi 4s in one home), or you want automatic container rescheduling on hardware failure. For multi-site, use separate instances per location instead.

### K3s/Kubernetes: Complex but capable

K3s deployments work well for **sophisticated single-site HA** with proper device plugin systems for USB coordinators, mature storage options (Longhorn for distributed storage, NFS for shared configs), and excellent service discovery via CoreDNS. The community has documented successful implementations using generic device plugins to dynamically allocate Zigbee/Z-Wave USB dongles to pods as they move between nodes.

**Multi-location pattern**: Deploy **separate K3s clusters per location** with fleet management via Rancher or GitOps tools (Flux/ArgoCD). This is technically viable but introduces significant complexity. You're essentially running independent clusters that happen to be centrally managed, similar to the Remote HA pattern but with Kubernetes overhead.

**Verdict**: Unless you're already deep in the Kubernetes ecosystem or managing dozens of services beyond HA, the complexity outweighs benefits for most multi-location smart home deployments. The learning curve is steep, you lose HA add-ons (container-only mode), and you must manually deploy all supporting services.

## Integration deployment at edge locations

### Zigbee2MQTT distributed architecture

Each physical location **requires its own Zigbee coordinator** due to radio range limitations (typically 30m indoors). You deploy a Raspberry Pi or similar device with a USB coordinator (Sonoff Zigbee 3.0, Aeotec Z-Stick 7) at each location running Zigbee2MQTT. These instances connect to a central MQTT broker (cloud-hosted or at your main site) or use a local MQTT broker with bridging.

**Bandwidth-efficient setup**: Point all remote Zigbee2MQTT instances to a cloud MQTT broker (HiveMQ, CloudMQTT free tier, or self-hosted on a VPS). Your central HA also connects to this broker. Data usage is minimal—kilobytes per day for sensor updates—making this viable even for metered connections like LTE.

**Configuration at remote site**:
```yaml
# Zigbee2MQTT configuration.yaml
mqtt:
  base_topic: zigbee2mqtt/warehouse
  server: mqtt://central-broker.yourdomain.com:1883
  user: warehouse
  password: !secret mqtt_password
```

Your central HA subscribes to these MQTT topics and controls devices via MQTT commands. This works across any network topology without VPN if you secure the MQTT broker properly.

### ESPHome: The resilience champion

ESPHome devices offer **exceptional local automation capabilities** that provide failover resilience. You can program on-device automations that continue operating even when central HA is unreachable. A documented heating controller maintains local control for 1 hour if HA API disconnects, falling back to preset setpoint using local temperature sensor.

**Critical automation pattern**:
```yaml
# ESPHome device config
binary_sensor:
  - platform: gpio
    pin: GPIO12
    name: "Motion Sensor"
    on_press:
      then:
        - switch.turn_on: relay_light
        - delay: 5min
        - switch.turn_off: relay_light

sensor:
  - platform: dht
    humidity:
      on_value_range:
        - above: 70.0
          then: [switch.turn_on: exhaust_fan]
        - below: 50.0
          then: [switch.turn_off: exhaust_fan]
```

**Network connectivity**: ESPHome devices connect via native API (preferred for same network/VPN), or MQTT for remote locations without VPN. The devices on remote subnets connected via site-to-site VPN work excellently with the native API, or you can configure them to use MQTT to avoid VPN complexity.

**Best practice**: Implement critical automations (lighting, HVAC, security) locally on ESP devices. Central HA updates parameters like setpoints and schedules but isn't required for basic operation.

### Z-Wave remote deployment

Z-Wave networks require physical proximity to coordinators, but you can run **multiple Z-Wave networks via Z-Wave JS WebSocket**. Deploy Raspberry Pi with Aeotec stick at each location running Z-Wave JS server (not Zigbee2MQTT—different software). Add multiple Z-Wave integrations in central HA pointing to different WebSocket endpoints.

**Setup process**: Install Z-Wave JS on remote Pi, configure WebSocket connection (port 3000 default), add Z-Wave integration in central HA (select "Z-Wave" type, not "add device"), provide remote IP:port. Community has documented successful deployments at distances up to half a mile over network.

**Performance note**: Update coordinator firmware for best performance. Initial firmware on new coordinators causes slowness; use Simplicity Studio 5 for Aeotec updates.

## MQTT versus WebSocket: The synchronization showdown

The 2025 consensus has **shifted away from MQTT-based synchronization** toward WebSocket approaches.

### MQTT Statestream: Legacy with issues

MQTT Statestream publishes one-way state changes to MQTT topics but has **known bugs in 2024** (filter issues in HA 2024.7+, race conditions on startup). It provides read-only synchronization; you need complex automations to achieve bidirectional control. The community reports that include/exclude filters don't work correctly—entities publish despite exclusion rules.

**When it's acceptable**: If you only need monitoring (sensors from remote site visible on central HA) and don't need to control devices from central, Statestream works. For anything more sophisticated, use Remote Home Assistant integration.

### MQTT EventStream: Actively discouraged

MQTT EventStream publishes **all events** between instances and causes performance degradation, event loops if misconfigured, and shutdown problems (particularly with Z-Wave/ZHA). Community guidance from 2024 is explicit: "Don't do it." The integration sees very low adoption (~144 active installations versus ~1,023 for Statestream).

### MQTT DiscoveryStream: Better MQTT option

If you must use MQTT, the **MQTT DiscoveryStream custom component** (available via HACS) extends Statestream with automatic entity discovery and built-in bidirectional control. It generates proper MQTT discovery messages so entities appear automatically on the central instance with full control capabilities.

### Remote Home Assistant via WebSocket: The winner

**Why WebSocket beats MQTT for multi-instance**: Bidirectional by default, no MQTT infrastructure needed, better error handling with errors propagated back to primary, automatic state tracking with entity removal on connection loss, and native integration with HA's internal APIs.

The community has decisively moved to this approach in 2024-2025. Forum discussions consistently recommend Remote HA over MQTT solutions for new deployments.

## The "single pane of glass" reality

Achieving truly unified management requires understanding current limitations and workarounds.

### What works today

Your **central HA dashboard shows entities from all locations** via Remote Home Assistant integration. You create automations on the central instance that control devices across sites. The integration handles routing service calls to appropriate remote instances automatically. You can build unified dashboards showing farm temperature sensors, office lighting, warehouse security cameras, and home HVAC in a single view.

### What doesn't work yet

The HA mobile apps (iOS/Android) **don't natively support multi-instance switching**. This is a frequently requested feature in the community but remains unimplemented as of 2025. Your app connects to your primary instance, which has access to all remote entities via Remote HA integration, so you have control. However, directly accessing individual remote instances requires web browser bookmarks.

**Workaround**: Use the native app for your primary/central home. Create Safari/Chrome bookmarks (add to home screen) for direct access to remote instances when needed. The primary instance remains your "single pane of glass" for cross-site automation and monitoring.

### Entity naming strategy

Use **location prefixes for all entities** to prevent conflicts: `farm_light.workshop`, `warehouse_sensor.temperature`, `office_switch.meeting_room`. The Remote HA integration supports automatic prefixing via the `entity_prefix` configuration option.

## Backup and redundancy patterns

### Per-location autonomy (recommended approach)

The **most robust pattern** is ensuring each physical location operates independently when network connectivity fails. Deploy full HA instance at each site with local automation for critical functions. Use ESPHome local automations for essential device control (lighting, HVAC). Implement Zigbee groups for coordinator-independent operation—switches bind directly to light groups, bypassing the coordinator entirely.

**Zigbee group pattern**: Create a group in Zigbee2MQTT or ZHA, add lights to group, bind physical switch to group (not coordinator). Commands from switch go directly to group devices even if HA/coordinator is offline. This provides millisecond response times with no network dependency.

### Single-site high availability (HAHA pattern)

For locations requiring 99.9% uptime, deploy the **HAHA (Highly Available Home Assistant) solution**: 3+ Raspberry Pi 4s in Docker Swarm with GlusterFS real-time config synchronization, MariaDB Galera cluster for database replication, and automated Ansible deployment. Node failure is detected by Docker Swarm, containers reschedule to healthy nodes, and typical recovery is 30-60 seconds.

**Extreme HA mode**: Run 2 HA replicas simultaneously. Both instances write to the same GlusterFS volume (handles concurrent writes gracefully) with zero downtime during failover. Cost is higher resource usage.

### Backup strategy

**Automated backups**: HA built-in backups (Settings → System → Backups) scheduled daily minimum, stored offsite on Samba share or cloud storage. Pre-upgrade backups are essential. Coordinator backups via ZHA/Zigbee2MQTT preserve network topology.

**State preservation**: Configure ESPHome devices with `restore_mode: RESTORE_DEFAULT_OFF` or similar options. Zigbee devices store state in firmware. Enable MQTT retained messages for last-known state persistence across broker restarts.

### Monitoring and alerting

Implement **connection status monitoring** with automations that alert when remote instances go offline. Monitor coordinator health via binary sensors, log device offline events, and use ping sensors for VPN tunnel monitoring. Test failover scenarios monthly—simulate HA restart, disconnect network, trigger full failover—measuring recovery time and validating critical functions.

## Network architecture considerations

### WireGuard VPN (most popular 2025)

**Site-to-site configuration**: Deploy WireGuard add-on in central HA, configure server settings with external IP or DuckDNS domain, forward port 51820 UDP on router, generate keys for each remote site, install WireGuard client on remote devices with persistent keepalive (25 seconds recommended).

**Behind CGNAT workaround**: Use Oracle Cloud free tier or VPS as VPN relay. Central HA and all remote sites connect as clients to relay. All communication routes through relay server.

**Why WireGuard wins**: Modern fast protocol, low overhead, built-in encryption, works behind CGNAT in client mode, easy HA add-on installation, and dominant community adoption in 2024-2025.

### Cloudflare Tunnel (zero-trust alternative)

**Zero port forwarding**: Outbound-only connections from HA to Cloudflare edge. Install Cloudflare add-on, create tunnel in Cloudflare dashboard, configure HA trusted proxies. This works behind any firewall or CGNAT without router configuration.

**Limitations**: Terms prohibit heavy video streaming. Best for accessing HA UI remotely, webhook receivers, and mobile app access. Not ideal for high-bandwidth camera feeds or frequent automation triggers. Potential latency for real-time control.

### Tailscale mesh VPN (increasingly popular)

Managed WireGuard-based mesh with zero-config peer-to-peer connections. Free tier sufficient for most home setups. All devices get private IPs, and connections are direct peer-to-peer when possible (falls back to relay servers). Excellent for non-technical users who want VPN benefits without VPN complexity.

## Better alternatives to multiple HA instances?

The question of "better orchestration tools" has a clear 2025 answer: **the community has tried many approaches, and separate HA instances with Remote HA integration has emerged as the most reliable pattern**.

### Why not single HA with edge agents?

A single central HA controlling lightweight edge agents (ESPHome devices, MQTT sensors) seems simpler but creates a single point of failure. Remote site loses all automation if central HA restarts or network connection drops. Zigbee/Z-Wave radio range limits force physical coordinators at remote locations anyway. Processing latency increases with centralization.

**When single HA works**: Closely located buildings with reliable networking (large property with detached garage 100m away), stable VPN connection between sites, remote site doesn't need autonomous operation, and low device count at remote locations.

### Why not manage everything through Portainer/K8s?

Orchestration platforms like Kubernetes or Docker Swarm are designed for **stateless services and distributed computing**, not multi-location IoT device management. Home Assistant's architecture (stateful, local hardware dependencies, device discovery protocols) doesn't align well with these patterns at scale.

**The overhead problem**: K3s requires significantly more resources than standalone Docker. You lose HA add-ons (container-only mode) and must manually deploy all supporting services. Configuration management via YAML manifests is complex for home automation scenarios. You're essentially recreating the Remote HA integration pattern with Kubernetes complexity on top.

### The winning pattern: Distributed instances with WebSocket federation

**Architectural principles**: Full HA instance at each physical location (autonomous operation), master instance at primary location (unified control), Remote Home Assistant integration links them (WebSocket API), and critical automations run locally (resilience).

This pattern aligns with network realities (WAN latency, unreliable connectivity), follows distributed systems principles (no single point of failure), leverages HA's strengths (local control, extensive integrations), and proves simplest to maintain (community consensus from years of experimentation).

## Practical implementation roadmap

For your specific scenario with Docker/Portainer experience wanting centralized control with distributed execution:

### Phase 1: Central instance setup
Deploy HA at your primary location (Intel NUC or HA Blue recommended over Raspberry Pi for reliability). Install Remote Home Assistant integration via HACS. Configure WireGuard VPN or Tailscale for site-to-site connectivity. Deploy Portainer if you want centralized Docker stack management.

### Phase 2: First remote location
Deploy full HA instance on Raspberry Pi 4 with SSD boot (not SD card—reliability issues). Install Remote Home Assistant integration (enable on remote, minimal config). Add remote instance to central HA via integration UI (provide IP, port, long-lived access token). Configure entity prefix (e.g., `farm_`). Test bidirectional control and verify autonomous operation by disconnecting network.

### Phase 3: Integration deployment
Deploy Zigbee2MQTT or Z-Wave JS at locations with physical coordinators needed. Configure ESPHome devices with local automations for critical functions. Point all integrations to central MQTT broker or enable direct WebSocket connections. Use Portainer Edge Agent if you want centralized monitoring of edge Docker containers.

### Phase 4: Scaling to additional locations
Repeat Phase 2 for each additional location (office, warehouse, etc.). Use consistent naming conventions with location prefixes. Build unified dashboards on central instance. Create cross-site automations (e.g., "turn off all warehouse lights when leaving office"). Implement monitoring automations for connection status.

### Phase 5: Resilience and backup
Configure automated backups on all instances (daily, stored offsite). Implement ESPHome local automations for critical functions. Set up Zigbee groups for coordinator-independent operation. Test failover scenarios (network loss, HA restart, hardware failure). Document recovery procedures.

## What actually matters most

After analyzing all approaches, **three factors determine success** for multi-location HA deployments:

**Local autonomy**: Each site must function when network connectivity fails. Critical automations (lighting, HVAC, security) run locally on the HA instance at that location or on ESPHome devices.

**Simple synchronization**: Remote Home Assistant integration via WebSocket provides the best balance of capability and simplicity. MQTT solutions work but add complexity and have more edge cases.

**Network resilience**: VPN is essential (WireGuard or Tailscale). Cloudflare Tunnel provides backup access method. Monitor connection health and alert on failures.

Your desire for simplicity and centralized management is best achieved not by orchestrating edge devices from a single central instance, but by deploying full instances at each location and linking them. This architecture has proven most reliable through years of community experimentation and aligns with Home Assistant's design philosophy of local control and privacy.

The Docker/Portainer experience you have is valuable for managing the infrastructure (deploying Edge Agents for remote Docker host management, creating Edge Stacks for consistent service deployment), but the "orchestration" happens at the Home Assistant integration layer (Remote HA integration) rather than the container orchestration layer (Swarm/K3s).

## Key GitHub repositories and resources

**Remote Home Assistant**: github.com/custom-components/remote_homeassistant (most critical resource)

**HAHA High Availability**: github.com/cvb941/HAHA (single-site HA via Docker Swarm)

**ESPHome Voice Satellites**: github.com/tetele/onju-voice-satellite (lightweight edge agents)

**MQTT Share Remote**: github.com/mikelawrence/mqtt-share-remote-hacs (better than built-in MQTT integrations)

**Community discussions**: community.home-assistant.io topics on "multi-site support" and "multiple property management" contain years of practical experience and validated patterns

The 2024-2025 community consensus is remarkably clear: separate HA instances per location, linked via Remote Home Assistant integration, connected via WireGuard VPN, with ESPHome providing local automation resilience. This architecture delivers your centralized control requirement while maintaining distributed execution and local autonomy.
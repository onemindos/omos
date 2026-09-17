# Anduril EagleEye: Tactical AR Platform for Warfighters

**Link:** https://www.anduril.com/eagleeye

**Vendor:** Anduril Industries  
**Product Type:** Warfighter-worn AR/wearable computing platform  
**Tagline:** "Superpowers for Superheroes"  
**Status:** In development / early deployment  
**Form Factor:** Combat eyewear + integrated compute  
**Architecture:** Modular, extensible, mission-configurable  

---

## What It Is

Augmented reality platform purpose-built for combat operations. Integrates AR display (eyewear), compute, sensors, and communications into a single warfighter-worn system. Designed to give individual operators real-time access to intelligence, communications, targeting, and mission planning tools.

**Core Purpose:** Transform individual warfighter into network node with full situational awareness, connectivity, and decision support — from individual gunfighter to integrated team.

---

## What It Can Do: Four Core Capabilities

### 1. Enhanced Perception
- **Real-time sensor fusion** — Integrate tactical data streams into visual display
- **Overlaid intelligence** — See enemy positions, friendly locations, terrain overlays
- **Thermal/night vision** — Integrated night-fighting capability
- **Target identification** — Auto-label threats (AI-powered target recognition)
- **Environmental sensing** — Atmospheric data, obstacles, threats displayed in real-time
- **ISR integration** — Display drone/overhead imagery in AR (know what air assets see)

### 2. Heightened Survivability
- **Ballistic protection** — Combat eyewear with military-grade ballistic rating
- **Integrated protection** — Weight + protection + compute balanced
- **Heads-up display** — Critical info without looking down (maintain situational awareness)
- **Threat alerts** — Real-time warnings (incoming fire, hazards, system status)
- **Communications redundancy** — Multi-path comms (cellular, radio, mesh)
- **Armor integration** — Works with existing body armor, tactical kit

### 3. Lethal Connectivity
- **Fire control integration** — Direct link to weapons systems (targeting, fire support)
- **Comms architecture** — Voice, data, video integrated into single system
- **Mission command** — Real-time orders, targeting data, fire missions
- **Peer-to-peer network** — Direct warfighter-to-warfighter comms (squad level)
- **Cloud connectivity** — Link to command, ISR, fires (when available)
- **Open API** — Integrate with existing tactical systems (TAK, ATAK, etc.)

### 4. Mission Planning
- **Real-time mission briefs** — Display objectives, routes, intelligence in AR
- **Route planning** — Heads-up turn-by-turn navigation (no looking at map)
- **Objective overlays** — See targets, waypoints, restricted areas overlaid on terrain
- **Contingency planning** — Pre-loaded alternate routes, backup objectives
- **Decision support** — AI-assisted analysis (threat assessment, optimal tactics)
- **After-action** — Record and review mission for lessons learned

---

## Technical Architecture

### Hardware Integration

**Display:**
- **Combat eyewear form factor** (Oakley Standard Issue partnership)
- **AR optics** — Full-color, heads-up display
- **Ballistic rating** — Military-grade protection
- **Field of view** — Wide enough for tactical awareness
- **Low-light capable** — Night vision / FLIR integration

**Compute:**
- **Qualcomm integration** — Small form factor, power-efficient processor
- **Real-time processing** — <100ms latency for combat effectiveness
- **AI/ML inference** — On-device threat detection, target recognition
- **Storage** — Local cache for offline capability
- **Thermal management** — Designed for extended field ops

**Sensors:**
- **Camera** — Visible spectrum + thermal
- **IMU** — Gyro, accelerometer, compass (orientation/movement)
- **GPS/GNSS** — Position awareness
- **LiDAR/depth** — 3D environment mapping
- **Microphone** — Voice control + comms
- **Extensible slots** — Add sensors as needed (radar, RF, etc.)

**Power:**
- **Battery** — Extended field life (8+ hours typical)
- **Quick-swap** — Carry multiple batteries
- **Power management** — Adaptive modes (high-performance vs. extended operation)
- **Solar option** — Possible add-on (future)

### Software Stack

**Core OS:**
- **Real-time kernel** — Low-latency responsiveness
- **Linux-based** (likely) — Open, hackable, military-familiar
- **Modular architecture** — Load/unload capabilities as needed

**Applications:**
- **Situational Awareness** — Real-time threat/friendly position map
- **Comms Controller** — Multi-channel voice/data management
- **Targeting** — Integration with weapon systems
- **Navigation** — Heads-up routing with terrain overlay
- **Recording** — Capture video/telemetry for review

**Integration Layer:**
- **TAK/ATAK compatibility** — Import/export position, intel, tasking
- **API framework** — Open for tactical app development
- **Mission upload** — Sync with command, pull updated intel
- **Plugin system** — Extend with customer software

### Network Stack

**Comms Modes:**
1. **Cellular (5G/4G LTE)** — Primary when available
2. **Military radio** — Integrated with squad/unit net
3. **WiFi / ad-hoc mesh** — Local team connectivity
4. **Satellite (optional)** — Beyond-cellular operations
5. **Physical hand-off** — Store-and-forward when disconnected

**Protocols:**
- **Encrypted channels** — E2E for sensitive comms
- **TAK protocol** — Compatible with TAK Server ecosystem
- **Open standards** — NATS integration possible
- **Fallback modes** — Operate locally if network lost

---

## Integration with OneMind OS / Tactical Operations

### EagleEye as TAK/ATAK Terminal
```
EagleEye warfighter
  ↓ AR display (see friendly positions)
  ↓ Integrated comms (talk to team)
  ↓ TAK protocol (sync with TAK Server)
  ↓ Mission planning (see objectives in AR)
  ↓ Targeting (call in fires)
  
TAK Server aggregates:
  ← Position from every EagleEye operator
  ← Situational awareness (threats, intel)
  ← Comms (voice integrated with position)
  ← Targeting data (coordinate strikes)
  
Command sees unified picture:
  All warfighters as nodes on COP
  Real-time team coordination
  Integrated fires + comms + movement
```

### EagleEye + Mumble Voice Integration
```
EagleEye audio path:
  Squad PTT (microphone on EagleEye)
  ↓ Mumble client (encrypted voice)
  ↓ Local Mumble server (air-gapped or connected)
  ↓ Squad channel (all hear comms)
  ↓ Overlay: see who's speaking (via AR display)

Enhancement: Who's speaking is highlighted on AR map
  - Real-time speaker identification
  - Know which teammate just called in contact
  - Instant awareness of team status changes
```

### EagleEye + DDF Intelligence Catalog
```
DDF metadata catalog (geospatial intelligence)
  ↓ REST API query from EagleEye
  ↓ Retrieve relevant intel for operational area
  ↓ Display in AR (enemy positions, terrain, restricted areas)
  ↓ Warfighter sees full intelligence picture overlaid on terrain

Example: "Show me all known enemy positions in 5km radius"
  Result: AR layer shows threat dots at known locations
  Warfighter can plan movement to avoid, or coordinate attack
```

### EagleEye + NATS Event Fabric
```
EagleEye publishes to NATS:
  - agents.tactical.warfighter.{id}.position (location stream)
  - agents.tactical.warfighter.{id}.status (battery, comms quality, readiness)
  - agents.tactical.warfighter.{id}.sensor (what EagleEye "sees" — threats detected)
  - agents.tactical.warfighter.{id}.comms (voice/data activity)

NATS subscribers:
  - TAK Server (updates COP in real-time)
  - Analytics (compute threat picture)
  - Alerting (commander notified of incidents)
  - Archival (record all ops for review)

Result: One unified event stream from edge warfighter to command fabric
```

### EagleEye in Air-Gapped Tactical Operations
```
Tactical TOC (air-gapped, local network):
  ├── KubeSolo compute (TAK Server running locally)
  ├── Mumble server (comms)
  ├── DDF catalog (local intelligence)
  └── NATS fabric (event bus)

EagleEye warfighters (in field):
  ├── Connect via WiFi/LAN to TOC
  ├── Receive mission brief (objectives, intel in AR)
  ├── Conduct ops (position sync, targeting, comms)
  ├── No external connectivity required
  ├── Full autonomy + team coordination
  
When comms available:
  → Sync mission results to cloud Legacy (mother)
  → Share intelligence / coordinate with other TOCs
  → Resume normal ops
```

---

## Partnership Model & Commercial Integration

### Strategic Partners

**1. Meta (Augmented Reality)**
- **Contribution:** AR tech, optical display, mixed-reality experience
- **Benefit to EagleEye:** State-of-the-art AR capabilities, proven consumer tech adapted to combat
- **Integration:** Meta's research in AR form factor + Anduril's tactical requirements
- **Cost reduction:** Leverage Meta's billion-dollar AR investment

**2. Oakley Standard Issue (Combat Eyewear)**
- **Contribution:** Combat eyewear design, ballistic protection, fit/function
- **Benefit to EagleEye:** Military-grade optics, durability tested in field
- **Integration:** Anduril compute + sensors integrated into Oakley eyewear frame
- **Cost reduction:** Use proven military eyewear design (lower risk)

**3. Qualcomm (Mobile Compute)**
- **Contribution:** Small form factor processors, power efficiency, AI/ML on-device
- **Benefit to EagleEye:** Leading-edge mobile SoCs optimized for real-time processing
- **Integration:** Snapdragon or custom Qualcomm chip in EagleEye compute module
- **Cost reduction:** Leverage Qualcomm's mobile chip R&D (billions in investment)

**4. Gentex (Ballistic Protection)**
- **Contribution:** Advanced materials, ballistic testing, helmet/face protection
- **Benefit to EagleEye:** Combat-proven protection materials
- **Integration:** Gentex protection standards applied to EagleEye eyewear
- **Cost reduction:** Use existing Gentex supply chain + testing

### Commercial Adaptation Strategy

**Why Partner with Commercial Leaders?**
- **Risk mitigation** — Use proven tech, not experimental
- **Cost reduction** — Commercial scale is cheaper than defense-only production
- **Speed to market** — Partners have existing R&D, manufacturing, supply chains
- **Continuous innovation** — Commercial updates drive warfighter capability improvements
- **Path forward** — Commercial advancements automatically flow to military versions

**Example: Smartphone ← Tactical Platform**
```
Consumer AR smartphone apps → Research + Development
     ↓
Meta invests $B in AR research
     ↓
Anduril takes those innovations → Adapt to tactical requirements
     ↓
Add military integration (TAK, encryption, ballistics)
     ↓
Qualcomm mobile chip → Fast, low-power, proven
     ↓
EagleEye uses same architecture (lower cost than custom chip)
     ↓
Commercial improvements → Automatically benefit EagleEye
     ↓
Continuous upgrade path (no obsolescence)
```

---

## Modular Architecture for Mission Adaptation

### Core Module (Always Included)
- **Display** — AR eyewear
- **Compute** — Mobile SoC
- **Comms** — Cellular + WiFi + radio
- **Sensors** — Camera, IMU, GPS, mic
- **Battery** — Field-life power

### Optional Add-ons (Mission-Specific)
- **Thermal module** — Swap in/out for night ops
- **Rangefinder** — Laser for distance measurement
- **Radar module** — Ground penetrating or short-range
- **Weapons integration** — Ballistic comp, fire control
- **Extra battery** — Extended ops
- **Rugged case** — Protective transport

### Software Plugins (Downloaded/Installed)
- **Target recognition** — AI for threat classification
- **Encryption suite** — Higher classification support
- **Offline maps** — Pre-load terrain for air-gapped ops
- **Custom apps** — Unit-specific mission software
- **Integration modules** — TAK, ATAK, DDF, custom systems

### Benefits of Modularity
✓ **Mission tailoring** — Squad ops ≠ sniper ops ≠ air support control  
✓ **Weight optimization** — Only carry needed sensors  
✓ **Cost control** — Don't pay for unused capabilities  
✓ **Rapid adaptation** — Swap hardware in hours (not months of re-engineering)  
✓ **Future-proof** — Upgrade individual modules as tech improves  

---

## Use Cases & Deployment Scenarios

### Scenario 1: Squad-Level Patrol (Urban COIN)
```
4-person squad, each with EagleEye:
  - See friendly positions (who's where on AR map)
  - Hear comms (Mumble PTT integrated)
  - Get real-time intel (DDF: known threat locations)
  - Navigate (heads-up routing via AR)
  - Call for support (comms to mortars/air)
  
Outcome: Coordinated squad, informed decision-making, rapid comms
```

### Scenario 2: Special Operations (Air-Gapped)
```
SEAL team, tactical ops center:
  - Air-gapped network (no external connectivity)
  - Mission brief uploaded to EagleEye (objectives, intel)
  - Operator sees targets in AR overlaid on terrain
  - Real-time position sync (no exfil position leakage)
  - Call in fire (laser targeting from EagleEye)
  - Record entire op (video + telemetry for review)
  
Outcome: Self-contained ops, full autonomy, after-action analysis
```

### Scenario 3: Multi-Domain Operations
```
Joint air/ground coordination:
  - EagleEye ground operator sees drone feed (ISR video)
  - Drone operator sees ground team positions
  - Common operating picture: all parties see same threat picture
  - Fast targeting cycle: ground spots threat → calls drone → strike
  
Outcome: Integrated fires, shortened kill chain
```

### Scenario 4: Forward Observer / Fire Control
```
FO with EagleEye:
  - Laser range finder (module add-on)
  - Target marked in AR (exact coordinates)
  - EagleEye sends targeting data (encrypted)
  - Artillery receives mission (via integrated comms)
  - Rounds on target in seconds
  
Outcome: Precision fires, reduced call-for-fire time
```

### Scenario 5: Training & Home Station
```
Unit trains with EagleEye (commercial AR gaming):
  - Load training scenarios (enemy positions, objectives)
  - Squad practices coordination
  - Comms procedure practice (Mumble on training net)
  - AI instructor provides feedback
  - Low cost (leverage commercial training software)
  
Outcome: Realistic prep, faster adoption, continuous learning
```

---

## Comparison with Alternatives

| Aspect | EagleEye | ATAK Tablet | Radio + Map | Existing HUD |
|--------|----------|-----------|-----------|--------------|
| **Hands-free** | ✓ Yes | ✗ Requires hands | ✗ No | ✓ Yes |
| **Real-time position** | ✓ Yes | ✓ Yes | ✗ No | ✗ No |
| **AR overlay** | ✓ Yes | ✗ No (flat map) | ✗ No | ✗ Limited |
| **Integrated comms** | ✓ Yes | ✓ Yes | ✓ Yes | ✗ No |
| **Targeting support** | ✓ Yes | ✗ Limited | ✗ No | ✗ No |
| **Night ops** | ✓ Yes (thermal) | ✓ Yes (if display) | ✗ No | ✓ Yes |
| **Weight** | Light | Medium (tablet) | Light | Light |
| **Battery life** | 8+ hrs | 4-6 hrs | N/A | N/A |
| **Cost** | $$ (high) | $ (lower) | $ (lowest) | $$ (high) |

---

## Procurement & Acquisition Model

### Defense Procurement Path
```
FY2024: Development contract (DoD / SOCOM)
  ↓
FY2025: Limited production (100-500 units)
  ↓
FY2026: Pilot deployment (select units)
  ↓
FY2027: Full-rate production (1000+ units/year)
  ↓
FY2028+: Fielding across force (integration with TAK Server, etc.)
```

### Commercial Partnerships as Force Multiplier
- **No need to re-invent AR display** — Use Meta's research
- **No need for custom mobile SoC** — Use Qualcomm's Snapdragon
- **No need for ballistic materials R&D** — Use Gentex's proven tech
- **No need for eyewear design** — Use Oakley's combat eyewear

**Result:** Faster development, lower cost, continuous innovation (commercial improvements auto-flow to military version)

### Cost Model (Estimated)
```
Commercial components (commercial pricing):
  - Qualcomm SoC: $200
  - AR optics (Meta): $300
  - Oakley eyewear frame: $300
  - Sensors (camera, IMU, etc.): $200
  - Battery + power mgmt: $100
  ──────────────────────────────
  Components total: ~$1,100

Anduril integration:
  - Tactical software stack: $400
  - Comms/encryption: $300
  - TAK/ATAK integration: $200
  - Manufacturing/assembly: $300
  ──────────────────────────────
  Integration: ~$1,200

Unit cost (R&D amortized + profit): ~$2,500-3,500 per unit
Military procurement cost: ~$5,000-8,000 per unit (with support)

Benefits: 10x cheaper than building all-new system from scratch
```

---

## Open Architecture & Extensibility

### Open Application System
- **Partner integrations** — Tactical app developers can write for EagleEye
- **Third-party apps** — Companies build mission-specific software
- **Government apps** — DoD agencies develop specialized capabilities
- **Continuous updates** — New apps roll out without new hardware

### Mission Command Architecture
- **API framework** — Standardized way to interact with EagleEye
- **Plug-in system** — Load custom comms, targeting, intel modules
- **Data interchange** — Share position/intel with other systems (TAK, DDF, etc.)
- **Backward compatibility** — New apps work with existing hardware

### Examples of Extensions
```
Third-party developer writes threat detection AI:
  - Trains on enemy vehicle imagery
  - Deploys as EagleEye plugin
  - Auto-labels vehicles in real-time
  - Warfighter sees colored icons (threat priority)
  
Another developer writes tactical comms relay:
  - Route voice through multiple networks
  - Auto-switch to best path
  - Transparent to user (seamless comms)
  
Another builds integration with commercial drone:
  - Syncs EagleEye position with drone
  - Displays drone view in EagleEye AR
  - Warfighter controls drone from EagleEye
  
All deployed via software update (no hardware change)
```

---

## Integration Roadmap with OneMind OS

### Phase 1: Core Tactical Integration (Year 1)
- [ ] EagleEye ↔ TAK Server sync (position, intel, tasking)
- [ ] Mumble comms integration (voice in EagleEye audio path)
- [ ] NATS client deployed to EagleEye (event publishing)
- [ ] DDF REST API client (query intelligence catalog)
- [ ] Basic mission command (receive objectives in AR)

### Phase 2: Advanced Features (Year 2)
- [ ] AI threat detection (trained on OneMind intel database)
- [ ] Autonomous routing (avoid surveillance, optimize for OPSEC)
- [ ] Mesh networking (EagleEye-to-EagleEye relay for extended range)
- [ ] ORP integration (3D tactical map visualization)
- [ ] Waveform integration (tactical radio compatibility)

### Phase 3: Force Multiplier (Year 3+)
- [ ] Swarm coordination (multiple EagleEye operators as unified force)
- [ ] ML-assisted targeting (AI predicts enemy movement)
- [ ] Sovereign identity integration (no external dependencies)
- [ ] Satellite backhaul (connectivity in denied areas)
- [ ] Quantum-resistant crypto (post-quantum NATS integration)

---

## Training & Readiness

### Operator Training (Individual Warfighter)
1. **Hardware familiarization** (30 min)
   - How to wear, adjust, power on/off
   - Battery swap, charging
   - Quick-reference controls

2. **Navigation & UI** (1 hour)
   - Heads-up display controls (eye tracking? hand gesture?)
   - Navigating between apps
   - Mission upload/download

3. **Tactical integration** (2 hours)
   - TAK/ATAK positioning
   - Comms (Mumble PTT)
   - Targeting procedures
   - Emergency modes (comms lost, battery low)

4. **Scenarios & drills** (4+ hours)
   - Live squad ops with EagleEye
   - Comms discipline practice
   - Failure mode recovery
   - Integration with supporting assets

**Total initial training: ~8 hours per operator**

### Unit-Level Integration
1. **Tech lead certification** (2 days)
   - System maintenance
   - Software updates
   - Troubleshooting
   - Security protocols

2. **Squad leader training** (1 day)
   - Tactical employment
   - Mission planning on EagleEye
   - Comms procedures
   - Leadership integration (how to use team feedback)

3. **Command team training** (1 day)
   - COP interpretation (what EagleEye data means)
   - Integration with TAK Server
   - Fire support coordination
   - Logistics/supply planning

---

## Security Considerations

### Operational Security (OPSEC)
- **Position discipline** — EagleEye transmits location; restrict in contested areas
- **Comms security** — Encrypted channels, but radio can be DF'd (direction finding)
- **Friendly fire prevention** — Position overlay prevents blue-on-blue (but not foolproof)
- **Cyber resilience** — System should operate in degraded/offline mode if compromised

### Equipment Security
- **Hardware tamper evidence** — Seals to detect unauthorized access
- **Encryption key management** — Secure provisioning + rotation
- **Device wipe** — Capability to remotely erase if captured
- **Cross-domain solution** — Can run classified + unclassified (with isolation)

### Personnel Security
- **User authentication** — Biometric (fingerprint, face) + PIN
- **Authorized users only** — Role-based access to classified features
- **Audit trail** — Log who accessed what data when
- **Training certification** — Only cleared personnel can operate

---

## Known Considerations & Trade-offs

### Advantages
✓ Real-time AR situational awareness (game-changer for COIN)  
✓ Hands-free operation (keep hands on weapon)  
✓ Integrated comms (talk + see same picture)  
✓ Commercial partnerships (faster, cheaper, continuous improvement)  
✓ Modular (adapt to any mission)  

### Challenges
⚠ **Weight** — Adds to combat load (but less than tablet + radio + map)  
⚠ **Battery life** — 8+ hours good, but not 24-hour ops  
⚠ **Connectivity** — Depends on network (fallback to local in denied areas)  
⚠ **Training burden** — Operators need to learn new system  
⚠ **Cost** — High per-unit cost (but less than equivalent tactical suite)  
⚠ **Adoption risk** — Untested on large scale (need pilot programs first)  

---

## Related Resources in Convergence

- **Samsung_S23_Tactical_Mobile_Platform.md** — Complementary mobile device (field communication hub)
- **Mumble.md** — Voice comms integration (EagleEye audio path)
- **TAK Server** (memory ref) — Mission command backend (EagleEye as client)
- **Mil-OSS.md** — Open source foundations (EagleEye builds on OSS)
- **NATS** (memory ref) — Event fabric (EagleEye publishes telemetry)
- **DDF.md** — Intelligence catalog (EagleEye queries for tactical intel)
- **KubeSolo.md** — Deploy TAK Server/Mumble locally (air-gapped TOC)

---

## Anduril's Product Ecosystem

**Anduril Platform Categories:**
- **Sea** — Naval/maritime domain
- **Land** — Ground warfare (EagleEye category)
- **Air** — Aerial systems
- **Space** — Space operations
- **Lattice** — Command & control / situational awareness backbone
- **Arsenal** — Effects (weapons, targeting, fires)

EagleEye is part of **Land** warfare + integration with **Lattice** (command backbone).

---

## Learning Resources

### Official
- **Product page:** https://www.anduril.com/eagleeye
- **Anduril site:** https://www.anduril.com/
- **Careers/tech talks:** Anduril newsroom + YouTube

### Related
- **TAK Server documentation** — For integration details
- **AR development** — Meta's mixed reality SDKs
- **Military comms** — NATO waveform standards (NETT Warrior references)
- **Tactical tech** — ATAK/ATAK CIV developer forums

---

**Captured:** 2026-07-09  
**Status:** In development / early deployment (not yet field-wide)  
**Relevance:** Next-gen tactical AR platform, warfighter-worn compute, integrated comms + targeting  
**Integration Priority:** Very High — potential cornerstone of OneMind edge warrior capability  
**Acquisition Timeline:** FY2024-2028 (pilot to full fielding)

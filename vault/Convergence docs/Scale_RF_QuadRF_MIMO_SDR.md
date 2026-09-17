# Scale RF QuadRF — 4x4 MIMO SDR Phased Array Kit

**Official Links:**
- **Crowd Supply Campaign:** https://www.crowdsupply.com/scale-rf/quadrf
- **Updates Page:** https://scalerf.com/updates/
- **Discord Community:** https://discord.gg/scalerf (via ScaleRF site)
- **Documentation:** https://scalerf.com/docs/
- **GitHub:** https://github.com/scalerf/ (Open-source SDR tools)

**Company:** Scale RF

**Type:** 4x4 MIMO Software-Defined Radio (SDR) tile with phased array capabilities

**Status:** Crowd Supply campaign LIVE (319% funded as of 2026-07-11)

**Funding Goal:** $100,000 (ACHIEVED in 22 hours)

**Campaign Period:** June 23 - August 6, 2026

**Price:** $499 - $594 USD

**Backers:** 484+ (as of 2026-07-11)

**Estimated Delivery:** First batch of 1000 units

---

## What Is QuadRF?

**QuadRF** is a **4x4 MIMO Software-Defined Radio tile** that brings **phased-array technology** into an accessible, affordable, education and development kit.

**Core Components:**
- 4 coherent full-duplex RF chains
- Dual-polarization antennas (4 total antennas)
- Integrated Raspberry Pi 5 (quad-core ARM A76)
- Real-time RF processing and visualization
- Open-source antenna architecture
- Web browser control interface
- Preloaded Linux/SDR software stack

**Key Innovation:** First mass-market product to move **software-defined radio beyond time and frequency into the spatial domain** — enabling real-time **RF spatial vision** (phased-array RF sensing).

---

## Core Capabilities

### 1. RF Augmented Reality (RF Vision)

**Real-time RF Camera Mode** (30 fps)
- Mount smartphone on mobile handle
- View live RF overlay directly on phone/laptop
- Visualize wireless signals as color-coded radio sources
- Pinpoint exact physical location of transmitters
- Detect hidden wireless cameras, rogue access points (through walls)

**Technology:**
- 4 coherent antennas measure signal arrival time differences
- Time-difference-of-arrival (TDOA) triangulation in real-time
- Augmented-reality heat-map overlay

### 2. Spatial RF Sensing

**Applications:**
- **Visualize Radiation:** Pinpoint transmitter locations precisely
- **Observe Environment:** Watch RF reflections, shadowing, absorption as people/materials move
- **Empower Robotics:** Give machines real-time spatial RF awareness (beacons, local infrastructure)
- **Communicate Long Range:** Beamforming for signal gain + interference rejection (moving relay links)
- **Develop Mesh Networks:** High-bandwidth wireless relays (point-to-point, point-to-multipoint with MIMO)
- **Sense Near-field:** Transmit probe signals to measure spatial responses (distance, absorption, reflection, polarization)
- **Test Antennas:** Instantly visualize how polarization/placement changes RF field
- **Enhance Education:** Teaching tool for MIMO, phased arrays, spatial wireless measurement

### 3. Phased Array Foundation

**Scalable Architecture:**
- Single QuadRF = complete 4x4 MIMO phased array
- Multiple QuadRFs interconnect via open-source structures
- Gateway to larger phased arrays (MoonRF = 240 antennas)
- Designed for tile-based manufacturing scaling

**Array Modes:**
- Single tile: 4x4 MIMO (4 TX, 4 RX simultaneously)
- Multiple tiles: Larger phased arrays with coherent phase alignment
- Networked mode: Distribute I/Q streams via Ethernet/USB/Wi-Fi

---

## Hardware Specifications

### RF Performance

| Spec | Value |
|------|-------|
| **Frequency Range** | Covers ISM + licensed bands (exact range TBD) |
| **Bandwidth** | 4 × 40 MHz full-duplex per chain (160 MHz total) |
| **MIMO Configuration** | 4x4 (4 transmit, 4 receive, full-duplex) |
| **Antennas** | 4 integrated dual-polarization antennas |
| **Phase Coherence** | Fully coherent across 4 chains (phased-array capable) |
| **Output Power** | TBD (low-power, no ham license required) |

### Computing

| Component | Spec |
|-----------|------|
| **SoC** | Raspberry Pi 5 (BCM2712) |
| **CPU** | Quad-core ARM A76 @ 2.4 GHz |
| **RAM** | 8 GB |
| **Storage** | microSD card (preloaded Linux + SDR tools) |
| **Connectivity** | Gigabit Ethernet, USB 3.0, Wi-Fi 6 |

### Interface

| Port | Capability |
|------|-----------|
| **Web Interface** | http://quadrf/ (browser-based control) |
| **Ethernet** | Gigabit (I/Q streaming, remote control) |
| **USB 3.0** | High-speed I/Q sample streaming |
| **Wi-Fi 6** | Remote control + I/Q streaming |
| **HDMI** | Local display output |

---

## Software Stack

### Preloaded on microSD Card

**100% Open-Source (GPLv2)**

1. **Linux Kernel** (Real-time patches available)
2. **Device Drivers** — RF chain control, calibration, antenna switching
3. **SDR Libraries** — SoapySDR, GNU Radio, Pothos
4. **Visualization Tools** — Real-time RF camera software
5. **Example Applications**:
   - Wi-Fi analyzer
   - Beacon tracker
   - RF signal decoder
   - Beamforming demonstrations
   - GNU Radio flowgraph examples

### Integration with External Tools

**Local Processing (on Pi 5):**
- GNU Radio
- Pothos
- QRadioLink
- Custom Python/C++ applications
- Real-time visualization

**Remote Processing (via Ethernet/USB/Wi-Fi):**
- I/Q sample streaming to:
  - SDR# (Windows)
  - SDRangel
  - GQRX
  - CubicSDR
  - Laptop-based signal processing

---

## Key Features

### 1. Plug-and-Play Setup

- Power on, connect Wi-Fi/Ethernet/USB
- Open browser to http://quadrf/
- Start experimenting in seconds
- Preloaded drivers + examples on microSD

### 2. Real-time Calibration

**Automated calibration for:**
- Phase alignment (4-channel coherence)
- Quadrature (I/Q accuracy)
- LO leakage (local oscillator suppression)
- Antenna impedance matching

### 3. Full-Duplex 4x4 MIMO

**Capabilities:**
- Transmit & receive simultaneously on same frequency
- 4 independent data streams (4 TX, 4 RX)
- Enable MIMO communication + spatial sensing
- Foundation for beamforming

### 4. Portable Mobile Kit

**Included:**
- QuadRF tile
- Raspberry Pi 5
- Mobile handle (tripod-mountable)
- Smartphone mount bracket
- Preloaded microSD card
- Cables, antennas, documentation

**Weight:** Approx. 2-3 kg with handle (portable for field deployment)

---

## Scalability: From Quad to Moon

### The QuadRF Building Block Architecture

**Single QuadRF (4 antennas):**
- Complete 4x4 MIMO SDR kit
- Functional phased array on its own
- Educational + prototyping platform

**Multiple QuadRFs (scale up):**
- 8 QuadRFs = 8x4 MIMO (32 antennas)
- 16 QuadRFs = 16x4 MIMO (64 antennas)
- **N QuadRFs = N × 4 antennas** (linear scaling)

**Open-Source Interconnect Structures:**
- RF board can be removed from enclosure
- Mount into open-source mechanical/RF structures
- Maintain full phase coherence across tiles
- Distribute I/Q streams via optical fiber (for large arrays)

### MoonRF Connection

**MoonRF = 240-antenna moon-bounce array**
- Built from scaled-up QuadRF architecture
- Each QuadRF tile is a building block
- QuadRF kit allows people to prototype MoonRF concepts locally
- First 1000 QuadRF units shipped → foundation for larger arrays

---

## Applications & Use Cases

### Education

- **University Teaching:** MIMO, phased arrays, spatial wireless measurement (hands-on lab)
- **Signal Processing:** Real-time visualization of RF phenomena
- **Antenna Theory:** Instant feedback on polarization + placement effects

### Commercial RF Development

- **Wi-Fi Analysis:** Visualize Wi-Fi signal propagation
- **Spectrum Monitoring:** Detect interference, unauthorized transmitters
- **Antenna Testing:** Rapid iteration on antenna designs
- **Beamforming R&D:** Develop beam steering algorithms
- **Mesh Networking:** Build + test multi-hop wireless relays

### Autonomous Systems

- **Drone Coordination:** RF-based ranging + bearing to other UAVs
- **Robotics:** Spatial RF awareness (beacon-based localization)
- **Swarm Operations:** Coordinate multiple robots via RF signals

### Tactical Applications

- **Direction Finding:** Locate RF emitters (covert surveillance detection)
- **Electronic Warfare:** Signal detection + analysis
- **Mesh Network Deployment:** Build resilient ad-hoc networks
- **Remote Communication:** Long-range beamformed links for harsh environments

---

## Campaign Status & Timeline

### Crowd Supply Campaign

**Campaign Start:** June 23, 2026

**Achievement Timeline:**
- **22 hours:** Funding goal ($100,000) achieved (319% total)
- **484+ backers** (as of July 11, 2026)
- **$319,813 raised** (as of July 11, 2026)
- **26 days remaining** (campaign ends August 6, 2026)

### Recent Updates (from scalerf.com/updates/)

#### July 9, 2026 — "Plans for Injection Molding, Wi-Fi Analysis, Teardown 2026 in July, and GRCon in September"
- Manufacturing prep underway
- Wi-Fi analyzer app in development
- Teardown 2026 demo scheduled (July)
- GNU Radio Conference 2026 talk (September)

#### June 30, 2026 — "This Product Is Happening!"
- Campaign crossed funding goal in record time (22 hours)
- Component pre-ordering to manage supply-chain risks
- First batch of 1000 units secured
- 100% delivery success rate expected (Crowd Supply track record)
- Updates through August 6, then production + shipping

#### June 23, 2026 — "Our Campaign Is Now Live!"
- Crowd Supply pre-launch accepted June 2
- Full campaign went live June 23
- Immediate positive response

### Production & Delivery

**Plan:**
- First batch: 1000 units (securing components now)
- Assembly: Q3/Q4 2026
- Shipping: Expected Q4 2026 / Q1 2027
- Crowd Supply handles: Distribution, shipping, VAT, export compliance

**Risk Mitigation:**
- Pre-ordering components to avoid supply-chain delays
- Crowd Supply 100% historical delivery success rate
- Manufacturing experience from MoonRF development

---

## Purchasing Options

### Tiers (on Crowd Supply)

| Tier | Price | Includes |
|------|-------|----------|
| **QuadRF Kit** | $499 | 1 × QuadRF + Pi 5 + all accessories |
| **QuadRF Kit (Early Bird)** | $449 (if available) | Same as above |
| **Multiple Units** | $594/unit (bulk discount) | Multiple kits for phased array experiments |

**Funding Link:** https://www.crowdsupply.com/scale-rf/quadrf

---

## Integration with OneMind/Convergence

### RF Sensing Layer for Autonomous Systems

**Pattern Recognition:**
- **PX4/ArduPilot (autopilots)** + **QuadRF (RF sensing)** = spatial-aware autonomous vehicles
- **DJI O4 (RF infrastructure)** ↔ **QuadRF (local RF camera)** = multi-layer RF network visualization
- **TAK Server (C2)** → **QuadRF (spatial RF data)** = real-time RF threat picture

### Mesh Network Foundation

**OneMind NATS mesh + QuadRF beamforming:**
- Robust RF links with beamforming + interference rejection
- Real-time RF visualization (who's transmitting, where, on what frequency)
- Scalable MIMO capability (multiple nodes on same frequency, full-duplex)

### Tactical Applications

1. **Direction Finding:** Locate hostile RF emitters (covert detection)
2. **Electronic Intelligence:** Spectrum analysis + threat identification
3. **Resilient Comms:** Mesh networks with beamforming (robust despite jamming/interference)
4. **Drone Coordination:** RF-based relative positioning between UAVs
5. **Swarm Awareness:** Spatial RF sensing for autonomous robot teams

---

## Why QuadRF Matters for Convergence

### Technical Significance

1. **First mass-market coherent MIMO SDR** — Moves RF from frequency-only to spatial domain
2. **Phased array democratization** — $500 for technology that cost $50k+ before
3. **Open-source foundation** — All software, antenna files, mechanical designs released
4. **Raspberry Pi integration** — Edge compute + RF sensing in single device
5. **Real-time visualization** — Makes RF phenomena tangible (RF augmented reality)

### Strategic Significance

- **Low-cost RF infrastructure** for distributed autonomous systems
- **Spatial awareness** for swarms, drones, mesh networks
- **Tactical resilience** — Beamforming + MIMO = interference-resistant comms
- **Scalable architecture** — From 4 antennas to 240+ (modular phased arrays)
- **Open-source sovereignty** — No vendor lock-in, full transparency

### Convergence Connection

**RF Sensing + OneMind Mesh:**
- NATS (NATS meshing)
- PLCs (deterministic control)
- TAK (C2 / COP)
- **+ QuadRF (spatial RF vision + resilient comms)**
= Complete autonomous intelligent infrastructure with electromagnetic awareness

---

## Resources

### Official Links

- **Campaign:** https://www.crowdsupply.com/scale-rf/quadrf
- **Updates:** https://scalerf.com/updates/
- **Docs:** https://scalerf.com/docs/
- **Discord:** https://discord.gg/scalerf
- **GitHub:** https://github.com/scalerf/

### Community & Learning

- **GNU Radio:** https://www.gnuradio.org/
- **SoapySDR:** https://github.com/pothosware/SoapySDR
- **Raspberry Pi:** https://www.raspberrypi.org/
- **RF Engineering:** MIMO fundamentals via book + academic papers

### Related Standards

- **IEEE 802.11 (Wi-Fi):** Visible with RF camera
- **Cellular Bands:** Detectable (ISM + licensed bands)
- **Phase Array Principles:** Open-source documentation

---

## Summary

**QuadRF is a $500 kit that brings phased-array RF sensing to commodity hardware.**

Key achievements:
- ✅ 4x4 MIMO SDR tile (coherent phase-aligned channels)
- ✅ Real-time RF augmented reality (30 fps spatial visualization)
- ✅ Integrated Raspberry Pi 5 (on-device signal processing)
- ✅ Open-source software stack (GPLv2, all drivers + examples)
- ✅ Scalable architecture (modules for 240+ antenna arrays)
- ✅ Crowd-funded success (319% funded, 484+ backers)
- ✅ Estimated delivery Q4 2026 / Q1 2027

**For OneMind Convergence:**
- RF sensing layer for autonomous systems
- Beamforming + MIMO for resilient mesh networks
- Spatial awareness for drone swarms + robot teams
- Open-source sovereignty (no vendor lock-in)
- Tactical resilience (interference-resistant communications)

---

**Captured:** 2026-07-11

**Source:** https://www.crowdsupply.com/scale-rf/quadrf + https://scalerf.com/updates/

**Status:** Production campaign active (ends August 6, 2026)

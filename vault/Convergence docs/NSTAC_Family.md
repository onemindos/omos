# NSTAC Family: Protocol-to-CoT Converters (ADSBcot, AISCot, PyTAK)

**Project Type:** Lightweight protocol-to-Cursor on Target (CoT) converters  
**License:** Apache 2.0 (open source, commercial-friendly)  
**Architecture:** No fusion, no NATS, direct TAK integration  
**Maintained by:** NSTAC (National Security Telecommunications Advisory Committee) community + open source contributors  
**Repository:** GitHub (NSTAC projects)  
**Status:** Production-ready, actively used in tactical operations  
**Primary Use:** Bridge external data streams directly into TAK/CloudTAK without middleware  

---

## What It Is: The NSTAC Family Overview

The **NSTAC family** comprises lightweight, purpose-built protocol converters that translate real-world sensor data formats into Cursor on Target (CoT) XML format for direct consumption by TAK (Team Awareness Kit) systems.

**Core Philosophy:**
- **No fusion layer** — Pass data through unchanged
- **No NATS dependency** — Direct TAK integration
- **Minimal processing** — Convert format only, preserve data fidelity
- **Resilience bypass** — Work even when ORP/central infrastructure is down
- **Plug-and-play** — Drop into TAK ecosystem immediately

**What each converter does:**

| Converter | Input | Output | Use Case |
|-----------|-------|--------|----------|
| **ADSBcot** | ADS-B (aircraft) | CoT position messages | Track aircraft |
| **AISCot** | AIS (maritime) | CoT position messages | Track ships |
| **PyTAK** | TAK protocol library | CoT/TAK streams | Custom converters |

---

## Core Components

### 1. ADSBcot: Aircraft Tracking (ADS-B to CoT)

#### What It Does

**Converts ADS-B (Automatic Dependent Surveillance-Broadcast) data into TAK CoT format.**

```
Real-world data flow:

Aircraft with ADS-B transponder
  ↓ (broadcasts position, altitude, speed, callsign)
ADS-B receiver hardware (RTL-SDR, USRP, commercial)
  ↓ (receives radio signals on 1090 MHz)
decode-bin (extracts frames, decodes positions)
  ↓ (outputs JSON/raw format)
ADSBcot converter
  ↓ (transforms to CoT XML)
TAK Server
  ↓
TAK/ATAK clients see aircraft on map
```

#### Key Features

**Aircraft data captured:**
- Aircraft position (latitude, longitude, altitude)
- Speed and heading
- Callsign (flight identifier: N12345, BA123, etc.)
- Altitude and climb rate
- ICAO address (unique aircraft identifier)
- Timestamp (position time)

**CoT output:**
```xml
<event version="2.0" uid="ICAO-8D406F" type="a-f-A-M-F-U" time="2026-07-09T12:34:56Z" start="2026-07-09T12:34:56Z" stale="2026-07-09T12:35:56Z" how="m-g">
  <point lat="40.7128" lon="-74.0060" hae="3000" ce="100" le="100"/>
  <contact callsign="N12345"/>
  <track speed="450" course="270"/>
</event>
```

#### Installation & Usage

**Setup (Ubuntu):**
```bash
# Install dump1090 (ADS-B decoder)
git clone https://github.com/flightaware/dump1090.git
cd dump1090
make
./dump1090 --interactive

# In another terminal, install ADSBcot
git clone https://github.com/nstac/adsbcot.git
cd adsbcot
pip install -r requirements.txt

# Run ADSBcot (connects to dump1090 + TAK Server)
python adsbcot.py \
  --host localhost \
  --port 5000 \
  --url tcp://tak-server:8087 \
  --tls-ca /path/to/ca.pem \
  --tls-cert /path/to/cert.pem \
  --tls-key /path/to/key.pem
```

**Hardware setup:**
```
Antenna → RTL-SDR receiver (USB stick, $20) 
          ↓
        dump1090 (software decoder)
          ↓
      localhost:5000 (raw ADS-B data)
          ↓
      ADSBcot (convert to CoT)
          ↓
      TAK Server (8087 TLS)
          ↓
      ATAK on devices (see aircraft on map)
```

#### Tactical Applications

**Scenario 1: ISR Asset Tracking**
```
Friendly UAVs transmit ADS-B signal
  → ADSBcot converts to CoT
  → TAK displays all friendly air assets
  → Commanders see real-time UAV positions
  → Coordinate with ground forces

Benefit: Automatic position update (no manual entry)
         All devices see same UAV picture
         Position accurate to ±20 feet (ADS-B accuracy)
```

**Scenario 2: Airspace Awareness**
```
Commercial and civilian aircraft in operational area
  → ADSBcot tracks all aircraft on TAK
  → Operators avoid friendly fire (know airspace is occupied)
  → Correlate aircraft with flight plans
  → Identify threats vs. friendlies

Benefit: Situational awareness of entire airspace
         Prevent blue-on-blue (aircraft collision)
         Know who owns the sky
```

**Scenario 3: Resilience Bypass (No ORP/Central Infrastructure)**
```
Scenario: ORP (central data processing) is down
          TAK Server still running locally

Without ADSBcot:
  - Lose aircraft tracking capability
  - No way to get position data into TAK
  
With ADSBcot:
  - Connect ADS-B receiver directly to TAK
  - Aircraft positions flow directly to TAK
  - Bypass ORP entirely
  - TAK continues operating at full capability

Result: Resilient, degraded-mode operations
        No single point of failure
```

---

### 2. AISCot: Maritime Tracking (AIS to CoT)

#### What It Does

**Converts AIS (Automatic Identification System) data into TAK CoT format.**

```
Real-world data flow:

Ship with AIS transponder
  ↓ (broadcasts position, course, speed, call sign, cargo)
AIS receiver hardware (software-defined radio, commercial)
  ↓ (receives radio signals on 161.975 MHz / 162.025 MHz)
AIS decoder (extracts frames, decodes ship positions)
  ↓ (outputs JSON or NMEA format)
AISCot converter
  ↓ (transforms to CoT XML)
TAK Server
  ↓
TAK clients see ships on map
```

#### Key Features

**Ship data captured:**
- Ship position (latitude, longitude)
- Course (heading)
- Speed (SOG - Speed Over Ground)
- Call sign (international maritime identifier)
- MMSI (Maritime Mobile Service Identity - unique ship ID)
- Ship type (cargo, tanker, fishing, passenger, etc.)
- Destination port
- Timestamp

**CoT output:**
```xml
<event version="2.0" uid="MMSI-211378040" type="a-f-S-U-C-N" time="2026-07-09T12:34:56Z" start="2026-07-09T12:34:56Z" stale="2026-07-09T12:35:56Z" how="m-g">
  <point lat="40.7000" lon="-74.0100" hae="0" ce="100" le="100"/>
  <contact callsign="PRINCESS OF THE SEAS"/>
  <track speed="15" course="270"/>
  <remarks>Ship type: Passenger, Destination: NY</remarks>
</event>
```

#### Installation & Usage

**Setup (Ubuntu):**
```bash
# Install rtl-sdr (software-defined radio for AIS)
sudo apt install rtl-sdr

# Install AIS decoder (gr-osmosdr or similar)
git clone https://github.com/dgiardina/rtl-ais.git
cd rtl-ais && make

# Run AIS receiver on 161.975 MHz and 162.025 MHz (AIS frequencies)
./rtl_ais -d 0 -f 161975000 -f 162025000

# Install AISCot
git clone https://github.com/nstac/aiscot.git
cd aiscot
pip install -r requirements.txt

# Run AISCot
python aiscot.py \
  --host localhost \
  --port 10001 \
  --url tcp://tak-server:8087 \
  --tls-ca /path/to/ca.pem
```

**Hardware setup:**
```
Antenna (VHF marine band) → RTL-SDR receiver (USB)
                              ↓
                          rtl-ais (decode AIS)
                              ↓
                        localhost:10001 (AIS data)
                              ↓
                          AISCot (convert to CoT)
                              ↓
                          TAK Server (8087 TLS)
                              ↓
                      ATAK on devices (see ships)
```

#### Tactical Applications

**Scenario 1: Naval Awareness**
```
Friendly naval assets + commercial maritime traffic
  → AISCot tracks all ships on TAK
  → Commanders see maritime picture
  → Identify merchant vessels vs. warships
  → Coordinate naval operations

Benefit: Real-time ship tracking (AIS broadcasts every 10 seconds)
         Know shipping lanes and merchant traffic
         Identify anomalies (ship in restricted area, dark ship)
```

**Scenario 2: Port Security**
```
Monitor port entrance for unauthorized vessels
  → AISCot tracks all ships entering/leaving
  → Compare against whitelist (authorized vessels)
  → Flag unknown ships for inspection
  → Real-time security monitoring

Benefit: Automatic threat detection
         No manual tracking needed
         Every ship logged and time-stamped
```

**Scenario 3: Search & Rescue**
```
Scenario: Aircraft ditches in ocean, survivors in water

Search area established
  → AISCot shows all ships in area (potential rescue)
  → SAR coordinator contacts nearest ships
  → Dispatch rescue helicopter
  → Coordinate multi-asset rescue

Benefit: Know exactly which ships are nearby
         Fast coordination of rescue assets
         Real-time ship positions on COP
```

---

### 3. PyTAK: TAK Protocol Library

#### What It Is

**Python library for creating custom TAK clients and converters.**

PyTAK provides the building blocks to:
- Create TAK CoT XML messages
- Connect to TAK Server (TCP/TLS)
- Subscribe to CoT streams
- Build custom converters for any protocol

#### Installation

```bash
pip install pytak
```

#### Basic Usage

```python
import pytak
import asyncio
from datetime import datetime

async def main():
    # Create TAK client
    client = pytak.TAKClient(
        event_queue=asyncio.Queue(),
        server='tak-server',
        port=8087,
        tls_client_cert='cert.pem',
        tls_client_key='key.pem',
        tls_ca_cert='ca.pem'
    )
    
    # Create CoT event (point)
    event = pytak.CoT(
        uid='custom-marker-123',
        cot_type='a-f-G-U-S',  # Ground unit, friendly
        lat=40.7128,
        lon=-74.0060,
        ce=100,  # Circular error
        le=100   # Linear error
    )
    event.set_callsign('Custom Unit')
    event.set_speed(25.5)
    event.set_course(270)
    
    # Queue event for transmission
    await client.event_queue.put(event)
    
    # Run client (connect and send)
    await client.run()

asyncio.run(main())
```

#### Custom Converter Template

```python
import pytak
import asyncio
import json

class CustomProtocolClient(pytak.TAKClient):
    """Convert custom protocol to CoT"""
    
    async def handle_custom_data(self, data):
        """Convert incoming custom data to CoT"""
        
        # Parse custom protocol
        msg = json.loads(data)
        
        # Create CoT event
        event = pytak.CoT(
            uid=f"device-{msg['id']}",
            cot_type='a-f-G-U-S',
            lat=msg['latitude'],
            lon=msg['longitude'],
            ce=msg.get('accuracy', 100),
            le=100
        )
        event.set_callsign(msg['name'])
        event.set_speed(msg.get('speed', 0))
        event.set_course(msg.get('heading', 0))
        
        # Queue for transmission
        await self.event_queue.put(event)

async def main():
    client = CustomProtocolClient(...)
    # Listen for custom data, convert, and send to TAK
    await client.run()

asyncio.run(main())
```

#### Use Cases

**Build custom converters for any protocol:**
```
- UAV telemetry (ArduPilot, PX4)
- Sensor data (Lidar, radar, sonar)
- Third-party systems (proprietary formats)
- Legacy systems (old communication protocols)
- IoT devices (temperature, pressure, etc.)

Result: Anything that moves can be tracked in TAK
```

---

## Architecture: How NSTAC Converters Work

### Converter Pipeline

```
External Data Source
    ↓
Receiver/Decoder Hardware
    ↓
Protocol Stream (ADS-B, AIS, custom JSON, etc.)
    ↓
Converter (ADSBcot, AISCot, or PyTAK-based)
    │
    ├─ Parse input format
    ├─ Extract key fields (position, speed, heading, ID)
    ├─ Create CoT XML event
    └─ Timestamp and finalize
    ↓
TAK Server Connection (TCP/TLS)
    ├─ Authenticate (certificates)
    ├─ Send CoT XML
    └─ Maintain connection
    ↓
TAK Server
    ├─ Receives CoT
    ├─ Stores in database
    ├─ Broadcasts to connected clients
    └─ Updates COP
    ↓
ATAK/TAK Clients
    ├─ Receive CoT updates
    ├─ Display on map
    └─ Update positions in real-time
```

### Key Design Principles

**1. No Fusion**
```
Converter passes data through unchanged:
  - No combining of sources
  - No conflict resolution
  - No confidence weighting
  - Data fidelity preserved

Advantage: Fast, simple, predictable
Disadvantage: Duplicates/conflicts not resolved upstream
Solution: Let TAK or application layer do fusion
```

**2. No NATS**
```
Direct connection to TAK Server:
  - Converter → TAK Server (only connection)
  - No NATS dependency
  - No message bus
  - No extra middleware

Advantage: Simple, direct, low latency
Disadvantage: No other services see the data
Solution: TAK Server can broadcast via NATS if needed
```

**3. Direct TAK Integration**
```
Converter speaks TAK protocol natively:
  - Sends CoT XML over TCP/TLS
  - No adapter layer needed
  - No translation layer
  - Plug-and-play into TAK ecosystem

Advantage: Works immediately with any TAK system
Disadvantage: TAK-specific (not multi-format)
```

**4. Resilience Bypass**
```
Scenario: Central infrastructure (ORP) down

Traditional path:
  Data → Central processor → TAK Server

NSTAC path (resilience):
  Data → Converter → TAK Server (direct)

Result:
  - ORP outage doesn't affect data ingestion
  - Tactical operations continue
  - Distributed architecture, no single point of failure
```

---

## Integration with OneMind OS / Tactical Operations

### Scenario 1: Resilience Bypass (ORP Down)

```
OneMind OS normal operation:
  Sensors → NATS → Legacy (ORP) → TAK Server → ATAK

Failure: ORP down (component failure, network issue)
  Sensors → ??? → TAK Server → ATAK (broken)

NSTAC resilience bypass:
  ADS-B sensors → ADSBcot → TAK Server → ATAK (works!)
  AIS sensors → AISCot → TAK Server → ATAK (works!)
  
Result: Core operational capability preserved
        Graceful degradation
        No tactical impact from ORP failure
        TAK continues operating
```

### Scenario 2: Multi-Domain Asset Tracking

```
Tactical scenario: Coordinate across air, sea, ground, cyber

Assets:
  - Aircraft (ADS-B tracked by ADSBcot)
  - Ships (AIS tracked by AISCot)
  - Ground vehicles (GPS sent via PyTAK converter)
  - Drones (telemetry via PyTAK custom converter)

All converters feed TAK Server directly:
  ADSBcot → TAK → Unified COP
  AISCot  → TAK → Unified COP
  PyTAK (vehicles) → TAK → Unified COP
  PyTAK (drones) → TAK → Unified COP

Result: All assets visible on single COP
        No fusion layer needed
        Direct data flow
        Low latency
        Commanders see complete picture
```

### Scenario 3: Backup TAK Server (Tactical Network)

```
Scenario: Deploy independent TAK network (no central ORP)

Architecture:
  ┌─ ADSBcot ──┐
  │            │
  ├─ AISCot ───┤ → TAK Server (local) → ATAK on devices
  │            │
  └─ PyTAK ────┘

Setup:
  1. Deploy TAK Server in tactical TOC
  2. Run converters on same network
  3. Connect receivers (ADS-B, AIS, sensors)
  4. Operators access TAK via WiFi or cellular
  5. No dependency on central ORP
  6. No dependency on OneMind cloud

Result:
  - Independent tactical network
  - Works in air-gapped/denied environment
  - Direct data flow to operators
  - Minimal infrastructure
  - Rapid deployment
```

### Scenario 4: Custom Sensor Integration (PyTAK)

```
Scenario: Integrate custom ground sensors into TAK

Sensor types:
  - Seismic sensors (earthquake detection)
  - Pressure sensors (building monitoring)
  - Temperature sensors (environmental)
  - Custom tactical radars
  - Proprietary threat detection systems

PyTAK converter:
  1. Read sensor data (serial port, network API, file)
  2. Parse sensor format
  3. Create CoT event
  4. Send to TAK Server
  
Example:
  Seismic network → Custom PyTAK converter → TAK → Earthquake mapped

Result:
  - Any sensor data integrates with TAK
  - Operators see all sensor feeds on COP
  - Unified tactical picture
  - Custom systems supported
```

---

## Comparison: NSTAC vs. OneMind Full Integration

| Aspect | NSTAC Converters | OneMind Full Stack |
|--------|------------------|-------------------|
| **Setup time** | Minutes (hardware + software) | Hours (full deployment) |
| **Dependency** | TAK Server only | TAK + NATS + ORP + CloudTAK |
| **Resilience** | Direct bypass | Through ORP |
| **Failure impact** | Component down (ADS-B) ≠ whole system | ORP down = whole system affected |
| **Latency** | Minimal (converter → TAK direct) | Slightly higher (through NATS/ORP) |
| **Flexibility** | Purpose-built converters | Universal architecture |
| **Fusion** | None (pass-through) | Full fusion + correlation |
| **Use case** | Resilience, backup, rapid deploy | Production, full integration |

---

## NSTAC Family Ecosystem

### ADSBcot

**Purpose:** Aircraft tracking  
**Input:** ADS-B signals (1090 MHz)  
**Hardware:** RTL-SDR receiver ($20-50)  
**Output:** CoT position messages  
**Use:** ISR assets, airspace awareness, collision avoidance  

### AISCot

**Purpose:** Maritime tracking  
**Input:** AIS signals (VHF marine band, 161.975/162.025 MHz)  
**Hardware:** RTL-SDR + marine antenna ($50-200)  
**Output:** CoT position messages  
**Use:** Naval operations, port security, vessel tracking  

### PyTAK

**Purpose:** TAK protocol library + custom converters  
**Input:** Any protocol (JSON, binary, serial, network)  
**Hardware:** Any (software library)  
**Output:** CoT XML messages  
**Use:** Custom sensor integration, legacy system bridges  

### Other NSTAC Projects (Common)

- **goatak** — Go-based TAK client
- **FreeTAKServer** — Open-source TAK Server alternative
- **TAK Product Center** — Official TAK documentation
- **ATAK Plugins** — Extensions for ATAK

---

## Deployment: Quick Start Guide

### 1. Hardware Setup (ADS-B Example)

```bash
# Step 1: Assemble hardware
Materials:
  - RTL-SDR USB receiver (~$25, Amazon)
  - Dipole antenna or discone antenna
  - USB cable (usually included)

# Step 2: Install software
sudo apt update && sudo apt install -y rtl-sdr dump1090

# Step 3: Test ADS-B reception
rtl_fm -f 1090000000 -s 2000000 - | aplay
# You should hear ADS-B transponder signals

# Step 4: Run dump1090
dump1090 --interactive --net

# Step 5: Verify position data
# Browser to http://localhost:8080 (should show aircraft map)
```

### 2. Converter Setup

```bash
# Step 1: Install ADSBcot
git clone https://github.com/nstac/adsbcot.git
cd adsbcot
pip install -r requirements.txt

# Step 2: Configure connection
cat > config.ini << EOF
[adsbcot]
host = localhost
port = 5000
tls_ca = /path/to/ca.pem
tls_cert = /path/to/cert.pem
tls_key = /path/to/key.pem
tak_server = 192.168.1.100
tak_port = 8087
EOF

# Step 3: Run converter
python adsbcot.py --config config.ini

# Step 4: Verify on TAK
# Open ATAK/TAK on device
# Should see aircraft appearing on map in real-time
```

### 3. Resilience Architecture

```
Main OneMind OS:
  Sensors → NATS → Legacy (ORP) → TAK → ATAK
  
Backup/Resilience layer:
  ADS-B receiver → ADSBcot → TAK → ATAK (direct)
  AIS receiver → AISCot → TAK → ATAK (direct)
  
Benefit: If ORP fails, direct converters keep tactical ops going
```

---

## Performance & Limitations

### Performance Characteristics

```
ADSBcot (ADS-B):
  - Latency: 1-3 seconds (ADS-B broadcast + reception + conversion)
  - Accuracy: ±20 feet (95% confidence)
  - Update rate: Every 0-2 seconds (depends on aircraft)
  - Range: 100-250 nm (RTL-SDR receiver with good antenna)
  - Scalability: 1000+ aircraft tracked simultaneously

AISCot (AIS):
  - Latency: 5-60 seconds (AIS broadcast every 10-600 seconds depending on speed)
  - Accuracy: ±10 meters (GPS-derived)
  - Update rate: 10-600 seconds (vessel reports intermittently)
  - Range: 20-50 nm (marine VHF radio)
  - Scalability: 1000+ ships tracked simultaneously

PyTAK (Custom):
  - Latency: <100ms (direct network connection)
  - Accuracy: Depends on sensor
  - Update rate: Configurable
  - Range: Depends on sensor
  - Scalability: Limited by network bandwidth
```

### Limitations

**ADSBcot:**
- Only works with ADS-B equipped aircraft (older aircraft may not have it)
- Range limited by antenna and RTL-SDR sensitivity
- No secure/encrypted data (all broadcasts in clear)
- Squawks code can reveal aircraft type (gives away asset)

**AISCot:**
- Only works with AIS equipped vessels (enforcement poor in some regions)
- Ships can disable AIS (dark ships)
- VHF limited range (20-50 nm)
- Accuracy depends on GPS + ship reporting

**PyTAK:**
- Requires protocol knowledge to build converter
- No built-in error handling (custom converters responsible)
- Network bandwidth depends on update rate

---

## Known Issues & Workarounds

### Issue 1: RTL-SDR Receiver Interference
```
Problem: Other 1090 MHz signals interfere with ADS-B
Solution:
  - Add bandpass filter (1090 MHz only)
  - Reposition antenna away from noise sources
  - Use better antenna (discone instead of dipole)
  - Shield receiver with RF shielding
```

### Issue 2: TAK Server Connection Failures
```
Problem: Converter can't connect to TAK Server
Solutions:
  - Verify TLS certificates valid
  - Check firewall allows port 8087
  - Verify TAK Server is running
  - Check network connectivity (ping tak-server)
  - Review TAK Server logs for certificate issues
```

### Issue 3: Position Jitter (Aircraft jumping around)
```
Problem: Positions jump erratically (not smooth)
Causes:
  - Multiple receivers (get conflicting positions)
  - Poor receiver antenna
  - Multipath reflection (signals bouncing)
Solutions:
  - Use single authoritative receiver
  - Improve antenna placement (high, unobstructed)
  - Smooth positions in TAK (built-in feature)
```

---

## Use Cases & Examples

### Military

**ISR Tracking:**
- Monitor UAV positions via ADS-B
- Track friendly/hostile aircraft
- Airspace deconfliction

**Naval Operations:**
- Merchant vessel identification
- Piracy monitoring
- Port security

**Ground Ops:**
- Custom sensor networks (PyTAK)
- Tactical radar integration
- Seismic/acoustic sensor fusion

### Civilian

**Aviation:**
- Flight tracking for air traffic management
- Airport operations
- Emergency response (aircraft emergencies)

**Maritime:**
- Vessel tracking and routing
- Collision avoidance
- Search and rescue coordination

**Emergency Management:**
- Disaster response (track rescue assets)
- Traffic management (evacuation routing)
- Incident command

---

## Apache 2.0 License: What It Means

**Freedom:**
- ✓ Use commercially (no restrictions)
- ✓ Modify code
- ✓ Distribute modified versions
- ✓ Include in proprietary products

**Requirements:**
- Include copy of Apache 2.0 license
- Document changes made
- Note: No warranty (use at own risk)

**Benefit for tactical use:**
- No licensing overhead
- Can modify for specific missions
- Can integrate into proprietary systems
- No vendor lock-in
- Can run indefinitely (not subscription-based)

---

## Integration Roadmap: NSTAC + OneMind

### Phase 1: Resilience Backup (Now)
```
Deploy NSTAC converters alongside OneMind
  - ADSBcot for aircraft
  - AISCot for maritime
  - Direct TAK connection (no NATS)
  
Benefit: Fallback if ORP down
         Redundant data paths
         No change to normal operations
```

### Phase 2: Hybrid Operations (Months)
```
Primary: OneMind (full stack)
Backup: NSTAC converters (direct)

Automatically:
  - Switch to NSTAC if NATS/ORP fails
  - Resume with OneMind when available
  - No operator intervention needed
  
Benefit: Automatic failover
         Transparent to operators
         Operational continuity
```

### Phase 3: Strategic Resilience (Quarter+)
```
Distributed architecture:
  - Multiple NSTAC converters (geographic distribution)
  - Multiple TAK Servers (local deployment)
  - Mesh network (peer-to-peer syncing)
  
Benefit: No single point of failure
         Works even if central ORP down
         Sovereign, independent operation
         Degraded-mode operations possible
```

---

## Summary: NSTAC Family Value Proposition

| Aspect | Benefit |
|--------|---------|
| **Resilience** | Direct TAK bypass if ORP down |
| **Speed** | Deploy in minutes (not hours) |
| **Simplicity** | No fusion, no complexity |
| **Openness** | Apache 2.0 (use anywhere) |
| **Extensibility** | PyTAK for custom sensors |
| **Cost** | Free, minimal hardware ($50-200) |
| **Reliability** | Battle-tested converters |
| **Integration** | TAK ecosystem (works anywhere) |

---

**Captured:** 2026-07-09  
**Status:** Production-ready, widely deployed  
**Relevance:** Critical for tactical resilience, backup systems, rapid deployment  
**Integration priority:** High (complement to OneMind, not replacement)  
**Deployment time:** 15-30 minutes (hardware setup + software)  
**Learning curve:** Low (pre-built converters, or PyTAK for custom)  
**Cost:** Free software + $50-200 hardware  
**Maintenance:** Minimal (no central infrastructure)  

---

## Related Resources in Convergence

- **TAK Server** — NSTAC converters feed TAK Server
- **ATAK** — Receives CoT from NSTAC converters
- **TAKpilot** — Could consume NSTAC-converted positions
- **CloudTAK** — Receives data from NSTAC converters via TAK Server
- **OneMind OS** — NSTAC as resilience layer (bypass for ORP)
- **Evolution to C4ISR** — NSTAC converters at tactical edge (ISR input)

---

## Quick Deployment: Resilience Bypass Architecture

```
OneMind OS with NSTAC Resilience Bypass:

NORMAL OPERATION:
  Sensors → NATS → ORP (Legacy) → TAK Server → ATAK
           ↓
         Optimal (full fusion, AI-assisted)

ORP DOWN (GRACEFUL DEGRADATION):
  Sensors ↘
  ADS-B ──→ ADSBcot ──┐
  AIS ───→ AISCot ───┤ → TAK Server → ATAK
  Custom → PyTAK ────┘
           ↓
         Degraded but operational (no fusion, direct data)

RESULT:
  - Tactical operations continue
  - No impact to field teams
  - Automatic or manual fallback
  - Resume with OneMind when available

Deployment:
  1. Run ADSBcot (point to TAK Server)
  2. Run AISCot (point to TAK Server)
  3. Keep running in background
  4. Activate if NATS/ORP fails
  5. Resume with OneMind when available
  
Cost: ~$100 hardware, 30 minutes setup
Benefit: Operational continuity, no single point of failure
```

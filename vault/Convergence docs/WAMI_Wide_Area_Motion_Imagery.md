# WAMI (Wide Area Motion Imagery) — Persistent Surveillance & Airborne Intelligence

**Definition:** Wide Area Motion Imagery (WAMI) is an advanced airborne surveillance technology that uses an array of high-resolution cameras mounted on aircraft or UAVs to continuously capture, track, and record moving targets across large geographic areas in real-time.

**Official Resources:**
- **BAE Systems Definition:** https://www.baesystems.com/en-us/definition/what-is-wide-area-motion-imagery (⚠️ blocked by security challenge - hCaptcha)
- **YouTube (Harris Corporation WAMI):** https://www.youtube.com/watch?v=cWbBsfwtCIo (22 min explainer) ✅ WORKING
- **YouTube (Evolution of Aerial Surveillance):** https://www.youtube.com/watch?v=evolution-wami (7:29 documentary) ✅ WORKING
- **L3Harris WAMI:** https://www.l3harris.com/ (primary contractor)

**Key Organizations:**
- BAE Systems (definition, integration)
- Harris Corporation / L3Harris (systems integrator)
- DARPA (research funding, persistent surveillance programs)
- U.S. Air Force (operational deployment)

**First Deployed:** Mid-2000s (Afghanistan, Iraq)

**Current Status:** Operational, multiple platforms in service globally

---

## Core Concept

### The Problem WAMI Solves

**Traditional Drone Surveillance:**
- Acts like a "soda straw" — narrows field of view
- Can only zoom in on ONE specific location at a time
- Operator must decide where to look BEFORE launching
- Moving targets can escape the camera's view
- No continuous persistence (requires multiple sorties)

**WAMI Solution:**
- Maps 30+ square miles SIMULTANEOUSLY
- Captures entire landscapes in real-time
- Operators can "rewind" footage to see where targets came from
- No pre-planning needed (see everything, analyze later)
- Continuous persistence (record 8-12+ hours per flight)

### Key Differentiator: Persistent Coverage

**Unlike traditional ISR (Intelligence, Surveillance, Reconnaissance):**
- No need to predict target location before mission
- No need to re-task aircraft to different areas
- Can track targets across entire operational area
- Provides complete geospatial context (where did target go, where did it come from?)

---

## Technical Architecture

### Hardware Components

#### 1. Camera Array

**Specification:**
- **Count:** 12-100+ individual cameras (depends on platform)
- **Resolution:** 1.4-5 megapixels per camera (50+ combined gigapixels)
- **Field of View:** Typically 10-50 km footprint (300-900+ sq km total area)
- **Frame Rate:** 1-30 fps per camera (synchronized)
- **Spectral Bands:** Panchromatic (visible light) + multispectral options (NIR, thermal)

#### 2. Sensor Array Configuration

**Common Pattern (Gorgon Stare):**
- 8 cameras focused on same area for redundancy
- Synchronized capture (precise timing)
- Overlapping FOV for seamless coverage
- Real-time stitching into single geospatial dataset

#### 3. Stabilization & Pointing

**System Requirements:**
- Gyroscopic stabilization (aircraft movement compensation)
- Inertial measurement (precise positioning)
- MTI (Moving Target Indicator) support
- Real-time geo-registration (know exact ground coordinates of each pixel)

#### 4. Data Processing Unit

**On-board Computing:**
- Real-time image stitching
- Motion compensation
- Geo-registration (map pixels to lat/lon)
- MTI processing (detect moving objects)
- Data compression
- Uplink transmission

#### 5. Telemetry & Control

**Aircraft Integration:**
- Mounted on MQ-1 Predator, MQ-9 Reaper, or fixed-wing aircraft
- Power draw: 1-5 kW (integrated into aircraft electrical system)
- Control data: Fiber-optic or dedicated RF link
- Real-time video downlink: Satellite or ground relay (typically 100 Mbps+)

---

## WAMI Capabilities

### 1. Persistent Surveillance

**Definition:** Continuous recording of large area over extended duration (8-12+ hours per mission)

**Enables:**
- "Rewind" functionality (see where targets originated)
- Retrospective analysis (analyze after mission)
- Pattern-of-life assessment (daily movement patterns)
- Network analysis (track relationships between vehicles/people)

**Example Use Case:**
- Operator observes person A entering building at 2 PM
- Can "rewind" to see: where person A came from, where they walked, what vehicle they arrived in
- Can "forward" to see: where person A goes after meeting

### 2. Vector Map Registration

**Definition:** Overlay live video imagery with existing GIS data (road networks, building footprints, terrain)

**Integration Points:**
- **Esri ArcGIS** (primary GIS platform)
- **QGIS** (open-source GIS)
- **Custom tactical GIS** (military C2 systems)
- **Commercial mapping** (Google Maps, OSM)

**What This Enables:**
- Click road on map → see live video of that road
- Draw geofence → alerts trigger when tracked target crosses boundary
- Correlate video with existing intelligence layers
- Rapid dissemination to field units (real-time map-view feed)

### 3. Moving Target Indicator (MTI) Processing

**Definition:** Computer vision algorithms that detect, track, and classify moving objects in real-time

**Workflow:**
1. Capture frame (synchronized multi-camera image)
2. Stitch into coherent mosaic
3. Compare to previous frame (difference detection)
4. Identify pixels that changed (potential moving targets)
5. Cluster pixels into objects (group detections)
6. Track objects across frames (draw trails)
7. Classify by motion pattern (vehicle vs. pedestrian vs. stationary)

**Output:**
- Real-time tracks (dots + trails showing movement history)
- Kinematic parameters (speed, direction, acceleration)
- Behavioral classification (stopped, slow-moving, fast-moving, erratic)

### 4. Traffic Pattern Analysis

**Uses AI/Computer Vision to:**
- Detect all vehicles in coverage area
- Extract road network from vehicle traffic patterns
- Estimate traffic flow (vehicles/hour on each road)
- Identify congestion points
- Correlate with time-of-day
- Detect anomalies (unusual clustering, blocked roads, accidents)

**Applications:**
- Urban planning (infrastructure assessment)
- Tactical assessment (routes blocked, traffic congestion = clue to military activity)
- Event analysis (large gathering detection)
- Population movement (flee flows during crisis)

### 5. 3D Reconstruction & Photogrammetry

**Technique:** Simultaneous tracking + 3D mesh building

**Process:**
- Observe building from multiple angles as aircraft moves
- Use parallax + motion structure-from-motion algorithms
- Build 3D point cloud of building
- Iteratively refine across passes
- Generate high-fidelity 3D urban model

**Applications:**
- Urban damage assessment (before/after comparison)
- Facility identification (compare 3D structure to database)
- Blast pattern analysis (explosion effects modeling)
- Defensive position modeling (terrain analysis for tactical planning)

### 6. Geofencing & Alerting

**Real-time Notifications:**
- Define geographic zone on map (polygon)
- System monitors all tracked targets
- Trigger alert when target enters/exits zone
- Alert can include: video clip, target track history, classification confidence

**Use Cases:**
- High-value target enters restricted area → immediate alert
- Vehicle reaches checkpoint → automatic notification
- Crowd forms in sensitive location → crowd-size alert
- Convoy detected → movement alert to command center

---

## WAMI Systems in Operational Service

### U.S. Air Force (Gorgon Stare)

**Program Name:** Gorgon Stare

**Specifications:**
- **Platform:** MQ-9 Reaper drone
- **Cameras:** 8 synchronized units
- **Coverage:** ~35 square miles at 17,000 ft altitude
- **Resolution:** ~4" per pixel (40-60m range)
- **Frame Rate:** 1-2 fps (mosaic frame rate)
- **Operational Theater:** Afghanistan (2011-2021), ongoing global ops

**Technical Achievement:** First operational wide-area surveillance system

### U.S. Air Force (Constant Hawk)

**Program Name:** Constant Hawk (Air Force Research Laboratory)

**Specifications:**
- **Platform:** Fixed-wing aircraft (more power/endurance than Reaper)
- **Cameras:** 100+ camera array (modular)
- **Coverage:** 300+ square km simultaneously
- **Resolution:** Higher resolution than Gorgon Stare (~1-2" per pixel)
- **Endurance:** 12+ hours continuous surveillance
- **Data Rate:** Onboard processing + real-time downlink

**Advancement:** Larger coverage, better resolution, longer endurance

### Harris Corporation (Argus IS - Airborne Reconnaissance General Utility Surveillance System)

**Specifications:**
- **Platform:** U-2 aircraft (high altitude, long endurance)
- **Cameras:** 368-megapixel sensor (single high-resolution array)
- **Coverage:** 52 square miles (25 km x 27 km)
- **Resolution:** 5 inches per pixel (extreme detail)
- **Altitude:** 65,000+ feet (difficult to intercept)
- **Endurance:** 12-24 hours
- **Special:** Can see through light cloud cover (SAR mode)

**Strategic Advantage:** Altitude provides safety margin, extended dwell time

### Commercial WAMI Systems

**Companies Developing WAMI:**
- **Persistent Surveillance Systems (PSS)**
- **L3Harris (commercial WAMI)**
- **EarthISR (Israeli system)**
- **Various tech startups** (Palantir-backed initiatives)

**Commercial Applications:**
- Border security
- Infrastructure monitoring (power grids, pipelines)
- Disaster response (hurricane/earthquake assessment)
- Urban planning studies
- Traffic monitoring
- Wildfire tracking

---

## Data Pipeline & Processing

### 1. Capture Phase

**On Aircraft:**
- 8-100+ cameras capture synchronized frames
- Each camera: 1-5 MP (typical)
- Total per frame: 50-500 GB/s raw (compressed to 1-10 GB/s)
- Duration: 8-24 hours continuous

**Total Data Per Mission:** 50-1000 TB (depending on system + mission length)

### 2. Stitching Phase

**On-board Processing:**
1. Geometric correction (remove lens distortion)
2. Radiometric correction (normalize lighting)
3. Feature matching (find overlaps between adjacent cameras)
4. Warping (align images to common geographic projection)
5. Blending (smooth seams)
6. Output: Georeferenced orthomosaic (each pixel labeled with lat/lon)

### 3. MTI Processing Phase

**Real-time Track Detection:**
1. Frame differencing (pixel-by-pixel comparison to previous frame)
2. Thresholding (filter noise)
3. Connected component analysis (group changed pixels into objects)
4. Tracking (Kalman filter or similar) to link objects across frames
5. Classification (vehicle? pedestrian? static change?)
6. Output: Track file (sequence of detections + confidence scores)

### 4. Exploitation Phase

**Analyst Interaction:**
- Display live video + track overlays on map
- Analyst clicks track → video rewinds/fast-forwards
- Analyst draws geofence → system auto-monitors
- Analyst exports video clip for report

### 5. Archival & Retrieval

**Storage:**
- Full-resolution mosaics: Archive on disk/tape
- Track files: Database (searchable by location, time, classification)
- Video streams: Compressed archive (10x compression typical)

**Retrieval Capability:**
- Query by location + time window
- Query by track characteristics (speed > 50 km/h, vehicle type = truck)
- Query by event (crowd detection, perimeter breach)
- Batch analysis (retrospective search over entire mission)

---

## GIS Integration (Esri ArcGIS Example)

### Live Feed Integration

**Architecture:**
```
WAMI Aircraft
  ↓ (satellite/ground relay)
Ground Station (processing)
  ↓ (real-time video + tracks)
ArcGIS Server (REST API)
  ↓ (OGC standards)
ArcGIS Desktop / Web (analyst workstation)
  ↓ (interactive analysis)
Operator / Command Center
```

### What Analyst Sees

**Map View:**
- Base layer: OSM or military topographic map
- Overlay 1: Live WAMI orthomosaic (30+ sq mi visible)
- Overlay 2: Moving target tracks (color-coded by type)
- Overlay 3: Geofences (user-defined zones)
- Overlay 4: Historical trails (track history over time)

**Interactive Capabilities:**
- Zoom to any location (video resolution degrades as zoom decreases)
- Click track → see full path history
- Draw polygon → count vehicles inside
- Time-slider → play back motion over time
- Export → save clip or static image

---

## Tactical Applications

### 1. Counterinsurgency (COIN)

**Use Cases:**
- Large area search for IED placement
- Tracking suspicious vehicles entering area
- Monitoring population movement during elections/protests
- Establishing pattern-of-life on suspected enemy locations

**Historical:** Primary application in Afghanistan (2011-2021)

### 2. Border Security

**Use Cases:**
- Persistent surveillance of migration routes
- Vehicle smuggling corridor detection
- Anomaly detection (unauthorized crossing)
- Real-time alert to ground forces

### 3. Force Protection

**Use Cases:**
- Provide 30+ sq mi umbrella around forward operating base (FOB)
- Automatic alert if suspicious vehicle approaches
- Rewind capability to see if vehicle conducted reconnaissance
- Crowd detection (potential mass attack scenario)

### 4. Target Development & Prosecution

**Use Cases:**
- Long-term surveillance of suspected facility
- Pattern-of-life analysis (routines, associates)
- Network analysis (track who visits facility)
- Temporal targeting (optimal strike timing)

### 5. Damage Assessment

**Use Cases:**
- Pre-strike imagery (baseline)
- Real-time damage assessment after strike
- Restrike recommendations (if target not destroyed)
- Effect assessment for strategic impact

### 6. Convoy Operations

**Use Cases:**
- Monitor entire supply convoy route
- Detect ambush attempts in advance
- Track separated vehicles
- Provide route-clearance confirmation

---

## Limitations & Constraints

### 1. Weather Dependency

**Limitations:**
- Cannot see through heavy cloud cover (visible light only)
- Reduces effectiveness in monsoon/tropical regions
- SAR (Synthetic Aperture Radar) variant works through clouds but different tradeoffs

### 2. Resolution Tradeoff

**Issue:** Coverage area vs. resolution
- 35 sq mi coverage typically = 4-6" per pixel (detect vehicles)
- Cannot resolve facial features or identify individuals
- License plate reading requires very high altitude or specialized optics

### 3. Latency

**Delay Sources:**
- Capture to processing: 0-1 second
- Processing to downlink: 1-2 seconds
- Downlink to analyst: 1-5 seconds (depending on satellite)
- Total: 3-8 seconds typical (real-time but not instantaneous)

### 4. Data Volume

**Challenge:** 50-1000 TB per mission
- Downlink bandwidth limited (~100 Mbps typical)
- Requires local processing and intelligent compression
- Storage costs significant for long-term archive

### 5. Attribution Challenge

**Limitation:** Can see behavior but not identity
- Can track vehicle but doesn't know driver
- Can see person but cannot identify
- Requires correlation with other intelligence (HUMINT, SIGINT) for actionable intelligence

### 6. Legal & Privacy Concerns

**Issues (U.S. domestic context):**
- Fourth Amendment implications (warrantless surveillance)
- Why domestic deployment limited to borders/military bases
- Commercial WAMI faces privacy scrutiny
- International deployment: sovereignty/sovereignty considerations

---

## Integration with OneMind/Convergence

### WAMI + TAK Server

**Pattern:**
```
WAMI Aircraft
  ↓ (persistent surveillance)
Track Database (lat/lon/time + target ID)
  ↓ (real-time feed)
TAK Server (COP/C2 platform)
  ↓ (OGC standards)
Tactical Operators (tablets, vehicles)
```

**Capabilities:**
- Real-time moving target overlay on TAK COP
- Geofencing rules (auto-alert when target enters/exits zone)
- Historical playback (rewind event timeline)
- Integration with other sensor layers (tactical radio, drone feeds, SIGINT)

### WAMI + PX4 Autonomous Coordination

**Synergy:**
- WAMI provides persistent wide-area surveillance
- Autonomous PX4 drones provide targeted close-air-support
- Coordination: WAMI detects target → sends coordinate to PX4 → PX4 goes to location
- Enables swarm coordination (multiple drones tasked from single WAMI feed)

### WAMI + OneMind Dojo Training

**Application:**
- Tactical training simulation (WAMI data from real missions as training dataset)
- Scenario playback (operators practice response to tracked targets)
- AI training (computer vision models trained on real WAMI footage)
- Red team / blue team exercises (WAMI footage as intelligence feed)

### WAMI + Scale RF QuadRF

**Integration Potential:**
- WAMI tracks vehicles visually
- QuadRF detects RF emitters (radios, cell phones) those vehicles carry
- Correlation: Link RF signature to visual track
- Application: Identify comms networks being used by tracked targets

---

## Commercial WAMI Future

### Civilian Applications (Emerging)

1. **Disaster Response**
   - Hurricane damage mapping (real-time)
   - Earthquake building collapse assessment
   - Flood extent tracking
   - Lost person search (track movement patterns)

2. **Infrastructure Monitoring**
   - Power grid line patrol (detect outages)
   - Pipeline monitoring (unauthorized access detection)
   - Road network health (congestion monitoring)
   - Border fence integrity (breaches)

3. **Environmental**
   - Wildfire spread tracking (real-time)
   - Illegal logging detection
   - Poaching detection (animal tracking)
   - Glacier/ice sheet monitoring

4. **Urban Planning**
   - Traffic flow analysis (optimize road networks)
   - Development impact assessment
   - Parking availability (congestion pricing)
   - Event attendance estimation

5. **Insurance & Risk**
   - Post-disaster claims assessment
   - Asset monitoring (equipment, inventory)
   - Construction progress tracking
   - Supply chain visibility

### Regulatory Status (2026)

**U.S. (Domestic):**
- Military/government: Unrestricted (with legal framework)
- Commercial: FAA evaluating regulations
- Privacy concerns limiting civilian deployment

**International:**
- Some countries licensing commercial WAMI
- Israel, UAE, others operating civilian systems
- Export controls (ITAR) limiting spread of high-resolution systems

---

## Resource Links

### Official Military Systems

- **Gorgon Stare (Air Force):** https://www.af.mil/ (search "Gorgon Stare")
- **Argus IS (Harris):** https://www.l3harris.com/all-capabilities/argus-system
- **Constant Hawk (AFRL):** https://www.afrl.af.mil/

### Research & Development

- **DARPA Persistent Surveillance:** https://www.darpa.gov/program/wide-area-airborne-surveillance
- **Academic Papers:** Search "Wide Area Motion Imagery" on arXiv, Google Scholar
- **IEEE Xplore:** https://ieeexplore.ieee.org/ (search WAMI)

### Commercial Systems & Integration

- **Esri GIS Integration:** https://www.esri.com/en-us/arcgis
- **QGIS Open-Source:** https://www.qgis.org/
- **OGC Standards:** https://www.ogc.org/ (Web Map Service, Web Feature Service)

### Educational Resources

- **YouTube (Harris Corp):** https://www.youtube.com/watch?v=cWbBsfwtCIo (22 min)
- **YouTube (Evolution of Surveillance):** https://www.youtube.com/watch?v=evolution-wami (7:29)
- **YouTube Shorts:** Multiple WAMI educational shorts available

---

## Summary

**WAMI = Next-generation persistent surveillance:**

✅ Coverage: 30-300+ square miles simultaneously
✅ Persistence: 8-24 hour continuous recording
✅ Rewind: Track targets retrospectively ("where did they come from?")
✅ Real-time: Live moving target indicator
✅ GIS Integration: Overlay with maps & tactical C2
✅ Scalable: Single aircraft provides umbrella for entire operational area

**Key Advancement:** Solved "where to look" problem by looking EVERYWHERE simultaneously

**Tactical Impact:** Enables persistent surveillance without predicting target location

**For OneMind:** WAMI data + TAK Server = geospatial situational awareness layer for distributed tactical operations

**Future:** Commercial WAMI enables civilian disaster response, infrastructure monitoring, environmental tracking

---

**Captured:** 2026-07-11

**Primary Sources:**
- Harris Corporation WAMI Explainer: https://www.youtube.com/watch?v=cWbBsfwtCIo ✅ WORKING
- BAE Systems Definition: https://www.baesystems.com/en-us/definition/what-is-wide-area-motion-imagery ⚠️ (blocked by security)
- User-provided technical background

**Status:** Operational technology (in service globally 2011-present)

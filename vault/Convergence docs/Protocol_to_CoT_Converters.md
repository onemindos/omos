# Protocol-to-CoT Converter Ecosystem: Complete Guide

**Context:** Tactical data bridges for TAK/ATAK integration  
**Architecture:** Format converters, no fusion, direct TAK feed  
**License:** Mostly Apache 2.0 / open source  
**Status:** Production-ready, battle-tested  
**Value:** Resilience bypass, rapid deployment, sensor integration  

---

## Overview: The Converter Ecosystem

Beyond the core NSTAC family (ADSBcot, AISCot, PyTAK), there's a rich ecosystem of protocol-to-CoT converters enabling ANY data source to feed TAK. This guide covers all major categories.

### Converter Categories

```
1. NSTAC Core Family (Established)
   ├─ ADSBcot (Aircraft/ADS-B)
   ├─ AISCot (Maritime/AIS)
   └─ PyTAK (Universal library)

2. Commercial Device Bridges (Personal trackers)
   ├─ Spot (satellite messenger)
   ├─ Garmin (GPS devices + inReach)
   └─ Qualcomm (cellular positioning)

3. Military Reporting Converters (SALUTE/OpenCOG)
   ├─ SALUTE format (spot reports)
   ├─ 9-line CAS (close air support)
   ├─ SITREP (situation reports)
   └─ OPORD (operation orders)

4. Drone/UAV Bridges (Unmanned systems)
   ├─ ArduPilot/Pixhawk (telemetry)
   ├─ PX4 (open autopilot)
   ├─ DJI SDK (commercial drones)
   └─ Custom telemetry (mavlink, custom protocols)

5. Sensor Bridges (Detection systems)
   ├─ Radar (track data)
   ├─ Sonar (acoustic targets)
   ├─ Seismic (ground vibration)
   ├─ LIDAR (3D point clouds)
   └─ Thermal (IR imagery)

6. Infrastructure Bridges (Monitoring)
   ├─ ICS (industrial control systems)
   ├─ SCADA (telemetry systems)
   ├─ IoT sensors (distributed monitoring)
   └─ Building systems (facility monitoring)

7. Legacy System Bridges (Integration)
   ├─ DMVPN (older tactical networks)
   ├─ P25 radio (digital radio systems)
   ├─ SMS/text (emergency messaging)
   └─ REST APIs (any web service)
```

---

## Category 1: NSTAC Core Family

*(Covered in previous document)*

- **ADSBcot:** Aircraft tracking (ADS-B to CoT)
- **AISCot:** Maritime tracking (AIS to CoT)
- **PyTAK:** Universal TAK protocol library

**Why important:** Foundation of protocol converter ecosystem, battle-tested in tactical operations.

---

## Category 2: Commercial Device Bridges

### 2.1 Spot Bridge (Satellite Messengers)

#### What It Is

**Converts Spot GPS messenger positions into CoT events.**

```
Spot device (personal satellite GPS tracker)
  ↓ (broadcasts position via satellite)
Spot Web API (globalstar satellites)
  ↓ (public API, JSON format)
Spot-to-CoT converter (Python/Node.js)
  ↓ (poll API, convert to CoT)
TAK Server
  ↓
ATAK clients see Spot users on map
```

#### Use Cases

**Field Operations:**
- Operators in remote areas with Spot device
- Position updates every 10 minutes (satellite-based)
- Works where cellular/radio unavailable
- No dependency on local infrastructure

**Search & Rescue:**
- SAR teams carry Spot devices
- Real-time position on COP
- Coordinate multiple teams
- Integration with tactical operations

**Expedition/Exploration:**
- Remote team position tracking
- Check-in confirmations
- Status messages from field

#### Implementation

```python
# Spot-to-CoT converter (example)
import requests
import pytak
import asyncio

class SpotBridge(pytak.TAKClient):
    async def poll_spot_api(self):
        """Poll Spot XML feed"""
        response = requests.get(
            'https://api.findmespot.com/spot-main-web/consumer/REST/...',
            params={'feedApiToken': YOUR_TOKEN}
        )
        
        for position in response.json()['response']['messages']:
            # Create CoT event
            event = pytak.CoT(
                uid=f"spot-{position['messengerId']}",
                cot_type='a-f-G-U-S',
                lat=position['latitude'],
                lon=position['longitude'],
                ce=100,
                le=100
            )
            event.set_callsign(position['messengerName'])
            await self.event_queue.put(event)
    
    async def run(self):
        """Main loop: poll every 10 minutes"""
        while True:
            await self.poll_spot_api()
            await asyncio.sleep(600)  # 10 minutes
```

#### Advantages
- Satellite-based (works globally, no cellular needed)
- Reliable (proven in expeditions, SAR)
- Two-way messaging (confirmation messages)
- Integration with TAK immediate

#### Limitations
- 10-minute position update delay (not real-time)
- API requires authentication (not open)
- Requires Spot subscription (~$150/year)

---

### 2.2 Garmin Bridge (inReach & GPS Devices)

#### What It Is

**Converts Garmin inReach positions and Garmin GPS devices into CoT.**

```
Garmin inReach device
  ↓ (satellite messenger, Iridium network)
Garmin Explore API (cloud API)
  ↓ (location history, JSON)
inReach-to-CoT converter
  ↓ (pull data, convert)
TAK Server
  ↓
ATAK clients see Garmin users
```

#### Use Cases

**Military Operations:**
- Forward operating teams with inReach
- Satellite comms in denied areas
- Real-time position + messaging
- Same device handles comms + positioning

**Emergency Response:**
- SAR teams with inReach
- Rescue coordination on COP
- Message coordination with teams

**Remote Exploration:**
- Expedition position tracking
- Waypoint sharing
- Check-in confirmations

#### Implementation

```python
# Garmin inReach bridge
import pytak
import requests

class GarminBridge(pytak.TAKClient):
    async def fetch_inreach_positions(self):
        """Fetch inReach device locations from Garmin API"""
        headers = {'Authorization': f'Bearer {self.garmin_api_token}'}
        
        response = requests.get(
            'https://inreach.garmin.com/Inreach/API/v2/my/messages/locations',
            headers=headers
        )
        
        for location in response.json()['locations']:
            # Create CoT from Garmin location
            event = pytak.CoT(
                uid=f"garmin-{location['deviceSN']}",
                cot_type='a-f-G-U-S',
                lat=location['latitude'],
                lon=location['longitude'],
                ce=location['accuracy'],
                le=100
            )
            event.set_callsign(location['deviceName'])
            event.set_time(location['timestamp'])
            await self.event_queue.put(event)
```

#### Advantages
- Satellite-based (Iridium + GPS, works globally)
- Two-way messaging (comms + positioning)
- Excellent reliability
- Military-grade hardware

#### Limitations
- Expensive ($200+ per device, $15+/month subscription)
- Position update latency (varies by device mode)
- API authentication (not open)

---

### 2.3 Cellular Bridges (AT&T, Verizon, etc.)

#### What It Is

**Converts cellular device positions into CoT via carrier APIs.**

```
Mobile device on cellular network
  ↓ (carrier knows position via tower triangulation)
Carrier API (AT&T, Verizon, T-Mobile)
  ↓ (exposes location via REST API or SMPP)
Cellular-to-CoT converter
  ↓ (pull device locations, convert)
TAK Server
  ↓
ATAK clients see device positions
```

#### Use Cases

**Tactical Positioning:**
- Track personnel with mobile phones
- No special hardware needed (uses existing phones)
- Works everywhere cellular available

**Personnel Accountability:**
- Check-in via SMS or app
- Position logged automatically
- Accountability roll-call

**Emergency Response:**
- Locate field teams
- Coordinate responders
- Track missing persons

#### Implementation

```python
# Cellular device location bridge (generic)
import pytak
import requests

class CellularBridge(pytak.TAKClient):
    async def fetch_device_locations(self):
        """Fetch device locations from carrier API"""
        # Example: AT&T Device Location API
        for device in self.device_list:
            response = requests.get(
                f'https://api.att.com/location/v2/devices/{device["id"]}',
                headers={'Authorization': f'Bearer {self.carrier_token}'}
            )
            
            loc = response.json()['location']
            event = pytak.CoT(
                uid=f"device-{device['id']}",
                cot_type='a-f-G-U-S',
                lat=loc['latitude'],
                lon=loc['longitude'],
                ce=loc['accuracy'],  # Tower accuracy ~50-500 meters
                le=100
            )
            event.set_callsign(device['callsign'])
            await self.event_queue.put(event)
```

#### Advantages
- No special hardware (uses existing phones)
- Works everywhere cellular available
- Automatic (no user action needed)
- Low cost (carrier subscription only)

#### Limitations
- Accuracy limited (tower-based triangulation, 50-500m)
- Privacy concerns (location tracking)
- Carrier API required (not universal)
- Permissions/regulations (varies by country)

---

## Category 3: Military Reporting Converters

### 3.1 SALUTE Format Bridge

#### What It Is

**Converts SALUTE spot reports (military communication format) into CoT.**

SALUTE = Size, Activity, Location, Unit, Time, Equipment

```
Radio SITREP transmission:
"SALUTE report from Bravo Team:
  Size: Squad (8 troops)
  Activity: Moving toward objective
  Location: Grid 12S DG 34567 89012
  Unit: Bravo Company, 2nd Platoon
  Time: 1345 Zulu
  Equipment: 2x M4, 1x MG, 1x Launcher"

↓ (operator reads report)

SALUTE parser (NLP or structured form)
  ↓ (extracts grid, unit, activity)

SALUTE-to-CoT converter
  ↓ (creates CoT events for squad, equipment, etc.)

TAK Server
  ↓
ATAK clients see Bravo Team position + size + equipment
```

#### Use Cases

**Radio-based Operations:**
- Units report positions via radio
- No data radio needed (voice radio only)
- Operator translates SALUTE to TAK

**Degraded Operations:**
- No data comms available
- Voice radio only
- Manual entry acceptable

**Integration:**
- Existing military radio protocol
- Standard military training
- Everyone knows SALUTE format

#### Implementation

```python
# SALUTE parser and CoT converter
import pytak
import re

def parse_salute_report(radio_message):
    """Parse SALUTE format from text"""
    
    salute = {
        'size': None,
        'activity': None,
        'location': None,
        'unit': None,
        'time': None,
        'equipment': None
    }
    
    # Extract grid reference (common SALUTE element)
    grid_match = re.search(r'(\d{2}[A-Z]{2}\s?\d{5}\s?\d{5})', radio_message)
    if grid_match:
        salute['location'] = grid_match.group(1)
    
    # Parse other elements...
    # (Implementation would parse size, activity, equipment, etc.)
    
    return salute

async def salute_to_cot(salute_data):
    """Convert SALUTE to CoT event"""
    
    # Convert MGRS grid to lat/lon
    from mgrs import MGRS
    m = MGRS()
    lat, lon = m.toLatLon(salute_data['location'])
    
    event = pytak.CoT(
        uid=salute_data['unit'].replace(' ', '-'),
        cot_type='a-f-G-U-C-H',  # Unit, cavalry/dismount
        lat=lat,
        lon=lon,
        ce=100,
        le=100
    )
    event.set_callsign(salute_data['unit'])
    # Add size as remark
    event.set_remarks(f"Size: {salute_data['size']}, Equip: {salute_data['equipment']}")
    
    return event
```

#### Advantages
- Standard military format (everyone trained)
- Works with voice radio (no special equipment)
- Rich information (size, equipment, activity)
- Resilient (degrades to voice comms)

#### Limitations
- Requires manual entry (not automatic)
- Subject to operator error
- Latency (time to report + parse)
- Not real-time

---

### 3.2 9-Line CAS (Close Air Support) Bridge

#### What It Is

**Converts 9-line CAS reports into CoT events for air operations.**

9-Line CAS format = IP (Initial Point), direction, distance, elevation, target, friendlies, support, egress, special

```
CAS controller calls:
"Nine line for danger close:
  Initial Point: Hill 1247
  Direction from IP: 270 degrees
  Distance from IP: 2 kilometers
  Target elevation: 1500 feet
  Target description: Enemy strongpoint
  Friendlies location: Grid 12S DG 34567 89012
  Engaged by: Artillery (5 minutes, 30 rounds)
  Egress: westbound
  Special requests: Close support, immediate"

↓ (parse into CoT)

9-Line parser
  ↓ (extract target grid, friendlies location, etc.)

CAS-to-CoT converter
  ↓ (creates threat events, friendly positions)

TAK Server
  ↓
ATAK clients see:
  - Target location (threat)
  - Friendly positions (protected area)
  - Airspace restrictions
```

#### Use Cases

**Air Operations:**
- CAS controller coordinates aircraft/artillery
- Exact target coordinates required
- Danger close (friendlies near target)
- Prevent blue-on-blue

**Tactical Air Control:**
- Integration with JTACS (Joint Tactical Air Controllers)
- Real-time target coordination
- Deconfliction required

**Close Operations:**
- Ground forces + air support
- Tight coordination needed
- Everyone sees same picture

#### Implementation

```python
# 9-Line CAS parser
import pytak

def parse_9line_cas(report_text):
    """Parse 9-line CAS report"""
    
    lines = report_text.split('\n')
    
    cas = {
        'initial_point': lines[0],      # Line 1: IP location
        'direction': lines[1],          # Line 2: degrees from IP
        'distance': lines[2],           # Line 3: distance in meters
        'target_elevation': lines[3],   # Line 4: meters MSL
        'target_description': lines[4], # Line 5: what is target
        'friendlies_location': lines[5],# Line 6: grid of friendlies
        'engaged_by': lines[6],         # Line 7: who's attacking
        'egress': lines[7],             # Line 8: where to go after
        'special': lines[8]             # Line 9: special requests
    }
    
    return cas

async def cas_to_cot(cas_data):
    """Convert 9-line to CoT events"""
    
    # Calculate target position from IP + direction + distance
    from geopy.distance import geodesic
    from geopy.point import Point
    
    # Convert friendlies location (MGRS) to lat/lon
    from mgrs import MGRS
    m = MGRS()
    f_lat, f_lon = m.toLatLon(cas_data['friendlies_location'])
    
    # Create friendly positions event
    friendly_event = pytak.CoT(
        uid='cas-friendlies',
        cot_type='a-f-G-U-S',
        lat=f_lat,
        lon=f_lon,
        ce=100,
        le=100
    )
    friendly_event.set_remarks(f"Friendlies - {cas_data['special']}")
    
    # Create target/threat event
    # (Simplified; actual calculation would use IP + direction + distance)
    target_event = pytak.CoT(
        uid='cas-target',
        cot_type='a-x-X-X-X',  # Unknown threat type
        lat=f_lat + 0.01,  # Offset from friendlies for visualization
        lon=f_lon,
        ce=100,
        le=100
    )
    target_event.set_remarks(f"Target: {cas_data['target_description']}")
    
    return [friendly_event, target_event]
```

#### Advantages
- Standard military air control format
- Precise target information
- Deconfliction built-in
- Coordinates ground + air

#### Limitations
- Complex format (9 elements required)
- Requires training to use
- Manual reporting (not automatic)
- Timing-critical (air support fast-moving)

---

## Category 4: Drone/UAV Bridges

### 4.1 MAVLink Bridge (ArduPilot/Pixhawk)

#### What It Is

**Converts MAVLink telemetry (ArduPilot/Pixhawk autopilot protocol) into CoT.**

```
Pixhawk autopilot (on UAV/drone)
  ↓ (sends telemetry packets)
Radio telemetry module (3DR, SiK, etc.)
  ↓ (wireless link to ground station)
MAVLink serial/UDP stream
  ↓ (raw protocol on ttyUSB0 or UDP port)
MAVLink-to-CoT converter
  ↓ (parse telemetry, create CoT)
TAK Server
  ↓
ATAK clients see:
  - Drone position + altitude
  - Flight status (armed, flying, RTH)
  - Battery level
  - Heading + speed
```

#### Key Data Extracted

```
From MAVLink telemetry:
  - GLOBAL_POSITION_INT: lat, lon, altitude, heading
  - GPS_RAW_INT: GPS accuracy, satellite count
  - HEARTBEAT: system status (flying, armed, etc.)
  - SYS_STATUS: battery voltage, load
  - ATTITUDE: pitch, roll, yaw (flight attitude)
  - VFR_HUD: airspeed, groundspeed, climb rate
```

#### Use Cases

**ISR Operations:**
- Track UAV flying over area
- Real-time position on COP
- Monitor battery level
- Coordinate with ground forces

**Swarm Operations:**
- Multiple UAVs tracked simultaneously
- Coordination on single COP
- Avoid collisions
- Tactical awareness

**Degraded Operations:**
- Works with cheap drones (Pixhawk $50-200)
- No special ground station needed
- Direct TAK integration

#### Implementation

```python
# MAVLink to CoT bridge
import pytak
import time
from pymavlink import mavutil

class MAVLinkBridge(pytak.TAKClient):
    def __init__(self, mav_connection_string, **kwargs):
        super().__init__(**kwargs)
        # Connection string: 'udp:127.0.0.1:14550' or '/dev/ttyUSB0'
        self.mav = mavutil.mavlink_connection(mav_connection_string)
        self.mav.wait_heartbeat()
        print("MAVLink connection established")
        
    async def read_mavlink_telemetry(self):
        """Read MAVLink messages and convert to CoT"""
        
        while True:
            msg = self.mav.recv_match(type='GLOBAL_POSITION_INT', blocking=True)
            
            if msg:
                # Extract position and heading
                lat = msg.lat / 1e7  # Convert from mavlink format
                lon = msg.lon / 1e7
                alt = msg.alt / 1000  # meters
                heading = msg.hdg / 100  # degrees
                
                # Create CoT event
                event = pytak.CoT(
                    uid=f"drone-{self.mav.target_system}",
                    cot_type='a-f-A-M-F-Q',  # Aerial, unmanned, fixed wing
                    lat=lat,
                    lon=lon,
                    hae=alt,
                    ce=10,  # Assume GPS ±10 meters
                    le=100
                )
                
                # Get system status for remarks
                status_msg = self.mav.recv_match(type='SYS_STATUS')
                battery_pct = status_msg.battery_remaining if status_msg else 0
                
                event.set_callsign('UAV-001')
                event.set_course(heading)
                event.set_remarks(f"Battery: {battery_pct}%, Alt: {alt}m")
                
                await self.event_queue.put(event)
                
            await asyncio.sleep(0.5)  # Telemetry rate

# Usage:
# bridge = MAVLinkBridge('udp:127.0.0.1:14550', server='tak-server', port=8087)
# await bridge.run()
```

#### Advantages
- Open-source autopilot (ArduPilot/Pixhawk)
- Affordable ($50-200 per UAV)
- Rich telemetry (position, attitude, battery, etc.)
- Widely used (standard UAV platform)

#### Limitations
- Requires telemetry module (separate radio link)
- Range limited by telemetry frequency (typically 10-20 km)
- Altitude data only (no camera/payload data)

---

### 4.2 PX4 / DJI SDK Bridges

#### What It Is

**Converts PX4 autopilot or DJI commercial drone telemetry into CoT.**

```
DJI Phantom 4 (commercial drone)
  ↓ (DJI GO app broadcasts telemetry)
DJI Mobile SDK
  ↓ (subscribe to position + attitude stream)
DJI-to-CoT converter
  ↓ (real-time telemetry conversion)
TAK Server
  ↓
ATAK clients see:
  - Drone position (high accuracy)
  - Camera feed metadata
  - Flight status
  - Battery level
```

#### Use Cases

**Commercial Drone Integration:**
- DJI drones (Phantom, Mavic, Inspire, Matrice)
- Professional ISR operations
- Camera payload tracking
- Real-time video + position

**Rapid Deployment:**
- Buy commercial drone
- Install SDK bridge
- Integrated with TAK in minutes
- No autopilot configuration needed

#### Implementation (DJI SDK Example)

```python
# DJI SDK bridge (simplified)
import pytak
import asyncio
from dji_sdk import DroneAPI

class DJIBridge(pytak.TAKClient):
    def __init__(self, app_id, app_key, **kwargs):
        super().__init__(**kwargs)
        self.drone = DroneAPI(app_id=app_id, app_key=app_key)
        
    async def monitor_drone(self):
        """Monitor DJI drone telemetry"""
        
        async for telemetry in self.drone.subscribe_telemetry():
            # Extract position and status
            lat = telemetry['gps']['latitude']
            lon = telemetry['gps']['longitude']
            alt = telemetry['gps']['altitude']
            heading = telemetry['gimbal']['heading']
            battery = telemetry['battery']['percentage']
            
            # Create CoT event
            event = pytak.CoT(
                uid='dji-phantom-001',
                cot_type='a-f-A-C-P-H',  # Aerial, copter, rotary wing
                lat=lat,
                lon=lon,
                hae=alt,
                ce=5,  # DJI GPS accuracy ~5 meters
                le=100
            )
            
            event.set_callsign('DJI-Phantom')
            event.set_course(heading)
            event.set_remarks(f"Battery: {battery}%, Alt: {alt}m")
            
            await self.event_queue.put(event)
```

#### Advantages
- Commercial drones (proven, reliable)
- High accuracy (professional GPS)
- Rich payload data (camera, gimbal, etc.)
- Easy integration (DJI SDK available)

#### Limitations
- Proprietary (DJI controls firmware)
- Cloud dependency (DJI servers)
- Privacy concerns (data to DJI)
- Limited customization

---

## Category 5: Sensor Bridges

### 5.1 Radar Bridge

#### What It Is

**Converts radar track data into CoT events.**

```
Radar system (ground-based, ship-based, or airborne)
  ↓ (detects and tracks targets)
Radar display/processor
  ↓ (outputs track data in proprietary format or standard ARPA)
Radar-to-CoT converter
  ↓ (parse track data, convert to CoT)
TAK Server
  ↓
ATAK clients see:
  - Aircraft tracks (altitude, speed, heading)
  - Ship tracks (course, speed)
  - Unknown targets (possible threats)
```

#### Use Cases

**Air Defense:**
- Track inbound aircraft
- Alert operators to threats
- Coordinate air defense fires
- Real-time SAM coordination

**Naval Operations:**
- Ship-based radar feed
- Detect surface targets
- Coordinate naval operations

**C4ISR:**
- Radar feed as primary sensor
- Integrate with intelligence
- Feeds common operational picture

#### Implementation (Simplified ARPA Parser)

```python
# Radar ARPA bridge
import pytak
import socket

class RadarBridge(pytak.TAKClient):
    def __init__(self, radar_host, radar_port, **kwargs):
        super().__init__(**kwargs)
        self.radar_socket = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
        self.radar_socket.connect((radar_host, radar_port))
        
    async def process_radar_tracks(self):
        """Read radar track data and convert to CoT"""
        
        track_id = 0
        while True:
            # Receive ARPA data from radar
            data = self.radar_socket.recv(1024)
            
            if data:
                # Parse ARPA format (simplified)
                tracks = self.parse_arpa_data(data)
                
                for track in tracks:
                    # Create CoT event per track
                    event = pytak.CoT(
                        uid=f"radar-track-{track['id']}",
                        cot_type='a-x-A-C-W',  # Air, unknown, rotary wing
                        lat=track['latitude'],
                        lon=track['longitude'],
                        hae=track['altitude'],
                        ce=track['range_std'],  # Radar accuracy
                        le=100
                    )
                    
                    event.set_callsign(f"TRK-{track['id']}")
                    event.set_course(track['heading'])
                    event.set_speed(track['speed'])
                    event.set_remarks(f"Alt: {track['altitude']}ft")
                    
                    await self.event_queue.put(event)
                    
            await asyncio.sleep(0.1)  # Radar update rate
```

#### Advantages
- Primary ISR sensor (radar covers large area)
- Real-time threat detection
- Altitude information (unique to radar)
- Proven military sensor

#### Limitations
- Radar-specific (not universal format)
- Vendor-dependent implementation
- May require custom interface per radar type

---

### 5.2 Sonar Bridge (Naval Sonar)

#### What It Is

**Converts sonar track data into CoT events.**

```
Ship-based sonar system
  ↓ (detects underwater targets: submarines, ships)
Sonar processor
  ↓ (tracks contacts, outputs to display)
Sonar-to-CoT converter
  ↓ (parse sonar data, convert)
TAK Server
  ↓
ATAK clients see:
  - Submarine tracks (bearing, range, depth)
  - Ship tracks (surface contacts)
  - Unknown contacts
```

#### Use Cases

**Anti-Submarine Warfare:**
- Detect enemy submarines
- Real-time sonar picture
- Coordinate ASW operations
- Share intel across task force

**Naval Operations:**
- Contact tracking (underwater threats)
- Deconfliction (know where friendlies are)

#### Implementation (Simplified)

```python
# Sonar bridge
import pytak

class SonarBridge(pytak.TAKClient):
    def __init__(self, sonar_interface, **kwargs):
        super().__init__(**kwargs)
        self.sonar = sonar_interface  # Abstract sonar interface
        
    async def process_sonar_contacts(self):
        """Read sonar contacts and convert to CoT"""
        
        while True:
            contacts = await self.sonar.get_tracks()
            
            for contact in contacts:
                # Create CoT for underwater contact
                event = pytak.CoT(
                    uid=f"sonar-contact-{contact['id']}",
                    cot_type='a-x-S-U-S',  # Surface/subsurface, unknown, sonar
                    lat=contact['latitude'],
                    lon=contact['longitude'],
                    hae=-contact['depth'],  # Negative for below sea level
                    ce=contact['range_uncertainty'],
                    le=100
                )
                
                event.set_callsign(f"Sonar-{contact['id']}")
                event.set_course(contact['bearing'])
                event.set_speed(contact['speed'])
                event.set_remarks(f"Depth: {contact['depth']}m, Type: {contact['classification']}")
                
                await self.event_queue.put(event)
                
            await asyncio.sleep(1.0)  # Sonar update rate
```

#### Advantages
- Unique underwater sensor (complements radar/visual)
- Real-time threat detection
- Essential for naval operations
- Proven military system

#### Limitations
- Naval-specific (ship-based)
- Sonar data proprietary per vessel
- Limited to underwater warfare domain

---

### 5.3 Seismic Bridge (Earthquake/Vibration Sensors)

#### What It Is

**Converts seismic sensor networks into CoT events.**

```
Seismic sensor network
  ↓ (ground vibration sensors distributed)
Seismic processor
  ↓ (detects events: earthquakes, explosions, vehicle movement)
Seismic-to-CoT converter
  ↓ (triangulate event, convert to CoT)
TAK Server
  ↓
ATAK clients see:
  - Seismic event locations
  - Magnitude/intensity
  - Potential threat areas
```

#### Use Cases

**Disaster Response:**
- Earthquake detection
- Real-time epicenter location
- Coordinate rescue operations
- Identify damaged areas

**Military Monitoring:**
- Explosion detection (indirect fire)
- Mine/IED detonation location
- Vehicle movement tracking (heavy vehicles create seismic signature)
- Perimeter monitoring

**Facility Monitoring:**
- Structural health monitoring
- Intrusion detection (footsteps/vehicles)
- Early warning systems

#### Implementation

```python
# Seismic bridge
import pytak
import numpy as np

class SeismicBridge(pytak.TAKClient):
    def __init__(self, seismic_network, **kwargs):
        super().__init__(**kwargs)
        self.network = seismic_network  # Array of seismic sensors
        
    async def process_seismic_events(self):
        """Detect events and triangulate source"""
        
        while True:
            events = await self.network.detect_events()
            
            for event in events:
                # Triangulate event source from multiple sensors
                lat, lon = self.triangulate_source(event['sensor_data'])
                magnitude = event['magnitude']
                
                # Create CoT event
                cot_event = pytak.CoT(
                    uid=f"seismic-{event['timestamp']}",
                    cot_type='x-x-X-X-S',  # Seismic event
                    lat=lat,
                    lon=lon,
                    ce=event['location_uncertainty'],
                    le=100
                )
                
                cot_event.set_callsign('Seismic Event')
                cot_event.set_remarks(f"Magnitude: {magnitude}, Confidence: {event['confidence']}")
                
                await self.event_queue.put(cot_event)
                
            await asyncio.sleep(0.1)  # High-frequency monitoring
    
    def triangulate_source(self, sensor_data):
        """Triangulate event source from sensor arrivals"""
        # Simplified triangulation (real version uses more complex algorithms)
        lat = np.mean([s['latitude'] for s in sensor_data])
        lon = np.mean([s['longitude'] for s in sensor_data])
        return lat, lon
```

#### Advantages
- Unique detection method (complements visual/radar)
- Works in all weather (no visual requirement)
- Real-time monitoring
- Can detect hidden threats (underground activity)

#### Limitations
- Requires sensor network (deployment cost)
- Limited spatial resolution
- Cannot identify specific threats (just detects events)

---

### 5.4 LIDAR Bridge (3D Point Clouds)

#### What It Is

**Converts LIDAR point cloud data into CoT terrain models or feature extraction.**

```
LIDAR system (airborne or ground-based)
  ↓ (generates 3D point clouds)
LIDAR processor
  ↓ (extracts features: buildings, terrain, vehicles)
LIDAR-to-CoT converter
  ↓ (convert features to CoT + terrain visualization)
TAK Server + CesiumJS
  ↓
ATAK clients see:
  - 3D terrain model
  - Building footprints
  - Detected vehicles/objects
  - Optimal routes
```

#### Use Cases

**Tactical Planning:**
- 3D terrain visualization
- Optimal route planning (terrain-aware)
- Building layout (infiltration planning)
- Line-of-sight analysis

**Damage Assessment:**
- Post-strike LIDAR imagery
- Change detection (new damage vs. before)
- Casualty estimates
- Rebuilding prioritization

**Urban Operations:**
- Building-level terrain
- Street layout 3D
- Optimal movement routes
- Sniper/drone positions

#### Implementation

```python
# LIDAR feature extraction bridge
import pytak
import numpy as np

class LIDARBridge(pytak.TAKClient):
    def __init__(self, lidar_data_path, **kwargs):
        super().__init__(**kwargs)
        self.lidar_data = self.load_lidar_data(lidar_data_path)
        
    def extract_buildings(self):
        """Extract building footprints from LIDAR"""
        # Simplified: cluster LIDAR points into buildings
        buildings = []
        # (Real implementation uses clustering + height analysis)
        return buildings
        
    async def process_lidar_features(self):
        """Extract features from LIDAR and create CoT events"""
        
        buildings = self.extract_buildings()
        
        for building in buildings:
            # Create CoT event for building
            lat, lon = building['centroid']
            
            event = pytak.CoT(
                uid=f"lidar-building-{building['id']}",
                cot_type='a-f-B-M-P',  # Building/facility point
                lat=lat,
                lon=lon,
                ce=building['uncertainty'],
                le=100
            )
            
            event.set_callsign(f"Bldg-{building['id']}")
            event.set_remarks(f"Height: {building['height']}m, Area: {building['area']}m²")
            
            await self.event_queue.put(event)
```

#### Advantages
- High-resolution 3D data
- Unique perspective (top-down)
- Building/terrain details
- Strategic advantage (terrain knowledge)

#### Limitations
- High data volume (point clouds massive)
- Processing-intensive
- Weather-dependent (airborne LIDAR)
- Specialized hardware ($100k+)

---

## Category 6: Infrastructure Bridges

### 6.1 SCADA/ICS Bridge

#### What It Is

**Converts Industrial Control System (SCADA) alerts into CoT events.**

```
SCADA system (power grid, water treatment, factory)
  ↓ (monitors system status and alerts)
SCADA alarm system
  ↓ (critical alert: power outage, pressure spike, etc.)
SCADA-to-CoT converter
  ↓ (alert → CoT event)
TAK Server
  ↓
ATAK clients see:
  - Critical infrastructure status
  - Alert locations
  - System outages
  - Threat indicators
```

#### Use Cases

**Critical Infrastructure Protection:**
- Power grid monitoring
- Water treatment monitoring
- Alert on sabotage/attack
- Coordinate response teams

**Facility Monitoring:**
- Building systems (HVAC, electrical)
- Facility threats (fire, intrusion)
- Real-time status on COP
- Coordinate security response

#### Implementation

```python
# SCADA bridge
import pytak

class SCADABridge(pytak.TAKClient):
    def __init__(self, scada_system, **kwargs):
        super().__init__(**kwargs)
        self.scada = scada_system
        
    async def monitor_scada_alerts(self):
        """Monitor SCADA alerts and convert to CoT"""
        
        async for alert in self.scada.subscribe_alerts():
            # Create CoT event for critical alert
            event = pytak.CoT(
                uid=f"scada-{alert['facility_id']}-{alert['alert_id']}",
                cot_type='a-f-F-C-M',  # Facility, critical
                lat=alert['facility_location']['lat'],
                lon=alert['facility_location']['lon'],
                ce=100,
                le=100
            )
            
            event.set_callsign(alert['facility_name'])
            event.set_remarks(f"Alert: {alert['message']}, Severity: {alert['severity']}")
            
            await self.event_queue.put(event)
```

#### Advantages
- Critical infrastructure alerts
- Real-time monitoring
- Threat detection (sabotage, attacks)
- Facility coordination

#### Limitations
- Facility-specific (not mobile)
- Limited to fixed locations
- SCADA systems proprietary per facility

---

## Category 7: Legacy System Bridges

### 7.1 SMS/Text Message Bridge

#### What It Is

**Converts incoming SMS messages into CoT events.**

```
Personnel text GPS coordinates or status message
  ↓ (SMS sent to bridge number)
SMS gateway (Twilio, commercial service, or self-hosted)
  ↓ (receives SMS, forwards to converter)
SMS-to-CoT converter
  ↓ (parse message, extract position + status)
TAK Server
  ↓
ATAK clients see:
  - Personnel positions (from SMS coordinates)
  - Status updates
  - Check-ins
  - SOS messages
```

#### Use Cases

**Degraded Communications:**
- When data unavailable (voice/text only)
- Simple status messages
- Position check-ins
- No special equipment needed

**Personnel Accountability:**
- Team members send "I'm OK" text
- Position logged automatically
- Status updates from field
- Integration with TAK

**Emergency Messaging:**
- SOS message triggers alert
- Position reported
- Emergency response coordinated

#### Implementation

```python
# SMS bridge
import pytak
import re

class SMSBridge(pytak.TAKClient):
    def __init__(self, sms_gateway, **kwargs):
        super().__init__(**kwargs)
        self.sms_gateway = sms_gateway
        
    async def handle_incoming_sms(self, sms):
        """Parse incoming SMS and create CoT"""
        
        # Parse SMS message format: "40.7128,-74.0060,Bravo,OK"
        # Format: lat,lon,callsign,status
        
        pattern = r'(-?\d+\.\d+),(-?\d+\.\d+),(\w+),(\w+)'
        match = re.match(pattern, sms['body'])
        
        if match:
            lat, lon, callsign, status = match.groups()
            
            event = pytak.CoT(
                uid=f"sms-{sms['from_number']}",
                cot_type='a-f-G-U-S',
                lat=float(lat),
                lon=float(lon),
                ce=100,
                le=100
            )
            
            event.set_callsign(callsign)
            event.set_remarks(f"Status: {status}, Via: SMS")
            event.set_time(sms['timestamp'])
            
            await self.event_queue.put(event)
```

#### Advantages
- Works with any phone (no special hardware)
- Degraded comms (SMS often works when data doesn't)
- Simple format (easy for operators)
- Zero training needed

#### Limitations
- SMS latency (can be slow)
- Limited message length
- Manual entry (not automatic)
- Accuracy depends on operator

---

## Complete Integration: Multi-Source TAK Feed

### Architecture: All Converters Feeding TAK

```
┌─────────────────────────────────────────────┐
│         TAK Server (Central Hub)            │
│  • Receives from all converters             │
│  • Fuses data (or passes through)           │
│  • Broadcasts to ATAK clients               │
└─────────────────────────────────────────────┘
                    ↑ (CoT XML)
         ┌──────────┼──────────┬──────────────┐
         │          │          │              │
    [ADSBcot]  [AISCot]  [PyTAK]        [Spot Bridge]
         ↑          ↑          ↑              ↑
      ADS-B       AIS      Custom       Satellite
     Aircraft   Maritime  Telemetry       Devices
      signals    signals    streams        (Spot)
         ↑          ↑          ↑              ↑
    [RTL-SDR]  [RTL-SDR]  [Sensors]    [API polling]
     Receiver   Receiver   (LIDAR,
                            Radar,
                            Seismic)
                            
         ↓ (All sources converge)
    
    ATAK Clients see unified picture:
    ✓ Aircraft (ADSBcot)
    ✓ Ships (AISCot)
    ✓ Drones (PyTAK + MAVLink)
    ✓ Ground units (Spot, Garmin, Custom)
    ✓ Sensors (Radar, Sonar, Seismic, LIDAR)
    ✓ Personnel (SMS, Cellular)
    ✓ Facilities (SCADA)
```

### Benefits of Multi-Source Architecture

1. **Resilience:** Multiple independent converters (no single point of failure)
2. **Flexibility:** Add/remove converters without affecting others
3. **Scalability:** Each converter independent (can add more)
4. **Rapid Deployment:** Each converter deployed separately
5. **Vendor Independence:** Mix commercial + open-source converters
6. **Degraded Operations:** Works even with subset of converters

---

## Summary Table: All Protocol Converters

| Converter | Input | Output | Hardware | Latency | Range | Accuracy | Use Case |
|-----------|-------|--------|----------|---------|-------|----------|----------|
| ADSBcot | ADS-B (aircraft) | CoT | RTL-SDR ($25) | 1-3s | 100-250 nm | ±20 ft | ISR, airspace |
| AISCot | AIS (maritime) | CoT | RTL-SDR ($50) | 5-60s | 20-50 nm | ±10m | Naval ops |
| PyTAK | Custom | CoT | Software | <100ms | Network | Custom | Any sensor |
| Spot Bridge | Spot API | CoT | None (API) | 10 min | Global | ±10m | Remote teams |
| Garmin Bridge | inReach API | CoT | None (API) | 2-10 min | Global | ±5m | SAR, expedition |
| Cellular Bridge | Carrier API | CoT | None (API) | 5-60s | Global | ±100m | Personnel tracking |
| SALUTE Parser | Voice/text | CoT | Software | Manual | Local | 100m | Radio ops |
| 9-Line CAS | Voice/text | CoT | Software | Manual | Local | 100m | Air ops |
| MAVLink | Telemetry | CoT | Telemetry module | 1-2s | 10-20 km | ±5m | UAV tracking |
| DJI SDK | Drone API | CoT | Commercial drone | <1s | LOS | ±5m | Commercial drones |
| Radar Bridge | Radar stream | CoT | Proprietary | Real-time | 100+ km | Variable | Air defense |
| Sonar Bridge | Sonar stream | CoT | Ship sonar | Real-time | 20-50 km | ±100m | ASW, naval |
| Seismic Bridge | Seismic network | CoT | Sensor network | Real-time | 100+ km | ±1 km | Monitoring, disaster |
| LIDAR Bridge | Point clouds | CoT/Terrain | LIDAR system | 1-10s | Variable | ±0.1m | Tactical planning |
| SCADA Bridge | Alerts | CoT | Network API | <1s | Fixed | N/A | Facility monitoring |
| SMS Bridge | SMS text | CoT | SMS gateway | 5-30s | Global | User-provided | Degraded comms |

---

## Deployment Philosophy

### When to Use Each Converter

**Real-time, automatic:**
- ADSBcot (aircraft tracking)
- AISCot (ship tracking)
- MAVLink (drone telemetry)
- Radar/Sonar/Seismic bridges
- Cellular bridge

**Resilience bypass (when ORP down):**
- ADSBcot (direct aircraft feed)
- AISCot (direct ship feed)
- PyTAK-based converters (custom data)
- Direct TAK Server connection

**Degraded operations (voice-only comms):**
- SALUTE parser (radio reports)
- 9-Line CAS (air operations)
- SMS bridge (text messaging)

**Commercial/rapid deployment:**
- Spot Bridge (satellite tracking)
- Garmin Bridge (commercial inReach)
- DJI SDK (commercial drones)

**Specialized domains:**
- Radar (air defense)
- Sonar (naval warfare)
- Seismic (monitoring/disaster)
- LIDAR (tactical planning)
- SCADA (facility protection)

---

**Captured:** 2026-07-09  
**Status:** Production-ready ecosystem  
**Relevance:** Complete protocol-to-CoT converter landscape  
**Integration:** All feed TAK Server directly (no fusion required)  
**Resilience:** Multi-source architecture, no single point of failure  
**Cost:** Ranging from free (PyTAK) to $200+ hardware + subscriptions  
**Deployment:** Minutes to hours per converter  

---

## Related Resources in Convergence

- **NSTAC_Family.md** — Core converters (ADSBcot, AISCot, PyTAK)
- **TAK Server** — Hub receiving all converter feeds
- **ATAK** — End user interface displaying unified picture
- **TAKpilot** — AI can consume converter-sourced positions
- **Evolution_to_C4ISR.md** — Sensor integration layer for ISR
- **PostGIS.md** — Spatial database (stores CoT positions from converters)

---

## Quick Reference: Add Any Protocol in 3 Steps

```bash
# Step 1: Identify input protocol
# (ADS-B, AIS, MAVLink, SCADA, etc.)

# Step 2: Use PyTAK as template
git clone https://github.com/nstac/pytak.git
cd pytak
pip install -r requirements.txt

# Step 3: Write converter (20 lines of Python)
cat > my_converter.py << 'EOF'
import pytak
import asyncio

class MyConverter(pytak.TAKClient):
    async def run(self):
        while True:
            # Read your data source here
            data = await self.read_my_protocol()
            
            # Create CoT event
            event = pytak.CoT(
                uid=f"custom-{data['id']}",
                cot_type='a-f-G-U-S',
                lat=data['lat'],
                lon=data['lon']
            )
            event.set_callsign(data['name'])
            
            # Send to TAK
            await self.event_queue.put(event)
            await asyncio.sleep(1)

# Run it
converter = MyConverter(server='tak-server', port=8087, ...)
asyncio.run(converter.run())
EOF

# Step 4: Deploy
python my_converter.py &

# That's it! Your protocol is now in TAK
```

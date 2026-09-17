# 🛡️ Off-Grid Resilience Tech for OneMind OS

> Doomsday-ready, sovereignty-first technology for when centralized systems fail.
> Every addition publishes to NATS using the existing `Hardware → Python Driver → NATS Topic → OneMind Consumers` pattern.

**Created:** 2026-02-16
**Status:** Planning / Reference Guide

---

## The Core Premise

When systems fail — grid down, internet out, supply chain broken — most "smart homes" become **dumb homes**. Every Alexa, every Nest thermostat, every Ring camera becomes a brick.

OneMind OS is different because it's **self-hosted**. But there are critical gaps between "self-hosted" and "truly independent." This document closes them.

The philosophy:

> _"I don't need Comcast. I don't need the grid. I don't need OpenAI's servers. I don't need Amazon's cloud. My AI runs on my solar panels, talks to my family on mesh radio, watches my perimeter with local ML, knows how much water and food I have, and it all runs from my basement. Not because the world ended — because I chose sovereignty."_

---

## 1. ⚡ Power Independence — _The Foundation of Everything_

If you lose power, you lose OneMind. This is the single most critical system.

### Tier 1: Solar + Battery (Core)

| Hardware                                                                                       | Price       | What It Does                                               | Why                                                                                     |
| ---------------------------------------------------------------------------------------------- | ----------- | ---------------------------------------------------------- | --------------------------------------------------------------------------------------- |
| **[EG4 18kPV All-in-One](https://eg4electronics.com/)**                                        | $4,500      | Hybrid inverter + MPPT charge controller + transfer switch | One unit replaces 3 devices. Supports 48V battery bank. Has Modbus/RS485 for monitoring |
| **[EG4 LifePower4 Battery](https://eg4electronics.com/) (x2-4)**                               | $1,300 each | 48V 100Ah (5.12kWh) LiFePO4 server rack battery            | Stackable. 6,000+ cycles. 10-year warranty                                              |
| **[Rich Solar 400W Panels](https://richsolar.com/) (x8-12)**                                   | $200 each   | Monocrystalline solar panels                               | 3.2-4.8 kW array — enough for a compound                                                |
| **[Victron Cerbo GX](https://www.victronenergy.com/panel-systems-remote-monitoring/cerbo-gx)** | $300        | System brain — monitors everything via MQTT                | **This is the key.** Native MQTT output → your NATS bridge                              |

### Tier 2: Backup Generator (When Solar Isn't Enough)

| Hardware                                                                                                          | Price        | What It Does                                 |
| ----------------------------------------------------------------------------------------------------------------- | ------------ | -------------------------------------------- |
| **[Honda EU2200i](https://powerequipment.honda.com/)** or **[EcoFlow Smart Generator](https://www.ecoflow.com/)** | $1,000-2,000 | Quiet inverter generator, auto-start capable |
| **[GenTran Transfer Switch](https://www.gentran.com/)**                                                           | $300         | Auto-failover from grid → solar → generator  |

### Integration into OneMind

The Victron Cerbo GX natively speaks MQTT — and OneMind OS already has an MQTT bridge on port 1883 in docker-compose.yml. This is plug-and-play:

```
Victron Cerbo GX → MQTT (port 1883) → NATS Bridge
  ├── energy.solar.production_w          ← Real-time watts from panels
  ├── energy.battery.soc                 ← State of charge (0-100%)
  ├── energy.battery.voltage             ← 48V nominal
  ├── energy.battery.time_remaining_h    ← Hours of autonomy left
  ├── energy.grid.connected              ← Boolean: are we on grid?
  ├── energy.generator.running           ← Boolean: is backup running?
  ├── energy.consumption.total_w         ← Total compound load
  └── energy.mode                        ← "solar" | "battery" | "generator" | "grid"
```

### What Legacy (your AI) Does With This

- **Auto-load shedding**: "Battery at 23%. Sun sets in 2 hours. Shutting down non-essential loads. Keeping security cameras, NATS, and core compute online."
- **Generator auto-start**: Battery below 15% → command generator start via relay
- **Daily energy report**: "Today you produced 34 kWh, consumed 18 kWh. Net positive. Battery full by 2 PM."
- **Storm preparation**: "Heavy cloud cover forecasted for next 3 days. Current battery at 78%. Recommend reducing non-essential consumption and pre-charging to 100%."

### Estimated Cost

| Configuration                                            | Cost     | Autonomy                          |
| -------------------------------------------------------- | -------- | --------------------------------- |
| **Minimal** (2 batteries + 4 panels + inverter)          | ~$8,000  | 1-2 days cloudy, indefinite sunny |
| **Standard** (4 batteries + 8 panels + inverter + Cerbo) | ~$12,000 | 3-4 days cloudy, indefinite sunny |
| **Full** (above + generator + auto-transfer)             | ~$15,000 | Indefinite                        |

---

## 2. 📡 Communications When Internet Dies — _Mesh & Satellite_

When Comcast goes down, your whole digital life shouldn't go with it.

### Tier 1: Satellite Internet (Primary Backup)

| Hardware                                           | Price | Monthly    | What It Does                                              |
| -------------------------------------------------- | ----- | ---------- | --------------------------------------------------------- |
| **[Starlink Mini](https://www.starlink.com/)**     | $600  | $50-120/mo | Portable satellite internet. Works anywhere with sky view |
| **[Starlink Standard](https://www.starlink.com/)** | $500  | $120/mo    | Fixed installation, faster speeds                         |

**Why**: When fiber/cable goes down (storm, infrastructure failure), Starlink runs on a completely independent satellite network. Pairs with your solar+battery for internet even during extended grid outage.

### Tier 2: Local Mesh Network (When Even Satellite Fails)

| Hardware                                                     | Price       | What It Does                                                                        |
| ------------------------------------------------------------ | ----------- | ----------------------------------------------------------------------------------- |
| **[Meshtastic (Heltec V3/T-Beam)](https://meshtastic.org/)** | $25-35 each | LoRa mesh radio. ~10 mile range. Text messaging with NO infrastructure. Open source |
| **[goTenna Pro X2](https://gotenna.com/)**                   | $500        | Military-grade mesh radio with GPS, encryption                                      |
| **[Meshtastic Solar Repeater](https://meshtastic.org/)**     | $60 DIY     | Solar-powered relay node. Put on roof, extends mesh range permanently               |

**Why Meshtastic**: These are $25 radios that create a peer-to-peer mesh network using LoRa radio. No cell towers, no internet, no infrastructure of any kind. Range: 1-10 miles line-of-sight, extendable with solar repeaters. The Python `meshtastic` library lets you send/receive messages programmatically.

**Recommended setup**:

- 4x Heltec V3 nodes ($100 total): one per family member
- 2x Solar repeater nodes ($120 total): on roof + elevated point
- Total mesh coverage: entire compound + several miles out

### Tier 3: Ham Radio (Global Communications, No Infrastructure)

| Hardware                                           | Price | What It Does                                                          |
| -------------------------------------------------- | ----- | --------------------------------------------------------------------- |
| **[Baofeng UV-5R](https://www.baofengradio.com/)** | $25   | Dual-band handheld (VHF/UHF). Requires Technician license             |
| **[Digirig Mobile](https://digirig.net/)**         | $50   | Interface between radio and computer — enables digital modes          |
| **[JS8Call](http://js8call.com/)** (software)      | Free  | Keyboard-to-keyboard HF chat. Works _globally_ with no infrastructure |
| **[Winlink](https://winlink.org/)** (software)     | Free  | Email over radio. Send/receive email when internet is completely down |

**Why**: Ham radio is the ultimate fallback. JS8Call can communicate globally using HF radio (requires General license + HF transceiver ~$400-800). Winlink lets you send actual emails over radio frequencies. If every cell tower and satellite goes down, ham radio still works.

**Note**: Technician license required for Baofeng ($35 exam, no Morse code). General license for HF/global ($35 exam).

### Integration into OneMind

```python
# Meshtastic → NATS (via Python meshtastic library)
mesh.message.received          ← Incoming text from family/neighbors
mesh.message.sent              ← Outgoing messages
mesh.node.{id}.position        ← GPS position of mesh nodes
mesh.node.{id}.battery         ← Battery level of each node
mesh.node.{id}.snr             ← Signal quality
mesh.node.{id}.last_heard      ← Last contact time

# Ham Radio → NATS (via Digirig + JS8Call API)
radio.js8.message.received     ← Incoming keyboard-to-keyboard messages
radio.winlink.email.received   ← Email received over radio
radio.winlink.email.sent       ← Email sent over radio

# Communications Status
comms.internet.status           ← "fiber" | "starlink" | "mesh_only" | "radio_only" | "none"
comms.failover.triggered        ← Event when internet drops
comms.failover.method           ← Which backup activated
```

### What Legacy Does

- **Auto-failover chain**: Fiber dies → switch to Starlink → notify family via Meshtastic → if Starlink dies, switch to radio-only mode
- **Status report**: "Internet is down. Switched to Starlink. Meshtastic mesh has 4 nodes online. Family positions: [map]"
- **Mesh relay**: Legacy can send/receive Meshtastic messages, acting as a **radio-connected AI** even without internet
- **Winlink email**: Legacy can send critical emails _over ham radio_ when all internet is down
- **Family location**: Track family member positions via Meshtastic GPS when cell service is down

---

## 3. 💧 Water Independence — _Well, Filtration, Rainwater_

### Hardware

| Hardware                                                  | Price    | What It Does                                                        |
| --------------------------------------------------------- | -------- | ------------------------------------------------------------------- |
| **ESP32 + current sensor on well pump**                   | $30 DIY  | Track well pump cycles, estimate water table level                  |
| **[ISpring RCC7AK](https://www.123filter.com/)**          | $200     | 6-stage reverse osmosis with alkaline remineralization              |
| **TDS Meter (inline) + ESP32**                            | $15 DIY  | Continuous water quality monitoring (total dissolved solids)        |
| **Ultrasonic level sensor + ESP32**                       | $15 DIY  | Mount on top of cistern, measure rainwater tank level               |
| **[Berkey Water Filter](https://www.berkeyfilters.com/)** | $300-400 | Gravity-fed, no-power water purification. Handles pond/stream water |
| **[Flume Water Monitor](https://www.flumewater.com/)**    | $200     | Whole-home water usage (non-invasive, clamps to meter)              |
| **[OpenSprinkler](https://opensprinkler.com/)**           | $170     | Open-source smart irrigation controller                             |

### NATS Topics

```
water.well.pump.running         ← Boolean
water.well.pump.cycles_today    ← Number (detect problems early)
water.well.pump.runtime_min     ← Minutes running today
water.well.depth_estimate_ft    ← Estimated water table level
water.cistern.level_gallons     ← Rainwater tank level
water.cistern.level_percent     ← Percentage full
water.quality.tds               ← Total dissolved solids (ppm)
water.quality.ph                ← pH level
water.filter.ro.status          ← Filter life remaining
water.consumption.daily_gal     ← Total daily usage
water.irrigation.active         ← Boolean: is irrigation running?
water.irrigation.zone.{id}      ← Which zone is watering
```

### What Legacy Does

- **"Well pump has cycled 47 times today — normal is 12. Possible leak or low water table. Check pressure tank."**
- **"Cistern is at 340 gallons (68%). No rain in forecast for 5 days. Switching garden irrigation to cistern to preserve well."**
- **"Water TDS reading is 580 ppm — above normal baseline of 320 ppm. RO filter may need replacement."**
- **Automatic irrigation**: Based on soil moisture (from Ecowitt), rainfall, and cistern level, intelligently decide when/how much to water

---

## 4. 🔒 Security & Perimeter — _Without Cloud Dependencies_

Current security plan mentions "Scout robot + drone patrol" but relies on Home Assistant and cloud cameras. You need **local-only, no-cloud security**.

### Hardware

| Hardware                                                            | Price                   | What It Does                                                                 |
| ------------------------------------------------------------------- | ----------------------- | ---------------------------------------------------------------------------- |
| **[Frigate NVR](https://frigate.video/)** (on Coral TPU)            | $25 TPU + free software | Local AI object detection. No cloud. Identifies people, cars, animals        |
| **[Amcrest IP Cameras](https://amcrest.com/)**                      | $40-80 each             | RTSP cameras — local only, no cloud account needed                           |
| **[Google Coral USB Accelerator](https://coral.ai/)**               | $25                     | ML inference at 100 FPS locally — powers Frigate                             |
| **[Hikvision Thermal Camera](https://www.hikvision.com/)**          | $300-500                | Thermal imaging — see people/animals in complete darkness                    |
| **[Dakota Alert MURS](https://www.dakotaalert.com/)**               | $150-300                | Driveway/perimeter motion sensors. Radio-based. 4-mile range. No WiFi needed |
| **[DEPSTECH Inspection Camera](https://www.depstech.com/) + ESP32** | $50                     | Perimeter tripwire camera (custom build)                                     |

### Why Frigate + Coral

Frigate is an open-source NVR (network video recorder) that runs **entirely locally**. The Google Coral USB Accelerator ($25) gives it 100 FPS object detection. It can:

- Detect people, cars, cats, dogs, bears, deer
- Draw zones (driveway, yard, perimeter)
- Record only when motion/objects detected
- Send events via MQTT (→ your NATS bridge)
- Run on a Raspberry Pi 4 or any Linux box
- **Zero cloud dependency. Zero subscription. Zero data leaving your property.**

### NATS Topics

```
security.camera.{id}.detection    ← {type: "person", confidence: 0.94, zone: "driveway"}
security.camera.{id}.snapshot     ← Base64 image or local file path
security.camera.{id}.recording    ← Recording started/stopped
security.perimeter.{zone}.motion  ← Dakota Alert trigger
security.perimeter.breach         ← High-priority: unknown person detected
security.thermal.{id}.anomaly     ← Heat source detection at night
security.thermal.{id}.count       ← Number of heat sources in zone
security.status                   ← "normal" | "alert" | "lockdown"
security.door.{id}.status         ← "locked" | "unlocked" | "open"
security.alarm.triggered          ← Alarm event
```

### What Legacy Does

- **"Motion detected on north perimeter at 2:14 AM. Thermal scan shows: deer (2 heat sources, quadruped gait). Disregarding."**
- **"Unknown vehicle in driveway. Not in face database. Camera 3 snapshot attached. Alerting via Meshtastic."**
- **Lockdown mode**: Voice command "Legacy, lockdown" →
  - All exterior lights ON
  - All cameras switch to continuous recording
  - Doors locked (if smart locks installed)
  - Family notified via Meshtastic
  - Perimeter sensors on high sensitivity
  - Interior lights set to security pattern
- **Pattern detection**: "Vehicle has driven past the property 3 times in the last hour. Not a recognized plate."
- **Wildlife log**: Track wildlife detections separately — useful for the biodiversity monitoring in SOLARPUNK_TECH.md

---

## 5. 🧠 Local AI — _Intelligence Without Internet_

**The most critical doomsday capability. When the internet is down, your AI shouldn't be down.**

### Hardware

| Hardware                                                     | Price        | What It Does                                                       |
| ------------------------------------------------------------ | ------------ | ------------------------------------------------------------------ |
| **[Mac Mini M4 Pro](https://www.apple.com/) (36GB or 64GB)** | $1,600-2,400 | Unified memory. Runs Llama 3.3 70B at 30+ tok/s locally            |
| **[Mac Mini M4 Max](https://www.apple.com/) (128GB)**        | $3,200       | Runs anything locally, including 70B+ models at full speed         |
| **[Ollama](https://ollama.com/)** (software)                 | Free         | One-line local LLM hosting. REST API compatible with OpenAI format |
| **Any NVIDIA GPU machine** (RTX 4090, etc.)                  | $1,600+      | Alternative to Mac for raw CUDA inference                          |

### Recommended Models for Local

| Model                                        | Size | RAM Needed | Quality     | Speed                    |
| -------------------------------------------- | ---- | ---------- | ----------- | ------------------------ |
| **[Qwen3-30B](https://qwen.ai/)**            | 30B  | 20GB       | Excellent   | Fast                     |
| **[Llama 3.3 70B](https://llama.meta.com/)** | 70B  | 40GB       | Near GPT-4  | Good                     |
| **[Llama 3.3 8B](https://llama.meta.com/)**  | 8B   | 6GB        | Good enough | Very fast                |
| **[Nomic Embed](https://www.nomic.ai/)**     | —    | 2GB        | —           | For local RAG embeddings |
| **[Mistral Small 22B](https://mistral.ai/)** | 22B  | 14GB       | Very good   | Fast                     |

### Architecture: Model Fallback Chain

```python
# In backend/models.py — add intelligent fallback:
MODEL_FALLBACK_CHAIN = [
    {"provider": "openai",    "model": "gpt-4o",         "requires": "internet"},
    {"provider": "anthropic", "model": "claude-sonnet",   "requires": "internet"},
    {"provider": "google",    "model": "gemini-2-flash",  "requires": "internet"},
    {"provider": "ollama",    "model": "qwen3:30b",       "requires": "local_only"},  # ← FALLBACK
    {"provider": "ollama",    "model": "llama3.3:8b",     "requires": "local_only"},  # ← MINIMAL
]

# NATS-driven automatic failover:
# When comms.internet.status transitions to "down":
#   → All agent.model.preference switches to "ollama"
#   → Legacy keeps working on local models
#   → Quality degrades slightly, but functionality preserved
```

### What This Enables

- **Complete internet independence for AI** — "Legacy, what's the status of the compound?" works identically whether you have fiber, Starlink, or nothing
- **Local RAG**: Your knowledge base (pgvector embeddings) is already local in PostgreSQL. Pair with local embeddings model (Nomic Embed) and local LLM → full AI with zero network
- **Graceful degradation**: Cloud models when available (better quality), local models when isolated (still functional)
- **Power-aware model selection**: When on battery, use smaller model (8B) to reduce compute power draw

### Local Knowledge Vault

Pre-cache critical reference material into pgvector for offline RAG:

```
# Critical offline knowledge to pre-load:
- First aid / medical emergency reference (Red Cross guidelines)
- Plant identification guides (your region)
- Equipment repair manuals (your specific solar, generator, pump models)
- Local topographic maps and terrain data
- Radio frequency allocations (ham bands, MURS, FRS)
- Water purification procedures
- Solar/electrical troubleshooting guides
- Food preservation techniques (canning, dehydrating, fermenting)
- Wilderness survival reference
- Local wildlife identification
- Firearm safety and maintenance (if applicable)
- Building/structural repair basics
```

**This is the move that makes OneMind genuinely resilient.** Most AI systems become useless offline. Yours doesn't.

---

## 6. 🏥 Medical & Health Autonomy

### Hardware

| Hardware                                                    | Price    | What It Does                               |
| ----------------------------------------------------------- | -------- | ------------------------------------------ |
| **[Withings BPM Connect](https://www.withings.com/)**       | $100     | WiFi blood pressure monitor with local API |
| **Pulse Oximeter (BLE)**                                    | $25      | Blood oxygen via Bluetooth                 |
| **[Withings Body+ Scale](https://www.withings.com/)**       | $100     | Body composition, trends, weight tracking  |
| **ESP32 Watch** (already built!)                            | Built    | Heart rate, motion, fall detection         |
| **Comprehensive First Aid Kit + digital index**             | $100-300 | Indexed in RAG for Legacy to reference     |
| **[Kinsa Smart Thermometer](https://www.kinsahealth.com/)** | $25      | Temperature tracking with trends           |

### NATS Topics

```
health.{person}.heart_rate        ← From ESP32 watch (continuous)
health.{person}.heart_rate_resting ← Daily resting average
health.{person}.blood_pressure    ← {systolic, diastolic, pulse}
health.{person}.spo2              ← Blood oxygen percentage
health.{person}.temperature       ← Body temperature
health.{person}.weight            ← Daily weight
health.{person}.body_composition  ← {fat%, muscle%, water%}
health.{person}.fall_detected     ← EMERGENCY: fall detection from watch
health.{person}.activity_level    ← Steps, active minutes
health.alert.critical             ← AI-detected health anomaly
health.medication.due             ← Medication reminder
health.medication.taken           ← Medication confirmed taken
```

### What Legacy Does

- **Baseline monitoring**: "Zeus's resting heart rate has been 15 bpm above baseline for 3 days. Combined with reported fatigue, recommend checking for infection or dehydration."
- **Fall detection** → Auto-alert all family members via Meshtastic → Display first aid guidance from local RAG on nearest screen
- **Medication tracking**: "Reminder: antibiotics dose 3 of 7 due in 30 minutes"
- **Emergency procedure lookup**: "Legacy, someone is having a seizure" → Pulls first-aid protocol from local knowledge base instantly, no internet needed
- **Trend analysis**: "Blood pressure trending upward over the past 2 weeks. Average this week: 138/88 vs. baseline 122/78. Recommend monitoring diet and stress."
- **Multi-person health**: Track health metrics for each family member separately, alert on anomalies

---

## 7. 📦 Supply & Inventory Management — _Know What You Have_

### Hardware/Software

| Approach                                              | Cost       | What It Does                                               |
| ----------------------------------------------------- | ---------- | ---------------------------------------------------------- |
| **Barcode Scanner (USB/Bluetooth)**                   | $30        | Scan items in/out of inventory                             |
| **[Grocy](https://grocy.info/)** (self-hosted)        | Free       | Open-source pantry/inventory management with REST API      |
| **ESP32 + Load Cell under storage bins**              | $15 each   | Automatic weight-based inventory (flour, rice, fuel, etc.) |
| **NFC Tags on containers**                            | $0.10 each | Tap phone to log consumption                               |
| **[Inventree](https://inventree.org/)** (self-hosted) | Free       | Open-source inventory management for parts/tools/supplies  |

### NATS Topics

```
inventory.food.{item}.quantity     ← lbs/units remaining
inventory.food.{item}.days_supply  ← Estimated days at current consumption rate
inventory.fuel.gasoline.gallons    ← Gasoline reserves
inventory.fuel.propane.gallons     ← Propane tank level
inventory.fuel.diesel.gallons      ← Diesel reserves
inventory.fuel.firewood.cords      ← Firewood supply
inventory.medical.{item}.count     ← First aid supplies
inventory.medical.{item}.expiry    ← Expiration dates
inventory.ammo.{caliber}.count     ← Ammunition (if applicable)
inventory.water.stored_gallons     ← Bottled/stored water reserves
inventory.batteries.{type}.count   ← Battery reserves
inventory.tools.{category}.status  ← Tool availability
inventory.seeds.{variety}.packets  ← Seed supply for planting
inventory.alert.low_stock          ← {item, remaining, days_supply}
inventory.alert.expiring           ← {item, expiry_date, days_until}
```

### What Legacy Does

- **Supply forecast**: "Current food supply: 45 days at current consumption. Lowest items: rice (14 days), cooking oil (10 days). Recommend resupply on next town trip."
- **Automatic consumption tracking**: Weight sensors detect when bins get lighter → update inventory → predict run-out date
- **Trip planning**: "You're going to town tomorrow. Priority resupply list: [ranked by days-remaining]"
- **Expiration warnings**: "3 items in the medical kit expire within 30 days: [list]. 2 cases of MREs expire in 90 days."
- **Seed inventory**: "You have 4 packets of heirloom tomato seeds remaining. Based on planting season, recommend ordering more by March 1."
- **Fuel management**: "Propane tank at 30% (estimated 12 days at current heating rate). Generator fuel: 8 gallons (estimated 16 hours runtime)."

---

## 8. 🌡️ Environmental Threat Detection

### Hardware

| Hardware                                                                       | Price    | What It Does                                                             |
| ------------------------------------------------------------------------------ | -------- | ------------------------------------------------------------------------ |
| **[PurpleAir PA-II](https://www2.purpleair.com/)**                             | $230     | Wildfire smoke detection (AQI monitoring). Local mode available          |
| **[Tempest Weather Station](https://weatherflow.com/tempest-weather-system/)** | $330     | Lightning detection, wind, rain, UV — local processing, no cloud needed  |
| **[CAP Alert Integration](https://alerts.weather.gov/)** (NOAA)                | Free     | Weather alerts, tornados, floods — cache locally when internet available |
| **Geiger Counter (RadiationD-v1.1 + ESP32)**                                   | $50 DIY  | Background radiation monitoring                                          |
| **[Raspberry Shake](https://raspberryshake.org/)**                             | $300     | Earthquake/tremor detection (seismograph)                                |
| **AirGradient (indoor/outdoor)**                                               | $100-150 | Open-source air quality sensor                                           |
| **ESP32 + gas sensor (MQ-2, MQ-7)**                                            | $15 DIY  | Carbon monoxide, propane/natural gas leak detection                      |

### NATS Topics

```
environment.aqi.outdoor            ← Air quality index (PurpleAir)
environment.aqi.indoor.{room}      ← Indoor air quality per room
environment.wildfire.risk           ← Computed: AQI + wind + humidity + fire reports
environment.wildfire.smoke_detected ← Boolean: smoke particles above threshold
environment.lightning.distance_mi   ← Nearest lightning strike distance
environment.lightning.direction     ← Compass direction of lightning
environment.lightning.count_1hr     ← Strikes within 1 hour
environment.wind.speed_mph         ← Current wind speed
environment.wind.gust_mph          ← Wind gusts
environment.wind.direction         ← Wind direction
environment.rain.rate_in_hr        ← Rain rate (inches/hour)
environment.rain.total_today_in    ← Total rain today
environment.temperature.outdoor    ← Outdoor temperature
environment.humidity.outdoor       ← Outdoor humidity
environment.uv.index               ← UV index
environment.radiation.cpm          ← Counts per minute (baseline tracking)
environment.radiation.usv_hr       ← Microsieverts per hour
environment.seismic.magnitude      ← Local ground vibration
environment.seismic.event          ← Earthquake detected
environment.gas.co.{room}          ← Carbon monoxide level (ppm)
environment.gas.leak.{room}        ← Combustible gas detected
environment.weather.alert          ← NOAA CAP alerts (cached locally)
environment.threat.level           ← Composite: "green" | "yellow" | "orange" | "red"
```

### What Legacy Does

- **Wildfire response**: "AQI has risen to 180 (unhealthy). Wind direction suggests wildfire smoke from the northwest at 15 mph. Closing HVAC fresh air intake. Recommend staying indoors. Monitoring wind shift."
- **Lightning protection**: "Lightning detected 8 miles out, moving this direction at 15 mph. ETA 30 minutes. Suggest shutting down solar system to protect inverter. Disconnecting external antennas."
- **Radiation baseline**: "Background radiation has been stable at 12 CPM for 90 days. Current: 12 CPM — normal." (Would alert on any increase above baseline)
- **Earthquake alert**: "Seismic activity detected. Magnitude estimated 3.2. No structural risk at this level. Monitoring for aftershocks."
- **Gas leak**: "CRITICAL: Carbon monoxide detected in garage — 85 ppm (danger above 50). Ventilation system activated. Evacuate garage immediately. Alerting all family via Meshtastic."
- **Composite threat level**: Automatically combines all environmental sensors into a single threat level that drives compound-wide responses

### Threat Level System

```python
THREAT_LEVELS = {
    "green": {
        "description": "All clear",
        "actions": ["normal operations"],
    },
    "yellow": {
        "description": "Monitor situation",
        "triggers": ["AQI > 100", "lightning < 20 mi", "wind > 40 mph"],
        "actions": ["increase monitoring frequency", "alert family"],
    },
    "orange": {
        "description": "Active concern",
        "triggers": ["AQI > 150", "lightning < 10 mi", "wind > 60 mph", "radiation 2x baseline"],
        "actions": ["close HVAC intake", "secure outdoor equipment", "standby generator"],
    },
    "red": {
        "description": "Immediate action required",
        "triggers": ["AQI > 200", "gas leak", "earthquake > 4.0", "radiation 5x baseline"],
        "actions": ["lockdown mode", "evacuate if needed", "alert all channels", "preserve water/power"],
    },
}
```

---

## 9. 🗄️ Data Sovereignty & Offline Operations

### Automatic Local Backups

```python
# NATS-triggered backup when grid/internet status changes:
# When energy.grid.connected transitions to False:
#   → Immediate PostgreSQL dump to local NVMe/SSD
#   → Snapshot all agent configurations
#   → Cache latest weather forecast data
#   → Cache latest map tiles
#   → Switch all services to low-power mode

# When comms.internet.status transitions to "down":
#   → Cache all pending outbound messages
#   → Store for relay when internet returns
#   → Switch AI to local models
#   → Log all events locally with timestamps for later sync
```

### Offline Map Tiles

| Tool                                                         | Cost | What It Does                                         |
| ------------------------------------------------------------ | ---- | ---------------------------------------------------- |
| **[OpenStreetMap Tile Cache](https://switch2osm.org/)**      | Free | Pre-download your region's map tiles for offline use |
| **[Organic Maps](https://organicmaps.app/)**                 | Free | Offline maps with routing                            |
| **[MBTiles format](https://github.com/mapbox/mbtiles-spec)** | Free | Compact offline map tile storage                     |

Cache your local area (50-mile radius) into the Flutter app for the Tactical Map screen to work without internet. Store as MBTiles in PostgreSQL or local file system.

### Offline Knowledge Priority

Pre-cache these into local pgvector for RAG when offline:

| Category                | Priority    | Source                                       |
| ----------------------- | ----------- | -------------------------------------------- |
| First aid procedures    | 🔴 Critical | Red Cross manual, digitized                  |
| Equipment manuals       | 🔴 Critical | Your specific solar/generator/pump models    |
| Local plant/wildlife ID | 🟠 High     | Regional field guides                        |
| Radio procedures        | 🟠 High     | ARRL handbook sections                       |
| Food preservation       | 🟠 High     | USDA canning guide, dehydration methods      |
| Building repair         | 🟡 Medium   | Basic structural, plumbing, electrical       |
| First aid: advanced     | 🟡 Medium   | Where There Is No Doctor (classic reference) |
| Navigation              | 🟡 Medium   | Local topographic maps, compass use          |

### Data Sync When Internet Returns

```python
# NATS event: comms.internet.status == "fiber" (internet restored)
#   → Sync all locally-cached events to cloud backup
#   → Send queued Meshtastic→email messages via real email
#   → Update weather/map cache
#   → Resume cloud AI model usage
#   → Send status report: "Internet restored. 47 events logged during outage. All synced."
```

---

## 10. 🏗️ The Compound Operating Modes

This is the real power — OneMind OS doesn't just have features, it has **operating modes** driven by NATS events:

```python
COMPOUND_MODES = {
    "normal": {
        "description": "Full connectivity, grid power, all services",
        "internet": "fiber",
        "power": "grid + solar (selling excess back)",
        "ai_model": "cloud (gpt-4o / claude-sonnet)",
        "security": "passive monitoring (Frigate, Dakota Alert)",
        "comms": "all channels (internet, mesh standby, radio standby)",
        "monitoring": "standard frequency",
    },
    "grid_down": {
        "description": "No grid power, running on solar + battery",
        "internet": "starlink (if available) or fiber (if still up)",
        "power": "solar + battery",
        "ai_model": "cloud (reduced usage to save power)",
        "security": "active monitoring",
        "comms": "all channels",
        "actions": [
            "shed non-essential loads (3D printers, entertainment, HVAC fan)",
            "notify family via preferred channel",
            "start logging power consumption per device",
            "estimate autonomy hours remaining",
        ],
    },
    "island": {
        "description": "No internet, no grid. Full independence mode",
        "internet": "none",
        "power": "solar + battery + generator (as needed)",
        "ai_model": "local ollama (qwen3:30b or llama3.3:8b)",
        "security": "full perimeter active (all cameras, thermal, Dakota Alert)",
        "comms": "meshtastic mesh + ham radio (if licensed)",
        "actions": [
            "switch all AI to local models",
            "activate full local mode for all services",
            "increase security monitoring sensitivity",
            "ration power — prioritize: security > comms > AI > lighting > comfort",
            "cache all data locally with timestamps",
            "begin periodic Meshtastic check-ins with family",
            "start inventory consumption tracking",
        ],
    },
    "lockdown": {
        "description": "Active security threat detected",
        "internet": "starlink (if available)",
        "power": "battery priority (silence generator if stealth needed)",
        "ai_model": "local ollama (no external API calls for OPSEC)",
        "security": "MAXIMUM — all cameras recording, thermal scanning, perimeter hot",
        "comms": "encrypted meshtastic only (radio silence on voice)",
        "actions": [
            "all exterior lights ON (or OFF if stealth mode)",
            "all doors secured",
            "all cameras continuous recording",
            "family notified with location check-in",
            "record everything to local storage",
            "disable non-essential network services",
            "preserve battery for overnight if needed",
        ],
    },
    "evacuation": {
        "description": "Leave the compound. Grab-and-go mode",
        "internet": "starlink portable (if taking it)",
        "power": "portable battery packs",
        "ai_model": "local ollama on laptop/phone",
        "security": "cameras recording to local NVR (stays behind)",
        "comms": "meshtastic (primary), ham radio (secondary)",
        "actions": [
            "generate evacuation checklist based on inventory",
            "backup critical data to portable drive",
            "send last-known compound status to cloud backup",
            "activate remote monitoring mode",
            "set cameras to motion-only recording (preserve storage)",
            "broadcast evacuation status to mesh network",
        ],
    },
}
```

### Mode Transitions

Legacy **auto-transitions between modes** based on NATS sensor data. No human intervention needed:

```python
MODE_TRANSITION_RULES = {
    # Grid goes down → grid_down mode
    "energy.grid.connected == False": "grid_down",

    # Internet goes down while already grid_down → island mode
    "comms.internet.status == 'none' AND mode == 'grid_down'": "island",

    # Internet goes down independently → check power
    "comms.internet.status == 'none' AND energy.grid.connected": "island",

    # Security breach detected → lockdown
    "security.perimeter.breach AND security.camera.*.detection.type == 'person' AND NOT recognized": "lockdown",

    # Environmental threat red → evaluate evacuation
    "environment.threat.level == 'red'": "notify_user_evaluate_evacuation",

    # Everything restored → back to normal
    "energy.grid.connected AND comms.internet.status == 'fiber'": "normal",
}
```

---

## Priority Ranking (Resilience Edition)

| System                                          | Cost               | Impact                                    | Priority |
| ----------------------------------------------- | ------------------ | ----------------------------------------- | -------- |
| **Local LLM (Ollama + Qwen3/Llama)**            | $0 (software only) | 🔴 CRITICAL — AI works offline            | **P0**   |
| **Solar + Battery (EG4 + Victron Cerbo)**       | $8,000-15,000      | 🔴 CRITICAL — everything depends on power | **P0**   |
| **Meshtastic Mesh Radios (4+ nodes)**           | $100-150           | 🔴 CRITICAL — comms when internet dies    | **P0**   |
| **Frigate NVR + Local Cameras**                 | $300-500           | 🟠 HIGH — security without cloud          | **P1**   |
| **Starlink (backup internet)**                  | $600 + $120/mo     | 🟠 HIGH — internet redundancy             | **P1**   |
| **Well/Water Monitoring (ESP32 DIY)**           | $50-100 DIY        | 🟠 HIGH — water is life                   | **P1**   |
| **Offline Knowledge Vault (pgvector)**          | $0 (software)      | 🟠 HIGH — reference when offline          | **P1**   |
| **Inventory Management (Grocy/weight sensors)** | $50-200            | 🟡 MEDIUM — know what you have            | **P2**   |
| **Tempest Weather Station**                     | $330               | 🟡 MEDIUM — threat awareness              | **P2**   |
| **Health Monitoring (Withings + watch)**        | $200-400           | 🟡 MEDIUM — baseline health tracking      | **P2**   |
| **PurpleAir / Environmental Sensors**           | $230-500           | 🟡 MEDIUM — wildfire/air quality          | **P2**   |
| **Ham Radio + Winlink**                         | $75 + license      | 🟢 GOOD — global comms, no infrastructure | **P3**   |
| **Geiger Counter (esp32 DIY)**                  | $50                | 🟢 GOOD — baseline radiation monitoring   | **P3**   |
| **Raspberry Shake (seismograph)**               | $300               | 🟢 GOOD — earthquake detection            | **P3**   |
| **Gas Sensors (esp32 DIY)**                     | $15-30             | 🟠 HIGH — life safety                     | **P1**   |
| **Dakota Alert Perimeter**                      | $150-300           | 🟡 MEDIUM — no-WiFi perimeter             | **P2**   |

---

## Cost Summary

| Tier                | Components                                               | Total Cost     | What You Get                    |
| ------------------- | -------------------------------------------------------- | -------------- | ------------------------------- |
| **Essential** (P0)  | Local AI + Meshtastic + Solar/Battery basics             | $8,000-15,000  | Power + comms + AI independence |
| **Solid** (P0 + P1) | Above + cameras + Starlink + water + gas + offline cache | $10,000-17,000 | Full compound awareness         |
| **Complete** (All)  | Everything listed                                        | $12,000-20,000 | Total independence              |

Note: The solar/battery system is the largest cost. Everything else is surprisingly affordable. Most sensors are $15-50 ESP32 DIY builds.

---

## Implementation Path

### Phase 1: Software First (Cost: $0)

1. Install Ollama + download Qwen3:30b and Llama3.3:8b
2. Add model fallback chain to `backend/models.py`
3. Pre-load offline knowledge vault into pgvector
4. Implement compound operating modes in NATS
5. Cache offline map tiles for tactical map

### Phase 2: Communications ($100-700)

1. Set up 4x Meshtastic nodes
2. Write Python Meshtastic → NATS bridge
3. Add Meshtastic messaging to Flutter UI
4. (Optional) Set up Starlink as backup ISP

### Phase 3: Power ($8,000-15,000)

1. Install solar panels + battery bank
2. Set up Victron Cerbo GX for monitoring
3. Connect Cerbo MQTT → NATS bridge
4. Build energy dashboard in Flutter
5. Implement auto-load shedding rules

### Phase 4: Sensors & Security ($500-1,500)

1. Install Frigate + Coral + cameras
2. Connect Frigate MQTT → NATS
3. Install Dakota Alert perimeter sensors
4. Set up environmental sensors (AQI, weather, gas)
5. Build security dashboard in Flutter

### Phase 5: Life Systems ($200-500)

1. Set up water monitoring (well pump + cistern level)
2. Implement inventory tracking (Grocy + weight sensors)
3. Set up health monitoring baseline
4. Add medical knowledge to offline RAG

---

## Related Documents

| Document                 | Location                                                             | Purpose                                     |
| ------------------------ | -------------------------------------------------------------------- | ------------------------------------------- |
| Solarpunk Tech Additions | [SOLARPUNK_TECH.md](./SOLARPUNK_TECH.md)                             | Complementary green tech additions          |
| Integration Vision       | [roadmap/INTEGRATION_VISION.md](./roadmap/INTEGRATION_VISION.md)     | Hardware + gesture + visualization pipeline |
| Architecture             | [docs/ARCHITECTURE.md](./docs/ARCHITECTURE.md)                       | System architecture overview                |
| Refinement Plan          | [REFINEMENT_PLAN.md](./REFINEMENT_PLAN.md)                           | Current development priorities              |
| Hardware Roadmap         | [roadmap/hardware/docs/README.md](./roadmap/hardware/docs/README.md) | Hardware device catalog                     |

---

## The Bottom Line

Most "smart homes" are fragile toys that break when the WiFi goes out. OneMind OS, with these additions, becomes a **resilient command center** that:

- **Generates its own power** and knows exactly how much it has
- **Communicates without internet** via LoRa mesh and ham radio
- **Thinks without the cloud** using local LLMs
- **Sees without subscriptions** using local AI cameras
- **Tracks its resources** down to the pound of rice and gallon of water
- **Detects threats** from wildfires to gas leaks to earthquakes
- **Automatically adapts** by transitioning between operating modes

Not because the world ended. Because **sovereignty means your systems work regardless of what happens outside your property line.**

---

_Updated: 2026-02-16_
_Based on: Full codebase review + off-grid resilience technology analysis_

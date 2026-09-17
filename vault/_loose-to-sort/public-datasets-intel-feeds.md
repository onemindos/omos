# Public Datasets → OneMind Intel Feeds
# Source: https://github.com/awesomedata/awesome-public-datasets

---

## Already Running as Live NATS Feeds

| Dataset | Feed Worker | NATS Subject |
|---|---|---|
| NOAA Climate / Weather Alerts | feed-intel-noaa | intel.weather.noaa.alert.* |
| USGS Earthquake Archives | feed-intel-usgs | intel.seismic.usgs.* |
| NASA FIRMS Active Fire | feed-intel-firms | intel.fire.nasa.firms.* |
| ACLED Conflict Events | feed-intel-acled | intel.conflict.acled.* |
| GDACS Global Disasters | feed-intel-gdacs | intel.disaster.gdacs.* |
| EMSC Seismic (European) | feed-intel-emsc | intel.seismic.emsc.* |

---

## High-Value Additions — TAK Operational Picture

### EarthScience
- **Smithsonian Global Volcano & Eruption Database** → intel.volcanic.smithsonian.*
  - Simple API, global coverage, real-time eruption alerts, high visual impact on map
- **USGS Stream Gauge / Flood Stage** → intel.flood.usgs.*
  - Real-time river levels, flood warnings, SAR-critical
- **Global Wind Atlas**
  - Wind layer for drone flight planning and fire spread modeling

### GIS — Directly Paintable on TAK
- **OpenStreetMap (OSM)** — buildings, roads, infrastructure as base overlays in CloudTAK
- **Natural Earth** — country/region boundaries, coastlines
- **Homeland Infrastructure Foundation-Level Data (HIFLD)** → intel.infrastructure.hifld.*
  - Power plants, cell towers, hospitals, bridges, water treatment
  - Makes your TAK map look like what FEMA runs
  - Static dataset, permanent overlay, extremely operationally relevant
- **Global Administrative Areas (GADM)** — precise boundary polygons, every country/state/county
- **Shadow Accrual Maps** — shadow coverage by time of day → drone/surveillance ops planning
- **TIGER/Line** — US road network for nav overlays

### Government
- State/county 911 dispatch feeds, road closures, evacuation orders → CoT-inject to TAK in real time
- Every US state has an open data portal — most publish incident/emergency data

### Energy
- **Global Power Plant Database** — every power plant on Earth, lat/lon, fuel type, capacity
  - Paint as infrastructure markers on TAK
  - Critical for SAR near power infrastructure, wildfire crews managing transmission lines
- **EIA** — real-time US grid status

### Biology / Environment
- **iNaturalist** — wildlife sightings with GPS coordinates → intel.wildlife.inaturalist.*
  - Directly relevant to RADS/homestead: animal intrusions, endangered species locations
  - "A mountain lion was sighted 3 miles away — it's on your map" is immediately real

### Agriculture → RADS / Homestead Digital Twin
- **USDA FoodData Central**
- **USDA PLANTS Database**
- **USDA NASS County Crop Yields**
- All feed into omos-rads soil/land agent context

### CyberSecurity → Guardian Agent
- OONI censorship/network interference data
- Rapid7 Sonar internet scan data
- Feed Guardian's threat awareness model

---

## Three Feeds to Build First (Priority Order)

### 1. HIFLD Critical Infrastructure
- Source: https://hifld-geoplatform.opendata.arcgis.com/
- Every hospital, power plant, cell tower, bridge, water treatment facility in the US
- Static GeoJSON → one-time load + periodic refresh
- Subject: intel.infrastructure.hifld.{type}.{id}
- Impact: TAK map immediately looks like a professional operations platform

### 2. Smithsonian Volcano Database
- Source: https://volcano.si.edu/
- Global volcano locations + real-time eruption alerts via GeoRSS/API
- Subject: intel.volcanic.smithsonian.{region}.{volcano_id}
- Impact: global coverage, visually dramatic, operationally real

### 3. iNaturalist Wildlife Sightings
- Source: https://www.inaturalist.org/pages/api+reference
- GPS-tagged wildlife observations by species
- Subject: intel.wildlife.inaturalist.{taxon}.{observation_id}
- Impact: connects homestead/RADS story, real value for land owners, unique on TAK

---

## Education Architecture

### Course: "Build a Live Sovereign Intelligence Platform from Public Data"

Each lesson follows the same pattern (mirrors the actual stack):
1. Pick a public dataset from the catalog
2. Write a feed worker in Python — publish to NATS intel.*
3. Watch it appear on TAK via sink-tak (automatic — no extra code)
4. Query history in ClickHouse (onemind.events table)
5. Ask Legacy to analyze what came in

Every student ends up running a mini version of the OneMind stack.
The dataset catalog is the supply chain for the curriculum.

### Why This Works as a Business
- Free data sources = zero cost to the curriculum
- Students get a working system, not a theory
- Each feed they build is a real skill: Python, NATS, JetStream, ClickHouse, TAK CoT
- The stack they build is the same stack used in real field operations
- Education → community → operators who run real missions using the platform

---

## The Bigger Picture

TAK Server is the front door.
NATS is what makes TAK a platform instead of an app.
Public datasets are the fuel.
AI agents (Legacy, Oracle, Guardian) are the analysts.
The community runs real missions and teaches each other.

### Who You Serve
- Search & rescue teams (already on TAK, need what's next)
- Wildland fire / incident command
- Rural land owners and farms (drone monitoring, soil sensors, perimeter)
- Off-grid and resilient communities (sovereign comms, no cloud dependency)
- First responder units (small departments, can't afford enterprise C2)
- Tactical preparedness community (large, organized, already on ATAK)

### Three Tiers of Teaching
- **Tier 1 — Sovereign Awareness**: Run your own TAK server + AI-powered data feeds
- **Tier 2 — AI-Augmented Operations**: Voice commands to AI agent acting on the map, mission planning with Legacy
- **Tier 3 — Autonomous Field Systems**: Drone MAVLink → NATS → TAK track, computer vision as CoT events, robotic sensors feeding ClickHouse, AR overlays

---

## Reference Links

- Awesome Public Datasets: https://github.com/awesomedata/awesome-public-datasets
- HIFLD Open Data: https://hifld-geoplatform.opendata.arcgis.com/
- Smithsonian Volcano: https://volcano.si.edu/
- iNaturalist API: https://www.inaturalist.org/pages/api+reference
- USGS Water Services: https://waterservices.usgs.gov/
- Global Power Plant Database: https://datasets.wri.org/dataset/globalpowerplantdatabase
- OpenStreetMap: https://www.openstreetmap.org/
- Natural Earth: https://www.naturalearthdata.com/
- GADM: https://gadm.org/
- TIGER/Line: https://www.census.gov/geographies/mapping-files/time-series/geo/tiger-line-file.html
- Shadow Accrual Maps: https://github.com/nicholasjclark/shadow_accrual_maps

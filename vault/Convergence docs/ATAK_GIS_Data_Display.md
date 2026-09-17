# ATAK & WinTAK: How To Display GIS Data — Complete Guide

**Link:** https://mappingsupport.com/p2/atak/display_gis_data_with_atak.html

**Author:** Joseph Elfelt (MappingSupport.com)  
**Focus:** Practical GIS data integration with ATAK/WinTAK  
**Resource Type:** Step-by-step tutorials + PDF guides  
**Format:** 7 comprehensive PDFs with examples  
**Target Users:** ATAK operators, mission planners, GIS analysts, military personnel  
**Status:** Active, regularly maintained  

---

## What It Is

A comprehensive knowledge base and step-by-step guide for displaying Geographic Information System (GIS) data in ATAK (Android Tactical Assault Kit) and WinTAK (Windows version). Focuses on practical methods for integrating various GIS data sources without requiring plugins.

**Core Purpose:** Enable tactical operators and planners to leverage government GIS data repositories (especially ArcGIS servers) for mission planning, terrain analysis, and situational awareness.

**Philosophy:** "No plug-in needed" — Use standard ATAK/WinTAK features to display any GIS data source.

---

## What It Covers: Seven Key Areas

### 1. **ATAK/WinTAK & GIS Overview** (PDF 1)
- Background on ATAK, WinTAK, and GIS concepts
- Options for displaying GIS data (basemaps vs overlays)
- How to find ArcGIS server addresses
- Using ArcGIS server table of contents
- Basic steps for adding/deleting basemaps and overlays
- **XML syntax examples** (add custom basemaps)
  - High-quality Open Street Map (OSM) basemap
  - Best USA road overlay
- **Key learning:** How ArcGIS servers expose data + how to consume it

### 2. **Summary Tables: Basemaps & Overlays** (PDF 2)
- Quick-reference decision tables
- Shows options based on:
  - **Basemap vs overlay choice**
  - **Do you have XML link or capabilities file?**
  - **Do you have XML content or capabilities content?**
  - **Do you want GIS attribute access?**
- Use table to decide fastest path to display data
- **Key learning:** Decision tree for fastest integration

### 3. **ArcGIS Servers: WMS, XYZ, & More** (PDF 3)
- Deep dive: How to use ArcGIS server table of contents
- Recognize whether data can be displayed via XML or capabilities file
- **XML syntax examples:**
  - Single ArcGIS data source
  - Multiple ArcGIS data sources combined
  - Custom styling (colors, transparency)
- How to identify WMS vs. XYZ vs. other service types
- **Key learning:** Master ArcGIS server integration

### 4. **ArcGIS Download via KML Network Links** (PDF 4)
- Use KML network link feature (built into ATAK/WinTAK)
- Download vector data directly from ArcGIS server
- Convert ArcGIS features → KML → ATAK overlay
- Practical workflow:
  1. Identify feature service in ArcGIS
  2. Create KML query URL
  3. Add as network link to ATAK
  4. Data syncs automatically
- **Key learning:** Live data updates from server

### 5. **XYZ, WMS, WMTS, & WFS Servers** (PDF 5)
- Comprehensive guide for non-ArcGIS sources
- **XYZ tile servers** (Slippy map format)
  - What they are: Simple grid of tile images
  - How to use: URL format, parameter substitution
  - Examples: OSM, various commercial providers
  
- **WMS servers** (Web Map Service, OGC standard)
  - What they are: Dynamic map generation
  - How to use: Request styled maps from server
  - Parameters: Layers, styles, transparency
  
- **WMTS servers** (Web Map Tile Service)
  - What they are: Pre-rendered tile cache
  - How to use: Reference tile matrix set
  - Examples: Government mapping services
  
- **WFS servers** (Web Feature Service)
  - What they are: Vector feature access
  - How to use: Query features, download as GeoJSON/KML
  - Parameters: Feature types, filters, bounding box
  
- **Key learning:** Understand standard OGC service types

### 6. **GIS Feature Requests & Improvements** (PDF 6)
- Author's suggestions for missing ATAK features
- "Big impact" enhancements that would improve GIS integration:
  - Better basemap management UI
  - Improved WMS capabilities parsing
  - Layer styling controls
  - Performance optimizations
  - Attribute query interface
- **Key learning:** Known limitations + requested improvements

### 7. **ATAK for Recreation** (PDF 7)
- Using ATAK for outdoor activities (hiking, climbing, etc.)
- Recreation GIS data sources:
  - **USFS** (U.S. Forest Service) — trails, campgrounds, wilderness areas
  - **BLM** (Bureau of Land Management) — public lands, recreational areas
  - **NPS** (National Park Service) — park boundaries, facilities
  - **USGS** (U.S. Geological Survey) — topographic maps, DEM
  - Commercial providers: Garmin, CalTopo, etc.
- **Key learning:** Civilian use cases, data availability

---

## GIS Data Sources by Type

### Government GIS Servers (ArcGIS-based)

**Common government ArcGIS servers:**
```
USGS: https://server.arcgisonline.com/ArcGIS/rest/services
NOAA: https://gis.ngdc.noaa.gov/arcgis/rest/services
FEMA: https://gis.fema.gov/arcgis/rest/services
USDA: https://apps.fs.usda.gov/arcx/rest/services
BLM: https://gis.blm.gov/arcgis/rest/services
DOD: https://services.arcgisonline.com/arcgis/rest/services/Military

Each server exposes:
  - Basemaps (satellite, topo, road)
  - Overlays (hazards, infrastructure, administrative boundaries)
  - Feature layers (queryable, styleable)
```

### Standard GIS Service Types

| Type | Use | Example |
|------|-----|---------|
| **XYZ tile** | Fast basemaps | OpenStreetMap, satellite imagery |
| **WMS** | Dynamic layers | Government overlays, analysis products |
| **WMTS** | Cached tiles | Pre-rendered high-quality maps |
| **WFS** | Vector features | Searchable points, boundaries, infrastructure |
| **KML** | Vector overlay | Custom shapes, routes, positions |
| **GeoJSON** | Web vector format | Modern standard, browser-friendly |

---

## Technical Integration Methods

### Method 1: Add Basemap (XML Configuration)

```xml
<!-- Example: Add Open Street Map basemap -->
<overlay name="Open Street Map" 
         type="online"
         url="https://tile.openstreetmap.org/{z}/{x}/{y}.png"
         maxZoom="18"
         attribution="© OpenStreetMap contributors">
</overlay>
```

**Steps:**
1. Get basemap URL (XYZ format)
2. Create XML entry
3. Add to ATAK/WinTAK basemap folder
4. Select in app

### Method 2: Add WMS Overlay

```xml
<!-- Example: Add WMS layer from government server -->
<overlay name="Road Overlay"
         type="wms"
         url="https://server.arcgisonline.com/arcgis/services/World_Street_Map"
         layers="0"
         format="image/png"
         transparent="true">
</overlay>
```

**Steps:**
1. Get WMS capabilities URL
2. Choose layers from table of contents
3. Create XML with layer selection
4. Add to ATAK/WinTAK

### Method 3: Add Network Link (KML)

```xml
<!-- Example: Add live vector data from ArcGIS feature service -->
<NetworkLink>
  <name>Enemy Positions</name>
  <Link>
    <href>https://services.arcgisonline.com/arcgis/rest/services/...
    /query?where=1%3D1&outFields=*&f=kml</href>
    <refreshMode>onInterval</refreshMode>
    <refreshInterval>300</refreshInterval> <!-- 5 minutes -->
  </Link>
</NetworkLink>
```

**Steps:**
1. Identify feature service in ArcGIS
2. Create feature query URL
3. Wrap in KML NetworkLink
4. Add to ATAK/WinTAK
5. Data updates automatically on interval

### Method 4: Custom Layer Styling

```xml
<!-- Example: Color buildings by height -->
<overlay name="Buildings"
         url="..."
         type="wms">
  <style name="height">
    <rule condition="height &gt; 30">
      <color>#FF0000</color> <!-- Red: tall buildings -->
    </rule>
    <rule condition="height &gt; 15">
      <color>#FFA500</color> <!-- Orange: medium -->
    </rule>
    <rule condition="height &lt;= 15">
      <color>#FFFF00</color> <!-- Yellow: short -->
    </rule>
  </style>
</overlay>
```

---

## Integration with OneMind OS / Tactical Operations

### Scenario 1: Mission Planning with Government Basemaps

```
Scenario: Plan patrol route in unfamiliar terrain

Workflow:
  1. Open ATAK
  2. Add basemap: USGS topographic map (XYZ)
  3. Add overlay: USFS trails and hazards (WMS)
  4. Add overlay: Urban infrastructure (ArcGIS feature service)
  5. Plan route: Draw line from start → objective
  6. Analyze: Terrain, obstacles, enemy threat rings
  7. Export: Route as KML to team

Result: 
  - Full situational awareness
  - Terrain-aware planning
  - Shared route on all devices
  - Coordinated movement
```

### Scenario 2: Real-Time Intelligence Overlay

```
Scenario: Display current threat picture + terrain

Workflow:
  1. DDF query (via REST API) → Get enemy positions
  2. Convert to GeoJSON
  3. Create ATAK overlay: "Threats"
  4. Add government satellite imagery (WMS)
  5. Add terrain (DEM overlay)
  
Result:
  - Operator sees: Threats + terrain + satellite
  - Can analyze: Line-of-sight, movements, strategy
  - Coordinated with: TAKpilot (AI chat control)
```

### Scenario 3: Weather & Environmental Overlays

```
Scenario: Plan operation considering environmental factors

Overlays:
  1. NOAA weather (WMS) — Real-time weather
  2. USGS flooding forecast (WMS) — High water areas
  3. NPS fire perimeter (WMS) — Active fires
  4. USDA smoke forecast (WMS) — Air quality
  
Result:
  - Full environmental picture
  - Route avoids hazards
  - Timing coordinated with weather
  - Personnel safety maximized
```

### Scenario 4: ISR Analysis with Multiple Layers

```
Scenario: Reconnaissance imagery analysis

Layers:
  1. High-res satellite (WMS) — Recent imagery
  2. Historical positions (KML network link) — Trend data
  3. Infrastructure (WFS) — Roads, buildings, utilities
  4. No-fly zones (KML overlay) — Airspace restrictions
  5. Friendly positions (NATS stream) — Real-time allies

Analysis:
  - Find optimal recon path
  - Avoid coverage (enemy radar + visual)
  - Integrate with friendly positions
  - Export as mission package
```

### Scenario 5: Disaster Response Operations

```
Scenario: Emergency response to natural disaster

Layers:
  1. Satellite damage assessment (WMS) — Building damage
  2. Road status (WFS) — Impassable roads highlighted
  3. Hospital locations (WFS) — Medical facilities
  4. Shelters (WFS) — Evacuation shelters, services
  5. No-go zones (KML) — Hazardous areas
  6. Resource positions (NATS) — Fire/EMS locations

Result:
  - Operators see full disaster picture
  - Route allocation optimized
  - Resource coordination efficient
  - Personnel safety prioritized
```

---

## Key Technical Concepts for Operators

### Basemaps vs. Overlays

**Basemap:**
- Background reference (satellite, topo, street map)
- Only one active at a time
- Fast rendering (XYZ tiles)
- Examples: Satellite, OpenStreetMap, topographic

**Overlay:**
- Data layer on top of basemap
- Multiple overlays can be active
- Can be toggled on/off
- Examples: Threats, infrastructure, weather, roads

### Raster vs. Vector Data

**Raster (image-based):**
- Satellite imagery, elevation models, weather maps
- Format: WMS, XYZ tiles, GeoTIFF
- Good for: Photography, continuous data (elevation, temperature)
- Size: Large (full coverage)
- Rendering: Fast (GPU-accelerated)

**Vector (feature-based):**
- Points (locations), lines (roads, boundaries), polygons (areas)
- Formats: GeoJSON, KML, WFS
- Good for: Discrete features (cities, roads, hazards)
- Size: Small (only features, not full coverage)
- Rendering: Fast, scalable

### Coordinate Systems

**Standard for ATAK:**
- WGS84 (GPS standard)
- Latitude/Longitude (degrees)
- Grid references: MGRS, UTM

**GIS servers use:**
- WGS84 (most common)
- Web Mercator (web mapping standard)
- Local projections (country-specific)

**ATAK handles conversion automatically** (no operator work needed)

---

## Common GIS Data Integration Workflows

### Workflow 1: Add USGS Satellite Basemap

```
1. Find USGS server: https://server.arcgisonline.com/arcgis/rest/services
2. Navigate: World_Imagery > Layers > 0 (satellite imagery layer)
3. Copy URL to ATAK basemap config
4. Restart ATAK
5. Select "USGS Satellite" as basemap
```

### Workflow 2: Add National Park Trails

```
1. Query NPS GIS server for trail features
2. Convert to KML or GeoJSON
3. Import as overlay in ATAK
4. Toggle on/off as needed
5. Download locally for offline use
```

### Workflow 3: Add Custom Feature Layer

```
1. Create GeoJSON file: points, lines, polygons
2. Host on local server or web server
3. Add WFS endpoint to ATAK
4. ATAK queries features on demand
5. Style as needed
```

### Workflow 4: Live Data Updates from Server

```
1. Set up KML NetworkLink with refresh interval
2. Server updates features (e.g., threat positions)
3. ATAK polls server every N minutes
4. New data automatically displayed
5. No operator interaction needed
```

---

## Performance Considerations

### Basemap Performance

```
XYZ tiles (fastest):
  - Rendering: <100ms (cached client-side)
  - Network: ~50KB per screen (~10-50 tiles)
  - Suitable for: Real-time COP, high-frequency updates

WMS (medium):
  - Rendering: 100-500ms (generated by server)
  - Network: 100KB-1MB per request (entire viewport)
  - Suitable for: Analysis, less frequent updates

Vector (fast for few features):
  - Rendering: <100ms (1000 features)
  - Network: 10KB-100KB per update
  - Suitable for: Specific feature layers
```

### Optimization Tips

```
✓ Use XYZ tiles for basemap (fastest)
✓ Limit WMS overlays (1-2 max for smooth performance)
✓ Cache basemaps locally (offline capability)
✓ Pre-filter vector data (don't fetch entire dataset)
✓ Use appropriate zoom levels (load only visible data)
✓ Compress overlays where possible
✓ Test on target hardware (ATAK on Android often CPU-limited)
```

---

## Data Sources Quick Reference

### United States Government

| Agency | Server | Data |
|--------|--------|------|
| **USGS** | https://server.arcgisonline.com/ | Satellite, topo, roads, geology |
| **NOAA** | https://gis.ngdc.noaa.gov/arcgis | Weather, ocean, coastal |
| **FEMA** | https://gis.fema.gov/arcgis | Hazards, flood zones, shelters |
| **USDA Forest Service** | https://apps.fs.usda.gov/arcx | Trails, forests, fire perimeters |
| **BLM** | https://gis.blm.gov/arcgis | Public lands, recreation |
| **DoD** | https://services.arcgisonline.com/Military | Military overlays |

### International / Open Source

| Source | Format | Data |
|--------|--------|------|
| **OpenStreetMap** | XYZ tiles | Streets, buildings, POI (crowdsourced) |
| **Natural Earth** | WMS, GeoJSON | Natural features, admin boundaries |
| **Copernicus** | WMS | Satellite imagery (Sentinel), land cover |
| **GEBCO** | GeoTIFF, WMS | Ocean bathymetry |

---

## Troubleshooting Common Issues

### Issue 1: Data Won't Display

**Diagnosis:**
1. Check URL is correct (test in browser)
2. Verify coordinate system (should be WGS84)
3. Check zoom level (some layers only at certain zoom)
4. Verify layer is visible (not turned off)

**Solution:**
- Copy working example XML
- Replace URL incrementally
- Test each step

### Issue 2: Slow Performance

**Diagnosis:**
1. Too many overlays active
2. Raster overlay too large
3. Vector data too dense
4. Network latency

**Solution:**
- Reduce overlays to 1-2 essential ones
- Switch to vector data (smaller)
- Increase map refresh interval
- Cache basemaps locally
- Test on device (don't assume desktop performance)

### Issue 3: Data Coordinates Wrong

**Diagnosis:**
1. Server using different coordinate system
2. Manual coordinate entry error
3. Projection mismatch

**Solution:**
- ATAK auto-converts WGS84 ↔ MGRS ↔ UTM
- Verify source data is WGS84
- Use "Copy coordinates" feature
- Test with known location

### Issue 4: Cannot Find Data Source

**Diagnosis:**
1. Server URL changed
2. GIS server offline/deprecated
3. Firewall blocking access
4. Authentication required

**Solution:**
- Try alternative source (USGS instead of agency-specific)
- Check ArcGIS Online for updated servers
- Test network connectivity
- Request authentication credentials

---

## Best Practices

### ✓ Do This

- **Verify data before mission** — Test all overlays
- **Customize for mission** — Keep only relevant layers
- **Cache offline** — Download basemaps for GPS-denied areas
- **Document sources** — Know which data came from where
- **Version control** — Keep backup XML configs
- **Train team** — Everyone knows layer meaning
- **Communicate layers** — Tell team what overlay means

### ✗ Don't Do This

- **Use stale data** — Verify freshness before op
- **Mix coordinate systems** — Keep everything WGS84
- **Add too many overlays** — More than 3-4 causes slowdown
- **Assume online availability** — Cache for offline
- **Skip testing** — Test on actual hardware + connectivity
- **Use non-standard formats** — Stick to KML, GeoJSON, WMS
- **Leave sensitive data unencrypted** — Secure downloaded overlays

---

## Related Products & Resources

### ATAK Companions

| Tool | Purpose | Integration |
|------|---------|-------------|
| **WinTAK** | Desktop (Windows) version | Same overlays as Android |
| **CloudTAK** | Web-based C2 | Browser-based alternative |
| **TAKServer** | Centralized position sync | Data source for ATAK |
| **TAKpilot** | AI chat control | Can add GIS layers via chat |

### GIS Tools for Pre-processing

```
Convert data → Overlay-ready format:
  - QGIS: Convert shapefiles → KML
  - ArcGIS Pro: Export to GeoJSON
  - gdal_translate: Command-line conversion
  - Leaflet.draw: Web interface for custom features
```

---

## Learning Path

### Beginner (1-2 weeks)
1. Download ATAK/WinTAK
2. Add one basemap (XYZ tiles)
3. Add one overlay (USGS topographic map)
4. Practice placing markers + routes
5. Export KML file

### Intermediate (2-4 weeks)
1. Add multiple overlays (satellite + roads + threats)
2. Learn WMS configuration
3. Create custom XML configs
4. Set up KML network links (live data)
5. Integrate with DDF intelligence

### Advanced (1-2 months)
1. Master ArcGIS server navigation
2. Create custom feature services
3. Build automated data pipelines (script → GeoJSON → ATAK)
4. Optimize performance for large deployments
5. Deploy on Team TAK infrastructure

### Expert (Ongoing)
1. Contribute to ATAK-CIV development
2. Build custom ATAK plugins (for advanced features)
3. Deploy for large-scale operations
4. Mentor others

---

## Community & Support

**Author:** Joseph Elfelt (MappingSupport.com)  
**Contact:** https://m.ai6yr.org/@mappingsupport (Mastodon)  
**Website:** https://mappingsupport.com/  

**ATAK Community:**
- TAK.gov (official ATAK documentation)
- ATAK-CIV (civilian fork)
- GitHub: TAK project repositories

---

## Summary: Value of GIS Integration with ATAK

| Benefit | Impact | Result |
|---------|--------|--------|
| **Situational awareness** | Operators see full geospatial picture | Better decisions |
| **Terrain analysis** | Elevation, slope, obstacles visible | Optimal routes |
| **Intelligence overlay** | Multiple data sources combined | Comprehensive threat picture |
| **Interoperability** | Standard formats (KML, GeoJSON) | Share with non-ATAK systems |
| **Scalability** | Can add unlimited overlays (one at a time) | Supports all mission types |
| **Offline capability** | Cache data locally | Works GPS-denied |
| **No plugins needed** | Built-in ATAK features | Simple deployment |

---

**Captured:** 2026-07-09  
**Status:** Active, maintained knowledge base  
**Relevance:** Essential for ATAK/WinTAK operators, ISR teams, mission planners  
**Learning curve:** Low to medium (requires GIS conceptual understanding)  
**Resource quality:** High (step-by-step, no fluff)  
**Integration priority:** High (foundational for tactical COP)  
**Key innovation:** Practical focus + government data source inventory  

---

## Related Resources in Convergence

- **TAK Server** — Backend position sync + data server
- **TAKpilot** — AI chat control (can manipulate GIS layers)
- **CesiumJS** — 3D geospatial visualization (complements ATAK 2D)
- **DDF** — Intelligence catalog (data source for overlays)
- **NATS** — Real-time position streaming (feed to ATAK)
- **Mumble** — Comms coordination (synced with tactical ops)

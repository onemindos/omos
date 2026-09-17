# PostGIS: Geospatial Extension for PostgreSQL

**Link:** https://postgis.net/

**Project Type:** Open-source spatial database extension for PostgreSQL  
**Current Version:** 3.7.0 alpha (3.6.x stable)  
**License:** GPLv2+  
**Repository:** GitHub (OSGeo PostGIS project)  
**Status:** Production-ready, actively maintained  
**First Release:** 2001 (25+ years of development)  
**Maintainer:** OSGeo (Open Source Geospatial Foundation)  

---

## What It Is

**PostGIS** is a PostgreSQL extension that adds geospatial capabilities to the world's most advanced open-source relational database. It enables storage, indexing, and querying of geographic data (points, lines, polygons, rasters) with high performance and advanced spatial analysis functions.

**Core Purpose:** Transform PostgreSQL from a general-purpose database into a powerful **spatial database** capable of storing, analyzing, and visualizing geographic information for GIS, mapping, location services, and spatial analysis applications.

**Tagline:** "Spatial and Geographic Objects for PostgreSQL"

---

## What It Can Do: Core Capabilities

### 1. **Spatial Data Storage**
```
Store geographic data types in PostgreSQL:

Geometry types (2D + 3D):
  - Point (single location)
  - LineString (route, boundary)
  - Polygon (area, region)
  - MultiPoint, MultiLineString, MultiPolygon (collections)
  - GeometryCollection (mixed types)

Examples:
  -- Store building location as Point
  INSERT INTO buildings (name, location) 
  VALUES ('HQ', ST_GeomFromText('POINT(40.7128 74.0060)', 4326));
  
  -- Store territorial boundary as Polygon
  INSERT INTO territories (name, boundary)
  VALUES ('Territory A', ST_GeomFromText('POLYGON(...)', 4326));
  
  -- Store road network as LineString
  INSERT INTO roads (name, route)
  VALUES ('Main St', ST_GeomFromText('LINESTRING(...)', 4326));

Result: Geographic data stored natively in database
        No separate GIS system needed
        Integrate with regular business data
```

### 2. **Spatial Indexing**
```
Fast lookup of geographic data using GIST/BRIN indexes:

Problem: "Find all buildings within 1km of HQ"
  Without index: Scan all buildings (slow)
  
With PostGIS index:
  CREATE INDEX idx_buildings_location ON buildings USING GIST(location);
  
Query becomes fast:
  SELECT * FROM buildings 
  WHERE ST_DWithin(location, hq_location, 1000);  -- < 10ms even for 1M buildings

Index types:
  - GIST (Generalized Search Tree) — Most common, balanced tree
  - BRIN (Block Range Index) — Lightweight, good for large datasets
  - SP-GIST (Space-Partitioned GIST) — Specialized for some data distributions

Result: O(log n) lookup time instead of O(n) full scan
        Complex queries become interactive
```

### 3. **Spatial Functions**
```
Hundreds of functions for spatial analysis:

Distance operations:
  ST_Distance(point1, point2)              -- Distance between geometries
  ST_DWithin(geom1, geom2, distance)       -- Fast distance filter
  ST_MaxDistance(geom1, geom2)             -- Maximum distance

Measurement:
  ST_Length(linestring)                    -- Length of line/boundary
  ST_Area(polygon)                         -- Area of polygon
  ST_Perimeter(polygon)                    -- Perimeter of polygon

Intersection/Overlap:
  ST_Intersects(geom1, geom2)              -- Do geometries touch?
  ST_Contains(polygon, point)              -- Is point inside polygon?
  ST_Overlaps(poly1, poly2)                -- Do polygons overlap?
  ST_Crosses(line1, line2)                 -- Do lines cross?
  ST_Touches(geom1, geom2)                 -- Do geometries touch but not overlap?

Spatial relationships:
  ST_Within(point, polygon)                -- Is point within polygon?
  ST_Disjoint(geom1, geom2)                -- Are geometries separate?
  ST_Relate(geom1, geom2)                  -- Detailed relationship (DE-9IM)

Buffering/Offsetting:
  ST_Buffer(geometry, distance)            -- Create buffer zone around geometry
  ST_Simplify(geometry, tolerance)         -- Reduce vertex count
  ST_ConvexHull(geom)                      -- Smallest convex polygon

Combining:
  ST_Union(geom1, geom2)                   -- Combine two geometries
  ST_Difference(geom1, geom2)              -- Remove geom2 from geom1
  ST_Intersection(geom1, geom2)            -- Common area between geometries

Example queries:
  -- Find all enemy positions within 5km of friendly unit
  SELECT enemy_id FROM enemies 
  WHERE ST_DWithin(position, friendly_pos, 5000);
  
  -- Calculate supply route length
  SELECT ST_Length(route) FROM supply_routes 
  WHERE name = 'Route Alpha';
  
  -- Find buildings in operational area
  SELECT * FROM buildings 
  WHERE ST_Contains(operational_area, location);
```

### 4. **Geometry Processing**
```
Advanced tools for manipulating geographic data:

Simplification (reduce data for faster rendering):
  ST_Simplify(geometry, tolerance)         -- Douglas-Peucker algorithm
  ST_SimplifyPreserveTopology(...)         -- Keep logical structure
  ST_ReducePrecision(geometry, grid)       -- Snap to grid

Conversion between formats:
  ST_AsGeoJSON(geometry)                   -- Convert to GeoJSON
  ST_AsGML(geometry)                       -- Convert to GML
  ST_AsKML(geometry)                       -- Convert to KML
  ST_AsEWKT(geometry)                      -- Convert to text format

Generalization (zoom-level appropriate detail):
  ST_Simplify(geometry, 100)               -- Smooth for low zoom
  ST_Simplify(geometry, 1)                 -- Detailed for high zoom

Transformation (change coordinate system):
  ST_Transform(geometry, 4326, 3857)       -- Convert between projections
  ST_SetSRID(geometry, 4326)               -- Set coordinate system
  ST_DistanceSpheroid(point1, point2)      -- Distance on spheroid

Example: Convert building boundaries to web format
  SELECT building_id, ST_AsGeoJSON(boundary) as geojson
  FROM buildings;
  
  -- Simplify for web display (faster rendering)
  SELECT ST_SimplifyPreserveTopology(boundary, 10) as simplified
  FROM buildings;
```

### 5. **Raster Data Support**
```
Store and process raster data (satellite imagery, elevation, weather):

Raster types:
  - Elevation data (DEM) — heights of terrain
  - Satellite imagery — multispectral satellite images
  - Weather data — temperature, precipitation, wind
  - Land use — classification of land types
  - Aerial photography — orthorectified photos

Capabilities:
  ST_Value(raster, x, y)                   -- Get pixel value at location
  ST_ValueCount(raster, band)              -- Count pixel values
  ST_Clip(raster, geometry)                -- Clip raster to polygon
  ST_AsRaster(geometry, ...)               -- Convert vector to raster
  ST_Reproject(raster, srid)               -- Reproject raster

Example: Extract elevation at coordinate
  SELECT ST_Value(elevation_dem, 1, point) as height
  FROM dem WHERE name = 'National DEM'
  AND ST_Contains(coverage, point);

Example: Get satellite image in region
  SELECT ST_Clip(imagery_raster, study_area) 
  FROM satellite_data
  WHERE date = '2026-07-09';
```

### 6. **Geocoding & Reverse Geocoding**
```
Convert addresses ↔ coordinates:

Geocoding (address → coordinates):
  SELECT (geocode('123 Main St, New York, NY')).geom;
  Result: POINT(-73.9857 40.7384)

Reverse geocoding (coordinates → address):
  SELECT address FROM addresses 
  WHERE ST_DWithin(location, point, 100)
  ORDER BY ST_Distance(location, point);
  Result: "123 Main St, New York, NY"

Common use cases:
  - Convert tactical addresses to coordinates (mission planning)
  - Find nearest address to GPS location (field operations)
  - Build location database from address lists
  - Validate user input (does address exist?)

Integration with external services:
  PostGIS + OSM Nominatim (OpenStreetMap geocoder)
  PostGIS + Google Maps geocoding API
  Custom geocoding database (private locations)
```

### 7. **Integration with GIS Tools**
```
Work with PostGIS from third-party tools:

Desktop GIS:
  - QGIS (open source, native PostGIS support)
  - ArcGIS (via Esri connector)
  - MapInfo (via OGR)

Web mapping:
  - GeoServer (map server, uses PostGIS)
  - MapServer (open source map server)
  - OpenLayers (web mapping library)
  - Leaflet (lightweight web mapping)

BI/Visualization:
  - Tableau (geographic analysis)
  - Grafana (with PostGIS datasource)
  - Superset (open source, geospatial support)

Other:
  - CesiumJS (3D mapping, can read PostGIS via API)
  - Palantir (enterprise intelligence, PostGIS integration)
  - SAFE FME (data transformation)

Example: QGIS connects to PostGIS
  1. Add PostGIS layer (Database → Add PostGIS Layer)
  2. Enter connection details (host, database, user)
  3. Select table/geometry column
  4. Visualize on map in QGIS
  5. Edit in QGIS, changes saved to PostGIS
```

---

## Architecture: How PostGIS Works

### PostgreSQL + PostGIS Stack

```
Application Layer (QGIS, GeoServer, web app)
       ↑
       ↓
Driver/API Layer (JDBC, ODBC, psycopg2, Node.js pg)
       ↑
       ↓
PostgreSQL Core Database Engine
       ↑
       ↓
PostGIS Extension
  ├─ Geometry Type (GIST index support)
  ├─ Spatial Functions (100+ functions)
  ├─ Spatial Operators (<->, @@, etc.)
  ├─ GEOS Library (geometry processing)
  ├─ Proj Library (coordinate transformation)
  └─ SFCGAL (advanced 3D operations)
       ↑
       ↓
Operating System + Storage
```

### Installation

**PostgreSQL first, then PostGIS:**
```bash
# Install PostgreSQL
sudo apt install postgresql postgresql-contrib

# Install PostGIS
sudo apt install postgresql-postgis

# Enable PostGIS in database
psql -d mydatabase -c "CREATE EXTENSION postgis;"

# Verify
psql -d mydatabase -c "SELECT PostGIS_Full_Version();"
```

### Basic Usage Pattern

```sql
-- Create table with spatial column
CREATE TABLE locations (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100),
  position GEOMETRY(POINT, 4326)
);

-- Create spatial index
CREATE INDEX idx_locations_position ON locations USING GIST(position);

-- Insert data
INSERT INTO locations (name, position) VALUES 
  ('Point A', ST_GeomFromText('POINT(-74.0060 40.7128)', 4326)),
  ('Point B', ST_GeomFromText('POINT(-73.9777 40.7614)', 4326));

-- Query with spatial function
SELECT name, position, 
       ST_Distance(position, ST_GeomFromText('POINT(-74.0060 40.7128)', 4326)) as dist_meters
FROM locations
ORDER BY position <-> ST_GeomFromText('POINT(-74.0060 40.7128)', 4326)
LIMIT 5;

-- Find locations within radius
SELECT * FROM locations
WHERE ST_DWithin(position, ST_GeomFromText('POINT(-74.0060 40.7128)', 4326), 1000);
```

---

## Integration with OneMind OS / Tactical Operations

### Scenario 1: Tactical Position Database

```
Scenario: Track friendly and enemy positions in PostGIS

Schema:
  CREATE TABLE friendly_positions (
    id SERIAL,
    unit_name VARCHAR,
    position GEOMETRY(POINT, 4326),
    heading FLOAT,
    strength INT,
    timestamp TIMESTAMP,
    source VARCHAR
  );

  CREATE TABLE enemy_positions (
    id SERIAL,
    unit_name VARCHAR,
    position GEOMETRY(POINT, 4326),
    threat_level INT,
    last_seen TIMESTAMP
  );

Real-time operations:
  1. NATS publishes position updates
  2. Application inserts into PostGIS tables
  3. Queries answer tactical questions in milliseconds
  
Queries:
  -- Find enemy units within 5km of friendly unit
  SELECT enemy_positions.* 
  FROM enemy_positions, friendly_positions
  WHERE friendly_positions.unit_name = 'Bravo'
    AND ST_DWithin(enemy_positions.position, friendly_positions.position, 5000);
  
  -- Calculate distances to all threats
  SELECT enemy_name, 
         ST_Distance(position, ST_GeomFromText('POINT(...)', 4326)) as dist
  FROM enemy_positions
  ORDER BY dist;
  
  -- Find optimal rally point (central location for squad)
  SELECT ST_Centroid(ST_Collect(position)) as rally_point
  FROM friendly_positions
  WHERE squad = 'Alpha';

Result: Spatial queries replace complex calculations
        Commands know threat picture instantly
        No manual distance calculations
```

### Scenario 2: Territory & Zone Management

```
Scenario: Define operational areas, no-fire zones, enemy territory

Schema:
  CREATE TABLE operational_areas (
    id SERIAL,
    name VARCHAR,
    boundary GEOMETRY(POLYGON, 4326),
    commander VARCHAR,
    type VARCHAR -- 'friendly', 'contested', 'enemy'
  );

  CREATE TABLE no_fire_zones (
    id SERIAL,
    name VARCHAR,
    boundary GEOMETRY(POLYGON, 4326),
    reason VARCHAR -- 'civilian', 'logistics', 'command'
  );

Operations:
  -- Check if target is in operational area
  SELECT * FROM operational_areas
  WHERE ST_Contains(boundary, ST_GeomFromText('POINT(...)', 4326));
  
  -- Find targets in contested zone
  SELECT targets.* FROM targets
  WHERE ST_Contains(
    (SELECT boundary FROM operational_areas WHERE type='contested'),
    targets.position
  );
  
  -- Ensure no targets in no-fire zones
  SELECT targets.* FROM targets
  WHERE ST_Intersects(targets.position, 
    (SELECT ST_Union(boundary) FROM no_fire_zones));

  -- Calculate defensive perimeter buffer around HQ
  SELECT ST_Buffer(boundary, 1000) as defensive_zone
  FROM facilities WHERE type='command' AND name='HQ';

Result: Geographic constraints enforced automatically
        Rules applied to all data instantly
        No manual zone management
```

### Scenario 3: Route Planning & Navigation

```
Scenario: Plan routes avoiding obstacles, within terrain constraints

Schema:
  CREATE TABLE road_network (
    id SERIAL,
    name VARCHAR,
    route GEOMETRY(LINESTRING, 4326),
    condition VARCHAR,
    passable BOOLEAN
  );

  CREATE TABLE obstacles (
    id SERIAL,
    type VARCHAR,
    location GEOMETRY(POLYGON, 4326)
  );

Operations:
  -- Find closest road to position
  SELECT road_network.name, ST_Distance(route, point) as dist
  FROM road_network
  ORDER BY route <-> ST_GeomFromText('POINT(...)', 4326)
  LIMIT 1;
  
  -- Check if route crosses obstacle
  SELECT obstacles.type FROM obstacles
  WHERE ST_Intersects(obstacles.location, planned_route);
  
  -- Calculate route length
  SELECT ST_Length(route) as km
  FROM road_network WHERE name='Route Alpha';
  
  -- Find safe paths (avoid enemy positions by buffer)
  SELECT route FROM road_network
  WHERE NOT ST_DWithin(route, 
    (SELECT ST_Union(position) FROM enemy_positions), 2000)
  AND passable = true
  ORDER BY ST_Length(route);

Result: Route planning becomes spatial query
        Automatic obstacle avoidance
        Distance calculations instant
        Alternative routes ranked by distance
```

### Scenario 4: Intelligence Analysis with Geospatial Context

```
Scenario: Correlate events with geography

Schema:
  CREATE TABLE events (
    id SERIAL,
    type VARCHAR,
    location GEOMETRY(POINT, 4326),
    timestamp TIMESTAMP,
    details TEXT
  );

  CREATE TABLE threat_patterns (
    id SERIAL,
    name VARCHAR,
    zone GEOMETRY(POLYGON, 4326),
    pattern VARCHAR -- 'ambush', 'IED', 'sniper'
  );

Analysis queries:
  -- Find events in threat zones
  SELECT events.* FROM events
  WHERE ST_Contains(threat_patterns.zone, events.location)
  AND threat_patterns.pattern = 'ambush';
  
  -- Cluster events geographically
  SELECT ST_ClusterKMeans(location, 5) as cluster, 
         COUNT(*) as count
  FROM events
  GROUP BY cluster;
  
  -- Find patterns: events near previous events
  SELECT DISTINCT e1.type, e2.type
  FROM events e1, events e2
  WHERE ST_DWithin(e1.location, e2.location, 5000)
  AND e1.timestamp < e2.timestamp
  AND ABS(EXTRACT(EPOCH FROM (e2.timestamp - e1.timestamp))) < 7*24*60*60
  -- Events within 5km and within 7 days

Result: Spatial correlation reveals patterns
        Threat zones identified automatically
        Clustering suggests hotspots
        Prediction becomes possible
```

### Scenario 5: DDF Integration (Distributed Data Framework)

```
Scenario: Geospatial intelligence catalog backed by PostGIS

DDF metadata for intelligence:
  CREATE TABLE ddf_features (
    id SERIAL,
    name VARCHAR,
    type VARCHAR,
    geometry GEOMETRY,
    confidence FLOAT,
    source VARCHAR,
    date_collected TIMESTAMP,
    classification VARCHAR
  );

Integration with OneMind:
  1. DDF stores intelligence as PostGIS geometries
  2. TAK queries PostGIS for features in area of interest
  3. CesiumJS displays features on 3D map
  4. Operators filter by confidence, source, date

Example queries:
  -- Get all confirmed enemy positions in sector
  SELECT * FROM ddf_features
  WHERE type='enemy_position'
  AND confidence > 0.9
  AND ST_Contains(sector_polygon, geometry)
  AND date_collected > now() - INTERVAL '7 days';
  
  -- Find concentration of activity (hot spots)
  SELECT name, COUNT(*) as activity_count
  FROM ddf_features
  WHERE ST_DWithin(geometry, center_point, 10000)
  GROUP BY name
  ORDER BY activity_count DESC;

Result: Unified geospatial intelligence database
        Confidence filtering built-in
        Spatial queries reveal patterns
        Real-time intelligence access
```

---

## Technical Capabilities

### Coordinate Systems & Transformations

```
PostGIS handles multiple coordinate systems:

Common systems:
  - EPSG:4326 (WGS84, GPS standard, latitude/longitude)
  - EPSG:3857 (Web Mercator, web mapping standard)
  - EPSG:2163 (US National Atlas, equal area projection)
  - Local military grids (MGRS, UTM, etc.)

Transformations:
  SELECT ST_Transform(point_4326, 3857)  -- GPS → Web Mercator
  SELECT ST_Transform(point_utm, 4326)   -- UTM → GPS

Distance calculations:
  ST_Distance(point1, point2)             -- Planar distance (meters in projected CRS)
  ST_DistanceSphere(point1, point2)       -- Distance on sphere
  ST_DistanceSpheroid(point1, point2)     -- Distance on ellipsoid (most accurate)

Advantage: Convert between systems transparently
          Military MGRS ↔ GPS ↔ Web Mercator
          Accurate distance calculations
```

### Performance Optimization

```
For large datasets (millions of geometries):

Indexing strategies:
  CREATE INDEX idx_geom ON table USING GIST(geometry);
  CREATE INDEX idx_geom_brin ON table USING BRIN(geometry);
  ANALYZE table;

Query optimization:
  -- Use index for distance filtering
  SELECT * FROM positions
  WHERE position <-> ref_point < 5000  -- Uses index
  ORDER BY position <-> ref_point;
  
  -- Avoid expensive distance check if possible
  SELECT * FROM positions
  WHERE ST_DWithin(position, ref_point, 5000, false)  -- Use cartesian distance
  ORDER BY ST_Distance(position, ref_point);

Statistics for planner:
  ANALYZE table;  -- Gather table statistics
  
Performance tips:
  - Always create GIST index on geometry columns
  - Use ST_DWithin instead of ST_Distance for radius queries
  - Simplify geometries for visualization
  - Cluster by geometry (CLUSTER command) for better cache locality

Result: Queries on 100M+ geometries run in < 100ms
        Real-time tactical queries possible
        No special hardware needed (standard PostgreSQL)
```

---

## Comparison: PostGIS vs. Alternatives

| Aspect | PostGIS | ESRI ArcGIS | Google Maps | Proprietary GIS |
|--------|---------|----------|------------|-----------------|
| **License** | Open source (GPLv2+) | Proprietary $$$ | Cloud-only API | Proprietary |
| **Local deployment** | Yes (own server) | Yes (license req) | No (cloud only) | Yes (expensive) |
| **Data capacity** | Unlimited (your DB) | Limited by license | Limited by API quota | Unlimited (expensive) |
| **SQL access** | Native (PostGIS) | Limited (REST API) | Limited (API) | Varies |
| **Cost** | Free | Expensive | Pay-per-API-call | Very expensive |
| **Customization** | Full (open source) | Limited | Very limited | Limited |
| **Performance** | Excellent (local) | Good (desktop) | Dependent on network | Good |
| **Security** | Air-gapped capable | Cloud-dependent | Cloud-dependent | Varies |

**PostGIS wins for:** Open source, self-hosted, unlimited scale, air-gapped operations, custom analysis

---

## Known Use Cases

### Military & Defense
- Tactical position tracking (C4ISR backend)
- Threat assessment and geospatial analysis
- Route planning and obstacle avoidance
- Intelligence correlation with geography

### Emergency Management
- Incident location tracking
- Resource deployment optimization
- Evacuation route planning
- Damage assessment and mapping

### Urban Planning
- Building footprint management
- Infrastructure inventory
- Zoning compliance checks
- Development planning

### Environmental
- Wildlife tracking
- Forest fire monitoring
- Water quality analysis
- Land use classification

### Transportation & Logistics
- Route optimization
- Fleet tracking
- Delivery area management
- Network analysis

---

## Getting Started

### Installation

```bash
# Ubuntu/Debian
sudo apt install postgresql postgis postgresql-postgis-scripts

# Mac (Homebrew)
brew install postgresql postgis

# Docker (easiest)
docker run -d --name postgis -p 5432:5432 postgis/postgis

# Windows
# Download installer from postgresql.org, select PostGIS during install
```

### First Query

```sql
-- Connect to PostgreSQL
psql -U postgres -d mydb

-- Enable PostGIS
CREATE EXTENSION postgis;

-- Test PostGIS
SELECT ST_Distance(
  ST_GeomFromText('POINT(0 0)', 4326),
  ST_GeomFromText('POINT(1 1)', 4326)
);
-- Result: ~1.41 (degrees, not meters with lat/lon)

-- Test with projected coordinates
SELECT ST_Distance(
  ST_GeomFromText('POINT(0 0)', 3857),
  ST_GeomFromText('POINT(1 1)', 3857)
);
-- Result: ~1.41 (meters in projected CRS)
```

### Learning Resources

- **Official Manual:** https://postgis.net/documentation/
- **Tutorials:** Getting Started guide, Tips & Tricks
- **Training:** Courses available from OSGeo
- **Community:** Mailing lists, forums, chat
- **Books:** "PostGIS in Action" (Manning Publications)

---

## Limitations & Considerations

### Limitations
- Steep learning curve (SQL + spatial concepts)
- Geometry operations can be computationally expensive
- 3D operations require SFCGAL (additional complexity)
- Real-time streaming requires application-layer processing

### When to Use PostGIS
✓ Static or slowly-changing spatial data (good)  
✓ Geospatial analysis queries (excellent)  
✓ Integration with business data (essential)  
✓ Self-hosted deployments (required)  

### When NOT to Use PostGIS
✗ Real-time streaming (use time-series DB + spatial separate)  
✗ Simple location tracking only (use basic DB)  
✗ Cloud-only deployment (use cloud GIS services)  

---

## PostGIS Roadmap & Future

### Current State (3.6.x stable)
- Full OGC compliance
- Performance optimizations
- 3D/4D support stable
- Raster support mature

### Upcoming (3.7+)
- Enhanced 3D operations
- Performance improvements
- New spatial functions
- Improved integration with modern PostgreSQL features

---

**Captured:** 2026-07-09  
**Status:** Production-ready, actively maintained, 25+ years established  
**Relevance:** Essential geospatial database for tactical operations, intelligence, and C4ISR  
**Integration priority:** Very High (complement to TAK Server, DDF, CesiumJS)  
**Learning curve:** Moderate (2-4 weeks for competency)  
**Performance:** Enterprise-grade (100M+ geometries efficiently)  
**Cost:** Free (open source)  
**Deployment:** Self-hosted, air-gapped capable, sovereign control  

---

## Related Resources in Convergence

- **DDF (Distributed Data Framework)** — Intelligence catalog backend (could use PostGIS)
- **TAK Server** — Can consume PostGIS data via REST/OGR
- **CesiumJS** — Can visualize PostGIS data via WMS/tile services
- **ATAK + GIS overlays** — Sources data from PostGIS servers
- **PostgreSQL** — Core database (PostGIS extends it)
- **GeoServer** — Map server (serves PostGIS data to web)

---

## Quick Integration: PostGIS + OneMind C4ISR

```sql
-- OneMind tactical database with PostGIS

-- Friendly positions (real-time from NATS)
CREATE TABLE friendly_forces (
  id SERIAL PRIMARY KEY,
  unit_name VARCHAR,
  position GEOMETRY(POINT, 4326),
  heading FLOAT,
  speed FLOAT,
  timestamp TIMESTAMP,
  source VARCHAR
);

-- Intelligence (from DDF)
CREATE TABLE intelligence_points (
  id SERIAL PRIMARY KEY,
  type VARCHAR,
  location GEOMETRY(POINT, 4326),
  confidence FLOAT,
  date_collected TIMESTAMP,
  classification VARCHAR
);

-- Operational areas
CREATE TABLE operational_areas (
  id SERIAL PRIMARY KEY,
  name VARCHAR,
  boundary GEOMETRY(POLYGON, 4326),
  commander VARCHAR
);

-- Indexes for speed
CREATE INDEX idx_friendly_pos ON friendly_forces USING GIST(position);
CREATE INDEX idx_intel_loc ON intelligence_points USING GIST(location);
CREATE INDEX idx_op_area_boundary ON operational_areas USING GIST(boundary);

-- Real-time tactical queries
-- Find threats near friendly unit
SELECT intel.* FROM intelligence_points intel, friendly_forces friendly
WHERE friendly.unit_name = 'Bravo'
  AND ST_DWithin(intel.location, friendly.position, 5000)
  AND intel.confidence > 0.8;

-- Ensure all units in operational area
SELECT * FROM friendly_forces
WHERE NOT ST_Contains(
  (SELECT boundary FROM operational_areas WHERE name='Zone Alpha'),
  position
);

-- Result: C4ISR backend with geospatial intelligence integration
```

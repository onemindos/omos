# Leica BLK — Reference Platform

## What It Is
Professional LIDAR laser scanning hardware and software ecosystem by Leica Geosystems (Hexagon). Captures high-resolution 3D point clouds and HDR imagery of physical environments. Used in surveying, construction, public safety, and military site mapping.

## Hardware Form Factors
- **BLK360** — compact terrestrial scanner, indoor/outdoor, 360° scan in minutes
- **BLK2GO** — handheld mobile scanner, walks through space while scanning
- **BLK2FLY** — autonomous flying scanner (drone-mounted), scans from air
- **BLK ARC** — robot-mounted scanner (Boston Dynamics Spot compatible)
- **RTC360** — high-accuracy tripod scanner for precise site mapping

## What It Does
- Captures millions of measured points per second as a 3D point cloud (LAS/LAZ/E57 format)
- Registers multiple scans into a unified 3D model
- Outputs: point clouds, orthophotos, floor plans, 3D meshes, Cesium 3D Tiles
- Software: Leica Cyclone (processing), Leica Register (alignment), Hexagon GeoCloud (storage)

## Data Formats
- Native: LGSx (Leica proprietary)
- Export: E57, LAS, LAZ (industry standard point cloud formats)
- Photogrammetry outputs: OBJ, FBX, glTF, Cesium 3D Tiles (via RealityCapture / Cyclone 3DR)

## Why We're Building Our Own
- Leica hardware is $15k–$100k per unit — we need a pipeline that works with any scanner
- Their software (Cyclone, GeoCloud) is proprietary and cloud-locked
- We build the processing pipeline once and it works with any LAS/LAZ/E57 source
- Our pipeline outputs directly to our geo namespace tile server — no third-party cloud

---

## OneMind Survey Layer — `omos-survey`

### What We Build
An open-source LIDAR processing pipeline that converts any point cloud (from any scanner) into Cesium 3D Tiles served from our geo namespace.

### Architecture
```
Field scanner (any device outputting LAS/LAZ/E57)
  → upload to MinIO (geo-tiles bucket) via field tablet / drone link
      → NATS trigger: evt.survey.upload.{scan_id}
          → omos-survey service picks up
              → PDAL pipeline:
                  filter noise → classify ground/building/vegetation
                  → colorize from HDR imagery
                  → normalize intensity
              → py3dtiles / entwine → Cesium 3D Tiles (pnts format)
              → output to geo namespace tile server
                  → NATS publish: evt.survey.ready.{scan_id}
                      → CloudTAK: deck.gl PointCloudLayer renders it
                      → ATAK: tileset.json URL pushed as Data Package
                      → TAK Server: distribute to all clients
```

### NATS Subjects
```
evt.survey.upload.{scan_id}     → new scan uploaded, trigger processing
evt.survey.processing.{scan_id} → processing in progress (progress %)
evt.survey.ready.{scan_id}      → tileset.json URL available
evt.survey.error.{scan_id}      → processing failed
```

### Output
```
http://geo.onemindos.dev/3dtiles/{scan_id}/tileset.json
```
Both ATAK and CloudTAK consume this URL.

### Tech Stack
- Python 3.12 with PDAL (Point Data Abstraction Library)
- py3dtiles or entwine for 3D Tiles conversion
- MinIO for raw scan storage (geo-tiles bucket)
- NATS for event-driven processing pipeline
- Martin or nginx for serving the output tiles

### Use Cases
- BLK2FLY drone scan of a facility → 3D point cloud on all TAK clients in minutes
- Ground team with BLK2GO scans a building interior → live site map for operators
- Any future scanner (DJI Zenmuse L2, Velodyne) — pipeline is format-agnostic

### Connection to geo namespace
- Upload endpoint: existing MinIO in geo namespace
- Processing: omos-survey pod with PDAL + GPU support optional
- Serving: existing Martin or TileServer GL for output tiles
- STAC catalog: scan metadata registered in existing STAC API for search by area/date

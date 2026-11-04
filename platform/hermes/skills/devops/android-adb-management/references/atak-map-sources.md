# ATAK Map Sources — Complete Catalog
> Confirmed loaded on S24 Ultra (ATAK-CIV 5.7.0.10), July 2026. 32 sources total.
> Files live at `/sdcard/atak/imagery/mobile/mapsources/` on device.
> Source XMLs at `/opt/data/atak_maps/*.xml` on Legacy pod.

---

## Google (4 sources)

| File | Name | Zoom | Tile | URL lyrs param |
|------|------|------|------|----------------|
| `google_satellite.xml` | Google Satellite | 0–21 | png | `lyrs=s` |
| `google_hybrid.xml` | Google Hybrid | 0–21 | png | `lyrs=y` |
| `google_roads.xml` | Google Roads | 0–21 | png | `lyrs=m` |
| `google_terrain.xml` | Google Terrain | 0–15 | png | `lyrs=p` |
| `google_terrain_only.xml` | Google Terrain Only | 0–15 | png | `lyrs=t` |

Base URL: `https://mt1.google.com/vt/lyrs=<X>&x={$x}&y={$y}&z={$z}`

**Notes:**
- `lyrs=s` = satellite only
- `lyrs=y` = satellite + labels (hybrid)
- `lyrs=m` = roads/streets
- `lyrs=p` = terrain + roads
- `lyrs=t` = terrain only (no roads)
- Google Maps ToS technically limits non-Google use; use ESRI as backup if restricted

---

## ESRI (8 sources)

All use URL pattern: `https://server.arcgisonline.com/ArcGIS/rest/services/<SERVICE>/MapServer/tile/{$z}/{$y}/{$x}`
**Note:** ESRI uses `{$z}/{$y}/{$x}` (row before column), opposite of standard XYZ.

| File | Name | Zoom | Tile | Service |
|------|------|------|------|---------|
| `esri_world.xml` | ESRI World Imagery | 0–19 | jpg | `World_Imagery` |
| `esri_topo.xml` | ESRI Topo | 0–19 | jpg | `World_Topo_Map` |
| `esri_street.xml` | ESRI Street Map | 0–19 | jpg | `World_Street_Map` |
| `esri_natgeo.xml` | ESRI NatGeo World Map | 0–16 | jpg | `NatGeo_World_Map` |
| `esri_ocean.xml` | ESRI Ocean Basemap | 0–16 | jpg | `Ocean/World_Ocean_Base` |
| `esri_shaded_relief.xml` | ESRI Shaded Relief | 0–13 | jpg | `World_Shaded_Relief` |
| `esri_gray_canvas.xml` | ESRI Dark Gray Canvas | 0–16 | png | `Canvas/World_Dark_Gray_Base` |
| `esri_satellite_labels.xml` | ESRI Satellite + Labels | 0–19 | png | `Reference/World_Boundaries_and_Places` |

**Notes:**
- ESRI tiles are free for non-commercial use (ArcGIS Online Terms)
- `esri_satellite_labels.xml` is an **overlay** (transparent background `#00000000`) — layer on top of satellite
- `esri_gray_canvas.xml` is useful as a dark tactical background with minimal clutter

---

## Topo / Terrain (5 sources)

| File | Name | Zoom | URL |
|------|------|------|-----|
| `osm_topo.xml` | OpenTopoMap | 0–17 | `https://a.tile.opentopomap.org/{$z}/{$x}/{$y}.png` |
| `stamen_terrain.xml` | Stadia Stamen Terrain | 0–18 | `https://tiles.stadiamaps.com/tiles/stamen_terrain/{$z}/{$x}/{$y}.png` |
| `stamen_toner.xml` | Stadia Stamen Toner (B&W TAC) | 0–20 | `https://tiles.stadiamaps.com/tiles/stamen_toner/{$z}/{$x}/{$y}.png` |
| `usgs_topo.xml` | USGS US Topo | 0–16 | `https://basemap.nationalmap.gov/arcgis/rest/services/USGSTopo/MapServer/tile/{$z}/{$y}/{$x}` |
| `usgs_imagery.xml` | USGS Imagery Topo | 0–16 | `https://basemap.nationalmap.gov/arcgis/rest/services/USGSImageryTopo/MapServer/tile/{$z}/{$y}/{$x}` |

**Notes:**
- Stamen Toner: high-contrast black-and-white — ideal for tactical overlays where color COP markers need to pop
- OpenTopoMap: best for elevation + trail detail, hiking-grade contours
- USGS tiles: official US government quads, CONUS only

---

## Government / Specialized (5 sources)

| File | Name | Zoom | URL |
|------|------|------|-----|
| `usgs_shadedrelief.xml` | USGS Shaded Relief | 0–13 | `https://basemap.nationalmap.gov/.../USGSShadedReliefOnly/...` |
| `noaa_charts.xml` | NOAA Nautical Charts | 4–16 | `https://tileservice.charts.noaa.gov/tiles/50000_1/{$z}/{$x}/{$y}.png` |
| `sec_aviation.xml` | FAA Sectional Charts | 4–11 | `https://nfdc.faa.gov/tpo/rest/tiles/Sectional_Current/{$z}/{$x}/{$y}.png` |
| `opensky_wms.xml` | OpenAIP VFR Overlay | 4–11 | `https://tiles.openaip.net/geowebcache/service/tms/...` |
| `mapbox_dark.xml` | OpenRailwayMap (Infra Overlay) | 2–19 | `https://a.tiles.openrailwaymap.org/standard/{$z}/{$x}/{$y}.png` |

**Notes:**
- NOAA charts: excellent for coastal/maritime ops, harbor navigation
- FAA Sectional: shows airspace classes, NOTAMs, restricted areas, airports, obstacles — critical for drone ops
- OpenAIP: transparent overlay showing airports, airspace, obstructions globally
- OpenRailwayMap: transparent infrastructure overlay, useful for route planning near rail corridors

---

## OSM / Open (3 sources)

| File | Name | Zoom | URL |
|------|------|------|-----|
| `osm.xml` | OpenStreetMap (Standard) | 0–19 | `https://tile.openstreetmap.org/{$z}/{$x}/{$y}.png` |
| `osm_humanitarian.xml` | OSM Humanitarian | 0–19 | `https://tile-a.openstreetmap.fr/hot/{$z}/{$x}/{$y}.png` |
| `osm_cycle.xml` | OSM CycleMap | 0–17 | `https://tile.waymarkedtrails.org/cycling/{$z}/{$x}/{$y}.png` |

---

## Dark / Night Ops (3 sources)

| File | Name | Zoom | URL |
|------|------|------|-----|
| `carto_dark.xml` | Carto Dark Matter | 0–19 | `https://cartodb-basemaps-a.global.ssl.fastly.net/dark_all/{$z}/{$x}/{$y}.png` |
| `satellite_night.xml` | NASA GIBS Blue Marble (Night) | 0–8 | `https://gibs.earthdata.nasa.gov/wmts/epsg3857/.../VIIRS_Black_Marble/...` |
| `nasa_worldview.xml` | NASA True Color (MODIS) | 0–8 | `https://gibs.earthdata.nasa.gov/wmts/epsg3857/.../MODIS_Terra_.../2024-01-01/...` |

**Notes:**
- Carto Dark Matter: best night ops base map — city outlines visible, minimal glare
- NASA Blue Marble Night: city lights from orbit, low zoom only (global/regional scale)
- MODIS True Color: actual satellite imagery by date, useful for vegetation/terrain assessment

---

## Other (3 sources)

| File | Name | Zoom | URL |
|------|------|------|-----|
| `carto_light.xml` | Carto Positron | 0–19 | `https://cartodb-basemaps-a.global.ssl.fastly.net/light_all/{$z}/{$x}/{$y}.png` |
| `thunderforest_outdoor.xml` | Thunderforest Outdoors | 0–22 | `https://tile.thunderforest.com/outdoors/{$z}/{$x}/{$y}.png?apikey=...` |
| `bing_aerial.xml` | Bing Aerial | 0–19 | `http://ecn.t3.tiles.virtualearth.net/tiles/a{$q}.jpeg?g=1` |

**Notes:**
- Thunderforest Outdoors includes an API key (`a5dd6a2f1c934394bce6b0fb077203eb`) — free tier, 150k tiles/month
- Bing requires `{$q}` quadkey format and `<xyzProvider>bing</xyzProvider>` in XML — different from standard XYZ tiles
- Carto Positron: clean light map for daytime planning sessions

---

## Supercharge Prefs (pushed as `supercharge.pref`)

Key settings confirmed working on S24 Ultra / ATAK-CIV 5.7.0.10:

```xml
<entry key="tileCacheSize" class="class java.lang.Integer">512000</entry>  <!-- 512MB tile cache -->
<entry key="reportInterval" class="class java.lang.Integer">3</entry>       <!-- SA every 3s -->
<entry key="staleInterval" class="class java.lang.Integer">30</entry>       <!-- contacts stale after 30s -->
<entry key="coord_display_pref" class="class java.lang.String">MGRS</entry>
<entry key="show_track_history" class="class java.lang.Boolean">true</entry>
<entry key="track_history_points" class="class java.lang.Integer">60</entry>  <!-- 3min @ 3s interval -->
<entry key="keep_screen_on" class="class java.lang.Boolean">true</entry>
<entry key="atakNightMode" class="class java.lang.Boolean">true</entry>
<entry key="speed_unit" class="class java.lang.String">mph</entry>
<entry key="alt_unit" class="class java.lang.String">feet</entry>
<entry key="range_unit" class="class java.lang.String">Metric</entry>
```

---

## Recommended Defaults by Mission Type

| Scenario | Primary Map | Overlay |
|----------|------------|---------|
| Urban ops / navigation | Google Hybrid | ESRI Street labels |
| Rural / trail | OpenTopoMap or USGS Topo | Stamen Terrain |
| Night ops | Carto Dark Matter | — |
| Maritime / coastal | ESRI Ocean Basemap | NOAA Nautical Charts |
| Aviation / drone | Google Hybrid | FAA Sectional |
| General field / quick SA | Google Satellite | — |
| Planning session (daylight) | ESRI Topo | — |
| Infrastructure recon | Google Hybrid | OpenRailwayMap |
| B&W tac (marker clarity) | Stamen Toner | — |

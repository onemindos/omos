# Tactical Map - Feature Documentation

> Military-grade location tracking and navigation for OneMind OS

## Overview

The Tactical Map is a real-time location tracking system built on Google Maps Flutter SDK. It provides GPS tracking, geofence management, place search, and turn-by-turn directions with a dark tactical theme.

**Location**: `frontend/lib/tactical/screens/tactical_map_screen.dart`
**Route**: `/tactical-map`
**Navigation**: OneMind OS > Tactical Map

---

## Current Features

### 1. Real-Time GPS Tracking

**What it does**: Continuously tracks user location with activity recognition (walking, running, driving, etc.)

**Implementation**:
```dart
// Uses LocationService from device_hub_providers
final locationService = ref.read(locationServiceProvider);
await locationService.startTracking(mode: LocationTrackingMode.balanced);

// Watches position updates via Riverpod
final positionAsync = ref.watch(currentPositionProvider);
```

**Key Components**:
- `DevicePosition` model with lat/lng, speed, heading, altitude, accuracy
- `ActivityType` enum: stationary, walking, running, cycling, driving
- Auto-follow mode that keeps map centered on user
- Activity icons that change based on detected movement type

**Status Card Displays**:
| Metric | Description | Format |
|--------|-------------|--------|
| Coordinates | Current lat/lng | `40.712800, -74.006000` |
| Activity | Detected movement | Walking, Driving, etc. |
| Speed | Current velocity | `25.5 mph` |
| Heading | Direction of travel | `45° NE` |
| Altitude | Elevation | `125m` |
| Accuracy | GPS precision | `±10m` (color-coded) |

---

### 2. Multiple Map Types

**What it does**: Toggle between 4 map visualization styles

**Implementation**:
```dart
MapType _mapType = MapType.normal;

void _toggleMapType() {
  switch (_mapType) {
    case MapType.normal: _mapType = MapType.satellite; break;
    case MapType.satellite: _mapType = MapType.terrain; break;
    case MapType.terrain: _mapType = MapType.hybrid; break;
    case MapType.hybrid: _mapType = MapType.normal; break;
  }
}
```

**Available Types**:
| Type | Description |
|------|-------------|
| Standard | Default road map with dark tactical styling |
| Satellite | Aerial imagery |
| Terrain | Topographic with elevation shading |
| Hybrid | Satellite + road labels |

---

### 3. Dark Tactical Map Style

**What it does**: Custom JSON styling for military/tactical aesthetic

**Implementation**:
```dart
const darkStyle = '''
[
  {"elementType": "geometry", "stylers": [{"color": "#1d1d1d"}]},
  {"elementType": "labels.text.fill", "stylers": [{"color": "#8ec3b9"}]},
  {"featureType": "road", "elementType": "geometry", "stylers": [{"color": "#304a7d"}]},
  {"featureType": "water", "elementType": "geometry", "stylers": [{"color": "#17263c"}]},
  // ... more styling rules
]
''';
await _mapController?.setMapStyle(darkStyle);
```

**Style Features**:
- Dark background (#1d1d1d)
- Cyan/teal labels (#8ec3b9)
- Navy roads (#304a7d, #2c3e50)
- Hidden POIs and transit
- Reduced visual clutter

---

### 4. Geofence Zones

**What it does**: Create and visualize circular geographic boundaries (home, work, gym, etc.)

**Implementation**:
```dart
// Geofence model
class Geofence {
  final String id;
  final String name;
  final double latitude;
  final double longitude;
  final double radius; // meters
}

// Build circles for map
Set<Circle> _buildGeofenceCircles(List<Geofence> geofences) {
  return geofences.map((g) {
    Color fillColor, strokeColor;
    // Color based on zone type
    if (g.name.toLowerCase().contains('home')) {
      fillColor = Colors.green.withOpacity(0.2);
      strokeColor = Colors.green;
    } else if (g.name.toLowerCase().contains('work')) {
      fillColor = Colors.blue.withOpacity(0.2);
      strokeColor = Colors.blue;
    }
    // ... more zone types

    return Circle(
      circleId: CircleId(g.id),
      center: LatLng(g.latitude, g.longitude),
      radius: g.radius,
      fillColor: fillColor,
      strokeColor: strokeColor,
      strokeWidth: 2,
    );
  }).toSet();
}
```

**Zone Colors**:
| Zone Type | Fill | Stroke |
|-----------|------|--------|
| Home | Green 20% | Green |
| Work | Blue 20% | Blue |
| Gym | Orange 20% | Orange |
| Other | Red 20% | Red (primary) |

**Add Zone Dialog**:
- Name input field
- Radius slider (50m - 500m)
- Uses current GPS location as center
- Stored via `geofenceListNotifierProvider`

---

### 5. Place Search (Google Places Autocomplete)

**What it does**: Search for any location with autocomplete suggestions

**API Used**: Google Places Autocomplete API

**Implementation**:
```dart
Future<void> _searchPlaces(String query) async {
  final url = Uri.parse(
    'https://maps.googleapis.com/maps/api/place/autocomplete/json'
    '?input=${Uri.encodeComponent(query)}'
    '&location=$currentLocation'
    '&radius=50000'
    '&key=$_apiKey',
  );

  final response = await http.get(url);
  final data = json.decode(response.body);
  _searchResults = data['predictions'].map((p) => {
    'place_id': p['place_id'],
    'description': p['description'],
    'main_text': p['structured_formatting']['main_text'],
  }).toList();
}

// Get place details for coordinates
Future<void> _selectPlace(String placeId, String name) async {
  final url = Uri.parse(
    'https://maps.googleapis.com/maps/api/place/details/json'
    '?place_id=$placeId'
    '&fields=geometry'
    '&key=$_apiKey',
  );
  // Extract lat/lng and animate map
}
```

**Features**:
- Real-time autocomplete (triggers after 3 characters)
- Location-biased results (50km radius from current position)
- Tap result to navigate and get directions

---

### 6. Directions & Routing

**What it does**: Calculate and display driving routes with distance/duration

**API Used**: Google Directions API

**Implementation**:
```dart
Future<void> _getDirections(LatLng origin, LatLng destination) async {
  final url = Uri.parse(
    'https://maps.googleapis.com/maps/api/directions/json'
    '?origin=${origin.latitude},${origin.longitude}'
    '&destination=${destination.latitude},${destination.longitude}'
    '&mode=driving'
    '&key=$_apiKey',
  );

  final response = await http.get(url);
  final data = json.decode(response.body);
  final route = data['routes'][0];
  final leg = route['legs'][0];

  // Decode polyline points
  _routePoints = _decodePolyline(route['overview_polyline']['points']);
  _routeInfo = '${leg['distance']['text']} • ${leg['duration']['text']}';
}

// Google polyline decoding algorithm
List<LatLng> _decodePolyline(String encoded) {
  // Standard polyline decoding implementation
}
```

**Route Display**:
- Red polyline on map (4px width)
- Info card showing: destination name, distance, duration
- Auto-zoom to fit entire route
- Clear route button

---

### 7. Markers

**What it does**: Display location markers for user, destination, and zones

**Marker Types**:
| Marker | Color | Description |
|--------|-------|-------------|
| User | Cyan | Current location with activity info |
| Destination | Red | Selected search result |
| Zone Centers | Violet | Geofence center points |

---

## Google APIs Required

| API | Purpose | Endpoint |
|-----|---------|----------|
| Maps SDK for iOS/Android | Map rendering | Native SDK |
| Places API | Autocomplete search | `places/autocomplete/json` |
| Places API | Place details | `places/details/json` |
| Directions API | Route calculation | `directions/json` |

**Setup Required**:
1. Enable APIs in Google Cloud Console
2. Set API key in environment: `GOOGLE_MAPS_API_KEY`
3. Add to Flutter build: `--dart-define=GOOGLE_MAPS_API_KEY=xxx`

---

## File Dependencies

```
tactical_map_screen.dart
├── google_maps_flutter          # Map widget
├── http                         # API calls
├── flutter_riverpod             # State management
├── tactical.dart (theme)        # TacticalColors, TacticalText
├── device_hub_providers.dart    # Location providers
│   ├── locationServiceProvider
│   ├── currentPositionProvider
│   ├── locationStateProvider
│   └── geofenceListNotifierProvider
├── models.dart                  # DevicePosition, Geofence, ActivityType
└── environment.dart             # API key config
```

---

## Future Features

### High Priority

#### 1. Offline Maps
**Description**: Download map tiles for offline use
**Implementation**: Use `flutter_map` with `mbtiles` or cached tiles
**Benefit**: Works without internet in remote areas

#### 2. Route Turn-by-Turn Navigation
**Description**: Voice-guided navigation with upcoming turn alerts
**Implementation**:
- Parse `steps` array from Directions API
- Use `flutter_tts` for voice prompts
- Show maneuver icons (turn left, merge, etc.)
**Benefit**: Full navigation experience

#### 3. Traffic Layer
**Description**: Real-time traffic conditions overlay
**Implementation**: Enable `trafficEnabled: true` on GoogleMap widget + Directions API with `departure_time`
**Benefit**: Accurate ETAs and route optimization

#### 4. Multiple Travel Modes
**Description**: Support walking, cycling, transit directions
**Implementation**: Add mode selector, change `&mode=` param
**Options**: driving, walking, bicycling, transit

#### 5. Waypoints/Multi-Stop Routes
**Description**: Add intermediate stops to routes
**Implementation**: Add `&waypoints=` param to Directions API
**UI**: Drag-and-drop route editing

---

### Medium Priority

#### 6. Location History / Breadcrumb Trail
**Description**: Show path traveled over time
**Implementation**: Store positions in local DB, draw polyline trail
**Options**: Last hour, today, this week

#### 7. Share Location
**Description**: Send current location to contacts/agents
**Implementation**: Generate shareable link or integrate with messaging

#### 8. Geofence Events/Notifications
**Description**: Trigger actions when entering/exiting zones
**Implementation**:
```dart
// Use geofencing plugin
geofenceService.onGeofenceStatusChanged.listen((status) {
  if (status.type == GeofenceType.enter) {
    // Trigger agent, notification, etc.
  }
});
```
**Use Cases**: "Arrived at work" automation, reminders

#### 9. Custom Markers/Icons
**Description**: Replace default markers with custom tactical icons
**Implementation**: Use `BitmapDescriptor.fromAssetImage()` or SVG
**Styles**: Military, tech, custom user icons

#### 10. Heatmaps
**Description**: Visualize location frequency/density
**Implementation**: `google_maps_flutter_heatmap` plugin
**Use Cases**: Where you spend most time, activity patterns

---

### Lower Priority

#### 11. 3D Buildings
**Description**: Enable 3D building extrusion
**Implementation**: `buildingsEnabled: true` (default on)

#### 12. Street View
**Description**: 360° panoramic street-level imagery
**Implementation**: `google_maps_flutter_web` or native Street View SDK

#### 13. Indoor Maps
**Description**: Floor plans for malls, airports, etc.
**Implementation**: Google Indoor Maps API (limited availability)

#### 14. AR Navigation
**Description**: Augmented reality walking directions
**Implementation**: ARCore/ARKit + camera overlay
**Complexity**: High

#### 15. Speed Alerts
**Description**: Notify when exceeding speed limit
**Implementation**: Compare `position.speed` to road speed limits (Roads API)

#### 16. Fuel/EV Station Finder
**Description**: Find nearby gas or charging stations
**Implementation**: Places API with `type=gas_station` or `type=electric_vehicle_charging_station`

#### 17. Weather Overlay
**Description**: Show weather conditions on map
**Implementation**: OpenWeatherMap tile layer or custom overlay

#### 18. Saved Places
**Description**: Bookmark favorite locations
**Implementation**: Local storage with custom marker display

#### 19. Distance Measurement Tool
**Description**: Measure distance between two points
**Implementation**: Tap two points, calculate haversine distance

#### 20. Export GPX/KML
**Description**: Export routes and tracks to standard formats
**Implementation**: Generate GPX/KML XML from route points

---

## Agent Integration Ideas

### Location-Aware Agents

```python
# Backend tool for agents to access location
@tool
def get_current_location() -> dict:
    """Get user's current GPS location"""
    return {
        "latitude": 40.7128,
        "longitude": -74.0060,
        "activity": "driving",
        "zone": "work"
    }

@tool
def get_directions(destination: str) -> dict:
    """Get directions to a destination"""
    # Call Google Directions API
    return {
        "distance": "5.2 mi",
        "duration": "12 mins",
        "route_summary": "via I-95 N"
    }

@tool
def check_geofence(zone_name: str) -> dict:
    """Check if user is in a geofence zone"""
    return {
        "in_zone": True,
        "zone_name": "Home",
        "since": "2 hours ago"
    }
```

### Use Cases
- **"Navigate to my next meeting"** - Agent gets calendar event, extracts location, triggers directions
- **"Tell me when I get home"** - Agent monitors geofence, sends notification on arrival
- **"Where am I?"** - Agent responds with current location context
- **"Find nearest coffee shop"** - Agent searches places, shows on map

---

## Performance Considerations

| Optimization | Implementation |
|--------------|----------------|
| Debounce search | Wait 300ms after typing stops |
| Cache API responses | Store recent searches |
| Limit markers | Cluster when > 50 markers |
| Battery optimization | Use `LocationTrackingMode.balanced` |
| Lazy load | Only fetch directions when needed |

---

## Testing Checklist

- [ ] Location permission request
- [ ] GPS tracking accuracy
- [ ] Map renders correctly
- [ ] Search autocomplete works
- [ ] Place selection navigates map
- [ ] Directions polyline displays
- [ ] Geofence circles render
- [ ] Add zone dialog works
- [ ] Map type toggle cycles
- [ ] Center on user button works
- [ ] Route info card appears/dismisses
- [ ] Status card shows all metrics
- [ ] Dark theme applied
- [ ] Works offline (graceful degradation)

---

## Changelog

| Date | Version | Changes |
|------|---------|---------|
| 2025-02-03 | 1.0 | Initial implementation with GPS tracking, map types, dark theme |
| 2025-02-03 | 1.1 | Added geofence zones, place search, directions |


---

## Related

- [[00-24 SO (Sovereign Operations)/06 Research (R&D) (The Lab)/06 Research|06 Research]]

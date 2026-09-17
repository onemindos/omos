# OneMind OS Vehicles System

> **Document Status:** CANONICAL
> **Created:** 2026-01-28
> **Section:** HIVE - Physical Interface
> **Related Issue:** [GitHub Issue #97](https://github.com/Zeus-Delacruz/OneMind-OS/issues/97)

---

## Overview

The Vehicles System is a core component of the HIVE section, providing fleet management capabilities for multiple vehicle types. It serves as the interface between OneMind OS and physical transportation assets, enabling real-time monitoring, remote control, and AI-powered diagnostics.

### Philosophy

Vehicles represent mobile extensions of the OneMind ecosystem. Unlike stationary IoT devices, vehicles operate across varying environments and connectivity conditions. The system is designed to:

1. **Unify** - Manage all vehicle types through a single interface
2. **Adapt** - Support multiple connection methods (API, OBD-II, Edge)
3. **Evolve** - Enable AI-powered software module installation for enhanced capabilities

---

## Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           VEHICLES SYSTEM ARCHITECTURE                       │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐         │
│  │   VEHICLES UI   │    │   PROVIDERS     │    │   DATA MODELS   │         │
│  │                 │    │                 │    │                 │         │
│  │ VehiclesScreen  │◄──►│ vehiclesProvider│◄──►│    Vehicle      │         │
│  │ VehicleDetail   │    │ vehicleProvider │    │ VehicleTelemetry│         │
│  │ VehicleCard     │    │ telemetryProvider    │   VehicleTrip   │         │
│  └─────────────────┘    │ tripsProvider   │    └─────────────────┘         │
│           │             └─────────────────┘             │                   │
│           │                      │                      │                   │
│           ▼                      ▼                      ▼                   │
│  ┌─────────────────────────────────────────────────────────────────┐       │
│  │                      CONNECTION LAYER (Future)                   │       │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐        │       │
│  │  │ Tesla API│  │ Ford Pass│  │  OBD-II  │  │   Edge   │        │       │
│  │  │          │  │          │  │ Bluetooth│  │  Device  │        │       │
│  │  └──────────┘  └──────────┘  └──────────┘  └──────────┘        │       │
│  └─────────────────────────────────────────────────────────────────┘       │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## File Structure

```
frontend/lib/hive/vehicles/
├── models/
│   ├── vehicle.dart           # Core data models with Freezed
│   ├── vehicle.freezed.dart   # Generated Freezed code
│   └── vehicle.g.dart         # Generated JSON serialization
├── providers/
│   └── vehicles_provider.dart # Riverpod state management
├── screens/
│   ├── vehicles_screen.dart   # Fleet list view
│   └── vehicle_detail_screen.dart # Individual vehicle view
└── widgets/
    └── vehicle_card.dart      # Vehicle card component
```

---

## Data Models

### Vehicle

The primary model representing a vehicle in the fleet.

```dart
@freezed
sealed class Vehicle with _$Vehicle {
  const factory Vehicle({
    required String id,
    required String name,
    required VehicleType type,
    required String make,
    required String model,
    required int year,
    String? vin,
    String? licensePlate,
    String? color,
    String? photoUrl,
    @Default(ConnectionType.manual) ConnectionType connectionType,
    @Default(ConnectionStatus.disconnected) ConnectionStatus connectionStatus,
    DateTime? lastSeen,
    String? edgeDeviceId,
    @Default([]) List<String> installedModules,
    VehicleTelemetry? latestTelemetry,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) = _Vehicle;

  factory Vehicle.fromJson(Map<String, dynamic> json) => _$VehicleFromJson(json);
}
```

#### Vehicle Types

| Type | Icon | Description |
|------|------|-------------|
| `car` | `Icons.directions_car` | Standard automobiles, trucks, SUVs |
| `motorcycle` | `Icons.two_wheeler` | Motorcycles, scooters |
| `boat` | `Icons.sailing` | Watercraft, yachts |
| `jet` | `Icons.flight` | Private aircraft |
| `rv` | `Icons.rv_hookup` | Recreational vehicles |
| `atv` | `Icons.terrain` | All-terrain vehicles |

#### Connection Types

| Type | Description | Use Case |
|------|-------------|----------|
| `obd2` | OBD-II Bluetooth/WiFi adapter | Universal vehicle diagnostics |
| `api` | Manufacturer API (Tesla, Ford, etc.) | Native cloud integration |
| `edge` | OneMind Edge device | Full telemetry with custom hardware |
| `manual` | Manual entry only | Vehicles without connectivity |

#### Connection Status

| Status | Visual Indicator | Description |
|--------|------------------|-------------|
| `connected` | Green "LIVE" badge | Real-time data streaming |
| `disconnected` | Gray "OFFLINE" badge | No current connection |
| `syncing` | Animated indicator | Data synchronization in progress |

---

### VehicleTelemetry

Real-time and historical telemetry data from vehicles.

```dart
@freezed
sealed class VehicleTelemetry with _$VehicleTelemetry {
  const factory VehicleTelemetry({
    required String vehicleId,
    required DateTime timestamp,

    // Location
    double? latitude,
    double? longitude,
    double? altitude,
    double? heading,
    double? speed,

    // Power
    double? batteryPercent,      // EV battery level (0-100)
    double? fuelLevel,           // ICE fuel level (0-100)
    double? rangeRemaining,      // Estimated range in miles
    bool? isCharging,

    // Vehicle State
    double? odometer,            // Total miles/km
    double? engineTemp,
    double? batteryVoltage,      // 12V system voltage
    bool? doorsLocked,
    bool? trunkOpen,
    bool? windowsOpen,
  }) = _VehicleTelemetry;

  factory VehicleTelemetry.fromJson(Map<String, dynamic> json) =>
      _$VehicleTelemetryFromJson(json);
}
```

#### Telemetry Fields by Vehicle Type

| Field | Car (EV) | Car (ICE) | Motorcycle | Boat | Jet |
|-------|----------|-----------|------------|------|-----|
| `batteryPercent` | ✓ | - | - | - | - |
| `fuelLevel` | - | ✓ | ✓ | ✓ | ✓ |
| `rangeRemaining` | ✓ | ✓ | ✓ | ✓ | ✓ |
| `doorsLocked` | ✓ | ✓ | - | - | - |
| `engineTemp` | ✓ | ✓ | ✓ | ✓ | ✓ |
| `batteryVoltage` | ✓ | ✓ | ✓ | ✓ | ✓ |

---

### VehicleTrip

Historical trip data for analytics and logging.

```dart
@freezed
sealed class VehicleTrip with _$VehicleTrip {
  const factory VehicleTrip({
    required String id,
    required String vehicleId,
    required DateTime startTime,
    DateTime? endTime,
    required double startLatitude,
    required double startLongitude,
    double? endLatitude,
    double? endLongitude,
    double? distance,            // Miles
    int? duration,               // Minutes
    double? avgSpeed,            // MPH
    double? maxSpeed,            // MPH
    String? startAddress,
    String? endAddress,
  }) = _VehicleTrip;

  factory VehicleTrip.fromJson(Map<String, dynamic> json) =>
      _$VehicleTripsFromJson(json);
}
```

---

## State Management

### Providers

The vehicles system uses Riverpod for state management with the following providers:

```dart
/// Provider for the list of vehicles
final vehiclesProvider = StateNotifierProvider<VehiclesNotifier, AsyncValue<List<Vehicle>>>;

/// Provider for a single vehicle by ID
final vehicleProvider = Provider.family<Vehicle?, String>;

/// Provider for vehicle telemetry history
final vehicleTelemetryProvider = FutureProvider.family<List<VehicleTelemetry>, String>;

/// Provider for vehicle trips
final vehicleTripsProvider = FutureProvider.family<List<VehicleTrip>, String>;
```

### VehiclesNotifier

The state notifier provides the following methods:

| Method | Description | Parameters |
|--------|-------------|------------|
| `refresh()` | Reload all vehicles from backend | - |
| `addVehicle(vehicle)` | Add a new vehicle | `Vehicle` |
| `updateVehicle(vehicle)` | Update existing vehicle | `Vehicle` |
| `deleteVehicle(id)` | Remove a vehicle | `String id` |
| `lockVehicle(id)` | Send lock command | `String id` |
| `unlockVehicle(id)` | Send unlock command | `String id` |
| `startClimate(id)` | Start climate control | `String id` |
| `honkHorn(id)` | Trigger horn | `String id` |

---

## UI Components

### VehiclesScreen

The main fleet view displaying all vehicles in a grid layout.

**Route:** `/vehicles`

**Features:**
- Stats header showing total vehicles, connected count, and types
- Filter chips for vehicle types
- Responsive grid layout (2 columns on mobile, 3+ on tablet)
- Pull-to-refresh functionality
- Empty state with add vehicle CTA

**Visual Structure:**
```
┌────────────────────────────────────────┐
│  VEHICLES                    [+ Add]   │
├────────────────────────────────────────┤
│  ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐  │
│  │  4   │ │  2   │ │  2   │ │  1   │  │
│  │Total │ │ Live │ │ Cars │ │Bikes │  │
│  └──────┘ └──────┘ └──────┘ └──────┘  │
├────────────────────────────────────────┤
│ [All] [Cars] [Motorcycles] [Boats]...  │
├────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐      │
│  │  [IMAGE]    │  │  [IMAGE]    │      │
│  │  Model S    │  │Street Glide │      │
│  │  78% ⚡285mi│  │  65% ⛽     │      │
│  │ [Lock][Map] │  │ [Lock][Map] │      │
│  └─────────────┘  └─────────────┘      │
│  ┌─────────────┐  ┌─────────────┐      │
│  │  [IMAGE]    │  │  [IMAGE]    │      │
│  │  F-150      │  │  Sea Ray    │      │
│  │  45% ⚡142mi│  │  82% ⛽     │      │
│  │ [Lock][Map] │  │ [Lock][Map] │      │
│  └─────────────┘  └─────────────┘      │
└────────────────────────────────────────┘
```

---

### VehicleDetailScreen

Detailed view for individual vehicle management.

**Route:** `/vehicles/:id`

**Tabs:**

| Tab | Contents |
|-----|----------|
| **STATUS** | Live telemetry, map placeholder, power/location/security stats |
| **CONTROLS** | Vehicle commands (Lock, Climate, Horn, Locate) |
| **TRIPS** | Trip history with distance, duration, addresses |
| **INFO** | Vehicle specs, VIN, installed modules |

**Visual Structure:**
```
┌────────────────────────────────────────┐
│  ← Model S                    [⋮]     │
├────────────────────────────────────────┤
│            ┌──────────────┐            │
│            │              │            │
│            │   [IMAGE]    │            │
│            │              │            │
│            └──────────────┘            │
│         2024 Tesla Model S Plaid       │
│              ● LIVE                    │
├────────────────────────────────────────┤
│  [STATUS]  [CONTROLS]  [TRIPS]  [INFO] │
├────────────────────────────────────────┤
│  ┌─────────────────────────────────┐   │
│  │         MAP PLACEHOLDER         │   │
│  │     33.4484, -112.0740         │   │
│  └─────────────────────────────────┘   │
│                                        │
│  ┌────────┐ ┌────────┐ ┌────────┐     │
│  │  78%   │ │ 285 mi │ │ 12.4k  │     │
│  │Battery │ │ Range  │ │Odometer│     │
│  └────────┘ └────────┘ └────────┘     │
└────────────────────────────────────────┘
```

---

### VehicleCard

Compact card widget for the fleet grid.

**Props:**
| Prop | Type | Description |
|------|------|-------------|
| `vehicle` | `Vehicle` | Vehicle data |
| `onTap` | `VoidCallback` | Navigate to detail |
| `onLock` | `VoidCallback` | Lock command |
| `onUnlock` | `VoidCallback` | Unlock command |

**Visual Elements:**
- Vehicle photo with gradient overlay
- Type badge (top-left, red background)
- Connection status badge (top-right, green/gray)
- Vehicle name and year/make/model
- Telemetry stats (battery/fuel, range, odometer)
- Quick action buttons (Lock/Unlock, Locate)

---

## Routes

| Route | Screen | Description |
|-------|--------|-------------|
| `/vehicles` | `VehiclesScreen` | Fleet list view |
| `/vehicles/:id` | `VehicleDetailScreen` | Individual vehicle detail |

**Router Configuration:**
```dart
GoRoute(
  path: '/vehicles',
  name: 'vehicles',
  builder: (context, state) => const VehiclesScreen(),
),
GoRoute(
  path: '/vehicles/:id',
  name: 'vehicle-detail',
  builder: (context, state) {
    final vehicleId = state.pathParameters['id']!;
    return VehicleDetailScreen(vehicleId: vehicleId);
  },
),
```

---

## Design System

### Colors

| Element | Color | Usage |
|---------|-------|-------|
| Connected Status | `#00FF00` | Live indicator, active lock |
| Disconnected Status | `#FFFFFF50` | Offline indicator |
| Vehicle Type Badge | `#FF0000` | Type label background |
| Card Background | `#151515` | Card surface |
| Card Border (Connected) | `#00FF0030` | Green glow effect |
| Card Border (Disconnected) | `#FFFFFF10` | Subtle border |

### Typography

| Element | Style |
|---------|-------|
| Vehicle Name | 18px, weight 700, white |
| Year/Make/Model | 12px, weight 400, white 70% |
| Status Badge | 9px, weight 700, letter-spacing 0.5 |
| Stat Value | 14px, weight 600, monospace |
| Stat Label | 9px, weight 400, white 40% |

---

## Mock Data

The current implementation includes mock data for testing:

| Vehicle | Type | Connection | Status |
|---------|------|------------|--------|
| Model S | Car (EV) | Tesla API | Connected |
| Street Glide | Motorcycle | OBD-II | Connected |
| F-150 | Car (EV) | Ford Pass API | Disconnected |
| Sea Ray | Boat | Edge Device | Disconnected |

---

## Future Integration

### Phase 1: Backend Integration
- REST API endpoints for vehicle CRUD
- WebSocket for real-time telemetry
- PostgreSQL storage for vehicles and trips

### Phase 2: Connection Integrations
- Tesla API integration
- Ford Pass API integration
- OBD-II Bluetooth adapter support

### Phase 3: Edge Deployment
- OneMind Edge device protocol
- Local telemetry processing
- Offline-first architecture

### Phase 4: AI Enhancement
- Automatic software module detection
- Predictive maintenance alerts
- Driving pattern analysis

See [GitHub Issue #97](https://github.com/Zeus-Delacruz/OneMind-OS/issues/97) for the complete roadmap.

---

## Related Documentation

- [Navigation System](./NAVIGATION-SYSTEM.md) - HIVE section definition
- [Edge Devices](./EDGE-DEVICES.md) - Edge deployment architecture (planned)

---

*Part of the HIVE section - Physical Interface layer of OneMind OS*

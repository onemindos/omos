# Apple Health

> **iOS/macOS central biometric hub aggregating all health data**

## Overview

Apple Health is the central aggregator for all biometric and fitness data on Apple devices. Pulls data from Oura, Strava, Trainerize, Apple Watch, and hundreds of apps into a unified timeline. HealthKit API enables programmatic access for HP domain automation.

## Status

| Attribute   | Value                          |
| ----------- | ------------------------------ |
| Status      | ✅ Active                       |
| Category    | Health                         |
| Type        | Health Data Platform           |
| Integration | HealthKit API, iOS App         |

## Use Cases in OneMind

1. **Data Aggregation** — Single source for all biometric data on iPhone/Mac
2. **HP Domain Feed** — Export health data for HP domain scoring
3. **Workout History** — Complete exercise log from all connected apps
4. **Trend Analysis** — Long-term health trend visualization
5. **Medical Records** — Store and access health records on iPhone

## Key Data Categories

| Category      | Sources                           |
| ------------- | --------------------------------- |
| Sleep         | Oura, Apple Watch                 |
| Activity      | Apple Watch, Strava, AllTrails    |
| Heart Rate    | Oura, Apple Watch                 |
| Workouts      | Trainerize, Strava                |
| Nutrition     | Manual entry, third-party apps    |
| Mental Health | Mood journaling apps              |

## HealthKit Access (iOS App)

```swift
// Request HealthKit authorization
let healthStore = HKHealthStore()
let types = Set([HKObjectType.quantityType(forIdentifier: .heartRate)!])
healthStore.requestAuthorization(toShare: nil, read: types) { success, error in }
```

## Links

- Website: [apple.com/health](https://www.apple.com/health)
- HealthKit Docs: [developer.apple.com/healthkit](https://developer.apple.com/healthkit)


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]

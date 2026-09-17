# Dart / Flutter

> **Google's cross-platform UI framework — build iOS, Android, and web from one codebase**

## Overview

Dart is the language; Flutter is the UI framework built on Dart. Together they enable building native-quality apps for iOS, Android, web, and desktop from a single codebase. Used in onemind-mobile as the cross-platform alternative to native Swift/Kotlin.

## Status

| Attribute   | Value                                      |
| ----------- | ------------------------------------------ |
| Status      | ✅ Active                                   |
| Category    | Development / Language + Framework         |
| Type        | Cross-Platform Mobile Framework            |
| Integration | Flutter SDK, Dart pub, iOS, Android, Web   |

## Use Cases in OneMind

1. **onemind-mobile** — Cross-platform iOS + Android app (VisionClaw Android base)
2. **Agent UI** — Mobile chat interface for OpenClaw agent interaction
3. **Dashboard Widgets** — HP domain health metrics on mobile
4. **Cross-Platform** — Single codebase for iOS and Android vs. SwiftUI-only

## Key Packages (pub.dev)

| Package       | Use                              |
| ------------- | -------------------------------- |
| `dio`         | HTTP client for REST APIs        |
| `riverpod`    | State management                 |
| `flutter_secure_storage` | Secure credential storage |
| `health`      | HealthKit + Google Fit bridge    |
| `web_socket_channel` | WebSocket for live agent chat |

## Quick Start

```bash
flutter create onemind_mobile
cd onemind_mobile
flutter pub get
flutter run                    # iOS Simulator or Android
flutter build ios --release    # App Store build
```

## Repo

- onemind-mobile: `/Users/zeuslegacy/Documents/onemind-mobile/`
- VisionClaw Android sample: `samples/VisionClaw-Android/`

## Links

- Flutter: [flutter.dev](https://flutter.dev)
- Dart: [dart.dev](https://dart.dev)
- pub.dev: [pub.dev](https://pub.dev)


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]

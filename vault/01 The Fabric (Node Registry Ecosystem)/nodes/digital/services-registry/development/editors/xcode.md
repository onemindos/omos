---
title: "Xcode"
type: reference
status: active
tags: [platform, development, ide, apple, ios, macos, swift]
---
# Xcode

> Apple's IDE. Required for iOS/macOS/watchOS/tvOS development. Swift, SwiftUI, Interface Builder.

## Key Features
| Feature | Notes |
|---------|-------|
| Languages | Swift, Objective-C, C/C++ |
| UI | SwiftUI, Interface Builder (Storyboards) |
| Simulator | iOS/watchOS/tvOS simulators |
| Instruments | Performance profiling |
| TestFlight | Beta distribution |
| Signing | Code signing + provisioning |
| CLI | xcodebuild, xcrun, simctl |

## CLI
```bash
xcode-select --install      # Install CLI tools
xcodebuild -list            # List schemes/targets
xcrun simctl list devices   # List simulators
```

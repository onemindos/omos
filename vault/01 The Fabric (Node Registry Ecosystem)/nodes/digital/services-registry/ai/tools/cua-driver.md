---
title: "CUA Driver"
type: reference
status: active
tags: [platform, ai, tools, automation, computer-use, macos]
---
# CUA Driver

> Computer-Use Automation driver for macOS. Screenshots, AX tree walking, click/type/scroll via accessibility APIs.

## Key Features
| Feature | Notes |
|---------|-------|
| AX tree | Walk any app's accessibility tree |
| Screenshots | Per-window capture via ScreenCaptureKit |
| Click | Element-indexed or pixel-addressed |
| Type | Text insertion via AX or CGEvent |
| Scroll | Keyboard-based scrolling |
| Hotkey | Key combinations to any pid |
| Launch | Background app launch (no focus steal) |
| Agent cursor | Visual overlay showing AI actions |

## Architecture
```
Hermes Agent → MCP → CUA Driver → macOS Accessibility APIs
                                 → ScreenCaptureKit
                                 → CGEvent (keyboard/mouse)
                                 → SkyLight SPI (pid-targeted)
```

## Permissions Required
- Accessibility (TCC)
- Screen Recording (TCC)

## Use Cases
- Browser automation (Safari, Chrome)
- Desktop app control
- Form filling
- Visual verification
- End-to-end testing

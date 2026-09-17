# HoloMat2 UI Concepts

## Radial App Launcher

Extracted from [Concept-Bytes/HoloMat2](https://github.com/Concept-Bytes/HoloMat2)

The HoloMat2 presents an interesting UI pattern: a **radial app launcher** with animated transitions. This could inspire OneMind's Flutter mobile interface.

---

## The Concept

```
        ┌─────────────────────────────────────────┐
        │                                         │
        │           [📊]         [📷]             │
        │          Analytics    Camera            │
        │                                         │
        │    [🏠]        ┌───┐        [🤖]        │
        │   Smart       │ L │        Robot       │
        │   Home        └───┘        Control     │
        │                ↑                        │
        │          Main Circle                    │
        │      (Tap to toggle)                    │
        │                                         │
        │    [🖨️]                      [🌡️]       │
        │   3D Print              Sensors         │
        │                                         │
        │           [📅]         [💬]             │
        │          Calendar     Chat              │
        │                                         │
        └─────────────────────────────────────────┘
```

---

## Key Features

### 1. Center Toggle

- Main circle acts as show/hide toggle
- Tap to expand app circles outward
- Tap again to collapse back to center

### 2. Smooth Animations

```python
# Animation logic from HoloMat2
if elapsed_time < 0.5:  # 500ms animation
    t = elapsed_time / 0.5
    if self.visible:
        # Expand: center → final position
        self.center = (
            int((1 - t) * center_x + t * final_x),
            int((1 - t) * center_y + t * final_y)
        )
    else:
        # Collapse: final position → center
        self.center = (
            int(t * center_x + (1 - t) * final_x),
            int(t * center_y + (1 - t) * final_y)
        )
```

### 3. Touch Feedback

- Hover effect: circle grows slightly
- Audio feedback on tap
- Different sounds for: startup, toggle, confirm, reject

### 4. Radial Layout

```python
angle_step = 360 / num_apps
for i in range(num_apps):
    angle = math.radians(angle_step * i)
    x = center_x + int(distance * math.cos(angle))
    y = center_y + int(distance * math.sin(angle))
    positions.append((x, y))
```

---

## Flutter Implementation Concept

### Widget Structure

```dart
class RadialAppLauncher extends StatefulWidget {
  final List<AppItem> apps;
  final VoidCallback? onExpand;
  final VoidCallback? onCollapse;

  @override
  _RadialAppLauncherState createState() => _RadialAppLauncherState();
}

class _RadialAppLauncherState extends State<RadialAppLauncher>
    with SingleTickerProviderStateMixin {

  late AnimationController _controller;
  bool _isExpanded = false;

  @override
  void initState() {
    super.initState();
    _controller = AnimationController(
      duration: Duration(milliseconds: 500),
      vsync: this,
    );
  }

  void _toggle() {
    if (_isExpanded) {
      _controller.reverse();
      widget.onCollapse?.call();
    } else {
      _controller.forward();
      widget.onExpand?.call();
    }
    setState(() => _isExpanded = !_isExpanded);
  }

  @override
  Widget build(BuildContext context) {
    return Stack(
      alignment: Alignment.center,
      children: [
        // App circles
        ...widget.apps.asMap().entries.map((entry) {
          final index = entry.key;
          final app = entry.value;
          final angle = (2 * pi / widget.apps.length) * index - pi / 2;

          return AnimatedBuilder(
            animation: _controller,
            builder: (context, child) {
              final distance = _controller.value * 120; // Max radius
              final offset = Offset(
                cos(angle) * distance,
                sin(angle) * distance,
              );

              return Transform.translate(
                offset: offset,
                child: Opacity(
                  opacity: _controller.value,
                  child: _AppCircle(app: app),
                ),
              );
            },
          );
        }),

        // Center button
        GestureDetector(
          onTap: _toggle,
          child: Container(
            width: 80,
            height: 80,
            decoration: BoxDecoration(
              shape: BoxShape.circle,
              gradient: LinearGradient(
                colors: [Color(0xFF667EEA), Color(0xFF764BA2)],
              ),
              boxShadow: [
                BoxShadow(
                  color: Color(0xFF667EEA).withOpacity(0.5),
                  blurRadius: 20,
                  spreadRadius: 2,
                ),
              ],
            ),
            child: Icon(
              _isExpanded ? Icons.close : Icons.apps,
              color: Colors.white,
              size: 32,
            ),
          ),
        ),
      ],
    );
  }
}
```

---

## OneMind OS Application

### Home Screen Concept

```
┌─────────────────────────────────────────────────────────┐
│  ● 9:41 AM                               🔋 85%         │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  "Good morning, Zeus."                                  │
│                                                         │
│  ┌─────────────────────────────────────────────────┐   │
│  │                                                  │   │
│  │            [Quick Actions Radial]                │   │
│  │                                                  │   │
│  │     🎤          👆           🏠                   │   │
│  │   Voice       Gesture     Smart                  │   │
│  │                                                  │   │
│  │          ┌─────────────┐                        │   │
│  │   🖨️     │     L       │     📊                  │   │
│  │  Print   │   (tap)     │   Status               │   │
│  │          └─────────────┘                        │   │
│  │                                                  │   │
│  │     🤖                          📅              │   │
│  │   Robots                     Calendar           │   │
│  │                                                  │   │
│  └─────────────────────────────────────────────────┘   │
│                                                         │
│  ┌───────────────────────────────────────────────────┐ │
│  │  Recent: "Printed phone stand - 2h ago"           │ │
│  │  Alert: "Server room temp elevated"              │ │
│  └───────────────────────────────────────────────────┘ │
│                                                         │
│               [  🎙️ Talk to Legacy  ]                  │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

### Key Apps for Radial Menu

1. **Voice** - Start voice conversation
2. **Gesture** - Enable gesture control mode
3. **Smart Home** - Home Assistant shortcuts
4. **3D Print** - Fabrication workflow
5. **Robots** - Drone/robot control
6. **Status** - System dashboard
7. **Calendar** - Today's schedule
8. **Custom** - User-defined action

---

## Audio Feedback

From HoloMat2's approach:

```
audio/
├── startup.wav      # App launch
├── home.wav         # Toggle radial menu
├── confirmation.wav # Action confirmed
└── reject.wav       # Action failed/cancelled
```

**For Flutter:**

```dart
import 'package:audioplayers/audioplayers.dart';

class SoundFeedback {
  static final AudioPlayer _player = AudioPlayer();

  static Future<void> play(FeedbackType type) async {
    final asset = switch (type) {
      FeedbackType.toggle => 'sounds/toggle.wav',
      FeedbackType.confirm => 'sounds/confirm.wav',
      FeedbackType.reject => 'sounds/reject.wav',
      FeedbackType.tap => 'sounds/tap.wav',
    };
    await _player.play(AssetSource(asset));
  }
}
```

---

## References

- Original: [HoloMat2/home_screen.py](https://github.com/Concept-Bytes/HoloMat2/blob/main/home_screen.py)
- Flutter animations: [flutter.dev/docs/development/ui/animations](https://flutter.dev/docs/development/ui/animations)
- Audio: [pub.dev/packages/audioplayers](https://pub.dev/packages/audioplayers)

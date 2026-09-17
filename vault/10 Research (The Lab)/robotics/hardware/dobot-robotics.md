# Dobot Robotics & OneMind Integration

**Type:** Hardware / Robotics  
**Category:** Physical Automation, Educational Robotics  
**Created:** 2026-01-20  
**Status:** Research & Integration Strategy  
**Website:** [dobot.cc](https://www.dobot.cc/)

---

## Overview

**Dobot** manufactures high-precision, desktop-sized robotic arms. Often used in education and light industry, they are the ideal hardware interface for bringing **OneMind** into the physical world.

Their flagship model, the **Dobot Magician**, is a 4-axis arm that can:

- 3D Print
- Laser Engrave
- Write/Draw
- Pick & Place (Suction or Gripper)

---

## Why Dobot for OneMind?

Most industrial robots (Kuka, Fanuc) are dangerous, expensive, and use proprietary languages. Dobot is:

1.  **Safe:** Desktop scale, low torque (won't break your hand).
2.  **Open:** Robust Python support (official + community drivers).
3.  **Versatile:** Interchangeable end-effectors allowing it to switch roles.

### The "Physical API" Concept

In the OneMind philosophy, we want to bridge the digital and physical.

- _Digital:_ "Activepieces detected a new order."
- _Physical:_ "Dobot picks up a pen and writes a Thank You note."

---

## Open Source Control (The Stack)

We will likely bypass the official proprietary studio software and use **Python** libraries to control the robot directly from our OneMind server (Ubuntu/Raspberry Pi).

### Top Libraries

1.  **`pydobot`** (Community Standard)
    - Simple, lightweight Python library.
    - Direct control over XYZ coordinates, suction cup, and gripper.
    - _Best for:_ Custom OneMind scripts.
    - [GitHub Link](https://github.com/luismesas/pydobot)

2.  **Official Python SDK**
    - Provided by Dobot, but often relies on DLLs (Windows-centric).
    - _Verdict:_ Use `pydobot` for Linux/headless integration.

---

## 🤖 Use Cases for OneMind

### 1. The "Physical Dashboard" (Ambient Notification)

Instead of a screen notification, the robot moves to indicate status.

- **Scenario:** Stock market crash or Server Down.
- **Action:** Robot picks up a "Red Flag" and places it on the center of the desk.
- **Connection:** `80 Treasury` -> Activepieces -> Python Script -> Dobot.

### 2. Automated Organization (The Sorter)

- **Scenario:** You dump SD cards, coins, or screws into a tray.
- **Hardware:** Dobot + Camera (OpenCV).
- **Action:** The OneMind Vision system identifies objects and sorts them into bins.
- **Connection:** `26 Mind` (Vision) -> Dobot.

### 3. "Handwritten" Notes (Legacy Evolution)

- **Scenario:** Sending personal letters to family/network.
- **Action:** You type a note in Obsidian. The robot writes it out with a fountain pen on high-quality paper.
- **Connection:** `50-74 LE` (Relationships) -> Dobot (Writing Kit).

---

## 🛠️ Integration Code Example

Here is how a Python script in the OneMind repo would control the robot:

```python
from serial.tools import list_ports
from pydobot import Dobot

# 1. Connect to Robot
port = list_ports.comports()[0].device
device = Dobot(port=port, verbose=True)

# 2. Define Positions
home = (200, 0, 0, 0)
bin_a = (200, 50, 0, 0)

# 3. Execute "Pick & Place"
device.move_to(*home)       # Go to neutral
device.suck(True)           # Turn on suction pump
device.move_to(*bin_a)      # Move to object
device.suck(False)          # Drop object

device.close()
```

---

## Future Vision: The "Laboratory"

As OneMind evolves, we can expand to a "Lab" setup (Project `84 Ventures`):

- **Dobot M1 Pro (SCARA):** Faster, stronger arm for actual light manufacturing.
- **Conveyor Belts:** Creating a mini assembly line for physical products.
- **AGVs (Carts):** Mobile bases moving items between stations.

**Conclusion:**
Dobot is the "Arm" of the Centaur. It allows the AI to touch the world.

---

**Tags:** `#robotics` `#python` `#automation` `#dobot` `#hardware`

**Related Docs:**

- `system-python-automation.md` (Coming soon)
- `reference-hybrid-civilization.md`


---

## Related

- [[00-24 SO (Sovereign Operations)/06 Research (R&D) (The Lab)/06 Research|06 Research]]

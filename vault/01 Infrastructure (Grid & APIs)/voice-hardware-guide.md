# Home Assistant Voice Hardware: ESP32-S3 Comparison

**Type:** Hardware Guide  
**Category:** One Mind OS / Interface  
**Created:** 2026-01-20  
**Status:** Hardware Selection  
**Pillar Connection:** `26 Mind` (Interface Layer)

---

## 🎯 The Goal

To establish a **Conversal Voice** interface for One Mind OS.
We need hardware that:

1.  Detects the wake word ("Hey Jarvis" / "Okay OneMind") reliably.
2.  Processes audio quickly (Low latency).
3.  Integrates seamlessly with our Home Assistant core.

---

## 🔝 The Hardware Contenders

We are comparing three devices based on the **ESP32-S3** chip.

### 1. Home Assistant Voice Preview Edition (The New Standard)

- **Best Feature:** **Dedicated XMOS Audio Processor.**
  - This is the "secret sauce." Unlike the generic Box-3, this device has a specialized audio chip for Echo Cancellation and Noise Reduction.
- **Controls:** Dial for volume, hardware mute switch.
- **Screen:** NO Screen (Just status LEDs).
- **OneMind Verdict:** **Recommended for Pure Voice.**
  - If you just want to talk to the room without looking at a device, this is the superior hardware.

### 2. ESP32-S3-BOX-3 (The All-Rounder)

- **Best Feature:** **Touchscreen Display.**
  - Allows you to _see_ the response or toggle buttons (e.g., turn off lights via screen if voice fails).
- **Expandability:** Has a dock for sensors (Temp/Humidity) and PCIe slot.
- **Drawback:** Lacks the dedicated XMOS chip, so wake word detection is _good_ but not _perfect_ in noisy rooms.
- **OneMind Verdict:** **Recommended for Desktop/Bedside.**
  - Use this where you might want visual feedback (e.g., "Show me the weather").

### 3. ESP32-S3-BOX-Lite (The Budget Option)

- **Feature:** Stripped down version.
- **Drawback:** Buttons are mechanical (not capacitive), fewer sensors, no dock.
- **Status:** Generally EOL (End of Life) and harder to find.
- **OneMind Verdict:** **Avoid.** The price difference doesn't justify the loss in features.

---

## 📊 Comparison Matrix

| Feature              | HA Voice Preview (VPE)    | ESP32-S3-BOX-3                | ESP32-S3-BOX-Lite      |
| :------------------- | :------------------------ | :---------------------------- | :--------------------- |
| **Wake Word Engine** | **Excellent** (XMOS Chip) | Good                          | Good                   |
| **Microphones**      | Dual Mic Array            | Dual Mic Array                | Dual Mic Array         |
| **Display**          | LED Ring                  | **2.4” Touchscreen**          | 2.4" LCD (Non-Touch\*) |
| **Visual Feedback**  | Minimal                   | Rich (Home Assistant Cards)   | Basic                  |
| **Audio Quality**    | Tuned for Speech          | General Purpose               | General Purpose        |
| **Expandability**    | Grove Port                | PCIe + Dock                   | Pmod Headers           |
| **Use Case**         | Ceiling/Shelf (Invisible) | Desk/Nightstand (Interactive) | Prototyping            |

_\*Note: Some Lite versions have touch, but builds vary._

---

## 🛠️ The Strategy for One Mind OS

We will use a **Hybrid Deployment**:

### Zone A: High Traffic (Living Room / Kitchen)

- **Device:** **HA Voice Preview Edition**
- **Why:** Noise cancellation is critical here (TV, Cooking, Talking). We need the XMOS chip to hear "Hey OneMind" over the background noise.

### Zone B: Interactive (Office / Bedroom)

- **Device:** **ESP32-S3-BOX-3**
- **Why:** You are sitting close to it. You might want to _tap_ a button to execute a "Focus Mode" routine or see the time/weather on the screen.

### Zone C: Physical Action (The Lab)

- **Integration:** The **Dobot** robot arm will be connected to the same Home Assistant instance.
- **Flow:**
  1.  You speak to the _Box-3_ on your desk: _"Initiate Clean Up."_
  2.  Home Assistant sends signal to _Dobot_.
  3.  Dobot executes the sorting script.

---

## 🔮 Configuration (ESPHome)

Both devices run on **ESPHome**, meaning we have full root control. No cloud dependency.

**Sample ESPHome YAML for Box-3:**

```yaml
substitutions:
  name: "onemind-office-voice"
  friendly_name: "OneMind Office"

esphome:
  name: ${name}
  name_add_mac_suffix: true
  project:
    name: esphome.voice-assistant
    version: "1.0"

esp32:
  board: esp32-s3-box-3
  framework:
    type: esp-idf

# Enable Microphone & Speaker
i2s_audio:
  i2s_lrclk_pin: 47
  i2s_bclk_pin: 17

microphone:
  - platform: i2s_audio
    i2s_din_pin: 16
    adc_type: external
    pdm: true

speaker:
  - platform: i2s_audio
    dac_type: external
    i2s_dout_pin: 46
    mode: stereo

# The Voice Assistant Pipeline
voice_assistant:
  microphone: my_microphone
  speaker: my_speaker
  on_listening:
    - light.turn_on: led_ring
  on_end:
    - light.turn_off: led_ring
```

---

**Tags:** `#hardware` `#voice-ai` `#home-assistant` `#esp32` `#interface`

**Related Docs:**

- `system-conversal-voice-architecture.md` (The software layer)
- `system-dobot-robotics.md` (The action layer)

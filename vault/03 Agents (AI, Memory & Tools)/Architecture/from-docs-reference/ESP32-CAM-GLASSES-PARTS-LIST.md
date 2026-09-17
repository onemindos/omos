# ESP32 DIY Camera Glasses - Parts List

> Shopping list for building Legacy-connected smart glasses

---

## Core Components

### Option A: XIAO ESP32S3 Sense (Recommended)

| Component | Qty | Price | Link |
|-----------|-----|-------|------|
| Seeed XIAO ESP32S3 Sense | 1 | $14 | [Seeed Studio](https://www.seeedstudio.com/XIAO-ESP32S3-Sense-p-5639.html) |

**Includes:** ESP32-S3 + OV2640 camera + PDM microphone

### Option B: Separate Components

| Component | Qty | Price | Link |
|-----------|-----|-------|------|
| ESP32-S3 DevKit | 1 | $10 | AliExpress/Amazon |
| OV2640 Camera Module | 1 | $5 | AliExpress/Amazon |
| INMP441 I2S Microphone | 1 | $5 | AliExpress/Amazon |

---

## Audio Output (Optional but Recommended)

| Component | Qty | Price | Notes |
|-----------|-----|-------|-------|
| MAX98357A I2S Amplifier | 1 | $5 | 3W mono amp |
| Small Speaker (8ohm, 1W) | 1 | $3 | 20-28mm diameter |

---

## Power

| Component | Qty | Price | Notes |
|-----------|-----|-------|-------|
| LiPo Battery 3.7V 500mAh | 1 | $8 | JST connector, ~30x20mm |
| TP4056 Charging Module | 1 | $2 | USB-C version preferred |
| Slide Switch | 1 | $1 | Power on/off |

---

## Frame / Enclosure

| Option | Price | Notes |
|--------|-------|-------|
| 3D Print (DIY) | ~$5 filament | Design needed |
| Safety Glasses (modify) | ~$10 | Drill/cut to mount |
| Existing Frames (modify) | ~$20 | Sunglasses/clear glasses |

---

## Wiring / Misc

| Component | Qty | Price |
|-----------|-----|-------|
| Jumper wires | 10+ | $3 |
| Heat shrink tubing | 1 set | $3 |
| Hot glue / epoxy | 1 | $5 |
| USB-C cable | 1 | $5 |

---

## Total Cost Estimate

| Build | Cost |
|-------|------|
| **Minimal (XIAO + battery)** | ~$25 |
| **Basic (+ audio out)** | ~$35 |
| **Full (+ nice frame)** | ~$50-60 |

---

## Recommended Vendors

| Vendor | Best For | Shipping |
|--------|----------|----------|
| Seeed Studio | XIAO boards | 1-2 weeks |
| AliExpress | Cheap components | 2-4 weeks |
| Amazon | Fast shipping | 1-2 days |
| Adafruit | Quality + tutorials | 3-5 days |
| SparkFun | Quality + tutorials | 3-5 days |
| DigiKey/Mouser | Pro components | 2-3 days |

---

## Pin Connections (XIAO ESP32S3 Sense)

```
XIAO ESP32S3 Sense Pinout for Glasses:

Camera: Built-in (OV2640)
Microphone: Built-in (PDM)

Audio Output (MAX98357A):
  - VIN  → 3.3V
  - GND  → GND
  - DIN  → GPIO 7 (I2S Data)
  - BCLK → GPIO 8 (I2S Clock)
  - LRC  → GPIO 9 (I2S Word Select)

Battery:
  - BAT+ → Battery positive
  - BAT- → GND

Power Switch:
  - Between battery and board
```

---

## Alternative: ESP32-CAM Pinout

```
ESP32-CAM Pinout for Glasses:

Camera: Built-in (OV2640)

External Microphone (INMP441):
  - VDD  → 3.3V
  - GND  → GND
  - SD   → GPIO 32 (I2S Data)
  - SCK  → GPIO 14 (I2S Clock)
  - WS   → GPIO 15 (I2S Word Select)
  - L/R  → GND (left channel)

Audio Output (MAX98357A):
  - VIN  → 3.3V
  - GND  → GND
  - DIN  → GPIO 25 (I2S Data Out)
  - BCLK → GPIO 26 (I2S Clock Out)
  - LRC  → GPIO 27 (I2S WS Out)

Note: ESP32-CAM has limited GPIO, may need creative pin sharing
```

---

## Next Steps

1. [ ] Order XIAO ESP32S3 Sense
2. [ ] Order audio components (MAX98357A + speaker)
3. [ ] Order LiPo battery + charger
4. [ ] Design/find 3D printable frame
5. [ ] Set up development environment (PlatformIO)
6. [ ] Flash test firmware
7. [ ] Test camera streaming
8. [ ] Test audio I/O
9. [ ] Assemble prototype
10. [ ] Integrate with Legacy
